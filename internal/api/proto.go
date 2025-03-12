package api

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unsafe"

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

type Handle[T any] uint32

func NewHandle[T any](v *T) Handle[T] {
	return Handle[T](uintptr(unsafe.Pointer(v)))
}

const (
	MethodConfigure           Method = "configure"
	MethodParseConfigFile     Method = "parseConfigFile"
	MethodLoadProject         Method = "loadProject"
	MethodGetSymbolAtPosition Method = "getSymbolAtPosition"
	MethodGetTypeOfSymbol     Method = "getTypeOfSymbol"
	MethodGetSourceFile       Method = "getSourceFile"
)

var unmarshalers = map[Method]func([]byte) (any, error){
	MethodParseConfigFile:     unmarshallerFor[ParseConfigFileParams],
	MethodLoadProject:         unmarshallerFor[LoadProjectParams],
	MethodGetSourceFile:       unmarshallerFor[GetSourceFileParams],
	MethodGetSymbolAtPosition: decodeGetSymbolAtPositionParams,
	MethodGetTypeOfSymbol:     decodeGetTypeOfSymbolParams,
}

type ConfigureParams struct {
	Callbacks []string `json:"callbacks"`
	LogFile   string   `json:"logFile"`
}

type ParseConfigFileParams struct {
	FileName string `json:"fileName"`
}

type LoadProjectParams struct {
	ConfigFileName string `json:"configFileName"`
}

type ProjectData struct {
	Id              int                   `json:"id"`
	ConfigFileName  string                `json:"configFileName"`
	RootFiles       []string              `json:"rootFiles"`
	CompilerOptions *core.CompilerOptions `json:"compilerOptions"`
}

func NewProjectData(project *project.Project, id int) *ProjectData {
	return &ProjectData{
		Id:              id,
		ConfigFileName:  project.Name(),
		RootFiles:       project.GetRootFileNames(),
		CompilerOptions: project.GetCompilerOptions(),
	}
}

type GetSymbolAtPositionParams struct {
	Project  uint16 `json:"project"`
	Position uint32 `json:"position"`
	FileName string `json:"fileName"`
}

type SymbolData struct {
	Id         Handle[ast.Symbol] `json:"id"`
	Name       string             `json:"name"`
	Flags      uint32             `json:"flags"`
	CheckFlags uint32             `json:"checkFlags"`
}

func NewSymbolData(symbol *ast.Symbol) *SymbolData {
	return &SymbolData{
		Id:         NewHandle(symbol),
		Name:       symbol.Name,
		Flags:      uint32(symbol.Flags),
		CheckFlags: uint32(symbol.CheckFlags),
	}
}

func decodeGetSymbolAtPositionParams(data []byte) (any, error) {
	var v GetSymbolAtPositionParams
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &v.Project); err != nil {
		return nil, fmt.Errorf("failed to read Project: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &v.Position); err != nil {
		return nil, fmt.Errorf("failed to read Position: %w", err)
	}
	fileName, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read fileName: %w", err)
	}
	v.FileName = string(fileName)
	return &v, nil
}

type GetTypeOfSymbolParams struct {
	Project uint16             `json:"project"`
	Symbol  Handle[ast.Symbol] `json:"symbol"`
}

func decodeGetTypeOfSymbolParams(data []byte) (any, error) {
	var v GetTypeOfSymbolParams
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &v.Project); err != nil {
		return nil, fmt.Errorf("failed to read Project: %w", err)
	}
	if err := binary.Read(reader, binary.LittleEndian, &v.Symbol); err != nil {
		return nil, fmt.Errorf("failed to read Symbol: %w", err)
	}
	return &v, nil
}

type TypeData struct {
	Id    Handle[checker.Type] `json:"id"`
	Flags uint32               `json:"flags"`
}

func NewTypeData(t *checker.Type) *TypeData {
	return &TypeData{
		Id:    NewHandle(t),
		Flags: uint32(t.Flags()),
	}
}

type GetSourceFileParams struct {
	Project  int    `json:"project"`
	FileName string `json:"fileName"`
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
