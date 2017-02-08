package jid

import "strings"
import "errors"

type JID struct {
	Node string
	Domain string
	Resource string
}

func New(jid_string string) (*JID, error) {
	var node, domain, resource string = "", "", ""

	a := strings.SplitN(jid_string, "/", 2)
	if len(a) == 2 {
		resource = a[1]
	}
	a = strings.SplitN(a[0], "@", 2)
	if len(a) == 2 {
		node = a[0]
		domain = a[1]
	} else {
		domain = a[0]
	}
	return &JID{
		Node : Nodeprep(node),
		Domain : Nameprep(domain),
		Resource : Resourceprep(resource),
	}, nil
}

func Nodeprep(s string) (string) {
	return string(nodeprep([]rune(s)))
}

func Resourceprep(s string) (string) {
	return string(resourceprep([]rune(s)))
}
	
func Nameprep(s string) (string) {
	return string(nameprep([]rune(s)))
}

func (jid *JID) GetIDN() (string, error) {
	// already nameprepared
	parts := strings.Split(jid.Domain, ".")
	var bad = false
	for i, part := range parts {
		if part == "" {
			return "", errors.New("Two dots")
		}
		rs := []rune(part)
		if rs[0] == 0x2D {
			return "", errors.New("'-' as first character")
		}
		for _, c := range rs {
			if c < 0x2D || (c > 0x2D && c < 0x30) || 
				(c > 0x39 && c < 0x41) || 
				(c > 0x5C && c < 0x61) || 
				(c > 0x7C && c < 0x80) {
				bad = true
				break
			}
		}
		if rs[len(rs)-1] == 0x2D {
			bad = true
		}
		if bad {
			return "", errors.New("bad symbols")
		}
		var contains_outside_ascii = false
		for _, c := range rs {
			if c > 0x7F {
				contains_outside_ascii = true
				break
			}
		}
		if contains_outside_ascii {
			if strings.HasPrefix(part, "xn--") {
				return "", errors.New("invalid use of ACE prefix")
			}
			var epart, err = punycode_encode(rs)
			if err != nil {
				return "", err
			}
			parts[i] = "xn--" + string(epart)
		}
	}
	return strings.Join(parts, "."), nil
}

	func (jid *JID) Bare() string {
		if jid.Node != "" {
			return jid.Node + "@" + jid.Domain
		} else {
			return jid.Domain
		}
}

	func (jid *JID) String() string {
		result := jid.Domain
		if jid.Node != "" {
			result = jid.Node + "@" + jid.Domain
		}
		if jid.Resource != "" {
			result += "/" + jid.Resource
		}
		return result
}

func (jid *JID) FromString(s string) (interface{}, error) {
	return New(s)
}

func (jid *JID) Server() *JID {
	return &JID{
		Domain: jid.Domain,
	}
}
