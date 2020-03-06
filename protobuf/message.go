package protobuf

import (
	"errors"
	"fmt"
)

type message struct {
	*label
	fields      []messageField
	definitions []Definition
	parent      definitionContainer
}

func (m message) NumFields() uint {
	return uint(len(m.fields))
}

func (m message) Field(i uint) messageField {
	return m.fields[i]
}

func (m *message) insertField(i uint, field messageField) error {
	if err := field.validateAsMessageField(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	m.fields = append(m.fields, nil)
	copy(m.fields[i+1:], m.fields[i:])
	m.fields[i] = field
	return nil
}

func (m *message) newTypedField(parent interface{}) typedField {
	return typedField{
		field: field{
			parent: m,
			label: &label{
				parent: m,
			},
			number: &number{
				parent: m,
				integer: integer{
					parent: parent,
				},
			},
		},
	}
}

func (m *message) NewField() *repeatableField {
	out := &repeatableField{
		parent: m,
	}
	out.typedField = m.newTypedField(out)
	return out
}

func (m *message) NewMap() *mapField {
	out := &mapField{
		parent: m,
	}
	out.typedField = m.newTypedField(out)
	return out
}

func (m *message) NewOneOf() *oneOf {
	out := &oneOf{
		parent: m,
		label: label{
			parent: m,
		},
	}
	return out
}

func (m *message) NewReservedNumbers() *reservedNumbers {
	return &reservedNumbers{
		parent: m,
	}
}

func (m *message) NewReservedLabels() *reservedLabels {
	return &reservedLabels{
		parent: m,
	}
}

func (m *message) NewMessage() Message {
	out := &message{
		parent: m,
		label: &label{
			parent: m,
		},
	}
	return out
}

func (m *message) NewEnum() Enum {
	out := &enum{
		parent: m,
		label: &label{
			parent: m,
		},
	}
	return out
}

// TODO: use a common implementation for definition containers
func (m message) NumDefinitions() uint {
	return uint(len(m.definitions))
}

func (m message) Definition(i uint) Definition {
	return m.definitions[i]
}

func (m *message) insertDefinition(i uint, d Definition) error {
	if err := m.validateAsDefinition(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	m.definitions = append(m.definitions, nil)
	copy(m.definitions[i+1:], m.definitions[i:])
	m.definitions[i] = d
	return nil
}

func (m *message) InsertIntoParent(i uint) (err error) {
	return m.parent.insertDefinition(i, m)
}

func (m message) validateLabel(l *label) error {
	// TODO: if the policy now develops such that everything is validated by its
	// parent, this should also be done by a function independent of the
	// identifier. this makes the whole extra type unnecessary.
	if err := l.validate(); err != nil {
		return err
	}
	for _, f := range m.fields {
		if f.hasLabel(l) {
			return fmt.Errorf("label %q already declared", l.value)
		}
	}
	// TODO: definitions and fields share a namespace
	return nil
}

func (m message) validateNumber(n fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch v := n.(type) {
	case Number:
		if v.GetValue() < 1 {
			return errors.New("message field number must be >= 1")
		}
	case NumberRange:
		if v.GetStart() < 1 {
			return errors.New("message field numbers must be >= 1")
		}
	default:
		panic(fmt.Sprintf("unhandled field number type %T", v))
	}
	for _, f := range m.fields {
		if f.hasNumber(n) {
			return fmt.Errorf("field number %s already in use", n)
		}
	}
	return nil
}

func (m *message) validateAsDefinition() (err error) {
	return m.parent.validateLabel(m.label)
}

func (m *message) _fieldType() {}
