package core

import (
	"fmt"
)

type oneOf struct {
	label  label
	fields map[*oneOfField]struct{}
	parent *message
}

func (o *oneOf) Label() Identifier {
	if o.label.parent == nil {
		o.label.parent = o
	}
	return &o.label
}

func (o *oneOf) NewField() OneOfField {
	v := &oneOfField{parent: o}
	v.typedField.label.parent = v
	v.typedField.number.parent = v
	v.typedField.deprecated.parent = v
	v.typedField._type.parent = v
	return v
}

func (o oneOf) Fields() (out []OneOfField) {
	out = make([]OneOfField, len(o.fields))
	i := 0
	for f := range o.fields {
		out[i] = f
	}
	return
}

func (o *oneOf) InsertIntoParent() error {
	return o.parent.insertField(o)
}

func (o *oneOf) Parent() Message {
	return o.parent
}

func (o *oneOf) insertField(f *oneOfField) error {
	if o.fields == nil {
		o.fields = make(map[*oneOfField]struct{})
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

func (o *oneOf) validateLabel(l *label) error {
	if o.hasLabel(l) {
		return fmt.Errorf("field label %s already in use", l)
	}
	return o.parent.validateLabel(l)
}

func (o oneOf) validateNumber(n FieldNumber) error {
	if o.hasNumber(n) {
		return fmt.Errorf("field number %s already in use", n)
	}
	return o.parent.validateNumber(n)
}

func (o *oneOf) validateAsMessageField() error {
	// TODO: let label self-validate
	if o.label.value == "" {
		return fmt.Errorf("oneof label not set")
	}
	if err := o.parent.validateLabel(&o.label); err != nil {
		return err
	}
	for f := range o.fields {
		if err := f.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (o oneOf) hasNumber(n FieldNumber) bool {
	for f := range o.fields {
		if f.hasNumber(n) {
			return true
		}
	}
	return false
}

func (o *oneOf) hasLabel(l *label) bool {
	for f := range o.fields {
		if f.hasLabel(l) {
			return true
		}
	}
	return o.label.hasLabel(l)
}

type oneOfField struct {
	typedField
	parent *oneOf
}

func (f *oneOfField) InsertIntoParent() error {
	return f.parent.insertField(f)
}

func (f *oneOfField) Parent() OneOf {
	return f.parent
}

func (f oneOfField) validateLabel(l *label) error {
	return f.parent.validateLabel(l)
}

func (f oneOfField) validateNumber(n FieldNumber) error {
	return f.parent.validateNumber(n)
}

func (f oneOfField) validateFlag(*flag) error {
	return nil
}
