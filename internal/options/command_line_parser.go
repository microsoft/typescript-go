package tsoptions

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

// type commandLineOptionOfTypes struct { // new - this is conbining commandLineOptionOfCustomType, commandLineOptionOfStringType, commandLineOptionOfNumberType, commandLineOptionOfBooleanType, tsConfigOnlyOption
// 	CommandLineOptionBase
// 	defaultValueDescription   defaultValueDescriptionType
// 	deprecatedKeys            *map[string]bool
// 	commandLineOptionBasetype *optionType
// 	elementOptions            *map[string]CommandLineOption
// 	//extraKeydiagnostics *DidYouMeanOptionsdiagnostics;
// }

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
// type CommandLineOptionKind string

// const (
// 	CommandLineOptionTypeString        CommandLineOptionKind = "string"
// 	CommandLineOptionTypeNumber        CommandLineOptionKind = "number"
// 	CommandLineOptionTypeBoolean       CommandLineOptionKind = "boolean"
// 	CommandLineOptionTypeObject        CommandLineOptionKind = "object"
// 	CommandLineOptionTypeList          CommandLineOptionKind = "list"
// 	CommandLineOptionTypeListOrElement CommandLineOptionKind = "listOrElement"
// 	CommandLineOptionTypeEnum          CommandLineOptionKind = "enum" //map
// )

// type CommandLineOption struct {
// 	kind            CommandLineOptionKind
// 	name, shortName string
// 	paramType       diagnostics.Message
// 	// used in parsing
// 	isFilePath        bool
// 	isTSConfigOnly    bool
// 	isCommandLineOnly bool

// 	// used in output
// 	description              *diagnostics.Message
// 	defaultValueDescription  any
// 	showInSimplifiedHelpView bool

// 	// used in output in serializing and generate tsconfig
// 	category *diagnostics.Message

// 	// defined once
// 	extraValidation *func(value core.CompilerOptionsValue) (d *diagnostics.Message, args []string)

// 	// true or undefined
// 	// used for configDirTemplateSubstitutionOptions
// 	allowConfigDirTemplateSubstitution,
// 	// used for filter in compilerrunner
// 	affectsDeclarationPath,
// 	affectsProgramStructure,
// 	affectsSemanticdiagnostics,
// 	affectsBuildInfo,
// 	affectsBinddiagnostics,
// 	affectsSourceFile,
// 	affectsModuleResolution,
// 	affectsEmit,

// 	allowJsFlag,
// 	strictFlag bool

// 	// transpileoptions worker
// 	transpileOptionValue core.Tristate
// 	// options[option.name] = option.transpileOptionValue;

// 	// used in listtype
// 	listPreserveFalsyValues bool
// 	disallowNullOrUndefined bool
// }

// // CommandLineOption.Elements()
// var commandLineOptionElements = map[string]*CommandLineOption{
// 	"lib": {
// 		name:                    "lib",
// 		kind:                    CommandLineOptionTypeEnum, // libMap,
// 		defaultValueDescription: core.TSUnknown,
// 	},
// 	"rootDirs": {
// 		name:       "rootDirs",
// 		kind:       CommandLineOptionTypeString,
// 		isFilePath: true,
// 	},
// 	"typeRoots": {
// 		name:       "typeRoots",
// 		kind:       CommandLineOptionTypeString,
// 		isFilePath: true,
// 	},
// 	"types": {
// 		name: "types",
// 		kind: CommandLineOptionTypeString,
// 	},
// 	"moduleSuffixes": {
// 		name: "suffix",
// 		kind: CommandLineOptionTypeString,
// 	},
// 	"customConditions": {
// 		name: "condition",
// 		kind: CommandLineOptionTypeString,
// 	},
// 	"plugins": {
// 		name: "plugin",
// 		kind: CommandLineOptionTypeObject,
// 	},
// }

// func (option *CommandLineOption) Elements() *CommandLineOption {
// 	if option.kind != CommandLineOptionTypeList && option.kind != CommandLineOptionTypeListOrElement {
// 		return nil
// 	}
// 	return commandLineOptionElements[option.name]
// }

// ***********************************************************************//

//var optionDeclarations []CommandLineOption = append(commonOptionsWithBuild, commandOptionsWithoutBuild...)

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

// var compilerOptionsDeclaration commandLineOptionOfTypes = commandLineOptionOfTypes{
// 	CommandLineOptionBase: CommandLineOptionBase{
// 		name:                      "compilerOptions",
// 		commandLineOptionBasetype: commandLineOptionBaseType{optionType: "object"},
// 	},
// 	//elementOptions: getCommandLineCompilerOptionsMap(),
// 	//extraKeydiagnostics: compilerOptionsDidYouMeandiagnostics,
// }

//	var typeAcquisitionDeclaration commandLineOptionOfTypes = commandLineOptionOfTypes{
//		CommandLineOptionBase: CommandLineOptionBase{
//			name:                      "typeAcquisition",
//			commandLineOptionBasetype: commandLineOptionBaseType{optionType: "object"},
//		},
//		// elementOptions: getCommandLineTypeAcquisitionMap(),
//		// extraKeydiagnostics: typeAcquisitionDidYouMeandiagnostics,
//	}
var tsconfigRootOptions *CommandLineOption //TsConfigOnlyOption

// func getTsconfigRootOptionsMap() CommandLineOption{ //TsConfigOnlyOption
// 	if tsconfigRootOptions == nil {
// 	    tsconfigRootOptions = CommandLineOption{
// 			CommandLineOptionBase: {

// 			},
// 	        //name: , // should never be needed since this is root
// 	        type: "object",
// 	        elementOptions: commandLineOptionsToMap([
// 	            compilerOptionsDeclaration,
// 	            watchOptionsDeclaration,
// 	            typeAcquisitionDeclaration,
// 	            extendsOptionDeclaration,
// 	            {
// 	                name: "references",
// 	                type: "list",
// 	                element: {
// 	                    name: "references",
// 	                    type: "object",
// 	                },
// 	                category: diagnostics.Projects,
// 	            },
// 	            {
// 	                name: "files",
// 	                type: "list",
// 	                element: {
// 	                    name: "files",
// 	                    type: "string",
// 	                },
// 	                category: diagnostics.File_Management,
// 	            },
// 	            {
// 	                name: "include",
// 	                type: "list",
// 	                element: {
// 	                    name: "include",
// 	                    type: "string",
// 	                },
// 	                category: diagnostics.File_Management,
// 	                defaultValueDescription: diagnostics.if_files_is_specified_otherwise_Asterisk_Asterisk_Slash_Asterisk,
// 	            },
// 	            {
// 	                name: "exclude",
// 	                type: "list",
// 	                element: {
// 	                    name: "exclude",
// 	                    type: "string",
// 	                },
// 	                category: diagnostics.File_Management,
// 	                defaultValueDescription: diagnostics.node_modules_bower_components_jspm_packages_plus_the_value_of_outDir_if_one_is_specified,
// 	            },
// 	            compileOnSaveCommandLineOption,
// 	        ]),
// 	    };
// 	} //todo
// 	return tsconfigRootOptions
// 	return nil
// }

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

func parseOwnConfigOfJsonSourceFile(
	sourceFile *ast.SourceFile,
	host ParseConfigHost,
	basePath string,
	configFileName *string,
	errors []*ast.Diagnostic,
) *ParsedTsconfig {
	options := getDefaultCompilerOptions(*configFileName)
	//var typeAcquisition *compiler.TypeAcquisition
	//var watchOptions *compiler.WatchOptions
	//var extendedConfigPath []string = []string{} // | string
	// var rootCompilerOptions []ast.PropertyName

	// rootOptions := getTsconfigRootOptionsMap()
	//var conversionNotifier JsonConversionNotifier

	// onPropertySet := func(
	// 	keyText string,
	// 	value any,
	// 	propertyAssignment ast.PropertyAssignment,
	// 	parentOption *commandLineOptionOfTypes, //TsConfigOnlyOption,
	// 	option *commandLineOption,
	// ) {
	// 	// Ensure value is verified except for extends which is handled in its own way for error reporting
	// 	if option != nil { //&& option != extendsOptionDeclaration {
	// 		value = convertJsonOption(*option, value, basePath, errors, &propertyAssignment, propertyAssignment.Initializer, sourceFile)
	// 	}
	// 	if parentOption.name != "" {
	// 		if option != nil {
	// 			//var currentOption
	// 			if parentOption == &compilerOptionsDeclaration {
	// 				// currentOption := options
	// 				// } else if parentOption == &watchOptionsDeclaration {
	// 				// 	if !watchOptions { //if watchOptions is null or undefined
	// 				// 		currentOption = watchOptions
	// 				// 	}
	// 				// }
	// 			} else if parentOption == &typeAcquisitionDeclaration {
	// 				// if typeAcquisition != nil {
	// 				// 	typeAcquisition = getDefaultTypeAcquisition(configFileName)
	// 				// }
	// 				//currentOption := typeAcquisition
	// 			} //else Debug.fail("Unknown option");
	// 			//currentOption := value //*currentOption[option] = value find a way to do this
	// 		} else if keyText != "" && parentOption != nil { //&& parentOption.extraKeydiagnostics {
	// 			if parentOption.elementOptions != nil {
	// 				// errors.push(createUnknownOptionError(
	// 				// 	keyText,
	// 				// 	parentOption.extraKeydiagnostics,
	// 				// 	/*unknownOptionErrorText*/ undefined,
	// 				// 	propertyAssignment.name,
	// 				// 	sourceFile,
	// 				// ));
	// 			}
	// 			// else {
	// 			//     errors.push(createDiagnosticForNodeInSourceFile(sourceFile, propertyAssignment.name, parentOption.extraKeydiagnostics.unknownOptionDiagnostic, keyText));
	// 			// }
	// 		}
	// 	} else if parentOption == rootOptions {
	// 		// t := option //here need to fix
	// 		// if option.CommandLineOptionOfListType == extendsOptionDeclaration {
	// 		// 	extendedConfigPath = getExtendsConfigPathOrArray(value, host, basePath, configFileName, errors, propertyAssignment, propertyAssignment.initializer, sourceFile)
	// 		// } else if !option {
	// 		// 	if keyText == "excludes" {
	// 		// 		//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, propertyAssignment.name, diagnostics.Unknown_option_excludes_Did_you_mean_exclude));
	// 		// 	}
	// 		// 	// if (compiler.Find(commandOptionsWithoutBuild, opt => opt.name === keyText)) {
	// 		// 	//     rootCompilerOptions = append(rootCompilerOptions, propertyAssignment.name);
	// 		// 	// }
	// 		// }
	// 	}
	// }

	// json := convertConfigFileToObject(
	// 	sourceFile,
	// 	errors,
	// 	JsonConversionNotifier{rootOptions, onPropertySet},
	// )
	// if typeAcquisition == nil {
	// 	typeAcquisition = getDefaultTypeAcquisition(configFileName)
	// }

	// if rootCompilerOptions != nil && json != nil { //&& json.compilerOptions == nil {
	// 	//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, rootCompilerOptions[0], diagnostics._0_should_be_set_inside_the_compilerOptions_object_of_the_config_json_file, getTextOfPropertyName(rootCompilerOptions[0]) as string));
	// }

	return &ParsedTsconfig{
		// raw:     json,
		options: &options,
		//watchOptions:    watchOptions,
		// typeAcquisition: typeAcquisition,
		//extendedConfigPath: extendedConfigPath,
	}

}

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
	rootOptions tsConfigOnlyOption //TsConfigOnlyOption
	//onPropertySet func(keyText string, value any, propertyAssignment ast.PropertyAssignment, parentOption commandLineOptionOfTypes, option CommandLineOption)
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
		var expectedType = string(option.Kind)
		k := reflect.TypeOf(value)
		fmt.Println(k)
		if CommandLineOptionTypeEnum == option.Kind {
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
	val string, //core.CompilerOptionsValue,
	errors []*ast.Diagnostic,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) any {
	if val == "" {
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
	mappedValue := core.Map(values, func(v string) any {
		if valueExpression != nil {
			expression = valueExpression.AsArrayLiteralExpression().Elements.Nodes[index]
		}
		var t = convertJsonOption(*option.Elements(), v, basePath, errors, propertyAssignment, expression, sourceFile)
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

func normalizeNonListOptionValue(option CommandLineOption, basePath string, value any) core.CompilerOptionsValue {
	if option.isFilePath {
		value = value
		if !startsWithConfigDirTemplate(value) {
			value = tspath.GetNormalizedAbsolutePath(value.(string), basePath)
		}
		//value = !startsWithConfigDirTemplate(value) ? getNormalizedAbsolutePath(value, basePath) : value;
		if value == "" {
			value = "."
		}
	}
	return core.CompilerOptionsValue{StringValue: value.(string)}
}

func convertJsonOption(
	opt CommandLineOption,
	value any,
	basePath string,
	errors []*ast.Diagnostic,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) any {
	if opt.isCommandLineOnly != false {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, propertyAssignment?.name, diagnostics.Option_0_can_only_be_specified_on_command_line, opt.name));
		return core.CompilerOptionsValue{}
	}
	if isCompilerOptionsValue(opt, value).BooleanValue {
		optType := opt.Kind
		_, ok := value.([]string)
		if (optType == "list") && ok {
			list := convertJsonOptionOfListType(opt, value.([]string), basePath, errors, propertyAssignment, valueExpression, sourceFile) //as ArrayLiteralExpression | undefined
			return list
		} else if optType == "listOrElement" {
			if ok {
				return convertJsonOptionOfListType(opt, value.([]string), basePath, errors, propertyAssignment, valueExpression, sourceFile)
			} else {
				return convertJsonOption(*opt.Elements(), value, basePath, errors, propertyAssignment, valueExpression, sourceFile)
			}
		} else if !(opt.Kind == "string") {
			return convertJsonOptionOfCustomType(opt, value.(string), errors, valueExpression, sourceFile)
		}
		validatedValue := validateJsonOptionValue(opt, value.(string), errors, valueExpression, sourceFile)
		if validatedValue != nil {
			return validatedValue
		} else {
			return normalizeNonListOptionValue(opt, basePath, validatedValue)
		}
	} else {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.Compiler_option_0_requires_a_value_of_type_1, opt.name, getCompilerOptionValueTypeString(opt)));
		return core.CompilerOptionsValue{}
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
	default:
		// Handle unknown options
		fmt.Printf("Unknown option: %s\n", key)
	}

	return options
}

func getOptionName(option CommandLineOption) string {
	return option.Name
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
	commandLineCompilerOptionsMapCache = commandLineOptionsToMap(optionsDeclarations)
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
			//lowerCaseKey := capitalizeFirstLetter(opt.Name)
			*defaultOptions = parseCompilerOptions(key, convertJsonOption(opt, value, basePath, errors, nil, nil, nil))
		}
	}
	// defaultOptions.Option = make(map[string]core.CompilerOptionsValue)
	// for key, value := range jsonOptions {
	// 	opt, ok := optionsNameMap[key]
	// 	if ok {
	// 		// name := opt.name "lib"
	// 		fmt.Println(capitalizeFirstLetter(opt.Name))
	// 		x := convertJsonOption(opt, value, basePath, errors, nil, nil, nil)
	// 		defaultOptions.Lib = parseCompilerOptions(x)
	// 		defaultOptions.Option[capitalizeFirstLetter(opt.Name)] = x
	// 	} else {
	// 		//errors.push(createUnknownOptionError(id, diagnostics));
	// 	}
	// }
	return *defaultOptions
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
	convertOptionsFromJson(getCommandLineCompilerOptionsMap(), jsonOptions, basePath, &options, errors)
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
