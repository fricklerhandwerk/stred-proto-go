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

// TODO: maybe this assumes to much about how the UI will behave, and we should instead have an interface consistent with the other containers: r.NewNumber(), r.NewNumberRange(), etc.
func (r *reservedNumbers) InsertNumber(i uint, n uint) error {
	// check self-consistency in case range was not yet added to parent
	if err := r.validateNumber(number(n)); err != nil {
		return err
	}
	r.numbers = append(r.numbers, nil)
	copy(r.numbers[i+1:], r.numbers[i:])
	r.numbers[i] = number(n)
	return nil
}

func (r reservedNumbers) validateNumber(n fieldNumber) error {
	for _, i := range r.numbers {
		if i.intersects(n) {
			var source string
			switch v := i.(type) {
			case number:
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

// TODO: this is a bad interface. the number in the list should instead be
// converted directly, otherwise we have no guarantee that the underlying value
// is actually of type `number`. but then we have to rework that type...
func (r *reservedNumbers) ToRange(i uint, end uint) error {
	var (
		start number
		ok    bool
	)
	if start, ok = r.numbers[i].(number); !ok {
		// TODO: probably panic here, should not happen
		return errors.New("reserved field must be a single number")
	}
	nr := &numberRange{
		parent: r,
		start:  start,
	}
	var empty *numberRange
	r.numbers[i] = empty // otherwise old value will be part of the check
	if err := nr.SetEnd(end); err != nil {
		r.numbers[i] = start
		return err
	}
	r.numbers[i] = nr
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
