package jsonsplit

import (
	"bytes"
	"encoding/json"
	"slices"

	"github.com/go-json-experiment/jsonsplit"
)

var codec jsonsplit.Codec

func init() {
	codec.SetMarshalCallMode(jsonsplit.CallBothButReturnV2)
	codec.SetUnmarshalCallMode(jsonsplit.CallBothButReturnV2)

	codec.ReportDifference = func(d jsonsplit.Difference) {
		options := slices.Collect(d.OptionNames())
		if len(options) == 1 && options[0] == "jsonv1.ReportErrorsWithLegacySemantics" {
			return
		}

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
