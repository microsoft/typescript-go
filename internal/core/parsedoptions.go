package core

type ParsedOptions struct {
	CompilerOptions *CompilerOptions `json:"CompilerOptions"`
	WatchOptions    *WatchOptions    `json:"WatchOptions"`

	FileNames         []string           `json:"FileNames"`
	ProjectReferences []ProjectReference `json:"ProjectReferences"`
}
