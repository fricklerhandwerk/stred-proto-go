package protobuf

import (
	"errors"
	"fmt"
)

type number struct {
	value  uint
	parent Numbered
}

func (n number) intersects(other FieldNumber) bool {
	switch o := other.(type) {
	case *number:
		return n.value == o.value
	case *numberRange:
		return n.value >= o.start.value && n.value <= o.end.value
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}

func (n number) Value() uint {
	return n.value
}

func (n number) Number() uint {
	return n.value
}

func (n *number) maybeValue() *uint {
	if n == nil {
		return nil
	}
	value := n.value
	return &value
}

func (n *number) maybeNumber() *uint {
	return n.maybeValue()
}

func (n *number) SetValue(v uint) error {
	old := n.value
	n.value = v
	if err := n.parent.validateNumber(n); err != nil {
		n.value = old
		return err
	}
	return nil
}

func (n *number) SetNumber(v uint) error {
	return n.SetValue(v)
}

type newNumberRange struct {
	numberRange *numberRange
}

func (r newNumberRange) MaybeStart() *uint {
	return r.numberRange.start.maybeValue()
}

func (r *newNumberRange) SetStart(s uint) (err error) {
	if r.numberRange.start == nil {
		r.numberRange.start = &number{
			parent: r.numberRange.parent,
		}
		defer func() {
			if err != nil {
				r.numberRange.start = nil
			}
		}()
	}

	if r.numberRange.end != nil {
		return r.numberRange.SetStart(s)
	}
	return r.numberRange.start.SetValue(s)
}

func (r newNumberRange) MaybeEnd() *uint {
	return r.numberRange.end.maybeValue()
}

func (r *newNumberRange) SetEnd(s uint) (err error) {
	if r.numberRange.end == nil {
		r.numberRange.end = &number{
			parent: r.numberRange.parent,
		}
		defer func() {
			if err != nil {
				r.numberRange.end = nil
			}
		}()
	}

	if r.numberRange.start != nil {
		return r.numberRange.SetEnd(s)
	}
	return r.numberRange.end.SetValue(s)
}

func (r *newNumberRange) InsertIntoParent(i uint) error {
	if err := r.numberRange.parent.insertNumber(i, r.numberRange); err != nil {
		return err
	}
	r.numberRange = &numberRange{
		parent: r.numberRange.parent,
	}
	return nil
}

type numberRange struct {
	start  *number
	end    *number
	parent *reservedNumbers
}

func (r numberRange) Start() uint {
	return r.start.value
}

func (r *numberRange) SetStart(s uint) (err error) {
	old := r.start.value
	r.start.value = s

	defer func() {
		if err != nil {
			r.start.value = old
		}
	}()

	if s >= r.end.value {
		return errors.New("end of number range must be greater than start")
	}
	return r.parent.validateNumber(r)
}

func (r numberRange) End() uint {
	return r.end.value
}

func (r *numberRange) SetEnd(e uint) (err error) {
	old := r.end.value
	r.end.value = e

	defer func() {
		if err != nil {
			r.end.value = old
		}
	}()

	if r.start.value >= e {
		return errors.New("end of number range must be greater than start")
	}
	return r.parent.validateNumber(r)
}

func (r *numberRange) intersects(other FieldNumber) bool {
	switch o := other.(type) {
	case *number:
		return o.intersects(r)
	case *numberRange:
		return o.start.intersects(r) || o.end.intersects(r)
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}
