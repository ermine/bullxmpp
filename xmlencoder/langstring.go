package xmlencoder

import (
	"encoding/xml"
)

type LangString map[string]string

func (l LangString) Get(lang string) string {
	if t, ok := l[lang]; ok { return t }
	if lang != "" {
		if t, ok := l[""]; ok { return t }
	}
	if lang != "en" {
		if t, ok := l["en"]; ok { return t }
	}
	if len(l) > 0 {
		for _, v := range l {
			return v
		}
	}
	return ""
}

func (l *LangString) Decode(d *Decoder, tag *xml.StartElement) error {
	lang := AttributeValue(tag.Attr, ns_xml, "lang")
	text, err := d.Text()
	if err != nil { return err }
	(*l)[lang] = text
	return nil
}

func (l *LangString) Encode(e *Encoder, space, local string) error {
	var err error
	for lang, value := range *l {
		if err = e.StartElement(space, local); err != nil { return err }
		if err = e.Attribute(ns_xml, "lang", lang); err != nil { return err }
		if err = e.Text(value); err != nil { return err }
		if err = e.EndElement(); err != nil { return err }
	}
	return nil
}
