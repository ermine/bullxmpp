package client

import (
	"github.com/ermine/bullxmpp/session"
	"github.com/ermine/bullxmpp/stream"
	"math/rand"
	"strconv"
)

type SessionFeature struct {
	session.Session
}

func (t *SessionFeature) IsMandatory(strm *stream.Stream) bool {
	return false
}

func (t *SessionFeature) NeedRestart() bool {
	return false
}

func (t *SessionFeature) Negotate(strm *stream.Stream) error {
	id := strconv.Itoa(rand.Intn(1000))
	typ := IqTypeSet
	iq := Iq{From:strm.Jid, Id:&id, Payload:t, Type: &typ}
	if err := strm.Send(&iq); err != nil { return err }
	resp, err := strm.Read()
	if err != nil { return err }
	if resp == nil { return stream.ErrUnexpectedEOF }
	switch resp := resp.(type) {
	case *Iq:
		if resp.Type != nil && *resp.Type == IqTypeResult {
			return nil
		}
	case *stream.Error:
		return resp
	}
	return stream.ErrMalformedProtocol
}
		
