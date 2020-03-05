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
		value:     n,
		container: r.parent,
		parent:    r,
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
	for _, i := range r.numbers {
		if i != n && i.intersects(n) {
			var source string
			switch v := i.(type) {
			case Number:
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
	case Enum:
		if err := r.validateAsEnumField(); err != nil {
			return err
		}
		p.insertField(i, r)
	case Message:
		if err := r.validateAsMessageField(); err != nil {
			return err
		}
		p.insertField(i, r)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
	return nil
}

func (r reservedNumbers) validateAsEnumField() error {
	panic("not implemented")
}

func (r reservedNumbers) validateAsMessageField() error {
	panic("not implemented")
}

func (r reservedNumbers) hasLabel(l string) bool {
	return false
}

func (r reservedNumbers) hasNumber(n fieldNumber) bool {
	for _, m := range r.numbers {
		if n.intersects(m) {
			return true
		}
	}
	return false
}

type reservedLabels struct {
	labels []identifier
	parent Definition
}

func (r reservedLabels) GetLabels() []string {
	panic("not implemented")
}

func (r reservedLabels) InsertLabel(index uint, n string) error {
	panic("not implemented")
}

func (r *reservedLabels) InsertIntoParent(i uint) error {
	switch p := r.parent.(type) {
	case Enum:
		if err := r.validateAsEnumField(); err != nil {
			return err
		}
		p.insertField(i, r)
	case Message:
		if err := r.validateAsMessageField(); err != nil {
			return err
		}
		p.insertField(i, r)
	default:
		panic(fmt.Sprintf("unhandled reservation parent type %T", p))
	}
	return nil
}

func (r reservedLabels) validateAsEnumField() error {
	panic("not implemented")
}

func (r reservedLabels) validateAsMessageField() error {
	panic("not implemented")
}

func (r reservedLabels) hasLabel(l string) bool {
	for _, s := range r.labels {
		if s.String() == l {
			return true
		}
	}
	return false
}

func (r reservedLabels) hasNumber(n fieldNumber) bool {
	return false
}
