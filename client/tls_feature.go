package client

import (
	"github.com/ermine/bullxmpp/starttls"
	"github.com/ermine/bullxmpp/stream"
	"crypto/tls"
	"fmt"
)

type TLSOptions interface {
	TLSUse() TLSOption
}

type TLSOption int

const (
	TLSEnabled TLSOption = 1 << iota
	TLSDisabled TLSOption = 1 << iota
	TLSRequired TLSOption = 1 << iota
)

var DefaultConfig tls.Config

type TLSFeature struct {
	starttls.Starttls
}

func (t *TLSFeature) IsMandatory(strm *stream.Stream) bool {
	mandatory := false
	if options, ok := strm.UData.(TLSOptions); ok {
		fmt.Println("options ", options)
		if options.TLSUse() == TLSRequired {
			mandatory = true
		}
	}
	return mandatory || t.Required
}

func (t *TLSFeature) NeedRestart() bool {
	return true
}

func (t *TLSFeature) Negotate(strm *stream.Stream) error {
	if err := strm.Send(t); err != nil { return err }
	data, err := strm.Read()
	if err != nil { return err }
	if data == nil { return stream.ErrUnexpectedEOF }
	switch resp := data.(type) {
	case *starttls.Proceed:
		return turnTLS(strm)
	case *starttls.Failure:
		return resp
	}
	return stream.ErrMalformedProtocol
}

func turnTLS(strm *stream.Stream) error {
	DefaultConfig.ServerName = strm.Jid.Domain
	tlsconn := tls.Client(strm.Conn, &DefaultConfig)
	if err := tlsconn.Handshake(); err != nil { return err }
	strm.Conn = tlsconn
	return nil
}












