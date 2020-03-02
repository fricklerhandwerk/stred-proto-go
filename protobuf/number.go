package protobuf

import (
	"errors"
	"fmt"
)

// TODO: maybe make this a struct with a parent, but this will require quite
// a lot of refactoring... not sure yet how exactly to proceed. having a common
// public interface for single numbers and ranges would be very good, because
// then the client will not need to define interfaces just to distinguish them.
// probably this will go in the direction of what is proposed in [1], such that
// structs will be returned by value, implementing exported interfaces, and
// collections will return their items through an index getter. this requires
// a lot of boilerplate, but is also the safest construction. there may still
// be surprises with nested pointers though...
// [1]: https://levelup.gitconnected.com/building-immutable-data-structures-in-go-56a1068c76b2
type number uint

func (n number) intersects(other fieldNumber) bool {
	switch o := other.(type) {
	case number:
		if n == other {
			return true
		}
	case *numberRange:
		if n >= o.start && n <= o.end {
			return true
		}
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
	return false
}

func (n number) GetValue() uint {
	return uint(n)
}

type numberRange struct {
	start  number
	end    number
	parent *reservedNumbers
}

func (r numberRange) GetStart() uint {
	return uint(r.start)
}

func (r *numberRange) SetStart(s uint) error {
	if s >= uint(r.end) {
		return errors.New("end of number range must be greater than start")
	}
	old := r.start
	r.start = number(s)
	if err := r.parent.validateNumber(r); err != nil {
		r.start = old
		return err
	}
	return nil
}

func (r numberRange) GetEnd() uint {
	return uint(r.end)
}

func (r *numberRange) SetEnd(e uint) error {
	if uint(r.start) >= e {
		return errors.New("end of number range must be greater than start")
	}
	old := r.end
	r.end = number(e)
	if err := r.parent.validateNumber(r); err != nil {
		r.end = old
		return err
	}
	return nil
}

func (r *numberRange) intersects(other fieldNumber) bool {
	switch o := other.(type) {
	case number:
		return o.intersects(r)
	case *numberRange:
		return o.start.intersects(r) || o.end.intersects(r)
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}
