package core

import "fmt"

// PanicWithStack wraps a recovered panic value along with the original stack trace
// captured at the site of the panic. This allows re-panicking while preserving the
// original stack context, rather than losing it at the re-panic site.
type PanicWithStack struct {
	Value any
	Stack []byte
}

func (p *PanicWithStack) String() string {
	return fmt.Sprintf("%v\n%s", p.Value, string(p.Stack))
}
