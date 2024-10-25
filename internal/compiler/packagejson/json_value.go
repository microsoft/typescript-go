package packagejson

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type JSONValueType int16

const (
	JSONValueTypeNull JSONValueType = iota
	JSONValueTypeString
	JSONValueTypeNumber
	JSONValueTypeBoolean
	JSONValueTypeArray
	JSONValueTypeObject
)

type JSONValue struct {
	Type  JSONValueType
	Value any
}

func (v *JSONValue) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		v.Type = JSONValueTypeNull
		return nil
	}
	if data[0] == '"' {
		v.Type = JSONValueTypeString
		return json.Unmarshal(data, &v.Value)
	}
	if data[0] == '[' {
		var elements []JSONValue
		if err := json.Unmarshal(data, &elements); err != nil {
			return err
		}
		v.Type = JSONValueTypeArray
		v.Value = elements
		return nil
	}
	if data[0] == '{' {
		var object map[string]JSONValue
		if err := json.Unmarshal(data, &object); err != nil {
			return err
		}
		v.Type = JSONValueTypeObject
		v.Value = object
		return nil
	}
	v.Type = JSONValueTypeNumber
	return json.Unmarshal(data, &v.Value)
}

func (v *JSONValue) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	switch dec.PeekKind() {
	case jsontext.Null.Kind():
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		v.Type = JSONValueTypeNull
		return nil
	case '"':
		v.Type = JSONValueTypeString
		if err := json.UnmarshalDecode(dec, &v.Value, opts); err != nil {
			return err
		}
	case '[':
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		var elements []JSONValue
		for dec.PeekKind() != jsontext.ArrayEnd.Kind() {
			var element JSONValue
			if err := json.UnmarshalDecode(dec, &element, opts); err != nil {
				return err
			}
			elements = append(elements, element)
		}
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		v.Type = JSONValueTypeArray
		v.Value = elements
	case '{':
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		object := make(map[string]JSONValue)
		for dec.PeekKind() != jsontext.ObjectEnd.Kind() {
			var key string
			var value JSONValue
			if err := json.UnmarshalDecode(dec, &key, opts); err != nil {
				return err
			}
			if err := json.UnmarshalDecode(dec, &value, opts); err != nil {
				return err
			}
			object[key] = value
		}
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		v.Type = JSONValueTypeObject
		v.Value = object
	default:
		v.Type = JSONValueTypeNumber
		if err := json.UnmarshalDecode(dec, &v.Value, opts); err != nil {
			return err
		}
	}
	return nil
}
