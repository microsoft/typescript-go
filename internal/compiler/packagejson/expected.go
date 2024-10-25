package packagejson

import "encoding/json"

type Expected[T any] struct {
	Valid bool
	Value T
}

func (e *Expected[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &e.Value); err != nil {
		e.Valid = false
		return nil
	}
	e.Valid = true
	return nil
}
