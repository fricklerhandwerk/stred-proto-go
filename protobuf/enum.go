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

type enumField interface {
	InsertIntoParent(uint) error
	validateAsEnumField() error
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

func (e enum) GetFields() []enumField {
	return e.fields
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

func (e enum) validateLabel(l identifier) error {
	if err := l.validate(); err != nil {
		return err
	}
	for _, f := range e.fields {
		switch field := f.(type) {
		case *enumeration:
			if field.GetLabel() == l.String() {
				return errors.New(fmt.Sprintf("label %s already declared", l.String()))
			}
		case *reservedLabels:
			for _, r := range field.GetLabels() {
				if r == l.String() {
					return errors.New(fmt.Sprintf("label %s already declared", l.String()))
				}
			}
		case *reservedNumbers:
			continue
		default:
			panic(fmt.Sprintf("unhandled enum field type %T", f))
		}
	}
	return nil
}

func (e enum) validateNumber(f fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch n := f.(type) {
	case number:
		return e.validateNumberSingle(n)
	case numberRange:
		return e.validateNumberRange(n)
	default:
		panic(fmt.Sprintf("unhandled field number type %T", f))
	}
}

func (e enum) validateNumberSingle(n number) error {
	for _, f := range e.fields {
		switch field := f.(type) {
		case *enumeration:
			if !e.allowAlias && field.GetNumber() == uint(n) {
				lines := []string{
					fmt.Sprintf("field number %d already in use.", uint(n)),
					fmt.Sprintf("set %q to allow multiple labels for one number.", "allow_alias = true"),
				}
				return errors.New(strings.Join(lines, " "))
			}
		case *reservedNumbers:
			panic("not implemented")
		case *reservedLabels:
			continue
		default:
			panic(fmt.Sprintf("unhandled field number type %T", f))
		}
	}
	return nil
}

func (e enum) validateNumberRange(n numberRange) error {
	panic("not implemented")
}

func (e *enum) InsertIntoParent(i uint) (err error) {
	err = e.parent.validateLabel(identifier(e.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return
	}
	e.parent.insertDefinition(i, e)
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

func (e enumeration) validateAsEnumField() (err error) {
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
