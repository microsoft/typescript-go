package options

import "strings"

func (parser *CommandLineParser) GetOptionsNameMap() *NameMap {
	if parser.workerDiagnostics.optionsNameMap != nil {
		return parser.workerDiagnostics.optionsNameMap
	}

	optionsNames := map[string]*CommandLineOption{}
	shortOptionNames := map[string]string{}
	for _, option := range parser.workerDiagnostics.didYouMean.OptionDeclarations {
		optionsNames[strings.ToLower(option.name)] = &option
		if option.shortName != "" {
			shortOptionNames[option.shortName] = option.name
		}
	}
	parser.workerDiagnostics.optionsNameMap = &NameMap{
		optionsNames:     optionsNames,
		shortOptionNames: shortOptionNames,
	}

	return parser.workerDiagnostics.optionsNameMap
}

type NameMap struct {
	optionsNames     map[string]*CommandLineOption
	shortOptionNames map[string]string
}

func (nm *NameMap) Get(name string) *CommandLineOption {
	return nm.optionsNames[strings.ToLower(name)]
}

func (nm *NameMap) GetShort(shortName string) *CommandLineOption {
	name, ok := nm.shortOptionNames[shortName]
	if !ok {
		return nil
	}
	return nm.optionsNames[name]
}
