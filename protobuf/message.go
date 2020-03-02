package protobuf

import (
	"errors"
	"fmt"
)

type message struct {
	label
	fields      []messageField
	definitions []Definition
	parent      definitionContainer
}

func (m message) NumFields() uint {
	return uint(len(m.fields))
}

func (m message) Field(i uint) messageField {
	return m.fields[i]
}

func (m *message) insertField(i uint, field messageField) {
	m.fields = append(m.fields, nil)
	copy(m.fields[i+1:], m.fields[i:])
	m.fields[i] = field
}

func (m *message) newTypedField() typedField {
	return typedField{
		field: field{
			parent: m,
			label: label{
				parent: m,
			},
		},
	}
}

func (m *message) NewField() *repeatableField {
	return &repeatableField{
		parent:     m,
		typedField: m.newTypedField(),
	}
}

func (m *message) NewMap() *mapField {
	return &mapField{
		parent:     m,
		typedField: m.newTypedField(),
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

func (m *message) NewReservedNumbers() *reservedNumbers {
	return &reservedNumbers{
		parent: m,
	}
}

func (m *message) NewReservedLabels() *reservedLabels {
	return &reservedLabels{
		parent: m,
	}
}

func (m *message) NewMessage() Message {
	return &message{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

func (m *message) NewEnum() Enum {
	return &enum{
		parent: m,
		label: label{
			parent: m,
		},
	}
}

// TODO: use a common implementation for definition containers
func (m message) NumDefinitions() uint {
	return uint(len(m.definitions))
}

func (m message) Definition(i uint) Definition {
	return m.definitions[i]
}

func (m *message) insertDefinition(i uint, d Definition) {
	m.definitions = append(m.definitions, nil)
	copy(m.definitions[i+1:], m.definitions[i:])
	m.definitions[i] = d
}

func (m *message) InsertIntoParent(i uint) (err error) {
	if err = m.validateAsDefinition(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return
	}
	m.parent.insertDefinition(i, m)
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
		if f != nil && f.hasLabel(l.String()) {
			return fmt.Errorf("label %q already declared", l.String())
		}
	}
	// TODO: definitions and fields share a namespace
	return nil
}

func (m message) validateNumber(n fieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	switch v := n.(type) {
	case Number:
		if v.GetValue() < 1 {
			return errors.New("message field number must be >= 1")
		}
	case NumberRange:
		if v.GetStart() < 1 {
			return errors.New("message field numbers must be >= 1")
		}
	default:
		panic(fmt.Sprintf("unhandled field number type %T", v))
	}
	for _, f := range m.fields {
		if f != nil && f.hasNumber(n) {
			return fmt.Errorf("field number %s already in use", n)
		}
	}
	return nil
}

func (m *message) validateAsDefinition() (err error) {
	if err = m.parent.validateLabel(m.label.label); err != nil {
		return
	}
	for i, f := range m.fields {
		// XXX: this only works because `f.validateAsMessageField()` calls
		// `f.parent.validateLabel()`, which in turn iterates over `m.fields`
		// under the assumption that `f` is not contained.
		// setting `m.fields[i] = nil` and *importantly* checking for `f != nil` in
		// `f.parent.validateLabel()` satisfies that assumption.
		// so we have O(n^2) runtime and a really confusing interdependency just to
		// check uniqueness of labels and numbers...
		// on the other hand with this approach we can check nested structures
		// conveniently: the parent does not need to know implementation details of
		// the children, except that the validation requires them to be out of the
		// collection.
		// unfortunately that way we delocalise the logic of excluding the current
		// field from the set of fields to validate against. it must be done
		// somewhere, after all, but it would be more intuitive and direct to
		// exclude the field itself from comparison *during* that comparison. note
		// that before insertion we actually want to validate every property
		// separately, so it is not a meaningful option to handle everything in the
		// field container. in the current setup the interfaces to handle
		// properties is somewhat generic, and due to reused embedded structs we
		// have it such that objects, which can carry one or more identifiers,
		// would need to inform that identifier that they carry it (by leaving
		// a member pointer).  but since just about anything has an identifier, the
		// type of identified object must essentially be `interface{}`. there is
		// nothing we need from these objects except their pointer, and there are
		// no useful methods on a possible sum type, so either the interface is
		// literally empty or we create a dummy sum type interface and implement it
		// on almost everything...
		// the other question would be how to pass the identifier's enclosing
		// object to the validator, and probably the safest way would be by setting
		// it already in the object's constructor. that is weird and cumbersome,
		// but centralised.
		// now all of this happens internally, so it does not matter a lot how ugly
		// and unsafe this is, but it needs to be somewhat easy to understand and
		// reason about. the other tradeoff is hard the semantics are encoded into
		// the types. here we weigh against each other a more type-safe variant,
		// where excluding the current field is done by convention, and a more
		// loosely typed one, where the exclusion is explicit and localised.
		m.fields[i] = nil
		defer func() { m.fields[i] = f }()
		if err = f.validateAsMessageField(); err != nil {
			return
		}
	}
	for i, d := range m.definitions {
		m.definitions[i] = nil
		defer func() { m.definitions[i] = d }()
		if err = d.validateAsDefinition(); err != nil {
			return
		}
	}
	return
}

func (m *message) _fieldType() {}
