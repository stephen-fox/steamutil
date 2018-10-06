package vdf

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

const (
	boolValue FieldValueType = iota
	idValue
	stringValue
	doubleQuoteString
	int32Value
	sliceValue
)

type FieldValueType int

type Field interface {
	Name() string
	ValueType() FieldValueType
	StringValue() string
	SliceValue() []string
	IdValue() int
	BoolValue() bool
	Int32Value() int32
	Append(sb *strings.Builder, version FormatVersion) error
}

type defaultField struct {
	name        string
	valueType   FieldValueType
	stringValue string
	sliceValue  []string
	idValue     int
	boolValue   bool
	int32Value  int32
}

func (o *defaultField) Name() string {
	return o.name
}

func (o *defaultField) ValueType() FieldValueType {
	return o.valueType
}

func (o *defaultField) StringValue() string {
	return o.stringValue
}

func (o *defaultField) SliceValue() []string {
	return o.sliceValue
}

func (o *defaultField) IdValue() int {
	return o.idValue
}

func (o *defaultField) BoolValue() bool {
	return o.boolValue
}

func (o *defaultField) Int32Value() int32 {
	return o.int32Value
}

func (o *defaultField) Append(sb *strings.Builder, version FormatVersion) error {
	switch version {
	case V1:
		o.appendFieldV1(sb)
	default:
		return errors.New("Format " + string(version) + " is not supported")
	}

	return nil
}

func (o defaultField) appendFieldV1(sb *strings.Builder) {
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

func (o defaultField) appendBoolV1(sb *strings.Builder) {
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

func (o defaultField) appendStringV1(sb *strings.Builder){
	sb.WriteString(stringField)
	sb.WriteString(o.name)
	sb.WriteString(null)

	if len(o.stringValue) > 0 {
		sb.WriteString(o.stringValue)
	}

	sb.WriteString(null)
}

func (o defaultField) appendDoubleQuoteStringV1(sb *strings.Builder) {
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

func (o defaultField) appendIdV1(sb *strings.Builder) {
	sb.WriteString(strconv.Itoa(o.idValue))
	sb.WriteString(null)
}

func (o defaultField) appendInt32V1(sb *strings.Builder) {
	sb.WriteString(intField)
	sb.WriteString(o.name)
	sb.WriteString(null)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(o.int32Value))
	sb.WriteString(string(b))
}

func (o defaultField) appendSliceV1(sb *strings.Builder) {
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