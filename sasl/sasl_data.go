package sasl

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "urn:ietf:params:xml:ns:xmpp-sasl"
type Mechanisms struct {
  Mechanism []string
}
type Auth struct {
  Mechanism *string
  Data []byte
}
type Success struct {
  Data []byte
}
type Challenge struct {
  Data []byte
}
type Response struct {
  Data []byte
}
type Failure struct {
  Text *xmlencoder.LangString
  Condition *FailureCondition
}
type FailureCondition string
const (
FailureConditionAborted FailureCondition = "aborted"
FailureConditionAccountDisabled FailureCondition = "account-disabled"
FailureConditionCredentialsExpired FailureCondition = "credentials-expired"
FailureConditionEncryptionRequired FailureCondition = "encryption-required"
FailureConditionIncorrectEncoding FailureCondition = "incorrect-encoding"
FailureConditionInvalidAuthzid FailureCondition = "invalid-authzid"
FailureConditionInvalidMechanism FailureCondition = "invalid-mechanism"
FailureConditionMalformedRequest FailureCondition = "malformed-request"
FailureConditionMechanismTooWeak FailureCondition = "mechanism-too-weak"
FailureConditionNotAuthorized FailureCondition = "not-authorized"
FailureConditionTemporaryAuthFailure FailureCondition = "temporary-auth-failure"
FailureConditionTransitionNeeded FailureCondition = "transition-needed"
)
func (elm *Mechanisms) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "mechanisms"); err != nil { return err }
for _, x := range elm.Mechanism {
if err = e.SimpleElement(NS, "mechanism", x); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Mechanisms) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "mechanism":
var s string
if s, err = d.Text(); err != nil { return err }
elm.Mechanism = append(elm.Mechanism, s)
}
}
}
return err
}

func (elm *Auth) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "auth"); err != nil { return err }
if elm.Mechanism != nil {
if err = e.Attribute("", "mechanism", *elm.Mechanism); err != nil { return err }
}
if elm.Data != nil {
if err = e.Bytes(elm.Data); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Auth) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
for _, x := range tag.Attr {
switch {
case x.Name.Space == "" && x.Name.Local == "mechanism":
elm.Mechanism = xmlencoder.Copystring(x.Value)
}
}
var bdata []byte
if bdata, err = d.Bytes(); err != nil { return err }
elm.Data = bdata
return err
}

func (elm *Success) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "success"); err != nil { return err }
if elm.Data != nil {
if err = e.Bytes(elm.Data); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Success) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var bdata []byte
if bdata, err = d.Bytes(); err != nil { return err }
elm.Data = bdata
return err
}

func (elm *Challenge) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "challenge"); err != nil { return err }
if elm.Data != nil {
if err = e.Bytes(elm.Data); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Challenge) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var bdata []byte
if bdata, err = d.Bytes(); err != nil { return err }
elm.Data = bdata
return err
}

func (elm *Response) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "response"); err != nil { return err }
if elm.Data != nil {
if err = e.Bytes(elm.Data); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Response) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var bdata []byte
if bdata, err = d.Bytes(); err != nil { return err }
elm.Data = bdata
return err
}

func (elm *Failure) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "failure"); err != nil { return err }
if elm.Text != nil {
elm.Text.Encode(e, NS, "text")
}
if elm.Condition != nil {
if err = e.StartElement(NS, string(*elm.Condition)); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Failure) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
*elm.Condition = FailureCondition(t.Name.Local)
if err = d.Skip(); err != nil { return err }
}
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "mechanisms"}, Mechanisms{}, true, false)
 xmlencoder.AddExtension(xml.Name{NS, "auth"}, Auth{}, false, true)
 xmlencoder.AddExtension(xml.Name{NS, "success"}, Success{}, true, false)
 xmlencoder.AddExtension(xml.Name{NS, "challenge"}, Challenge{}, true, false)
 xmlencoder.AddExtension(xml.Name{NS, "response"}, Response{}, false, true)
 xmlencoder.AddExtension(xml.Name{NS, "failure"}, Failure{}, true, false)
}
