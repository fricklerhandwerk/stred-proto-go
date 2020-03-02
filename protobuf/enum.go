package protobuf

import (
	"errors"
	"fmt"
	"strings"
)

type enum struct {
	label
	allowAlias bool
	fields     []enumField
	parent     definitionContainer
}

func (e *enum) SetAlias(b bool) error {
	if !b && e.allowAlias {
		// check if aliasing is in place
		numbers := make(map[uint]bool, len(e.fields))
		for _, field := range e.fields {
			switch f := field.(type) {
			case *enumeration:
				n := f.GetNumber()
				if numbers[n] {
					lines := []string{
						fmt.Sprintf("field number %d is used multiple times.", n),
						fmt.Sprintf("remove aliasing before setting %q.", "allow_alias = false"),
					}
					return errors.New(strings.Join(lines, " "))
				}
				numbers[n] = true
			case fieldNumber:
				continue
			default:
				panic(fmt.Sprintf("unhandled enum field type %T", f))
			}
		}
	}
	e.allowAlias = b
	return nil
}

func (e enum) GetAlias() bool {
	return e.allowAlias
}

func (e enum) NumFields() uint {
	return uint(len(e.fields))
}

func (e enum) Field(i uint) enumField {
	return e.fields[i]
}

func (e *enum) insertField(i uint, field enumField) {
	// <https://github.com/golang/go/wiki/SliceTricks#insert>
	// <https://stackoverflow.com/a/46130603/5147619>
	e.fields = append(e.fields, nil)
	copy(e.fields[i+1:], e.fields[i:])
	e.fields[i] = field
}

func (e *enum) NewField() *enumeration {
	return &enumeration{
		parent: e,
		field: field{
			parent: e,
			label: label{
				parent: e,
			},
		},
	}
}

func (e *enum) NewReservedNumbers() *reservedNumbers {
	return &reservedNumbers{
		parent: e,
	}
}

func (e *enum) NewReservedLabels() *reservedLabels {
	panic("not implemented")
}

func (e enum) validateLabel(l identifier) error {
	if err := l.validate(); err != nil {
		return err
	}
	for _, f := range e.fields {
		if f != nil && f.hasLabel(l.String()) {
			return fmt.Errorf("label %s already declared", l.String())
		}
	}
	return nil
}

func (e enum) validateNumber(n fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	for _, f := range e.fields {
		if f != nil && f.hasNumber(n) {
			switch f.(type) {
			case *enumeration:
				switch num := n.(type) {
				case Number:
					if e.allowAlias {
						// TODO: without additional information it is not clear which field
						// type this number belongs to. if it is a `reservedNumbers` field,
						// aliasing is still not allowed
						return nil
					}
					lines := []string{
						fmt.Sprintf("field number %d already in use.", num.GetValue()),
						fmt.Sprintf("set %q to allow multiple labels for one number.", "allow_alias = true"),
					}
					return errors.New(strings.Join(lines, " "))
				}
			}
			return fmt.Errorf("field number %s already in use", n)
		}
	}
	return nil
}

func (e *enum) InsertIntoParent(i uint) (err error) {
	err = e.validateAsDefinition()
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return
	}
	e.parent.insertDefinition(i, e)
	return
}

func (e *enum) validateAsDefinition() (err error) {
	if err = e.parent.validateLabel(e.label.label); err != nil {
		return
	}
	for i, f := range e.fields {
		e.fields[i] = nil
		defer func() { e.fields[i] = f }()
		if err = f.validateAsEnumField(); err != nil {
			return
		}
	}
	return
}
func (e *enum) _fieldType() {}

type enumeration struct {
	field
	parent *enum
}

func (e *enumeration) InsertIntoParent(i uint) error {
	if err := e.validateAsEnumField(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	e.parent.insertField(i, e)
	return nil
}

func (e *enumeration) validateAsEnumField() (err error) {
	err = e.parent.validateLabel(identifier(e.GetLabel()))
	if err != nil {
		return err
	}
	err = e.parent.validateNumber(number(e.GetNumber()))
	if err != nil {
		return err
	}
	return nil
}
