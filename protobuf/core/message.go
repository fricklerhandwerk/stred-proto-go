package core

import (
	"errors"
	"fmt"
)

type Message interface {
	Label() *Label
	Fields() []MessageField

	NewField() *Field
	NewMap() *Map
	NewOneOf() *OneOf
	NewReservedNumber() *ReservedNumber
	NewReservedRange() *ReservedRange
	NewReservedLabel() *ReservedLabel

	Messages() []Message
	NewMessage() *NewMessage

	Enums() []Enum
	NewEnum() *NewEnum

	Parent() DefinitionContainer
	Document() *Document
	String() string

	validate() error
	hasLabel(*Label) bool
	validateLabel(*Label) error
	validateNumber(FieldNumber) error
	insertField(MessageField) error

	addReference(MessageReference)
	removeReference(MessageReference)

	ValueType
}

type message struct {
	label      Label
	fields     map[MessageField]struct{}
	messages   map[*message]struct{}
	enums      map[*enum]struct{}
	references map[MessageReference]struct{}
	parent     DefinitionContainer

	ValueType
}

type MessageField interface {
	validateAsMessageField() error
	hasLabel(*Label) bool
	hasNumber(FieldNumber) bool
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

func (m *message) addReference(t MessageReference) {
	// TODO: check if reference is already in document, otherwise we would have
	// to do it in bulk when it gets relevant. it is actually even more
	// complicated: once something containing a reference is inserted, that
	// reference must be added, too. maybe there is no point in holding these
	// references after all, as managing them is too much of a hassle, and we
	// should collect them from the document on demand.
	if m.references == nil {
		m.references = make(map[MessageReference]struct{})
	}
	m.references[t] = struct{}{}
}

func (m *message) removeReference(t MessageReference) {
	delete(m.references, t)
}

type MessageReference interface {
	_isReference()
}

func (m *message) NewField() *Field {
	f := &Field{parent: m}
	f.label.parent = f
	f.number.parent = f
	f.deprecated.parent = f
	f._type.parent = f
	f.repeated.parent = f
	return f
}

func (m *message) NewMap() *Map {
	f := &Map{parent: m}
	f.label.parent = f
	f.number.parent = f
	f.deprecated.parent = f
	f._type.parent = f
	f.keyType.parent = f
	return f
}

func (m *message) NewOneOf() *OneOf {
	o := &OneOf{parent: m}
	o.label.parent = o
	return o
}

func (m *message) NewReservedRange() *ReservedRange {
	r := &ReservedRange{parent: m}
	r.start.parent = r
	r.end.parent = r
	return r
}

func (m *message) NewReservedNumber() *ReservedNumber {
	n := &ReservedNumber{parent: m}
	n.number.parent = n
	return n
}

func (m *message) NewReservedLabel() *ReservedLabel {
	l := &ReservedLabel{parent: m}
	l.label.parent = l
	return l
}

func (m *message) NewMessage() *NewMessage {
	n := &NewMessage{parent: m}
	n.label.parent = n
	return n
}

func (m *message) NewEnum() *NewEnum {
	e := &NewEnum{parent: m}
	e.label.parent = e
	return e
}

func (m *message) Parent() DefinitionContainer {
	return m.parent
}

func (m message) Document() *Document {
	return m.parent.Document()
}

func (m *message) String() string {
	return m.Document().Printer.Message(m)
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
	case *ReservedRange:
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

type NewMessage struct {
	label  Label
	parent DefinitionContainer
}

func (m *NewMessage) Label() *Label {
	return &m.label
}

func (m *NewMessage) InsertIntoParent() error {
	return m.parent.insertMessage(m.toMessage())
}

func (m *NewMessage) toMessage() *message {
	mm := &message{
		parent: m.parent,
		label: Label{
			value: m.label.Get(),
		},
	}
	mm.label.parent = mm
	return mm
}

func (m *NewMessage) Parent() DefinitionContainer {
	return m.parent
}

func (m *NewMessage) Document() *Document {
	return m.parent.Document()
}

func (m *NewMessage) String() string {
	return m.toMessage().String()
}

func (m *NewMessage) validateLabel(l *Label) error {
	return m.parent.validateLabel(l)
}
