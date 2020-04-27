package core

import (
	"errors"
	"fmt"
	"strings"
)

type Enum interface {
	Label() *Label
	AllowAlias() *Flag

	Fields() []EnumField

	NewVariant() *Variant
	NewReservedNumber() *ReservedNumber
	NewReservedRange() *ReservedRange
	NewReservedLabel() *ReservedLabel

	Parent() DefinitionContainer

	validate() error
	hasLabel(*Label) bool
	validateLabel(*Label) error
	validateNumber(FieldNumber) error
	insertField(EnumField) error

	addReference(*Type)
	removeReference(*Type)

	ValueType
}

type enum struct {
	label      Label
	allowAlias Flag
	fields     map[EnumField]struct{}
	references map[*Type]struct{}
	parent     DefinitionContainer

	ValueType
}

type EnumField interface {
	validateAsEnumField() error
	hasLabel(*Label) bool
	hasNumber(FieldNumber) bool
}

func (e *enum) Label() *Label {
	return &e.label
}

func (e *enum) AllowAlias() *Flag {
	return &e.allowAlias
}

func (e enum) validateFlag(*Flag) error {
	// check if aliasing is in place
	numbers := make(map[uint]bool, len(e.fields))
	for field := range e.fields {
		switch f := field.(type) {
		case *Variant:
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
		i++
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

func (e *enum) addReference(t *Type) {
	if e.references == nil {
		e.references = make(map[*Type]struct{})
	}
	e.references[t] = struct{}{}
}

func (e *enum) removeReference(t *Type) {
	delete(e.references, t)
}

func (e *enum) NewVariant() *Variant {
	v := &Variant{parent: e}
	v.field.label.parent = v
	v.field.number.parent = v
	v.field.deprecated.parent = v
	return v
}

func (e *enum) NewReservedRange() *ReservedRange {
	return &ReservedRange{parent: e}
}

func (e *enum) NewReservedNumber() *ReservedNumber {
	n := &ReservedNumber{parent: e}
	n.number.parent = n
	return n
}

func (e *enum) NewReservedLabel() *ReservedLabel {
	l := &ReservedLabel{parent: e}
	l.label.parent = l
	return l
}

func (e enum) Parent() DefinitionContainer {
	return e.parent
}

func (e *enum) hasLabel(l *Label) bool {
	return e.label.hasLabel(l)
}

func (e *enum) validateLabel(l *Label) error {
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
			case *Variant:
				switch n := n.(type) {
				case *Number:
					switch n.parent.(type) {
					case *Variant:
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
			case *ReservedRange:
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

type NewEnum struct {
	label  Label
	parent DefinitionContainer
}

func (e *NewEnum) InsertIntoParent() error {
	ee := &enum{
		parent: e.parent,
		label: Label{
			value: e.label.Get(),
		},
	}
	ee.label.parent = ee
	ee.allowAlias.parent = ee
	return e.parent.insertEnum(ee)
}

func (e *NewEnum) Label() *Label {
	if e.label.parent == nil {
		e.label.parent = e
	}
	return &e.label
}

func (e NewEnum) Parent() DefinitionContainer {
	return e.parent
}

func (e NewEnum) validateLabel(l *Label) error {
	return e.parent.validateLabel(l)
}

type Variant struct {
	field
	parent *enum
}

func (v *Variant) InsertIntoParent() error {
	vv := &Variant{
		field:  v.field,
		parent: v.parent,
	}
	vv.label.parent = vv
	vv.number.parent = vv
	vv.deprecated.parent = vv
	return v.parent.insertField(vv)
}

func (v Variant) Parent() Enum {
	return v.parent
}

func (v *Variant) validateAsEnumField() (err error) {
	return v.field.validate()
}

func (v Variant) validateLabel(l *Label) error {
	return v.parent.validateLabel(l)
}

func (v Variant) validateNumber(n FieldNumber) error {
	return v.parent.validateNumber(n)
}

func (v Variant) validateFlag(f *Flag) error {
	return v.parent.validateFlag(f)
}
