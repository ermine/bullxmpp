package iqtime

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
const NS = "jabber:iq:time"
type Time struct {
  Utc *string
  Tz *string
  Display *string
}
func (elm *Time) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if elm.Utc != nil {
if err = e.SimpleElement(NS, "utc", *elm.Utc); err != nil { return err }
}
if elm.Tz != nil {
if err = e.SimpleElement(NS, "tz", *elm.Tz); err != nil { return err }
}
if elm.Display != nil {
if err = e.SimpleElement(NS, "display", *elm.Display); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Time) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "utc":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Utc = s
case t.Name.Space == NS && t.Name.Local == "tz":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Tz = s
case t.Name.Space == NS && t.Name.Local == "display":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Display = s
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Time{}, true, true)
}
