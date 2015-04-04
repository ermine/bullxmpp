package stanza

import "encoding/xml"
import "github.com/ermine/bullxmpp/xmlencoder"
const NS = "urn:ietf:params:xml:ns:xmpp-stanzas"
type Error struct {
  Text *xmlencoder.LangString
  Condition struct {
  Name *ErrorConditionName
  Extra *string
}

}
type ErrorConditionName string
const (
ErrorConditionNameBadRequest ErrorConditionName = "bad-request"
ErrorConditionNameConflict ErrorConditionName = "conflict"
ErrorConditionNameFeatureNotImplemented ErrorConditionName = "feature-not-implemented"
ErrorConditionNameForbidden ErrorConditionName = "forbidden"
ErrorConditionNameGone ErrorConditionName = "gone"
ErrorConditionNameInternalServerError ErrorConditionName = "internal-server-error"
ErrorConditionNameItemNotFound ErrorConditionName = "item-not-found"
ErrorConditionNameJidMalformed ErrorConditionName = "jid-malformed"
ErrorConditionNameNotAcceptable ErrorConditionName = "not-acceptable"
ErrorConditionNameNotAllowed ErrorConditionName = "not-allowed"
ErrorConditionNameNotAuthorized ErrorConditionName = "not-authorized"
ErrorConditionNamePaymentRequired ErrorConditionName = "payment-required"
ErrorConditionNamePolicyViolation ErrorConditionName = "policy-violation"
ErrorConditionNameRecipientUnavailable ErrorConditionName = "recipient-unavailable"
ErrorConditionNameRedirect ErrorConditionName = "redirect"
ErrorConditionNameRegistrationRequired ErrorConditionName = "registration-required"
ErrorConditionNameRemoteServerNotFound ErrorConditionName = "remote-server-not-found"
ErrorConditionNameRemoteServerTimeout ErrorConditionName = "remote-server-timeout"
ErrorConditionNameResourceConstraint ErrorConditionName = "resource-constraint"
ErrorConditionNameServiceUnavailable ErrorConditionName = "service-unavailable"
ErrorConditionNameSubscriptionRequired ErrorConditionName = "subscription-required"
ErrorConditionNameUndefinedCondition ErrorConditionName = "undefined-condition"
ErrorConditionNameUnexpectedRequest ErrorConditionName = "unexpected-request"
)
func (elm *Error) Encode(e *xmlencoder.Encoder) error {
var err error
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
 xmlencoder.AddExtension(xml.Name{NS, "error"}, Error{}, true, true)
}
