package protobuf

import (
	"errors"
	"fmt"
	"regexp"
)

type Protocol struct {
	_package *identifier
	imports  []_import
	services []service
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

type identifier string

func (i identifier) validate() error {
	pattern := "[a-zA-Z]([0-9a-zA-Z_])*"
	regex := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
	if !regex.MatchString(i.String()) {
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
	rpcs []rpc
}

type rpc struct {
	declaration
	requestType    *Message
	streamRequest  bool
	responseType   *Message
	streamResponse bool
}

type declaration struct {
	label identifier
}

func (d declaration) GetLabel() string {
	return d.label.String()
}

func (d *declaration) SetLabel(label string) error {
	ident := identifier(label)
	if err := ident.validate(); err != nil {
		return err
	}
	d.label = ident
	return nil
}

type container interface {
	getDefinitions() []definition
	insertDefinition(index uint, def definition) error
}

type definition interface {
	GetFields() []definitionField
	InsertField(index uint, f definitionField) error
}

// sum type for fields in definitions
type definitionField interface {
	_definitionField()
}

type field struct {
	definitionField

	declaration
	number     uint
	deprecated bool
}

func (f field) _definitionField() {}

func (f field) GetNumber() uint {
	return f.number
}

func (f *field) SetNumber(n uint) error {
	// TODO: validate based on parent
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

func (n number) _fieldNumber() {}

type numberRange struct {
	fieldNumber
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

type fieldNumber interface {
	_fieldNumber()
}

type ReservedNumbers struct {
	messageField
	enumField

	numbers []fieldNumber
}

func (r ReservedNumbers) Insert(index uint, n fieldNumber) error {
	panic("not implemented")
}

type ReservedLabels struct {
	messageField
	enumField

	labels []identifier
}

func (r ReservedLabels) Insert(index uint, n string) error {
	panic("not implemented")
}

type Message struct {
	definition
	container
	fieldType

	declaration
	fields      []messageField
	definitions []definition
	// ...
}

func (m Message) GetFields() []definitionField {
	out := make([]definitionField, len(m.fields))
	for i, v := range m.fields {
		out[i] = v.(definitionField)
	}
	return out
}

func (m *Message) InsertField(i uint, value definitionField) error {
	var (
		field messageField
		ok    bool
	)
	if field, ok = value.(messageField); !ok {
		return errors.New(fmt.Sprintf("field must be suitable for message, but is %T", value))
	}
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

type TypedField struct {
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

func (f TypedField) _messageField() {}

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

type oneOf struct {
	declaration
	messageField

	fields []TypedField
	// ...
}

func (o oneOf) getFields() []TypedField {
	panic("not implemented")
}

func (o oneOf) insertField(i uint, f TypedField) error {
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

// sum type for message fields
type messageField interface {
	definitionField

	_messageField()
}

type enum struct {
	definition
	fieldType

	allowAlias bool
	fields     []enumField
}

type enumValue struct {
	enumField

	field
}

// sum type for enum fields
type enumField interface {
	definitionField

	_enumField()
}
