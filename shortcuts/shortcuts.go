package shortcuts

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

const (
	// null is the null ASCII character.
	null = "\x00"
	// soh is the 'Start of Header' ASCII character.
	soh = "\x01"
	// stx is the 'Start of Text' ASCII character.
	stx = "\x02"

	stringField = soh
	intField    = stx
	sliceField  = null

	appNameField            = "AppName"
	exePathField            = "Exe"
	startDirField           = "StartDir"
	iconPathField           = "icon"
	shortcutPathField       = "ShortcutPath"
	launchOptionsField      = "LaunchOptions"
	isHiddenField           = "IsHidden"
	allowDesktopConfigField = "AllowDesktopConfig"
	allowOverlayField       = "AllowOverlay"
	isOpenVrField           = "OpenVR"
	lastPlayTimeField       = "LastPlayTime"
	tagsField               = "tags"
)

const (
	boolValue valueType = iota
	idValue
	stringValue
	doubleQuoteString
	int32Value
	sliceValue
)

type valueType int

type Shortcut struct {
	Id                 int
	AppName            string
	ExePath            string
	StartDir           string
	IconPath           string
	ShortcutPath       string
	LaunchOptions      string
	IsHidden           bool
	AllowDesktopConfig bool
	AllowOverlay       bool
	IsOpenVr           bool
	LastPlayTimeEpoch  int32
	Tags               []string
}

func (o Shortcut) VdfV1String(footer []byte) string {
	sb := &strings.Builder{}

	for _, f := range o.fields() {
		f.appendFieldV1(sb)
	}

	if len(fileFooter) > 0 {
		sb.Write(footer)
	}

	return sb.String()
}

func (o Shortcut) fields() []vdfField {
	var fields []vdfField

	fields = append(fields, vdfField{
		valueType: idValue,
		idValue:   o.Id,
	})

	fields = append(fields, vdfField{
		name:        appNameField,
		valueType:   stringValue,
		stringValue: o.AppName,
	})

	fields = append(fields, vdfField{
		name:        exePathField,
		valueType:   doubleQuoteString,
		stringValue: o.ExePath,
	})

	fields = append(fields, vdfField{
		name:        startDirField,
		valueType:   doubleQuoteString,
		stringValue: o.StartDir,
	})

	fields = append(fields, vdfField{
		name:        iconPathField,
		valueType:   stringValue,
		stringValue: o.IconPath,
	})

	fields = append(fields, vdfField{
		name:        shortcutPathField,
		valueType:   stringValue,
		stringValue: o.ShortcutPath,
	})

	fields = append(fields, vdfField{
		name:        launchOptionsField,
		valueType:   stringValue,
		stringValue: o.LaunchOptions,
	})

	fields = append(fields, vdfField{
		name:      isHiddenField,
		valueType: boolValue,
		boolValue: o.IsHidden,
	})

	fields = append(fields, vdfField{
		name:      allowDesktopConfigField,
		valueType: boolValue,
		boolValue: o.AllowDesktopConfig,
	})

	fields = append(fields, vdfField{
		name:      allowOverlayField,
		valueType: boolValue,
		boolValue: o.AllowOverlay,
	})

	fields = append(fields, vdfField{
		name:      isOpenVrField,
		valueType: boolValue,
		boolValue: o.IsOpenVr,
	})

	fields = append(fields, vdfField{
		name:       lastPlayTimeField,
		valueType:  int32Value,
		int32Value: o.LastPlayTimeEpoch,
	})

	fields = append(fields, vdfField{
		name:       tagsField,
		valueType:  sliceValue,
		sliceValue: o.Tags,
	})

	return fields
}

type vdfField struct {
	name        string
	valueType   valueType
	stringValue string
	sliceValue  []string
	idValue     int
	boolValue   bool
	int32Value  int32
}

func (o vdfField) appendFieldV1(sb *strings.Builder) {
	switch o.valueType {
	case stringValue:
		o.appendStringV1(sb)
	case doubleQuoteString:
		o.appendDoubleQuoteStringV1(sb)
	case idValue:
		o.appendIdV1(sb)
	case boolValue:
		o.appendBoolV1(sb)
	case int32Value:
		o.appendInt32V1(sb)
	case sliceValue:
		o.appendSliceV1(sb)
	default:
		sb.WriteString(null)
	}
}

func (o vdfField) appendBoolV1(sb *strings.Builder) {
	sb.WriteString(intField)
	sb.WriteString(o.name)
	sb.WriteString(null)

	if o.boolValue {
		sb.WriteString(soh)
	} else {
		sb.WriteString(null)
	}

	sb.WriteString(strings.Repeat(null, 3))
}

func (o vdfField) appendStringV1(sb *strings.Builder){
	sb.WriteString(stringField)
	sb.WriteString(o.name)
	sb.WriteString(null)

	if len(o.stringValue) > 0 {
		sb.WriteString(o.stringValue)
	}

	sb.WriteString(null)
}

func (o vdfField) appendDoubleQuoteStringV1(sb *strings.Builder) {
	sb.WriteString(stringField)
	sb.WriteString(o.name)
	sb.WriteString(null)

	if len(o.stringValue) > 0 {
		sb.WriteString("\"")
		sb.WriteString(o.stringValue)
		sb.WriteString("\"")
	}

	sb.WriteString(null)
}

func (o vdfField) appendIdV1(sb *strings.Builder) {
	sb.WriteString(strconv.Itoa(o.idValue))
	sb.WriteString(null)
}

func (o vdfField) appendInt32V1(sb *strings.Builder) {
	sb.WriteString(intField)
	sb.WriteString(o.name)
	sb.WriteString(null)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(o.int32Value))
	sb.WriteString(string(b))
}

func (o vdfField) appendSliceV1(sb *strings.Builder) {
	sb.WriteString(null)
	sb.WriteString(o.name)

	for i, v := range o.sliceValue {
		sb.WriteString(null)
		sb.WriteString(stringField)
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(null)
		sb.WriteString(v)
	}

	sb.WriteString(null)
}

func WriteVdfV1(shortcuts []Shortcut, w io.Writer) error {
	_, err := w.Write(fileHeader)
	if err != nil {
		return err
	}

	for i, s := range shortcuts {
		var footer []byte

		if i < len(shortcuts) - 1 {
			footer = shortcutsDelim
		}

		_, err := w.Write([]byte(s.VdfV1String(footer)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write(fileFooter)
	if err != nil {
		return err
	}

	return nil
}
