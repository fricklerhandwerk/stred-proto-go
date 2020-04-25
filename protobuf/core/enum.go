package core

import (
	"errors"
	"fmt"
	"strings"
)

type newEnum struct {
	label  label
	parent DefinitionContainer
}

func (e *newEnum) InsertIntoParent() error {
	ee := &enum{
		parent: e.parent,
		label: label{
			value: e.label.Get(),
		},
	}
	ee.label.parent = ee
	ee.allowAlias.parent = ee
	return e.parent.insertEnum(ee)
}

func (e *newEnum) Label() Identifier {
	if e.label.parent == nil {
		e.label.parent = e
	}
	return &e.label
}

func (e newEnum) Parent() DefinitionContainer {
	return e.parent
}

func (e newEnum) validateLabel(l *label) error {
	return e.parent.validateLabel(l)
}

type enum struct {
	label      label
	allowAlias flag
	fields     map[EnumField]struct{}
	references map[*_type]struct{}
	parent     DefinitionContainer

	ValueType
}

func (e *enum) Label() Identifier {
	return &e.label
}

func (e *enum) AllowAlias() Flag {
	return &e.allowAlias
}

func (e enum) validateFlag(*flag) error {
	// check if aliasing is in place
	numbers := make(map[uint]bool, len(e.fields))
	for field := range e.fields {
		switch f := field.(type) {
		case *variant:
			n := *f.number.value
			if numbers[n] {
				// TODO: return error type with references to aliased fields
				lines := []string{
					fmt.Sprintf("field number %d is used multiple times.", n),
					fmt.Sprintf("remove aliasing before setting %q.", "allow_alias = false"),
				}
				return errors.New(strings.Join(lines, " "))
			}
			numbers[n] = true
		case FieldNumber:
			continue
		default:
			panic(fmt.Sprintf("unhandled enum field type %T", f))
		}
	}
	return nil
}

func (e *enum) Fields() (out []EnumField) {
	out = make([]EnumField, len(e.fields))
	i := 0
	for f := range e.fields {
		out[i] = f
	}
	return
}

func (e *enum) insertField(f EnumField) error {
	if e.fields == nil {
		e.fields = make(map[EnumField]struct{})
	}
	if _, ok := e.fields[f]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := f.validateAsEnumField(); err != nil {
		return err
	}
	e.fields[f] = struct{}{}
	return nil
}

func (e *enum) addReference(t *_type) {
	if e.references == nil {
		e.references = make(map[*_type]struct{})
	}
	e.references[t] = struct{}{}
}

func (e *enum) removeReference(t *_type) {
	delete(e.references, t)
}

func (e *enum) NewVariant() Variant {
	v := &variant{parent: e}
	v.field.label.parent = v
	v.field.number.parent = v
	v.field.deprecated.parent = v
	return v
}

func (e *enum) NewReservedRange() ReservedRange {
	return &reservedRange{parent: e}
}

func (e *enum) NewReservedNumber() ReservedNumber {
	n := &reservedNumber{parent: e}
	n.number.parent = n
	return n
}

func (e *enum) NewReservedLabel() ReservedLabel {
	l := &reservedLabel{parent: e}
	l.label.parent = l
	return l
}

func (e enum) Parent() DefinitionContainer {
	return e.parent
}

func (e *enum) hasLabel(l *label) bool {
	return e.label.hasLabel(l)
}

func (e *enum) validateLabel(l *label) error {
	switch l {
	case &e.label:
		return e.parent.validateLabel(l)
	default:
		for f := range e.fields {
			if f.hasLabel(l) {
				return fmt.Errorf("label %s already declared", l.value)
			}
		}
	}
	return nil
}

func (e enum) validateNumber(n FieldNumber) error {
	// TODO: check that 0 is present and not reserved
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	for f := range e.fields {
		if f.hasNumber(n) {
			switch f.(type) {
			case *variant:
				switch n := n.(type) {
				case *Number:
					switch n.parent.(type) {
					case *variant:
						if e.allowAlias.value {
							return nil
						}
						lines := []string{
							fmt.Sprintf("field number %d already in use.", *n.value),
							fmt.Sprintf("set %q to allow multiple labels for one number.", "allow_alias = true"),
						}
						return errors.New(strings.Join(lines, " "))
					}
				}
			}
			// TODO: return error type with instance of duplication, no need to be so
			// verbose
			var source string
			switch v := n.(type) {
			case *Number:
				source = fmt.Sprintf("field number %d", *v.value)
			case *reservedRange:
				source = fmt.Sprintf("range %d to %d", *v.start.value, *v.end.value)
			default:
				panic(fmt.Sprintf("unhandled number type %T", n))
			}
			return fmt.Errorf("%s already in use", source)
		}
	}
	return nil
}

func (e *enum) validate() (err error) {
	if e.label.value == "" {
		return fmt.Errorf("label not set")
	}
	return e.parent.validateLabel(&e.label)
}

type variant struct {
	field
	parent *enum
}

func (v *variant) InsertIntoParent() error {
	vv := &variant{
		field:  v.field,
		parent: v.parent,
	}
	vv.label.parent = vv
	vv.number.parent = vv
	vv.deprecated.parent = vv
	return v.parent.insertField(vv)
}

func (v variant) Parent() Enum {
	return v.parent
}

func (v *variant) validateAsEnumField() (err error) {
	return v.field.validate()
}

func (v variant) validateLabel(l *label) error {
	return v.parent.validateLabel(l)
}

func (v variant) validateNumber(n FieldNumber) error {
	return v.parent.validateNumber(n)
}

func (v variant) validateFlag(f *flag) error {
	return v.parent.validateFlag(f)
}
