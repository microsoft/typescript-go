// Package contentmapper defines the types describing an external content mapper: a plugin that
// transforms foreign file content (e.g. .vue) into TypeScript during program construction.
//
// A mapper is declared in tsconfig (Definition), its implementation is described by fields in its npm
// package's package.json (Manifest), and the two are combined once the package is resolved (Mapper).
// Resolution itself lives in the tsoptions package (it needs node module resolution), keeping this
// package free of dependencies so lower layers can reference these types.
package contentmapper

// Definition is a content mapper as declared in a tsconfig's "contentMappers": the npm package that
// implements the mapper and the foreign file extensions it registers.
type Definition struct {
	Package    string   `json:"package"`
	Extensions []string `json:"extensions"`
}

// Manifest is the content-mapper information read from a package's package.json: its name and version
// (which form the mapper's identity) and the argv used to run it.
type Manifest struct {
	Name    string
	Version string
	Exec    []string
}

// Mapper is a resolved content mapper: its tsconfig Definition combined with the Manifest resolved from
// the package's package.json, plus the package directory used as the mapper's working directory.
type Mapper struct {
	Definition
	Manifest         `json:"-"`
	PackageDirectory string `json:"-"`
}

// Identity returns the mapper's "name@version" identity, or just the name when it declares no version,
// or an empty string when the mapper has not been resolved to a name.
func (m *Mapper) Identity() string {
	switch {
	case m.Name == "":
		return ""
	case m.Version == "":
		return m.Name
	default:
		return m.Name + "@" + m.Version
	}
}
