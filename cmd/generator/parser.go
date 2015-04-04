package main

import (
	"fmt"
	"errors"
	"io/ioutil"
	"strings"
)

var (
	ErrUnexpectedEOF = errors.New("Unexpected EOF")
	ErrBadSyntax = errors.New("bad syntax")
)

type Buf struct {
	buf []byte
	len int
	idx int
}

func (schema *Schema) ParseFile(filename string) error {
	fmt.Println("Parsing file ", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil { return err }
	buf := &Buf{data, len(data), 0}
	if err = schema.parseSchema(buf); err != nil {
		line, col := getPos(buf)
		fmt.Printf("at line %d:%d: %s\n", line, col, err)
		return err
	}
	return nil
}

func (schema *Schema) parseSchema(buf *Buf) error {
	for {
		token := getToken(buf)
		switch token {
		case "":
			return nil
		case "\n":
		case "package":
			next := getToken(buf)
			schema.PackageName = next
		case "category":
			next := getToken(buf)
			schema.Props["category"] = next
		case "documentation":
			next := getToken(buf)
			schema.Props["documentation"] = next
		case "targetNamespace":
			space := getToken(buf)
			prefix := ""
			next := getToken(buf)
			switch next {
			case "": return ErrUnexpectedEOF
			case "{":
			default:
				prefix = space
				space = next
				next = getToken(buf)
			}
			if next != "{" {
				return ErrBadSyntax
			}
			target := &Target{}
			target.Props = make(map[string]string)
			target.Prefix = prefix
			target.Space = space
			target.parseDefs(buf)
			schema.Targets = append(schema.Targets, target)
		default:
			next := getRest(buf)
			schema.Props[token] = next
		}
	}
	return nil
}

func (target *Target) parseDefs(buf *Buf) error {
Loop:
	for {
		skipSpaceWithNewline(buf)
		ident := getToken(buf)
		switch ident {
		case "":
			return ErrUnexpectedEOF
		case "}":
			break Loop
		default:
			next := getToken(buf)
			switch next {
			case "":
				return ErrUnexpectedEOF
			case "=":
				rest := getRest(buf)
				target.Props[ident] = rest
			case "::=":
				var reciver string
				b := []rune(ident)[0]
				switch b {
				case '*': reciver = "both"
				case '+': reciver = "server"
				case '-': reciver = "client"
				}
				if reciver != "" {
					ident = string([]rune(ident)[1:])
				}
				field := Field{Name:ident, Reciver_type:reciver}
				if err := target.parseDef(buf, &field); err != nil { return err }
			default:
				fmt.Println("unknown token " + next)
				return ErrBadSyntax
			}
		}
	}
	return nil
}

func getPos(buf *Buf) (int, int) {
	var line, col int = 1,0
	for i := 0; i < buf.idx; i++ {
		col++
		if buf.buf[i] == '\n' {
			col = 0
			line++
		}
	}
	return line, col
}

func getToken(buf *Buf) string {
	skipSpace(buf)
	if buf.idx >= buf.len {
		return ""
	}
	var ident string
	i := buf.idx
	b := buf.buf[i]
	switch b {
	case '{':
		ident = "{"
		buf.idx++
		skipSpaceWithNewline(buf)
		i = buf.idx
	case '=':
		ident = string(b)
		i++
	case ':':
		i++
		if i < buf.len && buf.buf[i] == ':' {
			i++
			if i < buf.len && buf.buf[i] == '=' {
				ident = "::="
				i++
			} else {
				ident = "::"
			}
		}
	case '\n':
		ident = "\n"
		i++
	case '(', ')':
		ident = string(b)
		i++
	default:
		for ; i < buf.len; i++ {
			b := buf.buf[i]
			if b == ':' {
				if i+2 < buf.len && buf.buf[i+1] == ':' && buf.buf[i+2] == '=' {
					break
				} else {
					continue
				}
			}
			if b == ' ' || b == '{' || b == '\t' || b == '=' || b == '\n' ||
				b == '(' || b == ')' {
				break
			}
		}
		if i > buf.idx {
			ident = string(buf.buf[buf.idx:i])
		}
	}
	buf.idx = i
	skipSpace(buf)
	return ident
}

func skipSpace(buf *Buf) {
	if buf.idx < buf.len && buf.buf[buf.idx] == '#' {
		for ; buf.idx < buf.len; buf.idx++ {
			if buf.buf[buf.idx] == '\n' {
				break
			}
		}
	}
	for ; buf.idx < buf.len; buf.idx++ {
		b := buf.buf[buf.idx]
		if b != ' ' && b != '\t' {
			break
		}
	}
}

func skipSpaceWithNewline(buf *Buf) {
	for ; buf.idx < buf.len; buf.idx++ {
		b := buf.buf[buf.idx]
		if b == '#' {
			for ; buf.idx < buf.len; buf.idx++ {
				if buf.buf[buf.idx] == '\n' {
					b = '\n'
					break
				}
			}
		}
		if b != ' ' && b != '\t' && b != '\n' {
			break
		}
	}
}

func getRest(buf *Buf) string {
	old := buf.idx
	for i := buf.idx; i < buf.len; i++ {
		b := buf.buf[i]
		if b == '\n' {
			buf.idx = i+1
			break
		}
	}
	if old < buf.idx {
		rest := string(buf.buf[old:buf.idx])
		return strings.TrimSpace(rest)
	}
	return ""
}

func (target *Target) parseDef(buf *Buf, field *Field) error {
	next := getToken(buf)
	switch next {
	case "":
		return ErrUnexpectedEOF
	case "\n":
		return ErrBadSyntax
	case "extension":
		next = getToken(buf)
		switch next {
		case "":
			return ErrUnexpectedEOF
		case "\n":
			field.Type = Extension{}
		case "(":
			space := getToken(buf)
			local := getToken(buf)
			closed := getToken(buf)
			if closed != ")" {
				return ErrBadSyntax
			}
			field.Type = Extension{space, local}
			next = getToken(buf)
		}
	case "set":
		set, err := parseSet(buf)
		if err != nil { return err }
		field.Type = set
		next = getToken(buf)
	case "choice":
		choice, err := parseChoice(buf)
		if err != nil { return err }
		field.Type = choice
		next = getToken(buf)
	case "sequence":
		sequence, err := parseSequence(buf)
		if err != nil { return err }
		field.Type = sequence
		next = getToken(buf)
	case "enum":
		enums, err := getEnum(buf)
		if err != nil { return err }
		field.Type = Enum(enums)
		next = getToken(buf)
	default:
		field.Type = next
		next = getToken(buf)
	}
	var enc *Encoding
	switch next {
	case "":
		return ErrUnexpectedEOF
	case "\n":
		break
	default:
		enctype := next
		var encspace, encname string
		next = getToken(buf)
		if next == "(" {
			encname = getToken(buf)
			next = getToken(buf)
			if next != ")" {
				encspace = encname
				encname = next
				next = getToken(buf) // ")"
			}
			next = getToken(buf) // "\n"
		}
		enc = &Encoding{enctype, encspace, encname}
	}
	field.EncodingRule = enc
	if next != "\n" {
		fmt.Println("expected \\n ", next)
		return ErrBadSyntax
	}
	target.Fields = append(target.Fields, field)
	return nil
}

func getEnum(buf *Buf) ([]string, error) {
	bracket := getToken(buf)
	if bracket != "{" {
		return nil, nil
	}
	var values []string
Loop:
	for {
		token := getToken(buf)
		switch token {
		case "":
			return nil, errors.New("unexpected EOF")
		case "}":
			break Loop
		case "\n":
		default:
		values = append(values, token)
		}
	}
	return values, nil
}

func parseSequence(buf *Buf) (interface{}, error) {
	switch getToken(buf) {
	case "":
		return nil, ErrUnexpectedEOF
	case "of":
		next := getToken(buf)
		if next == "" {
			return nil, ErrUnexpectedEOF
		}
		return SequenceOf(next), nil
	case "{":
		fields, err := parseFields(buf)
		if err != nil { return nil, err }
		return Sequence(fields), nil
	}
	return nil, ErrBadSyntax
}

func parseSet(buf *Buf) (Set, error) {
	bracket := getToken(buf)
	if bracket == "" {
		return nil, ErrUnexpectedEOF
	}
	if bracket != "{" {
		return nil, ErrBadSyntax
	}
	fields, err := parseFields(buf)
	if err != nil { return nil, err }
	return Set(fields), nil
}

func parseChoice(buf *Buf) (Choice, error) {
	bracket := getToken(buf)
	if bracket == "" {
		return nil, ErrUnexpectedEOF
	}
	if bracket != "{" {
		return nil, ErrBadSyntax
	}
	fields, err := parseFields(buf)
	if err != nil { return nil, err }
	return Choice(fields), nil
}

func parseFields(buf *Buf) ([]*Field, error) {
	var fields []*Field

	for {
		skipSpaceWithNewline(buf)
		name := getToken(buf)
		if name == "" {
			return nil, ErrUnexpectedEOF
		}
		if name == "}" {
			break
		}
		var typ interface{}
		next := getToken(buf)
		switch next {
		case "":
			return nil, ErrUnexpectedEOF
		case "\n":
			fields = append(fields, &Field{Type:name})
			continue
		case "extension":
			next = getToken(buf)
			switch next {
			case "":
				return nil, ErrUnexpectedEOF
			case "\n":
				typ = Extension{}
			case "(":
				space := getToken(buf)
				local := getToken(buf)
				closed := getToken(buf)
				if closed != ")" {
					return nil, ErrBadSyntax
				}
				typ = Extension{space, local}
				next = getToken(buf)
			}
		case "set":
			set, err := parseSet(buf)
			if err != nil { return nil, err }
			typ = set
			next = getToken(buf)
		case "choice":
			choice, err := parseChoice(buf)
			if err != nil { return nil, err }
			typ = choice
			next = getToken(buf)
		case "sequence":
			sequence, err := parseSequence(buf)
			if err != nil { return nil, err }
			typ = sequence
			next = getToken(buf)
		case "enum":
			enums, err := getEnum(buf)
			if err != nil { return nil, err }
			typ = Enum(enums)
			next = getToken(buf)
		default:
			typ = next
			next = getToken(buf)
		}
		
		var encoding *Encoding
		required := false
		defaultValue := ""
		switch next {
		case "":
			return nil, ErrUnexpectedEOF
		case "\n":
			break
		default:
			enctype := next
			var encspace, encname string
			next = getToken(buf)
			if next == "(" {
				encname = getToken(buf)
				next = getToken(buf)
				if next != ")" {
					encspace = encname
					encname = next
					next = getToken(buf) // ")"
				}
			}
			encoding = &Encoding{enctype, encspace, encname}
		Rest:
			for {
				switch next {
				case "":
					return nil, ErrUnexpectedEOF
				case "\n":
					break Rest
				case "required":
					required = true
				case "default":
					next = getToken(buf)
					switch next {
					case "":
						return nil, ErrUnexpectedEOF
					case "\n":
						return nil, ErrBadSyntax
					default:
						defaultValue = next
					}
				default:
					return nil, errors.New("Unknown tag " + next)
				}
				next = getToken(buf)
			}
			if next != "\n" {
				return nil, ErrBadSyntax
			}
		}
		fields = append(fields, &Field{Name:name, Type:typ, Required:required,
			EncodingRule:encoding, DefaultValue:defaultValue})
	}
	return fields, nil
}

/*
func parseType(buf *Buf) (interface{}, error) {
	var typ interface{}
	next := getToken(buf)
	switch next {
	case "":
		return nil, ErrUnexpectedEOF
	case "set":
		set, err := parseSet(buf)
		if err != nil { return nil, err }
		typ = set
		case "choice":
		choice, err := parseChoice(buf)
		if err != nil { return nil, err }
		typ = choice
	case "sequence":
		sequence, err := parseSequence(buf)
		if err != nil { return nil, err }
		typ = sequence
	case "enum":
		enums, err := getEnum(buf)
		if err != nil { return nil, err }
		typ = Enum(enums)
	default:
		typ = next
	}
	return typ, nil
}
*/
