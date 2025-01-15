package tsoptions

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

func parseTristate(value any) core.Tristate {
	switch v := value.(type) {
	case bool:
		if v {
			return core.TSTrue
		}
		return core.TSFalse
	}
	return core.TSUnknown
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
		allOptions.Jsx = value.(core.JsxEmit)
	case "lib":
		if _, ok := value.([]string); ok {
			allOptions.Lib = value.([]string)
		} else {
			allOptions.Lib = parseStringArray(value)
		}
	case "module":
		allOptions.ModuleKind = value.(core.ModuleKind)
	case "moduleResolution":
		allOptions.ModuleResolution = value.(core.ModuleResolutionKind)
	case "moduleSuffixes":
		allOptions.ModuleSuffixes = parseStringArray(value)
	case "moduleDetection":
		allOptions.ModuleDetection = value.(core.ModuleDetectionKind)
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
		allOptions.Target = value.(core.ScriptTarget)
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
		allOptions.NewLine = value.(core.NewLineKind)
	}
	return nil
}
