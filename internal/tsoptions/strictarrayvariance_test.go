package tsoptions_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tsoptions/tsoptionstest"
	"gotest.tools/v3/assert"
)

func parseTsconfigAtWork(t *testing.T, host *tsoptionstest.VfsParseConfigHost, jsonText string) *tsoptions.ParsedCommandLine {
	t.Helper()
	const basePath = "/work"
	configFileName := "/work/tsconfig.json"
	path := tspath.ToPath(configFileName, basePath, true)
	parsed, parseErrs := tsoptions.ParseConfigFileTextToJson(configFileName, path, jsonText)
	assert.Assert(t, len(parseErrs) == 0, "ParseConfigFileTextToJson: %v", parseErrs)
	return tsoptions.ParseJsonConfigFileContent(parsed, host, basePath, nil, configFileName, nil, nil, nil)
}

func parseJsconfigAtWork(t *testing.T, host *tsoptionstest.VfsParseConfigHost, jsonText string) *tsoptions.ParsedCommandLine {
	t.Helper()
	const basePath = "/work"
	configFileName := "/work/jsconfig.json"
	path := tspath.ToPath(configFileName, basePath, true)
	parsed, parseErrs := tsoptions.ParseConfigFileTextToJson(configFileName, path, jsonText)
	assert.Assert(t, len(parseErrs) == 0, "ParseConfigFileTextToJson: %v", parseErrs)
	return tsoptions.ParseJsonConfigFileContent(parsed, host, basePath, nil, configFileName, nil, nil, nil)
}

func TestStrictArrayVarianceCommandLine(t *testing.T) {
	t.Parallel()

	host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
		"/proj/0.ts": "",
	}, "/proj", true)

	t.Run("explicit false", func(t *testing.T) {
		t.Parallel()
		res := tsoptions.ParseCommandLine([]string{"--strictArrayVariance", "false", "0.ts"}, host)
		assert.Assert(t, len(res.Errors) == 0, "got errors: %v", res.Errors)
		assert.Equal(t, core.TSFalse, res.CompilerOptions().StrictArrayVariance)
	})

	t.Run("explicit true", func(t *testing.T) {
		t.Parallel()
		res := tsoptions.ParseCommandLine([]string{"--strictArrayVariance", "true", "0.ts"}, host)
		assert.Assert(t, len(res.Errors) == 0, "got errors: %v", res.Errors)
		assert.Equal(t, core.TSTrue, res.CompilerOptions().StrictArrayVariance)
	})

	t.Run("implicit true", func(t *testing.T) {
		t.Parallel()
		res := tsoptions.ParseCommandLine([]string{"--strictArrayVariance"}, host)
		assert.Assert(t, len(res.Errors) == 0, "got errors: %v", res.Errors)
		assert.Equal(t, core.TSTrue, res.CompilerOptions().StrictArrayVariance)
	})

	t.Run("strict does not enable strictArrayVariance", func(t *testing.T) {
		t.Parallel()
		res := tsoptions.ParseCommandLine([]string{"--strict", "0.ts"}, host)
		assert.Assert(t, len(res.Errors) == 0, "got errors: %v", res.Errors)
		assert.Equal(t, core.TSTrue, res.CompilerOptions().Strict)
		assert.Equal(t, core.TSUnknown, res.CompilerOptions().StrictArrayVariance)
	})

	t.Run("non true false literal sets true and does not consume next token as flag value", func(t *testing.T) {
		t.Parallel()
		h := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/proj/0.ts":  "",
			"/proj/maybe": "",
		}, "/proj", true)
		res := tsoptions.ParseCommandLine([]string{"--strictArrayVariance", "maybe", "0.ts"}, h)
		assert.Assert(t, len(res.Errors) == 0, "got errors: %v", res.Errors)
		assert.Equal(t, core.TSTrue, res.CompilerOptions().StrictArrayVariance)
		assert.DeepEqual(t, []string{"maybe", "0.ts"}, res.FileNames())
	})
}

func TestStrictArrayVarianceTsConfigJson(t *testing.T) {
	t.Parallel()

	t.Run("compilerOptions.strictArrayVariance true", func(t *testing.T) {
		t.Parallel()
		jsonText := `{
  "compilerOptions": {
    "strictArrayVariance": true
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": jsonText,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseTsconfigAtWork(t, host, jsonText)
		assert.Assert(t, len(pc.Errors) == 0, "got errors: %v", pc.Errors)
		assert.Equal(t, core.TSTrue, pc.CompilerOptions().StrictArrayVariance)
	})

	t.Run("jsconfig compilerOptions.strictArrayVariance true", func(t *testing.T) {
		t.Parallel()
		jsonText := `{
  "compilerOptions": {
    "strictArrayVariance": true
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/jsconfig.json": jsonText,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseJsconfigAtWork(t, host, jsonText)
		assert.Assert(t, len(pc.Errors) == 0, "got errors: %v", pc.Errors)
		assert.Equal(t, core.TSTrue, pc.CompilerOptions().StrictArrayVariance)
	})

	t.Run("compilerOptions.strictArrayVariance false", func(t *testing.T) {
		t.Parallel()
		jsonText := `{
  "compilerOptions": {
    "strictArrayVariance": false
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": jsonText,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseTsconfigAtWork(t, host, jsonText)
		assert.Assert(t, len(pc.Errors) == 0, "got errors: %v", pc.Errors)
		assert.Equal(t, core.TSFalse, pc.CompilerOptions().StrictArrayVariance)
	})

	t.Run("strict true does not imply strictArrayVariance", func(t *testing.T) {
		t.Parallel()
		jsonText := `{
  "compilerOptions": {
    "strict": true
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": jsonText,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseTsconfigAtWork(t, host, jsonText)
		assert.Assert(t, len(pc.Errors) == 0, "got errors: %v", pc.Errors)
		assert.Equal(t, core.TSTrue, pc.CompilerOptions().Strict)
		assert.Equal(t, core.TSUnknown, pc.CompilerOptions().StrictArrayVariance)
	})

	t.Run("extends inherits strictArrayVariance from base", func(t *testing.T) {
		t.Parallel()
		baseJSON := `{"compilerOptions":{"strictArrayVariance":true}}`
		derivedJSON := `{
  "extends": "./base.json",
  "compilerOptions": {},
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/base.json":     baseJSON,
			"/work/tsconfig.json": derivedJSON,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseTsconfigAtWork(t, host, derivedJSON)
		assert.Assert(t, len(pc.Errors) == 0, "got errors: %v", pc.Errors)
		assert.Equal(t, core.TSTrue, pc.CompilerOptions().StrictArrayVariance)
	})

	t.Run("extends child strictArrayVariance overrides base", func(t *testing.T) {
		t.Parallel()
		baseJSON := `{"compilerOptions":{"strictArrayVariance":true}}`
		derivedJSON := `{
  "extends": "./base.json",
  "compilerOptions": {
    "strictArrayVariance": false
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/base.json":     baseJSON,
			"/work/tsconfig.json": derivedJSON,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseTsconfigAtWork(t, host, derivedJSON)
		assert.Assert(t, len(pc.Errors) == 0, "got errors: %v", pc.Errors)
		assert.Equal(t, core.TSFalse, pc.CompilerOptions().StrictArrayVariance)
	})

	t.Run("invalid strictArrayVariance type yields diagnostics", func(t *testing.T) {
		t.Parallel()
		jsonText := `{
  "compilerOptions": {
    "strictArrayVariance": "not-a-boolean"
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": jsonText,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseTsconfigAtWork(t, host, jsonText)
		assert.Assert(t, len(pc.Errors) > 0, "expected diagnostics for invalid strictArrayVariance, got none")
		assert.Equal(t, core.TSUnknown, pc.CompilerOptions().StrictArrayVariance)
	})

	t.Run("numeric strictArrayVariance is rejected as non boolean", func(t *testing.T) {
		t.Parallel()
		jsonText := `{
  "compilerOptions": {
    "strictArrayVariance": 1
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": jsonText,
			"/work/a.ts":          "",
		}, "/work", true)
		pc := parseTsconfigAtWork(t, host, jsonText)
		assert.Assert(t, len(pc.Errors) > 0, "expected diagnostics for numeric strictArrayVariance, got none")
		assert.Equal(t, core.TSUnknown, pc.CompilerOptions().StrictArrayVariance)
	})
}

func TestStrictArrayVarianceCommandLineOverridesTsconfig(t *testing.T) {
	t.Parallel()

	t.Run("CLI false overrides tsconfig true", func(t *testing.T) {
		t.Parallel()
		tsconfigText := `{
  "compilerOptions": {
    "strictArrayVariance": true
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": tsconfigText,
			"/work/a.ts":          "",
		}, "/work", true)

		cmdLine := tsoptions.ParseCommandLine([]string{"--project", "/work", "--strictArrayVariance", "false"}, host)
		assert.Assert(t, len(cmdLine.Errors) == 0, "ParseCommandLine errors: %v", cmdLine.Errors)

		rawMap, ok := cmdLine.Raw.(*collections.OrderedMap[string, any])
		assert.Assert(t, ok, "cmdLine.Raw should be *OrderedMap[string, any]")
		wrappedRaw := &collections.OrderedMap[string, any]{}
		wrappedRaw.Set("compilerOptions", rawMap)

		parsedConfig, diag := tsoptions.GetParsedCommandLineOfConfigFile(
			"/work/tsconfig.json",
			cmdLine.CompilerOptions(),
			wrappedRaw,
			host,
			nil,
		)
		assert.Assert(t, len(diag) == 0, "GetParsedCommandLineOfConfigFile: %v", diag)
		assert.Assert(t, len(parsedConfig.Errors) == 0, "parsed config errors: %v", parsedConfig.Errors)
		assert.Equal(t, core.TSFalse, parsedConfig.CompilerOptions().StrictArrayVariance)
	})

	t.Run("CLI true overrides tsconfig false", func(t *testing.T) {
		t.Parallel()
		tsconfigText := `{
  "compilerOptions": {
    "strictArrayVariance": false
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": tsconfigText,
			"/work/a.ts":          "",
		}, "/work", true)

		cmdLine := tsoptions.ParseCommandLine([]string{"--project", "/work", "--strictArrayVariance", "true"}, host)
		assert.Assert(t, len(cmdLine.Errors) == 0, "ParseCommandLine errors: %v", cmdLine.Errors)

		rawMap, ok := cmdLine.Raw.(*collections.OrderedMap[string, any])
		assert.Assert(t, ok, "cmdLine.Raw should be *OrderedMap[string, any]")
		wrappedRaw := &collections.OrderedMap[string, any]{}
		wrappedRaw.Set("compilerOptions", rawMap)

		parsedConfig, diag := tsoptions.GetParsedCommandLineOfConfigFile(
			"/work/tsconfig.json",
			cmdLine.CompilerOptions(),
			wrappedRaw,
			host,
			nil,
		)
		assert.Assert(t, len(diag) == 0, "GetParsedCommandLineOfConfigFile: %v", diag)
		assert.Assert(t, len(parsedConfig.Errors) == 0, "parsed config errors: %v", parsedConfig.Errors)
		assert.Equal(t, core.TSTrue, parsedConfig.CompilerOptions().StrictArrayVariance)
	})

	t.Run("project without CLI flag keeps tsconfig strictArrayVariance", func(t *testing.T) {
		t.Parallel()
		tsconfigText := `{
  "compilerOptions": {
    "strictArrayVariance": true
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/tsconfig.json": tsconfigText,
			"/work/a.ts":          "",
		}, "/work", true)

		cmdLine := tsoptions.ParseCommandLine([]string{"--project", "/work"}, host)
		assert.Assert(t, len(cmdLine.Errors) == 0, "ParseCommandLine errors: %v", cmdLine.Errors)

		rawMap, ok := cmdLine.Raw.(*collections.OrderedMap[string, any])
		assert.Assert(t, ok, "cmdLine.Raw should be *OrderedMap[string, any]")
		wrappedRaw := &collections.OrderedMap[string, any]{}
		wrappedRaw.Set("compilerOptions", rawMap)

		parsedConfig, diag := tsoptions.GetParsedCommandLineOfConfigFile(
			"/work/tsconfig.json",
			cmdLine.CompilerOptions(),
			wrappedRaw,
			host,
			nil,
		)
		assert.Assert(t, len(diag) == 0, "GetParsedCommandLineOfConfigFile: %v", diag)
		assert.Assert(t, len(parsedConfig.Errors) == 0, "parsed config errors: %v", parsedConfig.Errors)
		assert.Equal(t, core.TSTrue, parsedConfig.CompilerOptions().StrictArrayVariance)
	})

	t.Run("project explicit jsconfig path loads strictArrayVariance", func(t *testing.T) {
		t.Parallel()
		jsconfigText := `{
  "compilerOptions": {
    "strictArrayVariance": true
  },
  "files": ["a.ts"]
}`
		host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
			"/work/jsconfig.json": jsconfigText,
			"/work/a.ts":          "",
		}, "/work", true)

		cmdLine := tsoptions.ParseCommandLine([]string{"--project", "/work/jsconfig.json"}, host)
		assert.Assert(t, len(cmdLine.Errors) == 0, "ParseCommandLine errors: %v", cmdLine.Errors)

		rawMap, ok := cmdLine.Raw.(*collections.OrderedMap[string, any])
		assert.Assert(t, ok, "cmdLine.Raw should be *OrderedMap[string, any]")
		wrappedRaw := &collections.OrderedMap[string, any]{}
		wrappedRaw.Set("compilerOptions", rawMap)

		parsedConfig, diag := tsoptions.GetParsedCommandLineOfConfigFile(
			"/work/jsconfig.json",
			cmdLine.CompilerOptions(),
			wrappedRaw,
			host,
			nil,
		)
		assert.Assert(t, len(diag) == 0, "GetParsedCommandLineOfConfigFile: %v", diag)
		assert.Assert(t, len(parsedConfig.Errors) == 0, "parsed config errors: %v", parsedConfig.Errors)
		assert.Equal(t, core.TSTrue, parsedConfig.CompilerOptions().StrictArrayVariance)
	})
}

func TestStrictArrayVarianceParseBuildCommandLine(t *testing.T) {
	t.Parallel()
	host := tsoptionstest.NewVFSParseConfigHost(map[string]string{}, "/proj", true)

	// `ParseBuildCommandLine` only accepts build/watch options plus an allow-listed set of
	// compiler flags. `strictArrayVariance` is type-checking-only and must be set via
	// tsconfig (or non-build `tsc`), matching TS "may not be used with '--build'" behavior.
	t.Run("strictArrayVariance is rejected when parsing as build mode", func(t *testing.T) {
		t.Parallel()
		res := tsoptions.ParseBuildCommandLine([]string{"--strictArrayVariance", "true", "packages/a"}, host)
		assert.Assert(t, len(res.Errors) == 1, "errors: %v", res.Errors)
		assert.Equal(t, diagnostics.Compiler_option_0_may_not_be_used_with_build.Code(), res.Errors[0].Code())
		assert.Equal(t, core.TSUnknown, res.CompilerOptions.StrictArrayVariance)
	})
}

func TestStrictArrayVarianceFlipAffectsIncrementalOptionComparisons(t *testing.T) {
	t.Parallel()
	on := &core.CompilerOptions{StrictArrayVariance: core.TSTrue}
	off := &core.CompilerOptions{StrictArrayVariance: core.TSFalse}
	assert.Assert(t, tsoptions.CompilerOptionsAffectSemanticDiagnostics(on, off),
		"incremental reuse of semantic state must not cross strictArrayVariance=false vs true")
	assert.Assert(t, !tsoptions.CompilerOptionsAffectSemanticDiagnostics(on, on))
}
