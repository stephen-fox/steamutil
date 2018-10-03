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

	var valueFieldEndIndex int
	var nextStartIndex int

	switch currentValueType {
	case stringValue:
		valueFieldEndIndex = strings.Index(o.raw, null)
		nextStartIndex = valueFieldEndIndex + 1
	case intValue:
		valueFieldEndIndex = 4
		nextStartIndex = valueFieldEndIndex
	case sliceValue:
		// TODO: Jank.
		valueFieldEndIndex = strings.LastIndex(o.raw, null)
		nextStartIndex = valueFieldEndIndex + 1
	default:
		return o.wip, errors.New("Unknown field type - " + strconv.Itoa(int(currentValueType)))
	}

	if isIndexOutsideString(valueFieldEndIndex, o.raw) {
		return o.wip, errors.New("Value field is missing terminator")
	}

	value := o.raw[0:valueFieldEndIndex]

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

	if !o.deleteBeforeIndex(nextStartIndex) {
		// EOF.
		return o.wip, nil
	}

	return o.Parse()
}

func (o *defaultRawParser) parseId() error {
	i, err := strconv.Atoi(string(o.raw[0]))
	if err != nil {
		return errors.New("Failed to parse shortcut ID - " + err.Error())
	}

	// Drop the ID + null.
	if !o.deleteBeforeIndex(2) {
		return errors.New("Failed to cut ID field - index out of range")
	}

	o.wip.Id = i
	o.gotId = true

	return nil
}

func (o *defaultRawParser) parseCurrentValueType() (valueType, error) {
	var currentValueType valueType

	t := string(o.raw[0])
	switch t {
	case sliceField:
		currentValueType = sliceValue
	case intField:
		currentValueType = intValue
	case stringField:
		currentValueType = stringValue
	default:
		return stringValue, fmt.Errorf("%s, %x", "Invalid field type", t)
	}

	// Drop the type field.
	if !o.deleteBeforeIndex(len(t)) {
		return stringValue, errors.New("Failed to cut type field - index out of range")
	}

	return currentValueType, nil
}

func (o *defaultRawParser) parseFieldName() (name string, isEof bool, err error) {
	if !unicode.IsLetter(rune(o.raw[0])) {
		return "", false, errors.New("Field name does not start with a letter")
	}

	fieldNameEndIndex := strings.Index(o.raw, null)
	if fieldNameEndIndex < 0 {
		return "", false, errors.New("Field name is missing null terminator")
	}

	currentFieldName := string(o.raw[0:fieldNameEndIndex])
	// Drop the field name and the null terminator.
	if !o.deleteBeforeIndex(fieldNameEndIndex + 1) {
		// EOF.
		return currentFieldName, true, nil
	}

	return currentFieldName, false, nil
}

func (o *defaultRawParser) deleteBeforeIndex(startingIndex int) bool {
	if isIndexOutsideString(startingIndex, o.raw) {
		return false
	}

	o.raw = o.raw[startingIndex:]

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
