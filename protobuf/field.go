package protobuf

import "errors"

type field struct {
	*label
	*number
	deprecated bool
	parent     Definition
}

// TODO: maybe this should be a pointer for easier distinction in UI
func (f *field) GetNumber() uint {
	return f.number.GetValue()
}

func (f *field) SetNumber(v uint) error {
	// TODO: numbers should probably be initialised as `nil` so UI can display it
	// as not set, so check for that here
	return f.number.SetValue(v)
}

func (f field) GetDeprecated() bool {
	return f.deprecated
}

func (f *field) SetDeprecated(b bool) {
	f.deprecated = b
}

func (f field) hasNumber(n fieldNumber) bool {
	return n != f.number && n.intersects(f.number)
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
	return r.parent.insertField(i, r)
}

func (r *repeatableField) validateAsMessageField() (err error) {
	err = r.parent.validateLabel(r.label)
	if err != nil {
		return err
	}
	err = r.parent.validateNumber(r.number)
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
	return o.parent.insertField(i, o)
}

func (o *oneOf) validateAsMessageField() error {
	panic("not implemented")
}

func (o oneOf) hasNumber(n fieldNumber) bool {
	for _, f := range o.fields {
		if f.hasNumber(n) {
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
	return m.parent.insertField(i, m)
}

func (m *mapField) validateAsMessageField() error {
	err := m.parent.validateLabel(m.label)
	if err != nil {
		return err
	}
	err = m.parent.validateNumber(m.number)
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
