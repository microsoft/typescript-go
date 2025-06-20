package projectv2

import "github.com/microsoft/typescript-go/internal/tsoptions"

type Kind int

const (
	KindInferred Kind = iota
	KindConfigured
)

type Project struct {
	Kind        Kind
	CommandLine *tsoptions.ParsedCommandLine
}
