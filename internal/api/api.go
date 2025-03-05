package api

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type weakValuedMap[T any] map[Handle[T]]*T

type APIOptions struct {
	Logger *project.Logger
}

type API struct {
	host    APIHost
	options APIOptions

	documentRegistry *project.DocumentRegistry
	scriptInfosMu    sync.RWMutex
	scriptInfos      map[tspath.Path]*project.ScriptInfo

	projects map[tspath.Path]*project.Project
	symbols  weakValuedMap[ast.Symbol]
}

var _ project.ProjectHost = (*API)(nil)

func NewAPI(host APIHost, options APIOptions) *API {
	return &API{
		host:    host,
		options: options,
		documentRegistry: project.NewDocumentRegistry(tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		}),
		scriptInfos: make(map[tspath.Path]*project.ScriptInfo),
		projects:    make(map[tspath.Path]*project.Project),
		symbols:     make(weakValuedMap[ast.Symbol]),
	}
}

// DefaultLibraryPath implements ProjectHost.
func (api *API) DefaultLibraryPath() string {
	return api.host.DefaultLibraryPath()
}

// DocumentRegistry implements ProjectHost.
func (api *API) DocumentRegistry() *project.DocumentRegistry {
	return api.documentRegistry
}

// FS implements ProjectHost.
func (api *API) FS() vfs.FS {
	return api.host.FS()
}

// GetCurrentDirectory implements ProjectHost.
func (api *API) GetCurrentDirectory() string {
	return api.host.GetCurrentDirectory()
}

// GetOrCreateScriptInfoForFile implements ProjectHost.
func (api *API) GetOrCreateScriptInfoForFile(fileName string, path tspath.Path, scriptKind core.ScriptKind) *project.ScriptInfo {
	return api.getOrCreateScriptInfo(fileName, path, scriptKind)
}

// GetScriptInfoByPath implements ProjectHost.
func (api *API) GetScriptInfoByPath(path tspath.Path) *project.ScriptInfo {
	api.scriptInfosMu.RLock()
	defer api.scriptInfosMu.RUnlock()
	return api.scriptInfos[path]
}

// OnDiscoveredSymlink implements ProjectHost.
func (api *API) OnDiscoveredSymlink(info *project.ScriptInfo) {
	// !!!
}

// Log implements ProjectHost.
func (api *API) Log(s string) {
	api.options.Logger.Info(s)
}

// NewLine implements ProjectHost.
func (api *API) NewLine() string {
	return api.host.NewLine()
}

func (api *API) HandleRequest(id int, method string, payload json.RawMessage) (any, error) {
	params, err := unmarshalPayload(method, payload)
	if err != nil {
		return nil, err
	}

	switch Method(method) {
	case MethodParseConfigFile:
		return api.ParseConfigFile(params.(*ParseConfigFileParams).FileName)
	case MethodLoadProject:
		return api.LoadProject(params.(*LoadProjectParams).ConfigFileName)
	case MethodGetSymbolAtPosition:
		params := params.(*GetSymbolAtPositionParams)
		return api.GetSymbolAtPosition(api.toPath(params.Project), params.FileName, int(params.Position))
	case MethodGetTypeOfSymbol:
		params := params.(*GetTypeOfSymbolParams)
		return api.GetTypeOfSymbol(api.toPath(params.Project), params.Symbol)
	default:
		return nil, fmt.Errorf("unhandled API method %q", method)
	}
}

func (api *API) Close() {
	// !!!
}

func (api *API) ParseConfigFile(configFileName string) (*tsoptions.ParsedCommandLine, error) {
	configFileName = api.toAbsoluteFileName(configFileName)
	if configFileContent, ok := api.host.FS().ReadFile(configFileName); ok {
		configDir := tspath.GetDirectoryPath(configFileName)
		tsConfigSourceFile := tsoptions.NewTsconfigSourceFileFromFilePath(configFileName, api.toPath(configFileName), configFileContent)
		parsedCommandLine := tsoptions.ParseJsonSourceFileConfigFileContent(
			tsConfigSourceFile,
			api.host,
			configDir,
			nil, /*existingOptions*/
			configFileName,
			nil, /*resolutionStack*/
			nil, /*extraFileExtensions*/
			nil, /*extendedConfigCache*/
		)
		return parsedCommandLine, nil
	}
	return nil, fmt.Errorf("could not read file %q", configFileName)
}

func (api *API) LoadProject(configFileName string) (*ProjectData, error) {
	configFileName = api.toAbsoluteFileName(configFileName)
	configFilePath := api.toPath(configFileName)
	project := project.NewConfiguredProject(configFileName, configFilePath, api)
	if err := project.LoadConfig(); err != nil {
		return nil, err
	}
	api.projects[configFilePath] = project
	return NewProjectData(project), nil
}

func (api *API) GetSymbolAtPosition(projectPath tspath.Path, fileName string, position int) (*SymbolData, error) {
	project, ok := api.projects[projectPath]
	if !ok {
		return nil, fmt.Errorf("project %q not found", projectPath)
	}
	symbol, err := project.LanguageService().GetSymbolAtPosition(fileName, position)
	if err != nil || symbol == nil {
		return nil, err
	}
	data := NewSymbolData(symbol, project.Version())
	api.symbols[data.Id] = symbol
	return data, nil
}

func (api *API) GetTypeOfSymbol(projectPath tspath.Path, symbolHandle Handle[ast.Symbol]) (*TypeData, error) {
	project, ok := api.projects[projectPath]
	if !ok {
		return nil, fmt.Errorf("project %q not found", projectPath)
	}
	symbol, ok := api.symbols[symbolHandle]
	if !ok {
		return nil, fmt.Errorf("symbol %q not found", symbolHandle)
	}
	t := project.LanguageService().GetTypeOfSymbol(symbol)
	if t == nil {
		return nil, nil
	}
	return NewTypeData(t), nil
}

func (api *API) getOrCreateScriptInfo(fileName string, path tspath.Path, scriptKind core.ScriptKind) *project.ScriptInfo {
	api.scriptInfosMu.RLock()
	info, ok := api.scriptInfos[path]
	api.scriptInfosMu.RUnlock()
	if ok {
		return info
	}

	content, ok := api.host.FS().ReadFile(fileName)
	if !ok {
		return nil
	}
	info = project.NewScriptInfo(fileName, path, scriptKind)
	info.SetTextFromDisk(content)
	api.scriptInfosMu.Lock()
	defer api.scriptInfosMu.Unlock()
	api.scriptInfos[path] = info
	return info
}

func (api *API) toAbsoluteFileName(fileName string) string {
	return tspath.GetNormalizedAbsolutePath(fileName, api.host.GetCurrentDirectory())
}

func (api *API) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, api.host.GetCurrentDirectory(), api.host.FS().UseCaseSensitiveFileNames())
}
