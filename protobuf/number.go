package protobuf

import (
	"errors"
	"fmt"
)

type number struct {
	value  uint
	parent numberContainer
}

func (n number) intersects(other fieldNumber) bool {
	switch o := other.(type) {
	case *number:
		return n.value == o.value
	case *numberRange:
		return n.value >= o.start.value && n.value <= o.end.value
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}

// TODO: try `SetNumber()` so we can embed this in fields
func (n number) GetValue() uint {
	return n.value
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

type numberRange struct {
	start  *number
	end    *number
	parent *reservedNumbers
}

// TODO: maybe this should be a pointer for easier distinction in UI
func (r numberRange) GetStart() uint {
	return r.start.value
}

func (r *numberRange) SetStart(s uint) (err error) {
	if r.start != nil {
		old := r.start.value
		r.start.value = s

		defer func() {
			if err != nil {
				r.start.value = old
			}
		}()
	} else {
		r.start = &number{
			parent: r.parent,
			value:  s,
		}
		defer func() {
			if err != nil {
				r.start = nil
			}
		}()
	}

	if r.end != nil {
		if s >= r.end.value {
			return errors.New("end of number range must be greater than start")
		}
		return r.parent.validateNumber(r)
	}
	return r.parent.validateNumber(r.start)
}

func (r numberRange) GetEnd() uint {
	return r.end.value
}

func (r *numberRange) SetEnd(e uint) (err error) {
	if r.end != nil {
		old := r.end.value
		r.end.value = e

		defer func() {
			if err != nil {
				r.end.value = old
			}
		}()
	} else {
		r.end = &number{
			value:  e,
			parent: r.parent,
		}
		defer func() {
			if err != nil {
				r.end = nil
			}
		}()
	}

	if r.start != nil {
		if r.start.value >= e {
			return errors.New("end of number range must be greater than start")
		}
		return r.parent.validateNumber(r)
	}

	return r.parent.validateNumber(r.end)
}

func (r *numberRange) InsertIntoParent(i uint) error {
	return r.parent.insertNumber(i, r)
}

func (r *numberRange) intersects(other fieldNumber) bool {
	switch o := other.(type) {
	case *number:
		return o.intersects(r)
	case *numberRange:
		return o.start.intersects(r) || o.end.intersects(r)
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}
