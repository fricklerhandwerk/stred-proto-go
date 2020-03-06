package protobuf

import (
	"fmt"
	"regexp"
)

type label struct {
	value  string
	parent declarationContainer
}

func (l label) GetLabel() string {
	return l.value
}

func (l *label) SetLabel(label string) error {
	old := l.value
	l.value = label
	if err := l.parent.validateLabel(l); err != nil {
		l.value = old
		return err
	}
	return nil
}

func (l *label) hasLabel(other *label) bool {
	return l != other && l.value == other.value
}

// TODO: probably there is no need to have an extra type here, and validation
// can be done in a function
func (l label) validate() error {
	pattern := "[a-zA-Z]([0-9a-zA-Z_])*"
	regex := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
	if !regex.MatchString(l.value) {
		// TODO: there are at least two sources of errors which should be
		// differentiated by type: API caller and user. maybe API usage errors
		// should even result in a panic, since a nonsensical operation due to
		// broken implementation simply must not be allowed.
		return fmt.Errorf("Identifier must match %s", pattern)
	}
	return nil
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

func (k keyType) _fieldType() {}

type valueType string

const (
	Double valueType = "double"
	Float  valueType = "float"
	Bytes  valueType = "bytes"
)

func (v valueType) _fieldType() {}
