package shortcuts

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type RawParser interface {
	Parse() (Shortcut, error)
}

type defaultRawParser struct {
	gotId bool
	raw   string
	wip   Shortcut
}

func (o *defaultRawParser) Parse() (Shortcut, error) {
	if len(o.raw) == 0 {
		return o.wip, nil
	}

	if !o.gotId {
		err := o.parseId()
		if err != nil {
			return o.wip, err
		}
	}

	currentValueType, err := o.parseCurrentValueType()
	if err != nil {
		return o.wip, err
	}

	currentFieldName, isEnd, err := o.parseFieldName()
	if err != nil {
		return o.wip, err
	}

	if isEnd {
		return o.wip, nil
	}

	value, err := o.value(currentValueType)
	if err != nil {
		return o.wip, err
	}

	switch currentFieldName {
	case appNameField:
		o.wip.AppName = value
	case exePathField:
		o.wip.ExePath = trimDoubleQuote(value)
	case startDirField:
		o.wip.StartDir = trimDoubleQuote(value)
	case iconPathField:
		o.wip.IconPath = trimDoubleQuote(value)
	case shortcutPathField:
		o.wip.ShortcutPath = trimDoubleQuote(value)
	case launchOptionsField:
		o.wip.LaunchOptions = value
	case isHiddenField:
		o.wip.IsHidden = parseRawBoolValue(value)
	case allowDesktopConfigField:
		o.wip.AllowDesktopConfig = parseRawBoolValue(value)
	case allowOverlayField:
		o.wip.AllowOverlay = parseRawBoolValue(value)
	case isOpenVrField:
		o.wip.IsOpenVr = parseRawBoolValue(value)
	case lastPlayTimeField:
		o.wip.LastPlayTimeEpoch = parseRawInt32Value(value)
	case tagsField:
		// TODO: Parse tags field.
	}

	return o.Parse()
}

func (o *defaultRawParser) parseId() error {
	// Drop the ID + null.
	value, ok := o.get(2, false)
	if !ok {
		return errors.New("Failed to cut ID field - index out of range")
	}

	i, err := strconv.Atoi(string(value[0]))
	if err != nil {
		return errors.New("Failed to parse shortcut ID - " + err.Error())
	}

	o.wip.Id = i
	o.gotId = true

	return nil
}

func (o *defaultRawParser) parseCurrentValueType() (valueType, error) {
	value, ok := o.get(1, false)
	if !ok {
		return stringValue, errors.New("Failed to read type field - no bytes remaining")
	}

	var currentValueType valueType

	switch string(value[0]) {
	case sliceField:
		currentValueType = sliceValue
	case intField:
		currentValueType = intValue
	case stringField:
		currentValueType = stringValue
	default:
		return stringValue, fmt.Errorf("%s, %x", "Invalid field type", value)
	}

	return currentValueType, nil
}

func (o *defaultRawParser) parseFieldName() (name string, isEof bool, err error) {
	// Drop the field name and the null terminator.
	v, ok := o.get(strings.Index(o.raw, null) + 1, true)
	if !ok {
		return "", false, errors.New("Field name is missing null terminator")
	}

	if !unicode.IsLetter(rune(v[0])) {
		return "", false, errors.New("Field name does not start with a letter")
	}

	return v, false, nil
}

func (o *defaultRawParser) value(current valueType) (string, error) {
	var numToCopy int
	var excludeLastChar bool

	switch current {
	case stringValue:
		numToCopy = strings.Index(o.raw, null) + 1
		excludeLastChar = true
	case intValue:
		numToCopy = 4
	case sliceValue:
		// TODO: Jank.
		numToCopy = strings.LastIndex(o.raw, null) + 1
		excludeLastChar = true
	default:
		return "", errors.New("Unknown field type - " + strconv.Itoa(int(current)))
	}

	value, ok := o.get(numToCopy, excludeLastChar)
	if !ok {
		return "", errors.New("Failed to read value field")
	}

	return value, nil
}

func (o *defaultRawParser) get(numberOfBytes int, excludeLastChar bool) (string, bool) {
	if isIndexOutsideString(numberOfBytes - 1, o.raw) {
		return "", false
	}

	value := o.raw[0:numberOfBytes]

	o.raw = o.raw[numberOfBytes:]

	if excludeLastChar {
		valLen := len(value)
		if valLen > 0 {
			value = value[0:valLen-1]
		}
	}

	return value, true
}

func (o *defaultRawParser) deleteBeforeIndex(index int) bool {
	if isIndexOutsideString(index, o.raw) {
		return false
	}

	o.raw = o.raw[index:]

	return true
}

func isIndexOutsideString(index int, s string) bool {
	totalIndexes := len(s) - 1

	if totalIndexes - index < 0 {
		return true
	}

	return false
}

// TODO: Finish this POS.
func parseTags(raw string) ([]string, int) {
	var values []string
	expectInt := true
	for i, s := range strings.Split(raw, null + one) {
		if expectInt {
			if !unicode.IsDigit(rune(s[0])) {
				return values, i
			}
			expectInt = false
		} else {
			values = append(values, s)
		}
	}

	return values, 0
}

var (
	fileHeader     = []byte{0, 's', 'h', 'o', 'r', 't', 'c', 'u', 't', 's', 0, 0}
	shortcutsDelim = []byte{8, 8, 0}
	fileFooter     = []byte{8, 8, 8, 8}
)

func Shortcuts(r io.Reader) ([]Shortcut, error) {
	var shortcuts []Shortcut
	s := bufio.NewScanner(r)
	s.Split(splitConfigData)

	for s.Scan() {
		sc, err := NewShortcut(s.Text())
		if err != nil {
			return shortcuts, err
		}

		shortcuts = append(shortcuts, sc)
	}

	return shortcuts, nil
}

func splitConfigData(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, fileHeader); i >= 0 {
		return i + len(fileHeader), nil, nil
	}

	if i := bytes.Index(data, shortcutsDelim); i >= 0 {
		return i + len(shortcutsDelim), data[0:i], nil
	}

	if i := bytes.Index(data, fileFooter); i >= 0 {
		return i + len(fileFooter), data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func NewShortcut(rawData string) (Shortcut, error) {
	return NewRawParser(rawData).Parse()
}

func NewRawParser(rawData string) RawParser {
	return &defaultRawParser{
		raw: rawData,
	}
}

func parseRawInt32Value(raw string) int32 {
	var i int32

	err := binary.Read(strings.NewReader(raw), binary.LittleEndian, &i)
	if err != nil {
		return 0
	}

	return i
}

func parseRawBoolValue(raw string) bool {
	var b bool

	err := binary.Read(strings.NewReader(raw), binary.LittleEndian, &b)
	if err != nil {
		return false
	}

	return b
}

func trimDoubleQuote(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "\""), "\"")
}
