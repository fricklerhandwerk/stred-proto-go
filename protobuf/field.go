package protobuf

import (
	"errors"
	"fmt"
)

type field struct {
	*label
	*number
	deprecated bool
}

func (f field) Deprecated() bool {
	return f.deprecated
}

func (f *field) SetDeprecated(b bool) error {
	f.deprecated = b
	return nil
}

func (f field) hasNumber(n FieldNumber) bool {
	return n != f.number && n.intersects(f.number)
}

type newOneOfField struct {
	oneOfField *oneOfField
}

func (f *newOneOfField) MaybeLabel() *string {
	return f.oneOfField.maybeLabel()
}

func (f *newOneOfField) SetLabel(l string) (err error) {
	if f.oneOfField.label == nil {
		f.oneOfField.label = &label{
			parent: f.oneOfField.parent,
		}
		defer func() {
			if err != nil {
				f.oneOfField.label = nil
			}
		}()
	}
	return f.oneOfField.SetLabel(l)
}

func (f *newOneOfField) MaybeNumber() *uint {
	return f.oneOfField.maybeNumber()
}

func (f *newOneOfField) SetNumber(n uint) (err error) {
	if f.oneOfField.number == nil {
		f.oneOfField.number = &number{
			parent: f.oneOfField.parent,
		}
		defer func() {
			if err != nil {
				f.oneOfField.number = nil
			}
		}()
	}
	return f.oneOfField.SetNumber(n)
}

func (f *newOneOfField) MaybeType() FieldType {
	return f.oneOfField.maybeType()
}

func (f *newOneOfField) Deprecated() bool {
	return f.oneOfField.Deprecated()
}

func (f *newOneOfField) SetDeprecated(b bool) error {
	return f.oneOfField.SetDeprecated(b)
}

func (f *newOneOfField) SetType(t FieldType) error {
	return f.oneOfField.SetType(t)
}

func (f *newOneOfField) InsertIntoParent(i uint) error {
	if err := f.oneOfField.parent.insertField(i, f.oneOfField); err != nil {
		return err
	}
	f.oneOfField = &oneOfField{
		parent: f.oneOfField.parent,
	}
	return nil
}

type oneOfField struct {
	typedField
	parent *oneOf
}

func (o oneOfField) validateAsOneOfField() (err error) {
	err = o.parent.validateLabel(o.label)
	if err != nil {
		return err
	}

	if o.number == nil {
		return fmt.Errorf("field number not set")
	}
	err = o.parent.validateNumber(o.number)
	if err != nil {
		return err
	}
	if o._type == nil {
		return errors.New("message field type not set")
	}
	return nil
}

type typedField struct {
	field
	_type FieldType
}

func (f typedField) Type() FieldType {
	return f._type
}

func (f typedField) maybeType() FieldType {
	return f._type
}

func (f *typedField) SetType(t FieldType) error {
	f._type = t
	return nil
}

type labelContainer interface {
	validateLabel(Label) error
}

type fieldContainer interface {
	numberContainer
	labelContainer
}

type newRepeatableField struct {
	repeatableField *repeatableField
}

func (r *newRepeatableField) InsertIntoParent(i uint) error {
	if err := r.repeatableField.parent.insertField(i, r.repeatableField); err != nil {
		return err
	}
	r.repeatableField = &repeatableField{
		parent: r.repeatableField.parent,
	}
	return nil
}

func (r *newRepeatableField) MaybeLabel() *string {
	return r.repeatableField.maybeLabel()
}

func (r *newRepeatableField) SetLabel(l string) (err error) {
	if r.repeatableField.label == nil {
		r.repeatableField.label = &label{
			parent: r.repeatableField.parent,
		}
		defer func() {
			if err != nil {
				r.repeatableField.label = nil
			}
		}()
	}
	return r.repeatableField.SetLabel(l)
}

func (r *newRepeatableField) MaybeNumber() *uint {
	return r.repeatableField.maybeNumber()
}

func (r *newRepeatableField) SetNumber(n uint) (err error) {
	if r.repeatableField.number == nil {
		r.repeatableField.number = &number{
			parent: r.repeatableField.parent,
		}
		defer func() {
			if err != nil {
				r.repeatableField.number = nil
			}
		}()
	}
	return r.repeatableField.SetNumber(n)
}

func (r *newRepeatableField) MaybeType() FieldType {
	return r.repeatableField.maybeType()
}

func (r *newRepeatableField) SetType(t FieldType) error {
	return r.repeatableField.SetType(t)
}

func (r *newRepeatableField) Deprecated() bool {
	return r.repeatableField.Deprecated()
}

func (r *newRepeatableField) SetDeprecated(b bool) error {
	return r.repeatableField.SetDeprecated(b)
}

func (r *newRepeatableField) Repeated() bool {
	return r.repeatableField.Repeated()
}

func (r *newRepeatableField) SetRepeated(b bool) error {
	return r.repeatableField.SetRepeated(b)
}

type repeatableField struct {
	typedField

	parent   *message
	repeated bool
}

func (r repeatableField) Repeated() bool {
	return r.repeated
}

func (r repeatableField) SetRepeated(repeat bool) error {
	r.repeated = repeat
	return nil
}

func (r *repeatableField) validateAsMessageField() (err error) {
	err = r.parent.validateLabel(r.label)
	if err != nil {
		return err
	}

	if r.number == nil {
		return fmt.Errorf("field number not set")
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

type newOneOf struct {
	oneOf *oneOf
}

func (o *newOneOf) InsertIntoParent(i uint) error {
	if err := o.oneOf.parent.insertField(i, o.oneOf); err != nil {
		return err
	}
	o.oneOf = &oneOf{
		parent: o.oneOf.parent,
	}
	return nil
}

func (o *newOneOf) MaybeLabel() *string {
	return o.oneOf.maybeLabel()
}

func (o *newOneOf) SetLabel(l string) (err error) {
	if o.oneOf.label == nil {
		o.oneOf.label = &label{
			parent: o.oneOf.parent,
		}
		defer func() {
			if err != nil {
				o.oneOf.label = nil
			}
		}()
	}
	return o.oneOf.SetLabel(l)
}

func (o *newOneOf) NewField() NewOneOfField {
	return o.oneOf.NewField()
}

func (o newOneOf) NumFields() uint {
	return o.oneOf.NumFields()
}

func (o newOneOf) Field(i uint) OneOfField {
	return o.oneOf.Field(i)
}

type oneOf struct {
	*label
	fields []*oneOfField
	parent *message
}

func (o oneOf) NumFields() uint {
	return uint(len(o.fields))
}

func (o oneOf) Field(i uint) OneOfField {
	return o.fields[i]
}

func (o *oneOf) NewField() NewOneOfField {
	return &newOneOfField{
		oneOfField: &oneOfField{parent: o},
	}
}

func (o *oneOf) insertField(i uint, f *oneOfField) error {
	if err := f.validateAsOneOfField(); err != nil {
		return err
	}
	o.fields = append(o.fields, nil)
	copy(o.fields[i+1:], o.fields[i:])
	o.fields[i] = f
	return nil
}

func (o oneOf) validateLabel(l *label) error {
	return o.parent.validateLabel(l)
}

func (o oneOf) validateNumber(n FieldNumber) error {
	if o.hasNumber(n) {
		return fmt.Errorf("field number %s already reserved", n)
	}
	return o.parent.validateNumber(n)
}

func (o *oneOf) validateAsMessageField() error {
	return o.parent.validateLabel(o.label)
}

func (o oneOf) hasNumber(n FieldNumber) bool {
	for _, f := range o.fields {
		if f.hasNumber(n) {
			return true
		}
	}
	return false
}

func (o oneOf) hasLabel(l *label) bool {
	for _, f := range o.fields {
		if f.hasLabel(l) {
			return true
		}
	}
	return o.label.hasLabel(l)
}

type newMapField struct {
	mapField *mapField
}

func (m *newMapField) InsertIntoParent(i uint) error {
	if err := m.mapField.parent.insertField(i, m.mapField); err != nil {
		return err
	}
	m.mapField = &mapField{
		parent: m.mapField.parent,
	}
	return nil
}

func (m *newMapField) MaybeLabel() *string {
	return m.mapField.maybeLabel()
}

func (m *newMapField) SetLabel(l string) (err error) {
	if m.mapField.label == nil {
		m.mapField.label = &label{
			parent: m.mapField.parent,
		}
		defer func() {
			if err != nil {
				m.mapField.label = nil
			}
		}()
	}
	return m.mapField.SetLabel(l)
}

func (m *newMapField) MaybeNumber() *uint {
	return m.mapField.maybeNumber()
}

func (m *newMapField) SetNumber(n uint) (err error) {
	if m.mapField.number == nil {
		m.mapField.number = &number{
			parent: m.mapField.parent,
		}
		defer func() {
			if err != nil {
				m.mapField.number = nil
			}
		}()
	}
	return m.mapField.SetNumber(n)
}

func (m *newMapField) MaybeKeyType() KeyType {
	return m.mapField.key
}

func (m *newMapField) SetKeyType(t KeyType) error {
	return m.mapField.SetKeyType(t)
}

func (m *newMapField) MaybeValueType() FieldType {
	return m.mapField.maybeType()
}

func (m *newMapField) SetValueType(t FieldType) error {
	return m.mapField.SetValueType(t)
}

func (m *newMapField) Deprecated() bool {
	return m.mapField.Deprecated()
}

func (m *newMapField) SetDeprecated(b bool) error {
	return m.mapField.SetDeprecated(b)
}

type mapField struct {
	typedField

	parent *message
	key    keyType
}

func (m mapField) GetKeyType() KeyType {
	return m.key
}

func (m mapField) SetKeyType(k KeyType) error {
	m.key = k.(keyType)
	return nil
}

func (m mapField) ValueType() FieldType {
	return m._type
}

func (m mapField) SetValueType(k FieldType) error {
	m._type = k
	return nil
}

func (m *mapField) validateAsMessageField() error {
	err := m.parent.validateLabel(m.label)
	if err != nil {
		return err
	}

	if m.number == nil {
		return fmt.Errorf("field number not set")
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
