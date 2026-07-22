import {
    SpanMap,
    SpanMapFidelity,
    SpanMapKind,
} from "@typescript/native-preview/unstable/ast";
import assert from "node:assert";
import {
    describe,
    test,
} from "node:test";

describe("SpanMap", () => {
    const map = new SpanMap([
        { generatedStart: 2, generatedEnd: 6, originalStart: 10, originalEnd: 14, kind: SpanMapKind.Verbatim },
        { generatedStart: 8, generatedEnd: 11, originalStart: 20, originalEnd: 27, kind: SpanMapKind.Atom },
        { generatedStart: 14, generatedEnd: 18, originalStart: 30, originalEnd: 34, kind: SpanMapKind.Verbatim },
    ]);

    test("maps generated positions and ranges to original", () => {
        assert.deepEqual(map.mapPosition(4), { position: 12, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(map.mapSpan({ pos: 3, end: 5 }), { range: { pos: 11, end: 13 }, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(map.mapPosition(9), { position: 20, fidelity: SpanMapFidelity.Atom });
        assert.deepEqual(map.mapSpan({ pos: 8, end: 10 }), { range: { pos: 20, end: 27 }, fidelity: SpanMapFidelity.Atom });
        assert.deepEqual(map.mapSpan({ pos: 5, end: 15 }), { range: { pos: 13, end: 31 }, fidelity: SpanMapFidelity.Approximate });
    });

    test("maps synthesized gaps to insertion points", () => {
        assert.deepEqual(map.mapPosition(0), { position: 0, fidelity: SpanMapFidelity.None });
        assert.deepEqual(map.mapSpan({ pos: 6, end: 8 }), { range: { pos: 14, end: 14 }, fidelity: SpanMapFidelity.None });
        assert.deepEqual(map.mapPosition(19), { position: 34, fidelity: SpanMapFidelity.None });
    });

    test("maps original positions and ranges to generated", () => {
        assert.deepEqual(map.mapPositionToGenerated(12), { position: 4, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(map.mapRangeToGenerated({ pos: 21, end: 25 }), { range: { pos: 8, end: 11 }, fidelity: SpanMapFidelity.Atom });
        assert.deepEqual(map.mapRangeToGenerated({ pos: 13, end: 31 }), { range: { pos: 5, end: 15 }, fidelity: SpanMapFidelity.Approximate });
        assert.deepEqual(map.mapPositionToGenerated(15), { position: 6, fidelity: SpanMapFidelity.None });
    });

    test("sorts generated and original indexes independently", () => {
        const reordered = new SpanMap([
            { generatedStart: 0, generatedEnd: 2, originalStart: 10, originalEnd: 12, kind: SpanMapKind.Verbatim },
            { generatedStart: 2, generatedEnd: 4, originalStart: 0, originalEnd: 2, kind: SpanMapKind.Verbatim },
        ]);
        assert.deepEqual(reordered.mapPosition(3), { position: 1, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(reordered.mapPositionToGenerated(1), { position: 3, fidelity: SpanMapFidelity.Exact });
    });

    test("an empty map describes fully synthesized output", () => {
        const empty = new SpanMap([]);
        assert.deepEqual(empty.mapPosition(5), { position: 0, fidelity: SpanMapFidelity.None });
        assert.deepEqual(empty.mapSpan({ pos: 2, end: 7 }), { range: { pos: 0, end: 0 }, fidelity: SpanMapFidelity.None });
        assert.deepEqual(empty.mapPositionToGenerated(5), { position: 0, fidelity: SpanMapFidelity.None });
    });

    test("exposes fidelity predicates", () => {
        assert.equal(SpanMap.isExact(SpanMapFidelity.Exact), true);
        assert.equal(SpanMap.isSingleSegment(SpanMapFidelity.Exact), true);
        assert.equal(SpanMap.isSingleSegment(SpanMapFidelity.Atom), true);
        assert.equal(SpanMap.isSingleSegment(SpanMapFidelity.Approximate), false);
        assert.equal(SpanMap.isNone(SpanMapFidelity.None), true);
    });
});
