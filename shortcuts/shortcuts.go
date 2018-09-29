package shortcuts

import (
	"strconv"
	"strings"
)

const (
	shortcutsDelim = "\x08\x08"

	null = "\x00"
	one  = "\x01"
	two  = "\x02"

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
	boolValue  valueType = iota
	intValue
	stringValue
	doubleQuoteString
	dateValue
	sliceValue
)

type valueType int

type Shortcut struct {
	Id                 int
	AppName            string `shortcut:"AppName"`
	ExePath            string `shortcut:"Exe"`
	StartDir           string `shortcut:"StartDir"`
	IconPath           string `shortcut:"icon"`
	ShortcutPath       string `shortcut:"ShortcutPath"`
	LaunchOptions      string `shortcut:"LaunchOptions"`
	IsHidden           bool `shortcut:"IsHidden"`
	AllowDesktopConfig bool `shortcut:"AllowDesktopConfig"`
	AllowOverlay       bool `shortcut:"AllowOverlay"`
	IsOpenVr           bool `shortcut:"OpenVR"`
	LastPlayTime       string `shortcut:"LastPlayTime"`
	Tags               []string `shortcut:"tags"`
}

func (o Shortcut) VdfString() string {
	var fieldStrings []string

	for _, f := range o.fields() {
		fieldStrings = append(fieldStrings, f.string())
	}

	return strings.Join(fieldStrings, "")
}

func (o Shortcut) fields() ([]field) {
	var fields []field

	fields = append(fields, field{
		valueType: intValue,
		intValue:  o.Id,
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
		name:        lastPlayTimeField,
		valueType:   dateValue,
		stringValue: o.LastPlayTime,
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
	intValue    int
	boolValue   bool
}

func (o field) string() string {
	sb := &strings.Builder{}

	sb.WriteString(o.name)
	sb.WriteString(null)

	switch o.valueType {
	case stringValue:
		o.appendString(sb)
	case doubleQuoteString:
		o.appendDoubleQuoteString(sb)
	case intValue:
		o.appendInt(sb)
	case boolValue:
		o.appendBool(sb)
	case dateValue:
		o.appendDate(sb)
	case sliceValue:
		o.appendSlice(sb)
	default:
		sb.WriteString(null)
	}

	return sb.String()
}


func (o field) appendBool(sb *strings.Builder) *strings.Builder {
	if o.boolValue {
		sb.WriteString(one)
	} else {
		sb.WriteString(null)
	}

	sb.WriteString(strings.Repeat(null, 3))
	sb.WriteString(two)

	return sb
}

func (o field) appendString(sb *strings.Builder) *strings.Builder {
	if len(o.stringValue) > 0 {
		sb.WriteString(o.stringValue)
	}

	sb.WriteString(null)
	sb.WriteString(one)

	return sb
}

func (o field) appendDoubleQuoteString(sb *strings.Builder) *strings.Builder {
	if len(o.stringValue) > 0 {
		sb.WriteString("\"")
		sb.WriteString(o.stringValue)
		sb.WriteString("\"")
	}

	sb.WriteString(null)
	sb.WriteString(one)

	return sb
}

func (o field) appendInt(sb *strings.Builder) *strings.Builder {
	sb.WriteString(strconv.Itoa(o.intValue))
	sb.WriteString(null)
	sb.WriteString(one)

	return sb
}

// TODO: Actually do something.
func (o field) appendDate(sb *strings.Builder) *strings.Builder {
	sb.WriteString(strings.Repeat(null, 6))

	return sb
}

func (o field) appendSlice(sb *strings.Builder) *strings.Builder {
	for i, v := range o.sliceValue {
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(null)
		sb.WriteString(one)
		sb.WriteString(v)

		if i < len(o.sliceValue) - 1 {
			sb.WriteString(null)
		}
	}

	return sb
}
