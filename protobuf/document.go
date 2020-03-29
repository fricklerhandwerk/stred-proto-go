package protobuf

import (
	"fmt"
)

func NewDocument() Document {
	return &document{}
}

type document struct {
	// TODO: there is no fixed order in the syntax, so separating it this way is
	// an artificial constraint for the textual representation. we have the
	// following options:
	//
	// 1. just have one field `[]documentItem` and let the appropriate types
	//    implement `documentItem`.
	// 2. keep the field layout as is, but somehow keep track of the position of
	//    each item in the document.
	//
	// in any case we need to rework semantics of `InsertIntoParent(index uint)
	// error` to insert an object not at a position within a list of that type,
	// but within the whole parent. also the validation would look slightly
	// different, requiring a `validateAsDocumentItem()` method on `documentItem`
	// and maybe some specialised (private?) accessors on `Document`. the same
	// would apply for `Message` for example.
	// there should also be a way to move items around, since actually deleting
	// and re-inserting them is not generally possible, because references would
	// break.
	_package    *label
	imports     []_import
	services    []*service
	definitions []Definition
}

func (d document) MaybePackage() *string {
	return d._package.maybeLabel()
}

func (d *document) SetPackage(pkg string) (err error) {
	if err := validateIdentifier(pkg); err != nil {
		return err
	}
	if d._package == nil {
		d._package = &label{}
	}
	d._package.value = pkg
	return nil
}

func (d document) validateLabel(l *label) error {
	for _, def := range d.definitions {
		if def.hasLabel(l) {
			return fmt.Errorf("label %s already declared for other %T", l.value, def)
		}
	}
	for _, s := range d.services {
		if s.hasLabel(l) {
			return fmt.Errorf("label %s already declared for a service", l.value)
		}
	}

	return nil
}

func (d document) NumImports() uint {
	return uint(len(d.imports))
}

func (d document) Import(i uint) Import {
	return d.imports[i]
}

func (d document) NumServices() uint {
	return uint(len(d.services))
}

func (d document) Service(i uint) Service {
	return d.services[i]
}

func (d *document) insertService(i uint, s *service) error {
	if err := s.validateAsService(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	d.services = append(d.services, nil)
	copy(d.services[i+1:], d.services[i:])
	d.services[i] = s
	return nil
}

// TODO: use a common implementation for definition containers
func (d document) NumDefinitions() uint {
	return uint(len(d.definitions))
}

func (d document) Definition(i uint) Definition {
	return d.definitions[i]
}

func (d *document) insertDefinition(i uint, def Definition) error {
	if err := def.validateAsDefinition(); err != nil {
		// TODO: still counting on this becoming a panic instead
		return err
	}
	d.definitions = append(d.definitions, nil)
	copy(d.definitions[i+1:], d.definitions[i:])
	d.definitions[i] = def
	return nil
}

func (d *document) NewImport() NewImport {
	return &newImport{
		_import: &_import{
			parent: d,
		},
	}
}

func (d *document) NewService() NewService {
	return &newService{
		service: &service{
			parent: d,
		},
	}
}

func (d *document) NewMessage() NewMessage {
	return &newMessage{
		message: &message{parent: d},
	}
}

func (d *document) NewEnum() NewEnum {
	return &newEnum{
		enum: &enum{parent: d},
	}
}

type newImport struct {
	_import *_import
}

func (i newImport) MaybePath() *string {
	if i._import == nil {
		return nil
	}
	out := i._import.path
	return &out
}

func (i newImport) SetPath(string) error {
	panic("not implemented")
}

func (i newImport) Public() bool {
	return i._import.public
}

func (i newImport) SetPublic(b bool) error {
	panic("not implemented")
}

func (i newImport) InsertIntoParent(index uint) error {
	panic("not implemented")
}

type _import struct {
	parent *document
	path   string
	public bool
}

func (i _import) Path() string {
	return i.path
}

func (i _import) SetPath(string) error {
	panic("not implemented")
}

func (i _import) Public() bool {
	return i.public
}

func (i _import) SetPublic(b bool) error {
	panic("not implemented")
}
