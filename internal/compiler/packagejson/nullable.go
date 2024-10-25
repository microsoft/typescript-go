package packagejson

import "encoding/json"

type nullableState int16

const (
	NullableStateMissing nullableState = iota
	NullableStateNull
	NullableStatePresent
)

type Nullable[T any] struct {
	state nullableState
	Value T
}

func (n *Nullable[T]) IsNull() bool {
	return n.state == NullableStateNull
}

func (n *Nullable[T]) IsMissing() bool {
	return n.state == NullableStateMissing
}

func (n *Nullable[T]) IsPresent() bool {
	return n.state == NullableStatePresent
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.state = NullableStateNull
		return nil
	}
	if err := json.Unmarshal(data, &n.Value); err != nil {
		return err
	}
	n.state = NullableStatePresent
	return nil
}
