package protobuf

// this interface definition serves multiple purposes.
//
// 1. since the goal of this library is to be very hard to misuse, its
//    consumers are not supposed to create values of any of its types other
//    than an empty `Document`, which will properly initialise and produce all
//    subsidiary items. these are therefore only visible through the interfaces
//    they implement. the consumer is relieved of any additional setup
//    requirements. interacting with the interfaces by passing `string` and
//    `uint` parameters is almost everything there is.
//    admittedly this is not very idiomatic Go, but exporting types and writing
//    in documentation that they should not be used is of much less use than
//    a compiler-enforced constraint.
// 2. these interfaces do not only model external behavior, but also guide the
//    implementation by constraining internal relations to some extent. they
//    also reduce cognitive overhead by avoiding to define required behavior in
//    multiple locations, thus reducing opportunities for programmer error when
//    used as a weak replacement for strong types.
//    it is not necessary at all to it that way. most unexported methods can be
//    left out, as they are specific to the implementing types, and internally
//    we could instead specify interfaces for the required subsets of behavior.
//    except for `KeyType` and `FieldType`, public interfaces do not have to be
//    sealed to prevent misuse. but there would be little added convenience for
//    the effort to create a different facade. removing all other private
//    methods would merely improve readability for consumers, and only if they
//    look at this file. even if it were possible to completely replace the
//    implementation - which it is not, due to those two mandatory sealed
//    interfaces - there is not much point to it, as the library is
//    light-weight and free of side-effects.
//
// although there is significant overlap between items, all methods are
// declared explicitly. no non-trivial set of methods should be present more
// than twice anyway, such that deduplication would actually reduce
// readability.
// the layout follows a roughly hierarchic pattern: the root object is
// `Document`, from which we can create tentative subsidiary items, which are
// defined below. their name is prefixed with `New`, and for each we can set
// attributes, add children (which are likewise preliminary at creation), and
// insert the object into the parent it was created by, at a specified index in
// the parent's collection. insertion means deep-copy, and a tentative
// structure can be reused thereafter without side-effects - but it is
// recommended to fetch the same object back from its parent to have a more
// convenient interface. on the tentative variants, attribute accessors have
// names prefixed with `Maybe` and return pointers, signifying that a value may
// not have been set yet. if not `nil`, pointers point to a copy of the
// original value, as all mutation must go through the appropriate interface.
// where interfaces are returned, consumers must check for `nil`. objects
// retrieved from a parent will always return values or non-nil interfaces, as
// they have been validated prior.

type Document interface {
	MaybePackage() *string
	SetPackage(string) error

	NewService() NewService

	NumServices() uint
	Service(index uint) Service

	NumDefinitions() uint
	Definition(index uint) Definition

	NewMessage() NewMessage
	NewEnum() NewEnum

	validateLabel(Label) error

	insertService(index uint, service Service) error
	insertDefinition(index uint, definition Definition) error
}

type NewService interface {
	MaybeLabel() *string
	SetLabel(string) error

	NumRPCs() uint
	RPC(index uint) RPC

	NewRPC() NewRPC

	InsertIntoParent(index uint) error

	validateLabel(Label) error
	insertRPC(index uint, rpc RPC) error
}

type Service interface {
	Label() string
	SetLabel(string) error

	NumRPCs() uint
	RPC(index uint) RPC

	NewRPC() NewRPC

	validateLabel(Label) error
	insertRPC(index uint, rpc RPC) error

	hasLabel(Label) bool
	validateAsService() error
}

type NewRPC interface {
	MaybeLabel() *string
	SetLabel(string) error

	MaybeRequestType() Message
	SetRequestType(Message) error
	StreamRequest() bool
	SetStreamRequest(bool) error

	MaybeResponseType() Message
	SetResponseType(Message) error
	StreamResponse() bool
	SetStreamResponse(bool) error

	InsertIntoParent(index uint) error
}

type RPC interface {
	Label() string
	SetLabel(string) error

	RequestType() Message
	SetRequestType(Message) error
	StreamRequest() bool
	SetStreamRequest(bool) error

	ResponseType() Message
	SetResponseType(Message) error
	StreamResponse() bool
	SetStreamResponse(bool) error

	validateAsRPC() error
	hasLabel(Label) bool
}

type Definition interface {
	Label() string
	SetLabel(string) error

	NumFields() uint
	NewReservedNumbers() NewReservedNumbers
	NewReservedLabels() NewReservedLabels

	validateNumber(FieldNumber) error
	validateLabel(Label) error

	validateAsDefinition() error
	hasLabel(Label) bool
}

type NewMessage interface {
	MaybeLabel() *string
	SetLabel(string) error

	NumFields() uint
	Field(uint) MessageField

	NewField() NewField
	NewMap() NewMap
	NewOneOf() NewOneOf
	NewReservedNumbers() NewReservedNumbers
	NewReservedLabels() NewReservedLabels

	NumDefinitions() uint
	Definition(uint) Definition

	NewMessage() NewMessage
	NewEnum() NewEnum

	InsertIntoParent(index uint) error

	validateNumber(FieldNumber) error
	validateLabel(Label) error

	insertField(index uint, field MessageField) error
	insertDefinition(index uint, definition Definition) error
}

type Message interface {
	Label() string
	SetLabel(string) error

	NumFields() uint
	Field(uint) MessageField

	NewField() NewField
	NewMap() NewMap
	NewOneOf() NewOneOf
	NewReservedNumbers() NewReservedNumbers
	NewReservedLabels() NewReservedLabels

	NumDefinitions() uint
	Definition(uint) Definition

	NewMessage() NewMessage
	NewEnum() NewEnum

	validateNumber(FieldNumber) error
	validateLabel(Label) error

	insertField(uint, MessageField) error
	insertDefinition(index uint, def Definition) error

	validateAsDefinition() error
	hasLabel(Label) bool

	FieldType
}

type NewField interface {
	MaybeLabel() *string
	SetLabel(string) error

	MaybeNumber() *uint
	SetNumber(uint) error

	MaybeType() FieldType
	SetType(FieldType) error

	Repeated() bool
	SetRepeated(bool) error

	Deprecated() bool
	SetDeprecated(bool) error

	InsertIntoParent(index uint) error
}

type Field interface {
	Label() string
	SetLabel(string) error

	Number() uint
	SetNumber(uint) error

	Type() FieldType
	SetType(FieldType) error

	Repeated() bool
	SetRepeated(bool) error

	Deprecated() bool
	SetDeprecated(bool) error

	validateAsMessageField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type NewMap interface {
	MaybeLabel() *string
	SetLabel(string) error

	MaybeNumber() *uint
	SetNumber(uint) error

	MaybeKeyType() KeyType
	SetKeyType(KeyType) error

	MaybeValueType() FieldType
	SetValueType(FieldType) error

	Deprecated() bool
	SetDeprecated(bool) error

	InsertIntoParent(index uint) error
}

type Map interface {
	Label() string
	SetLabel(string) error

	Number() uint
	SetNumber(uint) error

	KeyType() keyType
	SetKeyType(keyType) error

	ValueType() FieldType
	SetValueType(FieldType) error

	Deprecated() bool
	SetDeprecated(bool) error

	validateAsMessageField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type NewOneOf interface {
	MaybeLabel() *string
	SetLabel(string) error

	NumFields() uint
	Field(uint) OneOfField

	NewField() NewOneOfField

	validateLabel(Label) error
	validateNumber(Number) error
	insertField(index uint, field OneOfField) error

	InsertIntoParent(index uint) error
}

type OneOf interface {
	Label() string
	SetLabel(string) error

	NewField() NewOneOfField

	validateLabel(Label) error
	validateNumber(Number) error
	insertField(index uint, field OneOfField) error

	validateAsMessageField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type NewOneOfField interface {
	MaybeLabel() *string
	SetLabel(string) error

	MaybeNumber() *uint
	SetNumber(uint) error

	MaybeType() FieldType
	SetType(FieldType) error

	Deprecated() bool
	SetDeprecated(bool) error

	InsertIntoParent(index uint) error
}

type OneOfField interface {
	Label() string
	SetLabel(string) error

	Number() uint
	SetNumber(uint) error

	Type() FieldType
	SetType(FieldType) error

	Deprecated() bool
	SetDeprecated(bool) error

	validateAsOneOfField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type NewReservedNumbers interface {
	NumNumbers() uint
	Number(index uint) FieldNumber

	InsertNumber(index, number uint) error
	NewNumberRange() NewNumberRange

	InsertIntoParent(index uint) error

	validateNumber(FieldNumber) error
	insertNumber(index uint, number FieldNumber) error
}

type ReservedNumbers interface {
	NumNumbers() uint
	Number(index uint) FieldNumber

	InsertNumber(index, number uint) error
	NewNumberRange() NewNumberRange

	validateNumber(FieldNumber) error
	insertNumber(index uint, number FieldNumber) error

	validateAsMessageField() error
	validateAsEnumField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type NewReservedLabels interface {
	NumLabels() uint
	Label(index uint) Label

	InsertLabel(index uint, label string) error

	validateLabel(Label) error

	InsertIntoParent(index uint) error
}

type ReservedLabels interface {
	NumLabels() uint
	Label(index uint) Label

	InsertLabel(index uint, label string) error

	validateLabel(Label) error

	validateAsMessageField() error
	validateAsEnumField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type MessageField interface {
	validateAsMessageField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type NewEnum interface {
	MaybeLabel() *string
	SetLabel(string) error

	AllowAlias() bool
	SetAllowAlias(bool) error

	NumFields() uint
	Field(uint) EnumField

	NewVariant() NewVariant
	NewReservedNumbers() NewReservedNumbers
	NewReservedLabels() NewReservedLabels

	InsertIntoParent(index uint) error

	validateLabel(Label) error
	validateNumber(FieldNumber) error
	insertField(uint, EnumField) error
}

type Enum interface {
	Label() string
	SetLabel(string) error

	AllowAlias() bool
	SetAllowAlias(bool) error

	NumFields() uint
	Field(uint) EnumField

	NewVariant() NewVariant
	NewReservedNumbers() NewReservedNumbers
	NewReservedLabels() NewReservedLabels

	validateLabel(Label) error
	validateNumber(FieldNumber) error
	insertField(uint, EnumField) error

	validateAsDefinition() error
	hasLabel(Label) bool

	FieldType
}

type EnumField interface {
	validateAsEnumField() error
	hasLabel(Label) bool
	hasNumber(FieldNumber) bool
}

type NewVariant interface {
	MaybeLabel() *string
	SetLabel(string) error

	MaybeNumber() *uint
	SetNumber(uint) error

	Deprecated() bool
	SetDeprecated(bool) error

	InsertIntoParent(index uint) error
}

type Variant interface {
	Label() string
	SetLabel(string) error

	Number() uint
	SetNumber(uint) error

	Deprecated() bool
	SetDeprecated(bool) error

	validateAsEnumField() error
	hasNumber(FieldNumber) bool
	hasLabel(Label) bool
}

type Label interface {
	Value() string
	SetValue(string) error
}

type Number interface {
	FieldNumber

	Value() uint
	SetValue(uint) error
}

type NewNumberRange interface {
	MaybeStart() *uint
	SetStart(uint) error
	MaybeEnd() *uint
	SetEnd(uint) error

	InsertIntoParent(index uint) error
}

type NumberRange interface {
	FieldNumber

	Start() uint
	SetStart(uint) error
	End() uint
	SetEnd(uint) error
}

type FieldNumber interface {
	intersects(FieldNumber) bool
}

type KeyType interface {
	_isKeyType()
}

type FieldType interface {
	_isFieldType()
}
