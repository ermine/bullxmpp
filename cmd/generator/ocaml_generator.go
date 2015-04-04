package main

import (
	"os"
	"path/filepath"
	"unicode"
	//	"fmt"
	"errors"
	"fmt"
)

func OcamlGenerate() error {
	if cfg.Ocaml.Outdir == "" {
		panic("no outdir")
	}
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if name, ok := target.Props["name"]; ok {
				target.Name = name
			}
		}
		ocaml_generate_package(schema)
	}
	// ocaml_generate_extensions(cfg.Ocaml.Extensionfile)
	ocaml_generate_mlpack(cfg.Ocaml.Mlpack)
	return nil
}

func ocaml_generate_package(schema *Schema) error {
	dir := cfg.Ocaml.Outdir
	extension := ""
	if category, ok := schema.Props["category"]; ok {
		if category == "extension" {
			extension = "xep_"
		}
	}
	filename := filepath.Join(dir, extension + schema.PackageName + "_data.ml")
	file, err := os.Create(filename)
	if err != nil { return err }
	defer file.Close()
	return ocaml_generate_file(file, schema)
}

func ocaml_generate_file(file *os.File, schema *Schema) error {
	file.WriteString("open Types\n")
	file.WriteString("open Xmlstream")
	file.WriteString("\n")
	for _, target := range schema.Targets {
		if target.Name != "" {
			file.WriteString("module " + ocaml_makeCapital(target.Name) + " = \n")
			file.WriteString("strct\n")
			ocaml_generate_module(file, target)
			file.WriteString("end\n")
		} else {
			if len(target.Fields) > 0 {
				ocaml_generate_module(file, target)
			}
		}
	}
	return nil
}

func ocaml_normalize(s string) string {	
	var r []rune
	for _, x := range s {
		if x == rune('-') {
			r = append(r, '_')
		} else {
			r = append(r, x)
		}
	}
	return string(r)
}

func ocaml_makeCapital(s string) string {
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

func ocaml_generate_module(file *os.File, target *Target) {
	file.WriteString("let ns : Xml.namespace = Some \"" + target.Space + "\"\n\n")
	first := true
	for _, field := range target.Fields {
		if first {
			file.WriteString("type t_")
			first = false
		} else {
			file.WriteString("and t_")
		}
		file.WriteString(ocaml_normalize(field.Name) + " = ")
		ocaml_generate_type(file, "", field, false)
		ocaml_check_inner_types(file, field)
		file.WriteString("\n")
	}
	ocaml_generate_encoders(file, target)
}

func ocaml_check_inner_types(file *os.File, field *Field) {
	switch typ := field.Type.(type) {
	case Set:
		fields := []*Field(typ)
		for _, x := range fields {
			switch x.Type.(type) {
			case Set, Enum, Sequence, Choice:
				file.WriteString("and t_" + ocaml_normalize(field.Name) + "_" +
					ocaml_normalize(x.Name) + " = ")
				ocaml_generate_type(file, ocaml_normalize(field.Name) + "_", x, false)
				if _, ok := x.Type.(Set); ok {
					ocaml_check_inner_types(file, x)
				}
			}
		}
	}
}

var builtin = map[string]string{
	"string": "string",
	"int": "int",
	"uint": "uint",
	"bytestring": "bytestring",
	"jid": "jid",
	"boolean": "boolean",
	"extension": "extension",
}

func ocaml_is_builtin(t string) bool {
	if _, ok := builtin[t]; ok {
		return true
	}
	return false
}

func ocaml_generate_type(file *os.File, parent string, field *Field, deep bool) {
	switch typ := field.Type.(type) {
	case Enum:
		if deep {
			file.WriteString("t_" + parent + ocaml_normalize(field.Name))
		} else {
			file.WriteString("\n")
			for _, x := range []string(typ) {
				file.WriteString(" | ")
				file.WriteString(ocaml_makeCapital(x))
				file.WriteString("\n")
			}
		}
	case string:
		if unicode.IsUpper([]rune(typ)[0]) {
			file.WriteString("t_" + typ)
		} else {
			file.WriteString(typ)
		}
	case SequenceOf:
		t := string(typ)
		if ocaml_is_builtin(string(typ)) {
			file.WriteString(t + " list")
		} else {
			file.WriteString("t_" + t + " list")
		}			
	case Sequence:
		if deep {
			file.WriteString("t_" + parent + ocaml_normalize(field.Name) + " list")
		} else {
			for _, x := range []*Field(typ) {
				file.WriteString(" | ")
				if x.Name != "" {
					file.WriteString(ocaml_makeCapital(x.Name) + " of t_")
					ocaml_generate_type(file, ocaml_normalize(field.Name) + "_", x, false)
				} else {
					if t, ok := x.Type.(string); ok {
						file.WriteString(ocaml_makeCapital(t) + " of t_" + ocaml_normalize(t))
					}
				}
				file.WriteString("\n")
			}
		}
	case Choice:
		if deep {
			file.WriteString("t_" + parent + ocaml_normalize(field.Name))
		} else {
			file.WriteString("[\n")
			for _, x := range []*Field(typ) {
				if x.Name == "" {
					if t, ok := x.Type.(string); ok {
						file.WriteString("| `" + ocaml_makeCapital(t) + " of t_" + t + "\n")
					} else {
						file.WriteString("unknown field\n")
					}
				} else {
					file.WriteString("| `" + ocaml_normalize(x.Name) + " of ")
					ocaml_generate_type(file, "", x, true)
					file.WriteString("\n")
				}
			}
			file.WriteString("]\n")
		}
	case Extension:
		if typ.Local == "" {
			file.WriteString("extension")
		} else {
			t, err := ocaml_getExternalType(typ.Space, typ.Local)
			if err == nil {
				file.WriteString(t)
			} else {
				file.WriteString(err.Error())
			}
		}
	case Set:
		if deep {
			file.WriteString("t_" + parent + ocaml_normalize(field.Name))
		} else {
			if len([]*Field(typ)) == 0 {
				file.WriteString("()\n")
			} else {
				file.WriteString("{\n")
				for _, subfield := range []*Field(typ) {
					if subfield.Name != "" {
						file.WriteString("  f_" + ocaml_normalize(subfield.Name) + " : ")
						ocaml_generate_type(file, ocaml_normalize(field.Name) + "_", subfield, true)
						file.WriteString(";\n")
					} else {
						if t, ok := subfield.Type.(string); ok {
							file.WriteString("  f_" + t + " : t_" + t + ";\n")
						} else {
							file.WriteString("  unknown\n")
						}
					}
				}
				file.WriteString("}\n")
			}
		}
	}
}

func ocaml_getSpaceAndName(target *Target, field *Field) (string, string) {
	if s, ok := field.Type.(string); ok {
		if s == "xmllang" {
			return ns_xml, "lang"
		}
	}
	var space, local string
	if field.Name == "" {
		field1 := ocaml_getFieldByName(target, field.Type.(string))
		if field1 == nil {
			fmt.Println("Cannot find field for ", field.Type)
		}
		field = field1
	}
	local = field.Name
	if field.EncodingRule != nil && field.EncodingRule.Name != "" {
		local = field.EncodingRule.Name
	}
	space = target.Space
	if field.EncodingRule != nil && field.EncodingRule.Space != "" {
		space = field.EncodingRule.Space
	}
	return space, local
}

func ocaml_getFieldByName(target *Target, f string) *Field {
	for _, x := range target.Fields {
		if x.Name == f {
			return x
		}
	}
	return nil
}

func ocaml_getExternalType(space, local string) (string, error) {
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if target.Space == space {
				for _, field := range target.Fields {
					aspace, alocal := ocaml_getSpaceAndName(target, field)
					if aspace == space && alocal == local {
						if target.Name == "" {
							return ocaml_makeCapital(schema.PackageName) + "_data.t_" +
								ocaml_normalize(field.Name), nil
						} else {
							return ocaml_makeCapital(schema.PackageName) + "_data." +
								ocaml_makeCapital(target.Name) + ".t_" + ocaml_normalize(field.Name), nil
						}
					}
				}
			}
		}
	}
	return "", errors.New("extenal type for " + space + " " + local + " not found")
}
					
func ocaml_generate_extensions(filename string) {
	file, err := os.Create(filename)
	if err != nil { fmt.Println(err) }
	defer file.Close()
	file.WriteString("type extension = [\n")
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if len(target.Fields) == 1 {
				field := target.Fields[0]
				file.WriteString(" | `" + ocaml_makeCapital(schema.PackageName))
				if target.Name != "" {
					file.WriteString(ocaml_makeCapital(target.Name))
				}
				file.WriteString(" of ")
				file.WriteString(ocaml_makeCapital(schema.PackageName) + "_data.")
				if target.Name != "" {
					file.WriteString(ocaml_makeCapital(target.Name) + ".")
				}
				file.WriteString("t_" + ocaml_normalize(field.Name) + "\n")
			} else {
				for _, field := range target.Fields {
					if field.EncodingRule != nil && field.EncodingRule.Type == "element" {
						file.WriteString(" | `" + ocaml_makeCapital(schema.PackageName))
						if target.Name != "" {
							file.WriteString(ocaml_makeCapital(target.Name))
						}
						file.WriteString(ocaml_makeCapital(field.Name))
						file.WriteString(" of ")
						file.WriteString(ocaml_makeCapital(schema.PackageName) + "_data.")
						if target.Name != "" {
							file.WriteString(ocaml_makeCapital(target.Name) + ".")
						}
						file.WriteString("t_" + ocaml_normalize(field.Name) + "\n")
					}
				}
			}
		}
	}
	file.WriteString("]\n")
}
					
func ocaml_generate_mlpack(filename string) {
	file, err := os.Create(filename)
	if err != nil { fmt.Println(err) }
	defer file.Close()
	for _, schema := range schemas {
		file.WriteString("coders/" + ocaml_makeCapital(schema.PackageName) + "_data\n")
	}
}

func ocaml_generate_encoders(file *os.File, target *Target) {
	for _, field := range target.Fields {
		if field.EncodingRule != nil {
			file.WriteString("let encode_" + ocaml_normalize(field.Name) + " ser t = \n")
			if target.Prefix != "" {
				file.WriteString("  let () = bind_prefix ser \"" + target.Prefix + "\" ns in\n")
			}
			switch field.EncodingRule.Type {
			case "element:cdata":
				name := ""
				file.WriteString("  Xmlstring.start_tag ser \"" + name + "\";\n")
				file.WriteString("  Xmlstring.text ser t;\n")
				file.WriteString("  Xmlstring.end_tag ser;\n")
			case "element:bool":
				name := ""
				file.WriteString("  Xmlstring.start_tag ser \"" + name + "\";\n")
				file.WriteString("  Xmlstring.end_tag ser;\n")
			case "startelement":
				name := ""
				file.WriteString("  Xmlstring.start_tag ser \"" + name + "\";\n")
				
			case "element":
				
			}
			
			file.WriteString("()")
		}
	}
}
		
