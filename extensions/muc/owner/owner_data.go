package owner

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
import "strconv"
import "time"
import "jabber.ru/xmpp/jid"
import "jabber.ru/xmpp/xdata"
const NS = "http://jabber.org/protocol/muc#owner"
type Configure xdata.X
func (elm *Configure) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if err = elm.Encode(e); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Configure) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.StartElement:
if t.Name.Space == "jabber:x:data" && t.Name.Local == "x"{
newel := &xdata.X{}
if err = newel.Decode(d, &t); err != nil { return err }
*elm = Configure(*newel)
} else {
if err = d.Skip(); err != nil { return err }
}
case xml.EndElement:
break Loop
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Configure{}, true, true)
}
