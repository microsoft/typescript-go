package spanmap_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"gotest.tools/v3/assert"
)

func TestGeneratedToOriginalSpanVerbatim(t *testing.T) {
	t.Parallel()

	// Generated [0,10) is a verbatim copy of original [100,110).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
	})

	got, fidelity := m.GeneratedToOriginalSpan(core.NewTextRange(3, 7))
	assert.Equal(t, got.Pos(), 103)
	assert.Equal(t, got.End(), 107)
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestGeneratedToOriginalSpanAtom(t *testing.T) {
	t.Parallel()

	// Generated [0,3) is a synthesized gap; [3,14) ("MyComponent") is an atom of the original [60,71).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 3, GenEnd: 14, OrigStart: 60, OrigEnd: 71, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeAll},
	})

	// A span inside the atom maps to the whole atom span.
	got, fidelity := m.GeneratedToOriginalSpan(core.NewTextRange(5, 9))
	assert.Equal(t, got.Pos(), 60)
	assert.Equal(t, got.End(), 71)
	assert.Equal(t, fidelity, spanmap.FidelityAtom)
}

func TestGeneratedToOriginalSpanSynthesizedGap(t *testing.T) {
	t.Parallel()

	// A gap between two verbatim segments is synthesized: it maps to the insertion point (the preceding
	// segment's original end) with no fidelity.
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
	})

	got, fidelity := m.GeneratedToOriginalSpan(core.NewTextRange(12, 15))
	assert.Equal(t, got.Pos(), 110)
	assert.Equal(t, got.End(), 110)
	assert.Equal(t, fidelity, spanmap.FidelityNone)
}

func TestGeneratedToOriginalSpanEmptyIsSynthesized(t *testing.T) {
	t.Parallel()

	// An empty map describes fully synthesized output: everything maps to the start with no fidelity.
	m := spanmap.New(nil)
	got, fidelity := m.GeneratedToOriginalSpan(core.NewTextRange(5, 10))
	assert.Equal(t, got.Pos(), 0)
	assert.Equal(t, got.End(), 0)
	assert.Equal(t, fidelity, spanmap.FidelityNone)
}

func TestGeneratedToOriginalSpanCrossingSegments(t *testing.T) {
	t.Parallel()

	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim},
		{GenStart: 10, GenEnd: 20, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim},
	})

	got, fidelity := m.GeneratedToOriginalSpan(core.NewTextRange(5, 15))
	assert.Equal(t, got.Pos(), 105)
	assert.Equal(t, got.End(), 205)
	assert.Equal(t, fidelity, spanmap.FidelityApproximate)
}

func TestGeneratedToOriginalSpanNilIdentity(t *testing.T) {
	t.Parallel()

	var m *spanmap.SpanMap
	got, fidelity := m.GeneratedToOriginalSpan(core.NewTextRange(3, 7))
	assert.Equal(t, got.Pos(), 3)
	assert.Equal(t, got.End(), 7)
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestGeneratedToOriginalPosition(t *testing.T) {
	t.Parallel()

	// Generated [0,10) is a verbatim copy of original [100,110); [10,20) is an atom of original [200,210).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeAll},
	})

	testCases := []struct {
		name     string
		pos      core.TextPos
		want     core.TextPos
		fidelity spanmap.Fidelity
	}{
		{"verbatim interpolates", 3, 103, spanmap.FidelityExact},
		{"atom maps to its start", 25, 200, spanmap.FidelityAtom},
		{"gap maps to insertion point", 15, 110, spanmap.FidelityNone},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, fidelity := m.GeneratedToOriginalPosition(tc.pos)
			assert.Equal(t, got, tc.want)
			assert.Equal(t, fidelity, tc.fidelity)
			// GeneratedToOriginalPosition must agree with GeneratedToOriginalSpan on a zero-length range.
			span, spanFidelity := m.GeneratedToOriginalSpan(core.NewTextRange(int(tc.pos), int(tc.pos)))
			assert.Equal(t, got, core.TextPos(span.Pos()))
			assert.Equal(t, fidelity, spanFidelity)
		})
	}
}

func TestMapPositionNilIdentity(t *testing.T) {
	t.Parallel()

	var m *spanmap.SpanMap
	got, fidelity := m.GeneratedToOriginalPosition(7)
	assert.Equal(t, got, core.TextPos(7))
	assert.Equal(t, fidelity, spanmap.FidelityExact)
}

func TestOriginalToGeneratedSpanVerbatim(t *testing.T) {
	t.Parallel()

	// Generated [0,10) is a verbatim copy of original [100,110).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
	})

	results := m.OriginalToGeneratedSpans(core.NewTextRange(103, 107), spanmap.PurposeAll)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Span.Pos(), 3)
	assert.Equal(t, results[0].Span.End(), 7)
	assert.Equal(t, results[0].Fidelity, spanmap.FidelityExact)
}

func TestOriginalToGeneratedSpanAtom(t *testing.T) {
	t.Parallel()

	// Generated [3,14) is an atom of the original [60,71).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 3, GenEnd: 14, OrigStart: 60, OrigEnd: 71, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeAll},
	})

	// A span inside the original atom maps to the whole generated span.
	results := m.OriginalToGeneratedSpans(core.NewTextRange(63, 67), spanmap.PurposeAll)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Span.Pos(), 3)
	assert.Equal(t, results[0].Span.End(), 14)
	assert.Equal(t, results[0].Fidelity, spanmap.FidelityAtom)
}

func TestOriginalToGeneratedSpanGap(t *testing.T) {
	t.Parallel()

	// An original range with no covering segment has no generated counterpart.
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
	})

	assert.Equal(t, len(m.OriginalToGeneratedSpans(core.NewTextRange(150, 160), spanmap.PurposeAll)), 0)
}

func TestOriginalToGeneratedSpanNilIdentity(t *testing.T) {
	t.Parallel()

	var m *spanmap.SpanMap
	results := m.OriginalToGeneratedSpans(core.NewTextRange(3, 7), spanmap.PurposeAll)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Span.Pos(), 3)
	assert.Equal(t, results[0].Span.End(), 7)
	assert.Equal(t, results[0].Fidelity, spanmap.FidelityExact)
}

func TestOriginalToGeneratedPositions(t *testing.T) {
	t.Parallel()

	// Original [100,110) is a verbatim copy of generated [0,10); [200,210) is an atom of generated [20,30).
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
		{GenStart: 20, GenEnd: 30, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeAll},
	})

	testCases := []struct {
		name     string
		pos      core.TextPos
		want     core.TextPos
		fidelity spanmap.Fidelity
	}{
		{"verbatim interpolates", 103, 3, spanmap.FidelityExact},
		{"atom maps to its start", 205, 20, spanmap.FidelityAtom},
		{"gap has no projection", 150, 0, spanmap.FidelityNone},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			positions := m.OriginalToGeneratedPositions(tc.pos, spanmap.PurposeAll)
			spans := m.OriginalToGeneratedSpans(core.NewTextRange(int(tc.pos), int(tc.pos)), spanmap.PurposeAll)
			if tc.fidelity.IsNone() {
				assert.Equal(t, len(positions), 0)
				assert.Equal(t, len(spans), 0)
				return
			}
			assert.Equal(t, len(positions), 1)
			assert.Equal(t, positions[0].Position, tc.want)
			assert.Equal(t, positions[0].Fidelity, tc.fidelity)
			assert.Equal(t, len(spans), 1)
			assert.Equal(t, core.TextPos(spans[0].Span.Pos()), tc.want)
			assert.Equal(t, spans[0].Fidelity, tc.fidelity)
		})
	}
}

func TestOriginalToGeneratedDuplicateGroup(t *testing.T) {
	t.Parallel()

	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 3, OrigStart: 10, OrigEnd: 13, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeNavigation},
		{GenStart: 10, GenEnd: 13, OrigStart: 10, OrigEnd: 13, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeSemantic},
		{GenStart: 20, GenEnd: 25, OrigStart: 10, OrigEnd: 13, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeNavigation},
	})

	semantic := m.OriginalToGeneratedPositions(11, spanmap.PurposeSemantic)
	assert.Equal(t, len(semantic), 1)
	assert.Equal(t, semantic[0].Position, core.TextPos(11))
	assert.Equal(t, semantic[0].Fidelity, spanmap.FidelityExact)

	navigation := m.OriginalToGeneratedPositions(11, spanmap.PurposeNavigation)
	assert.Equal(t, len(navigation), 2)
	assert.Equal(t, navigation[0].Position, core.TextPos(1))
	assert.Equal(t, navigation[0].Fidelity, spanmap.FidelityExact)
	assert.Equal(t, navigation[1].Position, core.TextPos(20))
	assert.Equal(t, navigation[1].Fidelity, spanmap.FidelityAtom)

	spans := m.OriginalToGeneratedSpans(core.NewTextRange(10, 13), spanmap.PurposeNavigation)
	assert.Equal(t, len(spans), 2)
	assert.Equal(t, spans[0].Span.Pos(), 0)
	assert.Equal(t, spans[0].Span.End(), 3)
	assert.Equal(t, spans[1].Span.Pos(), 20)
	assert.Equal(t, spans[1].Span.End(), 25)
}

func TestOriginalToGeneratedCrossGroupProjections(t *testing.T) {
	t.Parallel()

	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 2, OrigStart: 0, OrigEnd: 2, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeSemantic},
		{GenStart: 2, GenEnd: 4, OrigStart: 2, OrigEnd: 4, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeSemantic},
		{GenStart: 10, GenEnd: 12, OrigStart: 0, OrigEnd: 2, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeSemantic},
		{GenStart: 12, GenEnd: 14, OrigStart: 2, OrigEnd: 4, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeSemantic},
	})

	spans := m.OriginalToGeneratedSpans(core.NewTextRange(1, 3), spanmap.PurposeSemantic)
	assert.Equal(t, len(spans), 2)
	assert.Equal(t, spans[0].Span, core.NewTextRange(1, 3))
	assert.Equal(t, spans[1].Span, core.NewTextRange(11, 13))
	for _, mapped := range spans {
		assert.Equal(t, mapped.Fidelity, spanmap.FidelityApproximate)
	}
}

func TestOriginalToGeneratedExplicitZeroPurpose(t *testing.T) {
	t.Parallel()

	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 3, OrigStart: 10, OrigEnd: 13, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeNone},
	})

	assert.Equal(t, len(m.OriginalToGeneratedPositions(11, spanmap.PurposeSemantic)), 0)
	assert.Equal(t, len(m.OriginalToGeneratedPositions(11, spanmap.PurposeNavigation)), 0)
	assert.Equal(t, len(m.OriginalToGeneratedSpans(core.NewTextRange(10, 13), spanmap.PurposeSemantic)), 0)

	data, err := m.Marshal()
	assert.NilError(t, err)
	assert.Equal(t, string(data), "[[0,3,10,3,0,0]]")
	decoded, err := spanmap.Unmarshal(data)
	assert.NilError(t, err)
	segments := decoded.Segments()
	assert.Equal(t, segments[0].Purpose, spanmap.PurposeNone)

	legacy, err := spanmap.Unmarshal([]byte("[[0,3,10,3,0]]"))
	assert.NilError(t, err)
	assert.Equal(t, legacy.Segments()[0].Purpose, spanmap.PurposeAll)
	assert.Equal(t, len(legacy.OriginalToGeneratedPositions(11, spanmap.PurposeSemantic)), 1)
}

func TestOriginalToGeneratedSpanRoundTrip(t *testing.T) {
	t.Parallel()

	// Original spans are out of order relative to generated spans, exercising the reverse index.
	m := spanmap.New([]spanmap.Segment{
		{GenStart: 0, GenEnd: 10, OrigStart: 200, OrigEnd: 210, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
		{GenStart: 10, GenEnd: 20, OrigStart: 100, OrigEnd: 110, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeAll},
	})

	for _, r := range []core.TextRange{core.NewTextRange(2, 8), core.NewTextRange(12, 18)} {
		orig, fidelity := m.GeneratedToOriginalSpan(r)
		assert.Equal(t, fidelity, spanmap.FidelityExact)
		back := m.OriginalToGeneratedSpans(orig, spanmap.PurposeAll)
		assert.Equal(t, len(back), 1)
		assert.Equal(t, back[0].Fidelity, spanmap.FidelityExact)
		assert.Equal(t, back[0].Span.Pos(), r.Pos())
		assert.Equal(t, back[0].Span.End(), r.End())
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	t.Parallel()

	original := spanmap.New([]spanmap.Segment{
		{GenStart: 3, GenEnd: 14, OrigStart: 60, OrigEnd: 71, Kind: spanmap.KindAtom},
		{GenStart: 14, GenEnd: 24, OrigStart: 71, OrigEnd: 81, Kind: spanmap.KindVerbatim},
	})

	data, err := original.Marshal()
	assert.NilError(t, err)
	decoded, err := spanmap.Unmarshal(data)
	assert.NilError(t, err)

	for _, r := range []core.TextRange{core.NewTextRange(1, 2), core.NewTextRange(4, 10), core.NewTextRange(16, 20)} {
		wantRange, wantFidelity := original.GeneratedToOriginalSpan(r)
		gotRange, gotFidelity := decoded.GeneratedToOriginalSpan(r)
		assert.Equal(t, gotRange, wantRange)
		assert.Equal(t, gotFidelity, wantFidelity)
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	const transformed = "const greeting = 1;\n"
	const original = "<x>const greeting = 1;\n</x>"
	scriptStart := 3 // index of "const" in original

	testCases := []struct {
		name     string
		segs     []spanmap.Segment
		wantKind spanmap.MappingErrorKind
		wantOK   bool
	}{
		{
			name:   "valid verbatim",
			segs:   []spanmap.Segment{{GenStart: 0, GenEnd: core.TextPos(len(transformed)), OrigStart: core.TextPos(scriptStart), OrigEnd: core.TextPos(scriptStart + len(transformed)), Kind: spanmap.KindVerbatim}},
			wantOK: true,
		},
		{
			name:   "empty is valid",
			segs:   nil,
			wantOK: true,
		},
		{
			name:   "gap is allowed",
			segs:   []spanmap.Segment{{GenStart: 3, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom}},
			wantOK: true,
		},
		{
			name: "overlap",
			segs: []spanmap.Segment{
				{GenStart: 0, GenEnd: 10, OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom},
				{GenStart: 5, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom},
			},
			wantKind: spanmap.MappingErrorKindOverlap,
		},
		{
			name:     "original out of bounds",
			segs:     []spanmap.Segment{{GenStart: 0, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: core.TextPos(len(original) + 10), Kind: spanmap.KindAtom}},
			wantKind: spanmap.MappingErrorKindOutOfBounds,
		},
		{
			name:     "verbatim text mismatch",
			segs:     []spanmap.Segment{{GenStart: 0, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: core.TextPos(len(transformed)), Kind: spanmap.KindVerbatim}},
			wantKind: spanmap.MappingErrorKindVerbatimMismatch,
		},
		{
			name:     "unknown kind",
			segs:     []spanmap.Segment{{GenStart: 0, GenEnd: 1, OrigStart: 0, OrigEnd: 1, Kind: 2}},
			wantKind: spanmap.MappingErrorKindKind,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			problem := spanmap.New(tc.segs).Validate(transformed, original)
			if tc.wantOK {
				assert.Assert(t, problem == nil, "expected valid, got %+v", problem)
				return
			}
			assert.Assert(t, problem != nil, "expected a problem")
			assert.Equal(t, problem.Kind, tc.wantKind)
		})
	}
}

func TestValidateOriginalOverlapAndPurposes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		segments []spanmap.Segment
		wantKind spanmap.MappingErrorKind
		valid    bool
	}{
		{
			name: "identical duplicate group",
			segments: []spanmap.Segment{
				{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeNavigation},
				{GenStart: 3, GenEnd: 6, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindVerbatim, Purpose: spanmap.PurposeSemantic},
			},
			valid: true,
		},
		{
			name: "partial original overlap",
			segments: []spanmap.Segment{
				{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindAtom},
				{GenStart: 3, GenEnd: 6, OrigStart: 2, OrigEnd: 5, Kind: spanmap.KindAtom},
			},
			wantKind: spanmap.MappingErrorKindOriginalOverlap,
		},
		{
			name: "nested original overlap",
			segments: []spanmap.Segment{
				{GenStart: 0, GenEnd: 5, OrigStart: 0, OrigEnd: 5, Kind: spanmap.KindAtom},
				{GenStart: 5, GenEnd: 6, OrigStart: 1, OrigEnd: 4, Kind: spanmap.KindAtom},
			},
			wantKind: spanmap.MappingErrorKindOriginalOverlap,
		},
		{
			name: "duplicate without explicit purpose is tolerant",
			segments: []spanmap.Segment{
				{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindAtom},
				{GenStart: 3, GenEnd: 6, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeNavigation},
			},
			valid: true,
		},
		{
			name: "duplicate with two semantic members is tolerant",
			segments: []spanmap.Segment{
				{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeSemantic},
				{GenStart: 3, GenEnd: 6, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeSemantic | spanmap.PurposeNavigation},
			},
			valid: true,
		},
		{
			name: "purpose on sole cover is valid",
			segments: []spanmap.Segment{
				{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindAtom, Purpose: spanmap.PurposeNavigation},
			},
			valid: true,
		},
		{
			name: "unknown purpose flag",
			segments: []spanmap.Segment{
				{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 3, Kind: spanmap.KindAtom, Purpose: 4},
			},
			wantKind: spanmap.MappingErrorKindPurpose,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			problem := spanmap.New(test.segments).Validate("abcabc", "abcdef")
			if test.valid {
				assert.Assert(t, problem == nil, "expected valid, got %+v", problem)
				return
			}
			assert.Assert(t, problem != nil)
			assert.Equal(t, problem.Kind, test.wantKind)
		})
	}
}

func TestValidateNilIsValid(t *testing.T) {
	t.Parallel()
	var m *spanmap.SpanMap
	assert.Assert(t, m.Validate("abc", "abc") == nil)
}
