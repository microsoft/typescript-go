package core

// ContentMapper describes an external content mapper declared at the top level of a tsconfig.
// A content mapper registers foreign file extensions whose contents are transformed into
// TypeScript during program construction.
type ContentMapper struct {
	Command    []string `json:"command"`
	Extensions []string `json:"extensions"`
	// Name and Version identify the mapper implementation. They are reported by the mapper itself
	// (e.g. during an initialize handshake) rather than declared in the tsconfig, and are used to
	// attribute diagnostics to a specific mapper.
	Name    string `json:"-"`
	Version string `json:"-"`
}

// Identity returns a human-readable "name@version" identifier for the mapper, or just the name when no
// version is reported, or an empty string when the mapper has not identified itself.
func (m *ContentMapper) Identity() string {
	switch {
	case m.Name == "":
		return ""
	case m.Version == "":
		return m.Name
	default:
		return m.Name + "@" + m.Version
	}
}
