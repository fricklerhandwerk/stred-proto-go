package core

import (
	"fmt"
	"strings"
)

// Printer back-end for document items. The printer must ensure that `protoc`
// can compile a document.
type Printer interface {
	Document(*Document) string
	Package(*Package) string
	Import(*Import) string
	Service(*Service) string
	RPC(*RPC) string
	Message(Message) string
	Field(*Field) string
	Map(*Map) string
	OneOf(*OneOf) string
	OneOfField(*OneOfField) string
	Enum(Enum) string
	Variant(*Variant) string
	ReservedNumber(*ReservedNumber) string
	ReservedRange(*ReservedRange) string
	ReservedLabel(*ReservedLabel) string
	Label(*Label) string
	Number(*Number) string
	Type(*Type) string
	KeyType(*KeyType) string
}

type Print struct {
	Indent string
	Blank  string
}

var DefaultPrinter = Print{
	Indent: "  ",
	Blank:  "█",
}

func (p Print) Document(d *Document) string {
	numItems := 1 + len(d.Imports()) + len(d.Services()) + len(d.Messages()) + len(d.Enums())
	items := make([]string, 0, numItems)
	if d._package.label.value != "" {
		items = append(items, d._package.String())
	}
	for _, i := range d.Imports() {
		items = append(items, i.String())
	}
	for _, s := range d.Services() {
		items = append(items, s.String())
	}
	for _, e := range d.Enums() {
		items = append(items, e.String())
	}
	for _, m := range d.Messages() {
		items = append(items, m.String())
	}
	for i, item := range items {
		items[i] = fmt.Sprint("\n\n", item)
	}
	return fmt.Sprint("syntax = \"proto3\";", strings.Join(items, ""))
}

func (p Print) Package(pkg *Package) string {
	return fmt.Sprintf("package %s;", pkg.label)
}

func (p Print) Import(*Import) string         { return "" }
func (p Print) Service(*Service) string       { return "" }
func (p Print) RPC(*RPC) string               { return "" }
func (p Print) Message(Message) string        { return "" }
func (p Print) Field(*Field) string           { return "" }
func (p Print) Map(*Map) string               { return "" }
func (p Print) OneOf(*OneOf) string           { return "" }
func (p Print) OneOfField(*OneOfField) string { return "" }

func (p Print) Enum(e Enum) string {
	items := make([]string, 0, 1+len(e.Fields()))
	// `protoc` does not allow "unnecessary" declaration of `allow_alias = true`
	// when there is no aliasing in place.
	if e.AllowAlias().value && aliased(e) {
		items = append(items, "option allow_alias = true;")
	}
	for _, f := range e.Fields() {
		items = append(items, f.String())
	}
	for i, item := range items {
		items[i] = fmt.Sprintln(p.Indent, item)
	}

	return fmt.Sprintf("enum %s {\n%s}", e.Label(), strings.Join(items, ""))
}

func aliased(e Enum) bool {
	for _, a := range e.Aliases() {
		if len(a) > 1 {
			return true
		}
	}
	return false
}

func (p Print) Variant(v *Variant) string {
	var deprecated string
	if v.deprecated.value {
		deprecated = " [deprecated=true]"
	}
	return fmt.Sprintf("%s = %s%s;", v.label, v.number, deprecated)
}

func (p Print) ReservedNumber(n *ReservedNumber) string {
	return fmt.Sprintf("reserved %s;", n.number)
}

func (p Print) ReservedRange(r *ReservedRange) string {
	return fmt.Sprintf("reserved %s to %s;", r.start, r.end)
}

func (p Print) ReservedLabel(l *ReservedLabel) string {
	return fmt.Sprintf("reserved %s;", l.label)
}

func (p Print) Label(l *Label) string {
	if l.value != "" {
		return l.value
	}
	return p.Blank
}

func (p Print) Number(n *Number) string {
	if n.value != nil {
		return fmt.Sprint(*n.value)
	}
	return p.Blank
}

func (p Print) Type(*Type) string       { return "" }
func (p Print) KeyType(*KeyType) string { return "" }
