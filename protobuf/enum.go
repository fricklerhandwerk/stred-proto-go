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

func (e *enum) insertField(i uint, field enumField) error {
	if err := field.validateAsEnumField(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	// <https://github.com/golang/go/wiki/SliceTricks#insert>
	// <https://stackoverflow.com/a/46130603/5147619>
	e.fields = append(e.fields, nil)
	copy(e.fields[i+1:], e.fields[i:])
	e.fields[i] = field
	return nil
}

func (e *enum) NewField() *enumeration {
	out := &enumeration{
		parent: e,
		field: field{
			parent: e,
			label: label{
				parent: e,
			},
			number: &number{
				parent: e,
			},
		},
	}
	out.identifier.parent = out
	out.number.integer.parent = out
	return out
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
		if f != l.parent && f.hasLabel(l.String()) {
			return fmt.Errorf("label %s already declared", l.String())
		}
	}
	return nil
}

func (e enum) validateNumber(n fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	for _, f := range e.fields {
		if f != n.getParent() && f.hasNumber(n) {
			switch f.(type) {
			case *enumeration:
				switch n.getParent().(type) {
				case *enumeration:
					if e.allowAlias {
						return nil
					}
					lines := []string{
						fmt.Sprintf("field number %d already in use.", n.(Number).GetValue()),
						fmt.Sprintf("set %q to allow multiple labels for one number.", "allow_alias = true"),
					}
					return errors.New(strings.Join(lines, " "))
				}
			}
			var source string
			switch v := n.(type) {
			case Number:
				source = fmt.Sprintf("field number %d", v.GetValue())
			case NumberRange:
				source = fmt.Sprintf("range %d to %d", v.GetStart(), v.GetEnd())
			default:
				panic(fmt.Sprintf("unhandled number type %T", n))
			}
			return fmt.Errorf("%s already in use", source)
		}
	}
	return nil
}

func (e *enum) InsertIntoParent(i uint) (err error) {
	return e.parent.insertDefinition(i, e)
}

func (e *enum) validateAsDefinition() (err error) {
	return e.parent.validateLabel(e.identifier)
}

func (e *enum) _fieldType() {}

type enumeration struct {
	field
	parent *enum
}

func (e *enumeration) InsertIntoParent(i uint) error {
	return e.parent.insertField(i, e)
}

func (e *enumeration) validateAsEnumField() (err error) {
	err = e.parent.validateLabel(e.identifier)
	if err != nil {
		return err
	}
	err = e.parent.validateNumber(e.number)
	if err != nil {
		return err
	}
	return nil
}
