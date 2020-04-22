package core

import (
	"errors"
	"fmt"
)

type reservedNumber struct {
	number
	parent Definition
}

func (r *reservedNumber) InsertIntoParent() error {
	switch p := r.parent.(type) {
	case *enum:
		return p.insertField(r)
	case *message:
		return p.insertField(r)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
}

func (r reservedNumber) Parent() Definition {
	return r.parent
}

func (r *reservedNumber) validateNumber(n FieldNumber) error {
	return r.parent.validateNumber(n)
}

func (r *reservedNumber) validateAsEnumField() error {
	return r.number.validate()
}

func (r reservedNumber) validateAsMessageField() error {
	return r.validateAsEnumField()
}

func (r reservedNumber) hasLabel(l *label) bool {
	return false
}

type reservedRange struct {
	start  number
	end    number
	parent Definition
}

func (r *reservedRange) Start() Number {
	if r.start.parent == nil {
		r.start.parent = r
	}
	return &r.start
}

func (r *reservedRange) End() Number {
	if r.end.parent == nil {
		r.end.parent = r
	}
	return &r.end
}

func (r *reservedRange) InsertIntoParent() error {
	switch p := r.parent.(type) {
	case *enum:
		return p.insertField(r)
	case *message:
		return p.insertField(r)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
}

func (r *reservedRange) Parent() Definition {
	return r.parent
}

func (r *reservedRange) hasNumber(other FieldNumber) bool {
	return r != other && r.intersects(other)
}

func (r *reservedRange) intersects(other FieldNumber) bool {
	switch o := other.(type) {
	case *number:
		return o.intersects(r)
	case *reservedRange:
		return o.start.intersects(r) || o.end.intersects(r)
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}

func (r reservedRange) hasLabel(l *label) bool {
	return false
}

func (r *reservedRange) validateNumber(n FieldNumber) error {
	switch n {
	case &r.start:
		if r.end.value == nil {
			return r.parent.validateNumber(n)
		}
	case &r.end:
		if r.start.value == nil {
			return r.parent.validateNumber(n)
		}
	case r:
		if err := r.start.validate(); err != nil {
			return err
		}
		return r.end.validate()
	default:
		panic("number for validation must be start, end, or whole range")
	}
	switch {
	case *r.start.value >= *r.end.value:
		return errors.New("end of number range must be greater than start")
	default:
		return r.parent.validateNumber(r)
	}
}

func (r *reservedRange) validateAsEnumField() error {
	if err := r.validateNumber(r); err != nil {
		return err
	}
	return r.parent.validateNumber(r)
}

func (r *reservedRange) validateAsMessageField() error {
	return r.validateAsEnumField()
}

type reservedLabel struct {
	label
	parent Definition
}

func (r *reservedLabel) InsertIntoParent() error {
	switch p := r.parent.(type) {
	case *enum:
		return p.insertField(r)
	case *message:
		return p.insertField(r)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
}

func (r reservedLabel) validateLabel(l *label) error {
	return r.parent.validateLabel(l)
}

func (r reservedLabel) Parent() Definition {
	return r.parent
}

func (r *reservedLabel) validateAsEnumField() error {
	return r.label.validate()
}

func (r reservedLabel) validateAsMessageField() error {
	return r.validateAsEnumField()
}

func (r reservedLabel) hasNumber(n FieldNumber) bool {
	return false
}
