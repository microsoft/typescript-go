package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type APIOptions struct {
	Logger *project.Logger
}

type API struct {
	host    APIHost
	options APIOptions

	documentRegistry *project.DocumentRegistry
	scriptInfosMu    sync.RWMutex
	scriptInfos      map[tspath.Path]*project.ScriptInfo

	projects  []*project.Project
	symbolsMu sync.Mutex
	symbols   map[Handle[ast.Symbol]]*ast.Symbol
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
		symbols:     make(map[Handle[ast.Symbol]]*ast.Symbol),
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

func (api *API) HandleRequest(id int, method string, payload []byte) ([]byte, error) {
	now := time.Now()
	params, err := unmarshalPayload(method, payload)
	if err != nil {
		return nil, err
	}
	api.options.Logger.PerfTrace(fmt.Sprintf("%s unmarshal - %s", method, time.Since(now)))

	switch Method(method) {
	case MethodGetSourceFile:
		params := params.(*GetSourceFileParams)
		sourceFile, err := api.GetSourceFile(params.Project, params.FileName)
		if err != nil {
			return nil, err
		}
		return EncodeSourceFile(sourceFile)
	case MethodParseConfigFile:
		return encodeJSON(api.ParseConfigFile(params.(*ParseConfigFileParams).FileName))
	case MethodLoadProject:
		return encodeJSON(api.LoadProject(params.(*LoadProjectParams).ConfigFileName))
	case MethodGetSymbolAtPosition:
		params := params.(*GetSymbolAtPositionParams)
		// return encodeJSON(handleBatchableRequest(params, func(params *GetSymbolAtPositionParams) (any, error) {
		return api.GetSymbolAtPosition(int(params.Project), params.FileName, int(params.Position))
		// }))
	case MethodGetTypeOfSymbol:
		params := params.(*GetTypeOfSymbolParams)
		// return encodeJSON(handleBatchableRequest(params, func(params *GetTypeOfSymbolParams) (any, error) {
		return api.GetTypeOfSymbol(int(params.Project), params.Symbol)
		// }))
	default:
		return nil, fmt.Errorf("unhandled API method %q", method)
	}
}

func (api *API) Close() {
	api.options.Logger.Close()
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
	project.GetProgram()
	id := len(api.projects)
	api.projects = append(api.projects, project)
	return NewProjectData(project, id), nil
}

func (api *API) GetSymbolAtPosition(projectId int, fileName string, position int) ([]byte, error) {
	if projectId >= len(api.projects) {
		return nil, errors.New("project not found")
	}
	project := api.projects[projectId]
	symbol, err := project.LanguageService().GetSymbolAtPosition(fileName, position)
	if err != nil || symbol == nil {
		return nil, err
	}
	id := NewHandle(symbol)
	data, err := EncodeSymbolResponse(symbol, id)
	if err != nil {
		return nil, err
	}
	api.symbolsMu.Lock()
	defer api.symbolsMu.Unlock()
	api.symbols[id] = symbol
	return data, nil
}

func (api *API) GetTypeOfSymbol(projectId int, symbolHandle Handle[ast.Symbol]) ([]byte, error) {
	if projectId >= len(api.projects) {
		return nil, errors.New("project not found")
	}
	project := api.projects[projectId]
	symbol, ok := api.symbols[symbolHandle]
	if !ok {
		return nil, fmt.Errorf("symbol %q not found", symbolHandle)
	}
	t := project.LanguageService().GetTypeOfSymbol(symbol)
	if t == nil {
		return nil, nil
	}
	// return NewTypeData(t), nil
	data, err := EncodeTypeResponse(t, NewHandle(t))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *API) GetSourceFile(projectId int, fileName string) (*ast.SourceFile, error) {
	if projectId >= len(api.projects) {
		return nil, errors.New("project not found")
	}
	project := api.projects[projectId]
	sourceFile := project.GetProgram().GetSourceFile(fileName)
	if sourceFile == nil {
		return nil, fmt.Errorf("source file %q not found", fileName)
	}
	return sourceFile, nil
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

func handleBatchableRequest[T any](params any, executeRequest func(T) (any, error)) (any, error) {
	if batchParams, ok := params.(*[]T); !ok {
		return executeRequest(params.(T))
	} else {
		var err error
		batchParams := *batchParams
		results := make([]any, len(batchParams))
		wg := core.NewWorkGroup(true /*singleThreaded*/) // !!! TODO: make GetSymbolAtLocation et al. concurrency-safe
		for i, params := range batchParams {
			wg.Queue(func() {
				var result any
				if err != nil {
					return
				}
				result, err = executeRequest(params)
				results[i] = result
			})
		}
		wg.RunAndWait()
		if err != nil {
			return nil, err
		}
		return results, nil
	}
}

func encodeJSON(v any, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}
