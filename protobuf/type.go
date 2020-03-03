package protobuf

import (
	"fmt"
	"regexp"
)

// TODO: probably there is no need to have an extra type here, and validation
// can be done in a function
type identifier struct {
	value  string
	parent interface{}
}

func (i identifier) validate() error {
	pattern := "[a-zA-Z]([0-9a-zA-Z_])*"
	regex := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
	if !regex.MatchString(i.String()) {
		// TODO: there are at least two sources of errors which should be
		// differentiated by type: API caller and user. maybe API usage errors
		// should even result in a panic, since a nonsensical operation due to
		// broken implementation simply must not be allowed.
		return fmt.Errorf("Identifier must match %s", pattern)
	}
	return nil
}

func (i *identifier) String() string {
	if i == nil {
		return ""
	}
	return i.value
}

type label struct {
	identifier identifier
	parent     declarationContainer
}

func (d label) GetLabel() string {
	return d.identifier.String()
}

func (d *label) SetLabel(label string) error {
	old := d.identifier.value
	d.identifier.value = label
	if err := d.parent.validateLabel(d.identifier); err != nil {
		d.identifier.value = old
		return err
	}
	return nil
}

func (d *label) hasLabel(l string) bool {
	return l == d.GetLabel()
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
