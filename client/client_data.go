package client

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
import "github.com/ermine/bullxmpp/jid"
import "github.com/ermine/bullxmpp/stanza"
import "strconv"
const NS = "jabber:client"
type Iq struct {
  From *jid.JID
  To *jid.JID
  Id *string
  Type *IqType
  Lang *string
  Payload interface{}

  Error *stanza.Error

}
type Presence struct {
  From *jid.JID
  To *jid.JID
  Id *string
  Type *PresenceType
  Lang *string
  Show *PresenceShow
  Status *string
  Priority *int
  X []interface{}
  Error *stanza.Error

}
type Message struct {
  From *jid.JID
  To *jid.JID
  Id *string
  Type *MessageType
  Lang *string
  Thread *string
  Subject *xmlencoder.LangString
  Body *xmlencoder.LangString
  X []interface{}
  Error *stanza.Error

}
type IqType string
const (
IqTypeGet IqType = "get"
IqTypeSet IqType = "set"
IqTypeResult IqType = "result"
IqTypeError IqType = "error"
)
type PresenceType string
const (
PresenceTypeSubscribe PresenceType = "subscribe"
PresenceTypeSubscribed PresenceType = "subscribed"
PresenceTypeUnsubscribe PresenceType = "unsubscribe"
PresenceTypeUnsubscribed PresenceType = "unsubscribed"
PresenceTypeUnavailable PresenceType = "unavailable"
)
type PresenceShow string
const (
PresenceShowChat PresenceShow = "chat"
PresenceShowAway PresenceShow = "away"
PresenceShowXa PresenceShow = "xa"
PresenceShowDnd PresenceShow = "dnd"
)
type MessageType string
const (
MessageTypeNormal MessageType = "normal"
MessageTypeChat MessageType = "chat"
MessageTypeGroupchat MessageType = "groupchat"
MessageTypeHeadline MessageType = "headline"
)
func (elm *Iq) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "iq"); err != nil { return err }
if elm.From != nil {
if err = e.Attribute("", "from", elm.From.String()); err != nil { return err }
}
if elm.To != nil {
if err = e.Attribute("", "to", elm.To.String()); err != nil { return err }
}
if elm.Id != nil {
if err = e.Attribute("", "id", *elm.Id); err != nil { return err }
}
if elm.Type != nil {
if err = e.Attribute("", "type", string(*elm.Type)); err != nil { return err }
}
if elm.Lang != nil {
if err = e.Attribute("http://www.w3.org/XML/1998/namespace", "lang", string(*elm.Lang)); err != nil { return err }
}
if elm.Payload != nil {
if err = elm.Payload.(xmlencoder.Extension).Encode(e); err != nil { return err }
}
if elm.Error != nil {
if err = elm.Error.Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Iq) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "from":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.From = j
case x.Name.Space == "" && x.Name.Local == "to":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.To = j
case x.Name.Space == "" && x.Name.Local == "id":
elm.Id = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "type":
value := IqType(x.Value)
elm.Type = &value
case x.Name.Space == "http://www.w3.org/XML/1998/namespace" && x.Name.Local == "lang":
elm.Lang = & x.Value
}
}
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
default:
if newel, ok := xmlencoder.GetExtension(t.Name); ok {
if err = newel.(xmlencoder.Extension).Decode(d, &t); err != nil { return err }
elm.Payload = newel
} else {
if err = d.Skip(); err != nil { return err }
}
case t.Name.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" && t.Name.Local == "error":newel := &stanza.Error{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Error = newel
}
}
}
return err
}

func (elm *Presence) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "presence"); err != nil { return err }
if elm.From != nil {
if err = e.Attribute("", "from", elm.From.String()); err != nil { return err }
}
if elm.To != nil {
if err = e.Attribute("", "to", elm.To.String()); err != nil { return err }
}
if elm.Id != nil {
if err = e.Attribute("", "id", *elm.Id); err != nil { return err }
}
if elm.Type != nil {
if err = e.Attribute("", "type", string(*elm.Type)); err != nil { return err }
}
if elm.Lang != nil {
if err = e.Attribute("http://www.w3.org/XML/1998/namespace", "lang", string(*elm.Lang)); err != nil { return err }
}
if elm.Show != nil {
if err = e.SimpleElement(NS, "show", string(*elm.Show)); err != nil { return err }
}
if elm.Status != nil {
if err = e.SimpleElement(NS, "status", *elm.Status); err != nil { return err }
}
if elm.Priority != nil {
if err = e.SimpleElement(NS, "priority", strconv.FormatInt(int64(*elm.Priority), 10)); err != nil { return err }
}
for _, x := range elm.X {
if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err} 
}
if elm.Error != nil {
if err = elm.Error.Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Presence) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "from":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.From = j
case x.Name.Space == "" && x.Name.Local == "to":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.To = j
case x.Name.Space == "" && x.Name.Local == "id":
elm.Id = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "type":
value := PresenceType(x.Value)
elm.Type = &value
case x.Name.Space == "http://www.w3.org/XML/1998/namespace" && x.Name.Local == "lang":
elm.Lang = & x.Value
}
}
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
case t.Name.Space == NS && t.Name.Local == "show":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Show = PresenceShow(s)
case t.Name.Space == NS && t.Name.Local == "status":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Status = s
case t.Name.Space == NS && t.Name.Local == "priority":
var s string
if s, err = d.Text(); err != nil { return err }
var i int64
if i, err = strconv.ParseInt(s, 10, 0); err == nil {
*elm.Priority = int(i)
}
default:
if newel, ok := xmlencoder.GetExtension(t.Name); ok {
if err = newel.(xmlencoder.Extension).Decode(d, &t); err != nil { return err }
elm.X = append(elm.X, newel)
} else {
if err = d.Skip(); err != nil { return err }
}
case t.Name.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" && t.Name.Local == "error":newel := &stanza.Error{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Error = newel
}
}
}
return err
}

func (elm *Message) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "message"); err != nil { return err }
if elm.From != nil {
if err = e.Attribute("", "from", elm.From.String()); err != nil { return err }
}
if elm.To != nil {
if err = e.Attribute("", "to", elm.To.String()); err != nil { return err }
}
if elm.Id != nil {
if err = e.Attribute("", "id", *elm.Id); err != nil { return err }
}
if elm.Type != nil {
if err = e.Attribute("", "type", string(*elm.Type)); err != nil { return err }
}
if elm.Lang != nil {
if err = e.Attribute("http://www.w3.org/XML/1998/namespace", "lang", string(*elm.Lang)); err != nil { return err }
}
if elm.Thread != nil {
if err = e.SimpleElement(NS, "thread", *elm.Thread); err != nil { return err }
}
if elm.Subject != nil {
elm.Subject.Encode(e, NS, "subject")
}
if elm.Body != nil {
elm.Body.Encode(e, NS, "body")
}
for _, x := range elm.X {
if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err} 
}
if elm.Error != nil {
if err = elm.Error.Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Message) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "from":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.From = j
case x.Name.Space == "" && x.Name.Local == "to":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.To = j
case x.Name.Space == "" && x.Name.Local == "id":
elm.Id = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "type":
value := MessageType(x.Value)
elm.Type = &value
case x.Name.Space == "http://www.w3.org/XML/1998/namespace" && x.Name.Local == "lang":
elm.Lang = & x.Value
}
}
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
case t.Name.Space == NS && t.Name.Local == "thread":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Thread = s
case t.Name.Space == NS && t.Name.Local == "subject":
if err = elm.Subject.Decode(d, &t); err != nil { return err }
case t.Name.Space == NS && t.Name.Local == "body":
if err = elm.Body.Decode(d, &t); err != nil { return err }
default:
if newel, ok := xmlencoder.GetExtension(t.Name); ok {
if err = newel.(xmlencoder.Extension).Decode(d, &t); err != nil { return err }
elm.X = append(elm.X, newel)
} else {
if err = d.Skip(); err != nil { return err }
}
case t.Name.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" && t.Name.Local == "error":newel := &stanza.Error{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Error = newel
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "iq"}, Iq{}, true, true)
 xmlencoder.AddExtension(xml.Name{NS, "presence"}, Presence{}, true, true)
 xmlencoder.AddExtension(xml.Name{NS, "message"}, Message{}, true, true)
}
