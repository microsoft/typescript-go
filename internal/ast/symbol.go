package ast

import (
	"sync/atomic"
)

// Symbol

type Symbol struct {
	Flags            SymbolFlags
	CheckFlags       CheckFlags // Non-zero only in transient symbols created by Checker
	Name             SymbolNameKey
	Declarations     []*Node
	ValueDeclaration *Node
	Members          SymbolTable
	Exports          SymbolTable
	id               atomic.Uint64
	Parent           *Symbol
	ExportSymbol     *Symbol
}

func (s *Symbol) IsExternalModule() bool {
	return s.Flags&SymbolFlagsModule != 0 && len(s.Name) > 0 && s.Name[0] == '"'
}

func (s *Symbol) IsStatic() bool {
	if s.ValueDeclaration == nil {
		return false
	}
	modifierFlags := s.ValueDeclaration.ModifierFlags()
	return modifierFlags&ModifierFlagsStatic != 0
}

// See comment on `declareModuleMember` in `binder.go`.
func (s *Symbol) CombinedLocalAndExportSymbolFlags() SymbolFlags {
	if s.ExportSymbol != nil {
		return s.Flags | s.ExportSymbol.Flags
	}
	return s.Flags
}

// SymbolTable

type SymbolNameKey string

type SymbolTable map[SymbolNameKey]*Symbol

func (name SymbolNameKey) EscapedText() string {
	return string(name)
}

func InternalSymbolName(suffix string) SymbolNameKey {
	return SymbolNameKey(InternalSymbolNamePrefix + suffix)
}

const InternalSymbolNamePrefix = "__"

const (
	InternalSymbolNameCall                    SymbolNameKey = InternalSymbolNamePrefix + "call"                    // Call signatures
	InternalSymbolNameConstructor             SymbolNameKey = InternalSymbolNamePrefix + "constructor"             // Constructor implementations
	InternalSymbolNameNew                     SymbolNameKey = InternalSymbolNamePrefix + "new"                     // Constructor signatures
	InternalSymbolNameIndex                   SymbolNameKey = InternalSymbolNamePrefix + "index"                   // Index signatures
	InternalSymbolNameExportStar              SymbolNameKey = InternalSymbolNamePrefix + "export"                  // Module export * declarations
	InternalSymbolNameGlobal                  SymbolNameKey = InternalSymbolNamePrefix + "global"                  // Global self-reference
	InternalSymbolNameMissing                 SymbolNameKey = InternalSymbolNamePrefix + "missing"                 // Indicates missing symbol
	InternalSymbolNameType                    SymbolNameKey = InternalSymbolNamePrefix + "type"                    // Anonymous type literal symbol
	InternalSymbolNameObject                  SymbolNameKey = InternalSymbolNamePrefix + "object"                  // Anonymous object literal declaration
	InternalSymbolNameJSXAttributes           SymbolNameKey = InternalSymbolNamePrefix + "jsxAttributes"           // Anonymous JSX attributes object literal declaration
	InternalSymbolNameClass                   SymbolNameKey = InternalSymbolNamePrefix + "class"                   // Unnamed class expression
	InternalSymbolNameFunction                SymbolNameKey = InternalSymbolNamePrefix + "function"                // Unnamed function expression
	InternalSymbolNameComputed                SymbolNameKey = InternalSymbolNamePrefix + "computed"                // Computed property name declaration with dynamic name
	InternalSymbolNameAssignmentDeclaration   SymbolNameKey = InternalSymbolNamePrefix + "assignment"              // Assignment declarations
	InternalSymbolNameInstantiationExpression SymbolNameKey = InternalSymbolNamePrefix + "instantiationExpression" // Instantiation expressions
	InternalSymbolNameImportAttributes        SymbolNameKey = InternalSymbolNamePrefix + "importAttributes"
	InternalSymbolNameExportEquals            SymbolNameKey = "export=" // Export assignment symbol
	InternalSymbolNameDefault                 SymbolNameKey = "default" // Default export symbol (technically not wholly internal, but included here for usability)
	InternalSymbolNameThis                    SymbolNameKey = "this"
	InternalSymbolNameModuleExports           SymbolNameKey = "module.exports"
)

func SymbolName(symbol *Symbol) string {
	if symbol.ValueDeclaration != nil && IsPrivateIdentifierClassElementDeclaration(symbol.ValueDeclaration) {
		return symbol.ValueDeclaration.Name().Text()
	}
	return UnescapeLeadingUnderscores(symbol.Name)
}

func EscapeLeadingUnderscores(identifier string) SymbolNameKey {
	if len(identifier) >= 2 && identifier[0] == '_' && identifier[1] == '_' {
		return SymbolNameKey("_" + identifier)
	}
	return SymbolNameKey(identifier)
}

func UnescapeLeadingUnderscores(identifier SymbolNameKey) string {
	name := string(identifier)
	if len(name) >= 3 && name[0] == '_' && name[1] == '_' && name[2] == '_' {
		return name[1:]
	}
	return name
}
