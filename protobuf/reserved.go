package protobuf

type reservedNumbers struct {
	numbers []fieldNumber
}

type fieldNumber interface {
	intersects([]fieldNumber) bool
}

func (r reservedNumbers) Insert(index uint, n fieldNumber) error {
	panic("not implemented")
}

func (e reservedNumbers) validateAsEnumField() error {
	panic("not implemented")
}

func (e reservedNumbers) validateAsMessageField() error {
	panic("not implemented")
}

type reservedLabels struct {
	labels []identifier
}

func (r reservedLabels) Get() []string {
	panic("not implemented")
}

func (r reservedLabels) Insert(index uint, n string) error {
	panic("not implemented")
}

func (r reservedLabels) validateAsEnumField() error {
	panic("not implemented")
}

func (r reservedLabels) validateAsMessageField() error {
	panic("not implemented")
}
