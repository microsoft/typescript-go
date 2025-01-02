package core

type ParsedOptions struct {
	Options *CompilerOptions
	FileNames     []string
	Raw           any
	CompileOnSave *bool
	ProjectReferences []ProjectReference
}
