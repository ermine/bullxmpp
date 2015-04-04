package ping

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "urn:xmpp:ping"
type Ping struct {
}
func (elm *Ping) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "ping"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Ping) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
if err = d.Skip(); err != nil { return err }
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "ping"}, Ping{}, true, true)
}
