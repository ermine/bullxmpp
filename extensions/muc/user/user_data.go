package user

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
import "strconv"
import "time"
import "github.com/ermine/bullxmpp/jid"
import "github.com/ermine/bullxmpp/xdata"
const NS = "http://jabber.org/protocol/muc#user"
type Action struct {
  Decline *Decline
  Destroy *Destroy
  Invite []*Invite
  Item *Item
  Password *string
  Status []*Status
}
type Status struct {
  Code *int
}
type Decline struct {
  From *jid.JID
  To *jid.JID
  Reason *string
}
type Destroy struct {
  Jid *jid.JID
  Reason *string
}
type Invite struct {
  From *jid.JID
  To *jid.JID
  Reason *string
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
  Continue bool
}
type ItemAffiliation string
const (
userItemAffiliationAdmin ItemAffiliation = "admin"
userItemAffiliationMember ItemAffiliation = "member"
userItemAffiliationNone ItemAffiliation = "none"
userItemAffiliationOutcast ItemAffiliation = "outcast"
userItemAffiliationOwner ItemAffiliation = "owner"
)
type ItemRole string
const (
userItemRoleModerator ItemRole = "moderator"
userItemRoleNone ItemRole = "none"
userItemRoleParticipant ItemRole = "participant"
userItemRoleVisitor ItemRole = "visitor"
)
func (elm *Action) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "x"); err != nil { return err }
if elm.Decline != nil {
if err = elm.Decline.Encode(e); err != nil { return err }
}
if elm.Destroy != nil {
if err = elm.Destroy.Encode(e); err != nil { return err }
}
for _, x := range elm.Invite {
if err = x.Encode(e); err != nil { return err} 
}
if elm.Item != nil {
if err = elm.Item.Encode(e); err != nil { return err }
}
if elm.Password != nil {
if err = e.SimpleElement(NS, "password", *elm.Password); err != nil { return err }
}
for _, x := range elm.Status {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Action) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "invite":
newel := &Invite{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Invite = append(elm.Invite, newel)
case t.Name.Space == NS && t.Name.Local == "password":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Password = s
case t.Name.Space == NS && t.Name.Local == "status":
newel := &Status{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Status = append(elm.Status, newel)
}
}
}
return err
}

func (elm *Status) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "status"); err != nil { return err }
if elm.Code != nil {
if err = e.Attribute("", "code", strconv.FormatInt(int64(*elm.Code), 10)); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Status) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "code":
var i int64
i, err = strconv.ParseInt(x.Value, 10, 0)
if err == nil {
*elm.Code = int(i)
}
}
}
return err
}

func (elm *Decline) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "decline"); err != nil { return err }
if elm.From != nil {
if err = e.Attribute("", "from", elm.From.String()); err != nil { return err }
}
if elm.To != nil {
if err = e.Attribute("", "to", elm.To.String()); err != nil { return err }
}
if elm.Reason != nil {
if err = e.SimpleElement(NS, "reason", *elm.Reason); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Decline) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "reason":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Reason = s
}
}
}
return err
}

func (elm *Destroy) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "destroy"); err != nil { return err }
if elm.Jid != nil {
if err = e.Attribute("", "jid", elm.Jid.String()); err != nil { return err }
}
if elm.Reason != nil {
if err = e.SimpleElement(NS, "reason", *elm.Reason); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Destroy) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "jid":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.Jid = j
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
case t.Name.Space == NS && t.Name.Local == "reason":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Reason = s
}
}
}
return err
}

func (elm *Invite) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "invite"); err != nil { return err }
if elm.From != nil {
if err = e.Attribute("", "from", elm.From.String()); err != nil { return err }
}
if elm.To != nil {
if err = e.Attribute("", "to", elm.To.String()); err != nil { return err }
}
if elm.Reason != nil {
if err = e.SimpleElement(NS, "reason", *elm.Reason); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Invite) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "reason":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Reason = s
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
if elm.Continue {
if err = e.StartElement(NS, "continue"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
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
case t.Name.Space == NS && t.Name.Local == "continue":
elm.Continue = true
if err = d.Skip(); err != nil { return err }
continue
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "x"}, Action{}, true, true)
}
