package main

import (
	"os"
	"path/filepath"
	"unicode"
	"errors"
	"fmt"
	"strconv"
)

func CGenerate() error {
	var err error
	if cfg.C.Outdir == "" {
		panic("no outdir")
	}
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			target.Name = schema.PackageName
			if name, ok := target.Props["name"]; ok {
				target.Name += "_" + name
			}
		}
	}
	for _, schema := range schemas {
		if err = c_generate_h(schema); err != nil { return err }
		if err = c_generate_c(schema); err != nil { return err }
	}
	c_generate_extensions(cfg.C.Extensionfile)
	c_generate_extensions_types(cfg.C.ExtensionTypes)
	return nil
}

func c_generate_h(schema *Schema) error {
	dir := cfg.C.Outdir
	extension := ""
	if category, ok := schema.Props["category"]; ok {
		if category == "extension" {
			extension = "xep_"
		}
	}
	filename := extension + schema.PackageName + "_data"
	fullfilename := filepath.Join(dir, filename + ".h")
	file, err := os.Create(fullfilename)
	if err != nil { return err }
	defer file.Close()
	return c_generate_file_h(file, filename, schema)
}

func c_generate_file_h(file *os.File, filename string, schema *Schema) error {
	var err error
	imports := c_getImports_schema(schema)
	forwards := c_getForwards_schema(schema)
	file.WriteString("#ifndef _" + c_uppercase(filename) + "_H_\n")
	file.WriteString("#define  _" +  c_uppercase(filename) + "_H_\n\n")
	file.WriteString("#include \"xmlreader.h\"\n")
	file.WriteString("#include \"xmlwriter.h\"\n")
	file.WriteString("#include <string.h>\n")
	file.WriteString("#include \"xstream.h\"\n")
	file.WriteString("#include \"types.h\"\n")
	for _, i := range imports {
		file.WriteString("#include \"" + i + ".h\"\n")
	}
	file.WriteString("\n")
	for _, f := range forwards {
		file.WriteString("struct " + f + ";\n")
	}
	file.WriteString("\n");
	for _, target := range schema.Targets {
		file.WriteString("extern const char* ns_" + target.Name + ";\n")
	}
	file.WriteString("\n")

	for _, target := range schema.Targets {
		if err = c_generate_structs(file, target); err != nil { return err }
		c_generate_signatures(file, target);
	}
	file.WriteString("#endif\n")
	return nil
}
func c_getForwards_schema(schema *Schema) []string {
	var forwards []string
	for _, x := range schema.Targets {
		c_getForwards(schema, x.Fields, &forwards)
	}
	return forwards
}

func c_getForwards(schema *Schema, fields []*Field, forwards *[]string) {
	for _, x := range fields {
		switch typ := x.Type.(type) {
		case Set:
			fields = []*Field(typ)
			c_getForwards(schema, fields, forwards)
		case Sequence:
			fields = []*Field(typ)
			c_getForwards(schema, fields, forwards)
		case Choice:
			fields = []*Field(typ)
			c_getForwards(schema, fields, forwards)
		case Enum:
		case Extension:
		case SequenceOf:
			for _, target := range schema.Targets {
				for _, field := range target.Fields {
					if string(typ) == field.Name {
						append_import(forwards, c_makeType(target.Name, field.Name))
					}
				}
			}
		case string:
			for _, target := range schema.Targets {
				for _, field := range target.Fields {
					if typ == field.Name {
						switch field.Type.(type) {
						case Set, Sequence, Choice, SequenceOf:
							append_import(forwards, c_makeType(target.Name, field.Name))
						}
					}
				}
			}
		}
	}
}

func c_getImports_schema(schema *Schema) []string {
	var imports []string
	for _, x := range schema.Targets {
		c_getImports(x.Fields, &imports)
	}
	return imports
}

func c_getImports(fields []*Field, imports *[]string) {
	for _, x := range fields {
		switch typ := x.Type.(type) {
		case Set:
			fields = []*Field(typ)
			c_getImports(fields, imports)
		case Sequence:
			fields = []*Field(typ)
			c_getImports(fields, imports)
		case Choice:
			fields = []*Field(typ)
			c_getImports(fields, imports)
		case Enum:
		case Extension:
			if typ.Local != "" {
			Added:
				for _, schema := range schemas {
					for _, target := range schema.Targets {
						if target.Space == typ.Space {
							extension := ""
							if category, ok := schema.Props["category"]; ok {
								if category == "extension" {
									extension = "xep_"
								}
							}
							append_import(imports, extension + schema.PackageName + "_data")
							break Added
						}
					}
				}
			}
		case SequenceOf:
		case string:
			if typ == "jid" {
				append_import(imports, "jid/jid")
			}
		}
	}
}

/*
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
*/

func c_uppercase(s string) string {
	var r []rune
	for _, x := range s {
		if x == rune('-') {
			r = append(r, '_')
		} else {
			r = append(r, unicode.ToUpper(x))
		}
	}
	return string(r)
}

func c_collect_enums(target *Target, fields []*Field, enums *[]*Field, prefix string) {
	for _, x := range fields {
		switch t := x.Type.(type) {
		case Set:
			fields := []*Field(t)
			if prefix != "" {
				prefix += "_"
			}
			c_collect_enums(target, fields, enums, prefix + c_normalize(x.Name))
		case Sequence:
			fields := []*Field(t)
			if prefix != "" {
				prefix += "_"
			}
			c_collect_enums(target, fields, enums, prefix + c_normalize(x.Name))
		case Choice:
			fields := []*Field(t)
			if prefix != "" {
				prefix += "_"
			}
			c_collect_enums(target, fields, enums, prefix + c_normalize(x.Name))
			c_checkTypes(target, fields)
		case Enum:
			field := &Field {
				Name:x.Name,
				Type: x.Type,
				EncodingRule:x.EncodingRule,
				DefaultValue:x.DefaultValue,
				Required:x.Required,
			}
			if prefix != "" {
				typename := prefix + "_" + x.Name
				x.Type = typename
				field.Name = typename
			}
			*enums = append(*enums, field)
		}
	}
	// return enums
}

func c_generate_structs(file *os.File, target *Target) error {
	var err error
	var enums []*Field
	c_collect_enums(target, target.Fields, &enums, "")
	for _, x := range enums {
		exists := false
		for _, z := range target.Fields {
			if x.Name == z.Name {
				exists = true
				break
			}
		}
		if !exists {
			target.Fields = append(target.Fields, x)
		}
	}
	c_generate_enums(file, target, enums)
	file.WriteString("\n")
	for _, def := range target.Fields {
		if _, ok := def.Type.(Enum); ok {
			continue
		}
		switch t := def.Type.(type) {
		case Extension:
			if t.Local != "" {
				field, err := c_getExternalType(t.Space, t.Local)
				if err != nil { return err }
				file.WriteString("typedef " + field + " " + c_makeType(target.Name, def.Name) + ";\n")
			} else {
				file.WriteString("typedef extension_t* " + c_makeType(target.Name, def.Name) + ";\n")
			}
		case SequenceOf:
			// file.WriteString("array_t* " + c_makeIdent(def.Name) + ";\n")
			// file.WriteString("typedef array_t " + c_makeType(target.Name, def.Name) + ";\n")
		case string:
			var typ string
			typ = c_simpletype(t);
			if typ == t {
				field := c_getFieldByName(target, t)
				if field == nil {
					return errors.New("unknown type " + t)
				}
				typ = c_referenceType(target, field)
			}
				if typ[len(typ)-1] == '*' {
					typ = typ[:len(typ)-1]
				}
			file.WriteString("typedef " + typ + " " + c_makeType(target.Name, def.Name) + ";\n")
		case Set:
			file.WriteString("struct " + c_makeType(target.Name, def.Name) + " {\n")
			if err = c_generate_fields(file, target, []*Field(t)); err != nil { return err }
			file.WriteString("};\n")
		case Sequence:
			file.WriteString("struct " + c_makeType(target.Name, def.Name) + " {\n")
			file.WriteString("   size_t length;\n")
			file.WriteString("  struct " + c_makeType(target.Name, def.Name + "_sequence") + " {\n")
			file.WriteString("    int type;\n")
			file.WriteString("    union {\n")
			for _, x := range []*Field(t) {
				file.WriteString("    struct " + c_makeType(target.Name, x.Name) + "* " +
					c_makeIdent(x.Name) + ";\n")
			}
			file.WriteString("    } u;\n")
			file.WriteString("  } sequence_t* sequence;\n")
			file.WriteString("};")
		case Choice:
			def.Type = Set([]*Field{&Field{Name:"u", Type:t}})
			file.WriteString("struct " + c_makeType(target.Name, def.Name) +  " {\n")
			file.WriteString("  int type;\n")
			file.WriteString("  union {\n")
			for _, x := range []*Field(t) {
				name := x.Name
				if name == "" {
					name = x.Type.(string)
				}
				// file.WriteString("    struct " + c_makeType(target.Name, name) + "* " +
				// file.WriteString("    " + c_makeType(target.Name, name) + "* " +
				file.WriteString("    " + c_referenceType(target, x) + c_makeIdent(name) + ";\n")
			}
			file.WriteString("  } *u;\n")
			file.WriteString("};")
		}
		file.WriteString("\n\n")
	}
	return nil
}

func c_generate_enums(file *os.File, target *Target, enums []*Field) {
	for _, x := range enums {
		t := c_makeType(target.Name, x.Name)
		file.WriteString("enum " + t  + " {\n")
		enum := []string(x.Type.(Enum))
		for _, z := range enum {
			file.WriteString("  " + c_uppercase (target.Name + "_" + x.Name + "_" + z) + ",\n")
		}
		file.WriteString("};\n\n")
		v := t[:len(t)-2]
		file.WriteString("enum " + c_makeType(target.Name, x.Name) + " enum_" + v +
			"_from_string(const char* value);\n")
		file.WriteString("const char* enum_" + v + "_to_string(enum " + t + ");\n")
	}
}

func c_generate_fields(file *os.File, target *Target, fields []*Field) error {
	for _, x := range fields {
		if x.Name == "" {
			x.Name = x.Type.(string)
		}
		switch t := x.Type.(type) {
		case string:
			switch t {
				case "bytestring":
				file.WriteString("  unsigned char* " + c_makeIdent(x.Name) + ";\n")
				file.WriteString("  int " + c_makeIdent(x.Name) + "_len;\n")
			case "boolean":
				if x.EncodingRule != nil && (x.EncodingRule.Type == "attribute" ||
					x.EncodingRule.Type == "element:cdata") { 
					file.WriteString("  bool* " + c_makeIdent(x.Name) + ";\n")
				} else {
					file.WriteString("  bool " + c_makeIdent(x.Name) + ";\n")
				}
			default:
				var typ string
				typ = c_simpletype(t);
				if typ == t && typ != "int" {
					field := c_getFieldByName(target, t)
					if field == nil {
						return errors.New("unknown type " + t)
					}
					typ = c_referenceType(target, field)
				}
				if x.EncodingRule != nil && x.EncodingRule.Type == "element:bool" {
					file.WriteString("  " + typ + " " + c_makeIdent(x.Name) + ";\n")
				} else {
					file.WriteString("  " + typ + " " + c_makeIdent(x.Name) + ";\n")
				}
			}
		case SequenceOf, Sequence:
			file.WriteString("  array_t* " + c_makeIdent(x.Name) + ";\n")
		case Extension:
			if t.Local != "" {
				fieldtype, err := c_getExternalType(t.Space, t.Local)
				if err != nil { return err }
				file.WriteString(fieldtype + " " + c_makeIdent(x.Name) + ";\n")
			} else {
				file.WriteString("  extension_t* " + c_makeIdent(x.Name) + ";\n")
			}
		case Choice:
			file.WriteString(" extension_t* " + c_makeIdent(x.Name) + ";\n")
		case Set:
			fields := []*Field(t)
			file.WriteString(" struct " + c_makeType(target.Name, x.Name + "_set") + " {\n")
			c_generate_fields(file, target, fields)
			file.WriteString("} " + c_makeIdent(x.Name) + ";\n")
		default:
			fmt.Println("default1: ", x.EncodingRule)
		}
	}
	return nil
}

func c_makeType(prefix, s string) string {
	var r []rune
	for _, x := range s {
		if x == rune('-') {
			r = append(r, '_')
		} else {
			r = append(r, x)
		}
	}
	return prefix + "_" + string(r) + "_t"
}

func c_makeIdent(s string) string {
	var r []rune
	first := true
	for _, x := range s {
		if x == rune('-') {
			r = append(r, '_')
		} else {
			if first {
				r = append(r, unicode.ToUpper(x))
				first = false
			} else {
				r = append(r, x)
			}
		}
	}
	return "f" + string(r)
}

func c_normalize(s string) string {
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

func c_referenceType(target *Target, field *Field) string {
	name := field.Name
	if name == "" {
		name = field.Type.(string)
	}
	if _, ok := field.Type.(Choice); ok {
		return ""
	}
	for _, x := range target.Fields {
		if x.Name == name {
			switch x.Type.(type) {
			case Enum:
				return "enum " + c_makeType(target.Name, name)
			case string, Extension:
				return c_makeType(target.Name, name) + "*"
			default:
				return "struct " + c_makeType(target.Name, name) + "*"
			}
		}
	}
	return c_simpletype(name)
}

func c_makeConverterFromValue(target *Target, field *Field, varname string) string {
	name := field.Name
	if name == "" {
		name = field.Type.(string)
	} else if _, ok := field.Type.(string); ok {
				name = field.Type.(string)
	}
	for _, x := range target.Fields {
		if x == field {
			return "(" + c_referenceType(target, field) + ")" + varname
		}
		if x.Name == name {
			if _, ok := x.Type.(Enum); ok {
				t := c_makeType(target.Name, name)
				v := t[:len(t)-2]
				return "enum_" + v + "_from_string(" + varname + ")"
			} else {
				return "(" + c_makeType(target.Name, name) + "*)" + varname
			}
		}
	}
	// return "(" + c_simpletype(name) + ")"
	// return "(" + c_referenceType(target, field) + ")"
	return "(char*)" +varname
}

func c_makeConverterToValue(target *Target, field *Field, varname string) string {
	name := field.Name
	if name == "" {
		name = field.Type.(string)
	} else if _, ok := field.Type.(string); ok {
				name = field.Type.(string)
	}
	for _, x := range target.Fields {
		if x == field {
			return "(" + c_referenceType(target, field) + ")" + varname
		}
		if x.Name == name {
			if _, ok := x.Type.(Enum); ok {
				t := c_makeType(target.Name, name)
				v := t[:len(t)-2]
				return "enum_" + v + "_to_string(" + varname + ")"
			} else {
				return "(" + c_makeType(target.Name, name) + "*)" + varname
			}
		}
	}
	// return "(" + c_simpletype(name) + ")"
	// return "(" + c_referenceType(target, field) + ")"
	return c_generate_simplevalue_encoder (varname, field)
}


func c_simpletype(s string) string {
	switch s {
	case "jid": return "jid_t* "
	case "int": return "int *"
	case "uint": return "uint32_t *"
	case "string": return "char*"
	case "boolean": return "int *"
	case "langstring": return "langstring_t *"
	case "xmllang": return "char*"
	case "bytestring":
		return "unsigned char* ";
	case "extension": return "void *"
	case "datetime": return "struct tm*"
	}
	return s
}

func c_getFieldByName(target *Target, f string) *Field {
	for _, x := range target.Fields {
		if x.Name == f {
			return x
		}
	}
	return nil
}

func c_isStruct(field *Field) bool {
	_, ok := field.Type.(Set)
	return ok
}

func c_checkTypes(target *Target, fields []*Field) {
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

func c_getExtensionDecoder(space, local string) (string, error) {
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if target.Space == space || "ns_" + target.Name == space {
				for _, field := range target.Fields {
					_, alocal := c_getSpaceAndName(target, target.Space, field)
					if  alocal == local {
						return target.Name + "_" + c_normalize(field.Name) + "_decode", nil
					}
				}
			}
		}
	}
	return "", errors.New("extenal type for " + space + " " + local + " not found")
}	
	

func c_getExtensionEncoder(space, local string) (string, error) {
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if target.Space == space || "ns_" + target.Name == space {
				for _, field := range target.Fields {
					_, alocal := c_getSpaceAndName(target, target.Space, field)
					if  alocal == local {
						return target.Name + "_" + c_normalize(field.Name) + "_encode", nil
					}
				}
			}
		}
	}
	return "", errors.New("extenal type for " + space + " " + local + " not found")
}	

func c_getExtensionFree(space, local string) (string, error) {
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if target.Space == space || "ns_" + target.Name == space {
				for _, field := range target.Fields {
					_, alocal := c_getSpaceAndName(target, target.Space, field)
					if  alocal == local {
						return target.Name + "_" + c_normalize(field.Name) + "_free", nil
					}
				}
			}
		}
	}
	return "", errors.New("extenal type for " + space + " " + local + " not found")
}	


func c_getExternalType(space, local string) (string, error) {
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if target.Space == space {
				for _, field := range target.Fields {
					_, alocal := c_getSpaceAndName(target, target.Space, field)
					if target.Space == space && alocal == local {
						return "struct " + c_makeType(target.Name, field.Name) + "*", nil
					}
				}
			}
		}
	}
	return "", errors.New("extenal type for " + space + " " + local + " not found")
}

func c_getSpaceAndName(target *Target, targetNS string, field *Field) (string, string) {
	if s, ok := field.Type.(string); ok {
		if s == "xmllang" {
			return "ns_xml", "lang"
		}
	}
	var space, local string
	if field.Name == "" {
		field1 := c_getFieldByName(target, field.Type.(string))
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
	if space == targetNS {
		space = "ns_" + target.Name
	} else {
		space = "\"" + space + "\""
}
	return space, local
}

func c_generate_signatures(file *os.File, target *Target) {
	for _, field := range target.Fields {
		switch field.Type.(type) {
		case Enum:
		case SequenceOf:
			name := target.Name + "_" + c_normalize(field.Name)
			file.WriteString("array_t* " + name + "_decode(xmlreader_t* reader);\n")
			file.WriteString("int " + name + "_encode(xmlwriter_t* writer, array_t* data);\n")
			file.WriteString("void " + name + "_free (array_t* data);\n")
		case Set, Sequence, Choice:
			name := target.Name + "_" + c_normalize(field.Name)
			file.WriteString("struct " + c_makeType(target.Name, field.Name) + "* " +
				name + "_decode(xmlreader_t* reader);\n")
			file.WriteString("int " + name + "_encode(xmlwriter_t* writer, struct " +
				c_makeType(target.Name, field.Name) + "* data);\n")
			file.WriteString("void " + name + "_free (struct " +
								c_makeType(target.Name, field.Name) + "* data);\n")
		default:
			name := target.Name + "_" + c_normalize(field.Name)
			file.WriteString(c_makeType(target.Name, field.Name) + "* " +
				name + "_decode(xmlreader_t* reader);\n")
			file.WriteString("int " + name + "_encode(xmlwriter_t* writer, " +
				c_makeType(target.Name, field.Name) + "* data);\n")
			file.WriteString("void " + name + "_free (" +
				c_makeType(target.Name, field.Name) + "* data);\n")
		}
	}
}

func c_generate_c(schema *Schema) error {
		dir := cfg.C.Outdir
	extension := ""
	if category, ok := schema.Props["category"]; ok {
		if category == "extension" {
			extension = "xep_"
		}
	}
	filename := extension + schema.PackageName + "_data"
	fullfilename := filepath.Join(dir, filename + ".c")
	file, err := os.Create(fullfilename)
	if err != nil { return err }
	defer file.Close()
	return c_generate_file_c(file, filename, schema)
}

func c_generate_file_c(file *os.File, filename string, schema *Schema) error {
	var err error
	file.WriteString("#include \"" + filename + ".h\"\n")
	file.WriteString("#include \"helpers.h\"\n")
	file.WriteString("#include \"errors.h\"\n\n")
	for _, target := range schema.Targets {
		file.WriteString("const char* ns_" + target.Name + " = \"" + target.Space + "\";\n")
	}
	file.WriteString("\n")
	for _, target := range schema.Targets {
		for _, field := range target.Fields {
			switch field.Type.(type) {
			case Enum:
				c_generate_enum_decoder(file, target, field)
				c_generate_enum_encoder(file, target, field)

			case SequenceOf:
				name := target.Name + "_" + c_normalize(field.Name)
				file.WriteString("array_t* " + name + "_decode(xmlreader_t* reader) {\n")
				// todo: type of sizeof
				file.WriteString(" array_t* elm = array_new (sizeof (extension_t), 0);\n")
				if err := c_generate_element_decoder(file, target, "elm", field); err != nil { return err }
				file.WriteString("  return elm;\n")
				file.WriteString("}\n\n")
				file.WriteString("int " + name + "_encode(xmlwriter_t* writer, array_t* elm) {\n")
				file.WriteString("int err = 0;\n")
				if target.Prefix != "" {
					file.WriteString("err = xmlwriter_set_prefix (writer, \"" + target.Prefix + "\", ns_" +
						target.Name + ");\n")
					file.WriteString("if (err != 0) return err;\n")
				}
				if err := c_generate_element_encoder(file, target, "elm", field); err != nil { return err }
				file.WriteString("  return 0;\n")
				file.WriteString("}\n\n")
				file.WriteString("void " + name + "_free (array_t* data) {\n")
				file.WriteString("if (data == NULL)\n")
				file.WriteString("return;\n")
				file.WriteString("int len = array_length (data);\n")
				file.WriteString("int i = 0;\n")
				file.WriteString("for (i = 0; i < len; i++) {\n")
				t := field.Type.(SequenceOf)
				if (t == "extension") {
					file.WriteString("extension_t* ext = array_get (data, i);\n")
					file.WriteString("xstream_extension_free (ext);\n")
				} else {
					f := c_getFieldByName(target, string(t))
					if f == nil {
						return errors.New("unknown type " + string(t))
					}
					typ := c_referenceType(target, f)
					file.WriteString(typ + "* item = array_get (data, i);\n")
					freer, err := c_getExtensionFree(target.Space, string(t))
					if err != nil { return err }
					file.WriteString(freer + "(item);\n")
				}
				file.WriteString("}\n")
				file.WriteString("array_free (data);\n")
				file.WriteString("}\n\n")
				
			case Set, Sequence, Choice:
				name := target.Name + "_" + c_normalize(field.Name)
				file.WriteString("struct " + c_makeType(target.Name, field.Name) + "* " +
					name + "_decode(xmlreader_t* reader) {\n")
				t := c_referenceType(target, field)
				if t[len(t)-1] == '*' {
					t = t[:len(t)-1]
				}
				file.WriteString("  struct " + c_makeType(target.Name, field.Name) + " *elm = NULL;\n")
				file.WriteString(" elm = malloc (sizeof (" + t + "));\n")
				file.WriteString("if (elm == NULL)\n")
				file.WriteString("fatal (\"" + c_makeType (target.Name, field.Name) + ": malloc failed\");\n")
				file.WriteString("memset (elm, 0, sizeof (" + t + "));\n")
				if err := c_generate_element_decoder(file, target, "elm", field); err != nil { return err }
				file.WriteString("  return elm;\n")
				file.WriteString("}\n\n")
				file.WriteString("int " + name + "_encode(xmlwriter_t* writer, struct " +
					c_makeType(target.Name, field.Name) + "* elm) {\n")
				file.WriteString("int err = 0;\n")
				if target.Prefix != "" {
					file.WriteString("err = xmlwriter_set_prefix (writer, \"" + target.Prefix + "\", ns_" +
						target.Name + ");\n")
					file.WriteString("if (err != 0) return err;\n")
				}
				if err := c_generate_element_encoder(file, target, "elm", field); err != nil { return err }
				file.WriteString("  return 0;\n")
				file.WriteString("}\n\n")
				file.WriteString("void " + name + "_free (struct " + 
					c_makeType(target.Name, field.Name) + "* data) {\n")
				file.WriteString("if (data == NULL) return;\n")
				c_generate_free (file,  target, "data", field)
				file.WriteString("free (data);\n")
				file.WriteString("}\n\n")
			default:
				name := target.Name + "_" + c_normalize(field.Name)
				file.WriteString(c_makeType(target.Name, field.Name) + "* " +
					name + "_decode(xmlreader_t* reader) {\n")
				file.WriteString("  " + c_makeType(target.Name, field.Name) + " *elm = NULL;\n")
				c_generate_element_decoder(file, target, "elm", field)
				if _, ok := field.Type.(SequenceOf); ok {
					file.WriteString("  return *elm;\n")
				} else {
					file.WriteString("  return elm;\n")
				}
				file.WriteString("}\n\n")
				file.WriteString("int " + name + "_encode(xmlwriter_t* writer, " +
					c_makeType(target.Name, field.Name) + "* elm) {\n")
				file.WriteString("int err = 0;\n")
				if target.Prefix != "" {
					file.WriteString("err = xmlwriter_set_prefix (writer, \"" + target.Prefix + "\", ns_" +
						target.Name + ");\n")
					file.WriteString("if (err != 0) return err;\n")
				}
				if err := c_generate_element_encoder(file, target, "elm", field); err != nil { return err }
				file.WriteString("  return 0;\n")
				file.WriteString("}\n\n")
				file.WriteString("void " + name + "_free (" + 
					c_makeType(target.Name, field.Name) + "* data) {\n")
				file.WriteString("if (data == NULL) return;\n")
				c_generate_free (file,  target, "data", field)
				file.WriteString("free (data);\n")
				file.WriteString("}\n\n")
			}
		}
	}
	return err
}

func c_generate_element_decoder(file *os.File, target *Target, prefix string, element *Field) error {
	if element.EncodingRule == nil {
		fmt.Println("dont know how to generate decoder ", element.Name)
		return nil
	}
	switch element.EncodingRule.Type {
	case "element:cdata":
		c_generate_cdata_decoder(file, target, prefix, element)
	case "element:bool":
		file.WriteString("if (xmlreader_skip_element (reader) == 1) return NULL;\n")
	case "startelement", "element":
		switch typ := element.Type.(type) {
		case string:
			return errors.New("dont know what todo with " + typ)
		case Extension:
			file.WriteString("int type = 0;\n")
			file.WriteString("  while (1) {\n")
			file.WriteString("type = xmlreader_next (reader);\n")
			file.WriteString("if (type == XML_END_ELEMENT) break;\n")
			file.WriteString("if (type == XML_ERROR) return NULL;\n")
			file.WriteString("if (type == XML_START_ELEMENT) {\n")
			if typ.Local != "" {
				// fieldname, err := c_getExternalType(typ.Space, typ.Local)
				file.WriteString("const char* name = xmlreader_get_name (reader);\n")
				file.WriteString("const char* namespace = xmlreader_get_namespace (reader);\n")
				file.WriteString("if ((strcmp (namespace, \"" + typ.Space +
					"\") == 0) && (strcmp (name, \"" + typ.Local + "\") == 0)) {\n")
				decoder, err := c_getExtensionDecoder(typ.Space, typ.Local)
				if err != nil { return err }
				file.WriteString(c_referenceType(target, element) + " newel = " + decoder + "(reader);\n")
				file.WriteString("if (newel == NULL) return NULL;\n")
				if prefix == "elm" {
					file.WriteString(prefix + " = " +
						c_makeConverterToValue(target, element, "newel") + ";\n")
				} else {
					file.WriteString(prefix + " = newel;\n")
				}
				file.WriteString("} else {\n")
				file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
				file.WriteString("}\n")
			} else {
				file.WriteString("extension_t* newel = malloc (sizeof (extension_t));\n")
				file.WriteString("int err = xstream_extension_decode (reader, ext);\n")
				file.WriteString("if (reader->err != 0) return NULL;\n")
				file.WriteString("if (err == ERR_EXTENSION_NOT_FOUND) {\n")
				file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
				file.WriteString("} else {\n")
				if prefix == "elm" {
					file.WriteString(prefix + " = " +
						c_makeConverterToValue(target, element, "(newel") + ";\n")
				} else {
					file.WriteString(prefix + " = newel;\n")
				}
				file.WriteString("}\n")
			}
			file.WriteString("}\n")
//			file.WriteString("}\n")
			file.WriteString("}\n")
		case Set:
			fields := []*Field(typ)
			if len(fields) == 0 {
				file.WriteString("if (xmlreader_skip_element (reader) == -1)return NULL;\n")
			} else {
				close, err := c_generate_element_set_decoder(file, target, prefix, fields)
				if err != nil { return err }
				if element.EncodingRule.Type != "startelement" && close {
					file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
				}
			}
		case Choice:
			fields := []*Field(typ)
			file.WriteString("int type = 0;\n")
			file.WriteString("  while (1) {\n")
			file.WriteString("type = xmlreader_next (reader);\n")
			file.WriteString("if (type == XML_END_ELEMENT) break;\n")
			file.WriteString("if (type == XML_ERROR) return NULL;\n")
			file.WriteString("if (type == XML_START_ELEMENT) {\n")
			for _, z := range fields {
				space, local := c_getSpaceAndName(target, target.Space, z)
				file.WriteString(" if ((strcmp (namespace, " + space +
					") == 0) && (strcmp (name,  \"" + local + "\") == 0)) {\n")
				decoder, err := c_getExtensionDecoder(target.Space, local)
				if err != nil { return err }
				file.WriteString(c_referenceType(target, z) + " newel = " + decoder + "(reader);\n")
				file.WriteString("if (newel == NULL) return NULL;\n")
				file.WriteString(prefix + "->u = (" + c_referenceType(target, element) + "*) newel;\n")
				file.WriteString("  }\n")
			}
			file.WriteString("}\n")
			file.WriteString(" }\n")
			
		case SequenceOf:
			field := string(typ)
			file.WriteString("int type = 0;\n")
			file.WriteString("  while (1) {\n")
			file.WriteString("type = xmlreader_next (reader);\n")
			file.WriteString("if (type == XML_END_ELEMENT) break;\n")
			file.WriteString("if (type == XML_ERROR) return NULL;\n")
			file.WriteString("if (type == XML_START_ELEMENT) {\n")
			if field == "extension" {
				name := prefix + "->" + element.Name
				if prefix == "elm" {
					name = prefix
				}
				file.WriteString("extension_t ext;\n")
				file.WriteString("int err = xstream_extension_decode (reader, &ext);\n")
				file.WriteString("if (reader->err != 0) return NULL;\n")
				file.WriteString("if (err == ERR_EXTENSION_NOT_FOUND) {\n")
				file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
				file.WriteString("} else {\n")
				if prefix != "elm" {
					file.WriteString("if (" + name + " == NULL)\n")
					file.WriteString(name + " = array_new (sizeof (extension_t), 0);\n")
				}
				file.WriteString("array_append (" + name + ", &ext);\n")
				file.WriteString("}\n")
			} else {
				f := c_getFieldByName(target, field)
				if f == nil {
					// import from other packages?
					return errors.New("dont know what to do with " + field)
				}
				space, local := c_getSpaceAndName(target, target.Space, f)
				file.WriteString("const char* name = xmlreader_get_name (reader);\n")
				file.WriteString("const char* namespace = xmlreader_get_namespace (reader);\n")
				file.WriteString("if ((strcmp (namespace, " + space +
					") == 0) && (strcmp (name, \"" + local + "\") == 0)) {\n")
				decoder, err := c_getExtensionDecoder(space, local)
				if err != nil { return err }
				file.WriteString(c_referenceType(target, f) + " newel = " + decoder + "(reader);\n")
				file.WriteString("  if (newel == NULL) {\n    return NULL;\n  }\n")
				name := prefix
				if name != "elm" {
					name = "&" + name
				}
				if prefix != "elm" {
					file.WriteString("if (" + name + " == NULL)\n")
					file.WriteString(name + " = array_new (sizeof (extension_t), 0);\n")
				}
				file.WriteString("  array_append (" + name + ", newel);\n")
				file.WriteString("}\n")
			}
			file.WriteString("}\n")
			file.WriteString("}\n")
		}
	}
	return nil
}

func c_generate_cdata_decoder(file *os.File, target *Target, prefix string, field *Field) {
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
		file.WriteString("const char* value = xmlreader_text (reader);\n")
		file.WriteString("if (reader->err != 0) return NULL;\n")
		name := prefix
		if isarray {
			if prefix != "elm" {
				file.WriteString("if (" + name + " == NULL)\n")
				file.WriteString(name + " = array_new (sizeof (char*), 0);\n")
			}
			file.WriteString("  array_append (" + name + ", (void*) &value);\n")
		} else {
//			if prefix == "elm" {
			file.WriteString(prefix + " = " +  c_makeConverterFromValue(target, field, "value") + ";\n")
//			} else {
//				file.WriteString(prefix + " = value;\n")
//			}
		}
	case "jid":
		file.WriteString("const char* s = xmlreader_text (reader);\n")
		file.WriteString("if (reader->err != 0) return NULL;\n")
		file.WriteString("  jid_t *jid = jid_of_string((const char*) s);\n")
		if isarray {
			file.WriteString("if (" + prefix + " == NULL)\n")
			file.WriteString(prefix + " = array_new (sizeof (jid_t*), 0);\n")
			file.WriteString("  array_append (" + prefix + ", &jid);\n")
		} else {
			file.WriteString(prefix + " = jid;\n")
		}
	case "bytestring":
		if isarray {
			file.WriteString("todo();\n")
		} else {
			file.WriteString("if (xmlreader_base64 (reader, &" + prefix + ", &" + prefix + "_len) != 0) return NULL;\n")
		}
	case "uint":
		file.WriteString(" const char* s = xmlreader_text (reader);\n")
		file.WriteString("if (reader->err != 0) return NULL;\n")
		if isarray {
			file.WriteString("if (" + prefix + " == NULL)\n")
			file.WriteString(prefix + " = array_new (sizeof (uint32_t), 0);\n")
			file.WriteString("  array_append (" + prefix + "strconv_parse_uint64 (s), 0);\n")
		} else {
			file.WriteString("  " + prefix + " = strconv_parse_uint64 (s);\n")
		}
	case "int":
		file.WriteString(" const char* s = xmlreader_text (reader);\n")
		file.WriteString("if (reader->err != 0) return NULL;\n")
		if isarray {
			file.WriteString("if (" + prefix + " == NULL)\n")
			file.WriteString(prefix + " = array_new (sizeof (int), 0);\n")
			file.WriteString("  array_append (" + prefix + ", strconv_parse_int (s), 0);\n")
		} else {
			file.WriteString("  " + prefix + " = strconv_parse_int (s);\n")
		}
	case "datetime":
		file.WriteString(" const char* s = xmlreader_text (reader);\n")
		file.WriteString("if (reader->err != 0) return NULL;\n")
		file.WriteString("if tm, err = time.Parse(time.RFC3339, s); err != nil { return err }")
		if isarray {
			file.WriteString("if (" + prefix + " == NULL)\n")
			file.WriteString(prefix + " = array_new (sizeof (struct tm), 0);\n")
			file.WriteString("struct tm* datetime = datetime_parse(s);\n")
			file.WriteString("array_append (" + prefix + ", datetime);\n")
		} else {
			file.WriteString(prefix + " = datetime_parse(s);\n")
		}
	default:
		f := &Field{Type:field.Type.(string)}
		file.WriteString(" const char* s = xmlreader_text (reader);\n")
		file.WriteString("if (reader->err != 0) return NULL;\n")
		if isarray {
			file.WriteString(prefix + " = append(*" + prefix + ", " +
				c_referenceType(target, f) + "(s));\n")
		} else {
			file.WriteString(prefix + " = " + c_makeConverterFromValue(target, f, "s") + ";\n")
		}
	}
}

func c_generate_element_set_decoder(file *os.File, target *Target, prefix string, fields []*Field) (bool, error) {
	sep := "->"
	if prefix != "elm" {
		sep = "."
	}
	// var err error
	var attrs []*Field
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
			attrs = append(attrs, x)
		}
	}
	if len(attrs) > 0 {
		file.WriteString("  const char* avalue;\n")
	}
	for _, x := range attrs {
		space, local := c_getSpaceAndName(target, "", x)
		if space == ns_xml {
			local = "xml:" + local
		}
		file.WriteString("  avalue = xmlreader_attribute (reader, NULL, \"" + local + "\");\n")
		file.WriteString("  if (avalue != NULL) {\n")
		c_generate_simplevalue_decoder(file, target, prefix + sep + c_makeIdent(x.Name), "avalue", x)
		file.WriteString("  }\n")
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
		elseif := false
		file.WriteString("int type = 0;\n")
		file.WriteString("  while (1) {\n")
		file.WriteString("type = xmlreader_next (reader);\n")
		file.WriteString("if (type == XML_END_ELEMENT) break;\n")
		file.WriteString("if (type == XML_ERROR) return NULL;\n")
		file.WriteString("if (type == XML_START_ELEMENT) {\n")
		file.WriteString("const char* namespace = xmlreader_get_namespace (reader);\n")
		file.WriteString("const char* name = xmlreader_get_name (reader);\n")
		for _, x := range elems {
			if x.EncodingRule != nil {
				space, local := c_getSpaceAndName(target, target.Space, x)
				if !elseif {
					// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
					// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
				}				
				if elseif {
					file.WriteString("  else if")
				} else {
					elseif = true
					file.WriteString(" if")
				}
				file.WriteString(" ((strcmp (name, \"" + local +
					"\") == 0) && (strcmp (namespace, " + space + ") == 0)) {\n")
				switch x.EncodingRule.Type {
				case "element:bool":
					file.WriteString("    " + prefix + sep + c_makeIdent(x.Name) + " = true;\n")
					file.WriteString("  if (xmlreader_skip_element (reader) == -1) return NULL;\n")
					file.WriteString("  continue;\n")
				case "element:cdata":
					c_generate_cdata_decoder(file, target, prefix + sep + c_makeIdent(x.Name), x)
				case "element":
					switch typ := x.Type.(type) {
					case string:
						if typ == "langstring" {
							file.WriteString("    langstring_decode (reader, " +
								prefix + sep + c_makeIdent(x.Name) + ");\n")
						} else {
							file.WriteString(prefix + sep + c_makeIdent(x.Name) +
								" = xmlreader_text (reader);\n")
							file.WriteString("if (reader->err != 0) return NULL;\n")
						}
					case Set:
						fields := []*Field(typ)
						close, err := c_generate_element_set_decoder(file, target,
							prefix + sep + c_makeIdent(x.Name), fields)
						if err != nil { return false, err }
						if close {
							file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
						}
					case SequenceOf:
						field := c_getFieldByName(target, string(typ))
						if field == nil {
							return false, errors.New("Cannot find field " + string(typ))
						}
						space, local := c_getSpaceAndName(target, target.Space, field)
						file.WriteString("int type = 0;\n")
						file.WriteString("  while (1) {\n")
						file.WriteString("type = xmlreader_next (reader);\n")
						file.WriteString("if (type == XML_END_ELEMENT) break;\n")
						file.WriteString("if (type == XML_ERROR) return NULL;\n")
						file.WriteString("if (type == XML_START_ELEMENT) {\n")
						// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
						// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
						file.WriteString("  if ((strcmp(namespace, " + space +
							") == 0) && (strcmp (name, \"" + local + "\") == 0)) {\n")
						decoder, err := c_getExtensionDecoder(space, local)
						if err != nil { return false, err }
						file.WriteString(c_referenceType(target, field) + " newel = " + decoder + "(reader);\n")
						file.WriteString("  if (newel != NULL) {\n    return NULL;\n  }\n")
						name := prefix + sep + c_makeIdent(x.Name)
						file.WriteString("if (" + name + " == NULL)\n")
						file.WriteString(name + " = array_new (sizeof (extension_t), 0);\n")
						file.WriteString("  array_append (" + name + ", newel);\n")
						file.WriteString("}\n")
						file.WriteString("  }\n")
						file.WriteString("}\n")
					case Sequence:
						file.WriteString("int type = 0;\n")
						file.WriteString("  while (1) {\n")
						file.WriteString("type = xmlreader_next (reader);\n")
						file.WriteString("if (type == XML_END_ELEMENT) break;\n")
						file.WriteString("if (type == XML_ERROR) return NULL;\n")
						file.WriteString("if (type == XML_START_ELEMENT) {\n")
						for _, z := range []*Field(typ) {
							space, local := c_getSpaceAndName(target, target.Space, z)
							if !elseif {
								// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
								// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
							}
							if elseif {
								file.WriteString("  else if")
							} else {
								elseif = true
								file.WriteString(" if")
							}
							file.WriteString(" ((strcmp (space, " + space +
								") == 0) && (strcmp (name, \"" + local + "\") == 0)) {\n")
							decoder, err := c_getExtensionDecoder(space, local)
							if err != nil { return false, err }
							file.WriteString(c_referenceType(target, z) + "  newel = " + decoder + "(reader);\n")
							file.WriteString("  if (newel == NULL) {\n    return NULL;\n  }\n")
							name := prefix + sep + c_makeIdent(x.Name)
							file.WriteString("if (" + name + " == NULL)\n")
							file.WriteString(name + " = array_new (sizeof (extension_t), 0);\n")
							file.WriteString("  array_append (" + name + ", newel);\n")
							file.WriteString("  }\n")
						}
						file.WriteString("  }\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
					case Choice:
						file.WriteString("int type = 0;\n")
						file.WriteString("while (1) {\n")
						file.WriteString("type = xmlreader_next (reader);\n")
						file.WriteString("if (type == XML_END_ELEMENT) break;\n")
						file.WriteString("if (type == XML_ERROR) return NULL;\n")
						file.WriteString("if (type == XML_START_ELEMENT) {\n")
						for _, z := range []*Field(typ) {
							space, local := c_getSpaceAndName(target, target.Space, z)
							file.WriteString("if (strcmp (namespace, " + space +
								") == 0 && strcmp (name, \"" + local + "\") == 0)  {\n")
							decoder, err := c_getExtensionDecoder(space, local)
							if err != nil { return false, err }
							file.WriteString(c_referenceType(target, z) + " newel = " + decoder + "(reader);")
							file.WriteString("  if (newel == NULL) {\n    return NULL;\n  }\n")
							file.WriteString(prefix + sep + c_makeIdent(x.Name) + " = (" +
								c_referenceType(target, z) + ") newel;\n")
							file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
							file.WriteString("break;\n")
							file.WriteString("  }\n")
						}
						file.WriteString("}\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
					}
				}
				file.WriteString("  }\n")
			} else {
				switch typ := x.Type.(type) {
				case Extension:
					if typ.Local == "" {
						if elseif {
							file.WriteString("  else ")
						} else {
							elseif = true
						}
						file.WriteString("if (strcmp (namespace, ns_" + target.Name + ") != 0) {\n")
						file.WriteString("extension_t* newel = malloc (sizeof (extension_t));\n")
						file.WriteString("int err = xstream_extension_decode (reader, newel);\n")
						file.WriteString("if (reader->err != 0) return NULL;\n")
						file.WriteString("if (err == ERR_EXTENSION_NOT_FOUND) {\n")
						file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
						file.WriteString("} else {\n")
						file.WriteString(prefix + sep + c_makeIdent(x.Name) + " = newel;\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
					} else {
						typename, err := c_getExternalType(typ.Space, typ.Local)
						if err != nil { return false, err }
						was_elseif := false
						if elseif {
							file.WriteString("else ")
						}
						if !elseif {
							file.WriteString("{\n")
							was_elseif = true
							// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
							// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
						}				
						file.WriteString(" if ((strcmp (namespace, \"" + typ.Space +
							"\") == 0) && (strcmp (name, \"" + typ.Local + "\") == 0)) {\n")
						decoder, err := c_getExtensionDecoder(typ.Space, typ.Local)
						if err != nil { return false, err }
						file.WriteString(typename + " newel = " + decoder + "(reader);\n")
						file.WriteString("  if (newel == NULL) {\n    return NULL;\n  }\n")
						file.WriteString(prefix + sep + c_makeIdent(x.Name) + " = newel;\n")
						file.WriteString("}\n")
						if was_elseif {
							file.WriteString("}\n")
						}
					}
				case Set:
					fields := []*Field(typ)
					for _, z := range fields {
						space, local := c_getSpaceAndName(target, target.Space, z)
						if !elseif {
							// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
							// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
						}				
						if elseif {
							file.WriteString("  else if")
						} else {
							elseif = true
							file.WriteString(" if")
						}
						file.WriteString(" ((strcmp (namespace, " + space +
							") == 0) && (strcmp (name, \"" + local + "\") == 0)) {\n")
						c_generate_element_decoder(file, target, prefix + sep + c_makeIdent(x.Name) + "." +
							c_makeIdent(z.Name), z)
						file.WriteString("  }\n")
					}
				case SequenceOf:
					switch string(typ) {
					case "extension":
						if elseif {
							file.WriteString("else ")
						} else {
							elseif = true
						}
						file.WriteString("if (strcmp (namespace, ns_" + target.Name + ") != 0) {\n")
						file.WriteString("extension_t ext;\n")
						file.WriteString("int err = xstream_extension_decode (reader, &ext);\n")
						file.WriteString("if (reader->err != 0) return NULL;\n")
						file.WriteString("if (err == ERR_EXTENSION_NOT_FOUND) {\n")
						file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
						file.WriteString("} else{\n")
						file.WriteString("if (" + prefix + sep + c_makeIdent(x.Name) + " == NULL)\n")
						file.WriteString(prefix + sep + c_makeIdent(x.Name) +
							" = array_new (sizeof (extension_t), 0);\n")
						file.WriteString("array_append (" + prefix + sep + c_makeIdent(x.Name) +
							", &ext);\n")
						file.WriteString("}\n")
						file.WriteString("}\n")
						
					default:
						field := c_getFieldByName(target, string(typ))
						if field != nil {
							space, local := c_getSpaceAndName(target, target.Space, field)
							if !elseif {
								// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
								// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
							}				
							if elseif {
								file.WriteString("  else if")
							} else {
								elseif = true
								file.WriteString(" if")
							}
							file.WriteString(" ((strcmp (namespace, " + space +
								") == 0) && (strcmp (name, \"" + local + "\") == 0)) {\n")
							decoder, err := c_getExtensionDecoder(space, local)
							if err != nil { return false, err }
							file.WriteString(c_referenceType(target, field) +
								" newel = " + decoder + "(reader);\n")
							file.WriteString("  if (newel == NULL) {\n    return NULL;\n  }\n")
							file.WriteString("if (" + prefix + sep + c_makeIdent(x.Name) + " == NULL)\n")
							file.WriteString(prefix + sep + c_makeIdent(x.Name) + " = array_new (sizeof (extension_t), 0);\n")
							file.WriteString(
								"  array_append (" + prefix + sep + c_makeIdent(x.Name) +
									", newel);\n")
							file.WriteString("  }\n")
						} else {
							fmt.Println("dont know how to decode 111 ", typ)
						}
					}
				case Choice:
					fields := []*Field(typ)
					for _, z := range fields {
						space, local := c_getSpaceAndName(target, target.Space, z)
						if !elseif {
							// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
							// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
						}				
						if elseif {
							file.WriteString("  else if")
						} else {
							elseif = true
							file.WriteString(" if")
						}
						file.WriteString(" ((strcmp (namespace, " + space +
							") == 0) && (strcmp (name, \"" + local + "\") == 0)) {\n")
						decoder, err := c_getExtensionDecoder(space, local)
						if err != nil { return false, err }
						file.WriteString(c_referenceType(target, z) + " newel = " + decoder +
							"(reader);\n")
						file.WriteString("  if (newel == NULL) {\n    return NULL;\n  }\n")
						fname := z.Name
						if fname == "" {
							fname = z.Type.(string)
						}
						file.WriteString(prefix + "->type = EXTENSION_TYPE_" + c_uppercase(target.Name + "_" +
							c_normalize(z.Type.(string))) + ";\n")
						file.WriteString(prefix + sep + "u->" + c_makeIdent(z.Type.(string)) + " = newel;\n")
							// "(" + c_referenceType(target, x) + ") newel;\n")
						file.WriteString("  }\n")
					}
				}
			}
		}
		if any != nil {
			// if !elseif {
				// file.WriteString("    const char* namespace = xmlreader_get_namespace (reader);\n")
		// }				
			if elseif {
				file.WriteString("  else if")
			} else {
				elseif = true
				file.WriteString("  if ")
			}
			file.WriteString(" (strcmp (namespace, ns_" + target.Name + ") != 0) {\n")
			switch any.EncodingRule.Type {
			case "element:name":
				///+ enum
				typ := any.Type.(string)
				field := &Field{Type:typ}
				// file.WriteString("    const char* name = xmlreader_get_name (reader);\n")
				file.WriteString(prefix + sep + c_makeIdent(any.Name) +
					" = " + c_makeConverterFromValue(target, field, "name") + ";\n")
				file.WriteString("  if (xmlreader_skip_element (reader) == -1) return NULL;\n")
			case "name":
				file.WriteString(prefix + sep + c_makeIdent(any.Name) + " = name;\n")
			case "element":
				subfields := []*Field(any.Type.(Set))
				close, err := c_generate_element_set_decoder(file, target,
					prefix + sep + c_makeIdent(any.Name), subfields)
				if err != nil { return false, err }
				if close {
					file.WriteString("if (xmlreader_skip_element (reader) == -1) return NULL;\n")
				}
			}
			file.WriteString("      }\n")
		}
		file.WriteString ("}\n")
		file.WriteString("  }\n")
	}
	var cdata *Field
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "cdata" {
			cdata = x
			break
		}
	}
	if cdata != nil {
		c_generate_cdata_decoder(file, target, prefix + sep + c_makeIdent(cdata.Name), cdata)
	}
	if len(elems) == 0 && cdata == nil && any == nil {
		return true, nil
	}
	return false, nil
}

func c_generate_simplevalue_decoder(file *os.File, target *Target, prefix, varname string, field *Field) {
	var typ string
	switch t := field.Type.(type) {
	case SequenceOf:
		typ = string(t)
	case string:
		typ =t
	}
	switch typ {
	case "boolean":
		file.WriteString(prefix + " = strconv_parse_boolean(" + varname + ");\n")
	case "string":
		file.WriteString("  " + prefix + " = (char*) " + varname + ";\n")
	case "bytestring":
		file.WriteString(prefix + " = " + varname + ");\n")
	case "jid":
		file.WriteString("jid_t* jid = jid_of_string (" + varname + ");\n")
		file.WriteString("  " + prefix + " = jid;\n")
	case "uint":
		file.WriteString("  " + prefix + " = strconv_parse_uint (" + varname + ");\n")
	case "int":
		file.WriteString("  " + prefix + " = strconv_parse_int (" + varname + ");\n")
	case "datetime":
		file.WriteString(prefix + " = datetime_parse (" + varname + ");\n")
	case "xmllang":
		file.WriteString(prefix + " = (char*)" + varname + ";\n")
	default: // enums?
		file.WriteString(prefix + " = " + c_makeConverterFromValue(target, field, varname) + ";\n")
	}
}

func c_generate_enum_decoder(file *os.File, target *Target, field *Field) {
	t := c_makeType(target.Name, field.Name)
	v := t[:len(t)-2]
	file.WriteString("enum " + c_makeType(target.Name, field.Name) + " enum_" + v +
		"_from_string(const char* value) {\n")
	elseif := false
	enum := []string(field.Type.(Enum))
	for _, z := range enum {
		if elseif {
			file.WriteString("else ")
		}
		file.WriteString("if (strcmp (value, \"" + z + "\") == 0)\n")
		file.WriteString("return " + c_uppercase (target.Name + "_" + field.Name + "_" + z) + ";\n")
		elseif = true
	}
	file.WriteString("return 0;\n")
	file.WriteString("}\n")
}

func c_generate_enum_encoder(file *os.File, target *Target, field *Field) {
	t := c_makeType(target.Name, field.Name)
	v := t[:len(t)-2]
	file.WriteString("const char* enum_" + v + "_to_string(enum " + t + " value) {\n")
	file.WriteString("switch (value) {\n")
	enum := []string(field.Type.(Enum))
	for _, z := range enum {
		file.WriteString("case " + c_uppercase (target.Name + "_" + field.Name + "_" + z) + ":\n")
		file.WriteString("return \"" + z + "\";\n")
	}
	file.WriteString("}\n")
	file.WriteString("return NULL;\n")
	file.WriteString("}\n")
}




func c_hasChilds(target *Target, field *Field) bool {
	if field.EncodingRule != nil {
		switch field.EncodingRule.Type {
		case "element:cdata", "cdata", "name", "element:bool", "element:name", "attribute":
			return false
		case "element":
			switch typ := field.Type.(type) {
			case string:
				field := c_getFieldByName(target, string(typ))
				if field != nil { return c_hasReallyChilds(target, field) }
			case Set:
				fields := []*Field(typ)
				for _, x := range fields {
					if c_hasReallyChilds(target, x) { return true }
				}
			case SequenceOf:
				field := c_getFieldByName(target, string(typ))
				if field != nil { return c_hasReallyChilds(target, field) }
			case Sequence:
				fields := []*Field(typ)
				for _, x := range fields {
					if c_hasReallyChilds(target, x) { return true }
				}
			case Choice:
				fields := []*Field(typ)
				for _, x := range fields {
					if c_hasReallyChilds(target, x) { return true }
				}
			case Extension:
				return false
			}
		}
	}
	return false
}

func c_hasReallyChilds(target *Target, field *Field) bool {
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
			field := c_getFieldByName(target, typ)
			if field != nil { return c_hasReallyChilds(target, field)}
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				if c_hasReallyChilds(target, x) { return true }
			}
		case SequenceOf:
			field := c_getFieldByName(target, string(typ))
			if field != nil {
				return c_hasReallyChilds(target, field)
			}
		case Sequence:
			fields := []*Field(typ)
			for _, x := range fields {
				if c_hasReallyChilds(target, x) { return true }
			}			
		case Choice:
			fields := []*Field(typ)
			for _, x := range fields {
				if c_hasReallyChilds(target, x) { return true }
			}			
		case Extension:
			return false
		}
	}
	return false
	
}

func c_generate_simplevalue_encoder(name string, field *Field) string {
	isarray := false
	var typ string
	switch t := field.Type.(type) {
	case Enum:
		
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
				return name
			} else {
				return "(char*)" + name
			}
		}
		if name == "elm" {
			return name
		} else {
			return name
		}
	case "boolean":
		return "strconv_format_boolean(" + name + ")"
	case "jid":
		return "jid_to_string(" + name + ")"
	case "uint":
		return "strconv_format_uint(" + name + ")"
	case "int":
		return "strconv_format_int(" + name + ")"
	case "datetime":
		return "datetime_to_string(" + name + ")"
	}
	return name
}

func c_is_enum (target *Target, field *Field) bool {
	switch typ := field.Type.(type) {
	case string:
		for _, f := range target.Fields {
			if typ == f.Name {
				if _, ok := f.Type.(Enum); ok {
					return true
				}
			}
		}
	}
	return false
}

func c_generate_check_condition(file *os.File, name string, target *Target, field *Field) bool {
	switch typ := field.Type.(type) {
	case string:
		if typ == "boolean" {
			file.WriteString("if (" + name + ") {\n")
			return true
		}
		for _, f := range target.Fields {
			if typ == f.Name {
				if _, ok := f.Type.(Enum); ok {
					file.WriteString("if (" + name + " != 0) {\n")
					return true
				}
			}
		}
	case Sequence, SequenceOf, Set:
		return false
		//	case string, Extension:
	}
	file.WriteString("if (" + name + " != NULL) {\n")
	return true
}

func c_getElementName(fields []*Field) *Field {
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "name" {
			return x
		}
	}
	return nil
}

func c_generate_extensions(filename string) {
	fmt.Print("Generating extensions file\n")
	file, err := os.Create(filename)
	if err != nil { fmt.Println(err) }
	defer file.Close()
	file.WriteString("#include \"extensions.h\"\n")
	file.WriteString("#include \"xstream.h\"\n\n")
	for _, schema := range schemas {
		extension := ""
		if category, ok := schema.Props["category"]; ok {
			if category == "extension" {
				extension = "xep_"
			}
		}
		filename := extension + schema.PackageName + "_data.h"
		file.WriteString("#include \"" + filename + "\"\n")
	}
	file.WriteString("\n")
	extensions_len := 0
	file.WriteString("struct xstream_extension_t extensions[] = {\n")
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			for _, field := range target.Fields {
				if (field.EncodingRule != nil && (field.EncodingRule.Type == "element" ||
					field.EncodingRule.Type == "startelement")) || len(target.Fields) == 1 {
					_, local := c_getSpaceAndName(target, target.Space, field)
					name := target.Name + "_" + c_normalize(field.Name)
					file.WriteString(" {\"" + target.Space +
						"\", \"" + local + "\", EXTENSION_TYPE_" + c_uppercase(name) +
						", (void *(*)(xmlreader_t*)) " +
						name + "_decode, (int (*)(xmlwriter_t*, void*)) " +
						name + "_encode, (void (*)(void*)) " + name + "_free},\n")
					extensions_len++
				}
			}
		}
	}
	file.WriteString("};\n\n")
	file.WriteString("int extensions_len = " + strconv.FormatInt(int64(extensions_len), 10) + ";\n")
}

func c_generate_extensions_types(filename string) {
	fmt.Print("Generating extensions file\n")
	file, err := os.Create(filename)
	if err != nil { fmt.Println(err) }
	defer file.Close()
	file.WriteString("#ifndef _EXTENSIONS_H_\n")
	file.WriteString("#define _EXTENSIONS_H_\n\n")
	file.WriteString("enum extension_type {\n")
	
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			for _, field := range target.Fields {
				if (field.EncodingRule != nil && (field.EncodingRule.Type == "element" ||
					field.EncodingRule.Type == "startelement" ||
					field.EncodingRule.Type == "element:cdata")) || len(target.Fields) == 1 {
					name := target.Name + "_" + c_normalize(field.Name)
					file.WriteString("  EXTENSION_TYPE_" + c_uppercase(name) + ",\n")
				}
			}
		}
	}
	file.WriteString("};\n\n")
	file.WriteString("#endif\n")
}


		
func c_generate_element_encoder(file *os.File, target *Target, prefix string,
	element *Field) error {
	if element.EncodingRule == nil {
		switch typ := element.Type.(type) {
		case Extension:
			if typ.Local == "" {
				file.WriteString("err = xstream_extension_encode (writer, " + prefix +
					"->data, " + prefix + "->type);\n")
				file.WriteString("if (err != 0) return err;\n")
			} else {
				encoder, err := c_getExtensionEncoder(typ.Space, typ.Local)
				if err != nil {return err}
				file.WriteString("err = " + encoder + "(writer, " + prefix + ");\n")
				file.WriteString("if (err != 0) return err;\n")
			}
		case string:
			file.WriteString("err = " + target.Name + "_" + c_normalize(string(typ)) +
				"_encode(writer, " + prefix + ");\n")
			file.WriteString("if (err != 0) return err;\n")
		case SequenceOf:
			file.WriteString("{\n")
			arrname := prefix
			file.WriteString("int len = array_length (" + arrname + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			file.WriteString("extension_t* ext = array_get (" + arrname + ", i);\n")
			if string(typ) == "extension" {
				file.WriteString("err = xstream_extension_encode(writer, ext->data, ext->type);\n")
				file.WriteString("if (err != 0) return err;\n")
			} else {
				encoder, err := c_getExtensionEncoder(target.Space, string(typ))
				if err != nil { return err }
				file.WriteString("err = " + encoder + " (writer, ext->data);\n")
				file.WriteString("if (err != 0) return err;\n")
			}
			file.WriteString("}\n")
			file.WriteString("}\n")
		case Sequence:
			file.WriteString("{\n")
			arrname := prefix
			file.WriteString("int len = array_length (" + arrname + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			file.WriteString("extension_t* ext = array_get (" + arrname + ", i);\n")
			file.WriteString("err = xstream_extension_encode (writer, ext->data, ext->type);\n")
			file.WriteString("if (err != 0) return -1;\n")
			file.WriteString("}\n")
			file.WriteString("}\n")
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				close := c_generate_check_condition(file, prefix + "." + c_makeIdent(x.Name), target, x)
				c_generate_element_encoder(file, target, prefix + "." + c_makeIdent(x.Name), x)
				if close {
					file.WriteString("}\n")
				}
			}
		case Choice:
			file.WriteString("if (" + prefix + " != NULL) {\n")
			file.WriteString("if err = " + prefix + ".(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			file.WriteString("}\n")
		default:
			fmt.Println("dont know what to do ", element.Name, " ", element.Type)
		}
		return nil
	}
	fieldname := c_makeIdent(element.Name)
	space, local := c_getSpaceAndName(target, target.Space, element)
	switch element.EncodingRule.Type {
	case "element:name":
		file.WriteString("const char* name = " + c_makeConverterToValue(target, element, prefix) +
			";\n")
		file.WriteString("err = xmlwriter_simple_element (writer, " + space + ", name, NULL);\n")
		file.WriteString("if (err != 0) return err;\n")			
	case "element:cdata":
		varname := prefix
		isarray := false
		if _, ok := element.Type.(SequenceOf); ok {
			varname = "value"
			isarray = true
		}
		value := c_makeConverterToValue(target, element, varname)
		if isarray {
			arrname := prefix
			file.WriteString("int len = array_length (" + arrname + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			file.WriteString("char** value = array_get (" + arrname + ", i);\n")
		} 
		file.WriteString("err = xmlwriter_simple_element (writer, " + space + ", \"" + local + "\", " +
			value + ");\n")
		file.WriteString("if (err != 0) return err;\n")
		if isarray {
			file.WriteString("}\n")
		}
	case "element:bool":
		file.WriteString("err = xmlwriter_simple_element (writer, " + space + ", \"" + local +
			"\", NULL);\n")
		file.WriteString("if (err != 0) return err;\n")
	case "startelement", "element":
		switch typ := element.Type.(type) {
		case Extension:
			space, local := c_getSpaceAndName(target, target.Space, element)
			file.WriteString("err = xmlwriter_start_element (writer, " + space + ", \"" + local +
				"\");\n")
			file.WriteString("if (err != 0) return err;\n")
			if typ.Local != "" {
				encoder, err := c_getExtensionEncoder(typ.Space, typ.Local)
				if err != nil { return err }
				file.WriteString("err = " + encoder + "(writer, " + prefix + ");\n")
				file.WriteString("if (err != 0) return err;\n")
			} else {
				file.WriteString("todo();\n")
			}
			file.WriteString("err = xmlwriter_end_element(writer);\n")
		file.WriteString("if (err != 0) return err;\n")
		case Set:
			sep := "->"
			if prefix != "elm" {
				sep = "."
			}
			fields := []*Field(typ)
			specialFieldName := c_getElementName(fields)
			if specialFieldName != nil {
				fieldname = c_makeIdent(specialFieldName.Name)
				file.WriteString("const char* name = " +
					c_makeConverterToValue(target, specialFieldName, prefix + sep + fieldname) + ";\n")
				file.WriteString("err = xmlwriter_start_element (writer, " + space + ", name);\n")
		file.WriteString("if (err != 0) return err;\n")
			} else {
				file.WriteString("err = xmlwriter_start_element (writer, " + space + ", \"" +
					local + "\");\n")
				file.WriteString("if (err != 0) return err;\n")
			}
			for _, x := range fields {
				if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
					aspace, alocal := c_getSpaceAndName(target, "", x)
					if aspace == "ns_" + target.Name {
						aspace = "NULL"
					}
					close := c_generate_check_condition(file, prefix + sep + c_makeIdent(x.Name), target, x)
					value := c_makeConverterToValue(target, x, prefix + sep + c_makeIdent(x.Name))
					file.WriteString("err = xmlwriter_attribute (writer, " + aspace + ", \"" + alocal +
						"\", " + value + ");\n")
					file.WriteString("if (err != 0) return err;\n")
					if close {
						file.WriteString("}\n")
					}
				}
			}
			for _, x := range fields {
				if x.EncodingRule != nil &&
					(x.EncodingRule.Type == "name" ||
					x.EncodingRule.Type == "cdata" || x.EncodingRule.Type == "attribute") {
					continue
				}
				name := prefix
				if _, ok := x.Type.(Choice); ok {
					name += "->u"
					close := c_generate_check_condition(file, name, target, x)
					file.WriteString("err = xstream_extension_encode(writer, (void*)" + name + ", " + prefix +
						"->type);\n")
					file.WriteString("if (err != 0) return err;\n")
					if close {
						file.WriteString("}\n")
					}
				} else {
					name += sep + c_makeIdent(x.Name)
					close := c_generate_check_condition(file, name, target, x)
					c_generate_element_encoder(file, target, name, x)
					if close {
					file.WriteString("}\n")
					}
				}
			}
			for _, x := range fields {
				if x.EncodingRule != nil && x.EncodingRule.Type == "cdata" {
					file.WriteString("if (" + prefix + sep + c_makeIdent(x.Name) + " != NULL) {\n")
					typ := x.Type.(string)
					if typ == "bytestring" {
						file.WriteString("err = xmlwriter_base64 (writer, " +
							prefix + sep + c_makeIdent(x.Name) + ", " +
							prefix + sep + c_makeIdent(x.Name) + "_len);\n")
						file.WriteString("if (err != 0) return err;\n")
					} else {
						value := c_generate_simplevalue_encoder(prefix + sep + c_makeIdent(x.Name), x)
						file.WriteString("err = xmlwriter_text (writer, " + value + ");\n")
						file.WriteString("if (err != 0) return err;\n")
					}
					file.WriteString("}\n")
				}
			}
			if element.EncodingRule.Type == "element" {
				file.WriteString("err = xmlwriter_end_element(writer);\n")
				file.WriteString("if (err != 0) return err;\n")
			}
		case string:
			switch typ {
			case "langstring":
				file.WriteString("err = langstring_encode(writer, " + space + ", \"" + local + "\", " +
					prefix + ");\n")
				file.WriteString("if (err != 0) return err;\n")
			case "extension":
				file.WriteString(prefix + ".(xmlencoder.Extension).Encode(e)\n")
			default:
				file.WriteString(prefix + ".Encode(e)\n")
			}
		case Choice:
			file.WriteString("err = xmlwriter_start_element (writer, " + space + ", \"" +
				local + "\");\n")
			file.WriteString("if (err != 0) return err;\n")
			if prefix == "elm" {
				file.WriteString("if err = elm.Payload.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			} else {
				file.WriteString("if err = " + prefix + ".(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			}
			file.WriteString("err = xmlwriter_end_element(writer);\n")
			file.WriteString("if (err != 0) return err;\n")
		case SequenceOf:
			file.WriteString("err = xmlwriter_start_element (writer, " + space + ", \"" + local +
				"\");\n")
			file.WriteString("if (err != 0) return err;\n")
			arrname := prefix
			file.WriteString("int len = array_length (" + arrname + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			if typ == "extension" {
				file.WriteString("extension_t* ext = array_get (" + arrname + ", i);\n")
				file.WriteString("err = xstream_extension_encode (writer, ext->data, ext->type);\n")
				file.WriteString("if (err != 0) return err;\n")
			} else {
				field := c_getFieldByName(target, string(typ))
				if field == nil { return errors.New("unknown type " + string(typ)) }
				space, local := c_getSpaceAndName(target, target.Space, field)
				encoder, err := c_getExtensionEncoder(space, local)
				if err != nil { return err }
				file.WriteString(c_referenceType(target, field) + " data = array_get (" + arrname + ", i);\n")
				file.WriteString("err = " + encoder + "(writer, data);\n")
				file.WriteString("if (err != 0) return err;\n")
			}
			file.WriteString("}\n")
			file.WriteString("err = xmlwriter_end_element(writer);\n")
			file.WriteString("if (err != 0) return err;\n")
		case Sequence:
			file.WriteString("err = xmlwriter_start_element (writer, " + space + ", \"" + local +
				"\");\n")
			file.WriteString("if (err != 0) return err;\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i; i < " + prefix + "->length; i++) {\n")
			file.WriteString("if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			file.WriteString("}\n")
			file.WriteString("err = xmlwriter_end_element(writer);\n")
			file.WriteString("if (err != 0) return err;\n")
		}
	}
	return nil
}

func c_generate_free(file *os.File, target *Target, prefix string, element *Field) error {
	if c_is_enum(target, element) {
		return nil
	}
	if element.EncodingRule == nil {
		switch typ := element.Type.(type) {
		case Extension:
			if typ.Local == "" {
				file.WriteString("xstream_extension_free (" + prefix + ");\n")
				file.WriteString("free (" + prefix + ");\n")
			} else {
				freer, err := c_getExtensionFree(typ.Space, typ.Local)
				if err != nil {return err}
				file.WriteString(freer + "(" + prefix + ");\n")
			}
		case string:
			freer, err := c_getExtensionFree(target.Name, string(typ))
			if err != nil { return err }
			file.WriteString(freer + "(" + prefix + ");\n")
		case SequenceOf:
			file.WriteString("{\n")
			arrname := prefix
			file.WriteString("int len = array_length (" + arrname + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			if string(typ) == "extension" {
				file.WriteString("extension_t* ext = array_get (" + arrname + ", i);\n")
				file.WriteString("xstream_extension_free (ext);\n")
			} else {
				freer, err := c_getExtensionFree(target.Space, string(typ))
				if err != nil { return err }
				file.WriteString(freer + " (array_get (" + prefix + ", i));\n")
			}
			file.WriteString("}\n")
			file.WriteString("array_free (" + arrname + ");\n")
			file.WriteString("}\n")
		case Sequence:
			file.WriteString("{\n")
			arrname := prefix
			file.WriteString("int len = array_length (" + arrname + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			file.WriteString("extension_t* ext = array_get (" + arrname + ", i);\n")
			file.WriteString("xstream_extension_free (ext);\n")
			file.WriteString("}\n")
			file.WriteString("array_free (" + arrname + ");\n")
			file.WriteString("}\n")
		case Set:
			fields := []*Field(typ)
			sep := "."
			if prefix == "data" {
				sep = "->"
			}
			for _, x := range fields {
				close := c_generate_check_condition(file, prefix + sep + c_makeIdent(x.Name), target, x)
				c_generate_free(file, target, prefix + sep + c_makeIdent(x.Name), x)
				if close {
					file.WriteString("}\n")
				}
			}
		case Choice:
			file.WriteString("if (" + prefix + " != NULL) {\n")
			file.WriteString("if err = " + prefix + ".(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			file.WriteString("}\n")
		default:
			fmt.Println("dont know what to do ", element.Name, " ", element.Type)
		}
		return nil
	}
	// fieldname := c_makeIdent(element.Name)
	// space, local := c_getSpaceAndName(target, target.Space, element)
	switch element.EncodingRule.Type {
	case "element:name":
	case "element:cdata":
		funcname := "free"
		var typ string
		switch tt := element.Type.(type) {
		case string:
			typ = tt
		case SequenceOf:
			typ = string(tt)
		}
		switch typ {
		case "jid": funcname = "jid_free"
		}
		isarray := false
		if _, ok := element.Type.(SequenceOf); ok {
			isarray = true
		}
		if isarray {
			arrname := prefix
			file.WriteString("int len = array_length (" + arrname + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			file.WriteString("char** value = array_get (" + arrname + ", i);\n")
			file.WriteString(funcname +"(*value);\n")
			file.WriteString("}\n")
			file.WriteString("array_free (" + arrname + ");\n")
		} else {
			file.WriteString(funcname + " (" + prefix + ");\n")
		}
	case "element:bool":
	case "startelement", "element":
		switch typ := element.Type.(type) {
		case Extension:
			// space, local := c_getSpaceAndName(target, target.Space, element)
			if typ.Local != "" {
				freer, err := c_getExtensionFree(typ.Space, typ.Local)
				if err != nil { return err }
				file.WriteString(freer + "(" + prefix + ");\n")
			} else {
				file.WriteString("todo();\n")
			}
		case Set:
			sep := "->"
			if prefix != "data" {
				sep = "."
			}
			fields := []*Field(typ)
			for _, x := range fields {
				if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
					if c_is_enum(target, x) {
						continue;
					}
					close := c_generate_check_condition(file, prefix + sep + c_makeIdent(x.Name), target, x)
					varname := prefix + sep + c_makeIdent(x.Name)
					switch typ := x.Type.(type) {
					case string:
						switch typ {
						case "string", "boolean", "uint", "int", "datetime", "xmllang":
							file.WriteString("free (" + varname + ");\n")
						case "jid":
							file.WriteString("jid_free (" + varname + ");\n")
						// default:
							// file.WriteString("free (" + varname + ");\n")
						}
					}
					if close {
						file.WriteString("}\n")
					}
				}
			}
			for _, x := range fields {
				if x.EncodingRule != nil &&
					(x.EncodingRule.Type == "name" ||
					x.EncodingRule.Type == "cdata" || x.EncodingRule.Type == "attribute") {
					continue
				}
				name := prefix
				if _, ok := x.Type.(Choice); ok {
					name += "->u"
					close := c_generate_check_condition(file, name, target, x)
					file.WriteString("extension_t ext = {" + prefix + "->type, " + name + "};\n")
					file.WriteString("xstream_extension_free (&ext);\n")
					if close {
						file.WriteString("}\n")
					}
				} else {
					name += sep + c_makeIdent(x.Name)
					close := c_generate_check_condition(file, name, target, x)
					c_generate_free(file, target, name, x)
					if close {
					file.WriteString("}\n")
					}
				}
			}
			for _, x := range fields {
				if x.EncodingRule != nil && x.EncodingRule.Type == "cdata" {
					file.WriteString("if (" + prefix + sep + c_makeIdent(x.Name) + " != NULL) {\n")
					typ := x.Type.(string)
					if typ == "bytestring" {
						file.WriteString("free (" + prefix + sep + c_makeIdent(x.Name) + ");\n")
					} else {
						value := c_generate_simplevalue_encoder(prefix + sep + c_makeIdent(x.Name), x)
						file.WriteString("free (" + value + ");\n")
					}
					file.WriteString("}\n")
				}
			}
			if element.EncodingRule.Type == "element" {
			}
		case string:
			switch typ {
			case "langstring":
				file.WriteString("langstring_free (" + prefix + ");\n")
			case "extension":
				file.WriteString(prefix + ".(xmlencoder.Extension).Encode(e)\n")
			default:
				file.WriteString(prefix + ".Encode(e)\n")
			}
		case Choice:
			if prefix == "data" {
				file.WriteString("if err = elm.Payload.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			} else {
				file.WriteString("if err = " + prefix + ".(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			}
		case SequenceOf:
			file.WriteString("int len = array_length (" + prefix + ");\n")
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i = 0; i < len; i++) {\n")
			if typ == "extension" {
				file.WriteString("extension_t* ext = array_get (" + prefix + ", i);\n")
				file.WriteString("xstream_extension_free (ext);\n")
			} else {
				field := c_getFieldByName(target, string(typ))
				if field == nil { return errors.New("unknown type " + string(typ)) }
				space, local := c_getSpaceAndName(target, target.Space, field)
				freer, err := c_getExtensionFree(space, local)
				if err != nil { return err }
//++
				file.WriteString(c_referenceType(target, field) + " item = array_get (" + prefix + ", i);\n")
				file.WriteString(freer + "(item);\n")
			}
			file.WriteString("}\n")
			file.WriteString("array_free (" + prefix + ");\n")
		case Sequence:
			file.WriteString("int i = 0;\n")
			file.WriteString("for (i; i < " + prefix + "->length; i++) {\n")
			file.WriteString("if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }\n")
			file.WriteString("}\n")
		}
	}
	return nil
}
