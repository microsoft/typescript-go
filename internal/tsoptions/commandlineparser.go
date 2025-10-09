package tsoptions

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

func (p *commandLineParser) AlternateMode() *AlternateModeDiagnostics {
	return p.workerDiagnostics.didYouMean.alternateMode
}

func (p *commandLineParser) OptionsDeclarations() []*CommandLineOption {
	return p.workerDiagnostics.didYouMean.OptionDeclarations
}

func (p *commandLineParser) UnknownOptionDiagnostic() *diagnostics.Message {
	return p.workerDiagnostics.didYouMean.UnknownOptionDiagnostic
}

func (p *commandLineParser) UnknownDidYouMeanDiagnostic() *diagnostics.Message {
	return p.workerDiagnostics.didYouMean.UnknownDidYouMeanDiagnostic
}

type commandLineParser struct {
	workerDiagnostics *ParseCommandLineWorkerDiagnostics
	optionsMap        *NameMap
	fs                vfs.FS
	options           *collections.OrderedMap[string, any]
	fileNames         []string
	errors            []*ast.Diagnostic
}

func ParseCommandLine(
	commandLine []string,
	host ParseConfigHost,
) *ParsedCommandLine {
	if commandLine == nil {
		commandLine = []string{}
	}
	parser := parseCommandLineWorker(CompilerOptionsDidYouMeanDiagnostics, commandLine, host.FS())
	optionsWithAbsolutePaths := convertToOptionsWithAbsolutePaths(parser.options, CommandLineCompilerOptionsMap, host.GetCurrentDirectory())
	compilerOptions := convertMapToOptions(optionsWithAbsolutePaths, &compilerOptionsParser{&core.CompilerOptions{}}).CompilerOptions
	watchOptions := convertMapToOptions(optionsWithAbsolutePaths, &watchOptionsParser{&core.WatchOptions{}}).WatchOptions
	result := NewParsedCommandLine(compilerOptions, parser.fileNames, tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
		CurrentDirectory:          host.GetCurrentDirectory(),
	})
	result.ParsedConfig.WatchOptions = watchOptions
	result.Errors = parser.errors
	result.Raw = parser.options
	return result
}

func ParseBuildCommandLine(
	commandLine []string,
	host ParseConfigHost,
) *ParsedBuildCommandLine {
	if commandLine == nil {
		commandLine = []string{}
	}
	parser := parseCommandLineWorker(buildOptionsDidYouMeanDiagnostics, commandLine, host.FS())
	compilerOptions := &core.CompilerOptions{}
	for key, value := range parser.options.Entries() {
		buildOption := BuildNameMap.Get(key)
		if buildOption == &TscBuildOption || buildOption == CompilerNameMap.Get(key) {
			ParseCompilerOptions(key, value, compilerOptions)
		}
	}
	result := &ParsedBuildCommandLine{
		BuildOptions:    convertMapToOptions(parser.options, &buildOptionsParser{&core.BuildOptions{}}).BuildOptions,
		CompilerOptions: compilerOptions,
		WatchOptions:    convertMapToOptions(parser.options, &watchOptionsParser{&core.WatchOptions{}}).WatchOptions,
		Projects:        parser.fileNames,
		Errors:          parser.errors,

		comparePathsOptions: tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		},
	}

	if len(result.Projects) == 0 {
		// tsc -b invoked with no extra arguments; act as if invoked with "tsc -b ."
		result.Projects = append(result.Projects, ".")
	}

	// Nonsensical combinations
	if result.BuildOptions.Clean.IsTrue() && result.BuildOptions.Force.IsTrue() {
		result.Errors = append(result.Errors, ast.NewCompilerDiagnostic(diagnostics.Options_0_and_1_cannot_be_combined, "clean", "force"))
	}
	if result.BuildOptions.Clean.IsTrue() && result.BuildOptions.Verbose.IsTrue() {
		result.Errors = append(result.Errors, ast.NewCompilerDiagnostic(diagnostics.Options_0_and_1_cannot_be_combined, "clean", "verbose"))
	}
	if result.BuildOptions.Clean.IsTrue() && result.CompilerOptions.Watch.IsTrue() {
		result.Errors = append(result.Errors, ast.NewCompilerDiagnostic(diagnostics.Options_0_and_1_cannot_be_combined, "clean", "watch"))
	}
	if result.CompilerOptions.Watch.IsTrue() && result.BuildOptions.Dry.IsTrue() {
		result.Errors = append(result.Errors, ast.NewCompilerDiagnostic(diagnostics.Options_0_and_1_cannot_be_combined, "watch", "dry"))
	}

	return result
}

func parseCommandLineWorker(
	parseCommandLineWithDiagnostics *ParseCommandLineWorkerDiagnostics,
	commandLine []string,
	fs vfs.FS,
) *commandLineParser {
	parser := &commandLineParser{
		fs:                fs,
		workerDiagnostics: parseCommandLineWithDiagnostics,
		fileNames:         []string{},
		options:           &collections.OrderedMap[string, any]{},
		errors:            []*ast.Diagnostic{},
	}
	parser.optionsMap = GetNameMapFromList(parser.OptionsDeclarations())
	parser.parseStrings(commandLine)
	return parser
}

func (p *commandLineParser) parseStrings(args []string) {
	i := 0
	for i < len(args) {
		s := args[i]
		i++
		if s == "" {
			continue
		}
		switch s[0] {
		case '@':
			p.parseResponseFile(s[1:])
		case '-':
			inputOptionName := getInputOptionName(s)
			opt := p.optionsMap.GetOptionDeclarationFromName(inputOptionName, true /*allowShort*/)
			if opt != nil {
				i = p.parseOptionValue(args, i, opt, nil)
			} else {
				watchOpt := WatchNameMap.GetOptionDeclarationFromName(inputOptionName, true /*allowShort*/)
				if watchOpt != nil {
					i = p.parseOptionValue(args, i, watchOpt, watchOptionsDidYouMeanDiagnostics.OptionTypeMismatchDiagnostic)
				} else {
					p.errors = append(p.errors, p.createUnknownOptionError(inputOptionName, s, nil, nil))
				}
			}
		default:
			p.fileNames = append(p.fileNames, s)
		}
	}
}

func getInputOptionName(input string) string {
	// removes at most two leading '-' from the input string
	return strings.TrimPrefix(strings.TrimPrefix(input, "-"), "-")
}

func (p *commandLineParser) parseResponseFile(fileName string) {
	fileContents, errors := tryReadFile(fileName, func(fileName string) (string, bool) {
		if p.fs == nil {
			return "", false
		}
		read, err := p.fs.ReadFile(fileName)
		return read, err
	}, p.errors)
	p.errors = errors

	if fileContents == "" {
		return
	}

	var args []string
	text := []rune(fileContents)
	textLength := len(text)
	pos := 0
	for pos < textLength {
		for pos < textLength && text[pos] <= ' ' {
			pos++
		}
		if pos >= textLength {
			break
		}
		start := pos
		if text[pos] == '"' {
			pos++
			for pos < textLength && text[pos] != '"' {
				pos++
			}
			if pos < textLength {
				args = append(args, string(text[start+1:pos]))
				pos++
			} else {
				p.errors = append(p.errors, ast.NewCompilerDiagnostic(diagnostics.Unterminated_quoted_string_in_response_file_0, fileName))
			}
		} else {
			for text[pos] > ' ' {
				pos++
			}
			args = append(args, string(text[start:pos]))
		}
	}
	p.parseStrings(args)
}

func tryReadFile(fileName string, readFile func(string) (string, bool), errors []*ast.Diagnostic) (string, []*ast.Diagnostic) {
	// this function adds a compiler diagnostic if the file cannot be read
	text, e := readFile(fileName)

	if !e || text == "" {
		// !!! Divergence: the returned error will not give a useful message
		// errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Cannot_read_file_0_Colon_1, *e));
		text = ""
		errors = append(errors, ast.NewCompilerDiagnostic(diagnostics.Cannot_read_file_0, fileName))
	}
	return text, errors
}

func (p *commandLineParser) parseOptionValue(
	args []string,
	i int,
	opt *CommandLineOption,
	diag *diagnostics.Message,
) int {
	if opt.IsTSConfigOnly && i <= len(args) {
		optValue := ""
		if i < len(args) {
			optValue = args[i]
		}
		if optValue == "null" {
			p.options.Set(opt.Name, nil)
			i++
		} else if opt.Kind == "boolean" {
			if optValue == "false" {
				p.options.Set(opt.Name, false)
				i++
			} else {
				if optValue == "true" {
					i++
				}
				p.errors = append(p.errors, ast.NewCompilerDiagnostic(diagnostics.Option_0_can_only_be_specified_in_tsconfig_json_file_or_set_to_false_or_null_on_command_line, opt.Name))
			}
		} else {
			p.errors = append(p.errors, ast.NewCompilerDiagnostic(diagnostics.Option_0_can_only_be_specified_in_tsconfig_json_file_or_set_to_null_on_command_line, opt.Name))
			if len(optValue) != 0 && !strings.HasPrefix(optValue, "-") {
				i++
			}
		}
	} else {
		// Check to see if no argument was provided (e.g. "--locale" is the last command-line argument).
		if i >= len(args) {
			if opt.Kind != "boolean" {
				if diag == nil {
					diag = p.workerDiagnostics.OptionTypeMismatchDiagnostic
				}
				p.errors = append(p.errors, ast.NewCompilerDiagnostic(diag, opt.Name, getCompilerOptionValueTypeString(opt)))
				if opt.Kind == "list" {
					p.options.Set(opt.Name, []string{})
				} else if opt.Kind == "enum" {
					p.errors = append(p.errors, createDiagnosticForInvalidEnumType(opt, nil, nil))
				}
			} else {
				p.options.Set(opt.Name, true)
			}
			return i
		}
		if args[i] != "null" {
			switch opt.Kind {
			case "number":
				// !!! Make sure this parseInt matches JS parseInt
				num, e := strconv.Atoi(args[i])
				if e == nil {
					p.options.Set(opt.Name, num)
				}
				i++
			case "boolean":
				// boolean flag has optional value true, false, others
				optValue := args[i]

				// check next argument as boolean flag value
				if optValue == "false" {
					p.options.Set(opt.Name, false)
				} else {
					p.options.Set(opt.Name, true)
				}
				// try to consume next argument as value for boolean flag; do not consume argument if it is not "true" or "false"
				if optValue == "false" || optValue == "true" {
					i++
				}
			case "string":
				val, err := validateJsonOptionValue(opt, args[i], nil, nil)
				if err == nil {
					p.options.Set(opt.Name, val)
				} else {
					p.errors = append(p.errors, err...)
				}
				i++
			case "list":
				result, err := p.parseListTypeOption(opt, args[i])
				p.options.Set(opt.Name, result)
				p.errors = append(p.errors, err...)
				if len(result) > 0 || len(err) > 0 {
					i++
				}
			case "listOrElement":
				// If not a primitive, the possible types are specified in what is effectively a map of options.
				panic("listOrElement not supported here")
			default:
				val, err := convertJsonOptionOfEnumType(opt, strings.TrimFunc(args[i], stringutil.IsWhiteSpaceLike), nil, nil)
				p.options.Set(opt.Name, val)
				p.errors = append(p.errors, err...)
				i++
			}
		} else {
			p.options.Set(opt.Name, nil)
			i++
		}
	}
	return i
}

func (p *commandLineParser) parseListTypeOption(opt *CommandLineOption, value string) ([]any, []*ast.Diagnostic) {
	return ParseListTypeOption(opt, value)
}

func ParseListTypeOption(opt *CommandLineOption, value string) ([]any, []*ast.Diagnostic) {
	value = strings.TrimSpace(value)
	var errors []*ast.Diagnostic
	if strings.HasPrefix(value, "-") {
		return []any{}, errors
	}
	if opt.Kind == "listOrElement" && !strings.ContainsRune(value, ',') {
		val, err := validateJsonOptionValue(opt, value, nil, nil)
		if err != nil {
			return []any{}, err
		}
		return []any{val.(string)}, errors
	}
	if value == "" {
		return []any{}, errors
	}
	values := strings.Split(value, ",")
	switch opt.Elements().Kind {
	case "string":
		elements := core.MapFiltered(values, func(v string) (any, bool) {
			val, err := validateJsonOptionValue(opt.Elements(), v, nil, nil)
			if s, ok := val.(string); ok && len(err) == 0 && s != "" {
				return s, true
			}
			errors = append(errors, err...)
			return "", false
		})
		return elements, errors
	case "boolean", "object", "number":
		// do nothing: only string and enum/object types currently allowed as list entries
		// 				!!! we don't actually have number list options, so I didn't implement number list parsing
		panic("List of " + opt.Elements().Kind + " is not yet supported.")
	default:
		result := core.MapFiltered(values, func(v string) (any, bool) {
			val, err := convertJsonOptionOfEnumType(opt.Elements(), strings.TrimFunc(v, stringutil.IsWhiteSpaceLike), nil, nil)
			if s, ok := val.(string); ok && len(err) == 0 && s != "" {
				return s, true
			}
			errors = append(errors, err...)
			return "", false
		})
		return result, errors
	}
}

func convertJsonOptionOfEnumType(
	opt *CommandLineOption,
	value string,
	valueExpression *ast.Expression,
	sourceFile *ast.SourceFile,
) (any, []*ast.Diagnostic) {
	if value == "" {
		return nil, nil
	}
	key := strings.ToLower(value)
	typeMap := opt.EnumMap()
	if typeMap == nil {
		return nil, nil
	}
	val, ok := typeMap.Get(key)
	if ok {
		return validateJsonOptionValue(opt, val, valueExpression, sourceFile)
	}
	return nil, []*ast.Diagnostic{createDiagnosticForInvalidEnumType(opt, sourceFile, valueExpression)}
}

// GenerateTSConfig generates a tsconfig.json configuration when running command line "--init"
func GenerateTSConfig(options *core.CompilerOptions, newLine string) string {
	tab := "  "
	result := []string{}
	
	result = append(result, "{")
	result = append(result, tab+"// "+diagnostics.Visit_https_Colon_Slash_Slashaka_ms_Slashtsconfig_to_read_more_about_this_file.Format())
	result = append(result, tab+"\"compilerOptions\": {")
	
	emitHeader := func(header *diagnostics.Message) {
		result = append(result, tab+tab+"// "+header.Format())
	}
	
	newline := func() {
		result = append(result, "")
	}
	
	emitOption := func(setting string, defaultValue any, commented ...string) {
		commentMode := "never"
		if len(commented) > 0 {
			commentMode = commented[0]
		}
		
		var comment bool
		if commentMode == "always" {
			comment = true
		} else if commentMode == "never" {
			comment = false
		} else {
			// "optional" - always comment out for the template
			comment = true
		}
		
		formattedValue := formatValueOrArray(setting, defaultValue)
		
		if comment {
			result = append(result, tab+tab+"// \""+setting+"\": "+formattedValue+",")
		} else {
			result = append(result, tab+tab+"\""+setting+"\": "+formattedValue+",")
		}
	}
	
	// File Layout
	emitHeader(diagnostics.File_Layout)
	emitOption("rootDir", "./src", "optional")
	emitOption("outDir", "./dist", "optional")
	
	newline()
	
	// Environment Settings
	emitHeader(diagnostics.Environment_Settings)
	emitHeader(diagnostics.See_also_https_Colon_Slash_Slashaka_ms_Slashtsconfig_Slashmodule)
	emitOption("module", core.ModuleKindNodeNext)
	emitOption("target", core.ScriptTargetESNext)
	emitOption("types", []any{})
	if options.Lib != nil && len(options.Lib) > 0 {
		emitOption("lib", options.Lib)
	}
	emitHeader(diagnostics.For_nodejs_Colon)
	result = append(result, tab+tab+"// \"lib\": [\"esnext\"],")
	result = append(result, tab+tab+"// \"types\": [\"node\"],")
	emitHeader(diagnostics.X_and_npm_install_D_types_Slashnode)
	
	newline()
	
	// Other Outputs
	emitHeader(diagnostics.Other_Outputs)
	emitOption("sourceMap", true)
	emitOption("declaration", true)
	emitOption("declarationMap", true)
	
	newline()
	
	// Stricter Typechecking Options
	emitHeader(diagnostics.Stricter_Typechecking_Options)
	emitOption("noUncheckedIndexedAccess", true)
	emitOption("exactOptionalPropertyTypes", true)
	
	newline()
	
	// Style Options
	emitHeader(diagnostics.Style_Options)
	emitOption("noImplicitReturns", true, "optional")
	emitOption("noImplicitOverride", true, "optional")
	emitOption("noUnusedLocals", true, "optional")
	emitOption("noUnusedParameters", true, "optional")
	emitOption("noFallthroughCasesInSwitch", true, "optional")
	emitOption("noPropertyAccessFromIndexSignature", true, "optional")
	
	newline()
	
	// Recommended Options
	emitHeader(diagnostics.Recommended_Options)
	emitOption("strict", true)
	emitOption("jsx", core.JsxEmitReactJSX)
	emitOption("verbatimModuleSyntax", true)
	emitOption("isolatedModules", true)
	emitOption("noUncheckedSideEffectImports", true)
	emitOption("moduleDetection", core.ModuleDetectionKindForce)
	// Last option - no trailing comma
	result = append(result, tab+tab+"\"skipLibCheck\": true")
	
	result = append(result, tab+"}")
	result = append(result, "}")
	result = append(result, "")
	
	return strings.Join(result, newLine)
}

func formatValueOrArray(settingName string, value any) string {
	option := CommandLineCompilerOptionsMap.Get(settingName)
	if option == nil {
		// Fallback for unknown options
		return formatSingleValue(value, nil)
	}
	
	typeMap := option.EnumMap()
	
	// Handle array values
	if v, ok := value.([]any); ok {
		var elementMap *collections.OrderedMap[string, any]
		if element := option.Elements(); element != nil {
			elementMap = element.EnumMap()
		}
		formattedValues := make([]string, len(v))
		for i, item := range v {
			formattedValues[i] = formatSingleValue(item, elementMap)
		}
		return "[" + strings.Join(formattedValues, ", ") + "]"
	}
	
	// Handle string array (lib option)
	if v, ok := value.([]string); ok {
		formattedValues := make([]string, len(v))
		for i, item := range v {
			formattedValues[i] = formatSingleValue(item, nil)
		}
		return "[" + strings.Join(formattedValues, ", ") + "]"
	}
	
	return formatSingleValue(value, typeMap)
}

func formatSingleValue(value any, typeMap *collections.OrderedMap[string, any]) string {
	if typeMap != nil {
		// Find the key for this enum value
		for k, v := range typeMap.Entries() {
			if v == value {
				value = k
				break
			}
		}
	}
	
	// Handle different types
	switch v := value.(type) {
	case string:
		return strconv.Quote(v)
	case bool:
		return strconv.FormatBool(v)
	case int, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
	default:
		// For unknown types, use string representation
		return strconv.Quote(fmt.Sprintf("%v", v))
	}
}
