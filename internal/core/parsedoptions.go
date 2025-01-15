package core

import "github.com/microsoft/typescript-go/internal/collections"

type ParsedOptions struct {
	CompilerOptions   *CompilerOptions
	FileNames         collections.OrderedMap[string, string]
	ProjectReferences []ProjectReference
}
