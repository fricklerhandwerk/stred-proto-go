package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fricklerhandwerk/stred-proto/protobuf"
)

func TestProtocolSetPackage(t *testing.T) {
	p := protobuf.Protocol{}
	err := p.SetPackage("package")
	require.Nil(t, err)
	assert.Equal(t, "package", p.GetPackage().String())
}

func TestProtocolSetInvalidPackage(t *testing.T) {
	p := protobuf.Protocol{}
	err := p.SetPackage("package!")
	require.NotNil(t, err)
	assert.Nil(t, p.GetPackage())
}

func TestMessageSetLabel(t *testing.T) {
	m := protobuf.Message{}
	err := m.SetLabel("message")
	require.Nil(t, err)
	assert.Equal(t, "message", m.GetLabel())
}

func TestMessageSetInvalidPackage(t *testing.T) {
	m := protobuf.Message{}
	err := m.SetLabel("message!")
	require.NotNil(t, err)
	assert.Empty(t, m.GetLabel())
}

func TestTypedFieldSetProperties(t *testing.T) {
	f := protobuf.TypedField{}

	err := f.SetLabel("typedField")
	require.Nil(t, err)
	assert.Equal(t, "typedField", f.GetLabel())

	err = f.SetNumber(1)
	require.Nil(t, err)
	assert.EqualValues(t, 1, f.GetNumber())

	f.SetDeprecated(true)
	assert.Equal(t, true, f.GetDeprecated())

	f.SetType(protobuf.Int32)
	assert.Equal(t, protobuf.Int32, f.GetType())

	f.SetType(protobuf.Bytes)
	assert.Equal(t, protobuf.Bytes, f.GetType())
}
