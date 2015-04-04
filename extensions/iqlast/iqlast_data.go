package iqlast

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
import "strconv"
const NS = "jabber:iq:last"
type Last struct {
  Seconds *uint
  Extra *string
}
func (elm *Last) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if elm.Seconds != nil {
if err = e.Attribute("", "seconds", strconv.FormatUint(uint64(*elm.Seconds), 10)); err != nil { return err }
}
if elm.Extra != nil {
if err = e.Text(*elm.Extra); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Last) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "seconds":
var i uint64
i, err = strconv.ParseUint(x.Value, 10, 0)
if err == nil {
*elm.Seconds = uint(i)
}
}
}
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Extra = s
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Last{}, true, true)
}
