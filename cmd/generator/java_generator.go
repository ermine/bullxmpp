package main

import (
	"os"
	"path/filepath"
	"unicode"
	"fmt"
)

func JavaGenerate() error {
	if cfg.Java.Outdir == "" {
		panic("no outdir")
	}
	for _, schema := range schemas {
		dir := cfg.Java.Outdir
		if category, ok := schema.Props["category"]; ok {
			if category == "extension" {
				dir = filepath.Join(cfg.Java.Outdir, "extensions")
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
		java_generate_package(fdir, schema)
	}
	java_generate_adders()
	return nil
}

func java_generate_package(dir string, schema *Schema) error {
	for _, target := range schema.Targets {
		currdir := dir
		if target.Name != "" {
			currdir = filepath.Join(currdir, target.Name)
		}
		for _, field := range target.Fields {
			filename := filepath.Join(currdir, java_makeClassName(field.Name) + ".java")
			file, err := os.Create(filename)
			if err != nil { return err }
			defer file.Close()
			java_generate_class(file, schema, target, field)
		}
	}
	return nil
}

func java_generate_class(file *os.File, schema *Schema, target *Target, field *Field) {
	subname := ""
	if target.Name != "" {
		subname = "." + target.Name
	}
	file.WriteString("package " + cfg.Java.Package_prefix + ".")
	if category, ok := schema.Props["category"]; ok && category == "extension" {
		file.WriteString("extensions.")
	}
	file.WriteString(schema.PackageName + subname + ";\n\n")
	var imports []string
	java_getImports(target, field, &imports)
	file.WriteString("import java.io.IOException;\n")
	file.WriteString("import org.xmlpull.v1.XmlPullParser;\n")
	file.WriteString("import org.xmlpull.v1.XmlPullParserException;\n")
	file.WriteString("import org.xmlpull.v1.XmlSerializer;\n")
	file.WriteString("import ru.sulci.sb.xmlencoder.XmlEncoder;\n")
	for _, i := range imports {
		file.WriteString("import " + i + ";\n")
	}
	file.WriteString("\n")
	file.WriteString("public class " + java_makeClassName(field.Name) +
		" implements XmlEncoder {\n")
	if field.Reciver_type != "" {
		file.WriteString("  public static final String NS = \"" + target.Space + "\";\n\n")
	} else {
		file.WriteString("  private static final String NS = \"" + target.Space + "\";\n\n")
	}		
	java_generate_class_fields(file, target, field)
	file.WriteString("\n\n")
	file.WriteString("  public " + java_makeClassName(field.Name) + "() {}\n")
	file.WriteString("\n\n")
	java_generate_encoder(file, target, field)
	java_generate_decoder(file, target, field)
	file.WriteString("}\n")
}

func java_makeUppercase(s string) string {
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

func java_makeClassName(s string) string {
	return java_makeUppercase(s) + "Proto"
}

func java_makeIdent(prefix, s string) string {
	return "m" + prefix + java_makeUppercase(s)
}

func java_makeEnumName(s string) string {
	return java_makeUppercase(s) + "Enum"
}

func java_getImports(target *Target, field *Field, imports *[]string) {
	switch typ := field.Type.(type) {
	case string:
		if data, ok := javaSimpleTypes[typ]; ok {
			if data.Import != "" {
				appendImport(imports, data.Import)
			}
		}
	case SequenceOf:
		appendImport(imports, "java.util.ArrayList")
		if data, ok := javaSimpleTypes[string(typ)]; ok {
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
									pName := cfg.Java.Package_prefix + "." + s.PackageName
								if t.Name != "" {
									pName += "." + t.Name
								}
								pName += "." + java_makeClassName(f.Name)
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
		appendImport(imports, "java.util.ArrayList")
		for _, x := range []*Field(typ) {
				java_getImports(target, x, imports)
		}
	case Choice:
		for _, x := range []*Field(typ) {
			java_getImports(target, x, imports)
		}			
	case Set:
		for _, x := range []*Field(typ) {
			java_getImports(target, x, imports)
		}
	case Enum:
	}
}

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

var javaSimpleTypes = map[string]struct{
	Type string
	Import string
}{
	"boolean": {"boolean", ""},
	"string": {"String", ""},
	"bytestring": {"String", ""},
	"jid": {"JID", "ru.sulci.sb.jid.JID"},
	"datetime": {"Date", "java.util.Date"},
	"int": {"Integer", ""},
	"uint": {"Integer", ""},
	"xmllang": {"String", ""},
	"langstring": {"LangString", "ru.sulci.sb.xmlencoder.LangString"},
	"extension": {"XmlEncoder", ""},
}

func java_generate_class_fields(file *os.File, target *Target, field *Field) {
	prefix := ""
	switch typ := field.Type.(type) {
	case string:
		if t, ok := javaSimpleTypes[typ]; ok {
			file.WriteString("  " + t.Type + " " + java_makeIdent(prefix, field.Name) + ";\n")
		} else {
			name := field.Name
			if name == "" {
				name = typ
			}
			file.WriteString("  " + java_makeClassName(typ) + " " + java_makeIdent(prefix, name) + ";\n")
		}
	case Extension:
		if typ.Local == "" {
			file.WriteString("XmlEncoder mPayload;\n")
		} else {
			t := java_getExtensionType(typ.Space, typ.Local)
			file.WriteString("  " + t + " mPayload;\n")
		}
	case SequenceOf:
		t := string(typ)
		tt := "XmlEncoder"
		if t != "extension" {
			if s, ok := javaSimpleTypes[t]; ok {
				tt = s.Type
			} else {
				tt = java_makeClassName(t)
			}
		}
		file.WriteString("  ArrayList<" + tt + "> mPayloadSequence;\n")
		file.WriteString("  public ArrayList<" + tt + "> get" + java_makeUppercase(field.Name) +
			"() {\n")
		file.WriteString("    return mPayloadSequence;\n")
		file.WriteString("  }\n")
		file.WriteString("  public void add" + java_makeUppercase(field.Name) + "(" + tt +
			" value) {\n")
		file.WriteString("    mPayloadSequence.add(value);\n")
		file.WriteString("  }\n")
	case Sequence:
		file.WriteString("  ArrayList<XmlEncoder> mPayloadSequence;\n")
		file.WriteString("  public ArrayList<XmlEncoder> get" + java_makeUppercase(field.Name) +
			"() {\n")
		file.WriteString("    return mPayloadSequence;\n")
		file.WriteString("  }\n")
		file.WriteString("  public void add" + java_makeUppercase(field.Name) + "(XmlEncoder xe) {\n")
		file.WriteString("    mPayloadSequence.add(xe);\n")
		file.WriteString("  }\n")
	case Choice:
		file.WriteString("  XmlEncoder mPayloadChoice;\n")
	case Set:
		fields := []*Field(typ)
		for _, x := range fields {
			java_generate_field(file, "  ", target, "", x)
		}
	case Enum:
	}
}

func java_generate_field(file *os.File, ident string, target *Target, prefix string, field *Field) {
	switch typ := field.Type.(type) {
	case string:
		if t, ok := javaSimpleTypes[typ]; ok {
			file.WriteString(ident + t.Type + " " + java_makeIdent(prefix, field.Name) + ";\n")
			file.WriteString(ident + "public " + t.Type + " get" + prefix +
				java_makeUppercase(field.Name) + "() {\n")
			file.WriteString(ident + "  return " + java_makeIdent(prefix, field.Name) + ";\n")
			file.WriteString(ident + "}\n")
			file.WriteString(ident + "public void set" + prefix + java_makeUppercase(field.Name) +
				"(" + t.Type + " value) {\n")
			file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = value;\n")
			file.WriteString(ident + "}\n")
		} else {
			name := field.Name
			if name == "" {
				name = typ
			}
			file.WriteString(ident + java_makeClassName(typ) + " " + java_makeIdent(prefix, name) + ";\n")
			file.WriteString(ident + "public " + java_makeClassName(typ) + " get" + prefix +
				java_makeUppercase(name) + "() {\n")
			file.WriteString(ident + "  return " + java_makeIdent(prefix, name) + ";\n")
			file.WriteString(ident + "}\n")
			file.WriteString(ident + "public void set" + prefix + java_makeUppercase(name) + "(" +
				java_makeClassName(typ) + " value) {\n")
			file.WriteString(ident + "  " + java_makeIdent(prefix, name) + " = value;\n")
			file.WriteString(ident + "}\n")
		}
	case Extension:
		var t string
		if typ.Local == "" {
			t = "XmlEncoder"
		} else {
			t = java_getExtensionType(typ.Space, typ.Local)
		}			
		file.WriteString(ident + t + " " + java_makeIdent(prefix, field.Name) + ";\n")
		file.WriteString(ident + "public " + t + " get" + prefix + java_makeUppercase(field.Name) +
			"() {\n")
		file.WriteString(ident + "  return " + java_makeIdent(prefix, field.Name) + ";\n")
		file.WriteString(ident + "}\n")
		file.WriteString(ident + "public void set" + prefix + java_makeUppercase(field.Name) + "(" +
			t + " value) {\n")
		file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = value;\n")
		file.WriteString(ident + "}\n")
	case SequenceOf:
		t := string(typ)
		tt := "XmlEncoder"
		if t != "extension" {
			if s, ok := javaSimpleTypes[t]; ok {
				tt = s.Type
			} else {
				tt = java_makeClassName(t)
			}
		}
		file.WriteString(ident + "ArrayList<" + tt + "> " + java_makeIdent(prefix, field.Name) + ";\n")
	case Sequence:
		file.WriteString(ident + "ArrayList<XmlEncoder> " + java_makeIdent(prefix, field.Name) + ";\n")
		fields := []*Field(typ)
		for _, x := range fields {
			if x.Name != "" {
				file.WriteString("  class " + java_makeClassName(x.Name) + " implements XmlEncoder {\n")
				java_generate_class_fields(file, target, x)
				java_generate_encoder(file, target, x)
				java_generate_decoder(file, target, x)
				file.WriteString("  }\n")
			}
		}
	case Enum:
		file.WriteString(ident + java_makeEnumName(field.Name) + " " +
			java_makeIdent(prefix, field.Name) + ";\n\n")
		file.WriteString(ident + "public " + java_makeEnumName(field.Name) + " get" + prefix +
			java_makeUppercase(field.Name) + "() {\n")
		file.WriteString(ident + "  return " + java_makeIdent(prefix, field.Name) + ";\n")
		file.WriteString(ident + "}\n")
		file.WriteString(ident + "public void set" + prefix + java_makeUppercase(field.Name) +
			"(" + java_makeEnumName(field.Name) + " value) {\n")
		file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = value;\n")
		file.WriteString(ident + "}\n")
		java_generate_enum(file, typ, java_makeEnumName(field.Name))
	case Choice:
		file.WriteString(ident + "XmlEncoder " + java_makeIdent(prefix, field.Name) + ";\n")
	case Set:
		fields := []*Field(typ)
		for _, x := range fields {
			java_generate_field(file, ident, target, java_makeClassName(field.Name), x)
		}
	}
}

func java_generate_enum(file *os.File, enum Enum, name string) {
	file.WriteString("  public enum " + name + " {\n")
	variants := []string(enum)
	str, hashyphen := java_make_enum_string(variants[0])
	file.WriteString("    " + str)
	for _, x := range variants[1:] {
		file.WriteString(",\n")
		str, h := java_make_enum_string(x)
		if h {
			hashyphen = true
		}
		file.WriteString("    " + str)
	}
	file.WriteString(";\n\n")
	file.WriteString("    @Override\n")
	file.WriteString("    public String toString() {\n")
	if hashyphen {
		file.WriteString("      return super.toString().toLowerCase().replace('_', '-');\n")
	} else {
		file.WriteString("      return super.toString().toLowerCase();\n")
	}
	file.WriteString("    }\n\n")
	file.WriteString("    public static " + name + " fromString(String value) {\n")
	if hashyphen {
		file.WriteString("      return " + name + ".valueOf(value.toUpperCase().replace('-', '_'));\n")
	} else {
		file.WriteString("      return " + name + ".valueOf(value.toUpperCase());\n")
	}
	file.WriteString("    }\n")
	file.WriteString("  }\n\n")
}

func java_make_enum_string(s string) (string, bool) {
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

func java_resulveType(typ string) string {
	if data, ok := javaSimpleTypes[typ]; ok {
		return data.Type
	} else {
		if typ == "extension" {
			return "XmlEncoder"
		}
	}
	return java_makeClassName(typ)
}

func java_getExtensionType(space, local string) string {
	for _, s := range schemas {
		for _, t := range s.Targets {
			if t.Space == space {
				fmt.Println("found namespace ", space, " ", local)
				for _, f := range t.Fields {
					fmt.Println("checking ", f.Name, " ", f.EncodingRule)
					if (f.EncodingRule != nil && f.EncodingRule.Name == local) ||
						(f.EncodingRule != nil && f.EncodingRule.Name == "" && f.Name == local) ||
						(f.EncodingRule == nil && f.Name == local) {
						pkgName := cfg.Java.Package_prefix + "."
						if category, ok := s.Props["category"]; ok && category == "extension" {
							pkgName += "extensions."
						}
						pkgName += s.PackageName
						if t.Name != "" {
							pkgName += "." + t.Name
						}
						return pkgName + "." + java_makeClassName(f.Name)
					}
				}
			}
		}
	}
	fmt.Println("not found ", space, " ", local)
	return "UnknownNamespace" + space + local
}

func java_generate_encoder(file *os.File, target *Target, field *Field) {
	file.WriteString("  public void Encode(XmlSerializer xs) throws IOException {\n")
//	if target.Prefix != "" {
		file.WriteString("    xs.setPrefix(\"" + target.Prefix + "\", \"" + target.Space + "\");")
//	}
	java_generate_class_encoder(file, target, field)
	file.WriteString("  }\n\n")
}

func java_generate_class_encoder(file *os.File, target *Target, field *Field) {
	prefix := ""
	if field.EncodingRule != nil {
		space := "NS"
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
			file.WriteString("    xs.startTag(" + space + ", " + local + ");\n")
			file.WriteString("    xs.endTag(" + space + ", " + local + ");\n")
		case "element:cdata":
			file.WriteString("    xs.startTag(" + space + ", " + local + ");\n")
			file.WriteString("    xs.text(" + java_makeIdent(prefix, field.Name) + ");\n")
			file.WriteString("    xs.endTag(" + space + ", " + local + ");\n")
		case "element:name":
			local = "mPayload"
			file.WriteString("    xs.startTag(" + space + ", " + local + ");\n")
			file.WriteString("    xs.endTag(" + space + ", " + local + ");\n")
		case "startelement", "element":
			switch typ := field.Type.(type) {
			case Set:
				fields := []*Field(typ)
				for _, x := range fields {
					if x.EncodingRule != nil && x.EncodingRule.Type == "name" {
						// from enum
						local = java_makeIdent(prefix, x.Name) + ".toString()"
					}
				}
				file.WriteString("    xs.startTag(NS, " + local + ");\n")
				java_generate_attributes_encoder(file, prefix, fields)
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
						java_generate_element_encoder(file, "  ", target, java_makeClassName(name), x)
					} else {
						if x.EncodingRule != nil && x.EncodingRule.Type == "element:bool" {
							file.WriteString("    if (" + java_makeIdent(prefix, name) + ") {\n")
						} else {
							file.WriteString("    if (" + java_makeIdent(prefix, name) + " != null) {\n")
						}
						java_generate_element_encoder(file, "      ", target, prefix, x)
						file.WriteString("    }\n")
					}
				}
				for _, x := range fields {
					if x.EncodingRule != nil && x.EncodingRule.Type == "cdata" {
						file.WriteString("    if (" + java_makeIdent(prefix, x.Name) + " != null) {\n")
						file.WriteString("      xs.text(" + java_makeIdent(prefix, x.Name) + ");\n")
						file.WriteString("    }\n")
					}
				}
				if field.EncodingRule.Type == "element" {
					file.WriteString("    xs.endTag(" + space + ", " + local + ");\n")
				}
			case Sequence, SequenceOf:
				file.WriteString("    xs.startTag(NS, " + local + ");\n")
				file.WriteString("    for (final XmlEncoder xe : mPayloadSequence) {\n")
				file.WriteString("      xe.Encode(xs);\n")
				file.WriteString("    }\n")
			case Choice:
				file.WriteString("    xs.startTag(NS, " + local + ");\n")
				file.WriteString("    if (mPayloadChoice != null) {\n")
				file.WriteString("      mPayloadChoice.Encode(xs);\n")
				file.WriteString("    }\n")
			case string:
				file.WriteString("    string not implemented\n")
			case Enum:
				file.WriteString("    enum not implemented\n")
			}
		}
	}
}

func java_generate_attributes_encoder(file *os.File, prefix string, fields []*Field) {
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
				file.WriteString("    if (" + java_makeIdent(prefix, x.Name) + ") {\n")
				file.WriteString("      xs.attribute(\"" + space + "\", \"" + local + "\", \"true\");\n")
			} else {
				value := java_simplevalue(prefix, x)
				file.WriteString("    if (" + java_makeIdent(prefix, x.Name) + " != null) {\n")
				file.WriteString("      xs.attribute(\"" + space + "\", \"" + local + "\", " + value + ");\n")
			}
			file.WriteString("    }\n")
		}
	}
}

func java_generate_element_encoder(file *os.File, ident string, target *Target, prefix string, field *Field) {
	if field.EncodingRule != nil {
		space := "NS"
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
			file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
			file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
		case "element:cdata":
			if _, ok := field.Type.(SequenceOf); ok {
				file.WriteString(ident + "for (final String str : " + java_makeIdent(prefix, field.Name) + ") {\n")
				file.WriteString(ident + "  xs.startTag(" + space + ", " + local + ");\n")
				file.WriteString(ident + "  xs.text(str);\n")
				file.WriteString(ident + "  xs.endTag(" + space + ", " + local + ");\n")
				file.WriteString(ident + "}\n")
			} else {
				value := java_simplevalue(prefix, field)
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
				file.WriteString(ident + "xs.text(" + value + ");\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
			}
		case "element:name":
			// from enum
			local = java_makeIdent(prefix, field.Name) + ".toString()"
			file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
			file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
		case "element":
			switch typ := field.Type.(type) {
			case string:
				name := field.Name
				if name == "" {
					name = typ
				}
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
				file.WriteString(ident + java_makeIdent(prefix, name) + ".Encode(xs);\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
			case SequenceOf, Sequence:
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
				file.WriteString(ident + "for (final XmlEncoder xe : " +
					java_makeIdent(prefix, field.Name) + ") {\n")
				file.WriteString(ident + "  xe.Encode(xs);\n")
				file.WriteString(ident + "}\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
			case Extension:
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
				file.WriteString(ident + java_makeIdent(prefix, field.Name) + ".Encode(xs);\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
			case Choice:
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
				file.WriteString(ident + java_makeIdent(prefix, field.Name) + ".Encode(xs);\n")
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
			case Enum:
			case Set:
				fields := []*Field(typ)
				space := "NS"
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
						local = java_makeIdent(prefix, x.Name) + ".toString()"
					}
				}
				file.WriteString(ident + "xs.startTag(" + space + ", " + local + ");\n")
				java_generate_attributes_encoder(file, prefix, fields)
				for _, x := range fields {
					java_generate_element_encoder(file, ident +
						"  ", target, prefix + java_makeClassName(field.Name), x)
				}
				file.WriteString(ident + "xs.endTag(" + space + ", " + local + ");\n")
			}
		}
	} else {
		switch typ := field.Type.(type) {
		case string:
			name := field.Name
			if name == "" {
				name = typ
			}
			file.WriteString(ident + java_makeIdent(prefix, name) + ".Encode(xs);\n")
		case SequenceOf, Sequence:
			file.WriteString(ident + "for (final XmlEncoder xe : " +
				java_makeIdent(prefix, field.Name) + ") {\n")
			file.WriteString(ident + "  xe.Encode(xs);\n")
			file.WriteString(ident + "}\n")
		case Extension, Choice:
			file.WriteString(ident + java_makeIdent(prefix, field.Name) + ".Encode(xs);\n")
		case Enum:
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				java_generate_element_encoder(file, ident + "  ", target, prefix, x)
			}
		}
	}
}

func java_simplevalue(prefix string, field *Field) string {
	switch typ := field.Type.(type) {
	case string:
		switch typ {
		case "xmllang", "string":
			return java_makeIdent(prefix, field.Name)
//		case "jid":
//		case "uint", "int":
//		case "boolean":
//		case "datetime":
		}
	}
	return java_makeIdent(prefix, field.Name) + ".toString()"
}

func java_generate_decoder(file *os.File, target *Target, field *Field) {
	file.WriteString("  public void Decode(XmlPullParser xp) throws IOException, XmlPullParserException {\n")
	ident := "    "
	prefix := ""
	java_generate_element_decoder(file, target, 0, false, false, ident, prefix, field)
	file.WriteString("  }\n")
}

func java_generate_element_decoder(file *os.File, target *Target, depth int, elseif bool, decl bool,
	ident, prefix string, field *Field) {
	depth++
	if field.EncodingRule != nil {
		switch field.EncodingRule.Type {
		case "element:cdata":
			file.WriteString(ident + "String data = ru.sulci.sb.xmlencoder.XStream.getText(xp);\n")
			if _, ok := field.Type.(SequenceOf); ok {
				file.WriteString(ident + "if (" + java_makeIdent(prefix, field.Name) + " == null) {\n")
				file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = new ArrayList();\n")
				file.WriteString(ident + "}\n")
				file.WriteString(ident + java_makeIdent(prefix, field.Name) + ".add(data);\n")
			} else {
				java_simplevalue_decode(file, ident, prefix, "data", field)
				// file.WriteString(ident + java_makeIdent(prefix, field.Name) + " = data;\n")
			}
		case "element:bool":
			file.WriteString(ident + java_makeIdent(prefix, field.Name) + " = true;\n")
			file.WriteString(ident + "ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
		case "element:name":
			file.WriteString(ident + "String name = xp.getName();\n")
			file.WriteString(ident + java_makeIdent(prefix, field.Name) + ".fromString(name);\n")
		case "startelement", "element":
			if depth == 1 {
				file.WriteString(ident + "int ev;\n")
			}
			switch typ := field.Type.(type) {
			case string:
				if typ == "langstring" {
					file.WriteString(ident + "if (" + java_makeIdent(prefix, field.Name) + " == null) {\n")
					file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = new LangString();\n")
					file.WriteString(ident + "}\n")
					file.WriteString(ident + java_makeIdent(prefix, field.Name) + ".Decode(xp);\n")
				} else {
					file.WriteString("not implemented\n")
				}
			case Extension:
				file.WriteString(ident + "for (;;) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next();\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break;\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				if typ.Local == "" {
					file.WriteString(forident + "  XmlEncoder obj = ru.sulci.sb.xmlencoder.XStream.getExtension(xp.getNamespace(), xp.getName());\n")
				} else {
					space := "\"" + typ.Space + "\""
					local := "\"" + typ.Local + "\""
					file.WriteString(forident + "  if (xp.getNamespace() == " + space + " && xp.getName() == " +
						local + ") {\n")
					file.WriteString(forident + "    XmlEncoder obj = ru.sulci.sb.xmlencoder.XStream.getExtension(" + space + ", " + local + ");\n")
				}
				file.WriteString(forident + "    if (obj != null) {\n")
				file.WriteString(forident + "      obj.Decode(xp);\n")
				if typ.Local != "" {
					file.WriteString(forident + "      mPayload = (" +
						java_getExtensionType(typ.Space, typ.Local) + ") obj;\n")
				} else {
					file.WriteString(forident + "      mPayload = obj;\n")
				}
				file.WriteString(forident + "    } else {\n")
				file.WriteString(forident + "      ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
				file.WriteString(forident + "    }\n")
				if typ.Local != "" {
					file.WriteString(forident + "  } else {\n")
					file.WriteString(forident + "    ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
					file.WriteString(forident + "  }\n")
				}					
				file.WriteString(forident + "}\n")
				file.WriteString(ident + "}\n")
			case SequenceOf:
				file.WriteString(ident + "for (;;) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next();\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break;\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				if string(typ) == "extension" {
					file.WriteString(forident + "XmlEncoder obj = ru.sulci.sb.xmlencoder.XStream.getExtension(xp.getNamespace(), xp.getName());\n")
					file.WriteString(forident + "if (obj != null) {\n")
				} else {
					var tname string
					space := "NS"
					local := field.Name
					if t, ok := javaSimpleTypes[string(typ)]; ok {
						tname = t.Type
					} else {
						f := java_getFieldByName(target, string(typ))
						tname = java_makeClassName(f.Name)
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
					file.WriteString(forident + "  " + tname + " obj = new " + tname + "();\n")
				}
				file.WriteString(forident + "  obj.Decode(xp);\n")
				var name string
				if depth == 1 {
					name = "mPayloadSequence"
				} else {
					name = java_makeIdent(prefix, field.Name)
				}
				file.WriteString(forident + "if (" + name + " == null) {\n")
				file.WriteString(forident + "  " + name + " = new ArrayList();\n")
				file.WriteString(forident + "}\n")
				file.WriteString(forident + name + ".add(obj);\n")
				file.WriteString(forident + "} else {\n")
				file.WriteString(forident + "  ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
				file.WriteString(forident + "  }\n")
				file.WriteString(forident + "}\n")
				file.WriteString(ident + "}\n")
			case Sequence:
				file.WriteString(ident + "for (;;) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next();\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break;\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				fields := []*Field(typ)
				for _, x := range fields {
					name := x.Type.(string)
					f := java_getFieldByName(target, name)
					space := "NS"
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
					file.WriteString(forident + "    " + java_makeClassName(name) + " obj = new " +
						java_makeClassName(name) + "();\n")
					file.WriteString(forident + "    obj.Decode(xp);\n")
					if depth == 1 {
						file.WriteString(forident + "    mPayloadSequence.add(obj);\n")
					} else {
						file.WriteString(forident + "    " + java_makeIdent(prefix, field.Name) + ".add(obj);\n")
					}
				}
				file.WriteString(forident + "  }\n")
				file.WriteString(forident + "}\n")
				file.WriteString(ident + "}\n")
			case Choice:
				file.WriteString(ident + "for (;;) {\n")
				forident := ident + "  "
				file.WriteString(forident + "ev = xp.next();\n")
				file.WriteString(forident + "if (ev == XmlPullParser.END_TAG) {\n")
				file.WriteString(forident + "  break;\n")
				file.WriteString(forident + "} else if (ev == XmlPullParser.START_TAG) {\n")
				fields := []*Field(typ)
				elseif := false
				for _, x := range fields {
					name := x.Type.(string)
					f := java_getFieldByName(target, name)
					space := "NS"
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
					file.WriteString(forident + "    " + java_makeClassName(name) + " obj = new " +
						java_makeClassName(name) + "();\n")
					file.WriteString(forident + "    obj.Decode(xp);\n")
					if depth == 1 {
						file.WriteString(forident + "    mPayloadChoice = obj;\n")
					} else {
						file.WriteString(forident + "    " + java_makeIdent(prefix, field.Name) + " = obj;\n")
					}
				}
				file.WriteString(forident + "  }\n")
				file.WriteString(forident + "}\n")
				file.WriteString(ident + "}\n")
			case Set:
				fields := []*Field(typ)
				if len(fields) == 0 {
					file.WriteString(ident + "ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
					return 
				}
				if depth > 1 {
					prefix += java_makeClassName(field.Name)
				}
				decl = java_generate_attributes_decoder(file, decl, ident, prefix, fields)
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
					file.WriteString(ident + java_makeIdent(prefix, anyname.Name) + " = " +
						java_makeEnumName(anyname.Name) + ".fromString(xp.getName());\n")
				}
				if len(elems) > 0 || any != nil {
					file.WriteString(ident + "for (;;) {\n")
					file.WriteString(ident + "  ev = xp.next();\n")
					file.WriteString(ident + "  if (ev == XmlPullParser.END_TAG) {\n")
					file.WriteString(ident + "    break;\n")
					file.WriteString(ident + "  } else if (ev == XmlPullParser.START_TAG) {\n")
					forident := ident + "  "
					elseif := false
					for _, z := range elems {
						if z.EncodingRule != nil {
							space := "NS"
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
							java_generate_element_decoder(file, target, depth, elseif, decl,
								forident + "    ", prefix, z)
						} else {
							java_generate_element_decoder(file, target, depth, elseif, decl,
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
						file.WriteString("if (xp.getNamespace() == NS) {\n")
						switch any.EncodingRule.Type {
						case "element:name":
							var typ string
							if t, ok := any.Type.(string); ok {
								typ = java_makeClassName(t)
							} else {
								typ = java_makeEnumName(any.Name)
							}
							file.WriteString("  " + java_makeIdent(prefix, any.Name) + " = " +
								typ + ".fromString(xp.getName());\n")
						case "name":
							file.WriteString("not implemented;\n")
						case "element":
							java_generate_element_decoder(file, target, depth, elseif, decl,
								forident + "  ", prefix, any)
						}
					}
					if extension != nil {
						file.WriteString(forident + "  ")
						if elseif {
							file.WriteString("} else ")
						}
						file.WriteString("if (xp.getNamespace() != NS) {\n")
						java_generate_element_decoder(file, target, depth, elseif, decl,
								forident + "    ", prefix, extension)
						elseif = true
					}
					if len(elems) > 0 || any != nil || extension != nil {
						file.WriteString(forident + "  } else {\n")
						file.WriteString(forident + "    ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
						file.WriteString(forident + "  }\n")
					}
					file.WriteString(ident + "  }\n")
					file.WriteString(ident + "}\n")
				}
				if cdata != nil {
					file.WriteString(ident + "String data;\n")
					file.WriteString(ident + "data = ru.sulci.sb.xmlencoder.XStream.getText(xp);\n")
					java_simplevalue_decode(file, ident, prefix, "data", cdata)
				}
			}
		}
	} else {
		switch typ := field.Type.(type) {
		case string:
			name := field.Name
			f := java_getFieldByName(target, typ)
			if name == "" {
				name = typ
			}
			space := "NS"
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
			file.WriteString(ident + "  " + java_makeIdent(prefix, name) + " = new " +
				java_makeClassName(typ) + "();\n")
			file.WriteString(ident + "  " + java_makeIdent(prefix, name) + ".Decode(xp);\n")
		case SequenceOf:
			if string(typ) == "extension" {
				file.WriteString(ident + "XmlEncoder xe = ru.sulci.sb.xmlencoder.XStream.getExtension(xp.getNamespace(), xp.getName());\n")
				file.WriteString(ident + "if (xe != null) {\n")
				file.WriteString(ident + "  xe.Decode(xp);\n")
				file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + ".add(xe);\n")
				file.WriteString(ident + "} else {\n")
				file.WriteString(ident + "  ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
				file.WriteString(ident + "}\n")
			} else {
				f := java_getFieldByName(target, string(typ))
				space := "NS"
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
				file.WriteString(ident + "  " + java_makeClassName(string(typ)) + " obj = new " +
					java_makeClassName(string(typ)) + "();\n")
				file.WriteString(ident + "  obj.Decode(xp);\n")
				file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + ".add(obj);\n")
			}
		case Extension:
			if typ.Local == "" {
				file.WriteString(ident + "XmlEncoder xe = ru.sulci.sb.xmlencoder.XStream.getExtension(xp.getNamespace(), xp.getName());\n")
				file.WriteString(ident + "if (xe != null) {\n")
				file.WriteString(ident + "  xe.Decode(xp);\n")
				file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = xe;\n")
				file.WriteString(ident + "} else {\n")
				file.WriteString(ident + "  ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
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
				file.WriteString(ident + "  XmlEncoder xe = ru.sulci.sb.xmlencoder.XStream.getExtension(" +
					space + ", " + local + ");\n")
				file.WriteString(ident + "  if (xe != null) {\n")
				file.WriteString(ident + "    xe.Decode(xp);\n")
				file.WriteString(ident + "    " + java_makeIdent(prefix, field.Name) + " = (" +
					java_getExtensionType(typ.Space, typ.Local) + ") xe;\n")
				file.WriteString(ident + "  } else {\n")
				file.WriteString(ident + "    ru.sulci.sb.xmlencoder.XStream.getEndTag(xp);\n")
				file.WriteString(ident  + "  }\n")
			}
		case Sequence:
			fields := []*Field(typ)
			for _, x := range fields {
				f := x
				if x.Name == "" {
					f = java_getFieldByName(target, x.Type.(string))
				}
				space := "NS"
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
				file.WriteString(ident + "  " + java_makeClassName(f.Name) + " obj = new " +
					java_makeClassName(f.Name) + "();\n")
				file.WriteString(ident + "  obj.Decode(xp);\n")
				file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + ".add(obj);\n")
			}
		case Choice:
			fields := []*Field(typ)
			for _, x := range fields {
				f := java_getFieldByName(target, x.Type.(string))
				space := "NS"
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
				file.WriteString(ident + "  " + java_makeClassName(f.Name) + " obj = new " +
					java_makeClassName(f.Name) + "();\n")
				file.WriteString(ident + "  obj.Decode(xp);\n")
				file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = obj;\n")
				if !elseif {
					elseif = true
				}
			}
		case Set:
			fields := []*Field(typ)
			for _, x := range fields {
				f := x
				if x.Name == "" {
					f = java_getFieldByName(target, x.Type.(string))
				}
				space := "NS"
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
				java_generate_element_decoder(file, target, depth, elseif, decl,
					ident + "    ", prefix + java_makeClassName(field.Name), x)
/*					
				file.WriteString(ident + "  " + java_makeClassName(f.Name) + " obj = new " +
					java_makeClassName(f.Name) + "();\n")
				file.WriteString(ident + "  obj.Decode(xp);\n")
				name := x.Name
				if name == "" {
					name = x.Type.(string)
				}
				pr := prefix + java_makeClassName(field.Name)
				file.WriteString(ident + "  " + java_makeIdent(pr, name) + " = obj;\n")
*/
				if !elseif {
					elseif = true
				}
			}
			
		}
	}
}

func java_getFieldByName(target *Target, f string) *Field {
	for _, x := range target.Fields {
		if x.Name == f {
			return x
		}
	}
	return nil
}

func java_generate_attributes_decoder(file *os.File, decl bool, ident, prefix string, fields []*Field) (ret bool) {
	hasAttrs := false
	for _, x := range fields {
		if x.EncodingRule != nil && x.EncodingRule.Type == "attribute" {
			hasAttrs = true
			ret = true
			break
		}
	}
	if hasAttrs && !decl {
		file.WriteString(ident + "String value;\n")
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
			file.WriteString(ident + "value = xp.getAttributeValue(" + space + ", " + local + ");\n")
			file.WriteString(ident + "if (value != null) {\n")
			java_simplevalue_decode(file, ident + "  ", prefix, "value", x)
			file.WriteString(ident + "}\n")
		}
	}
	return
}

func java_simplevalue_decode(file *os.File, ident, prefix, varname string, field *Field) {
	switch typ := field.Type.(type) {
	case string:
		switch typ {
		case "boolean":
			file.WriteString(ident + java_makeIdent(prefix, field.Name) + " = Boolean.valueOf(" +
				varname + ");\n")
			return
		case "int", "uint":
			file.WriteString(ident + java_makeIdent(prefix, field.Name) + " = Integer.valueOf(" +
				varname + ");\n")
			return
		case "bytestring", "string", "xmllang":
			file.WriteString(ident + java_makeIdent(prefix, field.Name) + " = " + varname + ";\n")
			return
		case "jid":
			file.WriteString(ident + java_makeIdent(prefix, field.Name) + " = new JID(" +
				varname + ");\n")
			return
		case "datetime":
			file.WriteString(ident + "try {\n")
			file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) +
				" = ru.sulci.sb.xmlencoder.DateTime.parseRFC3339(" + varname + ");\n")
			file.WriteString(ident + "} catch (java.text.ParseException e) {}\n")
			return
		}
	case Enum:
		file.WriteString(ident + "try {\n")
		file.WriteString(ident + "  " + java_makeIdent(prefix, field.Name) + " = " +
			java_makeEnumName(field.Name) + ".fromString(" + varname + ");\n")
		file.WriteString(ident + "} catch (IllegalArgumentException iae) {}\n")
		return
	}
	file.WriteString("not implemented\n")
}

func java_generate_adders() error {
	filename := filepath.Join(cfg.Java.Outdir, "xmlencoder", "Extensions.java")
	file, err := os.Create(filename)
	if err != nil { return err }
	file.WriteString("package ru.sulci.sb.xmlencoder;\n\n")
	file.WriteString("import ru.sulci.sb.xmlencoder.QName;\n")
	file.WriteString("import java.util.HashMap;\n")
	file.WriteString("import java.util.Map;\n")
	file.WriteString("public class Extensions {\n")
	file.WriteString("  final static Map<QName, Data> mapData;\n")
	file.WriteString("  static {\n")
	file.WriteString("    mapData = new HashMap<QName, Data>();\n")
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
					file.WriteString("    mapData.put(new QName(\"" + target.Space + "\", \"" + local + "\"), ")
					file.WriteString(" new Data(ru.sulci.sb.")
					if category == "extension" {
						file.WriteString("extensions.")
					}
					file.WriteString(schema.PackageName + ".")
					if target.Name != "" {
						file.WriteString(target.Name + ".")
					}
					file.WriteString(java_makeClassName(field.Name) + ".class")
					switch field.Reciver_type {
					case "both":
						file.WriteString(", true, true)")
					case "server":
						file.WriteString(", true, false)")
					case "client":
						file.WriteString(", false, true)")
					}
					file.WriteString(");\n")
				}
			}
		}
	}
	file.WriteString(";\n")
	file.WriteString("  }\n")
	file.WriteString("  final static class Data {\n")
	file.WriteString("    Class className;\n")
	file.WriteString("    Boolean forServer;\n")
	file.WriteString("    Boolean forClient;\n")
	file.WriteString("    Boolean Use;\n")
	file.WriteString("\n")
	file.WriteString("    public Data(Class className, Boolean forServer, Boolean forClient) {\n")
	file.WriteString("      this.className = className;\n")
	file.WriteString("      this.forServer = forServer;\n")
	file.WriteString("      this.forClient = forClient;\n")
	file.WriteString("      this.Use = false;\n")
	file.WriteString("    }\n")
	file.WriteString("  }\n")
	file.WriteString("  public static Object getExtension(String space, String local) {\n")
	file.WriteString("    Object obj = null;\n")
	file.WriteString("    Data data = mapData.get(new QName(space, local));\n")
	file.WriteString("    if (data != null) {\n")
	file.WriteString("      try {\n")
	file.WriteString("        obj = data.className.newInstance();\n")
	file.WriteString("      } catch (Exception e) {}\n")
	file.WriteString("    }\n")
	file.WriteString("    return obj;\n")
	file.WriteString("  }\n")
	file.WriteString("}\n")
	file.Close()
	return nil
}
