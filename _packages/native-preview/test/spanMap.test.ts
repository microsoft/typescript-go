import {
    SpanMap,
    SpanMapFeature,
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
        assert.equal(map.segments[0].features, SpanMapFeature.All);
        assert.deepEqual(map.generatedToOriginalPosition(4), { position: 12, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(map.generatedToOriginalSpan({ pos: 3, end: 5 }), { range: { pos: 11, end: 13 }, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(map.generatedToOriginalPosition(9), { position: 20, fidelity: SpanMapFidelity.Atom });
        assert.deepEqual(map.generatedToOriginalSpan({ pos: 8, end: 10 }), { range: { pos: 20, end: 27 }, fidelity: SpanMapFidelity.Atom });
        assert.deepEqual(map.generatedToOriginalSpan({ pos: 5, end: 15 }), { range: { pos: 13, end: 31 }, fidelity: SpanMapFidelity.Approximate });
    });

    test("maps aliases with atom geometry", () => {
        const alias = new SpanMap([
            { generatedStart: 0, generatedEnd: 3, originalStart: 0, originalEnd: 1, kind: SpanMapKind.Alias },
        ]);
        assert.deepEqual(alias.generatedToOriginalSpan({ pos: 0, end: 3 }), {
            range: { pos: 0, end: 1 },
            fidelity: SpanMapFidelity.Atom,
        });
    });

    test("maps synthesized gaps to insertion points", () => {
        assert.deepEqual(map.generatedToOriginalPosition(0), { position: 0, fidelity: SpanMapFidelity.None });
        assert.deepEqual(map.generatedToOriginalSpan({ pos: 6, end: 8 }), { range: { pos: 14, end: 14 }, fidelity: SpanMapFidelity.None });
        assert.deepEqual(map.generatedToOriginalPosition(19), { position: 34, fidelity: SpanMapFidelity.None });
    });

    test("maps original positions and ranges to generated", () => {
        assert.deepEqual(map.originalToGeneratedPositions(12, SpanMapFeature.All), [{ position: 4, fidelity: SpanMapFidelity.Exact }]);
        assert.deepEqual(map.originalToGeneratedSpans({ pos: 21, end: 25 }, SpanMapFeature.All), [{ range: { pos: 8, end: 11 }, fidelity: SpanMapFidelity.Atom }]);
        assert.deepEqual(map.originalToGeneratedSpans({ pos: 13, end: 31 }, SpanMapFeature.All), [{ range: { pos: 5, end: 15 }, fidelity: SpanMapFidelity.Approximate }]);
        assert.deepEqual(map.originalToGeneratedPositions(15, SpanMapFeature.All), []);
    });

    test("maps segment endpoints", () => {
        assert.deepEqual(map.generatedToOriginalPosition(18), { position: 34, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(map.originalToGeneratedPositions(34, SpanMapFeature.All), [{ position: 18, fidelity: SpanMapFidelity.Exact }]);
        assert.deepEqual(map.generatedToOriginalPosition(6), { position: 14, fidelity: SpanMapFidelity.None });
        assert.deepEqual(map.originalToGeneratedPositions(14, SpanMapFeature.All), [{ position: 6, fidelity: SpanMapFidelity.Exact }]);

        const adjacent = new SpanMap([
            { generatedStart: 20, generatedEnd: 23, originalStart: 10, originalEnd: 13, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Hover },
            { generatedStart: 2, generatedEnd: 5, originalStart: 13, originalEnd: 16, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Completion },
            { generatedStart: 30, generatedEnd: 33, originalStart: 20, originalEnd: 25, kind: SpanMapKind.Atom },
        ]);
        assert.deepEqual(adjacent.originalToGeneratedPositions(13, SpanMapFeature.All), [
            { position: 2, fidelity: SpanMapFidelity.Exact },
            { position: 23, fidelity: SpanMapFidelity.Exact },
        ]);
        assert.deepEqual(adjacent.originalToGeneratedPositions(13, SpanMapFeature.Completion), [
            { position: 2, fidelity: SpanMapFidelity.Exact },
        ]);
        assert.deepEqual(adjacent.originalToGeneratedPositions(25, SpanMapFeature.All), [
            { position: 33, fidelity: SpanMapFidelity.Atom },
        ]);
    });

    test("sorts generated and original indexes independently", () => {
        const reordered = new SpanMap([
            { generatedStart: 0, generatedEnd: 2, originalStart: 10, originalEnd: 12, kind: SpanMapKind.Verbatim },
            { generatedStart: 2, generatedEnd: 4, originalStart: 0, originalEnd: 2, kind: SpanMapKind.Verbatim },
        ]);
        assert.deepEqual(reordered.generatedToOriginalPosition(3), { position: 1, fidelity: SpanMapFidelity.Exact });
        assert.deepEqual(reordered.originalToGeneratedPositions(1, SpanMapFeature.All), [{ position: 3, fidelity: SpanMapFidelity.Exact }]);
    });

    test("an empty map describes fully synthesized output", () => {
        const empty = new SpanMap([]);
        assert.deepEqual(empty.generatedToOriginalPosition(5), { position: 0, fidelity: SpanMapFidelity.None });
        assert.deepEqual(empty.generatedToOriginalSpan({ pos: 2, end: 7 }), { range: { pos: 0, end: 0 }, fidelity: SpanMapFidelity.None });
        assert.deepEqual(empty.originalToGeneratedPositions(5, SpanMapFeature.All), []);
    });

    test("maps duplicate groups by features", () => {
        const duplicates = new SpanMap([
            { generatedStart: 0, generatedEnd: 3, originalStart: 10, originalEnd: 13, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Definition },
            { generatedStart: 10, generatedEnd: 13, originalStart: 10, originalEnd: 13, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Hover },
            { generatedStart: 14, generatedEnd: 17, originalStart: 10, originalEnd: 13, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Hover },
            { generatedStart: 20, generatedEnd: 25, originalStart: 10, originalEnd: 13, kind: SpanMapKind.Atom, features: SpanMapFeature.Definition },
        ]);

        assert.deepEqual(duplicates.originalToGeneratedPositions(11, SpanMapFeature.Hover), [
            { position: 11, fidelity: SpanMapFidelity.Exact },
            { position: 15, fidelity: SpanMapFidelity.Exact },
        ]);
        assert.deepEqual(duplicates.originalToGeneratedPositions(11, SpanMapFeature.Definition), [
            { position: 1, fidelity: SpanMapFidelity.Exact },
            { position: 20, fidelity: SpanMapFidelity.Atom },
        ]);
        assert.deepEqual(duplicates.originalToGeneratedPositions(13, SpanMapFeature.Hover), [
            { position: 13, fidelity: SpanMapFidelity.Exact },
            { position: 17, fidelity: SpanMapFidelity.Exact },
        ]);
    });

    test("maps minimal cross-group projections", () => {
        const projections = new SpanMap([
            { generatedStart: 0, generatedEnd: 2, originalStart: 0, originalEnd: 2, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Hover },
            { generatedStart: 2, generatedEnd: 4, originalStart: 2, originalEnd: 4, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Hover },
            { generatedStart: 10, generatedEnd: 12, originalStart: 0, originalEnd: 2, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Hover },
            { generatedStart: 12, generatedEnd: 14, originalStart: 2, originalEnd: 4, kind: SpanMapKind.Verbatim, features: SpanMapFeature.Hover },
        ]);

        assert.deepEqual(projections.originalToGeneratedSpans({ pos: 1, end: 3 }, SpanMapFeature.Hover), [
            { range: { pos: 1, end: 3 }, fidelity: SpanMapFidelity.Approximate },
            { range: { pos: 11, end: 13 }, fidelity: SpanMapFidelity.Approximate },
        ]);
    });

    test("explicit zero features disables original-to-generated mapping", () => {
        const disabled = new SpanMap([
            { generatedStart: 0, generatedEnd: 3, originalStart: 10, originalEnd: 13, kind: SpanMapKind.Verbatim, features: SpanMapFeature.None },
        ]);

        assert.deepEqual(disabled.originalToGeneratedPositions(11, SpanMapFeature.Hover), []);
        assert.deepEqual(disabled.originalToGeneratedPositions(11, SpanMapFeature.Definition), []);
        assert.deepEqual(disabled.originalToGeneratedSpans({ pos: 10, end: 13 }, SpanMapFeature.Hover), []);
    });

    test("exposes fidelity predicates", () => {
        assert.equal(SpanMap.isExact(SpanMapFidelity.Exact), true);
        assert.equal(SpanMap.isSingleSegment(SpanMapFidelity.Exact), true);
        assert.equal(SpanMap.isSingleSegment(SpanMapFidelity.Atom), true);
        assert.equal(SpanMap.isSingleSegment(SpanMapFidelity.Approximate), false);
        assert.equal(SpanMap.isNone(SpanMapFidelity.None), true);
    });
});
