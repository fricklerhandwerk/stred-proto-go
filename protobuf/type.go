package protobuf

import (
	"fmt"
	"regexp"
)

type declarationContainer interface {
	validateLabel(*label) error
}

type label struct {
	value  string
	parent declarationContainer
}

func (l label) Value() string {
	return l.value
}

func (l label) Label() string {
	return l.value
}

func (l *label) maybeLabel() *string {
	if l == nil {
		return nil
	}
	value := l.value
	return &value
}

func (l *label) SetValue(label string) error {
	if err := validateIdentifier(label); err != nil {
		return err
	}
	old := l.value
	l.value = label
	if err := l.parent.validateLabel(l); err != nil {
		l.value = old
		return err
	}
	return nil
}

func (l *label) SetLabel(label string) error {
	return l.SetValue(label)
}

func validateIdentifier(value string) (err error) {
	pattern := "[a-zA-Z]([0-9a-zA-Z_])*"
	regex := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
	if !regex.MatchString(value) {
		err = fmt.Errorf("Identifier must match %s", pattern)
	}
	return
}

func (l *label) hasLabel(other *label) bool {
	return l != other && l.value == other.value
}

type keyType string

const (
	Int32    keyType = "int32"
	Int64    keyType = "int64"
	Uint32   keyType = "uint32"
	Uint64   keyType = "uint64"
	Sint32   keyType = "sint32"
	Sint64   keyType = "sint64"
	Fixed32  keyType = "fixed32"
	Fixed64  keyType = "fixed64"
	Sfixed32 keyType = "sfixed32"
	Sfixed64 keyType = "sfixed64"
	Bool     keyType = "bool"
	String   keyType = "string"
)

func (k keyType) _isFieldType() {}
func (k keyType) _isKeyType()   {}

type valueType string

const (
	Double valueType = "double"
	Float  valueType = "float"
	Bytes  valueType = "bytes"
)

func (v valueType) _isFieldType() {}
