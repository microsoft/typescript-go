package execute

import (
	"fmt"
	"io"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func getFormatOptsOfSys(sys System) *diagnosticwriter.FormattingOptions {
	return &diagnosticwriter.FormattingOptions{
		NewLine: "\n",
		ComparePathsOptions: tspath.ComparePathsOptions{
			CurrentDirectory:          sys.GetCurrentDirectory(),
			UseCaseSensitiveFileNames: sys.FS().UseCaseSensitiveFileNames(),
		},
	}
}

type diagnosticReporter = func(*ast.Diagnostic)

func quietDiagnosticReporter(diagnostic *ast.Diagnostic) {}

func createDiagnosticReporter(sys System, options *core.CompilerOptions) diagnosticReporter {
	if options.Quiet.IsTrue() {
		return quietDiagnosticReporter
	}

	formatOpts := getFormatOptsOfSys(sys)
	writeDiagnostic := core.IfElse(shouldBePretty(sys, options), diagnosticwriter.FormatDiagnosticWithColorAndContext, diagnosticwriter.WriteFormatDiagnostic)
	return func(diagnostic *ast.Diagnostic) {
		writeDiagnostic(sys.Writer(), diagnostic, formatOpts)
		fmt.Fprint(sys.Writer(), formatOpts.NewLine)
	}
}

func shouldBePretty(sys System, options *core.CompilerOptions) bool {
	if options == nil || options.Pretty.IsTrueOrUnknown() {
		// todo: return defaultIsPretty(sys);
		return true
	}
	return false
}

type diagnosticsReporter = func(diagnostics []*ast.Diagnostic)

func quietDiagnosticsReporter(diagnostics []*ast.Diagnostic) {}

func createReportErrorSummary(sys System, options *core.CompilerOptions) diagnosticsReporter {
	if !options.Quiet.IsTrue() && shouldBePretty(sys, options) {
		formatOpts := getFormatOptsOfSys(sys)
		return func(diagnostics []*ast.Diagnostic) {
			diagnosticwriter.WriteErrorSummaryText(sys.Writer(), diagnostics, formatOpts)
		}
	}
	return quietDiagnosticsReporter
}

func createBuilderStatusReporter(sys System, options *core.CompilerOptions) diagnosticReporter {
	if options.Quiet.IsTrue() {
		return quietDiagnosticReporter
	}

	formatOpts := getFormatOptsOfSys(sys)
	writeStatus := core.IfElse(shouldBePretty(sys, options), diagnosticwriter.FormatDiagnosticsStatusWithColorAndTime, diagnosticwriter.FormatDiagnosticsStatusAndTime)
	return func(diagnostic *ast.Diagnostic) {
		writeStatus(sys.Writer(), sys.Now().Format("03:04:05 PM"), diagnostic, formatOpts)
		fmt.Fprint(sys.Writer(), formatOpts.NewLine, formatOpts.NewLine)
	}
}

type statistics struct {
	isAggregate      bool
	projects         int
	projectsBuilt    int
	timestampUpdates int
	files            int
	lines            int
	identifiers      int
	symbols          int
	types            int
	instantiations   int
	memoryUsed       uint64
	memoryAllocs     uint64
	compileTimes     *compileTimes
}

func statisticsFromProgram(program *compiler.Program, compileTimes *compileTimes, memStats *runtime.MemStats) *statistics {
	return &statistics{
		files:          len(program.SourceFiles()),
		lines:          program.LineCount(),
		identifiers:    program.IdentifierCount(),
		symbols:        program.SymbolCount(),
		types:          program.TypeCount(),
		instantiations: program.InstantiationCount(),
		memoryUsed:     memStats.Alloc,
		memoryAllocs:   memStats.Mallocs,
		compileTimes:   compileTimes,
	}
}

func (p *statistics) report(ioWriter io.Writer, testing bool) {
	if testing {
		return
	}
	var stats table
	var prefix string

	if p.isAggregate {
		prefix = "Aggregate "
		stats.add("Projects in scope", p.projects)
		stats.add("Projects built", p.projectsBuilt)
		stats.add("Timestamps only updates", p.timestampUpdates)
	}
	stats.add(prefix+"Files", p.files)
	stats.add(prefix+"Lines", p.lines)
	stats.add(prefix+"Identifiers", p.identifiers)
	stats.add(prefix+"Symbols", p.symbols)
	stats.add(prefix+"Types", p.types)
	stats.add(prefix+"Instantiations", p.instantiations)
	stats.add(prefix+"Memory used", fmt.Sprintf("%vK", p.memoryUsed/1024))
	stats.add(prefix+"Memory allocs", strconv.FormatUint(p.memoryAllocs, 10))
	if p.compileTimes.configTime != 0 {
		stats.add(prefix+"Config time", p.compileTimes.configTime)
	}
	if p.compileTimes.buildInfoReadTime != 0 {
		stats.add(prefix+"BuildInfo read time", p.compileTimes.buildInfoReadTime)
	}
	stats.add(prefix+"Parse time", p.compileTimes.parseTime)
	if p.compileTimes.bindTime != 0 {
		stats.add(prefix+"Bind time", p.compileTimes.bindTime)
	}
	if p.compileTimes.checkTime != 0 {
		stats.add(prefix+"Check time", p.compileTimes.checkTime)
	}
	if p.compileTimes.emitTime != 0 {
		stats.add(prefix+"Emit time", p.compileTimes.emitTime)
	}
	if p.compileTimes.changesComputeTime != 0 {
		stats.add(prefix+"Changes compute time", p.compileTimes.changesComputeTime)
	}
	stats.add(prefix+"Total time", p.compileTimes.totalTime)
	stats.print(ioWriter)
}

func printVersion(sys System) {
	fmt.Fprintln(sys.Writer(), diagnostics.Version_0.Format(core.Version()))
}

func printHelp(sys System, commandLine *tsoptions.ParsedCommandLine) {
	if commandLine.CompilerOptions().All.IsFalseOrUnknown() {
		printEasyHelp(sys, getOptionsForHelp(commandLine))
	} else {
		// !!! printAllHelp(sys, getOptionsForHelp(commandLine))
	}
}

func getOptionsForHelp(commandLine *tsoptions.ParsedCommandLine) []*tsoptions.CommandLineOption {
	// Sort our options by their names, (e.g. "--noImplicitAny" comes before "--watch")
	opts := slices.Clone(tsoptions.OptionsDeclarations)
	opts = append(opts, &tsoptions.TscBuildOption)

	if commandLine.CompilerOptions().All.IsTrue() {
		slices.SortFunc(opts, func(a, b *tsoptions.CommandLineOption) int {
			return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		})
		return opts
	} else {
		return core.Filter(opts, func(opt *tsoptions.CommandLineOption) bool {
			return opt.ShowInSimplifiedHelpView
		})
	}
}

func getHeader(sys System, message string) []string {
	// !!! const colors = createColors(sys);
	var header []string
	// !!! terminalWidth := sys.GetWidthOfTerminal?.() ?? 0
	const tsIconLength = 5

	//     const tsIconFirstLine = colors.blueBackground("".padStart(tsIconLength));
	//     const tsIconSecondLine = colors.blueBackground(colors.brightWhite("TS ".padStart(tsIconLength)));
	//     // If we have enough space, print TS icon.
	//     if (terminalWidth >= message.length + tsIconLength) {
	//         // right align of the icon is 120 at most.
	//         const rightAlign = terminalWidth > 120 ? 120 : terminalWidth;
	//         const leftAlign = rightAlign - tsIconLength;
	//         header.push(message.padEnd(leftAlign) + tsIconFirstLine + sys.newLine);
	//         header.push("".padStart(leftAlign) + tsIconSecondLine + sys.newLine);
	//     }
	//     else {
	header = append(header, message+"\n", "\n")
	//     }
	return header
}

func printEasyHelp(sys System, simpleOptions []*tsoptions.CommandLineOption) {
	// !!! const colors = createColors(sys);
	var output []string
	example := func(examples []string, desc *diagnostics.Message) {
		for _, example := range examples {
			// !!! colors
			// output.push("  " + colors.blue(example) + sys.newLine);
			output = append(output, "  ", example, "\n")
		}
		output = append(output, "  ", desc.Format(), "\n", "\n")
	}

	msg := diagnostics.X_tsc_Colon_The_TypeScript_Compiler.Format() + " - " + diagnostics.Version_0.Format(core.Version())
	output = append(output, getHeader(sys, msg)...)

	output = append(output /*colors.bold(*/, diagnostics.COMMON_COMMANDS.Format() /*)*/, "\n", "\n")

	example([]string{"tsc"}, diagnostics.Compiles_the_current_project_tsconfig_json_in_the_working_directory)
	example([]string{"tsc app.ts util.ts"}, diagnostics.Ignoring_tsconfig_json_compiles_the_specified_files_with_default_compiler_options)
	example([]string{"tsc -b"}, diagnostics.Build_a_composite_project_in_the_working_directory)
	example([]string{"tsc --init"}, diagnostics.Creates_a_tsconfig_json_with_the_recommended_settings_in_the_working_directory)
	example([]string{"tsc -p ./path/to/tsconfig.json"}, diagnostics.Compiles_the_TypeScript_project_located_at_the_specified_path)
	example([]string{"tsc --help --all"}, diagnostics.An_expanded_version_of_this_information_showing_all_possible_compiler_options)
	example([]string{"tsc --noEmit", "tsc --target esnext"}, diagnostics.Compiles_the_current_project_with_additional_settings)

	var cliCommands []*tsoptions.CommandLineOption
	var configOpts []*tsoptions.CommandLineOption
	for _, opt := range simpleOptions {
		if opt.IsCommandLineOnly || opt.Category == diagnostics.Command_line_Options {
			cliCommands = append(cliCommands, opt)
		} else {
			configOpts = append(configOpts, opt)
		}
	}

	output = append(output, generateSectionOptionsOutput(sys, diagnostics.COMMAND_LINE_FLAGS.Format(), cliCommands /*subCategory*/, false /*beforeOptionsDescription*/, nil /*afterOptionsDescription*/, nil)...)

	after := diagnostics.You_can_learn_about_all_of_the_compiler_options_at_0.Format("https://aka.ms/tsc")
	output = append(output, generateSectionOptionsOutput(sys, diagnostics.COMMON_COMPILER_OPTIONS.Format(), configOpts /*subCategory*/, false /*beforeOptionsDescription*/, nil,
		// !!! locale formatMessage(Diagnostics.You_can_learn_about_all_of_the_compiler_options_at_0, "https://aka.ms/tsc")),
		&after)...)

	for _, chunk := range output {
		fmt.Fprint(sys.Writer(), chunk)
	}
}

func printBuildHelp(sys System, buildOptions []*tsoptions.CommandLineOption) {
	var output []string
	output = append(output, getHeader(sys, diagnostics.X_tsc_Colon_The_TypeScript_Compiler.Format()+" - "+diagnostics.Version_0.Format(core.Version()))...)
	output = append(output, generateSectionOptionsOutput(
		sys,
		diagnostics.BUILD_OPTIONS.Format(),
		core.Filter(buildOptions, func(option *tsoptions.CommandLineOption) bool {
			return option != &tsoptions.TscBuildOption
		}),
		false,
		nil, // !!! locale formatMessage(Diagnostics.Using_build_b_will_make_tsc_behave_more_like_a_build_orchestrator_than_a_compiler_This_is_used_to_trigger_building_composite_projects_which_you_can_learn_more_about_at_0, "https://aka.ms/tsc-composite-builds")),
		nil,
	)...)

	for _, chunk := range output {
		fmt.Fprint(sys.Writer(), chunk)
	}
}

func generateSectionOptionsOutput(
	sys System,
	sectionName string,
	options []*tsoptions.CommandLineOption,
	subCategory bool,
	beforeOptionsDescription,
	afterOptionsDescription *string,
) (output []string) {
	// !!! color
	output = append(output /*createColors(sys).bold(*/, sectionName /*)*/, "\n", "\n")

	if beforeOptionsDescription != nil {
		output = append(output, *beforeOptionsDescription, "\n", "\n")
	}
	if !subCategory {
		output = append(output, generateGroupOptionOutput(sys, options)...)
		if afterOptionsDescription != nil {
			output = append(output, *afterOptionsDescription, "\n", "\n")
		}
		return output
	}
	categoryMap := make(map[string][]*tsoptions.CommandLineOption)
	for _, option := range options {
		if option.Category == nil {
			continue
		}
		curCategory := option.Category.Format()
		categoryMap[curCategory] = append(categoryMap[curCategory], option)
	}
	for key, value := range categoryMap {
		output = append(output, "### ", key, "\n", "\n")
		output = append(output, generateGroupOptionOutput(sys, value)...)
	}
	if afterOptionsDescription != nil {
		output = append(output, *afterOptionsDescription, "\n", "\n")
	}

	return output
}

func generateGroupOptionOutput(sys System, optionsList []*tsoptions.CommandLineOption) []string {
	var maxLength int
	for _, option := range optionsList {
		curLenght := len(getDisplayNameTextOfOption(option))
		maxLength = max(curLenght, maxLength)
	}

	// left part should be right align, right part should be left align

	// assume 2 space between left margin and left part.
	rightAlignOfLeftPart := maxLength + 2
	// assume 2 space between left and right part
	leftAlignOfRightPart := rightAlignOfLeftPart + 2

	var lines []string
	for _, option := range optionsList {
		tmp := generateOptionOutput(sys, option, rightAlignOfLeftPart, leftAlignOfRightPart)
		lines = append(lines, tmp...)
	}

	// make sure always a blank line in the end.
	if len(lines) < 2 || lines[len(lines)-2] != "\n" {
		lines = append(lines, "\n")
	}

	return lines
}

func generateOptionOutput(
	sys System,
	option *tsoptions.CommandLineOption,
	rightAlignOfLeftPart, leftAlignOfRightPart int,
) []string {
	var text []string
	// !!! const colors = createColors(sys);

	// name and description
	name := getDisplayNameTextOfOption(option)

	// value type and possible value
	valueCandidates := getValueCandidate(option)

	var defaultValueDescription string
	if msg, ok := option.DefaultValueDescription.(*diagnostics.Message); ok && msg != nil {
		defaultValueDescription = msg.Format()
	} else {
		defaultValueDescription = formatDefaultValue(
			option.DefaultValueDescription,
			core.IfElse(
				option.Kind == tsoptions.CommandLineOptionTypeList || option.Kind == tsoptions.CommandLineOptionTypeListOrElement,
				option.Elements(), option,
			),
		)
	}

	var terminalWidth int
	// !!! const terminalWidth = sys.getWidthOfTerminal?.() ?? 0;

	// Note: child_process might return `terminalWidth` as undefined.
	if terminalWidth >= 80 {
		// !!!     let description = "";
		// !!!     if (option.description) {
		// !!!         description = getDiagnosticText(option.description);
		// !!!     }
		// !!!     text.push(...getPrettyOutput(name, description, rightAlignOfLeft, leftAlignOfRight, terminalWidth, /*colorLeft*/ true), sys.newLine);
		// !!!     if (showAdditionalInfoOutput(valueCandidates, option)) {
		// !!!         if (valueCandidates) {
		// !!!             text.push(...getPrettyOutput(valueCandidates.valueType, valueCandidates.possibleValues, rightAlignOfLeft, leftAlignOfRight, terminalWidth, /*colorLeft*/ false), sys.newLine);
		// !!!         }
		// !!!         if (defaultValueDescription) {
		// !!!             text.push(...getPrettyOutput(getDiagnosticText(Diagnostics.default_Colon), defaultValueDescription, rightAlignOfLeft, leftAlignOfRight, terminalWidth, /*colorLeft*/ false), sys.newLine);
		// !!!         }
		// !!!     }
		// !!!     text.push(sys.newLine);
	} else {
		text = append(text /* !!! colors.blue(name) */, name, "\n")
		if option.Description != nil {
			text = append(text, option.Description.Format())
		}
		text = append(text, "\n")
		if showAdditionalInfoOutput(valueCandidates, option) {
			if valueCandidates != nil {
				text = append(text, valueCandidates.valueType, " ", valueCandidates.possibleValues)
			}
			if defaultValueDescription != "" {
				if valueCandidates != nil {
					text = append(text, "\n")
				}
				text = append(text, diagnostics.X_default_Colon.Format(), " ", defaultValueDescription)
			}

			text = append(text, "\n")
		}
		text = append(text, "\n")
	}

	return text
}

func formatDefaultValue(defaultValue any, option *tsoptions.CommandLineOption) string {
	if defaultValue == nil || defaultValue == core.TSUnknown {
		return "undefined"
	}

	if option.Kind == tsoptions.CommandLineOptionTypeEnum {
		// e.g. ScriptTarget.ES2015 -> "es6/es2015"
		var names []string
		for name, value := range option.EnumMap().Entries() {
			if value == defaultValue {
				names = append(names, name)
			}
		}
		return strings.Join(names, "/")
	}
	return fmt.Sprintf("%v", defaultValue)
}

type valueCandidate struct {
	// "one or more" or "any of"
	valueType      string
	possibleValues string
}

func showAdditionalInfoOutput(valueCandidates *valueCandidate, option *tsoptions.CommandLineOption) bool {
	if option.Category == diagnostics.Command_line_Options {
		return false
	}
	if valueCandidates != nil && valueCandidates.possibleValues == "string" &&
		(option.DefaultValueDescription == nil ||
			option.DefaultValueDescription == "false" ||
			option.DefaultValueDescription == "n/a") {
		return false
	}
	return true
}

func getValueCandidate(option *tsoptions.CommandLineOption) *valueCandidate {
	// option.type might be "string" | "number" | "boolean" | "object" | "list" | Map<string, number | string>
	// string -- any of: string
	// number -- any of: number
	// boolean -- any of: boolean
	// object -- null
	// list -- one or more: , content depends on `option.element.type`, the same as others
	// Map<string, number | string> -- any of: key1, key2, ....
	if option.Kind == tsoptions.CommandLineOptionTypeObject {
		return nil
	}

	res := &valueCandidate{}
	if option.Kind == tsoptions.CommandLineOptionTypeListOrElement {
		// assert(option.type !== "listOrElement")
		panic("no value candidate for list or element")
	}

	switch option.Kind {
	case tsoptions.CommandLineOptionTypeString,
		tsoptions.CommandLineOptionTypeNumber,
		tsoptions.CommandLineOptionTypeBoolean:
		res.valueType = diagnostics.X_type_Colon.Format()
	case tsoptions.CommandLineOptionTypeList:
		res.valueType = diagnostics.X_one_or_more_Colon.Format()
	default:
		res.valueType = diagnostics.X_one_of_Colon.Format()
	}

	res.possibleValues = getPossibleValues(option)

	return res
}

func getPossibleValues(option *tsoptions.CommandLineOption) string {
	switch option.Kind {
	case tsoptions.CommandLineOptionTypeString,
		tsoptions.CommandLineOptionTypeNumber,
		tsoptions.CommandLineOptionTypeBoolean:
		return string(option.Kind)
	case tsoptions.CommandLineOptionTypeList,
		tsoptions.CommandLineOptionTypeListOrElement:
		return getPossibleValues(option.Elements())
	case tsoptions.CommandLineOptionTypeObject:
		return ""
	default:
		// Map<string, number | string>
		// Group synonyms: es6/es2015
		enumMap := option.EnumMap()
		inverted := collections.NewOrderedMapWithSizeHint[any, []string](enumMap.Size())
		deprecatedKeys := option.DeprecatedKeys()

		for name, value := range enumMap.Entries() {
			if deprecatedKeys == nil || !deprecatedKeys.Has(name) {
				inverted.Set(value, append(inverted.GetOrZero(value), name))
			}
		}
		var syns []string
		for synonyms := range inverted.Values() {
			syns = append(syns, strings.Join(synonyms, "/"))
		}
		return strings.Join(syns, ", ")
	}
}

func getDisplayNameTextOfOption(option *tsoptions.CommandLineOption) string {
	return "--" + option.Name + core.IfElse(option.ShortName != "", ", -"+option.ShortName, "")
}

type tableRow struct {
	name  string
	value string
}

type table struct {
	rows []tableRow
}

func (t *table) add(name string, value any) {
	if d, ok := value.(time.Duration); ok {
		value = formatDuration(d)
	}
	t.rows = append(t.rows, tableRow{name, fmt.Sprint(value)})
}

func (t *table) print(w io.Writer) {
	nameWidth := 0
	valueWidth := 0
	for _, r := range t.rows {
		nameWidth = max(nameWidth, len(r.name))
		valueWidth = max(valueWidth, len(r.value))
	}

	for _, r := range t.rows {
		fmt.Fprintf(w, "%-*s %*s\n", nameWidth+1, r.name+":", valueWidth, r.value)
	}
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3fs", d.Seconds())
}

func identifierCount(p *compiler.Program) int {
	count := 0
	for _, file := range p.SourceFiles() {
		count += file.IdentifierCount
	}
	return count
}
