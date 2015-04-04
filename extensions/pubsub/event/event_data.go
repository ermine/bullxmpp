package event

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
import "github.com/ermine/bullxmpp/xdata"
import "time"
import "github.com/ermine/bullxmpp/jid"
const NS = "http://jabber.org/protocol/pubsub#event"
type Event struct {
 Payload interface{}
}
type Collection struct {
  Node *string
  Type struct {
  Node *string
  Type *CollectionTypeType
}

}
type Configuration struct {
  Node *string
  Xdata *xdata.X

}
type Delete struct {
  Node *string
  Redirect struct {
  Url *string
}

}
type Items struct {
  Node *string
  Items []*Item
  Retracts []*Retract
}
type Item struct {
  Id *string
  Node *string
  Publisher *string
  Event interface{}

}
type Purge struct {
  Node *string
}
type Retract struct {
  Id *string
}
type Subscription struct {
  Expiry *time.Time
  Jid *jid.JID
  Node *string
  Subid *string
  Subscription *SubscriptionSubscription
}
type CollectionTypeType string
const (
eventCollectionTypeTypeAssociate CollectionTypeType = "associate"
eventCollectionTypeTypeDisassociate CollectionTypeType = "disassociate"
)
type SubscriptionSubscription string
const (
eventSubscriptionSubscriptionNone SubscriptionSubscription = "none"
eventSubscriptionSubscriptionPending SubscriptionSubscription = "pending"
eventSubscriptionSubscriptionSubscribed SubscriptionSubscription = "subscribed"
eventSubscriptionSubscriptionUnconfigured SubscriptionSubscription = "unconfigured"
)
func (elm *Event) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "event"); err != nil { return err }
if elm.Payload != nil {
if elm.Payload != nil {
if err = elm.Payload.(xmlencoder.Extension).Encode(e); err != nil { return err }
}
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Event) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "collection":
newel := &Collection{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
case t.Name.Space == NS && t.Name.Local == "configuration":
newel := &Configuration{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
case t.Name.Space == NS && t.Name.Local == "delete":
newel := &Delete{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
case t.Name.Space == NS && t.Name.Local == "items":
newel := &Items{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
case t.Name.Space == NS && t.Name.Local == "purge":
newel := &Purge{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
case t.Name.Space == NS && t.Name.Local == "subscription":
newel := &Subscription{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
}
}
}
return err
}

func (elm *Collection) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "collection"); err != nil { return err }
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
if err = e.StartElement(NS, string(*elm.Type.Type)); err != nil { return err }
if elm.Type.Node != nil {
if err = e.Attribute("", "node", *elm.Type.Node); err != nil { return err }
}
if elm.Type.Type != nil {
}
if err = e.EndElement(); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Collection) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
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
if t.Name.Space == NS {
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "node":
elm.Type.Node = xmlencoder.Copystring(x.Value)
}
}
}
}
}
}
return err
}

func (elm *Configuration) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "configuration"); err != nil { return err }
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
if elm.Xdata != nil {
if err = elm.Xdata.Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Configuration) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
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
case t.Name.Space == "jabber:x:data" && t.Name.Local == "x":newel := &xdata.X{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Xdata = newel
}
}
}
return err
}

func (elm *Delete) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "delete"); err != nil { return err }
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
if err = e.StartElement(NS, "redirect"); err != nil { return err }
if elm.Redirect.Url != nil {
if err = e.Attribute("", "url", *elm.Redirect.Url); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Delete) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
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
case t.Name.Space == NS && t.Name.Local == "redirect":
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "url":
elm.Redirect.Url = xmlencoder.Copystring(x.Value)
}
}
}
}
}
return err
}

func (elm *Items) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "items"); err != nil { return err }
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
for _, x := range elm.Items {
if err = x.Encode(e); err != nil { return err} 
}
for _, x := range elm.Retracts {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Items) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
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
case t.Name.Space == NS && t.Name.Local == "retract":
newel := &Retract{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Retracts = append(elm.Retracts, newel)
}
}
}
return err
}

func (elm *Item) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "item"); err != nil { return err }
if elm.Id != nil {
if err = e.Attribute("", "id", *elm.Id); err != nil { return err }
}
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
if elm.Publisher != nil {
if err = e.Attribute("", "publisher", *elm.Publisher); err != nil { return err }
}
if elm.Event != nil {
if err = elm.Event.(xmlencoder.Extension).Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Item) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "id":
elm.Id = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "publisher":
elm.Publisher = xmlencoder.Copystring(x.Value)
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
elm.Event = newel
} else {
if err = d.Skip(); err != nil { return err }
}
}
}
}
return err
}

func (elm *Purge) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "purge"); err != nil { return err }
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Purge) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
}
}
return err
}

func (elm *Retract) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "retract"); err != nil { return err }
if elm.Id != nil {
if err = e.Attribute("", "id", *elm.Id); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Retract) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "id":
elm.Id = xmlencoder.Copystring(x.Value)
}
}
return err
}

func (elm *Subscription) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "subscription"); err != nil { return err }
if elm.Expiry != nil {
if err = e.Attribute("", "expiry", elm.Expiry.String()); err != nil { return err }
}
if elm.Jid != nil {
if err = e.Attribute("", "jid", elm.Jid.String()); err != nil { return err }
}
if elm.Node != nil {
if err = e.Attribute("", "node", *elm.Node); err != nil { return err }
}
if elm.Subid != nil {
if err = e.Attribute("", "subid", *elm.Subid); err != nil { return err }
}
if elm.Subscription != nil {
if err = e.Attribute("", "subscription", string(*elm.Subscription)); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Subscription) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "expiry":
*elm.Expiry, err = time.Parse(time.RFC3339, x.Value)
if err != nil { return err }
case x.Name.Space == "" && x.Name.Local == "jid":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.Jid = j
case x.Name.Space == "" && x.Name.Local == "node":
elm.Node = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "subid":
elm.Subid = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "subscription":
value := SubscriptionSubscription(x.Value)
elm.Subscription = &value
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "event"}, Event{}, true, true)
}
