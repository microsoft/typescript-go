// @noEmit: true
// @strict: true

// Regression test for https://github.com/microsoft/TypeScript/issues/63441
// The pass-through methods on the recovery-boundary wrapping tracker
// (ReportInferenceFallback, ReportTruncationError, ReportNonlocalAugmentation,
// Push/PopErrorFallbackNode) used to delegate through every nested
// SymbolTrackerImpl/wrappingTracker pair, adding two stack frames per
// nested tryReuseExistingNodeHelper boundary. For deeply recursive types
// (like the one below) this overflowed the goroutine stack.

export interface CustomNode<P> {
    getNextNode: () => CustomNode<P>;
}

export declare const createNode: () => {
    getNextNode: <T>() => CustomNode<T>;
};

function wrapNode<T>(getNode: () => CustomNode<T>) {
    return getNode;
}

wrapNode(() => {
    const node = createNode();

    return wrapNode<typeof node.getNextNode<any>>(node.getNextNode);
});
