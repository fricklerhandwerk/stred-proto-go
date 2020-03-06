package protobuf

import (
	"fmt"
)

func NewDocument() Document {
	return &document{}
}

type document struct {
	_package    *label
	imports     []_import
	services    []*service
	definitions []Definition
}

func (p document) GetPackage() *string {
	// TODO: this mess is just another argument to just store the identifier as a string (pointer)
	if p._package == nil {
		return nil
	}
	s := p._package.value
	return &s
}

func (p *document) SetPackage(pkg string) error {
	label := &label{
		value:  pkg,
		parent: p,
	}
	if err := label.validate(); err != nil {
		return err
	}
	p._package = label
	return nil
}

func (p document) validateLabel(l *label) error {
	if err := l.validate(); err != nil {
		return err
	}
	for _, d := range p.definitions {
		if d.hasLabel(l) {
			return fmt.Errorf("label %s already declared for other %T", l.value, d)
		}
	}
	for _, s := range p.services {
		if s.hasLabel(l) {
			return fmt.Errorf("label %s already declared for a service", l.value)
		}
	}

	return nil
}

// TODO: use a common implementation for definition containers
func (p document) NumDefinitions() uint {
	return uint(len(p.definitions))
}

func (p document) Definition(i uint) Definition {
	return p.definitions[i]
}

func (p *document) insertDefinition(i uint, d Definition) error {
	if err := d.validateAsDefinition(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	p.definitions = append(p.definitions, nil)
	copy(p.definitions[i+1:], p.definitions[i:])
	p.definitions[i] = d
	return nil
}

func (p *document) NewService() Service {
	out := &service{
		label: label{
			parent: p,
		},
	}
	return out
}

func (p *document) NewMessage() Message {
	out := &message{
		parent: p,
		label: &label{
			parent: p,
		},
	}
	return out
}

func (p *document) NewEnum() Enum {
	out := &enum{
		parent: p,
		label: &label{
			parent: p,
		},
	}
	return out
}

type _import struct {
	path   string
	public bool
}

func (i _import) SetPath(string) error {
	panic("not implemented")
}

func (i _import) SetPublic(b bool) {
	panic("not implemented")
}
