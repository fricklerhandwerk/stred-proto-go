package protobuf

import "fmt"

type reservedNumbers struct {
	numbers []fieldNumber
	parent  definition
}

type fieldNumber interface {
	intersects([]fieldNumber) bool
}

func (r reservedNumbers) InsertNumber(index uint, n fieldNumber) error {
	panic("not implemented")
}

func (r *reservedNumbers) InsertIntoParent(i uint) error {
	switch p := r.parent.(type) {
	case *enum:
		if err := r.validateAsEnumField(); err != nil {
			return err
		}
		p.insertField(i, r)
	case *message:
		if err := r.validateAsMessageField(); err != nil {
			return err
		}
		p.insertField(i, r)
	default:
		panic(fmt.Sprintf("unhandled parent type %T", p))
	}
	return nil
}

func (e reservedNumbers) validateAsEnumField() error {
	panic("not implemented")
}

func (e reservedNumbers) validateAsMessageField() error {
	panic("not implemented")
}

type reservedLabels struct {
	labels []identifier
	parent definition
}

func (r reservedLabels) GetLabels() []string {
	panic("not implemented")
}

func (r reservedLabels) InsertLabel(index uint, n string) error {
	panic("not implemented")
}

func (r *reservedLabels) InsertIntoParent(i uint) error {
	switch p := r.parent.(type) {
	case *enum:
		if err := r.validateAsEnumField(); err != nil {
			return err
		}
		p.insertField(i, r)
	case *message:
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
