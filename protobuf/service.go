package protobuf

type service struct {
	label
	rpcs []rpc
}

type rpc struct {
	// TODO: rpc labels and rpc argument/return types share a namespace with
	// *unqualified* message/enum labels within a service. you can only have "rpc
	// Foo" and use "message Foo" as an argument/return type in one of the same
	// service's rpcs if you use the qualified message label "rpc Foo
	// (package.Foo)", but then you *must* have a package name
	label
	requestType    *message
	streamRequest  bool
	responseType   *message
	streamResponse bool
}
