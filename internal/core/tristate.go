package core

//go:generate go run golang.org/x/tools/cmd/stringer -type=Tristate -output=tristate_stringer_generated.go

// Tristate

type Tristate byte

const (
	TSUnknown Tristate = iota
	TSFalse
	TSTrue
)
