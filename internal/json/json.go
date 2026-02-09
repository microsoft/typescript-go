//nolint:depguard
package json

import (
	"io"
	"slices"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var allowInvalid []jsonv2.Options = slices.Clip([]jsonv2.Options{jsontext.AllowInvalidUTF8(true)})

func Marshal(in any, opts ...jsonv2.Options) (out []byte, err error) {
	if len(opts) == 0 {
		opts = allowInvalid
	} else {
		opts = append(allowInvalid, opts...)
	}
	return jsonv2.Marshal(in, opts...)
}

func MarshalEncode(out *jsontext.Encoder, in any, opts ...jsonv2.Options) (err error) {
	if len(opts) == 0 {
		opts = allowInvalid
	} else {
		opts = append(allowInvalid, opts...)
	}
	return jsonv2.MarshalEncode(out, in, opts...)
}

func MarshalWrite(out io.Writer, in any, opts ...jsonv2.Options) (err error) {
	if len(opts) == 0 {
		opts = allowInvalid
	} else {
		opts = append(allowInvalid, opts...)
	}
	return jsonv2.MarshalWrite(out, in, opts...)
}

func MarshalIndent(in any, prefix, indent string) (out []byte, err error) {
	if prefix == "" && indent == "" {
		// WithIndentPrefix and WithIndent imply multiline output, so skip them.
		return Marshal(in)
	}
	return Marshal(in, jsontext.WithIndentPrefix(prefix), jsontext.WithIndent(indent))
}

func MarshalIndentWrite(out io.Writer, in any, prefix, indent string) (err error) {
	if prefix == "" && indent == "" {
		// WithIndentPrefix and WithIndent imply multiline output, so skip them.
		return MarshalWrite(out, in)
	}
	return MarshalWrite(out, in, jsontext.WithIndentPrefix(prefix), jsontext.WithIndent(indent))
}

func Unmarshal(in []byte, out any, opts ...jsonv2.Options) (err error) {
	return jsonv2.Unmarshal(in, out, opts...)
}

func UnmarshalDecode(in *jsontext.Decoder, out any, opts ...jsonv2.Options) (err error) {
	return jsonv2.UnmarshalDecode(in, out, opts...)
}

func UnmarshalRead(in io.Reader, out any, opts ...jsonv2.Options) (err error) {
	return jsonv2.UnmarshalRead(in, out, opts...)
}

func AllowDuplicateNames(allow bool) jsonv2.Options {
	return jsontext.AllowDuplicateNames(allow)
}

func WithIndent(indent string) jsonv2.Options {
	return jsontext.WithIndent(indent)
}

type (
	Value           = jsontext.Value
	UnmarshalerFrom = jsonv2.UnmarshalerFrom
	MarshalerTo     = jsonv2.MarshalerTo
	Decoder         = jsontext.Decoder
	Encoder         = jsontext.Encoder
)

var (
	BeginObject = jsontext.BeginObject
	EndObject   = jsontext.EndObject
	Null        = jsontext.Null
	BeginArray  = jsontext.BeginArray
	EndArray    = jsontext.EndArray
)
