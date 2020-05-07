package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	d := NewDocument()
	d.Package().Set("somePackage")
	expected := "package somePackage;"
	assert.Equal(t, expected, d.Package().String())
}

func TestEnum(t *testing.T) {
	d := NewDocument()
	ne := d.NewEnum()
	ne.Label().Set("frooble")
	expected := "enum frooble {}"
	assert.Equal(t, expected, ne.String())
	ne.InsertIntoParent()
	e := d.Enums()[0]
	assert.Equal(t, expected, e.String())
}

func TestVariant(t *testing.T) {
	d := NewDocument()
	ne := d.NewEnum()
	ne.Label().Set("frooble")
	ne.InsertIntoParent()
	e := d.Enums()[0]
	v := e.NewVariant()
	v.Label().Set("foo")
	v.Number().Set(0)
	v.InsertIntoParent()
	expected := "enum frooble {\n  foo = 0;\n}"
	assert.Equal(t, expected, e.String())
}

func TestEnumAlias(t *testing.T) {
	d := NewDocument()
	ne := d.NewEnum()
	ne.Label().Set("frooble")
	ne.InsertIntoParent()
	e := d.Enums()[0]
	e.AllowAlias().Set(true)
	v := e.NewVariant()
	v.Label().Set("foo")
	v.Number().Set(0)
	v.InsertIntoParent()
	v = e.NewVariant()
	v.Label().Set("bar")
	v.Number().Set(1)
	v.InsertIntoParent()
	// output order is non-deterministic, won't implement sorting
	assert.Contains(t, e.String(), "enum frooble")
	assert.Contains(t, e.String(), "foo = 0;")
	assert.Contains(t, e.String(), "bar = 1;")
	assert.NotContains(t, e.String(), "allow_alias")
	v.Number().Set(0)
	assert.Contains(t, e.String(), "bar = 0;")
	assert.Contains(t, e.String(), "option allow_alias = true;")
}

func TestMessage(t *testing.T) {
	d := NewDocument()
	ne := d.NewMessage()
	ne.Label().Set("groogle")
	expected := "message groogle {}"
	assert.Equal(t, expected, ne.String())
	ne.InsertIntoParent()
	e := d.Messages()[0]
	assert.Equal(t, expected, e.String())
}

func TestDocument(t *testing.T) {
	d := NewDocument()
	d.Package().Set("somePackage")
	ne := d.NewEnum()
	ne.Label().Set("frooble")
	ne.InsertIntoParent()
	e := d.Enums()[0]
	v := e.NewVariant()
	v.Label().Set("foo")
	v.Number().Set(0)
	v.InsertIntoParent()
	v = e.NewVariant()
	v.Label().Set("bar")
	v.Number().Set(1)
	v.InsertIntoParent()
	e.AllowAlias().Set(true)
	v = e.NewVariant()
	v.Label().Set("baz")
	v.Number().Set(1)
	v.InsertIntoParent()
	nm := d.NewMessage()
	nm.Label().Set("groogle")
	nm.InsertIntoParent()
	m := d.Messages()[0]
	f := m.NewField()
	f.Label().Set("someField")
	f.Number().Set(1)
	f.Type().Set(e)
	f.Repeated().Set(true)
	f.Type().Set(m)
	f.Repeated().Set(false)
	f.Deprecated().Set(true)
	f.InsertIntoParent()
	mm := m.NewMap()
	mm.Label().Set("someMap")
	mm.Number().Set(4)
	mm.KeyType().Set(Int64)
	mm.Type().Set(e)
	mm.InsertIntoParent()
	o := m.NewOneOf()
	o.Label().Set("oneYeah")
	of := o.NewField()
	of.Label().Set("boing")
	of.Number().Set(5)
	of.Type().Set(String)
	of.InsertIntoParent()
	o.InsertIntoParent()

	i := d.NewImport()
	i.Path().Set("fooBar")
	i.InsertIntoParent()
	i = d.NewImport()
	i.Path().Set("barbooz")
	i.Public().Set(true)
	i.InsertIntoParent()

	s := d.NewService()
	s.Label().Set("someService")
	s.InsertIntoParent()
	r := s.NewRPC()
	r.Label().Set("GetFoo")
	r.Request().Set(m)
	r.Response().Set(m)
	r.Response().Stream().Set(true)
	r.InsertIntoParent()

	nm = m.NewMessage()
	nm.Label().Set("nestedMessage")
	nm.InsertIntoParent()
	m2 := m.Messages()[0]
	f2 := m2.NewField()
	f2.Label().Set("foild")
	f2.Number().Set(1)
	f2.Type().Set(Uint64)
	f2.InsertIntoParent()

	ne = m.NewEnum()
	ne.Label().Set("nestedEnum")
	ne.InsertIntoParent()
	e2 := m.Enums()[0]
	v2 := e2.NewVariant()
	v2.Label().Set("vorient")
	v2.Number().Set(0)
	v2.InsertIntoParent()
}
