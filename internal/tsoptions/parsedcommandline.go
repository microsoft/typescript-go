package tsoptions

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

type ParsedCommandLine struct {
	ParsedOptions *core.ParsedOptions
	// WatchOptions WatchOptions

	ConfigFile *ast.SourceFile // TsConfigSourceFile, used in Program and ExecuteCommandLine
	Errors     []*ast.Diagnostic
	Raw        any
	// WildcardDirectories map[string]watchDirectoryFlags
	CompileOnSave *bool
	// TypeAquisition *core.TypeAcquisition
}

func NewParsedCommandLine(
	options *core.ParsedOptions,
	configFile *ast.SourceFile,
	errors []*ast.Diagnostic,
	raw any,
	compileOnSave *bool,
) ParsedCommandLine {
	return ParsedCommandLine{
		ParsedOptions: options,
		ConfigFile:    configFile,
		Errors:        errors,
		Raw:           raw,
		CompileOnSave: compileOnSave,
	}
}

func (p *ParsedCommandLine) SetParsedOptions(o *core.ParsedOptions) {
	p.ParsedOptions = o
}

func (p *ParsedCommandLine) SetCompilerOptions(o *core.CompilerOptions) {
	p.ParsedOptions.CompilerOptions = o
}

func (p *ParsedCommandLine) CompilerOptions() *core.CompilerOptions {
	return p.ParsedOptions.CompilerOptions
}

func (p *ParsedCommandLine) FileNames() []string {
	return p.ParsedOptions.FileNames
}

func (p *ParsedCommandLine) ProjectReferences() []core.ProjectReference {
	return p.ParsedOptions.ProjectReferences
}

func (p *ParsedCommandLine) GetConfigFileParsingDiagnostics() []*ast.Diagnostic {
	if p.ConfigFile != nil {
		// todo: !!! should be ConfigFile.ParseDiagnostics, check if they are the same
		return slices.Concat(p.ConfigFile.Diagnostics(), p.Errors)
	}
	return p.Errors
}
