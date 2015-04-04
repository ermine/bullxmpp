package client

import (
	"jabber.ru/xmpp/jid"
	"jabber.ru/xmpp/sasl"
	"jabber.ru/xmpp/stream"
	"encoding/base64"
	"math/rand"
	"crypto/md5"
	"crypto/hmac"
	"crypto/sha1"
	"text/scanner"
	"strings"
	"strconv"
	"errors"
	"bytes"
	"fmt"
)

var (
	ErrNoKnownSASLMethod = errors.New("No known SASL method")
)

type SASLOptions interface {
	GetUserPassword() string
}

type SASLFeature struct {
	sasl.Mechanisms
}

type MechanismImpl struct {
	Name string
	Callback func(strm *stream.Stream, password string) error
}

var Mechanisms = []MechanismImpl {
//	MechanismImpl{"SCRAM-SHA-1", sasl_SCRAM_SHA_1},
//	MechanismImpl{"DIGEST-MD5", sasl_DIGEST_MD5},
	MechanismImpl{"PLAIN", sasl_PLAIN},
}


func (t *SASLFeature) IsMandatory(strm *stream.Stream) bool {
	return true
}

func (t *SASLFeature) NeedRestart() bool {
	return true
}

func (t *SASLFeature) Negotate(strm *stream.Stream) error {
	if options, ok := strm.UData.(SASLOptions); ok {
		for _, x := range Mechanisms {
			for _, z := range t.Mechanism {
				if x.Name == z {
					return x.Callback(strm, options.GetUserPassword())
				}
			}
		}
		return ErrNoKnownSASLMethod
	}
	return errors.New("password not found")
}

func sasl_PLAIN(strm *stream.Stream, password string) error {
	var err error
	data := "\x00" + strm.Jid.Node + "\x00" + password
	method := "PLAIN"
	if err = strm.Send(&sasl.Auth{&method, base64_encode(data)}); err != nil { return err }
	resp, err := strm.Read()
	if err != nil { return err }
	switch resp := resp.(type) {
	case *sasl.Success:
		return nil
	case *sasl.Failure:
		return resp
	case *stream.Error:
		return resp
	}
	return stream.ErrMalformedProtocol
}

func sasl_DIGEST_MD5(strm *stream.Stream, password string) error {
	var err error
	var resp interface {}
	method := "DIGEST-MD5"
	if err = strm.Send(&sasl.Auth{&method, []byte{}}); err != nil { return err }
	resp, err = strm.Read()
	if err != nil { return err }
	if resp == nil { return stream.ErrUnexpectedEOF }
	switch rsp := resp.(type) {
	case *sasl.Challenge:
		var t []byte
		if t, err = base64.StdEncoding.DecodeString(string(rsp.Data)); err != nil { return err }
		var tokens map[string]string
		if tokens, err = get_map(t); err != nil { return err }
		nonce := tokens["nonce"]
		cnonce := make_cnonce()
		nc := "00000001"
		digest_uri := "xmpp/" + strm.Jid.Domain
		realm := strm.Jid.Domain
		qop_list := strings.Split(tokens["qop"], ",")
		exists := false
		fmt.Println(qop_list)
		for _, x := range qop_list {
			if x == "auth" {
				exists = true
				break
			}
		}
		if !exists { return errors.New("No known qop method") }
		qop := "auth"
		a0 := md5_sum([]byte(strm.Jid.Node + ":" + realm + ":" + password))
		a1 := string(a0) + ":" + nonce + ":" + cnonce
		a2 := "AUTHENTICATE:" + digest_uri
		t = []byte(hex(md5_sum([]byte(a1))) + ":" + nonce + ":" + nc + ":" + cnonce + 
			":" + qop + ":" + hex(md5_sum([]byte(a2))))
		response := hex(md5_sum(t))
		response = `charset=utf-8,username="` + strm.Jid.Node + `",realm="` + realm + 
			`",nonce=` + nonce + ",cnonce=" + cnonce + ",nc=" + nc +
			`,qop="` + qop + `",digest-uri="` + digest_uri + `",response=` + response
		if err = strm.Send(&sasl.Response{base64_encode(response)}); err != nil { return err }
		resp, err = strm.Read()
		if err != nil { return err }
		if resp == nil { return stream.ErrUnexpectedEOF }
		switch rsp := resp.(type) {
		case *sasl.Challenge:
			if err = strm.Send(&sasl.Response{}); err != nil { return err }
			if resp, err = strm.Read(); err != nil { return err }
			if resp == nil { return stream.ErrUnexpectedEOF }
			switch rsp := resp.(type) {
			case *sasl.Success:
				return nil
			case *sasl.Failure:
				return rsp
			case *stream.Error:
				return rsp
			}
		case *sasl.Failure:
			return rsp
		case *stream.Error:
			return rsp
		}
	}
	return stream.ErrMalformedProtocol
}

func sasl_SCRAM_SHA_1(strm *stream.Stream, password string) error {
	var err error
	var resp interface{}
	cnonce := make_cnonce()
	clientFirstMessageBare := "n=" + strm.Jid.Node + ",r=" + cnonce
	data := "n,," + clientFirstMessageBare
	method := "SCRAM-SHA-1"
	if err = strm.Send(&sasl.Auth{&method, base64_encode(data)}); err != nil { return err }
	if resp, err = strm.Read(); err != nil { return err }
	if resp == nil { return stream.ErrUnexpectedEOF }
	switch rsp := resp.(type) {
	case *sasl.Challenge:
		var t []byte
		if t, err = base64.StdEncoding.DecodeString(string(rsp.Data)); err != nil { return err }
		ch := string(t)
		serverFirstMessage := ch
		tokens := scram_attributes(ch)
		var salt []byte
		if salt, err = base64.StdEncoding.DecodeString(tokens["s"]); err != nil { return err }
		i, _ := strconv.Atoi(tokens["i"])
		r := tokens["r"]
		saltedpassword := scram_Hi([]byte(jid.Resourceprep(password)), salt, i)
		clientKey := scram_HMAC(saltedpassword, []byte("Client Key"))
		storedKey := scram_H(clientKey)
		clientFinalMessageWithoutProof := "c=biws,r=" + r
		authMessage := []byte(clientFirstMessageBare + "," + serverFirstMessage + "," +
			clientFinalMessageWithoutProof)
		clientSignature := scram_HMAC(storedKey, []byte(authMessage))
		clientProof := scram_XOR(clientKey, clientSignature)
		serverKey := scram_HMAC(saltedpassword, []byte("Server Key"))
		serverSignature := scram_HMAC(serverKey, authMessage)
		response := clientFinalMessageWithoutProof + ",p=" + string(base64_encode(string(clientProof)))
		if err = strm.Send(&sasl.Response{base64_encode(response)}); err != nil { return err }
		if resp, err = strm.Read(); err != nil { return err }
		if resp == nil { return stream.ErrUnexpectedEOF }
		switch rsp := resp.(type) {
		case *sasl.Success:
			fmt.Println("success ", string(rsp.Data))
			var t []byte
			if t, err = base64.StdEncoding.DecodeString(string(rsp.Data)); err != nil { return err }
			fmt.Println(string(t))
			tokens = scram_attributes(string(t))
			v := tokens["v"]
			fmt.Println("v ", v)
			if t, err = base64.StdEncoding.DecodeString(v); err != nil { return err }
			if string(t) != string(serverSignature) {
				return errors.New("Server Signature mismatch")
			}
			
			return nil
		case *sasl.Failure:
			return rsp
		case *stream.Error:
			return rsp
		}
	case *sasl.Failure:
		return rsp
	case *stream.Error:
		return rsp
	}
	return stream.ErrMalformedProtocol
}

func get_map(data []byte) (map[string]string, error) {
	src := bytes.NewReader(data)
	var s scanner.Scanner
	var tok rune
	s.Init(src)
	var key, value string
	m := map[string]string{}
	for {
		if tok = s.Scan(); tok == scanner.Ident {
			key = s.TokenText()
			if tok = s.Scan(); tok == rune('=') {
				tok = s.Scan()
				if tok == scanner.String {
					v := []rune(s.TokenText())
					value = string(v[1:len(v)-1])
				} else {
					if tok == scanner.Ident {
						value = s.TokenText()
						for {
							if tok = s.Peek(); tok == rune(',') || tok == scanner.EOF {
								break
							}
							tok = s.Scan()
							value += s.TokenText()
						}
					}
				}
				m[key] = value
				tok = s.Scan()
				if tok == scanner.EOF {
					break
				}
				if tok == rune(',') {
					continue
				}
			}
		}
		return nil, errors.New("Failed to parse SASL challenge string")
	}
	return m, nil
}

func make_cnonce() string {
	array := make([]byte, 8)
	for i := 0; i < 8; i++ {
		array[i] = byte(rand.Intn(255))
	}
	return hex(array)
}

func md5_sum(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

func hex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

func base64_encode(str string) []byte {
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(str)))
	base64.StdEncoding.Encode(enc, []byte(str))
	return enc
}

func scram_attributes(str string) map[string]string {
	kvs := strings.Split(str, ",")
	tokens := map[string]string{}
	for _, x := range kvs {
		if x == "" {
			continue
		}
		if len(x) == 1 {
			tokens[x] = ""
			continue
		}
		str := []rune(x)
		key := str[0]
		value := str[2:]
		tokens[string(key)] = string(value)
	}
	return tokens
}

func scram_H(str []byte) []byte {
	sha := sha1.New()
	sha.Write(str)
	return sha.Sum(nil)
}

func scram_Hi(str, salt []byte, i int) []byte {
	u1 := scram_HMAC(str,  bytes.Join([][]byte{salt,[]byte{0,0,0,1}}, []byte{}))
	x := u1
	uprev := u1
	for ; i > 1; i-- {
		u2 := scram_HMAC(str, uprev)
		uprev = u2
		x = scram_XOR(x, u2)
	}
	return x
}

func scram_HMAC(str1, str2 []byte) []byte {
	mac := hmac.New(sha1.New, str1)
	mac.Write(str2)
	return mac.Sum(nil)
}

func scram_XOR(x1, x2 []byte) []byte {
	l := len(x1)
	for i := 0; i < l; i++ {
		x1[i] ^= x2[i]
	}
	return x1
}

