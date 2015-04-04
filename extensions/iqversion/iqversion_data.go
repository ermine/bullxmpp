package iqversion

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "jabber:iq:version"
type Version struct {
  Name *string
  Version *string
  Os *string
}
func (elm *Version) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if elm.Name != nil {
if err = e.SimpleElement(NS, "name", *elm.Name); err != nil { return err }
}
if elm.Version != nil {
if err = e.SimpleElement(NS, "version", *elm.Version); err != nil { return err }
}
if elm.Os != nil {
if err = e.SimpleElement(NS, "os", *elm.Os); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Version) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "name":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Name = s
case t.Name.Space == NS && t.Name.Local == "version":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Version = s
case t.Name.Space == NS && t.Name.Local == "os":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Os = s
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Version{}, true, true)
}
