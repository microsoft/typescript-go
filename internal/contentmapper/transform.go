package contentmapper

import (
	"errors"
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
)

// TransformAndParse runs the given content mapper's transform for a foreign file and
// parses the resulting TypeScript, preserving the original file name and retaining the untransformed text
// on the source file. The mapper is supplied by the caller (which also owns the failure accounting) so it
// is neither re-resolved nor substituted here. It returns an error if the transform fails or the mapper
// produces invalid position mappings (a *spanmap.MappingError); the caller decides how to report the failure
// and what placeholder file to substitute. It is the shared implementation behind
// CompilerHost.GetContentMappedSourceFile.
func TransformAndParse(
	parseOptions ast.SourceFileParseOptions,
	content string,
	mapper *Mapper,
	compilerOptions *core.CompilerOptions,
	host Host,
) (*ast.SourceFile, error) {
	if host == nil {
		panic(fmt.Sprintf("content mapper host is required to load content-mapped file %q", parseOptions.FileName))
	}
	result, err := host.Transform(mapper, Request{
		FileName:        parseOptions.FileName,
		Content:         content,
		ConfigFileName:  compilerOptions.ConfigFilePath,
		CompilerOptions: compilerOptions,
	})
	if err != nil {
		return nil, err
	}
	if result.Mappings == nil {
		return nil, errors.New("content mapper host returned a successful transform without position mappings")
	}
	if problem := result.Mappings.Validate(result.Text, content); problem != nil {
		return nil, problem
	}
	sourceFile := parser.ParseSourceFile(parseOptions, result.Text, result.ScriptKind)
	sourceFile.SetOriginalText(content)
	sourceFile.SetSpanMap(result.Mappings)
	sourceFile.SetContentMapper(mapper.Identity())
	if len(result.Diagnostics) > 0 {
		// The runner produces diagnostics without a source file (it doesn't have one yet); associate
		// them with the file now so they are reported against it.
		for _, diagnostic := range result.Diagnostics {
			diagnostic.SetFile(sourceFile)
		}
		sourceFile.SetDiagnostics(append(sourceFile.Diagnostics(), result.Diagnostics...))
	}
	return sourceFile, nil
}
