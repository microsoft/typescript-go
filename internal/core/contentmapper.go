package core

// ContentMapper describes an external content mapper declared at the top level of a tsconfig.
// A content mapper registers foreign file extensions whose contents are transformed into
// TypeScript during program construction.
type ContentMapper struct {
	Command    []string `json:"command"`
	Extensions []string `json:"extensions"`
}
