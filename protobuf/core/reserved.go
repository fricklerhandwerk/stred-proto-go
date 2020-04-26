package core

import (
	"errors"
	"fmt"
)

type reservedNumber struct {
	number Number
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

func (r reservedNumber) Get() *uint {
	return r.number.Get()
}

func (r *reservedNumber) Set(value uint) error {
	return r.number.Set(value)
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

func (r reservedNumber) hasLabel(l *Label) bool {
	return false
}

func (r *reservedNumber) hasNumber(n FieldNumber) bool {
	return r.number.hasNumber(n)
}

func (r reservedNumber) intersects(other FieldNumber) bool {
	return r.number.intersects(other)
}

type ReservedRange struct {
	start  Number
	end    Number
	parent Definition
}

func (r *ReservedRange) Start() *Number {
	if r.start.parent == nil {
		r.start.parent = r
	}
	return &r.start
}

func (r *ReservedRange) End() *Number {
	if r.end.parent == nil {
		r.end.parent = r
	}
	return &r.end
}

func (r *ReservedRange) InsertIntoParent() error {
	switch p := r.parent.(type) {
	case *enum:
		return p.insertField(r)
	case *message:
		return p.insertField(r)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
}

func (r *ReservedRange) Parent() Definition {
	return r.parent
}

func (r *ReservedRange) hasNumber(other FieldNumber) bool {
	return r != other && r.intersects(other)
}

func (r *ReservedRange) intersects(other FieldNumber) bool {
	switch o := other.(type) {
	case *Number:
		return o.intersects(r)
	case *ReservedRange:
		return o.start.intersects(r) || o.end.intersects(r)
	default:
		panic(fmt.Sprintf("unhandled fieldNumber type %T", o))
	}
}

func (r ReservedRange) hasLabel(l *Label) bool {
	return false
}

func (r *ReservedRange) validateNumber(n FieldNumber) error {
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

func (r *ReservedRange) validateAsEnumField() error {
	if err := r.validateNumber(r); err != nil {
		return err
	}
	return r.parent.validateNumber(r)
}

func (r *ReservedRange) validateAsMessageField() error {
	return r.validateAsEnumField()
}

type reservedLabel struct {
	label  Label
	parent Definition
}

func (r reservedLabel) Get() string {
	return r.label.Get()
}

func (r *reservedLabel) Set(value string) error {
	return r.label.Set(value)
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

func (r reservedLabel) validateLabel(l *Label) error {
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

func (r *reservedLabel) hasLabel(l *Label) bool {
	return r.label.hasLabel(l)
}

func (r reservedLabel) hasNumber(n FieldNumber) bool {
	return false
}
