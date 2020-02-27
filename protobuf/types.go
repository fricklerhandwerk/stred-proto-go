package protobuf

import (
	"errors"
	"fmt"
	"regexp"
)

type declaration interface {
	GetLabel() string
	SetLabel(string) error
}

type declarationContainer interface {
	validateLabel(identifier) error
}

type definitionContainer interface {
	declarationContainer

	GetDefinitions() []definition
	insertDefinition(index uint, def definition)
}

type definition interface {
	declaration
	declarationContainer

	validateNumber(fieldNumber) error
	InsertIntoParent(uint) error
}

// TODO: probably there is no need to have an extra type here, and validation
// can be done in a function
type identifier string

func (i identifier) validate() error {
	pattern := "[a-zA-Z]([0-9a-zA-Z_])*"
	regex := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
	if !regex.MatchString(i.String()) {
		// TODO: there are at least two sources of errors which should be
		// differentiated by type: API caller and user. maybe API usage errors
		// should even result in a panic, since a nonsensical operation due to
		// broken implementation simply must not be allowed.
		return errors.New(fmt.Sprintf("Identifier must match %s", pattern))
	}
	return nil
}

func (i *identifier) String() string {
	if i == nil {
		return ""
	}
	return string(*i)
}

type label struct {
	label  identifier
	parent declarationContainer
}

func (d label) GetLabel() string {
	return d.label.String()
}

func (d *label) SetLabel(label string) error {
	ident := identifier(label)
	if err := d.parent.validateLabel(ident); err != nil {
		return err
	}
	d.label = ident
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
