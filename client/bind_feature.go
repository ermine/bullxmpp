package client

import (
	"jabber.ru/xmpp/bind"
	"jabber.ru/xmpp/stream"
	"math/rand"
	"strconv"
)

type BindFeature struct {
	bind.Bind
}

func (t *BindFeature) IsMandatory(strm *stream.Stream) bool {
	return true
}

func (t *BindFeature) NeedRestart() bool {
	return false
}

func (t *BindFeature) Negotate(strm *stream.Stream) error {
	if strm.Jid.Resource != "" {
		t.Resource = &strm.Jid.Resource
	}
	id := strconv.Itoa(rand.Intn(1000))
	iq := Iq{To:strm.Jid.Server(), Id:&id, Payload:t}
	typ := IqTypeSet
	iq.Type = &typ
	if err := strm.Send(&iq); err != nil { return err }
	resp, err := strm.Read()
	if err != nil { return err }
	if resp == nil { return stream.ErrUnexpectedEOF }
	switch rsp := resp.(type) {
	case *Iq:
		if rsp.Type != nil && *rsp.Type == IqTypeResult && rsp.Payload != nil {
			if newt, ok := rsp.Payload.(*BindFeature); ok {
				if newt.Jid != nil {
					strm.Jid = newt.Jid
					return nil
				}
			} else {
				return stream.ErrMalformedProtocol
			}
		}
	case *stream.Error:
		return rsp
	}
	return stream.ErrMalformedProtocol
}




