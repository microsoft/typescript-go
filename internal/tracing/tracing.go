package tracing

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"
	"time"

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

// traceStackEntry records the state of a Push call so Pop can
// either emit a matching "E" event (separateBeginAndEnd) or a sampled "X" event.
type traceStackEntry struct {
	phase               Phase
	name                string
	argsJSON            string
	startTime           time.Time // wall-clock time for sampling decisions
	separateBeginAndEnd bool
}

// sampleInterval matches TypeScript's 10ms sampling interval.
// Events with separateBeginAndEnd=false are only recorded if their
// duration crosses a 10ms sampling boundary.
const sampleInterval = 10 * time.Millisecond

// Tracing manages the overall tracing session including all checkers
type Tracing struct {
	fs               vfs.FS
	traceDir         string
	configFilePath   string
	legend           []TraceRecord
	tracers          []*typeTracer
	checkerTracings  []*CheckerTracing
	traceContent     strings.Builder
	traceStarted     bool
	deterministic    bool   // when true, use monotonic counter instead of real time
	timestampCounter uint64 // only used in deterministic mode
	startTime        time.Time
	eventStack       []traceStackEntry
	mu               sync.Mutex
}

// Phase represents a tracing phase
type Phase string

const (
	PhaseParse      Phase = "parse"
	PhaseProgram    Phase = "program"
	PhaseBind       Phase = "bind"
	PhaseCheck      Phase = "check"
	PhaseCheckTypes Phase = "checkTypes"
	PhaseEmit       Phase = "emit"
	PhaseSession    Phase = "session"
)

// StartTracing creates a new tracing session.
// When deterministic is true, timestamps use a monotonic counter instead of
// real wall-clock time, producing stable output for test baselines.
func StartTracing(fs vfs.FS, traceDir string, configFilePath string, deterministic bool) (*Tracing, error) {
	tr := &Tracing{
		fs:             fs,
		traceDir:       traceDir,
		configFilePath: configFilePath,
		legend:         []TraceRecord{},
		tracers:        []*typeTracer{},
		traceStarted:   true,
		deterministic:  deterministic,
		startTime:      time.Now(),
	}

	// Write the trace file header with metadata events
	tr.traceContent.WriteString("[\n")

	// Write metadata events (matching TypeScript's format)
	metaTs := tr.timestamp()
	tr.writeEventRaw("M", "__metadata", metaTs, "process_name", "{\"name\":\"tsgo\"}")
	tr.traceContent.WriteString(",\n")
	tr.writeEventRaw("M", "__metadata", metaTs, "thread_name", "{\"name\":\"Main\"}")
	tr.traceContent.WriteString(",\n")
	tr.writeEventRaw("M", "disabled-by-default-devtools.timeline", metaTs, "TracingStartedInBrowser", "")

	return tr, nil
}

// timestamp returns the current timestamp in microseconds.
// In deterministic mode it returns a monotonically increasing counter;
// otherwise it returns the real elapsed wall-clock time since tracing started,
// matching TypeScript's 1000 * timestamp() (microseconds).
func (tr *Tracing) timestamp() float64 {
	if tr.deterministic {
		tr.timestampCounter++
		return float64(tr.timestampCounter)
	}
	return float64(time.Since(tr.startTime).Nanoseconds()) / 1000.0
}

// writeEventTo writes a trace event to the given builder with deterministic field ordering.
// argsJSON should be pre-formatted JSON for the args object contents, or empty string for no args.
// extras should be pre-formatted JSON key-value pair(s) like `"dur":123`, or empty string.
func writeEventTo(buf *strings.Builder, ph string, cat string, ts float64, name string, argsJSON string, extras ...string) {
	buf.WriteString("{\"pid\":1,\"tid\":1,\"ph\":\"")
	buf.WriteString(ph)
	buf.WriteString("\"")
	if cat != "" {
		buf.WriteString(",\"cat\":\"")
		buf.WriteString(cat)
		buf.WriteString("\"")
	}
	buf.WriteString(",\"ts\":")
	writeNumberTo(buf, ts)
	if name != "" {
		buf.WriteString(",\"name\":\"")
		buf.WriteString(name)
		buf.WriteString("\"")
	}
	for _, extra := range extras {
		if extra != "" {
			buf.WriteString(",")
			buf.WriteString(extra)
		}
	}
	if argsJSON != "" {
		buf.WriteString(",\"args\":")
		buf.WriteString(argsJSON)
	}
	buf.WriteString("}")
}

// writeNumberTo writes a number to the builder, using integer format when
// the value has no fractional part to produce cleaner output.
func writeNumberTo(buf *strings.Builder, v float64) {
	if v == float64(int64(v)) {
		fmt.Fprintf(buf, "%d", int64(v))
	} else {
		fmt.Fprintf(buf, "%.4f", v)
	}
}

// writeEventRaw writes a trace event to the shared trace content.
func (tr *Tracing) writeEventRaw(ph string, cat string, ts float64, name string, argsJSON string, extras ...string) {
	writeEventTo(&tr.traceContent, ph, cat, ts, name, argsJSON, extras...)
}

// writeArgsJSON constructs a deterministic JSON object from key-value pairs.
// Keys and values are written in the order provided. Values are JSON-escaped strings.
func writeArgsJSON(pairs ...string) string {
	if len(pairs) == 0 || len(pairs)%2 != 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("{")
	for i := 0; i < len(pairs); i += 2 {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("\"")
		b.WriteString(pairs[i])
		b.WriteString("\":\"")
		b.WriteString(pairs[i+1])
		b.WriteString("\"")
	}
	b.WriteString("}")
	return b.String()
}

// Instant records an instant event in the trace.
// Safe to call on nil receiver.
func (tr *Tracing) Instant(phase Phase, name string, args ...string) {
	if tr == nil || !tr.traceStarted {
		return
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	ts := tr.timestamp()
	tr.traceContent.WriteString(",\n")
	tr.writeEventRaw("I", string(phase), ts, name, writeArgsJSON(args...))
}

// pushTo implements the Push logic, writing to the given buffer and event stack.
func pushTo(buf *strings.Builder, stack *[]traceStackEntry, ts float64, phase Phase, name string, separateBeginAndEnd bool, args ...string) {
	argsJSON := writeArgsJSON(args...)
	*stack = append(*stack, traceStackEntry{
		phase:               phase,
		name:                name,
		argsJSON:            argsJSON,
		startTime:           time.Now(),
		separateBeginAndEnd: separateBeginAndEnd,
	})

	if separateBeginAndEnd {
		buf.WriteString(",\n")
		writeEventTo(buf, "B", string(phase), ts, name, argsJSON)
	}
}

// popFrom implements the Pop logic, writing to the given buffer and event stack.
func popFrom(buf *strings.Builder, stack *[]traceStackEntry, ts float64, startTime time.Time) {
	if len(*stack) == 0 {
		return
	}
	n := len(*stack)
	entry := (*stack)[n-1]
	*stack = (*stack)[:n-1]

	if entry.separateBeginAndEnd {
		buf.WriteString(",\n")
		writeEventTo(buf, "E", string(entry.phase), ts, entry.name, entry.argsJSON)
	} else {
		now := time.Now()
		startMicros := float64(entry.startTime.Sub(startTime).Nanoseconds()) / 1000.0
		endMicros := float64(now.Sub(startTime).Nanoseconds()) / 1000.0
		intervalMicros := float64(sampleInterval.Nanoseconds()) / 1000.0
		dur := endMicros - startMicros

		if intervalMicros-math.Mod(startMicros, intervalMicros) <= dur {
			buf.WriteString(",\n")
			writeEventTo(buf, "X", string(entry.phase), startMicros, entry.name, entry.argsJSON, fmt.Sprintf("\"dur\":%.4f", dur))
		}
	}
}

// Push starts a trace event block on the shared trace buffer.
// Safe to call on nil receiver.
//
// When separateBeginAndEnd is true, a "B" (begin) event is written immediately and
// Pop will write a matching "E" (end) event. This is used for events that must always
// appear in the trace (e.g. checkSourceFile, createProgram, emit).
//
// When separateBeginAndEnd is false (the default in TypeScript), the event is only
// recorded if its duration crosses a 10ms sampling boundary, matching TypeScript's
// behavior of sampling short-lived events to avoid trace bloat.
//
// args are key-value pairs: "key1", "value1", "key2", "value2", ...
func (tr *Tracing) Push(phase Phase, name string, separateBeginAndEnd bool, args ...string) {
	if tr == nil || !tr.traceStarted {
		return
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	var ts float64
	if separateBeginAndEnd {
		ts = tr.timestamp()
	}
	pushTo(&tr.traceContent, &tr.eventStack, ts, phase, name, separateBeginAndEnd, args...)
}

// Pop ends the most recent trace event block on the shared trace buffer.
// Safe to call on nil receiver.
func (tr *Tracing) Pop() {
	if tr == nil || !tr.traceStarted {
		return
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	if len(tr.eventStack) == 0 {
		return
	}
	var ts float64
	if tr.eventStack[len(tr.eventStack)-1].separateBeginAndEnd {
		ts = tr.timestamp()
	}
	popFrom(&tr.traceContent, &tr.eventStack, ts, tr.startTime)
}

// CheckerTracing manages per-checker trace events. Each checker gets its own
// CheckerTracing whose events are combined with the shared global events to
// produce checker-specific trace_N.json files.
type CheckerTracing struct {
	tr         *Tracing
	content    strings.Builder
	eventStack []traceStackEntry
	mu         sync.Mutex
}

// Push starts a trace event block on the per-checker trace buffer.
// Safe to call on nil receiver.
func (ct *CheckerTracing) Push(phase Phase, name string, separateBeginAndEnd bool, args ...string) {
	if ct == nil {
		return
	}

	ct.mu.Lock()
	defer ct.mu.Unlock()

	var ts float64
	if separateBeginAndEnd {
		ts = ct.tr.timestamp()
	}
	pushTo(&ct.content, &ct.eventStack, ts, phase, name, separateBeginAndEnd, args...)
}

// Pop ends the most recent trace event block on the per-checker trace buffer.
// Safe to call on nil receiver.
func (ct *CheckerTracing) Pop() {
	if ct == nil {
		return
	}

	ct.mu.Lock()
	defer ct.mu.Unlock()

	if len(ct.eventStack) == 0 {
		return
	}
	var ts float64
	if ct.eventStack[len(ct.eventStack)-1].separateBeginAndEnd {
		ts = ct.tr.timestamp()
	}
	popFrom(&ct.content, &ct.eventStack, ts, ct.tr.startTime)
}

// Instant records an instant event on the per-checker trace buffer.
// Safe to call on nil receiver.
func (ct *CheckerTracing) Instant(phase Phase, name string, args ...string) {
	if ct == nil {
		return
	}

	ct.mu.Lock()
	defer ct.mu.Unlock()

	ts := ct.tr.timestamp()
	ct.content.WriteString(",\n")
	writeEventTo(&ct.content, "I", string(phase), ts, name, writeArgsJSON(args...))
}

// NewCheckerTracing creates a per-checker tracing handle.
// Each checker should have its own CheckerTracing so that checker-specific events
// are written to separate trace files.
func (tr *Tracing) NewCheckerTracing(checkerIndex int) *CheckerTracing {
	if tr == nil {
		return nil
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	ct := &CheckerTracing{tr: tr}
	for len(tr.checkerTracings) <= checkerIndex {
		tr.checkerTracings = append(tr.checkerTracings, nil)
	}
	tr.checkerTracings[checkerIndex] = ct
	return ct
}

// NewTypeTracer creates a new tracer for a specific checker.
// The checkerIndex is used to create unique filenames for each checker's output.
func (tr *Tracing) NewTypeTracer(checkerIndex int) Tracer {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	tracePath := tspath.CombinePaths(tr.traceDir, fmt.Sprintf("trace_%d.json", checkerIndex))
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
		TracePath:      tracePath,
		TypesPath:      typesPath,
	})
	return tracer
}

// StopTracing finalizes the tracing session and writes all output files
func (tr *Tracing) StopTracing() error {
	// Dump types from all tracers BEFORE acquiring the lock, because
	// DumpTypes → buildTypeDescriptor → Display() → TypeToString can
	// re-enter the checker which calls Push/Pop (which need tr.mu).
	for _, tracer := range tr.tracers {
		if err := tracer.DumpTypes(); err != nil {
			return fmt.Errorf("failed to dump types for checker %d: %w", tracer.checkerIndex, err)
		}
	}

	tr.mu.Lock()
	defer tr.mu.Unlock()

	// Close the trace file(s)
	if tr.traceStarted {
		sharedContent := tr.traceContent.String()

		if len(tr.checkerTracings) > 0 {
			// Write per-checker trace files: shared global events + checker-specific events
			for i, ct := range tr.checkerTracings {
				tracePath := tspath.CombinePaths(tr.traceDir, fmt.Sprintf("trace_%d.json", i))
				var full strings.Builder
				full.WriteString(sharedContent)
				if ct != nil {
					full.WriteString(ct.content.String())
				}
				full.WriteString("\n]\n")
				if err := tr.fs.WriteFile(tracePath, full.String(), false); err != nil {
					return fmt.Errorf("failed to write trace file: %w", err)
				}
			}
		} else {
			// No per-checker tracings: write shared content to trace_0.json
			tracePath := tspath.CombinePaths(tr.traceDir, "trace_0.json")
			if err := tr.fs.WriteFile(tracePath, sharedContent+"\n]\n", false); err != nil {
				return fmt.Errorf("failed to write trace file: %w", err)
			}
		}
		tr.traceStarted = false
	}

	// Sort legend entries by typesPath for deterministic output
	slices.SortFunc(tr.legend, func(a, b TraceRecord) int {
		return strings.Compare(a.TypesPath, b.TypesPath)
	})

	// Write the legend file
	legendPath := tspath.CombinePaths(tr.traceDir, "legend.json")
	legendData, err := json.MarshalIndent(tr.legend, "", "  ")
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
