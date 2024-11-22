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
	kind            CommandLineOptionKind
	name, shortName string

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
	allowConfigDirTemplateSubstitution,
	// used for filter in compilerrunner
	affectsDeclarationPath,
	affectsProgramStructure,
	affectsSemanticDiagnostics,
	affectsBuildInfo,
	affectsBindDiagnostics,
	affectsSourceFile,
	affectsModuleResolution,
	affectsEmit,

	allowJsFlag,
	strictFlag bool

	// transpileoptions worker
	transpileOptionValue core.Tristate // i think this can be reduced to boolean
	// options[option.name] = option.transpileOptionValue;

	// used in listtype
	listPreserveFalsyValues bool
}

func (option *CommandLineOption) DeprecatedKeys() map[string]bool {
	if option.kind != CommandLineOptionTypeEnum {
		return nil
	}
	return CommandLineOptionDeprecated[option.name]
}
func (option *CommandLineOption) TypeMap() *collections.OrderedMap[string, any] {
	if option.kind != CommandLineOptionTypeEnum {
		return nil
	}
	return CommandLineOptionCustomType[option.name]
}
func (option *CommandLineOption) Elements() *CommandLineOption {
	if option.kind != CommandLineOptionTypeList && option.kind != CommandLineOptionTypeListOrElement {
		return nil
	}
	return CommandLineOptionElements[option.name]
}

func (option *CommandLineOption) DisallowNullOrUndefined() bool {
	return option.name == "extends"
}

// elements *CommandLineOption
var CommandLineOptionElements = map[string]*CommandLineOption{
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

// typeMap *map[string]string
var CommandLineOptionCustomType = map[string]*(collections.OrderedMap[string, any]){
	"lib":              libMap,
	"moduleResolution": moduleResolutionOptionMap,
	"module":           moduleOptionMap,
	"target":           targetOptionMap,
	"moduleDetection":  moduleDetectionOptionMap,
	"jsx":              jsxOptionMap,
	"newLine":          newLineOptionMap,
}

// deprecatedKeys map[string]bool
var CommandLineOptionDeprecated = map[string](map[string]bool){
	"moduleResolution": map[string]bool{"node": true},
	"target":           map[string]bool{"es3": true},
}

type CompilerOptionsValue any
type CustomValueType string
