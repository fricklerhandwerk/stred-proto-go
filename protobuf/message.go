package protobuf

import (
	"errors"
	"fmt"
)

type message struct {
	label
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

func (m *message) insertField(i uint, field messageField) {
	m.fields = append(m.fields, nil)
	copy(m.fields[i+1:], m.fields[i:])
	m.fields[i] = field
}

func (m *message) newTypedField() typedField {
	return typedField{
		field: field{
			parent: m,
			label: label{
				parent: m,
			},
		},
	}
}

func (m *message) NewField() *repeatableField {
	return &repeatableField{
		parent:     m,
		typedField: m.newTypedField(),
	}
}

func (m *message) NewMap() *mapField {
	return &mapField{
		parent:     m,
		typedField: m.newTypedField(),
	}
}

func (m *message) NewOneOf() *oneOf {
	return &oneOf{
		parent: m,
		label: label{
			parent: m,
		},
	}
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
	return &message{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

func (m *message) NewEnum() Enum {
	return &enum{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

// TODO: use a common implementation for definition containers
func (m message) NumDefinitions() uint {
	return uint(len(m.definitions))
}

func (m message) Definition(i uint) Definition {
	return m.definitions[i]
}

func (m *message) insertDefinition(i uint, d Definition) {
	m.definitions = append(m.definitions, nil)
	copy(m.definitions[i+1:], m.definitions[i:])
	m.definitions[i] = d
}

func (m *message) InsertIntoParent(i uint) (err error) {
	err = m.parent.validateLabel(identifier(m.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return
	}
	m.parent.insertDefinition(i, m)
	return
}

func (m message) validateLabel(l identifier) error {
	// TODO: if the policy now develops such that everything is validated by its
	// parent, this should also be done by a function independent of the
	// identifier. this makes the whole extra type unnecessary.
	if err := l.validate(); err != nil {
		return err
	}
	// TODO: definitions and fields share a namespace
	// TODO: the assumption is that the declaration with this identifier was not
	// inserted yet. but we must also be able to perform a recursive validation
	// before inserting the message into a definition container...
	for _, f := range m.fields {
		switch field := f.(type) {
		case *repeatableField:
			if field.GetLabel() == l.String() {
				return fmt.Errorf("label %q already declared", l.String())
			}
		case *mapField:
			panic("not implemented")
		case *oneOf:
			panic("not implemented")
		case *reservedLabels:
			panic("not implemented")
		case *reservedNumbers:
			continue
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (m message) validateNumber(f fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch n := f.(type) {
	case number:
		return m.validateNumberSingle(n)
	case *numberRange:
		return m.validateNumberRange(n)
	default:
		panic(fmt.Sprintf("unhandled field number type %T", f))
	}
}

func (m message) validateNumberSingle(n number) error {
	if n < 1 {
		return errors.New("message field number must be >= 1")
	}
	for _, f := range m.fields {
		switch field := f.(type) {
		case *repeatableField:
			if field.GetNumber() == uint(n) {
				return fmt.Errorf("field number %d already in use", uint(n))
			}
		case *mapField:
			panic("not implemented")
		case *oneOf:
			panic("not implemented")
		case *reservedNumbers:
			panic("not implemented")
		case *reservedLabels:
			continue
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (m message) validateNumberRange(n *numberRange) error {
	if n.GetStart() < 1 {
		return errors.New("message field numbers must be >= 1")
	}
	for _, f := range m.fields {
		switch field := f.(type) {
		case *repeatableField:
			if n.intersects(number(field.GetNumber())) {
				return fmt.Errorf("field number %d already in use", field.GetNumber())
			}
		case *mapField:
			panic("not implemented")
		case *oneOf:
			panic("not implemented")
		case *reservedNumbers:
			panic("not implemented")
		case *reservedLabels:
			continue
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (m message) validateAsDefinition() (err error) {
	if err = m.parent.validateLabel(m.label.label); err != nil {
		return
	}
	for _, f := range m.fields {
		if err = f.validateAsMessageField(); err != nil {
			return
		}
	}
	for _, d := range m.definitions {
		if err = d.validateAsDefinition(); err != nil {
			return
		}
	}
	return
}

func (m *message) _fieldType() {}
