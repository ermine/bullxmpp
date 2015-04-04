package session

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "urn:ietf:params:xml:ns:xmpp-session"
type Session struct {
}
func (elm *Session) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "session"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Session) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
if err = d.Skip(); err != nil { return err }
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "session"}, Session{}, true, true)
}
