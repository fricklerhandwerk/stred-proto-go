package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fricklerhandwerk/stred-proto/protobuf"
)

func TestProtocolSetPackage(t *testing.T) {
	p := protobuf.NewDocument()
	err := p.SetPackage("package")
	require.Nil(t, err)
	assert.Equal(t, "package", p.GetPackage().String())
}

func TestProtocolSetInvalidPackage(t *testing.T) {
	p := protobuf.NewDocument()
	err := p.SetPackage("package!")
	require.NotNil(t, err)
	assert.Nil(t, p.GetPackage())
}

func TestMessageSetLabel(t *testing.T) {
	p := protobuf.NewDocument()
	m := p.NewMessage()
	err := m.SetLabel("message")
	require.Nil(t, err)
	assert.Equal(t, "message", m.GetLabel())
}

func TestMessageSetInvalidPackage(t *testing.T) {
	p := protobuf.NewDocument()
	m := p.NewMessage()

	err := m.SetLabel("message!")
	require.NotNil(t, err)
	assert.Empty(t, m.GetLabel())
}

func TestTypedFieldSetProperties(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()
	f := m.NewField()

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

func TestMessageAddField(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()
	err := m.SetLabel("message")
	require.Nil(t, err)

	f := m.NewField()
	err = f.SetLabel("messageField")
	require.Nil(t, err)
	err = f.SetNumber(1)
	require.Nil(t, err)
	f.SetType(protobuf.Bool)

	err = m.InsertField(0, f)
	require.Nil(t, err)
	assert.NotEmpty(t, m.GetFields())
}

func TestEnumAddField(t *testing.T) {
	e := protobuf.NewDocument().NewEnum()

	f := e.NewField()
	err := f.SetLabel("enumValue")
	require.Nil(t, err)
	err = f.SetNumber(1)
	require.Nil(t, err)

	err = e.InsertField(0, f)
	require.Nil(t, err)
	assert.NotEmpty(t, e.GetFields())
}

func TestEnumerationSetInvalidProperties(t *testing.T) {
	e := protobuf.NewDocument().NewEnum()

	f1 := e.NewField()
	err := f1.SetLabel("enumValue1")
	require.Nil(t, err)
	err = f1.SetNumber(1)
	require.Nil(t, err)

	// do not forget to add the first field to the enum!
	err = e.InsertField(0, f1)
	require.Nil(t, err)

	// duplicate label
	f2 := e.NewField()
	err = f2.SetLabel("enumValue1")
	require.NotNil(t, err)
	err = f2.SetLabel("enumValue2")
	require.Nil(t, err)

	//duplicate field number
	err = f2.SetNumber(1)
	require.NotNil(t, err)
	err = f2.SetNumber(2)
	require.Nil(t, err)

	// duplicate field number with "allow_alias = true"
	err = e.SetAlias(true)
	require.Nil(t, err)

	// try to disable aliasing with duplicate field numbers
	err = f2.SetNumber(1)
	require.Nil(t, err)
	e.InsertField(1, f2)
	err = e.SetAlias(false)
	require.NotNil(t, err)
}

func TestTypedFieldSetInvalidProperties(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()

	f1 := m.NewField()
	err := f1.SetLabel("messageField")
	require.Nil(t, err)
	err = f1.SetNumber(1)
	require.Nil(t, err)
	f1.SetType(protobuf.Bool)

	m.InsertField(0, f1)

	// duplicate label
	f2 := m.NewField()
	err = f2.SetLabel("messageField")
	require.NotNil(t, err)
	err = f2.SetNumber(2)
	assert.Nil(t, err)
	f2.SetType(protobuf.Bool)

	// duplicate field number
	err = f2.SetLabel("messageField2")
	require.Nil(t, err)
	err = f2.SetNumber(1)
	assert.NotNil(t, err)
}

func TestEnumAddInvalidField(t *testing.T) {
	e := protobuf.NewDocument().NewEnum()
	f1 := e.NewField()

	// label not set
	err := e.InsertField(0, f1)
	require.NotNil(t, err)

	// duplicate field number not checked
	err = f1.SetLabel("someLabel")
	require.Nil(t, err)
	err = e.InsertField(0, f1)
	require.Nil(t, err)

	f2 := e.NewField()
	err = f2.SetLabel("anotherLabel")
	require.Nil(t, err)
	err = e.InsertField(1, f2)
	require.NotNil(t, err)
}

func TestMessageInsertInvalidField(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()

	f1 := m.NewField()

	// label not set
	err := m.InsertField(0, f1)
	require.NotNil(t, err)

	f2 := m.NewMap()
	err = f2.SetLabel("someLabel")
	require.Nil(t, err)

	// field number not checked
	err = m.InsertField(0, f2)
	require.NotNil(t, err)
}
