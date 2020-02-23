package main

import (
	"testing"

	"github.com/fricklerhandwerk/stred-proto/protobuf"
)

func TestProtocolSetPackage(t *testing.T) {
	p := protobuf.Protocol{}
	err := p.SetPackage("package")
	if err != nil {
		t.Fatalf("failed to set package name: %s", err)
	}

	if p.GetPackage().String() != "package" {
		t.Fatalf("expected package name %s, found %s", "someFoo", p.GetPackage().String())
	}
}

func TestProtocolSetInvalidPackage(t *testing.T) {
	p := protobuf.Protocol{}
	err := p.SetPackage("package!")
	if err == nil {
		t.Fatal("expected to fail setting invalid package name")
	}

	if p.GetPackage() != nil {
		t.Fatalf("package is not nil")
	}
}

func TestMessageSetLabel(t *testing.T) {
	m := protobuf.Message{}
	err := m.SetLabel("message")
	if err != nil {
		t.Fatalf("failed to set message label: %s", err)
	}
	if m.GetLabel() != "message" {
		t.Fatalf("expected message label %s, found %s", "message", m.GetLabel())
	}
}

func TestMessageSetInvalidPackage(t *testing.T) {
	m := protobuf.Message{}
	err := m.SetLabel("message!")
	if err == nil {
		t.Fatal("expected to fail setting invalid package name")
	}

	if m.GetLabel() != "" {
		t.Fatalf("package is not nil")
	}
}

func TestTypedFieldSetProperties(t *testing.T) {
	f := protobuf.TypedField{}
	err := f.SetLabel("typedField")
	if err != nil {
		t.Fatalf("failed to set message label: %s", err)
	}
	if f.GetLabel() != "typedField" {
		t.Fatalf("expected field label %s, found %s", "typedField", f.GetLabel())
	}
	err = f.SetNumber(1)
	if err != nil {
		t.Fatalf("failed to set field number: %s", err)
	}
	if f.GetNumber() != 1 {
		t.Fatalf("expected field number %d, found %d", 1, f.GetNumber())
	}
	f.SetDeprecated(true)
	if f.GetDeprecated() != true {
		t.Fatalf("expected field option \"deprecated\" %v, found %v", true, f.GetDeprecated())
	}
	f.SetType(protobuf.Int32)
	if f.GetType() != protobuf.Int32 {
		t.Fatalf("expected field type %v, found %v", protobuf.Int32, f.GetType())
	}
	f.SetType(protobuf.Bytes)
	if f.GetType() != protobuf.Bytes {
		t.Fatalf("expected field type %v, found %v", protobuf.Bytes, f.GetType())
	}
}
