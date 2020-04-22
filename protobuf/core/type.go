package core

import "fmt"

type _type struct {
	value  ValueType
	parent Typed
}

func (t *_type) Get() ValueType {
	return t.value

}

func (t *_type) Set(value ValueType) error {
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

func (t *_type) Parent() Typed {
	return t.parent
}

func (t *_type) validate() error {
	if t.value == nil {
		return fmt.Errorf("type must not be nil")
	}
	return nil
}

type repeatableType struct {
	_type
	repeated flag
	parent   *repeatableField
}

func (t *repeatableType) Type() Type {
	if t._type.parent == nil {
		t._type.parent = t
	}
	return &t._type
}

func (t *repeatableType) Repeated() Flag {
	if t.repeated.parent == nil {
		t.repeated.parent = t
	}
	return &t.repeated
}

func (t *repeatableType) validateFlag(f *flag) error {
	switch f {
	case &t.repeated:
		// TODO: if we ever have "safe mode" to prevent backwards-incompatible
		// changes, that is where errors whould happen
	}
	// deprecated
	return nil
}

func (t *repeatableType) Parent() Field {
	return t.parent
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
