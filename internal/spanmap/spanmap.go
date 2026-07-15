// Package spanmap provides a span-aware mapping from positions in a content mapper's transformed
// output back to positions in the original, untransformed source. Unlike a source map, which records
// point correspondences and leaves spans and "no origin" implicit, a SpanMap tiles the generated text
// into explicit segments so that a generated span can be mapped to an original span with a known
// fidelity. All positions are absolute offsets (core.TextPos), matching the compiler's TextRange model.
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
	// KindSynthesized segments are generated content with no original counterpart (e.g. a synthesized
	// import). Their OrigStart == OrigEnd marks the point in the original text where the content was
	// inserted, used as a fallback location.
	KindSynthesized
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

// SpanMap is a sorted, gap-free tiling of a content mapper's generated text.
type SpanMap struct {
	segments []Segment
}

// Validation failures. A content mapper is required to provide a valid span map; these describe the
// ways a map can be malformed, so the compiler can attribute the failure to the mapper precisely and
// point the mapper's author at the offending location.
type MappingErrorKind int

const (
	// MappingErrorKindMissing means no mappings were provided at all.
	MappingErrorKindMissing MappingErrorKind = iota
	// MappingErrorKindCoverage means the segments do not tile the entire transformed text (a gap, overlap, or a
	// segment extending past the end of the transformed text).
	MappingErrorKindCoverage
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
	case MappingErrorKindCoverage:
		return fmt.Sprintf("content mapper position mappings do not cover the transformed output near output offset %d", p.GenPos)
	case MappingErrorKindOutOfBounds:
		return fmt.Sprintf("content mapper position mapping points outside the original content at original offset %d", p.OrigPos)
	case MappingErrorKindVerbatimMismatch:
		return fmt.Sprintf("content mapper verbatim mapping does not match the original content at output offset %d, original offset %d", p.GenPos, p.OrigPos)
	default:
		return "content mapper did not provide the required position mappings"
	}
}

// Validate enforces the content-mapper span map contract against the transformed and original text: the
// segments must tile the whole transformed text with no gaps or overlaps, every original span must lie
// within the original text, and every verbatim segment's text must match the original exactly. It
// returns the first violation found, or nil if the map is valid.
func (m *SpanMap) Validate(transformed, original string) *MappingError {
	if m == nil || len(m.segments) == 0 {
		return &MappingError{Kind: MappingErrorKindMissing}
	}
	genLen := core.TextPos(len(transformed))
	origLen := core.TextPos(len(original))
	var expectedGenStart core.TextPos
	for i := range m.segments {
		s := &m.segments[i]
		if s.GenStart != expectedGenStart || s.GenEnd < s.GenStart || s.GenEnd > genLen {
			return &MappingError{Kind: MappingErrorKindCoverage, GenPos: expectedGenStart}
		}
		expectedGenStart = s.GenEnd
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
	if expectedGenStart != genLen {
		return &MappingError{Kind: MappingErrorKindCoverage, GenPos: expectedGenStart}
	}
	return nil
}

// New builds a SpanMap from segments, which are sorted by generated start. The segments are expected to
// tile the generated text; gaps and overlaps are tolerated but map on a best-effort basis.
func New(segments []Segment) *SpanMap {
	sorted := slices.Clone(segments)
	slices.SortFunc(sorted, func(a, b Segment) int {
		return int(a.GenStart - b.GenStart)
	})
	return &SpanMap{segments: sorted}
}

// MapSpan maps a generated range to an original range, along with the fidelity of the result. A nil or
// empty SpanMap maps identically.
func (m *SpanMap) MapSpan(r core.TextRange) (core.TextRange, Fidelity) {
	if m == nil || len(m.segments) == 0 {
		return r, FidelityExact
	}
	genStart := core.TextPos(r.Pos())
	genEnd := max(core.TextPos(r.End()), genStart)

	startIdx := m.findSegmentIndex(genStart)
	endProbe := genEnd
	if genEnd > genStart {
		endProbe = genEnd - 1
	}
	endIdx := m.findSegmentIndex(endProbe)

	startSeg := &m.segments[startIdx]
	endSeg := &m.segments[endIdx]

	origStart := mapLow(genStart, startSeg)
	origEnd := max(mapHigh(genEnd, endSeg), origStart)

	var fidelity Fidelity
	if startIdx == endIdx {
		switch startSeg.Kind {
		case KindVerbatim:
			fidelity = FidelityExact
		case KindAtom:
			fidelity = FidelityAtom
		default:
			fidelity = FidelityNone
		}
	} else {
		fidelity = FidelityApproximate
	}

	return core.NewTextRange(int(origStart), int(origEnd)), fidelity
}

// findSegmentIndex returns the index of the segment containing pos, clamping to the first or last
// segment when pos lies outside the tiling.
func (m *SpanMap) findSegmentIndex(pos core.TextPos) int {
	idx, found := slices.BinarySearchFunc(m.segments, pos, func(s Segment, p core.TextPos) int {
		return int(s.GenStart - p)
	})
	if found {
		return idx
	}
	if idx == 0 {
		return 0
	}
	return idx - 1
}

func mapLow(pos core.TextPos, seg *Segment) core.TextPos {
	if seg.Kind == KindVerbatim {
		return clamp(seg.OrigStart+(pos-seg.GenStart), seg.OrigStart, seg.OrigEnd)
	}
	return seg.OrigStart
}

func mapHigh(pos core.TextPos, seg *Segment) core.TextPos {
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
