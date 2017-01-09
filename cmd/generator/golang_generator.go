package main

import (
	"os"
	"errors"
	"path/filepath"
	"fmt"
	"unicode"	
)

const	ns_xml = "http://www.w3.org/XML/1998/namespace"

func GolangGenerate() error {
	if cfg.Golang.Outdir == "" {
		panic("no outdir")
	}

	for _, schema := range schemas {
		dir := cfg.Golang.Outdir
		if category, ok := schema.Props["category"]; ok {
			if category == "extension" {
				dir = filepath.Join(cfg.Golang.Outdir, "extensions")
			}
		}
		fdir := filepath.Join(dir, schema.PackageName)
		fmt.Println("making directory " , fdir)
		os.MkdirAll(fdir, 0755)
		for _, target := range schema.Targets {
			if name, ok := target.Props["name"]; ok {
				target.Name = name
			}
			if target.Name != "" {
				subdir := filepath.Join(fdir, target.Name)
				fmt.Println("Creating subdir ", subdir)
				os.MkdirAll(subdir, 0755)
			}
		}
		golang_generate_package(fdir, schema)
	}
	return nil
}

var mapSimpleTypes = map[string]struct{
	Type string
	Import string
}{
	"jid": {"jid.JID", "github.com/ermine/bullxmpp/jid"},
	"int": {"int", "strconv"},
	"uint": {"uint", "strconv"},
	"string": {"string", ""},
	"boolean": {"bool", ""},
	"langstring": {"xmlencoder.LangString", ""},
	"xmllang": {"string", ""},
	"bytestring": {"[]byte", ""},
	"extension": {"interface{}", ""},
	"datetime": {"time.Time", "time"},
}

func golang_simpletype(s string) string {
	if t, ok := mapSimpleTypes[s]; ok {
		return t.Type
	}
	return s
}

func golang_generate_package(dir string, schema *Schema) error {
	for _, target := range schema.Targets {
		currdir := dir
		var filename string
		if target.Name != "" {
			currdir = filepath.Join(currdir, target.Name)
			filename = filepath.Join(currdir, target.Name + "_data.go")
		} else {
			filename = filepath.Join(currdir, schema.PackageName + "_data.go")
		}			
		fmt.Println("filename " + filename)
		file, err := os.Create(filename)
		if err != nil { return err }
		defer file.Close()
		imports := golang_getImports_schema(schema)
		if target.Name != "" {
			file.WriteString("package " + target.Name + "\n")
		} else {
			file.WriteString("package " + schema.PackageName + "\n")
		}
		file.WriteString("\n")
		file.WriteString("import \"encoding/xml\"\n")
		file.WriteString("import \"github.com/ermine/bullxmpp/xmlencoder\"\n")
		for _, i := range imports {
			file.WriteString("import \"" + i + "\"\n")
		}
		file.WriteString("const NS = \"" + target.Space + "\"\n")
		if err = golang_generate_structs(file, target); err != nil { return err }
		if err = golang_generate_encoders(file, target); err != nil { return err }
		golang_generate_init(file, target)
	}
	return nil
}

func golang_getImports_schema(schema *Schema) []string {
	var imports []string
	for _, x := range schema.Targets {
		golang_getImports(x.Fields, &imports)
	}
	return imports
}

func golang_getImports(fields []*Field, imports *[]string) {
	for _, x := range fields {
		switch typ := x.Type.(type) {
		case Set:
			fields = []*Field(typ)
			golang_getImports(fields, imports)
		case Sequence:
			fields = []*Field(typ)
			golang_getImports(fields, imports)
		case Choice:
			fields = []*Field(typ)
			golang_getImports(fields, imports)
		case Enum:
		case Extension:
			if typ.Local != "" {
			Added:
				for _, schema := range schemas {
					for _, target := range schema.Targets {
						if target.Space == typ.Space {
							if target.Name != "" {
								append_import(imports, "github.com/ermine/bullxmpp/" + schema.PackageName + "/" + target.Name)
							} else {
								append_import(imports, "github.com/ermine/bullxmpp/" + schema.PackageName)
							}
							break Added
						}
					}
				}
			}
		case SequenceOf:
		case string:
			if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" && typ == "boolean" {
				append_import(imports, "strconv")
			} else if t, ok := mapSimpleTypes[typ]; ok {
				if t.Import != "" {
					append_import(imports, t.Import)
				}
			}
		}
	}
}

func append_import(imports *[]string, i string) {
	found := false
Found:
	for _, j := range *imports {
		if j == i {
			found = true
			break Found
		}
	}
	if !found {
		*imports = append(*imports, i)
	}
}

func golang_referenceType(target *Target, field *Field) string {
	name := field.Name
	if name == "" {
		name = field.Type.(string)
	}
	for _, x := range target.Fields {
		if x.Name == name {
			return golang_makeIdent(name)
		}
	}
	return golang_simpletype(name)
}

func golang_generate_structs(file *os.File, target *Target) error {
	var err error
	var enums []*Field
	golang_collect_enums(target, target.Fields, &enums, "")
	target.Fields = append(target.Fields, enums...)
	for _, def := range target.Fields {
		if _, ok := def.Type.(Enum); ok {
			continue
		}
		file.WriteString("type " + golang_makeIdent(def.Name))
		switch t := def.Type.(type) {
		case Extension:
			if t.Local != "" {
				field, err := golang_getExternalType(t.Space, t.Local)
				if err != nil { return err }
				file.WriteString(" " + field + "\n")
			} else {
				file.WriteString(" interface{}\n")
			}
		case SequenceOf:
			if t == "extension" {
				file.WriteString(" []interface{}\n")
			} else {
				field := golang_getFieldByName(target, string(t))
				file.WriteString(" []" + golang_referenceType(target, field) + "\n")
			}
		case string:
			var typ string
			if z, ok := mapSimpleTypes[t]; ok {
				typ = z.Type
			} else {
				field := golang_getFieldByName(target, z.Type)
				if field == nil {
					return errors.New("unknown type " + t)
				}
				typ = golang_referenceType(target, field)
			}
			file.WriteString(" " + typ + "\n")
		case Set:
			file.WriteString(" struct {\n")
			if err = golang_generate_fields(file, target, []*Field(t)); err != nil { return err }
			file.WriteString("}\n")
		case Sequence:
			file.WriteString(" []interface{}\n}")
		case Choice:
			def.Type = Set([]*Field{&Field{Name:"Payload", Type:t}})
			file.WriteString(" struct {\n")
			file.WriteString(" Payload interface{}\n")
			file.WriteString("}\n")
		}
	}
	golang_generate_enums(file, target, enums)
	return nil
}

func golang_generate_enums(file *os.File, target *Target, enums []*Field) {
	for _, x := range enums {
		typ := golang_makeIdent(x.Name)
		file.WriteString("type " + typ + " string\n")
		enum := []string(x.Type.(Enum))
		file.WriteString("const (\n")
		for _, z := range enum {
			file.WriteString(target.Name + x.Name + golang_makeIdent(z) + " " + typ + " = \"" + z + "\"\n")
		}
		file.WriteString(")\n")
	}
}

func golang_collect_enums(target *Target, fields []*Field, enums *[]*Field, prefix string) {
	for _, x := range fields {
		switch t := x.Type.(type) {
		case Set:
			fields := []*Field(t)
			golang_collect_enums(target, fields, enums, prefix + golang_makeIdent(x.Name))
		case Sequence:
			fields := []*Field(t)
			golang_collect_enums(target, fields, enums, prefix + golang_makeIdent(x.Name))
		case Choice:
			fields := []*Field(t)
			golang_collect_enums(target, fields, enums, prefix + golang_makeIdent(x.Name))
			golang_checkTypes(target, fields)
		case Enum:
			field := &Field {
				Name:x.Name,
				Type: x.Type,
				EncodingRule:x.EncodingRule,
				DefaultValue:x.DefaultValue,
				Required:x.Required,
			}
			if prefix != "" {
				typename := prefix + golang_makeIdent(x.Name)
				x.Type = typename
				field.Name = typename
			}
			*enums = append(*enums, field)
		}
	}
	// return enums
}

func golang_generate_fields(file *os.File, target *Target, fields []*Field) error {
	for _, x := range fields {
		if x.Name == "" {
			x.Name = x.Type.(string)
		}
		file.WriteString("  " + golang_makeIdent(x.Name))
		switch t := x.Type.(type) {
		case string:
			switch t {
				case "bytestring":
				file.WriteString(" []byte")
			case "boolean":
				file.WriteString(" bool")
			default:
				var typ string
				if z, ok := mapSimpleTypes[t]; ok {
					typ = z.Type
				} else {
					field := golang_getFieldByName(target, t)
					if field == nil {
						return errors.New("unknown type " + t)
					}
					typ = golang_referenceType(target, field)
				}
				if x.EncodingRule != nil && x.EncodingRule.Type == "element:bool" {
					file.WriteString(" " + typ)
				} else {
					file.WriteString(" *" + typ)
				}
			}
		case SequenceOf:
			switch t {
			case "extension":
				file.WriteString(" []interface{}")
			default:
				field := golang_getFieldByName(target, string(t))
				if field == nil {
					file.WriteString(" []" + golang_simpletype(string(t)))
				} else {
					if golang_isStruct(field) {
						file.WriteString(" []*" + golang_referenceType(target, field))
					} else {
						file.WriteString(" []" + golang_referenceType(target, field))
					}
				}
			}
		case Extension:
			if t.Local != "" {
				fieldtype, err := golang_getExternalType(t.Space, t.Local)
				if err != nil { return err }
				file.WriteString(" *" + fieldtype +  "\n")
			} else {
				file.WriteString(" interface{}\n")
			}
		case Sequence:
			file.WriteString(" []interface{}")
		case Choice:
			file.WriteString(" interface{}")
		case Set:
			fields := []*Field(t)
			file.WriteString(" struct {\n")
			golang_generate_fields(file, target, fields)
			file.WriteString("}\n")
		default:
			fmt.Println("default1: ", x.EncodingRule)
		}
		file.WriteString("\n")
	}
	return nil
}

func golang_checkTypes(target *Target, fields []*Field) {
	for _, x := range fields {
		found := false
		for _, z := range target.Fields {
			if z.Name == x.Name {
				found = true
				break
			}
		}
		if !found && x.EncodingRule != nil {
			target.Fields = append(target.Fields, x)
		}
	}
}

func golang_generate_encoders(file *os.File, target *Target) error {
	var err error
	for _, x := range target.Fields {
		if _, ok := x.Type.(Enum); ok {
			continue
		}
		if err = golang_generate_encoder(file, target, x); err != nil { return err }
		if err = golang_generate_decoder(file, target, x); err != nil { return err }
	}
	return nil
}

func golang_getSpaceAndName(target *Target, targetNS string, field *Field) (string, string) {
	if s, ok := field.Type.(string); ok {
		if s == "xmllang" {
			return "\"" +  ns_xml + "\"" , "lang"
		}
	}
	var space, local string
	if field.Name == "" {
		field1 := golang_getFieldByName(target, field.Type.(string))
		if field1 == nil {
			fmt.Println("Cannot find field for ", field.Type)
		}
		field = field1
	}
	local = field.Name
	if field.EncodingRule != nil && field.EncodingRule.Name != "" {
		local = field.EncodingRule.Name
	}
	space = targetNS
	if field.EncodingRule != nil && field.EncodingRule.Space != "" {
		space = field.EncodingRule.Space
	}
	if space != "" && space == targetNS {
		space = "NS"
	} else {
		space = "\"" + space + "\""
	}
	// local = "\"" + local + "\""
	return space, local
}

func golang_getElementName(fields []*Field) *Field {
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "name" {
			return x
		}
	}
	return nil
}

func golang_generate_encoder(file *os.File, target *Target, element *Field) error {
	file.WriteString("func (elm *" + golang_makeIdent(element.Name) + ") Encode(e *xmlencoder.Encoder) error {\n")
	file.WriteString("var err error\n")
	if target.Prefix != "" {
		file.WriteString("if err = e.SetPrefix(\"" + target.Prefix +
			"\", \"" + target.Space + "\"); err != nil { return err }\n")
	}
	golang_generate_element_encoder(file, target, "elm", element)
	file.WriteString("return nil\n")
	file.WriteString("}\n\n")
	return nil
}

func golang_generate_element_encoder(file *os.File, target *Target, prefix string,
	element *Field) error {
	if element.EncodingRule == nil {
		switch typ := element.Type.(type) {
		case Extension:
			if typ.Local == "" {
				file.WriteString("if err = " + prefix + ".(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			} else {
				file.WriteString("if err = " + prefix + ".Encode(e); err != nil { return err }\n")
			}
		case string:
			file.WriteString("if err = " + prefix + ".Encode(e); err != nil { return err }\n")
		case SequenceOf:
			file.WriteString("for _, x := range " + prefix + " {\n")
			if string(typ) == "extension" {
				file.WriteString("if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err} \n")
			} else {
				file.WriteString("if err = x.Encode(e); err != nil { return err} \n")
			}
			file.WriteString("}\n")
		case Sequence:
			file.WriteString("for _, x := range " + prefix + " {\n")
			file.WriteString("if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			file.WriteString("}\n")
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				close := golang_generate_check_condition(file, prefix + "." + golang_makeIdent(x.Name), x)
				golang_generate_element_encoder(file, target, prefix + "." + golang_makeIdent(x.Name), x)
				if close {
					file.WriteString("}\n")
				}
			}
		case Choice:
			file.WriteString("if " + prefix + " != nil {\n")
			file.WriteString("if err = " + prefix + ".(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			file.WriteString("}\n")
		default:
			fmt.Println("dont know what to do ", element.Name, " ", element.Type)
		}
		return nil
	}
	fieldname := golang_makeIdent(element.Name)
	space, local := golang_getSpaceAndName(target, target.Space, element)
	switch element.EncodingRule.Type {
	case "element:name":
		file.WriteString("if err = e.StartElement(" + space + ", string(*" + prefix + ")); err != nil { return err }\n")
		file.WriteString("if err = e.EndElement(); err != nil { return err }\n")
	case "element:cdata":
		varname := prefix
		isarray := false
		if _, ok := element.Type.(SequenceOf); ok {
			varname = "x"
			isarray = true
		}
		value := golang_generate_simplevalue_encoder(varname, element)
		if isarray {
			file.WriteString("for _, x := range " + prefix + " {\n")
		}
		file.WriteString("if err = e.SimpleElement(" + space + ", \"" + local + "\", " +
			value + "); err != nil { return err }\n")
		if isarray {
			file.WriteString("}\n")
		}
	case "element:bool":
		file.WriteString("if err = e.StartElement(" + space + ", \"" +
			local + "\"); err != nil { return err }\n")
		file.WriteString("if err = e.EndElement(); err != nil { return err }\n")
	case "startelement", "element":
		switch typ := element.Type.(type) {
		case Extension:
			space, local := golang_getSpaceAndName(target, target.Space, element)
			file.WriteString("if err = e.StartElement(" + space + ", \"" + local + "\"); err != nil { return err }\n")
			file.WriteString("if err = " + prefix + ".Encode(e); err != nil { return err }\n")
			file.WriteString("if err = e.EndElement(); err != nil { return err }\n")
		case Set:
			fields := []*Field(typ)
			specialFieldName := golang_getElementName(fields)
			if specialFieldName != nil {
				fieldname = golang_makeIdent(specialFieldName.Name)
				file.WriteString("if err = e.StartElement(" + space + ", string(*" +
					prefix + "." + fieldname + ")); err != nil { return err }\n")
			} else {
				file.WriteString("if err = e.StartElement(" + space +
					", \"" + local + "\"); err != nil { return err }\n")
			}
			for _, x := range fields {
				if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
					aspace, alocal := golang_getSpaceAndName(target, "", x)
					close := golang_generate_check_condition(file, prefix + "." + golang_makeIdent(x.Name), x)
					value := golang_generate_simplevalue_encoder(prefix + "." + golang_makeIdent(x.Name), x)
					file.WriteString("if err = e.Attribute(" + aspace + ", \"" +
						alocal + "\", " + value + "); err != nil { return err }\n")
					if close {
						file.WriteString("}\n")
					}
				}
			}
			for _, x := range fields {
				if x.EncodingRule != nil &&
					(x.EncodingRule.Type == "cdata" || x.EncodingRule.Type == "attribute") {
					continue
				}
				close := golang_generate_check_condition(file, prefix + "." + golang_makeIdent(x.Name), x)
				golang_generate_element_encoder(file, target, prefix + "." + golang_makeIdent(x.Name), x)
				if close {
					file.WriteString("}\n")
				}
			}
			for _, x := range fields {
				if x.EncodingRule != nil && x.EncodingRule.Type == "cdata" {
					file.WriteString("if " + prefix + "." + golang_makeIdent(x.Name) + " != nil {\n")
					typ := x.Type.(string)
					if typ == "bytestring" {
						file.WriteString("if err = e.Bytes(" + prefix + "." +
							golang_makeIdent(x.Name) + "); err != nil { return err }\n")
					} else {
						value := golang_generate_simplevalue_encoder(prefix + "." + golang_makeIdent(x.Name), x)
						file.WriteString("if err = e.Text(" + value + "); err != nil { return err }\n")
					}
					file.WriteString("}\n")
				}
			}
			if element.EncodingRule.Type == "element" {
				file.WriteString("if err = e.EndElement(); err != nil { return err }\n")
			}
		case string:
			switch typ {
			case "langstring":
				file.WriteString(prefix + ".Encode(e, " + space + ", \"" + local + "\")\n")
			case "extension":
				file.WriteString(prefix + ".(xmlencoder.Extension).Encode(e)\n")
			default:
				file.WriteString(prefix + ".Encode(e)\n")
			}
		case Choice:
			file.WriteString("if err = e.StartElement(" + space + ", \"" + local + "\"); err != nil { return err }\n")
			if prefix == "elm" {
				file.WriteString("if err = elm.Payload.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			} else {
				file.WriteString("if err = " + prefix + ".(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			}
			file.WriteString("if err = e.EndElement(); err != nil { return err }\n")
		case SequenceOf:
			file.WriteString("if err = e.StartElement(" + space + ", \"" + local + "\"); err != nil { return err }\n")
			if prefix == "elm" {
				file.WriteString("for _, x := range *elm {\n")
			} else {
				file.WriteString("for _, x := range " + prefix + " {\n")
			}
			if typ == "extension" {
				file.WriteString("if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			} else {
				file.WriteString("if err = x.Encode(e); err != nil { return err }\n")
			}
			file.WriteString("}\n")
			file.WriteString("if err = e.EndElement(); err != nil { return err }\n")
		case Sequence:
			file.WriteString("if err = e.StartElement(" + space + ", \"" + local + "\"); err != nil { return err }\n")
			file.WriteString("for _, x := range *" + prefix + " {\n")
			file.WriteString("if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			file.WriteString("}\n")
			file.WriteString("if err = e.EndElement(); err != nil { return err }\n")
		}
	}
	return nil
}

func golang_hasChilds(target *Target, field *Field) bool {
	if field.EncodingRule != nil {
		switch field.EncodingRule.Type {
		case "element:cdata", "cdata", "name", "element:bool", "element:name", "attribute":
			return false
		case "element":
			switch typ := field.Type.(type) {
			case string:
				field := golang_getFieldByName(target, string(typ))
				if field != nil { return golang_hasReallyChilds(target, field) }
			case Set:
				fields := []*Field(typ)
				for _, x := range fields {
					if golang_hasReallyChilds(target, x) { return true }
				}
			case SequenceOf:
				field := golang_getFieldByName(target, string(typ))
				if field != nil { return golang_hasReallyChilds(target, field) }
			case Sequence:
				fields := []*Field(typ)
				for _, x := range fields {
					if golang_hasReallyChilds(target, x) { return true }
				}
			case Choice:
				fields := []*Field(typ)
				for _, x := range fields {
					if golang_hasReallyChilds(target, x) { return true }
				}
			case Extension:
				return false
			}
		}
	}
	return false
}

func golang_hasReallyChilds(target *Target, field *Field) bool {
	if field.EncodingRule != nil {
		switch field.EncodingRule.Type {
		case "element:cdata", "element:bool", "element:name", "element":
			if field.EncodingRule.Space == "" || field.EncodingRule.Space == target.Space {
				return true
			}
		}
	} else {
		switch typ := field.Type.(type) {
		case string:
			field := golang_getFieldByName(target, typ)
			if field != nil { return golang_hasReallyChilds(target, field)}
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				if golang_hasReallyChilds(target, x) { return true }
			}
		case SequenceOf:
			field := golang_getFieldByName(target, string(typ))
			if field != nil {
				return golang_hasReallyChilds(target, field)
			}
		case Sequence:
			fields := []*Field(typ)
			for _, x := range fields {
				if golang_hasReallyChilds(target, x) { return true }
			}			
		case Choice:
			fields := []*Field(typ)
			for _, x := range fields {
				if golang_hasReallyChilds(target, x) { return true }
			}			
		case Extension:
			return false
		}
	}
	return false
	
}

func golang_generate_decoder(file *os.File, target *Target, element *Field) error {
	file.WriteString("func (elm *" + golang_makeIdent(element.Name) + ") Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {\n")
	file.WriteString("var err error\n")
	/*
	if golang_hasChilds(target, element) {
		file.WriteString("ns := \"" + target.Space + "\"\n")
	}
*/
	golang_generate_element_decoder(file, target, "elm", element)
	file.WriteString("return err\n")
	file.WriteString("}\n\n")
	return nil
}

func golang_generate_element_decoder(file *os.File, target *Target, prefix string, element *Field) error {
	// var err error
	if element.EncodingRule == nil {
		fmt.Println("dont know how to generate decoder ", element.Name)
		return nil
	}
	switch element.EncodingRule.Type {
	case "element:cdata":
		golang_generate_cdata_decoder(file, target, prefix, element)
	case "element:bool":
		file.WriteString("if err = d.Skip(); err != nil { return err }\n")
	case "startelement", "element":
		switch typ := element.Type.(type) {
		case string:
			return errors.New("dont know what todo with " + typ)
		case Extension:
			file.WriteString("var t xml.Token\n")
			file.WriteString("Loop:\n")
			file.WriteString("for {\n")
			file.WriteString("if t, err = d.Token(); err != nil { return err }\n")
			file.WriteString("switch t := t.(type) {\n")
			file.WriteString("case xml.StartElement:\n")
			if typ.Local != "" {
				fieldname, err := golang_getExternalType(typ.Space, typ.Local)
				if err != nil { return err }
				file.WriteString("if t.Name.Space == \"" + typ.Space + "\" && t.Name.Local == \"" +
					typ.Local + "\"{\n")
				file.WriteString("newel := &" + fieldname + "{}\n")
			} else {
				file.WriteString("if newel, ok := xmlencoder.GetExtension(t.Name); ok {\n")
			}
			file.WriteString("if err = newel.Decode(d, &t); err != nil { return err }\n")
			if prefix == "elm" {
				ctype := golang_referenceType(target, element)
				file.WriteString("*" + prefix + " = " + ctype + "(*newel)\n")
			} else {
				file.WriteString("*" + prefix + " = newel\n")
			}
			file.WriteString("} else {\n")
			file.WriteString("if err = d.Skip(); err != nil { return err }\n")
			file.WriteString("}\n")
			file.WriteString("case xml.EndElement:\n")
			file.WriteString("break Loop\n")
			file.WriteString("}\n")
			file.WriteString("}\n")

		case Set:
			fields := []*Field(typ)
			if len(fields) == 0 {
				file.WriteString("if err = d.Skip(); err != nil { return err }\n")
			} else {
				golang_generate_element_set_decoder(file, target, prefix, fields)
			}
		case Choice:
			fields := []*Field(typ)
			file.WriteString("var t xml.Token\n")
			file.WriteString("Loop:\n")
			file.WriteString("for {\n")
			file.WriteString("if t, err = d.Token(); err != nil { return err }\n")
			file.WriteString("switch t := t.(type) {\n")
			file.WriteString("case xml.StartElement:\n")
			file.WriteString("switch {\n")
			for _, z := range fields {
				space, local := golang_getSpaceAndName(target, target.Space, z)
				file.WriteString("case t.Name.Space == " + space + " && t.Name.Local == \"" +
					local + "\":\n")
				file.WriteString("newel := &" + golang_referenceType(target, z) + "{}\n")
				file.WriteString("if err = newel.Decode(d, &t); err != nil { return err}\n")
				file.WriteString("*" + prefix + " = newel\n")
			}
			file.WriteString("}\n")
			file.WriteString("case xml.EndElement:\n")
			file.WriteString("break Loop\n")
			file.WriteString("}\n")
			file.WriteString("}\n")
			
		case SequenceOf:
			field := string(typ)
			file.WriteString("var t xml.Token\n")
			file.WriteString("Loop:\n")
			file.WriteString("for {\n")
			file.WriteString("if t, err = d.Token(); err != nil { return err }\n")
			file.WriteString("switch t := t.(type) {\n")
			file.WriteString("case xml.StartElement:\n")
			if field == "extension" {
				file.WriteString("if newel, ok := xmlencoder.GetExtension(t.Name); ok {\n")
				file.WriteString("if err = newel.(xmlencoder.Extension).Decode(d, &t); err != nil { return err }\n")
				file.WriteString("*" + prefix + " = append(*" + prefix + ", newel)\n")
				file.WriteString("} else {\n")
				file.WriteString("if err = d.Skip(); err != nil { return err }\n")
				file.WriteString("}\n")
			} else {
				f := golang_getFieldByName(target, field)
				if f == nil {
					// import from other packages?
					return errors.New("dont know what to do with " + field)
				}
				space, local := golang_getSpaceAndName(target, target.Space, f)
				file.WriteString("if t.Name.Space == " + space + " && t.Name.Local == \"" + local + "\" {\n")
				if golang_isStruct(f) {
					file.WriteString("newel := &" + golang_referenceType(target, f) + "{}\n")
				} else {
					file.WriteString("var newel " + golang_referenceType(target, f) + "\n")
				}					
				file.WriteString("if err = newel.Decode(d, &t); err != nil { return err }\n")
				if golang_isStruct(f) {
					file.WriteString(prefix + " = append(" + prefix + ", newel)\n")
				} else {
					file.WriteString("*" + prefix + " = append(*" + prefix + ", newel)\n")
				}
				file.WriteString("}\n")
			}
			file.WriteString("case xml.EndElement:\n")
			file.WriteString("break Loop\n")
			file.WriteString("}\n")
			file.WriteString("}\n")
		}
	}
	return nil
}

func golang_getFieldByName(target *Target, f string) *Field {
	for _, x := range target.Fields {
		if x.Name == f {
			return x
		}
	}
	return nil
}

func golang_generate_element_set_decoder(file *os.File, target *Target, prefix string, fields []*Field) error {
	var err error
	var attrs []*Field
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
			attrs = append(attrs, x)
		}
	}
	if len(attrs) > 0 {
		file.WriteString("for _, x := range tag.Attr {\n")
		file.WriteString("switch {\n")
		for _, x := range attrs {
			space, local := golang_getSpaceAndName(target, "", x)
			file.WriteString("case x.Name.Space == " + space + 
				" && x.Name.Local == \"" + local + "\":\n")
			golang_generate_simplevalue_decoder(file, target, prefix + "." + golang_makeIdent(x.Name), "x.Value", x)
		}
		file.WriteString("}\n")
		file.WriteString("}\n")
	}
	var elems []*Field
	var any *Field
Loop:
	for _, x := range fields {
		if x.EncodingRule == nil {
			elems = append(elems, x)
		} else {
			switch x.EncodingRule.Type {
			case "element:name":
				any = x
			case "element":
				if set, ok := x.Type.(Set); ok {
					fields := []*Field(set)
					for _, z := range fields {
						if z.EncodingRule.Type == "name" {
							any = x
							continue Loop
						}
					}
				}
				elems = append(elems, x)
			case "element:cdata", "", "element:bool":
				elems = append(elems, x)
			}
		}
	}
	if len(elems) > 0 || any != nil {
		file.WriteString("var t xml.Token\n")
		file.WriteString("Loop:\n")
		file.WriteString("for {\n")
		file.WriteString("if t, err = d.Token(); err != nil { return err }\n")
		file.WriteString("switch t := t.(type) {\n")
		file.WriteString("case xml.EndElement:\n")
		file.WriteString("break Loop\n")
		file.WriteString("case xml.StartElement:\n")
		file.WriteString("switch {\n")
		for _, x := range elems {
			if x.EncodingRule != nil {
				space, local := golang_getSpaceAndName(target, target.Space, x)
				file.WriteString("case t.Name.Space == " + space + " && t.Name.Local == \"" +
					local + "\":\n")
				switch x.EncodingRule.Type {
				case "element:bool":
					file.WriteString(prefix + "." + golang_makeIdent(x.Name) + " = true\n")
					file.WriteString("if err = d.Skip(); err != nil { return err }\n")
					file.WriteString("continue\n")
				case "element:cdata":
					golang_generate_cdata_decoder(file, target, prefix + "." + golang_makeIdent(x.Name), x)
				case "element":
					switch typ := x.Type.(type) {
					case string:
						file.WriteString("if err = " + prefix + "." + golang_makeIdent(x.Name) +
							".Decode(d, &t); err != nil { return err }\n")
					case Set:
						fields := []*Field(typ)
						if err = golang_generate_element_set_decoder(file, target, prefix + "." + golang_makeIdent(x.Name),
							fields); err != nil { return err }
					case SequenceOf:
						field := golang_getFieldByName(target, string(typ))
						if field == nil {
							return errors.New("Cannot find field " + string(typ))
						}
						space, local := golang_getSpaceAndName(target, target.Space, field)
						file.WriteString("var t xml.Token\n")
						file.WriteString("InLoop:\n")
						file.WriteString("for {\n")
						file.WriteString("if t, err = d.Token(); err != nil { return err }\n")
						file.WriteString("switch t := t.(type) {\n")
						file.WriteString("case xml.StartElement:\n")
						file.WriteString("if t.Name.Space == " + space + " && t.Name.Local == \"" +
							local + "\" {\n")
						file.WriteString("newel := &" + golang_referenceType(target, field) + "{}\n")
						file.WriteString("if err = newel.Decode(d, &t); err != nil { return err }\n")
						file.WriteString(prefix + "." + golang_makeIdent(x.Name) + " = append(" +
							prefix + "." + golang_makeIdent(x.Name) + ", newel)\n")
						file.WriteString("}\n")
						file.WriteString("case xml.EndElement:\n")
						file.WriteString("break InLoop\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
					case Sequence:
						file.WriteString("InLoop:\n")
						file.WriteString("for {\n")
						file.WriteString("if t, err = d.Token(); err != nil { return err }\n")
						file.WriteString("switch t := t.(type) {\n")
						file.WriteString("case xml.StartElement:\n")
						file.WriteString("switch {\n")
						for _, z := range []*Field(typ) {
							space, local := golang_getSpaceAndName(target, target.Space, z)
							file.WriteString("case t.Name.Space == " + space + " && t.Name.Local == \"" +
								local + "\":\n")
							file.WriteString("newel := &" + golang_referenceType(target, z) + "{}\n")
							file.WriteString("if err = newel.Decode(d, &t); err != nil { return err }\n")
							file.WriteString(prefix + "." + golang_makeIdent(x.Name) + " = append(" +
								prefix + "." + golang_makeIdent(x.Name) + ", newel)\n")
						}
						file.WriteString("case xml.EndElement:\n")
						file.WriteString("break InLoop\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
					case Choice:
						file.WriteString("InLoop:\n")
						file.WriteString("for {\n")
						file.WriteString("switch {\n")
						file.WriteString("if t, err = d.Token(); err != nil { return err }\n")
						file.WriteString("switch t := t.(type) {\n")
						file.WriteString("case xml.StartElement:\n")
						for _, z := range []*Field(typ) {
							space, local := golang_getSpaceAndName(target, target.Space, z)
							file.WriteString("case t.Name.Space == " + space + " && t.Name.Local == \"" +
								local + "\":\n")
							file.WriteString("newel := &" + golang_referenceType(target, z) + "{}\n")
							file.WriteString("if err = newel.Decode(d, &t); err != nil { return err }\n")
							file.WriteString(prefix + "." + golang_makeIdent(x.Name) + " = newel\n")
							file.WriteString("if err = d.Skip(); err != nil { return err }\n")
							file.WriteString("break InLoop\n")
						}
						file.WriteString("case xml.EndElement:\n")
						file.WriteString("break InLoop\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
					}
				}
			} else {
				switch typ := x.Type.(type) {
				case Extension:
					if typ.Local == "" {
						file.WriteString("default:\n")
						file.WriteString("if newel, ok := xmlencoder.GetExtension(t.Name); ok {\n")
						file.WriteString("if err = newel.(xmlencoder.Extension).Decode(d, &t); err != nil { return err }\n")
						file.WriteString(prefix + "." + golang_makeIdent(x.Name) + " = newel\n")
						file.WriteString("} else {\n")
						file.WriteString("if err = d.Skip(); err != nil { return err }\n")
						file.WriteString("}\n")
					} else {
						typename, err := golang_getExternalType(typ.Space, typ.Local)
						if err != nil { return err }
						file.WriteString("case t.Name.Space == \"" + typ.Space + "\" && t.Name.Local == \"" +
							typ.Local + "\":")
						file.WriteString("newel := &" + typename + "{}\n")
						file.WriteString("if err = newel.Decode(d, &t); err != nil { return err }\n")
						file.WriteString(prefix + "." + golang_makeIdent(x.Name) + " = newel\n")
					}
				case Set:
					fields := []*Field(typ)
					for _, z := range fields {
						space, local := golang_getSpaceAndName(target, target.Space, z)
						file.WriteString("case t.Name.Space == " + space + " && t.Name.Local == \"" +
							local + "\":\n")
						golang_generate_element_decoder(file, target, prefix + "." + golang_makeIdent(x.Name) + "." +
							golang_makeIdent(z.Name), z)
					}
				case SequenceOf:
					switch string(typ) {
					case "extension":
						file.WriteString("default:\n")
						file.WriteString("if newel, ok := xmlencoder.GetExtension(t.Name); ok {\n")
						file.WriteString("if err = newel.(xmlencoder.Extension).Decode(d, &t); err != nil { return err }\n")
						file.WriteString(prefix + "." + golang_makeIdent(x.Name) +
							" = append(" + prefix + "." + golang_makeIdent(x.Name) + ", newel)\n")
						file.WriteString("} else {\n")
						file.WriteString("if err = d.Skip(); err != nil { return err }\n")
						file.WriteString("}\n")
					default:
						field := golang_getFieldByName(target, string(typ))
						if field != nil {
							space, local := golang_getSpaceAndName(target, target.Space, field)
							file.WriteString("case t.Name.Space == " + space + " && t.Name.Local == \"" +
								local + "\":\n")
							file.WriteString("newel := &" + golang_referenceType(target, field) + "{}\n")
							file.WriteString("if err = newel.Decode(d, &t); err != nil { return err}\n")
							file.WriteString(prefix + "." + golang_makeIdent(x.Name) +
								" = append(" + prefix + "." + golang_makeIdent(x.Name) + ", newel)\n")
						} else {
							fmt.Println("dont know how to decode 111 ", typ)
						}
					}
				case Choice:
					fields := []*Field(typ)
					for _, z := range fields {
						space, local := golang_getSpaceAndName(target, target.Space, z)
						file.WriteString("case t.Name.Space == " + space + " && t.Name.Local == \"" +
							local + "\":\n")
						file.WriteString("newel := &" + golang_referenceType(target, z) + "{}\n")
						file.WriteString("if err = newel.Decode(d, &t); err != nil { return err}\n")
						file.WriteString(prefix + "." + golang_makeIdent(x.Name) + " = newel\n")
					}
				}
			}
		}
		if any != nil {
			file.WriteString("default:\n")
			file.WriteString("if t.Name.Space == NS" + " {\n")
			switch any.EncodingRule.Type {
			case "element:name":
				typ := any.Type.(string)
				field := &Field{Type:typ}
				file.WriteString("*" + prefix + "." + golang_makeIdent(any.Name) +
					" = " + golang_referenceType(target, field) + "(t.Name.Local)\n")
				file.WriteString("if err = d.Skip(); err != nil { return err }\n")
			case "name":
				file.WriteString(prefix + "." + golang_makeIdent(any.Name) + " = t.Name.Local\n")
			case "element":
				subfields := []*Field(any.Type.(Set))
				golang_generate_element_set_decoder(file, target, prefix + "." + golang_makeIdent(any.Name), subfields)
			}
			file.WriteString("}\n")
		}
		file.WriteString("}\n")
		file.WriteString("}\n")
		file.WriteString("}\n")
	}
	var cdata *Field
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "cdata" {
			cdata = x
			break
		}
	}
	if cdata != nil {
		golang_generate_cdata_decoder(file, target, prefix + "." + golang_makeIdent(cdata.Name), cdata)
	}
	return nil
}

func golang_generate_cdata_decoder(file *os.File, target *Target, prefix string, field *Field) {
	isarray := false
	var typ string
	switch t := field.Type.(type) {
	case SequenceOf:
		isarray = true
		typ = string(t)
	case string:
		typ = string(t)
	}
	switch typ {
	case "string":
		file.WriteString("var s string\n")
		file.WriteString("if s, err = d.Text(); err != nil { return err }\n")
		if isarray {
			if prefix == "elm" {
				file.WriteString("*" + prefix + " = append(*" + prefix + ", " +
					golang_referenceType(target, field) + "(s))\n")
			} else {
				file.WriteString(prefix + " = append(" + prefix + ", s)\n")
			}
		} else {
			if prefix == "elm" {
				file.WriteString("*" + prefix + " = " +  golang_referenceType(target, field) + "(s)\n")
			} else {
				file.WriteString("*" + prefix + " = s\n")
			}
		}
	case "jid":
		file.WriteString("var s string\n")
		file.WriteString("if s, err = d.Text(); err != nil { return err }\n")
		file.WriteString("var j *jid.JID\n")
		file.WriteString("if j, err = jid.New(s); err != nil { return err }\n")
		if isarray {
			file.WriteString(prefix + " = append(" + prefix + ", j)\n")
		} else {
			file.WriteString(prefix + " = j\n")
		}
	case "bytestring":
		file.WriteString("var bdata []byte\n")
		file.WriteString("if bdata, err = d.Bytes(); err != nil { return err }\n")
		if isarray {
			file.WriteString("*" + prefix + " = append(*" + prefix + ", bdata)\n")
		} else {
			file.WriteString(prefix + " = bdata\n")
		}
	case "uint":
		file.WriteString("var s string\n")
		file.WriteString("if s, err = d.Text(); err != nil { return err }\n")
		file.WriteString("var i uint64\n")
		file.WriteString("if i, err = strconv.ParseUint(s, 10, 0); err == nil {\n")
		if isarray {
			file.WriteString("*" + prefix + " = append(*" + prefix + ", uint(i))\n")
		} else {
			file.WriteString("*" + prefix + " = uint(i)\n")
		}
		file.WriteString("}\n")
	case "int":
		file.WriteString("var s string\n")
		file.WriteString("if s, err = d.Text(); err != nil { return err }\n")
		file.WriteString("var i int64\n")
		file.WriteString("if i, err = strconv.ParseInt(s, 10, 0); err == nil {\n")
		if isarray {
			file.WriteString("*" + prefix + " = append(*" + prefix + ", int(i))\n")
		} else {
			file.WriteString("*" + prefix + " = int(i)\n")
		}
		file.WriteString("}\n")
	case "datetime":
		file.WriteString("var s string\n")
		file.WriteString("if s, err = d.Text(); err != nil { return err }\n")
		file.WriteString("var tm time.Time\n")
		file.WriteString("if tm, err = time.Parse(time.RFC3339, s); err != nil { return err }")
		if isarray {
			file.WriteString("*" + prefix + " = append(*" + prefix + ", tm)\n")
		} else {
			file.WriteString("*" + prefix + " = tm\n")
		}
	default:
		f := &Field{Type:field.Type.(string)}
		file.WriteString("var s string\n")
		file.WriteString("if s, err = d.Text(); err != nil { return err }\n")
		if isarray {
			file.WriteString("*" + prefix + " = append(*" + prefix + ", " +
				golang_referenceType(target, f) + "(s))\n")
		} else {
			file.WriteString("*" + prefix + " = " + golang_referenceType(target, f) + "(s)\n")			
		}
	}
}

func golang_getExternalType(space, local string) (string, error) {
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if target.Space == space {
				for _, field := range target.Fields {
					aspace, alocal := golang_getSpaceAndName(target, target.Space, field)
					if aspace == "NS" && alocal == local {
						if target.Name == "" {
							return schema.PackageName + "." + golang_makeIdent(field.Name), nil
						} else {
							return target.Name + "." + golang_makeIdent(field.Name), nil
						}
					}
				}
			}
		}
	}
	return "", errors.New("extenal type for " + space + " " + local + " not found")
}

func golang_isStruct(field *Field) bool {
	_, ok := field.Type.(Set)
	return ok
}

func golang_generate_check_condition(file *os.File, name string, field *Field) bool {
	switch typ := field.Type.(type) {
	case string:
		if typ == "boolean" {
			file.WriteString("if " + name + " {\n")
			return true
		}
	case Sequence, SequenceOf, Set:
		return false
		//	case string, Extension:
	}
	file.WriteString("if " + name + " != nil {\n")
	return true
}

func golang_makeIdent(s string) string {
	var r []rune
	first := true
	for _, x := range s {
		if x == rune('-') {
			first = true
			continue
		}
		if first {
			r = append(r, unicode.ToUpper(x))
			first = false
		} else {
			r = append(r, x)
		}
	}
	return string(r)
}

func golang_generate_simplevalue_encoder(name string, field *Field) string {
	isarray := false
	var typ string
	switch t := field.Type.(type) {
	case SequenceOf:
		isarray = true
		typ = string(t)
	case string:
		typ = t
	}
	switch typ {
	case "string":
		if isarray {
			if name == "elm" {
				return "string(" + name + ")"
			} else {
				return name
			}
		}
		if name == "elm" {
			return "string(*" + name + ")"
		} else {
			return "*" + name
		}
	case "boolean":
		return "strconv.FormatBool(" + name + ")"
	case "jid":
		return name + ".String()"
	case "uint":
		return "strconv.FormatUint(uint64(*" + name + "), 10)"
	case "int":
		return "strconv.FormatInt(int64(*" + name + "), 10)"
	case "datetime":
		return name + ".String()"
	}
	return "string(*" + name + ")"
}

func golang_generate_simplevalue_decoder(file *os.File, target *Target, prefix, varname string, field *Field) {
	var typ string
	switch t := field.Type.(type) {
	case SequenceOf:
		typ = string(t)
	case string:
		typ =t
	}
	switch typ {
	case "boolean":
		file.WriteString("var b bool\n")
		file.WriteString("b, err = strconv.ParseBool(" + varname + ")\n")
		file.WriteString("if err == nil {\n")
		file.WriteString(prefix + " = b\n")
		file.WriteString("}\n")
	case "string":
		file.WriteString(prefix + " = xmlencoder.Copystring(" + varname + ")\n")
	case "bytestring":
		file.WriteString(prefix + " = []byte(" + varname + ")\n")
	case "jid":
		file.WriteString("var j *jid.JID\n")
		file.WriteString("if j, err = jid.New(" + varname + "); err != nil { return err }\n")
		file.WriteString(prefix + " = j\n")
	case "uint":
		file.WriteString("var i uint64\n")
		file.WriteString("i, err = strconv.ParseUint(" + varname + ", 10, 0)\n")
		file.WriteString("if err == nil {\n")
		file.WriteString("*" + prefix + " = uint(i)\n")
		file.WriteString("}\n")
	case "int":
		file.WriteString("var i int64\n")
		file.WriteString("i, err = strconv.ParseInt(" + varname + ", 10, 0)\n")
		file.WriteString("if err == nil {\n")
		file.WriteString("*" + prefix + " = int(i)\n")
		file.WriteString("}\n")
	case "datetime":
		file.WriteString("*" + prefix + ", err = time.Parse(time.RFC3339, " + varname + ")\n")
		file.WriteString("if err != nil { return err }\n")
	case "xmllang":
		file.WriteString(prefix + " = & " + varname + "\n")
	default: // enums?
		value := golang_makeIdent(typ) + "(" + varname + ")"
		file.WriteString("value := " + value + "\n")
		file.WriteString(prefix + " = &value\n")
	}
}

func golang_generate_init(file *os.File, target *Target) {
	file.WriteString("func init() {\n")
	for _, x := range target.Fields {
		if isForClient(x) || isForServer(x) {
			local := x.Name
			if x.EncodingRule != nil && x.EncodingRule.Name != "" {
				local = x.EncodingRule.Name
			}
			file.WriteString(" xmlencoder.AddExtension(xml.Name{NS, \"" + local + "\"}, ")
			file.WriteString(golang_makeIdent(x.Name) + "{}, ")
			switch {
			case isForClient(x) && isForServer(x):
				file.WriteString("true, true)\n")
			case isForClient(x):
				file.WriteString("true, false)\n")
			case isForServer(x):
				file.WriteString("false, true)\n")
			}
		}
	}
	file.WriteString("}\n")
}

	
