package admin

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
import "strconv"
import "time"
import "jabber.ru/xmpp/jid"
import "jabber.ru/xmpp/xdata"
const NS = "http://jabber.org/protocol/muc#admin"
type Query struct {
  Items []*Item
}
type Item struct {
  Affiliation *ItemAffiliation
  Jid *jid.JID
  Nick *string
  Role *ItemRole
  Actor struct {
  Jid *jid.JID
}

  Reason *string
}
type ItemAffiliation string
const (
adminItemAffiliationAdmin ItemAffiliation = "admin"
adminItemAffiliationMember ItemAffiliation = "member"
adminItemAffiliationNone ItemAffiliation = "none"
adminItemAffiliationOutcast ItemAffiliation = "outcast"
adminItemAffiliationOwner ItemAffiliation = "owner"
)
type ItemRole string
const (
adminItemRoleModerator ItemRole = "moderator"
adminItemRoleNone ItemRole = "none"
adminItemRoleParticipant ItemRole = "participant"
adminItemRoleVisitor ItemRole = "visitor"
)
func (elm *Query) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
for _, x := range elm.Items {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Query) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "item":
newel := &Item{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Items = append(elm.Items, newel)
}
}
}
return err
}

func (elm *Item) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "item"); err != nil { return err }
if elm.Affiliation != nil {
if err = e.Attribute("", "affiliation", string(*elm.Affiliation)); err != nil { return err }
}
if elm.Jid != nil {
if err = e.Attribute("", "jid", elm.Jid.String()); err != nil { return err }
}
if elm.Nick != nil {
if err = e.Attribute("", "nick", *elm.Nick); err != nil { return err }
}
if elm.Role != nil {
if err = e.Attribute("", "role", string(*elm.Role)); err != nil { return err }
}
if err = e.StartElement(NS, "actor"); err != nil { return err }
if elm.Actor.Jid != nil {
if err = e.Attribute("", "jid", elm.Actor.Jid.String()); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
if elm.Reason != nil {
if err = e.SimpleElement(NS, "reason", *elm.Reason); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Item) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "affiliation":
value := ItemAffiliation(x.Value)
elm.Affiliation = &value
case x.Name.Space == "" && x.Name.Local == "jid":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.Jid = j
case x.Name.Space == "" && x.Name.Local == "nick":
elm.Nick = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "role":
value := ItemRole(x.Value)
elm.Role = &value
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
case t.Name.Space == NS && t.Name.Local == "actor":
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "jid":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.Actor.Jid = j
}
}
case t.Name.Space == NS && t.Name.Local == "reason":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Reason = s
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Query{}, true, true)
}
