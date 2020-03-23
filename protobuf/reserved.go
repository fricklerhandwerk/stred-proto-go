package protobuf

import (
	"errors"
	"fmt"
)

type newReservedNumbers struct {
	reservedNumbers *reservedNumbers
}

func (r *newReservedNumbers) InsertIntoParent(i uint) error {
	if len(r.reservedNumbers.numbers) < 1 {
		return errors.New("reserved numbers need at least one entry")
	}
	switch p := r.reservedNumbers.parent.(type) {
	case *enum:
		return p.insertField(i, r.reservedNumbers)
	case *message:
		return p.insertField(i, r.reservedNumbers)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
}

func (r *newReservedNumbers) NumNumbers() uint {
	return r.reservedNumbers.NumNumbers()
}

func (r *newReservedNumbers) Number(i uint) FieldNumber {
	return r.reservedNumbers.Number(i)
}

func (r *newReservedNumbers) InsertNumber(i, n uint) error {
	return r.reservedNumbers.InsertNumber(i, n)
}

func (r *newReservedNumbers) NewNumberRange() NewNumberRange {
	return r.reservedNumbers.NewNumberRange()
}

type reservedNumbers struct {
	numbers []FieldNumber
	parent  Definition
}

func (r reservedNumbers) NumNumbers() uint {
	return uint(len(r.numbers))
}

func (r reservedNumbers) Number(i uint) FieldNumber {
	return r.numbers[i]
}

func (r *reservedNumbers) NewNumberRange() NewNumberRange {
	return &newNumberRange{
		numberRange: &numberRange{parent: r},
	}
}

func (r *reservedNumbers) InsertNumber(i uint, n uint) error {
	num := &number{
		parent: r,
		value:  n,
	}
	return r.insertNumber(i, num)
}

func (r *reservedNumbers) insertNumber(i uint, n FieldNumber) error {
	if err := r.validateNumber(n); err != nil {
		return err
	}
	r.numbers = append(r.numbers, nil)
	copy(r.numbers[i+1:], r.numbers[i:])
	r.numbers[i] = n
	return nil
}

func (r *reservedNumbers) validateNumber(n FieldNumber) error {
	switch v := n.(type) {
	case *numberRange:
		if v.start == nil {
			return fmt.Errorf("number range has no start")
		}
		if v.end == nil {
			return fmt.Errorf("number range has no end")
		}
	}
	for _, i := range r.numbers {
		if i != n && i.intersects(n) {
			var source string
			switch v := i.(type) {
			case *number:
				source = fmt.Sprintf("field number %d", v)
			case *numberRange:
				source = fmt.Sprintf("range %d to %d", v.Start(), v.End())
			default:
				panic(fmt.Sprintf("unhandled number type %T", i))
			}
			return fmt.Errorf("%s already reserved", source)
		}
	}
	if err := r.parent.validateNumber(n); err != nil {
		return err
	}
	return nil
}

func (r reservedNumbers) validateAsEnumField() error {
	panic("not implemented")
}

func (r reservedNumbers) validateAsMessageField() error {
	panic("not implemented")
}

func (r reservedNumbers) hasLabel(l *label) bool {
	return false
}

func (r reservedNumbers) hasNumber(n FieldNumber) bool {
	for _, m := range r.numbers {
		if m != n && n.intersects(m) {
			return true
		}
	}
	return false
}

type newReservedLabels struct {
	reservedLabels *reservedLabels
}

func (r *newReservedLabels) InsertIntoParent(i uint) error {
	switch p := r.reservedLabels.parent.(type) {
	case *enum:
		return p.insertField(i, r.reservedLabels)
	case *message:
		return p.insertField(i, r.reservedLabels)
	default:
		panic(fmt.Sprintf("unhandled reservation parent type %T", p))
	}
}

func (r *newReservedLabels) InsertLabel(i uint, l string) error {
	return r.reservedLabels.InsertLabel(i, l)
}

func (r *newReservedLabels) NumLabels() uint {
	return r.reservedLabels.NumLabels()
}
func (r *newReservedLabels) Label(i uint) Label {
	return r.reservedLabels.Label(i)
}

type reservedLabels struct {
	labels []*label
	parent declarationContainer
}

func (r reservedLabels) NumLabels() uint {
	return uint(len(r.labels))
}

func (r reservedLabels) Label(i uint) Label {
	return r.labels[i]
}

func (r *reservedLabels) InsertLabel(i uint, s string) error {
	if err := validateIdentifier(s); err != nil {
		return err
	}
	l := &label{
		parent: r,
		value:  s,
	}
	if err := r.validateLabel(l); err != nil {
		return err
	}
	r.labels = append(r.labels, nil)
	copy(r.labels[i+1:], r.labels[i:])
	r.labels[i] = l
	return nil
}

func (r reservedLabels) validateLabel(l *label) error {
	if r.hasLabel(l) {
		return fmt.Errorf("label %q already declared", l.value)
	}
	return r.parent.validateLabel(l)
}

func (r reservedLabels) validateAsEnumField() error {
	if len(r.labels) < 1 {
		return errors.New("reserved labels need at least one entry")
	}
	for _, l := range r.labels {
		if err := l.parent.validateLabel(l); err != nil {
			return err
		}
	}
	return nil
}

func (r reservedLabels) validateAsMessageField() error {
	return r.validateAsEnumField()
}

func (r reservedLabels) hasLabel(l *label) bool {
	for _, s := range r.labels {
		if s.hasLabel(l) {
			return true
		}
	}
	return false
}

func (r reservedLabels) hasNumber(n FieldNumber) bool {
	return false
}
