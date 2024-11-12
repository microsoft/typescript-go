package ast

// Ids

type NodeId uint32
type SymbolId uint32
type MergeId uint32

// Symbol

type Symbol struct {
	Flags                        SymbolFlags
	CheckFlags                   CheckFlags // Non-zero only in transient symbols created by Checker
	ConstEnumOnlyModule          bool       // True if module contains only const enums or other modules with only const enums
	IsReplaceableByMethod        bool
	Name                         string
	Declarations                 []*any
	ValueDeclaration             *any
	Members                      SymbolTable
	Exports                      SymbolTable
	Id                           SymbolId
	MergeId                      MergeId // Assigned once symbol is merged somewhere
	Parent                       *Symbol
	ExportSymbol                 *Symbol
	AssignmentDeclarationMembers map[NodeId]*any // Set of detected assignment declarations
	GlobalExports                SymbolTable      // Conditional global UMD exports
}

// SymbolTable

type SymbolTable map[string]*Symbol