package stream

import (
	"github.com/ermine/bullxmpp/jid"
	"github.com/ermine/bullxmpp/logger"
	"github.com/ermine/bullxmpp/xmlencoder"
	"errors"
	"net"
	"reflect"
)

const (
  ns_streams = "http://etherx.jabber.org/streams"
  ns_server = "jabber:server"
	ns_xmpp_stanzas = "urn:ietf:params:xml:ns:xmpp-stanzas"
)

var (
	ErrUnexpectedEOF = errors.New("Unexpected end of stream")
	ErrMalformedProtocol = errors.New("Malformed Protocol")
	ErrUnsupportedType = errors.New("unsupported type")
	ErrUnsupportedVersion = errors.New("Unsupported XMPP version")
)

type Stream struct {
	UData interface{}
	Jid *jid.JID
	Conn net.Conn
	lang string
	logging bool
	content_namespace string
	encoder *xmlencoder.Encoder
	decoder *xmlencoder.Decoder
	handlers map[reflect.Type]HandlerFunc
}

type HandlerFunc func(strm *Stream, data interface{}) error

func New(conn net.Conn, jid *jid.JID, content_namespace string, 
	lang string, logging bool, userdata interface{}) *Stream {
	strm := Stream {
		UData: userdata,
		Jid: jid, 
		Conn: conn,
		lang: lang,
		logging : logging,
		content_namespace: content_namespace,
	}		
	if strm.logging {
		wr := logger.New(conn, conn)
		strm.encoder = xmlencoder.NewEncoder(wr)
		strm.decoder = xmlencoder.NewDecoder(wr)
	} else {
		strm.encoder = xmlencoder.NewEncoder(conn)
		strm.decoder = xmlencoder.NewDecoder(conn)
	}
	return &strm
}

type InitiatorStreamFeature interface {
	IsMandatory(strm *Stream) bool
	NeedRestart() bool
	Negotate(strm *Stream) error
}

func (strm *Stream) InitiatorSetup(knownFeatures []InitiatorStreamFeature) error {
	var err error
	if err = strm.encoder.StartStream(); err != nil { return err }
	if err = strm.encoder.SetPrefix("", strm.content_namespace); err != nil { return err }
	version := "1.0"
	s := Start{
		Version:  &version,
		Lang: &strm.lang,
		To: strm.Jid.Server(),
	}
	if err = strm.Send(&s); err != nil { return err }
	
	var r interface{}
	if r, err = strm.Read(); err != nil { return err }
	if r == nil { return ErrUnexpectedEOF }
	switch r := r.(type) {
	case *Start:
		if r.Version == nil || *r.Version != version {
			return ErrUnsupportedVersion
		}
		if r.Lang != nil {
			strm.lang = *r.Lang
		}
	default:
		return ErrMalformedProtocol
	}

	if r, err = strm.decoder.Decode(); err != nil { return err }
	if r == nil { return ErrUnexpectedEOF }
	switch r := r.(type) {
	case *Features:
		return strm.handleInitiatorFeatures(r, knownFeatures)
	case *Error:
		return r
	}
	return ErrMalformedProtocol
}

func (strm *Stream) handleInitiatorFeatures(sf *Features, knownFeatures []InitiatorStreamFeature) error {
	var newknownFeatures []InitiatorStreamFeature
	for i, k := range knownFeatures {
		kType := reflect.TypeOf(k)
		for _, x := range *sf {
			if kType == reflect.TypeOf(x) {
				f := x.(InitiatorStreamFeature)
				if i < len(knownFeatures) - 1 {
					newknownFeatures = append(newknownFeatures, knownFeatures[i+1:]...)
				}
				if err := f.Negotate(strm); err != nil { return err }
				if f.NeedRestart() {
					if strm.logging {
						wr := logger.New(strm.Conn, strm.Conn)
						strm.decoder = xmlencoder.NewDecoder(wr)
						strm.encoder = xmlencoder.NewEncoder(wr)
					} else {
						strm.decoder = xmlencoder.NewDecoder(strm.Conn)
						strm.encoder = xmlencoder.NewEncoder(strm.Conn)
					}
					return strm.InitiatorSetup(newknownFeatures)
				} else {
					return strm.handleInitiatorFeatures(sf, newknownFeatures)
				}
			}
		}
		if k.IsMandatory(strm) {
			return errors.New("Unable to negotate mandatory feature")
		}
		newknownFeatures = append(newknownFeatures, k)
	}
	return nil
}

func (strm *Stream) CloseStream() error {
	return strm.encoder.EndStream()
}

func (strm *Stream) Disconnect() error {
	return strm.Conn.Close()
}

func (strm *Stream) Send(data xmlencoder.Extension) error {
	var err error
	if err =  data.Encode(strm.encoder); err != nil { return err }
	if err = strm.encoder.Flush(); err != nil { return err }
	return nil
}

func (strm *Stream) Read() (interface{}, error) {
	return strm.decoder.Decode()
	// TODO: handle stream error?
}

func (se Error) Error() string {
	return se.Text.Get("en")
}
