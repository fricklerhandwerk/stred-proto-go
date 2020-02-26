package protobuf

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// TODO: maybe do not export any types at all, but just constructors such as
// `NewProtocol() protocol`. that way we can enforce setting a parent on types
// which need one, thus avoiding one source of API usage errors.
type Protocol struct {
	_package *identifier
	imports  []_import
	// weird naming rules...
	// 1. service labels share a namespace with message/enum labels
	// 2. rpc labels and rpc argument/return types share a namespace with
	//    *unqualified* message/enum labels.  you can only have "rpc Foo" and use
	//    "message Foo" as an argument/return type in one of your rpcs if you use
	//    the qualified message label "rpc Foo (package.Foo)", but then you
	//    *must* have a package name
	services    []service
	definitions []definition
}

func (p Protocol) GetPackage() *identifier {
	return p._package
}

func (p *Protocol) SetPackage(pkg string) error {
	ident := identifier(pkg)
	if err := ident.validate(); err != nil {
		return err
	}
	p._package = &ident
	return nil
}

func (p Protocol) validateLabel(l identifier) error {
	for _, d := range p.definitions {
		if d.GetLabel() == l.String() {
			return errors.New(fmt.Sprintf("label %s already declared", l.String()))
		}
	}
	return nil
}

func (p Protocol) getDefinitions() []definition {
	out := make([]definition, len(p.definitions))
	for i, v := range p.definitions {
		out[i] = v.(definition)
	}
	return out
}

func (p *Protocol) insertDefinition(i uint, d definition) error {
	panic("not implemented")
}

// TODO: probably there is no need to have an extra type here, and validation
// can be done in a function
type identifier string

func (i identifier) validate() error {
	pattern := "[a-zA-Z]([0-9a-zA-Z_])*"
	regex := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
	if !regex.MatchString(i.String()) {
		// TODO: there are at least two sources of errors which should be
		// differentiated by type: API caller and user. maybe API usage errors
		// should even result in a panic, since a nonsensical operation due to
		// broken implementation simply must not be allowed.
		return errors.New(fmt.Sprintf("Identifier must match %s", pattern))
	}
	return nil
}

func (i *identifier) String() string {
	if i == nil {
		return ""
	}
	return string(*i)
}

type _import struct {
	path   string
	public bool
}

func (i _import) SetPath(string) error {
	panic("not implemented")
}

func (i _import) SetPublic(b bool) {
	panic("not implemented")
}

type service struct {
	label
	rpcs []rpc
}

type rpc struct {
	label
	requestType    *Message
	streamRequest  bool
	responseType   *Message
	streamResponse bool
}

type declarationContainer interface {
	validateLabel(identifier) error
}

type declaration interface {
	GetLabel() string
	SetLabel(string) error
}

type label struct {
	label  identifier
	parent declarationContainer
}

func (d label) GetLabel() string {
	return d.label.String()
}

func (d *label) SetLabel(label string) error {
	ident := identifier(label)
	if err := ident.validate(); err != nil {
		return err
	}
	if d.parent == nil {
		return errors.New("declaration has no parent")
	}
	if err := d.parent.validateLabel(ident); err != nil {
		return err
	}
	d.label = ident
	return nil
}

type definitionContainer interface {
	declarationContainer

	getDefinitions() []definition
	insertDefinition(index uint, def definition) error
}

type definition interface {
	declaration
	declarationContainer

	SetParent(p definitionContainer) error
	validateNumber(fieldNumber) error
}

type field struct {
	label
	number     uint
	deprecated bool
	parent     definition
}

func (f *field) SetParent(d definition) error {
	if d == nil {
		return errors.New("parent must not be nil")
	}
	f.label.parent = d
	f.parent = d
	return nil
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

type number uint

func (n number) intersects(other []fieldNumber) bool {
	for _, o := range other {
		switch v := o.(type) {
		case number:
			if n == o {
				return true
			}
		case numberRange:
			if uint(n) >= v.start || uint(n) <= v.end {
				return true
			}
		default:
			panic(fmt.Sprintf("unhandled fieldNumber type %T", v))
		}
	}
	return false
}

type numberRange struct {
	start uint
	end   uint
}

func (r numberRange) GetStart() uint {
	return r.start
}

func (r numberRange) SetStart(s uint) error {
	r.start = s
	// TODO: validate
	return nil
}

func (r numberRange) GetEnd() uint {
	return r.end
}

func (r numberRange) SetEnd(e uint) error {
	r.end = e
	// TODO: validate
	return nil
}

func (r numberRange) intersects(other []fieldNumber) bool {
	for _, o := range other {
		switch v := o.(type) {
		case number:
			if uint(v) >= r.start || uint(v) <= r.end {
				return true
			}
		case numberRange:
			if (v.start >= r.start && v.start <= r.end) || (v.end >= r.start && v.end <= r.end) {
				return true
			}
		default:
			panic(fmt.Sprintf("unhandled fieldNumber type %T", v))
		}
	}
	return false
}

type fieldNumber interface {
	intersects([]fieldNumber) bool
}

type ReservedNumbers struct {
	numbers []fieldNumber
}

func (r ReservedNumbers) Insert(index uint, n fieldNumber) error {
	panic("not implemented")
}

func (e ReservedNumbers) validateAsEnumField() error {
	panic("not implemented")
}

func (e ReservedNumbers) validateAsMessageField() error {
	panic("not implemented")
}

type ReservedLabels struct {
	labels []identifier
}

func (r ReservedLabels) Get() []string {
	panic("not implemented")
}

func (r ReservedLabels) Insert(index uint, n string) error {
	panic("not implemented")
}

func (r ReservedLabels) validateAsEnumField() error {
	panic("not implemented")
}

func (r ReservedLabels) validateAsMessageField() error {
	panic("not implemented")
}

type Message struct {
	label
	fields      []messageField
	definitions []definition
	parent      definitionContainer
}

func (m *Message) SetParent(d definitionContainer) error {
	if d == nil {
		return errors.New("parent must not be nil")
	}
	if v, ok := d.(*Message); ok && m == v {
		return errors.New("message cannot be parent of itself")
	}
	m.label.parent = d
	m.parent = d
	return nil
}

func (m Message) GetFields() []messageField {
	return m.fields
}

func (m *Message) InsertField(i uint, field messageField) error {
	// TODO: let field self-validate
	switch f := field.(type) {
	case TypedField:
		// <https://github.com/golang/go/wiki/SliceTricks#insert>
		// <https://stackoverflow.com/a/46130603/5147619>
		m.fields = append(m.fields, nil)
		copy(m.fields[i+1:], m.fields[i:])
		m.fields[i] = f
	default:
		panic(fmt.Sprintf("unhandled message field type %T", f))
	}
	return nil
}

func (m Message) getDefinitions() []definition {
	panic("not implemented")
}

func (m *Message) insertDefinition(i uint, d definition) error {
	panic("not implemented")
}

func (m Message) validateLabel(l identifier) error {
	for _, f := range m.fields {
		switch field := f.(type) {
		case TypedField:
			if field.GetLabel() == l.String() {
				return errors.New(fmt.Sprintf("label %q already declared", l.String()))
			}
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (m Message) validateNumber(f fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch n := f.(type) {
	case number:
		return m.validateNumberSingle(n)
	case numberRange:
		return m.validateNumberRange(n)
	default:
		panic(fmt.Sprintf("unhandled field number type %T", f))
	}
}

func (m Message) validateNumberSingle(n number) error {
	for _, f := range m.fields {
		switch field := f.(type) {
		case TypedField:
			if field.GetNumber() == uint(n) {
				return errors.New(fmt.Sprintf("field number %d already in use", uint(n)))
			}
		case ReservedNumbers:
			panic("not implemented")
		case ReservedLabels:
			continue
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (m Message) validateNumberRange(n numberRange) error {
	panic("not implemented")
}

func (m *Message) _fieldType() {}

type TypedField struct {
	messageField

	field
	_type fieldType
}

func (f TypedField) GetType() fieldType {
	return f._type
}

func (f *TypedField) SetType(t fieldType) {
	f._type = t
}

type fieldType interface {
	_fieldType()
}

type repeatableField struct {
	TypedField
	messageField

	repeated bool
}

func (r repeatableField) setRepeated(repeat bool) {
	r.repeated = repeat
}

func (r repeatableField) getRepeated() bool {
	return r.repeated
}

type OneOf struct {
	messageField

	label
	fields []TypedField
	// ...
}

func (o OneOf) GetFields() []TypedField {
	panic("not implemented")
}

func (o OneOf) InsertField(i uint, f TypedField) error {
	panic("not implemented")
}

type mapField struct {
	messageField
	TypedField

	key keyType
}

func (m mapField) getKeyType() keyType {
	return m.key
}

func (m mapField) setKeyType(k keyType) error {
	panic("not implemented")
}

type keyType string

const (
	Int32    keyType = "int32"
	Int64    keyType = "int64"
	Uint32   keyType = "uint32"
	Uint64   keyType = "uint64"
	Sint32   keyType = "sint32"
	Sint64   keyType = "sint64"
	Fixed32  keyType = "fixed32"
	Fixed64  keyType = "fixed64"
	Sfixed32 keyType = "sfixed32"
	Sfixed64 keyType = "sfixed64"
	Bool     keyType = "bool"
	String   keyType = "string"
)

func (k keyType) _fieldType() {}

type valueType string

const (
	Double valueType = "double"
	Float  valueType = "float"
	Bytes  valueType = "bytes"
)

func (v valueType) _fieldType() {}

type messageField interface {
	validateAsMessageField() error
}

type Enum struct {
	label
	allowAlias bool
	fields     []enumField
	parent     definitionContainer
}

func (e *Enum) AllowAlias(b bool) error {
	if !b && e.allowAlias {
		// check if aliasing is in place
		numbers := make(map[uint]bool, len(e.fields))
		for _, field := range e.fields {
			switch f := field.(type) {
			case Enumeration:
				n := f.GetNumber()
				if numbers[n] {
					lines := []string{fmt.Sprintf(
						"field number %d is used multiple times.", n),
						"remove aliasing before disallowing it.",
					}
					return errors.New(strings.Join(lines, " "))
				}
				numbers[n] = true
			default:
				continue
			}
		}
	}
	e.allowAlias = b
	return nil
}

func (e Enum) GetFields() []enumField {
	return e.fields
}

func (e *Enum) InsertField(i uint, field enumField) error {
	// TODO: let field self-validate
	switch f := field.(type) {
	case Enumeration:
		// <https://github.com/golang/go/wiki/SliceTricks#insert>
		// <https://stackoverflow.com/a/46130603/5147619>
		e.fields = append(e.fields, nil)
		copy(e.fields[i+1:], e.fields[i:])
		e.fields[i] = f
	default:
		panic(fmt.Sprintf("unhandled message field type %T", f))
	}
	return nil
}

func (e *Enum) SetParent(d definitionContainer) error {
	if d == nil {
		return errors.New("parent is nil")
	}
	e.label.parent = d
	e.parent = d
	return nil
}

func (e Enum) validateLabel(l identifier) error {
	for _, f := range e.fields {
		switch field := f.(type) {
		case Enumeration:
			if field.GetLabel() == l.String() {
				return errors.New(fmt.Sprintf("label %s already declared", l.String()))
			}
		case ReservedLabels:
			for _, r := range field.Get() {
				if r == l.String() {
					return errors.New(fmt.Sprintf("label %s already declared", l.String()))
				}
			}
		case ReservedNumbers:
			continue
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (e Enum) validateNumber(f fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch n := f.(type) {
	case number:
		return e.validateNumberSingle(n)
	case numberRange:
		return e.validateNumberRange(n)
	default:
		panic(fmt.Sprintf("unhandled field number type %T", f))
	}
}

func (e Enum) validateNumberSingle(n number) error {
	for _, f := range e.fields {
		switch field := f.(type) {
		case Enumeration:
			if !e.allowAlias && field.GetNumber() == uint(n) {
				lines := []string{
					fmt.Sprintf("field number %d already in use.", uint(n)),
					fmt.Sprintf("set %q to allow multiple labels for one number.", "allow_alias = true"),
				}
				return errors.New(strings.Join(lines, " "))
			}
		case ReservedNumbers:
			panic("not implemented")
		case ReservedLabels:
			continue
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (e Enum) validateNumberRange(n numberRange) error {
	panic("not implemented")
}

func (e *Enum) _fieldType() {}

type Enumeration struct {
	field
}

func (e Enumeration) validateAsEnumField() error {
	panic("not implemented")
}

type enumField interface {
	validateAsEnumField() error
}
