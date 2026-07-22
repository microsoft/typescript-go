// Package spanmap provides a span-aware mapping from positions in a content mapper's transformed
// output back to positions in the original, untransformed source. Unlike a source map, which records
// point correspondences and leaves spans and "no origin" implicit, a SpanMap records explicit segments
// for the parts of the generated text that correspond to the original; positions not covered by any
// segment are synthesized (generated content with no original counterpart). All positions are absolute
// offsets (core.TextPos), matching the compiler's TextRange model.
package spanmap

import (
	"fmt"
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/json"
)

// Kind describes how positions inside a segment relate the generated span to the original span.
type Kind int32

const (
	// KindVerbatim segments are length-preserving: the generated and original spans have the same
	// length and interior positions map 1:1 (origPos = pos - GenStart + OrigStart). A generated span
	// fully within a verbatim segment maps to an exact original span.
	KindVerbatim Kind = iota
	// KindAtom segments map a generated span to an original span as a whole; interior positions are not
	// interpolatable (the lengths may differ), so positions within clamp to the segment's endpoints.
	// Used for renamed identifiers or short expressions.
	KindAtom
)

// Fidelity describes how faithfully a mapped span reflects the original.
type Fidelity int32

const (
	// FidelityExact means the span fell entirely within a single verbatim segment and maps precisely.
	FidelityExact Fidelity = iota
	// FidelityAtom means the span fell within a single atom segment and maps to that atom's span.
	FidelityAtom
	// FidelityApproximate means the span crossed segment boundaries; its endpoints were mapped and clamped.
	FidelityApproximate
	// FidelityNone means the span had no original counterpart (it was entirely synthesized).
	FidelityNone
)

// IsExact reports whether the mapping was fully faithful — the input fell within a single verbatim span —
// so the result maps 1:1 and can host a text edit written back to the original.
func (f Fidelity) IsExact() bool {
	return f == FidelityExact
}

// IsSingleSegment reports whether the input fell within one segment, verbatim or atom, so the result is a
// concrete location rather than a best-effort approximation across boundaries or a synthesized gap.
func (f Fidelity) IsSingleSegment() bool {
	return f == FidelityExact || f == FidelityAtom
}

// IsNone reports whether the input had no original counterpart, meaning the mapped result is a synthesized
// gap that does not correspond to any location in the original text.
func (f Fidelity) IsNone() bool {
	return f == FidelityNone
}

// Segment maps a contiguous generated span to an original span.
type Segment struct {
	GenStart  core.TextPos
	GenEnd    core.TextPos
	OrigStart core.TextPos
	OrigEnd   core.TextPos
	Kind      Kind
}

// SpanMap is a sparse, ordered set of segments over a content mapper's generated text. Segments do not
// need to cover the whole output: any generated position not inside a segment is synthesized (it has no
// original counterpart). An empty SpanMap therefore describes fully synthesized output.
type SpanMap struct {
	segments []Segment

	// origOnce guards lazy construction of origSorted, the segments ordered by OrigStart, used by
	// MapRangeToGenerated for reverse (original -> generated) lookups.
	origOnce   sync.Once
	origSorted []Segment
}

// Validation failures. A content mapper is required to provide a valid span map; these describe the
// ways a map can be malformed, so the compiler can attribute the failure to the mapper precisely and
// point the mapper's author at the offending location.
type MappingErrorKind int

const (
	// MappingErrorKindOverlap means the segments overlap, run backwards, or extend past the end of the
	// transformed text (they must be ordered and disjoint in generated space).
	MappingErrorKindOverlap MappingErrorKind = iota
	// MappingErrorKindOutOfBounds means a segment's original span lies outside the original text.
	MappingErrorKindOutOfBounds
	// MappingErrorKindVerbatimMismatch means a verbatim segment's generated and original text differ.
	MappingErrorKindVerbatimMismatch
)

// MappingError describes a single span map validation failure, including the offsets involved so the mapper's
// author can locate it. GenPos is an offset into the transformed output; OrigPos is an offset into the
// original content. Either may be unused (zero) depending on Kind.
type MappingError struct {
	Kind    MappingErrorKind
	GenPos  core.TextPos
	OrigPos core.TextPos
}

func (p *MappingError) Error() string {
	switch p.Kind {
	case MappingErrorKindOverlap:
		return fmt.Sprintf("content mapper position mappings overlap or are out of order near output offset %d", p.GenPos)
	case MappingErrorKindOutOfBounds:
		return fmt.Sprintf("content mapper position mapping points outside the original content at original offset %d", p.OrigPos)
	case MappingErrorKindVerbatimMismatch:
		return fmt.Sprintf("content mapper verbatim mapping does not match the original content at output offset %d, original offset %d", p.GenPos, p.OrigPos)
	default:
		return "content mapper produced an invalid position mapping"
	}
}

// Validate enforces the content-mapper span map contract against the transformed and original text: the
// segments must be ordered and disjoint in generated space and stay within the transformed text, every
// original span must lie within the original text, and every verbatim segment's text must match the
// original exactly. Gaps are allowed (they map as synthesized) and an empty map is valid. It returns the
// first violation found, or nil if the map is valid.
func (m *SpanMap) Validate(transformed, original string) *MappingError {
	if m == nil {
		return nil
	}
	genLen := core.TextPos(len(transformed))
	origLen := core.TextPos(len(original))
	var prevGenEnd core.TextPos
	for i := range m.segments {
		s := &m.segments[i]
		if s.GenStart < prevGenEnd || s.GenEnd < s.GenStart || s.GenEnd > genLen {
			return &MappingError{Kind: MappingErrorKindOverlap, GenPos: s.GenStart}
		}
		prevGenEnd = s.GenEnd
		if s.OrigStart < 0 || s.OrigEnd < s.OrigStart || s.OrigEnd > origLen {
			return &MappingError{Kind: MappingErrorKindOutOfBounds, GenPos: s.GenStart, OrigPos: s.OrigEnd}
		}
		if s.Kind == KindVerbatim {
			if s.GenEnd-s.GenStart != s.OrigEnd-s.OrigStart ||
				transformed[s.GenStart:s.GenEnd] != original[s.OrigStart:s.OrigEnd] {
				return &MappingError{Kind: MappingErrorKindVerbatimMismatch, GenPos: s.GenStart, OrigPos: s.OrigStart}
			}
		}
	}
	return nil
}

// New builds a SpanMap from segments, sorted by generated start. Segments describe only the parts of the
// generated text that correspond to the original; anything not covered maps as synthesized.
func New(segments []Segment) *SpanMap {
	sorted := slices.Clone(segments)
	slices.SortFunc(sorted, func(a, b Segment) int {
		return int(a.GenStart - b.GenStart)
	})
	return &SpanMap{segments: sorted}
}

// Segments returns the map's segments ordered by generated start.
func (m *SpanMap) Segments() []Segment {
	if m == nil {
		return nil
	}
	return slices.Clone(m.segments)
}

// MapSpan maps a generated range to an original range, along with the fidelity of the result. A generated
// range that lies entirely in a gap between segments (or in an empty map) is synthesized: it maps to the
// insertion point in the original with FidelityNone. A nil SpanMap maps identically.
func (m *SpanMap) MapSpan(r core.TextRange) (core.TextRange, Fidelity) {
	if m == nil {
		return r, FidelityExact
	}
	genStart := core.TextPos(r.Pos())
	genEnd := max(core.TextPos(r.End()), genStart)

	startIdx, startIn := m.segmentIndexAt(genStart)
	endProbe := genEnd
	if genEnd > genStart {
		endProbe = genEnd - 1
	}
	endIdx, endIn := m.segmentIndexAt(endProbe)

	if startIdx == endIdx && startIn == endIn {
		if startIn {
			seg := &m.segments[startIdx]
			if seg.Kind == KindVerbatim {
				origStart := clamp(seg.OrigStart+(genStart-seg.GenStart), seg.OrigStart, seg.OrigEnd)
				origEnd := clamp(seg.OrigStart+(genEnd-seg.GenStart), origStart, seg.OrigEnd)
				return core.NewTextRange(int(origStart), int(origEnd)), FidelityExact
			}
			return core.NewTextRange(int(seg.OrigStart), int(seg.OrigEnd)), FidelityAtom
		}
		// Entirely within a single synthesized gap.
		pos := m.insertionPoint(startIdx)
		return core.NewTextRange(int(pos), int(pos)), FidelityNone
	}

	origStart := m.mapLow(genStart, startIdx, startIn)
	origEnd := max(m.mapHigh(genEnd, endIdx, endIn), origStart)
	return core.NewTextRange(int(origStart), int(origEnd)), FidelityApproximate
}

// MapPosition maps a single generated position to the corresponding original position, along with the
// fidelity of the result. It is the single-position analog of MapSpan: a position in a gap (or in an empty
// map) is synthesized and maps to the insertion point with FidelityNone. A nil SpanMap maps identically.
func (m *SpanMap) MapPosition(pos core.TextPos) (core.TextPos, Fidelity) {
	if m == nil {
		return pos, FidelityExact
	}
	idx, in := m.segmentIndexAt(pos)
	if !in {
		return m.insertionPoint(idx), FidelityNone
	}
	seg := &m.segments[idx]
	if seg.Kind == KindVerbatim {
		return clamp(seg.OrigStart+(pos-seg.GenStart), seg.OrigStart, seg.OrigEnd), FidelityExact
	}
	return seg.OrigStart, FidelityAtom
}

// segmentIndexAt returns the index of the segment containing pos and true, or, when pos lies in a gap,
// the index of the segment immediately before pos (-1 if none) and false.
func (m *SpanMap) segmentIndexAt(pos core.TextPos) (int, bool) {
	idx, found := slices.BinarySearchFunc(m.segments, pos, func(s Segment, p core.TextPos) int {
		return int(s.GenStart - p)
	})
	if found {
		return idx, true
	}
	prev := idx - 1
	if prev >= 0 && pos < m.segments[prev].GenEnd {
		return prev, true
	}
	return prev, false
}

// insertionPoint returns the original offset where synthesized content following segment prev sits: the
// original end of that segment, or 0 before the first segment.
func (m *SpanMap) insertionPoint(prev int) core.TextPos {
	if prev < 0 {
		return 0
	}
	return m.segments[prev].OrigEnd
}

func (m *SpanMap) mapLow(pos core.TextPos, idx int, in bool) core.TextPos {
	if !in {
		return m.insertionPoint(idx)
	}
	seg := &m.segments[idx]
	if seg.Kind == KindVerbatim {
		return clamp(seg.OrigStart+(pos-seg.GenStart), seg.OrigStart, seg.OrigEnd)
	}
	return seg.OrigStart
}

func (m *SpanMap) mapHigh(pos core.TextPos, idx int, in bool) core.TextPos {
	if !in {
		return m.insertionPoint(idx)
	}
	seg := &m.segments[idx]
	if seg.Kind == KindVerbatim {
		return clamp(seg.OrigStart+(pos-seg.GenStart), seg.OrigStart, seg.OrigEnd)
	}
	return seg.OrigEnd
}

// MapRangeToGenerated maps an original range to the corresponding range in the generated text, along with
// the fidelity of the result; it is the inverse of MapSpan. An original range lying entirely in a gap
// (no segment covers it) is synthesized-in-reverse: it maps to the insertion point in the generated text
// with FidelityNone. A nil SpanMap maps identically.
func (m *SpanMap) MapRangeToGenerated(r core.TextRange) (core.TextRange, Fidelity) {
	if m == nil {
		return r, FidelityExact
	}
	segs := m.origIndex()
	origStart := core.TextPos(r.Pos())
	origEnd := max(core.TextPos(r.End()), origStart)

	startIdx, startIn := reverseIndexAt(segs, origStart)
	endProbe := origEnd
	if origEnd > origStart {
		endProbe = origEnd - 1
	}
	endIdx, endIn := reverseIndexAt(segs, endProbe)

	if startIdx == endIdx && startIn == endIn {
		if startIn {
			seg := &segs[startIdx]
			if seg.Kind == KindVerbatim {
				genStart := clamp(seg.GenStart+(origStart-seg.OrigStart), seg.GenStart, seg.GenEnd)
				genEnd := clamp(seg.GenStart+(origEnd-seg.OrigStart), genStart, seg.GenEnd)
				return core.NewTextRange(int(genStart), int(genEnd)), FidelityExact
			}
			return core.NewTextRange(int(seg.GenStart), int(seg.GenEnd)), FidelityAtom
		}
		// Entirely within a single gap: no generated counterpart.
		pos := reverseInsertionPoint(segs, startIdx)
		return core.NewTextRange(int(pos), int(pos)), FidelityNone
	}

	genStart := reverseMapLow(segs, origStart, startIdx, startIn)
	genEnd := max(reverseMapHigh(segs, origEnd, endIdx, endIn), genStart)
	return core.NewTextRange(int(genStart), int(genEnd)), FidelityApproximate
}

// MapPositionToGenerated maps a single original position to the corresponding generated position, along
// with the fidelity of the result; it is the single-position inverse of MapPosition. A position in a gap
// maps to the insertion point in the generated text with FidelityNone. A nil SpanMap maps identically.
func (m *SpanMap) MapPositionToGenerated(pos core.TextPos) (core.TextPos, Fidelity) {
	if m == nil {
		return pos, FidelityExact
	}
	segs := m.origIndex()
	idx, in := reverseIndexAt(segs, pos)
	if !in {
		return reverseInsertionPoint(segs, idx), FidelityNone
	}
	seg := &segs[idx]
	if seg.Kind == KindVerbatim {
		return clamp(seg.GenStart+(pos-seg.OrigStart), seg.GenStart, seg.GenEnd), FidelityExact
	}
	return seg.GenStart, FidelityAtom
}

// origIndex returns the segments ordered by OrigStart, building it once on first use.
func (m *SpanMap) origIndex() []Segment {
	m.origOnce.Do(func() {
		m.origSorted = slices.Clone(m.segments)
		slices.SortFunc(m.origSorted, func(a, b Segment) int {
			return int(a.OrigStart - b.OrigStart)
		})
	})
	return m.origSorted
}

// reverseIndexAt returns the index (within a slice ordered by OrigStart) of the segment whose original
// span contains pos and true, or, when pos lies in a gap, the index of the segment immediately before pos
// (-1 if none) and false.
func reverseIndexAt(segs []Segment, pos core.TextPos) (int, bool) {
	idx, found := slices.BinarySearchFunc(segs, pos, func(s Segment, p core.TextPos) int {
		return int(s.OrigStart - p)
	})
	if found {
		return idx, true
	}
	prev := idx - 1
	if prev >= 0 && pos < segs[prev].OrigEnd {
		return prev, true
	}
	return prev, false
}

// reverseInsertionPoint returns the generated offset where content following segment prev sits: the
// generated end of that segment, or 0 before the first segment.
func reverseInsertionPoint(segs []Segment, prev int) core.TextPos {
	if prev < 0 {
		return 0
	}
	return segs[prev].GenEnd
}

func reverseMapLow(segs []Segment, pos core.TextPos, idx int, in bool) core.TextPos {
	if !in {
		return reverseInsertionPoint(segs, idx)
	}
	seg := &segs[idx]
	if seg.Kind == KindVerbatim {
		return clamp(seg.GenStart+(pos-seg.OrigStart), seg.GenStart, seg.GenEnd)
	}
	return seg.GenStart
}

func reverseMapHigh(segs []Segment, pos core.TextPos, idx int, in bool) core.TextPos {
	if !in {
		return reverseInsertionPoint(segs, idx)
	}
	seg := &segs[idx]
	if seg.Kind == KindVerbatim {
		return clamp(seg.GenStart+(pos-seg.OrigStart), seg.GenStart, seg.GenEnd)
	}
	return seg.GenEnd
}

func clamp(v, lo, hi core.TextPos) core.TextPos {
	return max(lo, min(v, hi))
}

// wireSegment is the JSON tuple form exchanged with an out-of-process content mapper:
// [genStart, genLength, origStart, origLength, kind].
type wireSegment [5]int32

// Unmarshal decodes a SpanMap from the JSON tuple form produced by an out-of-process content mapper.
func Unmarshal(data []byte) (*SpanMap, error) {
	var tuples []wireSegment
	if err := json.Unmarshal(data, &tuples); err != nil {
		return nil, err
	}
	segments := make([]Segment, len(tuples))
	for i, t := range tuples {
		segments[i] = Segment{
			GenStart:  core.TextPos(t[0]),
			GenEnd:    core.TextPos(t[0] + t[1]),
			OrigStart: core.TextPos(t[2]),
			OrigEnd:   core.TextPos(t[2] + t[3]),
			Kind:      Kind(t[4]),
		}
	}
	return New(segments), nil
}

// Marshal encodes a SpanMap into the JSON tuple form.
func (m *SpanMap) Marshal() ([]byte, error) {
	tuples := make([]wireSegment, len(m.segments))
	for i, s := range m.segments {
		tuples[i] = wireSegment{
			int32(s.GenStart),
			int32(s.GenEnd - s.GenStart),
			int32(s.OrigStart),
			int32(s.OrigEnd - s.OrigStart),
			int32(s.Kind),
		}
	}
	return json.Marshal(tuples)
}
