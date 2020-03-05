package protobuf

type Document interface {
	definitionContainer

	GetPackage() *string
	SetPackage(string) error
	NewService() Service
}

type Service interface {
	declaration
}

type declaration interface {
	GetLabel() string
	SetLabel(string) error
}

type declarationContainer interface {
	validateLabel(ident identifier) error
}

type definitionContainer interface {
	declarationContainer

	NumDefinitions() uint
	Definition(uint) Definition
	NewMessage() Message
	NewEnum() Enum

	insertDefinition(index uint, def Definition) error
}

type Definition interface {
	declaration
	numberContainer

	InsertIntoParent(uint) error
	NumFields() uint
	NewReservedNumbers() *reservedNumbers
	NewReservedLabels() *reservedLabels

	validateAsDefinition() error
}

type numberContainer interface {
	validateNumber(n fieldNumber) error
}

type Message interface {
	Definition
	definitionContainer
	fieldType

	Field(uint) messageField
	NewField() *repeatableField
	NewMap() *mapField
	NewOneOf() *oneOf

	insertField(uint, messageField) error
}

type messageField interface {
	InsertIntoParent(uint) error

	validateAsMessageField() error
	hasLabel(string) bool
	hasNumber(fieldNumber) bool
}

type Enum interface {
	Definition
	fieldType

	GetAlias() bool
	SetAlias(bool) error
	Field(uint) enumField

	NewField() *enumeration

	insertField(uint, enumField) error
}

type enumField interface {
	InsertIntoParent(uint) error
	validateAsEnumField() error
	hasLabel(string) bool
	hasNumber(fieldNumber) bool
}

type fieldNumber interface {
	intersects(fieldNumber) bool
	getParent() interface{}
}

type Number interface {
	fieldNumber

	GetValue() uint
	SetValue(uint) error
}

type NumberRange interface {
	fieldNumber

	GetStart() uint
	SetStart(uint) error
	GetEnd() uint
	SetEnd(uint) error
	InsertIntoParent(uint) error
}

type fieldType interface {
	_fieldType()
}
