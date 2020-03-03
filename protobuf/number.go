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
type number struct {
	value     uint
	container numberContainer
	parent    interface{}
}

func (n number) intersects(other fieldNumber) bool {
	switch o := other.(type) {
	case *number:
		if n.value == o.value {
			return true
		}
	case *numberRange:
		if n.value >= o.start.value && n.value <= o.end.value {
			return true
		}
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
	return false
}

func (n number) getParent() interface{} {
	return n.parent
}

func (n number) GetValue() uint {
	return n.value
}

func (n *number) SetValue(v uint) error {
	old := n.value
	n.value = v
	if err := n.container.validateNumber(n); err != nil {
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

func (r numberRange) getParent() interface{} {
	return r.parent
}

func (r numberRange) GetStart() uint {
	return r.start.value
}

func (r *numberRange) SetStart(s uint) error {
	if s >= r.end.value {
		return errors.New("end of number range must be greater than start")
	}
	old := r.start.value
	r.start.value = s
	if err := r.parent.validateNumber(r); err != nil {
		r.start.value = old
		return err
	}
	return nil
}

func (r numberRange) GetEnd() uint {
	return r.end.value
}

func (r *numberRange) SetEnd(e uint) error {
	if r.start.value >= e {
		return errors.New("end of number range must be greater than start")
	}
	old := r.end.value
	r.end.value = e
	if err := r.parent.validateNumber(r); err != nil {
		r.end.value = old
		return err
	}
	return nil
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
