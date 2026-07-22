import { SpanMapFidelity } from "#enums/spanMapFidelity";
import { SpanMapKind } from "#enums/spanMapKind";
import type { ReadonlyTextRange } from "./ast.ts";

export { SpanMapFidelity, SpanMapKind };

export interface SpanMapSegment {
    readonly generatedStart: number;
    readonly generatedEnd: number;
    readonly originalStart: number;
    readonly originalEnd: number;
    readonly kind: SpanMapKind;
}

export interface MappedPosition {
    readonly position: number;
    readonly fidelity: SpanMapFidelity;
}

export interface MappedRange {
    readonly range: ReadonlyTextRange;
    readonly fidelity: SpanMapFidelity;
}

export class SpanMap {
    readonly segments: readonly SpanMapSegment[];
    private originalSegments: readonly SpanMapSegment[] | undefined;

    constructor(segments: readonly SpanMapSegment[]) {
        this.segments = [...segments].sort((left, right) => left.generatedStart - right.generatedStart);
    }

    static isExact(fidelity: SpanMapFidelity): boolean {
        return fidelity === SpanMapFidelity.Exact;
    }

    static isSingleSegment(fidelity: SpanMapFidelity): boolean {
        return fidelity === SpanMapFidelity.Exact || fidelity === SpanMapFidelity.Atom;
    }

    static isNone(fidelity: SpanMapFidelity): boolean {
        return fidelity === SpanMapFidelity.None;
    }

    mapSpan(range: ReadonlyTextRange): MappedRange {
        return this.mapRange(range, this.segments, false);
    }

    mapPosition(position: number): MappedPosition {
        return this.mapPoint(position, this.segments, false);
    }

    mapRangeToGenerated(range: ReadonlyTextRange): MappedRange {
        return this.mapRange(range, this.getOriginalSegments(), true);
    }

    mapPositionToGenerated(position: number): MappedPosition {
        return this.mapPoint(position, this.getOriginalSegments(), true);
    }

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

    private getOriginalSegments(): readonly SpanMapSegment[] {
        return this.originalSegments ??= [...this.segments].sort((left, right) => left.originalStart - right.originalStart);
    }
}

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

function insertionPoint(segments: readonly SpanMapSegment[], previous: number, reverse: boolean): number {
    if (previous < 0) return 0;
    return reverse ? segments[previous].generatedEnd : segments[previous].originalEnd;
}

function mapVerbatimPosition(segment: SpanMapSegment, position: number, reverse: boolean): number {
    const sourceStart = reverse ? segment.originalStart : segment.generatedStart;
    const targetStart = reverse ? segment.generatedStart : segment.originalStart;
    const targetEnd = reverse ? segment.generatedEnd : segment.originalEnd;
    return clamp(targetStart + position - sourceStart, targetStart, targetEnd);
}

function mapBoundary(segments: readonly SpanMapSegment[], position: number, index: number, inside: boolean, reverse: boolean, high: boolean): number {
    if (!inside) return insertionPoint(segments, index, reverse);
    const segment = segments[index];
    if (segment.kind === SpanMapKind.Verbatim) return mapVerbatimPosition(segment, position, reverse);
    if (reverse) return high ? segment.generatedEnd : segment.generatedStart;
    return high ? segment.originalEnd : segment.originalStart;
}

function targetRange(segment: SpanMapSegment, reverse: boolean): ReadonlyTextRange {
    return reverse
        ? { pos: segment.generatedStart, end: segment.generatedEnd }
        : { pos: segment.originalStart, end: segment.originalEnd };
}

function clamp(value: number, low: number, high: number): number {
    return Math.max(low, Math.min(value, high));
}
