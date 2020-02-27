package protobuf

import (
	"errors"
	"fmt"
)

type message struct {
	label
	fields      []messageField
	definitions []definition
	parent      definitionContainer
}

type messageField interface {
	validateAsMessageField() error
}

func (m message) GetFields() []messageField {
	return m.fields
}

// TODO: this is a bad interface, as it requires checking that the parent of
// the inserted field is really this message. instead wie should have
// `field.Insert(uint) error {}`, which may call its parent, which is the right
// thing by construction, to do the work.
// doing it that way has the added benifit that self-validation semantics are
// contained in the child type instead of calling the child from here, which
// calls the parent again.
func (m *message) InsertField(i uint, field messageField) error {
	if err := field.validateAsMessageField(); err != nil {
		return err
	}
	m.fields = append(m.fields, nil)
	copy(m.fields[i+1:], m.fields[i:])
	m.fields[i] = field

	return nil
}

func (m *message) newTypedField() *typedField {
	return &typedField{
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
		typedField: *m.newTypedField(),
	}
}

func (m *message) NewMap() *mapField {
	return &mapField{
		typedField: *m.newTypedField(),
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

func (m *message) NewMessage() *message {
	return &message{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

func (m *message) NewEnum() *enum {
	return &enum{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

func (m message) GetDefinitions() []definition {
	panic("not implemented")
}

func (m *message) InsertDefinition(i uint, d definition) error {
	panic("not implemented")
}

func (m message) validateAsDefinition() (err error) {
	err = m.parent.validateLabel(identifier(m.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return
	}
	return
}

func (m message) validateLabel(l identifier) error {
	// TODO: if the policy now develops such that everything is validated by its
	// parent, this should also be done by a function independent of the
	// identifier. this makes the whole extra type unnecessary.
	if err := l.validate(); err != nil {
		return err
	}
	for _, f := range m.fields {
		switch field := f.(type) {
		case *repeatableField:
			if field.GetLabel() == l.String() {
				return errors.New(fmt.Sprintf("label %q already declared", l.String()))
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
	case numberRange:
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
				return errors.New(fmt.Sprintf("field number %d already in use", uint(n)))
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

func (m message) validateNumberRange(n numberRange) error {
	panic("not implemented")
}

func (m *message) _fieldType() {}
