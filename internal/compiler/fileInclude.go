package compiler

import (
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type fileIncludeKind int

const (
	// References from file
	fileIncludeKindImport = iota
	fileIncludeKindReferenceFile
	fileIncludeKindTypeReferenceDirective
	fileIncludeKindLibReferenceDirective

	fileIncludeKindRootFile
	fileIncludeKindSourceFromProjectReference
	fileIncludeKindOutputFromProjectReference
	fileIncludeKindLibFile
	fileIncludeKindAutomaticTypeDirectiveFile
)

type fileIncludeReason struct {
	kind fileIncludeKind
	data any

	diag     *ast.Diagnostic
	diagOnce sync.Once
}

type referencedFileData struct {
	file         tspath.Path
	index        int
	synthetic    *ast.Node
	location     *referenceFileLocation
	locationOnce sync.Once
}

type referenceFileLocation struct {
	file      *ast.SourceFile
	node      *ast.Node
	ref       *ast.FileReference
	packageId module.PackageId
}

func (r *referenceFileLocation) text() string {
	if r.node != nil {
		return r.node.Text()
	} else {
		return r.file.Text()[r.ref.Pos():r.ref.End()]
	}
}

type automaticTypeDirectiveFileData struct {
	typeReference string
	packageId     module.PackageId
}

func (r *fileIncludeReason) asIndex() int {
	return r.data.(int)
}

func (r *fileIncludeReason) asLibFileIndex() (int, bool) {
	index, ok := r.data.(int)
	return index, ok
}

func (r *fileIncludeReason) isReferencedFile() bool {
	return r.kind <= fileIncludeKindLibReferenceDirective
}

func (r *fileIncludeReason) asReferencedFileData() *referencedFileData {
	return r.data.(*referencedFileData)
}

func (r *fileIncludeReason) asAutomaticTypeDirectiveFileData() *automaticTypeDirectiveFileData {
	return r.data.(*automaticTypeDirectiveFileData)
}

func (r *fileIncludeReason) toDiagnostic(program *Program, toFileName func(string) string) *ast.Diagnostic {
	r.diagOnce.Do(func() {
		if r.isReferencedFile() {
			r.diag = r.toReferenceDiagnostic(program, toFileName)
			return
		}
		switch r.kind {
		case fileIncludeKindRootFile:
			if program.opts.Config.ConfigFile != nil {
				config := program.opts.Config
				fileName := tspath.GetNormalizedAbsolutePath(config.FileNames()[r.asIndex()], program.GetCurrentDirectory())
				if matchedFileSpec := config.GetMatchedFileSpec(fileName); matchedFileSpec != "" {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Part_of_files_list_in_tsconfig_json, matchedFileSpec, toFileName(fileName))
				} else if matchedIncludeSpec, isDefaultIncludeSpec := config.GetMatchedIncludeSpec(fileName); matchedIncludeSpec != "" {
					if isDefaultIncludeSpec {
						r.diag = ast.NewCompilerDiagnostic(diagnostics.Matched_by_default_include_pattern_Asterisk_Asterisk_Slash_Asterisk)
					} else {
						r.diag = ast.NewCompilerDiagnostic(diagnostics.Matched_by_include_pattern_0_in_1, matchedIncludeSpec, toFileName(config.ConfigName()))
					}
				} else {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Root_file_specified_for_compilation)
				}
			} else {
				r.diag = ast.NewCompilerDiagnostic(diagnostics.Root_file_specified_for_compilation)
			}
		case fileIncludeKindSourceFromProjectReference:
		case fileIncludeKindOutputFromProjectReference:
			diag := core.IfElse(
				r.kind == fileIncludeKindOutputFromProjectReference,
				diagnostics.Output_from_referenced_project_0_included_because_module_is_specified_as_none,
				diagnostics.Source_from_referenced_project_0_included_because_module_is_specified_as_none,
			)
			referencedResolvedRef := program.projectReferenceFileMapper.getResolvedProjectReferences()[r.asIndex()]
			r.diag = ast.NewCompilerDiagnostic(diag, toFileName(referencedResolvedRef.ConfigName()))
		case fileIncludeKindAutomaticTypeDirectiveFile:
			data := r.asAutomaticTypeDirectiveFileData()
			if program.Options().Types != nil {
				if data.packageId.Name != "" {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Entry_point_of_type_library_0_specified_in_compilerOptions_with_packageId_1, data.typeReference, data.packageId.String())
				} else {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Entry_point_of_type_library_0_specified_in_compilerOptions, data.typeReference)
				}
			} else {
				if data.packageId.Name != "" {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Entry_point_for_implicit_type_library_0_with_packageId_1, data.typeReference, data.packageId.String())
				} else {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Entry_point_for_implicit_type_library_0, data.typeReference)
				}
			}
		case fileIncludeKindLibFile:
			index, ok := r.asLibFileIndex()
			if ok {
				r.diag = ast.NewCompilerDiagnostic(diagnostics.Library_0_specified_in_compilerOptions, program.Options().Lib[index])
			} else {
				target := program.Options().GetEmitScriptTarget().String()
				if target != "" {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Default_library_for_target_0, target)
				} else {
					r.diag = ast.NewCompilerDiagnostic(diagnostics.Default_library)
				}
			}
		default:
			panic(fmt.Sprintf("unknown reason: %v", r.kind))
		}
	})
	return r.diag
}

func (r *fileIncludeReason) toReferenceDiagnostic(program *Program, toFileName func(string) string) *ast.Diagnostic {
	referenceLocation := r.getReferencedLocation(program)
	referenceText := referenceLocation.text()
	switch r.kind {
	case fileIncludeKindImport:
		if specifier, ok := program.importHelpersImportSpecifiers[referenceLocation.file.Path()]; ok && specifier == referenceLocation.node {
			if referenceLocation.packageId.Name != "" {
				return ast.NewCompilerDiagnostic(diagnostics.Imported_via_0_from_file_1_with_packageId_2_to_import_importHelpers_as_specified_in_compilerOptions, referenceText, toFileName(referenceLocation.file.FileName()), referenceLocation.packageId.String())
			} else {
				return ast.NewCompilerDiagnostic(diagnostics.Imported_via_0_from_file_1_to_import_importHelpers_as_specified_in_compilerOptions, referenceText, toFileName(referenceLocation.file.FileName()))
			}
		} else if jsxSpecifier, ok := program.jsxRuntimeImportSpecifiers[referenceLocation.file.Path()]; ok && jsxSpecifier.specifier == referenceLocation.node {
			if referenceLocation.packageId.Name != "" {
				return ast.NewCompilerDiagnostic(diagnostics.Imported_via_0_from_file_1_with_packageId_2_to_import_jsx_and_jsxs_factory_functions, referenceText, toFileName(referenceLocation.file.FileName()), referenceLocation.packageId.String())
			} else {
				return ast.NewCompilerDiagnostic(diagnostics.Imported_via_0_from_file_1_to_import_jsx_and_jsxs_factory_functions, referenceText, toFileName(referenceLocation.file.FileName()))
			}
		} else {
			if referenceLocation.packageId.Name != "" {
				return ast.NewCompilerDiagnostic(diagnostics.Imported_via_0_from_file_1_with_packageId_2, referenceText, toFileName(referenceLocation.file.FileName()), referenceLocation.packageId.String())
			} else {
				return ast.NewCompilerDiagnostic(diagnostics.Imported_via_0_from_file_1, referenceText, toFileName(referenceLocation.file.FileName()))
			}
		}
	case fileIncludeKindReferenceFile:
		return ast.NewCompilerDiagnostic(diagnostics.Referenced_via_0_from_file_1, referenceText, toFileName(referenceLocation.file.FileName()))
	case fileIncludeKindTypeReferenceDirective:
		if referenceLocation.packageId.Name != "" {
			return ast.NewCompilerDiagnostic(diagnostics.Type_library_referenced_via_0_from_file_1_with_packageId_2, referenceText, toFileName(referenceLocation.file.FileName()), referenceLocation.packageId.String())
		} else {
			return ast.NewCompilerDiagnostic(diagnostics.Type_library_referenced_via_0_from_file_1, referenceText, toFileName(referenceLocation.file.FileName()))
		}
	case fileIncludeKindLibReferenceDirective:
		return ast.NewCompilerDiagnostic(diagnostics.Library_referenced_via_0_from_file_1, referenceText, toFileName(referenceLocation.file.FileName()))

	default:
		panic(fmt.Sprintf("unknown reason: %v", r.kind))
	}
}

func (r *fileIncludeReason) getReferencedLocation(program *Program) *referenceFileLocation {
	ref := r.asReferencedFileData()
	ref.locationOnce.Do(func() {
		file := program.GetSourceFileByPath(ref.file)
		switch r.kind {
		case fileIncludeKindImport:
			var specifier *ast.Node
			if ref.synthetic != nil {
				specifier = ref.synthetic
			} else if ref.index < len(file.Imports()) {
				specifier = file.Imports()[ref.index]
			} else {
				augIndex := len(file.Imports())
				for _, imp := range file.ModuleAugmentations {
					if imp.Kind == ast.KindStringLiteral {
						if augIndex == ref.index {
							specifier = imp
							break
						}
						augIndex++
					}
				}
			}
			resolution := program.GetResolvedModuleFromModuleSpecifier(file, specifier)
			ref.location = &referenceFileLocation{
				file:      file,
				node:      specifier,
				packageId: resolution.PackageId,
			}
		case fileIncludeKindReferenceFile:
			ref.location = &referenceFileLocation{
				file: file,
				ref:  file.ReferencedFiles[ref.index],
			}
		case fileIncludeKindTypeReferenceDirective:
			ref.location = &referenceFileLocation{
				file: file,
				ref:  file.TypeReferenceDirectives[ref.index],
			}
		case fileIncludeKindLibReferenceDirective:
			ref.location = &referenceFileLocation{
				file: file,
				ref:  file.LibReferenceDirectives[ref.index],
			}
		default:
			panic(fmt.Sprintf("unknown reason: %v", r.kind))
		}
	})
	return ref.location
}
