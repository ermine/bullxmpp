package starttls

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "urn:ietf:params:xml:ns:xmpp-tls"
type Starttls struct {
  Required bool
}
type Proceed struct {
}
type Failure struct {
}
func (elm *Starttls) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "starttls"); err != nil { return err }
if elm.Required {
if err = e.StartElement(NS, "required"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Starttls) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "required":
elm.Required = true
if err = d.Skip(); err != nil { return err }
continue
}
}
}
return err
}

func (elm *Proceed) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "proceed"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Proceed) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
if err = d.Skip(); err != nil { return err }
return err
}

func (elm *Failure) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "failure"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Failure) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
if err = d.Skip(); err != nil { return err }
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "starttls"}, Starttls{}, true, true)
 xmlencoder.AddExtension(xml.Name{NS, "proceed"}, Proceed{}, true, false)
 xmlencoder.AddExtension(xml.Name{NS, "failure"}, Failure{}, true, false)
}
