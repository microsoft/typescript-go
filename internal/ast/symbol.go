package ast

import (
	"strings"
	"sync/atomic"
)

// Symbol

type Symbol struct {
	Flags            SymbolFlags
	CheckFlags       CheckFlags // Non-zero only in transient symbols created by Checker
	Name             string
	Declarations     []*Node
	ValueDeclaration *Node
	id               atomic.Uint64
	Parent           *Symbol
	extra            *symbolExtra // Rarely populated, so boxed to keep Symbol small
}

type symbolExtra struct {
	members      SymbolTable
	exports      SymbolTable
	exportSymbol *Symbol
}

func (s *Symbol) getExtra() *symbolExtra {
	if s.extra == nil {
		s.extra = &symbolExtra{}
	}
	return s.extra
}

func (s *Symbol) Members() SymbolTable {
	if s.extra != nil {
		return s.extra.members
	}
	return nil
}

func (s *Symbol) SetMembers(members SymbolTable) {
	if members == nil && s.extra == nil {
		return
	}
	s.getExtra().members = members
}

func (s *Symbol) Exports() SymbolTable {
	if s.extra != nil {
		return s.extra.exports
	}
	return nil
}

func (s *Symbol) SetExports(exports SymbolTable) {
	if exports == nil && s.extra == nil {
		return
	}
	s.getExtra().exports = exports
}

func (s *Symbol) ExportSymbol() *Symbol {
	if s.extra != nil {
		return s.extra.exportSymbol
	}
	return nil
}

func (s *Symbol) SetExportSymbol(exportSymbol *Symbol) {
	if exportSymbol == nil && s.extra == nil {
		return
	}
	s.getExtra().exportSymbol = exportSymbol
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
	if exportSymbol := s.ExportSymbol(); exportSymbol != nil {
		return s.Flags | exportSymbol.Flags
	}
	return s.Flags
}

// SymbolTable

type SymbolTable map[string]*Symbol

const InternalSymbolNamePrefix = "\xFE" // Invalid UTF8 sequence, will never occur as IdentifierName

const (
	InternalSymbolNameCall                    = InternalSymbolNamePrefix + "call"                    // Call signatures
	InternalSymbolNameConstructor             = InternalSymbolNamePrefix + "constructor"             // Constructor implementations
	InternalSymbolNameNew                     = InternalSymbolNamePrefix + "new"                     // Constructor signatures
	InternalSymbolNameIndex                   = InternalSymbolNamePrefix + "index"                   // Index signatures
	InternalSymbolNameExportStar              = InternalSymbolNamePrefix + "export"                  // Module export * declarations
	InternalSymbolNameGlobal                  = InternalSymbolNamePrefix + "global"                  // Global self-reference
	InternalSymbolNameMissing                 = InternalSymbolNamePrefix + "missing"                 // Indicates missing symbol
	InternalSymbolNameType                    = InternalSymbolNamePrefix + "type"                    // Anonymous type literal symbol
	InternalSymbolNameObject                  = InternalSymbolNamePrefix + "object"                  // Anonymous object literal declaration
	InternalSymbolNameJSXAttributes           = InternalSymbolNamePrefix + "jsxAttributes"           // Anonymous JSX attributes object literal declaration
	InternalSymbolNameClass                   = InternalSymbolNamePrefix + "class"                   // Unnamed class expression
	InternalSymbolNameFunction                = InternalSymbolNamePrefix + "function"                // Unnamed function expression
	InternalSymbolNameComputed                = InternalSymbolNamePrefix + "computed"                // Computed property name declaration with dynamic name
	InternalSymbolNameAssignmentDeclaration   = InternalSymbolNamePrefix + "assignment"              // Assignment declarations
	InternalSymbolNameInstantiationExpression = InternalSymbolNamePrefix + "instantiationExpression" // Instantiation expressions
	InternalSymbolNameImportAttributes        = InternalSymbolNamePrefix + "importAttributes"
	InternalSymbolNameExportEquals            = "export=" // Export assignment symbol
	InternalSymbolNameDefault                 = "default" // Default export symbol (technically not wholly internal, but included here for usability)
	InternalSymbolNameThis                    = "this"
	InternalSymbolNameModuleExports           = "module.exports"
)

func SymbolName(symbol *Symbol) string {
	if symbol.ValueDeclaration != nil && IsPrivateIdentifierClassElementDeclaration(symbol.ValueDeclaration) {
		return symbol.ValueDeclaration.Name().Text()
	}
	return symbol.Name
}

// EscapeAllInternalSymbolNames replaces internal symbol name markers ("\xFE") with "__".
func EscapeAllInternalSymbolNames(name string) string {
	return strings.ReplaceAll(name, InternalSymbolNamePrefix, "__")
}

func EscapeInternalSymbolName(name string) string {
	if rest, ok := strings.CutPrefix(name, InternalSymbolNamePrefix); ok {
		return "__" + rest
	}
	return name
}

// EscapeSymbolName converts a binder symbol name into its escaped "__String"
// form. Internal names (prefixed with the "\xFE" sentinel) become "__"-prefixed,
// and user names that already begin with "__" gain an extra leading underscore
// so they can be distinguished from internal names.
func EscapeSymbolName(name string) string {
	if rest, ok := strings.CutPrefix(name, InternalSymbolNamePrefix); ok {
		return "__" + rest
	}
	if len(name) >= 2 && name[0] == '_' && name[1] == '_' {
		return "_" + name
	}
	return name
}
