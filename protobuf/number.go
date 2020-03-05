package protobuf

import (
	"errors"
	"fmt"
)

// TODO: if you compare `number` with `label`, you see that `number` here
// serves two purposes at once: it implements `fieldNumber` and is also
// a numeric identifier base type, this is why we have an overlap between
// `parent` and `container`. to have it consistent, it should look something
// like this:
//   type number struct {
//     integer integer
//     parent numberContainer
//   }
//   type integer {
//     value uint
//     parent interface{}
//   }

type integer struct {
	value  uint
	parent interface{}
}

type number struct {
	integer
	parent numberContainer
}

func (i integer) intersects(other fieldNumber) bool {
	switch o := other.(type) {
	case *number:
		return i.value == o.value
	case *numberRange:
		return i.value >= o.start.value && i.value <= o.end.value
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}

func (n number) getParent() interface{} {
	return n.integer.parent
}

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
	start  *integer
	end    *integer
	parent *reservedNumbers
}

func (r numberRange) getParent() interface{} {
	return r.parent
}

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
		r.start = &integer{
			parent: r,
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
	n := &number{
		parent:  r.parent,
		integer: *r.start,
	}
	return r.parent.validateNumber(n)
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
		r.end = &integer{
			value:  e,
			parent: r,
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

	n := &number{
		parent:  r.parent,
		integer: *r.end,
	}
	return r.parent.validateNumber(n)
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
