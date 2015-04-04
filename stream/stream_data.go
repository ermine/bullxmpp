package stream

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
import "jabber.ru/xmpp/jid"
const NS = "http://etherx.jabber.org/streams"
type Start struct {
  To *jid.JID
  From *jid.JID
  Id *string
  Version *string
  Lang *string
}
type Features []interface{}
type Error struct {
  Text *xmlencoder.LangString
  Condition struct {
  Name *ErrorConditionName
  Extra *string
}

}
type ErrorConditionName string
const (
ErrorConditionNameBadFormat ErrorConditionName = "bad-format"
ErrorConditionNameBadNamespacePrefix ErrorConditionName = "bad-namespace-prefix"
ErrorConditionNameConflict ErrorConditionName = "conflict"
ErrorConditionNameConnectionTimeout ErrorConditionName = "connection-timeout"
ErrorConditionNameHostGone ErrorConditionName = "host-gone"
ErrorConditionNameHostUnknown ErrorConditionName = "host-unknown"
ErrorConditionNameImproperAddressing ErrorConditionName = "improper-addressing"
ErrorConditionNameInternalServerError ErrorConditionName = "internal-server-error"
ErrorConditionNameInvalidFrom ErrorConditionName = "invalid-from"
ErrorConditionNameInvalidNamespace ErrorConditionName = "invalid-namespace"
ErrorConditionNameInvalidXml ErrorConditionName = "invalid-xml"
ErrorConditionNameNotAuthorized ErrorConditionName = "not-authorized"
ErrorConditionNameNotWellFormed ErrorConditionName = "not-well-formed"
ErrorConditionNamePolicyViolation ErrorConditionName = "policy-violation"
ErrorConditionNameRemoteConnectionFailed ErrorConditionName = "remote-connection-failed"
ErrorConditionNameReset ErrorConditionName = "reset"
ErrorConditionNameResourceConstraint ErrorConditionName = "resource-constraint"
ErrorConditionNameRestrictedXml ErrorConditionName = "restricted-xml"
ErrorConditionNameSeeOtherHost ErrorConditionName = "see-other-host"
ErrorConditionNameSystemShutdown ErrorConditionName = "system-shutdown"
ErrorConditionNameUndefinedCondition ErrorConditionName = "undefined-condition"
ErrorConditionNameUnsupportedEncoding ErrorConditionName = "unsupported-encoding"
ErrorConditionNameUnsupportedFeature ErrorConditionName = "unsupported-feature"
ErrorConditionNameUnsupportedStanzaType ErrorConditionName = "unsupported-stanza-type"
ErrorConditionNameUnsupportedVersion ErrorConditionName = "unsupported-version"
)
func (elm *Start) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SetPrefix("stream", "http://etherx.jabber.org/streams"); err != nil { return err }
if err = e.StartElement(NS, "stream"); err != nil { return err }
if elm.To != nil {
if err = e.Attribute("", "to", elm.To.String()); err != nil { return err }
}
if elm.From != nil {
if err = e.Attribute("", "from", elm.From.String()); err != nil { return err }
}
if elm.Id != nil {
if err = e.Attribute("", "id", *elm.Id); err != nil { return err }
}
if elm.Version != nil {
if err = e.Attribute("", "version", *elm.Version); err != nil { return err }
}
if elm.Lang != nil {
if err = e.Attribute("http://www.w3.org/XML/1998/namespace", "lang", string(*elm.Lang)); err != nil { return err }
}
return nil
}

func (elm *Start) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "to":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.To = j
case x.Name.Space == "" && x.Name.Local == "from":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.From = j
case x.Name.Space == "" && x.Name.Local == "id":
elm.Id = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "version":
elm.Version = xmlencoder.Copystring(x.Value)
case x.Name.Space == "http://www.w3.org/XML/1998/namespace" && x.Name.Local == "lang":
elm.Lang = & x.Value
}
}
return err
}

func (elm *Features) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SetPrefix("stream", "http://etherx.jabber.org/streams"); err != nil { return err }
if err = e.StartElement(NS, "features"); err != nil { return err }
for _, x := range *elm {
if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Features) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.StartElement:
if newel, ok := xmlencoder.GetExtension(t.Name); ok {
if err = newel.(xmlencoder.Extension).Decode(d, &t); err != nil { return err }
*elm = append(*elm, newel)
} else {
if err = d.Skip(); err != nil { return err }
}
case xml.EndElement:
break Loop
}
}
return err
}

func (elm *Error) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SetPrefix("stream", "http://etherx.jabber.org/streams"); err != nil { return err }
if err = e.StartElement(NS, "error"); err != nil { return err }
if elm.Text != nil {
elm.Text.Encode(e, NS, "text")
}
if err = e.StartElement(NS, string(*elm.Condition.Name)); err != nil { return err }
if elm.Condition.Name != nil {
}
if elm.Condition.Extra != nil {
if err = e.Text(*elm.Condition.Extra); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Error) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "text":
if err = elm.Text.Decode(d, &t); err != nil { return err }
default:
if t.Name.Space == NS {
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Condition.Extra = s
}
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "stream"}, Start{}, true, true)
 xmlencoder.AddExtension(xml.Name{NS, "features"}, Features{}, true, true)
 xmlencoder.AddExtension(xml.Name{NS, "error"}, Error{}, true, true)
}
