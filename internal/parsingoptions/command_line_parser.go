package parsingoptions

import (

	//"slices"

	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/compiler/module"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
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
type CommandLineOptionBase struct { //???
	name                      string
	commandLineOptionBasetype commandLineOptionBaseType //type "string" | "number" | "boolean" | "object" | "list" | "listOrElement" | Map<string, number | string>;    // a value of a primitive type, or an object literal mapping named values to actual values
	isFilePath                *bool                     // True if option value is a path or fileName
	shortName                 *string                   // A short mnemonic for convenience - for instance, 'h' can be used in place of 'help'
	description               *diagnostics.Message      //DM			                       // The message describing what the command line switch does.
	//defaultValueDescription?: string | number | boolean | DiagnosticMessage | undefined;   // The message describing what the dafault value is. string type is prepared for fixed chosen like "false" which do not need I18n.
	paramType                  *diagnostics.Message //DM                       // The name to be used for a non-boolean option's parameter
	isTSConfigOnly             *bool                // True if option can only be specified via tsconfig.json file
	isCommandLineOnly          *bool
	showInSimplifiedHelpView   *bool
	category                   *diagnostics.Message //DM
	strictFlag                 *bool                // true if the option is one of the flag under strict
	allowJsFlag                *bool
	affectsSourceFile          *bool // true if we should recreate SourceFiles after this option changes
	affectsModuleResolution    *bool // currently same effect as `affectsSourceFile`
	affectsBinddiagnostics     *bool // true if this affects binding (currently same effect as `affectsSourceFile`)
	affectsSemanticdiagnostics *bool // true if option affects semantic diagnostics
	affectsEmit                *bool // true if the options affects emit
	affectsProgramStructure    *bool // true if program should be reconstructed from root files if option changes and does not affect module resolution as affectsModuleResolution indirectly means program needs to reconstructed
	affectsDeclarationPath     *bool // true if the options affects declaration file path computed
	affectsBuildInfo           *bool // true if this options should be emitted in buildInfo
	transpileOptionValue       *bool // If set this means that the option should be set to this value when transpiling
	//extraValidation?: (value: CompilerOptionsValue) => [DiagnosticMessage, ...string[]] | undefined; // Additional validation to be performed for the value to be valid
	disallowNullOrUndefined            *bool // If set option does not allow setting null
	allowConfigDirTemplateSubstitution *bool // If set option allows substitution of `${configDir}` in the value
}

type commandLineOptionOfTypes struct { // new - this is conbining commandLineOptionOfCustomType, commandLineOptionOfStringType, commandLineOptionOfNumberType, commandLineOptionOfBooleanType, tsConfigOnlyOption
	CommandLineOptionBase
	defaultValueDescription   defaultValueDescriptionType
	deprecatedKeys            *map[string]bool
	commandLineOptionBasetype *optionType
	elementOptions            *map[string]CommandLineOption
	//extraKeydiagnostics *DidYouMeanOptionsdiagnostics;
}

// type commandLineOptionOfCustomType struct {
// 	CommandLineOptionBase
// 	//commandLineOptionType map[string]string // an object literal mapping named values to actual values //todo was originally Map<string, number | string>;
// 	defaultValueDescription defaultValueDescriptionType
// 	deprecatedKeys          map[string]bool
// }

// type commandLineOptionOfStringType struct {
// 	CommandLineOptionBase
// 	//commandLineOptionType optionType //"string"
// 	defaultValueDescription defaultValueDescriptionType // should only be string | DiagnosticMessage;
// }

// type commandLineOptionOfNumberType struct {
// 	CommandLineOptionBase
// 	commandLineOptionBasetype optionType //"number:"
// 	defaultValueDescription   defaultValueDescriptionType
// }

// type commandLineOptionOfBooleanType struct {
// 	CommandLineOptionBase
// 	commandLineOptionBasetype optionType //"boolean";
// 	defaultValueDescription   defaultValueDescriptionType
// }

// type CommandLineOptionOfListType struct { //new changes it a little
// 	CommandLineOptionBase
// 	element                 commandLineOptionOfTypes
// 	listPreserveFalsyValues *bool
// }

// ***********************************************************************//
// this will be merged with Isabeel's pr
type CommandLineOptionKind string

const (
	CommandLineOptionTypeString        CommandLineOptionKind = "string"
	CommandLineOptionTypeNumber        CommandLineOptionKind = "number"
	CommandLineOptionTypeBoolean       CommandLineOptionKind = "boolean"
	CommandLineOptionTypeObject        CommandLineOptionKind = "object"
	CommandLineOptionTypeList          CommandLineOptionKind = "list"
	CommandLineOptionTypeListOrElement CommandLineOptionKind = "listOrElement"
	CommandLineOptionTypeEnum          CommandLineOptionKind = "enum" //map
)

type CommandLineOption struct {
	kind            CommandLineOptionKind
	name, shortName string
	paramType       diagnostics.Message
	// used in parsing
	isFilePath        bool
	isTSConfigOnly    bool
	isCommandLineOnly bool

	// used in output
	description              *diagnostics.Message
	defaultValueDescription  any
	showInSimplifiedHelpView bool

	// used in output in serializing and generate tsconfig
	category *diagnostics.Message

	// defined once
	extraValidation *func(value core.CompilerOptionsValue) (d *diagnostics.Message, args []string)

	// true or undefined
	// used for configDirTemplateSubstitutionOptions
	allowConfigDirTemplateSubstitution,
	// used for filter in compilerrunner
	affectsDeclarationPath,
	affectsProgramStructure,
	affectsSemanticdiagnostics,
	affectsBuildInfo,
	affectsBinddiagnostics,
	affectsSourceFile,
	affectsModuleResolution,
	affectsEmit,

	allowJsFlag,
	strictFlag bool

	// transpileoptions worker
	transpileOptionValue core.Tristate
	// options[option.name] = option.transpileOptionValue;

	// used in listtype
	listPreserveFalsyValues bool
	disallowNullOrUndefined bool
}

// CommandLineOption.Elements()
var commandLineOptionElements = map[string]*CommandLineOption{
	"lib": {
		name:                    "lib",
		kind:                    CommandLineOptionTypeEnum, // libMap,
		defaultValueDescription: core.TSUnknown,
	},
	"rootDirs": {
		name:       "rootDirs",
		kind:       CommandLineOptionTypeString,
		isFilePath: true,
	},
	"typeRoots": {
		name:       "typeRoots",
		kind:       CommandLineOptionTypeString,
		isFilePath: true,
	},
	"types": {
		name: "types",
		kind: CommandLineOptionTypeString,
	},
	"moduleSuffixes": {
		name: "suffix",
		kind: CommandLineOptionTypeString,
	},
	"customConditions": {
		name: "condition",
		kind: CommandLineOptionTypeString,
	},
	"plugins": {
		name: "plugin",
		kind: CommandLineOptionTypeObject,
	},
}

func (option *CommandLineOption) Elements() *CommandLineOption {
	if option.kind != CommandLineOptionTypeList && option.kind != CommandLineOptionTypeListOrElement {
		return nil
	}
	return commandLineOptionElements[option.name]
}

// ***********************************************************************//

var optionDeclarations []CommandLineOption = append(commonOptionsWithBuild, commandOptionsWithoutBuild...)

type tsConfigOnlyOption struct {
	CommandLineOptionBase
	commandLineOptionBasetype optionType //"object";
	elementOptions            *map[string]CommandLineOption
	//extraKeydiagnostics?: DidYouMeanOptionsdiagnostics;
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

// var extendsOptionDeclaration CommandLineOptionOfListType = CommandLineOptionOfListType{
// 	CommandLineOptionBase: CommandLineOptionBase{
// 		name:                      "extends",
// 		commandLineOptionBasetype: commandLineOptionBaseType{optionType: "listOrElement"},
// 		//category: compiler.diagnostics.File_Management,
// 		disallowNullOrUndefined: func(b bool) *bool { return &b }(true), //need to check this
// 	},
// 	element: commandLineOptionOfTypes{
// 		CommandLineOptionBase: CommandLineOptionBase{
// 			name:                      "extends",
// 			commandLineOptionBasetype: commandLineOptionBaseType{optionType: "string"},
// 		},
// 	},
// }

var compilerOptionsDeclaration commandLineOptionOfTypes = commandLineOptionOfTypes{
	CommandLineOptionBase: CommandLineOptionBase{
		name:                      "compilerOptions",
		commandLineOptionBasetype: commandLineOptionBaseType{optionType: "object"},
	},
	//elementOptions: getCommandLineCompilerOptionsMap(),
	//extraKeydiagnostics: compilerOptionsDidYouMeandiagnostics,
}

var typeAcquisitionDeclaration commandLineOptionOfTypes = commandLineOptionOfTypes{
	CommandLineOptionBase: CommandLineOptionBase{
		name:                      "typeAcquisition",
		commandLineOptionBasetype: commandLineOptionBaseType{optionType: "object"},
	},
	// elementOptions: getCommandLineTypeAcquisitionMap(),
	// extraKeydiagnostics: typeAcquisitionDidYouMeandiagnostics,
}
var tsconfigRootOptions commandLineOptionOfTypes //TsConfigOnlyOption

func getTsconfigRootOptionsMap() *commandLineOptionOfTypes { //TsConfigOnlyOption
	// if tsconfigRootOptions == undefined {
	//     tsconfigRootOptions = {
	//         name: undefined!, // should never be needed since this is root
	//         type: "object",
	//         elementOptions: commandLineOptionsToMap([
	//             compilerOptionsDeclaration,
	//             watchOptionsDeclaration,
	//             typeAcquisitionDeclaration,
	//             extendsOptionDeclaration,
	//             {
	//                 name: "references",
	//                 type: "list",
	//                 element: {
	//                     name: "references",
	//                     type: "object",
	//                 },
	//                 category: diagnostics.Projects,
	//             },
	//             {
	//                 name: "files",
	//                 type: "list",
	//                 element: {
	//                     name: "files",
	//                     type: "string",
	//                 },
	//                 category: diagnostics.File_Management,
	//             },
	//             {
	//                 name: "include",
	//                 type: "list",
	//                 element: {
	//                     name: "include",
	//                     type: "string",
	//                 },
	//                 category: diagnostics.File_Management,
	//                 defaultValueDescription: diagnostics.if_files_is_specified_otherwise_Asterisk_Asterisk_Slash_Asterisk,
	//             },
	//             {
	//                 name: "exclude",
	//                 type: "list",
	//                 element: {
	//                     name: "exclude",
	//                     type: "string",
	//                 },
	//                 category: diagnostics.File_Management,
	//                 defaultValueDescription: diagnostics.node_modules_bower_components_jspm_packages_plus_the_value_of_outDir_if_one_is_specified,
	//             },
	//             compileOnSaveCommandLineOption,
	//         ]),
	//     };
	// } //todo
	//return tsconfigRootOptions
	return nil
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

// var compilerOptionsDeclaration {
//     name: "compilerOptions",
//     type: "object",
//     elementOptions: getCommandLineCompilerOptionsMap(),
//     extraKeydiagnostics: compilerOptionsDidYouMeandiagnostics,
// };

type ParseConfigHost struct {
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
	scriptKind     *core.ScriptKind
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

/**
 * Parse the contents of a config file (tsconfig.json).
 * @param jsonNode The contents of the config file to parse
 * @param host Instance of ParseConfigHost used to enumerate files in folder.
 * @param basePath A root directory to resolve relative path entries in the config
 *    file to. e.g. outDir
 */
// func ParseJsonSourceFileConfigFileContent(
// 	sourceFile *ast.SourceFile,
// 	host ParseConfigHost,
// 	basePath string,
// 	existingOptions *core.CompilerOptions,
// 	configFileName *string,
// 	resolutionStack *[]tspath.Path,
// 	extraFileExtenstions *[]FileExtensionInfo,
// 	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
// 	//existingWatchOptions *compiler.WatchOptions) *compiler.ParsedCommandLine {
// )*compiler.ParsedCommandLine {
// 	//tracing?.push(tracing.Phase.Parse, "parseJsonSourceFileConfigFileContent", { path: sourceFile.fileName });
// 	result := parseJsonConfigFileContentWorker( /*json*/ nil, sourceFile, host, basePath, *existingOptions, configFileName, *resolutionStack, *extraFileExtenstions, extendedConfigCache)
// 	//tracing?.pop();
// 	return result
// }

// type readFile func(fileName string) string

// func tryReadFile(fileName string, fn readFile) { //return string | Diagnostic
// 	//text := fn(fileName)
// 	// if err != nil {
// 	// 	//return createCompilerDiagnostic(diagnostics.Cannot_read_file_0_Colon_1, fileName, err.Error())
// 	// }
// 	// if text == nil {
// 	// 	//return createCompilerDiagnostic(diagnostics.Cannot_read_file_0, fileName)
// 	// }
// 	//return text
// 	// catch (e) {
// 	//     //return createCompilerDiagnostic(diagnostics.Cannot_read_file_0_Colon_1, fileName, e.message);
// 	// }
// 	return
// }

// func getBaseFileName(path string, extensions *[]string, ignoreCase *bool) string {
// 	path = tspath.NormalizeSlashes(path)

// 	// if the path provided is itself the root, then it has not file name.
// 	rootLength := tspath.GetRootLength(path)
// 	if rootLength == len(path) {
// 		return ""
// 	}

// 	// return the trailing portion of the path starting after the last (non-terminal) directory
// 	// separator but not including any trailing directory separator.
// 	path = tspath.RemoveTrailingDirectorySeparator(path)
// 	//name :=  path[int(math.Max(float64(compiler.GetRootLength(path)),float64(strings.LastIndex(path,compiler.DirectorySeparator)+1)))]//path.slice(Math.max(getRootLength(path), path.lastIndexOf(directorySeparator) + 1));
// 	// var extension string
// 	// if extensions != nil && ignoreCase != nil {
// 	//     extension = getAnyExtensionFromPath(name, extensions, ignoreCase)
// 	// }
// 	// if extension != nil {
// 	//     return name[0:len(name) - len(extension)]
// 	// }
// 	// return name
// 	return ""
// }

// func parseOwnConfigOfJsonSourceFile(
// 	sourceFile *ast.SourceFile,
// 	host ParseConfigHost,
// 	basePath string,
// 	configFileName *string,
// 	errors []*ast.Diagnostic,
// ) *ParsedTsconfig {
// 	options := getDefaultCompilerOptions(configFileName)
// 	//var typeAcquisition *compiler.TypeAcquisition
// 	//var watchOptions *compiler.WatchOptions
// 	//var extendedConfigPath []string = []string{} // | string
// 	// var rootCompilerOptions []ast.PropertyName

// 	// rootOptions := getTsconfigRootOptionsMap()
// 	//var conversionNotifier JsonConversionNotifier

// 	// onPropertySet := func(
// 	// 	keyText string,
// 	// 	value any,
// 	// 	propertyAssignment ast.PropertyAssignment,
// 	// 	parentOption *commandLineOptionOfTypes, //TsConfigOnlyOption,
// 	// 	option *commandLineOption,
// 	// ) {
// 	// 	// Ensure value is verified except for extends which is handled in its own way for error reporting
// 	// 	if option != nil { //&& option != extendsOptionDeclaration {
// 	// 		value = convertJsonOption(*option, value, basePath, errors, &propertyAssignment, propertyAssignment.Initializer, sourceFile)
// 	// 	}
// 	// 	if parentOption.name != "" {
// 	// 		if option != nil {
// 	// 			//var currentOption
// 	// 			if parentOption == &compilerOptionsDeclaration {
// 	// 				// currentOption := options
// 	// 				// } else if parentOption == &watchOptionsDeclaration {
// 	// 				// 	if !watchOptions { //if watchOptions is null or undefined
// 	// 				// 		currentOption = watchOptions
// 	// 				// 	}
// 	// 				// }
// 	// 			} else if parentOption == &typeAcquisitionDeclaration {
// 	// 				// if typeAcquisition != nil {
// 	// 				// 	typeAcquisition = getDefaultTypeAcquisition(configFileName)
// 	// 				// }
// 	// 				//currentOption := typeAcquisition
// 	// 			} //else Debug.fail("Unknown option");
// 	// 			//currentOption := value //*currentOption[option] = value find a way to do this
// 	// 		} else if keyText != "" && parentOption != nil { //&& parentOption.extraKeydiagnostics {
// 	// 			if parentOption.elementOptions != nil {
// 	// 				// errors.push(createUnknownOptionError(
// 	// 				// 	keyText,
// 	// 				// 	parentOption.extraKeydiagnostics,
// 	// 				// 	/*unknownOptionErrorText*/ undefined,
// 	// 				// 	propertyAssignment.name,
// 	// 				// 	sourceFile,
// 	// 				// ));
// 	// 			}
// 	// 			// else {
// 	// 			//     errors.push(createDiagnosticForNodeInSourceFile(sourceFile, propertyAssignment.name, parentOption.extraKeydiagnostics.unknownOptionDiagnostic, keyText));
// 	// 			// }
// 	// 		}
// 	// 	} else if parentOption == rootOptions {
// 	// 		// t := option //here need to fix
// 	// 		// if option.CommandLineOptionOfListType == extendsOptionDeclaration {
// 	// 		// 	extendedConfigPath = getExtendsConfigPathOrArray(value, host, basePath, configFileName, errors, propertyAssignment, propertyAssignment.initializer, sourceFile)
// 	// 		// } else if !option {
// 	// 		// 	if keyText == "excludes" {
// 	// 		// 		//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, propertyAssignment.name, diagnostics.Unknown_option_excludes_Did_you_mean_exclude));
// 	// 		// 	}
// 	// 		// 	// if (compiler.Find(commandOptionsWithoutBuild, opt => opt.name === keyText)) {
// 	// 		// 	//     rootCompilerOptions = append(rootCompilerOptions, propertyAssignment.name);
// 	// 		// 	// }
// 	// 		// }
// 	// 	}
// 	// }

// 	// json := convertConfigFileToObject(
// 	// 	sourceFile,
// 	// 	errors,
// 	// 	JsonConversionNotifier{rootOptions, onPropertySet},
// 	// )
// 	// if typeAcquisition == nil {
// 	// 	typeAcquisition = getDefaultTypeAcquisition(configFileName)
// 	// }

// 	// if rootCompilerOptions != nil && json != nil { //&& json.compilerOptions == nil {
// 	// 	//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, rootCompilerOptions[0], diagnostics._0_should_be_set_inside_the_compilerOptions_object_of_the_config_json_file, getTextOfPropertyName(rootCompilerOptions[0]) as string));
// 	// }

// 	return &ParsedTsconfig{
// 		// raw:     json,
// 		options: &options,
// 		//watchOptions:    watchOptions,
// 		// typeAcquisition: typeAcquisition,
// 		//extendedConfigPath: extendedConfigPath,
// 	}

// }

func getExtendedConfig(
	sourceFile *ast.SourceFile,
	extendedConfigPath string,
	host ParseConfigHost,
	resolutionStack []string,
	errors []*ast.Diagnostic,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
	result extendsResult,
) ParsedTsconfig {
	var path string
	if host.useCaseSensitiveFileNames {
		path = extendedConfigPath
	} else {
		path = tspath.ToFileNameLowerCase(extendedConfigPath)
	}
	var value ExtendedConfigCacheEntry
	var extendedResult *ast.SourceFile
	var extendedConfig ParsedTsconfig

	value = (*extendedConfigCache)[path]
	if extendedConfigCache != nil && value == (ExtendedConfigCacheEntry{}) {
		extendedResult = value.extendedResult
		extendedConfig = value.extendedConfig
	} else {
		extendedResult = readJsonConfigFile(extendedConfigPath, host.readFile) //probably readfile will give undefined
		if extendedResult != nil {                                             //parsediagnostics.length { //come back
			extendedConfig = parseConfig(nil, extendedResult, host, tspath.GetDirectoryPath(extendedConfigPath), tspath.GetBaseFileName(extendedConfigPath), resolutionStack, errors, extendedConfigCache)
		}
		if extendedConfigCache != nil {
			(*extendedConfigCache)[path] = ExtendedConfigCacheEntry{extendedResult, extendedConfig}
		}
	}
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
	*ast.SourceFile
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
	rootOptions   tsConfigOnlyOption //TsConfigOnlyOption
	onPropertySet func(keyText string, value any, propertyAssignment ast.PropertyAssignment, parentOption commandLineOptionOfTypes, option CommandLineOption)
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

// func isNullOrUndefined(x any) { // eslint-disable-line no-restricted-syntax
// 	return x == nil // eslint-disable-line no-restricted-syntax
// }

func isCompilerOptionsValue(option CommandLineOption, value any) core.CompilerOptionsValue {
	if option != (CommandLineOption{}) {
		//if compiler.Checker.IsNullOrUndefined(value) {
		if value == nil {
			option.disallowNullOrUndefined = false // All options are undefinable/nullable
			return core.CompilerOptionsValue{BooleanValue: option.disallowNullOrUndefined}
		}
		if option.kind == "list" {
			_, ok := value.([]string)
			return core.CompilerOptionsValue{BooleanValue: ok}
		}
		if option.kind == "listOrElement" {
			_, ok := value.([]string)
			return core.CompilerOptionsValue{BooleanValue: ok}
			//isCompilerOptionsValue(option.element, value);
		}
		var expectedType = string(option.kind)
		k := reflect.TypeOf(value)
		fmt.Println(k)
		if CommandLineOptionTypeEnum == option.kind {
			expectedType = "string"
		}
		return core.CompilerOptionsValue{BooleanValue: reflect.TypeOf(value).String() == expectedType}
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
	value core.CompilerOptionsValue,
	errors []ast.Diagnostic,
	valueExpression ast.Expression,
	sourceFile tsConfigSourceFile,
) core.CompilerOptionsValue {
	// if value == (compilerOptionsValue{}) {
	// 	return value
	// } //if (isNullOrUndefined(value)) return undefined;
	d, _ := (*opt.extraValidation)(value)
	if d == nil {
		return value
	}
	//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, ...d));
	return core.CompilerOptionsValue{}
}

func convertJsonOptionOfCustomType(
	opt CommandLineOption,
	value string,
	errors []*ast.Diagnostic,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) core.CompilerOptionsValue {
	if value == "" {
		return core.CompilerOptionsValue{}
	}
	key := strings.ToLower(value)
	val := commandLineOptionElements[key] //const val = opt.type.get(key);
	if val != nil {
		//return validateJsonOptionValue(opt, val, errors, valueExpression, sourceFile)
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
) []core.CompilerOptionsValue {
	// y := core.Filter(core.Map(values, func(v interface{}) compilerOptionsValue {
	// 	return convertJsonOption(*option.Elements(), v, basePath, errors, propertyAssignment, valueExpression.Elements.Nodes[0], sourceFile)
	// }), func(v compilerOptionsValue) bool {
	// 	if option.listPreserveFalsyValues {
	// 		return true
	// 	} else {
	// 		//return v != nil && v != false && v != 0 && v != ""
	// 		return false
	// 	}
	// })
	index := 0 //need to be changed
	var expression *ast.Node
	mappedValue := core.Map(values, func(v string) core.CompilerOptionsValue {
		if valueExpression != nil {
			expression = valueExpression.AsArrayLiteralExpression().Elements.Nodes[index]
		}
		var t core.CompilerOptionsValue = convertJsonOption(*option.Elements(), v, basePath, errors, propertyAssignment, expression, sourceFile)
		index++
		return t
	})
	filteredValues := core.Filter(mappedValue, func(v core.CompilerOptionsValue) bool {
		if option.listPreserveFalsyValues {
			return true
		} else {
			//return v != nil && v != false && v != 0 && v != ""
			return false
		}
	})
	return filteredValues
}

func convertJsonOption(
	opt CommandLineOption,
	value any,
	basePath string,
	errors []*ast.Diagnostic,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) core.CompilerOptionsValue {
	if opt.isCommandLineOnly != false {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, propertyAssignment?.name, diagnostics.Option_0_can_only_be_specified_on_command_line, opt.name));
		return core.CompilerOptionsValue{}
	}
	if isCompilerOptionsValue(opt, value).BooleanValue {
		optType := opt.kind
		_, ok := value.([]string)
		if (optType == "list") && ok {
			list := convertJsonOptionOfListType(opt, value.([]string), basePath, errors, propertyAssignment, valueExpression, sourceFile) //as ArrayLiteralExpression | undefined
			return core.CompilerOptionsValue{CompilerOptionsValueSlice: list}
		} else if optType == "listOrElement" {
			if ok {
				return core.CompilerOptionsValue{CompilerOptionsValueSlice: convertJsonOptionOfListType(opt, value.([]string), basePath, errors, propertyAssignment, valueExpression, sourceFile)}
			} else {
				return convertJsonOption(*opt.Elements(), value, basePath, errors, propertyAssignment, valueExpression, sourceFile)
			}
		} else if !(opt.kind == "string") {
			return convertJsonOptionOfCustomType(opt, value.(string), errors, valueExpression, sourceFile)
		}
		//const validatedValue = validateJsonOptionValue(opt, value, errors, valueExpression, sourceFile)
		//return isNullOrUndefined(validatedValue) ? validatedValue : normalizeNonListOptionValue(opt, basePath, validatedValue);
	} else {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.Compiler_option_0_requires_a_value_of_type_1, opt.name, getCompilerOptionValueTypeString(opt)));
	}
	return core.CompilerOptionsValue{}
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

// func getExtendsConfigPathOrArray(
// 	value compilerOptionsValue,
// 	host ParseConfigHost,
// 	basePath string,
// 	configFileName string,
// 	errors []ast.Diagnostic,
// 	propertyAssignment ast.PropertyAssignment,
// 	valueExpression ast.Expression,
// 	sourceFile *ast.SourceFile,
// ) {
// 	var extendedConfigPath []string
// 	var newBase string
// 	if configFileName != "" {
// 		newBase = directoryOfCombinedPath(configFileName, basePath)
// 	} else {
// 		newBase = basePath
// 	}
// 	if compiler.IsString(value) {
// 		extendedConfigPath = []string{getExtendsConfigPath(
// 			value.stringValue,
// 			host,
// 			newBase,
// 			errors,
// 			valueExpression,
// 			sourceFile,
// 		)}
// 	} else if compiler.IsSlice(value) {
// 		extendedConfigPath = []string{}
// 		// for index := 0; index < len(value); index++ {
// 		// 	fileName := value[index]
// 		// 	if compiler.IsString(fileName) {
// 		// 		extendedConfigPath = append(
// 		// 			extendedConfigPath,
// 		// 			getExtendsConfigPath(
// 		// 				fileName,
// 		// 				host,
// 		// 				newBase,
// 		// 				errors,
// 		// 				valueExpression, //(valueExpression as ArrayLiteralExpression | undefined)?.elements[index],
// 		// 				sourceFile,
// 		// 			),
// 		// 		)
// 		// 	} else {
// 		// 		convertJsonOption(extendsOptionDeclaration.element, value, basePath, errors, propertyAssignment, valueExpression.elements[index], sourceFile)
// 		// 	}
// 		// }
// 	} else {
// 		//convertJsonOption(extendsOptionDeclaration, value, basePath, errors, propertyAssignment, valueExpression, sourceFile)
// 	}
// 	//return extendedConfigPath
// }

// var commandLineTypeAcquisitionMapCache map[string]commandLineOption

// func getCommandLineTypeAcquisitionMap() {
// 	// if commandLineTypeAcquisitionMapCache != nil {
// 	// 	return commandLineTypeAcquisitionMap
// 	// }
// 	// commandLineTypeAcquisitionMapCache = commandLineOptionsToMap(typeAcquisitionDeclarations) //todo
// 	// return commandLineTypeAcquisitionMapCache
// }

var commonOptionsWithBuild = []CommandLineOption{
	{
		name:                     "help",
		shortName:                "h",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		isCommandLineOnly:        true,
		category:                 diagnostics.Command_line_Options,
		description:              diagnostics.Print_this_message,
		defaultValueDescription:  false,
	},
	{
		name:                    "help",
		shortName:               "?",
		kind:                    "boolean",
		isCommandLineOnly:       true,
		category:                diagnostics.Command_line_Options,
		defaultValueDescription: false,
	},
	{
		name:                     "watch",
		shortName:                "w",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		isCommandLineOnly:        true,
		category:                 diagnostics.Command_line_Options,
		description:              diagnostics.Watch_input_files,
		defaultValueDescription:  false,
	},
	{
		name:                     "preserveWatchOutput",
		kind:                     "boolean",
		showInSimplifiedHelpView: false,
		category:                 diagnostics.Output_Formatting,
		description:              diagnostics.Disable_wiping_the_console_in_watch_mode,
		defaultValueDescription:  false,
	},
	{
		name: "listFiles",
		kind: "boolean",
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Print_all_of_the_files_read_during_the_compilation,
		defaultValueDescription: false,
	},
	{
		name: "explainFiles",
		kind: "boolean",
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Print_files_read_during_the_compilation_including_why_it_was_included,
		defaultValueDescription: false,
	},
	{
		name: "listEmittedFiles",
		kind: "boolean",
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Print_the_names_of_emitted_files_after_a_compilation,
		defaultValueDescription: false,
	},
	{
		name:                     "pretty",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Output_Formatting,
		description:              diagnostics.Enable_color_and_formatting_in_TypeScript_s_output_to_make_compiler_errors_easier_to_read,
		defaultValueDescription:  true,
	},
	{
		name: "traceResolution",
		kind: "boolean",
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Log_paths_used_during_the_moduleResolution_process,
		defaultValueDescription: false,
	},
	{
		name: "diagnostics",
		kind: "boolean",
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Output_compiler_performance_information_after_building,
		defaultValueDescription: false,
	},
	{
		name: "extendeddiagnostics",
		kind: "boolean",
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Output_more_detailed_compiler_performance_information_after_building,
		defaultValueDescription: false,
	},
	{
		name:       "generateCpuProfile",
		kind:       "string",
		isFilePath: true,
		paramType:  *diagnostics.FILE_OR_DIRECTORY,
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Emit_a_v8_CPU_profile_of_the_compiler_run_for_debugging,
		defaultValueDescription: "profile.cpuprofile",
	},
	{
		name:       "generateTrace",
		kind:       "string",
		isFilePath: true,
		paramType:  *diagnostics.DIRECTORY,
		//category: diagnostics.Compiler_diagnostics,
		description: diagnostics.Generates_an_event_trace_and_a_list_of_types,
	},
	{
		name:                    "incremental",
		shortName:               "i",
		kind:                    "boolean",
		category:                diagnostics.Projects,
		description:             diagnostics.Save_tsbuildinfo_files_to_allow_for_incremental_compilation_of_projects,
		transpileOptionValue:    0,
		defaultValueDescription: diagnostics.X_false_unless_composite_is_set, //check
	},
	{
		name:      "declaration",
		shortName: "d",
		kind:      "boolean",
		// Not setting affectsEmit because we calculate this flag might not affect full emit
		affectsBuildInfo:         true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		transpileOptionValue:     0,
		description:              diagnostics.Generate_d_ts_files_from_TypeScript_and_JavaScript_files_in_your_project,
		defaultValueDescription:  diagnostics.X_false_unless_composite_is_set,
	},
	{
		name: "declarationMap",
		kind: "boolean",
		// Not setting affectsEmit because we calculate this flag might not affect full emit
		affectsBuildInfo:         true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		defaultValueDescription:  false,
		description:              diagnostics.Create_sourcemaps_for_d_ts_files,
	},
	{
		name: "emitDeclarationOnly",
		kind: "boolean",
		// Not setting affectsEmit because we calculate this flag might not affect full emit
		affectsBuildInfo:         true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		description:              diagnostics.Only_output_d_ts_files_and_not_JavaScript_files,
		transpileOptionValue:     0,
		defaultValueDescription:  false,
	},
	{
		name: "sourceMap",
		kind: "boolean",
		// Not setting affectsEmit because we calculate this flag might not affect full emit
		affectsBuildInfo:         true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		defaultValueDescription:  false,
		description:              diagnostics.Create_source_map_files_for_emitted_JavaScript_files,
	},
	{
		name: "inlineSourceMap",
		kind: "boolean",
		// Not setting affectsEmit because we calculate this flag might not affect full emit
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		description:             diagnostics.Include_sourcemap_files_inside_the_emitted_JavaScript,
		defaultValueDescription: false,
	},
	{
		name:                     "noCheck",
		kind:                     "boolean",
		showInSimplifiedHelpView: false,
		//category: diagnostics.Compiler_diagnostics,
		description:             diagnostics.Disable_full_type_checking_only_critical_parse_and_emit_errors_will_be_reported,
		transpileOptionValue:    0,
		defaultValueDescription: false,
		// Not setting affectsSemanticdiagnostics or affectsBuildInfo because we dont want all diagnostics to go away, its handled in builder
	},
	{
		name:                     "noEmit",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		description:              diagnostics.Disable_emitting_files_from_a_compilation,
		transpileOptionValue:     2,
		defaultValueDescription:  false,
	},
	{
		name:                       "assumeChangesOnlyAffectDirectDependencies",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsEmit:                true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Watch_and_Build_Modes,
		description:                diagnostics.Have_recompiles_in_projects_that_use_incremental_and_watch_mode_assume_that_changes_within_a_file_will_only_affect_files_directly_depending_on_it,
		defaultValueDescription:    false,
	},
	{
		name:                    "locale",
		kind:                    "string",
		category:                diagnostics.Command_line_Options,
		isCommandLineOnly:       true,
		description:             diagnostics.Set_the_language_of_the_messaging_from_TypeScript_This_does_not_affect_emit,
		defaultValueDescription: diagnostics.Platform_specific,
	},
	{
		name: "lib",
		kind: CommandLineOptionTypeList,
		// elements: &CommandLineOption{
		// 	name:                    "lib",
		// 	kind:                   CommandLineOptionTypeEnum, //libMap,
		// 	defaultValueDescription: core.TSUnknown,
		// },
		affectsProgramStructure:  true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Language_and_Environment,
		description:              diagnostics.Specify_a_set_of_bundled_library_declaration_files_that_describe_the_target_runtime_environment,
		transpileOptionValue:     core.TSUnknown,
	},
}
var targetOptionDeclaration = CommandLineOption{
	name:      "target",
	shortName: "t",
	// kind: new Map(Object.entries({
	//     es3: ScriptTarget.ES3,
	//     es5: ScriptTarget.ES5,
	//     es6: ScriptTarget.ES2015,
	//     es2015: ScriptTarget.ES2015,
	//     es2016: ScriptTarget.ES2016,
	//     es2017: ScriptTarget.ES2017,
	//     es2018: ScriptTarget.ES2018,
	//     es2019: ScriptTarget.ES2019,
	//     es2020: ScriptTarget.ES2020,
	//     es2021: ScriptTarget.ES2021,
	//     es2022: ScriptTarget.ES2022,
	//     es2023: ScriptTarget.ES2023,
	//     es2024: ScriptTarget.ES2024,
	//     esnext: ScriptTarget.ESNext,
	// })),
	affectsSourceFile:       true,
	affectsModuleResolution: true,
	affectsEmit:             true,
	affectsBuildInfo:        true,
	//deprecatedKeys: new Set(["es3"]),
	paramType:                *diagnostics.VERSION,
	showInSimplifiedHelpView: true,
	category:                 diagnostics.Language_and_Environment,
	description:              diagnostics.Set_the_JavaScript_language_version_for_emitted_JavaScript_and_include_compatible_library_declarations,
	//defaultValueDescription: core.ScriptTarget.ScriptTargetES5,
}
var moduleOptionDeclaration = CommandLineOption{
	name:      "module",
	shortName: "m",
	// kind: new Map(Object.entries({
	//     none: ModuleKind.None,
	//     commonjs: ModuleKind.CommonJS,
	//     amd: ModuleKind.AMD,
	//     system: ModuleKind.System,
	//     umd: ModuleKind.UMD,
	//     es6: ModuleKind.ES2015,
	//     es2015: ModuleKind.ES2015,
	//     es2020: ModuleKind.ES2020,
	//     es2022: ModuleKind.ES2022,
	//     esnext: ModuleKind.ESNext,
	//     node16: ModuleKind.Node16,
	//     nodenext: ModuleKind.NodeNext,
	//     preserve: ModuleKind.Preserve,
	// })),
	affectsSourceFile:        true,
	affectsModuleResolution:  true,
	affectsEmit:              true,
	affectsBuildInfo:         true,
	paramType:                *diagnostics.KIND,
	showInSimplifiedHelpView: true,
	category:                 diagnostics.Modules,
	description:              diagnostics.Specify_what_module_code_is_generated,
	defaultValueDescription:  nil,
}
var commandOptionsWithoutBuild = []CommandLineOption{
	{
		name:                     "all",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Command_line_Options,
		description:              diagnostics.Show_all_compiler_options,
		defaultValueDescription:  false,
	},
	{
		name:                     "version",
		shortName:                "v",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Command_line_Options,
		description:              diagnostics.Print_the_compiler_s_version,
		defaultValueDescription:  false,
	},
	{
		name:                     "init",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Command_line_Options,
		description:              diagnostics.Initializes_a_TypeScript_project_and_creates_a_tsconfig_json_file,
		defaultValueDescription:  false,
	},
	{
		name:                     "project",
		shortName:                "p",
		kind:                     "string",
		isFilePath:               true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Command_line_Options,
		paramType:                *diagnostics.FILE_OR_DIRECTORY,
		description:              diagnostics.Compile_the_project_given_the_path_to_its_configuration_file_or_to_a_folder_with_a_tsconfig_json,
	},
	{
		name:                     "showConfig",
		kind:                     "boolean",
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Command_line_Options,
		isCommandLineOnly:        true,
		description:              diagnostics.Print_the_final_configuration_instead_of_building,
		defaultValueDescription:  false,
	},
	{
		name:                    "listFilesOnly",
		kind:                    "boolean",
		category:                diagnostics.Command_line_Options,
		isCommandLineOnly:       true,
		description:             diagnostics.Print_names_of_files_that_are_part_of_the_compilation_and_then_stop_processing,
		defaultValueDescription: false,
	},

	// Basic
	targetOptionDeclaration,
	moduleOptionDeclaration,
	{
		name: "lib",
		kind: "list",
		// element: {
		//     name: "lib",
		//     kind: libMap,
		//     defaultValueDescription: undefined,
		// },
		affectsProgramStructure:  true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Language_and_Environment,
		description:              diagnostics.Specify_a_set_of_bundled_library_declaration_files_that_describe_the_target_runtime_environment,
		transpileOptionValue:     0,
	},
	{
		name:                     "allowJs",
		kind:                     "boolean",
		allowJsFlag:              true,
		affectsBuildInfo:         true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.JavaScript_Support,
		description:              diagnostics.Allow_JavaScript_files_to_be_a_part_of_your_program_Use_the_checkJS_option_to_get_errors_from_these_files,
		defaultValueDescription:  false,
	},
	{
		name:                       "checkJs",
		kind:                       "boolean",
		affectsModuleResolution:    true,
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		showInSimplifiedHelpView:   true,
		category:                   diagnostics.JavaScript_Support,
		description:                diagnostics.Enable_error_reporting_in_type_checked_JavaScript_files,
		defaultValueDescription:    false,
	},
	{
		name: "jsx",
		//kind: jsxOptionMap,
		affectsSourceFile:       true,
		affectsEmit:             true,
		affectsBuildInfo:        true,
		affectsModuleResolution: true,
		// The checker emits an error when it sees JSX but this option is not set in compilerOptions.
		// This is effectively a semantic error, so mark this option as affecting semantic diagnostics
		// so we know to refresh errors when this option is changed.
		affectsSemanticdiagnostics: true,
		paramType:                  *diagnostics.KIND,
		showInSimplifiedHelpView:   true,
		category:                   diagnostics.Language_and_Environment,
		description:                diagnostics.Specify_what_JSX_code_is_generated,
		defaultValueDescription:    nil,
	},
	{
		name:                     "outFile",
		kind:                     "string",
		affectsEmit:              true,
		affectsBuildInfo:         true,
		affectsDeclarationPath:   true,
		isFilePath:               true,
		paramType:                *diagnostics.FILE,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		description:              diagnostics.Specify_a_file_that_bundles_all_outputs_into_one_JavaScript_file_If_declaration_is_true_also_designates_a_file_that_bundles_all_d_ts_output,
		transpileOptionValue:     0,
	},
	{
		name:                     "outDir",
		kind:                     "string",
		affectsEmit:              true,
		affectsBuildInfo:         true,
		affectsDeclarationPath:   true,
		isFilePath:               true,
		paramType:                *diagnostics.DIRECTORY,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		description:              diagnostics.Specify_an_output_folder_for_all_emitted_files,
	},
	{
		name:                    "rootDir",
		kind:                    "string",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		affectsDeclarationPath:  true,
		isFilePath:              true,
		paramType:               *diagnostics.LOCATION,
		category:                diagnostics.Modules,
		description:             diagnostics.Specify_the_root_folder_within_your_source_files,
		defaultValueDescription: diagnostics.Computed_from_the_list_of_input_files,
	},
	{
		name: "composite",
		kind: "boolean",
		// Not setting affectsEmit because we calculate this flag might not affect full emit
		affectsBuildInfo:        true,
		isTSConfigOnly:          true,
		category:                diagnostics.Projects,
		transpileOptionValue:    0,
		defaultValueDescription: false,
		description:             diagnostics.Enable_constraints_that_allow_a_TypeScript_project_to_be_used_with_project_references,
	},
	{
		name:                    "tsBuildInfoFile",
		kind:                    "string",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		isFilePath:              true,
		paramType:               *diagnostics.FILE,
		category:                diagnostics.Projects,
		transpileOptionValue:    0,
		defaultValueDescription: ".tsbuildinfo",
		description:             diagnostics.Specify_the_path_to_tsbuildinfo_incremental_compilation_file,
	},
	{
		name:                     "removeComments",
		kind:                     "boolean",
		affectsEmit:              true,
		affectsBuildInfo:         true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Emit,
		defaultValueDescription:  false,
		description:              diagnostics.Disable_emitting_comments,
	},
	{
		name:                    "importHelpers",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		affectsSourceFile:       true,
		category:                diagnostics.Emit,
		description:             diagnostics.Allow_importing_helper_functions_from_tslib_once_per_project_instead_of_including_them_per_file,
		defaultValueDescription: false,
	},
	{
		name: "importsNotUsedAsValues",
		// kind: new Map(Object.entries({
		//     remove: ImportsNotUsedAsValues.Remove,
		//     preserve: ImportsNotUsedAsValues.Preserve,
		//     error: ImportsNotUsedAsValues.Error,
		// })),
		affectsEmit:                true,
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Backwards_Compatibility,
		description:                diagnostics.Specify_emit_Slashchecking_behavior_for_imports_that_are_only_used_for_types,
		//defaultValueDescription: ImportsNotUsedAsValues.Remove,
	},
	{
		name:                    "downlevelIteration",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		description:             diagnostics.Emit_more_compliant_but_verbose_and_less_performant_JavaScript_for_iteration,
		defaultValueDescription: false,
	},
	{
		name:                    "isolatedModules",
		kind:                    "boolean",
		category:                diagnostics.Interop_Constraints,
		description:             diagnostics.Ensure_that_each_file_can_be_safely_transpiled_without_relying_on_other_imports,
		transpileOptionValue:    2,
		defaultValueDescription: false,
	},
	{
		name:                       "verbatimModuleSyntax",
		kind:                       "boolean",
		affectsEmit:                true,
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Interop_Constraints,
		description:                diagnostics.Do_not_transform_or_elide_any_imports_or_exports_not_marked_as_type_only_ensuring_they_are_written_in_the_output_file_s_format_based_on_the_module_setting,
		defaultValueDescription:    false,
	},
	{
		name:                       "isolatedDeclarations",
		kind:                       "boolean",
		category:                   diagnostics.Interop_Constraints,
		description:                diagnostics.Require_sufficient_annotation_on_exports_so_other_tools_can_trivially_generate_declaration_files,
		defaultValueDescription:    false,
		affectsBuildInfo:           true,
		affectsSemanticdiagnostics: true,
	},

	// Strict Type Checks
	{
		name: "strict",
		kind: "boolean",
		// Though this affects semantic diagnostics, affectsSemanticdiagnostics is not set here
		// The value of each strictFlag depends on own strictFlag value or this and never accessed directly.
		// But we need to store `strict` in builf info, even though it won't be examined directly, so that the
		// flags it controls (e.g. `strictNullChecks`) will be retrieved correctly
		affectsBuildInfo:         true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Type_Checking,
		description:              diagnostics.Enable_all_strict_type_checking_options,
		defaultValueDescription:  false,
	},
	{
		name:                       "noImplicitAny",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Enable_error_reporting_for_expressions_and_declarations_with_an_implied_any_type,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                       "strictNullChecks",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.When_type_checking_take_into_account_null_and_undefined,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                       "strictFunctionTypes",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.When_assigning_functions_check_to_ensure_parameters_and_the_return_values_are_subtype_compatible,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                       "strictBindCallApply",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Check_that_the_arguments_for_bind_call_and_apply_methods_match_the_original_function,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                       "strictPropertyInitialization",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Check_for_class_properties_that_are_declared_but_not_set_in_the_constructor,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                       "strictBuiltinIteratorReturn",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Built_in_iterators_are_instantiated_with_a_TReturn_type_of_undefined_instead_of_any,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                       "noImplicitThis",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Enable_error_reporting_when_this_is_given_the_type_any,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                       "useUnknownInCatchVariables",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		strictFlag:                 true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Default_catch_clause_variables_as_unknown_instead_of_any,
		defaultValueDescription:    diagnostics.X_false_unless_strict_is_set,
	},
	{
		name:                    "alwaysStrict",
		kind:                    "boolean",
		affectsSourceFile:       true,
		affectsEmit:             true,
		affectsBuildInfo:        true,
		strictFlag:              true,
		category:                diagnostics.Type_Checking,
		description:             diagnostics.Ensure_use_strict_is_always_emitted,
		defaultValueDescription: diagnostics.X_false_unless_strict_is_set,
	},

	// Additional Checks
	{
		name:                       "noUnusedLocals",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Enable_error_reporting_when_local_variables_aren_t_read,
		defaultValueDescription:    false,
	},
	{
		name:                       "noUnusedParameters",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Raise_an_error_when_a_function_parameter_isn_t_read,
		defaultValueDescription:    false,
	},
	{
		name:                       "exactOptionalPropertyTypes",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Interpret_optional_property_types_as_written_rather_than_adding_undefined,
		defaultValueDescription:    false,
	},
	{
		name:                       "noImplicitReturns",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Enable_error_reporting_for_codepaths_that_do_not_explicitly_return_in_a_function,
		defaultValueDescription:    false,
	},
	{
		name:                       "noFallthroughCasesInSwitch",
		kind:                       "boolean",
		affectsBinddiagnostics:     true,
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Enable_error_reporting_for_fallthrough_cases_in_switch_statements,
		defaultValueDescription:    false,
	},
	{
		name:                       "noUncheckedIndexedAccess",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Add_undefined_to_a_type_when_accessed_using_an_index,
		defaultValueDescription:    false,
	},
	{
		name:                       "noImplicitOverride",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Ensure_overriding_members_in_derived_classes_are_marked_with_an_override_modifier,
		defaultValueDescription:    false,
	},
	{
		name:                       "noPropertyAccessFromIndexSignature",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		showInSimplifiedHelpView:   false,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Enforces_using_indexed_accessors_for_keys_declared_using_an_indexed_type,
		defaultValueDescription:    false,
	},

	// Module Resolution
	{
		name: "moduleResolution",
		// kind: new Map(Object.entries({
		//     // N.B. The first entry specifies the value shown in `tsc --init`
		//     node10: ModuleResolutionKind.Node10,
		//     node: ModuleResolutionKind.Node10,
		//     classic: ModuleResolutionKind.Classic,
		//     node16: ModuleResolutionKind.Node16,
		//     nodenext: ModuleResolutionKind.NodeNext,
		//     bundler: ModuleResolutionKind.Bundler,
		// })),
		// deprecatedKeys: new Set(["node"]),
		affectsSourceFile:       true,
		affectsModuleResolution: true,
		paramType:               *diagnostics.STRATEGY,
		category:                diagnostics.Modules,
		description:             diagnostics.Specify_how_TypeScript_looks_up_a_file_from_a_given_module_specifier,
		defaultValueDescription: diagnostics.X_module_AMD_or_UMD_or_System_or_ES6_then_Classic_Otherwise_Node,
	},
	{
		name:                    "baseUrl",
		kind:                    "string",
		affectsModuleResolution: true,
		isFilePath:              true,
		category:                diagnostics.Modules,
		description:             diagnostics.Specify_the_base_directory_to_resolve_non_relative_module_names,
	},
	{
		// this option can only be specified in tsconfig.json
		// use type = object to copy the value as-is
		name:                               "paths",
		kind:                               "object",
		affectsModuleResolution:            true,
		allowConfigDirTemplateSubstitution: true,
		isTSConfigOnly:                     true,
		category:                           diagnostics.Modules,
		description:                        diagnostics.Specify_a_set_of_entries_that_re_map_imports_to_additional_lookup_locations,
		transpileOptionValue:               0,
	},
	{
		// this option can only be specified in tsconfig.json
		// use type = object to copy the value as-is
		name:           "rootDirs",
		kind:           "list",
		isTSConfigOnly: true,
		// element: {
		//     name: "rootDirs",
		//     kind: "string",
		//     isFilePath: true,
		// },
		affectsModuleResolution:            true,
		allowConfigDirTemplateSubstitution: true,
		category:                           diagnostics.Modules,
		description:                        diagnostics.Allow_multiple_folders_to_be_treated_as_one_when_resolving_modules,
		transpileOptionValue:               0,
		defaultValueDescription:            diagnostics.Computed_from_the_list_of_input_files,
	},
	{
		name: "typeRoots",
		kind: "list",
		// element: {
		//     name: "typeRoots",
		//     kind: "string",
		//     isFilePath: true,
		// },
		affectsModuleResolution:            true,
		allowConfigDirTemplateSubstitution: true,
		category:                           diagnostics.Modules,
		description:                        diagnostics.Specify_multiple_folders_that_act_like_Slashnode_modules_Slash_types,
	},
	{
		name: "types",
		kind: "list",
		// element: {
		//     name: "types",
		//     kind: "string",
		// },
		affectsProgramStructure:  true,
		showInSimplifiedHelpView: true,
		category:                 diagnostics.Modules,
		description:              diagnostics.Specify_type_package_names_to_be_included_without_being_referenced_in_a_source_file,
		transpileOptionValue:     0,
	},
	{
		name:                       "allowSyntheticDefaultImports",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Interop_Constraints,
		description:                diagnostics.Allow_import_x_from_y_when_a_module_doesn_t_have_a_default_export,
		defaultValueDescription:    diagnostics.X_module_system_or_esModuleInterop,
	},
	{
		name:                       "esModuleInterop",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsEmit:                true,
		affectsBuildInfo:           true,
		showInSimplifiedHelpView:   true,
		category:                   diagnostics.Interop_Constraints,
		description:                diagnostics.Emit_additional_JavaScript_to_ease_support_for_importing_CommonJS_modules_This_enables_allowSyntheticDefaultImports_for_type_compatibility,
		defaultValueDescription:    false,
	},
	{
		name:                    "preserveSymlinks",
		kind:                    "boolean",
		category:                diagnostics.Interop_Constraints,
		description:             diagnostics.Disable_resolving_symlinks_to_their_realpath_This_correlates_to_the_same_flag_in_node,
		defaultValueDescription: false,
	},
	{
		name:                       "allowUmdGlobalAccess",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Modules,
		description:                diagnostics.Allow_accessing_UMD_globals_from_modules,
		defaultValueDescription:    false,
	},
	{
		name: "moduleSuffixes",
		kind: "list",
		// element: {
		//     name: "suffix",
		//     kind: "string",
		// },
		listPreserveFalsyValues: true,
		affectsModuleResolution: true,
		category:                diagnostics.Modules,
		description:             diagnostics.List_of_file_name_suffixes_to_search_when_resolving_a_module,
	},
	{
		name:                       "allowImportingTsExtensions",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Modules,
		description:                diagnostics.Allow_imports_to_include_TypeScript_file_extensions_Requires_moduleResolution_bundler_and_either_noEmit_or_emitDeclarationOnly_to_be_set,
		defaultValueDescription:    false,
		transpileOptionValue:       0,
	},
	{
		name:                       "rewriteRelativeImportExtensions",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Modules,
		//description: diagnostics.Rewrite_ts_tsx_mts_and_cts_file_extensions_in_relative_import_paths_to_their_JavaScript_equivalent_in_output_files,
		defaultValueDescription: false,
	},
	{
		name:                    "resolvePackageJsonExports",
		kind:                    "boolean",
		affectsModuleResolution: true,
		category:                diagnostics.Modules,
		description:             diagnostics.Use_the_package_json_exports_field_when_resolving_package_imports,
		defaultValueDescription: diagnostics.X_true_when_moduleResolution_is_node16_nodenext_or_bundler_otherwise_false,
	},
	{
		name:                    "resolvePackageJsonImports",
		kind:                    "boolean",
		affectsModuleResolution: true,
		category:                diagnostics.Modules,
		description:             diagnostics.Use_the_package_json_imports_field_when_resolving_imports,
		defaultValueDescription: diagnostics.X_true_when_moduleResolution_is_node16_nodenext_or_bundler_otherwise_false,
	},
	{
		name: "customConditions",
		kind: "list",
		// element: {
		//     name: "condition",
		//     kind: "string",
		// },
		affectsModuleResolution: true,
		category:                diagnostics.Modules,
		description:             diagnostics.Conditions_to_set_in_addition_to_the_resolver_specific_defaults_when_resolving_imports,
	},
	{
		name:                       "noUncheckedSideEffectImports",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Modules,
		description:                diagnostics.Check_side_effect_imports,
		defaultValueDescription:    false,
	},

	// Source Maps
	{
		name:             "sourceRoot",
		kind:             "string",
		affectsEmit:      true,
		affectsBuildInfo: true,
		paramType:        *diagnostics.LOCATION,
		category:         diagnostics.Emit,
		description:      diagnostics.Specify_the_root_path_for_debuggers_to_find_the_reference_source_code,
	},
	{
		name:             "mapRoot",
		kind:             "string",
		affectsEmit:      true,
		affectsBuildInfo: true,
		paramType:        *diagnostics.LOCATION,
		category:         diagnostics.Emit,
		description:      diagnostics.Specify_the_location_where_debugger_should_locate_map_files_instead_of_generated_locations,
	},
	{
		name:                    "inlineSources",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		description:             diagnostics.Include_source_code_in_the_sourcemaps_inside_the_emitted_JavaScript,
		defaultValueDescription: false,
	},

	// Experimental
	{
		name:                       "experimentalDecorators",
		kind:                       "boolean",
		affectsEmit:                true,
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Language_and_Environment,
		description:                diagnostics.Enable_experimental_support_for_legacy_experimental_decorators,
		defaultValueDescription:    false,
	},
	{
		name:                       "emitDecoratorMetadata",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsEmit:                true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Language_and_Environment,
		description:                diagnostics.Emit_design_type_metadata_for_decorated_declarations_in_source_files,
		defaultValueDescription:    false,
	},

	// Advanced
	{
		name:                    "jsxFactory",
		kind:                    "string",
		category:                diagnostics.Language_and_Environment,
		description:             diagnostics.Specify_the_JSX_factory_function_used_when_targeting_React_JSX_emit_e_g_React_createElement_or_h,
		defaultValueDescription: "`React.createElement`",
	},
	{
		name:                    "jsxFragmentFactory",
		kind:                    "string",
		category:                diagnostics.Language_and_Environment,
		description:             diagnostics.Specify_the_JSX_Fragment_reference_used_for_fragments_when_targeting_React_JSX_emit_e_g_React_Fragment_or_Fragment,
		defaultValueDescription: "React.Fragment",
	},
	{
		name:                       "jsxImportSource",
		kind:                       "string",
		affectsSemanticdiagnostics: true,
		affectsEmit:                true,
		affectsBuildInfo:           true,
		affectsModuleResolution:    true,
		affectsSourceFile:          true,
		category:                   diagnostics.Language_and_Environment,
		description:                diagnostics.Specify_module_specifier_used_to_import_the_JSX_factory_functions_when_using_jsx_Colon_react_jsx_Asterisk,
		defaultValueDescription:    "react",
	},
	{
		name:                    "resolveJsonModule",
		kind:                    "boolean",
		affectsModuleResolution: true,
		category:                diagnostics.Modules,
		description:             diagnostics.Enable_importing_json_files,
		defaultValueDescription: false,
	},
	{
		name:                    "allowArbitraryExtensions",
		kind:                    "boolean",
		affectsProgramStructure: true,
		category:                diagnostics.Modules,
		description:             diagnostics.Enable_importing_files_with_any_extension_provided_a_declaration_file_is_present,
		defaultValueDescription: false,
	},

	{
		name:                   "out",
		kind:                   "string",
		affectsEmit:            true,
		affectsBuildInfo:       true,
		affectsDeclarationPath: true,
		isFilePath:             false, // This is intentionally broken to support compatibility with existing tsconfig files
		// for correct behaviour, please use outFile
		category:             diagnostics.Backwards_Compatibility,
		paramType:            *diagnostics.FILE,
		transpileOptionValue: 0,
		description:          diagnostics.Deprecated_setting_Use_outFile_instead,
	},
	{
		name:                    "reactNamespace",
		kind:                    "string",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Language_and_Environment,
		description:             diagnostics.Specify_the_object_invoked_for_createElement_This_only_applies_when_targeting_react_JSX_emit,
		defaultValueDescription: "`React`",
	},
	{
		name: "skipDefaultLibCheck",
		kind: "boolean",
		// We need to store these to determine whether `lib` files need to be rechecked
		affectsBuildInfo:        true,
		category:                diagnostics.Completeness,
		description:             diagnostics.Skip_type_checking_d_ts_files_that_are_included_with_TypeScript,
		defaultValueDescription: false,
	},
	{
		name:                    "charset",
		kind:                    "string",
		category:                diagnostics.Backwards_Compatibility,
		description:             diagnostics.No_longer_supported_In_early_versions_manually_set_the_text_encoding_for_reading_files,
		defaultValueDescription: "utf8",
	},
	{
		name:                    "emitBOM",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		description:             diagnostics.Emit_a_UTF_8_Byte_Order_Mark_BOM_in_the_beginning_of_output_files,
		defaultValueDescription: false,
	},
	{
		name: "newLine",
		// kind: new Map(Object.entries({
		//     crlf: NewLineKind.CarriageReturnLineFeed,
		//     lf: NewLineKind.LineFeed,
		// })),
		affectsEmit:             true,
		affectsBuildInfo:        true,
		paramType:               *diagnostics.NEWLINE,
		category:                diagnostics.Emit,
		description:             diagnostics.Set_the_newline_character_for_emitting_files,
		defaultValueDescription: "lf",
	},
	{
		name:                       "noErrorTruncation",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Output_Formatting,
		description:                diagnostics.Disable_truncating_types_in_error_messages,
		defaultValueDescription:    false,
	},
	{
		name:                    "noLib",
		kind:                    "boolean",
		category:                diagnostics.Language_and_Environment,
		affectsProgramStructure: true,
		description:             diagnostics.Disable_including_any_library_files_including_the_default_lib_d_ts,
		// We are not returning a sourceFile for lib file when asked by the program,
		// so pass --noLib to avoid reporting a file not found error.
		transpileOptionValue:    2,
		defaultValueDescription: false,
	},
	{
		name:                    "noResolve",
		kind:                    "boolean",
		affectsModuleResolution: true,
		category:                diagnostics.Modules,
		description:             diagnostics.Disallow_import_s_require_s_or_reference_s_from_expanding_the_number_of_files_TypeScript_should_add_to_a_project,
		// We are not doing a full typecheck, we are not resolving the whole context,
		// so pass --noResolve to avoid reporting missing file errors.
		transpileOptionValue:    2,
		defaultValueDescription: false,
	},
	{
		name:                    "stripInternal",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		description:             diagnostics.Disable_emitting_declarations_that_have_internal_in_their_JSDoc_comments,
		defaultValueDescription: false,
	},
	{
		name:                    "disableSizeLimit",
		kind:                    "boolean",
		affectsProgramStructure: true,
		category:                diagnostics.Editor_Support,
		description:             diagnostics.Remove_the_20mb_cap_on_total_source_code_size_for_JavaScript_files_in_the_TypeScript_language_server,
		defaultValueDescription: false,
	},
	{
		name:                    "disableSourceOfProjectReferenceRedirect",
		kind:                    "boolean",
		isTSConfigOnly:          true,
		category:                diagnostics.Projects,
		description:             diagnostics.Disable_preferring_source_files_instead_of_declaration_files_when_referencing_composite_projects,
		defaultValueDescription: false,
	},
	{
		name:                    "disableSolutionSearching",
		kind:                    "boolean",
		isTSConfigOnly:          true,
		category:                diagnostics.Projects,
		description:             diagnostics.Opt_a_project_out_of_multi_project_reference_checking_when_editing,
		defaultValueDescription: false,
	},
	{
		name:                    "disableReferencedProjectLoad",
		kind:                    "boolean",
		isTSConfigOnly:          true,
		category:                diagnostics.Projects,
		description:             diagnostics.Reduce_the_number_of_projects_loaded_automatically_by_TypeScript,
		defaultValueDescription: false,
	},
	{
		name:                       "noImplicitUseStrict",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Backwards_Compatibility,
		description:                diagnostics.Disable_adding_use_strict_directives_in_emitted_JavaScript_files,
		defaultValueDescription:    false,
	},
	{
		name:                    "noEmitHelpers",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		description:             diagnostics.Disable_generating_custom_helper_functions_like_extends_in_compiled_output,
		defaultValueDescription: false,
	},
	{
		name:                    "noEmitOnError",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		transpileOptionValue:    0,
		description:             diagnostics.Disable_emitting_files_if_any_type_checking_errors_are_reported,
		defaultValueDescription: false,
	},
	{
		name:                    "preserveConstEnums",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Emit,
		description:             diagnostics.Disable_erasing_const_enum_declarations_in_generated_code,
		defaultValueDescription: false,
	},
	{
		name:                   "declarationDir",
		kind:                   "string",
		affectsEmit:            true,
		affectsBuildInfo:       true,
		affectsDeclarationPath: true,
		isFilePath:             true,
		paramType:              *diagnostics.DIRECTORY,
		category:               diagnostics.Emit,
		transpileOptionValue:   0,
		description:            diagnostics.Specify_the_output_directory_for_generated_declaration_files,
	},
	{
		name: "skipLibCheck",
		kind: "boolean",
		// We need to store these to determine whether `lib` files need to be rechecked
		affectsBuildInfo:        true,
		category:                diagnostics.Completeness,
		description:             diagnostics.Skip_type_checking_all_d_ts_files,
		defaultValueDescription: false,
	},
	{
		name:                       "allowUnusedLabels",
		kind:                       "boolean",
		affectsBinddiagnostics:     true,
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Disable_error_reporting_for_unused_labels,
		defaultValueDescription:    0,
	},
	{
		name:                       "allowUnreachableCode",
		kind:                       "boolean",
		affectsBinddiagnostics:     true,
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Type_Checking,
		description:                diagnostics.Disable_error_reporting_for_unreachable_code,
		defaultValueDescription:    0,
	},
	{
		name:                       "suppressExcessPropertyErrors",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Backwards_Compatibility,
		description:                diagnostics.Disable_reporting_of_excess_property_errors_during_the_creation_of_object_literals,
		defaultValueDescription:    false,
	},
	{
		name:                       "suppressImplicitAnyIndexErrors",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Backwards_Compatibility,
		description:                diagnostics.Suppress_noImplicitAny_errors_when_indexing_objects_that_lack_index_signatures,
		defaultValueDescription:    false,
	},
	{
		name:                    "forceConsistentCasingInFileNames",
		kind:                    "boolean",
		affectsModuleResolution: true,
		category:                diagnostics.Interop_Constraints,
		description:             diagnostics.Ensure_that_casing_is_correct_in_imports,
		defaultValueDescription: true,
	},
	{
		name:                    "maxNodeModuleJsDepth",
		kind:                    "number",
		affectsModuleResolution: true,
		category:                diagnostics.JavaScript_Support,
		description:             diagnostics.Specify_the_maximum_folder_depth_used_for_checking_JavaScript_files_from_node_modules_Only_applicable_with_allowJs,
		defaultValueDescription: 0,
	},
	{
		name:                       "noStrictGenericChecks",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Backwards_Compatibility,
		description:                diagnostics.Disable_strict_checking_of_generic_signatures_in_function_types,
		defaultValueDescription:    false,
	},
	{
		name:                       "useDefineForClassFields",
		kind:                       "boolean",
		affectsSemanticdiagnostics: true,
		affectsEmit:                true,
		affectsBuildInfo:           true,
		category:                   diagnostics.Language_and_Environment,
		description:                diagnostics.Emit_ECMAScript_standard_compliant_class_fields,
		defaultValueDescription:    diagnostics.X_true_for_ES2022_and_above_including_ESNext,
	},
	{
		name:                    "preserveValueImports",
		kind:                    "boolean",
		affectsEmit:             true,
		affectsBuildInfo:        true,
		category:                diagnostics.Backwards_Compatibility,
		description:             diagnostics.Preserve_unused_imported_values_in_the_JavaScript_output_that_would_otherwise_be_removed,
		defaultValueDescription: false,
	},

	{
		name:                    "keyofStringsOnly",
		kind:                    "boolean",
		category:                diagnostics.Backwards_Compatibility,
		description:             diagnostics.Make_keyof_only_return_strings_instead_of_string_numbers_or_symbols_Legacy_option,
		defaultValueDescription: false,
	},
	{
		// A list of plugins to load in the language service
		name:           "plugins",
		kind:           "list",
		isTSConfigOnly: true,
		// element: {
		//     name: "plugin",
		//     kind: "object",
		// },
		description: diagnostics.Specify_a_list_of_language_service_plugins_to_include,
		category:    diagnostics.Editor_Support,
	},
	{
		name: "moduleDetection",
		// kind: new Map(Object.entries({
		//     auto: ModuleDetectionKind.Auto,
		//     legacy: ModuleDetectionKind.Legacy,
		//     force: ModuleDetectionKind.Force,
		// })),
		affectsSourceFile:       true,
		affectsModuleResolution: true,
		description:             diagnostics.Control_what_method_is_used_to_detect_module_format_JS_files,
		category:                diagnostics.Language_and_Environment,
		defaultValueDescription: diagnostics.X_auto_Colon_Treat_files_with_imports_exports_import_meta_jsx_with_jsx_Colon_react_jsx_or_esm_format_with_module_Colon_node16_as_modules,
	},
	{
		name:                    "ignoreDeprecations",
		kind:                    "string",
		defaultValueDescription: 0,
	},
}

func getOptionName(option CommandLineOption) string {
	return option.name
}

func commandLineOptionsToMap(options []CommandLineOption) map[string]CommandLineOption {
	// var result = collections.NewOrderedMapFromList([]collections.MapEntry[string, CommandLineOption]{
	// {Key: getOptionName(options[0]), Value: options[0]},
	// {Key: getOptionName(options[1]), Value: options[1]},
	// {Key: getOptionName(options[2]), Value: options[2]},
	// {Key: getOptionName(options[3]), Value: options[3]},
	// {Key: getOptionName(options[4]), Value: options[4]},
	// {Key: getOptionName(options[5]), Value: options[5]},
	// {Key: getOptionName(options[6]), Value: options[6]},
	// {Key: getOptionName(options[7]), Value: options[7]},
	// {Key: getOptionName(options[8]), Value: options[8]},
	// {Key: getOptionName(options[9]), Value: options[9]},
	// {Key: getOptionName(options[10]), Value: options[10]},
	// {Key: getOptionName(options[11]), Value: options[11]},
	// {Key: getOptionName(options[12]), Value: options[12]},
	// {Key: getOptionName(options[13]), Value: options[13]},
	// {Key: getOptionName(options[14]), Value: options[14]},
	// {Key: getOptionName(options[15]), Value: options[15]},
	// {Key: getOptionName(options[16]), Value: options[16]},
	// {Key: getOptionName(options[17]), Value: options[17]},
	// {Key: getOptionName(options[18]), Value: options[18]},
	// {Key: getOptionName(options[19]), Value: options[19]},
	// {Key: getOptionName(options[20]), Value: options[20]},
	// {Key: getOptionName(options[21]), Value: options[21]},
	// {Key: getOptionName(options[22]), Value: options[22]},
	// {Key: getOptionName(options[23]), Value: options[23]},
	// })
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
	commandLineCompilerOptionsMapCache = commandLineOptionsToMap(optionDeclarations)
	return commandLineCompilerOptionsMapCache
}

// var commandLineWatchOptionsMapCache map[string]commandLineOption

// func getCommandLineWatchOptionsMap() {
// 	// if commandLineWatchOptionsMapCache != nil {
// 	// 	return commandLineWatchOptionsMapCache
// 	// }
// 	// commandLineWatchOptionsMapCache = commandLineOptionsToMap(optionsForWatch) //todo need to add watch related options
// 	// return commandLineWatchOptionsMapCache
// }

// func convertCompileOnSaveOptionFromJson(json any, basePath string, errors []ast.Diagnostic) bool {
// 	return false //todo
// }
// func convertWatchOptionsFromJsonWorker(jsonOptions any, basePath string, errors []ast.Diagnostic) compiler.WatchOptions {
// 	//return convertOptionsFromJson(getCommandLineWatchOptionsMap(), jsonOptions, basePath /*defaultOptions*/, undefined, watchOptionsDidYouMeandiagnostics, errors)
// }

// func getDefaultTypeAcquisition(configFileName *string) *compiler.TypeAcquisition {
// 	var options compiler.TypeAcquisition
// 	//options.enable = !!configFileName && getBaseFileName(configFileName) == "jsconfig.json"
// 	options.include = []string{}
// 	options.exclude = []string{}
// 	return options
// }

// func convertTypeAcquisitionFromJsonWorker(jsonOptions any, basePath string, errors []ast.Diagnostic, configFileName *string) compiler.TypeAcquisition {
// 	// const options = getDefaultTypeAcquisition(configFileName)
// 	// convertOptionsFromJson(getCommandLineTypeAcquisitionMap(), jsonOptions, basePath, options, typeAcquisitionDidYouMeandiagnostics, errors)
// 	// return options
// }

// type defaultOptions struct {
// 	core.CompilerOptions
// 	compiler.TypeAcquisition
// 	compiler.WatchOptions
// }

func convertOptionsFromJson(optionsNameMap map[string]CommandLineOption, jsonOptions map[string]interface{}, basePath string, defaultOptions core.CompilerOptions, errors []*ast.Diagnostic) core.CompilerOptions {
	if jsonOptions == nil {
		return core.CompilerOptions{}
	}
	for key, value := range jsonOptions {
		opt, ok := optionsNameMap[key]
		if ok {
			// name := opt.name "lib"
			defaultOptions.Option[opt.name] = convertJsonOption(opt, value, basePath, errors, nil, nil, nil)
		} else {
			//errors.push(createUnknownOptionError(id, diagnostics));
		}
	}
	return defaultOptions
}

// // func getDefaultCompilerOptions(configFileName *string) {
// // 	var options compiler.CompilerOptions
// // 	if configFileName != nil && compiler.GetBaseFileName(*configFileName) == "jsconfig.json" {
// // 		options.allowJs = true
// // 		maxNodeModuleJsDepth = 2
// // 		allowSyntheticDefaultImports = true
// // 		skipLibCheck = true
// // 		noEmit = true
// // 	}
// // 	return options
// // }

// type TsConfigOnlyOption struct {
// 	CommandLineOptionBase
// 	//commandLineOptionBasetype "object";
// 	ElementOptions *map[string]commandLineOption
// 	//extraKeydiagnostics?: DidYouMeanOptionsdiagnostics;
// }

func convertArrayLiteralExpressionToJson(
	elements []*ast.Expression,
	elementOption *tsConfigOnlyOption, //*commandLineOption,
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

// type JsonConversionNotifier struct {
// 	rootOptions TsConfigOnlyOption
// 	// onPropertySet(
// 	//     keyText: string,
// 	//     value: any,
// 	//     propertyAssignment: PropertyAssignment,
// 	//     parentOption: TsConfigOnlyOption | undefined,
// 	//     option: CommandLineOption | undefined,
// 	// ): void;

// }

type optionsBaseValue struct {
	core.CompilerOptionsValue
	*ast.SourceFile
}
type optionsBase struct {
	options map[string]optionsBaseValue
}

// func handleOptionConfigDirTemplateSubstitution(
// 	options optionsBase,
// 	optionDeclarations []CommandLineOption,
// 	basePath string,
// ) {
// 	// if options == (optionsBase{}) { // if options == (optionsBase{}) {
// 	// 	return options
// 	// }
// 	var result optionsBase
// 	for _, option := range optionDeclarations {
// 		if options.options.get(option.name) != nil {
// 			const value = options[option.name]
// 			// switch (option.type) { //todo need to fix option.type prob 11/13/24
// 			//     case "string":
// 			//         Debug.assert(option.isFilePath);
// 			//         if (startsWithConfigDirTemplate(value)) {
// 			//             setOptionValue(option, getSubstitutedPathWithConfigDirTemplate(value, basePath));
// 			//         }
// 			//         break;
// 			//     case "list":
// 			//         Debug.assert(option.element.isFilePath);
// 			//         const listResult = getSubstitutedStringArrayWithConfigDirTemplate(value as string[], basePath);
// 			//         if (listResult) setOptionValue(option, listResult);
// 			//         break;
// 			//     case "object":
// 			//         Debug.assert(option.name === "paths");
// 			//         const objectResult = getSubstitutedMapLikeOfStringArrayWithConfigDirTemplate(value as MapLike<string[]>, basePath);
// 			//         if (objectResult) setOptionValue(option, objectResult);
// 			//         break;
// 			//     default:
// 			//         Debug.fail("option type not supported");
// 			// }
// 		}
// 	}
// 	return result || options

// 	// func setOptionValue(option: CommandLineOption, value: CompilerOptionsValue) {
// 	//     (result ??= assign({}, options))[option.name] = value;
// 	// }
// }

func directoryOfCombinedPath(fileName string, basePath string) string {
	// Use the `getNormalizedAbsolutePath` function to avoid canonicalizing the path, as it must remain noncanonical
	// until consistent casing errors are reported
	return tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(fileName, basePath))
}

// var defaultIncludeSpec = "**/*"

// func validateSpecs(specs []string, errors []ast.Diagnostic, disallowTrailingRecursion bool, jsonSourceFile *ast.SourceFile, specKey string) []string {
// 	// return specs.filter(spec => {
// 	//     if (!isString(spec)) return false;
// 	//     const diag = specToDiagnostic(spec, disallowTrailingRecursion);
// 	//     if (diag !== undefined) {
// 	//         errors.push(createDiagnostic(...diag));
// 	//     }
// 	//     return diag === undefined;
// 	// });

// 	// function createDiagnostic(message: DiagnosticMessage, spec: string): Diagnostic {
// 	//     const element = getTsConfigPropArrayElementValue(jsonSourceFile, specKey, spec);
// 	//     return createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(jsonSourceFile, element, message, spec);
// 	// }
// }

// func getSubstitutedStringArrayWithConfigDirTemplate(list []string, basePath string) {
// 	if list == nil {
// 		return list
// 	}
// 	var result []string
// 	// list.forEach((element, index) => {
// 	//     if (!startsWithConfigDirTemplate(element)) return;
// 	//     (result ??= list.slice())[index] = getSubstitutedPathWithConfigDirTemplate(element, basePath);
// 	// });
// 	return result
// }

// func setConfigFileInOptions(options core.CompilerOptions, configFile *ast.SourceFile) {
// 	// if (configFile) {
// 	//     Object.defineProperty(options, "configFile", { enumerable: false, writable: false, value: configFile });
// 	// }
// }

// /**
//  * Gets the file names from the provided config file specs that contain, files, include, exclude and
//  * other properties needed to resolve the file names
//  * @param configFileSpecs The config file specs extracted with file names to include, wildcards to include/exclude and other details
//  * @param basePath The base path for any relative file specifications.
//  * @param options Compiler options.
//  * @param host The host used to resolve files and directories.
//  * @param extraFileExtensions optionaly file extra file extension information from host
//  *
//  * @internal
//  */
// func getFileNamesFromConfigSpecs(
// 	configFileSpecs configFileSpecs,
// 	basePath string,
// 	options core.CompilerOptions,
// 	host ParseConfigHost,
// 	extraFileExtensions []FileExtensionInfo,
// ) []string {
// 	basePath = tspath.NormalizePath(basePath)

// 	//const keyMapper = createGetCanonicalFileName(host.useCaseSensitiveFileNames);// core.ts

// 	// Literal file names (provided via the "files" array in tsconfig.json) are stored in a
// 	// file map with a possibly case insensitive key. We use this map later when when including
// 	// wildcard paths.
// 	literalFileMap := make(map[string]string)

// 	// Wildcard paths (provided via the "includes" array in tsconfig.json) are stored in a
// 	// file map with a possibly case insensitive key. We use this map to store paths matched
// 	// via wildcard, and to handle extension priority.
// 	//wildcardFileMap := make(map[string]string)

// 	// Wildcard paths of json files (provided via the "includes" array in tsconfig.json) are stored in a
// 	// file map with a possibly case insensitive key. We use this map to store paths matched
// 	// via wildcard of *.json kind
// 	//wildCardJsonFileMap := make(map[string]string)
// 	validatedFilesSpec := configFileSpecs.validatedFilesSpec
// 	validatedIncludeSpecs := configFileSpecs.validatedIncludeSpecs
// 	validatedExcludeSpecs := configFileSpecs.validatedExcludeSpecs

// 	// Rather than re-query this for each file and filespec, we query the supported extensions
// 	// once and store it on the expansion context.
// 	//const supportedExtensions = getSupportedExtensions(options, extraFileExtensions); // in utilities.ts
// 	//const supportedExtensionsWithJsonIfResolveJsonModule = getSupportedExtensionsWithJsonIfResolveJsonModule(options, supportedExtensions)

// 	// Literal files are always included verbatim. An "include" or "exclude" specification cannot
// 	// remove a literal file.
// 	if validatedFilesSpec != nil {
// 		for _, fileName := range validatedFilesSpec {
// 			file := tspath.GetNormalizedAbsolutePath(fileName, basePath)
// 			literalFileMap.set(keyMapper(file), file)
// 		}
// 	}

// 	//var jsonOnlyIncludeRegexes []regexp.Regexp
// 	// if validatedIncludeSpecs && len(validatedIncludeSpecs) > 0 {
// 	// 	for _, file := range host.readDirectory(basePath, flatten(supportedExtensionsWithJsonIfResolveJsonModule), validatedExcludeSpecs, validatedIncludeSpecs /*depth*/, undefined) {
// 	// 		if compiler.FileExtensionIs(file, Extension.Json) {
// 	// 			// Valid only if *.json specified
// 	// 			// if (!jsonOnlyIncludeRegexes) {
// 	// 			//     const includes = validatedIncludeSpecs.filter(s => endsWith(s, Extension.Json));
// 	// 			//     const includeFilePatterns = map(getRegularExpressionsForWildcards(includes, basePath, "files"), pattern => `^${pattern}$`);
// 	// 			//     jsonOnlyIncludeRegexes = includeFilePatterns ? includeFilePatterns.map(pattern => getRegexFromPattern(pattern, host.useCaseSensitiveFileNames)) : emptyArray;
// 	// 			// }
// 	// 			// const includeIndex = findIndex(jsonOnlyIncludeRegexes, re => re.test(file));
// 	// 			// if (includeIndex !== -1) {
// 	// 			//     const key = keyMapper(file);
// 	// 			//     if (!literalFileMap.has(key) && !wildCardJsonFileMap.has(key)) {
// 	// 			//         wildCardJsonFileMap.set(key, file);
// 	// 			//     }
// 	// 			// }
// 	// 			// continue;
// 	// 		}
// 	// 		// If we have already included a literal or wildcard path with a
// 	// 		// higher priority extension, we should skip this file.
// 	// 		//
// 	// 		// This handles cases where we may encounter both <file>.ts and
// 	// 		// <file>.d.ts (or <file>.js if "allowJs" is enabled) in the same
// 	// 		// directory when they are compilation outputs.
// 	// 		// if (hasFileWithHigherPriorityExtension(file, literalFileMap, wildcardFileMap, supportedExtensions, keyMapper)) {
// 	// 		//     continue;
// 	// 		// }

// 	// 		// We may have included a wildcard path with a lower priority
// 	// 		// extension due to the user-defined order of entries in the
// 	// 		// "include" array. If there is a lower priority extension in the
// 	// 		// same directory, we should remove it.
// 	// 		// removeWildcardFilesWithLowerPriorityExtension(file, wildcardFileMap, supportedExtensions, keyMapper);

// 	// 		// const key = keyMapper(file);
// 	// 		// if (!literalFileMap.has(key) && !wildcardFileMap.has(key)) {
// 	// 		//     wildcardFileMap.set(key, file);
// 	// 		// }
// 	// 	}
// 	// }

// 	// const literalFiles = arrayFrom(literalFileMap.values());
// 	// const wildcardFiles = arrayFrom(wildcardFileMap.values());

// 	// return literalFiles.concat(wildcardFiles, arrayFrom(wildCardJsonFileMap.values()));
// }

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

func convertObjectLiteralExpressionToJson(
	returnValue bool,
	node *ast.ObjectLiteralExpression,
	objectOption *tsConfigOnlyOption,
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
			// textOfKey = getTextOfPropertyName(element.Name)
		}
		var keyText = textOfKey //&& unescapeLeadingUnderscores(textOfKey);
		//var option any
		if keyText != nil {
			if objectOption.elementOptions != nil {
				//option = struct{}{}
				//option := objectOption.elementOptions[keyText]
			} else {
				//option = nil
				//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, element.name, diagnostics.Unknown_option_0, keyText));
			}

		}
		// todo
		keyText = element.AsPropertyAssignment().Name().Text()                                                                       //"exclude"
		var value = convertPropertyValueToJson(element.AsPropertyAssignment().Initializer, nil, returnValue, jsonConversionNotifier) // this needs to be element.initializer need to come back
		if keyText != "undefined" {
			if returnValue {
				result[keyText.(string)] = value
			}

			// Notify key value set, if user asked for it
			if jsonConversionNotifier != nil {
				//jsonConversionNotifier.onPropertySet(keyText, value, element, objectOption, option)
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
	var jsonConversionNotifierValue *tsConfigOnlyOption
	if jsonConversionNotifier != nil {
		jsonConversionNotifierValue = &jsonConversionNotifier.rootOptions
	}
	fmt.Println("in convertToJson", rootExpression.Kind)
	return convertPropertyValueToJson(rootExpression, jsonConversionNotifierValue, returnValue, jsonConversionNotifier)
}

func isDoubleQuotedString(node *ast.Node) bool {
	return ast.IsStringLiteral(node) //&& isStringDoubleQuoted(node, sourceFile);
}

func convertPropertyValueToJson(valueExpression *ast.Expression, option *tsConfigOnlyOption, returnValue bool, jsonConversionNotifier *JsonConversionNotifier) any {
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
func ParseJsonConfigFileContent(json map[string]interface{}, host ParseConfigHost, basePath string, existingOptions *core.CompilerOptions, configFileName *string, resolutionStack *[]tspath.Path, extraFileExtensions *[]FileExtensionInfo, extendedConfigCache *map[string]ExtendedConfigCacheEntry) module.ParsedCommandLine {
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
	} else {
		options = core.CompilerOptions{}
	}
	return options
}

type defaultOptions struct {
	core.CompilerOptions
	//TypeAcquisition
	//WatchOptions
}

// func convertOptionsFromJson(optionsNameMap map[string]CommandLineOption, jsonOptions optionsBase, basePath string, defaultOptions defaultOptions, diagnostics diagnostics, errors []diagnostics) { //diagnostics DidYouMeanOptionsDiagnostics
//     if (!jsonOptions) {
//         return;
//     }

//     for _ , id := range jsonOptions.options {
//         const opt = optionsNameMap.get(id.booleanValue);
//         if (opt) {
//             (defaultOptions || (defaultOptions = {}))[opt.name] = convertJsonOption(opt, jsonOptions[id], basePath, errors);
//         }
//         else {
//             errors.push(createUnknownOptionError(id, diagnostics));
//         }
//     }
//     return defaultOptions;
// }

func convertCompilerOptionsFromJsonWorker(jsonOptions map[string]interface{}, basePath string, errors []*ast.Diagnostic, configFileName string) core.CompilerOptions {
	options := getDefaultCompilerOptions(configFileName)
	convertOptionsFromJson(getCommandLineCompilerOptionsMap(), jsonOptions, basePath, options, errors)
	if configFileName != "" {
		options.ConfigFilePath = tspath.NormalizeSlashes(configFileName)
	}
	return options
}

func parseOwnConfigOfJson(
	json map[string]interface{},
	host ParseConfigHost,
	basePath string,
	configFileName string,
	errors []*ast.Diagnostic,
) *ParsedTsconfig {
	// if (hasProperty(json, "excludes")) {
	//     errors.push(createCompilerDiagnostic(diagnostics.Unknown_option_excludes_Did_you_mean_exclude));
	// }
	var options core.CompilerOptions
	for k, v := range json {
		if k == "compilerOptions" {
			options = convertCompilerOptionsFromJsonWorker(v.(map[string]interface{}), basePath, errors, configFileName)
		}
	}
	fmt.Println("options", options)
	// typeAcquisition := convertTypeAcquisitionFromJsonWorker(json.typeAcquisition, basePath, errors, configFileName)
	// watchOptions := convertWatchOptionsFromJsonWorker(json.watchOptions, basePath, errors)
	// json.compileOnSave = convertCompileOnSaveOptionFromJson(json, basePath, errors)
	// var extendedConfigPath string
	// if json.extends != nil || json.extends == "" {
	// 	extendedConfigPath = getExtendsConfigPathOrArray(json.extends, host, basePath, configFileName, errors)
	// }
	var parsedConfig = &ParsedTsconfig{
		raw:     json,
		options: &options,
	}
	return parsedConfig
}

/**
 * This *just* extracts options/include/exclude/files out of a config file.
 * It does *not* resolve the included files.
 */
func parseConfig(
	json map[string]interface{},
	sourceFile *ast.SourceFile,
	host ParseConfigHost,
	basePath string,
	configFileName string,
	resolutionStack []string,
	errors []*ast.Diagnostic,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
) ParsedTsconfig {
	basePath = tspath.NormalizeSlashes(basePath)
	resolvedPath := tspath.GetNormalizedAbsolutePath(configFileName, basePath)

	if slices.Contains(resolutionStack, resolvedPath) {
		var result *ParsedTsconfig
		//errors = append(errors, compiler.NewDiagnostic(resolvePath, 0, 0, compiler.DiagnosticCategory.Error, "Circularity detected while resolving configuration file "+resolvePath))
		if json != nil {
			result.raw = json
			return *result
		} else {
			convertToObject(sourceFile, errors)
		}
	}
	var ownConfig *ParsedTsconfig
	if json != nil {
		ownConfig = parseOwnConfigOfJson(json, host, basePath, configFileName, errors)
	}
	fmt.Println("ownConfig", ownConfig)
	// else { //!!!
	// 	parseOwnConfigOfJsonSourceFile(sourceFile, host, basePath, &configFileName, errors)
	// }

	applyExtendedConfig := func(result extendsResult, extendedConfigPath []string) { //here2
		extendedConfig := getExtendedConfig(sourceFile, extendedConfigPath[0], host, resolutionStack, errors, extendedConfigCache, result) //check
		if extendedConfig != (ParsedTsconfig{}) && isSuccessfulParsedTsconfig(extendedConfig) {
			// extendsRaw := extendedConfig.raw
			// var relativeDifference string
			// setPropertyInResultIfNotUndefined := func(propertyName string) {
			// 	if ownConfig.raw != nil { // ownConfig.raw[propertyName] {
			// 		return
			// 	} // No need to calculate if already set in own config
			// 	// if extendsRaw != nil {
			// 	// 	fn := func(path string) {
			// 	// 		if startsWithConfigDirTemplate(path) || isRootedDiskPath(path) {
			// 	// 			return path
			// 	// 		} else if relativeDifference == "" {
			// 	// 			relativeDifference = compiler.ConvertToRelativePath(compiler.GetDirectoryPath(extendedConfigPath), basePath, host.useCaseSensitiveFileNames)
			// 	// 			return compiler.CombinePaths(relativeDifference, path)
			// 	// 		}
			// 	// 	}
			// 	// 	result[propertyName] = map[extendsRaw[propertyName]]fn
			// 	// }
			// }
			// setPropertyInResultIfNotUndefined("include")
			// setPropertyInResultIfNotUndefined("exclude")
			// setPropertyInResultIfNotUndefined("files")
			// if extendsRaw.compileOnSave != nil {
			// 	result.compileOnSave = extendsRaw.compileOnSave
			// }
			// assign(result.options, extendedConfig.options) //assign is a function in core.ts
			// result.watchOptions = result.watchOptions && extendedConfig.watchOptions ?
			// 	assignWatchOptions(result, extendedConfig.watchOptions) :
			// 	result.watchOptions || extendedConfig.watchOptions;
			// TODO extend type typeAcquisition
		}
		// function assignWatchOptions(result: ExtendsResult, watchOptions: WatchOptions) {
		// 	if (result.watchOptionsCopied) return assign(result.watchOptions!, watchOptions);
		// 	result.watchOptionsCopied = true;
		// 	return assign({}, result.watchOptions, watchOptions);
		// }
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
		for _, extendedConfigPath := range *ownConfig.extendedConfigPath {
			applyExtendedConfig(result, []string{extendedConfigPath})
		}
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

	return *ownConfig
}

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
	sourceFile *ast.SourceFile,
	host ParseConfigHost,
	basePath string,
	existingOptions *core.CompilerOptions, //should default to an empty object
	//existingWatchOptions compiler.WatchOptions,
	configFileName *string,
	resolutionStack *[]tspath.Path,
	extraFileExtensions *[]FileExtensionInfo,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
) module.ParsedCommandLine {
	//Debug.assert((json === undefined && sourceFile !== undefined) || (json !== undefined && sourceFile === undefined));
	var errors []*ast.Diagnostic = []*ast.Diagnostic{}
	resolutionStackString := []string{}
	parsedConfig := parseConfig(json, sourceFile, host, basePath, *configFileName, resolutionStackString, errors, extendedConfigCache)
	fmt.Println("parsedConfig", parsedConfig)
	//var raw = parsedConfig.raw
	// const options = handleOptionConfigDirTemplateSubstitution(
	// 	extend(existingOptions, parsedConfig.options), //function in core.ts
	// 	configDirTemplateSubstitutionOptions,
	// 	basePath,
	// )
	// // const watchOptions = handleWatchOptionsConfigDirTemplateSubstitution(
	// //     existingWatchOptions && parsedConfig.watchOptions ?
	// //         extend(existingWatchOptions, parsedConfig.watchOptions) :
	// //         parsedConfig.watchOptions || existingWatchOptions,
	// //     basePath,
	// // );
	var options = parsedConfig.options
	if *configFileName != "" {
		options.ConfigFilePath = *configFileName
	} else {
		options.ConfigFilePath = tspath.NormalizeSlashes(*configFileName)
	}
	// var basePathForFileNames string
	// if configFileName != nil {
	// 	basePathForFileNames = tspath.NormalizePath(directoryOfCombinedPath(*configFileName, basePath))
	// } else {
	// 	basePathForFileNames = tspath.NormalizePath(basePath)
	// }

	// type validateElement func(value any) bool
	// type propOfRaw[T any] struct {
	// 	array    *[]T
	// 	notArray *string
	// 	noProp   *string
	// }
	// getPropFromRaw := func(prop string, validate validateElement, elementTypeName string) propOfRaw {
	// 	// if (hasProperty(raw, prop) && !isNullOrUndefined(raw[prop])) { hasProperty is a function in core.ts
	// 	//     if (isArray(raw[prop])) {
	// 	//         const result = raw[prop] as T[];
	// 	//         if (!sourceFile && !every(result, validateElement)) {
	// 	//             errors.push(createCompilerDiagnostic(diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop, elementTypeName));
	// 	//         }
	// 	//         return result;
	// 	//     }
	// 	//     else {
	// 	//         createCompilerDiagnosticOnlyIfJson(diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop, "Array");
	// 	//         return "not-array";
	// 	//     }
	// 	// }
	// 	return "no-prop"
	// }

	// toPropValue := func(specResult propOfRaw) {
	// 	if specResult.array != nil {
	// 		return specResult
	// 	}
	// }

	// getSpecsFromRaw := func(prop string) propOfRaw { //prop: "files" | "include" | "exclude"
	// 	return getPropFromRaw(prop, isString, "string")
	// }

	// getFileNames := func(basePath string) []string {
	// 	var fileNames = getFileNamesFromConfigSpecs(configFileSpecs, basePath, options, host, extraFileExtensions)
	// 	// if shouldReportNoInputFiles(fileNames, canJsonReportNoInputFiles(raw), resolutionStack) {
	// 	// 	errors.push(getErrorForNoInputFiles(configFileSpecs, configFileName))
	// 	// }
	// 	return fileNames
	// }

	// getProjectReferences := func(basePath string) []compiler.ProjectReference {
	// 	var projectReferences = []compiler.ProjectReference{}
	// 	const referencesOfRaw = getPropFromRaw("references", validateElement("onject"), "object")
	// 	if compiler.IsSlice(referencesOfRaw) {
	// 		for _, ref := range referencesOfRaw {
	// 			if ref.path != "string" { //typeof ref.path !== "string"
	// 				//createCompilerDiagnosticOnlyIfJson(diagnostics.Compiler_option_0_requires_a_value_of_type_1, "reference.path", "string");
	// 			} else {
	// 				projectReferences = append(projectReferences, compiler.ProjectReference{
	// 					path:         tspath.getNormalizedAbsolutePath(ref.path, basePath),
	// 					originalPath: ref.path,
	// 					prepend:      ref.prepend,
	// 					circular:     ref.circular,
	// 				})
	// 			}
	// 		}
	// 	}
	// 	return projectReferences
	// }

	// createCompilerDiagnosticOnlyIfJson := func(message []diagnostics.Message, args compiler.DiagnosticAndArguments) { //todo full
	// 	// if (!sourceFile) {
	// 	//     errors.push(createCompilerDiagnostic(message, ...args));
	// 	// }
	// }

	// getConfigFileSpecs := func() configFileSpecs {
	// 	referencesOfRaw := getPropFromRaw("references", validateElement(), "object") // come back to validateElement
	// 	filesSpecs := toPropValue(getSpecsFromRaw("files"))
	// 	if filesSpecs {
	// 		hasZeroOrNoReferences := referencesOfRaw == "no-prop" || compiler.IsSlice(referencesOfRaw) && referencesOfRaw.length == 0
	// 		hasExtends := hasProperty(raw, "extends") //hasProperty is a function in core.ts
	// 		if filesSpecs.length == 0 && hasZeroOrNoReferences && !hasExtends {
	// 			if sourceFile {
	// 				fileName := configFileName || "tsconfig.json"
	// 				//diagnosticMessage := diagnostics.The_files_list_in_config_file_0_is_empty;
	// 				//nodeValue := forEachTsConfigPropArray(sourceFile, "files", property => property.initializer);
	// 				//const error = createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, nodeValue, diagnosticMessage, fileName);
	// 				//errors.push(error);
	// 			} else {
	// 				//createCompilerDiagnosticOnlyIfJson(diagnostics.The_files_list_in_config_file_0_is_empty, configFileName || "tsconfig.json");
	// 			}
	// 		}
	// 	}

	// 	includeSpecs := toPropValue(getSpecsFromRaw("include"))

	// 	excludeOfRaw := getSpecsFromRaw("exclude")
	// 	isDefaultIncludeSpec := false
	// 	excludeSpecs := toPropValue(excludeOfRaw)
	// 	if excludeOfRaw == "no-prop" {
	// 		outDir := options.outDir
	// 		declarationDir := options.declarationDir

	// 		if outDir || declarationDir {
	// 			//excludeSpecs = filter([outDir, declarationDir], d => !!d) as string[];//filter is function in core.ts
	// 		}
	// 	}

	// 	if filesSpecs == nil && includeSpecs == nil {
	// 		includeSpecs = []string{defaultIncludeSpec}
	// 		isDefaultIncludeSpec = true
	// 	}
	// 	var validatedIncludeSpecsBeforeSubstitution []string
	// 	var alidatedExcludeSpecsBeforeSubstitution []string
	// 	var validatedIncludeSpecs []string
	// 	var validatedExcludeSpecs []string

	// 	// The exclude spec list is converted into a regular expression, which allows us to quickly
	// 	// test whether a file or directory should be excluded before recursively traversing the
	// 	// file system.

	// 	if includeSpecs {
	// 		validatedIncludeSpecsBeforeSubstitution = validateSpecs(includeSpecs, errors /*disallowTrailingRecursion*/, true, sourceFile, "include")
	// 		validatedIncludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
	// 			validatedIncludeSpecsBeforeSubstitution,
	// 			basePathForFileNames,
	// 		) || validatedIncludeSpecsBeforeSubstitution
	// 	}

	// 	if excludeSpecs {
	// 		validatedExcludeSpecsBeforeSubstitution = validateSpecs(excludeSpecs, errors /*disallowTrailingRecursion*/, false, sourceFile, "exclude")
	// 		validatedExcludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
	// 			validatedExcludeSpecsBeforeSubstitution,
	// 			basePathForFileNames,
	// 		) || validatedExcludeSpecsBeforeSubstitution
	// 	}

	// 	//validatedFilesSpecBeforeSubstitution := filter(filesSpecs, isString) //filter is a function in core.ts
	// 	validatedFilesSpec := getSubstitutedStringArrayWithConfigDirTemplate(
	// 		validatedFilesSpecBeforeSubstitution,
	// 		basePathForFileNames,
	// 	) || validatedFilesSpecBeforeSubstitution

	// 	return configFileSpecs{
	// 		filesSpecs,
	// 		includeSpecs,
	// 		excludeSpecs,
	// 		validatedFilesSpec,
	// 		validatedIncludeSpecs,
	// 		validatedExcludeSpecs,
	// 		validatedFilesSpecBeforeSubstitution,
	// 		validatedIncludeSpecsBeforeSubstitution,
	// 		validatedExcludeSpecsBeforeSubstitution,
	// 		isDefaultIncludeSpec,
	// 	}
	// }

	// configFileSpecs := getConfigFileSpecs()
	// if sourceFile {
	// 	sourceFile.configFileSpecs = configFileSpecs
	// }
	// setConfigFileInOptions(options, sourceFile)
	// // result := compiler.ParsedCommandLine {
	// //     options: options,
	// //     watchOptions: watchOptions,
	// //     fileNames: getFileNames(basePathForFileNames),
	// //     projectReferences: getProjectReferences(basePathForFileNames),
	// //     typeAcquisition: parsedConfig.typeAcquisition || getDefaultTypeAcquisition(),
	// //     raw: raw,
	// //     errors: errors,
	// //     // Wildcard directories (provided as part of a wildcard path) are stored in a
	// //     // file map that marks whether it was a regular wildcard match (with a `*` or `?` token),
	// //     // or a recursive directory. This information is used by filesystem watchers to monitor for
	// //     // new entries in these paths.
	// //     wildcardDirectories: getWildcardDirectories(configFileSpecs, basePathForFileNames, host.useCaseSensitiveFileNames),
	// //     compileOnSave: !!raw.compileOnSave,
	// // }

	var t = module.ParsedCommandLine{
		Options: options,
	}
	return t
}
