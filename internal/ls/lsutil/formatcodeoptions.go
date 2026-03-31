package lsutil

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

type IndentStyle int

const (
	IndentStyleNone IndentStyle = iota
	IndentStyleBlock
	IndentStyleSmart
)

func parseIndentStyle(v any) IndentStyle {
	switch s := v.(type) {
	case string:
		switch strings.ToLower(s) {
		case "none":
			return IndentStyleNone
		case "block":
			return IndentStyleBlock
		case "smart":
			return IndentStyleSmart
		}
	case float64:
		return IndentStyle(int(s))
	case int:
		return IndentStyle(s)
	}
	return IndentStyleSmart
}

type SemicolonPreference string

const (
	SemicolonPreferenceIgnore SemicolonPreference = "ignore"
	SemicolonPreferenceInsert SemicolonPreference = "insert"
	SemicolonPreferenceRemove SemicolonPreference = "remove"
)

func parseSemicolonPreference(v any) SemicolonPreference {
	if s, ok := v.(string); ok {
		switch strings.ToLower(s) {
		case "ignore":
			return SemicolonPreferenceIgnore
		case "insert":
			return SemicolonPreferenceInsert
		case "remove":
			return SemicolonPreferenceRemove
		}
	}
	return SemicolonPreferenceIgnore
}

type EditorSettings struct {
	BaseIndentSize         int         `raw:"baseIndentSize" config:"format.baseIndentSize"`
	IndentSize             int         `raw:"indentSize" config:"format.indentSize"`
	TabSize                int         `raw:"tabSize" config:"format.tabSize"`
	NewLineCharacter       string      `raw:"newLineCharacter" config:"format.newLineCharacter"`
	ConvertTabsToSpaces    bool        `raw:"convertTabsToSpaces" config:"format.convertTabsToSpaces"`
	IndentStyle            IndentStyle `raw:"indentStyle" config:"format.indentStyle"`
	TrimTrailingWhitespace bool        `raw:"trimTrailingWhitespace" config:"format.trimTrailingWhitespace"`
}

type FormatCodeSettings struct {
	EditorSettings
	InsertSpaceAfterCommaDelimiter                              core.Tristate       `raw:"insertSpaceAfterCommaDelimiter" config:"format.insertSpaceAfterCommaDelimiter"`
	InsertSpaceAfterSemicolonInForStatements                    core.Tristate       `raw:"insertSpaceAfterSemicolonInForStatements" config:"format.insertSpaceAfterSemicolonInForStatements"`
	InsertSpaceBeforeAndAfterBinaryOperators                    core.Tristate       `raw:"insertSpaceBeforeAndAfterBinaryOperators" config:"format.insertSpaceBeforeAndAfterBinaryOperators"`
	InsertSpaceAfterConstructor                                 core.Tristate       `raw:"insertSpaceAfterConstructor" config:"format.insertSpaceAfterConstructor"`
	InsertSpaceAfterKeywordsInControlFlowStatements             core.Tristate       `raw:"insertSpaceAfterKeywordsInControlFlowStatements" config:"format.insertSpaceAfterKeywordsInControlFlowStatements"`
	InsertSpaceAfterFunctionKeywordForAnonymousFunctions        core.Tristate       `raw:"insertSpaceAfterFunctionKeywordForAnonymousFunctions" config:"format.insertSpaceAfterFunctionKeywordForAnonymousFunctions"`
	InsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis  core.Tristate       `raw:"insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis" config:"format.insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis"`
	InsertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets     core.Tristate       `raw:"insertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets" config:"format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets"`
	InsertSpaceAfterOpeningAndBeforeClosingNonemptyBraces       core.Tristate       `raw:"insertSpaceAfterOpeningAndBeforeClosingNonemptyBraces" config:"format.insertSpaceAfterOpeningAndBeforeClosingNonemptyBraces"`
	InsertSpaceAfterOpeningAndBeforeClosingEmptyBraces          core.Tristate       `raw:"insertSpaceAfterOpeningAndBeforeClosingEmptyBraces" config:"format.insertSpaceAfterOpeningAndBeforeClosingEmptyBraces"`
	InsertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces core.Tristate       `raw:"insertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces" config:"format.insertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces"`
	InsertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces  core.Tristate       `raw:"insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces" config:"format.insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces"`
	InsertSpaceAfterTypeAssertion                               core.Tristate       `raw:"insertSpaceAfterTypeAssertion" config:"format.insertSpaceAfterTypeAssertion"`
	InsertSpaceBeforeFunctionParenthesis                        core.Tristate       `raw:"insertSpaceBeforeFunctionParenthesis" config:"format.insertSpaceBeforeFunctionParenthesis"`
	PlaceOpenBraceOnNewLineForFunctions                         core.Tristate       `raw:"placeOpenBraceOnNewLineForFunctions" config:"format.placeOpenBraceOnNewLineForFunctions"`
	PlaceOpenBraceOnNewLineForControlBlocks                     core.Tristate       `raw:"placeOpenBraceOnNewLineForControlBlocks" config:"format.placeOpenBraceOnNewLineForControlBlocks"`
	InsertSpaceBeforeTypeAnnotation                             core.Tristate       `raw:"insertSpaceBeforeTypeAnnotation" config:"format.insertSpaceBeforeTypeAnnotation"`
	IndentMultiLineObjectLiteralBeginningOnBlankLine            core.Tristate       `raw:"indentMultiLineObjectLiteralBeginningOnBlankLine" config:"format.indentMultiLineObjectLiteralBeginningOnBlankLine"`
	Semicolons                                                  SemicolonPreference `raw:"semicolons" config:"format.semicolons"`
	IndentSwitchCase                                            core.Tristate       `raw:"indentSwitchCase" config:"format.indentSwitchCase"`
}

func FromLSFormatOptions(f *FormatCodeSettings, opt *lsproto.FormattingOptions) *FormatCodeSettings {
	updatedSettings := f.Copy()
	updatedSettings.TabSize = int(opt.TabSize)
	updatedSettings.IndentSize = int(opt.TabSize)
	updatedSettings.ConvertTabsToSpaces = opt.InsertSpaces
	if opt.TrimTrailingWhitespace != nil {
		updatedSettings.TrimTrailingWhitespace = *opt.TrimTrailingWhitespace
	}
	return updatedSettings
}

func (settings *FormatCodeSettings) ToLSFormatOptions() *lsproto.FormattingOptions {
	return &lsproto.FormattingOptions{
		TabSize:                uint32(settings.TabSize),
		InsertSpaces:           settings.ConvertTabsToSpaces,
		TrimTrailingWhitespace: &settings.TrimTrailingWhitespace,
	}
}

func (settings *FormatCodeSettings) ParseEditorSettings(editorSettings map[string]any) *FormatCodeSettings {
	if editorSettings == nil {
		return settings
	}
	for name, value := range editorSettings {
		switch strings.ToLower(name) {
		case "baseindentsize", "indentsize", "tabsize", "newlinecharacter", "converttabstospaces", "indentstyle", "trimtrailingwhitespace":
			settings.Set(name, value)
		}
	}
	return settings
}

func (settings *FormatCodeSettings) Parse(prefs any) bool {
	formatSettingsMap, ok := prefs.(map[string]any)
	formatSettingsParsed := false
	if !ok {
		return false
	}
	for name, value := range formatSettingsMap {
		formatSettingsParsed = settings.Set(name, value) || formatSettingsParsed
	}
	return formatSettingsParsed
}

func (settings *FormatCodeSettings) Set(name string, value any) bool {
	switch strings.ToLower(name) {
	case "baseindentsize":
		settings.BaseIndentSize = parseIntWithDefault(value, 0)
	case "indentsize":
		settings.IndentSize = parseIntWithDefault(value, printer.GetDefaultIndentSize())
	case "tabsize":
		settings.TabSize = parseIntWithDefault(value, printer.GetDefaultIndentSize())
	case "newlinecharacter":
		settings.NewLineCharacter = core.GetNewLineKind(tsoptions.ParseString(value)).GetNewLineCharacter()
	case "converttabstospaces":
		settings.ConvertTabsToSpaces = parseBoolWithDefault(value, true)
	case "indentstyle":
		settings.IndentStyle = parseIndentStyle(value)
	case "trimtrailingwhitespace":
		settings.TrimTrailingWhitespace = parseBoolWithDefault(value, true)
	case "insertspaceaftercommadelimiter":
		settings.InsertSpaceAfterCommaDelimiter = tsoptions.ParseTristate(value)
	case "insertspaceaftersemicoloninformstatements":
		settings.InsertSpaceAfterSemicolonInForStatements = tsoptions.ParseTristate(value)
	case "insertspacebeforeandafterbinaryoperators":
		settings.InsertSpaceBeforeAndAfterBinaryOperators = tsoptions.ParseTristate(value)
	case "insertspaceafterconstructor":
		settings.InsertSpaceAfterConstructor = tsoptions.ParseTristate(value)
	case "insertspaceafterkeywordsincontrolflowstatements":
		settings.InsertSpaceAfterKeywordsInControlFlowStatements = tsoptions.ParseTristate(value)
	case "insertspaceafterfunctionkeywordforanonymousfunctions":
		settings.InsertSpaceAfterFunctionKeywordForAnonymousFunctions = tsoptions.ParseTristate(value)
	case "insertspaceafteropeningandbeforeclosingnonemptyparenthesis":
		settings.InsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis = tsoptions.ParseTristate(value)
	case "insertspaceafteropeningandbeforeclosingnonemptybrackets":
		settings.InsertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets = tsoptions.ParseTristate(value)
	case "insertspaceafteropeningandbeforeclosingnonemptybraces":
		settings.InsertSpaceAfterOpeningAndBeforeClosingNonemptyBraces = tsoptions.ParseTristate(value)
	case "insertspaceafteropeningandbeforeclosingemptybraces":
		settings.InsertSpaceAfterOpeningAndBeforeClosingEmptyBraces = tsoptions.ParseTristate(value)
	case "insertspaceafteropeningandbeforeclosingtemplatesttringbraces":
		settings.InsertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces = tsoptions.ParseTristate(value)
	case "insertspaceafteropeningandbeforeclosingjsxexpressionbraces":
		settings.InsertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces = tsoptions.ParseTristate(value)
	case "insertspaceaftertypeassertion":
		settings.InsertSpaceAfterTypeAssertion = tsoptions.ParseTristate(value)
	case "insertspacebeforefunctionparenthesis":
		settings.InsertSpaceBeforeFunctionParenthesis = tsoptions.ParseTristate(value)
	case "placeopenbraceonnewlineforfunctions":
		settings.PlaceOpenBraceOnNewLineForFunctions = tsoptions.ParseTristate(value)
	case "placeopenbraceonnewlineforcontrolblocks":
		settings.PlaceOpenBraceOnNewLineForControlBlocks = tsoptions.ParseTristate(value)
	case "insertspacebeforetypeannotation":
		settings.InsertSpaceBeforeTypeAnnotation = tsoptions.ParseTristate(value)
	case "indentmultilineobjectliteralbeginningonblankline":
		settings.IndentMultiLineObjectLiteralBeginningOnBlankLine = tsoptions.ParseTristate(value)
	case "semicolons":
		settings.Semicolons = parseSemicolonPreference(value)
	case "indentswitchcase":
		settings.IndentSwitchCase = tsoptions.ParseTristate(value)
	default:
		return false
	}
	return true
}

func (settings *FormatCodeSettings) Copy() *FormatCodeSettings {
	if settings == nil {
		return nil
	}
	copied := *settings
	return &copied
}

func GetDefaultFormatCodeSettings() *FormatCodeSettings {
	return &FormatCodeSettings{
		EditorSettings: EditorSettings{
			IndentSize:             printer.GetDefaultIndentSize(),
			TabSize:                printer.GetDefaultIndentSize(),
			NewLineCharacter:       "\n",
			ConvertTabsToSpaces:    true,
			IndentStyle:            IndentStyleSmart,
			TrimTrailingWhitespace: true,
		},
		InsertSpaceAfterConstructor:                                 core.TSFalse,
		InsertSpaceAfterCommaDelimiter:                              core.TSTrue,
		InsertSpaceAfterSemicolonInForStatements:                    core.TSTrue,
		InsertSpaceBeforeAndAfterBinaryOperators:                    core.TSTrue,
		InsertSpaceAfterKeywordsInControlFlowStatements:             core.TSTrue,
		InsertSpaceAfterFunctionKeywordForAnonymousFunctions:        core.TSFalse,
		InsertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesis:  core.TSFalse,
		InsertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets:     core.TSFalse,
		InsertSpaceAfterOpeningAndBeforeClosingNonemptyBraces:       core.TSTrue,
		InsertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces: core.TSFalse,
		InsertSpaceAfterOpeningAndBeforeClosingJsxExpressionBraces:  core.TSFalse,
		InsertSpaceBeforeFunctionParenthesis:                        core.TSFalse,
		PlaceOpenBraceOnNewLineForFunctions:                         core.TSFalse,
		PlaceOpenBraceOnNewLineForControlBlocks:                     core.TSFalse,
		Semicolons:                                                  SemicolonPreferenceIgnore,
		IndentSwitchCase:                                            core.TSTrue,
	}
}
