package xmlencoder

import (
	"encoding/xml"
	"io"
	"bufio"
	"errors"
)

const (
	ns_xml = "http://www.w3.org/XML/1998/namespace"
)

type Extension interface {
	Decode(d *Decoder, tag *xml.StartElement) error
	Encode(e *Encoder) error
}

type Encoder struct {
	w *bufio.Writer
	namespaces map[string]string
	openStartTag bool
	local_namespaces []*[]string
	depth int
	tags []string
}

func NewEncoder(w io.Writer) *Encoder {
	enc := Encoder{w: bufio.NewWriter(w), 
		namespaces : map[string]string{},
		depth : 0,
		local_namespaces : []*[]string{&[]string{}}}
	enc.namespaces[ns_xml] = "xml"
	return &enc
}

func (e *Encoder) SetPrefix(prefix string, uri string) (err error) {
	if e.openStartTag {
		if prefix == "" {
			return errors.New("empty prefix in opened tag")
		}
		_, err = e.w.WriteString("xmlns:" + prefix + `="` + uri + `"`)
		if err != nil { return err }
	}
	e.namespaces[uri] = prefix
	if len(e.local_namespaces) <= e.depth {
		e.local_namespaces = append(e.local_namespaces, &[]string{})
	}
	lns := e.local_namespaces[e.depth]
	*lns = append(*lns, uri)
	return
}

func (e *Encoder) UnsetPrefix(prefix string) {
	for u, p := range e.namespaces {
		if prefix == p {
			delete(e.namespaces, u)
		}
	}
}

func (e *Encoder) StartStream() error {
	_, err := e.w.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	return err
}

func (e *Encoder) EndStream() error {
	if len(e.tags) > 0 {
		for i := e.depth; i > 0; i-- {
			if err := e.EndElement(); err != nil { return err }
		}
	}
	return e.w.Flush()
}

func (e *Encoder) close_start_tag() (err error) {
	if e.openStartTag {
		_, err = e.w.WriteString(">")
		e.openStartTag = false
		e.depth++
	}
	return
}

func (e *Encoder) Flush() error {
	if err := e.close_start_tag(); err != nil { return err }
	return e.w.Flush()
}

func (e *Encoder) EndElement() (err error) {
	if e.openStartTag {
		if _, err = e.w.WriteString("/>"); err != nil { return }
		e.openStartTag = false
	} else {
		e.depth--
		tag := e.tags[e.depth]
		if _, err = e.w.WriteString("</" + tag + ">"); err != nil { return }
	}
	lns := e.local_namespaces[e.depth]
	for _, x := range *lns {
		delete(e.namespaces, x)
	}
	e.local_namespaces = e.local_namespaces[:e.depth]
	e.tags = e.tags[:e.depth]
	return nil
}

func (e *Encoder) StartElement(space, local string) error {
	var err error
	name := local
	if err = e.close_start_tag(); err != nil { return err }
	if len(e.local_namespaces) < e.depth+1 {
		e.local_namespaces = append(e.local_namespaces, &[]string{})
	}
	lns := e.local_namespaces[e.depth]
	if space != "" {
		prefix, ok := e.namespaces[space]
		if !ok {
			e.namespaces[space] = ""
			*lns = append(*lns, space)
		}
		if prefix != "" {
			name = prefix + ":" + local
		}
	}
	e.tags = append(e.tags, name)
	if _, err = e.w.WriteString("<" + name); err != nil { return err }
	e.openStartTag = true
	for _, x := range *lns {
		prefix := e.namespaces[x]
		var aname string
		if prefix != "" {
			aname = " xmlns:" + prefix
		} else {
			aname = " xmlns"
		}
		if _, err = e.w.WriteString(aname + `="` + x + `"`); err != nil { return err }
	}
	return nil
}

func (e *Encoder) SimpleElement(space, local, value string) error {
	var err error
	if err = e.StartElement(space, local); err != nil { return err }
	if err = e.Text(value); err != nil { return err }
	if err = e.EndElement(); err != nil { return err }
	return nil
}

func (e *Encoder) Text(text string) (err error) {
	if err = e.close_start_tag(); err != nil { return err }
	_, err = e.w.WriteString(text)
	return
}

func (e *Encoder) Bytes(data []byte) (err error) {
	if err = e.close_start_tag(); err != nil { return err }
	if _, err = e.w.Write(data); err != nil { return err }
	return
}

func (e *Encoder) Attribute(space, local, value string) error {
	var err error
	if e.openStartTag {
		aname := local
		if space != "" {
			prefix, ok := e.namespaces[space]
			if !ok || prefix == "" {
				return errors.New("no prefix for namespace " + space)
			}
			aname = prefix + ":" + local
		}
		if _, err = e.w.WriteString(" " + aname + `="`); err == nil {
			if err = xml.EscapeText(e.w, []byte(value)); err == nil {
				_, err = e.w.WriteString(`"`)
			}
		}
	} else {
		err = errors.New(`attribute after ">"`)
	}
	return err
}

