import { SpanMapFidelity } from "#enums/spanMapFidelity";
import { SpanMapKind } from "#enums/spanMapKind";
import { SpanMapPurpose } from "#enums/spanMapPurpose";
import type { ReadonlyTextRange } from "./ast.ts";

export { SpanMapFidelity, SpanMapKind, SpanMapPurpose };

// Keep this in sync with spanmap.go

/** Maps one half-open generated range to one half-open original range. */
export interface SpanMapSegment {
    readonly generatedStart: number;
    readonly generatedEnd: number;
    readonly originalStart: number;
    readonly originalEnd: number;
    readonly kind: SpanMapKind;
    readonly purpose?: SpanMapPurpose;
}

/** Internal segment representation after an omitted purpose has been normalized to `All`. */
type NormalizedSpanMapSegment = SpanMapSegment & { readonly purpose: SpanMapPurpose; };

/** One generated projection of an original position and its mapping fidelity. */
export interface MappedPosition {
    readonly position: number;
    readonly fidelity: SpanMapFidelity;
}

/** One generated projection of an original range and its mapping fidelity. */
export interface MappedRange {
    readonly range: ReadonlyTextRange;
    readonly fidelity: SpanMapFidelity;
}

/** Provides bidirectional span-aware mapping between generated and original text. */
export class SpanMap {
    readonly segments: readonly NormalizedSpanMapSegment[];
    private originalSegments: readonly NormalizedSpanMapSegment[] | undefined;

    /** Copies and sorts segments by generated start, normalizing omitted purposes to `All`. */
    constructor(segments: readonly SpanMapSegment[]) {
        this.segments = segments
            .map(segment => ({ ...segment, purpose: segment.purpose ?? SpanMapPurpose.All }))
            .sort((left, right) => left.generatedStart - right.generatedStart);
    }

    /** Reports whether a mapping is a precise, edit-safe projection through one verbatim segment. */
    static isExact(fidelity: SpanMapFidelity): boolean {
        return fidelity === SpanMapFidelity.Exact;
    }

    /** Reports whether a mapping lies in one verbatim or atom segment. */
    static isSingleSegment(fidelity: SpanMapFidelity): boolean {
        return fidelity === SpanMapFidelity.Exact || fidelity === SpanMapFidelity.Atom;
    }

    /** Reports whether the input had no counterpart in the target text. */
    static isNone(fidelity: SpanMapFidelity): boolean {
        return fidelity === SpanMapFidelity.None;
    }

    /**
     * Maps a generated range to original text. Gaps map to insertion points with `None` fidelity,
     * and ranges crossing segment boundaries map their endpoints with `Approximate` fidelity.
     */
    generatedToOriginalSpan(range: ReadonlyTextRange): MappedRange {
        return this.mapRange(range, this.segments, false);
    }

    /** Maps a generated position to original text, using `None` fidelity for synthesized gaps. */
    generatedToOriginalPosition(position: number): MappedPosition {
        return this.mapPoint(position, this.segments, false);
    }

    /**
     * Returns every generated projection of an original position whose segment supports `purpose`.
     * Results are ordered by generated start; uncovered or disabled positions produce no results.
     */
    originalToGeneratedPositions(position: number, purpose: SpanMapPurpose): readonly MappedPosition[] {
        const segments = originalSegmentsAt(this.getOriginalSegments(), position);
        if (!segments) return [];
        return segments
            .filter(segment => supportsPurpose(segment, purpose))
            .map(segment =>
                segment.kind === SpanMapKind.Verbatim
                    ? { position: mapVerbatimPosition(segment, position, true), fidelity: SpanMapFidelity.Exact }
                    : { position: segment.generatedStart, fidelity: SpanMapFidelity.Atom }
            );
    }

    /**
     * Returns every purpose-compatible generated projection of an original range.
     * A range contained by one duplicate group produces one exact or atom result per matching group member.
     *
     * A range that starts in one group and ends in another can have several possible generated ranges. For
     * example, suppose two original segments are each copied twice into the generated text:
     *
     * ```text
     * original:   [ A ][ B ]
     *                [---)       range from inside A to inside B
     *
     * generated:  [ A ][ B ]      [ A ][ B ]
     *                ^   ^          ^   ^
     *              start end      start end
     *                1   3          11  13
     * ```
     *
     * The map says that the range may start at 1 or 11 and end at 3 or 13, but it does not say which copy of A
     * belongs with which copy of B. We choose the smallest range around each possible location, producing [1,3)
     * and [11,13). We do not return [1,13), because it contains both smaller candidates and would include code
     * that may be unrelated to the original range. These cross-group results have approximate fidelity.
     */
    originalToGeneratedSpans(range: ReadonlyTextRange, purpose: SpanMapPurpose): readonly MappedRange[] {
        const results = this.originalToGeneratedSpansInGroup(range, purpose);
        if (results !== undefined) return results;
        const start = range.pos;
        const end = Math.max(range.end, start);
        const starts = this.originalStartProjections(start, purpose);
        const ends = this.originalEndProjections(start, end, purpose);
        if (starts.length === 0 || ends.length === 0) return [];
        return starts.flatMap((generatedStart, index) => {
            const generatedEnd = ends.find(end => end >= generatedStart);
            return generatedEnd === undefined || index + 1 < starts.length && starts[index + 1] <= generatedEnd
                ? []
                : [{ range: { pos: generatedStart, end: generatedEnd }, fidelity: SpanMapFidelity.Approximate }];
        });
    }

    /**
     * Maps the inclusive start of an original range through every matching segment. Verbatim segments preserve
     * the offset within the segment; atoms map to their generated start.
     *
     * ```text
     * original:       [---------)
     *                    ^ start
     *
     * generated:  [---------)   [---------)
     *                ^             ^
     *              result        result
     * ```
     */
    private originalStartProjections(start: number, purpose: SpanMapPurpose): readonly number[] {
        const segments = originalSegmentsAt(this.getOriginalSegments(), start);
        if (!segments) return [];
        return segments
            .filter(segment => supportsPurpose(segment, purpose))
            .map(segment =>
                segment.kind === SpanMapKind.Verbatim
                    ? mapVerbatimPosition(segment, start, true)
                    : segment.generatedStart
            );
    }

    /**
     * Maps the exclusive end of an original range through every matching segment. For a non-empty range,
     * `end - 1` identifies the segment containing its final character, while `end` remains the mapped boundary.
     *
     * ```text
     * original:       [---------)[ next segment )
     *                          ^`-- end
     *                          `--- end - 1
     *
     * generated:  [---------)   [---------)
     *                       ^             ^
     *                     result        result
     * ```
     */
    private originalEndProjections(start: number, end: number, purpose: SpanMapPurpose): readonly number[] {
        const lastCharacter = end > start ? end - 1 : end;
        const segments = originalSegmentsAt(this.getOriginalSegments(), lastCharacter);
        if (!segments) return [];
        return segments
            .filter(segment => supportsPurpose(segment, purpose))
            .map(segment =>
                segment.kind === SpanMapKind.Verbatim
                    ? mapVerbatimPosition(segment, end, true)
                    : segment.generatedEnd
            );
    }

    /**
     * Maps a range contained by one duplicate group. `undefined` means the range crosses groups; an empty array
     * means the range is in one group but no member supports `purpose`.
     */
    private originalToGeneratedSpansInGroup(range: ReadonlyTextRange, purpose: SpanMapPurpose): readonly MappedRange[] | undefined {
        const start = range.pos;
        const end = Math.max(range.end, start);
        const startSegments = originalSegmentsAt(this.getOriginalSegments(), start);
        const endSegments = originalSegmentsAt(this.getOriginalSegments(), end > start ? end - 1 : end);
        if (!startSegments || !endSegments || startSegments[0].originalStart !== endSegments[0].originalStart || startSegments[0].originalEnd !== endSegments[0].originalEnd) {
            return undefined;
        }
        return startSegments
            .filter(segment => supportsPurpose(segment, purpose))
            .map(segment => {
                if (segment.kind === SpanMapKind.Verbatim) {
                    const mappedStart = mapVerbatimPosition(segment, start, true);
                    const mappedEnd = Math.max(mappedStart, mapVerbatimPosition(segment, end, true));
                    return { range: { pos: mappedStart, end: mappedEnd }, fidelity: SpanMapFidelity.Exact };
                }
                return { range: { pos: segment.generatedStart, end: segment.generatedEnd }, fidelity: SpanMapFidelity.Atom };
            });
    }

    /** Maps one range through an ordered segment index in the direction selected by `reverse`. */
    private mapRange(range: ReadonlyTextRange, segments: readonly SpanMapSegment[], reverse: boolean): MappedRange {
        const start = range.pos;
        const end = Math.max(range.end, start);
        const [startIndex, startInside] = segmentIndexAt(segments, start, reverse);
        const endProbe = end > start ? end - 1 : end;
        const [endIndex, endInside] = segmentIndexAt(segments, endProbe, reverse);

        if (startIndex === endIndex && startInside === endInside) {
            if (startInside) {
                const segment = segments[startIndex];
                if (segment.kind === SpanMapKind.Verbatim) {
                    const mappedStart = mapVerbatimPosition(segment, start, reverse);
                    const mappedEnd = Math.max(mappedStart, mapVerbatimPosition(segment, end, reverse));
                    return { range: { pos: mappedStart, end: mappedEnd }, fidelity: SpanMapFidelity.Exact };
                }
                return { range: targetRange(segment, reverse), fidelity: SpanMapFidelity.Atom };
            }
            const position = insertionPoint(segments, startIndex, reverse);
            return { range: { pos: position, end: position }, fidelity: SpanMapFidelity.None };
        }

        const mappedStart = mapBoundary(segments, start, startIndex, startInside, reverse, false);
        const mappedEnd = Math.max(mappedStart, mapBoundary(segments, end, endIndex, endInside, reverse, true));
        return { range: { pos: mappedStart, end: mappedEnd }, fidelity: SpanMapFidelity.Approximate };
    }

    /** Maps one position through an ordered segment index in the direction selected by `reverse`. */
    private mapPoint(position: number, segments: readonly SpanMapSegment[], reverse: boolean): MappedPosition {
        const [index, inside] = segmentIndexAt(segments, position, reverse);
        if (!inside) {
            return { position: insertionPoint(segments, index, reverse), fidelity: SpanMapFidelity.None };
        }
        const segment = segments[index];
        if (segment.kind === SpanMapKind.Verbatim) {
            return { position: mapVerbatimPosition(segment, position, reverse), fidelity: SpanMapFidelity.Exact };
        }
        return {
            position: reverse ? segment.generatedStart : segment.originalStart,
            fidelity: SpanMapFidelity.Atom,
        };
    }

    /** Returns the lazily built segment index ordered by original start. */
    private getOriginalSegments(): readonly NormalizedSpanMapSegment[] {
        return this.originalSegments ??= [...this.segments].sort((left, right) => left.originalStart - right.originalStart);
    }
}

/**
 * Returns the complete duplicate group containing `position`. Segment ends are exclusive; starts, including
 * zero-length segment starts, are included. `segments` must be ordered by original start.
 */
function originalSegmentsAt(segments: readonly NormalizedSpanMapSegment[], position: number): readonly NormalizedSpanMapSegment[] | undefined {
    const first = segments.findIndex(segment => segment.originalStart === position || (segment.originalStart < position && position < segment.originalEnd));
    if (first < 0) return undefined;
    let end = first + 1;
    while (end < segments.length && segments[end].originalStart === segments[first].originalStart && segments[end].originalEnd === segments[first].originalEnd) end++;
    return segments.slice(first, end);
}

/** Reports whether a segment participates in an original-to-generated query for `purpose`. */
function supportsPurpose(segment: NormalizedSpanMapSegment, purpose: SpanMapPurpose): boolean {
    return (segment.purpose & purpose) !== 0;
}

/**
 * Finds the segment containing `position`, or the preceding segment when `position` is in a gap.
 * The boolean distinguishes containment from a gap; `reverse` selects original rather than generated coordinates.
 */
function segmentIndexAt(segments: readonly SpanMapSegment[], position: number, reverse: boolean): [number, boolean] {
    let low = 0;
    let high = segments.length;
    while (low < high) {
        const middle = (low + high) >>> 1;
        const start = reverse ? segments[middle].originalStart : segments[middle].generatedStart;
        if (start < position) low = middle + 1;
        else high = middle;
    }
    if (low < segments.length && (reverse ? segments[low].originalStart : segments[low].generatedStart) === position) {
        return [low, true];
    }
    const previous = low - 1;
    if (previous >= 0) {
        const end = reverse ? segments[previous].originalEnd : segments[previous].generatedEnd;
        if (position < end) return [previous, true];
    }
    return [previous, false];
}

/** Returns the target insertion point for a gap following `previous`, or zero before the first segment. */
function insertionPoint(segments: readonly SpanMapSegment[], previous: number, reverse: boolean): number {
    if (previous < 0) return 0;
    return reverse ? segments[previous].generatedEnd : segments[previous].originalEnd;
}

/** Linearly maps and clamps a position within a length-preserving verbatim segment. */
function mapVerbatimPosition(segment: SpanMapSegment, position: number, reverse: boolean): number {
    const sourceStart = reverse ? segment.originalStart : segment.generatedStart;
    const targetStart = reverse ? segment.generatedStart : segment.originalStart;
    const targetEnd = reverse ? segment.generatedEnd : segment.originalEnd;
    return clamp(targetStart + position - sourceStart, targetStart, targetEnd);
}

/** Maps a range boundary, using insertion points for gaps and the selected endpoint for atoms. */
function mapBoundary(segments: readonly SpanMapSegment[], position: number, index: number, inside: boolean, reverse: boolean, high: boolean): number {
    if (!inside) return insertionPoint(segments, index, reverse);
    const segment = segments[index];
    if (segment.kind === SpanMapKind.Verbatim) return mapVerbatimPosition(segment, position, reverse);
    if (reverse) return high ? segment.generatedEnd : segment.generatedStart;
    return high ? segment.originalEnd : segment.originalStart;
}

/** Returns the complete target range of a segment in the selected direction. */
function targetRange(segment: SpanMapSegment, reverse: boolean): ReadonlyTextRange {
    return reverse
        ? { pos: segment.generatedStart, end: segment.generatedEnd }
        : { pos: segment.originalStart, end: segment.originalEnd };
}

/** Confines `value` to the inclusive interval [`low`, `high`]. */
function clamp(value: number, low: number, high: number): number {
    return Math.max(low, Math.min(value, high));
}
