package tsoptions

import (

	//"slices"

	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/compiler/module"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type parsedConfigFileTextToJsonResult struct {
	config interface{}
	error  *ast.Diagnostic
}

// func convertConfigFileToObject(sourceFile *SourceFile, )

const (
	typeString        = "string"
	typeNumber        = "number"
	typeBoolean       = "boolean"
	typeObject        = "object"
	typeList          = "list"
	typeListOrElement = "listOrElement"
)

type optionType string
type optionMap map[string]int //should be number | string
type commandLineOptionBaseType struct {
	optionType
	optionMap
}

type extendsResult struct {
	options core.CompilerOptions
	//watchOptions        compiler.WatchOptions
	watchOptionsCopied  bool
	include             *[]string
	exclude             *[]string
	files               *[]string
	compileOnSave       *bool
	extendedSourceFiles *map[string]struct{} //*Set<string>;
}

var tsconfigRootOptions *CommandLineOption //TsConfigOnlyOption

func getTsconfigRootOptionsMap() CommandLineOption { //TsConfigOnlyOption
	if tsconfigRootOptions == nil {
		tsconfigRootOptions = &CommandLineOption{
			Name: "undefined", //undefined! // should never be needed since this is root
			Kind: CommandLineOptionTypeObject,
			ElementOptions: commandLineOptionsToMap([]CommandLineOption{
				compilerOptionsDeclaration,
				{
					Name: "references",
					Kind: CommandLineOptionTypeList, //should be a list of projectReference
					//Category: diagnostics.Projects,
				},
				{
					Name: "files",
					Kind: CommandLineOptionTypeList,
					//Category: diagnostics.File_Management,
				},
				{
					Name: "include",
					Kind: CommandLineOptionTypeList,
					//Category: diagnostics.File_Management,
					//DefaultValueDescription: diagnostics.if_files_is_specified_otherwise_Asterisk_Asterisk_Slash_Asterisk,
				},
				{
					Name: "exclude",
					Kind: CommandLineOptionTypeList,
					//Category: diagnostics.File_Management,
					//DefaultValueDescription: diagnostics.Node_modules_bower_components_jspm_packages_plus_the_value_of_outDir_if_one_is_specified,
				},
				compileOnSaveCommandLineOption,
			}),
		}
	}
	tsconfigRootOptions.Elements()
	return *tsconfigRootOptions
}

type configFileSpecs struct {
	filesSpecs []string
	/**
	 * Present to report errors (user specified specs), validatedIncludeSpecs are used for file name matching
	 */
	includeSpecs []string
	/**
	 * Present to report errors (user specified specs), validatedExcludeSpecs are used for file name matching
	 */
	excludeSpecs                            []string
	validatedFilesSpec                      []string
	validatedIncludeSpecs                   []string
	validatedExcludeSpecs                   []string
	validatedFilesSpecBeforeSubstitution    []string
	validatedIncludeSpecsBeforeSubstitution []string
	validatedExcludeSpecsBeforeSubstitution []string
	isDefaultIncludeSpec                    bool
}

type ParseConfigHost struct {
	module.ResolutionHost
	useCaseSensitiveFileNames bool
	readDirectory             func(rootDir string, extensions []string, excludes []string, includes []string, depth int) []string
	/**
	 * Gets a value indicating whether the specified path exists and is a file.
	 * @param path The path to test.
	 */
	fileExists func(path string) bool
	readFile   func(path string) string
	trace      func(s string)
}

type FileExtensionInfo struct {
	extension      string
	isMixedContent bool
	scriptKind     core.ScriptKind
}
type ExtendedConfigCacheEntry struct {
	extendedResult *ast.SourceFile
	extendedConfig ParsedTsconfig
}
type ParsedTsconfig struct {
	raw     any
	options *core.CompilerOptions
	//watchOptions    *compiler.WatchOptions
	//typeAcquisition *compiler.TypeAcquisition
	// Note that the case of the config path has not yet been normalized, as no files have been imported into the project yet
	extendedConfigPath *[]string
}

func isSuccessfulParsedTsconfig(value ParsedTsconfig) bool {
	return value.options != nil
}

func parseOwnConfigOfJsonSourceFile(
	sourceFile *tsConfigSourceFile,
	host VfsParseConfigHost,
	basePath string,
	configFileName *string,
	errors []*ast.Diagnostic,
) (*ParsedTsconfig, []*ast.Diagnostic) {
	options := getDefaultCompilerOptions(*configFileName)
	//var typeAcquisition *compiler.TypeAcquisition
	//var watchOptions *compiler.WatchOptions
	//var extendedConfigPath []string = []string{} // | string
	// var rootCompilerOptions []ast.PropertyName

	rootOptions := getTsconfigRootOptionsMap()
	onPropertySet := func(
		keyText string,
		value any,
		propertyAssignment ast.PropertyAssignment,
		parentOption CommandLineOption, //TsConfigOnlyOption,
		option *CommandLineOption,
	) {
		// Ensure value is verified except for extends which is handled in its own way for error reporting
		if option != nil && option != &extendsOptionDeclaration { //&& option != extendsOptionDeclaration {
			value, _ = convertJsonOption(*option, value, basePath, errors, &propertyAssignment, propertyAssignment.Initializer, sourceFile)
		}
		if parentOption.Name != "undefined" && value != nil {
			if option != nil && option.Name != "" {
				option.Name = value.(string)
			} else if keyText != "" { //&& parentOption.extraKeydiagnostics {
				if parentOption.ElementOptions != nil {
					errors = append(errors, compiler.NewDiagnosticForNode(&sourceFile.sourceFile.Node, diagnostics.Option_build_must_be_the_first_command_line_argument, keyText))
				} else {
					// errors = append(errors, compiler.NewDiagnosticForNode(&sourceFile.sourceFile.Node, diagnostics.Unknown_option_0_Did_you_mean_1, keyText, core.FindKey(parentOption.ElementOptions, keyText)))
				}
			}
		} else if parentOption.Name == rootOptions.Name { //need to compare both structs
			if option == &extendsOptionDeclaration { //todo in 2nd iteration
				//extendedConfigPath = getExtendsConfigPathOrArray(value, host, basePath, configFileName, errors, propertyAssignment, propertyAssignment.initializer, sourceFile)
			} else if option.Name == "" { //option == nil
				if keyText == "excludes" {
					errors = append(errors, compiler.NewDiagnosticForNode(&sourceFile.sourceFile.Node, diagnostics.Unknown_option_excludes_Did_you_mean_exclude, keyText))
				}
				core.Find(optionsDeclarations, func(option CommandLineOption) bool {
					return option.Name == keyText
				})
			}
		}
	}

	json := convertConfigFileToObject(
		sourceFile.sourceFile,
		errors,
		&JsonConversionNotifier{
			rootOptions,
			onPropertySet,
		},
	)

	// if rootCompilerOptions != nil && json != nil { //&& json.compilerOptions == nil {
	// 	//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, rootCompilerOptions[0], diagnostics._0_should_be_set_inside_the_compilerOptions_object_of_the_config_json_file, getTextOfPropertyName(rootCompilerOptions[0]) as string));
	// }

	return &ParsedTsconfig{
		raw:     json,
		options: &options,
		//watchOptions:    watchOptions,
		// typeAcquisition: typeAcquisition,
		//extendedConfigPath: extendedConfigPath,
	}, errors

}

func getExtendedConfig(
	sourceFile *tsConfigSourceFile,
	extendedConfigPath string,
	host VfsParseConfigHost,
	resolutionStack []string,
	errors []*ast.Diagnostic,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
	result extendsResult,
) ParsedTsconfig {
	var path string
	if host.FS().UseCaseSensitiveFileNames() {
		path = extendedConfigPath
	} else {
		path = tspath.ToFileNameLowerCase(extendedConfigPath)
	}
	var value ExtendedConfigCacheEntry
	var extendedResult *tsConfigSourceFile
	var extendedConfig ParsedTsconfig

	value = (*extendedConfigCache)[path]
	if extendedConfigCache != nil && value == (ExtendedConfigCacheEntry{}) {
		extendedResult.sourceFile = value.extendedResult
		extendedConfig = value.extendedConfig
	}
	// else {
	// 	contents, _ := host.FS().ReadFile(extendedConfigPath)
	// 	extendedResult.SourceFile = contents//readJsonConfigFile(extendedConfigPath, contents) //probably readfile will give undefined
	// 	if extendedResult != nil {                                                             //parsediagnostics.length { //come back
	// 		extendedConfig = parseConfig(nil, extendedResult, host, tspath.GetDirectoryPath(extendedConfigPath), tspath.GetBaseFileName(extendedConfigPath), resolutionStack, errors, extendedConfigCache)
	// 	}
	// 	if extendedConfigCache != nil {
	// 		(*extendedConfigCache)[path] = ExtendedConfigCacheEntry{extendedResult.SourceFile, extendedConfig}
	// 	}
	// }
	if sourceFile != nil {
		// if (extendedResult.extendedSourceFiles == nil) { //todo currently sourcefile does not have extendedSourceFiles
		// 	extendedResult.extendedSourceFiles = make(map[string]struct{});
		// }
		// if (extendedResult.extendedSourceFiles) {
		// 	for _, extenedSourceFile := range extendedResult.extendedSourceFiles {
		// 		result.extendedSourceFiles[extenedSourceFile] = struct{}
		// 	}
		// }
	}
	// if (extendedResult.parsediagnostics.length) {
	//     //errors.push(...extendedResult.parsediagnostics);
	//     return undefined;
	// }
	return extendedConfig //extendedConfig!
}

type readFile func(path string) string

func tryReadFile(fileName string, readFile readFile) (string, diagnostics.Message) {
	var text string
	if readFile(fileName) != "" {
		text = readFile(fileName)
	} else {
		//return "", createCompilerDiagnostic(Diagnostics.Cannot_read_file_0, fileName)
	}

	if text == "" {
		//createCompilerDiagnostic(Diagnostics.Cannot_read_file_0, fileName)
		return text, diagnostics.Message{} //remove later
	} else {
		return text, diagnostics.Message{}
	}
}

type tsConfigSourceFile struct {
	extendedSourceFiles []string
	configFileSpecs     *configFileSpecs
	sourceFile          *ast.SourceFile
}

/**
 * Read tsconfig.json file
 * @param fileName The path to the config file
 */
func readJsonConfigFile(fileName string, readFile readFile) *ast.SourceFile {
	var text, _ = tryReadFile(fileName, readFile)
	if text != "" {
		return compiler.ParseJSONText(fileName, text)
	} else {
		// result := tsConfigSourceFile {
		// 	//fileName: fileName,
		// 	SourceFile: []ast.Diagnostic{diagnostic},
		// }
		return nil
	}
}

type JsonConversionNotifier struct {
	rootOptions   CommandLineOption //TsConfigOnlyOption
	onPropertySet func(keyText string, value any, propertyAssignment ast.PropertyAssignment, parentOption CommandLineOption, option *CommandLineOption)
}

type defaultValueDescriptionType struct {
	valueString     string
	valueNumber     int
	valueDiagnostic ast.Diagnostic // todo should be DiagnosticMessage
}

type listType string

const listTypeList = listType("list")
const listTypeListOrElement = listType("listOrElement")

type configFileToObject struct {
	compilerOptions *core.CompilerOptions
	//watchOptions    *compiler.WatchOptions
	//typeAcquisition *compiler.TypeAcquisition
	include *[]string
	exclude *[]string
}

func convertConfigFileToObject(
	sourceFile *ast.SourceFile,
	errors []*ast.Diagnostic,
	jsonConversionNotifier *JsonConversionNotifier,
) any {
	var rootExpression *ast.Expression
	if len(sourceFile.Statements.Nodes) > 0 { //check
		rootExpression = sourceFile.Statements.Nodes[0].AsExpressionStatement().Expression
	}
	if rootExpression != nil && rootExpression.Kind != ast.KindObjectLiteralExpression {
		// errors.push(createDiagnosticForNodeInSourceFile(
		//     sourceFile,
		//     rootExpression,
		//     diagnostics.The_root_value_of_a_0_file_must_be_an_object,
		//     getBaseFileName(sourceFile.fileName) === "jsconfig.json" ? "jsconfig.json" : "tsconfig.json",
		// ));
		// Last-ditch error recovery. Somewhat useful because the JSON parser will recover from some parse errors by
		// synthesizing a top-level array literal expression. There's a reasonable chance the first element of that
		// array is a well-formed configuration object, made into an array element by stray characters.
		if ast.IsArrayLiteralExpression(rootExpression) {
			var firstObject = core.Find(rootExpression.AsArrayLiteralExpression().Elements.Nodes, ast.IsObjectLiteralExpression)
			if firstObject != nil {
				return convertToJson(sourceFile, firstObject, errors /*returnValue*/, true, jsonConversionNotifier)
			}
		}
		return nil
	}
	return convertToJson(sourceFile, rootExpression, errors, true, jsonConversionNotifier)
}

type pluginImport struct {
	name string
}

func isCompilerOptionsValue(option CommandLineOption, value any) core.CompilerOptionsValue {
	if option.Name != "" || option.Kind != "" {
		//if compiler.Checker.IsNullOrUndefined(value) {
		if value == nil {
			return core.CompilerOptionsValue{BooleanValue: !option.DisallowNullOrUndefined()} // All options are undefinable/nullable
		}
		if option.Kind == "list" {
			_, ok := value.([]string)
			return core.CompilerOptionsValue{BooleanValue: ok}
		}
		if option.Kind == "listOrElement" {
			_, ok := value.([]string)
			return core.CompilerOptionsValue{BooleanValue: ok}
			//isCompilerOptionsValue(option.element, value);
		}
		if option.Kind == "string" {
			_, ok := value.(string)
			return core.CompilerOptionsValue{BooleanValue: ok}
		}
		if option.Kind == "object" {
			return core.CompilerOptionsValue{BooleanValue: true}
		}
		// todo: other types need to be checked
		return core.CompilerOptionsValue{BooleanValue: false}
	}
	return core.CompilerOptionsValue{BooleanValue: false}
}

func mapValues(values []interface{}, transform func(interface{}, int) interface{}) []interface{} {
	result := make([]interface{}, len(values))
	for i, v := range values {
		result[i] = transform(v, i)
	}
	return result
}

func validateJsonOptionValue(
	opt CommandLineOption,
	val any,
	errors []*ast.Diagnostic,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) any {
	if val == nil || val == "" {
		return core.CompilerOptionsValue{}
	}
	d := (opt.extraValidation)
	if d == nil {
		return val
	} else {
		//d = opt.extraValidation.val
	}
	//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, d));
	return core.CompilerOptionsValue{}
}

func convertJsonOptionOfCustomType(
	opt CommandLineOption,
	value string,
	errors []*ast.Diagnostic,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) any {
	if value == "" {
		return core.CompilerOptionsValue{}
	}
	key := strings.ToLower(value)
	typeMap := opt.EnumMap()
	if typeMap == nil {
		return core.CompilerOptionsValue{}
	}
	val, b := typeMap.Get(key)
	if (val != nil) && (val != "" || b) { //need to check
		return validateJsonOptionValue(opt, val.(string), errors, valueExpression, sourceFile)
	}
	// else {
	//     errors.push(createDiagnosticForInvalidCustomType(opt, (message, ...args) => createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, message, ...args)));
	// }
	return core.CompilerOptionsValue{}
}

func convertJsonOptionOfListType(
	option CommandLineOption,
	values []string, //readonly
	basePath string,
	errors []*ast.Diagnostic,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Node,
	sourceFile *tsConfigSourceFile,
) []any {
	index := 0 //need to be changed
	var expression *ast.Node
	mappedValue := core.Map(values, func(v string) any {
		if valueExpression != nil {
			expression = valueExpression.AsArrayLiteralExpression().Elements.Nodes[index]
		}
		var t, _ = convertJsonOption(*option.Elements(), v, basePath, errors, propertyAssignment, expression, sourceFile)
		index++
		return t
	})
	filteredValues := core.Filter(mappedValue, func(v any) bool {
		if option.listPreserveFalsyValues {
			return true
		} else {
			return (v != nil && v != false && v != 0 && v != "") //!!v
		}
	})
	return filteredValues
}

const configDirTemplate = "${configDir}"

func startsWithConfigDirTemplate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return strings.HasPrefix(strings.ToLower(str), strings.ToLower(configDirTemplate))
}

func normalizeNonListOptionValue(option CommandLineOption, basePath string, value any) any {
	if option.isFilePath {
		// value = tspath.NormalizeSlashes(value) //what is value is not a string
		if !startsWithConfigDirTemplate(value) {
			value = tspath.GetNormalizedAbsolutePath(value.(string), basePath)
		}
		//value = !startsWithConfigDirTemplate(value) ? getNormalizedAbsolutePath(value, basePath) : value;
		if value == "" {
			value = "."
		}
	}
	return value
}

func convertJsonOption(
	opt CommandLineOption,
	value any,
	basePath string,
	errors []*ast.Diagnostic,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) (any, []*ast.Diagnostic) {
	if opt.isCommandLineOnly != false {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, propertyAssignment?.name, diagnostics.Option_0_can_only_be_specified_on_command_line, opt.name));
		return core.CompilerOptionsValue{}, errors
	}
	if isCompilerOptionsValue(opt, value).BooleanValue {
		optType := opt.Kind
		_, ok := value.([]string)
		if (optType == "list") && ok {
			list := convertJsonOptionOfListType(opt, value.([]string), basePath, errors, propertyAssignment, valueExpression, sourceFile) //as ArrayLiteralExpression | undefined
			return list, errors
		} else if optType == "listOrElement" {
			if ok {
				return convertJsonOptionOfListType(opt, value.([]string), basePath, errors, propertyAssignment, valueExpression, sourceFile), errors
			} else {
				return convertJsonOption(*opt.Elements(), value, basePath, errors, propertyAssignment, valueExpression, sourceFile)
			}
		} else if !(reflect.TypeOf(optType).Kind() == reflect.String) {
			return convertJsonOptionOfCustomType(opt, value.(string), errors, valueExpression, sourceFile), errors
		}
		validatedValue := validateJsonOptionValue(opt, value, errors, valueExpression, sourceFile)
		if validatedValue == nil {
			return validatedValue, errors
		} else {
			return normalizeNonListOptionValue(opt, basePath, validatedValue), errors
		}
	} else {
		errors = append(errors, compiler.NewDiagnosticForNode(&sourceFile.sourceFile.Node, diagnostics.Compiler_option_0_requires_a_value_of_type_1, opt.Name, opt.Kind))
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.Compiler_option_0_requires_a_value_of_type_1, opt.name, getCompilerOptionValueTypeString(opt)));
		return nil, errors
	}
}

func c(configFileName *string) { //compiler.TypeAcquisition {
	//enable := configFileName != nil && getBaseFileName(*configFileName, nil, nil) == "jsconfig.json"
	//return compiler.TypeAcquisition{Enable: &enable, Include: &[]string{}, Exclude: &[]string{}}
	return
}

func getExtendsConfigPath(
	extendedConfig string,
	host ParseConfigHost,
	basePath string,
	errors []ast.Diagnostic,
	valueExpression ast.Expression,
	sourceFile *ast.SourceFile,
) string {
	extendedConfig = tspath.NormalizeSlashes(extendedConfig)
	if tspath.IsRootedDiskPath(extendedConfig) { // || tspath.StartsWith(extendedConfig, "./", nil) || compiler.StartsWith(extendedConfig, "../", nil) {
		extendedConfigPath := tspath.GetNormalizedAbsolutePath(extendedConfig, basePath)
		if !host.fileExists(extendedConfigPath) { //&& !compiler.EndsWith(extendedConfigPath, Extension.Json) { //need to define Extension.Json
			extendedConfigPath = `${extendedConfigPath}.json`
			if !host.fileExists(extendedConfigPath) {
				//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.File_0_not_found, extendedConfig));
				return ""
			}
		}
		return extendedConfigPath
	}
	// If the path isn't a rooted or relative path, resolve like a module
	//const resolved = nodeNextJsonConfigResolver(extendedConfig, combinePaths(basePath, "tsconfig.json"), host);
	// if (resolved.resolvedModule) {
	//     return resolved.resolvedModule.resolvedFileName;
	// }
	if extendedConfig == "" {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.Compiler_option_0_cannot_be_given_an_empty_string, "extends"));
	} else {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.File_0_not_found, extendedConfig));
	}
	return ""
}

func parseTristate(value interface{}) core.Tristate {
	switch v := value.(type) {
	case bool:
		if v {
			return 2
		}
		return 1
	case string:
		if v == "true" {
			return 2
		} else if v == "false" {
			return 1
		}
	}
	return 0
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
	return nil
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

func parseCompilerOptions(key string, value any) core.CompilerOptions {
	var options core.CompilerOptions
	//options.Option = make(map[string]core.CompilerOptionsValue)
	switch key {
	case "allowJs":
		options.AllowJs = parseTristate(value)
	case "allowSyntheticDefaultImports":
		options.AllowSyntheticDefaultImports = parseTristate(value)
	case "allowUmdGlobalAccess":
		options.AllowUmdGlobalAccess = parseTristate(value)
	case "allowUnreachableCode":
		options.AllowUnreachableCode = parseTristate(value)
	case "allowUnusedLabels":
		options.AllowUnusedLabels = parseTristate(value)
	case "checkJs":
		options.CheckJs = parseTristate(value)
	case "customConditions":
		options.CustomConditions = parseStringArray(value)
	case "declarationDir":
		options.DeclarationDir = parseString(value)
	case "esModuleInterop":
		options.ESModuleInterop = parseTristate(value)
	case "exactOptionalPropertyTypes":
		options.ExactOptionalPropertyTypes = parseTristate(value)
	case "experimentalDecorators":
		options.ExperimentalDecorators = parseTristate(value)
	case "isolatedModules":
		options.IsolatedModules = parseTristate(value)
	// case "jsx":
	//     options.Jsx = parseJsxEmit(value)
	case "lib":
		options.Lib = parseStringArray(value)
	case "legacyDecorators":
		options.LegacyDecorators = parseTristate(value)
	// case "module":
	//     options.ModuleKind = parseModuleKind(value)
	// case "moduleResolution":
	//     options.ModuleResolution = parseModuleResolutionKind(value)
	case "moduleSuffixes":
		options.ModuleSuffixes = parseStringArray(value)
	// case "moduleDetectionKind":
	//     options.ModuleDetection = parseModuleDetectionKind(value)
	case "noFallthroughCasesInSwitch":
		options.NoFallthroughCasesInSwitch = parseTristate(value)
	case "noImplicitAny":
		options.NoImplicitAny = parseTristate(value)
	case "noImplicitThis":
		options.NoImplicitThis = parseTristate(value)
	case "noPropertyAccessFromIndexSignature":
		options.NoPropertyAccessFromIndexSignature = parseTristate(value)
	case "noUncheckedIndexedAccess":
		options.NoUncheckedIndexedAccess = parseTristate(value)
	case "paths":
		options.Paths = parseStringMap(value)
	case "preserveConstEnums":
		options.PreserveConstEnums = parseTristate(value)
	case "preserveSymlinks":
		options.PreserveSymlinks = parseTristate(value)
	case "resolveJsonModule":
		options.ResolveJsonModule = parseTristate(value)
	case "resolvePackageJsonExports":
		options.ResolvePackageJsonExports = parseTristate(value)
	case "resolvePackageJsonImports":
		options.ResolvePackageJsonImports = parseTristate(value)
	case "strict":
		options.Strict = parseTristate(value)
	case "strictBindCallApply":
		options.StrictBindCallApply = parseTristate(value)
	case "strictFunctionTypes":
		options.StrictFunctionTypes = parseTristate(value)
	case "strictNullChecks":
		options.StrictNullChecks = parseTristate(value)
	case "strictPropertyInitialization":
		options.StrictPropertyInitialization = parseTristate(value)
	// case "target":
	//     options.Target = parseScriptTarget(value)
	case "traceResolution":
		options.TraceResolution = parseTristate(value)
	case "typeRoots":
		options.TypeRoots = parseStringArray(value)
	case "types":
		options.Types = parseStringArray(value)
	case "useDefineForClassFields":
		options.UseDefineForClassFields = parseTristate(value)
	case "useUnknownInCatchVariables":
		options.UseUnknownInCatchVariables = parseTristate(value)
	case "verbatimModuleSyntax":
		options.VerbatimModuleSyntax = parseTristate(value)
	case "maxNodeModuleJsDepth":
		options.MaxNodeModuleJsDepth = parseTristate(value)
	case "skipLibCheck":
		options.SkipLibCheck = parseTristate(value)
	case "noEmit":
		options.NoEmit = parseTristate(value)
	case "configFilePath":
		options.ConfigFilePath = parseString(value)
	case "noDtsResolution":
		options.NoDtsResolution = parseTristate(value)
	case "pathsBasePath":
		options.PathsBasePath = parseString(value)
	case "outDir":
		options.OutDir = parseString(value)
	default:
		// Handle unknown options
		fmt.Printf("Unknown option: %s\n", key)
	}

	return options
}

type tsConfigOptions struct {
	prop                map[string][]string
	compilerOptionsProp core.CompilerOptions
	references          []compiler.ProjectReference
	notDefined          string
}

func parseProjectReference(json any) []compiler.ProjectReference {
	var result []compiler.ProjectReference
	if arr, ok := json.([]interface{}); ok {
		for _, v := range arr {
			if m, ok := v.(map[string]interface{}); ok {
				var reference compiler.ProjectReference
				if v, ok := m["path"]; ok {
					reference.Path = v.(string)
				}
				if v, ok := m["originalPath"]; ok {
					reference.OriginalPath = v.(string)
				}
				if v, ok := m["circular"]; ok {
					reference.Circular = v.(bool)
				}
				result = append(result, reference)
			}
		}
	}
	return result
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
			options.prop["exclude"] = parseRawStringArray(v) //possible parseStringsArray if it is not a string?
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
			options.compilerOptionsProp = convertCompilerOptionsFromJsonWorker(v.(map[string]interface{}), basePath, errors, configFileName)
		}
	}
	return options
}
func getOptionName(option CommandLineOption) string {
	return option.Name
}

func commandLineOptionsToMap(options []CommandLineOption) map[string]CommandLineOption {
	result := make(map[string]CommandLineOption)
	for i := 0; i < len(options); i++ {
		result[getOptionName(options[i])] = options[i]
	}
	return result
}

var commandLineCompilerOptionsMapCache map[string]CommandLineOption

func getCommandLineCompilerOptionsMap() map[string]CommandLineOption {
	// if commandLineCompilerOptionsMapCache {
	// 	return commandLineCompilerOptionsMapCache
	// }
	commandLineCompilerOptionsMapCache = commandLineOptionsToMap(optionsDeclarations)
	return commandLineCompilerOptionsMapCache
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
func convertOptionsFromJson(optionsNameMap map[string]CommandLineOption, jsonOptions map[string]interface{}, basePath string, defaultOptions *core.CompilerOptions, errors []*ast.Diagnostic) core.CompilerOptions {
	if jsonOptions == nil {
		return core.CompilerOptions{}
	}
	for key, value := range jsonOptions {
		opt, ok := optionsNameMap[key]
		if ok {
			convertJson, _ := convertJsonOption(opt, value, basePath, errors, nil, nil, nil)
			*defaultOptions = parseCompilerOptions(key, convertJson)
		}
	}
	return *defaultOptions
}

func convertArrayLiteralExpressionToJson(
	elements []*ast.Expression,
	elementOption *CommandLineOption, //*commandLineOption,
	returnValue bool,
) []string {
	if !returnValue {
		for _, element := range elements {
			convertPropertyValueToJson(element, elementOption, returnValue, nil)
		}
		return nil
	}

	// Filter out invalid values
	var convertedElements []string
	for _, element := range elements {
		var convertedValue string = convertPropertyValueToJson(element, elementOption, returnValue, nil).(string)
		convertedElements = append(convertedElements, convertedValue)
	}
	filteredElements := core.Filter(convertedElements, func(v string) bool {
		return v != ""
	})
	return filteredElements
}

type optionsBaseValue struct {
	core.CompilerOptionsValue
	*ast.SourceFile
}
type optionsBase struct {
	options map[string]optionsBaseValue
}

func directoryOfCombinedPath(fileName string, basePath string) string {
	// Use the `getNormalizedAbsolutePath` function to avoid canonicalizing the path, as it must remain noncanonical
	// until consistent casing errors are reported
	return tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(fileName, basePath))
}

// ******************************************************************* //
// This is for baselineParseResult test cases and probably some parts of json without source file

/**
 * Parse the text of the tsconfig.json file
 * @param fileName The path to the config file
 * @param jsonText The text of the config file
 */
func ParseConfigFileTextToJson(fileName string, jsonText string) parsedConfigFileTextToJsonResult {
	jsonSourceFile := compiler.ParseJSONText(fileName, jsonText)
	config := convertConfigFileToObject(jsonSourceFile, jsonSourceFile.Diagnostics() /*jsonConversionNotifier*/, nil)

	var error *ast.Diagnostic
	if len(jsonSourceFile.Diagnostics()) > 0 {
		error = jsonSourceFile.Diagnostics()[0]
	} else {
		error = nil
	}
	return parsedConfigFileTextToJsonResult{config, error}
}

type VfsParseConfigHost struct {
	fs               vfs.FS
	currentDirectory string
}

func (h *VfsParseConfigHost) FS() vfs.FS {
	return h.fs
}
func ParseJsonSourceFileConfigFileContent(sourceFile *tsConfigSourceFile, host VfsParseConfigHost, basePath string, existingOptions *core.CompilerOptions, configFileName string, resolutionStack []tspath.Path, extraFileExtensions []FileExtensionInfo, extendedConfigCache *map[string]ExtendedConfigCacheEntry) module.ParsedCommandLine {
	//tracing?.push(tracing.Phase.Parse, "parseJsonSourceFileConfigFileContent", { path: sourceFile.fileName });
	result := parseJsonConfigFileContentWorker( /*json*/ nil, sourceFile, host, basePath, existingOptions, configFileName, resolutionStack, extraFileExtensions, extendedConfigCache)
	//tracing?.pop();
	return result
}

func convertObjectLiteralExpressionToJson(
	returnValue bool,
	node *ast.ObjectLiteralExpression,
	objectOption *CommandLineOption,
	jsonConversionNotifier *JsonConversionNotifier,
) map[string]any {
	fmt.Println("convertObjectLiteralExpressionToJson")
	var result map[string]any
	if returnValue {
		result = make(map[string]any)
	} else {
		result = nil
	}
	for _, element := range node.Properties.Nodes {
		fmt.Println("element", element)
		if element.Kind != ast.KindPropertyAssignment {
			//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, element, diagnostics.Property_assignment_expected));
			continue
		}

		// if (element.questionToken) { //related to isGrammarError for PropertyAssignment
		//     errors.push(createDiagnosticForNodeInSourceFile(sourceFile, element.questionToken, diagnostics.The_0_modifier_can_only_be_used_in_TypeScript_files, "?"));
		// }
		if isDoubleQuotedString(element.Name()) == false {
			//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, element.name, diagnostics.String_literal_with_double_quotes_expected));
		}

		var textOfKey any
		if compiler.IsComputedNonLiteralName(element.Name()) {
			textOfKey = nil
		} else {
			textOfKey, _ = compiler.TryGetTextOfPropertyName(element.Name())
		}
		var keyText = textOfKey //&& unescapeLeadingUnderscores(textOfKey);
		var option CommandLineOption
		if keyText != nil {
			if objectOption != nil && objectOption.ElementOptions != nil {
				option = objectOption.ElementOptions[keyText.(string)]
			} else {
				option = CommandLineOption{}
				//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, element.name, diagnostics.Unknown_option_0, keyText));
			}

		}
		// todo
		//keyText = element.AsPropertyAssignment().Name().Text()                                                                       //"exclude"
		var value = convertPropertyValueToJson(element.AsPropertyAssignment().Initializer, &option, returnValue, jsonConversionNotifier) // this needs to be element.initializer need to come back
		if keyText != "undefined" {
			if returnValue {
				result[keyText.(string)] = value
			}

			// Notify key value set, if user asked for it
			if jsonConversionNotifier != nil {
				jsonConversionNotifier.onPropertySet(keyText.(string), value, *element.AsPropertyAssignment(), *objectOption, &option)
			}
		}
	}
	return result
}

/**
 * Convert the json syntax tree into the json value and report errors
 * This returns the json value (apart from checking errors) only if returnValue provided is true.
 * Otherwise it just checks the errors and returns undefined
 *
 * @internal
 */
type emptyStruct struct{}

func convertToJson(
	sourceFile *ast.SourceFile,
	rootExpression *ast.Expression,
	errors []*ast.Diagnostic,
	returnValue bool,
	jsonConversionNotifier *JsonConversionNotifier,
) any {
	if rootExpression == nil {
		if returnValue {
			return emptyStruct{}
		} else {
			return nil
		}
	}
	var jsonConversionNotifierValue *CommandLineOption
	if jsonConversionNotifier != nil {
		jsonConversionNotifierValue = &jsonConversionNotifier.rootOptions
	}
	fmt.Println("in convertToJson", rootExpression.Kind)
	return convertPropertyValueToJson(rootExpression, jsonConversionNotifierValue, returnValue, jsonConversionNotifier)
}

func isDoubleQuotedString(node *ast.Node) bool {
	return ast.IsStringLiteral(node) //&& isStringDoubleQuoted(node, sourceFile);
}

func convertPropertyValueToJson(valueExpression *ast.Expression, option *CommandLineOption, returnValue bool, jsonConversionNotifier *JsonConversionNotifier) any {
	fmt.Println("valueExpression    ", valueExpression.Kind)
	switch valueExpression.Kind {
	case ast.KindTrueKeyword:
		return true

	case ast.KindFalseKeyword:
		return false

	case ast.KindNullKeyword:
		return nil // eslint-disable-line no-restricted-syntax

	case ast.KindStringLiteral:
		if !isDoubleQuotedString(valueExpression) {
			//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, valueExpression, diagnostics.String_literal_with_double_quotes_expected));
		}
		return (valueExpression.AsStringLiteral()).Text

	case ast.KindNumericLiteral:
		return valueExpression.AsNumericLiteral().Text

	case ast.KindPrefixUnaryExpression:
		if valueExpression.AsPrefixUnaryExpression().Operator != ast.KindMinusToken || valueExpression.AsPrefixUnaryExpression().Operand.Kind != ast.KindNumericLiteral {
			break // not valid JSON syntax
		}
		return (valueExpression.AsPrefixUnaryExpression().Operand).AsNumericLiteral().Text

	case ast.KindObjectLiteralExpression:
		objectLiteralExpression := valueExpression.AsObjectLiteralExpression()

		// Currently having element option declaration in the tsconfig with type "object"
		// determines if it needs onSetValidOptionKeyValueInParent callback or not
		// At moment there are only "compilerOptions", "typeAcquisition" and "typingOptions"
		// that satisfies it and need it to modify options set in them (for normalizing file paths)
		// vs what we set in the json
		// If need arises, we can modify this interface and callbacks as needed
		return convertObjectLiteralExpressionToJson(returnValue, objectLiteralExpression, option, jsonConversionNotifier)
	case ast.KindArrayLiteralExpression:
		fmt.Println("array literal expression")
		result := convertArrayLiteralExpressionToJson(
			(valueExpression.AsArrayLiteralExpression()).Elements.Nodes,
			nil, //option && (option.(CommandLineOptionOfListType)).element,
			returnValue,
		)
		return result
	}
	// Not in expected format
	// if option {
	// 	//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, valueExpression, diagnostics.Compiler_option_0_requires_a_value_of_type_1, option.name, getCompilerOptionValueTypeString(option)));
	// } else {
	// 	//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, valueExpression, diagnostics.Property_value_can_only_be_string_literal_numeric_literal_true_false_null_object_literal_or_array_literal));
	// }
	return nil
}

// ******************************************************************* //
/**
 * Parse the contents of a config file (tsconfig.json).
 * @param jsonNode The contents of the config file to parse
 * @param host Instance of ParseConfigHost used to enumerate files in folder.
 * @param basePath A root directory to resolve relative path entries in the config
 *    file to. e.g. outDir
 */
func ParseJsonConfigFileContent(json map[string]interface{}, host VfsParseConfigHost, basePath string, existingOptions *core.CompilerOptions, configFileName string, resolutionStack []tspath.Path, extraFileExtensions []FileExtensionInfo, extendedConfigCache *map[string]ExtendedConfigCacheEntry) module.ParsedCommandLine {
	result := parseJsonConfigFileContentWorker(json /*sourceFile*/, nil, host, basePath, existingOptions, configFileName, resolutionStack, extraFileExtensions, extendedConfigCache)
	return result
}

/**
 * Convert the json syntax tree into the json value
 */
func convertToObject(sourceFile *ast.SourceFile, errors []*ast.Diagnostic) any {
	var rootExpression *ast.Expression
	if sourceFile.Statements != nil {
		rootExpression = sourceFile.Statements.Nodes[0].AsExpressionStatement().Expression
	}
	return convertToJson(sourceFile, rootExpression, errors /*returnValue*/, true /*jsonConversionNotifier*/, nil)
}

func getDefaultCompilerOptions(configFileName string) core.CompilerOptions {
	var options core.CompilerOptions
	if configFileName != "" && tspath.GetBaseFileName(configFileName) == "jsconfig.json" {
		options = core.CompilerOptions{
			AllowJs:                      2,
			MaxNodeModuleJsDepth:         2,
			AllowSyntheticDefaultImports: 2,
			SkipLibCheck:                 2,
			NoEmit:                       2,
		}
	}
	return options
}

type defaultOptions struct {
	core.CompilerOptions
	//TypeAcquisition
	//WatchOptions
}

type propFromRaw string

const (
	files           propFromRaw = "files"
	include         propFromRaw = "include"
	exclude         propFromRaw = "exclude"
	extends         propFromRaw = "extends"
	compilerOptions propFromRaw = "compilerOptions"
	references      propFromRaw = "references"
	noProp          propFromRaw = "no-prop"
)

func convertCompilerOptionsFromJsonWorker(jsonOptions map[string]interface{}, basePath string, errors []*ast.Diagnostic, configFileName string) core.CompilerOptions {
	options := getDefaultCompilerOptions(configFileName)
	convertOptionsFromJson(getCommandLineCompilerOptionsMap(), jsonOptions, basePath, &options, errors)
	if configFileName != "" {
		options.ConfigFilePath = tspath.NormalizeSlashes(configFileName)
	}
	return options
}

func parseOwnConfigOfJson(
	json map[string]interface{},
	host VfsParseConfigHost,
	basePath string,
	configFileName string,
	errors []*ast.Diagnostic,
) *ParsedTsconfig {
	if json["excludes"] != nil {
		errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Unknown_option_excludes_Did_you_mean_exclude))
	}
	var options core.CompilerOptions
	for k, v := range json {
		if k == "compilerOptions" {
			options = convertCompilerOptionsFromJsonWorker(v.(map[string]interface{}), basePath, errors, configFileName)
		}
	}
	// typeAcquisition := convertTypeAcquisitionFromJsonWorker(json.typeAcquisition, basePath, errors, configFileName)
	// watchOptions := convertWatchOptionsFromJsonWorker(json.watchOptions, basePath, errors)
	// json.compileOnSave = convertCompileOnSaveOptionFromJson(json, basePath, errors)
	// var extendedConfigPath string
	// if json.extends != nil || json.extends == "" {
	// 	extendedConfigPath = getExtendsConfigPathOrArray(json.extends, host, basePath, configFileName, errors)
	// }
	var parsedConfig = &ParsedTsconfig{
		raw: json,
		options: &core.CompilerOptions{
			Lib: options.Lib,
		},
	}
	return parsedConfig
}

/**
 * This *just* extracts options/include/exclude/files out of a config file.
 * It does *not* resolve the included files.
 */
func parseConfig(
	json map[string]interface{},
	sourceFile *tsConfigSourceFile,
	host VfsParseConfigHost,
	basePath string,
	configFileName string,
	resolutionStack []string,
	errors []*ast.Diagnostic,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
) (*ParsedTsconfig, []*ast.Diagnostic) {
	basePath = tspath.NormalizeSlashes(basePath)
	resolvedPath := tspath.GetNormalizedAbsolutePath(configFileName, basePath)

	if slices.Contains(resolutionStack, resolvedPath) {
		var result *ParsedTsconfig
		errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Circularity_detected_while_resolving_configuration_Colon_0))
		if json != nil {
			result = &ParsedTsconfig{raw: json}
		} else {
			result = &ParsedTsconfig{raw: convertToObject(sourceFile.sourceFile, errors)}
		}
		return result, errors
	}
	var ownConfig *ParsedTsconfig
	if json != nil {
		ownConfig = parseOwnConfigOfJson(json, host, basePath, configFileName, errors)
	} else {
		ownConfig, errors = parseOwnConfigOfJsonSourceFile(sourceFile, host, basePath, &configFileName, errors)
	}

	if ownConfig.options != nil && ownConfig.options.Paths != nil {
		// If we end up needing to resolve relative paths from 'paths' relative to
		// the config file location, we'll need to know where that config file was.
		// Since 'paths' can be inherited from an extended config in another directory,
		// we wouldn't know which directory to use unless we store it here.
		ownConfig.options.PathsBasePath = basePath
	}
	if ownConfig.extendedConfigPath != nil {
		// copy the resolution stack so it is never reused between branches in potential diamond-problem scenarios.
		resolutionStack = append(resolutionStack, resolvedPath) //resolutionStack.concat([resolvedPath]); //here
		var result = extendsResult{
			options: core.CompilerOptions{},
		}
		// if compiler.IsString(ownConfig.extendedConfigPath) {
		// 	applyExtendedConfig(result, *ownConfig.extendedConfigPath)
		// } else {
		// for _, extendedConfigPath := range *ownConfig.extendedConfigPath {
		// 	//applyExtendedConfig(result, []string{extendedConfigPath})
		// }
		if result.include != nil {
			ownConfig.raw = result.include //ownConfig.raw.include = result.include
		}
		if result.exclude != nil {
			ownConfig.raw = result.exclude // ownConfig.raw.exclude = result.exclude
		}
		if result.files != nil {
			ownConfig.raw = result.files //ownConfig.raw.files = result.files
		}

		if ownConfig.raw == nil && result.compileOnSave != nil { //ownConfig.raw.compileOnSave == nil && result.compileOnSave
			ownConfig.raw = result.compileOnSave // ownConfig.raw.compileOnSave = result.compileOnSave
		}
		if sourceFile != nil && result.extendedSourceFiles != nil {
			//sourceFile.extendedSourceFiles = arrayFrom(result.extendedSourceFiles.keys()) //todo extendedSourceFile does not exist in sourcefile
		}

		// ownConfig.options = assign(result.options, ownConfig.options);
		// ownConfig.watchOptions = ownConfig.watchOptions && result.watchOptions ?
		//     assignWatchOptions(result, ownConfig.watchOptions) :
		//     ownConfig.watchOptions || result.watchOptions;
	}

	return ownConfig, errors
}

// func handleOptionConfigDirTemplateSubstitution() {
// 	setOptionValue := func(option CommandLineOption, value core.CompilerOptionsValue) {
//         (result ??= assign({}, options))[option.name] = value;
//     }
// }

const defaultIncludeSpec = "**/*"

/**
 * Parse the contents of a config file from json or json source file (tsconfig.json).
 * @param json The contents of the config file to parse
 * @param sourceFile sourceFile corresponding to the Json
 * @param host Instance of ParseConfigHost used to enumerate files in folder.
 * @param basePath A root directory to resolve relative path entries in the config
 *    file to. e.g. outDir
 * @param resolutionStack Only present for backwards-compatibility. Should be empty.
 */
func parseJsonConfigFileContentWorker(
	json map[string]interface{},
	sourceFile *tsConfigSourceFile,
	host VfsParseConfigHost,
	basePath string,
	existingOptions *core.CompilerOptions, //should default to an empty object
	configFileName string,
	resolutionStack []tspath.Path,
	extraFileExtensions []FileExtensionInfo,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
) module.ParsedCommandLine {
	//Debug.assert((json === undefined && sourceFile !== undefined) || (json !== undefined && sourceFile === undefined));
	var errors []*ast.Diagnostic
	resolutionStackString := []string{}
	// sf := &tsConfigSourceFile{}
	// sf.sourceFile = sourceFile.sourceFile
	parsedConfig, errors := parseConfig(json, sourceFile, host, basePath, configFileName, resolutionStackString, errors, extendedConfigCache)
	// const options = handleOptionConfigDirTemplateSubstitution(
	// 	extend(existingOptions, parsedConfig.options), //function in core.ts
	// 	configDirTemplateSubstitutionOptions,
	// 	basePath,
	// )
	options := parsedConfig.options
	rawConfig := ParseRawConfig(parsedConfig.raw, basePath, errors, configFileName)
	var basePathForFileNames string
	if configFileName != "" {
		rawConfig.compilerOptionsProp.ConfigFilePath = tspath.NormalizeSlashes(configFileName)
		basePathForFileNames = tspath.NormalizePath(directoryOfCombinedPath(configFileName, basePath))
	} else {
		basePathForFileNames = tspath.NormalizePath(basePath)
	}

	getPropFromRaw := func(prop propFromRaw, validateElement func(value string) bool) []string {
		value, exists := rawConfig.prop[string(prop)]
		if exists {
			if len(value) >= 0 {
				result := rawConfig.prop[string(prop)]
				if sourceFile == nil && !core.Every(result, validateElement) {
					errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop)) // , elementTypeName
				}
				return result
			} else {
				errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop, "Array"))
			}
		}
		return []string{string(noProp)}
	}

	getConfigFileSpecs := func() configFileSpecs {
		referencesOfRaw := getPropFromRaw(references, func(element string) bool { return element == "object" })
		fileSpecs := getPropFromRaw(files, func(element string) bool { return reflect.TypeOf(element).Kind() == reflect.String })
		if len(fileSpecs) == 0 || fileSpecs[0] != "no-prop" {
			hasZeroOrNoReferences := false
			if len(referencesOfRaw) == 0 || referencesOfRaw[0] == "no-prop" {
				hasZeroOrNoReferences = true
			}
			hasExtends := rawConfig.prop[string(extends)]

			if len(fileSpecs) == 0 && hasZeroOrNoReferences && hasExtends == nil {
				if sourceFile != nil {
					var fileName string
					if configFileName != "" {
						fileName = configFileName
					} else {
						fileName = "tsconfig.json"
					}
					diagnosticMessage := diagnostics.The_files_list_in_config_file_0_is_empty
					nodeValue := compiler.ForEachTsConfigPropArray(sourceFile.sourceFile, "files", func(property ast.PropertyAssignment) *ast.Node { return property.Initializer })
					errors = append(errors, ast.NewCompilerDiagnostic(diagnosticMessage, fileName, nodeValue))
				} else {
					errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.The_files_list_in_config_file_0_is_empty))
				}
			}
		}

		includeSpecs := getPropFromRaw(include, func(element string) bool { return reflect.TypeOf(element).Kind() == reflect.String })

		excludeSpecs := getPropFromRaw(exclude, func(element string) bool { return reflect.TypeOf(element).Kind() == reflect.String })
		isDefaultIncludeSpec := false
		if len(excludeSpecs) != 0 && excludeSpecs[0] == "no-prop" {
			outDir := options.OutDir
			declarationDir := options.DeclarationDir

			if outDir != "" || declarationDir != "" {
				excludeSpecs = core.Filter([]string{outDir, declarationDir}, func(d string) bool { return d != "" })
			}
		}

		if len(fileSpecs) != 0 && fileSpecs[0] == "no-prop" && len(includeSpecs) != 0 && includeSpecs[0] == "no-prop" {
			includeSpecs = []string{defaultIncludeSpec}
			isDefaultIncludeSpec = true
		}
		var validatedIncludeSpecsBeforeSubstitution []string
		var validatedExcludeSpecsBeforeSubstitution []string
		var validatedIncludeSpecs []string
		var validatedExcludeSpecs []string

		// The exclude spec list is converted into a regular expression, which allows us to quickly
		// test whether a file or directory should be excluded before recursively traversing the
		// file system.

		if len(includeSpecs) != 0 && includeSpecs[0] != "no-prop" {
			validatedIncludeSpecsBeforeSubstitution = validateSpecs(includeSpecs, errors /*disallowTrailingRecursion*/, true, sourceFile, "include")
			validatedIncludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
				validatedIncludeSpecsBeforeSubstitution,
				basePathForFileNames,
			)
			if validatedIncludeSpecs == nil {
				validatedIncludeSpecs = validatedIncludeSpecsBeforeSubstitution
			}
		}

		if len(excludeSpecs) != 0 && excludeSpecs[0] != "no-prop" {
			validatedExcludeSpecsBeforeSubstitution = validateSpecs(excludeSpecs, errors /*disallowTrailingRecursion*/, false, sourceFile, "exclude")
			validatedExcludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
				validatedExcludeSpecsBeforeSubstitution,
				basePathForFileNames,
			)
			if validatedExcludeSpecs == nil {
				validatedExcludeSpecs = validatedExcludeSpecsBeforeSubstitution
			}
		}

		validatedFilesSpecBeforeSubstitution := core.Filter(fileSpecs, func(spec string) bool { return reflect.TypeOf(spec).Kind() == reflect.String })
		validatedFilesSpec := getSubstitutedStringArrayWithConfigDirTemplate(
			validatedFilesSpecBeforeSubstitution,
			basePathForFileNames,
		)
		if validatedFilesSpec == nil && len(validatedFilesSpecBeforeSubstitution) != 0 && validatedFilesSpecBeforeSubstitution[0] != "no-prop" {
			validatedFilesSpec = validatedFilesSpecBeforeSubstitution
		}

		return configFileSpecs{
			fileSpecs,
			includeSpecs,
			excludeSpecs,
			validatedFilesSpec,
			validatedIncludeSpecs,
			validatedExcludeSpecs,
			validatedFilesSpecBeforeSubstitution,
			validatedIncludeSpecsBeforeSubstitution,
			validatedExcludeSpecsBeforeSubstitution,
			isDefaultIncludeSpec,
		}
	}

	configFileSpecs := getConfigFileSpecs()

	getFileNames := func(basePath string) []string {
		fileNames := getFileNamesFromConfigSpecs(configFileSpecs, basePath, options, host.fs, extraFileExtensions)
		if shouldReportNoInputFiles(fileNames, canJsonReportNoInputFiles(rawConfig), resolutionStack) {
			errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.No_inputs_were_found_in_config_file_0_Specified_include_paths_were_1_and_exclude_paths_were_2, configFileName, configFileSpecs.includeSpecs, configFileSpecs.excludeSpecs))
		}
		return fileNames
	}
	return module.ParsedCommandLine{
		Options:   options,
		FileNames: getFileNames(basePathForFileNames),
		Raw:       parsedConfig.raw,
		Errors:    errors,
	}

}

func canJsonReportNoInputFiles(rawConfig tsConfigOptions) bool {
	_, filesExists := rawConfig.prop["files"]
	_, referencesExists := rawConfig.prop["references"]
	return !filesExists && !referencesExists
}

func shouldReportNoInputFiles(fileNames []string, canJsonReportNoInutFiles bool, resolutionStack []tspath.Path) bool {
	return len(fileNames) == 0 && canJsonReportNoInutFiles && (resolutionStack != nil || len(resolutionStack) == 0)
}

func validateSpecs(specs []string, errors []*ast.Diagnostic, disallowTrailingRecursion bool, jsonSourceFile *tsConfigSourceFile, specKey string) []string {
	createDiagnostic := func(message *diagnostics.Message, spec string) *ast.Diagnostic {
		element := getTsConfigPropArrayElementValue(jsonSourceFile, specKey, spec)
		return ast.NewCompilerDiagnostic(message, element)
	}

	return core.Filter(specs, func(spec string) bool {
		if spec == "" {
			return false
		}
		diag, _ := specToDiagnostic(spec, disallowTrailingRecursion)
		if diag != nil {
			errors = append(errors, createDiagnostic(diag, spec))
			// errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(jsonSourceFile, getTsConfigPropArrayElementValue(jsonSourceFile, specKey, spec), diag.message, spec));
		}
		return diag == nil
	})

}

// // func createDiagnostic(message, spec string) Diagnostic {
// // 	element := getTsConfigPropArrayElementValue(nil, "", spec)
// // 	return createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(nil, element, message, spec)
// // }

func specToDiagnostic(spec string, disallowTrailingRecursion bool) (*diagnostics.Message, string) {
	if disallowTrailingRecursion {
		if ok, _ := regexp.MatchString(invalidTrailingRecursionPattern, spec); ok {
			return diagnostics.File_specification_cannot_end_in_a_recursive_directory_wildcard_Asterisk_Asterisk_Colon_0, spec
		}
	} else if invalidDotDotAfterRecursiveWildcard(spec) {
		return diagnostics.File_specification_cannot_contain_a_parent_directory_that_appears_after_a_recursive_directory_wildcard_Asterisk_Asterisk_Colon_0, spec
	}
	return nil, ""
}

func invalidDotDotAfterRecursiveWildcard(s string) bool {
	// We used to use the regex /(^|\/)\*\*\/(.*\/)?\.\.($|\/)/ to check for this case, but
	// in v8, that has polynomial performance because the recursive wildcard match - **/ -
	// can be matched in many arbitrary positions when multiple are present, resulting
	// in bad backtracking (and we don't care which is matched - just that some /.. segment
	// comes after some **/ segment).
	var wildcardIndex int
	if strings.HasPrefix(s, "**/") {
		wildcardIndex = 0
	} else {
		wildcardIndex = strings.Index(s, "/**/")
	}
	if wildcardIndex == -1 {
		return false
	}
	var lastDotIndex int
	if strings.HasSuffix(s, "/..") {
		lastDotIndex = len(s)
	} else {
		lastDotIndex = strings.LastIndex(s, "/../")
	}
	return lastDotIndex > wildcardIndex
}

/**
 * Tests for a path that ends in a recursive directory wildcard.
 * Matches **, \**, **\, and \**\, but not a**b.
 *
 * NOTE: used \ in place of / above to avoid issues with multiline comments.
 *
 * Breakdown:
 *  (^|\/)      # matches either the beginning of the string or a directory separator.
 *  \*\*        # matches the recursive directory wildcard "**".
 *  \/?$        # matches an optional trailing directory separator at the end of the string.
 */
const invalidTrailingRecursionPattern = `(?:^|\/)\*\*\/?$`

func getTsConfigPropArrayElementValue(tsConfigSourceFile *tsConfigSourceFile, propKey string, elementValue string) *ast.StringLiteral {
	return forEachTsConfigPropArray(tsConfigSourceFile, propKey, func(property ast.PropertyAssignment) *ast.StringLiteral {
		if ast.IsArrayLiteralExpression(property.Initializer) {
			t := core.Find(property.Initializer.AsArrayLiteralExpression().Elements.Nodes, func(element *ast.Node) bool {
				return ast.IsStringLiteral(element) && element.AsStringLiteral().Text == elementValue
			}).AsStringLiteral()
			return t
		}
		return nil
	})
}

// if (isStringLiteral(property.initializer) && property.initializer.text === elementValue) {
// 	return property.initializer;
// }

// })
// return forEachTsConfigPropArray(tsConfigSourceFile, propKey, property =>
//     isArrayLiteralExpression(property.initializer) ?
//         find(property.initializer.elements, (element): element is StringLiteral => isStringLiteral(element) && element.text === elementValue) :
//         undefined);

func forEachTsConfigPropArray[T any](tsConfigSourceFile *tsConfigSourceFile, propKey string, callback func(property ast.PropertyAssignment) T) T {
	if tsConfigSourceFile != nil {
		return forEachPropertyAssignment(*getTsConfigObjectLiteralExpression(tsConfigSourceFile), propKey, callback)
	}
	return interface{}(nil).(T)
}
func forEachPropertyAssignment[T any](objectLiteral ast.ObjectLiteralExpression, key string, callback func(property ast.PropertyAssignment) T, key2 ...string) T {
	if objectLiteral != (ast.ObjectLiteralExpression{}) {
		for _, property := range objectLiteral.Properties.Nodes {
			if !ast.IsPropertyAssignment(property) {
				continue
			}
			if propName, ok := compiler.TryGetTextOfPropertyName(property.Name()); ok {
				if propName == key || (len(key2) > 0 && key2[0] == propName) {
					result := callback(*property.AsPropertyAssignment())
					return result
				}
			}
		}
	}
	return interface{}(nil).(T)
}

func getTsConfigObjectLiteralExpression(tsConfigSourceFile *tsConfigSourceFile) *ast.ObjectLiteralExpression {
	if tsConfigSourceFile != nil && tsConfigSourceFile.sourceFile.Statements != nil && len(tsConfigSourceFile.sourceFile.Statements.Nodes) > 0 {
		expression := tsConfigSourceFile.sourceFile.Statements.Nodes[0].AsExpressionStatement().Expression
		return expression.AsObjectLiteralExpression()
	}
	return nil
}

func getSubstitutedPathWithConfigDirTemplate(value string, basePath string) string {
	return tspath.GetNormalizedAbsolutePath(strings.ReplaceAll(value, configDirTemplate, "./"), basePath)
}
func getSubstitutedStringArrayWithConfigDirTemplate(list []string, basePath string) []string {
	if list == nil {
		return nil
	}
	var result []string
	for _, element := range list {
		if !startsWithConfigDirTemplate(element) {
			return nil
		} else {
			result = append(result, getSubstitutedPathWithConfigDirTemplate(element, basePath))
		}
	}
	return result
}

/**
 * Determines whether a literal or wildcard file has already been included that has a higher
 * extension priority.
 *
 * @param file The path to the file.
 */
func hasFileWithHigherPriorityExtension(file string, literalFiles map[string]string, wildcardFiles map[string]string, extensions [][]string, keyMapper func(value string) string) bool {
	var extensionGroup []string
	for _, group := range extensions {
		if tspath.FileExtensionIsOneOf(file, group) {
			extensionGroup = append(extensionGroup, group...)
		}
	}
	if extensionGroup == nil {
		return false
	}
	for _, ext := range extensionGroup {
		// d.ts files match with .ts extension and with case sensitive sorting the file order for same files with ts tsx and dts extension is
		// d.ts, .ts, .tsx in that order so we need to handle tsx and dts of same same name case here and in remove files with same extensions
		// So dont match .d.ts files with .ts extension
		if tspath.FileExtensionIs(file, ext) && (ext != tspath.ExtensionTs || !tspath.FileExtensionIs(file, tspath.ExtensionDts)) {
			return false
		}
		higherPriorityPath := keyMapper(tspath.ChangeExtension(file, extensionGroup[0]))
		if literalFiles[higherPriorityPath] != "" || wildcardFiles[higherPriorityPath] != "" {
			if ext == tspath.ExtensionDts && (tspath.FileExtensionIs(file, tspath.ExtensionJs) || tspath.FileExtensionIs(file, tspath.ExtensionJsx)) {
				// LEGACY BEHAVIOR: An off-by-one bug somewhere in the extension priority system for wildcard module loading allowed declaration
				// files to be loaded alongside their js(x) counterparts. We regard this as generally undesirable, but retain the behavior to
				// prevent breakage.
				continue
			}
			return true
		}
	}
	return false
}

/**
 * Removes files included via wildcard expansion with a lower extension priority that have
 * already been included.
 *
 * @param file The path to the file.
 */
func removeWildcardFilesWithLowerPriorityExtension(file string, wildcardFiles map[string]string, extensions [][]string, keyMapper func(value string) string) {
	var extensionGroup []string
	for _, group := range extensions {
		if tspath.FileExtensionIsOneOf(file, group) {
			extensionGroup = append(extensionGroup, group...)
		}
	}
	if extensionGroup == nil {
		return
	}

	for i := len(extensionGroup) - 1; i >= 0; i-- {
		ext := extensionGroup[i]
		if tspath.FileExtensionIs(file, ext) {
			return
		}
		lowerPriorityPath := keyMapper(tspath.ChangeExtension(file, ext))
		delete(wildcardFiles, lowerPriorityPath)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

/**
 * Gets the file names from the provided config file specs that contain, files, include, exclude and
 * other properties needed to resolve the file names
 * @param configFileSpecs The config file specs extracted with file names to include, wildcards to include/exclude and other details
 * @param basePath The base path for any relative file specifications.
 * @param options Compiler options.
 * @param host The host used to resolve files and directories.
 * @param extraFileExtensions optionaly file extra file extension information from host
 *
 * @internal
 */
func getFileNamesFromConfigSpecs(
	configFileSpecs configFileSpecs,
	basePath string, //considering this is the current directory
	options *core.CompilerOptions,
	host vfs.FS,
	extraFileExtensions []FileExtensionInfo,
) []string {
	extraFileExtensions = []FileExtensionInfo{}
	basePath = tspath.NormalizePath(basePath)

	// Literal file names (provided via the "files" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map later when when including
	// wildcard paths.
	var literalFileMap = make(map[string]string)

	// Wildcard paths (provided via the "includes" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map to store paths matched
	// via wildcard, and to handle extension priority.
	var wildcardFileMap = make(map[string]string)

	// Wildcard paths of json files (provided via the "includes" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map to store paths matched
	// via wildcard of *.json kind
	var wildCardJsonFileMap = make(map[string]string)
	validatedFilesSpec := configFileSpecs.validatedFilesSpec
	validatedIncludeSpecs := configFileSpecs.validatedIncludeSpecs
	validatedExcludeSpecs := configFileSpecs.validatedExcludeSpecs

	// Rather than re-query this for each file and filespec, we query the supported extensions
	// once and store it on the expansion context.
	var supportedExtensions = getSupportedExtensions(options, extraFileExtensions)
	var supportedExtensionsWithJsonIfResolveJsonModule = getSupportedExtensionsWithJsonIfResolveJsonModule(options, supportedExtensions)

	// Literal files are always included verbatim. An "include" or "exclude" specification cannot
	// remove a literal file.
	if validatedFilesSpec != nil {
		for _, fileName := range validatedFilesSpec {
			file := tspath.GetNormalizedAbsolutePath(fileName, basePath)
			literalFileMap[tspath.GetCanonicalFileName(fileName, host.UseCaseSensitiveFileNames())] = file
		}
	}

	var jsonOnlyIncludeRegexes []*regexp2.Regexp
	if validatedIncludeSpecs != nil && len(validatedIncludeSpecs) > 0 { // In place of process.cwd, I'm doing basePath which is the current directory
		files := compiler.ReadDirectory(host, basePath, basePath, core.Flatten(supportedExtensionsWithJsonIfResolveJsonModule), validatedExcludeSpecs, validatedIncludeSpecs, -1)
		for _, file := range files {
			if tspath.FileExtensionIs(file, tspath.ExtensionJson) {
				if jsonOnlyIncludeRegexes != nil {
					includes := core.Filter(validatedIncludeSpecs, func(include string) bool { return strings.HasSuffix(include, tspath.ExtensionJson) })
					var includeFilePatterns []string = core.Map(compiler.GetRegularExpressionsForWildcards(includes, basePath, "files"), func(pattern string) string { return fmt.Sprintf("^%s$", pattern) })
					if includeFilePatterns != nil {
						jsonOnlyIncludeRegexes = core.Map(includeFilePatterns, func(pattern string) *regexp2.Regexp {
							return compiler.GetRegexFromPattern(pattern, host.UseCaseSensitiveFileNames())
						})
					} else {
						jsonOnlyIncludeRegexes = nil
					}
					includeIndex := core.FindIndex(jsonOnlyIncludeRegexes, func(re *regexp2.Regexp) bool { return must(re.MatchString(file)) })
					if includeIndex != -1 {
						key := tspath.GetCanonicalFileName(file, host.UseCaseSensitiveFileNames())
						if literalFileMap[key] != "" && wildCardJsonFileMap[key] != "" {
							wildCardJsonFileMap[key] = file
						}
					}
					continue
				}
			}
			// If we have already included a literal or wildcard path with a
			// higher priority extension, we should skip this file.
			//
			// This handles cases where we may encounter both <file>.ts and
			// <file>.d.ts (or <file>.js if "allowJs" is enabled) in the same
			// directory when they are compilation outputs.
			if hasFileWithHigherPriorityExtension(file, literalFileMap, wildcardFileMap, supportedExtensions, func(value string) string {
				return tspath.GetCanonicalFileName(value, host.UseCaseSensitiveFileNames())
			}) {
				continue
			}

			// We may have included a wildcard path with a lower priority
			// extension due to the user-defined order of entries in the
			// "include" array. If there is a lower priority extension in the
			// same directory, we should remove it.
			removeWildcardFilesWithLowerPriorityExtension(file, wildcardFileMap, supportedExtensions, func(value string) string {
				return tspath.GetCanonicalFileName(value, host.UseCaseSensitiveFileNames())
			})
			key := tspath.GetCanonicalFileName(file, host.UseCaseSensitiveFileNames())
			if literalFileMap[key] == "" && wildcardFileMap[key] == "" {

				wildcardFileMap[key] = file
			}
		}
	}
	var literalFiles []string
	for _, file := range literalFileMap {
		literalFiles = append(literalFiles, file)
	}
	var wildcardFiles []string
	for _, file := range wildcardFileMap {
		wildcardFiles = append(wildcardFiles, file)
	}
	var wildCardJsonFiles []string
	for _, file := range wildCardJsonFileMap {
		wildCardJsonFiles = append(wildCardJsonFiles, file)
	}
	return slices.Concat(literalFiles, wildcardFiles, wildCardJsonFiles)
}

var allSupportedExtensions = [][]string{{tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionDts, tspath.ExtensionJs, tspath.ExtensionJsx}, {tspath.ExtensionCts, tspath.ExtensionDcts, tspath.ExtensionCjs}, {tspath.ExtensionMts, tspath.ExtensionDmts, tspath.ExtensionMjs}}
var supportedTSExtensions = [][]string{{tspath.ExtensionTs, tspath.ExtensionTsx, tspath.ExtensionDts}, {tspath.ExtensionCts, tspath.ExtensionDcts}, {tspath.ExtensionMts, tspath.ExtensionDmts}}
var allSupportedExtensionsWithJson = [][]string(slices.Concat(allSupportedExtensions, ([][]string{{tspath.ExtensionJson}})))
var supportedTSExtensionsWithJson = [][]string(slices.Concat(supportedTSExtensions, ([][]string{{tspath.ExtensionJson}})))

func getAllowJSCompilerOption(compilerOptions *core.CompilerOptions) core.Tristate {
	return core.ComputedOptions["allowJs"].ComputeValue(compilerOptions).(core.Tristate)
}
func getResolveJsonModule(compilerOptions *core.CompilerOptions) bool {
	return core.ComputedOptions["resolveJsonModule"].ComputeValue(compilerOptions).(bool)
}
func getSupportedExtensions(options *core.CompilerOptions, extraFileExtensions []FileExtensionInfo) [][]string {
	needJsExtensions := getAllowJSCompilerOption(options) == 2

	if extraFileExtensions == nil || len(extraFileExtensions) == 0 {
		if needJsExtensions {
			return allSupportedExtensions
		} else {
			return supportedTSExtensions
		}
	}
	var builtins [][]string
	if needJsExtensions {
		builtins = allSupportedExtensions
	} else {
		builtins = supportedTSExtensions
	}
	var flatBuiltins = core.Flatten(builtins)
	result := core.Map(extraFileExtensions, func(x FileExtensionInfo) []string {
		if x.scriptKind == core.ScriptKindDeferred || (needJsExtensions && (x.scriptKind == core.ScriptKindJS || x.scriptKind == core.ScriptKindJSX) && !slices.Contains(flatBuiltins, x.extension)) {
			return []string{x.extension}
		}
		return nil
	})
	var extensions = slices.Concat(builtins, result)
	return extensions
}

func getSupportedExtensionsWithJsonIfResolveJsonModule(options *core.CompilerOptions, supportedExtensions [][]string) [][]string {
	if options != nil || !getResolveJsonModule(options) {
		return supportedExtensions
	}
	compareExtensions := func(a, b [][]string) bool {
		if len(a) != len(b) {
			return false
		}

		for i := range a {
			if !slices.Equal(a[i], b[i]) {
				return false
			}
		}
		return true
	}
	if compareExtensions(supportedExtensions, allSupportedExtensions) {
		return allSupportedExtensionsWithJson

	}
	if compareExtensions(supportedExtensions, supportedTSExtensions) {
		return supportedTSExtensionsWithJson
	}
	return [][]string(slices.Concat(supportedExtensions, ([][]string{{tspath.ExtensionJson}})))
}
