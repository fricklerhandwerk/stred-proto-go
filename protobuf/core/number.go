package core

import (
	"fmt"
)

type number struct {
	value  *uint
	parent Numbered
}

func (n number) Get() *uint {
	if n.value == nil {
		return nil
	}
	out := *n.value
	return &out
}

func (n *number) Set(value uint) (err error) {
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

func (n *number) validate() error {
	if n.value == nil {
		return fmt.Errorf("number not set")
	}
	return n.parent.validateNumber(n)
}

func (n number) Parent() Numbered {
	return n.parent
}

func (n *number) hasNumber(other FieldNumber) bool {
	return n != other && n.intersects(other)
}

func (n *number) intersects(other FieldNumber) bool {
	switch o := other.(type) {
	case *number:
		return *n.value == *o.value
	case *reservedRange:
		return *n.value >= *o.start.value && *n.value <= *o.end.value
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}
