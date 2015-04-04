package xdata

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "jabber:x:data"
type X struct {
  Type *XType
  Title *string
  Reported []*Field
  Fields []interface{}
}
type Field struct {
  Label *string
  Type *FieldType
  Var *string
  Desc *string
  Required bool
  Value *string
  Option []*Option
}
type Option struct {
  Label *string
  Value *string
}
type XType string
const (
XTypeCancel XType = "cancel"
XTypeForm XType = "form"
XTypeResult XType = "result"
XTypeSubmit XType = "submit"
)
type FieldType string
const (
FieldTypeBoolean FieldType = "boolean"
FieldTypeFixed FieldType = "fixed"
FieldTypeHidden FieldType = "hidden"
FieldTypeJidMulti FieldType = "jid-multi"
FieldTypeJidSingle FieldType = "jid-single"
FieldTypeListMulti FieldType = "list-multi"
FieldTypeListSingle FieldType = "list-single"
FieldTypeTextMulti FieldType = "text-multi"
FieldTypeTextPrivate FieldType = "text-private"
FieldTypeTextSingle FieldType = "text-single"
)
func (elm *X) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "x"); err != nil { return err }
if elm.Type != nil {
if err = e.Attribute("", "type", string(*elm.Type)); err != nil { return err }
}
if elm.Title != nil {
if err = e.SimpleElement(NS, "title", *elm.Title); err != nil { return err }
}
if err = e.StartElement(NS, "reported"); err != nil { return err }
for _, x := range elm.Reported {
if err = x.Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
for _, x := range elm.Fields {
if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *X) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "type":
value := XType(x.Value)
elm.Type = &value
}
}
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
case t.Name.Space == NS && t.Name.Local == "title":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Title = s
case t.Name.Space == NS && t.Name.Local == "reported":
var t xml.Token
InLoop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.StartElement:
if t.Name.Space == NS && t.Name.Local == "field" {
newel := &Field{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Reported = append(elm.Reported, newel)
}
case xml.EndElement:
break InLoop
}
}
}
}
}
return err
}

func (elm *Field) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "field"); err != nil { return err }
if elm.Label != nil {
if err = e.Attribute("", "label", *elm.Label); err != nil { return err }
}
if elm.Type != nil {
if err = e.Attribute("", "type", string(*elm.Type)); err != nil { return err }
}
if elm.Var != nil {
if err = e.Attribute("", "var", *elm.Var); err != nil { return err }
}
if elm.Desc != nil {
if err = e.SimpleElement(NS, "desc", *elm.Desc); err != nil { return err }
}
if elm.Required {
if err = e.StartElement(NS, "required"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.Value != nil {
if err = e.SimpleElement(NS, "value", *elm.Value); err != nil { return err }
}
for _, x := range elm.Option {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Field) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "label":
elm.Label = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "type":
value := FieldType(x.Value)
elm.Type = &value
case x.Name.Space == "" && x.Name.Local == "var":
elm.Var = xmlencoder.Copystring(x.Value)
}
}
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
case t.Name.Space == NS && t.Name.Local == "desc":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Desc = s
case t.Name.Space == NS && t.Name.Local == "required":
elm.Required = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "value":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Value = s
case t.Name.Space == NS && t.Name.Local == "option":
newel := &Option{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Option = append(elm.Option, newel)
}
}
}
return err
}

func (elm *Option) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "option"); err != nil { return err }
if elm.Label != nil {
if err = e.Attribute("", "label", *elm.Label); err != nil { return err }
}
if elm.Value != nil {
if err = e.SimpleElement(NS, "value", *elm.Value); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Option) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "label":
elm.Label = xmlencoder.Copystring(x.Value)
}
}
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
case t.Name.Space == NS && t.Name.Local == "value":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Value = s
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "x"}, X{}, true, true)
}
