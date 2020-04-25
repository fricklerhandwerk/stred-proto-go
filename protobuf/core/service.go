package core

import "fmt"

type Service struct {
	label  Label
	rpcs   map[*RPC]struct{}
	parent *Document
}

func (s *Service) Label() *Label {
	if s.label.parent == nil {
		s.label.parent = s
	}
	return &s.label
}

func (s *Service) RPCs() (out []*RPC) {
	out = make([]*RPC, len(s.rpcs))
	for r := range s.rpcs {
		out = append(out, r)
	}
	return
}

func (s *Service) NewRPC() *RPC {
	panic("not implemented")
}

func (s *Service) InsertIntoParent() error {
	return s.parent.insertService(s)
}

func (s *Service) Parent() *Document {
	return s.parent
}

func (s *Service) hasLabel(l *Label) bool {
	return s.label.hasLabel(l)
}

func (s *Service) insertRPC(r *RPC) error {
	if _, ok := s.rpcs[r]; ok {
		return fmt.Errorf("already inserted")
	}
	if err := r.validate(); err != nil {
		return err
	}
	s.rpcs[r] = struct{}{}
	return nil
}

func (s *Service) validate() error {
	if s.label.value == "" {
		return fmt.Errorf("label not set")
	}
	if err := s.parent.validateLabel(&s.label); err != nil {
		return err
	}
	for r := range s.rpcs {
		if err := r.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) validateLabel(l *Label) error {
	// TODO: rpc labels and rpc argument/return types share a namespace with
	// *unqualified* message/enum labels within a service. you can only have "rpc
	// Foo" and use "message Foo" as an argument/return type in one of the same
	// service's rpcs if you use the qualified message label "rpc Foo
	// (package.Foo)", but then you *must* have a package name
	switch l {
	case &s.label:
		return s.parent.validateLabel(l)
	default:
		for r := range s.rpcs {
			// TODO: return error type with reference to other declaration
			if r.label.hasLabel(l) {
				return fmt.Errorf("label %q already declared", l.value)
			}
		}
	}
	return nil
}

type RPC struct {
	label    Label
	request  MessageType
	response MessageType
	parent   *Service
}

func (r *RPC) Label() *Label {
	return &r.label
}

func (r *RPC) Request() MessageType {
	return r.request
}

func (r *RPC) Response() MessageType {
	return r.response
}

func (r *RPC) InsertIntoParent() error {
	return r.parent.insertRPC(r)
}

func (r *RPC) Parent() *Service {
	return r.parent
}

func (r *RPC) hasLabel(*Label) bool {
	panic("not implemented")
}

func (r *RPC) validate() error {
	panic("not implemented")
}
