package core

import "github.com/microsoft/typescript-go/internal/contentmapper"

type ParsedOptions struct {
	CompilerOptions *CompilerOptions `json:"compilerOptions"`
	WatchOptions    *WatchOptions    `json:"watchOptions"`
	TypeAcquisition *TypeAcquisition `json:"typeAcquisition"`

	FileNames         []string                `json:"fileNames"`
	ProjectReferences []*ProjectReference     `json:"projectReferences"`
	ContentMappers    []*contentmapper.Mapper `json:"contentMappers"`
}
