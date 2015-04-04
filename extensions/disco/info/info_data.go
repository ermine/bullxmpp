package info

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
const NS = "http://jabber.org/protocol/disco#info"
type Info struct {
  Node *string
  Identities []*Identity
  Features []*Feature
}
type Identity struct {
  Category *string
  Type *string
}
type Feature struct {
  Var *string
}
func (elm *Info) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
for _, x := range elm.Identities {
if err = x.Encode(e); err != nil { return err} 
}
for _, x := range elm.Features {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Info) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
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
case t.Name.Space == NS && t.Name.Local == "identity":
newel := &Identity{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Identities = append(elm.Identities, newel)
case t.Name.Space == NS && t.Name.Local == "feature":
newel := &Feature{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Features = append(elm.Features, newel)
}
}
}
return err
}

func (elm *Identity) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "identity"); err != nil { return err }
if elm.Category != nil {
if err = e.Attribute("", "category", *elm.Category); err != nil { return err }
}
if elm.Type != nil {
if err = e.Attribute("", "type", *elm.Type); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Identity) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "category":
elm.Category = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "type":
elm.Type = xmlencoder.Copystring(x.Value)
}
}
return err
}

func (elm *Feature) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "feature"); err != nil { return err }
if elm.Var != nil {
if err = e.Attribute("", "var", *elm.Var); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Feature) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "var":
elm.Var = xmlencoder.Copystring(x.Value)
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Info{}, true, true)
}
