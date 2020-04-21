package protobuf

// the goal of this package's interface is to reduce input parameters to
// identifier strings, field numbers, boolean flags, built-in and valid
// user-defined message types, and positive array indices for ordering. it must
// be impossible to update the document with an invalid value, and ideally
// impossible to even supply an invald value.

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
// except for `KeyType` and `FieldType`, public interfaces do not have to be
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

type Identifier interface {
	Get() *string
	Set(string) error

	Parent() Labelled
}

type Labelled interface {
	validateLabel(*label) error
}

type Number interface {
	Get() *uint
	Set(uint) error

	Parent() Numbered

	FieldNumber
}

type Numbered interface {
	validateNumber(FieldNumber) error
}

type Flag interface {
	Get() bool
	Set(bool) error

	Parent() Flagged
}

type Flagged interface {
	_hasFlag()
}

type Type interface {
	Get() TypeID
	Set(TypeID) error

	Parent() Typed
}

type Typed interface {
	Type() Type
}

type KeyType interface {
	Get() KeyTypeID
	Set(KeyTypeID) error

	Parent() Map
}

// Field type identifier. Valid values are all built-in types and `Message`s.
// Note that API consumers cannot create types which implement interface
// `Message`, and the package will only emit properly constructed `Message`s.
type TypeID interface {
	_isFieldType()
}

// Map key type identifier.
// Valid values are all built-in types except:
// - `double`
// - `float`
// - `bytes`
type KeyTypeID interface {
	_isKeyType()
}

type Document interface {
	// Copy of list of top-level declarations in user-defined order
	Declarations() []TopLevelDeclaration

	Package() Identifier

	// Copy of list of imports, no guaranteed order
	Imports() []Import
	NewImport() Import

	// Copy of list of services, no guaranteed order
	Services() []Service
	NewService() Service

	// Copy of list of messages, no guaranteed order
	Messages() []Message
	NewMessage() NewMessage

	// Copy of list of enums, no guaranteed order
	Enums() []Enum
	NewEnum() NewEnum
}

type TopLevelDeclaration interface {
	_isDeclaration()
}

type Import interface {
	Path() Identifier
	Public() Flag

	InsertIntoParent(index uint) error
	Parent() Document
}

type Service interface {
	Label() Identifier

	// Copy of list of RPCs in user-defined order
	RPCs() []RPC
	NewRPC() RPC

	InsertIntoParent(index uint) error
	Parent() Document
}

type RPC interface {
	Label() Identifier

	Request() MessageType
	Response() MessageType

	InsertIntoParent(index uint) error
	Parent() Service
}

type MessageType interface {
	Get() Message
	Set(Message) error
	Stream() Flag

	Parent() RPC
}

type Definition interface {
	Label() Identifier

	NewReservedNumbers() ReservedNumbers
	NewReservedLabels() ReservedLabels

	InsertIntoParent(index uint) error
	Parent() DefinitionContainer

	hasLabel(*label) bool
}

type DefinitionContainer interface {
	Messages() []Message
	NewMessage() NewMessage

	Enums() []Enum
	NewEnum() NewEnum
}

type NewMessage interface {
	Label() Identifier

	InsertIntoParent(index uint) error
	Parent() DefinitionContainer
}

type Message interface {
	Label() Identifier

	Declarations() []MessageDeclaration

	Fields() []MessageField
	NewField() Field
	NewMap() Map
	NewOneOf() OneOf
	NewReservedNumbers() ReservedNumbers
	NewReservedLabels() ReservedLabels

	Messages() []Message
	NewMessage() NewMessage

	Enums() []Enum
	NewEnum() NewEnum

	InsertIntoParent(index uint) error
	Parent() DefinitionContainer

	TypeID
}

type MessageDeclaration interface {
	_isMessageDeclaration()
}

type Field interface {
	Label() Identifier
	Number() Number
	Type() Type
	Repeated() Flag
	Deprecated() Flag

	InsertIntoParent(index uint) error
	Parent() Message
}

type Map interface {
	Label() Identifier
	Number() Number
	KeyType() KeyType
	Type() Type
	Deprecated() Flag

	InsertIntoParent(index uint) error
	Parent() Message
}

type OneOf interface {
	Label() Identifier

	Fields() []OneOfField
	NewField() OneOfField

	InsertIntoParent(index uint) error
	Parent() Message
}

type OneOfField interface {
	Label() Identifier
	Number() Number
	Type() Type
	Deprecated() Flag

	InsertIntoParent(index uint) error
	Parent() OneOf
}

type ReservedNumbers interface {
	Numbers() []FieldNumber

	NewNumberRange() ReservedRange
	NewNumber() ReservedNumber

	InsertIntoParent(index uint) error
	Parent() Definition
}

type ReservedRange interface {
	Start() Number
	End() Number

	InsertIntoParent(index uint) error
	Parent() ReservedNumbers

	FieldNumber
}

type ReservedNumber interface {
	Get() *uint
	Set(uint) error

	InsertIntoParent(index uint) error
	Parent() ReservedNumbers
}

type ReservedLabels interface {
	Labels() []Identifier

	NewLabel() ReservedLabel

	InsertIntoParent(index uint) error
	Parent() Definition
}

type ReservedLabel interface {
	Get() *string
	Set(string) error

	InsertIntoParent(index uint) error
	Parent() ReservedLabels
}

type MessageField interface {
	_isMessageField()
}

type NewEnum interface {
	Label() Identifier
	AllowAlias() Flag

	InsertIntoParent(index uint) error
	Parent() DefinitionContainer
}

type Enum interface {
	Label() Identifier
	AllowAlias() Flag

	Fields() []EnumField

	NewVariant() Variant
	NewReservedNumbers() ReservedNumbers
	NewReservedLabels() ReservedLabels

	Parent() DefinitionContainer

	TypeID
}

type EnumField interface {
	_isEnumField()
}

type Variant interface {
	Label() Identifier
	Number() Number
	Deprecated() Flag

	InsertIntoParent(index uint) error
	Parent() Enum
}

type FieldNumber interface {
	intersects(FieldNumber) bool
}
