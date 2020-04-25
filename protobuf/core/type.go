package core

import "fmt"

// ValueType identifier. Valid values are all built-in types and `Message`s.
// Note that API consumers cannot create types which implement interface
// `Message`, and the package will only emit properly constructed `Message`s.
type ValueType interface {
	_isValueType()
}

// MapKeyType identifier.
// Valid values are all built-in types except:
// - `double`
// - `float`
// - `bytes`
type MapKeyType interface {
	_isKeyType()
}

type Type struct {
	value  ValueType
	parent Typed
}

type Typed interface {
	Type() *Type
}

func (t *Type) Get() ValueType {
	return t.value
}

func (t *Type) Set(value ValueType) error {
	old := t.value
	t.value = value
	if err := t.validate(); err != nil {
		t.value = old
		return err
	}
	switch old := t.value.(type) {
	case Definition:
		old.removeReference(t)
		switch v := value.(type) {
		case Definition:
			v.addReference(t)
		}
	}
	return nil
}

func (t *Type) Parent() Typed {
	return t.parent
}

func (t *Type) validate() error {
	if t.value == nil {
		return fmt.Errorf("type must not be nil")
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

func (k keyType) _isValueType()  {}
func (k keyType) _isMapKeyType() {}

type valueType string

const (
	Double valueType = "double"
	Float  valueType = "float"
	Bytes  valueType = "bytes"
)

func (v valueType) _isValueType() {}
