package protobuf

import "errors"

type field struct {
	label
	number     uint
	deprecated bool
	parent     Definition
}

func (f field) GetNumber() uint {
	return f.number
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

func (f field) hasNumber(n fieldNumber) bool {
	return n.intersects(number(f.GetNumber()))
}

type typedField struct {
	field
	_type fieldType
}

func (f typedField) GetType() fieldType {
	return f._type
}

func (f *typedField) SetType(t fieldType) {
	f._type = t
}

type repeatableField struct {
	typedField

	parent   *message
	repeated bool
}

func (r repeatableField) SetRepeated(repeat bool) {
	r.repeated = repeat
}

func (r repeatableField) GetRepeated() bool {
	return r.repeated
}

func (r *repeatableField) InsertIntoParent(i uint) error {
	if err := r.validateAsMessageField(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	r.parent.insertField(i, r)
	return nil
}

func (r *repeatableField) validateAsMessageField() (err error) {
	err = r.parent.validateLabel(identifier(r.GetLabel()))
	if err != nil {
		return err
	}
	err = r.parent.validateNumber(number(r.GetNumber()))
	if err != nil {
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

func (o *oneOf) InsertIntoParent(i uint) error {
	if err := o.validateAsMessageField(); err != nil {
		return err
	}
	o.parent.insertField(i, o)
	return nil
}

func (o *oneOf) validateAsMessageField() error {
	panic("not implemented")
}

func (o oneOf) hasNumber(n fieldNumber) bool {
	for _, f := range o.fields {
		if n.intersects(number(f.GetNumber())) {
			return true
		}
	}
	return false
}

type mapField struct {
	typedField

	parent *message
	key    keyType
}

func (m mapField) GetKeyType() keyType {
	return m.key
}

func (m mapField) SetKeyType(k keyType) {
	m.key = k
}

func (m *mapField) InsertIntoParent(i uint) error {
	if err := m.validateAsMessageField(); err != nil {
		return err
	}
	m.parent.insertField(i, m)
	return nil
}

func (m *mapField) validateAsMessageField() error {
	err := m.parent.validateLabel(identifier(m.GetLabel()))
	if err != nil {
		return err
	}
	err = m.parent.validateNumber(number(m.GetNumber()))
	if err != nil {
		return err
	}
	if m._type == nil {
		return errors.New("map value type not set")
	}
	if m.key == "" {
		return errors.New("map key type not set")
	}
	return nil
}
