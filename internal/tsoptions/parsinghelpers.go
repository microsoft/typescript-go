package tsoptions

import (
	"reflect"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

func parseTristate(value any) core.Tristate {
	if value == nil {
		return core.TSUnknown
	}
	if value == true {
		return core.TSTrue
	} else {
		return core.TSFalse
	}
}

func parseStringArray(value any) []string {
	if arr, ok := value.([]any); ok {
		var result []string
		for _, v := range arr {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return nil
}

func parseRawStringArray(value any) []string {
	if arr, ok := value.([]string); ok {
		return arr
	}
	return []string{}
}

func parseStringMap(value any) map[string][]string {
	if m, ok := value.(map[string]any); ok {
		result := make(map[string][]string)
		for k, v := range m {
			result[k] = parseStringArray(v)
		}
		return result
	}
	return nil
}

func parseString(value any) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func parseProjectReference(json any) []core.ProjectReference {
	var result []core.ProjectReference
	if v, ok := json.(map[string]any); ok {
		var reference core.ProjectReference
		if v, ok := v["path"]; ok {
			reference.Path = v.(string)
		}
		if v, ok := v["originalPath"]; ok {
			reference.OriginalPath = v.(string)
		}
		if v, ok := v["circular"]; ok {
			reference.Circular = v.(bool)
		}
		result = append(result, reference)
	}
	return result
}

func parseJsonToStringKey(json any) map[string]any {
	result := make(map[string]any)
	if m, ok := json.(map[string]any); ok {
		if v, ok := m["include"]; ok {
			if arr, ok := v.([]string); ok {
				if len(arr) == 0 {
					result["include"] = []any{}
				}
			} else {
				result["include"] = v
			}
		}
		if v, ok := m["exclude"]; ok {
			if arr, ok := v.([]string); ok {
				if len(arr) == 0 {
					result["exclude"] = []any{}
				}
			} else {
				result["exclude"] = v
			}
		}
		if v, ok := m["files"]; ok {
			if arr, ok := v.([]string); ok {
				if len(arr) == 0 {
					result["files"] = []any{}
				}
			} else {
				result["files"] = v
			}
		}
		if v, ok := m["references"]; ok {
			if arr, ok := v.([]string); ok {
				if len(arr) == 0 {
					result["references"] = []any{}
				}
			} else {
				result["references"] = v
			}
		}
		if v, ok := m["extends"]; ok {
			if arr, ok := v.([]string); ok {
				if len(arr) == 0 {
					result["extends"] = []any{}
				}
			} else if str, ok := v.(string); ok {
				result["extends"] = []any{str}
			} else {
				result["extends"] = v
			}
		}
		if v, ok := m["compilerOptions"]; ok {
			result["compilerOptions"] = v
		}
	}
	return result
}

func parseCompilerOptions(key string, value any, allOptions *core.CompilerOptions) []*ast.Diagnostic {
	if allOptions == nil {
		return nil
	}
	switch key {
	case "allowJs":
		allOptions.AllowJs = parseTristate(value)
	case "allowSyntheticDefaultImports":
		allOptions.AllowSyntheticDefaultImports = parseTristate(value)
	case "allowNonTsExtensions":
		allOptions.AllowNonTsExtensions = parseTristate(value)
	case "allowUmdGlobalAccess":
		allOptions.AllowUmdGlobalAccess = parseTristate(value)
	case "allowUnreachableCode":
		allOptions.AllowUnreachableCode = parseTristate(value)
	case "allowUnusedLabels":
		allOptions.AllowUnusedLabels = parseTristate(value)
	case "allowArbitraryExtensions":
		allOptions.AllowArbitraryExtensions = parseTristate(value)
	case "alwaysStrict":
		allOptions.AlwaysStrict = parseTristate(value)
	case "assumeChangesOnlyAffectDirectDependencies":
		allOptions.AssumeChangesOnlyAffectDirectDependencies = parseTristate(value)
	case "baseUrl":
		allOptions.BaseUrl = parseString(value)
	case "build":
		allOptions.Build = parseTristate(value)
	case "checkJs":
		allOptions.CheckJs = parseTristate(value)
	case "customConditions":
		allOptions.CustomConditions = parseStringArray(value)
	case "composite":
		allOptions.Composite = parseTristate(value)
	case "declarationDir":
		allOptions.DeclarationDir = parseString(value)
	case "diagnostics":
		allOptions.Diagnostics = parseTristate(value)
	case "disableSizeLimit":
		allOptions.DisableSizeLimit = parseTristate(value)
	case "disableSourceOfProjectReferenceRedirect":
		allOptions.DisableSourceOfProjectReferenceRedirect = parseTristate(value)
	case "disableSolutionSearching":
		allOptions.DisableSolutionSearching = parseTristate(value)
	case "disableReferencedProjectLoad":
		allOptions.DisableReferencedProjectLoad = parseTristate(value)
	case "declarationMap":
		allOptions.DeclarationMap = parseTristate(value)
	case "declaration":
		allOptions.Declaration = parseTristate(value)
	case "extendedDiagnostics":
		allOptions.ExtendedDiagnostics = parseTristate(value)
	case "emitDecoratorMetadata":
		allOptions.EmitDecoratorMetadata = parseTristate(value)
	case "esModuleInterop":
		allOptions.ESModuleInterop = parseTristate(value)
	case "exactOptionalPropertyTypes":
		allOptions.ExactOptionalPropertyTypes = parseTristate(value)
	case "explainFiles":
		allOptions.ExplainFiles = parseTristate(value)
	case "experimentalDecorators":
		allOptions.ExperimentalDecorators = parseTristate(value)
	case "forceConsistentCasingInFileNames":
		allOptions.ForceConsistentCasingInFileNames = parseTristate(value)
	case "generateCpuProfile":
		allOptions.GenerateCpuProfile = parseString(value)
	case "generateTrace":
		allOptions.GenerateTrace = parseString(value)
	case "isolatedModules":
		allOptions.IsolatedModules = parseTristate(value)
	case "ignoreDeprecations":
		allOptions.IgnoreDeprecations = parseString(value)
	case "importHelpers":
		allOptions.ImportHelpers = parseTristate(value)
	case "incremental":
		allOptions.Incremental = parseTristate(value)
	case "init":
		allOptions.Init = parseTristate(value)
	case "inlineSourceMap":
		allOptions.InlineSourceMap = parseTristate(value)
	case "inlineSources":
		allOptions.InlineSources = parseTristate(value)
	case "isolatedDeclarations":
		allOptions.IsolatedDeclarations = parseTristate(value)
	case "jsx":
		allOptions.Jsx = value.(core.JsxEmit)
	case "jsxFactory":
		allOptions.JsxFactory = parseString(value)
	case "jsxFragmentFactory":
		allOptions.JsxFragmentFactory = parseString(value)
	case "jsxImportSource":
		allOptions.JsxImportSource = parseString(value)
	case "keyofStringsOnly":
		allOptions.KeyofStringsOnly = parseTristate(value)
	case "lib":
		if _, ok := value.([]string); ok {
			allOptions.Lib = value.([]string)
		} else {
			allOptions.Lib = parseStringArray(value)
		}
	case "listEmittedFiles":
		allOptions.ListEmittedFiles = parseTristate(value)
	case "listFiles":
		allOptions.ListFiles = parseTristate(value)
	case "listFilesOnly":
		allOptions.ListFilesOnly = parseTristate(value)
	case "locale":
		allOptions.Locale = parseString(value)
	case "mapRoot":
		allOptions.MapRoot = parseString(value)
	case "module":
		allOptions.ModuleKind = value.(core.ModuleKind)
	case "moduleResolution":
		allOptions.ModuleResolution = value.(core.ModuleResolutionKind)
	case "moduleSuffixes":
		allOptions.ModuleSuffixes = parseStringArray(value)
	case "moduleDetection":
		allOptions.ModuleDetection = value.(core.ModuleDetectionKind)
	case "noCheck":
		allOptions.NoCheck = parseTristate(value)
	case "noFallthroughCasesInSwitch":
		allOptions.NoFallthroughCasesInSwitch = parseTristate(value)
	case "noEmitForJsFiles":
		allOptions.NoEmitForJsFiles = parseTristate(value)
	case "noImplicitAny":
		allOptions.NoImplicitAny = parseTristate(value)
	case "noImplicitThis":
		allOptions.NoImplicitThis = parseTristate(value)
	case "noPropertyAccessFromIndexSignature":
		allOptions.NoPropertyAccessFromIndexSignature = parseTristate(value)
	case "noUncheckedIndexedAccess":
		allOptions.NoUncheckedIndexedAccess = parseTristate(value)
	case "noEmitHelpers":
		allOptions.NoEmitHelpers = parseTristate(value)
	case "noEmitOnError":
		allOptions.NoEmitOnError = parseTristate(value)
	case "noImplicitReturns":
		allOptions.NoImplicitReturns = parseTristate(value)
	case "noUnusedLocals":
		allOptions.NoUnusedLocals = parseTristate(value)
	case "noUnusedParameters":
		allOptions.NoUnusedParameters = parseTristate(value)
	case "noImplicitOverride":
		allOptions.NoImplicitOverride = parseTristate(value)
	case "noUncheckedSideEffectImports":
		allOptions.NoUncheckedSideEffectImports = parseTristate(value)
	case "out":
		allOptions.Out = parseString(value)
	case "outFile":
		allOptions.OutFile = parseString(value)
	case "noResolve":
		allOptions.NoResolve = parseTristate(value)
	case "paths":
		allOptions.Paths = parseStringMap(value)
	case "preserveWatchOutput":
		allOptions.PreserveWatchOutput = parseTristate(value)
	case "preserveConstEnums":
		allOptions.PreserveConstEnums = parseTristate(value)
	case "preserveSymlinks":
		allOptions.PreserveSymlinks = parseTristate(value)
	case "project":
		allOptions.Project = parseString(value)
	case "pretty":
		allOptions.Pretty = parseTristate(value)
	case "resolveJsonModule":
		allOptions.ResolveJsonModule = parseTristate(value)
	case "resolvePackageJsonExports":
		allOptions.ResolvePackageJsonExports = parseTristate(value)
	case "resolvePackageJsonImports":
		allOptions.ResolvePackageJsonImports = parseTristate(value)
	case "reactNamespace":
		allOptions.ReactNamespace = parseString(value)
	case "rootDir":
		allOptions.RootDir = parseString(value)
	case "rootDirs":
		allOptions.RootDirs = parseStringArray(value)
	case "removeComments":
		allOptions.RemoveComments = parseTristate(value)
	case "strict":
		allOptions.Strict = parseTristate(value)
	case "strictBindCallApply":
		allOptions.StrictBindCallApply = parseTristate(value)
	case "strictFunctionTypes":
		allOptions.StrictFunctionTypes = parseTristate(value)
	case "strictNullChecks":
		allOptions.StrictNullChecks = parseTristate(value)
	case "strictPropertyInitialization":
		allOptions.StrictPropertyInitialization = parseTristate(value)
	case "skipDefaultLibCheck":
		allOptions.SkipDefaultLibCheck = parseTristate(value)
	case "sourceMap":
		allOptions.SourceMap = parseTristate(value)
	case "sourceRoot":
		allOptions.SourceRoot = parseString(value)
	case "stripInternal":
		allOptions.StripInternal = parseTristate(value)
	case "suppressOutputPathCheck":
		allOptions.SuppressOutputPathCheck = parseTristate(value)
	case "target":
		allOptions.Target = value.(core.ScriptTarget)
	case "traceResolution":
		allOptions.TraceResolution = parseTristate(value)
	case "tsBuildInfoFile":
		allOptions.TsBuildInfoFile = parseString(value)
	case "typeRoots":
		allOptions.TypeRoots = parseStringArray(value)
	case "tscBuild":
		allOptions.TscBuild = parseTristate(value)
	case "types":
		allOptions.Types = parseStringArray(value)
	case "useDefineForClassFields":
		allOptions.UseDefineForClassFields = parseTristate(value)
	case "useUnknownInCatchVariables":
		allOptions.UseUnknownInCatchVariables = parseTristate(value)
	case "verbatimModuleSyntax":
		allOptions.VerbatimModuleSyntax = parseTristate(value)
	case "version":
		allOptions.Version = parseTristate(value)
	case "maxNodeModuleJsDepth":
		allOptions.MaxNodeModuleJsDepth = parseTristate(value)
	case "skipLibCheck":
		allOptions.SkipLibCheck = parseTristate(value)
	case "noEmit":
		allOptions.NoEmit = parseTristate(value)
	case "showConfig":
		allOptions.ShowConfig = parseTristate(value)
	case "configFilePath":
		allOptions.ConfigFilePath = parseString(value)
	case "noDtsResolution":
		allOptions.NoDtsResolution = parseTristate(value)
	case "pathsBasePath":
		allOptions.PathsBasePath = parseString(value)
	case "outDir":
		allOptions.OutDir = parseString(value)
	case "newLine":
		allOptions.NewLine = value.(core.NewLineKind)
	case "watch":
		allOptions.Watch = parseTristate(value)
	}
	return nil
}

func mergeCompilerOptions(existingOptions, newOptions *core.CompilerOptions) {
	if existingOptions == nil {
		return
	}
	values := reflect.ValueOf(*newOptions)
	types := values.Type()
	for i := range values.NumField() {
		compareAndMergeCompilerOptions(types.Field(i).Name, existingOptions, newOptions)
	}
}

func compareAndMergeCompilerOptions(field string, existingOptions *core.CompilerOptions, newOptions *core.CompilerOptions) {
	switch field {
	case "AllowJs":
		if existingOptions.AllowJs != newOptions.AllowJs && newOptions.AllowJs == core.TSUnknown {
			newOptions.AllowJs = existingOptions.AllowJs
		}
	case "AllowSyntheticDefaultImports":
		if existingOptions.AllowSyntheticDefaultImports != newOptions.AllowSyntheticDefaultImports && newOptions.AllowSyntheticDefaultImports == core.TSUnknown {
			newOptions.AllowSyntheticDefaultImports = existingOptions.AllowSyntheticDefaultImports
		}
	case "AllowNonTsExtensions":
		if existingOptions.AllowNonTsExtensions != newOptions.AllowNonTsExtensions && newOptions.AllowNonTsExtensions == core.TSUnknown {
			newOptions.AllowNonTsExtensions = existingOptions.AllowNonTsExtensions
		}
	case "AllowUmdGlobalAccess":
		if existingOptions.AllowUmdGlobalAccess != newOptions.AllowUmdGlobalAccess && newOptions.AllowUmdGlobalAccess == core.TSUnknown {
			newOptions.AllowUmdGlobalAccess = existingOptions.AllowUmdGlobalAccess
		}
	case "AllowUnreachableCode":
		if existingOptions.AllowUnreachableCode != newOptions.AllowUnreachableCode && newOptions.AllowUnreachableCode == core.TSUnknown {
			newOptions.AllowUnreachableCode = existingOptions.AllowUnreachableCode
		}
	case "AllowUnusedLabels":
		if existingOptions.AllowUnusedLabels != newOptions.AllowUnusedLabels && newOptions.AllowUnusedLabels == core.TSUnknown {
			newOptions.AllowUnusedLabels = existingOptions.AllowUnusedLabels
		}
	case "AllowArbitraryExtensions":
		if existingOptions.AllowArbitraryExtensions != newOptions.AllowArbitraryExtensions && newOptions.AllowArbitraryExtensions == core.TSUnknown {
			newOptions.AllowArbitraryExtensions = existingOptions.AllowArbitraryExtensions
		}
	case "AlwaysStrict":
		if existingOptions.AlwaysStrict != newOptions.AlwaysStrict && newOptions.AlwaysStrict == core.TSUnknown {
			newOptions.AlwaysStrict = existingOptions.AlwaysStrict
		}
	case "AssumeChangesOnlyAffectDirectDependencies":
		if existingOptions.AssumeChangesOnlyAffectDirectDependencies != newOptions.AssumeChangesOnlyAffectDirectDependencies && newOptions.AssumeChangesOnlyAffectDirectDependencies == core.TSUnknown {
			newOptions.AssumeChangesOnlyAffectDirectDependencies = existingOptions.AssumeChangesOnlyAffectDirectDependencies
		}
	case "BaseUrl":
		if existingOptions.BaseUrl != newOptions.BaseUrl && newOptions.BaseUrl == "" {
			newOptions.BaseUrl = existingOptions.BaseUrl
		}
	case "Build":
		if existingOptions.Build != newOptions.Build && newOptions.Build == core.TSUnknown {
			newOptions.Build = existingOptions.Build
		}
	case "CheckJs":
		if existingOptions.CheckJs != newOptions.CheckJs && newOptions.CheckJs == core.TSUnknown {
			newOptions.CheckJs = existingOptions.CheckJs
		}
	case "CustomConditions":
		if !reflect.DeepEqual(existingOptions.CustomConditions, newOptions.CustomConditions) && len(newOptions.CustomConditions) == 0 {
			newOptions.CustomConditions = existingOptions.CustomConditions
		}
	case "Composite":
		if existingOptions.Composite != newOptions.Composite && newOptions.Composite == core.TSUnknown {
			newOptions.Composite = existingOptions.Composite
		}
	case "Declaration":
		if existingOptions.Declaration != newOptions.Declaration && newOptions.Declaration == core.TSUnknown {
			newOptions.Declaration = existingOptions.Declaration
		}
	case "DeclarationMap":
		if existingOptions.DeclarationMap != newOptions.DeclarationMap && newOptions.DeclarationMap == core.TSUnknown {
			newOptions.DeclarationMap = existingOptions.DeclarationMap
		}
	case "DeclarationDir":
		if existingOptions.DeclarationDir != newOptions.DeclarationDir && newOptions.DeclarationDir == "" {
			newOptions.DeclarationDir = existingOptions.DeclarationDir
		}
	case "Diagnostics":
		if existingOptions.Diagnostics != newOptions.Diagnostics && newOptions.Diagnostics == core.TSUnknown {
			newOptions.Diagnostics = existingOptions.Diagnostics
		}
	case "DisableSizeLimit":
		if existingOptions.DisableSizeLimit != newOptions.DisableSizeLimit && newOptions.DisableSizeLimit == core.TSUnknown {
			newOptions.DisableSizeLimit = existingOptions.DisableSizeLimit
		}
	case "DisableSourceOfProjectReferenceRedirect":
		if existingOptions.DisableSourceOfProjectReferenceRedirect != newOptions.DisableSourceOfProjectReferenceRedirect && newOptions.DisableSourceOfProjectReferenceRedirect == core.TSUnknown {
			newOptions.DisableSourceOfProjectReferenceRedirect = existingOptions.DisableSourceOfProjectReferenceRedirect
		}
	case "DisableSolutionSearching":
		if existingOptions.DisableSolutionSearching != newOptions.DisableSolutionSearching && newOptions.DisableSolutionSearching == core.TSUnknown {
			newOptions.DisableSolutionSearching = existingOptions.DisableSolutionSearching
		}
	case "DisableReferencedProjectLoad":
		if existingOptions.DisableReferencedProjectLoad != newOptions.DisableReferencedProjectLoad && newOptions.DisableReferencedProjectLoad == core.TSUnknown {
			newOptions.DisableReferencedProjectLoad = existingOptions.DisableReferencedProjectLoad
		}
	case "ExtendedDiagnostics":
		if existingOptions.ExtendedDiagnostics != newOptions.ExtendedDiagnostics && newOptions.ExtendedDiagnostics == core.TSUnknown {
			newOptions.ExtendedDiagnostics = existingOptions.ExtendedDiagnostics
		}
	case "EmitDecoratorMetadata":
		if existingOptions.EmitDecoratorMetadata != newOptions.EmitDecoratorMetadata && newOptions.EmitDecoratorMetadata == core.TSUnknown {
			newOptions.EmitDecoratorMetadata = existingOptions.EmitDecoratorMetadata
		}
	case "ESModuleInterop":
		if existingOptions.ESModuleInterop != newOptions.ESModuleInterop && newOptions.ESModuleInterop == core.TSUnknown {
			newOptions.ESModuleInterop = existingOptions.ESModuleInterop
		}
	case "ExactOptionalPropertyTypes":
		if existingOptions.ExactOptionalPropertyTypes != newOptions.ExactOptionalPropertyTypes && newOptions.ExactOptionalPropertyTypes == core.TSUnknown {
			newOptions.ExactOptionalPropertyTypes = existingOptions.ExactOptionalPropertyTypes
		}
	case "ExplainFiles":
		if existingOptions.ExplainFiles != newOptions.ExplainFiles && newOptions.ExplainFiles == core.TSUnknown {
			newOptions.ExplainFiles = existingOptions.ExplainFiles
		}
	case "ExperimentalDecorators":
		if existingOptions.ExperimentalDecorators != newOptions.ExperimentalDecorators && newOptions.ExperimentalDecorators == core.TSUnknown {
			newOptions.ExperimentalDecorators = existingOptions.ExperimentalDecorators
		}
	case "ForceConsistentCasingInFileNames":
		if existingOptions.ForceConsistentCasingInFileNames != newOptions.ForceConsistentCasingInFileNames && newOptions.ForceConsistentCasingInFileNames == core.TSUnknown {
			newOptions.ForceConsistentCasingInFileNames = existingOptions.ForceConsistentCasingInFileNames
		}
	case "GenerateCpuProfile":
		if existingOptions.GenerateCpuProfile != newOptions.GenerateCpuProfile && newOptions.GenerateCpuProfile == "" {
			newOptions.GenerateCpuProfile = existingOptions.GenerateCpuProfile
		}
	case "GenerateTrace":
		if existingOptions.GenerateTrace != newOptions.GenerateTrace && newOptions.GenerateTrace == "" {
			newOptions.GenerateTrace = existingOptions.GenerateTrace
		}
	case "IsolatedModules":
		if existingOptions.IsolatedModules != newOptions.IsolatedModules && newOptions.IsolatedModules == core.TSUnknown {
			newOptions.IsolatedModules = existingOptions.IsolatedModules
		}
	case "IgnoreDeprecations":
		if existingOptions.IgnoreDeprecations != newOptions.IgnoreDeprecations && newOptions.IgnoreDeprecations == "" {
			newOptions.IgnoreDeprecations = existingOptions.IgnoreDeprecations
		}
	case "ImportHelpers":
		if existingOptions.ImportHelpers != newOptions.ImportHelpers && newOptions.ImportHelpers == core.TSUnknown {
			newOptions.ImportHelpers = existingOptions.ImportHelpers
		}
	case "Incremental":
		if existingOptions.Incremental != newOptions.Incremental && newOptions.Incremental == core.TSUnknown {
			newOptions.Incremental = existingOptions.Incremental
		}
	case "Init":
		if existingOptions.Init != newOptions.Init && newOptions.Init == core.TSUnknown {
			newOptions.Init = existingOptions.Init
		}
	case "InlineSourceMap":
		if existingOptions.InlineSourceMap != newOptions.InlineSourceMap && newOptions.InlineSourceMap == core.TSUnknown {
			newOptions.InlineSourceMap = existingOptions.InlineSourceMap
		}
	case "InlineSources":
		if existingOptions.InlineSources != newOptions.InlineSources && newOptions.InlineSources == core.TSUnknown {
			newOptions.InlineSources = existingOptions.InlineSources
		}
	case "IsolatedDeclarations":
		if existingOptions.IsolatedDeclarations != newOptions.IsolatedDeclarations && newOptions.IsolatedDeclarations == core.TSUnknown {
			newOptions.IsolatedDeclarations = existingOptions.IsolatedDeclarations
		}
	case "Jsx":
		if existingOptions.Jsx != newOptions.Jsx && newOptions.Jsx == core.JsxEmitNone {
			newOptions.Jsx = existingOptions.Jsx
		}
	case "JsxFactory":
		if existingOptions.JsxFactory != newOptions.JsxFactory && newOptions.JsxFactory == "" {
			newOptions.JsxFactory = existingOptions.JsxFactory
		}
	case "JsxFragmentFactory":
		if existingOptions.JsxFragmentFactory != newOptions.JsxFragmentFactory && newOptions.JsxFragmentFactory == "" {
			newOptions.JsxFragmentFactory = existingOptions.JsxFragmentFactory
		}
	case "JsxImportSource":
		if existingOptions.JsxImportSource != newOptions.JsxImportSource && newOptions.JsxImportSource == "" {
			newOptions.JsxImportSource = existingOptions.JsxImportSource
		}
	case "KeyofStringsOnly":
		if existingOptions.KeyofStringsOnly != newOptions.KeyofStringsOnly && newOptions.KeyofStringsOnly == core.TSUnknown {
			newOptions.KeyofStringsOnly = existingOptions.KeyofStringsOnly
		}
	case "Lib":
		if !reflect.DeepEqual(existingOptions.Lib, newOptions.Lib) && len(newOptions.Lib) == 0 {
			newOptions.Lib = existingOptions.Lib
		}
	case "ListEmittedFiles":
		if existingOptions.ListEmittedFiles != newOptions.ListEmittedFiles && newOptions.ListEmittedFiles == core.TSUnknown {
			newOptions.ListEmittedFiles = existingOptions.ListEmittedFiles
		}
	case "ListFiles":
		if existingOptions.ListFiles != newOptions.ListFiles && newOptions.ListFiles == core.TSUnknown {
			newOptions.ListFiles = existingOptions.ListFiles
		}
	case "ListFilesOnly":
		if existingOptions.ListFilesOnly != newOptions.ListFilesOnly && newOptions.ListFilesOnly == core.TSUnknown {
			newOptions.ListFilesOnly = existingOptions.ListFilesOnly
		}
	case "Locale":
		if existingOptions.Locale != newOptions.Locale && newOptions.Locale == "" {
			newOptions.Locale = existingOptions.Locale
		}
	case "MapRoot":
		if existingOptions.MapRoot != newOptions.MapRoot && newOptions.MapRoot == "" {
			newOptions.MapRoot = existingOptions.MapRoot
		}
	case "ModuleKind":
		if existingOptions.ModuleKind != newOptions.ModuleKind && newOptions.ModuleKind == core.ModuleKindNone {
			newOptions.ModuleKind = existingOptions.ModuleKind
		}
	case "ModuleResolution":
		if existingOptions.ModuleResolution != newOptions.ModuleResolution && newOptions.ModuleResolution == core.ModuleResolutionKindUnknown {
			newOptions.ModuleResolution = existingOptions.ModuleResolution
		}
	case "ModuleSuffixes":
		if !reflect.DeepEqual(existingOptions.ModuleSuffixes, newOptions.ModuleSuffixes) && len(newOptions.ModuleSuffixes) == 0 {
			newOptions.ModuleSuffixes = existingOptions.ModuleSuffixes
		}
	case "ModuleDetection":
		if existingOptions.ModuleDetection != newOptions.ModuleDetection && newOptions.ModuleDetection == core.ModuleDetectionKindNone {
			newOptions.ModuleDetection = existingOptions.ModuleDetection
		}
	case "NoCheck":
		if existingOptions.NoCheck != newOptions.NoCheck && newOptions.NoCheck == core.TSUnknown {
			newOptions.NoCheck = existingOptions.NoCheck
		}
	case "NoFallthroughCasesInSwitch":
		if existingOptions.NoFallthroughCasesInSwitch != newOptions.NoFallthroughCasesInSwitch && newOptions.NoFallthroughCasesInSwitch == core.TSUnknown {
			newOptions.NoFallthroughCasesInSwitch = existingOptions.NoFallthroughCasesInSwitch
		}
	case "NoEmitForJsFiles":
		if existingOptions.NoEmitForJsFiles != newOptions.NoEmitForJsFiles && newOptions.NoEmitForJsFiles == core.TSUnknown {
			newOptions.NoEmitForJsFiles = existingOptions.NoEmitForJsFiles
		}
	case "NoImplicitAny":
		if existingOptions.NoImplicitAny != newOptions.NoImplicitAny && newOptions.NoImplicitAny == core.TSUnknown {
			newOptions.NoImplicitAny = existingOptions.NoImplicitAny
		}
	case "NoImplicitThis":
		if existingOptions.NoImplicitThis != newOptions.NoImplicitThis && newOptions.NoImplicitThis == core.TSUnknown {
			newOptions.NoImplicitThis = existingOptions.NoImplicitThis
		}
	case "NoPropertyAccessFromIndexSignature":
		if existingOptions.NoPropertyAccessFromIndexSignature != newOptions.NoPropertyAccessFromIndexSignature && newOptions.NoPropertyAccessFromIndexSignature == core.TSUnknown {
			newOptions.NoPropertyAccessFromIndexSignature = existingOptions.NoPropertyAccessFromIndexSignature
		}
	case "NoUncheckedIndexedAccess":
		if existingOptions.NoUncheckedIndexedAccess != newOptions.NoUncheckedIndexedAccess && newOptions.NoUncheckedIndexedAccess == core.TSUnknown {
			newOptions.NoUncheckedIndexedAccess = existingOptions.NoUncheckedIndexedAccess
		}
	case "NoEmitHelpers":
		if existingOptions.NoEmitHelpers != newOptions.NoEmitHelpers && newOptions.NoEmitHelpers == core.TSUnknown {
			newOptions.NoEmitHelpers = existingOptions.NoEmitHelpers
		}
	case "NoEmitOnError":
		if existingOptions.NoEmitOnError != newOptions.NoEmitOnError && newOptions.NoEmitOnError == core.TSUnknown {
			newOptions.NoEmitOnError = existingOptions.NoEmitOnError
		}
	case "NoImplicitReturns":
		if existingOptions.NoImplicitReturns != newOptions.NoImplicitReturns && newOptions.NoImplicitReturns == core.TSUnknown {
			newOptions.NoImplicitReturns = existingOptions.NoImplicitReturns
		}
	case "NoUnusedLocals":
		if existingOptions.NoUnusedLocals != newOptions.NoUnusedLocals && newOptions.NoUnusedLocals == core.TSUnknown {
			newOptions.NoUnusedLocals = existingOptions.NoUnusedLocals
		}
	case "NoUnusedParameters":
		if existingOptions.NoUnusedParameters != newOptions.NoUnusedParameters && newOptions.NoUnusedParameters == core.TSUnknown {
			newOptions.NoUnusedParameters = existingOptions.NoUnusedParameters
		}
	case "NoImplicitOverride":
		if existingOptions.NoImplicitOverride != newOptions.NoImplicitOverride && newOptions.NoImplicitOverride == core.TSUnknown {
			newOptions.NoImplicitOverride = existingOptions.NoImplicitOverride
		}
	case "NoUncheckedSideEffectImports":
		if existingOptions.NoUncheckedSideEffectImports != newOptions.NoUncheckedSideEffectImports && newOptions.NoUncheckedSideEffectImports == core.TSUnknown {
			newOptions.NoUncheckedSideEffectImports = existingOptions.NoUncheckedSideEffectImports
		}
	case "Out":
		if existingOptions.Out != newOptions.Out && newOptions.Out == "" {
			newOptions.Out = existingOptions.Out
		}
	case "OutFile":
		if existingOptions.OutFile != newOptions.OutFile && newOptions.OutFile == "" {
			newOptions.OutFile = existingOptions.OutFile
		}
	case "NoResolve":
		if existingOptions.NoResolve != newOptions.NoResolve && newOptions.NoResolve == core.TSUnknown {
			newOptions.NoResolve = existingOptions.NoResolve
		}
	case "Paths":
		if !reflect.DeepEqual(existingOptions.Paths, newOptions.Paths) && newOptions.Paths == nil {
			newOptions.Paths = existingOptions.Paths
		}
	case "PreserveWatchOutput":
		if existingOptions.PreserveWatchOutput != newOptions.PreserveWatchOutput && newOptions.PreserveWatchOutput == core.TSUnknown {
			newOptions.PreserveWatchOutput = existingOptions.PreserveWatchOutput
		}
	case "PreserveConstEnums":
		if existingOptions.PreserveConstEnums != newOptions.PreserveConstEnums && newOptions.PreserveConstEnums == core.TSUnknown {
			newOptions.PreserveConstEnums = existingOptions.PreserveConstEnums
		}
	case "PreserveSymlinks":
		if existingOptions.PreserveSymlinks != newOptions.PreserveSymlinks && newOptions.PreserveSymlinks == core.TSUnknown {
			newOptions.PreserveSymlinks = existingOptions.PreserveSymlinks
		}
	case "Project":
		if existingOptions.Project != newOptions.Project && newOptions.Project == "" {
			newOptions.Project = existingOptions.Project
		}
	case "Pretty":
		if existingOptions.Pretty != newOptions.Pretty && newOptions.Pretty == core.TSUnknown {
			newOptions.Pretty = existingOptions.Pretty
		}
	case "ResolveJsonModule":
		if existingOptions.ResolveJsonModule != newOptions.ResolveJsonModule && newOptions.ResolveJsonModule == core.TSUnknown {
			newOptions.ResolveJsonModule = existingOptions.ResolveJsonModule
		}
	case "ResolvePackageJsonExports":
		if existingOptions.ResolvePackageJsonExports != newOptions.ResolvePackageJsonExports && newOptions.ResolvePackageJsonExports == core.TSUnknown {
			newOptions.ResolvePackageJsonExports = existingOptions.ResolvePackageJsonExports
		}
	case "ResolvePackageJsonImports":
		if existingOptions.ResolvePackageJsonImports != newOptions.ResolvePackageJsonImports && newOptions.ResolvePackageJsonImports == core.TSUnknown {
			newOptions.ResolvePackageJsonImports = existingOptions.ResolvePackageJsonImports
		}
	case "ReactNamespace":
		if existingOptions.ReactNamespace != newOptions.ReactNamespace && newOptions.ReactNamespace == "" {
			newOptions.ReactNamespace = existingOptions.ReactNamespace
		}
	case "RootDir":
		if existingOptions.RootDir != newOptions.RootDir && newOptions.RootDir == "" {
			newOptions.RootDir = existingOptions.RootDir
		}
	case "RootDirs":
		if !reflect.DeepEqual(existingOptions.RootDirs, newOptions.RootDirs) && len(newOptions.RootDirs) == 0 {
			newOptions.RootDirs = existingOptions.RootDirs
		}
	case "RemoveComments":
		if existingOptions.RemoveComments != newOptions.RemoveComments && newOptions.RemoveComments == core.TSUnknown {
			newOptions.RemoveComments = existingOptions.RemoveComments
		}
	case "Strict":
		if existingOptions.Strict != newOptions.Strict && newOptions.Strict == core.TSUnknown {
			newOptions.Strict = existingOptions.Strict
		}
	case "StrictBindCallApply":
		if existingOptions.StrictBindCallApply != newOptions.StrictBindCallApply && newOptions.StrictBindCallApply == core.TSUnknown {
			newOptions.StrictBindCallApply = existingOptions.StrictBindCallApply
		}
	case "StrictFunctionTypes":
		if existingOptions.StrictFunctionTypes != newOptions.StrictFunctionTypes && newOptions.StrictFunctionTypes == core.TSUnknown {
			newOptions.StrictFunctionTypes = existingOptions.StrictFunctionTypes
		}
	case "StrictNullChecks":
		if existingOptions.StrictNullChecks != newOptions.StrictNullChecks && newOptions.StrictNullChecks == core.TSUnknown {
			newOptions.StrictNullChecks = existingOptions.StrictNullChecks
		}
	case "StrictPropertyInitialization":
		if existingOptions.StrictPropertyInitialization != newOptions.StrictPropertyInitialization && newOptions.StrictPropertyInitialization == core.TSUnknown {
			newOptions.StrictPropertyInitialization = existingOptions.StrictPropertyInitialization
		}
	case "SkipDefaultLibCheck":
		if existingOptions.SkipDefaultLibCheck != newOptions.SkipDefaultLibCheck && newOptions.SkipDefaultLibCheck == core.TSUnknown {
			newOptions.SkipDefaultLibCheck = existingOptions.SkipDefaultLibCheck
		}
	case "SourceMap":
		if existingOptions.SourceMap != newOptions.SourceMap && newOptions.SourceMap == core.TSUnknown {
			newOptions.SourceMap = existingOptions.SourceMap
		}
	case "SourceRoot":
		if existingOptions.SourceRoot != newOptions.SourceRoot && newOptions.SourceRoot == "" {
			newOptions.SourceRoot = existingOptions.SourceRoot
		}
	case "StripInternal":
		if existingOptions.StripInternal != newOptions.StripInternal && newOptions.StripInternal == core.TSUnknown {
			newOptions.StripInternal = existingOptions.StripInternal
		}
	case "SuppressOutputPathCheck":
		if existingOptions.SuppressOutputPathCheck != newOptions.SuppressOutputPathCheck && newOptions.SuppressOutputPathCheck == core.TSUnknown {
			newOptions.SuppressOutputPathCheck = existingOptions.SuppressOutputPathCheck
		}
	case "Target":
		if existingOptions.Target != newOptions.Target && newOptions.Target == core.ScriptTargetNone {
			newOptions.Target = existingOptions.Target
		}
	case "TraceResolution":
		if existingOptions.TraceResolution != newOptions.TraceResolution && newOptions.TraceResolution == core.TSUnknown {
			newOptions.TraceResolution = existingOptions.TraceResolution
		}
	case "TsBuildInfoFile":
		if existingOptions.TsBuildInfoFile != newOptions.TsBuildInfoFile && newOptions.TsBuildInfoFile == "" {
			newOptions.TsBuildInfoFile = existingOptions.TsBuildInfoFile
		}
	case "TypeRoots":
		if !reflect.DeepEqual(existingOptions.TypeRoots, newOptions.TypeRoots) && len(newOptions.TypeRoots) == 0 {
			newOptions.TypeRoots = existingOptions.TypeRoots
		}
	case "TscBuild":
		if existingOptions.TscBuild != newOptions.TscBuild && newOptions.TscBuild == core.TSUnknown {
			newOptions.TscBuild = existingOptions.TscBuild
		}
	case "Types":
		if !reflect.DeepEqual(existingOptions.Types, newOptions.Types) && len(newOptions.Types) == 0 {
			newOptions.Types = existingOptions.Types
		}
	case "UseDefineForClassFields":
		if existingOptions.UseDefineForClassFields != newOptions.UseDefineForClassFields && newOptions.UseDefineForClassFields == core.TSUnknown {
			newOptions.UseDefineForClassFields = existingOptions.UseDefineForClassFields
		}
	case "UseUnknownInCatchVariables":
		if existingOptions.UseUnknownInCatchVariables != newOptions.UseUnknownInCatchVariables && newOptions.UseUnknownInCatchVariables == core.TSUnknown {
			newOptions.UseUnknownInCatchVariables = existingOptions.UseUnknownInCatchVariables
		}
	case "VerbatimModuleSyntax":
		if existingOptions.VerbatimModuleSyntax != newOptions.VerbatimModuleSyntax && newOptions.VerbatimModuleSyntax == core.TSUnknown {
			newOptions.VerbatimModuleSyntax = existingOptions.VerbatimModuleSyntax
		}
	case "Version":
		if existingOptions.Version != newOptions.Version && newOptions.Version == core.TSUnknown {
			newOptions.Version = existingOptions.Version
		}
	case "MaxNodeModuleJsDepth":
		if existingOptions.MaxNodeModuleJsDepth != newOptions.MaxNodeModuleJsDepth && newOptions.MaxNodeModuleJsDepth == core.TSUnknown {
			newOptions.MaxNodeModuleJsDepth = existingOptions.MaxNodeModuleJsDepth
		}
	case "SkipLibCheck":
		if existingOptions.SkipLibCheck != newOptions.SkipLibCheck && newOptions.SkipLibCheck == core.TSUnknown {
			newOptions.SkipLibCheck = existingOptions.SkipLibCheck
		}
	case "NoEmit":
		if existingOptions.NoEmit != newOptions.NoEmit && newOptions.NoEmit == core.TSUnknown {
			newOptions.NoEmit = existingOptions.NoEmit
		}
	case "ShowConfig":
		if existingOptions.ShowConfig != newOptions.ShowConfig && newOptions.ShowConfig == core.TSUnknown {
			newOptions.ShowConfig = existingOptions.ShowConfig
		}
	case "ConfigFilePath":
		if existingOptions.ConfigFilePath != newOptions.ConfigFilePath && newOptions.ConfigFilePath == "" {
			newOptions.ConfigFilePath = existingOptions.ConfigFilePath
		}
	case "NoDtsResolution":
		if existingOptions.NoDtsResolution != newOptions.NoDtsResolution && newOptions.NoDtsResolution == core.TSUnknown {
			newOptions.NoDtsResolution = existingOptions.NoDtsResolution
		}
	case "PathsBasePath":
		if existingOptions.PathsBasePath != newOptions.PathsBasePath && newOptions.PathsBasePath == "" {
			newOptions.PathsBasePath = existingOptions.PathsBasePath
		}
	case "OutDir":
		if existingOptions.OutDir != newOptions.OutDir && newOptions.OutDir == "" {
			newOptions.OutDir = existingOptions.OutDir
		}
	case "NewLine":
		if existingOptions.NewLine != newOptions.NewLine && newOptions.NewLine == core.NewLineKindNone {
			newOptions.NewLine = existingOptions.NewLine
		}
	case "Watch":
		if existingOptions.Watch != newOptions.Watch && newOptions.Watch == core.TSUnknown {
			newOptions.Watch = existingOptions.Watch
		}
	}
}
