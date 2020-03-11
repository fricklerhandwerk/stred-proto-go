package protobuf

import (
	"errors"
	"fmt"
)

type reservedNumbers struct {
	numbers []fieldNumber
	parent  Definition
}

func (r reservedNumbers) NumNumbers() uint {
	return uint(len(r.numbers))
}

func (r reservedNumbers) Number(i uint) fieldNumber {
	return r.numbers[i]
}

func (r *reservedNumbers) NewRange() NumberRange {
	return &numberRange{
		parent: r,
	}
}

func (r *reservedNumbers) InsertNumber(i uint, n uint) error {
	num := &number{
		parent: r,
		value:  n,
	}
	return r.insertNumber(i, num)
}

func (r *reservedNumbers) insertNumber(i uint, n fieldNumber) error {
	if err := r.validateNumber(n); err != nil {
		return err
	}
	r.numbers = append(r.numbers, nil)
	copy(r.numbers[i+1:], r.numbers[i:])
	r.numbers[i] = n
	return nil
}

func (r *reservedNumbers) validateNumber(n fieldNumber) error {
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
				source = fmt.Sprintf("range %d to %d", v.GetStart(), v.GetEnd())
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

func (r *reservedNumbers) InsertIntoParent(i uint) error {
	if len(r.numbers) < 1 {
		return errors.New("reserved numbers need at least one entry")
	}
	switch p := r.parent.(type) {
	case *enum:
		return p.insertField(i, r)
	case *message:
		return p.insertField(i, r)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
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

func (r reservedNumbers) hasNumber(n fieldNumber) bool {
	for _, m := range r.numbers {
		if m != n && n.intersects(m) {
			return true
		}
	}
	return false
}

type reservedLabels struct {
	labels []*label
	parent declarationContainer
}

func (r reservedLabels) NumLabels() uint {
	return uint(len(r.labels))
}

func (r reservedLabels) Label(i uint) declaration {
	return r.labels[i]
}

func (r *reservedLabels) InsertLabel(i uint, s string) error {
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
	if err := l.validate(); err != nil {
		return err
	}
	if r.hasLabel(l) {
		return fmt.Errorf("label %q already declared", l.value)
	}
	return r.parent.validateLabel(l)
}

func (r *reservedLabels) InsertIntoParent(i uint) error {
	switch p := r.parent.(type) {
	case *enum:
		return p.insertField(i, r)
	case *message:
		return p.insertField(i, r)
	default:
		panic(fmt.Sprintf("unhandled reservation parent type %T", p))
	}
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

func (r reservedLabels) hasNumber(n fieldNumber) bool {
	return false
}
