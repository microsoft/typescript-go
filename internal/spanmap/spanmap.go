// Package spanmap provides bidirectional span-aware mapping between a content mapper's transformed
// output and its original, untransformed source. Unlike a source map, which records
// point correspondences and leaves spans and "no origin" implicit, a SpanMap records explicit segments
// for the parts of the generated text that correspond to the original; positions not covered by any
// segment are synthesized (generated content with no original counterpart). All positions are absolute
// offsets (core.TextPos), matching the compiler's TextRange model.
package spanmap

// Keep this in sync with spanMap.ts

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
	// KindAlias has atom geometry, but additionally asserts that the generated and original texts are
	// names for the same logical entity. Diagnostic presentation may substitute the original name.
	KindAlias
)

// Feature selects which language-service operations may use a segment. Diagnostics are intentionally not
// represented: generated diagnostics may not opt out of reporting. Text edits additionally require exact
// verbatim geometry regardless of feature participation.
type Feature int32

const (
	FeatureHover Feature = 1 << iota
	FeatureSignatureHelp
	FeatureCompletion
	FeatureDefinition
	FeatureTypeDefinition
	FeatureImplementation
	FeatureSourceDefinition
	FeatureReferences
	FeatureDocumentHighlights
	FeatureRename
	FeatureCallHierarchy
	FeatureCodeActions
	FeatureFormatting
	FeatureInlayHints
	FeatureSemanticTokens
	FeatureFoldingRanges
	FeatureSelectionRanges
	FeatureLinkedEditing
	FeatureAutoInsert
	FeatureDocumentSymbols
	FeatureCodeLens
	FeatureNone Feature = 0
	FeatureAll  Feature = (FeatureCodeLens << 1) - 1
)

const featureMask = FeatureAll

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

// Segment maps the half-open generated range [GenStart, GenEnd) to the half-open original range
// [OrigStart, OrigEnd). Features controls language-service participation; diagnostics and exact edit mapping
// deliberately bypass it.
type Segment struct {
	GenStart  core.TextPos
	GenEnd    core.TextPos
	OrigStart core.TextPos
	OrigEnd   core.TextPos
	Kind      Kind
	Features  Feature
}

// MappedPosition is one generated projection of an original position and its mapping fidelity.
type MappedPosition struct {
	Position core.TextPos
	Fidelity Fidelity
}

// MappedSpan is one generated projection of an original range and its mapping fidelity.
type MappedSpan struct {
	Span     core.TextRange
	Fidelity Fidelity
}

// SpanMap is a sparse, ordered set of segments over a content mapper's generated text. Segments do not
// need to cover the whole output: any generated position not inside a segment is synthesized (it has no
// original counterpart). An empty SpanMap therefore describes fully synthesized output.
type SpanMap struct {
	segments []Segment

	// origOnce guards lazy construction of origSorted, the segments ordered by OrigStart, used for
	// original-to-generated lookups.
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
	// MappingErrorKindKind means a segment uses an unsupported mapping kind.
	MappingErrorKindKind
	// MappingErrorKindOriginalOverlap means original spans partially overlap or contain one another.
	MappingErrorKindOriginalOverlap
	// MappingErrorKindFeature means a feature annotation contains unsupported flags.
	MappingErrorKindFeature
)

// MappingError describes a single span map validation failure, including the offsets involved so the mapper's
// author can locate it. GenPos is an offset into the transformed output; OrigPos is an offset into the
// original content. Either may be unused (zero) depending on Kind.
type MappingError struct {
	Kind    MappingErrorKind
	GenPos  core.TextPos
	OrigPos core.TextPos
}

// Error describes the invalid mapping and the coordinate at which it was detected.
func (p *MappingError) Error() string {
	switch p.Kind {
	case MappingErrorKindOverlap:
		return fmt.Sprintf("content mapper position mappings overlap or are out of order near output offset %d", p.GenPos)
	case MappingErrorKindOutOfBounds:
		return fmt.Sprintf("content mapper position mapping points outside the original content at original offset %d", p.OrigPos)
	case MappingErrorKindVerbatimMismatch:
		return fmt.Sprintf("content mapper verbatim mapping does not match the original content at output offset %d, original offset %d", p.GenPos, p.OrigPos)
	case MappingErrorKindKind:
		return fmt.Sprintf("content mapper position mapping has an invalid kind at output offset %d", p.GenPos)
	case MappingErrorKindOriginalOverlap:
		return fmt.Sprintf("content mapper position mappings partially overlap in the original content near offset %d", p.OrigPos)
	case MappingErrorKindFeature:
		return fmt.Sprintf("content mapper position mappings have invalid features near original offset %d", p.OrigPos)
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
		if s.Kind != KindVerbatim && s.Kind != KindAtom && s.Kind != KindAlias {
			return &MappingError{Kind: MappingErrorKindKind, GenPos: s.GenStart, OrigPos: s.OrigStart}
		}
		if s.Kind == KindVerbatim {
			if s.GenEnd-s.GenStart != s.OrigEnd-s.OrigStart ||
				transformed[s.GenStart:s.GenEnd] != original[s.OrigStart:s.OrigEnd] {
				return &MappingError{Kind: MappingErrorKindVerbatimMismatch, GenPos: s.GenStart, OrigPos: s.OrigStart}
			}
		}
		if s.Features&^featureMask != 0 {
			return &MappingError{Kind: MappingErrorKindFeature, GenPos: s.GenStart, OrigPos: s.OrigStart}
		}
	}
	originalSegments := m.origIndex()
	for i := 0; i < len(originalSegments); {
		groupEnd := i + 1
		for groupEnd < len(originalSegments) && originalSegments[groupEnd].OrigStart == originalSegments[i].OrigStart && originalSegments[groupEnd].OrigEnd == originalSegments[i].OrigEnd {
			groupEnd++
		}
		if i > 0 && originalSegments[i].OrigStart < originalSegments[i-1].OrigEnd {
			return &MappingError{Kind: MappingErrorKindOriginalOverlap, GenPos: originalSegments[i].GenStart, OrigPos: originalSegments[i].OrigStart}
		}
		i = groupEnd
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

// GeneratedToOriginalSpan maps a generated range to an original range, along with the fidelity of the result. A generated
// range that lies entirely in a gap between segments (or in an empty map) is synthesized: it maps to the
// insertion point in the original with FidelityNone. A nil SpanMap maps identically.
func (m *SpanMap) GeneratedToOriginalSpan(r core.TextRange) (core.TextRange, Fidelity) {
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

// GeneratedToOriginalSpanForFeature maps r only when every generated position in the non-empty range is
// covered by contiguous segments participating in feature. A zero-length range requires its containing
// segment to participate. Diagnostics and edit write-back intentionally use GeneratedToOriginalSpan instead.
func (m *SpanMap) GeneratedToOriginalSpanForFeature(r core.TextRange, feature Feature) (core.TextRange, Fidelity) {
	mapped, fidelity := m.GeneratedToOriginalSpan(r)
	if m == nil || m.generatedSpanSupportsFeature(r, feature) {
		return mapped, fidelity
	}
	return mapped, FidelityNone
}

func (m *SpanMap) generatedSpanSupportsFeature(r core.TextRange, feature Feature) bool {
	start := core.TextPos(r.Pos())
	end := max(core.TextPos(r.End()), start)
	if start == end {
		index, inside := m.segmentIndexAt(start)
		return inside && supportsFeature(m.segments[index], feature)
	}
	index, inside := m.segmentIndexAt(start)
	if !inside {
		return false
	}
	coveredThrough := start
	for index < len(m.segments) && coveredThrough < end {
		segment := m.segments[index]
		if segment.GenStart > coveredThrough || segment.GenEnd <= coveredThrough || !supportsFeature(segment, feature) {
			return false
		}
		coveredThrough = segment.GenEnd
		index++
	}
	return coveredThrough >= end
}

// GeneratedToOriginalPosition maps a single generated position to the corresponding original position, along with the
// fidelity of the result. It is the single-position analog of GeneratedToOriginalSpan: a position in a gap (or in an empty
// map) is synthesized and maps to the insertion point with FidelityNone. A nil SpanMap maps identically.
func (m *SpanMap) GeneratedToOriginalPosition(pos core.TextPos) (core.TextPos, Fidelity) {
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

// GeneratedToOriginalPositionForFeature maps pos only when its generated segment participates in feature.
// Diagnostics and edit write-back intentionally use GeneratedToOriginalPosition instead.
func (m *SpanMap) GeneratedToOriginalPositionForFeature(pos core.TextPos, feature Feature) (core.TextPos, Fidelity) {
	mapped, fidelity := m.GeneratedToOriginalPosition(pos)
	if m == nil {
		return mapped, fidelity
	}
	index, inside := m.segmentIndexAt(pos)
	if !inside || !supportsFeature(m.segments[index], feature) {
		return mapped, FidelityNone
	}
	return mapped, fidelity
}

// AliasForGeneratedSpan returns the alias segment exactly covering r. Partial overlap does not qualify:
// diagnostic text may be substituted only when the diagnostic identifies the complete generated alias.
func (m *SpanMap) AliasForGeneratedSpan(r core.TextRange) (Segment, bool) {
	if m == nil {
		return Segment{}, false
	}
	index, inside := m.segmentIndexAt(core.TextPos(r.Pos()))
	if !inside {
		return Segment{}, false
	}
	segment := m.segments[index]
	return segment, segment.Kind == KindAlias && r.Pos() == int(segment.GenStart) && r.End() == int(segment.GenEnd)
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

// mapLow maps a generated lower range boundary to original coordinates. A boundary in a synthesized
// gap uses that gap's insertion point; an atom uses its original start.
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

// mapHigh maps a generated upper range boundary to original coordinates. A boundary in a synthesized
// gap uses that gap's insertion point; an atom uses its original end.
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

// OriginalToGeneratedPositions returns every generated projection of an original position whose segment
// participates in feature. Segment ends are inclusive for point mapping, so a position shared by adjacent
// original spans returns projections from both sides. Results are ordered by generated position. It returns
// no results for an uncovered position or when all touching segments reject feature. A nil SpanMap maps identically.
func (m *SpanMap) OriginalToGeneratedPositions(pos core.TextPos, feature Feature) []MappedPosition {
	if m == nil {
		return []MappedPosition{{Position: pos, Fidelity: FidelityExact}}
	}
	groups := originalSegmentGroupsAtPoint(m.origIndex(), pos)
	if len(groups) == 0 {
		return nil
	}
	var results []MappedPosition
	for _, group := range groups {
		for _, segment := range group.segments {
			if !supportsFeature(segment, feature) {
				continue
			}
			mapped := MappedPosition{Fidelity: FidelityAtom}
			if segment.Kind == KindVerbatim {
				mapped.Position = clamp(segment.GenStart+(pos-segment.OrigStart), segment.GenStart, segment.GenEnd)
				mapped.Fidelity = FidelityExact
			} else if group.atEnd {
				mapped.Position = segment.GenEnd
			} else {
				mapped.Position = segment.GenStart
			}
			if !slices.Contains(results, mapped) {
				results = append(results, mapped)
			}
		}
	}
	slices.SortFunc(results, func(a, b MappedPosition) int { return int(a.Position - b.Position) })
	return results
}

// OriginalToGeneratedSpans returns every feature-compatible generated projection of an original range.
// A range contained by one duplicate group produces one exact or atom result per matching group member.
//
// A range that starts in one group and ends in another can have several possible generated ranges. For
// example, suppose two original segments are each copied twice into the generated text:
//
//	original:   [ A ][ B ]
//	               [---)       range from inside A to inside B
//
//	generated:  [ A ][ B ]      [ A ][ B ]
//	               ^   ^          ^   ^
//	             start end      start end
//	               1   3          11  13
//
// The map says that the range may start at 1 or 11 and end at 3 or 13, but it does not say which copy of A
// belongs with which copy of B. We choose the smallest range around each possible location, producing [1,3)
// and [11,13). We do not return [1,13), because it contains both smaller candidates and would include code
// that may be unrelated to the original range. These cross-group results have approximate fidelity.
// If either boundary is uncovered or disabled for feature, there are no results. A nil SpanMap maps identically.
func (m *SpanMap) OriginalToGeneratedSpans(r core.TextRange, feature Feature) []MappedSpan {
	if m == nil {
		return []MappedSpan{{Span: r, Fidelity: FidelityExact}}
	}
	start := core.TextPos(r.Pos())
	end := max(core.TextPos(r.End()), start)
	lastCharacter := end
	if end > start {
		lastCharacter--
	}
	originalSegments := m.origIndex()
	startSegments, startInside := originalSegmentsAt(originalSegments, start)
	endSegments, endInside := originalSegmentsAt(originalSegments, lastCharacter)
	if !startInside || !endInside {
		return nil
	}
	if sameOriginalRange(startSegments[0], endSegments[0]) {
		return originalToGeneratedSpansInGroup(startSegments, start, end, feature)
	}
	starts := originalStartProjections(startSegments, start, feature)
	ends := originalEndProjections(endSegments, end, feature)
	if len(starts) == 0 || len(ends) == 0 {
		return nil
	}
	results := make([]MappedSpan, 0, min(len(starts), len(ends)))
	for i, genStart := range starts {
		endIndex, _ := slices.BinarySearch(ends, genStart)
		if endIndex == len(ends) || i+1 < len(starts) && starts[i+1] <= ends[endIndex] {
			continue
		}
		results = append(results, MappedSpan{
			Span:     core.NewTextRange(int(genStart), int(ends[endIndex])),
			Fidelity: FidelityApproximate,
		})
	}
	return results
}

// originalStartProjections maps the inclusive start of an original range through every matching segment.
// Verbatim segments preserve the offset within the segment; atoms map to their generated start.
//
// For duplicate verbatim segments, the start keeps the same relative offset in every copy:
//
//	original:       [---------)
//	                   ^ start
//
//	generated:  [---------)   [---------)
//	               ^             ^
//	             result        result
func originalStartProjections(segments []Segment, start core.TextPos, feature Feature) []core.TextPos {
	results := make([]core.TextPos, 0, len(segments))
	for _, segment := range segments {
		if !supportsFeature(segment, feature) {
			continue
		}
		if segment.Kind == KindVerbatim {
			results = append(results, clamp(segment.GenStart+(start-segment.OrigStart), segment.GenStart, segment.GenEnd))
		} else {
			results = append(results, segment.GenStart)
		}
	}
	return results
}

// originalEndProjections maps the exclusive end of an original range through every matching segment.
// The caller uses end-1 to find the segment containing the final character, while this helper maps the end
// boundary itself. Verbatim segments preserve that boundary; atoms map to their generated end.
//
// The lookup uses end-1 so an end at a segment boundary selects the segment on its left, not the next one:
//
//	original:       [---------)[ next segment )
//	                         ^`-- end
//	                         `--- end-1
//
//	generated:  [---------)   [---------)
//	                      ^             ^
//	                    result        result
func originalEndProjections(segments []Segment, end core.TextPos, feature Feature) []core.TextPos {
	results := make([]core.TextPos, 0, len(segments))
	for _, segment := range segments {
		if !supportsFeature(segment, feature) {
			continue
		}
		if segment.Kind == KindVerbatim {
			results = append(results, clamp(segment.GenStart+(end-segment.OrigStart), segment.GenStart, segment.GenEnd))
		} else {
			results = append(results, segment.GenEnd)
		}
	}
	return results
}

// originalToGeneratedSpansInGroup maps a range whose boundaries are known to lie in segments.
func originalToGeneratedSpansInGroup(segments []Segment, start core.TextPos, end core.TextPos, feature Feature) []MappedSpan {
	results := make([]MappedSpan, 0, len(segments))
	for _, segment := range segments {
		if !supportsFeature(segment, feature) {
			continue
		}
		if segment.Kind == KindVerbatim {
			genStart := clamp(segment.GenStart+(start-segment.OrigStart), segment.GenStart, segment.GenEnd)
			genEnd := clamp(segment.GenStart+(end-segment.OrigStart), genStart, segment.GenEnd)
			results = append(results, MappedSpan{Span: core.NewTextRange(int(genStart), int(genEnd)), Fidelity: FidelityExact})
		} else {
			results = append(results, MappedSpan{Span: core.NewTextRange(int(segment.GenStart), int(segment.GenEnd)), Fidelity: FidelityAtom})
		}
	}
	return results
}

// sameOriginalRange reports whether two segments belong to the same duplicate group.
func sameOriginalRange(left Segment, right Segment) bool {
	return left.OrigStart == right.OrigStart && left.OrigEnd == right.OrigEnd
}

// origIndex returns the segments ordered by OrigStart, building it once on first use.
func (m *SpanMap) origIndex() []Segment {
	m.origOnce.Do(func() {
		m.origSorted = slices.Clone(m.segments)
		slices.SortFunc(m.origSorted, func(a, b Segment) int {
			if c := int(a.OrigStart - b.OrigStart); c != 0 {
				return c
			}
			if c := int(a.OrigEnd - b.OrigEnd); c != 0 {
				return c
			}
			return int(a.GenStart - b.GenStart)
		})
	})
	return m.origSorted
}

// originalSegmentsAt returns the complete duplicate group containing pos from a slice ordered by original
// start, original end, and generated start. Segment ends are exclusive; a segment start, including a zero-length
// segment, is considered contained. It finds a candidate in O(log n), then scans only the duplicate group.
// The boolean reports whether any group contains pos.
func originalSegmentsAt(segments []Segment, pos core.TextPos) ([]Segment, bool) {
	index, found := slices.BinarySearchFunc(segments, pos, func(segment Segment, position core.TextPos) int {
		return int(segment.OrigStart - position)
	})
	if !found {
		index--
	}
	if index < 0 || !(segments[index].OrigStart == pos || pos < segments[index].OrigEnd) {
		return nil, false
	}
	start := index
	for start > 0 && sameOriginalRange(segments[start-1], segments[index]) {
		start--
	}
	end := start + 1
	for end < len(segments) && sameOriginalRange(segments[end], segments[start]) {
		end++
	}
	return segments[start:end], true
}

type originalSegmentGroupAtPoint struct {
	segments []Segment
	atEnd    bool
}

// originalSegmentGroupsAtPoint returns the duplicate original-range groups containing or touching pos.
// Interior positions return one group. At a boundary between adjacent groups, both the group ending at pos
// and the group starting at pos are returned. segments must be ordered by OrigStart, OrigEnd, then GenStart.
func originalSegmentGroupsAtPoint(segments []Segment, pos core.TextPos) []originalSegmentGroupAtPoint {
	index, startsAtPosition := slices.BinarySearchFunc(segments, pos, func(segment Segment, position core.TextPos) int {
		return int(segment.OrigStart - position)
	})
	if startsAtPosition {
		right, _ := originalSegmentsAt(segments, pos)
		var groups []originalSegmentGroupAtPoint
		if index > 0 {
			leftIndex := index - 1
			if segments[leftIndex].OrigEnd == pos {
				leftStart := leftIndex
				for leftStart > 0 && sameOriginalRange(segments[leftStart-1], segments[leftIndex]) {
					leftStart--
				}
				groups = append(groups, originalSegmentGroupAtPoint{segments: segments[leftStart:index], atEnd: true})
			}
		}
		return append(groups, originalSegmentGroupAtPoint{segments: right})
	}
	if index == 0 {
		return nil
	}
	leftIndex := index - 1
	segment := segments[leftIndex]
	if pos > segment.OrigEnd {
		return nil
	}
	start := leftIndex
	for start > 0 && sameOriginalRange(segments[start-1], segment) {
		start--
	}
	return []originalSegmentGroupAtPoint{{segments: segments[start:index], atEnd: pos == segment.OrigEnd}}
}

// supportsFeature reports whether segment participates in feature.
func supportsFeature(segment Segment, feature Feature) bool {
	return segment.Features&feature != 0
}

// clamp confines v to the inclusive interval [lo, hi].
func clamp(v, lo, hi core.TextPos) core.TextPos {
	return max(lo, min(v, hi))
}

// Unmarshal decodes a SpanMap from the JSON tuple form produced by an out-of-process content mapper.
// Five-element tuples omit features and are normalized to FeatureAll; six-element tuples preserve the
// explicit feature mask, including FeatureNone.
func Unmarshal(data []byte) (*SpanMap, error) {
	var tuples [][]int32
	if err := json.Unmarshal(data, &tuples); err != nil {
		return nil, err
	}
	segments := make([]Segment, len(tuples))
	for i, t := range tuples {
		if len(t) != 5 && len(t) != 6 {
			return nil, fmt.Errorf("span map segment %d: expected 5 or 6 values, got %d", i, len(t))
		}
		segments[i] = Segment{
			GenStart:  core.TextPos(t[0]),
			GenEnd:    core.TextPos(t[0] + t[1]),
			OrigStart: core.TextPos(t[2]),
			OrigEnd:   core.TextPos(t[2] + t[3]),
			Kind:      Kind(t[4]),
			Features:  FeatureAll,
		}
		if len(t) == 6 {
			segments[i].Features = Feature(t[5])
		}
	}
	return New(segments), nil
}

// Marshal encodes a SpanMap into the JSON tuple form. FeatureAll uses the backward-compatible five-element
// tuple; every other feature mask is emitted as a sixth element.
func (m *SpanMap) Marshal() ([]byte, error) {
	tuples := make([][]int32, len(m.segments))
	for i, s := range m.segments {
		tuples[i] = []int32{
			int32(s.GenStart),
			int32(s.GenEnd - s.GenStart),
			int32(s.OrigStart),
			int32(s.OrigEnd - s.OrigStart),
			int32(s.Kind),
		}
		if s.Features != FeatureAll {
			tuples[i] = append(tuples[i], int32(s.Features))
		}
	}
	return json.Marshal(tuples)
}
