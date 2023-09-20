package jsong

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

// Decoder decodes JSONG from JSON.
type Decoder struct {
	r *bufio.Reader
}

// NewDecoder creates a new JSONG decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (d *Decoder) skipWhitespace() error {
	for {
		r, _, err := d.r.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			return d.r.UnreadRune()
		}
	}
}

func (d *Decoder) Decode() (any, error) {
	if err := d.skipWhitespace(); err != nil {
		return nil, err
	}
	b, err := d.r.Peek(1)
	if err != nil {
		return nil, err
	}
	switch b[0] {
	case 'n':
		return d.decodeNull()
	case 'f':
		return d.decodeFalse()
	case 't':
		return d.decodeTrue()
	case '-':
		d.r.Discard(1)
		n, err := d.decodeNumber()
		if err != nil {
			return nil, err
		}
		return -n.(num), nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.decodeNumber()
	case '"':
		return d.decodeString()
	case '[':
		return d.decodeArray()
	case '{':
		return d.decodeObject()
	default:
		return nil, fmt.Errorf("invalid character %q looking for beginning of value", b[0])
	}
}

var (
	nullData  = []byte("null")
	falseData = []byte("false")
	trueData  = []byte("true")
)

func (d *Decoder) decodeNull() (any, error) {
	bs, err := d.r.Peek(4)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(bs, nullData) {
		d.r.Discard(4)
		return null{}, nil
	}
	var j int
	for i := range bs {
		if bs[i] != nullData[i] {
			j = i
		}
	}
	return nil, fmt.Errorf("invalid character %q in literal null (expecting %q)", bs[j], nullData[j])
}

func (d *Decoder) decodeFalse() (any, error) {
	bs, err := d.r.Peek(5)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(bs, falseData) {
		d.r.Discard(5)
		return boolean(false), nil
	}
	var j int
	for i := range bs {
		if bs[i] != falseData[i] {
			j = i
		}
	}
	return nil, fmt.Errorf("invalid character %q in literal false (expecting %q)", bs[j], falseData[j])
}

func (d *Decoder) decodeTrue() (any, error) {
	bs, err := d.r.Peek(4)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(bs, trueData) {
		d.r.Discard(4)
		return boolean(true), nil
	}
	var j int
	for i := range bs {
		if bs[i] != trueData[i] {
			j = i
		}
	}
	return nil, fmt.Errorf("invalid character %q in literal true (expecting %q)", bs[j], trueData[j])
}

func (d *Decoder) decodeNumber() (any, error) {
	var dot bool
	var exp bool
	var expSign bool
	buf := bytes.NewBuffer(make([]byte, 0, 5))
loop:
	for {
		b, err := d.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch {
		case b == '-', b == '+':
			if dot && !exp {
				return nil, fmt.Errorf("invalid character %q after decimal point in numeric literal", b)
			}
			if expSign {
				return nil, fmt.Errorf("invalid character %q in exponent of numeric literal", b)
			}
			expSign = true
		case b == 'e', b == 'E':
			if exp {
				return nil, fmt.Errorf("invalid character %q in exponent of numeric literal", b)
			}
			exp = true
		case b == '.':
			if exp {
				return nil, fmt.Errorf("invalid character %q in exponent of numeric literal", b)
			}
			if dot {
				return nil, fmt.Errorf("invalid character %q after decimal point in numeric literal", b)
			}
			dot = true
		case '0' <= b && b <= '9':
		default:
			d.r.UnreadByte()
			break loop
		}
		buf.WriteByte(b)
	}

	v, err := strconv.ParseFloat(buf.String(), 64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode number %q: %v", buf.String(), err) // Catch parsing issues here.
	}
	return num(v), nil
}

func (d *Decoder) decodeString() (any, error) {
	d.r.Discard(1) // '"'
	buf := bytes.NewBuffer(make([]byte, 0, 16))
	buf.WriteByte('"')
	for {
		bs, err := d.r.ReadSlice('"')
		if err != nil {
			return nil, err
		}
		buf.Write(bs)
		if len(bs) == 1 || bs[len(bs)-2] != '\\' {
			break
		}
	}
	data, err := unquoteInPlace(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to unquote string %q: %v", buf.String(), err)
	}
	// Using string first helps the compiler find the optimization.
	return str(string(data)), nil
}

func (d *Decoder) decodeArray() (any, error) {
	res := array{}
	d.r.Discard(1) // '['
	for i := 0; ; i++ {
		if err := d.skipWhitespace(); err != nil {
			return nil, err
		}
		if i > 0 {
			b, err := d.r.ReadByte()
			if err != nil {
				return nil, err
			}
			if b == ']' {
				break
			}
			if b != ',' {
				return nil, fmt.Errorf("invalid character %q after array element", b)
			}
		} else {
			b, err := d.r.Peek(1)
			if err != nil {
				return nil, err
			}
			if b[0] == ']' {
				d.r.Discard(1)
				break
			}
		}
		e, err := d.Decode()
		if err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (d *Decoder) decodeObject() (any, error) {
	res := object{}
	d.r.Discard(1) // '{'
	for i := 0; ; i++ {
		if err := d.skipWhitespace(); err != nil {
			return nil, err
		}
		bs, err := d.r.Peek(1)
		if err != nil {
			return nil, err
		}
		if bs[0] == '}' {
			d.r.Discard(1)
			break
		}
		if i > 0 {
			if bs[0] != ',' {
				return nil, fmt.Errorf("invalid character %q after array element", bs[0])
			}
			d.r.Discard(1)
			if err := d.skipWhitespace(); err != nil {
				return nil, err
			}
			bs, err = d.r.Peek(1)
			if err != nil {
				return nil, err
			}
		}
		if bs[0] != '"' {
			return nil, fmt.Errorf("invalid character %q looking for beginning of object key string", bs[0])
		}
		k, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		if err := d.skipWhitespace(); err != nil {
			return nil, err
		}
		b, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b != ':' {
			return nil, fmt.Errorf("invalid character %q after object key", b)
		}
		v, err := d.Decode()
		if err != nil {
			return nil, err
		}
		res[string(k.(str))] = v // Repeat keys are ok.
	}
	return res, nil
}

var emptyStr = []byte(`""`)

func unquoteInPlace(b []byte) ([]byte, error) {
	if len(b) < 2 {
		return nil, strconv.ErrSyntax
	}
	if len(b) == 2 && bytes.Equal(b, emptyStr) {
		return []byte{}, nil
	}
	b[0] = b[1]
	i := 2
	if b[0] == '\\' {
		i++
		b[0] = b[2]
	}
	end := 1
	escape := false
	for ; i < len(b); i++ {
		switch b[i] {
		case '\\':
			if escape {
				escape = false
				b[end] = b[i]
				end++
				continue
			}
			escape = true
		case '"':
			if escape {
				escape = false
				b[end] = b[i]
				end++
				continue
			}
			return b[:end], nil
		default:
			escape = false
			b[end] = b[i]
			end++
		}
	}
	return nil, io.EOF
}
