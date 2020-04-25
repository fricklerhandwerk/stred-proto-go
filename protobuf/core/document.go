package core

import (
	"fmt"
)

func NewDocument() *Document {
	return &Document{}
}

type Document struct {
	_package Package
	imports  map[*Import]struct{}
	services map[*Service]struct{}
	messages map[*message]struct{}
	enums    map[*enum]struct{}
}

func (d *Document) Package() *Package {
	if d._package.parent == nil {
		d._package.parent = d
		d._package.label.parent = &d._package
	}
	return &d._package
}

func (d Document) Imports() []*Import {
	out := make([]*Import, len(d.imports))
	j := 0
	for i := range d.imports {
		out[j] = i
		j++
	}
	return out
}

func (d *Document) NewImport() *Import {
	return &Import{
		parent: d,
	}
}

func (d Document) Services() (out []*Service) {
	out = make([]*Service, len(d.services))
	i := 0
	for s := range d.services {
		out[i] = s
		i++
	}
	return
}

func (d *Document) NewService() *Service {
	return &Service{
		parent: d,
	}
}

func (d Document) Messages() (out []Message) {
	out = make([]Message, len(d.messages))
	i := 0
	for m := range d.messages {
		out[i] = m
		i++
	}
	return
}

func (d *Document) NewMessage() *NewMessage {
	return &NewMessage{parent: d}
}

func (d Document) Enums() (out []Enum) {
	out = make([]Enum, len(d.enums))
	i := 0
	for e := range d.enums {
		out[i] = e
		i++
	}
	return
}

func (d *Document) NewEnum() *NewEnum {
	return &NewEnum{parent: d}
}

func (d *Document) insertImport(i *Import) (err error) {
	if d.imports == nil {
		d.imports = make(map[*Import]struct{})
	}
	if _, ok := d.imports[i]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := i.validate(); err != nil {
		return err
	}
	d.imports[i] = struct{}{}
	return nil
}

func (d *Document) insertService(s *Service) (err error) {
	if _, ok := d.services[s]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := s.validate(); err != nil {
		return err
	}
	d.services[s] = struct{}{}
	return nil
}

func (d *Document) insertMessage(m *message) (err error) {
	if d.messages == nil {
		d.messages = make(map[*message]struct{})
	}
	if _, ok := d.messages[m]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := m.validate(); err != nil {
		return err
	}
	d.messages[m] = struct{}{}
	return nil
}

func (d *Document) insertEnum(e *enum) (err error) {
	if d.enums == nil {
		d.enums = make(map[*enum]struct{})
	}
	if _, ok := d.enums[e]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := e.validate(); err != nil {
		return err
	}
	d.enums[e] = struct{}{}
	return nil
}

func (d Document) validateLabel(l *Label) error {
	for s := range d.services {
		if s.hasLabel(l) {
			// TODO: return error type which contains other declaration
			return fmt.Errorf("label %s already declared for a service", l.value)
		}
	}
	for m := range d.messages {
		if m.hasLabel(l) {
			// TODO: return error type which contains other declaration
			return fmt.Errorf("label %s already declared for other message", l.value)
		}
	}
	for e := range d.enums {
		if e.hasLabel(l) {
			// TODO: return error type which contains other declaration
			return fmt.Errorf("label %s already declared for other enum", l.value)
		}
	}
	return nil
}

type Package struct {
	label  Label
	parent *Document
}

func (p Package) Get() string {
	return p.label.Get()
}

func (p *Package) Set(value string) error {
	return p.label.Set(value)
}

func (p *Package) Unset() error {
	// TODO: check if there is a condition where unsetting is impossible
	p.label.value = ""
	return nil
}

func (p *Package) Parent() *Document {
	return p.parent
}

func (p *Package) validateLabel(l *Label) error {
	return nil
}

type Import struct {
	parent *Document
	path   Label
	public *Flag
}

func (i Import) Path() *Label {
	return &i.path
}

func (i Import) Public() *Flag {
	return i.public
}

func (i *Import) InsertIntoParent() error {
	return i.parent.insertImport(i)
}

func (i Import) Parent() *Document {
	return i.parent
}

func (i Import) validate() error {
	return i.parent.validateLabel(&i.path)
}
