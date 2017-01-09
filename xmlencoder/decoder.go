package xmlencoder

import (
	"encoding/xml"
	"reflect"
	"io"
)

type Decoder struct {
	*xml.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{xml.NewDecoder(r)}
}

func (d *Decoder) Decode() (interface{}, error) {
	var err error
	var t xml.Token
	for {
		if t, err = d.Token(); err != nil { return nil, err }
		switch t := t.(type) {
		case xml.StartElement:
			if r, ok := GetExtension(t.Name); ok {
				if err = r.(Extension).Decode(d, &t); err != nil {
					return nil, err }
				return r, nil
			} else {
				if err = d.Skip(); err != nil { return nil, err }
			}
		case xml.EndElement:
			return nil, nil
		}
	}
	return nil, nil
}

func (d Decoder) Text() (string, error) {
	var err error
	var text string
	var t xml.Token
Loop:
	for {
		if t, err = d.Token(); err != nil { return "", err }
		switch t := t.(type) {
		case xml.EndElement:
			break Loop
		case xml.StartElement:
			if err = d.Skip(); err != nil { return "", err }
		case xml.CharData:
			text += string(t)
		}
	}
	return text, nil
}

func (d *Decoder) Bytes() ([]byte, error) {
	var err error
	var data []byte
	var t xml.Token
Loop:
	for {
		if t, err = d.Token(); err != nil { return data, err }
		switch t := t.(type) {
		case xml.EndElement:
			break Loop
		case xml.StartElement:
			if err = d.Skip(); err != nil { return data, err }
		case xml.CharData:
			data = append(data, t.Copy()...)
		}
	}
	return data, nil
}

func AttributeValue(attrs []xml.Attr, space, local string) string {
	for _, x := range attrs {
		if x.Name.Space == space && x.Name.Local == local {
			return x.Value
		}
	}
	return ""
}

func Copystring(s string) *string {
	return &s
}

type ExtensionType struct {
	Type reflect.Type
	ForClient bool
	ForServer bool
}

var Extensions = map[xml.Name]*ExtensionType{}

func AddExtension(xmlname xml.Name, i interface{}, forServer, forClient bool) {
	Extensions[xmlname] = &ExtensionType{reflect.TypeOf(i), forServer, forClient}
}

func ReplaceExtensionStruct(xmlname xml.Name, i interface{}) {
	if t, ok := Extensions[xmlname]; ok {
		t.Type = reflect.TypeOf(i)
	}
}

func GetExtension(xmlname xml.Name) (interface{}, bool) {
	if typ, ok := Extensions[xmlname]; ok {
		return reflect.New(typ.Type).Interface(), true
	}
	return nil, false
}
