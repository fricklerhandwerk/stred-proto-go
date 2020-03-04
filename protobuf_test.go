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
	pkg := p.GetPackage()
	require.NotNil(t, pkg)
	assert.Equal(t, "package", *pkg)
}

func TestProtocolSetInvalidPackage(t *testing.T) {
	p := protobuf.NewDocument()
	err := p.SetPackage("package!")
	require.NotNil(t, err)
	assert.Nil(t, p.GetPackage())
}

func TestProtocolDuplicateLabelsForMessageAndService(t *testing.T) {
	p := protobuf.NewDocument()
	m := p.NewMessage()
	err := m.SetLabel("foo")
	require.Nil(t, err)
	err = m.InsertIntoParent(0)
	require.Nil(t, err)
	s := p.NewService()
	err = s.SetLabel("foo")
	assert.NotNil(t, err)
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

	err = f.InsertIntoParent(0)
	require.Nil(t, err)
	assert.EqualValues(t, 1, m.NumFields())
}

func TestEnumAddField(t *testing.T) {
	e := protobuf.NewDocument().NewEnum()

	f := e.NewField()
	err := f.SetLabel("enumValue")
	require.Nil(t, err)
	err = f.SetNumber(1)
	require.Nil(t, err)

	err = f.InsertIntoParent(0)
	require.Nil(t, err)
	assert.EqualValues(t, 1, e.NumFields())
}

func TestEnumerationSetInvalidProperties(t *testing.T) {
	e := protobuf.NewDocument().NewEnum()

	f1 := e.NewField()
	err := f1.SetLabel("enumValue1")
	require.Nil(t, err)
	err = f1.SetNumber(1)
	require.Nil(t, err)

	// do not forget to add the first field to the enum!
	err = f1.InsertIntoParent(0)
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
	err = f2.InsertIntoParent(1)
	require.Nil(t, err)
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

	err = f1.InsertIntoParent(0)
	require.Nil(t, err)

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
	err := f1.InsertIntoParent(0)
	require.NotNil(t, err)

	// duplicate field number not checked
	err = f1.SetLabel("someLabel")
	require.Nil(t, err)
	err = f1.InsertIntoParent(0)
	require.Nil(t, err)

	f2 := e.NewField()
	err = f2.SetLabel("anotherLabel")
	require.Nil(t, err)
	err = f2.InsertIntoParent(1)
	require.NotNil(t, err)
}

func TestMessageInsertInvalidField(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()

	f1 := m.NewField()

	// label not set
	err := f1.InsertIntoParent(0)
	require.NotNil(t, err)

	f2 := m.NewMap()
	err = f2.SetLabel("someLabel")
	require.Nil(t, err)

	// field number not checked
	err = f2.InsertIntoParent(0)
	require.NotNil(t, err)
}

func TestMessageValidateReservedNumber(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()
	r := m.NewReservedNumbers()

	err := r.InsertNumber(0, 0)
	require.NotNil(t, err)

	err = r.InsertNumber(0, 1)
	require.Nil(t, err)
	require.EqualValues(t, 1, r.NumNumbers())

	err = r.ToRange(0, 10)
	require.Nil(t, err)

	err = r.InsertNumber(1, 2)
	require.NotNil(t, err)
	require.EqualValues(t, 1, r.NumNumbers())

	err = r.InsertNumber(1, 12)
	require.Nil(t, err)

	err = r.ToRange(1, 11)
	require.NotNil(t, err)

	err = r.ToRange(1, 20)
	require.Nil(t, err)

	nr, ok := r.Number(1).(protobuf.NumberRange)
	require.True(t, ok)
	err = nr.SetStart(9)
	require.NotNil(t, err)
	err = nr.SetStart(13)
	require.Nil(t, err)

	err = nr.SetEnd(13)
	require.NotNil(t, err)

	err = nr.SetEnd(22)
	require.Nil(t, err)
}

func TestEnumValidateReservedNumber(t *testing.T) {
	e := protobuf.NewDocument().NewEnum()
	r := e.NewReservedNumbers()
	err := r.InsertNumber(0, 0)
	require.Nil(t, err)
	require.EqualValues(t, 1, r.NumNumbers())

	err = r.ToRange(0, 10)
	require.Nil(t, err)

	err = r.InsertNumber(1, 2)
	require.NotNil(t, err)
	require.EqualValues(t, 1, r.NumNumbers())

	err = r.InsertNumber(1, 12)
	require.Nil(t, err)

	err = r.ToRange(1, 11)
	require.NotNil(t, err)

	err = r.ToRange(1, 20)
	require.Nil(t, err)

	nr, ok := r.Number(1).(protobuf.NumberRange)
	require.True(t, ok)
	err = nr.SetStart(9)
	require.NotNil(t, err)
	err = nr.SetStart(13)
	require.Nil(t, err)

	err = nr.SetEnd(13)
	require.NotNil(t, err)

	err = nr.SetEnd(22)
	require.Nil(t, err)
}

func TestMessageValidateRecursive(t *testing.T) {
	p := protobuf.NewDocument()

	m := p.NewMessage()
	err := m.SetLabel("myMessage")
	require.Nil(t, err)

	f1 := m.NewField()
	err = f1.SetLabel("myField")
	require.Nil(t, err)
	err = f1.SetNumber(1)
	require.Nil(t, err)
	f1.SetType(m)
	err = f1.InsertIntoParent(0)
	require.Nil(t, err)
	require.EqualValues(t, 1, m.NumFields())

	f2 := m.NewField()
	err = f2.SetLabel("myNewField")
	require.Nil(t, err)
	err = f2.SetNumber(2)
	require.Nil(t, err)
	f2.SetType(m)
	err = f2.InsertIntoParent(1)
	require.Nil(t, err)
	require.EqualValues(t, 2, m.NumFields())

	err = m.InsertIntoParent(0)
	require.Nil(t, err)
}

func TestEnumValidateRecursive(t *testing.T) {
	p := protobuf.NewDocument()

	e := p.NewEnum()
	err := e.SetLabel("myEnum")
	require.Nil(t, err)

	f1 := e.NewField()
	err = f1.SetLabel("myField")
	require.Nil(t, err)
	err = f1.SetNumber(0)
	require.Nil(t, err)
	err = f1.InsertIntoParent(0)
	require.Nil(t, err)
	require.EqualValues(t, 1, e.NumFields())

	f2 := e.NewField()
	err = f2.SetLabel("myNewField")
	require.Nil(t, err)
	err = f2.SetNumber(1)
	require.Nil(t, err)
	err = f2.InsertIntoParent(1)
	require.Nil(t, err)
	require.EqualValues(t, 2, e.NumFields())

	err = e.InsertIntoParent(0)
	require.Nil(t, err)
}