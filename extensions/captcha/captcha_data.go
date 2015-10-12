package captcha

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
import "github.com/ermine/bullxmpp/xdata"
const NS = "urn:xmpp:captcha"
type Captcha struct {
  Xdata *xdata.X

}
func (elm *Captcha) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "captcha"); err != nil { return err }
if elm.Xdata != nil {
if err = elm.Xdata.Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Captcha) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == "jabber:x:data" && t.Name.Local == "x":newel := &xdata.X{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Xdata = newel
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "captcha"}, Captcha{}, true, true)
}
