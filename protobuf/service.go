package protobuf

type newService struct {
	service *service
}

func (s *newService) InsertIntoParent(i uint) error {
	if err := s.service.parent.insertService(i, s.service); err != nil {
		return err
	}
	s.service = &service{
		parent: s.service.parent,
	}
	return nil
}

func (s *newService) MaybeLabel() *string {
	return s.service.maybeLabel()
}

func (s *newService) SetLabel(l string) (err error) {
	if s.service.label == nil {
		s.service.label = &label{
			parent: s.service.parent,
		}
		defer func() {
			if err != nil {
				s.service.label = nil
			}
		}()
	}
	return s.service.SetLabel(l)
}

func (s *newService) NewRPC() NewRPC {
	return s.service.NewRPC()
}

func (s *newService) NumRPCs() uint {
	return s.service.NumRPCs()
}

func (s *newService) RPC(i uint) RPC {
	return s.service.RPC(i)
}

type service struct {
	parent *document
	*label
	rpcs []*rpc
}

func (s *service) NumRPCs() uint {
	return uint(len(s.rpcs))
}

func (s *service) RPC(i uint) RPC {
	return s.rpcs[i]
}

func (s *service) NewRPC() NewRPC {
	panic("not implemented")
}

func (s *service) hasLabel(l *label) bool {
	for _, r := range s.rpcs {
		if r.hasLabel(l) {
			return true
		}
	}
	return false
}

func (s *service) insertRPC(i uint, r *rpc) error {
	panic("not implemented")
}

func (s *service) validateAsService() error {
	panic("not implemented")
}

func (s *service) validateLabel(l *label) error {
	panic("not implemented")
}

type rpc struct {
	// TODO: rpc labels and rpc argument/return types share a namespace with
	// *unqualified* message/enum labels within a service. you can only have "rpc
	// Foo" and use "message Foo" as an argument/return type in one of the same
	// service's rpcs if you use the qualified message label "rpc Foo
	// (package.Foo)", but then you *must* have a package name
	*label
	requestType    Message
	streamRequest  bool
	responseType   Message
	streamResponse bool
}

func (r *rpc) RequestType() Message {
	return r.requestType
}

func (r *rpc) SetRequestType(t Message) error {
	r.requestType = t
	return nil
}

func (r *rpc) StreamRequest() bool {
	return r.streamRequest
}

func (r *rpc) SetStreamRequest(b bool) error {
	r.streamRequest = b
	return nil
}

func (r *rpc) ResponseType() Message {
	return r.responseType
}

func (r *rpc) SetResponseType(t Message) error {
	r.responseType = t
	return nil
}

func (r *rpc) StreamResponse() bool {
	return r.streamResponse
}

func (r *rpc) SetStreamResponse(b bool) error {
	r.streamResponse = b
	return nil
}

func (r *rpc) validateAsRPC() error {
	panic("not implemented")
}
