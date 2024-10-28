package packagejson

type HeaderFields struct {
	Name    Expected[string] `json:"name"`
	Version Expected[string] `json:"version"`
	Type    Expected[string] `json:"type"`
}

type PathFields struct {
	TSConfig      Expected[string]           `json:"tsconfig"`
	Main          Expected[string]           `json:"main"`
	Types         Expected[string]           `json:"types"`
	Typings       Nullable[Expected[string]] `json:"typings"`
	TypesVersions any                        `json:"typesVersions"`
	Imports       any                        `json:"imports"`
	Exports       Exports                    `json:"exports"`
}

type DependencyFields struct {
	Dependencies         Expected[map[string]string] `json:"dependencies"`
	PeerDependencies     Expected[map[string]string] `json:"peerDependencies"`
	OptionalDependencies Expected[map[string]string] `json:"optionalDependencies"`
}

type Fields struct {
	HeaderFields
	PathFields
	DependencyFields
}
