package core

import (
	"fmt"
)

type Number struct {
	value  *uint
	parent Numbered
}

type Numbered interface {
	validateNumber(FieldNumber) error
}

func (n Number) Get() *uint {
	if n.value == nil {
		return nil
	}
	out := *n.value
	return &out
}

func (n *Number) Set(value uint) (err error) {
	if n.value == nil {
		n.value = &value
		defer func() {
			if err != nil {
				n.value = nil
			}
		}()
	} else {
		old := *n.value
		*n.value = value
		defer func() {
			if err != nil {
				*n.value = old
			}
		}()
	}
	return n.validate()
}

func (n *Number) validate() error {
	if n.value == nil {
		return fmt.Errorf("number not set")
	}
	return n.parent.validateNumber(n)
}

func (n Number) Parent() Numbered {
	return n.parent
}

func (n *Number) hasNumber(other FieldNumber) bool {
	return n != other && n.intersects(other)
}

func (n *Number) intersects(other FieldNumber) bool {
	switch o := other.(type) {
	case *Number:
		return *n.value == *o.value
	case *reservedRange:
		return *n.value >= *o.start.value && *n.value <= *o.end.value
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}
