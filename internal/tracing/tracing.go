package tracing

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// Tracer is an interface for recording types during type checking.
// Each checker should have its own Tracer instance to avoid sharing types between checkers.
type Tracer interface {
	// RecordType records a type for later dumping.
	RecordType(t TracedType)
	// DumpTypes writes all recorded types to disk.
	DumpTypes() error
}

// TracedType is an interface that represents a type that can be traced.
// This allows the tracing package to work with types from the checker package
// without creating a circular dependency.
type TracedType interface {
	Id() uint32
	FormatFlags() []string
	IsConditional() bool
	Symbol() *ast.Symbol
	AliasSymbol() *ast.Symbol
	AliasTypeArguments() []TracedType

	// Type-specific data accessors
	IntrinsicName() string
	UnionTypes() []TracedType
	IntersectionTypes() []TracedType
	IndexType() TracedType
	IndexedAccessObjectType() TracedType
	IndexedAccessIndexType() TracedType
	ConditionalCheckType() TracedType
	ConditionalExtendsType() TracedType
	ConditionalTrueType() TracedType
	ConditionalFalseType() TracedType
	SubstitutionBaseType() TracedType
	SubstitutionConstraintType() TracedType
	ReferenceTarget() TracedType
	ReferenceTypeArguments() []TracedType
	ReferenceNode() *ast.Node
	ReverseMappedSourceType() TracedType
	ReverseMappedMappedType() TracedType
	ReverseMappedConstraintType() TracedType
	EvolvingArrayElementType() TracedType
	EvolvingArrayFinalType() TracedType
	IsTuple() bool
	Pattern() *ast.Node
	RecursionIdentity() any

	// Display is an optional string representation of the type
	Display() string
}

// TraceRecord represents metadata about a single trace file
type TraceRecord struct {
	ConfigFilePath string `json:"configFilePath,omitzero"`
	TracePath      string `json:"tracePath,omitzero"`
	TypesPath      string `json:"typesPath,omitzero"`
}

// Tracing manages the overall tracing session including all checkers
type Tracing struct {
	fs             vfs.FS
	traceDir       string
	configFilePath string
	legend         []TraceRecord
	tracers        []*typeTracer
	mu             sync.Mutex
}

// StartTracing creates a new tracing session
func StartTracing(fs vfs.FS, traceDir string, configFilePath string) (*Tracing, error) {
	return &Tracing{
		fs:             fs,
		traceDir:       traceDir,
		configFilePath: configFilePath,
		legend:         []TraceRecord{},
		tracers:        []*typeTracer{},
	}, nil
}

// NewTypeTracer creates a new tracer for a specific checker.
// The checkerIndex is used to create unique filenames for each checker's output.
func (tr *Tracing) NewTypeTracer(checkerIndex int) Tracer {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	typesPath := tspath.CombinePaths(tr.traceDir, fmt.Sprintf("types_%d.json", checkerIndex))
	tracer := &typeTracer{
		fs:           tr.fs,
		checkerIndex: checkerIndex,
		typesPath:    typesPath,
		types:        []TracedType{},
	}
	tr.tracers = append(tr.tracers, tracer)
	tr.legend = append(tr.legend, TraceRecord{
		ConfigFilePath: tr.configFilePath,
		TypesPath:      typesPath,
	})
	return tracer
}

// StopTracing finalizes the tracing session and writes all output files
func (tr *Tracing) StopTracing() error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	// Dump types from all tracers
	for _, tracer := range tr.tracers {
		if err := tracer.DumpTypes(); err != nil {
			return fmt.Errorf("failed to dump types for checker %d: %w", tracer.checkerIndex, err)
		}
	}

	// Sort legend entries by typesPath for deterministic output
	slices.SortFunc(tr.legend, func(a, b TraceRecord) int {
		return strings.Compare(a.TypesPath, b.TypesPath)
	})

	// Write the legend file
	legendPath := tspath.CombinePaths(tr.traceDir, "legend.json")
	legendData, err := json.Marshal(tr.legend)
	if err != nil {
		return fmt.Errorf("failed to marshal legend file: %w", err)
	}
	if err := tr.fs.WriteFile(legendPath, string(legendData), false); err != nil {
		return fmt.Errorf("failed to write legend file: %w", err)
	}

	return nil
}

// typeTracer is the per-checker tracer implementation
type typeTracer struct {
	fs           vfs.FS
	checkerIndex int
	typesPath    string
	types        []TracedType
	mu           sync.Mutex
}

func (t *typeTracer) RecordType(typ TracedType) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.types = append(t.types, typ)
}

func (t *typeTracer) DumpTypes() error {
	// Copy the types slice under lock, then release so Display() calls during
	// buildTypeDescriptor don't deadlock when they create new types
	t.mu.Lock()
	types := make([]TracedType, len(t.types))
	copy(types, t.types)
	t.mu.Unlock()

	if len(types) == 0 {
		return nil
	}

	var sb strings.Builder
	// Write opening bracket (no newline so type ID matches line number)
	sb.WriteString("[")

	recursionIdentityMap := make(map[any]int)

	for i, typ := range types {
		descriptor := t.buildTypeDescriptor(typ, recursionIdentityMap)

		data, err := json.Marshal(descriptor)
		if err != nil {
			return fmt.Errorf("failed to marshal type %d: %w", typ.Id(), err)
		}

		sb.Write(data)

		if i < len(types)-1 {
			sb.WriteString(",\n")
		}
	}

	sb.WriteString("]\n")

	return t.fs.WriteFile(t.typesPath, sb.String(), false)
}

// TypeDescriptor represents a type in the output JSON
type TypeDescriptor struct {
	ID                      uint32   `json:"id"`
	IntrinsicName           string   `json:"intrinsicName,omitzero"`
	SymbolName              string   `json:"symbolName,omitzero"`
	RecursionID             *int     `json:"recursionId,omitzero"`
	IsTuple                 bool     `json:"isTuple,omitzero"`
	UnionTypes              []uint32 `json:"unionTypes,omitzero"`
	IntersectionTypes       []uint32 `json:"intersectionTypes,omitzero"`
	AliasTypeArguments      []uint32 `json:"aliasTypeArguments,omitzero"`
	KeyofType               *uint32  `json:"keyofType,omitzero"`
	IndexedAccessObjectType *uint32  `json:"indexedAccessObjectType,omitzero"`
	IndexedAccessIndexType  *uint32  `json:"indexedAccessIndexType,omitzero"`
	ConditionalCheckType    *uint32  `json:"conditionalCheckType,omitzero"`
	ConditionalExtendsType  *uint32  `json:"conditionalExtendsType,omitzero"`
	// ConditionalTrueType and ConditionalFalseType are *int32 (not *uint32) because
	// unresolved conditional branches are serialized as -1, matching TypeScript's behavior.
	ConditionalTrueType         *int32    `json:"conditionalTrueType,omitzero"`
	ConditionalFalseType        *int32    `json:"conditionalFalseType,omitzero"`
	SubstitutionBaseType        *uint32   `json:"substitutionBaseType,omitzero"`
	ConstraintType              *uint32   `json:"constraintType,omitzero"`
	InstantiatedType            *uint32   `json:"instantiatedType,omitzero"`
	TypeArguments               []uint32  `json:"typeArguments,omitzero"`
	ReferenceLocation           *Location `json:"referenceLocation,omitzero"`
	ReverseMappedSourceType     *uint32   `json:"reverseMappedSourceType,omitzero"`
	ReverseMappedMappedType     *uint32   `json:"reverseMappedMappedType,omitzero"`
	ReverseMappedConstraintType *uint32   `json:"reverseMappedConstraintType,omitzero"`
	EvolvingArrayElementType    *uint32   `json:"evolvingArrayElementType,omitzero"`
	EvolvingArrayFinalType      *uint32   `json:"evolvingArrayFinalType,omitzero"`
	DestructuringPattern        *Location `json:"destructuringPattern,omitzero"`
	FirstDeclaration            *Location `json:"firstDeclaration,omitzero"`
	Flags                       []string  `json:"flags"`
	Display                     string    `json:"display,omitzero"`
}

// Location represents a source code location
type Location struct {
	Path  string       `json:"path"`
	Start *LineAndChar `json:"start,omitzero"`
	End   *LineAndChar `json:"end,omitzero"`
}

// LineAndChar represents a line and character position (1-indexed)
type LineAndChar struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// escapeSymbolName converts internal symbol names to a JSON-safe format.
// The Go codebase uses \xFE as a prefix for internal symbol names (which is invalid UTF-8),
// while TypeScript uses "___" (three underscores). This function converts the Go format
// to "__" (two underscores) for JSON output, matching TypeScript's unescapeLeadingUnderscores behavior.
func escapeSymbolName(name string) string {
	if strings.HasPrefix(name, ast.InternalSymbolNamePrefix) {
		return "__" + name[len(ast.InternalSymbolNamePrefix):]
	}
	return name
}

func (t *typeTracer) buildTypeDescriptor(typ TracedType, recursionIdentityMap map[any]int) TypeDescriptor {
	symbol := typ.Symbol()
	aliasSymbol := typ.AliasSymbol()

	desc := TypeDescriptor{
		ID:    typ.Id(),
		Flags: typ.FormatFlags(),
	}

	// Assign a unique integer token per recursion identity, matching TypeScript's behavior.
	// This lets trace analysis tools detect which types share the same recursion identity.
	if identity := typ.RecursionIdentity(); identity != nil {
		token, ok := recursionIdentityMap[identity]
		if !ok {
			token = len(recursionIdentityMap)
			recursionIdentityMap[identity] = token
		}
		desc.RecursionID = &token
	}

	// Intrinsic name
	if name := typ.IntrinsicName(); name != "" {
		desc.IntrinsicName = name
	}

	// Symbol name - escape the internal symbol name prefix for valid JSON
	if sym := aliasSymbol; sym != nil {
		desc.SymbolName = escapeSymbolName(sym.Name)
	} else if symbol != nil {
		desc.SymbolName = escapeSymbolName(symbol.Name)
	}

	// Tuple flag
	if typ.IsTuple() {
		desc.IsTuple = true
	}

	// Union types
	if types := typ.UnionTypes(); len(types) > 0 {
		desc.UnionTypes = mapTypeIds(types)
	}

	// Intersection types
	if types := typ.IntersectionTypes(); len(types) > 0 {
		desc.IntersectionTypes = mapTypeIds(types)
	}

	// Alias type arguments
	if args := typ.AliasTypeArguments(); len(args) > 0 {
		desc.AliasTypeArguments = mapTypeIds(args)
	}

	// Index type (keyof)
	if indexType := typ.IndexType(); indexType != nil {
		id := indexType.Id()
		desc.KeyofType = &id
	}

	// Indexed access type
	if objType := typ.IndexedAccessObjectType(); objType != nil {
		id := objType.Id()
		desc.IndexedAccessObjectType = &id
	}
	if idxType := typ.IndexedAccessIndexType(); idxType != nil {
		id := idxType.Id()
		desc.IndexedAccessIndexType = &id
	}

	// Conditional type
	if typ.IsConditional() {
		if checkType := typ.ConditionalCheckType(); checkType != nil {
			id := checkType.Id()
			desc.ConditionalCheckType = &id
		}
		if extendsType := typ.ConditionalExtendsType(); extendsType != nil {
			id := extendsType.Id()
			desc.ConditionalExtendsType = &id
		}
		if trueType := typ.ConditionalTrueType(); trueType != nil {
			id := int32(trueType.Id())
			desc.ConditionalTrueType = &id
		} else {
			id := int32(-1)
			desc.ConditionalTrueType = &id
		}
		if falseType := typ.ConditionalFalseType(); falseType != nil {
			id := int32(falseType.Id())
			desc.ConditionalFalseType = &id
		} else {
			id := int32(-1)
			desc.ConditionalFalseType = &id
		}
	}

	// Substitution type
	if baseType := typ.SubstitutionBaseType(); baseType != nil {
		id := baseType.Id()
		desc.SubstitutionBaseType = &id
	}
	if constraint := typ.SubstitutionConstraintType(); constraint != nil {
		id := constraint.Id()
		desc.ConstraintType = &id
	}

	// Reference type
	if target := typ.ReferenceTarget(); target != nil {
		id := target.Id()
		desc.InstantiatedType = &id
	}
	if args := typ.ReferenceTypeArguments(); len(args) > 0 {
		desc.TypeArguments = mapTypeIds(args)
	}
	if node := typ.ReferenceNode(); node != nil {
		desc.ReferenceLocation = getLocation(node)
	}

	// Reverse mapped type
	if sourceType := typ.ReverseMappedSourceType(); sourceType != nil {
		id := sourceType.Id()
		desc.ReverseMappedSourceType = &id
	}
	if mappedType := typ.ReverseMappedMappedType(); mappedType != nil {
		id := mappedType.Id()
		desc.ReverseMappedMappedType = &id
	}
	if constraintType := typ.ReverseMappedConstraintType(); constraintType != nil {
		id := constraintType.Id()
		desc.ReverseMappedConstraintType = &id
	}

	// Evolving array type
	if elemType := typ.EvolvingArrayElementType(); elemType != nil {
		id := elemType.Id()
		desc.EvolvingArrayElementType = &id
	}
	if finalType := typ.EvolvingArrayFinalType(); finalType != nil {
		id := finalType.Id()
		desc.EvolvingArrayFinalType = &id
	}

	// Pattern (destructuring)
	if pattern := typ.Pattern(); pattern != nil {
		desc.DestructuringPattern = getLocation(pattern)
	}

	// First declaration
	if symbol != nil && len(symbol.Declarations) > 0 {
		desc.FirstDeclaration = getLocation(symbol.Declarations[0])
	}

	// Display text
	if display := typ.Display(); display != "" {
		desc.Display = display
	}

	return desc
}

func mapTypeIds(types []TracedType) []uint32 {
	if len(types) == 0 {
		return nil
	}
	ids := make([]uint32, len(types))
	for i, t := range types {
		if t != nil {
			ids[i] = t.Id()
		}
	}
	return ids
}

func getLocation(node *ast.Node) *Location {
	if node == nil {
		return nil
	}
	file := ast.GetSourceFileOfNode(node)
	if file == nil {
		return nil
	}

	startLine, startChar := scanner.GetECMALineAndCharacterOfPosition(file, node.Pos())
	endLine, endChar := scanner.GetECMALineAndCharacterOfPosition(file, node.End())

	return &Location{
		Path: string(tspath.ToPath(file.FileName(), "", false)),
		Start: &LineAndChar{
			Line:      startLine + 1,
			Character: startChar + 1,
		},
		End: &LineAndChar{
			Line:      endLine + 1,
			Character: endChar + 1,
		},
	}
}
