package project

import (
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type SymbolData struct {
	Name           string `json:"name"`
	Flags          uint32 `json:"flags"`
	CheckFlags     uint32 `json:"checkFlags"`
	ProjectVersion uint32 `json:"projectVersion"`
}

func NewSymbolData(symbol *ast.Symbol, projectVersion int) *SymbolData {
	return &SymbolData{
		Name:           symbol.Name,
		Flags:          uint32(symbol.Flags),
		CheckFlags:     uint32(symbol.CheckFlags),
		ProjectVersion: uint32(projectVersion),
	}
}

type ProjectData struct {
	ConfigFileName  string                `json:"configFileName"`
	RootFiles       []string              `json:"rootFiles"`
	CompilerOptions *core.CompilerOptions `json:"compilerOptions"`
}

func NewProjectData(project *Project) *ProjectData {
	return &ProjectData{
		ConfigFileName:  project.configFileName,
		RootFiles:       project.GetRootFileNames(),
		CompilerOptions: project.GetCompilerOptions(),
	}
}

type API struct {
	host    ServiceHost
	options ServiceOptions

	documentRegistry *documentRegistry
	scriptInfosMu    sync.RWMutex
	scriptInfos      map[tspath.Path]*ScriptInfo

	projects map[tspath.Path]*Project
}

var _ ProjectHost = (*API)(nil)

func NewAPI(host ServiceHost, options ServiceOptions) *API {
	return &API{
		host:    host,
		options: options,
		documentRegistry: newDocumentRegistry(tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		}),
		scriptInfos: make(map[tspath.Path]*ScriptInfo),
		projects:    make(map[tspath.Path]*Project),
	}
}

// DefaultLibraryPath implements ProjectHost.
func (api *API) DefaultLibraryPath() string {
	return api.options.DefaultLibraryPath
}

// DocumentRegistry implements ProjectHost.
func (api *API) DocumentRegistry() *documentRegistry {
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
func (api *API) GetOrCreateScriptInfoForFile(fileName string, path tspath.Path, scriptKind core.ScriptKind) *ScriptInfo {
	return api.getOrCreateScriptInfo(fileName, path, scriptKind)
}

// GetScriptInfoByPath implements ProjectHost.
func (api *API) GetScriptInfoByPath(path tspath.Path) *ScriptInfo {
	api.scriptInfosMu.RLock()
	defer api.scriptInfosMu.RUnlock()
	return api.scriptInfos[path]
}

// OnDiscoveredSymlink implements ProjectHost.
func (api *API) OnDiscoveredSymlink(info *ScriptInfo) {
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

func (api *API) HandleRequest(id *lsproto.ID, params *lsproto.APIRequestParams) (any, error) {
	switch params.Method {
	case lsproto.APIMethodParseConfigFile:
		return api.ParseConfigFile(params.Params.(*lsproto.APIParseConfigFileParams).FileName)
	case lsproto.APIMethodLoadProject:
		return api.LoadProject(params.Params.(*lsproto.APILoadProjectParams).ConfigFileName)
	case lsproto.APIMethodGetSymbolAtPosition:
		params := params.Params.(*lsproto.APIGetSymbolAtPositionParams)
		return api.GetSymbolAtPosition(api.toPath(params.Project), params.FileName, int(params.Position))
	default:
		return nil, fmt.Errorf("unhandled API method %q", params.Method)
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
	project := NewConfiguredProject(configFileName, configFilePath, api)
	if err := project.loadConfig(); err != nil {
		return nil, err
	}
	api.projects[configFilePath] = project
	return NewProjectData(project), nil
}

func (api *API) GetSymbolAtPosition(project tspath.Path, fileName string, position int) (*SymbolData, error) {
	if project, ok := api.projects[project]; ok {
		symbol := project.LanguageService().GetSymbolAtPosition(fileName, position)
		if symbol == nil {
			return nil, nil
		}
		data := NewSymbolData(symbol, project.version)
		return data, nil
	}
	return nil, fmt.Errorf("project %q not found", project)
}

func (api *API) getOrCreateScriptInfo(fileName string, path tspath.Path, scriptKind core.ScriptKind) *ScriptInfo {
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
	info = newScriptInfo(fileName, path, scriptKind)
	info.setTextFromDisk(content)
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
