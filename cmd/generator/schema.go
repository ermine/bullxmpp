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
	Parent *Field
	Type interface{}
	EncodingRule *Encoding
	DefaultValue string
	Required bool
	Annotations []*Annotation
}

type Annotation struct {
	Name string
	Params []string
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

func isForClient(field *Field) bool {
	if field.Annotations == nil {
		return false
	}
	for _, a := range field.Annotations {
		if a.Name == "decode" {
			if a.Params == nil || len(a.Params) == 0 {
				return true
			}
			for  _, p := range a.Params {
				if p == "client" {
					return true
				}
			}
			return false
		}
	}
	return false
}

func isForServer(field *Field) bool {
	if field.Annotations == nil {
		return false
	}
	for _, a := range field.Annotations {
		if a.Name == "decode" {
			if a.Params == nil || len(a.Params) == 0 {
				return true
			}
			for  _, p := range a.Params {
				if p == "server" {
					return true
				}
			}
			return false
		}
	}
	return false
}
