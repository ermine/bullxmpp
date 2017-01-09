package client

import (
	"github.com/ermine/bullxmpp/jid"
	"github.com/ermine/bullxmpp/stream"
	"github.com/ermine/bullxmpp/starttls"
	"github.com/ermine/bullxmpp/bind"
	"github.com/ermine/bullxmpp/sasl"
	"github.com/ermine/bullxmpp/session"
	"github.com/ermine/bullxmpp/xmlencoder"
	"net"
	"strconv"
	"errors"
	//	"math/rand"
	"encoding/xml"
)

const ns_client = "jabber:client"

var ErrorNotFound = errors.New("NotFound")

type Client struct {
	*stream.Stream
	sid int
}

func New(jid_s, password, lang string, secure bool, logging bool) (cli *Client, err error) {
	jid, err := jid.New(jid_s)
	if err != nil { return }
	return NewWithJID(jid, password, lang, secure, logging)
}

func NewWithJID(jid_ *jid.JID, password, lang string, secure bool, logging bool) (cli *Client, err error) {
	host, err := jid_.GetIDN()
	if err != nil {
		return
	}
	conn, err := Dial(host)
	if err != nil {
		return
	}

	xmlencoder.ReplaceExtensionStruct(xml.Name{starttls.NS, "starttls"}, TLSFeature{})
	xmlencoder.ReplaceExtensionStruct(xml.Name{sasl.NS, "mechanisms"}, SASLFeature{})
	xmlencoder.ReplaceExtensionStruct(xml.Name{bind.NS, "bind"}, BindFeature{})
	xmlencoder.ReplaceExtensionStruct(xml.Name{session.NS, "session"}, SessionFeature{})
	
	var options Options
	options.tlsUse = TLSDisabled
	if secure { options.tlsUse = TLSEnabled }
	options.password = password
	var knownFeatures []stream.InitiatorStreamFeature
	knownFeatures = append(knownFeatures, &TLSFeature{})
	knownFeatures = append(knownFeatures, &SASLFeature{})
	knownFeatures = append(knownFeatures, &BindFeature{})
	knownFeatures = append(knownFeatures, &SessionFeature{})

	strm := stream.New(conn, jid_, ns_client, lang, logging, options)
	cli = &Client{strm, 0}
	if err = cli.InitiatorSetup(knownFeatures); err != nil { return }
	return
}

type Options struct {
	tlsUse TLSOption
	password string
}

func (o Options) TLSUse() TLSOption {
	return o.tlsUse
}

func (o Options) GetUserPassword() string {
	return o.password
}

func Dial(host string) (net.Conn,  error) {
	_, srvs, _ := net.LookupSRV("xmpp-client", "tcp", host)
	var conn net.Conn
	var err error
	for _, x := range srvs {
		conn, err = net.Dial("tcp", x.Target + ":" + strconv.Itoa(int(x.Port)))
		if err == nil {
			return conn, nil
		}
	}
	return net.Dial("tcp", host + ":5222")
}

func (cli *Client) Disconnect() error {
	return cli.Disconnect()
}
