package roster

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
import "strconv"
import "github.com/ermine/bullxmpp/jid"
const NS = "jabber:iq:roster"
type Roster struct {
  Ver *string
  Items []*Item
}
type Item struct {
  Approved bool
  Ask *ItemAsk
  Jid *jid.JID
  Name *string
  Subscription *ItemSubscription
  Group []string
}
type ItemAsk string
const (
ItemAskSubscribe ItemAsk = "subscribe"
)
type ItemSubscription string
const (
ItemSubscriptionBoth ItemSubscription = "both"
ItemSubscriptionFrom ItemSubscription = "from"
ItemSubscriptionNone ItemSubscription = "none"
ItemSubscriptionRemove ItemSubscription = "remove"
ItemSubscriptionTo ItemSubscription = "to"
)
func (elm *Roster) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if elm.Ver != nil {
if err = e.Attribute("", "ver", *elm.Ver); err != nil { return err }
}
for _, x := range elm.Items {
if err = x.Encode(e); err != nil { return err} 
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Roster) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "ver":
elm.Ver = xmlencoder.Copystring(x.Value)
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
if elm.Approved {
if err = e.Attribute("", "approved", strconv.FormatBool(elm.Approved)); err != nil { return err }
}
if elm.Ask != nil {
if err = e.Attribute("", "ask", string(*elm.Ask)); err != nil { return err }
}
if elm.Jid != nil {
if err = e.Attribute("", "jid", elm.Jid.String()); err != nil { return err }
}
if elm.Name != nil {
if err = e.Attribute("", "name", *elm.Name); err != nil { return err }
}
if elm.Subscription != nil {
if err = e.Attribute("", "subscription", string(*elm.Subscription)); err != nil { return err }
}
for _, x := range elm.Group {
if err = e.SimpleElement(NS, "group", x); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Item) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "approved":
var b bool
b, err = strconv.ParseBool(x.Value)
if err == nil {
elm.Approved = b
}
case x.Name.Space == "" && x.Name.Local == "ask":
value := ItemAsk(x.Value)
elm.Ask = &value
case x.Name.Space == "" && x.Name.Local == "jid":
var j *jid.JID
if j, err = jid.New(x.Value); err != nil { return err }
elm.Jid = j
case x.Name.Space == "" && x.Name.Local == "name":
elm.Name = xmlencoder.Copystring(x.Value)
case x.Name.Space == "" && x.Name.Local == "subscription":
value := ItemSubscription(x.Value)
elm.Subscription = &value
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
case t.Name.Space == NS && t.Name.Local == "group":
var s string
if s, err = d.Text(); err != nil { return err }
elm.Group = append(elm.Group, s)
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Roster{}, true, true)
}
