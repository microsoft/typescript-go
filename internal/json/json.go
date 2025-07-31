package json

import (
	"encoding/json"

	jsonv2 "github.com/go-json-experiment/json"
)

func Marshal(in any, opts ...jsonv2.Options) (out []byte, err error) {
	return jsonv2.Marshal(in, opts...)
}

func Unmarshal(in []byte, out any, opts ...jsonv2.Options) (err error) {
	return jsonv2.Unmarshal(in, out, opts...)
}

type RawMessage = json.RawMessage
