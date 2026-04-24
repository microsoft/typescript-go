package tsoptions_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

func TestCompilerOptionsDeclaration(t *testing.T) {
	t.Parallel()

	decls := make(map[string]*tsoptions.CommandLineOption)

	for _, decl := range tsoptions.OptionsDeclarations {
		decls[strings.ToLower(decl.Name)] = decl
	}

	internalOptions := []string{
		"allowNonTsExtensions",
		"build",
		"configFilePath",
		"noDtsResolution",
		"noEmitForJsFiles",
		"pathsBasePath",
		"suppressOutputPathCheck",
		"build",
	}

	internalOptionsMap := make(map[string]string)
	for _, opt := range internalOptions {
		internalOptionsMap[strings.ToLower(opt)] = opt
	}

	compilerOptionsType := reflect.TypeFor[core.CompilerOptions]()
	for field := range compilerOptionsType.Fields() {
		if !field.IsExported() {
			continue
		}

		lowerName := strings.ToLower(field.Name)

		decl := decls[lowerName]
		if decl == nil {
			if name, ok := internalOptionsMap[lowerName]; ok {
				checkCompilerOptionJsonTagName(t, field, name)
				continue
			}
			t.Errorf("CompilerOptions.%s has no options declaration", field.Name)
			continue
		}
		delete(decls, lowerName)

		checkCompilerOptionJsonTagName(t, field, decl.Name)
	}

	skippedOptions := []string{
		"plugins",
	}

	for _, opt := range skippedOptions {
		delete(decls, strings.ToLower(opt))
	}

	for _, decl := range decls {
		t.Errorf("Option declaration %s is not present in CompilerOptions", decl.Name)
	}
}

func checkCompilerOptionJsonTagName(t *testing.T, field reflect.StructField, name string) {
	t.Helper()
	want := name + ",omitzero"
	got := field.Tag.Get("json")
	if got != want {
		t.Errorf("Field %s has json tag %s, but the option declaration has name %s", field.Name, got, name)
	}
}

func TestStrictArrayVarianceOptionDeclarationShape(t *testing.T) {
	t.Parallel()

	var opt *tsoptions.CommandLineOption
	for _, o := range tsoptions.OptionsDeclarations {
		if o.Name == "strictArrayVariance" {
			opt = o
			break
		}
	}
	if opt == nil {
		t.Fatal("strictArrayVariance missing from OptionsDeclarations")
	}
	if opt.Kind != tsoptions.CommandLineOptionTypeBoolean {
		t.Errorf("strictArrayVariance kind = %v, want boolean", opt.Kind)
	}
	if opt.IsCommandLineOnly || opt.IsTSConfigOnly {
		t.Errorf("strictArrayVariance should be a normal compiler option (CLI + tsconfig), got IsCommandLineOnly=%v IsTSConfigOnly=%v", opt.IsCommandLineOnly, opt.IsTSConfigOnly)
	}
	if !opt.AffectsSemanticDiagnostics || !opt.AffectsBuildInfo {
		t.Errorf("strictArrayVariance should affect semantic diagnostics and build info, got semantic=%v buildInfo=%v", opt.AffectsSemanticDiagnostics, opt.AffectsBuildInfo)
	}
	// strictFlag must stay false in declscompiler.go so `--strict` does not enable this option.
}
