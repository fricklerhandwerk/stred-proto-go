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
	numItems := 1 + len(d.imports) + len(d.services) + len(d.messages) + len(d.enums)
	items := make([]string, 0, numItems)
	if d._package.label.value != "" {
		items = append(items, d._package.String())
	}

	if len(d.imports) > 0 {
		imports := make([]string, 0, len(d.imports))
		for i := range d.imports {
			imports = append(imports, i.String())
		}
		items = append(items, strings.Join(imports, "\n"))
	}

	for s := range d.services {
		items = append(items, s.String())
	}
	for e := range d.enums {
		items = append(items, e.String())
	}
	for m := range d.messages {
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

func (p Print) Import(i *Import) string {
	var public string
	if i.public.value {
		public = "public "
	}
	return fmt.Sprintf("import %s%s;", public, i.path)
}

func (p Print) Service(s *Service) string {
	rpcs := make([]string, 0, len(s.rpcs))
	for r := range s.rpcs {
		rpcs = append(rpcs, r.String())
	}
	return fmt.Sprintf("service %s {\n%s\n}", s.label, p.indent(strings.Join(rpcs, "\n")))
}

func (p Print) RPC(r *RPC) string {
	return fmt.Sprintf("rpc %s (%s) returns (%s);", r.label, r.request.value.Label(), r.response.value.Label())
}

func (p Print) Message(m Message) string {
	items := make([]string, 0, 3)
	if len(m.Fields()) > 0 {
		items = append(items, p.messageFields(m))
	}
	if len(m.Enums()) > 0 {
		items = append(items, p.messageEnums(m))
	}
	if len(m.Messages()) > 0 {
		items = append(items, p.messageMessages(m))
	}
	return fmt.Sprintf("message %s {\n%s\n}", m.Label(), p.indent(strings.Join(items, "\n\n")))
}

func (p Print) messageFields(m Message) string {
	fields := make([]string, 0, len(m.Fields()))
	for _, f := range m.Fields() {
		fields = append(fields, f.String())
	}
	return strings.Join(fields, "\n")
}

func (p Print) messageEnums(m Message) string {
	enums := make([]string, 0, len(m.Enums()))
	for _, e := range m.Enums() {
		enums = append(enums, e.String())
	}
	return strings.Join(enums, "\n")
}

func (p Print) messageMessages(m Message) string {
	messages := make([]string, 0, len(m.Messages()))
	for _, n := range m.Messages() {
		messages = append(messages, n.String())
	}
	return strings.Join(messages, "\n")
}

func (p Print) Field(f *Field) string {
	var repeated string
	if f.repeated.value {
		repeated = "repeated "
	}
	return fmt.Sprintf("%s%s %s = %s%s;", repeated, f._type, f.label, f.number, deprecated(f.deprecated))
}

func (p Print) Map(m *Map) string {
	return fmt.Sprintf("map <%s,%s> %s = %s%s;", m.keyType, m._type, m.label, m.number, deprecated(m.deprecated))
}

func (p Print) OneOf(o *OneOf) string {
	items := make([]string, 0, len(o.fields))
	for f := range o.fields {
		items = append(items, f.String())
	}
	return fmt.Sprintf("oneof %s {\n%s\n}", o.Label(), p.indent(strings.Join(items, "\n")))
}

func (p Print) OneOfField(f *OneOfField) string {
	var deprecated string
	if f.deprecated.value {
		deprecated = " [deprecated=true]"
	}
	return fmt.Sprintf("%s %s = %s%s;", f._type, f.label, f.number, deprecated)
}

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

	return fmt.Sprintf("enum %s {\n%s\n}", e.Label(), p.indent(strings.Join(items, "\n")))
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
	return fmt.Sprintf("%s = %s%s;", v.label, v.number, deprecated(v.deprecated))
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
	if l.value == "" {
		return p.Blank
	}
	return l.value
}

func (p Print) Number(n *Number) string {
	if n.value == nil {
		return p.Blank
	}
	return fmt.Sprint(*n.value)
}

func (p Print) Type(t *Type) string {
	if t.value == nil {
		return p.Blank
	}
	switch v := t.value.(type) {
	case Message:
		return v.Label().String()
	case Enum:
		return v.Label().String()
	default:
		return fmt.Sprint(v)
	}
}

func (p Print) KeyType(k *KeyType) string {
	if k.value == nil {
		return p.Blank
	}
	return fmt.Sprint(k.value)
}

func (p Print) indent(in string) string {
	lines := strings.Split(in, "\n")
	for i, l := range lines {
		indent := ""
		if l != "" {
			indent = p.Indent
		}
		lines[i] = fmt.Sprint(indent, l)
	}
	return strings.Join(lines, "\n")
}

func deprecated(f Flag) string {
	if f.value {
		return " [deprecated=true]"
	}
	return ""
}
