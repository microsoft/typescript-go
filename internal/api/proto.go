package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/project"
)

var (
	ErrInvalidRequest = errors.New("api: invalid request")
	ErrClientError    = errors.New("api: client error")
)

type Method string

type Handle[T any] string

func ProjectHandle(p *project.Project) Handle[project.Project] {
	return createHandle[project.Project]("p", p.Name())
}

func SymbolHandle(symbol *ast.Symbol) Handle[ast.Symbol] {
	return createHandle[ast.Symbol]("s", symbol.Id.Load())
}

func TypeHandle(t *checker.Type) Handle[checker.Type] {
	return createHandle[checker.Type]("t", t.Id())
}

func createHandle[T any](prefix string, id any) Handle[T] {
	return Handle[T](fmt.Sprintf("%s%016x", prefix, id))
}

const (
	MethodConfigure Method = "configure"
	MethodRelease   Method = "release"

	MethodParseConfigFile      Method = "parseConfigFile"
	MethodLoadProject          Method = "loadProject"
	MethodGetSymbolAtPosition  Method = "getSymbolAtPosition"
	MethodGetSymbolAtPositions Method = "getSymbolAtPositions"
	MethodGetTypeOfSymbol      Method = "getTypeOfSymbol"
	MethodGetSourceFile        Method = "getSourceFile"
)

var unmarshalers = map[Method]func([]byte) (any, error){
	MethodParseConfigFile:      unmarshallerFor[ParseConfigFileParams],
	MethodRelease:              batchEnabledUnmarshallerFor[string],
	MethodLoadProject:          unmarshallerFor[LoadProjectParams],
	MethodGetSourceFile:        unmarshallerFor[GetSourceFileParams],
	MethodGetSymbolAtPosition:  batchEnabledUnmarshallerFor[GetSymbolAtPositionParams],
	MethodGetSymbolAtPositions: unmarshallerFor[GetSymbolAtPositionsParams],
	MethodGetTypeOfSymbol:      batchEnabledUnmarshallerFor[GetTypeOfSymbolParams],
}

type ConfigureParams struct {
	Callbacks []string `json:"callbacks"`
	LogFile   string   `json:"logFile"`
}

type ParseConfigFileParams struct {
	FileName string `json:"fileName"`
}

type ConfigFileResponse struct {
	FileNames []string              `json:"fileNames"`
	Options   *core.CompilerOptions `json:"options"`
}

type LoadProjectParams struct {
	ConfigFileName string `json:"configFileName"`
}

type ProjectResponse struct {
	Id              Handle[project.Project] `json:"id"`
	ConfigFileName  string                  `json:"configFileName"`
	RootFiles       []string                `json:"rootFiles"`
	CompilerOptions *core.CompilerOptions   `json:"compilerOptions"`
}

func NewProjectResponse(project *project.Project) *ProjectResponse {
	return &ProjectResponse{
		Id:              ProjectHandle(project),
		ConfigFileName:  project.Name(),
		RootFiles:       project.GetRootFileNames(),
		CompilerOptions: project.GetCompilerOptions(),
	}
}

type GetSymbolAtPositionParams struct {
	Project  Handle[project.Project] `json:"project"`
	FileName string                  `json:"fileName"`
	Position uint32                  `json:"position"`
}

type GetSymbolAtPositionsParams struct {
	Project   Handle[project.Project] `json:"project"`
	FileName  string                  `json:"fileName"`
	Positions []uint32                `json:"positions"`
}

type SymbolResponse struct {
	Id         Handle[ast.Symbol] `json:"id"`
	Name       string             `json:"name"`
	Flags      uint32             `json:"flags"`
	CheckFlags uint32             `json:"checkFlags"`
}

func NewSymbolResponse(symbol *ast.Symbol, projectVersion int) *SymbolResponse {
	return &SymbolResponse{
		Id:         SymbolHandle(symbol),
		Name:       symbol.Name,
		Flags:      uint32(symbol.Flags),
		CheckFlags: uint32(symbol.CheckFlags),
	}
}

type GetTypeOfSymbolParams struct {
	Project Handle[project.Project] `json:"project"`
	Symbol  Handle[ast.Symbol]      `json:"symbol"`
}

type TypeResponse struct {
	Id    Handle[checker.Type] `json:"id"`
	Flags uint32               `json:"flags"`
}

func NewTypeData(t *checker.Type) *TypeResponse {
	return &TypeResponse{
		Id:    TypeHandle(t),
		Flags: uint32(t.Flags()),
	}
}

type GetSourceFileParams struct {
	Project  Handle[project.Project] `json:"project"`
	FileName string                  `json:"fileName"`
}

func unmarshalPayload(method string, payload json.RawMessage) (any, error) {
	unmarshaler, ok := unmarshalers[Method(method)]
	if !ok {
		return nil, fmt.Errorf("unknown API method %q", method)
	}
	return unmarshaler(payload)
}

func batchEnabledUnmarshallerFor[T any](data []byte) (any, error) {
	if data[0] != '[' {
		return unmarshallerFor[T](data)
	}
	var v []*T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %T: %w", (*T)(nil), err)
	}
	return &v, nil
}

func unmarshallerFor[T any](data []byte) (any, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %T: %w", (*T)(nil), err)
	}
	return &v, nil
}
