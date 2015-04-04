package bind

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
import "jabber.ru/xmpp/jid"
const NS = "urn:ietf:params:xml:ns:xmpp-bind"
type Bind struct {
  Resource *string
  Jid *jid.JID
}
func (elm *Bind) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "bind"); err != nil { return err }
if elm.Resource != nil {
if err = e.SimpleElement(NS, "resource", *elm.Resource); err != nil { return err }
}
if elm.Jid != nil {
if err = e.SimpleElement(NS, "jid", elm.Jid.String()); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Bind) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "resource":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Resource = s
case t.Name.Space == NS && t.Name.Local == "jid":
var s string
if s, err = d.Text(); err != nil { return err }
var j *jid.JID
if j, err = jid.New(s); err != nil { return err }
elm.Jid = j
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "bind"}, Bind{}, true, true)
}
