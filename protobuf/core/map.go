package core

import (
	"errors"
)

type KeyType struct {
	value  MapKeyType
	parent *mapField
}

func (t KeyType) Get() MapKeyType {
	return t.value
}

func (t *KeyType) Set(value MapKeyType) error {
	t.value = value
	// TODO: checks in "safe mode"
	return nil
}

func (t KeyType) Parent() Map {
	return t.parent
}

type mapField struct {
	typedField
	keyType KeyType

	parent *message
}

func (m *mapField) KeyType() *KeyType {
	if m.keyType.parent == nil {
		m.keyType.parent = m
	}
	return &m.keyType
}

func (m *mapField) InsertIntoParent() error {
	return m.parent.insertField(m)
}

func (m *mapField) Parent() Message {
	return m.parent
}

func (m *mapField) validateAsMessageField() error {
	if err := m.validate(); err != nil {
		return err
	}
	if m.keyType.value == nil {
		return errors.New("map value type not set")
	}
	return nil
}

func (m *mapField) validateLabel(l *Label) error {
	return m.parent.validateLabel(l)
}

func (m *mapField) validateNumber(n FieldNumber) error {
	return m.parent.validateNumber(n)
}

func (m *mapField) validateFlag(*Flag) error {
	return nil
}
