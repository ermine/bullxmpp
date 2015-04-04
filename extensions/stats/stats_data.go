package stats

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "http://jabber.org/protocol/stats"
type Stats []Stat
type Stat struct {
  Name *string
  Units *string
  Value *string
}
func (elm *Stats) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
for _, x := range *elm {
if err = x.Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Stats) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.StartElement:
if t.Name.Space == NS && t.Name.Local == "stat" {
newel := &Stat{}
if err = newel.Decode(d, &t); err != nil { return err }
elm = append(elm, newel)
}
case xml.EndElement:
break Loop
}
}
return err
}

func (elm *Stat) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "stat"); err != nil { return err }
if elm.Name != nil {
if err = e.Attribute("", "name", *elm.Name); err != nil { return err }
}
if elm.Units != nil {
if err = e.Attribute("", "units", *elm.Units); err != nil { return err }
}
if elm.Value != nil {
if err = e.Attribute("", "value", *elm.Value); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Stat) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "name":
elm.Name = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "units":
elm.Units = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "value":
elm.Value = xmlencoder.Copystring(x.Value)
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Stats{}, true, true)
}
