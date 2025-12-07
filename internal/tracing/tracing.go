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
	// Close releases any resources held by the tracer.
	Close() error
}

// TracedType is an interface that represents a type that can be traced.
// This allows the tracing package to work with types from the checker package
// without creating a circular dependency.
type TracedType interface {
	Id() uint32
	Flags() uint32
	ObjectFlags() uint32
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

	// Display is an optional string representation of the type
	Display() string
}

// TraceRecord represents metadata about a single trace file
type TraceRecord struct {
	ConfigFilePath string `json:"configFilePath,omitempty"`
	TracePath      string `json:"tracePath,omitempty"`
	TypesPath      string `json:"typesPath,omitempty"`
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

func (t *typeTracer) Close() error {
	return nil
}

// TypeDescriptor represents a type in the output JSON
type TypeDescriptor struct {
	ID                          uint32    `json:"id"`
	IntrinsicName               string    `json:"intrinsicName,omitempty"`
	SymbolName                  string    `json:"symbolName,omitempty"`
	RecursionID                 *int      `json:"recursionId,omitempty"`
	IsTuple                     bool      `json:"isTuple,omitempty"`
	UnionTypes                  []uint32  `json:"unionTypes,omitempty"`
	IntersectionTypes           []uint32  `json:"intersectionTypes,omitempty"`
	AliasTypeArguments          []uint32  `json:"aliasTypeArguments,omitempty"`
	KeyofType                   *uint32   `json:"keyofType,omitempty"`
	IndexedAccessObjectType     *uint32   `json:"indexedAccessObjectType,omitempty"`
	IndexedAccessIndexType      *uint32   `json:"indexedAccessIndexType,omitempty"`
	ConditionalCheckType        *uint32   `json:"conditionalCheckType,omitempty"`
	ConditionalExtendsType      *uint32   `json:"conditionalExtendsType,omitempty"`
	ConditionalTrueType         *int32    `json:"conditionalTrueType,omitempty"`
	ConditionalFalseType        *int32    `json:"conditionalFalseType,omitempty"`
	SubstitutionBaseType        *uint32   `json:"substitutionBaseType,omitempty"`
	ConstraintType              *uint32   `json:"constraintType,omitempty"`
	InstantiatedType            *uint32   `json:"instantiatedType,omitempty"`
	TypeArguments               []uint32  `json:"typeArguments,omitempty"`
	ReferenceLocation           *Location `json:"referenceLocation,omitempty"`
	ReverseMappedSourceType     *uint32   `json:"reverseMappedSourceType,omitempty"`
	ReverseMappedMappedType     *uint32   `json:"reverseMappedMappedType,omitempty"`
	ReverseMappedConstraintType *uint32   `json:"reverseMappedConstraintType,omitempty"`
	EvolvingArrayElementType    *uint32   `json:"evolvingArrayElementType,omitempty"`
	EvolvingArrayFinalType      *uint32   `json:"evolvingArrayFinalType,omitempty"`
	DestructuringPattern        *Location `json:"destructuringPattern,omitempty"`
	FirstDeclaration            *Location `json:"firstDeclaration,omitempty"`
	Flags                       []string  `json:"flags"`
	Display                     string    `json:"display,omitempty"`
}

// Location represents a source code location
type Location struct {
	Path  string       `json:"path"`
	Start *LineAndChar `json:"start,omitempty"`
	End   *LineAndChar `json:"end,omitempty"`
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
	flags := typ.Flags()
	objectFlags := typ.ObjectFlags()
	symbol := typ.Symbol()
	aliasSymbol := typ.AliasSymbol()

	desc := TypeDescriptor{
		ID:    typ.Id(),
		Flags: formatTypeFlags(flags),
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
	if flags&TypeFlagsUnion != 0 {
		desc.UnionTypes = mapTypeIds(typ.UnionTypes())
	}

	// Intersection types
	if flags&TypeFlagsIntersection != 0 {
		desc.IntersectionTypes = mapTypeIds(typ.IntersectionTypes())
	}

	// Alias type arguments
	if args := typ.AliasTypeArguments(); len(args) > 0 {
		desc.AliasTypeArguments = mapTypeIds(args)
	}

	// Index type (keyof)
	if flags&TypeFlagsIndex != 0 {
		if indexType := typ.IndexType(); indexType != nil {
			id := indexType.Id()
			desc.KeyofType = &id
		}
	}

	// Indexed access type
	if flags&TypeFlagsIndexedAccess != 0 {
		if objType := typ.IndexedAccessObjectType(); objType != nil {
			id := objType.Id()
			desc.IndexedAccessObjectType = &id
		}
		if idxType := typ.IndexedAccessIndexType(); idxType != nil {
			id := idxType.Id()
			desc.IndexedAccessIndexType = &id
		}
	}

	// Conditional type
	if flags&TypeFlagsConditional != 0 {
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
	if flags&TypeFlagsSubstitution != 0 {
		if baseType := typ.SubstitutionBaseType(); baseType != nil {
			id := baseType.Id()
			desc.SubstitutionBaseType = &id
		}
		if constraint := typ.SubstitutionConstraintType(); constraint != nil {
			id := constraint.Id()
			desc.ConstraintType = &id
		}
	}

	// Reference type (Object with Reference flag)
	if flags&TypeFlagsObject != 0 && objectFlags&ObjectFlagsReference != 0 {
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
	}

	// Reverse mapped type
	if flags&TypeFlagsObject != 0 && objectFlags&ObjectFlagsReverseMapped != 0 {
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
	}

	// Evolving array type
	if flags&TypeFlagsObject != 0 && objectFlags&ObjectFlagsEvolvingArray != 0 {
		if elemType := typ.EvolvingArrayElementType(); elemType != nil {
			id := elemType.Id()
			desc.EvolvingArrayElementType = &id
		}
		if finalType := typ.EvolvingArrayFinalType(); finalType != nil {
			id := finalType.Id()
			desc.EvolvingArrayFinalType = &id
		}
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

// TypeFlags constants (copied from checker to avoid circular dependency)
const (
	TypeFlagsNone            uint32 = 0
	TypeFlagsAny             uint32 = 1 << 0
	TypeFlagsUnknown         uint32 = 1 << 1
	TypeFlagsUndefined       uint32 = 1 << 2
	TypeFlagsNull            uint32 = 1 << 3
	TypeFlagsVoid            uint32 = 1 << 4
	TypeFlagsString          uint32 = 1 << 5
	TypeFlagsNumber          uint32 = 1 << 6
	TypeFlagsBigInt          uint32 = 1 << 7
	TypeFlagsBoolean         uint32 = 1 << 8
	TypeFlagsESSymbol        uint32 = 1 << 9
	TypeFlagsStringLiteral   uint32 = 1 << 10
	TypeFlagsNumberLiteral   uint32 = 1 << 11
	TypeFlagsBigIntLiteral   uint32 = 1 << 12
	TypeFlagsBooleanLiteral  uint32 = 1 << 13
	TypeFlagsUniqueESSymbol  uint32 = 1 << 14
	TypeFlagsEnumLiteral     uint32 = 1 << 15
	TypeFlagsEnum            uint32 = 1 << 16
	TypeFlagsNonPrimitive    uint32 = 1 << 17
	TypeFlagsNever           uint32 = 1 << 18
	TypeFlagsTypeParameter   uint32 = 1 << 19
	TypeFlagsObject          uint32 = 1 << 20
	TypeFlagsIndex           uint32 = 1 << 21
	TypeFlagsTemplateLiteral uint32 = 1 << 22
	TypeFlagsStringMapping   uint32 = 1 << 23
	TypeFlagsSubstitution    uint32 = 1 << 24
	TypeFlagsIndexedAccess   uint32 = 1 << 25
	TypeFlagsConditional     uint32 = 1 << 26
	TypeFlagsUnion           uint32 = 1 << 27
	TypeFlagsIntersection    uint32 = 1 << 28

	TypeFlagsLiteral = TypeFlagsStringLiteral | TypeFlagsNumberLiteral | TypeFlagsBigIntLiteral | TypeFlagsBooleanLiteral
)

// ObjectFlags constants
const (
	ObjectFlagsNone          uint32 = 0
	ObjectFlagsClass         uint32 = 1 << 0
	ObjectFlagsInterface     uint32 = 1 << 1
	ObjectFlagsReference     uint32 = 1 << 2
	ObjectFlagsTuple         uint32 = 1 << 3
	ObjectFlagsAnonymous     uint32 = 1 << 4
	ObjectFlagsMapped        uint32 = 1 << 5
	ObjectFlagsInstantiated  uint32 = 1 << 6
	ObjectFlagsEvolvingArray uint32 = 1 << 8
	ObjectFlagsReverseMapped uint32 = 1 << 10
)

func formatTypeFlags(flags uint32) []string {
	var result []string

	flagNames := []struct {
		flag uint32
		name string
	}{
		{TypeFlagsAny, "Any"},
		{TypeFlagsUnknown, "Unknown"},
		{TypeFlagsUndefined, "Undefined"},
		{TypeFlagsNull, "Null"},
		{TypeFlagsVoid, "Void"},
		{TypeFlagsString, "String"},
		{TypeFlagsNumber, "Number"},
		{TypeFlagsBigInt, "BigInt"},
		{TypeFlagsBoolean, "Boolean"},
		{TypeFlagsESSymbol, "ESSymbol"},
		{TypeFlagsStringLiteral, "StringLiteral"},
		{TypeFlagsNumberLiteral, "NumberLiteral"},
		{TypeFlagsBigIntLiteral, "BigIntLiteral"},
		{TypeFlagsBooleanLiteral, "BooleanLiteral"},
		{TypeFlagsUniqueESSymbol, "UniqueESSymbol"},
		{TypeFlagsEnumLiteral, "EnumLiteral"},
		{TypeFlagsEnum, "Enum"},
		{TypeFlagsNonPrimitive, "NonPrimitive"},
		{TypeFlagsNever, "Never"},
		{TypeFlagsTypeParameter, "TypeParameter"},
		{TypeFlagsObject, "Object"},
		{TypeFlagsIndex, "Index"},
		{TypeFlagsTemplateLiteral, "TemplateLiteral"},
		{TypeFlagsStringMapping, "StringMapping"},
		{TypeFlagsSubstitution, "Substitution"},
		{TypeFlagsIndexedAccess, "IndexedAccess"},
		{TypeFlagsConditional, "Conditional"},
		{TypeFlagsUnion, "Union"},
		{TypeFlagsIntersection, "Intersection"},
	}

	for _, fn := range flagNames {
		if flags&fn.flag != 0 {
			result = append(result, fn.name)
		}
	}

	if len(result) == 0 {
		result = append(result, "None")
	}

	return result
}
