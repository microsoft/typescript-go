package tsoptions

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

func parseTristate(value interface{}) core.Tristate {
	switch v := value.(type) {
	case bool:
		if v {
			return core.TSTrue
		}
		if !v {
			return core.TSFalse
		}
	}
	return core.TSUnknown
}

func parseStringArray(value interface{}) []string {
	if arr, ok := value.([]interface{}); ok {
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

func parseRawStringArray(value interface{}) []string {
	if arr, ok := value.([]string); ok {
		return arr
	}
	return []string{}
}

func parseStringMap(value interface{}) map[string][]string {
	if m, ok := value.(map[string]interface{}); ok {
		result := make(map[string][]string)
		for k, v := range m {
			result[k] = parseStringArray(v)
		}
		return result
	}
	return nil
}

func parseString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func parseScriptTarget(json any) core.ScriptTarget {
	var result core.ScriptTarget
	if target, ok := json.(string); ok {
		target = strings.ToLower(target)
		switch target {
		case "es3":
			result = core.ScriptTargetES3
		case "es5":
			result = core.ScriptTargetES5
		case "es2015":
			result = core.ScriptTargetES2015
		case "es2016":
			result = core.ScriptTargetES2016
		case "es2017":
			result = core.ScriptTargetES2017
		case "es2018":
			result = core.ScriptTargetES2018
		case "es2019":
			result = core.ScriptTargetES2019
		case "es2020":
			result = core.ScriptTargetES2020
		case "es2021":
			result = core.ScriptTargetES2021
		case "es2022":
			result = core.ScriptTargetES2022
		case "es2023":
			result = core.ScriptTargetES2023
		case "esnext":
			result = core.ScriptTargetESNext
		default:
			result = core.ScriptTargetNone
		}
	}
	return result
}

func parseJsxEmit(json any) core.JsxEmit {
	var result core.JsxEmit
	if jsx, ok := json.(string); ok {
		switch jsx {
		case "preserve":
			result = core.JsxEmitPreserve
		case "react":
			result = core.JsxEmitReact
		case "react-native":
			result = core.JsxEmitReactNative
		case "react-jsx":
			result = core.JsxEmitReactJSX
		case "react-jsxdev":
			result = core.JsxEmitReactJSXDev
		default:
			result = core.JsxEmitNone
		}
	}
	return result
}

func parseModuleDetectionKind(json any) core.ModuleDetectionKind {
	var result core.ModuleDetectionKind
	if module, ok := json.(string); ok {
		module = strings.ToLower(module)
		switch module {
		case "auto":
			result = core.ModuleDetectionKindAuto
		case "legacy":
			result = core.ModuleDetectionKindLegacy
		case "force":
			result = core.ModuleDetectionKindForce
		default:
			result = core.ModuleDetectionKindNone
		}
	}
	return result
}

func parseModuleKind(json any) core.ModuleKind {
	var result core.ModuleKind
	if module, ok := json.(string); ok {
		module = strings.ToLower(module)
		switch module {
		case "none":
			result = core.ModuleKindNone
		case "commonjs":
			result = core.ModuleKindCommonJS
		case "amd":
			result = core.ModuleKindAMD
		case "umd":
			result = core.ModuleKindUMD
		case "system":
			result = core.ModuleKindSystem
		case "es2015":
			result = core.ModuleKindES2015
		case "es2020":
			result = core.ModuleKindES2020
		case "es2022":
			result = core.ModuleKindES2022
		case "node16":
			result = core.ModuleKindNode16
		case "esnext":
			result = core.ModuleKindESNext
		case "nodenext":
			result = core.ModuleKindNodeNext
		case "preserve":
			result = core.ModuleKindPreserve
		default:
			result = core.ModuleKindNone
		}
	}
	return result
}

func parseNewLineKind(json any) core.NewLineKind {
	var result core.NewLineKind
	if newline, ok := json.(string); ok {
		switch newline {
		case "crlf":
			result = core.NewLineKindCRLF
		case "lf":
			result = core.NewLineKindLF
		}
	}
	return result
}
func parseModuleResolutionKind(json any) core.ModuleResolutionKind {
	var result core.ModuleResolutionKind
	if module, ok := json.(string); ok {
		module = strings.ToLower(module)
		switch module {
		case "node":
			result = core.ModuleResolutionKindNode16
		case "classic":
			result = core.ModuleResolutionKindNodeNext
		case "bundler":
			result = core.ModuleResolutionKindBundler
		default:
			result = core.ModuleResolutionKindUnknown
		}
	}
	return result
}
func parseProjectReference(json any) []core.ProjectReference {
	var result []core.ProjectReference
	if arr, ok := json.([]map[string]interface{}); ok {
		for _, v := range arr {
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
	}
	return result
}

func parseJsonToStringKey(json any) map[string]interface{} {
	result := make(map[string]interface{})
	if m, ok := json.(map[string]interface{}); ok {
		if v, ok := m["include"]; ok {
			result["include"] = v
		}
		if v, ok := m["exclude"]; ok {
			result["exclude"] = v
		}
		if v, ok := m["files"]; ok {
			result["files"] = v
		}
		if v, ok := m["references"]; ok {
			result["references"] = v
		}
		if v, ok := m["extends"]; ok {
			result["extends"] = v
		}
		if v, ok := m["compilerOptions"]; ok {
			result["compilerOptions"] = v
		}
	}
	return result
}

func parseCompilerOptions(key string, value any, allOptions *core.CompilerOptions) *core.CompilerOptions {
	if allOptions == nil {
		return nil
	}
	switch key {
	case "allowJs":
		allOptions.AllowJs = parseTristate(value)
	case "allowSyntheticDefaultImports":
		allOptions.AllowSyntheticDefaultImports = parseTristate(value)
	case "allowUmdGlobalAccess":
		allOptions.AllowUmdGlobalAccess = parseTristate(value)
	case "allowUnreachableCode":
		allOptions.AllowUnreachableCode = parseTristate(value)
	case "allowUnusedLabels":
		allOptions.AllowUnusedLabels = parseTristate(value)
	case "checkJs":
		allOptions.CheckJs = parseTristate(value)
	case "customConditions":
		allOptions.CustomConditions = parseStringArray(value)
	case "declarationDir":
		allOptions.DeclarationDir = parseString(value)
	case "esModuleInterop":
		allOptions.ESModuleInterop = parseTristate(value)
	case "exactOptionalPropertyTypes":
		allOptions.ExactOptionalPropertyTypes = parseTristate(value)
	case "experimentalDecorators":
		allOptions.ExperimentalDecorators = parseTristate(value)
	case "isolatedModules":
		allOptions.IsolatedModules = parseTristate(value)
	case "jsx":
		allOptions.Jsx = parseJsxEmit(value)
	case "lib":
		allOptions.Lib = parseStringArray(value)
	case "legacyDecorators":
		allOptions.LegacyDecorators = parseTristate(value)
	case "module":
		allOptions.ModuleKind = parseModuleKind(value)
	case "moduleResolution":
		allOptions.ModuleResolution = parseModuleResolutionKind(value)
	case "moduleSuffixes":
		allOptions.ModuleSuffixes = parseStringArray(value)
	case "moduleDetectionKind":
		allOptions.ModuleDetection = parseModuleDetectionKind(value)
	case "noFallthroughCasesInSwitch":
		allOptions.NoFallthroughCasesInSwitch = parseTristate(value)
	case "noImplicitAny":
		allOptions.NoImplicitAny = parseTristate(value)
	case "noImplicitThis":
		allOptions.NoImplicitThis = parseTristate(value)
	case "noPropertyAccessFromIndexSignature":
		allOptions.NoPropertyAccessFromIndexSignature = parseTristate(value)
	case "noUncheckedIndexedAccess":
		allOptions.NoUncheckedIndexedAccess = parseTristate(value)
	case "paths":
		allOptions.Paths = parseStringMap(value)
	case "preserveConstEnums":
		allOptions.PreserveConstEnums = parseTristate(value)
	case "preserveSymlinks":
		allOptions.PreserveSymlinks = parseTristate(value)
	case "resolveJsonModule":
		allOptions.ResolveJsonModule = parseTristate(value)
	case "resolvePackageJsonExports":
		allOptions.ResolvePackageJsonExports = parseTristate(value)
	case "resolvePackageJsonImports":
		allOptions.ResolvePackageJsonImports = parseTristate(value)
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
	case "target":
		allOptions.Target = parseScriptTarget(value)
	case "traceResolution":
		allOptions.TraceResolution = parseTristate(value)
	case "typeRoots":
		allOptions.TypeRoots = parseStringArray(value)
	case "types":
		allOptions.Types = parseStringArray(value)
	case "useDefineForClassFields":
		allOptions.UseDefineForClassFields = parseTristate(value)
	case "useUnknownInCatchVariables":
		allOptions.UseUnknownInCatchVariables = parseTristate(value)
	case "verbatimModuleSyntax":
		allOptions.VerbatimModuleSyntax = parseTristate(value)
	case "maxNodeModuleJsDepth":
		allOptions.MaxNodeModuleJsDepth = parseTristate(value)
	case "skipLibCheck":
		allOptions.SkipLibCheck = parseTristate(value)
	case "noEmit":
		allOptions.NoEmit = parseTristate(value)
	case "configFilePath":
		allOptions.ConfigFilePath = parseString(value)
	case "noDtsResolution":
		allOptions.NoDtsResolution = parseTristate(value)
	case "pathsBasePath":
		allOptions.PathsBasePath = parseString(value)
	case "outDir":
		allOptions.OutDir = parseString(value)
	case "newLine":
		allOptions.NewLine = parseNewLineKind(value)
	default:
		// Handle unknown options
		fmt.Printf("Unknown option: %s\n", key)
	}
	return allOptions
}

func ParseRawConfig(json any, basePath string, errors []*ast.Diagnostic, configFileName string) tsConfigOptions {
	options := tsConfigOptions{
		prop: make(map[string][]string),
	}
	if json == nil {
		return options
	}
	if m, ok := json.(map[string]interface{}); ok {
		if v, ok := m["include"]; ok {
			options.prop["include"] = parseRawStringArray(v)
		}
		if v, ok := m["exclude"]; ok {
			options.prop["exclude"] = parseRawStringArray(v)
		}
		if v, ok := m["files"]; ok {
			options.prop["files"] = parseRawStringArray(v)
		}
		if v, ok := m["references"]; ok {
			options.references = parseProjectReference(v)
		}
		if v, ok := m["extends"]; ok {
			options.prop["extends"] = parseRawStringArray(v)
		}
		if v, ok := m["compilerOptions"]; ok {
			var option *core.CompilerOptions = &core.CompilerOptions{}
			if vMap, ok := v.(map[string]interface{}); ok {
				for key, value := range vMap {
					parseCompilerOptions(key, value, option)
				}
				options.compilerOptionsProp = *option
			}
		}
	}
	return options
}
