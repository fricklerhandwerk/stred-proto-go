package core

// this package captures the semantics of protobuf documents, disregarding
// presentational matters entierly.

// the design philosphy is to reduce input parameters to identifier strings,
// field numbers, boolean flags, built-in and valid user-defined message types.
// it must be impossible to update the document with an invalid value, and
// ideally impossible to even supply an invald value.

// a primary design decision following this principle is that once-valid values
// cannot be invalidated. for example setting an unset field to some valid
// value means that this field's value can only be changed to another valid
// one, but not unset any more.

// additionally consumers are not supposed to create values of any of these
// types other than an empty `Document`, which will properly initialise and
// produce all subsidiary items. these are therefore only visible through the
// interfaces they implement. the consumer is relieved of any additional setup
// requirements. admittedly this is not very idiomatic Go, but exporting types
// and writing in documentation that they should not be used is of much less
// use than a compiler-enforced constraint.

// these interfaces do not only model external behavior, but also guide the
// implementation by constraining internal relations to some extent. they
// also reduce cognitive overhead by centralising definitions of required
// behavior in one location, thus reducing opportunities for programmer
// error when used as a weak replacement for strong types.
// it is not necessary at all to it that way. most unexported methods can be
// left out, as they are specific to the implementing types, and internally
// we could instead specify interfaces for the required subsets of behavior.
// except for `KeyType` and `ValueType`, public interfaces do not have to be
// sealed to prevent misuse. but there would be little added convenience for
// the effort to create a different facade. removing all other private
// methods would merely improve readability for consumers, and only if they
// look at this file. even if it were possible to completely replace the
// implementation - which it is not, due to those two mandatory sealed
// interfaces - there is not much point to it, as the library is
// light-weight and free of side-effects.
//
// the layout follows a roughly hierarchic pattern: the root object is
// `Document`, from which we can create tentative subsidiary items, which are
// defined below. their name is prefixed with `New`, and for each we can set
// required attributes, and insert the object into the parent it was created
// by, at a specified index in the parent's collection.
// tentative objects are not tracked by the document and thus may not have
// children. first, doing that would complicate the necessary types to a degree
// Go does not afford to practically handle.  second, it prevents a problem
// with dangling references, since by design the document does not account for
// tentative elements - the document only cares about semantically correct
// structures, so tentative ones cannot even be referenced other than on
// creation, and cannot be passed into the document.
// insertion into the parent means emptying and thus invalidating the
// tentative object. this is mostly for implemenation simplicity; we could also
// deep-copy the new structure into its parent to have side-effect free
// handling of the new object, but this is much additional work without much
// value for now. to continue operating on the object you have to fetch it
// back from its parent, which will give you a more convenient interface.
// on the tentative variants, attribute accessors have names prefixed with
// `Maybe` and return pointers, signifying that a value may not have been set
// yet. if not `nil`, pointers point to a copy of the original value, as all
// mutation must go through the appropriate interface. where interfaces are
// returned, consumers must check for `nil`. objects retrieved from a parent
// will always return values or non-nil interfaces, as they have necessarily
// been validated prior.

type RepeatableType interface {
	Get() ValueType
	Set(ValueType) error
	Repeated() *Flag

	Parent() Field
}

type KeyType interface {
	Get() MapKeyType
	Set(MapKeyType) error

	Parent() Map
}

// ValueType identifier. Valid values are all built-in types and `Message`s.
// Note that API consumers cannot create types which implement interface
// `Message`, and the package will only emit properly constructed `Message`s.
type ValueType interface {
	_isValueType()
}

// MapKeyType identifier.
// Valid values are all built-in types except:
// - `double`
// - `float`
// - `bytes`
type MapKeyType interface {
	_isKeyType()
}

type Document interface {
	Package() *Package

	Imports() []*Import
	NewImport() *Import

	Services() []*Service
	NewService() *Service

	Messages() []Message
	NewMessage() NewMessage

	Enums() []Enum
	NewEnum() NewEnum

	validateLabel(*Label) error
	insertImport(*Import) error
	insertService(*Service) error
	insertMessage(*message) error
	insertEnum(*enum) error
}

type MessageType interface {
	Get() Message
	Set(Message) error
	Stream() *Flag

	Parent() *RPC
}

type Definition interface {
	Label() *Label

	NewReservedNumber() ReservedNumber
	NewReservedRange() ReservedRange
	NewReservedLabel() ReservedLabel

	Parent() DefinitionContainer

	validate() error
	hasLabel(*Label) bool
	validateLabel(*Label) error
	validateNumber(FieldNumber) error

	addReference(*Type)
	removeReference(*Type)
}

type DefinitionContainer interface {
	Messages() []Message
	NewMessage() NewMessage

	Enums() []Enum
	NewEnum() NewEnum

	validateLabel(*Label) error
	insertMessage(*message) error
	insertEnum(*enum) error
}

type NewMessage interface {
	Label() *Label

	InsertIntoParent() error
	Parent() DefinitionContainer
}

type Message interface {
	Label() *Label
	Fields() []MessageField

	NewField() Field
	NewMap() Map
	NewOneOf() OneOf
	NewReservedNumber() ReservedNumber
	NewReservedRange() ReservedRange
	NewReservedLabel() ReservedLabel

	Messages() []Message
	NewMessage() NewMessage

	Enums() []Enum
	NewEnum() NewEnum

	Parent() DefinitionContainer

	validate() error
	hasLabel(*Label) bool
	validateLabel(*Label) error
	validateNumber(FieldNumber) error
	insertField(MessageField) error

	addReference(*Type)
	removeReference(*Type)

	ValueType
}

type Field interface {
	Label() *Label
	Number() *Number
	Deprecated() *Flag
	Type() *Type
	Repeated() *Flag

	InsertIntoParent() error
	Parent() Message

	MessageField
}

type Map interface {
	Label() *Label
	Number() *Number
	Deprecated() *Flag
	KeyType() KeyType
	Type() *Type

	InsertIntoParent() error
	Parent() Message

	MessageField
}

type OneOf interface {
	Label() *Label

	Fields() []OneOfField
	NewField() OneOfField

	InsertIntoParent() error
	Parent() Message

	MessageField
}

type OneOfField interface {
	Label() *Label
	Number() *Number
	Type() *Type
	Deprecated() *Flag

	InsertIntoParent() error
	Parent() OneOf
}

type ReservedRange interface {
	Start() *Number
	End() *Number

	InsertIntoParent() error
	Parent() Definition

	// TODO: self-validate
	//validate() error

	FieldNumber

	validateAsMessageField() error
	validateAsEnumField() error
	hasLabel(*Label) bool
	hasNumber(FieldNumber) bool
}

type ReservedNumber interface {
	Get() *uint
	Set(uint) error

	InsertIntoParent() error
	Parent() Definition

	FieldNumber

	validateAsMessageField() error
	validateAsEnumField() error
	hasLabel(*Label) bool
	hasNumber(FieldNumber) bool
}

type ReservedLabel interface {
	Get() string
	Set(string) error

	InsertIntoParent() error
	Parent() Definition

	validateAsMessageField() error
	validateAsEnumField() error
	hasLabel(*Label) bool
	hasNumber(FieldNumber) bool
}

type MessageField interface {
	validateAsMessageField() error
	hasLabel(*Label) bool
	hasNumber(FieldNumber) bool
}

type NewEnum interface {
	Label() *Label

	InsertIntoParent() error
	Parent() DefinitionContainer
}

type Enum interface {
	Label() *Label
	AllowAlias() *Flag

	Fields() []EnumField

	NewVariant() Variant
	NewReservedNumber() ReservedNumber
	NewReservedRange() ReservedRange
	NewReservedLabel() ReservedLabel

	Parent() DefinitionContainer

	validate() error
	hasLabel(*Label) bool
	validateLabel(*Label) error
	validateNumber(FieldNumber) error
	insertField(EnumField) error

	addReference(*Type)
	removeReference(*Type)

	ValueType
}

type EnumField interface {
	validateAsEnumField() error
	hasLabel(*Label) bool
	hasNumber(FieldNumber) bool
}

type Variant interface {
	Label() *Label
	Number() *Number
	Deprecated() *Flag

	InsertIntoParent() error
	Parent() Enum

	EnumField
}

type FieldNumber interface {
	intersects(FieldNumber) bool
}
