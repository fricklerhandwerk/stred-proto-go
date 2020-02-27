package protobuf

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func NewDocument() *document {
	return &document{}
}

type document struct {
	_package    *identifier
	imports     []_import
	services    []service
	definitions []definition
}

func (p document) GetPackage() *identifier {
	return p._package
}

func (p *document) SetPackage(pkg string) error {
	ident := identifier(pkg)
	if err := ident.validate(); err != nil {
		return err
	}
	p._package = &ident
	return nil
}

func (p document) validateLabel(l identifier) error {
	if err := l.validate(); err != nil {
		return err
	}
	for _, d := range p.definitions {
		if d.GetLabel() == l.String() {
			return errors.New(fmt.Sprintf("label %s already declared for other %T", l.String(), d))
		}
	}
	for _, s := range p.services {
		if s.GetLabel() == l.String() {
			return errors.New(fmt.Sprintf("label %s already declared for a service", l.String()))
		}
	}

	return nil
}

func (p document) GetDefinitions() []definition {
	out := make([]definition, len(p.definitions))
	for i, v := range p.definitions {
		out[i] = v.(definition)
	}
	return out
}

func (p *document) InsertDefinition(i uint, d definition) error {
	if err := d.validateAsDefinition(); err != nil {
		return err
	}
	p.definitions = append(p.definitions, nil)
	copy(p.definitions[i+1:], p.definitions[i:])
	p.definitions[i] = d
	return nil
}

func (p *document) NewService() *service {
	return &service{
		label: label{
			parent: p,
		},
	}
}

func (p *document) NewMessage() *message {
	return &message{
		parent: p,
		label: label{
			parent: p,
		},
	}
}

func (p *document) NewEnum() enum {
	return enum{
		parent: p,
		label: label{
			parent: p,
		},
	}
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
	// TODO: rpc labels and rpc argument/return types share a namespace with
	// *unqualified* message/enum labels within a service. you can only have "rpc
	// Foo" and use "message Foo" as an argument/return type in one of the same
	// service's rpcs if you use the qualified message label "rpc Foo
	// (package.Foo)", but then you *must* have a package name
	label
	requestType    *message
	streamRequest  bool
	responseType   *message
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
	if err := d.parent.validateLabel(ident); err != nil {
		return err
	}
	d.label = ident
	return nil
}

type definitionContainer interface {
	declarationContainer

	GetDefinitions() []definition
	InsertDefinition(index uint, def definition) error
}

type definition interface {
	declaration
	declarationContainer

	validateNumber(fieldNumber) error
	validateAsDefinition() error
}

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

type reservedNumbers struct {
	numbers []fieldNumber
}

func (r reservedNumbers) Insert(index uint, n fieldNumber) error {
	panic("not implemented")
}

func (e reservedNumbers) validateAsEnumField() error {
	panic("not implemented")
}

func (e reservedNumbers) validateAsMessageField() error {
	panic("not implemented")
}

type reservedLabels struct {
	labels []identifier
}

func (r reservedLabels) Get() []string {
	panic("not implemented")
}

func (r reservedLabels) Insert(index uint, n string) error {
	panic("not implemented")
}

func (r reservedLabels) validateAsEnumField() error {
	panic("not implemented")
}

func (r reservedLabels) validateAsMessageField() error {
	panic("not implemented")
}

type message struct {
	label
	fields      []messageField
	definitions []definition
	parent      definitionContainer
}

func (m message) GetFields() []messageField {
	return m.fields
}

// TODO: this is a bad interface, as it requires checking that the parent of
// the inserted field is really this message. instead wie should have
// `field.Insert(uint) error {}`, which may call its parent, which is the right
// thing by construction, to do the work.
// doing it that way has the added benifit that self-validation semantics are
// contained in the child type instead of calling the child from here, which
// calls the parent again.
func (m *message) InsertField(i uint, field messageField) error {
	if err := field.validateAsMessageField(); err != nil {
		return err
	}
	m.fields = append(m.fields, nil)
	copy(m.fields[i+1:], m.fields[i:])
	m.fields[i] = field

	return nil
}

func (m *message) NewField() *typedField {
	return &typedField{
		field: field{
			parent: m,
			label: label{
				parent: m,
			},
		},
	}
}

func (m *message) NewMap() *mapField {
	return &mapField{
		typedField: *m.NewField(),
	}
}

func (m *message) NewOneOf() *oneOf {
	return &oneOf{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

func (m *message) NewMessage() *message {
	return &message{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

func (m *message) NewEnum() *enum {
	return &enum{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

func (m message) GetDefinitions() []definition {
	panic("not implemented")
}

func (m *message) InsertDefinition(i uint, d definition) error {
	panic("not implemented")
}

func (m message) validateAsDefinition() (err error) {
	err = m.parent.validateLabel(identifier(m.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return
	}
	return
}

func (m message) validateLabel(l identifier) error {
	// TODO: if the policy now develops such that everything is validated by its
	// parent, this should also be done by a function independent of the
	// identifier. this makes the whole extra type unnecessary.
	if err := l.validate(); err != nil {
		return err
	}
	for _, f := range m.fields {
		switch field := f.(type) {
		case *typedField:
			if field.GetLabel() == l.String() {
				return errors.New(fmt.Sprintf("label %q already declared", l.String()))
			}
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (m message) validateNumber(f fieldNumber) error {
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

func (m message) validateNumberSingle(n number) error {
	if n < 1 {
		return errors.New("message field number must be >= 1")
	}
	for _, f := range m.fields {
		switch field := f.(type) {
		case *typedField:
			if field.GetNumber() == uint(n) {
				return errors.New(fmt.Sprintf("field number %d already in use", uint(n)))
			}
		case *reservedNumbers:
			panic("not implemented")
		case *reservedLabels:
			continue
		default:
			panic(fmt.Sprintf("unhandled message field type %T", f))
		}
	}
	return nil
}

func (m message) validateNumberRange(n numberRange) error {
	panic("not implemented")
}

func (m *message) _fieldType() {}

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

func (f typedField) validateAsMessageField() (err error) {
	err = f.parent.validateLabel(identifier(f.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	err = f.parent.validateNumber(number(f.GetNumber()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	if f._type == nil {
		return errors.New("message field type not set")
	}
	return nil
}

type fieldType interface {
	_fieldType()
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

func (o oneOf) validateAsMessageField() error {
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

type enum struct {
	label
	allowAlias bool
	fields     []enumField
	parent     definitionContainer
}

func (e *enum) SetAlias(b bool) error {
	if !b && e.allowAlias {
		// check if aliasing is in place
		numbers := make(map[uint]bool, len(e.fields))
		for _, field := range e.fields {
			switch f := field.(type) {
			case *enumeration:
				n := f.GetNumber()
				if numbers[n] {
					lines := []string{
						fmt.Sprintf("field number %d is used multiple times.", n),
						fmt.Sprintf("remove aliasing before setting %q.", "allow_alias = false"),
					}
					return errors.New(strings.Join(lines, " "))
				}
				numbers[n] = true
			case fieldNumber:
				continue
			default:
				panic(fmt.Sprintf("unhandled enum field type %T", f))
			}
		}
	}
	e.allowAlias = b
	return nil
}

func (e enum) GetAlias() bool {
	return e.allowAlias
}

func (e enum) GetFields() []enumField {
	return e.fields
}

func (e *enum) InsertField(i uint, field enumField) error {
	if err := field.validateAsEnumField(); err != nil {
		return err
	}
	// <https://github.com/golang/go/wiki/SliceTricks#insert>
	// <https://stackoverflow.com/a/46130603/5147619>
	e.fields = append(e.fields, nil)
	copy(e.fields[i+1:], e.fields[i:])
	e.fields[i] = field

	return nil
}

func (e *enum) NewField() *enumeration {
	return &enumeration{
		field: field{
			parent: e,
			label: label{
				parent: e,
			},
		},
	}
}

func (e enum) validateLabel(l identifier) error {
	if err := l.validate(); err != nil {
		return err
	}
	for _, f := range e.fields {
		switch field := f.(type) {
		case *enumeration:
			if field.GetLabel() == l.String() {
				return errors.New(fmt.Sprintf("label %s already declared", l.String()))
			}
		case *reservedLabels:
			for _, r := range field.Get() {
				if r == l.String() {
					return errors.New(fmt.Sprintf("label %s already declared", l.String()))
				}
			}
		case *reservedNumbers:
			continue
		default:
			panic(fmt.Sprintf("unhandled enum field type %T", f))
		}
	}
	return nil
}

func (e enum) validateNumber(f fieldNumber) error {
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

func (e enum) validateNumberSingle(n number) error {
	for _, f := range e.fields {
		switch field := f.(type) {
		case *enumeration:
			if !e.allowAlias && field.GetNumber() == uint(n) {
				lines := []string{
					fmt.Sprintf("field number %d already in use.", uint(n)),
					fmt.Sprintf("set %q to allow multiple labels for one number.", "allow_alias = true"),
				}
				return errors.New(strings.Join(lines, " "))
			}
		case *reservedNumbers:
			panic("not implemented")
		case *reservedLabels:
			continue
		default:
			panic(fmt.Sprintf("unhandled field number type %T", f))
		}
	}
	return nil
}

func (e enum) validateNumberRange(n numberRange) error {
	panic("not implemented")
}

func (e enum) validateAsDefinition() (err error) {
	err = e.parent.validateLabel(identifier(e.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return
	}
	return
}

func (e *enum) _fieldType() {}

type enumeration struct {
	field
}

func (e enumeration) validateAsEnumField() (err error) {
	err = e.parent.validateLabel(identifier(e.GetLabel()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	err = e.parent.validateNumber(number(e.GetNumber()))
	if err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	return nil
}

type enumField interface {
	validateAsEnumField() error
}
