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
	// TODO: probably this should take a `declaration` so we can check if the
	// carrier of the identifier is already in the container. then it does not
	// need to fail a duplication check against itself
	validateLabel(identifier) error
}

type definitionContainer interface {
	declarationContainer

	NumDefinitions() uint
	Definition(uint) Definition
	NewMessage() Message
	NewEnum() Enum

	insertDefinition(index uint, def Definition)
}

type Definition interface {
	declaration

	InsertIntoParent(uint) error
	NumFields() uint
	NewReservedNumbers() *reservedNumbers
	NewReservedLabels() *reservedLabels

	validateNumber(fieldNumber) error
	validateAsDefinition() error
}

type Message interface {
	Definition
	definitionContainer

	Field(uint) messageField
	NewField() *repeatableField
	NewMap() *mapField
	NewOneOf() *oneOf

	insertField(uint, messageField)
}

type messageField interface {
	InsertIntoParent(uint) error

	validateAsMessageField() error
}

type Enum interface {
	Definition

	GetAlias() bool
	SetAlias(bool) error
	Field(uint) enumField

	NewField() *enumeration

	insertField(uint, enumField)
}

type enumField interface {
	InsertIntoParent(uint) error
	validateAsEnumField() error
}

type fieldNumber interface {
	intersects(fieldNumber) bool
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
}

type fieldType interface {
	_fieldType()
}
