package protobuf

import (
	"fmt"
)

func NewDocument() Document {
	return &document{}
}

type document struct {
	_package    *identifier
	imports     []_import
	services    []service
	definitions []Definition
}

func (p document) GetPackage() *string {
	// TODO: this mess is just another argument to just store the identifier as a string (pointer)
	if p._package == nil {
		return nil
	}
	s := p._package.String()
	return &s
}

func (p *document) SetPackage(pkg string) error {
	ident := identifier(pkg)
	if err := ident.validate(); err != nil {
		return err
	}
	p._package = &ident
	return nil
}

func (p document) validateLabel(l identifier) error {
	if err := l.validate(); err != nil {
		return err
	}
	for _, d := range p.definitions {
		if d.GetLabel() == l.String() {
			return fmt.Errorf("label %s already declared for other %T", l.String(), d)
		}
	}
	for _, s := range p.services {
		if s.GetLabel() == l.String() {
			return fmt.Errorf("label %s already declared for a service", l.String())
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

func (p *document) insertDefinition(i uint, d Definition) {
	p.definitions = append(p.definitions, nil)
	copy(p.definitions[i+1:], p.definitions[i:])
	p.definitions[i] = d
}

func (p *document) NewService() Service {
	return &service{
		label: label{
			parent: p,
		},
	}
}

func (p *document) NewMessage() Message {
	return &message{
		parent: p,
		label: label{
			parent: p,
		},
	}
}

func (p *document) NewEnum() Enum {
	return &enum{
		parent: p,
		label: label{
			parent: p,
		},
	}
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
