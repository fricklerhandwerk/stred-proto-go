package protobuf

import (
	"errors"
	"fmt"
)

type definitionContainer interface {
	insertDefinition(index uint, definition Definition) error
	validateLabel(*label) error
}

type newMessage struct {
	message *message
}

func (m *newMessage) InsertIntoParent(i uint) error {
	if err := m.message.parent.insertDefinition(i, m.message); err != nil {
		return err
	}
	m.message = &message{
		parent: m.message.parent,
	}
	return nil
}

func (m *newMessage) MaybeLabel() *string {
	return m.message.maybeLabel()
}

func (m *newMessage) SetLabel(l string) (err error) {
	if m.message.label == nil {
		m.message.label = &label{
			parent: m.message.parent,
		}
		defer func() {
			if err != nil {
				m.message.label = nil
			}
		}()
	}
	return m.message.SetLabel(l)
}

func (m *newMessage) NumDefinitions() uint {
	return m.message.NumDefinitions()
}

func (m *newMessage) Definition(i uint) Definition {
	return m.message.Definition(i)
}

func (m *newMessage) NumFields() uint {
	return m.message.NumFields()
}

func (m *newMessage) Field(i uint) MessageField {
	return m.message.Field(i)
}

func (m *newMessage) NewEnum() NewEnum {
	return m.message.NewEnum()
}

func (m *newMessage) NewMessage() NewMessage {
	return m.message.NewMessage()
}

func (m *newMessage) NewField() NewField {
	return m.message.NewField()
}

func (m *newMessage) NewMap() NewMap {
	return m.message.NewMap()
}

func (m *newMessage) NewOneOf() NewOneOf {
	return m.message.NewOneOf()
}

func (m *newMessage) NewReservedLabels() NewReservedLabels {
	return m.message.NewReservedLabels()
}

func (m *newMessage) NewReservedNumbers() NewReservedNumbers {
	return m.message.NewReservedNumbers()
}

type message struct {
	*label
	fields      []MessageField
	definitions []Definition
	parent      definitionContainer
}

func (m message) NumFields() uint {
	return uint(len(m.fields))
}

func (m message) Field(i uint) MessageField {
	return m.fields[i]
}

func (m *message) insertField(i uint, field MessageField) error {
	if err := field.validateAsMessageField(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	m.fields = append(m.fields, nil)
	copy(m.fields[i+1:], m.fields[i:])
	m.fields[i] = field
	return nil
}

func (m *message) NewField() NewField {
	return &newRepeatableField{
		repeatableField: &repeatableField{parent: m},
	}
}

func (m *message) NewMap() NewMap {
	return &newMapField{
		mapField: &mapField{parent: m},
	}
}

func (m *message) NewOneOf() NewOneOf {
	return &newOneOf{
		oneOf: &oneOf{parent: m},
	}
}

func (m *message) NewReservedNumbers() NewReservedNumbers {
	return &newReservedNumbers{
		reservedNumbers: &reservedNumbers{parent: m},
	}
}

func (m *message) NewReservedLabels() NewReservedLabels {
	return &newReservedLabels{
		reservedLabels: &reservedLabels{parent: m},
	}
}

func (m *message) NewMessage() NewMessage {
	return &newMessage{
		message: &message{parent: m},
	}
}

func (m *message) NewEnum() NewEnum {
	return &newEnum{
		enum: &enum{parent: m},
	}
}

// TODO: use a common implementation for definition containers
func (m message) NumDefinitions() uint {
	return uint(len(m.definitions))
}

func (m message) Definition(i uint) Definition {
	return m.definitions[i]
}

func (m *message) insertDefinition(i uint, d Definition) error {
	if err := d.validateAsDefinition(); err != nil {
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
	if l == nil {
		return fmt.Errorf("label not set")
	}
	for _, f := range m.fields {
		if f.hasLabel(l) {
			return fmt.Errorf("label %q already declared", l.value)
		}
	}
	// TODO: definitions and fields share a namespace
	return nil
}

func (m message) validateNumber(n FieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch v := n.(type) {
	case *number:
		if v.Value() < 1 {
			return errors.New("message field number must be >= 1")
		}
	case *numberRange:
		if v.Start() < 1 {
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

func (m *message) _isFieldType() {}
