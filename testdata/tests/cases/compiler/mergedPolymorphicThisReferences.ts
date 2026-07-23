// @noEmit: true
// @noTypesAndSymbols: true
// @strict: true

// @filename: class.ts
declare class MergedBase {}

// @filename: interface.ts
type Owned<T> = { owner: T };

interface MergedBase {
    clone(): this;
    isOwned(): this is this & Owned<this>;
    queried(value: this, box: { value: typeof value }): void;
    conditional(): this extends MergedDerived ? "derived" : "base";
}

declare class MergedDerived extends MergedBase {
    derived: true;
}

declare const base: MergedBase;
declare const derived: MergedDerived;

const cloned: MergedDerived = derived.clone();
const conditionalValue: "derived" = derived.conditional();

// @ts-expect-error `typeof value` must be instantiated from MergedBase to MergedDerived.
derived.queried(derived, { value: base });

if (derived.isOwned()) {
    const owner: MergedDerived = derived.owner;
}
