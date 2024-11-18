package parsingoptions

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// func ParseConfigFileTextToJson(fileName string, jsonText string) (any, Diagnostic) {
// 	jsonSourceFile := ParseJSONText(fileName, jsonText)
// 	// Parse the config file text into a JSON string
// 	return "", nil
// }

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
	description               *compiler.Diagnostic      //DM			                       // The message describing what the command line switch does.
	//defaultValueDescription?: string | number | boolean | DiagnosticMessage | undefined;   // The message describing what the dafault value is. string type is prepared for fixed chosen like "false" which do not need I18n.
	paramType                  *compiler.Diagnostic //DM                       // The name to be used for a non-boolean option's parameter
	isTSConfigOnly             *bool                // True if option can only be specified via tsconfig.json file
	isCommandLineOnly          *bool
	showInSimplifiedHelpView   *bool
	category                   *compiler.Diagnostic //DM
	strictFlag                 *bool                // true if the option is one of the flag under strict
	allowJsFlag                *bool
	affectsSourceFile          *bool // true if we should recreate SourceFiles after this option changes
	affectsModuleResolution    *bool // currently same effect as `affectsSourceFile`
	affectsBindDiagnostics     *bool // true if this affects binding (currently same effect as `affectsSourceFile`)
	affectsSemanticDiagnostics *bool // true if option affects semantic diagnostics
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
	elementOptions            *map[string]commandLineOption
	//extraKeyDiagnostics *DidYouMeanOptionsDiagnostics;
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

type CommandLineOptionOfListType struct { //new changes it a little
	CommandLineOptionBase
	element                 commandLineOptionOfTypes
	listPreserveFalsyValues *bool
}
type commandLineOption struct {
	CommandLineOptionOfListType
	commandLineOptionOfTypes
}

// type tsConfigOnlyOption struct {
// 	CommandLineOptionBase
// 	commandLineOptionBasetype optionType //"object";
// 	elementOptions            *map[string]commandLineOption
// 	//extraKeyDiagnostics?: DidYouMeanOptionsDiagnostics;
// }

type extendsResult struct {
	options             core.CompilerOptions
	watchOptions        compiler.WatchOptions
	watchOptionsCopied  bool
	include             *[]string
	exclude             *[]string
	files               *[]string
	compileOnSave       *bool
	extendedSourceFiles *map[string]struct{} //*Set<string>;
}

var extendsOptionDeclaration CommandLineOptionOfListType = CommandLineOptionOfListType{
	CommandLineOptionBase: CommandLineOptionBase{
		name:                      "extends",
		commandLineOptionBasetype: commandLineOptionBaseType{optionType: "listOrElement"},
		//category: compiler.Diagnostics.File_Management,
		disallowNullOrUndefined: func(b bool) *bool { return &b }(true), //need to check this
	},
	element: commandLineOptionOfTypes{
		CommandLineOptionBase: CommandLineOptionBase{
			name:                      "extends",
			commandLineOptionBasetype: commandLineOptionBaseType{optionType: "string"},
		},
	},
}

var compilerOptionsDeclaration commandLineOptionOfTypes = commandLineOptionOfTypes{
	CommandLineOptionBase: CommandLineOptionBase{
		name:                      "compilerOptions",
		commandLineOptionBasetype: commandLineOptionBaseType{optionType: "object"},
	},
	//elementOptions: getCommandLineCompilerOptionsMap(),
	//extraKeyDiagnostics: compilerOptionsDidYouMeanDiagnostics,
}

var typeAcquisitionDeclaration commandLineOptionOfTypes = commandLineOptionOfTypes{
	CommandLineOptionBase: CommandLineOptionBase{
		name:                      "typeAcquisition",
		commandLineOptionBasetype: commandLineOptionBaseType{optionType: "object"},
	},
	// elementOptions: getCommandLineTypeAcquisitionMap(),
	// extraKeyDiagnostics: typeAcquisitionDidYouMeanDiagnostics,
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
	//                 category: Diagnostics.Projects,
	//             },
	//             {
	//                 name: "files",
	//                 type: "list",
	//                 element: {
	//                     name: "files",
	//                     type: "string",
	//                 },
	//                 category: Diagnostics.File_Management,
	//             },
	//             {
	//                 name: "include",
	//                 type: "list",
	//                 element: {
	//                     name: "include",
	//                     type: "string",
	//                 },
	//                 category: Diagnostics.File_Management,
	//                 defaultValueDescription: Diagnostics.if_files_is_specified_otherwise_Asterisk_Asterisk_Slash_Asterisk,
	//             },
	//             {
	//                 name: "exclude",
	//                 type: "list",
	//                 element: {
	//                     name: "exclude",
	//                     type: "string",
	//                 },
	//                 category: Diagnostics.File_Management,
	//                 defaultValueDescription: Diagnostics.node_modules_bower_components_jspm_packages_plus_the_value_of_outDir_if_one_is_specified,
	//             },
	//             compileOnSaveCommandLineOption,
	//         ]),
	//     };
	// } //todo
	return tsconfigRootOptions
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
//     extraKeyDiagnostics: compilerOptionsDidYouMeanDiagnostics,
// };
/**
 * Parse the contents of a config file (tsconfig.json).
 * @param jsonNode The contents of the config file to parse
 * @param host Instance of ParseConfigHost used to enumerate files in folder.
 * @param basePath A root directory to resolve relative path entries in the config
 *    file to. e.g. outDir
 */
type ParseConfigHost struct { // should this be an interface? but how with useCaseSensitiveFileNames bool
	compiler.ModuleResolutionHost
	useCaseSensitiveFileNames bool
}

func (p ParseConfigHost) readDirectory(rootDir string, extensions []string, excludes []string, includes []string, depth int) []string {
	return nil //todo
}

// Gets a value indicating whether the specified path exists and is a file.
func (p ParseConfigHost) fileExists(path string) bool {
	return false //todo
}
func (p ParseConfigHost) readFile(path string) string {
	return "" //todo
}
func (p ParseConfigHost) trace(s *string) {
	return //todo
} //probably?? supposed to be trace?(s: string): void

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
	raw             any
	options         *core.CompilerOptions
	watchOptions    *compiler.WatchOptions
	typeAcquisition *compiler.TypeAcquisition
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
func ParseJsonSourceFileConfigFileContent(
	sourceFile *ast.SourceFile,
	host ParseConfigHost,
	basePath string,
	existingOptions *core.CompilerOptions,
	configFileName *string,
	resolutionStack *[]tspath.Path,
	extraFileExtenstions *[]FileExtensionInfo,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
	existingWatchOptions *compiler.WatchOptions) *compiler.ParsedCommandLine {
	//tracing?.push(tracing.Phase.Parse, "parseJsonSourceFileConfigFileContent", { path: sourceFile.fileName });
	result := parseJsonConfigFileContentWorker( /*json*/ nil, sourceFile, host, basePath, *existingOptions, *existingWatchOptions, configFileName, *resolutionStack, *extraFileExtenstions, extendedConfigCache)
	//tracing?.pop();
	return result
}

type readFile func(fileName string) string

func tryReadFile(fileName string, fn readFile) { //return string | Diagnostic
	//text := fn(fileName)
	// if err != nil {
	// 	//return createCompilerDiagnostic(Diagnostics.Cannot_read_file_0_Colon_1, fileName, err.Error())
	// }
	// if text == nil {
	// 	//return createCompilerDiagnostic(Diagnostics.Cannot_read_file_0, fileName)
	// }
	//return text
	// catch (e) {
	//     //return createCompilerDiagnostic(Diagnostics.Cannot_read_file_0_Colon_1, fileName, e.message);
	// }
	return
}

func getBaseFileName(path string, extensions *[]string, ignoreCase *bool) string {
	path = tspath.NormalizeSlashes(path)

	// if the path provided is itself the root, then it has not file name.
	rootLength := tspath.GetRootLength(path)
	if rootLength == len(path) {
		return ""
	}

	// return the trailing portion of the path starting after the last (non-terminal) directory
	// separator but not including any trailing directory separator.
	path = tspath.RemoveTrailingDirectorySeparator(path)
	//name :=  path[int(math.Max(float64(compiler.GetRootLength(path)),float64(strings.LastIndex(path,compiler.DirectorySeparator)+1)))]//path.slice(Math.max(getRootLength(path), path.lastIndexOf(directorySeparator) + 1));
	// var extension string
	// if extensions != nil && ignoreCase != nil {
	//     extension = getAnyExtensionFromPath(name, extensions, ignoreCase)
	// }
	// if extension != nil {
	//     return name[0:len(name) - len(extension)]
	// }
	// return name
	return ""
}

func parseOwnConfigOfJsonSourceFile(
	sourceFile *ast.SourceFile,
	host ParseConfigHost,
	basePath string,
	configFileName *string,
	errors []ast.Diagnostic,
) *ParsedTsconfig {
	options := getDefaultCompilerOptions(configFileName)
	var typeAcquisition *compiler.TypeAcquisition
	//var watchOptions *compiler.WatchOptions
	//var extendedConfigPath []string = []string{} // | string
	var rootCompilerOptions []ast.PropertyName

	rootOptions := getTsconfigRootOptionsMap()
	//var conversionNotifier JsonConversionNotifier

	onPropertySet := func(
		keyText string,
		value any,
		propertyAssignment ast.PropertyAssignment,
		parentOption *commandLineOptionOfTypes, //TsConfigOnlyOption,
		option *commandLineOption,
	) {
		// Ensure value is verified except for extends which is handled in its own way for error reporting
		if option != nil { //&& option != extendsOptionDeclaration {
			value = convertJsonOption(option, value, basePath, errors, &propertyAssignment, propertyAssignment.Initializer, sourceFile)
		}
		if parentOption.name != "" {
			if option != nil {
				//var currentOption
				if parentOption == &compilerOptionsDeclaration {
					// currentOption := options
					// } else if parentOption == &watchOptionsDeclaration {
					// 	if !watchOptions { //if watchOptions is null or undefined
					// 		currentOption = watchOptions
					// 	}
					// }
				} else if parentOption == &typeAcquisitionDeclaration {
					if typeAcquisition != nil {
						typeAcquisition = getDefaultTypeAcquisition(configFileName)
					}
					//currentOption := typeAcquisition
				} //else Debug.fail("Unknown option");
				//currentOption := value //*currentOption[option] = value find a way to do this
			} else if keyText != "" && parentOption != nil { //&& parentOption.extraKeyDiagnostics {
				if parentOption.elementOptions != nil {
					// errors.push(createUnknownOptionError(
					// 	keyText,
					// 	parentOption.extraKeyDiagnostics,
					// 	/*unknownOptionErrorText*/ undefined,
					// 	propertyAssignment.name,
					// 	sourceFile,
					// ));
				}
				// else {
				//     errors.push(createDiagnosticForNodeInSourceFile(sourceFile, propertyAssignment.name, parentOption.extraKeyDiagnostics.unknownOptionDiagnostic, keyText));
				// }
			}
		} else if parentOption == rootOptions {
			// t := option //here need to fix
			// if option.CommandLineOptionOfListType == extendsOptionDeclaration {
			// 	extendedConfigPath = getExtendsConfigPathOrArray(value, host, basePath, configFileName, errors, propertyAssignment, propertyAssignment.initializer, sourceFile)
			// } else if !option {
			// 	if keyText == "excludes" {
			// 		//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, propertyAssignment.name, Diagnostics.Unknown_option_excludes_Did_you_mean_exclude));
			// 	}
			// 	// if (compiler.Find(commandOptionsWithoutBuild, opt => opt.name === keyText)) {
			// 	//     rootCompilerOptions = append(rootCompilerOptions, propertyAssignment.name);
			// 	// }
			// }
		}
	}

	json := convertConfigFileToObject(
		sourceFile,
		errors,
		JsonConversionNotifier{rootOptions, onPropertySet},
	)
	if typeAcquisition == nil {
		typeAcquisition = getDefaultTypeAcquisition(configFileName)
	}

	if rootCompilerOptions != nil && json != nil { //&& json.compilerOptions == nil {
		//errors.push(createDiagnosticForNodeInSourceFile(sourceFile, rootCompilerOptions[0], Diagnostics._0_should_be_set_inside_the_compilerOptions_object_of_the_config_json_file, getTextOfPropertyName(rootCompilerOptions[0]) as string));
	}

	return &ParsedTsconfig{
		raw:     json,
		options: &options,
		//watchOptions:    watchOptions,
		typeAcquisition: typeAcquisition,
		//extendedConfigPath: extendedConfigPath,
	}

}

func getExtendedConfig(
	sourceFile *ast.SourceFile,
	extendedConfigPath string,
	host ParseConfigHost,
	resolutionStack []string,
	errors []ast.Diagnostic,
	extendedConfigCache map[string]ExtendedConfigCacheEntry,
	result extendsResult,
) ParsedTsconfig {
	var path string
	if host.useCaseSensitiveFileNames {
		path = extendedConfigPath
	} else {
		//path = toFileNameLowerCase(extendedConfigPath)
	}
	var value ExtendedConfigCacheEntry
	var extendedResult *ast.SourceFile
	var extendedConfig ParsedTsconfig

	value = extendedConfigCache[path]
	if extendedConfigCache != nil && value == (ExtendedConfigCacheEntry{}) {
		extendedResult = value.extendedResult
		extendedConfig = value.extendedConfig
	} else {
		//extendedResult = readJsonConfigFile(extendedConfigPath, host.readFile(path))
		if extendedResult != nil { //parseDiagnostics.length { //come back
			extendedConfig = parseConfig(nil, extendedResult, host, tspath.GetDirectoryPath(extendedConfigPath), getBaseFileName(extendedConfigPath, nil, nil), resolutionStack, errors, &extendedConfigCache)
		}
		if extendedConfigCache != nil {
			extendedConfigCache[path] = ExtendedConfigCacheEntry{extendedResult, extendedConfig}
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
	// if (extendedResult.parseDiagnostics.length) {
	//     //errors.push(...extendedResult.parseDiagnostics);
	//     return undefined;
	// }
	return extendedConfig //extendedConfig!
}

/**
 * This *just* extracts options/include/exclude/files out of a config file.
 * It does *not* resolve the included files.
 */
func parseConfig(
	json any,
	sourceFile *ast.SourceFile,
	host ParseConfigHost,
	basePath string,
	configFileName string,
	resolutionStack []string,
	errors []ast.Diagnostic,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
) ParsedTsconfig {
	basePath = tspath.NormalizeSlashes(basePath)
	resolvePath := tspath.GetNormalizedAbsolutePath(configFileName, basePath)

	if slices.Contains(resolutionStack, resolvePath) {
		var result *ParsedTsconfig
		//errors = append(errors, compiler.NewDiagnostic(resolvePath, 0, 0, compiler.DiagnosticCategory.Error, "Circularity detected while resolving configuration file "+resolvePath))
		if json != nil {
			result.raw = json
			return *result
		} else {
			ConvertToObject(sourceFile, errors) //brb
		}
	}
	var ownConfig *ParsedTsconfig
	if json != nil {
		parseOwnConfigOfJson(json, host, basePath, configFileName, errors)
	} else {
		parseOwnConfigOfJsonSourceFile(sourceFile, host, basePath, &configFileName, errors)
	}

	applyExtendedConfig := func(result extendsResult, extendedConfigPath []string) { //here
		extendedConfig := getExtendedConfig(sourceFile, extendedConfigPath, host, resolutionStack, errors, *extendedConfigCache, result)
		if extendedConfig != (ParsedTsconfig{}) && isSuccessfulParsedTsconfig(extendedConfig) {
			//extendsRaw := extendedConfig//.raw //tofo sourcefile does not have raw
			// var relativeDifference string
			setPropertyInResultIfNotUndefined := func(propertyName string) {
				if ownConfig != nil { // ownConfig.raw[propertyName] {
					return
				} // No need to calculate if already set in own config
				// if extendsRaw[propertyName] {
				// 	fn := func(path string) {
				// 		if startsWithConfigDirTemplate(path) || isRootedDiskPath(path) {
				// 			return path
				// 		} else if relativeDifference == "" {
				// 			relativeDifference = compiler.ConvertToRelativePath(compiler.GetDirectoryPath(extendedConfigPath), basePath, host.useCaseSensitiveFileNames)
				// 			return compiler.CombinePaths(relativeDifference, path)
				// 		}
				// 	}
				// 	result[propertyName] = map[extendsRaw[propertyName]]fn
				// }
			}
			setPropertyInResultIfNotUndefined("include")
			setPropertyInResultIfNotUndefined("exclude")
			setPropertyInResultIfNotUndefined("files")
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

	// if ownConfig.options && ownConfig.options.paths {
	// 	// If we end up needing to resolve relative paths from 'paths' relative to
	// 	// the config file location, we'll need to know where that config file was.
	// 	// Since 'paths' can be inherited from an extended config in another directory,
	// 	// we wouldn't know which directory to use unless we store it here.
	// 	ownConfig.options.pathsBasePath = basePath
	// }
	if ownConfig.extendedConfigPath != nil {
		// copy the resolution stack so it is never reused between branches in potential diamond-problem scenarios.
		resolutionStack = append(resolutionStack, resolvePath) //resolutionStack.concat([resolvedPath]); //here
		result := extendsResult{options: core.CompilerOptions{}}
		if compiler.IsString(ownConfig.extendedConfigPath) {
			applyExtendedConfig(result, *ownConfig.extendedConfigPath)
		} else {
			for _, extendedConfigPath := range *ownConfig.extendedConfigPath {
				applyExtendedConfig(result, []string{extendedConfigPath})
			}
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

		//ownConfig.options = assign(result.options, ownConfig.options);
		// ownConfig.watchOptions = ownConfig.watchOptions && result.watchOptions ?
		//     assignWatchOptions(result, ownConfig.watchOptions) :
		//     ownConfig.watchOptions || result.watchOptions;
	}

	return *ownConfig
}

/**
 * Read tsconfig.json file
 * @param fileName The path to the config file
 */
func readJsonConfigFile(fileName string, fn ParseConfigHost) *ast.SourceFile {
	// const textOrDiagnostic = tryReadFile(fileName, fn.readFile)
	// if compiler.IsString(textOrDiagnostic) {
	// 	return compiler.ParseJSONText(fileName, textOrDiagnostic)
	// }
	// result := compiler.SourceFile{
	// 	fileName: fileName,
	// 	//parseDiagnostics: [textOrDiagnostic]
	// }
	// return result
	return nil
}

type JsonConversionNotifier struct {
	rootOptions   commandLineOptionOfTypes //TsConfigOnlyOption
	onPropertySet func(keyText string, value any, propertyAssignment ast.PropertyAssignment, parentOption commandLineOptionOfTypes, option commandLineOption)
}

type defaultValueDescriptionType struct {
	valueString     string
	valueNumber     int
	valueDiagnostic ast.Diagnostic // todo should be DiagnosticMessage
}

type listType string

const listTypeList = listType("list")
const listTypeListOrElement = listType("listOrElement")

func convertConfigFileToObject(
	sourceFile *ast.SourceFile,
	errors []ast.Diagnostic,
	jsonConversionNotifier JsonConversionNotifier,
) any {
	// t := sourceFile.statements
	// var rootExpression compiler.Expression = sourceFile.statements[0].expression // ???
	// //const rootExpression: Expression | undefined = sourceFile.statements[0]?.expression;
	// if rootExpression && rootExpression.kind != compiler.SyntaxKindObjectLiteralExpression {
	// errors.push(createDiagnosticForNodeInSourceFile(
	//     sourceFile,
	//     rootExpression,
	//     Diagnostics.The_root_value_of_a_0_file_must_be_an_object,
	//     getBaseFileName(sourceFile.fileName) === "jsconfig.json" ? "jsconfig.json" : "tsconfig.json",
	// ));
	// Last-ditch error recovery. Somewhat useful because the JSON parser will recover from some parse errors by
	// synthesizing a top-level array literal expression. There's a reasonable chance the first element of that
	// array is a well-formed configuration object, made into an array element by stray characters.
	// if (isArrayLiteralExpression(rootExpression)) {
	//     const firstObject = find(rootExpression.elements, isObjectLiteralExpression);
	//     if (firstObject) {
	//         return convertToJson(sourceFile, firstObject, errors, /*returnValue*/ true, jsonConversionNotifier);
	//     }
	// }
	// return {}
	// 	return
	// }
	//return convertToJson(sourceFile, rootExpression, errors /*returnValue*/, true, jsonConversionNotifier)
}

type pluginImport struct {
	name string
}
type compilerOptionsValue struct {
	stringValue           string
	numberValue           float64 //number
	booleanValue          bool
	StringArrayValue      []string
	NumberArrayValue      []float64 //number
	MapLikeValue          *map[string][]string
	PluginImportArray     *[]pluginImport
	ProjectReferenceArray *[]compiler.ProjectReference
	NullValue             bool
	UndefinedValue        bool
	//(string | number)[]
}

// func isNullOrUndefined(x any) { // eslint-disable-line no-restricted-syntax
// 	return x == nil // eslint-disable-line no-restricted-syntax
// }

func isCompilerOptionsValue(option commandLineOption, value any) compilerOptionsValue {
	// if option != (commandLineOption{}) {
	// 	if isNullOrUndefined(value) {
	// 		return !option.disallowNullOrUndefined // All options are undefinable/nullable
	// 	}
	// 	if option.commandLineOptionOfCustomType.commandLineOptionType == "list" {
	// 		return isArray(value) //todo fix
	// 	}
	// 	// if (option.type === "listOrElement") { //todo fix
	// 	//     return isArray(value) || isCompilerOptionsValue(option.element, value);
	// 	// }
	// 	// const expectedType = isString(option.type) ? option.type : "string";
	// 	// return typeof value === expectedType;
	// }
	// return false
}

func convertJsonOptionOfListType(
	option CommandLineOptionOfListType,
	values []any, //readonly
	basePath string,
	errors []ast.Diagnostic,
	propertyAssignment ast.PropertyAssignment,
	valueExpression ast.ArrayLiteralExpression,
	sourceFile *ast.SourceFile,
) []any { //todo
	//return filter(map(values, (v, index) => convertJsonOption(option.element, v, basePath, errors, propertyAssignment, valueExpression?.elements[index], sourceFile)), v => option.listPreserveFalsyValues ? true : !!v);
	return []any{}
}

func convertJsonOption(
	opt commandLineOption,
	value any,
	basePath string,
	errors []ast.Diagnostic,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Expression,
	sourceFile *ast.SourceFile,
) compilerOptionsValue {
	if opt.commandLineOptionOfTypes.CommandLineOptionBase.isCommandLineOnly != nil {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, propertyAssignment?.name, Diagnostics.Option_0_can_only_be_specified_on_command_line, opt.name));
		return compilerOptionsValue{}
	}
	// if isCompilerOptionsValue(opt, value) != (compilerOptionsValue{}) {
	// 	const optType = opt.commandLineOptionBaseType.optionType
	// 	if (optType == "list") && isArray(value) { //todo something like isSlice??
	// 		return convertJsonOptionOfListType(opt, value, basePath, errors, propertyAssignment, valueExpression, sourceFile) //as ArrayLiteralExpression | undefined
	// 	}
	// 	// else if (optType === "listOrElement") {
	// 	//     return isArray(value) ?
	// 	//         convertJsonOptionOfListType(opt, value, basePath, errors, propertyAssignment, valueExpression as ArrayLiteralExpression | undefined, sourceFile) :
	// 	//         convertJsonOption(opt.element, value, basePath, errors, propertyAssignment, valueExpression, sourceFile);
	// 	// }
	// 	// else if (!isString(opt.type)) {
	// 	//     return convertJsonOptionOfCustomType(opt as CommandLineOptionOfCustomType, value as string, errors, valueExpression, sourceFile);
	// 	// }
	// 	const validatedValue = validateJsonOptionValue(opt, value, errors, valueExpression, sourceFile)
	// 	//return isNullOrUndefined(validatedValue) ? validatedValue : normalizeNonListOptionValue(opt, basePath, validatedValue);
	// }
	// else {
	//     errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, Diagnostics.Compiler_option_0_requires_a_value_of_type_1, opt.name, getCompilerOptionValueTypeString(opt)));
	// }
	return compilerOptionsValue{}
}

func getDefaultCompilerOptions(configFileName *string) core.CompilerOptions {
	var options core.CompilerOptions
	if configFileName != nil && getBaseFileName(*configFileName, nil, nil) == "jsconfig.json" {
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

func c(configFileName *string) compiler.TypeAcquisition {
	enable := configFileName != nil && getBaseFileName(*configFileName, nil, nil) == "jsconfig.json"
	return compiler.TypeAcquisition{Enable: &enable, Include: &[]string{}, Exclude: &[]string{}}
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
	if tspath.IsRootedDiskPath(extendedConfig) || tspath.StartsWith(extendedConfig, "./", nil) || compiler.StartsWith(extendedConfig, "../", nil) {
		extendedConfigPath := tspath.GetNormalizedAbsolutePath(extendedConfig, basePath)
		if !host.fileExists(extendedConfigPath) && !compiler.EndsWith(extendedConfigPath, Extension.Json) { //need to define Extension.Json
			extendedConfigPath = `${extendedConfigPath}.json`
			if !host.fileExists(extendedConfigPath) {
				//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, Diagnostics.File_0_not_found, extendedConfig));
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
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, Diagnostics.Compiler_option_0_cannot_be_given_an_empty_string, "extends"));
	} else {
		//errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, Diagnostics.File_0_not_found, extendedConfig));
	}
	return ""
}

func getExtendsConfigPathOrArray(
	value compilerOptionsValue,
	host ParseConfigHost,
	basePath string,
	configFileName string,
	errors []ast.Diagnostic,
	propertyAssignment ast.PropertyAssignment,
	valueExpression ast.Expression,
	sourceFile *ast.SourceFile,
) {
	var extendedConfigPath []string
	var newBase string
	if configFileName != "" {
		newBase = directoryOfCombinedPath(configFileName, basePath)
	} else {
		newBase = basePath
	}
	if compiler.IsString(value) {
		extendedConfigPath = []string{getExtendsConfigPath(
			value.stringValue,
			host,
			newBase,
			errors,
			valueExpression,
			sourceFile,
		)}
	} else if compiler.IsSlice(value) {
		extendedConfigPath = []string{}
		// for index := 0; index < len(value); index++ {
		// 	fileName := value[index]
		// 	if compiler.IsString(fileName) {
		// 		extendedConfigPath = append(
		// 			extendedConfigPath,
		// 			getExtendsConfigPath(
		// 				fileName,
		// 				host,
		// 				newBase,
		// 				errors,
		// 				valueExpression, //(valueExpression as ArrayLiteralExpression | undefined)?.elements[index],
		// 				sourceFile,
		// 			),
		// 		)
		// 	} else {
		// 		convertJsonOption(extendsOptionDeclaration.element, value, basePath, errors, propertyAssignment, valueExpression.elements[index], sourceFile)
		// 	}
		// }
	} else {
		//convertJsonOption(extendsOptionDeclaration, value, basePath, errors, propertyAssignment, valueExpression, sourceFile)
	}
	//return extendedConfigPath
}

func parseOwnConfigOfJson(
	json any,
	host ParseConfigHost,
	basePath string,
	configFileName string,
	errors []ast.Diagnostic,
) *ParsedTsconfig {
	// if (hasProperty(json, "excludes")) {
	//     errors.push(createCompilerDiagnostic(Diagnostics.Unknown_option_excludes_Did_you_mean_exclude));
	// }
	options := convertCompilerOptionsFromJsonWorker(json, basePath, errors, configFileName)
	typeAcquisition := convertTypeAcquisitionFromJsonWorker(json.typeAcquisition, basePath, errors, configFileName)
	watchOptions := convertWatchOptionsFromJsonWorker(json.watchOptions, basePath, errors)
	json.compileOnSave = convertCompileOnSaveOptionFromJson(json, basePath, errors)
	var extendedConfigPath string
	if json.extends != nil || json.extends == "" {
		extendedConfigPath = getExtendsConfigPathOrArray(json.extends, host, basePath, configFileName, errors)
	}
	var parsedConfig *ParsedTsconfig
	parsedConfig.raw = json
	parsedConfig.options = options
	parsedConfig.watchOptions = watchOptions
	parsedConfig.typeAcquisition = typeAcquisition
	parsedConfig.extendedConfigPath = extendedConfigPath
	return parsedConfig
}

var commandLineTypeAcquisitionMapCache map[string]commandLineOption

func getCommandLineTypeAcquisitionMap() {
	// if commandLineTypeAcquisitionMapCache != nil {
	// 	return commandLineTypeAcquisitionMap
	// }
	// commandLineTypeAcquisitionMapCache = commandLineOptionsToMap(typeAcquisitionDeclarations) //todo
	// return commandLineTypeAcquisitionMapCache
}

var commandLineCompilerOptionsMapCache map[string]commandLineOption

func getCommandLineCompilerOptionsMap() {
	// if commandLineCompilerOptionsMapCache != nil {
	// 	return commandLineCompilerOptionsMapCache
	// }
	// commandLineCompilerOptionsMapCache = commandLineOptionsToMap(compiler.optionDeclarations) //todo
	// return commandLineCompilerOptionsMapCache
}

var commandLineWatchOptionsMapCache map[string]commandLineOption

func getCommandLineWatchOptionsMap() {
	// if commandLineWatchOptionsMapCache != nil {
	// 	return commandLineWatchOptionsMapCache
	// }
	// commandLineWatchOptionsMapCache = commandLineOptionsToMap(optionsForWatch) //todo need to add watch related options
	// return commandLineWatchOptionsMapCache
}

func convertCompileOnSaveOptionFromJson(json any, basePath string, errors []ast.Diagnostic) bool {
	return false //todo
}
func convertWatchOptionsFromJsonWorker(jsonOptions any, basePath string, errors []ast.Diagnostic) compiler.WatchOptions {
	//return convertOptionsFromJson(getCommandLineWatchOptionsMap(), jsonOptions, basePath /*defaultOptions*/, undefined, watchOptionsDidYouMeanDiagnostics, errors)
}

func getDefaultTypeAcquisition(configFileName *string) *compiler.TypeAcquisition {
	var options compiler.TypeAcquisition
	//options.enable = !!configFileName && getBaseFileName(configFileName) == "jsconfig.json"
	options.include = []string{}
	options.exclude = []string{}
	return options
}

func convertTypeAcquisitionFromJsonWorker(jsonOptions any, basePath string, errors []ast.Diagnostic, configFileName *string) compiler.TypeAcquisition {
	// const options = getDefaultTypeAcquisition(configFileName)
	// convertOptionsFromJson(getCommandLineTypeAcquisitionMap(), jsonOptions, basePath, options, typeAcquisitionDidYouMeanDiagnostics, errors)
	// return options
}

func convertCompilerOptionsFromJsonWorker(jsonOptions any, basePath string, errors []ast.Diagnostic, configFileName *string) core.CompilerOptions {
	// options := getDefaultCompilerOptions(configFileName)
	// convertOptionsFromJson(getCommandLineCompilerOptionsMap(), jsonOptions, basePath, options, compilerOptionsDidYouMeanDiagnostics, errors)
	// if configFileName {
	// 	options.configFilePath = tspath.NormalizeSlashes(configFileName)
	// }
	// return options
}

type defaultOptions struct {
	core.CompilerOptions
	compiler.TypeAcquisition
	compiler.WatchOptions
}

func convertOptionsFromJson(optionsNameMap map[string]commandLineOption, jsonOptions any, basePath string, defaultOptions defaultOptions, diagnostics DidYouMeanOptionsDiagnostics, errors []ast.Diagnostic) {
	//todo
}

// func getDefaultCompilerOptions(configFileName *string) {
// 	var options compiler.CompilerOptions
// 	if configFileName != nil && compiler.GetBaseFileName(*configFileName) == "jsconfig.json" {
// 		options.allowJs = true
// 		maxNodeModuleJsDepth = 2
// 		allowSyntheticDefaultImports = true
// 		skipLibCheck = true
// 		noEmit = true
// 	}
// 	return options
// }

/**
 * Convert the json syntax tree into the json value
 */
func ConvertToObject(sourceFile *ast.SourceFile, errors []ast.Diagnostic) any {
	return convertToJson(sourceFile, sourceFile.Statements, errors /*returnValue*/, true /*jsonConversionNotifier*/, nil)
	//brb
}

/**
 * Convert the json syntax tree into the json value and report errors
 * This returns the json value (apart from checking errors) only if returnValue provided is true.
 * Otherwise it just checks the errors and returns undefined
 *
 * @internal
 */

// todo - all json work needs to be done
func convertToJson(
	sourceFile *ast.SourceFile,
	rootExpression ast.Expression,
	errors []ast.Diagnostic,
	returnValue bool,
	jsonConversionNotifier JsonConversionNotifier,
) any {
	if rootExpression == (ast.Expression{}) {
		return nil
	}
	//return convertPropertyValueToJson(rootExpression, jsonConversionNotifier.rootOptions) //todo
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
	compilerOptionsValue
	*ast.SourceFile
}
type optionsBase struct {
	options map[string]optionsBaseValue
}

func handleOptionConfigDirTemplateSubstitution(
	options optionsBase,
	optionDeclarations []commandLineOption,
	basePath string,
) {
	if len(options.options) == 0 { // if options == (optionsBase{}) {
		return options
	}
	var result optionsBase
	for _, option := range optionDeclarations {
		if options[option.name] != nil {
			const value = options[option.name]
			// switch (option.type) { //todo need to fix option.type prob 11/13/24
			//     case "string":
			//         Debug.assert(option.isFilePath);
			//         if (startsWithConfigDirTemplate(value)) {
			//             setOptionValue(option, getSubstitutedPathWithConfigDirTemplate(value, basePath));
			//         }
			//         break;
			//     case "list":
			//         Debug.assert(option.element.isFilePath);
			//         const listResult = getSubstitutedStringArrayWithConfigDirTemplate(value as string[], basePath);
			//         if (listResult) setOptionValue(option, listResult);
			//         break;
			//     case "object":
			//         Debug.assert(option.name === "paths");
			//         const objectResult = getSubstitutedMapLikeOfStringArrayWithConfigDirTemplate(value as MapLike<string[]>, basePath);
			//         if (objectResult) setOptionValue(option, objectResult);
			//         break;
			//     default:
			//         Debug.fail("option type not supported");
			// }
		}
	}
	return result || options

	// func setOptionValue(option: CommandLineOption, value: CompilerOptionsValue) {
	//     (result ??= assign({}, options))[option.name] = value;
	// }
}

func directoryOfCombinedPath(fileName string, basePath string) string {
	// Use the `getNormalizedAbsolutePath` function to avoid canonicalizing the path, as it must remain noncanonical
	// until consistent casing errors are reported
	return tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(fileName, basePath))
}

var defaultIncludeSpec = "**/*"

func validateSpecs(specs []string, errors []ast.Diagnostic, disallowTrailingRecursion bool, jsonSourceFile *ast.SourceFile, specKey string) []string {
	// return specs.filter(spec => {
	//     if (!isString(spec)) return false;
	//     const diag = specToDiagnostic(spec, disallowTrailingRecursion);
	//     if (diag !== undefined) {
	//         errors.push(createDiagnostic(...diag));
	//     }
	//     return diag === undefined;
	// });

	// function createDiagnostic(message: DiagnosticMessage, spec: string): Diagnostic {
	//     const element = getTsConfigPropArrayElementValue(jsonSourceFile, specKey, spec);
	//     return createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(jsonSourceFile, element, message, spec);
	// }
}

func getSubstitutedStringArrayWithConfigDirTemplate(list []string, basePath string) {
	if list == nil {
		return list
	}
	var result []string
	// list.forEach((element, index) => {
	//     if (!startsWithConfigDirTemplate(element)) return;
	//     (result ??= list.slice())[index] = getSubstitutedPathWithConfigDirTemplate(element, basePath);
	// });
	return result
}

func setConfigFileInOptions(options core.CompilerOptions, configFile *ast.SourceFile) {
	// if (configFile) {
	//     Object.defineProperty(options, "configFile", { enumerable: false, writable: false, value: configFile });
	// }
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
	basePath string,
	options core.CompilerOptions,
	host ParseConfigHost,
	extraFileExtensions []FileExtensionInfo,
) []string {
	basePath = tspath.NormalizePath(basePath)

	//const keyMapper = createGetCanonicalFileName(host.useCaseSensitiveFileNames);// core.ts

	// Literal file names (provided via the "files" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map later when when including
	// wildcard paths.
	literalFileMap := make(map[string]string)

	// Wildcard paths (provided via the "includes" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map to store paths matched
	// via wildcard, and to handle extension priority.
	//wildcardFileMap := make(map[string]string)

	// Wildcard paths of json files (provided via the "includes" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map to store paths matched
	// via wildcard of *.json kind
	//wildCardJsonFileMap := make(map[string]string)
	validatedFilesSpec := configFileSpecs.validatedFilesSpec
	validatedIncludeSpecs := configFileSpecs.validatedIncludeSpecs
	validatedExcludeSpecs := configFileSpecs.validatedExcludeSpecs

	// Rather than re-query this for each file and filespec, we query the supported extensions
	// once and store it on the expansion context.
	//const supportedExtensions = getSupportedExtensions(options, extraFileExtensions); // in utilities.ts
	//const supportedExtensionsWithJsonIfResolveJsonModule = getSupportedExtensionsWithJsonIfResolveJsonModule(options, supportedExtensions)

	// Literal files are always included verbatim. An "include" or "exclude" specification cannot
	// remove a literal file.
	if validatedFilesSpec != nil {
		for _, fileName := range validatedFilesSpec {
			file := tspath.GetNormalizedAbsolutePath(fileName, basePath)
			literalFileMap.set(keyMapper(file), file)
		}
	}

	//var jsonOnlyIncludeRegexes []regexp.Regexp
	// if validatedIncludeSpecs && len(validatedIncludeSpecs) > 0 {
	// 	for _, file := range host.readDirectory(basePath, flatten(supportedExtensionsWithJsonIfResolveJsonModule), validatedExcludeSpecs, validatedIncludeSpecs /*depth*/, undefined) {
	// 		if compiler.FileExtensionIs(file, Extension.Json) {
	// 			// Valid only if *.json specified
	// 			// if (!jsonOnlyIncludeRegexes) {
	// 			//     const includes = validatedIncludeSpecs.filter(s => endsWith(s, Extension.Json));
	// 			//     const includeFilePatterns = map(getRegularExpressionsForWildcards(includes, basePath, "files"), pattern => `^${pattern}$`);
	// 			//     jsonOnlyIncludeRegexes = includeFilePatterns ? includeFilePatterns.map(pattern => getRegexFromPattern(pattern, host.useCaseSensitiveFileNames)) : emptyArray;
	// 			// }
	// 			// const includeIndex = findIndex(jsonOnlyIncludeRegexes, re => re.test(file));
	// 			// if (includeIndex !== -1) {
	// 			//     const key = keyMapper(file);
	// 			//     if (!literalFileMap.has(key) && !wildCardJsonFileMap.has(key)) {
	// 			//         wildCardJsonFileMap.set(key, file);
	// 			//     }
	// 			// }
	// 			// continue;
	// 		}
	// 		// If we have already included a literal or wildcard path with a
	// 		// higher priority extension, we should skip this file.
	// 		//
	// 		// This handles cases where we may encounter both <file>.ts and
	// 		// <file>.d.ts (or <file>.js if "allowJs" is enabled) in the same
	// 		// directory when they are compilation outputs.
	// 		// if (hasFileWithHigherPriorityExtension(file, literalFileMap, wildcardFileMap, supportedExtensions, keyMapper)) {
	// 		//     continue;
	// 		// }

	// 		// We may have included a wildcard path with a lower priority
	// 		// extension due to the user-defined order of entries in the
	// 		// "include" array. If there is a lower priority extension in the
	// 		// same directory, we should remove it.
	// 		// removeWildcardFilesWithLowerPriorityExtension(file, wildcardFileMap, supportedExtensions, keyMapper);

	// 		// const key = keyMapper(file);
	// 		// if (!literalFileMap.has(key) && !wildcardFileMap.has(key)) {
	// 		//     wildcardFileMap.set(key, file);
	// 		// }
	// 	}
	// }

	// const literalFiles = arrayFrom(literalFileMap.values());
	// const wildcardFiles = arrayFrom(wildcardFileMap.values());

	// return literalFiles.concat(wildcardFiles, arrayFrom(wildCardJsonFileMap.values()));
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
	json any,
	sourceFile *ast.SourceFile,
	host ParseConfigHost,
	basePath string,
	existingOptions core.CompilerOptions, //should default to an empty object
	existingWatchOptions compiler.WatchOptions,
	configFileName *string,
	resolutionStack []tspath.Path,
	extraFileExtensions []FileExtensionInfo,
	extendedConfigCache *map[string]ExtendedConfigCacheEntry,
) *compiler.ParsedCommandLine {
	//Debug.assert((json === undefined && sourceFile !== undefined) || (json !== undefined && sourceFile === undefined));
	var errors []ast.Diagnostic
	parsedConfig := parseConfig(json, sourceFile, host, basePath, configFileName, resolutionStack, extendedConfigCache)
	var raw = parsedConfig.raw
	// const options = handleOptionConfigDirTemplateSubstitution(
	// 	extend(existingOptions, parsedConfig.options), //function in core.ts
	// 	configDirTemplateSubstitutionOptions,
	// 	basePath,
	// )
	// const watchOptions = handleWatchOptionsConfigDirTemplateSubstitution(
	//     existingWatchOptions && parsedConfig.watchOptions ?
	//         extend(existingWatchOptions, parsedConfig.watchOptions) :
	//         parsedConfig.watchOptions || existingWatchOptions,
	//     basePath,
	// );
	options.configFilePath = configFileName && tspath.NormalizeSlashes(configFileName)
	var basePathForFileNames string
	if configFileName != nil {
		basePathForFileNames = tspath.NormalizePath(directoryOfCombinedPath(configFileName, basePath))
	} else {
		basePathForFileNames = tspath.NormalizePath(basePath)
	}

	type validateElement func(value any) bool
	type propOfRaw[T any] struct {
		array    *[]T
		notArray *string
		noProp   *string
	}
	getPropFromRaw := func(prop string, validate validateElement, elementTypeName string) propOfRaw {
		// if (hasProperty(raw, prop) && !isNullOrUndefined(raw[prop])) { hasProperty is a function in core.ts
		//     if (isArray(raw[prop])) {
		//         const result = raw[prop] as T[];
		//         if (!sourceFile && !every(result, validateElement)) {
		//             errors.push(createCompilerDiagnostic(Diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop, elementTypeName));
		//         }
		//         return result;
		//     }
		//     else {
		//         createCompilerDiagnosticOnlyIfJson(Diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop, "Array");
		//         return "not-array";
		//     }
		// }
		return "no-prop"
	}

	toPropValue := func(specResult propOfRaw) {
		if specResult.array != nil {
			return specResult
		}
	}

	getSpecsFromRaw := func(prop string) propOfRaw { //prop: "files" | "include" | "exclude"
		return getPropFromRaw(prop, isString, "string")
	}

	getFileNames := func(basePath string) []string {
		var fileNames = getFileNamesFromConfigSpecs(configFileSpecs, basePath, options, host, extraFileExtensions)
		// if shouldReportNoInputFiles(fileNames, canJsonReportNoInputFiles(raw), resolutionStack) {
		// 	errors.push(getErrorForNoInputFiles(configFileSpecs, configFileName))
		// }
		return fileNames
	}

	getProjectReferences := func(basePath string) []compiler.ProjectReference {
		var projectReferences = []compiler.ProjectReference{}
		const referencesOfRaw = getPropFromRaw("references", validateElement("onject"), "object")
		if compiler.IsSlice(referencesOfRaw) {
			for _, ref := range referencesOfRaw {
				if ref.path != "string" { //typeof ref.path !== "string"
					//createCompilerDiagnosticOnlyIfJson(Diagnostics.Compiler_option_0_requires_a_value_of_type_1, "reference.path", "string");
				} else {
					projectReferences = append(projectReferences, compiler.ProjectReference{
						path:         tspath.getNormalizedAbsolutePath(ref.path, basePath),
						originalPath: ref.path,
						prepend:      ref.prepend,
						circular:     ref.circular,
					})
				}
			}
		}
		return projectReferences
	}

	createCompilerDiagnosticOnlyIfJson := func(message []diagnostics.Message, args compiler.DiagnosticAndArguments) { //todo full
		// if (!sourceFile) {
		//     errors.push(createCompilerDiagnostic(message, ...args));
		// }
	}

	getConfigFileSpecs := func() configFileSpecs {
		referencesOfRaw := getPropFromRaw("references", validateElement(), "object") // come back to validateElement
		filesSpecs := toPropValue(getSpecsFromRaw("files"))
		if filesSpecs {
			hasZeroOrNoReferences := referencesOfRaw == "no-prop" || compiler.IsSlice(referencesOfRaw) && referencesOfRaw.length == 0
			hasExtends := hasProperty(raw, "extends") //hasProperty is a function in core.ts
			if filesSpecs.length == 0 && hasZeroOrNoReferences && !hasExtends {
				if sourceFile {
					fileName := configFileName || "tsconfig.json"
					//diagnosticMessage := Diagnostics.The_files_list_in_config_file_0_is_empty;
					//nodeValue := forEachTsConfigPropArray(sourceFile, "files", property => property.initializer);
					//const error = createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, nodeValue, diagnosticMessage, fileName);
					//errors.push(error);
				} else {
					//createCompilerDiagnosticOnlyIfJson(Diagnostics.The_files_list_in_config_file_0_is_empty, configFileName || "tsconfig.json");
				}
			}
		}

		includeSpecs := toPropValue(getSpecsFromRaw("include"))

		excludeOfRaw := getSpecsFromRaw("exclude")
		isDefaultIncludeSpec := false
		excludeSpecs := toPropValue(excludeOfRaw)
		if excludeOfRaw == "no-prop" {
			outDir := options.outDir
			declarationDir := options.declarationDir

			if outDir || declarationDir {
				//excludeSpecs = filter([outDir, declarationDir], d => !!d) as string[];//filter is function in core.ts
			}
		}

		if filesSpecs == nil && includeSpecs == nil {
			includeSpecs = []string{defaultIncludeSpec}
			isDefaultIncludeSpec = true
		}
		var validatedIncludeSpecsBeforeSubstitution []string
		var alidatedExcludeSpecsBeforeSubstitution []string
		var validatedIncludeSpecs []string
		var validatedExcludeSpecs []string

		// The exclude spec list is converted into a regular expression, which allows us to quickly
		// test whether a file or directory should be excluded before recursively traversing the
		// file system.

		if includeSpecs {
			validatedIncludeSpecsBeforeSubstitution = validateSpecs(includeSpecs, errors /*disallowTrailingRecursion*/, true, sourceFile, "include")
			validatedIncludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
				validatedIncludeSpecsBeforeSubstitution,
				basePathForFileNames,
			) || validatedIncludeSpecsBeforeSubstitution
		}

		if excludeSpecs {
			validatedExcludeSpecsBeforeSubstitution = validateSpecs(excludeSpecs, errors /*disallowTrailingRecursion*/, false, sourceFile, "exclude")
			validatedExcludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
				validatedExcludeSpecsBeforeSubstitution,
				basePathForFileNames,
			) || validatedExcludeSpecsBeforeSubstitution
		}

		//validatedFilesSpecBeforeSubstitution := filter(filesSpecs, isString) //filter is a function in core.ts
		validatedFilesSpec := getSubstitutedStringArrayWithConfigDirTemplate(
			validatedFilesSpecBeforeSubstitution,
			basePathForFileNames,
		) || validatedFilesSpecBeforeSubstitution

		return configFileSpecs{
			filesSpecs,
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
	if sourceFile {
		sourceFile.configFileSpecs = configFileSpecs
	}
	setConfigFileInOptions(options, sourceFile)
	// result := compiler.ParsedCommandLine {
	//     options: options,
	//     watchOptions: watchOptions,
	//     fileNames: getFileNames(basePathForFileNames),
	//     projectReferences: getProjectReferences(basePathForFileNames),
	//     typeAcquisition: parsedConfig.typeAcquisition || getDefaultTypeAcquisition(),
	//     raw: raw,
	//     errors: errors,
	//     // Wildcard directories (provided as part of a wildcard path) are stored in a
	//     // file map that marks whether it was a regular wildcard match (with a `*` or `?` token),
	//     // or a recursive directory. This information is used by filesystem watchers to monitor for
	//     // new entries in these paths.
	//     wildcardDirectories: getWildcardDirectories(configFileSpecs, basePathForFileNames, host.useCaseSensitiveFileNames),
	//     compileOnSave: !!raw.compileOnSave,
	// }

	return nil

}
