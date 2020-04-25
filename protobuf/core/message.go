package core

import (
	"errors"
	"fmt"
)

type NewMessage struct {
	label  Label
	parent DefinitionContainer
}

func (m *NewMessage) Label() *Label {
	if m.label.parent == nil {
		m.label.parent = m
	}
	return &m.label
}

func (m *NewMessage) InsertIntoParent() error {
	mm := &message{
		parent: m.parent,
		label: Label{
			value: m.label.Get(),
		},
	}
	mm.label.parent = mm
	return m.parent.insertMessage(mm)
}

func (m *NewMessage) Parent() DefinitionContainer {
	return m.parent
}

func (m *NewMessage) validateLabel(l *Label) error {
	return m.parent.validateLabel(l)
}

type message struct {
	label      Label
	fields     map[MessageField]struct{}
	messages   map[*message]struct{}
	enums      map[*enum]struct{}
	references map[*Type]struct{}
	parent     DefinitionContainer

	ValueType
}

func (m *message) Label() *Label {
	return &m.label
}

func (m message) Fields() (out []MessageField) {
	out = make([]MessageField, len(m.fields))
	i := 0
	for f := range m.fields {
		out[i] = f
	}
	return
}

func (m message) Messages() (out []Message) {
	out = make([]Message, len(m.messages))
	i := 0
	for d := range m.messages {
		out[i] = d
	}
	return
}

func (m message) Enums() (out []Enum) {
	out = make([]Enum, len(m.enums))
	i := 0
	for e := range m.enums {
		out[i] = e
	}
	return
}

func (m *message) insertField(f MessageField) error {
	if m.fields == nil {
		m.fields = make(map[MessageField]struct{})
	}
	if _, ok := m.fields[f]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := f.validateAsMessageField(); err != nil {
		return err
	}
	m.fields[f] = struct{}{}
	return nil
}

func (m *message) insertEnum(e *enum) error {
	if _, ok := m.enums[e]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := e.validate(); err != nil {
		return err
	}
	m.enums[e] = struct{}{}
	return nil
}

func (m *message) insertMessage(n *message) error {
	if _, ok := m.messages[n]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := n.validate(); err != nil {
		return err
	}
	m.messages[n] = struct{}{}
	return nil
}

func (m *message) addReference(t *Type) {
	if m.references == nil {
		m.references = make(map[*Type]struct{})
	}
	m.references[t] = struct{}{}
}

func (m *message) removeReference(t *Type) {
	delete(m.references, t)
}

func (m *message) NewField() Field {
	f := &repeatableField{parent: m}
	f.label.parent = f
	f.number.parent = f
	f.deprecated.parent = f
	return f
}

func (m *message) NewMap() Map {
	f := &mapField{parent: m}
	f.label.parent = f
	f.number.parent = f
	f.deprecated.parent = f
	f._type.parent = f
	return f
}

func (m *message) NewOneOf() OneOf {
	return &oneOf{parent: m}
}

func (m *message) NewReservedRange() ReservedRange {
	return &reservedRange{parent: m}
}

func (m *message) NewReservedNumber() ReservedNumber {
	n := &reservedNumber{parent: m}
	n.number.parent = n
	return n

}

func (m *message) NewReservedLabel() ReservedLabel {
	l := &reservedLabel{parent: m}
	l.label.parent = l
	return l
}

func (m *message) NewMessage() *NewMessage {
	return &NewMessage{parent: m}
}

func (m *message) NewEnum() *NewEnum {
	return &NewEnum{parent: m}
}

func (m *message) Parent() DefinitionContainer {
	return m.parent
}

func (m message) hasLabel(l *Label) bool {
	return m.label.hasLabel(l)
}

func (m *message) validateLabel(l *Label) error {
	switch l {
	case &m.label:
		return m.parent.validateLabel(l)
	default:
		for f := range m.fields {
			if f.hasLabel(l) {
				// TODO: return error type with reference to other declaration
				return fmt.Errorf("label %q already declared", l.value)
			}
		}
		for d := range m.messages {
			// TODO: return error type with reference to other declaration
			if d.label.hasLabel(l) {
				return fmt.Errorf("label %q already declared", l.value)
			}
		}
		for e := range m.enums {
			// TODO: return error type with reference to other declaration
			if e.label.hasLabel(l) {
				return fmt.Errorf("label %q already declared", l.value)
			}
		}
	}
	return nil
}

func (m message) validateNumber(n FieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch v := n.(type) {
	case *Number:
		if *v.value < 1 {
			return errors.New("message field number must be >= 1")
		}
	case *reservedRange:
		if *v.start.value < 1 {
			return errors.New("message field numbers must be >= 1")
		}
		if *v.end.value < 1 {
			return errors.New("message field numbers must be >= 1")
		}
	default:
		panic(fmt.Sprintf("unhandled field number type %T", v))
	}
	for f := range m.fields {
		if f.hasNumber(n) {
			return fmt.Errorf("field number %s already in use", n)
		}
	}
	return nil
}

func (m *message) validate() (err error) {
	return m.label.validate()
}
