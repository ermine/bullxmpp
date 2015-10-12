package main

import (
	"os"
	"path/filepath"
	"unicode"
	"fmt"
	"errors"
)

func KotlinGenerate() error {
	if cfg.Kotlin.Outdir == "" {
		panic("no outdir")
	}
	for _, schema := range schemas {
		dir := cfg.Kotlin.Outdir
		if category, ok := schema.Props["category"]; ok {
			if category == "extension" {
				dir = filepath.Join(cfg.Kotlin.Outdir, "extensions")
			}
		}
		dir = filepath.Join(dir, schema.PackageName)
		fmt.Println("making directory " , dir)
		os.MkdirAll(dir, 0755)

		for _, target := range schema.Targets {
			if name, ok := target.Props["name"]; ok {
				target.Name = name
			}
				}
		for _, target := range schema.Targets {
			var filename string
			if target.Name != "" {
				filename = filepath.Join(dir, schema.PackageName + "_" + target.Name + "_data.kt")
			} else {
				filename = filepath.Join(dir, schema.PackageName + "_data.kt")
			}					
			file, err := os.Create(filename)
			if err != nil { return err }
			defer file.Close()
			kotlin_generate_package(file, schema, target)
		}
	}

	kotlin_generate_adders()
	return nil
}

func kotlin_generate_package(file *os.File, schema *Schema, target *Target) error {
	file.WriteString("package " + cfg.Kotlin.Package_prefix + ".")
	if category, ok := schema.Props["category"]; ok && category == "extension" {
		file.WriteString("extensions.")
	}
	file.WriteString(schema.PackageName)
	if target.Name != "" {
		file.WriteString("." + target.Name)
	}
	file.WriteString("\n\n")
	
	var imports []string
	for _, field := range target.Fields {
		kotlin_getImports(target, field, &imports)
	}
	file.WriteString("import org.xmlpull.v1.XmlPullParser\n")
	// file.WriteString("import org.xmlpull.v1.XmlPullParserException\n")
	file.WriteString("import org.xmlpull.v1.XmlSerializer\n")
	file.WriteString("import " + cfg.Kotlin.Package_prefix + ".xmlencoder.XmlEncoder\n")
	for _, i := range imports {
		file.WriteString("import " + i + "\n")
	}
	file.WriteString("\n")

	file.WriteString("val " + kotlin_ns(target) + " = \"" + target.Space + "\"\n\n")
	for _, field := range target.Fields {
		switch typ := field.Type.(type) {
		case Enum:
			kotlin_generate_enum(file, typ, kotlin_makeEnumName(field.Name))
		default:
			kotlin_generate_class(file, schema, target, field)
		}
		file.WriteString("\n")
	}
	return nil
}

func kotlin_generate_class(file *os.File, schema *Schema, target *Target, field *Field) {
	file.WriteString("class " + kotlin_makeClassName(field.Name) + "() : XmlEncoder {\n")
	kotlin_generate_class_fields(file, target, field)
	file.WriteString("\n")
	// file.WriteString("  fun " + kotlin_makeClassName(field.Name) + "() {}\n")
	// file.WriteString("\n\n")
	kotlin_generate_encoder(file, target, field)
	kotlin_generate_decoder(file, target, field)
	file.WriteString("}\n")
}

func kotlin_normalize(s string) string {
	var r []rune
	first := false
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

func kotlin_makeUppercase(s string) string {
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

func kotlin_makeClassName(s string) string {
	return kotlin_makeUppercase(s)
}

func kotlin_makeIdent(prefix, s string) string {
	switch s {
	case "var": s = "var_"
	case "continue": s = "continue_"
	}
	return prefix + kotlin_normalize (s)
}

func kotlin_makeEnumName(s string) string {
	return kotlin_makeUppercase(s) + "Enum"
}

func kotlin_getImports(target *Target, field *Field, imports *[]string) {
	switch typ := field.Type.(type) {
	case string:
		if data, ok := kotlinSimpleTypes[typ]; ok {
			if data.Import != "" {
				appendImport(imports, data.Import)
			}
		}
	case SequenceOf:
		appendImport(imports, "java.util.ArrayList")
		if data, ok := kotlinSimpleTypes[string(typ)]; ok {
			if data.Import != "" {
				appendImport(imports, data.Import)
			}
		}
		/*
	case Extension:
		if typ.Local != "" {
		Done:
			for _, s := range schemas {
				for _, t := range s.Targets {
					if t.Space == typ.Space {
						for _, f := range t.Fields {
							if (f.EncodingRule != nil && f.EncodingRule.Name == typ.Local) ||
								(f.EncodingRule != nil && f.EncodingRule.Name == "" && f.Name == typ.Local) ||
								(f.EncodingRule == nil && f.Name == typ.Local) {
									pName := cfg.Kotlin.Package_prefix + "." + s.PackageName
								if t.Name != "" {
									pName += "." + t.Name
								}
								pName += "." + kotlin_makeClassName(f.Name)
								appendImport(imports, pName)
								break Done
							}
						}
					}
				}
			}
		}
*/
	case Sequence:
		appendImport(imports, "java.util.*")
		for _, x := range []*Field(typ) {
				kotlin_getImports(target, x, imports)
		}
	case Choice:
		for _, x := range []*Field(typ) {
			kotlin_getImports(target, x, imports)
		}			
	case Set:
		for _, x := range []*Field(typ) {
			kotlin_getImports(target, x, imports)
		}
	case Enum:
	}
}

/*
func appendImport(imports *[]string, i string) {
	found := false
	for _, x := range *imports {
		if x == i {
			found = true
			break
		}
	}
	if !found {
		*imports = append(*imports, i)
	}
}		
*/

var kotlinSimpleTypes = map[string]struct{
	Type string
	Import string
}{
	"boolean": {"Boolean", ""},
	"string": {"String", ""},
	"bytestring": {"String", ""},
	//	"jid": {"JID", cfg.Kotlin.Package_prefix + ".jid.JID"},
	"jid": {"JID", "ru.sulci.kotlin.xmpp.jid.JID"},
	"datetime": {"Date", "java.util.Date"},
	"int": {"Int", ""},
	"uint": {"Int", ""},
	"xmllang": {"String", ""},
//	"langstring": {"LangString", cfg.Kotlin.Package_prefix + ".xmlencoder.LangString"},
	"langstring": {"LangString", "ru.sulci.kotlin.xmpp.xmlencoder.LangString"},
	"extension": {"XmlEncoder", ""},
}

func kotlin_generate_class_fields(file *os.File, target *Target, field *Field) {
	prefix := ""
	switch typ := field.Type.(type) {
	case string:
		if t, ok := kotlinSimpleTypes[typ]; ok {
			file.WriteString("  var " + kotlin_makeIdent(prefix, field.Name) + ": " + t.Type +
				"? = null\n")
		} else {
			name := field.Name
			if name == "" {
				name = typ
			}
			file.WriteString("  var " + kotlin_makeIdent(prefix, name) + ": " +
				kotlin_makeClassName(typ) + "? = null")
		}
	case Extension:
		if typ.Local == "" {
			file.WriteString("var payload: XmlEncoder? = null\n")
		} else {
			t := kotlin_getExtensionType(typ.Space, typ.Local)
			file.WriteString("  var payload: " + t + "? = null\n")
		}
	case SequenceOf:
		t := string(typ)
		tt := "XmlEncoder"
		if t != "extension" {
			if s, ok := kotlinSimpleTypes[t]; ok {
				tt = s.Type
			} else {
				tt = kotlin_makeClassName(t)
			}
		}
		file.WriteString("  var payloadSequence: ArrayList<" + tt + ">? = null\n")
	case Sequence:
		file.WriteString("  var payloadSequence: ArrayList<XmlEncoder>? = null\n")
	case Choice:
		file.WriteString("  var payload: XmlEncoder? = null\n")
	case Set:
		fields := []*Field(typ)
		for _, x := range fields {
			kotlin_generate_field(file, "  ", target, "", x)
		}
	case Enum:
	}
}

func kotlin_generate_field(file *os.File, ident string, target *Target, prefix string, field *Field) {
	switch typ := field.Type.(type) {
	case string:
		if t, ok := kotlinSimpleTypes[typ]; ok {
			file.WriteString(ident + "var " + kotlin_makeIdent(prefix, field.Name) + " : " + t.Type +
				"? = null\n")
		} else {
			name := field.Name
			if name == "" {
				name = typ
			}
			is_enum := false
			for _, f := range target.Fields {
				if typ == f.Name {
					if _, ok := f.Type.(Enum); ok {
						is_enum = true
						break
					}
				}
			}
			if is_enum {
				file.WriteString(ident + "var " + kotlin_makeIdent(prefix, name) + ": " +
					kotlin_makeClassName(typ) + "Enum? = null\n")
			} else {
				file.WriteString(ident + "var " + kotlin_makeIdent(prefix, name) + ": " +
					kotlin_makeClassName(typ) + "? = null\n")
			}
		}
	case Extension:
		var t string
		if typ.Local == "" {
			t = "XmlEncoder"
		} else {
			t = kotlin_getExtensionType(typ.Space, typ.Local)
		}			
		file.WriteString(ident + "var "+ kotlin_makeIdent(prefix, field.Name) + ": " + t + "? = null\n")
	case SequenceOf:
		t := string(typ)
		tt := "XmlEncoder"
		if t != "extension" {
			if s, ok := kotlinSimpleTypes[t]; ok {
				tt = s.Type
			} else {
				tt = kotlin_makeClassName(t)
			}
		}
		file.WriteString(ident + "var " + kotlin_makeIdent(prefix, field.Name) +
			": ArrayList<" + tt + ">? = null\n")
	case Sequence:
		file.WriteString(ident + "var " + kotlin_makeIdent(prefix, field.Name) +
			": ArrayList<XmlEncoder>? = null\n")
		fields := []*Field(typ)
		for _, x := range fields {
			if x.Name != "" {
				file.WriteString("  class " + kotlin_makeClassName(x.Name) + " : XmlEncoder {\n")
				kotlin_generate_class_fields(file, target, x)
				kotlin_generate_encoder(file, target, x)
				kotlin_generate_decoder(file, target, x)
				file.WriteString("  }\n")
			}
		}
	case Enum:
		file.WriteString(ident + "var "  + kotlin_makeIdent(prefix, field.Name) + ": " +
			kotlin_makeEnumName(field.Name) + "? = null\n\n")
		kotlin_generate_enum(file, typ, kotlin_makeEnumName(field.Name))
	case Choice:
		file.WriteString(ident + "XmlEncoder " + kotlin_makeIdent(prefix, field.Name) + "\n")
	case Set:
		fields := []*Field(typ)
		for _, x := range fields {
			kotlin_generate_field(file, ident, target, kotlin_makeClassName(field.Name), x)
		}
	}
}

func kotlin_generate_enum(file *os.File, enum Enum, name string) {
	file.WriteString("  enum class " + name + " {\n")
	variants := []string(enum)
	str, hashyphen := kotlin_make_enum_string(variants[0])
	file.WriteString("    " + str)
	for _, x := range variants[1:] {
		file.WriteString(",\n")
		str, h := kotlin_make_enum_string(x)
		if h {
			hashyphen = true
		}
		file.WriteString("    " + str)
	}
	file.WriteString(";\n\n")
	file.WriteString("    override fun toString() : String {\n")
	if hashyphen {
		file.WriteString("      return this.name().toLowerCase().replace('_', '-')\n")
	} else {
		file.WriteString("      return this.name().toLowerCase()\n")
	}
	file.WriteString("    }\n\n")
	file.WriteString("  }\n\n")
		
	file.WriteString("  fun " + name + ".valueOf(value: String) : " + name + " {\n")
	if hashyphen {
		file.WriteString("      return valueOf(value.toUpperCase().replace('-', '_'))\n")
	} else {
		file.WriteString("      return valueOf(value.toUpperCase())\n")
	}
	file.WriteString("    }\n")
	
		
}

func kotlin_make_enum_string(s string) (string, bool) {
	hashyphen := false
	var r []rune
	for _, x := range s {
		if x == rune('-') {
			hashyphen = true
			r = append(r, '_')
		} else {
			r = append(r, unicode.ToUpper(x))
		}
	}
	return string(r), hashyphen
}

func kotlin_resulveType(typ string) string {
	if data, ok := kotlinSimpleTypes[typ]; ok {
		return data.Type
	} else {
		if typ == "extension" {
			return "XmlEncoder"
		}
	}
	return kotlin_makeClassName(typ)
}

func kotlin_getFullClassName(space, local string) (string, error) {
	for _, schema := range schemas {
		for _, target := range schema.Targets {
			if target.Space == space || kotlin_ns(target) == space {
				for _, field := range target.Fields {
					_, alocal := c_getSpaceAndName(target, target.Space, field)
					if  alocal == local {
						pname := cfg.Kotlin.Package_prefix + "."
						if category, ok := schema.Props["category"]; ok && category == "extension" {
							pname += "extensions."
						}
						pname += schema.PackageName + "." +  kotlin_makeClassName(field.Name)
						return pname, nil
					}
				}
			}
		}
	}
	return "", errors.New("extenal type for " + space + " " + local + " not found")
}	
	




func kotlin_getExtensionType(space, local string) string {
	for _, s := range schemas {
		for _, t := range s.Targets {
			if t.Space == space {
				fmt.Println("found namespace ", space, " ", local)
				for _, f := range t.Fields {
					fmt.Println("checking ", f.Name, " ", f.EncodingRule)
					if (f.EncodingRule != nil && f.EncodingRule.Name == local) ||
						(f.EncodingRule != nil && f.EncodingRule.Name == "" && f.Name == local) ||
						(f.EncodingRule == nil && f.Name == local) {
						pkgName := cfg.Kotlin.Package_prefix + "."
						if category, ok := s.Props["category"]; ok && category == "extension" {
							pkgName += "extensions."
						}
						pkgName += s.PackageName
						return pkgName + "." + kotlin_makeClassName(f.Name)
					}
				}
			}
		}
	}
	fmt.Println("not found ", space, " ", local)
	return "UnknownNamespace" + space + local
}

func kotlin_generate_encoder(file *os.File, target *Target, field *Field) {
	file.WriteString("  override fun encode(xs: XmlSerializer) {\n")
//	if target.Prefix != "" {
	file.WriteString("    xs.setPrefix(\"" + target.Prefix + "\", " + kotlin_ns(target) + ")\n")
//	}
	kotlin_generate_class_encoder(file, target, field)
	file.WriteString("  }\n\n")
}

func kotlin_generate_class_encoder(file *os.File, target *Target, field *Field) {
	prefix := ""
	if field.EncodingRule != nil {
		space := kotlin_ns(target)
		if field.EncodingRule.Space != "" {
			space = "\"" + field.EncodingRule.Space + "\""
		}
		local := field.Name
		if field.EncodingRule.Name != "" {
			local = field.EncodingRule.Name
		}
		local = "\"" + local + "\""
		switch field.EncodingRule.Type {
		case "element:bool":
			file.WriteString("    xs.startTag(" + space + ", " + local + ")\n")
			file.WriteString("    xs.endTag(" + space + ", " + local + ")\n")
		case "element:cdata":
			file.WriteString("    xs.startTag(" + space + ", " + local + ")\n")
			file.WriteString("    xs.text(" + kotlin_makeIdent(prefix, field.Name) + ")\n")
			file.WriteString("    xs.endTag(" + space + ", " + local + ")\n")
		case "element:name":
			local = "payload"
			file.WriteString("    xs.startTag(" + space + ", " + local + ")\n")
			file.WriteString("    xs.endTag(" + space + ", " + local + ")\n")
		case "startelement", "element":
			switch typ := field.Type.(type) {
			case Set:
				fields := []*Field(typ)
				for _, x := range fields {
					if x.EncodingRule != nil && x.EncodingRule.Type == "name" {
						// from enum
						local = kotlin_makeIdent(prefix, x.Name) + ".toString()"
					}
				}
				file.WriteString("    xs.startTag(" + kotlin_ns(target) + ", " + local + ")\n")
				kotlin_generate_attributes_encoder(file, prefix, fields)
				for _, x := range fields {
					if x.EncodingRule != nil && (x.EncodingRule.Type == "attribute" ||
						x.EncodingRule.Type == "name" || x.EncodingRule.Type == "cdata") {
						continue
					}
					name := x.Name
					if name == "" {
						name = x.Type.(string)
					}
					if _, ok := x.Type.(Set); ok {
						kotlin_generate_element_encoder(file, "  ", target, kotlin_makeClassName(name), x)
					} else {
						if x.EncodingRule != nil && x.EncodingRule.Type == "element:bool" {
							file.WriteString("    if (" + kotlin_makeIdent(prefix, name) + " != null && " +
								kotlin_makeIdent(prefix, name) + " == true) {\n")
							kotlin_generate_element_encoder(file, "      ", target, prefix, x)
							file.WriteString("    }\n")
						} else {
							kotlin_generate_element_encoder(file, "      ", target, prefix, x)
						}
					}
				}
				for _, x := range fields {
					if x.EncodingRule != nil && x.EncodingRule.Type == "cdata" {
						file.WriteString("    if (" + kotlin_makeIdent(prefix, x.Name) + " != null) {\n")
						file.WriteString("      xs.text(" + kotlin_makeIdent(prefix, x.Name) + ")\n")
						file.WriteString("    }\n")
					}
				}
				if field.EncodingRule.Type == "element" {
					file.WriteString("    xs.endTag(" + space + ", " + local + ")\n")
				}
			case Sequence, SequenceOf:
				file.WriteString("    xs.startTag(" + kotlin_ns(target) + ", " + local + ")\n")
				
				file.WriteString("    payloadSequence?.forEach { it.encode(xs) }\n")
			case Choice:
				file.WriteString("    xs.startTag(" + kotlin_ns(target) + ", " + local + ")\n")
				file.WriteString("      payload?.encode(xs)\n")
			case string:
				file.WriteString("    string not implemented\n")
			case Enum:
				file.WriteString("    enum not implemented\n")
			}
		}
	}
}

func kotlin_generate_attributes_encoder(file *os.File, prefix string, fields []*Field) {
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
			space := x.EncodingRule.Space
			local := x.Name
			if x.EncodingRule.Name != "" {
				local = x.EncodingRule.Name
			}
			if s, ok := x.Type.(string); ok && s == "xmllang" {
				space = ns_xml
				local = "lang"
			}
			if b, ok := x.Type.(string); ok && b == "boolean" {
				file.WriteString("    if (" + kotlin_makeIdent(prefix, x.Name) + " != null && " +
					kotlin_makeIdent(prefix, x.Name) + " == true) {\n")
				file.WriteString("      xs.attribute(\"" + space + "\", \"" + local + "\", \"true\")\n")
			} else {
				value := kotlin_simplevalue(prefix, x)
				file.WriteString("    if (" + kotlin_makeIdent(prefix, x.Name) + " != null) {\n")
				file.WriteString("      xs.attribute(\"" + space + "\", \"" + local + "\", " + value + ")\n")
			}
			file.WriteString("    }\n")
		}
	}
}

func kotlin_generate_element_encoder(file *os.File, ident string, target *Target, prefix string, field *Field) {
	if field.EncodingRule != nil {
		space := kotlin_ns(target)
		if field.EncodingRule.Space != "" {
			space = "\"" + field.EncodingRule.Space + "\""
		}
		local := field.Name
		if field.EncodingRule.Name != "" {
			local = field.EncodingRule.Name
		}
		local = "\"" + local + "\""
		switch field.EncodingRule.Type {
		case "element:bool":
			file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
			file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
		case "element:cdata":
			if _, ok := field.Type.(SequenceOf); ok {
				file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "?.forEach {\n")
				file.WriteString(ident + "  xs.startTag(" + space + ", " + local + ")\n")
				file.WriteString(ident + "  xs.text(it)\n")
				file.WriteString(ident + "  xs.endTag(" + space + ", " + local + ")\n")
				file.WriteString(ident + "}\n")
			} else {
				value := kotlin_simplevalue(prefix, field)
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
				file.WriteString(ident + "xs.text(" + value + ")\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
			}
		case "element:name":
			// from enum
			local = kotlin_makeIdent(prefix, field.Name) + ".toString()"
			file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
			file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
		case "element":
			switch typ := field.Type.(type) {
			case string:
				name := field.Name
				if name == "" {
					name = typ
				}
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
				file.WriteString(ident + kotlin_makeIdent(prefix, name) + "?.encode(xs)\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
			case SequenceOf, Sequence:
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
				file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "?.forEach { it.encode(xs) }\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
			case Extension:
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
				file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "?.encode(xs)\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
			case Choice:
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
				file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "?.encode(xs)\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
			case Enum:
			case Set:
				fields := []*Field(typ)
				space := kotlin_ns(target)
				if field.EncodingRule.Space != "" {
					space = "\"" + field.EncodingRule.Space + "\""
				}
				local := field.Name
				if field.EncodingRule.Name != "" {
					local = field.EncodingRule.Name
				}
				local = "\"" + local + "\""
				for _, x := range fields {
					if x.EncodingRule != nil && x.EncodingRule.Type == "name" {
						// from enum
						local = kotlin_makeIdent(prefix, x.Name) + ".toString()"
					}
				}
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ")\n")
				kotlin_generate_attributes_encoder(file, prefix, fields)
				for _, x := range fields {
					kotlin_generate_element_encoder(file, ident +
						"  ", target, prefix + kotlin_makeClassName(field.Name), x)
				}
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ")\n")
			}
		}
	} else {
		switch typ := field.Type.(type) {
		case string:
			name := field.Name
			if name == "" {
				name = typ
			}
			file.WriteString(ident + kotlin_makeIdent(prefix, name) + "?.encode(xs)\n")
		case SequenceOf, Sequence:
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "?.forEach { it.encode(xs) }\n")
		case Extension, Choice:
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "?.encode(xs)\n")
		case Enum:
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				kotlin_generate_element_encoder(file, ident + "  ", target, prefix, x)
			}
		}
	}
}

func kotlin_simplevalue(prefix string, field *Field) string {
	switch typ := field.Type.(type) {
	case string:
		switch typ {
		case "xmllang", "string":
			return kotlin_makeIdent(prefix, field.Name)
//		case "jid":
//		case "uint", "int":
//		case "boolean":
//		case "datetime":
		}
	}
	return kotlin_makeIdent(prefix, field.Name) + ".toString()"
}

func kotlin_generate_decoder(file *os.File, target *Target, field *Field) {
	file.WriteString("  override fun decode(xp: XmlPullParser) {\n")
	ident := "    "
	prefix := ""
	kotlin_generate_element_decoder(file, target, 0, false, false, ident, prefix, field)
	file.WriteString("  }\n")
}

func kotlin_generate_element_decoder(file *os.File, target *Target, depth int, elseif bool, decl bool,
	ident, prefix string, field *Field) {
	depth++
	if field.EncodingRule != nil {
		switch field.EncodingRule.Type {
		case "element:cdata":
			file.WriteString(ident + "val _data_: String = " + cfg.Kotlin.Package_prefix + ".xmlencoder.getText(xp)\n")
			if _, ok := field.Type.(SequenceOf); ok {
				file.WriteString(ident + "if (" + kotlin_makeIdent(prefix, field.Name) + " == null)\n")
				file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + " = ArrayList()\n")
				file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "!!.add(_data_)\n")
			} else {
				kotlin_simplevalue_decode(file, ident, prefix, "_data_", field)
				// file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + " = data\n")
			}
		case "element:bool":
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + " = true\n")
			file.WriteString(ident + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
		case "element:name":
			file.WriteString(ident + "String name = xp.getName()\n")
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "(name)\n")
		case "startelement", "element":
			if depth == 1 {
				haschilds := false
				if fields, ok := field.Type.(Set); ok {
					for _, f := range []*Field(fields) {
						if f.EncodingRule != nil &&
							(f.EncodingRule.Type == "attribute" || f.EncodingRule.Type == "cdata" ||
							f.EncodingRule.Type == "name") {
							continue
						} else {
							haschilds = true
							break
						}
					}
				} else {
					haschilds = true
				}
				if haschilds {
					file.WriteString(ident + "var ev: Int\n")
				}
			}
			switch typ := field.Type.(type) {
			case string:
				if typ == "langstring" {
					file.WriteString(ident + "if (" + kotlin_makeIdent(prefix, field.Name) + " == null) {\n")
					file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + " = LangString()\n")
					file.WriteString(ident + "}\n")
					file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + "!!.decode(xp)\n")
				} else {
					file.WriteString("not implemented\n")
				}
			case Extension:
				file.WriteString(ident + "while (true) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next()\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				if typ.Local == "" {
					file.WriteString(forident +
						"payload = " + cfg.Kotlin.Package_prefix + ".xmlencoder.getExtension(xp.getNamespace(), xp.getName())\n")
					file.WriteString("if (payload != null) {\n")
					file.WriteString("payload.decode(xp)\n")
					file.WriteString("} else {\n")
					file.WriteString(forident + "  " + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
					file.WriteString(forident + "    }\n")
				} else {
					decoder, err := kotlin_getFullClassName(typ.Space, typ.Local)
					if err != nil { panic(err) }
					file.WriteString("payload = " + decoder + "()\n")
					file.WriteString("payload!!.decode(xp)\n")
				}
				if typ.Local != "" {
					file.WriteString(forident + "  } else {\n")
					file.WriteString(forident + "    " + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
					file.WriteString(forident + "  }\n")
				}					
				file.WriteString(forident + "}\n")
			case SequenceOf:
				file.WriteString(ident + "while (true) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next()\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				if string(typ) == "extension" {
					file.WriteString(forident + "val obj: XmlEncoder? = " + cfg.Kotlin.Package_prefix + ".xmlencoder.getExtension(xp.getNamespace(), xp.getName())\n")
					file.WriteString(forident + "if (obj != null) {\n")
				} else {
					var tname string
					space := kotlin_ns(target)
					local := field.Name
					if t, ok := kotlinSimpleTypes[string(typ)]; ok {
						tname = t.Type
					} else {
						f := kotlin_getFieldByName(target, string(typ))
						tname = kotlin_makeClassName(f.Name)
						if f.EncodingRule != nil && f.EncodingRule.Space != "" {
							space = "\"" + f.EncodingRule.Space + "\""
						}
						local = f.Name
						if f.EncodingRule != nil && f.EncodingRule.Name != "" {
							local = f.EncodingRule.Name
						}
					}
					local = "\"" + local + "\""
					file.WriteString(forident + "if (xp.getNamespace() == " + space + " && xp.getName() == " +
						local + ") {\n")
					file.WriteString(forident + "  val obj: " + tname + " = " + tname + "()\n")
				}
				file.WriteString(forident + "  obj.decode(xp)\n")
				var name string
				if depth == 1 {
					name = "payloadSequence"
				} else {
					name = kotlin_makeIdent(prefix, field.Name)
				}
				file.WriteString(forident + "if (" + name + " == null) {\n")
				file.WriteString(forident + "  " + name + " = ArrayList()\n")
				file.WriteString(forident + "}\n")
				file.WriteString("if (" + name + " == null)\n")
				file.WriteString(name + " = ArrayList()\n")
				file.WriteString(forident + name + "?.add(obj)\n")
				file.WriteString(forident + "} else {\n")
				file.WriteString(forident + "  " + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
				file.WriteString(forident + "  }\n")
				file.WriteString(forident + "}\n")
				file.WriteString(ident + "}\n")
			case Sequence:
				file.WriteString(ident + "while (true) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next()\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				fields := []*Field(typ)
				for _, x := range fields {
					name := x.Type.(string)
					f := kotlin_getFieldByName(target, name)
					space := kotlin_ns(target)
					if f.EncodingRule != nil && f.EncodingRule.Space != "" {
						space = "\"" + f.EncodingRule.Space + "\""
					}
					local := f.Name
					if f.EncodingRule != nil && f.EncodingRule.Name != "" {
						local = f.EncodingRule.Name
					}
					local = "\"" + local + "\""
					if !elseif {
						elseif = true
						file.WriteString(forident + "  ")
					} else {
						file.WriteString(forident + "  } else ")
					}
					file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
						local + ") {\n")
					file.WriteString(forident + "    val obj: " + kotlin_makeClassName(name) + " = " +
						kotlin_makeClassName(name) + "()\n")
					file.WriteString(forident + "    obj.decode(xp)\n")
					if depth == 1 {
						file.WriteString("if (payloadSequence == null)\n")
						file.WriteString("payloadSequence = ArrayList()\n")
						file.WriteString(forident + "    payloadSequence?.add(obj)\n")
					} else {
						file.WriteString("if (" + kotlin_makeIdent(prefix, field.Name) + " == null)\n")
						file.WriteString(kotlin_makeIdent(prefix, field.Name) + " = ArrayList()\n")
						file.WriteString(forident + "    " + kotlin_makeIdent(prefix, field.Name) + "?.add(obj)\n")
					}
				}
				file.WriteString(forident + "  }\n")
				file.WriteString(forident + "}\n")
				file.WriteString(ident + "}\n")
			case Choice:
				file.WriteString(ident + "while (true) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next()\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				fields := []*Field(typ)
				elseif := false
				for _, x := range fields {
					name := x.Type.(string)
					f := kotlin_getFieldByName(target, name)
					space := kotlin_ns(target)
					if f.EncodingRule != nil && f.EncodingRule.Space != "" {
						space = "\"" + f.EncodingRule.Space + "\""
					}
					local := f.Name
					if f.EncodingRule != nil && f.EncodingRule.Name != "" {
						local = f.EncodingRule.Name
					}
					local = "\"" + local + "\""
					if !elseif {
						elseif = true
						file.WriteString(forident + "  ")
					} else {
						file.WriteString(forident + "  } else ")
					}
					file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
						local + ") {\n")
					file.WriteString(forident + "    val obj: " + kotlin_makeClassName(name) + " = " +
						kotlin_makeClassName(name) + "()\n")
					file.WriteString(forident + "    obj.decode(xp)\n")
					if depth == 1 {
						file.WriteString(forident + "    payload = obj\n")
					} else {
						file.WriteString(forident + "    " + kotlin_makeIdent(prefix, field.Name) + " = obj\n")
					}
				}
				file.WriteString(forident + "  }\n")
				file.WriteString(forident + "}\n")
				file.WriteString(ident + "}\n")
			case Set:
				fields := []*Field(typ)
				if len(fields) == 0 {
					file.WriteString(ident + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
					return 
				}
				if depth > 1 {
					prefix += kotlin_makeClassName(field.Name)
				}
				decl = kotlin_generate_attributes_decoder(file, decl, ident, prefix, fields)
				var any, anyname, cdata, extension *Field
				var elems []*Field
			Loop:
				for _, x := range fields {
					if x.EncodingRule == nil {
						switch typ := x.Type.(type) {
						case string:
							if typ == "extension" {
								extension = x
							} else {
								elems = append(elems, x)
							}
						case SequenceOf:
							if string(typ) == "extension" {
								extension = x
							} else {
								elems = append(elems, x)
							}
						case Extension:
							if typ.Local == "" {
								extension = x
							} else {
								elems = append(elems, x)
							}
						default:							
							elems = append(elems, x)
						}
					} else {
						switch x.EncodingRule.Type {
						case "name": anyname = x
						case "cdata": cdata = x
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
				if anyname != nil {
					file.WriteString(ident + kotlin_makeIdent(prefix, anyname.Name) + " = " +
						kotlin_makeEnumName(anyname.Name) + ".valueOf(xp.getName())\n")
				}
				if len(elems) > 0 || any != nil || extension != nil {
					file.WriteString(ident + "while (true) {\n")
					file.WriteString(ident + "  ev = xp.next()\n")
					file.WriteString(ident + "  if (ev == XmlPullParser.END_TAG) {\n")
					file.WriteString(ident + "    break\n")
					file.WriteString(ident + "  } else if (ev == XmlPullParser.START_TAG) {\n")
					forident := ident + "  "
					elseif := false
					for _, z := range elems {
						if z.EncodingRule != nil {
							space := kotlin_ns(target)
							if z.EncodingRule != nil && z.EncodingRule.Space != "" {
								space = "\"" + z.EncodingRule.Space + "\""
							}
							local := z.Name
							if z.EncodingRule != nil && z.EncodingRule.Name != "" {
								local = z.EncodingRule.Name
							}
							local = "\"" + local + "\""
							if !elseif {
								file.WriteString(forident + "  ")
								elseif = true
							} else {
								file.WriteString(forident + "  } else ")
							}
							file.WriteString("if (xp.getNamespace() == " + space +
								" && xp.getName() == " + local + ") {\n")
							kotlin_generate_element_decoder(file, target, depth, elseif, decl,
								forident + "    ", prefix, z)
						} else {
							kotlin_generate_element_decoder(file, target, depth, elseif, decl,
								forident + "  ", prefix, z)
							if !elseif {
								elseif = true
							}
						}
					}
					if any != nil {
						file.WriteString(ident)
						if elseif {
							file.WriteString(" } else ")
						}
						file.WriteString("if (xp.getNamespace() == " + kotlin_ns(target) + ") {\n")
						switch any.EncodingRule.Type {
						case "element:name":
							var typ string
							if t, ok := any.Type.(string); ok {
								is_enum := false
								for _, f := range target.Fields {
									if t == f.Name {
										if _, ok := f.Type.(Enum); ok {
											is_enum = true
										}
										break
									}
								}
								if is_enum {
									typ = kotlin_makeEnumName(t)
								} else {
									typ = kotlin_makeClassName(t)
								}
							} else {
								typ = kotlin_makeEnumName(any.Name)
							}
							file.WriteString("  " + kotlin_makeIdent(prefix, any.Name) + " = " +
								typ + ".valueOf(xp.getName())\n")
						case "name":
							file.WriteString("not implemented\n")
						case "element":
							kotlin_generate_element_decoder(file, target, depth, elseif, decl,
								forident + "  ", prefix, any)
						}
					}
					if extension != nil {
						file.WriteString(forident + "  ")
						if elseif {
							file.WriteString("} else ")
						}
						file.WriteString("if (xp.getNamespace() != " + kotlin_ns(target) + ") {\n")
						kotlin_generate_element_decoder(file, target, depth, elseif, decl,
								forident + "    ", prefix, extension)
						elseif = true
					}
					if len(elems) > 0 || any != nil || extension != nil {
						file.WriteString(forident + "  } else {\n")
						file.WriteString(forident + "    " + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
						file.WriteString(forident + "  }\n")
					}
					file.WriteString(ident + "  }\n")
					file.WriteString(ident + "}\n")
				}
				if cdata != nil {
					file.WriteString(ident + "val _data_: String = " +
						cfg.Kotlin.Package_prefix + ".xmlencoder.getText(xp)\n")
					kotlin_simplevalue_decode(file, ident, prefix, "_data_", cdata)
				}
			}
		}
	} else {
		switch typ := field.Type.(type) {
		case string:
			name := field.Name
			f := kotlin_getFieldByName(target, typ)
			if name == "" {
				name = typ
			}
			space := kotlin_ns(target)
			local := f.Name
			if f.EncodingRule != nil && f.EncodingRule.Space != "" {
				space = "\"" + f.EncodingRule.Space + "\""
			}
			if f.EncodingRule != nil && f.EncodingRule.Name != "" {
				local = f.EncodingRule.Name
			}
			local = "\"" + local + "\""
			file.WriteString(ident)
			if elseif {
				file.WriteString("} else ")
			}
			file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
				local + ") {\n")
			file.WriteString(ident + "  " + kotlin_makeIdent(prefix, name) + " = " +
				kotlin_makeClassName(typ) + "()\n")
			file.WriteString(ident + "  " + kotlin_makeIdent(prefix, name) + "!!.decode(xp)\n")
		case SequenceOf:
			if string(typ) == "extension" {
				file.WriteString(ident +
					"val xe: XmlEncoder? = " + cfg.Kotlin.Package_prefix + ".xmlencoder.getExtension(xp.getNamespace(), xp.getName())\n")
				file.WriteString(ident + "if (xe != null) {\n")
				file.WriteString(ident + "  xe.decode(xp)\n")
				file.WriteString("if (" + kotlin_makeIdent(prefix, field.Name) + " == null)\n")
				file.WriteString(kotlin_makeIdent(prefix, field.Name) + " = ArrayList()\n")
				file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + "?.add(xe)\n")
				file.WriteString(ident + "} else {\n")
				file.WriteString(ident + "  " + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
				file.WriteString(ident + "}\n")
			} else {
				f := kotlin_getFieldByName(target, string(typ))
				space := kotlin_ns(target)
				local := f.Name
				if f.EncodingRule != nil && f.EncodingRule.Space != "" {
					space = "\"" + f.EncodingRule.Space + "\""
				}
				if f.EncodingRule != nil && f.EncodingRule.Name != "" {
					local = f.EncodingRule.Name
				}
				local = "\"" + local + "\""
				file.WriteString(ident)
				if elseif {
					file.WriteString("} else ")
				}
				file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
					local + ") {\n")
				file.WriteString(ident + "  val obj: " + kotlin_makeClassName(string(typ)) + " = " +
					kotlin_makeClassName(string(typ)) + "()\n")
				file.WriteString(ident + "  obj.decode(xp)\n")
				file.WriteString("if (" + kotlin_makeIdent(prefix, field.Name) + " == null)\n")
				file.WriteString(kotlin_makeIdent(prefix, field.Name) + " = ArrayList()\n")
				file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + "?.add(obj)\n")
			}
		case Extension:
			if typ.Local == "" {
				file.WriteString(ident + "val xe: XmlEncoder? = " + cfg.Kotlin.Package_prefix + ".xmlencoder.getExtension(xp.getNamespace(), xp.getName())\n")
				file.WriteString(ident + "if (xe != null) {\n")
				file.WriteString(ident + "  xe.decode(xp)\n")
				file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + " = xe\n")
				file.WriteString(ident + "} else {\n")
				file.WriteString(ident + "  " + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
				file.WriteString(ident  + "}\n")
			} else {
				file.WriteString(ident)
				if elseif {
					file.WriteString("} else ")
				}
				space := "\"" + typ.Space + "\""
				local := "\"" + typ.Local + "\""
				file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
					local + ") {\n")
				file.WriteString(ident + "  val xe: XmlEncoder? = " + cfg.Kotlin.Package_prefix + ".xmlencoder.getExtension(" +
					space + ", " + local + ")\n")
				file.WriteString(ident + "  if (xe != null && xe is " +
					kotlin_getExtensionType(typ.Space, typ.Local) + ") {\n")
				file.WriteString(ident + "    xe.decode(xp)\n")
					file.WriteString(ident + "    " + kotlin_makeIdent(prefix, field.Name) + " = xe\n")
				file.WriteString(ident + "  } else {\n")
				file.WriteString(ident + "    " + cfg.Kotlin.Package_prefix + ".xmlencoder.getEndTag(xp)\n")
				file.WriteString(ident  + "  }\n")
			}
		case Sequence:
			fields := []*Field(typ)
			for _, x := range fields {
				f := x
				if x.Name == "" {
					f = kotlin_getFieldByName(target, x.Type.(string))
				}
				space := kotlin_ns(target)
				if f.EncodingRule != nil && f.EncodingRule.Space != "" {
					space = "\"" + f.EncodingRule.Space + "\""
				}
				local := f.Name
				if f.EncodingRule != nil && f.EncodingRule.Name != "" {
					local = f.EncodingRule.Name
				}
				local = "\"" + local + "\""
				file.WriteString(ident)
				if elseif {
					file.WriteString("} else ")
				}
				file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
					local + ") {\n")
				file.WriteString(ident + "  val obj: " + kotlin_makeClassName(f.Name) + " = " +
					kotlin_makeClassName(f.Name) + "()\n")
				file.WriteString(ident + "  obj.decode(xp)\n")
				file.WriteString("if (" + kotlin_makeIdent(prefix, field.Name) + " == null)\n")
				file.WriteString(kotlin_makeIdent(prefix, field.Name) + " = ArrayList()\n")
				file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + "?.add(obj)\n")
			}
		case Choice:
			fields := []*Field(typ)
			for _, x := range fields {
				f := kotlin_getFieldByName(target, x.Type.(string))
				space := kotlin_ns(target)
				if f.EncodingRule != nil && f.EncodingRule.Space != "" {
					space = "\"" + f.EncodingRule.Space + "\""
				}
				local := f.Name
				if f.EncodingRule != nil && f.EncodingRule.Name != "" {
					local = f.EncodingRule.Name
				}
				local = "\"" + local + "\""
				file.WriteString(ident)
				if elseif {
					file.WriteString("} else ")
				}
				file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
					local + ") {\n")
				file.WriteString(ident + "  val obj: " + kotlin_makeClassName(f.Name) + " = " +
					kotlin_makeClassName(f.Name) + "()\n")
				file.WriteString(ident + "  obj.decode(xp)\n")
				file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + " = obj\n")
				if !elseif {
					elseif = true
				}
			}
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				f := x
				if x.Name == "" {
					f = kotlin_getFieldByName(target, x.Type.(string))
				}
				space := kotlin_ns(target)
				if f.EncodingRule != nil && f.EncodingRule.Space != "" {
					space = "\"" + f.EncodingRule.Space + "\""
				}
				local := f.Name
				if f.EncodingRule != nil && f.EncodingRule.Name != "" {
					local = f.EncodingRule.Name
				}
				local = "\"" + local + "\""
				file.WriteString(ident)
				if elseif {
					file.WriteString("} else ")
				}
				file.WriteString("if (xp.getNamespace() == " + space + " && xp.getName() == " +
					local + ") {\n")
				kotlin_generate_element_decoder(file, target, depth, elseif, decl,
					ident + "    ", prefix + kotlin_makeClassName(field.Name), x)
/*					
				file.WriteString(ident + "  val obj: " + kotlin_makeClassName(f.Name) + " = " +
					kotlin_makeClassName(f.Name) + "()\n")
				file.WriteString(ident + "  obj.decode(xp)\n")
				name := x.Name
				if name == "" {
					name = x.Type.(string)
				}
				pr := prefix + kotlin_makeClassName(field.Name)
				file.WriteString(ident + "  " + kotlin_makeIdent(pr, name) + " = obj\n")
*/
				if !elseif {
					elseif = true
				}
			}
			
		}
	}
}

func kotlin_getFieldByName(target *Target, f string) *Field {
	for _, x := range target.Fields {
		if x.Name == f {
			return x
		}
	}
	return nil
}

func kotlin_generate_attributes_decoder(file *os.File, decl bool, ident, prefix string, fields []*Field) (ret bool) {
	hasAttrs := false
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
			hasAttrs = true
			ret = true
			break
		}
	}
	if hasAttrs && !decl {
		file.WriteString(ident + "var _value_: String\n")
	}
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
			space := ""
			if x.EncodingRule.Space != "" {
				space = x.EncodingRule.Space
			}
			local := x.Name
			if x.EncodingRule.Name != "" {
				local = x.EncodingRule.Name
			}
			if str, ok := x.Type.(string); ok && str == "xmlstring" {
				space = ns_xml
				local = "lang"
			}
			space = "\"" + space + "\""
			local = "\"" + local + "\""
			
			file.WriteString(ident + "_value_ = xp.getAttributeValue(" + space + ", " + local + ")\n")
			file.WriteString("if (_value_ != null)\n")
			kotlin_simplevalue_decode(file, ident + "  ", prefix, "_value_", x)
		}
	}
	return
}

func kotlin_simplevalue_decode(file *os.File, ident, prefix, varname string, field *Field) {
	switch typ := field.Type.(type) {
	case string:
		switch typ {
		case "boolean":
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + " = " + varname +
				".toBoolean()\n")
			return
		case "int", "uint":
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) +
				" = " + varname + ".toInt()\n")
			return
		case "bytestring", "string", "xmllang":
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + " = " + varname + "\n")
			return
		case "jid":
			file.WriteString(ident + kotlin_makeIdent(prefix, field.Name) + " = JID(" + varname + ")\n")
			return
		case "datetime":
			file.WriteString(ident + "try {\n")
			file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) +
				" = " + cfg.Kotlin.Package_prefix + ".xmlencoder.parseRFC3339(" + varname + ")\n")
			file.WriteString(ident + "} catch (e: java.text.ParseException) {}\n")
			return
		}
	case Enum:
		file.WriteString(ident + "try {\n")
		file.WriteString(ident + "  " + kotlin_makeIdent(prefix, field.Name) + " = " +
			kotlin_makeEnumName(field.Name) + ".valueOf(" + varname + ")\n")
		file.WriteString(ident + "} catch (iae: IllegalArgumentException) {}\n")
		return
	}
	file.WriteString("not implemented\n")
}

func kotlin_generate_adders() error {
	filename := filepath.Join(cfg.Kotlin.Outdir, "xmlencoder", "extensions.kt")
	file, err := os.Create(filename)
	if err != nil { return err }
	file.WriteString("package " + cfg.Kotlin.Package_prefix + ".xmlencoder\n\n")
	file.WriteString("import org.xmlpull.v1.XmlPullParser\n")
	file.WriteString("import " + cfg.Kotlin.Package_prefix + ".xmlencoder.QName\n\n")
	
//	file.WriteString("val extensions = object {\n")
	//	file.WriteString("  val mapData: Map<QName, (() -> XmlEncoder)>\n\n")
	// file.WriteString("  init {\n")
	// file.WriteString("    mapData = mapOf(\n")
	file.WriteString("val extensions : Map<QName, (() -> XmlEncoder)> = hashMapOf(\n")
	first := false
	for _, schema := range schemas {
		category := ""
		if cat, ok := schema.Props["category"]; ok {
			category = cat
		}
		for _, target := range schema.Targets {
			for _, field := range target.Fields {
				if field.Reciver_type != "" {
					local := field.Name
					if field.EncodingRule != nil && field.EncodingRule.Name != "" {
						local = field.EncodingRule.Name
					}
					if first {
						file.WriteString(",\n")
					}
					first = true
					file.WriteString("    QName(\"" + target.Space + "\", \"" + local + "\") to ")
					file.WriteString("{ " + cfg.Kotlin.Package_prefix + ".")
					if category == "extension" {
						file.WriteString("extensions.")
					}
					file.WriteString(schema.PackageName + ".")
					if target.Name != "" {
						file.WriteString(target.Name + ".")
					}
					file.WriteString(kotlin_makeClassName(field.Name) + "()}")
				}
			}
		}
	}
	file.WriteString("    )\n\n")

	file.WriteString("fun getExtension(space: String, local: String) : XmlEncoder? {\n")
	file.WriteString("  val xe: (() -> XmlEncoder)? = extensions.get(QName(space, local))\n")
	file.WriteString("  if (xe != null)\n")
	file.WriteString("     return xe()\n")
	file.WriteString("  return null\n")
	file.WriteString("}\n")
//	file.WriteString("  }\n")
	file.Close()
	return nil
}

func kotlin_ns(target *Target) string {
	if target.Name != "" {
		return "ns_" + target.Name
	}
	return "ns"
}
