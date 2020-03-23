package protobuf

import (
	"errors"
	"fmt"
	"strings"
)

type newEnum struct {
	enum *enum
}

func (e *newEnum) InsertIntoParent(i uint) error {
	if err := e.enum.parent.insertDefinition(i, e.enum); err != nil {
		return err
	}
	e.enum = &enum{
		parent: e.enum.parent,
	}
	return nil
}

func (e *newEnum) MaybeLabel() *string {
	return e.enum.maybeLabel()
}

func (e *newEnum) SetLabel(l string) (err error) {
	if e.enum.label == nil {
		e.enum.label = &label{
			parent: e.enum.parent,
		}
		defer func() {
			if err != nil {
				e.enum.label = nil
			}
		}()
	}
	return e.enum.SetLabel(l)
}

func (e *newEnum) AllowAlias() bool {
	return e.enum.AllowAlias()
}

func (e *newEnum) SetAllowAlias(b bool) error {
	return e.enum.SetAllowAlias(b)
}

func (e *newEnum) NumFields() uint {
	return e.enum.NumFields()
}

func (e *newEnum) Field(i uint) EnumField {
	return e.enum.Field(i)
}

func (e *newEnum) NewVariant() NewVariant {
	return e.enum.NewVariant()
}

func (e *newEnum) NewReservedNumbers() NewReservedNumbers {
	return e.enum.NewReservedNumbers()
}

func (e *newEnum) NewReservedLabels() NewReservedLabels {
	return e.enum.NewReservedLabels()
}

type enum struct {
	*label
	allowAlias bool
	fields     []EnumField
	parent     definitionContainer
}

func (e enum) AllowAlias() bool {
	return e.allowAlias
}

func (e *enum) SetAllowAlias(b bool) error {
	if !b && e.allowAlias {
		// check if aliasing is in place
		numbers := make(map[uint]bool, len(e.fields))
		for _, field := range e.fields {
			switch f := field.(type) {
			case *variant:
				n := f.Number()
				if numbers[n] {
					lines := []string{
						fmt.Sprintf("field number %d is used multiple times.", n),
						fmt.Sprintf("remove aliasing before setting %q.", "allow_alias = false"),
					}
					return errors.New(strings.Join(lines, " "))
				}
				numbers[n] = true
			case FieldNumber:
				continue
			default:
				panic(fmt.Sprintf("unhandled enum field type %T", f))
			}
		}
	}
	e.allowAlias = b
	return nil
}

func (e enum) NumFields() uint {
	return uint(len(e.fields))
}

func (e enum) Field(i uint) EnumField {
	return e.fields[i]
}

func (e *enum) insertField(i uint, field EnumField) error {
	if err := field.validateAsEnumField(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	// <https://github.com/golang/go/wiki/SliceTricks#insert>
	// <https://stackoverflow.com/a/46130603/5147619>
	e.fields = append(e.fields, nil)
	copy(e.fields[i+1:], e.fields[i:])
	e.fields[i] = field
	return nil
}

func (e *enum) NewVariant() NewVariant {
	return &newVariant{
		variant: &variant{parent: e},
	}
}

func (e *enum) NewReservedNumbers() NewReservedNumbers {
	return &newReservedNumbers{
		reservedNumbers: &reservedNumbers{parent: e},
	}
}

func (e *enum) NewReservedLabels() NewReservedLabels {
	return &newReservedLabels{
		reservedLabels: &reservedLabels{parent: e},
	}
}

func (e enum) validateLabel(l *label) error {
	if l == nil {
		return fmt.Errorf("label not set")
	}
	for _, f := range e.fields {
		if f.hasLabel(l) {
			return fmt.Errorf("label %s already declared", l.value)
		}
	}
	return nil
}

func (e enum) validateNumber(n FieldNumber) error {
	// TODO: check valid values
	// https://developers.google.com/protocol-buffers/docs/proto3#assigning-field-numbers
	for _, f := range e.fields {
		if f.hasNumber(n) {
			switch f.(type) {
			case *variant:
				switch n := n.(type) {
				case *number:
					switch n.parent.(type) {
					case *variant:
						if e.allowAlias {
							return nil
						}
						lines := []string{
							fmt.Sprintf("field number %d already in use.", n.value),
							fmt.Sprintf("set %q to allow multiple labels for one number.", "allow_alias = true"),
						}
						return errors.New(strings.Join(lines, " "))
					}
				}
			}
			var source string
			switch v := n.(type) {
			case *number:
				source = fmt.Sprintf("field number %d", v.Value())
			case *numberRange:
				source = fmt.Sprintf("range %d to %d", v.Start(), v.End())
			default:
				panic(fmt.Sprintf("unhandled number type %T", n))
			}
			return fmt.Errorf("%s already in use", source)
		}
	}
	return nil
}

func (e *enum) validateAsDefinition() (err error) {
	return e.parent.validateLabel(e.label)
}

func (e *enum) _isFieldType() {}

type newVariant struct {
	variant *variant
}

func (v *newVariant) InsertIntoParent(i uint) error {
	if err := v.variant.parent.insertField(i, v.variant); err != nil {
		return err
	}
	v.variant = &variant{
		parent: v.variant.parent,
	}
	return nil
}

func (v newVariant) MaybeLabel() *string {
	return v.variant.maybeLabel()
}

func (v newVariant) SetLabel(l string) (err error) {
	if v.variant.label == nil {
		v.variant.label = &label{
			parent: v.variant.parent,
		}
		defer func() {
			if err != nil {
				v.variant.label = nil
			}
		}()
	}
	return v.variant.SetLabel(l)
}

func (v newVariant) MaybeNumber() *uint {
	return v.variant.maybeNumber()
}

func (v newVariant) SetNumber(n uint) (err error) {
	if v.variant.number == nil {
		v.variant.number = &number{
			// XXX: here the number's parent is the variant itself so we
			// can check aliasing
			parent: v.variant,
		}
		defer func() {
			if err != nil {
				v.variant.number = nil
			}
		}()
	}
	return v.variant.SetNumber(n)
}

func (v newVariant) Deprecated() bool {
	return v.variant.Deprecated()
}

func (v newVariant) SetDeprecated(b bool) error {
	return v.variant.SetDeprecated(b)
}

type variant struct {
	field
	parent *enum
}

func (v *variant) validateAsEnumField() (err error) {
	err = v.parent.validateLabel(v.label)
	if err != nil {
		return err

	}

	if v.number == nil {
		return fmt.Errorf("field number not set")
	}
	err = v.parent.validateNumber(v.number)
	if err != nil {
		return err
	}
	return nil
}

// XXX: this is a hack to be able to trace the parent of
// a number, so we can validate enum aliasing
func (v variant) validateNumber(n FieldNumber) error {
	return v.parent.validateNumber(n)
}
