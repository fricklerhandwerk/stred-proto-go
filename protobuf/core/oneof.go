package core

import (
	"fmt"
)

type OneOf struct {
	label  Label
	fields map[*OneOfField]struct{}
	parent *message
}

func (o *OneOf) Label() *Label {
	if o.label.parent == nil {
		o.label.parent = o
	}
	return &o.label
}

func (o *OneOf) NewField() *OneOfField {
	v := &OneOfField{parent: o}
	v.typedField.label.parent = v
	v.typedField.number.parent = v
	v.typedField.deprecated.parent = v
	v.typedField._type.parent = v
	return v
}

func (o OneOf) Fields() (out []*OneOfField) {
	out = make([]*OneOfField, len(o.fields))
	i := 0
	for f := range o.fields {
		out[i] = f
	}
	return
}

func (o *OneOf) InsertIntoParent() error {
	return o.parent.insertField(o)
}

func (o *OneOf) Parent() Message {
	return o.parent
}

func (o *OneOf) insertField(f *OneOfField) error {
	if o.fields == nil {
		o.fields = make(map[*OneOfField]struct{})
	}
	if _, ok := o.fields[f]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := f.validate(); err != nil {
		return err
	}
	o.fields[f] = struct{}{}
	return nil
}

func (o *OneOf) validateLabel(l *Label) error {
	if o.hasLabel(l) {
		return fmt.Errorf("field label %s already in use", l)
	}
	return o.parent.validateLabel(l)
}

func (o OneOf) validateNumber(n FieldNumber) error {
	if o.hasNumber(n) {
		return fmt.Errorf("field number %s already in use", n)
	}
	return o.parent.validateNumber(n)
}

func (o *OneOf) validateAsMessageField() error {
	if err := o.label.validate(); err != nil {
		return err
	}
	for f := range o.fields {
		if err := f.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (o OneOf) hasNumber(n FieldNumber) bool {
	for f := range o.fields {
		if f.hasNumber(n) {
			return true
		}
	}
	return false
}

func (o *OneOf) hasLabel(l *Label) bool {
	for f := range o.fields {
		if f.hasLabel(l) {
			return true
		}
	}
	return o.label.hasLabel(l)
}

type OneOfField struct {
	typedField
	parent *OneOf
}

func (f *OneOfField) InsertIntoParent() error {
	return f.parent.insertField(f)
}

func (f *OneOfField) Parent() *OneOf {
	return f.parent
}

func (f OneOfField) validateLabel(l *Label) error {
	return f.parent.validateLabel(l)
}

func (f OneOfField) validateNumber(n FieldNumber) error {
	return f.parent.validateNumber(n)
}

func (f OneOfField) validateFlag(*Flag) error {
	return nil
}
