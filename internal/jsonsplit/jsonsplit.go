package jsonsplit

import (
	"bytes"
	"encoding/json"

	"github.com/go-json-experiment/jsonsplit"
)

var codec jsonsplit.Codec

func init() {
	codec.SetMarshalCallMode(jsonsplit.CallBothButReturnV2)
	codec.SetUnmarshalCallMode(jsonsplit.CallBothButReturnV2)

	codec.ReportDifference = func(d jsonsplit.Difference) {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "    ")
		_ = enc.Encode(d)

		panic("jsonsplit: " + buf.String())
	}

	codec.AutoDetectOptions = true
}

var (
	Marshal   = codec.Marshal
	Unmarshal = codec.Unmarshal
)
