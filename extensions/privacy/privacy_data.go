package privacy

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
import "strconv"
const NS = "jabber:iq:privacy"
type Privacy struct {
  Active *Active
  Default *Default
  List []*List
}
type Active struct {
  Name *string
  Extra *string
}
type Default struct {
  Name *string
  Extra *string
}
type List struct {
  Name *string
  Items []*Item
}
type Item struct {
  Action *ItemAction
  Order *uint
  Type *ItemType
  Value *string
  Iq bool
  Message bool
  PresenceIn bool
  PresenceOut bool
}
type ItemAction string
const (
ItemActionAllow ItemAction = "allow"
ItemActionDeny ItemAction = "deny"
)
type ItemType string
const (
ItemTypeGroup ItemType = "group"
ItemTypeJid ItemType = "jid"
ItemTypeSubscription ItemType = "subscription"
)
func (elm *Privacy) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if elm.Active != nil {
if err = elm.Active.Encode(e); err != nil { return err }
}
if elm.Default != nil {
if err = elm.Default.Encode(e); err != nil { return err }
}
for _, x := range elm.List {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Privacy) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "list":
newel := &List{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.List = append(elm.List, newel)
}
}
}
return err
}

func (elm *Active) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "active"); err != nil { return err }
if elm.Name != nil {
if err = e.Attribute("", "name", *elm.Name); err != nil { return err }
}
if elm.Extra != nil {
if err = e.Text(*elm.Extra); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Active) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "name":
elm.Name = xmlencoder.Copystring(x.Value)
}
}
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Extra = s
return err
}

func (elm *Default) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "default"); err != nil { return err }
if elm.Name != nil {
if err = e.Attribute("", "name", *elm.Name); err != nil { return err }
}
if elm.Extra != nil {
if err = e.Text(*elm.Extra); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Default) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "name":
elm.Name = xmlencoder.Copystring(x.Value)
}
}
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Extra = s
return err
}

func (elm *List) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "list"); err != nil { return err }
if elm.Name != nil {
if err = e.Attribute("", "name", *elm.Name); err != nil { return err }
}
for _, x := range elm.Items {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *List) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "name":
elm.Name = xmlencoder.Copystring(x.Value)
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
if elm.Action != nil {
if err = e.Attribute("", "action", string(*elm.Action)); err != nil { return err }
}
if elm.Order != nil {
if err = e.Attribute("", "order", strconv.FormatUint(uint64(*elm.Order), 10)); err != nil { return err }
}
if elm.Type != nil {
if err = e.Attribute("", "type", string(*elm.Type)); err != nil { return err }
}
if elm.Value != nil {
if err = e.Attribute("", "value", *elm.Value); err != nil { return err }
}
if elm.Iq {
if err = e.StartElement(NS, "iq"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.Message {
if err = e.StartElement(NS, "message"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PresenceIn {
if err = e.StartElement(NS, "presence-in"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PresenceOut {
if err = e.StartElement(NS, "presence-out"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Item) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "action":
value := ItemAction(x.Value)
elm.Action = &value
case x.Name.Space == "" && x.Name.Local == "order":
var i uint64
i, err = strconv.ParseUint(x.Value, 10, 0)
if err == nil {
*elm.Order = uint(i)
}
case x.Name.Space == "" && x.Name.Local == "type":
value := ItemType(x.Value)
elm.Type = &value
case x.Name.Space == "" && x.Name.Local == "value":
elm.Value = xmlencoder.Copystring(x.Value)
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
case t.Name.Space == NS && t.Name.Local == "iq":
elm.Iq = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "message":
elm.Message = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "presence-in":
elm.PresenceIn = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "presence-out":
elm.PresenceOut = true
if err = d.Skip(); err != nil { return err }
continue
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Privacy{}, true, true)
}
