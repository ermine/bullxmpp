package iqregister

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
import "github.com/ermine/bullxmpp/xdata"
import "github.com/ermine/bullxmpp/xoob"
const NS = "jabber:iq:register"
type Query struct {
  Fields struct {
  Registered bool
  Instructions *string
  Username *string
  Nick *string
  Password *string
  Name *string
  First *string
  Last *string
  Email *string
  Address *string
  City *string
  State *string
  Zip *string
  Phone *string
  Url *string
  Date *string
  Misc *string
  Text *string
  Key *string
}

  Remove bool
  Xdata *xdata.X

  Xoob *xoob.X

}
func (elm *Query) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "query"); err != nil { return err }
if elm.Fields.Registered {
if err = e.StartElement(NS, "registered"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.Fields.Instructions != nil {
if err = e.SimpleElement(NS, "instructions", *elm.Fields.Instructions); err != nil { return err }
}
if elm.Fields.Username != nil {
if err = e.SimpleElement(NS, "username", *elm.Fields.Username); err != nil { return err }
}
if elm.Fields.Nick != nil {
if err = e.SimpleElement(NS, "nick", *elm.Fields.Nick); err != nil { return err }
}
if elm.Fields.Password != nil {
if err = e.SimpleElement(NS, "password", *elm.Fields.Password); err != nil { return err }
}
if elm.Fields.Name != nil {
if err = e.SimpleElement(NS, "name", *elm.Fields.Name); err != nil { return err }
}
if elm.Fields.First != nil {
if err = e.SimpleElement(NS, "first", *elm.Fields.First); err != nil { return err }
}
if elm.Fields.Last != nil {
if err = e.SimpleElement(NS, "last", *elm.Fields.Last); err != nil { return err }
}
if elm.Fields.Email != nil {
if err = e.SimpleElement(NS, "email", *elm.Fields.Email); err != nil { return err }
}
if elm.Fields.Address != nil {
if err = e.SimpleElement(NS, "address", *elm.Fields.Address); err != nil { return err }
}
if elm.Fields.City != nil {
if err = e.SimpleElement(NS, "city", *elm.Fields.City); err != nil { return err }
}
if elm.Fields.State != nil {
if err = e.SimpleElement(NS, "state", *elm.Fields.State); err != nil { return err }
}
if elm.Fields.Zip != nil {
if err = e.SimpleElement(NS, "zip", *elm.Fields.Zip); err != nil { return err }
}
if elm.Fields.Phone != nil {
if err = e.SimpleElement(NS, "phone", *elm.Fields.Phone); err != nil { return err }
}
if elm.Fields.Url != nil {
if err = e.SimpleElement(NS, "url", *elm.Fields.Url); err != nil { return err }
}
if elm.Fields.Date != nil {
if err = e.SimpleElement(NS, "date", *elm.Fields.Date); err != nil { return err }
}
if elm.Fields.Misc != nil {
if err = e.SimpleElement(NS, "misc", *elm.Fields.Misc); err != nil { return err }
}
if elm.Fields.Text != nil {
if err = e.SimpleElement(NS, "text", *elm.Fields.Text); err != nil { return err }
}
if elm.Fields.Key != nil {
if err = e.SimpleElement(NS, "key", *elm.Fields.Key); err != nil { return err }
}
if elm.Remove {
if err = e.StartElement(NS, "remove"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.Xdata != nil {
if err = elm.Xdata.Encode(e); err != nil { return err }
}
if elm.Xoob != nil {
if err = elm.Xoob.Encode(e); err != nil { return err }
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
case t.Name.Space == NS && t.Name.Local == "registered":
if err = d.Skip(); err != nil { return err }
case t.Name.Space == NS && t.Name.Local == "instructions":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Instructions = s
case t.Name.Space == NS && t.Name.Local == "username":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Username = s
case t.Name.Space == NS && t.Name.Local == "nick":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Nick = s
case t.Name.Space == NS && t.Name.Local == "password":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Password = s
case t.Name.Space == NS && t.Name.Local == "name":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Name = s
case t.Name.Space == NS && t.Name.Local == "first":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.First = s
case t.Name.Space == NS && t.Name.Local == "last":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Last = s
case t.Name.Space == NS && t.Name.Local == "email":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Email = s
case t.Name.Space == NS && t.Name.Local == "address":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Address = s
case t.Name.Space == NS && t.Name.Local == "city":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.City = s
case t.Name.Space == NS && t.Name.Local == "state":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.State = s
case t.Name.Space == NS && t.Name.Local == "zip":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Zip = s
case t.Name.Space == NS && t.Name.Local == "phone":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Phone = s
case t.Name.Space == NS && t.Name.Local == "url":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Url = s
case t.Name.Space == NS && t.Name.Local == "date":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Date = s
case t.Name.Space == NS && t.Name.Local == "misc":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Misc = s
case t.Name.Space == NS && t.Name.Local == "text":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Text = s
case t.Name.Space == NS && t.Name.Local == "key":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Fields.Key = s
case t.Name.Space == NS && t.Name.Local == "remove":
elm.Remove = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == "jabber:x:data" && t.Name.Local == "x":newel := &xdata.X{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Xdata = newel
case t.Name.Space == "jabber:x:oob" && t.Name.Local == "x":newel := &xoob.X{}
if err = newel.Decode(d, &t); err != nil { return err }
elm.Xoob = newel
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "query"}, Query{}, true, true)
}
