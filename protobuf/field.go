package protobuf

import "errors"

type field struct {
	label
	number     uint
	deprecated bool
	parent     definition
}

func (f field) GetNumber() uint {
	return uint(f.number)
}

func (f *field) SetNumber(n uint) error {
	if err := f.parent.validateNumber(number(n)); err != nil {
		return err
	}
	f.number = n
	return nil
}

func (f field) GetDeprecated() bool {
	return f.deprecated
}

func (f *field) SetDeprecated(b bool) {
	f.deprecated = b
}

type typedField struct {
	field
	_type fieldType
}

type fieldType interface {
	_fieldType()
}

func (f typedField) GetType() fieldType {
	return f._type
}

func (f *typedField) SetType(t fieldType) {
	f._type = t
}

type repeatableField struct {
	typedField

	repeated bool
}

func (r repeatableField) setRepeated(repeat bool) {
	r.repeated = repeat
}

func (r repeatableField) getRepeated() bool {
	return r.repeated
}

func (r repeatableField) validateAsMessageField() (err error) {
	err = r.parent.validateLabel(identifier(r.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	err = r.parent.validateNumber(number(r.GetNumber()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	if r._type == nil {
		return errors.New("message field type not set")
	}
	return nil
}

type oneOf struct {
	label
	fields []typedField
	parent *message
}

func (o oneOf) GetFields() []typedField {
	panic("not implemented")
}

func (o oneOf) InsertField(i uint, f typedField) error {
	panic("not implemented")
}

func (o *oneOf) validateAsMessageField() error {
	panic("not implemented")
}

type mapField struct {
	typedField

	key keyType
}

func (m mapField) GetKeyType() keyType {
	return m.key
}

func (m mapField) SetKeyType(k keyType) {
	m.key = k
}

func (m *mapField) validateAsMessageField() error {
	err := m.parent.validateLabel(identifier(m.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	err = m.parent.validateNumber(number(m.GetNumber()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	if m._type == nil {
		return errors.New("message field type not set")
	}
	return nil
}
