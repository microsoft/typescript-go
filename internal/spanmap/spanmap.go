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
