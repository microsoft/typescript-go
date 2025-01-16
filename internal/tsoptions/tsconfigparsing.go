package tsoptions

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/dlclark/regexp2"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/compiler/module"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type extendsResult struct {
	options *core.CompilerOptions
	// watchOptions        compiler.WatchOptions
	watchOptionsCopied  bool
	include             []string
	exclude             []string
	files               []string
	compileOnSave       bool
	extendedSourceFiles map[string]struct{}
}

var (
	tsconfigRootOptions       *CommandLineOption
	getTsconfigRootOptionsMap = sync.OnceValue(func() CommandLineOption {
		if tsconfigRootOptions == nil {
			tsconfigRootOptions = &CommandLineOption{
				Name: "undefined", // should never be needed since this is root
				Kind: CommandLineOptionTypeObject,
				ElementOptions: commandLineOptionsToMap([]*CommandLineOption{
					compilerOptionsDeclaration,
					extendsOptionDeclaration,
					{
						Name: "references",
						Kind: CommandLineOptionTypeList, // should be a list of projectReference
						// Category: diagnostics.Projects,
					},
					{
						Name: "files",
						Kind: CommandLineOptionTypeList,
						// Category: diagnostics.File_Management,
					},
					{
						Name: "include",
						Kind: CommandLineOptionTypeList,
						// Category: diagnostics.File_Management,
						// DefaultValueDescription: diagnostics.if_files_is_specified_otherwise_Asterisk_Asterisk_Slash_Asterisk,
					},
					{
						Name: "exclude",
						Kind: CommandLineOptionTypeList,
						// Category: diagnostics.File_Management,
						// DefaultValueDescription: diagnostics.Node_modules_bower_components_jspm_packages_plus_the_value_of_outDir_if_one_is_specified,
					},
					compileOnSaveCommandLineOption,
				}),
			}
		}
		return *tsconfigRootOptions
	})
)

type configFileSpecs struct {
	filesSpecs any
	// Present to report errors (user specified specs), validatedIncludeSpecs are used for file name matching
	includeSpecs any
	// Present to report errors (user specified specs), validatedExcludeSpecs are used for file name matching
	excludeSpecs                            any
	validatedFilesSpec                      []string
	validatedIncludeSpecs                   []string
	validatedExcludeSpecs                   []string
	validatedFilesSpecBeforeSubstitution    []string
	validatedIncludeSpecsBeforeSubstitution []string
	validatedExcludeSpecsBeforeSubstitution []string
	isDefaultIncludeSpec                    bool
}
type fileExtensionInfo struct {
	extension      string
	isMixedContent bool
	scriptKind     core.ScriptKind
}
type extendedConfigCacheEntry struct {
	extendedResult *tsConfigSourceFile
	extendedConfig *parsedTsconfig
}
type parsedTsconfig struct {
	raw     any
	options *core.CompilerOptions
	// watchOptions    *compiler.WatchOptions
	// typeAcquisition *compiler.TypeAcquisition
	// Note that the case of the config path has not yet been normalized, as no files have been imported into the project yet
	extendedConfigPath any
}

func parseOwnConfigOfJsonSourceFile(
	sourceFile *tsConfigSourceFile,
	host ParseConfigHost,
	basePath string,
	configFileName string,
) (*parsedTsconfig, []*ast.Diagnostic) {
	options := getDefaultCompilerOptions(configFileName)
	// var typeAcquisition *compiler.TypeAcquisition
	// var watchOptions *compiler.WatchOptions
	var extendedConfigPath any
	var rootCompilerOptions []*ast.PropertyName
	var errors []*ast.Diagnostic
	rootOptions := getTsconfigRootOptionsMap()
	onPropertySet := func(
		keyText string,
		value any,
		propertyAssignment *ast.PropertyAssignment,
		parentOption CommandLineOption, // TsConfigOnlyOption,
		option *CommandLineOption,
	) (any, []*ast.Diagnostic) {
		// Ensure value is verified except for extends which is handled in its own way for error reporting
		var propertySetErrors []*ast.Diagnostic
		if option != nil && option != extendsOptionDeclaration {
			value, propertySetErrors = convertJsonOption(option, value, basePath, propertyAssignment, propertyAssignment.Initializer, sourceFile)
		}
		if parentOption.Name != "undefined" && value != nil {
			if option != nil && option.Name != "" {
				commandLineOptionEnumMapVal := option.EnumMap()
				if commandLineOptionEnumMapVal != nil {
					val, ok := commandLineOptionEnumMapVal.Get(strings.ToLower(value.(string)))
					if ok {
						propertySetErrors = append(propertySetErrors, parseCompilerOptions(option.Name, val, options)...)
					}
				} else {
					propertySetErrors = append(propertySetErrors, parseCompilerOptions(option.Name, value, options)...)
				}
			} else if keyText != "" {
				if parentOption.ElementOptions != nil {
					propertySetErrors = append(propertySetErrors, ast.NewCompilerDiagnostic(diagnostics.Option_build_must_be_the_first_command_line_argument, keyText))
				} else {
					// errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Unknown_compiler_option_0_Did_you_mean_1, keyText, core.FindKey(parentOption.ElementOptions, keyText)))
				}
			}
		} else if parentOption.Name == rootOptions.Name {
			if option.Name == extendsOptionDeclaration.Name {
				configPath, err := getExtendsConfigPathOrArray(value, host, basePath, configFileName, propertyAssignment, propertyAssignment.Initializer, sourceFile)
				extendedConfigPath = configPath
				propertySetErrors = append(propertySetErrors, err...)
			} else if option == nil {
				if keyText == "excludes" {
					propertySetErrors = append(propertySetErrors, ast.NewCompilerDiagnostic(diagnostics.Unknown_option_excludes_Did_you_mean_exclude))
				}
				if core.Find(optionsDeclarations, func(option *CommandLineOption) bool {
					return option.Name == keyText
				}) != nil {
					rootCompilerOptions = append(rootCompilerOptions, propertyAssignment.Name())
				}
			}
		}
		return value, propertySetErrors
	}

	json, err := convertConfigFileToObject(
		sourceFile.sourceFile,
		&jsonConversionNotifier{
			rootOptions,
			onPropertySet,
		},
	)
	errors = append(errors, err...)
	return &parsedTsconfig{
		raw:     json,
		options: options,
		// watchOptions:    watchOptions,
		// typeAcquisition: typeAcquisition,
		extendedConfigPath: extendedConfigPath,
	}, errors
}

type tsConfigSourceFile struct {
	extendedSourceFiles []string
	configFileSpecs     *configFileSpecs
	sourceFile          *ast.SourceFile
}
type jsonConversionNotifier struct {
	rootOptions   CommandLineOption
	onPropertySet func(keyText string, value any, propertyAssignment *ast.PropertyAssignment, parentOption CommandLineOption, option *CommandLineOption) (any, []*ast.Diagnostic)
}

func convertConfigFileToObject(
	sourceFile *ast.SourceFile,
	jsonConversionNotifier *jsonConversionNotifier,
) (any, []*ast.Diagnostic) {
	var rootExpression *ast.Expression
	if len(sourceFile.Statements.Nodes) > 0 {
		rootExpression = sourceFile.Statements.Nodes[0].AsExpressionStatement().Expression
	}
	if rootExpression != nil && rootExpression.Kind != ast.KindObjectLiteralExpression {
		baseFileName := "tsconfig.json"
		if tspath.GetBaseFileName(sourceFile.FileName()) == "jsconfig.json" {
			baseFileName = "jsconfig.json"
		}
		errors := []*ast.Diagnostic{ast.NewCompilerDiagnostic(diagnostics.The_root_value_of_a_0_file_must_be_an_object, baseFileName)}
		// Last-ditch error recovery. Somewhat useful because the JSON parser will recover from some parse errors by
		// synthesizing a top-level array literal expression. There's a reasonable chance the first element of that
		// array is a well-formed configuration object, made into an array element by stray characters.
		if ast.IsArrayLiteralExpression(rootExpression) {
			firstObject := core.Find(rootExpression.AsArrayLiteralExpression().Elements.Nodes, ast.IsObjectLiteralExpression)
			if firstObject != nil {
				return convertToJson(sourceFile, firstObject /*returnValue*/, true, jsonConversionNotifier)
			}
		}
		return make(map[string]any), errors
	}
	return convertToJson(sourceFile, rootExpression, true, jsonConversionNotifier)
}

func isCompilerOptionsValue(option *CommandLineOption, value any) bool {
	if option.Name != "" || option.Kind != "" {
		if value == nil {
			return !option.DisallowNullOrUndefined() // All options are undefinable/nullable
		}
		if option.Kind == "list" {
			return reflect.TypeOf(value).Kind() == reflect.Slice
		}
		if option.Kind == "listOrElement" {
			if reflect.TypeOf(value).Kind() == reflect.Slice {
				return true
			} else {
				return isCompilerOptionsValue(option.Elements(), value)
			}
		}
		if option.Kind == "string" {
			return reflect.TypeOf(value).Kind() == reflect.String
		}
		if option.Kind == "boolean" {
			return reflect.TypeOf(value).Kind() == reflect.Bool
		}
		if option.Kind == "number" {
			return reflect.TypeOf(value).Kind() == reflect.Int
		}
		if option.Kind == "object" {
			return reflect.TypeOf(value).Kind() == reflect.Map
		}
		if option.Kind == "enum" {
			return true
		}
	}
	return false
}

func validateJsonOptionValue(
	opt *CommandLineOption,
	val any,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) any {
	if val == nil || val == "" {
		return nil
	}
	d := (opt.extraValidation)
	if d == nil {
		return val
	} else {
		// d = opt.extraValidation.val
	}
	// errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, d));
	return nil
}

func convertJsonOptionOfCustomType(
	opt *CommandLineOption,
	value string,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) any {
	if value == "" {
		return nil
	}
	key := strings.ToLower(value)
	typeMap := opt.EnumMap()
	if typeMap == nil {
		return nil
	}
	val, b := typeMap.Get(key)
	if (val != nil) && (val != "" || b) {
		return validateJsonOptionValue(opt, val.(string), valueExpression, sourceFile)
	}
	// else {
	//     errors.push(createDiagnosticForInvalidCustomType(opt, (message, ...args) => createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, message, ...args)));
	// }
	return nil
}

func convertJsonOptionOfListType(
	option *CommandLineOption,
	values any,
	basePath string,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Node,
	sourceFile *tsConfigSourceFile,
) ([]any, []*ast.Diagnostic) {
	var expression *ast.Node
	var errors []*ast.Diagnostic
	index := 0
	if _, ok := values.([]any); ok {
		mappedValue := core.Map(values.([]any), func(v any) any {
			if valueExpression != nil {
				expression = valueExpression.AsArrayLiteralExpression().Elements.Nodes[index]
			}
			result, err := convertJsonOption(option.Elements(), v, basePath, propertyAssignment, expression, sourceFile)
			index++
			errors = append(errors, err...)
			return result
		})
		filteredValues := core.Filter(mappedValue, func(v any) bool {
			if option.listPreserveFalsyValues {
				return true
			} else {
				return (v != nil && v != false && v != 0 && v != "")
			}
		})
		return filteredValues, errors
	}
	return nil, errors
}

const configDirTemplate = "${configDir}"

func startsWithConfigDirTemplate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return strings.HasPrefix(strings.ToLower(str), strings.ToLower(configDirTemplate))
}

func normalizeNonListOptionValue(option *CommandLineOption, basePath string, value any) any {
	if option.isFilePath {
		value = tspath.NormalizeSlashes(value.(string))
		if !startsWithConfigDirTemplate(value) {
			value = tspath.GetNormalizedAbsolutePath(value.(string), basePath)
		}
		if value == "" {
			value = "."
		}
	}
	return value
}

func convertJsonOption(
	opt *CommandLineOption,
	value any,
	basePath string,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) (any, []*ast.Diagnostic) {
	var errors []*ast.Diagnostic
	if opt.isCommandLineOnly {
		var nodeValue *ast.Node
		if propertyAssignment != nil {
			nodeValue = propertyAssignment.Name()
		}
		if sourceFile == nil && nodeValue == nil {
			errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Option_0_can_only_be_specified_on_command_line, opt.Name))
		} else {
			errors = append(errors, ast.NewDiagnostic(sourceFile.sourceFile, core.NewTextRange(scanner.SkipTrivia(sourceFile.sourceFile.Text, nodeValue.Loc.Pos()), nodeValue.End()), diagnostics.Option_0_can_only_be_specified_on_command_line, opt.Name))
		}
		return false, errors
	}
	if isCompilerOptionsValue(opt, value) {
		optType := opt.Kind
		if optType == "list" {
			list, err := convertJsonOptionOfListType(opt, value, basePath, propertyAssignment, valueExpression, sourceFile) // as ArrayLiteralExpression | undefined
			return list, append(errors, err...)
		} else if optType == "listOrElement" {
			if reflect.TypeOf(value).Kind() == reflect.Slice {
				listOrElement, err := convertJsonOptionOfListType(opt, value, basePath, propertyAssignment, valueExpression, sourceFile)
				errors = append(errors, err...)
				return listOrElement, errors
			} else {
				return convertJsonOption(opt.Elements(), value, basePath, propertyAssignment, valueExpression, sourceFile)
			}
		} else if !(reflect.TypeOf(optType).Kind() == reflect.String) {
			return convertJsonOptionOfCustomType(opt, value.(string), valueExpression, sourceFile), errors
		}
		validatedValue := validateJsonOptionValue(opt, value, valueExpression, sourceFile)
		if validatedValue == nil {
			return validatedValue, errors
		} else {
			return normalizeNonListOptionValue(opt, basePath, validatedValue), errors
		}
	} else {
		errors = append(errors, ast.NewDiagnostic(sourceFile.sourceFile, core.NewTextRange(scanner.SkipTrivia(sourceFile.sourceFile.Text, valueExpression.Loc.Pos()), valueExpression.End()), diagnostics.Compiler_option_0_requires_a_value_of_type_1, opt.Name, getCompilerOptionValueTypeString(opt)))
		return nil, errors
	}
}

func getExtendsConfigPathOrArray(
	value CompilerOptionsValue,
	host ParseConfigHost,
	basePath string,
	configFileName string,
	propertyAssignment *ast.PropertyAssignment,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) ([]string, []*ast.Diagnostic) {
	var extendedConfigPathArray []string
	var newBase string
	if configFileName != "" {
		newBase = directoryOfCombinedPath(configFileName, basePath)
	} else {
		newBase = basePath
	}

	var errors []*ast.Diagnostic
	if reflect.TypeOf(value).Kind() == reflect.String {
		val, err := getExtendsConfigPath(value.(string), host, newBase, valueExpression, sourceFile)
		if val != "" {
			extendedConfigPathArray = append(extendedConfigPathArray, val)
		}
		errors = append(errors, err...)
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		for index, fileName := range value.([]any) {
			if reflect.TypeOf(fileName).Kind() == reflect.String {
				var expression *ast.Expression = nil
				if valueExpression != nil {
					expression = valueExpression.AsArrayLiteralExpression().Elements.Nodes[index]
				}
				val, err := getExtendsConfigPath(fileName.(string), host, newBase, expression, sourceFile)
				if val != "" {
					extendedConfigPathArray = append(extendedConfigPathArray, val)
				}
				errors = append(errors, err...)
			} else {
				var err []*ast.Diagnostic
				_, err = convertJsonOption(extendsOptionDeclaration.Elements(), value, basePath, propertyAssignment, valueExpression.AsArrayLiteralExpression().Elements.Nodes[index], sourceFile) // check
				errors = append(errors, err...)
			}
		}
	} else {
		_, errors = convertJsonOption(extendsOptionDeclaration, value, basePath, propertyAssignment, valueExpression, sourceFile) // check
	}
	return extendedConfigPathArray, errors
}

func getExtendsConfigPath(
	extendedConfig string,
	host ParseConfigHost,
	basePath string,
	valueExpression *ast.Expression,
	sourceFile *tsConfigSourceFile,
) (string, []*ast.Diagnostic) {
	extendedConfig = tspath.NormalizeSlashes(extendedConfig)
	var errors []*ast.Diagnostic
	if tspath.IsRootedDiskPath(extendedConfig) || strings.HasPrefix(extendedConfig, "./") || strings.HasPrefix(extendedConfig, "../") {
		extendedConfigPath := tspath.GetNormalizedAbsolutePath(extendedConfig, basePath)
		if !host.FS().FileExists(extendedConfigPath) && !strings.HasSuffix(extendedConfigPath, tspath.ExtensionJson) {
			extendedConfigPath = extendedConfigPath + tspath.ExtensionJson
			if !host.FS().FileExists(extendedConfigPath) {
				// errors.push(createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.File_0_not_found, extendedConfig));
				errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.File_0_not_found, extendedConfig))
				return "", errors
			}
		}
		return extendedConfigPath, errors
	}
	// If the path isn't a rooted or relative path, resolve like a module
	resolverHost := &resolverHost{host}
	if resolved := module.ResolveConfig(extendedConfig, tspath.CombinePaths(basePath, "tsconfig.json"), resolverHost); resolved.IsResolved() {
		return resolved.ResolvedFileName, errors
	}
	if extendedConfig == "" {
		errors = append(errors, createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(nil, nil, diagnostics.Compiler_option_0_cannot_be_given_an_empty_string, "extends"))
	} else {
		errors = append(errors, createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(nil, nil, diagnostics.File_0_not_found, extendedConfig))
	}
	return "", errors
}

type tsConfigOptions struct {
	prop       map[string][]string
	references []core.ProjectReference
	notDefined string
}

func commandLineOptionsToMap(options []*CommandLineOption) map[string]*CommandLineOption {
	result := make(map[string]*CommandLineOption)
	for i := range options {
		result[(options[i]).Name] = options[i]
	}
	return result
}

var commandLineCompilerOptionsMapCache map[string]*CommandLineOption

func getCommandLineCompilerOptionsMap() map[string]*CommandLineOption {
	if commandLineCompilerOptionsMapCache != nil {
		return commandLineCompilerOptionsMapCache
	}
	commandLineCompilerOptionsMapCache = commandLineOptionsToMap(optionsDeclarations)
	return commandLineCompilerOptionsMapCache
}

func convertOptionsFromJson(optionsNameMap map[string]*CommandLineOption, jsonOptions any, basePath string, defaultOptions *core.CompilerOptions) (*core.CompilerOptions, []*ast.Diagnostic) {
	if jsonOptions == nil {
		return nil, nil
	}
	var errors []*ast.Diagnostic
	if _, ok := jsonOptions.(map[string]any); ok {
		for key, value := range jsonOptions.(map[string]any) {
			opt, ok := optionsNameMap[key]
			commandLineOptionEnumMapVal := opt.EnumMap()
			if commandLineOptionEnumMapVal != nil {
				val, ok := commandLineOptionEnumMapVal.Get(strings.ToLower(value.(string)))
				if ok {
					errors = parseCompilerOptions(key, val, defaultOptions)
				}
			} else if ok {
				convertJson, err := convertJsonOption(opt, value, basePath, nil, nil, nil)
				errors = append(errors, err...)
				compilerOptionsErr := parseCompilerOptions(key, convertJson, defaultOptions)
				errors = append(errors, compilerOptionsErr...)
			}
			// else {
			//     errors.push(createUnknownOptionError(id, diagnostics));
			// }
		}
	}
	return defaultOptions, errors
}

func convertArrayLiteralExpressionToJson(
	sourceFile *ast.SourceFile,
	elements []*ast.Expression,
	elementOption *CommandLineOption,
	returnValue bool,
) (any, []*ast.Diagnostic) {
	if !returnValue {
		for _, element := range elements {
			convertPropertyValueToJson(sourceFile, element, elementOption, returnValue, nil)
		}
		return nil, nil
	}
	// Filter out invalid values
	if len(elements) == 0 {
		return []string{}, nil
	}
	var errors []*ast.Diagnostic
	var value []any
	for _, element := range elements {
		convertedValue, err := convertPropertyValueToJson(sourceFile, element, elementOption, returnValue, nil)
		errors = append(errors, err...)
		value = append(value, convertedValue)
	}
	return value, errors
}

func directoryOfCombinedPath(fileName string, basePath string) string {
	// Use the `getNormalizedAbsolutePath` function to avoid canonicalizing the path, as it must remain noncanonical
	// until consistent casing errors are reported
	return tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(fileName, basePath))
}

// ParseConfigFileTextToJson parses the text of the tsconfig.json file
// fileName is the path to the config file
// jsonText is the text of the config file
func ParseConfigFileTextToJson(fileName string, basePath string, jsonText string) (any, []*ast.Diagnostic) {
	jsonSourceFile := parser.ParseJSONText(fileName, jsonText)
	config, errors := convertConfigFileToObject(jsonSourceFile /*jsonConversionNotifier*/, nil)
	if len(jsonSourceFile.Diagnostics()) > 0 {
		errors = []*ast.Diagnostic{jsonSourceFile.Diagnostics()[0]}
	}
	return config, errors
}

type ParseConfigHost interface {
	FS() vfs.FS
	GetCurrentDirectory() string
}

type resolverHost struct {
	ParseConfigHost
}

func (r *resolverHost) Trace(msg string) {}

func ParseJsonSourceFileConfigFileContent(sourceFile *tsConfigSourceFile, host ParseConfigHost, basePath string, existingOptions *core.CompilerOptions, configFileName string, resolutionStack []tspath.Path, extraFileExtensions []fileExtensionInfo, extendedConfigCache map[string]*extendedConfigCacheEntry) ParsedCommandLine {
	// tracing?.push(tracing.Phase.Parse, "parseJsonSourceFileConfigFileContent", { path: sourceFile.fileName });
	result := parseJsonConfigFileContentWorker( /*json*/ nil, sourceFile, host, basePath, existingOptions, configFileName, resolutionStack, extraFileExtensions, extendedConfigCache)
	// tracing?.pop();
	return result
}

func convertObjectLiteralExpressionToJson(
	sourceFile *ast.SourceFile,
	returnValue bool,
	node *ast.ObjectLiteralExpression,
	objectOption *CommandLineOption,
	jsonConversionNotifier *jsonConversionNotifier,
) (map[string]any, []*ast.Diagnostic) {
	var result map[string]any
	if returnValue {
		result = make(map[string]any)
	} else {
		result = nil
	}
	var errors []*ast.Diagnostic
	for _, element := range node.Properties.Nodes {
		if element.Kind != ast.KindPropertyAssignment {
			errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Property_assignment_expected))
			continue
		}

		if ast.IsQuestionToken(element) {
			errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Property_assignment_expected))
		}
		if element.Name() != nil && !isDoubleQuotedString(element.Name()) {
			errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.String_literal_with_double_quotes_expected))
		}

		var textOfKey any
		if ast.IsComputedNonLiteralName(element.Name()) {
			textOfKey = nil
		} else {
			textOfKey, _ = ast.TryGetTextOfPropertyName(element.Name())
		}
		keyText := textOfKey
		var option *CommandLineOption
		if keyText != nil {
			if objectOption != nil && objectOption.ElementOptions != nil {
				option = objectOption.ElementOptions[keyText.(string)]
			} else {
				option = &CommandLineOption{}
				// errors.push(createDiagnosticForNodeInSourceFile(sourceFile, element.name, diagnostics.Unknown_option_0, keyText));
			}
		}
		value, err := convertPropertyValueToJson(sourceFile, element.AsPropertyAssignment().Initializer, option, returnValue, jsonConversionNotifier)
		errors = append(errors, err...)
		if keyText != "undefined" {
			if returnValue {
				result[keyText.(string)] = value
			}
			// Notify key value set, if user asked for it
			if jsonConversionNotifier != nil {
				_, err := jsonConversionNotifier.onPropertySet(keyText.(string), value, element.AsPropertyAssignment(), *objectOption, option)
				errors = append(errors, err...)
			}
		}
	}
	return result, errors
}

// convertToJson converts the json syntax tree into the json value and report errors
// This returns the json value (apart from checking errors) only if returnValue provided is true.
// Otherwise it just checks the errors and returns undefined
func convertToJson(
	sourceFile *ast.SourceFile,
	rootExpression *ast.Expression,
	returnValue bool,
	jsonConversionNotifier *jsonConversionNotifier,
) (any, []*ast.Diagnostic) {
	if rootExpression == nil {
		if returnValue {
			return struct{}{}, nil
		} else {
			return nil, nil
		}
	}
	var jsonConversionNotifierValue *CommandLineOption
	if jsonConversionNotifier != nil {
		jsonConversionNotifierValue = &jsonConversionNotifier.rootOptions
	}
	return convertPropertyValueToJson(sourceFile, rootExpression, jsonConversionNotifierValue, returnValue, jsonConversionNotifier)
}

func isDoubleQuotedString(node *ast.Node) bool {
	return ast.IsStringLiteral(node)
}

func convertPropertyValueToJson(sourceFile *ast.SourceFile, valueExpression *ast.Expression, option *CommandLineOption, returnValue bool, jsonConversionNotifier *jsonConversionNotifier) (any, []*ast.Diagnostic) {
	switch valueExpression.Kind {
	case ast.KindTrueKeyword:
		return true, nil
	case ast.KindFalseKeyword:
		return false, nil
	case ast.KindNullKeyword: // todo: how to manage null
		return nil, nil

	case ast.KindStringLiteral:
		if !isDoubleQuotedString(valueExpression) {
			return (valueExpression.AsStringLiteral()).Text, []*ast.Diagnostic{ast.NewCompilerDiagnostic(diagnostics.String_literal_with_double_quotes_expected)}
		}
		return (valueExpression.AsStringLiteral()).Text, nil

	case ast.KindNumericLiteral:
		return valueExpression.AsNumericLiteral().Text, nil
	case ast.KindPrefixUnaryExpression:
		if valueExpression.AsPrefixUnaryExpression().Operator != ast.KindMinusToken || valueExpression.AsPrefixUnaryExpression().Operand.Kind != ast.KindNumericLiteral {
			break // not valid JSON syntax
		}
		return (valueExpression.AsPrefixUnaryExpression().Operand).AsNumericLiteral().Text, nil
	case ast.KindObjectLiteralExpression:
		objectLiteralExpression := valueExpression.AsObjectLiteralExpression()
		// Currently having element option declaration in the tsconfig with type "object"
		// determines if it needs onSetValidOptionKeyValueInParent callback or not
		// At moment there are only "compilerOptions", "typeAcquisition" and "typingOptions"
		// that satisfies it and need it to modify options set in them (for normalizing file paths)
		// vs what we set in the json
		// If need arises, we can modify this interface and callbacks as needed
		return convertObjectLiteralExpressionToJson(sourceFile, returnValue, objectLiteralExpression, option, jsonConversionNotifier)
	case ast.KindArrayLiteralExpression:
		result, errors := convertArrayLiteralExpressionToJson(
			sourceFile,
			(valueExpression.AsArrayLiteralExpression()).Elements.Nodes,
			option, // option && (option.(CommandLineOptionOfListType)).element,
			returnValue,
		)
		return result, errors
	}
	// Not in expected format
	var errors []*ast.Diagnostic
	if option != nil {
		errors = []*ast.Diagnostic{createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.Compiler_option_0_requires_a_value_of_type_1, option.Name, getCompilerOptionValueTypeString(option))}
	} else {
		errors = []*ast.Diagnostic{createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile, valueExpression, diagnostics.Property_value_can_only_be_string_literal_numeric_literal_true_false_null_object_literal_or_array_literal)}
	}
	return nil, errors
}

// ParseJsonConfigFileContent parses the contents of a config file (tsconfig.json).
// jsonNode: The contents of the config file to parse
// host: Instance of ParseConfigHost used to enumerate files in folder.
// basePath: A root directory to resolve relative path entries in the config file to. e.g. outDir
func ParseJsonConfigFileContent(json any, host ParseConfigHost, basePath string, existingOptions *core.CompilerOptions, configFileName string, resolutionStack []tspath.Path, extraFileExtensions []fileExtensionInfo, extendedConfigCache map[string]*extendedConfigCacheEntry) ParsedCommandLine {
	result := parseJsonConfigFileContentWorker(parseJsonToStringKey(json) /*sourceFile*/, nil, host, basePath, existingOptions, configFileName, resolutionStack, extraFileExtensions, extendedConfigCache)
	return result
}

// convertToObject converts the json syntax tree into the json value
func convertToObject(sourceFile *ast.SourceFile) (any, []*ast.Diagnostic) {
	var rootExpression *ast.Expression
	if sourceFile.Statements != nil {
		rootExpression = sourceFile.Statements.Nodes[0].AsExpressionStatement().Expression
	}
	return convertToJson(sourceFile, rootExpression, true /*jsonConversionNotifier*/, nil)
}

func getDefaultCompilerOptions(configFileName string) *core.CompilerOptions {
	var options *core.CompilerOptions = &core.CompilerOptions{}
	if configFileName != "" && tspath.GetBaseFileName(configFileName) == "jsconfig.json" {
		options = &core.CompilerOptions{
			AllowJs:                      core.TSTrue,
			MaxNodeModuleJsDepth:         core.TSTrue,
			AllowSyntheticDefaultImports: core.TSTrue,
			SkipLibCheck:                 core.TSTrue,
			NoEmit:                       core.TSTrue,
		}
	}
	return options
}

func convertCompilerOptionsFromJsonWorker(jsonOptions any, basePath string, configFileName string) (*core.CompilerOptions, []*ast.Diagnostic) {
	options := getDefaultCompilerOptions(configFileName)
	_, errors := convertOptionsFromJson(getCommandLineCompilerOptionsMap(), jsonOptions, basePath, options)
	if configFileName != "" {
		options.ConfigFilePath = tspath.NormalizeSlashes(configFileName)
	}
	return options, errors
}

func parseOwnConfigOfJson(
	json map[string]any,
	host ParseConfigHost,
	basePath string,
	configFileName string,
) (*parsedTsconfig, []*ast.Diagnostic) {
	var options *core.CompilerOptions
	var errors []*ast.Diagnostic
	options, errors = convertCompilerOptionsFromJsonWorker(json["compilerOptions"], basePath, configFileName)
	// typeAcquisition := convertTypeAcquisitionFromJsonWorker(json.typeAcquisition, basePath, errors, configFileName)
	// watchOptions := convertWatchOptionsFromJsonWorker(json.watchOptions, basePath, errors)
	// json.compileOnSave = convertCompileOnSaveOptionFromJson(json, basePath, errors)
	var extendedConfigPath []string
	if json["extends"] != nil || json["extends"] == "" {
		extendedConfigPath, errors = getExtendsConfigPathOrArray(json["extends"], host, basePath, configFileName, nil, nil, nil)
	} else {
		extendedConfigPath = nil
	}
	parsedConfig := &parsedTsconfig{
		raw:                json,
		options:            options,
		extendedConfigPath: extendedConfigPath,
	}
	return parsedConfig, errors
}

func readJsonConfigFile(fileName string, readFile func(fileName string) (string, bool)) (*tsConfigSourceFile, []*ast.Diagnostic) {
	text, diagnostic := TryReadFile(fileName, readFile, []*ast.Diagnostic{})
	if text != "" {
		return &tsConfigSourceFile{
			sourceFile: parser.ParseJSONText(fileName, text),
		}, diagnostic
	} else {
		file := &tsConfigSourceFile{
			sourceFile: (&ast.NodeFactory{}).NewSourceFile("", fileName, nil).AsSourceFile(),
		}
		file.sourceFile.SetDiagnostics(diagnostic)
		return file, diagnostic
	}
}

func getExtendedConfig(
	sourceFile *tsConfigSourceFile,
	extendedConfigPath string,
	host ParseConfigHost,
	resolutionStack []string,
	extendedConfigCache map[string]*extendedConfigCacheEntry,
	result *extendsResult,
) (*parsedTsconfig, []*ast.Diagnostic) {
	var path string
	if host.FS().UseCaseSensitiveFileNames() {
		path = extendedConfigPath
	} else {
		path = tspath.ToFileNameLowerCase(extendedConfigPath)
	}
	var value *extendedConfigCacheEntry
	var extendedResult *tsConfigSourceFile
	var extendedConfig *parsedTsconfig
	var errors []*ast.Diagnostic
	value = extendedConfigCache[path]
	if extendedConfigCache != nil && value != nil {
		extendedResult = value.extendedResult
		extendedConfig = value.extendedConfig
	} else {
		var err []*ast.Diagnostic
		extendedResult, err = readJsonConfigFile(extendedConfigPath, host.FS().ReadFile)
		errors = append(errors, err...)
		if extendedResult.sourceFile.Diagnostics() == nil {
			extendedConfig, err = parseConfig(nil, extendedResult, host, tspath.GetDirectoryPath(extendedConfigPath), tspath.GetBaseFileName(extendedConfigPath), resolutionStack, extendedConfigCache)
			errors = append(errors, err...)
		}
		if extendedConfigCache != nil {
			extendedConfigCache[path] = &extendedConfigCacheEntry{
				extendedResult: extendedResult,
				extendedConfig: extendedConfig,
			}
		}
	}
	if sourceFile != nil {
		if result.extendedSourceFiles == nil {
			result.extendedSourceFiles = make(map[string]struct{})
			result.extendedSourceFiles[extendedResult.sourceFile.FileName()] = struct{}{}
		}
		if len(extendedResult.extendedSourceFiles) != 0 {
			for _, extenedSourceFile := range extendedResult.extendedSourceFiles {
				result.extendedSourceFiles[extenedSourceFile] = struct{}{}
			}
		}
	}
	if extendedResult.sourceFile.Diagnostics() != nil {
		errors = append(errors, extendedResult.sourceFile.Diagnostics()...)
		return nil, errors
	}
	return extendedConfig, errors
}

// parseConfig just extracts options/include/exclude/files out of a config file.
// It does not resolve the included files.
func parseConfig(
	json map[string]any,
	sourceFile *tsConfigSourceFile,
	host ParseConfigHost,
	basePath string,
	configFileName string,
	resolutionStack []string,
	extendedConfigCache map[string]*extendedConfigCacheEntry,
) (*parsedTsconfig, []*ast.Diagnostic) {
	basePath = tspath.NormalizeSlashes(basePath)
	resolvedPath := tspath.GetNormalizedAbsolutePath(configFileName, basePath)
	var errors []*ast.Diagnostic
	if slices.Contains(resolutionStack, resolvedPath) {
		var result *parsedTsconfig
		errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Circularity_detected_while_resolving_configuration_Colon_0))
		if len(json) == 0 {
			result = &parsedTsconfig{raw: json}
		} else {
			rawResult, err := convertToObject(sourceFile.sourceFile)
			errors = append(errors, err...)
			result = &parsedTsconfig{raw: rawResult}
		}
		return result, errors
	}

	var ownConfig *parsedTsconfig
	if json != nil {
		config, err := parseOwnConfigOfJson(json, host, basePath, configFileName)
		ownConfig = config
		errors = append(errors, err...)
	} else {
		config, err := parseOwnConfigOfJsonSourceFile(sourceFile, host, basePath, configFileName)
		ownConfig = config
		errors = append(errors, err...)
	}
	if ownConfig.options != nil && ownConfig.options.Paths != nil {
		// If we end up needing to resolve relative paths from 'paths' relative to
		// the config file location, we'll need to know where that config file was.
		// Since 'paths' can be inherited from an extended config in another directory,
		// we wouldn't know which directory to use unless we store it here.
		ownConfig.options.PathsBasePath = basePath
	}

	applyExtendedConfig := func(result *extendsResult, extendedConfigPath string) {
		extendedConfig, err := getExtendedConfig(sourceFile, extendedConfigPath, host, resolutionStack, extendedConfigCache, result)
		errors = append(errors, err...)
		if extendedConfig != nil && extendedConfig.options != nil {
			extendsRaw := extendedConfig.raw
			var relativeDifference string
			setPropertyValue := func(propertyName string) {
				if rawMap, ok := ownConfig.raw.(map[string]any); ok && rawMap[propertyName] != nil {
					return
				}
				if rawMap, ok := extendedConfig.raw.(map[string]any); ok && rawMap[propertyName] != nil {
					if propertyName == "include" {
						result.include = core.Map(rawMap[propertyName].([]string), func(path string) string {
							if startsWithConfigDirTemplate(path) || tspath.IsRootedDiskPath(path) {
								return path
							} else {
								if relativeDifference == "" {
									t := tspath.ComparePathsOptions{
										UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
										CurrentDirectory:          host.GetCurrentDirectory(),
									}
									relativeDifference = tspath.ConvertToRelativePath(basePath, t)
								}
								return tspath.CombinePaths(relativeDifference, path)
							}
						})
					}
				}
			}

			setPropertyValue("include")
			setPropertyValue("exclude")
			setPropertyValue("files")
			if extendedRawMap, ok := extendsRaw.(map[string]any); ok && extendedRawMap["compileOnSave"] != nil {
				if compileOnSave, ok := extendedRawMap["compileOnSave"].(bool); ok {
					result.compileOnSave = compileOnSave
				}
			}
			mergeCompilerOptions(extendedConfig.options, result.options)
		}
	}

	if ownConfig.extendedConfigPath != nil {
		// copy the resolution stack so it is never reused between branches in potential diamond-problem scenarios.
		resolutionStack = append(resolutionStack, resolvedPath)
		var result *extendsResult = &extendsResult{
			options: &core.CompilerOptions{},
		}
		if reflect.TypeOf(ownConfig.extendedConfigPath).Kind() == reflect.String {
			applyExtendedConfig(result, ownConfig.extendedConfigPath.(string))
		} else if configPath, ok := ownConfig.extendedConfigPath.([]string); ok {
			for _, extendedConfigPath := range configPath {
				applyExtendedConfig(result, extendedConfigPath)
			}
		}
		if result.include != nil {
			ownConfig.raw = result.include
		}
		if result.exclude != nil {
			ownConfig.raw = result.exclude
		}
		if result.files != nil {
			ownConfig.raw = result.files
		}
		if ownConfig.raw == nil {
			if raw, ok := ownConfig.raw.(map[string]any); ok && raw["compileOnSave"] != nil {
				ownConfig.raw.(map[string]any)["compileOnSave"] = result.compileOnSave
			}
		}
		if sourceFile != nil && result.extendedSourceFiles != nil {
			for extendedSourceFile := range result.extendedSourceFiles {
				sourceFile.extendedSourceFiles = append(sourceFile.extendedSourceFiles, extendedSourceFile)
			}
		}
		mergeCompilerOptions(result.options, ownConfig.options)
		// ownConfig.watchOptions = ownConfig.watchOptions && result.watchOptions ?
		//     assignWatchOptions(result, ownConfig.watchOptions) :
		//     ownConfig.watchOptions || result.watchOptions;
	}
	return ownConfig, errors
}

const defaultIncludeSpec = "**/*"

type PropOfRaw struct {
	sliceValue any
	wrongValue string
}

// parseJsonConfigFileContentWorker parses the contents of a config file from json or json source file (tsconfig.json).
// json: The contents of the config file to parse
// sourceFile: sourceFile corresponding to the Json
// host: Instance of ParseConfigHost used to enumerate files in folder.
// basePath: A root directory to resolve relative path entries in the config file to. e.g. outDir
// resolutionStack: Only present for backwards-compatibility. Should be empty.
func parseJsonConfigFileContentWorker(
	json map[string]any,
	sourceFile *tsConfigSourceFile,
	host ParseConfigHost,
	basePath string,
	existingOptions *core.CompilerOptions,
	configFileName string,
	resolutionStack []tspath.Path,
	extraFileExtensions []fileExtensionInfo,
	extendedConfigCache map[string]*extendedConfigCacheEntry,
) ParsedCommandLine {
	// Debug.assert((json === undefined && sourceFile !== undefined) || (json !== undefined && sourceFile === undefined));
	var errors []*ast.Diagnostic
	resolutionStackString := []string{}
	parsedConfig, errors := parseConfig(json, sourceFile, host, basePath, configFileName, resolutionStackString, extendedConfigCache)
	mergeCompilerOptions(existingOptions, parsedConfig.options)
	// const options = handleOptionConfigDirTemplateSubstitution(
	// 	extend(existingOptions, parsedConfig.options), //function in core.ts
	// 	configDirTemplateSubstitutionOptions,
	// 	basePath,
	// )
	rawConfig := parseJsonToStringKey(parsedConfig.raw)
	var basePathForFileNames string
	if configFileName != "" {
		if parsedConfig.options != nil {
			parsedConfig.options.ConfigFilePath = tspath.NormalizeSlashes(configFileName)
		}
		basePathForFileNames = tspath.NormalizePath(directoryOfCombinedPath(configFileName, basePath))
	} else {
		basePathForFileNames = tspath.NormalizePath(basePath)
	}
	getPropFromRaw := func(prop string, validateElement func(value any) bool, elementTypeName string) PropOfRaw {
		value, exists := rawConfig[prop]
		if exists {
			if reflect.TypeOf(value).Kind() == reflect.Slice {
				result := rawConfig[prop]
				if _, ok := result.([]any); ok {
					if sourceFile == nil && !core.Every(result.([]any), validateElement) {
						errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop, elementTypeName))
					}
				}
				return PropOfRaw{sliceValue: result}
			} else {
				errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Compiler_option_0_requires_a_value_of_type_1, prop, "Array"))
				return PropOfRaw{sliceValue: nil, wrongValue: "not-array"}
			}
		}
		return PropOfRaw{sliceValue: nil, wrongValue: "no-prop"}
	}
	getConfigFileSpecs := func() configFileSpecs {
		referencesOfRaw := getPropFromRaw("references", func(element any) bool { return reflect.TypeOf(element).Kind() == reflect.Map }, "object")
		fileSpecs := getPropFromRaw("files", func(element any) bool { return reflect.TypeOf(element).Kind() == reflect.String }, "string")
		if fileSpecs.sliceValue != nil || fileSpecs.wrongValue == "" {
			hasZeroOrNoReferences := false
			if referencesOfRaw.wrongValue == "no-prop" || referencesOfRaw.wrongValue == "not-array" || len(referencesOfRaw.sliceValue.([]any)) == 0 {
				hasZeroOrNoReferences = true
			}
			hasExtends := rawConfig[string("extends")]
			if fileSpecs.sliceValue != nil && len(fileSpecs.sliceValue.([]any)) == 0 && hasZeroOrNoReferences && hasExtends == nil {
				if sourceFile != nil {
					var fileName string
					if configFileName != "" {
						fileName = configFileName
					} else {
						fileName = "tsconfig.json"
					}
					diagnosticMessage := diagnostics.The_files_list_in_config_file_0_is_empty
					nodeValue := forEachTsConfigPropArray(sourceFile, "files", func(property *ast.PropertyAssignment) *ast.Node { return property.Initializer })
					errors = append(errors, ast.NewDiagnostic(sourceFile.sourceFile, core.NewTextRange(scanner.SkipTrivia(sourceFile.sourceFile.Text, nodeValue.Pos()), nodeValue.End()), diagnosticMessage, fileName))
				} else {
					errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.The_files_list_in_config_file_0_is_empty, configFileName))
				}
			}
		}
		includeSpecs := getPropFromRaw("include", func(element any) bool { return reflect.TypeOf(element).Kind() == reflect.String }, "string")
		excludeSpecs := getPropFromRaw("exclude", func(element any) bool { return reflect.TypeOf(element).Kind() == reflect.String }, "string")
		isDefaultIncludeSpec := false
		if excludeSpecs.wrongValue == "no-prop" && parsedConfig.options != nil {
			outDir := parsedConfig.options.OutDir
			declarationDir := parsedConfig.options.DeclarationDir
			if outDir != "" || declarationDir != "" {
				values := []any{}
				if outDir != "" {
					values = append(values, outDir)
				}
				if declarationDir != "" {
					values = append(values, declarationDir)
				}
				excludeSpecs = PropOfRaw{sliceValue: values}
			}
		}
		if fileSpecs.sliceValue == nil && includeSpecs.sliceValue == nil {
			includeSpecs = PropOfRaw{sliceValue: []any{defaultIncludeSpec}}
			isDefaultIncludeSpec = true
		}
		var validatedIncludeSpecsBeforeSubstitution []string
		var validatedExcludeSpecsBeforeSubstitution []string
		var validatedFilesSpecBeforeSubstitution []string
		var validatedIncludeSpecs []string
		var validatedExcludeSpecs []string
		var validatedFilesSpec []string
		// The exclude spec list is converted into a regular expression, which allows us to quickly
		// test whether a file or directory should be excluded before recursively traversing the
		// file system.
		if includeSpecs.sliceValue != nil {
			var err []*ast.Diagnostic
			validatedIncludeSpecsBeforeSubstitution, err = validateSpecs(includeSpecs.sliceValue /*disallowTrailingRecursion*/, true, sourceFile, "include")
			errors = append(errors, err...)
			validatedIncludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
				validatedIncludeSpecsBeforeSubstitution,
				basePathForFileNames,
			)
			if validatedIncludeSpecs == nil {
				validatedIncludeSpecs = validatedIncludeSpecsBeforeSubstitution
			}
		}
		if excludeSpecs.sliceValue != nil {
			var err []*ast.Diagnostic
			validatedExcludeSpecsBeforeSubstitution, err = validateSpecs(excludeSpecs.sliceValue /*disallowTrailingRecursion*/, false, sourceFile, "exclude")
			errors = append(errors, err...)
			validatedExcludeSpecs = getSubstitutedStringArrayWithConfigDirTemplate(
				validatedExcludeSpecsBeforeSubstitution,
				basePathForFileNames,
			)
			if validatedExcludeSpecs == nil {
				validatedExcludeSpecs = validatedExcludeSpecsBeforeSubstitution
			}
		}
		if fileSpecs.sliceValue != nil {
			if _, ok := fileSpecs.sliceValue.([]any); ok {
				fileSpecs := core.Filter(fileSpecs.sliceValue.([]any), func(spec any) bool { return reflect.TypeOf(spec).Kind() == reflect.String })
				for _, spec := range fileSpecs {
					if spec, ok := spec.(string); ok {
						validatedFilesSpecBeforeSubstitution = append(validatedFilesSpecBeforeSubstitution, spec)
					}
				}
			}
			validatedFilesSpec = getSubstitutedStringArrayWithConfigDirTemplate(
				validatedFilesSpecBeforeSubstitution,
				basePathForFileNames,
			)
		}
		if validatedFilesSpec == nil {
			validatedFilesSpec = validatedFilesSpecBeforeSubstitution
		}
		return configFileSpecs{
			fileSpecs.sliceValue,
			includeSpecs.sliceValue,
			excludeSpecs.sliceValue,
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
	if sourceFile != nil {
		sourceFile.configFileSpecs = &configFileSpecs
	}

	getFileNames := func(basePath string) collections.OrderedMap[string, string] {
		var parsedConfigOptions *core.CompilerOptions = nil
		if parsedConfig.options != nil {
			parsedConfigOptions = parsedConfig.options
		}
		fileNames := getFileNamesFromConfigSpecs(configFileSpecs, basePath, parsedConfigOptions, host.FS(), extraFileExtensions)
		if shouldReportNoInputFiles(fileNames, canJsonReportNoInputFiles(rawConfig), resolutionStack) {
			includeSpecs := configFileSpecs.includeSpecs
			excludeSpecs := configFileSpecs.excludeSpecs
			if includeSpecs == nil {
				includeSpecs = []string{}
			}
			if excludeSpecs == nil {
				excludeSpecs = []string{}
			}
			errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.No_inputs_were_found_in_config_file_0_Specified_include_paths_were_1_and_exclude_paths_were_2, configFileName, core.Must(core.StringifyJson(includeSpecs)), core.Must(core.StringifyJson(excludeSpecs))))
		}
		return fileNames
	}

	getProjectReferences := func(basePath string) []core.ProjectReference {
		var projectReferences []core.ProjectReference = []core.ProjectReference{}
		referencesOfRaw := getPropFromRaw("references", func(element any) bool { return reflect.TypeOf(element).Kind() == reflect.Map }, "object")
		if referencesOfRaw.sliceValue != nil {
			for _, reference := range referencesOfRaw.sliceValue.([]any) {
				for _, ref := range parseProjectReference(reference) {
					if reflect.TypeOf(ref.Path).Kind() != reflect.String {
						errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Compiler_option_0_requires_a_value_of_type_1, "reference.path", "string"))
					} else {
						projectReferences = append(projectReferences, core.ProjectReference{
							Path:         tspath.GetNormalizedAbsolutePath(ref.Path, basePath),
							OriginalPath: ref.Path,
							Circular:     ref.Circular,
						})
					}
				}
			}
		}
		return projectReferences
	}

	return ParsedCommandLine{
		ParsedOptions: &core.ParsedOptions{
			CompilerOptions:   parsedConfig.options,
			FileNames:         getFileNames(basePathForFileNames),
			ProjectReferences: getProjectReferences(basePathForFileNames),
		},
		Raw:    parsedConfig.raw,
		Errors: errors,
	}
}

func canJsonReportNoInputFiles(rawConfig map[string]any) bool {
	_, filesExists := rawConfig["files"]
	_, referencesExists := rawConfig["references"]
	return !filesExists && !referencesExists
}

func shouldReportNoInputFiles(fileNames collections.OrderedMap[string, string], canJsonReportNoInutFiles bool, resolutionStack []tspath.Path) bool {
	return fileNames.Size() == 0 && canJsonReportNoInutFiles && (resolutionStack != nil || len(resolutionStack) == 0)
}

func validateSpecs(specs any, disallowTrailingRecursion bool, jsonSourceFile *tsConfigSourceFile, specKey string) ([]string, []*ast.Diagnostic) {
	createDiagnostic := func(message *diagnostics.Message, spec string) *ast.Diagnostic {
		element := getTsConfigPropArrayElementValue(jsonSourceFile, specKey, spec)
		return ast.NewCompilerDiagnostic(message, element)
	}
	var errors []*ast.Diagnostic
	var finalSpecs []string
	for _, spec := range specs.([]any) {
		if reflect.TypeOf(spec).Kind() != reflect.String {
			continue
		}
		diag, _ := specToDiagnostic(spec.(string), disallowTrailingRecursion)
		if diag != nil {
			errors = append(errors, createDiagnostic(diag, spec.(string)))
		}
		finalSpecs = append(finalSpecs, spec.(string))
	}
	return finalSpecs, errors
}

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

// Tests for a path that ends in a recursive directory wildcard.
//
//	Matches **, \**, **\, and \**\, but not a**b.
//	NOTE: used \ in place of / above to avoid issues with multiline comments.
//
// Breakdown:
//
//	(^|\/)      # matches either the beginning of the string or a directory separator.
//	\*\*        # matches the recursive directory wildcard "**".
//	\/?$        # matches an optional trailing directory separator at the end of the string.
const invalidTrailingRecursionPattern = `(?:^|\/)\*\*\/?$`

func getTsConfigPropArrayElementValue(tsConfigSourceFile *tsConfigSourceFile, propKey string, elementValue string) *ast.StringLiteral {
	return forEachTsConfigPropArray(tsConfigSourceFile, propKey, func(property *ast.PropertyAssignment) *ast.StringLiteral {
		if ast.IsArrayLiteralExpression(property.Initializer) {
			value := core.Find(property.Initializer.AsArrayLiteralExpression().Elements.Nodes, func(element *ast.Node) bool {
				return ast.IsStringLiteral(element) && element.AsStringLiteral().Text == elementValue
			})
			if value != nil {
				return value.AsStringLiteral()
			}
		}
		return nil
	})
}

func forEachTsConfigPropArray[T any](tsConfigSourceFile *tsConfigSourceFile, propKey string, callback func(property *ast.PropertyAssignment) *T) *T {
	if tsConfigSourceFile != nil {
		return forEachPropertyAssignment(getTsConfigObjectLiteralExpression(tsConfigSourceFile), propKey, callback)
	}
	return nil
}

func forEachPropertyAssignment[T any](objectLiteral *ast.ObjectLiteralExpression, key string, callback func(property *ast.PropertyAssignment) T, key2 ...string) T {
	if objectLiteral != nil {
		for _, property := range objectLiteral.Properties.Nodes {
			if !ast.IsPropertyAssignment(property) {
				continue
			}
			if propName, ok := ast.TryGetTextOfPropertyName(property.Name()); ok {
				if propName == key || (len(key2) > 0 && key2[0] == propName) {
					result := callback(property.AsPropertyAssignment())
					return result
				}
			}
		}
	}
	return *new(T)
}

func getTsConfigObjectLiteralExpression(tsConfigSourceFile *tsConfigSourceFile) *ast.ObjectLiteralExpression {
	if tsConfigSourceFile != nil && tsConfigSourceFile.sourceFile.Statements != nil && len(tsConfigSourceFile.sourceFile.Statements.Nodes) > 0 {
		expression := tsConfigSourceFile.sourceFile.Statements.Nodes[0].AsExpressionStatement().Expression
		return expression.AsObjectLiteralExpression()
	}
	return nil
}

func getSubstitutedPathWithConfigDirTemplate(value string, basePath string) string {
	return tspath.GetNormalizedAbsolutePath(strings.Replace(value, configDirTemplate, "./", 1), basePath)
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

// hasFileWithHigherPriorityExtension determines whether a literal or wildcard file has already been included that has a higher extension priority.
// file is the path to the file.
func hasFileWithHigherPriorityExtension(file string, literalFiles collections.OrderedMap[string, string], wildcardFiles collections.OrderedMap[string, string], extensions [][]string, keyMapper func(value string) string) bool {
	var extensionGroup []string
	for _, group := range extensions {
		if tspath.FileExtensionIsOneOf(file, group) {
			extensionGroup = append(extensionGroup, group...)
		}
	}
	if len(extensionGroup) == 0 {
		return false
	}
	for _, ext := range extensionGroup {
		// d.ts files match with .ts extension and with case sensitive sorting the file order for same files with ts tsx and dts extension is
		// d.ts, .ts, .tsx in that order so we need to handle tsx and dts of same same name case here and in remove files with same extensions
		// So dont match .d.ts files with .ts extension
		if tspath.FileExtensionIs(file, ext) && (ext != tspath.ExtensionTs || !tspath.FileExtensionIs(file, tspath.ExtensionDts)) {
			return false
		}
		higherPriorityPath := keyMapper(tspath.ChangeExtension(file, ext))
		if literalFiles.Has(higherPriorityPath) || wildcardFiles.Has(higherPriorityPath) {
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

// Removes files included via wildcard expansion with a lower extension priority that have already been included.
// file is the path to the file.
func removeWildcardFilesWithLowerPriorityExtension(file string, wildcardFiles collections.OrderedMap[string, string], extensions [][]string, keyMapper func(value string) string) {
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
		wildcardFiles.Delete(lowerPriorityPath)
	}
}

// getFileNamesFromConfigSpecs gets the file names from the provided config file specs that contain, files, include, exclude and
// other properties needed to resolve the file names
// configFileSpecs is the config file specs extracted with file names to include, wildcards to include/exclude and other details
// basePath is the base path for any relative file specifications.
// options is the Compiler options.
// host is the host used to resolve files and directories.
// extraFileExtensions optionaly file extra file extension information from host

func getFileNamesFromConfigSpecs(
	configFileSpecs configFileSpecs,
	basePath string, // considering this is the current directory
	options *core.CompilerOptions,
	host vfs.FS,
	extraFileExtensions []fileExtensionInfo,
) collections.OrderedMap[string, string] {
	extraFileExtensions = []fileExtensionInfo{}
	basePath = tspath.NormalizePath(basePath)
	// Literal file names (provided via the "files" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map later when when including
	// wildcard paths.
	var literalFileMap collections.OrderedMap[string, string]
	// Wildcard paths (provided via the "includes" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map to store paths matched
	// via wildcard, and to handle extension priority.
	var wildcardFileMap collections.OrderedMap[string, string]
	// Wildcard paths of json files (provided via the "includes" array in tsconfig.json) are stored in a
	// file map with a possibly case insensitive key. We use this map to store paths matched
	// via wildcard of *.json kind
	var wildCardJsonFileMap collections.OrderedMap[string, string]
	validatedFilesSpec := configFileSpecs.validatedFilesSpec
	validatedIncludeSpecs := configFileSpecs.validatedIncludeSpecs
	validatedExcludeSpecs := configFileSpecs.validatedExcludeSpecs
	// Rather than re-query this for each file and filespec, we query the supported extensions
	// once and store it on the expansion context.
	supportedExtensions := getSupportedExtensions(options, extraFileExtensions)
	supportedExtensionsWithJsonIfResolveJsonModule := getSupportedExtensionsWithJsonIfResolveJsonModule(options, supportedExtensions)
	// Literal files are always included verbatim. An "include" or "exclude" specification cannot
	// remove a literal file.
	for _, fileName := range validatedFilesSpec {
		file := tspath.GetNormalizedAbsolutePath(fileName, basePath)
		literalFileMap.Set(tspath.GetCanonicalFileName(fileName, host.UseCaseSensitiveFileNames()), file)
	}

	var jsonOnlyIncludeRegexes []*regexp2.Regexp
	if len(validatedIncludeSpecs) > 0 {
		files := readDirectory(host, basePath, basePath, core.Flatten(supportedExtensionsWithJsonIfResolveJsonModule), validatedExcludeSpecs, validatedIncludeSpecs, nil)
		for _, file := range files {
			if tspath.FileExtensionIs(file, tspath.ExtensionJson) {
				if jsonOnlyIncludeRegexes == nil {
					includes := core.Filter(validatedIncludeSpecs, func(include string) bool { return strings.HasSuffix(include, tspath.ExtensionJson) })
					var includeFilePatterns []string = core.Map(getRegularExpressionsForWildcards(includes, basePath, "files"), func(pattern string) string { return fmt.Sprintf("^%s$", pattern) })
					if includeFilePatterns != nil {
						jsonOnlyIncludeRegexes = core.Map(includeFilePatterns, func(pattern string) *regexp2.Regexp {
							return getRegexFromPattern(pattern, host.UseCaseSensitiveFileNames())
						})
					} else {
						jsonOnlyIncludeRegexes = nil
					}
					includeIndex := core.FindIndex(jsonOnlyIncludeRegexes, func(re *regexp2.Regexp) bool { return core.Must(re.MatchString(file)) })
					if includeIndex != -1 {
						key := tspath.GetCanonicalFileName(file, host.UseCaseSensitiveFileNames())
						if !literalFileMap.Has(key) && !wildCardJsonFileMap.Has(key) {
							wildCardJsonFileMap.Set(key, file)
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
			if !literalFileMap.Has(key) && !wildcardFileMap.Has(key) {
				wildcardFileMap.Set(key, file)
			}
		}
	}
	var files collections.OrderedMap[string, string]
	fileMaps := slices.Collect(literalFileMap.Values())
	for _, file := range fileMaps {
		files.Set(file, "")
	}
	wildcardMaps := slices.Collect(wildcardFileMap.Values())
	for _, file := range wildcardMaps {
		files.Set(file, "")
	}
	wildcarJsonMaps := slices.Collect(wildCardJsonFileMap.Values())
	for _, file := range wildcarJsonMaps {
		files.Set(file, "")
	}
	return files
}

func getSupportedExtensions(options *core.CompilerOptions, extraFileExtensions []fileExtensionInfo) [][]string {
	needJsExtensions := options.GetAllowJs()
	if len(extraFileExtensions) == 0 {
		if needJsExtensions {
			return tspath.AllSupportedExtensions
		} else {
			return tspath.SupportedTSExtensions
		}
	}
	var builtins [][]string
	if needJsExtensions {
		builtins = tspath.AllSupportedExtensions
	} else {
		builtins = tspath.SupportedTSExtensions
	}
	flatBuiltins := core.Flatten(builtins)
	var result [][]string
	for _, x := range extraFileExtensions {
		if x.scriptKind == core.ScriptKindDeferred || (needJsExtensions && (x.scriptKind == core.ScriptKindJS || x.scriptKind == core.ScriptKindJSX)) && !slices.Contains(flatBuiltins, x.extension) {
			result = append(result, []string{x.extension})
		}
	}
	extensions := slices.Concat(builtins, result)
	return extensions
}

func getSupportedExtensionsWithJsonIfResolveJsonModule(options *core.CompilerOptions, supportedExtensions [][]string) [][]string {
	if options != nil || options.GetResolveJsonModule() {
		return supportedExtensions
	}
	if core.Same(supportedExtensions, tspath.AllSupportedExtensions) {
		return tspath.AllSupportedExtensionsWithJson
	}
	if core.Same(supportedExtensions, tspath.SupportedTSExtensions) {
		return tspath.SupportedTSExtensionsWithJson
	}
	return slices.Concat(supportedExtensions, [][]string{{tspath.ExtensionJson}})
}

func createDiagnosticForNodeInSourceFileOrCompilerDiagnostic(sourceFile *ast.SourceFile, node *ast.Node, message *diagnostics.Message, args ...any) *ast.Diagnostic {
	if sourceFile != nil && node != nil {
		return ast.NewDiagnostic(sourceFile, core.NewTextRange(scanner.SkipTrivia(sourceFile.Text, node.Loc.Pos()), node.End()), message, args)
	}
	return ast.NewCompilerDiagnostic(message, args...)
}
