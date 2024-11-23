package options

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
)

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
	Name, shortName string
	Kind            CommandLineOptionKind

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
	extraValidation *func(value CompilerOptionsValue) (d *diagnostics.Message, args []string)

	// true or undefined
	// used for configDirTemplateSubstitutionOptions
	allowConfigDirTemplateSubstitution bool

	// used for filter in compilerrunner
	affectsDeclarationPath     bool
	affectsProgramStructure    bool
	affectsSemanticDiagnostics bool
	affectsBuildInfo           bool
	affectsBindDiagnostics     bool
	affectsSourceFile          bool
	affectsModuleResolution    bool
	affectsEmit                bool

	allowJsFlag bool
	strictFlag  bool

	// used in transpileoptions worker
	// todo: revisit to see if this can be reduced to boolean
	transpileOptionValue core.Tristate

	// used in listtype
	listPreserveFalsyValues bool
}

func (option *CommandLineOption) DeprecatedKeys() map[string]bool {
	if option.Kind != CommandLineOptionTypeEnum {
		return nil
	}
	return commandLineOptionDeprecated[option.Name]
}
func (option *CommandLineOption) EnumMap() *collections.OrderedMap[string, any] {
	if option.Kind != CommandLineOptionTypeEnum {
		return nil
	}
	return commandLineOptionEnumMap[option.Name]
}
func (option *CommandLineOption) Elements() *CommandLineOption {
	if option.Kind != CommandLineOptionTypeList && option.Kind != CommandLineOptionTypeListOrElement {
		return nil
	}
	return commandLineOptionElements[option.Name]
}

func (option *CommandLineOption) DisallowNullOrUndefined() bool {
	return option.Name == "extends"
}

// CommandLineOption.Elements()
var commandLineOptionElements = map[string]*CommandLineOption{
	"lib": {
		Name:                    "lib",
		Kind:                    CommandLineOptionTypeEnum, // libMap,
		defaultValueDescription: core.TSUnknown,
	},
	"rootDirs": {
		Name:       "rootDirs",
		Kind:       CommandLineOptionTypeString,
		isFilePath: true,
	},
	"typeRoots": {
		Name:       "typeRoots",
		Kind:       CommandLineOptionTypeString,
		isFilePath: true,
	},
	"types": {
		Name: "types",
		Kind: CommandLineOptionTypeString,
	},
	"moduleSuffixes": {
		Name: "suffix",
		Kind: CommandLineOptionTypeString,
	},
	"customConditions": {
		Name: "condition",
		Kind: CommandLineOptionTypeString,
	},
	"plugins": {
		Name: "plugin",
		Kind: CommandLineOptionTypeObject,
	},
}

// CommandLineOption.EnumMap()
var commandLineOptionEnumMap = map[string]*(collections.OrderedMap[string, any]){
	"lib":              libMap,
	"moduleResolution": moduleResolutionOptionMap,
	"module":           moduleOptionMap,
	"target":           targetOptionMap,
	"moduleDetection":  moduleDetectionOptionMap,
	"jsx":              jsxOptionMap,
	"newLine":          newLineOptionMap,
}

// CommandLineOption.DeprecatedKeys()
var commandLineOptionDeprecated = map[string](map[string]bool){
	"moduleResolution": map[string]bool{"node": true},
	"target":           map[string]bool{"es3": true},
}

// todo: revisit to see if this can be improved
type CompilerOptionsValue any
