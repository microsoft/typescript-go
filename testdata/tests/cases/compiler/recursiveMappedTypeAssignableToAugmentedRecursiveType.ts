// @strict: true
// @noEmit: true

type Transform<T> = {
    [K in keyof T]: T[K] extends object ? Transform<T[K]> : T[K];
};

type Child = ChildMap[keyof ChildMap];

interface ChildMap { root: Root; }
interface Root { child: Child; }

interface ChildMap { element: XMLElement; }
interface XMLElement { parent: Root; fn(): void; }

declare const root: Transform<Transform<Root>>;
const result: Root = root;
