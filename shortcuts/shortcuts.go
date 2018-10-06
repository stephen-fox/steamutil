package shortcuts

import (
	"encoding/binary"
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

func (o Shortcut) VdfString() string {
	sb := &strings.Builder{}

	for _, f := range o.fields() {
		f.append(sb)
	}

	return sb.String()
}

func (o Shortcut) fields() []field {
	var fields []field

	fields = append(fields, field{
		valueType: idValue,
		idValue:   o.Id,
	})

	fields = append(fields, field{
		name:        appNameField,
		valueType:   stringValue,
		stringValue: o.AppName,
	})

	fields = append(fields, field{
		name:        exePathField,
		valueType:   doubleQuoteString,
		stringValue: o.ExePath,
	})

	fields = append(fields, field{
		name:        startDirField,
		valueType:   doubleQuoteString,
		stringValue: o.StartDir,
	})

	fields = append(fields, field{
		name:        iconPathField,
		valueType:   stringValue,
		stringValue: o.IconPath,
	})

	fields = append(fields, field{
		name:        shortcutPathField,
		valueType:   stringValue,
		stringValue: o.ShortcutPath,
	})

	fields = append(fields, field{
		name:        launchOptionsField,
		valueType:   stringValue,
		stringValue: o.LaunchOptions,
	})

	fields = append(fields, field{
		name:      isHiddenField,
		valueType: boolValue,
		boolValue: o.IsHidden,
	})

	fields = append(fields, field{
		name:      allowDesktopConfigField,
		valueType: boolValue,
		boolValue: o.AllowDesktopConfig,
	})

	fields = append(fields, field{
		name:      allowOverlayField,
		valueType: boolValue,
		boolValue: o.AllowOverlay,
	})

	fields = append(fields, field{
		name:      isOpenVrField,
		valueType: boolValue,
		boolValue: o.IsOpenVr,
	})

	fields = append(fields, field{
		name:       lastPlayTimeField,
		valueType:  int32Value,
		int32Value: o.LastPlayTimeEpoch,
	})

	fields = append(fields, field{
		name:       tagsField,
		valueType:  sliceValue,
		sliceValue: o.Tags,
	})

	return fields
}

type field struct {
	name        string
	valueType   valueType
	stringValue string
	sliceValue  []string
	idValue     int
	boolValue   bool
	int32Value  int32
}

func (o field) append(sb *strings.Builder) {
	switch o.valueType {
	case stringValue:
		o.appendString(sb)
	case doubleQuoteString:
		o.appendDoubleQuoteString(sb)
	case idValue:
		o.appendId(sb)
	case boolValue:
		o.appendBool(sb)
	case int32Value:
		o.appendInt32(sb)
	case sliceValue:
		o.appendSlice(sb)
	default:
		sb.WriteString(null)
	}
}

func (o field) appendBool(sb *strings.Builder) {
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

func (o field) appendString(sb *strings.Builder){
	sb.WriteString(stringField)
	sb.WriteString(o.name)
	sb.WriteString(null)

	if len(o.stringValue) > 0 {
		sb.WriteString(o.stringValue)
	}

	sb.WriteString(null)
}

func (o field) appendDoubleQuoteString(sb *strings.Builder) {
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

func (o field) appendId(sb *strings.Builder) {
	sb.WriteString(null)
	sb.WriteString(strconv.Itoa(o.idValue))
	sb.WriteString(null)
}

func (o field) appendInt32(sb *strings.Builder) {
	sb.WriteString(intField)
	sb.WriteString(o.name)
	sb.WriteString(null)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(o.int32Value))
	sb.WriteString(string(b))
}

func (o field) appendSlice(sb *strings.Builder) {
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
