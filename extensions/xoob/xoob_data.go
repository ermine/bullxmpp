package xoob

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
const NS = "jabber:x:oob"
type X struct {
  Url *string
  Desc *string
}
func (elm *X) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "x"); err != nil { return err }
if elm.Url != nil {
if err = e.SimpleElement(NS, "url", *elm.Url); err != nil { return err }
}
if elm.Desc != nil {
if err = e.SimpleElement(NS, "desc", *elm.Desc); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *X) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
case t.Name.Space == NS && t.Name.Local == "url":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Url = s
case t.Name.Space == NS && t.Name.Local == "desc":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Desc = s
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "x"}, X{}, true, true)
}
