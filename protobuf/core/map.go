package core

import (
	"errors"
)

type KeyType struct {
	value  MapKeyType
	parent *Map
}

func (t KeyType) Get() MapKeyType {
	return t.value
}

func (t *KeyType) Set(value MapKeyType) error {
	t.value = value
	// TODO: checks in "safe mode"
	return nil
}

func (t KeyType) Parent() *Map {
	return t.parent
}

type Map struct {
	typedField
	keyType KeyType

	parent *message
}

func (m *Map) KeyType() *KeyType {
	return &m.keyType
}

func (m *Map) InsertIntoParent() error {
	return m.parent.insertField(m)
}

func (m *Map) Parent() Message {
	return m.parent
}

func (m *Map) Document() *Document {
	return m.parent.Document()
}

func (m *Map) validateAsMessageField() error {
	if err := m.validate(); err != nil {
		return err
	}
	if m.keyType.value == nil {
		return errors.New("map value type not set")
	}
	return nil
}

func (m *Map) validateLabel(l *Label) error {
	return m.parent.validateLabel(l)
}

func (m *Map) validateNumber(n FieldNumber) error {
	return m.parent.validateNumber(n)
}

func (m *Map) validateFlag(*Flag) error {
	return nil
}
