package protobuf

import (
	"errors"
	"fmt"
)

func NewDocument() *document {
	return &document{}
}

type document struct {
	_package    *identifier
	imports     []_import
	services    []service
	definitions []definition
}

func (p document) GetPackage() *identifier {
	return p._package
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
			return errors.New(fmt.Sprintf("label %s already declared for other %T", l.String(), d))
		}
	}
	for _, s := range p.services {
		if s.GetLabel() == l.String() {
			return errors.New(fmt.Sprintf("label %s already declared for a service", l.String()))
		}
	}

	return nil
}

func (p document) GetDefinitions() []definition {
	out := make([]definition, len(p.definitions))
	for i, v := range p.definitions {
		out[i] = v.(definition)
	}
	return out
}

func (p *document) InsertDefinition(i uint, d definition) error {
	if err := d.validateAsDefinition(); err != nil {
		return err
	}
	p.definitions = append(p.definitions, nil)
	copy(p.definitions[i+1:], p.definitions[i:])
	p.definitions[i] = d
	return nil
}

func (p *document) NewService() *service {
	return &service{
		label: label{
			parent: p,
		},
	}
}

func (p *document) NewMessage() *message {
	return &message{
		parent: p,
		label: label{
			parent: p,
		},
	}
}

func (p *document) NewEnum() enum {
	return enum{
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