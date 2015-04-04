package vcard

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
const NS = "vcard-temp"
type Vcard struct {
  VERSION *string
  FN *string
  N *N
  Fields []interface{}
}
type N struct {
  FAMILY *string
  GIVEN *string
  NIDDLE *string
  PREFIX *string
  SUFFIX *string
}
type NICKNAME string
type PHOTO struct {
  TYPE *string
  BINVAL *string
  EXTVAL *string
}
type BDAY string
type ADR struct {
  HOME bool
  WORK bool
  POSTAL bool
  PARCEL bool
  DOMINTL *DOMINTL
  PREF bool
  POBOX *string
  EXTADD *string
  STREET *string
  LOCALITY *string
  REGION *string
  PCODE *string
  CTRY *string
}
type LABEL struct {
  HOME bool
  WORK bool
  POSTAL bool
  PARCEL bool
  DOMINTL *DOMINTL
  PREF bool
  LINE []string
}
type TEL struct {
  HOME bool
  WORK bool
  VOICE bool
  FAX bool
  PAGER bool
  MSG bool
  CELL bool
  VIDEO bool
  BBS bool
  MODEM bool
  ISDN bool
  PCS bool
  PREF bool
  NUMBER *string
}
type EMAIL struct {
  HOME bool
  WORK bool
  INTERNET bool
  PREF bool
  X400 bool
  USERID *string
}
type JABBERID string
type MAILER string
type TZ string
type GEO struct {
  LAT *string
  LON *string
}
type TITLE string
type ROLE string
type LOGO struct {
  TYPE *string
  BINVAL *string
  EXTVAL *string
}
type EXTVAL string
type AGENT struct {
 Payload interface{}
}
type ORG struct {
  ORGNAME *string
  ORGUNIT []string
}
type CATEGORIES struct {
  KEYWORD []string
}
type NOTE string
type PRODID string
type REV string
type SORTSTRING string
type SOUND struct {
  Type *SOUNDType
  Value *string
}
type PHONETIC string
type UID string
type URL string
type DESC string
type CLASS struct {
  Type *CLASSType
}
type KEY struct {
  TYPE *string
  CRED *string
}
type DOMINTL string
const (
DOMINTLDOM DOMINTL = "DOM"
DOMINTLINTL DOMINTL = "INTL"
)
type SOUNDType string
const (
SOUNDTypePHONETIC SOUNDType = "PHONETIC"
SOUNDTypeBINVAL SOUNDType = "BINVAL"
SOUNDTypeEXTVAL SOUNDType = "EXTVAL"
)
type CLASSType string
const (
CLASSTypePUBLIC CLASSType = "PUBLIC"
CLASSTypePRIVATE CLASSType = "PRIVATE"
CLASSTypeCONFIDENTIAL CLASSType = "CONFIDENTIAL"
)
func (elm *Vcard) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "vCard"); err != nil { return err }
if elm.VERSION != nil {
if err = e.SimpleElement(NS, "VERSION", *elm.VERSION); err != nil { return err }
}
if elm.FN != nil {
if err = e.SimpleElement(NS, "FN", *elm.FN); err != nil { return err }
}
if elm.N != nil {
if err = elm.N.Encode(e); err != nil { return err }
}
for _, x := range elm.Fields {
if err = x.(xmlencoder.Extension).Encode(e); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Vcard) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "VERSION":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.VERSION = s
case t.Name.Space == NS && t.Name.Local == "FN":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.FN = s
}
}
}
return err
}

func (elm *N) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "N"); err != nil { return err }
if elm.FAMILY != nil {
if err = e.SimpleElement(NS, "FAMILY", *elm.FAMILY); err != nil { return err }
}
if elm.GIVEN != nil {
if err = e.SimpleElement(NS, "GIVEN", *elm.GIVEN); err != nil { return err }
}
if elm.NIDDLE != nil {
if err = e.SimpleElement(NS, "NIDDLE", *elm.NIDDLE); err != nil { return err }
}
if elm.PREFIX != nil {
if err = e.SimpleElement(NS, "PREFIX", *elm.PREFIX); err != nil { return err }
}
if elm.SUFFIX != nil {
if err = e.SimpleElement(NS, "SUFFIX", *elm.SUFFIX); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *N) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "FAMILY":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.FAMILY = s
case t.Name.Space == NS && t.Name.Local == "GIVEN":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.GIVEN = s
case t.Name.Space == NS && t.Name.Local == "NIDDLE":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.NIDDLE = s
case t.Name.Space == NS && t.Name.Local == "PREFIX":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.PREFIX = s
case t.Name.Space == NS && t.Name.Local == "SUFFIX":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.SUFFIX = s
}
}
}
return err
}

func (elm *NICKNAME) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "NICKNAME", string(*elm)); err != nil { return err }
return nil
}

func (elm *NICKNAME) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = NICKNAME(s)
return err
}

func (elm *PHOTO) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "PHOTO"); err != nil { return err }
if elm.TYPE != nil {
if err = e.SimpleElement(NS, "TYPE", *elm.TYPE); err != nil { return err }
}
if elm.BINVAL != nil {
if err = e.SimpleElement(NS, "BINVAL", *elm.BINVAL); err != nil { return err }
}
if elm.EXTVAL != nil {
if err = e.SimpleElement(NS, "EXTVAL", *elm.EXTVAL); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *PHOTO) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "TYPE":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.TYPE = s
case t.Name.Space == NS && t.Name.Local == "BINVAL":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.BINVAL = s
case t.Name.Space == NS && t.Name.Local == "EXTVAL":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.EXTVAL = s
}
}
}
return err
}

func (elm *BDAY) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "BDAY", string(*elm)); err != nil { return err }
return nil
}

func (elm *BDAY) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = BDAY(s)
return err
}

func (elm *ADR) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "ADR"); err != nil { return err }
if elm.HOME {
if err = e.StartElement(NS, "HOME"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.WORK {
if err = e.StartElement(NS, "WORK"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.POSTAL {
if err = e.StartElement(NS, "POSTAL"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PARCEL {
if err = e.StartElement(NS, "PARCEL"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.DOMINTL != nil {
if err = e.StartElement(NS, string(*elm.DOMINTL)); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PREF {
if err = e.StartElement(NS, "PREF"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.POBOX != nil {
if err = e.SimpleElement(NS, "POBOX", *elm.POBOX); err != nil { return err }
}
if elm.EXTADD != nil {
if err = e.SimpleElement(NS, "EXTADD", *elm.EXTADD); err != nil { return err }
}
if elm.STREET != nil {
if err = e.SimpleElement(NS, "STREET", *elm.STREET); err != nil { return err }
}
if elm.LOCALITY != nil {
if err = e.SimpleElement(NS, "LOCALITY", *elm.LOCALITY); err != nil { return err }
}
if elm.REGION != nil {
if err = e.SimpleElement(NS, "REGION", *elm.REGION); err != nil { return err }
}
if elm.PCODE != nil {
if err = e.SimpleElement(NS, "PCODE", *elm.PCODE); err != nil { return err }
}
if elm.CTRY != nil {
if err = e.SimpleElement(NS, "CTRY", *elm.CTRY); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *ADR) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "HOME":
elm.HOME = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "WORK":
elm.WORK = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "POSTAL":
elm.POSTAL = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PARCEL":
elm.PARCEL = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PREF":
elm.PREF = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "POBOX":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.POBOX = s
case t.Name.Space == NS && t.Name.Local == "EXTADD":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.EXTADD = s
case t.Name.Space == NS && t.Name.Local == "STREET":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.STREET = s
case t.Name.Space == NS && t.Name.Local == "LOCALITY":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.LOCALITY = s
case t.Name.Space == NS && t.Name.Local == "REGION":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.REGION = s
case t.Name.Space == NS && t.Name.Local == "PCODE":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.PCODE = s
case t.Name.Space == NS && t.Name.Local == "CTRY":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.CTRY = s
default:
if t.Name.Space == NS {
*elm.DOMINTL = DOMINTL(t.Name.Local)
if err = d.Skip(); err != nil { return err }
}
}
}
}
return err
}

func (elm *LABEL) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "LABEL"); err != nil { return err }
if elm.HOME {
if err = e.StartElement(NS, "HOME"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.WORK {
if err = e.StartElement(NS, "WORK"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.POSTAL {
if err = e.StartElement(NS, "POSTAL"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PARCEL {
if err = e.StartElement(NS, "PARCEL"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.DOMINTL != nil {
if err = e.StartElement(NS, string(*elm.DOMINTL)); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PREF {
if err = e.StartElement(NS, "PREF"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
for _, x := range elm.LINE {
if err = e.SimpleElement(NS, "LINE", x); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *LABEL) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "HOME":
elm.HOME = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "WORK":
elm.WORK = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "POSTAL":
elm.POSTAL = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PARCEL":
elm.PARCEL = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PREF":
elm.PREF = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "LINE":
var s string
if s, err = d.Text(); err != nil { return err }
elm.LINE = append(elm.LINE, s)
default:
if t.Name.Space == NS {
*elm.DOMINTL = DOMINTL(t.Name.Local)
if err = d.Skip(); err != nil { return err }
}
}
}
}
return err
}

func (elm *TEL) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "TEL"); err != nil { return err }
if elm.HOME {
if err = e.StartElement(NS, "HOME"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.WORK {
if err = e.StartElement(NS, "WORK"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.VOICE {
if err = e.StartElement(NS, "VOICE"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.FAX {
if err = e.StartElement(NS, "FAX"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PAGER {
if err = e.StartElement(NS, "PAGER"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.MSG {
if err = e.StartElement(NS, "MSG"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.CELL {
if err = e.StartElement(NS, "CELL"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.VIDEO {
if err = e.StartElement(NS, "VIDEO"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.BBS {
if err = e.StartElement(NS, "BBS"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.MODEM {
if err = e.StartElement(NS, "MODEM"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.ISDN {
if err = e.StartElement(NS, "ISDN"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PCS {
if err = e.StartElement(NS, "PCS"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PREF {
if err = e.StartElement(NS, "PREF"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.NUMBER != nil {
if err = e.SimpleElement(NS, "NUMBER", *elm.NUMBER); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *TEL) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "HOME":
elm.HOME = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "WORK":
elm.WORK = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "VOICE":
elm.VOICE = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "FAX":
elm.FAX = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PAGER":
elm.PAGER = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "MSG":
elm.MSG = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "CELL":
elm.CELL = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "VIDEO":
elm.VIDEO = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "BBS":
elm.BBS = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "MODEM":
elm.MODEM = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "ISDN":
elm.ISDN = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PCS":
elm.PCS = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PREF":
elm.PREF = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "NUMBER":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.NUMBER = s
}
}
}
return err
}

func (elm *EMAIL) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "EMAIL"); err != nil { return err }
if elm.HOME {
if err = e.StartElement(NS, "HOME"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.WORK {
if err = e.StartElement(NS, "WORK"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.INTERNET {
if err = e.StartElement(NS, "INTERNET"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.PREF {
if err = e.StartElement(NS, "PREF"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.X400 {
if err = e.StartElement(NS, "X400"); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.USERID != nil {
if err = e.SimpleElement(NS, "USERID", *elm.USERID); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *EMAIL) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "HOME":
elm.HOME = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "WORK":
elm.WORK = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "INTERNET":
elm.INTERNET = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "PREF":
elm.PREF = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "X400":
elm.X400 = true
if err = d.Skip(); err != nil { return err }
continue
case t.Name.Space == NS && t.Name.Local == "USERID":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.USERID = s
}
}
}
return err
}

func (elm *JABBERID) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "JABBERID", string(*elm)); err != nil { return err }
return nil
}

func (elm *JABBERID) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = JABBERID(s)
return err
}

func (elm *MAILER) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "MAILER", string(*elm)); err != nil { return err }
return nil
}

func (elm *MAILER) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = MAILER(s)
return err
}

func (elm *TZ) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "TZ", string(*elm)); err != nil { return err }
return nil
}

func (elm *TZ) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = TZ(s)
return err
}

func (elm *GEO) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "GEO"); err != nil { return err }
if elm.LAT != nil {
if err = e.SimpleElement(NS, "LAT", *elm.LAT); err != nil { return err }
}
if elm.LON != nil {
if err = e.SimpleElement(NS, "LON", *elm.LON); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *GEO) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "LAT":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.LAT = s
case t.Name.Space == NS && t.Name.Local == "LON":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.LON = s
}
}
}
return err
}

func (elm *TITLE) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "TITLE", string(*elm)); err != nil { return err }
return nil
}

func (elm *TITLE) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = TITLE(s)
return err
}

func (elm *ROLE) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "ROLE", string(*elm)); err != nil { return err }
return nil
}

func (elm *ROLE) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = ROLE(s)
return err
}

func (elm *LOGO) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "LOGO"); err != nil { return err }
if elm.TYPE != nil {
if err = e.SimpleElement(NS, "TYPE", *elm.TYPE); err != nil { return err }
}
if elm.BINVAL != nil {
if err = e.SimpleElement(NS, "BINVAL", *elm.BINVAL); err != nil { return err }
}
if elm.EXTVAL != nil {
if err = e.SimpleElement(NS, "EXTVAL", *elm.EXTVAL); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *LOGO) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "TYPE":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.TYPE = s
case t.Name.Space == NS && t.Name.Local == "BINVAL":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.BINVAL = s
case t.Name.Space == NS && t.Name.Local == "EXTVAL":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.EXTVAL = s
}
}
}
return err
}

func (elm *EXTVAL) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "EXTVAL", string(*elm)); err != nil { return err }
return nil
}

func (elm *EXTVAL) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = EXTVAL(s)
return err
}

func (elm *AGENT) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "AGENT"); err != nil { return err }
if elm.Payload != nil {
if elm.Payload != nil {
if err = elm.Payload.(xmlencoder.Extension).Encode(e); err != nil { return err }
}
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *AGENT) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "vCard":
newel := &Vcard{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
case t.Name.Space == NS && t.Name.Local == "EXTVAL":
newel := &EXTVAL{}
if err = newel.Decode(d, &t); err != nil { return err}
elm.Payload = newel
}
}
}
return err
}

func (elm *ORG) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "ORG"); err != nil { return err }
if elm.ORGNAME != nil {
if err = e.SimpleElement(NS, "ORGNAME", *elm.ORGNAME); err != nil { return err }
}
for _, x := range elm.ORGUNIT {
if err = e.SimpleElement(NS, "ORGUNIT", x); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *ORG) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "ORGNAME":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.ORGNAME = s
case t.Name.Space == NS && t.Name.Local == "ORGUNIT":
var s string
if s, err = d.Text(); err != nil { return err }
elm.ORGUNIT = append(elm.ORGUNIT, s)
}
}
}
return err
}

func (elm *CATEGORIES) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "CATEGORIES"); err != nil { return err }
for _, x := range elm.KEYWORD {
if err = e.SimpleElement(NS, "KEYWORD", x); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *CATEGORIES) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "KEYWORD":
var s string
if s, err = d.Text(); err != nil { return err }
elm.KEYWORD = append(elm.KEYWORD, s)
}
}
}
return err
}

func (elm *NOTE) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "NOTE", string(*elm)); err != nil { return err }
return nil
}

func (elm *NOTE) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = NOTE(s)
return err
}

func (elm *PRODID) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "PRODID", string(*elm)); err != nil { return err }
return nil
}

func (elm *PRODID) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = PRODID(s)
return err
}

func (elm *REV) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "REV", string(*elm)); err != nil { return err }
return nil
}

func (elm *REV) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = REV(s)
return err
}

func (elm *SORTSTRING) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "SORT-STRING", string(*elm)); err != nil { return err }
return nil
}

func (elm *SORTSTRING) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = SORTSTRING(s)
return err
}

func (elm *SOUND) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, string(*elm.Type)); err != nil { return err }
if elm.Type != nil {
}
if elm.Value != nil {
if err = e.Text(*elm.Value); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *SOUND) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Value = s
return err
}

func (elm *PHONETIC) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "PHONETIC", string(*elm)); err != nil { return err }
return nil
}

func (elm *PHONETIC) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = PHONETIC(s)
return err
}

func (elm *UID) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "UID", string(*elm)); err != nil { return err }
return nil
}

func (elm *UID) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = UID(s)
return err
}

func (elm *URL) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "URL", string(*elm)); err != nil { return err }
return nil
}

func (elm *URL) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = URL(s)
return err
}

func (elm *DESC) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.SimpleElement(NS, "DESC", string(*elm)); err != nil { return err }
return nil
}

func (elm *DESC) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var s string
if s, err = d.Text(); err != nil { return err }
*elm = DESC(s)
return err
}

func (elm *CLASS) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "CLASS"); err != nil { return err }
if elm.Type != nil {
if err = e.StartElement(NS, string(*elm.Type)); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *CLASS) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
default:
if t.Name.Space == NS {
*elm.Type = CLASSType(t.Name.Local)
if err = d.Skip(); err != nil { return err }
}
}
}
}
return err
}

func (elm *KEY) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "KEY"); err != nil { return err }
if elm.TYPE != nil {
if err = e.SimpleElement(NS, "TYPE", *elm.TYPE); err != nil { return err }
}
if elm.CRED != nil {
if err = e.SimpleElement(NS, "CRED", *elm.CRED); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *KEY) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
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
case t.Name.Space == NS && t.Name.Local == "TYPE":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.TYPE = s
case t.Name.Space == NS && t.Name.Local == "CRED":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.CRED = s
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "vCard"}, Vcard{}, true, true)
}
