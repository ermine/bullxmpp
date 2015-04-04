package captcha

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
import "jabber.ru/xmpp/xdata"
const NS = "urn:xmpp:captcha"
type Captcha struct {
  Xdata *xdata.X

}
func (elm *Captcha) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "captcha"); err != nil { return err }
if elm.Xdata != nil {
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Captcha) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "captcha"}, Captcha{}, true, true)
}
