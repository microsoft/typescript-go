package tsoptions

import (
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ParsedCommandLine struct {
	ParsedConfig *core.ParsedOptions `json:"parsedConfig"`

	ConfigFile    *TsConfigSourceFile `json:"configFile"` // TsConfigSourceFile, used in Program and ExecuteCommandLine
	Errors        []*ast.Diagnostic   `json:"errors"`
	Raw           any                 `json:"raw"`
	CompileOnSave *bool               `json:"compileOnSave"`
	// TypeAquisition *core.TypeAcquisition

	comparePathsOptions     tspath.ComparePathsOptions
	wildcardDirectoriesOnce sync.Once
	wildcardDirectories     map[string]bool
}

// WildcardDirectories returns the cached wildcard directories, initializing them if needed
func (p *ParsedCommandLine) WildcardDirectories() map[string]bool {
	p.wildcardDirectoriesOnce.Do(func() {
		p.wildcardDirectories = getWildcardDirectories(
			p.ConfigFile.configFileSpecs.validatedIncludeSpecs,
			p.ConfigFile.configFileSpecs.validatedExcludeSpecs,
			p.comparePathsOptions,
		)
	})

	return p.wildcardDirectories
}

func (p *ParsedCommandLine) SetParsedOptions(o *core.ParsedOptions) {
	p.ParsedConfig = o
}

func (p *ParsedCommandLine) SetCompilerOptions(o *core.CompilerOptions) {
	p.ParsedConfig.CompilerOptions = o
}

func (p *ParsedCommandLine) CompilerOptions() *core.CompilerOptions {
	return p.ParsedConfig.CompilerOptions
}

func (p *ParsedCommandLine) FileNames() []string {
	return p.ParsedConfig.FileNames
}

func (p *ParsedCommandLine) ProjectReferences() []core.ProjectReference {
	return p.ParsedConfig.ProjectReferences
}

func (p *ParsedCommandLine) GetConfigFileParsingDiagnostics() []*ast.Diagnostic {
	if p.ConfigFile != nil {
		// todo: !!! should be ConfigFile.ParseDiagnostics, check if they are the same
		return slices.Concat(p.ConfigFile.SourceFile.Diagnostics(), p.Errors)
	}
	return p.Errors
}

func (p *ParsedCommandLine) MatchesFileName(fileName string, comparePathsOptions tspath.ComparePathsOptions) bool {
	path := tspath.ToPath(fileName, comparePathsOptions.CurrentDirectory, comparePathsOptions.UseCaseSensitiveFileNames)
	if slices.ContainsFunc(p.FileNames(), func(f string) bool {
		return path == tspath.ToPath(f, comparePathsOptions.CurrentDirectory, comparePathsOptions.UseCaseSensitiveFileNames)
	}) {
		return true
	}

	if p.ConfigFile == nil {
		return false
	}

	if slices.ContainsFunc(p.ConfigFile.configFileSpecs.validatedFilesSpec, func(f string) bool {
		return path == tspath.ToPath(f, comparePathsOptions.CurrentDirectory, comparePathsOptions.UseCaseSensitiveFileNames)
	}) {
		return true
	}

	if len(p.ConfigFile.configFileSpecs.validatedIncludeSpecs) == 0 {
		return false
	}

	if p.ConfigFile.configFileSpecs.matchesExclude(fileName, comparePathsOptions) {
		return false
	}

	return p.ConfigFile.configFileSpecs.matchesInclude(fileName, comparePathsOptions)
}
