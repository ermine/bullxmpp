package muc

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
import "strconv"
import "time"
import "github.com/ermine/bullxmpp/jid"
import "github.com/ermine/bullxmpp/xdata"
const NS = "http://jabber.org/protocol/muc"
type Enter struct {
  History struct {
  Maxchars *int
  Maxstanzas *int
  Seconds *int
  Since *time.Time
}

  Password *string
}
func (elm *Enter) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "x"); err != nil { return err }
if err = e.StartElement(NS, "history"); err != nil { return err }
if elm.History.Maxchars != nil {
if err = e.Attribute("", "maxchars", strconv.FormatInt(int64(*elm.History.Maxchars), 10)); err != nil { return err }
}
if elm.History.Maxstanzas != nil {
if err = e.Attribute("", "maxstanzas", strconv.FormatInt(int64(*elm.History.Maxstanzas), 10)); err != nil { return err }
}
if elm.History.Seconds != nil {
if err = e.Attribute("", "seconds", strconv.FormatInt(int64(*elm.History.Seconds), 10)); err != nil { return err }
}
if elm.History.Since != nil {
if err = e.Attribute("", "since", elm.History.Since.String()); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
if elm.Password != nil {
if err = e.SimpleElement(NS, "password", *elm.Password); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Enter) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "history":
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "maxchars":
var i int64
i, err = strconv.ParseInt(x.Value, 10, 0)
if err == nil {
*elm.History.Maxchars = int(i)
}
case x.Name.Space == "" && x.Name.Local == "maxstanzas":
var i int64
i, err = strconv.ParseInt(x.Value, 10, 0)
if err == nil {
*elm.History.Maxstanzas = int(i)
}
case x.Name.Space == "" && x.Name.Local == "seconds":
var i int64
i, err = strconv.ParseInt(x.Value, 10, 0)
if err == nil {
*elm.History.Seconds = int(i)
}
case x.Name.Space == "" && x.Name.Local == "since":
*elm.History.Since, err = time.Parse(time.RFC3339, x.Value)
if err != nil { return err }
}
}
case t.Name.Space == NS && t.Name.Local == "password":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Password = s
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "x"}, Enter{}, true, true)
}
