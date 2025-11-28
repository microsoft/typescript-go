package ast

import (
	"iter"
	"maps"
	"strings"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// Symbol

type Symbol struct {
	Flags                        SymbolFlags
	CheckFlags                   CheckFlags // Non-zero only in transient symbols created by Checker
	Name                         string
	Declarations                 []*Node
	ValueDeclaration             *Node
	Members                      SymbolTable
	Exports                      SymbolTable
	id                           atomic.Uint64
	Parent                       *Symbol
	ExportSymbol                 *Symbol
	AssignmentDeclarationMembers collections.Set[*Node] // Set of detected assignment declarations
	GlobalExports                SymbolTable            // Conditional global UMD exports
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

// SymbolTable

// type SymbolTable map[string]*Symbol

type SymbolTable interface {
	Get(name string) *Symbol
	Get2(name string) (*Symbol, bool)
	Set(name string, symbol *Symbol)
	Delete(name string)
	Keys() iter.Seq[string]
	Values() iter.Seq[*Symbol]
	Each(func(name string, symbol *Symbol))
	Iter() iter.Seq2[string, *Symbol]
	Len() int
	Clone() SymbolTable
	Find(predicate func(*Symbol) bool) *Symbol
}

type SymbolMap struct {
	m map[string]*Symbol
}

func (m *SymbolMap) Find(predicate func(*Symbol) bool) *Symbol {
	for _, symbol := range m.m {
		if predicate(symbol) {
			return symbol
		}
	}
	return nil
}

func (m *SymbolMap) Clone() SymbolTable {
	return &SymbolMap{m: maps.Clone(m.m)}
}

func (m *SymbolMap) Len() int {
	return len(m.m)
}

func (m *SymbolMap) Iter() iter.Seq2[string, *Symbol] {
	return func(yield func(string, *Symbol) bool) {
		for name, symbol := range m.m {
			if !yield(name, symbol) {
				return
			}
		}
	}
}

func (m *SymbolMap) Get(name string) *Symbol {
	return m.m[name]
}

func (m *SymbolMap) Get2(name string) (*Symbol, bool) {
	symbol, ok := m.m[name]
	return symbol, ok
}

func (m *SymbolMap) Set(name string, symbol *Symbol) {
	m.m[name] = symbol
}

func (m *SymbolMap) Delete(name string) {
	delete(m.m, name)
}

func (m *SymbolMap) Keys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for name := range m.m {
			if !yield(name) {
				return
			}
		}
	}
}

func (m *SymbolMap) Values() iter.Seq[*Symbol] {
	return func(yield func(*Symbol) bool) {
		for _, symbol := range m.m {
			if !yield(symbol) {
				return
			}
		}
	}
}

func (m *SymbolMap) Each(fn func(name string, symbol *Symbol)) {
	for name, symbol := range m.m {
		fn(name, symbol)
	}
}

func NewSymbolTable() SymbolTable {
	return &SymbolMap{m: make(map[string]*Symbol)}
}

func NewSymbolTableWithCapacity(capacity int) SymbolTable {
	return &SymbolMap{m: make(map[string]*Symbol, capacity)}
}

func NewSymbolTableFromMap(m map[string]*Symbol) SymbolTable {
	return &SymbolMap{m: m}
}

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

type CombinedSymbolTable struct {
	firstTable  SymbolTable
	secondTable SymbolTable
}

// Clone implements SymbolTable.
func (c *CombinedSymbolTable) Clone() SymbolTable {
	return &CombinedSymbolTable{
		firstTable:  c.firstTable.Clone(),
		secondTable: c.secondTable.Clone(),
	}
}

// Delete implements SymbolTable.
func (c *CombinedSymbolTable) Delete(name string) {
	if c.firstTable.Get(name) != nil {
		c.firstTable.Delete(name)
	} else {
		c.secondTable.Delete(name)
	}
}

// Each implements SymbolTable.
func (c *CombinedSymbolTable) Each(fn func(name string, symbol *Symbol)) {
	c.firstTable.Each(func(name string, symbol *Symbol) {
		fn(name, symbol)
	})
	c.secondTable.Each(func(name string, symbol *Symbol) {
		fn(name, symbol)
	})
}

// Find implements SymbolTable.
func (c *CombinedSymbolTable) Find(predicate func(*Symbol) bool) *Symbol {
	ret := c.firstTable.Find(predicate)
	if ret != nil {
		return ret
	}
	return c.secondTable.Find(predicate)
}

// Get implements SymbolTable.
func (c *CombinedSymbolTable) Get(name string) *Symbol {
	ret := c.firstTable.Get(name)
	if ret != nil {
		return ret
	}
	return c.secondTable.Get(name)
}

// Get2 implements SymbolTable.
func (c *CombinedSymbolTable) Get2(name string) (*Symbol, bool) {
	if value, ok := c.firstTable.Get2(name); ok {
		return value, ok
	}
	return c.secondTable.Get2(name)
}

// Iter implements SymbolTable.
func (c *CombinedSymbolTable) Iter() iter.Seq2[string, *Symbol] {
	seen := make(map[string]struct{})
	return func(yield func(string, *Symbol) bool) {
		for name, symbol := range c.firstTable.Iter() {
			if _, ok := seen[name]; !ok {
				seen[name] = struct{}{}
				if !yield(name, symbol) {
					break
				}
			}
		}
		for name, symbol := range c.secondTable.Iter() {
			if _, ok := seen[name]; !ok {
				seen[name] = struct{}{}
				if !yield(name, symbol) {
					return
				}
			}
		}
	}
}

// Keys implements SymbolTable.
func (c *CombinedSymbolTable) Keys() iter.Seq[string] {
	return func(yield func(string) bool) {
		seen := make(map[string]struct{})
		for name := range c.firstTable.Keys() {
			if _, ok := seen[name]; !ok {
				seen[name] = struct{}{}
				if !yield(name) {
					break
				}
			}
		}

		for name := range c.secondTable.Keys() {
			if _, ok := seen[name]; !ok {
				seen[name] = struct{}{}
				if !yield(name) {
					return
				}
			}
		}
	}
}

// Len implements SymbolTable.
func (c *CombinedSymbolTable) Len() int {
	len := 0
	for k := range c.Iter() {
		_ = k
		len++
	}
	return len
}

// Set implements SymbolTable.
func (c *CombinedSymbolTable) Set(name string, symbol *Symbol) {
	c.firstTable.Set(name, symbol)
}

// Values implements SymbolTable.
func (c *CombinedSymbolTable) Values() iter.Seq[*Symbol] {
	return func(yield func(*Symbol) bool) {
		c.Iter()(func(name string, symbol *Symbol) bool {
			return yield(symbol)
		})
	}
}

var _ SymbolTable = (*CombinedSymbolTable)(nil)

type DenoForkContextInfo struct {
	TypesNodeIgnorableNames *collections.Set[string]
	NodeOnlyGlobalNames     *collections.Set[string]
}

type DenoForkContext struct {
	globals          SymbolTable
	nodeGlobals      SymbolTable
	combinedGlobals  SymbolTable
	mergeSymbol      func(target *Symbol, source *Symbol, unidirectional bool) *Symbol
	getMergedSymbol  func(source *Symbol) *Symbol
	isNodeSourceFile func(path tspath.Path) bool
	info             DenoForkContextInfo
}

func NewDenoForkContext(
	globals SymbolTable,
	nodeGlobals SymbolTable,
	mergeSymbol func(target *Symbol, source *Symbol, unidirectional bool) *Symbol,
	getMergedSymbol func(source *Symbol) *Symbol,
	isNodeSourceFile func(path tspath.Path) bool,
	info DenoForkContextInfo,
) *DenoForkContext {
	return &DenoForkContext{
		globals:     globals,
		nodeGlobals: nodeGlobals,
		combinedGlobals: &CombinedSymbolTable{
			firstTable:  nodeGlobals,
			secondTable: globals,
		},
		mergeSymbol:      mergeSymbol,
		getMergedSymbol:  getMergedSymbol,
		isNodeSourceFile: isNodeSourceFile,
		info:             info,
	}
}

func (c *DenoForkContext) GetGlobalsForName(name string) SymbolTable {
	if c.info.NodeOnlyGlobalNames.Has(name) {
		return c.nodeGlobals
	} else {
		return c.globals
	}
}

func isTypesNodePkgPath(path tspath.Path) bool {
	return strings.HasSuffix(string(path), ".d.ts") && strings.Contains(string(path), "/@types/node/")
}

func symbolHasAnyTypesNodePkgDecl(symbol *Symbol, hasNodeSourceFile func(*Node) bool) bool {
	if symbol == nil || symbol.Declarations == nil {
		return false
	}
	for _, decl := range symbol.Declarations {
		sourceFile := GetSourceFileOfNode(decl)
		if sourceFile != nil && hasNodeSourceFile(decl) && isTypesNodePkgPath(sourceFile.Path()) {
			return true
		}
	}
	return false
}

func (c *DenoForkContext) MergeGlobalSymbolTable(node *Node, source SymbolTable, unidirectional bool) {
	sourceFile := GetSourceFileOfNode(node)
	isNodeFile := c.HasNodeSourceFile(node)
	isTypesNodeSourceFile := isNodeFile && isTypesNodePkgPath(sourceFile.Path())

	for id, sourceSymbol := range source.Iter() {
		var target SymbolTable
		if isNodeFile {
			target = c.GetGlobalsForName(id)
		} else {
			target = c.globals
		}
		targetSymbol := target.Get(id)
		if isTypesNodeSourceFile {
		}
		if isTypesNodeSourceFile && targetSymbol != nil && c.info.TypesNodeIgnorableNames.Has(id) && !symbolHasAnyTypesNodePkgDecl(targetSymbol, c.HasNodeSourceFile) {
			continue
		}
		var merged *Symbol
		if targetSymbol != nil {
			merged = c.mergeSymbol(targetSymbol, sourceSymbol, unidirectional)
		} else {
			merged = c.getMergedSymbol(sourceSymbol)
		}

		target.Set(id, merged)
	}
}

func (c *DenoForkContext) CombinedGlobals() SymbolTable {
	return c.combinedGlobals
}

func (c *DenoForkContext) HasNodeSourceFile(node *Node) bool {
	if node == nil || c == nil || c.isNodeSourceFile == nil {
		return false
	}
	sourceFile := GetSourceFileOfNode(node)
	if sourceFile == nil {
		return false
	}
	if c.isNodeSourceFile(sourceFile.Path()) {
		return true
	}
	return false
}
