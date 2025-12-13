package lsproto

import (
	//nolint
	"encoding/json"
	"fmt"
)

const (
	MethodHandleCustomLspApiCommand Method = "$/handleCustomLspApiCommand"
)

type HandleCustomLspApiCommandParams struct {
	LspApiCommand LspApiCommand `json:"lspApiCommand"`
	Arguments     interface{}   `json:"args"`
}

type LspApiCommand string

const (
	CommandGetElementType             LspApiCommand = "getElementType"
	CommandGetSymbolType              LspApiCommand = "getSymbolType"
	CommandGetTypeProperties          LspApiCommand = "getTypeProperties"
	CommandGetTypeProperty            LspApiCommand = "getTypeProperty"
	CommandAreTypesMutuallyAssignable LspApiCommand = "areTypesMutuallyAssignable"
	CommandGetResolvedSignature       LspApiCommand = "getResolvedSignature"
)

type TypeRequestKind string

const (
	TypeRequestKindDefault               TypeRequestKind = "Default"
	TypeRequestKindContextual            TypeRequestKind = "Contextual"
	TypeRequestKindContextualCompletions TypeRequestKind = "ContextualCompletions"
)

type GetElementTypeArguments struct {
	File            DocumentUri     `json:"file"`
	Range           Range           `json:"range"`
	TypeRequestKind TypeRequestKind `json:"typeRequestKind"`
	ProjectFileName *DocumentUri    `json:"projectFileName,omitempty"`
	ForceReturnType bool            `json:"forceReturnType"`
}

type GetSymbolTypeArguments struct {
	TypeCheckerId int `json:"typeCheckerId"`
	ProjectId     int `json:"projectId"`
	SymbolId      int `json:"symbolId"`
}

type GetTypePropertiesArguments struct {
	TypeCheckerId int `json:"typeCheckerId"`
	ProjectId     int `json:"projectId"`
	TypeId        int `json:"typeId"`
}

type GetTypePropertyArguments struct {
	TypeCheckerId int    `json:"typeCheckerId"`
	ProjectId     int    `json:"projectId"`
	TypeId        int    `json:"typeId"`
	PropertyName  string `json:"propertyName"`
}

type AreTypesMutuallyAssignableArguments struct {
	TypeCheckerId int `json:"typeCheckerId"`
	ProjectId     int `json:"projectId"`
	Type1Id       int `json:"type1Id"`
	Type2Id       int `json:"type2Id"`
}

type GetResolvedSignatureArguments struct {
	File            DocumentUri  `json:"file"`
	Range           Range        `json:"range"`
	ProjectFileName *DocumentUri `json:"projectFileName"`
}

func (p *HandleCustomLspApiCommandParams) UnmarshalJSON(data []byte) error {
	// First unmarshal into a temporary structure to get the lspApiCommand
	type TempParams struct {
		LspApiCommand LspApiCommand   `json:"lspApiCommand"`
		Arguments     json.RawMessage `json:"args"`
	}

	var temp TempParams
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal HandleCustomTsServerCommandParams: %w", err)
	}

	// Set the LspApiCommand
	p.LspApiCommand = temp.LspApiCommand

	// Based on LspApiCommand, unmarshal args into the appropriate type
	var args interface{}
	switch temp.LspApiCommand {
	case CommandGetElementType:
		var typedArgs GetElementTypeArguments
		if err := json.Unmarshal(temp.Arguments, &typedArgs); err != nil {
			return fmt.Errorf("failed to unmarshal GetElementTypeArguments: %w", err)
		}
		args = &typedArgs

	case CommandGetSymbolType:
		var typedArgs GetSymbolTypeArguments
		if err := json.Unmarshal(temp.Arguments, &typedArgs); err != nil {
			return fmt.Errorf("failed to unmarshal GetSymbolTypeArguments: %w", err)
		}
		args = &typedArgs

	case CommandGetTypeProperties:
		var typedArgs GetTypePropertiesArguments
		if err := json.Unmarshal(temp.Arguments, &typedArgs); err != nil {
			return fmt.Errorf("failed to unmarshal GetTypePropertiesArguments: %w", err)
		}
		args = &typedArgs

	case CommandGetTypeProperty:
		var typedArgs GetTypePropertyArguments
		if err := json.Unmarshal(temp.Arguments, &typedArgs); err != nil {
			return fmt.Errorf("failed to unmarshal GetTypePropertyArguments: %w", err)
		}
		args = &typedArgs

	case CommandAreTypesMutuallyAssignable:
		var typedArgs AreTypesMutuallyAssignableArguments
		if err := json.Unmarshal(temp.Arguments, &typedArgs); err != nil {
			return fmt.Errorf("failed to unmarshal AreTypesMutuallyAssignableArguments: %w", err)
		}
		args = &typedArgs

	case CommandGetResolvedSignature:
		var typedArgs GetResolvedSignatureArguments
		if err := json.Unmarshal(temp.Arguments, &typedArgs); err != nil {
			return fmt.Errorf("failed to unmarshal GetResolvedSignatureArguments: %w", err)
		}
		args = &typedArgs

	default:
		return fmt.Errorf("unknown LspApiCommand: %s", temp.LspApiCommand)
	}

	p.Arguments = args
	return nil
}
