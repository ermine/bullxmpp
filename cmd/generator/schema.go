package main

type Schema struct {
	PackageName string
	Props map[string]string
	Targets []*Target
}

type Target struct {
	Name string
	Prefix string
	Space string
	Props map[string]string
	Fields []*Field
}

type Field struct {
	Name string
	Type interface{}
	EncodingRule *Encoding
	DefaultValue string
	Required bool
	Reciver_type string
}

type Encoding struct {
	Type string
	Space string
	Name string
}
// startelement
// attribute
// element
// element:name
// name
// element:bool
// element:cdata
// cdata

type Sequence []*Field
type SequenceOf string
type Choice []*Field
type Set []*Field
type Enum []string
type Extension struct {
	Space string
	Local string
}
