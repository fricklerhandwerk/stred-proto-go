package core

type field struct {
	label      Label
	number     Number
	deprecated Flag
}

func (f *field) Label() *Label {
	return &f.label
}

func (f *field) Number() *Number {
	return &f.number
}

func (f *field) Deprecated() *Flag {
	return &f.deprecated
}

func (f *field) hasLabel(other *Label) bool {
	return f.label.hasLabel(other)
}

func (f *field) hasNumber(n FieldNumber) bool {
	return f.number.hasNumber(n)
}

func (f *field) validate() (err error) {
	err = f.label.validate()
	if err != nil {
		return
	}
	err = f.number.validate()
	if err != nil {
		return
	}
	return
}

type typedField struct {
	field
	_type Type
}

func (f *typedField) Type() *Type {
	return &f._type
}

func (f *typedField) validate() (err error) {
	if err = f.field.validate(); err != nil {
		return
	}
	if err = f._type.validate(); err != nil {
		return
	}
	return
}

type Field struct {
	field
	_type    Type
	repeated Flag
	parent   *message
}

func (r *Field) Type() *Type {
	if r._type.parent == nil {
		r._type.parent = r
	}
	return &r._type
}

func (r *Field) Repeated() *Flag {
	if r.repeated.parent == nil {
		r.repeated.parent = r
	}
	return &r.repeated
}

func (r *Field) InsertIntoParent() error {
	return r.parent.insertField(r)
}

func (r *Field) Parent() Message {
	return r.parent
}

func (r *Field) Document() *Document {
	return r.parent.Document()
}

func (r *Field) validateAsMessageField() (err error) {
	if err = r.field.validate(); err != nil {
		return
	}
	if err = r._type.validate(); err != nil {
		return
	}
	return
}

func (r *Field) validateLabel(l *Label) error {
	return r.parent.validateLabel(l)
}

func (r *Field) validateNumber(n FieldNumber) error {
	return r.parent.validateNumber(n)
}

func (r *Field) validateFlag(f *Flag) error {
	// TODO: if we ever have "safe mode" to prevent backwards-incompatible
	// changes, that is where errors whould happen
	switch f {
	case &r.deprecated:
		return nil
	case &r.repeated:
		return nil
	}
	return nil
}
