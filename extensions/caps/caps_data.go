package caps

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
const NS = "http://jabber.org/protocol/caps"
type Caps struct {
  Ext *string
  Hash *string
  Node *string
  Ver *string
}
func (elm *Caps) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "c"); err != nil { return err }
if elm.Ext != nil {
if err = e.Attribute("", "ext", *elm.Ext); err != nil { return err }
}
if elm.Hash != nil {
if err = e.Attribute("", "hash", *elm.Hash); err != nil { return err }
}
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
if elm.Ver != nil {
if err = e.Attribute("", "ver", *elm.Ver); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Caps) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "ext":
elm.Ext = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "hash":
elm.Hash = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "ver":
elm.Ver = xmlencoder.Copystring(x.Value)
}
}
return err
}

func init() {
}
