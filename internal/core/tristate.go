package core

type Tristate byte

const (
	TSUnknown Tristate = iota
	TSFalse
	TSTrue
)

func (t *Tristate) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "true":
		*t = TSTrue
	case "false":
		*t = TSFalse
	default:
		*t = TSUnknown
	}
	return nil
}
