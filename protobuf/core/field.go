package core

import (
	"fmt"
)

type field struct {
	label      label
	number     Number
	deprecated flag
}

func (f *field) Label() Identifier {
	return &f.label
}

func (f *field) Number() *Number {
	return &f.number
}

func (f *field) Deprecated() Flag {
	return &f.deprecated
}

func (f *field) hasLabel(other *label) bool {
	return f.label.hasLabel(other)
}

func (f *field) hasNumber(n FieldNumber) bool {
	return f.number.hasNumber(n)
}

func (f *field) validate() (err error) {
	err = f.label.validate()
	if err != nil {
		return
	}
	// TODO: let number self-validate
	if f.number.value == nil {
		return fmt.Errorf("field number not set")
	}
	err = f.number.parent.validateNumber(&f.number)
	if err != nil {
		return
	}
	return
}

type typedField struct {
	field
	_type _type
}

func (f *typedField) Type() Type {
	return &f._type
}

func (f *typedField) validate() (err error) {
	if err = f.field.validate(); err != nil {
		return
	}
	if err = f._type.validate(); err != nil {
		return
	}
	return
}

type repeatableField struct {
	field
	repeatableType
	parent *message
}

func (r *repeatableField) InsertIntoParent() error {
	return r.parent.insertField(r)
}

func (r *repeatableField) Parent() Message {
	return r.parent
}

func (r *repeatableField) validateAsMessageField() (err error) {
	if err = r.field.validate(); err != nil {
		return
	}
	if err = r.repeatableType.validate(); err != nil {
		return
	}
	return
}

func (r *repeatableField) validateLabel(l *label) error {
	return r.parent.validateLabel(l)
}

func (r *repeatableField) validateNumber(n FieldNumber) error {
	return r.parent.validateNumber(n)
}
