package protobuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	protobuf "github.com/fricklerhandwerk/stred-proto/protobuf/core"
)

func TestProtocolSetPackage(t *testing.T) {
	p := protobuf.NewDocument().Package()
	err := p.Set("package")
	require.Nil(t, err)
	assert.Equal(t, "package", p.Get())
	err = p.Unset()
	require.Nil(t, err)
	assert.Empty(t, p.Get())
}

func TestProtocolSetInvalidPackage(t *testing.T) {
	p := protobuf.NewDocument().Package()
	err := p.Set("package!")
	require.NotNil(t, err)
	assert.Empty(t, p.Get())
}

func TestProtocolDuplicateLabels(t *testing.T) {
	p := protobuf.NewDocument()
	m := p.NewMessage()
	err := m.Label().Set("foo")
	require.Nil(t, err)
	err = m.InsertIntoParent()
	require.Nil(t, err)

	e := p.NewEnum()
	err = e.Label().Set("foo")
	assert.NotNil(t, err)
	s := p.NewService()
	err = s.Label().Set("foo")
	assert.NotNil(t, err)
}

func TestMessageSetLabel(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()
	err := m.Label().Set("message")
	require.Nil(t, err)
	assert.Equal(t, "message", m.Label().Get())
}

func TestMessageInsertValid(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()
	err := m.Label().Set("message")
	require.Nil(t, err)
	err = m.InsertIntoParent()
	assert.Nil(t, err)
}

func TestMessageSetInvalidLabel(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()

	err := m.Label().Set("message!")
	require.NotNil(t, err)
	assert.Empty(t, m.Label().Get())
}

func TestMessageInsertInvalid(t *testing.T) {
	m := protobuf.NewDocument().NewMessage()
	err := m.InsertIntoParent()
	assert.NotNil(t, err)
}

func TestTypedFieldSetProperties(t *testing.T) {
	p := protobuf.NewDocument()
	nm := p.NewMessage()
	err := nm.Label().Set("message")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	m := p.Messages()[0]

	f := m.NewField()

	err = f.Label().Set("typedField")
	require.Nil(t, err)
	assert.Equal(t, "typedField", f.Label().Get())

	err = f.Number().Set(1)
	require.Nil(t, err)
	require.NotNil(t, 1, f.Number().Get())
	assert.EqualValues(t, 1, *f.Number().Get())

	err = f.Deprecated().Set(true)
	require.Nil(t, err)
	assert.Equal(t, true, f.Deprecated().Get())

	err = f.Type().Set(protobuf.Int32)
	require.Nil(t, err)
	assert.Equal(t, protobuf.Int32, f.Type().Get())

	err = f.Repeated().Set(true)
	require.Nil(t, err)
	assert.Equal(t, true, f.Repeated().Get())
}

func TestMessageAddField(t *testing.T) {
	p := protobuf.NewDocument()
	nm := p.NewMessage()
	err := nm.Label().Set("message")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	// TODO: maybe the signature should be
	// `NewMessage.InsertIntoParent() (Message, error)` so we can get
	// a handle on the new object, since as opposed to all other types
	// we cannot fetch it from the parent by pointer comparison
	m := p.Messages()[0]

	f := m.NewField()
	err = f.Label().Set("messageField")
	require.Nil(t, err)
	err = f.Number().Set(1)
	require.Nil(t, err)
	err = f.Type().Set(protobuf.Bool)
	require.Nil(t, err)

	err = f.InsertIntoParent()
	require.Nil(t, err)
	assert.NotEmpty(t, m.Fields())
}

func TestEnumAddField(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	f := e.NewVariant()
	err = f.Label().Set("enumValue")
	require.Nil(t, err)
	err = f.Number().Set(1)
	require.Nil(t, err)

	err = f.InsertIntoParent()
	require.Nil(t, err)
	assert.NotEmpty(t, e.Fields())
}

func TestEnumSetLabel(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	err = e.Label().Set("foo")
	require.Nil(t, err)
	assert.Equal(t, "foo", e.Label().Get())
}

func TestVariantSetInvalidProperties(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	f1 := e.NewVariant()
	err = f1.Label().Set("enumValue1")
	require.Nil(t, err)
	err = f1.Number().Set(1)
	require.Nil(t, err)
	require.NotNil(t, f1.Number().Get())
	require.EqualValues(t, 1, *f1.Number().Get())

	// do not forget to add the first field to the enum!
	err = f1.InsertIntoParent()
	require.Nil(t, err)

	// duplicate label
	f2 := e.NewVariant()
	err = f2.Label().Set("enumValue1")
	require.NotNil(t, err)
	err = f2.Label().Set("enumValue2")
	require.Nil(t, err)

	//duplicate field number
	err = f2.Number().Set(1)
	require.NotNil(t, err)
	err = f2.Number().Set(2)
	require.Nil(t, err)

	// duplicate field number with "allow_alias = true"
	err = e.AllowAlias().Set(true)
	require.Nil(t, err)

	// try to disable aliasing with duplicate field numbers
	err = f2.Number().Set(1)
	require.Nil(t, err)
	err = f2.InsertIntoParent()
	require.Nil(t, err)
	err = e.AllowAlias().Set(false)
	require.NotNil(t, err)
	// keep aliasing enabled
	err = e.AllowAlias().Set(true)
	require.Nil(t, err)

	// disable aliasing after none is left
	err = f2.Number().Set(2)
	require.Nil(t, err)
	err = e.AllowAlias().Set(false)
	require.Nil(t, err)
}

func TestTypedFieldSetInvalidProperties(t *testing.T) {
	p := protobuf.NewDocument()
	nm := p.NewMessage()
	err := nm.Label().Set("message")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	m := p.Messages()[0]

	f1 := m.NewField()
	err = f1.Label().Set("messageField")
	require.Nil(t, err)
	err = f1.Number().Set(1)
	require.Nil(t, err)
	err = f1.Type().Set(protobuf.Bool)
	require.Nil(t, err)

	err = f1.InsertIntoParent()
	require.Nil(t, err)

	// duplicate label
	f2 := m.NewField()
	err = f2.Label().Set("messageField")
	require.NotNil(t, err)
	err = f2.Number().Set(2)
	assert.Nil(t, err)
	f2.Type().Set(protobuf.Bool)

	// duplicate field number
	err = f2.Label().Set("messageField2")
	require.Nil(t, err)
	err = f2.Number().Set(1)
	assert.NotNil(t, err)
}

func TestEnumAddInvalidField(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	f1 := e.NewVariant()

	// label not set
	err = f1.Number().Set(0)
	require.Nil(t, err)
	err = f1.InsertIntoParent()
	require.NotNil(t, err)

	// number not set
	f2 := e.NewVariant()
	err = f2.Label().Set("someLabel")
	require.Nil(t, err)
	err = f2.InsertIntoParent()
	require.NotNil(t, err)

	// duplicate field number not checked
	err = f1.Label().Set("anotherLabel")
	require.Nil(t, err)
	err = f1.InsertIntoParent()
	require.Nil(t, err)

	err = f2.Number().Set(0)
	require.NotNil(t, err)
}

func TestMessageInsertInvalidField(t *testing.T) {
	p := protobuf.NewDocument()
	nm := p.NewMessage()
	err := nm.Label().Set("message")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	m := p.Messages()[0]

	f1 := m.NewField()

	// label not set
	err = f1.InsertIntoParent()
	require.NotNil(t, err)

	f2 := m.NewMap()
	err = f2.Label().Set("someLabel")
	require.Nil(t, err)

	// field number not checked
	err = f2.InsertIntoParent()
	require.NotNil(t, err)
}

func TestMessageValidateReservedNumber(t *testing.T) {
	p := protobuf.NewDocument()
	nm := p.NewMessage()
	err := nm.Label().Set("message")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	m := p.Messages()[0]

	r := m.NewReservedNumber()

	err = r.InsertIntoParent()
	require.NotNil(t, err)
	err = r.Set(0)
	require.NotNil(t, err)
	err = r.Set(1)
	require.Nil(t, err)
	require.EqualValues(t, 1, *r.Get())
	err = r.InsertIntoParent()
	require.Nil(t, err)

	nr := m.NewReservedRange()
	err = nr.Start().Set(0)
	require.NotNil(t, err)
	err = nr.Start().Set(1)
	require.NotNil(t, err)
	err = nr.Start().Set(2)
	require.Nil(t, err)
	err = nr.End().Set(10)
	require.Nil(t, err)
	require.EqualValues(t, 2, *nr.Start().Get())
	require.EqualValues(t, 10, *nr.End().Get())
	err = nr.InsertIntoParent()
	require.Nil(t, err)

	nr2 := m.NewReservedRange()
	err = nr2.Start().Set(11)
	require.EqualValues(t, 11, *nr2.Start().Get())
	require.Nil(t, err)
	err = nr2.End().Set(20)
	require.Nil(t, err)
	err = nr.InsertIntoParent()
	require.NotNil(t, err)

	err = nr2.Start().Set(10)
	require.NotNil(t, err)

	err = r.Set(21)
	require.Nil(t, err)

	err = nr2.End().Set(21)
	require.NotNil(t, err)

	err = r.InsertIntoParent()
	// already inserted
	assert.NotNil(t, err)
}

func TestEnumValidateReservedNumber(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	r := e.NewReservedNumber()

	err = r.InsertIntoParent()
	require.NotNil(t, err)
	// enum may have 0 as field number
	err = r.Set(0)
	require.Nil(t, err)
	require.EqualValues(t, 0, *r.Get())
	err = r.InsertIntoParent()
	require.Nil(t, err)

	nr := e.NewReservedRange()
	err = nr.Start().Set(0)
	require.NotNil(t, err)
	err = nr.Start().Set(1)
	require.Nil(t, err)
	err = nr.Start().Set(2)
	require.Nil(t, err)
	err = nr.End().Set(10)
	require.Nil(t, err)
	require.EqualValues(t, 2, *nr.Start().Get())
	require.EqualValues(t, 10, *nr.End().Get())
	err = nr.InsertIntoParent()
	require.Nil(t, err)

	nr2 := e.NewReservedRange()
	err = nr2.Start().Set(11)
	require.EqualValues(t, 11, *nr2.Start().Get())
	require.Nil(t, err)
	err = nr2.End().Set(20)
	require.Nil(t, err)
	err = nr.InsertIntoParent()
	require.NotNil(t, err)

	err = nr2.Start().Set(10)
	require.NotNil(t, err)

	err = r.Set(21)
	require.Nil(t, err)

	err = nr2.End().Set(21)
	require.NotNil(t, err)

	err = r.InsertIntoParent()
	// already inserted
	assert.NotNil(t, err)
}

func TestMessageValidateDefinition(t *testing.T) {
	p := protobuf.NewDocument()
	nm := p.NewMessage()
	err := nm.Label().Set("myMessage")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	m := p.Messages()[0]

	f1 := m.NewField()
	err = f1.Label().Set("myField")
	require.Nil(t, err)
	err = f1.Number().Set(1)
	require.Nil(t, err)
	f1.Type().Set(m)
	err = f1.InsertIntoParent()
	require.Nil(t, err)
	require.EqualValues(t, 1, len(m.Fields()))

	f2 := m.NewField()
	err = f2.Label().Set("myNewField")
	require.Nil(t, err)
	err = f2.Number().Set(2)
	require.Nil(t, err)
	f2.Type().Set(m)
	err = f2.InsertIntoParent()
	require.Nil(t, err)
	require.EqualValues(t, 2, len(m.Fields()))
}

func TestEnumValidateDefinition(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	err = e.Label().Set("myEnum")
	require.Nil(t, err)

	f1 := e.NewVariant()
	err = f1.Label().Set("myField")
	require.Nil(t, err)
	err = f1.Number().Set(0)
	require.Nil(t, err)
	err = f1.InsertIntoParent()
	require.Nil(t, err)
	require.EqualValues(t, 1, len(e.Fields()))

	f2 := e.NewVariant()
	err = f2.Label().Set("myNewField")
	require.Nil(t, err)
	err = f2.Number().Set(1)
	require.Nil(t, err)
	err = f2.InsertIntoParent()
	require.Nil(t, err)
	require.EqualValues(t, 2, len(e.Fields()))
	for _, f := range e.Fields() {
		switch v := f.(type) {
		case *protobuf.Variant:
			require.Contains(t, []string{"myField", "myNewField"}, v.Label().Get())
		default:
			t.Errorf("unexpected enum field %v", v)
		}
	}
}

func TestInsertIncompleteRange(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	r := e.NewReservedRange()
	err = r.Start().Set(1)
	require.Nil(t, err)
	err = r.InsertIntoParent()
	require.NotNil(t, err)

	r = e.NewReservedRange()
	err = r.End().Set(1)
	require.Nil(t, err)
	err = r.InsertIntoParent()
	require.NotNil(t, err)
}

func TestValidateReservedLabels(t *testing.T) {
	p := protobuf.NewDocument()
	ne := p.NewEnum()
	err := ne.Label().Set("enum")
	require.Nil(t, err)
	err = ne.InsertIntoParent()
	require.Nil(t, err)
	e := p.Enums()[0]

	rl := e.NewReservedLabel()
	err = rl.InsertIntoParent()
	require.NotNil(t, err)
	err = rl.Set("invalid!")
	require.NotNil(t, err)
	err = rl.InsertIntoParent()
	require.NotNil(t, err)
	err = rl.Set("someLabel1")
	require.Nil(t, err)
	err = rl.InsertIntoParent()
	require.Nil(t, err)
	err = rl.Set("someLabel")
	require.Nil(t, err)
	// already inserted
	err = rl.InsertIntoParent()
	require.NotNil(t, err)

	rl = e.NewReservedLabel()
	err = rl.InsertIntoParent()
	require.NotNil(t, err)
	// label in use
	err = rl.Set("someLabel")
	require.NotNil(t, err)
	err = rl.InsertIntoParent()
	require.NotNil(t, err)
	err = rl.Set("someLabel2")
	require.Nil(t, err)
	err = rl.InsertIntoParent()
	require.Nil(t, err)
	err = rl.Set("someLabel")
	require.NotNil(t, err)
	err = rl.Set("someOtherLabel")
	require.Nil(t, err)
}

func TestOneOf(t *testing.T) {
	p := protobuf.NewDocument()
	nm := p.NewMessage()
	err := nm.Label().Set("myMessage")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	m := p.Messages()[0]

	o1 := m.NewOneOf()
	err = o1.Label().Set("invalid!")
	require.NotNil(t, err)
	// label not set
	err = o1.InsertIntoParent()
	require.NotNil(t, err)

	err = o1.Label().Set("myOneOf")
	require.Nil(t, err)
	f1 := o1.NewField()
	err = f1.Label().Set("foo")
	require.Nil(t, err)
	err = f1.Number().Set(1)
	require.Nil(t, err)
	// type not set
	err = f1.InsertIntoParent()
	require.NotNil(t, err)
	err = f1.Type().Set(protobuf.Uint32)
	require.Nil(t, err)
	err = f1.InsertIntoParent()
	require.Nil(t, err)
	require.EqualValues(t, 1, len(o1.Fields()))

	err = o1.InsertIntoParent()
	require.Nil(t, err)
	require.EqualValues(t, 1, len(m.Fields()))

	o2 := m.NewOneOf()
	// duplicate label
	err = o2.Label().Set("myOneOf")
	require.NotNil(t, err)
	err = o2.Label().Set("myOneOf2")
	require.Nil(t, err)
	f2 := o2.NewField()
	// duplicate field label
	err = f2.Label().Set("foo")
	require.NotNil(t, err)
	err = f2.Label().Set("bar")
	require.Nil(t, err)
	// duplicate field number
	err = f2.Number().Set(1)
	require.NotNil(t, err)
	err = f2.Number().Set(2)
	require.Nil(t, err)

	err = f2.Type().Set(protobuf.Bytes)
	require.Nil(t, err)
	err = f2.InsertIntoParent()
	require.Nil(t, err)
}

func TestServiceWithRPC(t *testing.T) {
	d := protobuf.NewDocument()
	nm := d.NewMessage()
	err := nm.Label().Set("Foo")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)
	nm = d.NewMessage()
	err = nm.Label().Set("Bar")
	require.Nil(t, err)
	err = nm.InsertIntoParent()
	require.Nil(t, err)

	s := d.NewService()
	err = s.Label().Set("FooService")
	require.Nil(t, err)
	err = s.InsertIntoParent()
	require.Nil(t, err)
	r := s.NewRPC()
	r.Label().Set("GetFoo")
	err = r.Request().Set(d.Messages()[0])
	require.Nil(t, err)
	err = r.Request().Stream().Set(true)
	require.Nil(t, err)
	err = r.Response().Set(d.Messages()[1])
	require.Nil(t, err)
	err = r.Response().Stream().Set(true)
	require.Nil(t, err)
	err = r.InsertIntoParent()
	require.Nil(t, err)
	err = r.Response().Set(d.Messages()[0])
	require.Nil(t, err)
	err = r.Request().Set(d.Messages()[0])
	require.Nil(t, err)
	err = r.Request().Stream().Set(false)
	require.Nil(t, err)
	err = r.Response().Stream().Set(false)
	require.Nil(t, err)
}
