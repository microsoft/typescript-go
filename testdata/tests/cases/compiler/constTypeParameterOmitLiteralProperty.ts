// https://github.com/microsoft/typescript-go/issues/4391
// @moduleResolution: bundler
// @module: preserve
// @target: esnext
// @strict: true
// @skipLibCheck: true
// @noEmit: true
// @jsx: react-jsx
// @noTypesAndSymbols: true

// @filename: /node_modules/@tanstack/react-router/package.json
{"name":"@tanstack/react-router","type":"module","types":"dist/esm/index.d.ts","exports":{".":{"import":{"types":"./dist/esm/index.d.ts"}}}}

// @filename: /node_modules/@tanstack/router-core/package.json
{"name":"@tanstack/router-core","type":"module","types":"dist/esm/index.d.ts","exports":{".":{"import":{"types":"./dist/esm/index.d.ts"}}}}

// @filename: /node_modules/@tanstack/history/package.json
{"name":"@tanstack/history","type":"module","types":"dist/esm/index.d.ts","exports":{".":{"import":{"types":"./dist/esm/index.d.ts"}}}}

// @filename: /node_modules/react/package.json
{"name":"react","type":"module","types":"index.d.ts","exports":{".":{"import":{"types":"./index.d.ts"}}}}

// @filename: /node_modules/@tanstack/history/dist/esm/index.d.ts
export interface NavigateOptions {
    ignoreBlocker?: boolean;
}
type SubscriberHistoryAction = {
    type: Exclude<HistoryAction, 'GO'>;
} | {
    type: 'GO';
    index: number;
};
type SubscriberArgs = {
    location: HistoryLocation;
    action: SubscriberHistoryAction;
};
export interface RouterHistory {
    location: HistoryLocation;
    length: number;
    subscribers: Set<(opts: SubscriberArgs) => void>;
    subscribe: (cb: (opts: SubscriberArgs) => void) => () => void;
    push: (path: string, state?: any, navigateOpts?: NavigateOptions) => void;
    replace: (path: string, state?: any, navigateOpts?: NavigateOptions) => void;
    go: (index: number, navigateOpts?: NavigateOptions) => void;
    back: (navigateOpts?: NavigateOptions) => void;
    forward: (navigateOpts?: NavigateOptions) => void;
    canGoBack: () => boolean;
    createHref: (href: string) => string;
    block: (blocker: NavigationBlocker) => () => void;
    flush: () => void;
    destroy: () => void;
    notify: (action: SubscriberHistoryAction) => void;
    _ignoreSubscribers?: boolean;
}
export interface HistoryLocation extends ParsedPath {
    state: ParsedHistoryState;
}
export interface ParsedPath {
    href: string;
    pathname: string;
    search: string;
    hash: string;
}
export interface HistoryState {
}
export type ParsedHistoryState = HistoryState & {
    key?: string;
    __TSR_key?: string;
    __TSR_index: number;
};
type ShouldAllowNavigation = any;
export type HistoryAction = 'PUSH' | 'REPLACE' | 'FORWARD' | 'BACK' | 'GO';
export type BlockerFnArgs = {
    currentLocation: HistoryLocation;
    nextLocation: HistoryLocation;
    action: HistoryAction;
};
export type BlockerFn = (args: BlockerFnArgs) => Promise<ShouldAllowNavigation> | ShouldAllowNavigation;
export type NavigationBlocker = {
    blockerFn: BlockerFn;
    enableBeforeUnload?: (() => boolean) | boolean;
};
export declare function createHistory(opts: {
    getLocation: () => HistoryLocation;
    getLength: () => number;
    pushState: (path: string, state: any) => void;
    replaceState: (path: string, state: any) => void;
    go: (n: number) => void;
    back: (ignoreBlocker: boolean) => void;
    forward: (ignoreBlocker: boolean) => void;
    createHref: (path: string) => string;
    flush?: () => void;
    destroy?: () => void;
    onBlocked?: () => void;
    getBlockers?: () => Array<NavigationBlocker>;
    setBlockers?: (blockers: Array<NavigationBlocker>) => void;
    notifyOnIndexChange?: boolean;
}): RouterHistory;
/**
 * Creates a history object that can be used to interact with the browser's
 * navigation. This is a lightweight API wrapping the browser's native methods.
 * It is designed to work with TanStack Router, but could be used as a standalone API as well.
 * IMPORTANT: This API implements history throttling via a microtask to prevent
 * excessive calls to the history API. In some browsers, calling history.pushState or
 * history.replaceState in quick succession can cause the browser to ignore subsequent
 * calls. This API smooths out those differences and ensures that your application
 * state will *eventually* match the browser state. In most cases, this is not a problem,
 * but if you need to ensure that the browser state is up to date, you can use the
 * `history.flush` method to immediately flush all pending state changes to the browser URL.
 * @param opts
 * @param opts.getHref A function that returns the current href (path + search + hash)
 * @param opts.createHref A function that takes a path and returns a href (path + search + hash)
 * @returns A history instance
 */
export declare function createBrowserHistory(opts?: {
    parseLocation?: () => HistoryLocation;
    createHref?: (path: string) => string;
    window?: any;
}): RouterHistory;
/**
 * Create a hash-based history implementation.
 * Useful for static hosts or environments without server URL rewriting.
 * @link https://tanstack.com/router/latest/docs/framework/react/guide/history-types
 */
export declare function createHashHistory(opts?: {
    window?: any;
}): RouterHistory;
/**
 * Create an in-memory history implementation.
 * Ideal for server rendering, tests, and non-DOM environments.
 * @link https://tanstack.com/router/latest/docs/framework/react/guide/history-types
 */
export declare function createMemoryHistory(opts?: {
    initialEntries: Array<string>;
    initialIndex?: number;
}): RouterHistory;
export declare function parseHref(href: string, state: ParsedHistoryState | undefined): HistoryLocation;
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/fileRoute.d.ts
import { Register } from './router.js';
import { AnyContext, AnyPathParams, AnyRoute, FileBaseRouteOptions, ResolveParams, Route, RouteConstraints, UpdatableRouteOptions } from './route.js';
import { AnyValidator } from './validators.js';
export interface FileRouteTypes {
    fileRoutesByFullPath: any;
    fullPaths: any;
    to: any;
    fileRoutesByTo: any;
    id: any;
    fileRoutesById: any;
}
export type InferFileRouteTypes<TRouteTree extends AnyRoute> = unknown extends TRouteTree['types']['fileRouteTypes'] ? never : TRouteTree['types']['fileRouteTypes'] extends FileRouteTypes ? TRouteTree['types']['fileRouteTypes'] : never;
export interface FileRoutesByPath {
}
export interface FileRouteOptions<TRegister, TFilePath extends string, TParentRoute extends AnyRoute, TId extends RouteConstraints['TId'], TPath extends RouteConstraints['TPath'], TFullPath extends RouteConstraints['TFullPath'], TSearchValidator = undefined, TParams = ResolveParams<TPath>, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TSSR = unknown, TServerMiddlewares = unknown, THandlers = undefined> extends FileBaseRouteOptions<TRegister, TParentRoute, TId, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, AnyContext, TRouteContextFn, TBeforeLoadFn, AnyContext, TSSR, TServerMiddlewares, THandlers>, UpdatableRouteOptions<TParentRoute, TId, TFullPath, TParams, TSearchValidator, TLoaderFn, TLoaderDeps, AnyContext, TRouteContextFn, TBeforeLoadFn> {
}
export type CreateFileRoute<TFilePath extends string, TParentRoute extends AnyRoute, TId extends RouteConstraints['TId'], TPath extends RouteConstraints['TPath'], TFullPath extends RouteConstraints['TFullPath']> = <TRegister = Register, TSearchValidator = undefined, TParams = ResolveParams<TPath>, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TSSR = unknown, TServerMiddlewares = unknown, THandlers = undefined>(options?: FileRouteOptions<TRegister, TFilePath, TParentRoute, TId, TPath, TFullPath, TSearchValidator, TParams, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TSSR, TServerMiddlewares, THandlers>) => Route<TRegister, TParentRoute, TPath, TFullPath, TFilePath, TId, TSearchValidator, TParams, AnyContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, unknown, unknown, TSSR, TServerMiddlewares, THandlers>;
export type LazyRouteOptions = Pick<UpdatableRouteOptions<AnyRoute, string, string, AnyPathParams, AnyValidator, {}, AnyContext, AnyContext, AnyContext, AnyContext>, 'component' | 'errorComponent' | 'pendingComponent' | 'notFoundComponent'>;
export interface LazyRoute<in out TRoute extends AnyRoute> {
    options: {
        id: string;
    } & LazyRouteOptions;
}
export type CreateLazyFileRoute<TRoute extends AnyRoute> = (opts: LazyRouteOptions) => LazyRoute<TRoute>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/utils.d.ts
import { RouteIds } from './routeInfo.js';
import { AnyRouter } from './router.js';
export type Awaitable<T> = T | Promise<T>;
export type NoInfer<T> = [T][T extends any ? 0 : never];
export type IsAny<TValue, TYesResult, TNoResult = TValue> = 1 extends 0 & TValue ? TYesResult : TNoResult;
export type PickAsRequired<TValue, TKey extends keyof TValue> = Omit<TValue, TKey> & Required<Pick<TValue, TKey>>;
export type PickRequired<T> = {
    [K in keyof T as undefined extends T[K] ? never : K]: T[K];
};
export type PickOptional<T> = {
    [K in keyof T as undefined extends T[K] ? K : never]: T[K];
};
export type WithoutEmpty<T> = T extends any ? ({} extends T ? never : T) : never;
export type Expand<T> = T extends object ? T extends infer O ? O extends Function ? O : {
    [K in keyof O]: O[K];
} : never : T;
export type DeepPartial<T> = T extends object ? {
    [P in keyof T]?: DeepPartial<T[P]>;
} : T;
export type MakeDifferenceOptional<TLeft, TRight> = keyof TLeft & keyof TRight extends never ? TRight : Omit<TRight, keyof TLeft & keyof TRight> & {
    [K in keyof TLeft & keyof TRight]?: TRight[K];
};
export type IsUnion<T, U extends T = T> = (T extends any ? (U extends T ? false : true) : never) extends false ? false : true;
export type IsNonEmptyObject<T> = T extends object ? keyof T extends never ? false : true : false;
export type Assign<TLeft, TRight> = TLeft extends any ? TRight extends any ? IsNonEmptyObject<TLeft> extends false ? TRight : IsNonEmptyObject<TRight> extends false ? TLeft : keyof TLeft & keyof TRight extends never ? TLeft & TRight : Omit<TLeft, keyof TRight> & TRight : never : never;
export type IntersectAssign<TLeft, TRight> = TLeft extends any ? TRight extends any ? IsNonEmptyObject<TLeft> extends false ? TRight : IsNonEmptyObject<TRight> extends false ? TLeft : TRight & TLeft : never : never;
export type Timeout = ReturnType<typeof setTimeout>;
export type Updater<TPrevious, TResult = TPrevious> = TResult | ((prev?: TPrevious) => TResult);
export type NonNullableUpdater<TPrevious, TResult = TPrevious> = TResult | ((prev: TPrevious) => TResult);
export type ExtractObjects<TUnion> = TUnion extends MergeAllPrimitive ? never : TUnion;
export type PartialMergeAllObject<TUnion> = ExtractObjects<TUnion> extends infer TObj ? [TObj] extends [never] ? never : {
    [TKey in TObj extends any ? keyof TObj : never]?: TObj extends any ? TKey extends keyof TObj ? TObj[TKey] : never : never;
} : never;
export type MergeAllPrimitive = ReadonlyArray<any> | number | string | bigint | boolean | symbol | undefined | null;
export type ExtractPrimitives<TUnion> = TUnion extends MergeAllPrimitive ? TUnion : TUnion extends object ? never : TUnion;
export type PartialMergeAll<TUnion> = ExtractPrimitives<TUnion> | PartialMergeAllObject<TUnion>;
export type Constrain<T, TConstraint, TDefault = TConstraint> = (T extends TConstraint ? T : never) | TDefault;
export type ConstrainLiteral<T, TConstraint, TDefault = TConstraint> = (T & TConstraint) | TDefault;
/**
 * To be added to router types
 */
export type UnionToIntersection<T> = (T extends any ? (arg: T) => any : never) extends (arg: infer T) => any ? T : never;
/**
 * Merges everything in a union into one object.
 * This mapped type is homomorphic which means it preserves stuff! :)
 */
export type MergeAllObjects<TUnion, TIntersected = UnionToIntersection<ExtractObjects<TUnion>>> = [keyof TIntersected] extends [never] ? never : {
    [TKey in keyof TIntersected]: TUnion extends any ? TUnion[TKey & keyof TUnion] : never;
};
export type MergeAll<TUnion> = MergeAllObjects<TUnion> | ExtractPrimitives<TUnion>;
export type ValidateJSON<T> = ((...args: Array<any>) => any) extends T ? unknown extends T ? never : 'Function is not serializable' : {
    [K in keyof T]: ValidateJSON<T[K]>;
};
export type LooseReturnType<T> = T extends (...args: Array<any>) => infer TReturn ? TReturn : never;
export type LooseAsyncReturnType<T> = T extends (...args: Array<any>) => infer TReturn ? TReturn extends Promise<infer TReturn> ? TReturn : TReturn : never;
/**
 * Return the last element of an array.
 * Intended for non-empty arrays used within router internals.
 */
export declare function last<T>(arr: ReadonlyArray<T>): T | undefined;
/**
 * Apply a value-or-updater to a previous value.
 * Accepts either a literal value or a function of the previous value.
 */
export declare function functionalUpdate<TPrevious, TResult = TPrevious>(updater: Updater<TPrevious, TResult> | NonNullableUpdater<TPrevious, TResult>, previous: TPrevious): TResult;
export declare const hasOwn: (v: PropertyKey) => boolean;
export declare function hasKeys(obj: Record<string, unknown>): boolean;
export declare const createNull: () => any;
export declare const nullReplaceEqualDeep: typeof replaceEqualDeep;
/**
 * This function returns `prev` if `_next` is deeply equal.
 * If not, it will replace any deeply equal children of `b` with those of `a`.
 * This can be used for structural sharing between immutable JSON values for example.
 * Do not use this with signals
 */
export declare function replaceEqualDeep<T>(prev: any, _next: T, _makeObj?: () => {}, _depth?: number): T;
export declare function isPlainObject(o: any): boolean;
/**
 * Check if a value is a "plain" array (no extra enumerable keys).
 */
export declare function isPlainArray(value: unknown): value is Array<unknown>;
/**
 * Perform a deep equality check with options for partial comparison and
 * ignoring `undefined` values. Optimized for router state comparisons.
 */
export declare function deepEqual(a: any, b: any, opts?: {
    partial?: boolean;
    ignoreUndefined?: boolean;
}): boolean;
export type StringLiteral<T> = T extends string ? string extends T ? string : T : never;
export type ThrowOrOptional<T, TThrow extends boolean> = TThrow extends true ? T : T | undefined;
export type StrictOrFrom<TRouter extends AnyRouter, TFrom, TStrict extends boolean = true> = TStrict extends false ? {
    from?: never;
    strict: TStrict;
} : {
    from: ConstrainLiteral<TFrom, RouteIds<TRouter['routeTree']>>;
    strict?: TStrict;
};
export type ThrowConstraint<TStrict extends boolean, TThrow extends boolean> = TStrict extends false ? (TThrow extends true ? never : TThrow) : TThrow;
export type ControlledPromise<T> = Promise<T> & {
    resolve: (value: T) => void;
    reject: (value: any) => void;
    status: 'pending' | 'resolved' | 'rejected';
    value?: T;
};
/**
 * Create a promise with exposed resolve/reject and status fields.
 * Useful for coordinating async router lifecycle operations.
 */
export declare function createControlledPromise<T>(onResolve?: (value: T) => void): ControlledPromise<T>;
/**
 * Heuristically detect dynamic import "module not found" errors
 * across major browsers for lazy route component handling.
 */
export declare function isModuleNotFoundError(error: any): boolean;
export declare function isPromise<T>(value: Promise<Awaited<T>> | T): value is Promise<Awaited<T>>;
export declare function findLast<T>(array: ReadonlyArray<T>, predicate: (item: T) => boolean): T | undefined;
/**
 * Default list of URL protocols to allow in links, redirects, and navigation.
 * Any absolute URL protocol not in this list is treated as dangerous by default.
 */
export declare const DEFAULT_PROTOCOL_ALLOWLIST: string[];
/**
 * Check if a URL string uses a protocol that is not in the allowlist.
 * Returns true for blocked protocols like javascript:, blob:, data:, etc.
 *
 * The URL constructor correctly normalizes:
 * - Mixed case (JavaScript: → javascript:)
 * - Whitespace/control characters (java\nscript: → javascript:)
 * - Leading whitespace
 *
 * For relative URLs (no protocol), returns false (safe).
 *
 * @param url - The URL string to check
 * @param allowlist - Set of protocols to allow
 * @returns true if the URL uses a protocol that is not allowed
 */
export declare function isDangerousProtocol(url: string, allowlist: Set<string>): boolean;
/**
 * Escape HTML special characters in a string to prevent XSS attacks
 * when embedding strings in script tags during SSR.
 *
 * This is essential for preventing XSS vulnerabilities when user-controlled
 * content is embedded in inline scripts.
 */
export declare function escapeHtml(str: string): string;
export declare function decodePath(path: string): {
    path: string;
    handledProtocolRelativeURL: boolean;
};
/**
 * Encodes a path the same way `new URL()` would, but without the overhead of full URL parsing.
 *
 * This function encodes:
 * - Whitespace characters (spaces → %20, tabs → %09, etc.)
 * - Non-ASCII/Unicode characters (emojis, accented characters, etc.)
 *
 * It preserves:
 * - Already percent-encoded sequences (won't double-encode %2F, %25, etc.)
 * - ASCII special characters valid in URL paths (@, $, &, +, etc.)
 * - Forward slashes as path separators
 *
 * Used to generate proper href values for SSR without constructing URL objects.
 *
 * @example
 * encodePathLikeUrl('/path/file name.pdf') // '/path/file%20name.pdf'
 * encodePathLikeUrl('/path/日本語') // '/path/%E6%97%A5%E6%9C%AC%E8%AA%9E'
 * encodePathLikeUrl('/path/already%20encoded') // '/path/already%20encoded' (preserved)
 */
export declare function encodePathLikeUrl(path: string): string;
/**
 * Builds the dev-mode CSS styles URL for route-scoped CSS collection.
 * Used by HeadContent components in all framework implementations to construct
 * the URL for the `/@tanstack-start/styles.css` endpoint.
 *
 * @param basepath - The router's basepath (may or may not have leading slash)
 * @param routeIds - Array of matched route IDs to include in the CSS collection
 * @returns The full URL path for the dev styles CSS endpoint
 */
export declare function buildDevStylesUrl(basepath: string, routeIds: Array<string>): string;
export declare function arraysEqual<T>(a: Array<T>, b: Array<T>): boolean;


// @filename: /node_modules/@tanstack/router-core/dist/esm/link.d.ts
import { HistoryState, ParsedHistoryState } from '@tanstack/history';
import { AllParams, CatchAllPaths, CurrentPath, FullSearchSchema, FullSearchSchemaInput, ParentPath, RouteByPath, RouteByToPath, RoutePaths, RouteToPath, ToPath } from './routeInfo.js';
import { AnyRouter, RegisteredRouter, ViewTransitionOptions } from './router.js';
import { ConstrainLiteral, Expand, MakeDifferenceOptional, NoInfer, NonNullableUpdater, Updater } from './utils.js';
import { ParsedLocation } from './location.js';
export type IsRequiredParams<TParams> = Record<never, never> extends TParams ? never : true;
export interface ParsePathParamsResult<in out TRequired, in out TOptional, in out TRest> {
    required: TRequired;
    optional: TOptional;
    rest: TRest;
}
export type AnyParsePathParamsResult = ParsePathParamsResult<string, string, string>;
export type ParsePathParamsBoundaryStart<T extends string> = T extends `${infer TLeft}{-${infer TRight}` ? ParsePathParamsResult<ParsePathParams<TLeft>['required'], ParsePathParams<TLeft>['optional'] | ParsePathParams<TRight>['required'] | ParsePathParams<TRight>['optional'], ParsePathParams<TRight>['rest']> : T extends `${infer TLeft}{${infer TRight}` ? ParsePathParamsResult<ParsePathParams<TLeft>['required'] | ParsePathParams<TRight>['required'], ParsePathParams<TLeft>['optional'] | ParsePathParams<TRight>['optional'], ParsePathParams<TRight>['rest']> : never;
export type ParsePathParamsSymbol<T extends string> = T extends `${string}$${infer TRight}` ? TRight extends `${string}/${string}` ? TRight extends `${infer TParam}/${infer TRest}` ? TParam extends '' ? ParsePathParamsResult<ParsePathParams<TRest>['required'], '_splat' | ParsePathParams<TRest>['optional'], ParsePathParams<TRest>['rest']> : ParsePathParamsResult<TParam | ParsePathParams<TRest>['required'], ParsePathParams<TRest>['optional'], ParsePathParams<TRest>['rest']> : never : TRight extends '' ? ParsePathParamsResult<never, '_splat', never> : ParsePathParamsResult<TRight, never, never> : never;
export type ParsePathParamsBoundaryEnd<T extends string> = T extends `${infer TLeft}}${infer TRight}` ? ParsePathParamsResult<ParsePathParams<TLeft>['required'] | ParsePathParams<TRight>['required'], ParsePathParams<TLeft>['optional'] | ParsePathParams<TRight>['optional'], ParsePathParams<TRight>['rest']> : never;
export type ParsePathParamsEscapeStart<T extends string> = T extends `${infer TLeft}[${infer TRight}` ? ParsePathParamsResult<ParsePathParams<TLeft>['required'] | ParsePathParams<TRight>['required'], ParsePathParams<TLeft>['optional'] | ParsePathParams<TRight>['optional'], ParsePathParams<TRight>['rest']> : never;
export type ParsePathParamsEscapeEnd<T extends string> = T extends `${string}]${infer TRight}` ? ParsePathParams<TRight> : never;
export type ParsePathParams<T extends string> = T extends `${string}[${string}` ? ParsePathParamsEscapeStart<T> : T extends `${string}]${string}` ? ParsePathParamsEscapeEnd<T> : T extends `${string}}${string}` ? ParsePathParamsBoundaryEnd<T> : T extends `${string}{${string}` ? ParsePathParamsBoundaryStart<T> : T extends `${string}$${string}` ? ParsePathParamsSymbol<T> : never;
export type AddTrailingSlash<T> = T extends `${string}/` ? T : `${T & string}/`;
export type RemoveTrailingSlashes<T> = T & `${string}/` extends never ? T : T extends `${infer R}/` ? R : T;
export type AddLeadingSlash<T> = T & `/${string}` extends never ? `/${T & string}` : T;
export type RemoveLeadingSlashes<T> = T & `/${string}` extends never ? T : T extends `/${infer R}` ? R : T;
type JoinPath<TLeft extends string, TRight extends string> = TRight extends '' ? TLeft : TLeft extends '' ? TRight : `${RemoveTrailingSlashes<TLeft>}/${RemoveLeadingSlashes<TRight>}`;
type RemoveLastSegment<T extends string, TAcc extends string = ''> = T extends `${infer TSegment}/${infer TRest}` ? TRest & `${string}/${string}` extends never ? TRest extends '' ? TAcc : `${TAcc}${TSegment}` : RemoveLastSegment<TRest, `${TAcc}${TSegment}/`> : TAcc;
export type ResolveCurrentPath<TFrom extends string, TTo extends string> = TTo extends '.' ? TFrom : TTo extends './' ? AddTrailingSlash<TFrom> : TTo & `./${string}` extends never ? never : TTo extends `./${infer TRest}` ? AddLeadingSlash<JoinPath<TFrom, TRest>> : never;
export type ResolveParentPath<TFrom extends string, TTo extends string> = TTo extends '../' | '..' ? TFrom extends '' | '/' ? never : AddLeadingSlash<RemoveLastSegment<TFrom>> : TTo & `../${string}` extends never ? AddLeadingSlash<JoinPath<TFrom, TTo>> : TFrom extends '' | '/' ? never : TTo extends `../${infer ToRest}` ? ResolveParentPath<RemoveLastSegment<TFrom>, ToRest> : AddLeadingSlash<JoinPath<TFrom, TTo>>;
export type ResolveRelativePath<TFrom, TTo = '.'> = string extends TFrom ? TTo : string extends TTo ? TFrom : undefined extends TTo ? TFrom : TTo extends string ? TFrom extends string ? TTo extends `/${string}` ? TTo : TTo extends `..${string}` ? ResolveParentPath<TFrom, TTo> : TTo extends `.${string}` ? ResolveCurrentPath<TFrom, TTo> : AddLeadingSlash<JoinPath<TFrom, TTo>> : never : never;
export type FindDescendantToPaths<TRouter extends AnyRouter, TPrefix extends string> = `${TPrefix}/${string}` & RouteToPath<TRouter>;
export type InferDescendantToPaths<TRouter extends AnyRouter, TPrefix extends string, TPaths = FindDescendantToPaths<TRouter, TPrefix>> = TPaths extends `${TPrefix}/` ? never : TPaths extends `${TPrefix}/${infer TRest}` ? TRest : never;
export type RelativeToPath<TRouter extends AnyRouter, TTo extends string, TResolvedPath extends string> = (TResolvedPath & RouteToPath<TRouter> extends never ? never : ToPath<TRouter, TTo>) | `${RemoveTrailingSlashes<TTo>}/${InferDescendantToPaths<TRouter, RemoveTrailingSlashes<TResolvedPath>>}`;
export type RelativeToParentPath<TRouter extends AnyRouter, TFrom extends string, TTo extends string, TResolvedPath extends string = ResolveRelativePath<TFrom, TTo>> = RelativeToPath<TRouter, TTo, TResolvedPath> | (TTo extends `${string}..` | `${string}../` ? TResolvedPath extends '/' | '' ? never : FindDescendantToPaths<TRouter, RemoveTrailingSlashes<TResolvedPath>> extends never ? never : `${RemoveTrailingSlashes<TTo>}/${ParentPath<TRouter>}` : never);
export type RelativeToCurrentPath<TRouter extends AnyRouter, TFrom extends string, TTo extends string, TResolvedPath extends string = ResolveRelativePath<TFrom, TTo>> = RelativeToPath<TRouter, TTo, TResolvedPath> | CurrentPath<TRouter>;
export type AbsoluteToPath<TRouter extends AnyRouter, TFrom extends string> = (string extends TFrom ? CurrentPath<TRouter> : TFrom extends `/` ? never : CurrentPath<TRouter>) | (string extends TFrom ? ParentPath<TRouter> : TFrom extends `/` ? never : ParentPath<TRouter>) | RouteToPath<TRouter> | (TFrom extends '/' ? never : string extends TFrom ? never : InferDescendantToPaths<TRouter, RemoveTrailingSlashes<TFrom>>);
export type RelativeToPathAutoComplete<TRouter extends AnyRouter, TFrom extends string, TTo extends string> = string extends TTo ? string : string extends TFrom ? AbsoluteToPath<TRouter, TFrom> : TTo & `..${string}` extends never ? TTo & `.${string}` extends never ? AbsoluteToPath<TRouter, TFrom> : RelativeToCurrentPath<TRouter, TFrom, TTo> : RelativeToParentPath<TRouter, TFrom, TTo>;
export type NavigateOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = '.', TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = ToOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & NavigateOptionProps;
/**
 * The NavigateOptions type is used to describe the options that can be used when describing a navigation action in TanStack Router.
 * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType)
 */
export interface NavigateOptionProps {
    /**
     * If set to `true`, the router will scroll the element with an id matching the hash into view with default `ScrollIntoViewOptions`.
     * If set to `false`, the router will not scroll the element with an id matching the hash into view.
     * If set to `ScrollIntoViewOptions`, the router will scroll the element with an id matching the hash into view with the provided options.
     * @default true
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#hashscrollintoview)
     * @see [MDN](https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView)
     */
    hashScrollIntoView?: boolean | ScrollIntoViewOptions;
    /**
     * `replace` is a boolean that determines whether the navigation should replace the current history entry or push a new one.
     * @default false
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#replace)
     */
    replace?: boolean;
    /**
     * Defaults to `true` so that the scroll position will be reset to 0,0 after the location is committed to the browser history.
     * If `false`, the scroll position will not be reset to 0,0 after the location is committed to history.
     * @default true
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#resetscroll)
     */
    resetScroll?: boolean;
    /** @deprecated All navigations now use startTransition under the hood */
    startTransition?: boolean;
    /**
     * If set to `true`, the router will wrap the resulting navigation in a `document.startViewTransition()` call.
     * If `ViewTransitionOptions`, route navigations will be called using `document.startViewTransition({update, types})`
     * where `types` will be the strings array passed with `ViewTransitionOptions["types"]`.
     * If the browser does not support viewTransition types, the navigation will fall back to normal `document.startTransition()`, same as if `true` was passed.
     *
     * If the browser does not support this api, this option will be ignored.
     * @default false
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#viewtransition)
     * @see [MDN](https://developer.mozilla.org/en-US/docs/Web/API/Document/startViewTransition)
     * @see [Google](https://developer.chrome.com/docs/web-platform/view-transitions/same-document#view-transition-types)
     */
    viewTransition?: boolean | ViewTransitionOptions;
    /**
     * If `true`, navigation will ignore any blockers that might prevent it.
     * @default false
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#ignoreblocker)
     */
    ignoreBlocker?: boolean;
    /**
     * If `true`, navigation to a route inside of router will trigger a full page load instead of the traditional SPA navigation.
     * @default false
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#reloaddocument)
     */
    reloadDocument?: boolean;
    /**
     * This can be used instead of `to` to navigate to a fully built href, e.g. pointing to an external target.
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#href)
     */
    href?: string;
}
export type ToOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = '.', TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = ToSubOptions<TRouter, TFrom, TTo> & MaskOptions<TRouter, TMaskFrom, TMaskTo>;
export interface MaskOptions<in out TRouter extends AnyRouter, in out TMaskFrom extends string, in out TMaskTo extends string> {
    _fromLocation?: ParsedLocation;
    mask?: ToMaskOptions<TRouter, TMaskFrom, TMaskTo>;
}
export type ToMaskOptions<TRouter extends AnyRouter = RegisteredRouter, TMaskFrom extends string = string, TMaskTo extends string = '.'> = ToSubOptions<TRouter, TMaskFrom, TMaskTo> & {
    unmaskOnReload?: boolean;
};
export type ToSubOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = '.'> = ToSubOptionsProps<TRouter, TFrom, TTo> & SearchParamOptions<TRouter, TFrom, TTo> & PathParamOptions<TRouter, TFrom, TTo>;
export interface RequiredToOptions<in out TRouter extends AnyRouter, in out TFrom extends string, in out TTo extends string | undefined> {
    /**
     * The internal route path to navigate to. This should be a relative or absolute path within your application.
     * For external URLs, use the `href` property instead.
     * @example "/dashboard" or "../profile"
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#href)
     */
    to: ToPathOption<TRouter, TFrom, TTo> & {};
}
export interface OptionalToOptions<in out TRouter extends AnyRouter, in out TFrom extends string, in out TTo extends string | undefined> {
    /**
     * The internal route path to navigate to. This should be a relative or absolute path within your application.
     * For external URLs, use the `href` property instead.
     * @example "/dashboard" or "../profile"
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType#href)
     */
    to?: ToPathOption<TRouter, TFrom, TTo> & {};
}
export type MakeToRequired<TRouter extends AnyRouter, TFrom extends string, TTo extends string | undefined> = string extends TFrom ? string extends TTo ? OptionalToOptions<TRouter, TFrom, TTo> : TTo & CatchAllPaths<TRouter> extends never ? RequiredToOptions<TRouter, TFrom, TTo> : OptionalToOptions<TRouter, TFrom, TTo> : OptionalToOptions<TRouter, TFrom, TTo>;
export type ToSubOptionsProps<TRouter extends AnyRouter = RegisteredRouter, TFrom extends RoutePaths<TRouter['routeTree']> | string = string, TTo extends string | undefined = '.'> = MakeToRequired<TRouter, TFrom, TTo> & {
    hash?: true | Updater<string>;
    state?: true | NonNullableUpdater<ParsedHistoryState, HistoryState>;
    from?: FromPathOption<TRouter, TFrom> & {};
    unsafeRelative?: 'path';
};
export type ParamsReducerFn<in out TRouter extends AnyRouter, in out TParamVariant extends ParamVariant, in out TFrom, in out TTo> = (current: Expand<ResolveFromParams<TRouter, TParamVariant, TFrom>>) => Expand<ResolveRelativeToParams<TRouter, TParamVariant, TFrom, TTo>>;
type ParamsReducer<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = Expand<ResolveRelativeToParams<TRouter, TParamVariant, TFrom, TTo>> | (ParamsReducerFn<TRouter, TParamVariant, TFrom, TTo> & {});
type ParamVariant = 'PATH' | 'SEARCH';
export type ResolveRoute<TRouter extends AnyRouter, TFrom, TTo, TPath = ResolveRelativePath<TFrom, TTo>> = TPath extends string ? TFrom extends TPath ? RouteByPath<TRouter['routeTree'], TPath> : RouteByToPath<TRouter, TPath> : never;
type ResolveFromParamType<TParamVariant extends ParamVariant> = TParamVariant extends 'PATH' ? 'allParams' : 'fullSearchSchema';
type ResolveFromAllParams<TRouter extends AnyRouter, TParamVariant extends ParamVariant> = TParamVariant extends 'PATH' ? AllParams<TRouter['routeTree']> : FullSearchSchema<TRouter['routeTree']>;
type ResolveFromParams<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom> = string extends TFrom ? ResolveFromAllParams<TRouter, TParamVariant> : RouteByPath<TRouter['routeTree'], TFrom>['types'][ResolveFromParamType<TParamVariant>];
type ResolveToParamType<TParamVariant extends ParamVariant> = TParamVariant extends 'PATH' ? 'allParams' : 'fullSearchSchemaInput';
type ResolveAllToParams<TRouter extends AnyRouter, TParamVariant extends ParamVariant> = TParamVariant extends 'PATH' ? AllParams<TRouter['routeTree']> : FullSearchSchemaInput<TRouter['routeTree']>;
export type ResolveToParams<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = ResolveRelativePath<TFrom, TTo> extends infer TPath ? undefined extends TPath ? never : string extends TPath ? ResolveAllToParams<TRouter, TParamVariant> : TPath extends CatchAllPaths<TRouter> ? ResolveAllToParams<TRouter, TParamVariant> : ResolveRoute<TRouter, TFrom, TTo>['types'][ResolveToParamType<TParamVariant>] : never;
type ResolveRelativeToParams<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo, TToParams = ResolveToParams<TRouter, TParamVariant, TFrom, TTo>> = TParamVariant extends 'SEARCH' ? TToParams : string extends TFrom ? TToParams : MakeDifferenceOptional<ResolveFromParams<TRouter, TParamVariant, TFrom>, TToParams>;
export interface MakeOptionalSearchParams<in out TRouter extends AnyRouter, in out TFrom, in out TTo> {
    search?: true | (ParamsReducer<TRouter, 'SEARCH', TFrom, TTo> & {});
}
export interface MakeOptionalPathParams<in out TRouter extends AnyRouter, in out TFrom, in out TTo> {
    params?: true | (ParamsReducer<TRouter, 'PATH', TFrom, TTo> & {});
}
type MakeRequiredParamsReducer<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = (string extends TFrom ? never : ResolveFromParams<TRouter, TParamVariant, TFrom> extends ResolveRelativeToParams<TRouter, TParamVariant, TFrom, TTo> ? true : never) | (ParamsReducer<TRouter, TParamVariant, TFrom, TTo> & {});
export interface MakeRequiredPathParams<in out TRouter extends AnyRouter, in out TFrom, in out TTo> {
    params: MakeRequiredParamsReducer<TRouter, 'PATH', TFrom, TTo> & {};
}
export interface MakeRequiredSearchParams<in out TRouter extends AnyRouter, in out TFrom, in out TTo> {
    search: MakeRequiredParamsReducer<TRouter, 'SEARCH', TFrom, TTo> & {};
}
export type IsRequired<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = ResolveRelativePath<TFrom, TTo> extends infer TPath ? undefined extends TPath ? never : TPath extends CatchAllPaths<TRouter> ? never : IsRequiredParams<ResolveRelativeToParams<TRouter, TParamVariant, TFrom, TTo>> : never;
export type SearchParamOptions<TRouter extends AnyRouter, TFrom, TTo> = IsRequired<TRouter, 'SEARCH', TFrom, TTo> extends never ? MakeOptionalSearchParams<TRouter, TFrom, TTo> : MakeRequiredSearchParams<TRouter, TFrom, TTo>;
export type PathParamOptions<TRouter extends AnyRouter, TFrom, TTo> = IsRequired<TRouter, 'PATH', TFrom, TTo> extends never ? MakeOptionalPathParams<TRouter, TFrom, TTo> : MakeRequiredPathParams<TRouter, TFrom, TTo>;
export type ToPathOption<TRouter extends AnyRouter = AnyRouter, TFrom extends string = string, TTo extends string | undefined = string> = ConstrainLiteral<TTo, RelativeToPathAutoComplete<TRouter, NoInfer<TFrom> extends string ? NoInfer<TFrom> : '', NoInfer<TTo> & string>>;
export type FromPathOption<TRouter extends AnyRouter, TFrom> = ConstrainLiteral<TFrom, RoutePaths<TRouter['routeTree']>>;
/**
 * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/navigation#active-options)
 */
export interface ActiveOptions {
    /**
     * If true, the link will be active if the current route matches the `to` route path exactly (no children routes)
     * @default false
     */
    exact?: boolean;
    /**
     * If true, the link will only be active if the current URL hash matches the `hash` prop
     * @default false
     */
    includeHash?: boolean;
    /**
     * If true, the link will only be active if the current URL search params inclusively match the `search` prop
     * @default true
     */
    includeSearch?: boolean;
    /**
     * This modifies the `includeSearch` behavior.
     * If true,  properties in `search` that are explicitly `undefined` must NOT be present in the current URL search params for the link to be active.
     * @default false
     */
    explicitUndefined?: boolean;
}
export interface LinkOptionsProps {
    /**
     * The standard anchor tag target attribute
     */
    target?: HTMLAnchorElement['target'];
    /**
     * Configurable options to determine if the link should be considered active or not
     * @default {exact:true,includeHash:true}
     */
    activeOptions?: ActiveOptions;
    /**
     * The preloading strategy for this link
     * - `false` - No preloading
     * - `'intent'` - Preload the linked route on hover and cache it for this many milliseconds in hopes that the user will eventually navigate there.
     * - `'viewport'` - Preload the linked route when it enters the viewport
     */
    preload?: false | 'intent' | 'viewport' | 'render';
    /**
     * When a preload strategy is set, this delays the preload by this many milliseconds.
     * If the user exits the link before this delay, the preload will be cancelled.
     */
    preloadDelay?: number;
    /**
     * Control whether the link should be disabled or not
     * If set to `true`, the link will be rendered without an `href` attribute
     * @default false
     */
    disabled?: boolean;
    /**
     * When the preload strategy is set to `intent`, this controls the proximity of the link to the cursor before it is preloaded.
     * If the user exits this proximity before this delay, the preload will be cancelled.
     */
    preloadIntentProximity?: number;
}
export type LinkOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = '.', TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & LinkOptionsProps;
export declare const preloadWarning = "Error preloading route! \u261D\uFE0F";
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/routeInfo.d.ts
import { InferFileRouteTypes } from './fileRoute.js';
import { AddTrailingSlash, RemoveTrailingSlashes } from './link.js';
import { AnyRoute } from './route.js';
import { AnyRouter, TrailingSlashOption } from './router.js';
import { PartialMergeAll } from './utils.js';
export type ParseRoute<TRouteTree, TAcc = TRouteTree> = TRouteTree extends {
    types: {
        children: infer TChildren;
    };
} ? unknown extends TChildren ? TAcc : TChildren extends ReadonlyArray<any> ? ParseRoute<TChildren[number], TAcc | TChildren[number]> : ParseRoute<TChildren[keyof TChildren], TAcc | TChildren[keyof TChildren]> : TAcc;
export type ParseRouteWithoutBranches<TRouteTree> = ParseRoute<TRouteTree> extends infer TRoute extends AnyRoute ? TRoute extends any ? unknown extends TRoute['types']['children'] ? TRoute : TRoute['types']['children'] extends ReadonlyArray<any> ? '/' extends TRoute['types']['children'][number]['path'] ? never : TRoute : '/' extends TRoute['types']['children'][keyof TRoute['types']['children']]['path'] ? never : TRoute : never : never;
export type CodeRoutesById<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? {
    [K in TRoutes as K['id']]: K;
} : never;
export type RoutesById<TRouteTree extends AnyRoute> = InferFileRouteTypes<TRouteTree> extends never ? CodeRoutesById<TRouteTree> : InferFileRouteTypes<TRouteTree>['fileRoutesById'];
export type RouteById<TRouteTree extends AnyRoute, TId> = Extract<RoutesById<TRouteTree>[TId & keyof RoutesById<TRouteTree>], AnyRoute>;
export type CodeRouteIds<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? TRoutes['id'] : never;
export type RouteIds<TRouteTree extends AnyRoute> = InferFileRouteTypes<TRouteTree> extends never ? CodeRouteIds<TRouteTree> : InferFileRouteTypes<TRouteTree>['id'];
export type ParentPath<TRouter extends AnyRouter> = TrailingSlashOptionByRouter<TRouter> extends 'always' ? '../' : TrailingSlashOptionByRouter<TRouter> extends 'never' ? '..' : '../' | '..';
export type CurrentPath<TRouter extends AnyRouter> = TrailingSlashOptionByRouter<TRouter> extends 'always' ? './' : TrailingSlashOptionByRouter<TRouter> extends 'never' ? '.' : './' | '.';
export type ToPath<TRouter extends AnyRouter, TTo extends string> = TrailingSlashOptionByRouter<TRouter> extends 'always' ? AddTrailingSlash<TTo> : TrailingSlashOptionByRouter<TRouter> extends 'never' ? RemoveTrailingSlashes<TTo> : AddTrailingSlash<TTo> | RemoveTrailingSlashes<TTo>;
export type CatchAllPaths<TRouter extends AnyRouter> = CurrentPath<TRouter> | ParentPath<TRouter>;
export type CodeRoutesByPath<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? {
    [K in TRoutes as K['fullPath']]: K;
} : never;
export type RoutesByPath<TRouteTree extends AnyRoute> = InferFileRouteTypes<TRouteTree> extends never ? CodeRoutesByPath<TRouteTree> : InferFileRouteTypes<TRouteTree>['fileRoutesByFullPath'];
export type RouteByPath<TRouteTree extends AnyRoute, TPath> = Extract<RoutesByPath<TRouteTree>[TPath & keyof RoutesByPath<TRouteTree>], AnyRoute>;
export type CodeRoutePaths<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? TRoutes['fullPath'] : never;
export type RoutePaths<TRouteTree extends AnyRoute> = unknown extends TRouteTree ? string : (InferFileRouteTypes<TRouteTree> extends never ? CodeRoutePaths<TRouteTree> : InferFileRouteTypes<TRouteTree>['fullPaths']) | '/';
export type RouteToPathAlwaysTrailingSlash<TRoute extends AnyRoute> = TRoute['path'] extends '/' ? TRoute['fullPath'] : TRoute['fullPath'] extends '/' ? TRoute['fullPath'] : `${TRoute['fullPath']}/`;
export type RouteToPathNeverTrailingSlash<TRoute extends AnyRoute> = TRoute['path'] extends '/' ? TRoute['fullPath'] extends '/' ? TRoute['fullPath'] : RemoveTrailingSlashes<TRoute['fullPath']> : TRoute['fullPath'];
export type RouteToPathPreserveTrailingSlash<TRoute extends AnyRoute> = RouteToPathNeverTrailingSlash<TRoute> | RouteToPathAlwaysTrailingSlash<TRoute>;
export type RouteToPathByTrailingSlashOption<TRoute extends AnyRoute> = {
    always: RouteToPathAlwaysTrailingSlash<TRoute>;
    preserve: RouteToPathPreserveTrailingSlash<TRoute>;
    never: RouteToPathNeverTrailingSlash<TRoute>;
};
export type TrailingSlashOptionByRouter<TRouter extends AnyRouter> = TrailingSlashOption extends TRouter['options']['trailingSlash'] ? 'never' : NonNullable<TRouter['options']['trailingSlash']>;
export type RouteToByRouter<TRouter extends AnyRouter, TRoute extends AnyRoute> = RouteToPathByTrailingSlashOption<TRoute>[TrailingSlashOptionByRouter<TRouter>];
export type CodeRouteToPath<TRouter extends AnyRouter> = ParseRouteWithoutBranches<TRouter['routeTree']> extends infer TRoute extends AnyRoute ? TRoute extends any ? RouteToByRouter<TRouter, TRoute> : never : never;
export type FileRouteToPath<TRouter extends AnyRouter, TTo = InferFileRouteTypes<TRouter['routeTree']>['to'], TTrailingSlashOption = TrailingSlashOptionByRouter<TRouter>> = 'never' extends TTrailingSlashOption ? TTo : 'always' extends TTrailingSlashOption ? AddTrailingSlash<TTo> : TTo | AddTrailingSlash<TTo>;
export type RouteToPath<TRouter extends AnyRouter> = unknown extends TRouter ? string : InferFileRouteTypes<TRouter['routeTree']> extends never ? CodeRouteToPath<TRouter> : FileRouteToPath<TRouter>;
export type CodeRoutesByToPath<TRouter extends AnyRouter> = ParseRouteWithoutBranches<TRouter['routeTree']> extends infer TRoutes extends AnyRoute ? {
    [TRoute in TRoutes as RouteToByRouter<TRouter, TRoute>]: TRoute;
} : never;
export type RoutesByToPath<TRouter extends AnyRouter> = InferFileRouteTypes<TRouter['routeTree']> extends never ? CodeRoutesByToPath<TRouter> : InferFileRouteTypes<TRouter['routeTree']>['fileRoutesByTo'];
export type CodeRouteByToPath<TRouter extends AnyRouter, TTo> = Extract<RoutesByToPath<TRouter>[TTo & keyof RoutesByToPath<TRouter>], AnyRoute>;
export type FileRouteByToPath<TRouter extends AnyRouter, TTo> = 'never' extends TrailingSlashOptionByRouter<TRouter> ? CodeRouteByToPath<TRouter, TTo> : 'always' extends TrailingSlashOptionByRouter<TRouter> ? TTo extends '/' ? CodeRouteByToPath<TRouter, TTo> : TTo extends `${infer TPath}/` ? CodeRouteByToPath<TRouter, TPath> : never : CodeRouteByToPath<TRouter, TTo extends '/' ? TTo : RemoveTrailingSlashes<TTo>>;
export type RouteByToPath<TRouter extends AnyRouter, TTo> = InferFileRouteTypes<TRouter['routeTree']> extends never ? CodeRouteByToPath<TRouter, TTo> : FileRouteByToPath<TRouter, TTo>;
export type FullSearchSchema<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? PartialMergeAll<TRoutes['types']['fullSearchSchema']> : never;
export type FullSearchSchemaInput<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? PartialMergeAll<TRoutes['types']['fullSearchSchemaInput']> : never;
export type AllParams<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? PartialMergeAll<TRoutes['types']['allParams']> : never;
export type AllContext<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? PartialMergeAll<TRoutes['types']['allContext']> : never;
export type AllLoaderData<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? PartialMergeAll<TRoutes['types']['loaderData']> : never;


// @filename: /node_modules/@tanstack/router-core/dist/esm/not-found.d.ts
import { RouteIds } from './routeInfo.js';
import { RegisteredRouter } from './router.js';
export type NotFoundError = {
    /**
      @deprecated
      Use `routeId: rootRouteId` instead
    */
    global?: boolean;
    /**
      @private
      Do not use this. It's used internally to indicate a path matching error
    */
    _global?: boolean;
    data?: any;
    throw?: boolean;
    routeId?: RouteIds<RegisteredRouter['routeTree']>;
    headers?: HeadersInit;
};
/**
 * Create a not-found error object recognized by TanStack Router.
 *
 * Throw this from loaders/actions to trigger the nearest `notFoundComponent`.
 * Use `routeId` to target a specific route's not-found boundary. If `throw`
 * is true, the error is thrown instead of returned.
 *
 * @param options Optional settings including `routeId`, `headers`, and `throw`.
 * @returns A not-found error object that can be thrown or returned.
 * @link https://tanstack.com/router/latest/docs/router/framework/react/api/router/notFoundFunction
 */
export declare function notFound(options?: NotFoundError): NotFoundError;
/** Determine if a value is a TanStack Router not-found error. */
export declare function isNotFound(obj: any): obj is NotFoundError;


// @filename: /node_modules/@tanstack/router-core/dist/esm/redirect.d.ts
import { NavigateOptions } from './link.js';
import { AnyRouter, RegisteredRouter } from './router.js';
export type AnyRedirect = Redirect<any, any, any, any, any>;
/**
 * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RedirectType)
 */
export type Redirect<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = undefined, TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = Response & {
    options: NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & {};
    redirectHandled?: boolean;
};
export type RedirectOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = undefined, TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = {
    href?: string;
    /**
     * @deprecated Use `statusCode` instead
     **/
    code?: number;
    /**
     * The HTTP status code to use when redirecting.
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RedirectType#statuscode-property)
     */
    statusCode?: number;
    /**
     * If provided, will throw the redirect object instead of returning it. This can be useful in places where `throwing` in a function might cause it to have a return type of `never`. In that case, you can use `redirect({ throw: true })` to throw the redirect object instead of returning it.
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RedirectType#throw-property)
     */
    throw?: any;
    /**
     * The HTTP headers to use when redirecting.
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RedirectType#headers-property)
     */
    headers?: HeadersInit;
} & NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>;
export type ResolvedRedirect<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string = '', TMaskFrom extends string = TFrom, TMaskTo extends string = ''> = Redirect<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>;
/**
 * Options for route-bound redirect, where 'from' is automatically set.
 * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RedirectType)
 */
export type RedirectOptionsRoute<TDefaultFrom extends string = string, TRouter extends AnyRouter = RegisteredRouter, TTo extends string | undefined = undefined, TMaskTo extends string = ''> = Omit<RedirectOptions<TRouter, TDefaultFrom, TTo, TDefaultFrom, TMaskTo>, 'from'>;
/**
 * A redirect function bound to a specific route, with 'from' pre-set to the route's fullPath.
 * This enables relative redirects like `Route.redirect({ to: './overview' })`.
 * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RedirectType)
 */
export interface RedirectFnRoute<in out TDefaultFrom extends string = string> {
    <TRouter extends AnyRouter = RegisteredRouter, const TTo extends string | undefined = undefined, const TMaskTo extends string = ''>(opts: RedirectOptionsRoute<TDefaultFrom, TRouter, TTo, TMaskTo>): Redirect<TRouter, TDefaultFrom, TTo, TDefaultFrom, TMaskTo>;
}
/**
 * Create a redirect Response understood by TanStack Router.
 *
 * Use from route `loader`/`beforeLoad` or server functions to trigger a
 * navigation. If `throw: true` is set, the redirect is thrown instead of
 * returned. When an absolute `href` is supplied and `reloadDocument` is not
 * set, a full-document navigation is inferred.
 *
 * @param opts Options for the redirect. Common fields:
 * - `href`: absolute URL for external redirects; infers `reloadDocument`.
 * - `statusCode`: HTTP status code to use (defaults to 307).
 * - `headers`: additional headers to include on the Response.
 * - Standard navigation options like `to`, `params`, `search`, `replace`,
 *   and `reloadDocument` for internal redirects.
 * @returns A Response augmented with router navigation options.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/redirectFunction
 */
export declare function redirect<TRouter extends AnyRouter = RegisteredRouter, const TTo extends string | undefined = '.', const TFrom extends string = string, const TMaskFrom extends string = TFrom, const TMaskTo extends string = ''>(opts: RedirectOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>): Redirect<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>;
/** Check whether a value is a TanStack Router redirect Response. */
/** Check whether a value is a TanStack Router redirect Response. */
export declare function isRedirect(obj: any): obj is AnyRedirect;
/** True if value is a redirect with a resolved `href` location. */
/** True if value is a redirect with a resolved `href` location. */
export declare function isResolvedRedirect(obj: any): obj is AnyRedirect & {
    options: {
        href: string;
    };
};
/** Parse a serialized redirect object back into a redirect Response. */
/** Parse a serialized redirect object back into a redirect Response. */
export declare function parseRedirect(obj: any): Redirect<AnyRouter, string, ".", string, ""> | undefined;


// @filename: /node_modules/@tanstack/router-core/dist/esm/Matches.d.ts
import { AnyRoute, StaticDataRouteOption } from './route.js';
import { AllContext, AllLoaderData, AllParams, FullSearchSchema, ParseRoute, RouteById, RouteIds } from './routeInfo.js';
import { AnyRouter, RegisteredRouter, SSROption } from './router.js';
import { Constrain, ControlledPromise } from './utils.js';
export type AnyMatchAndValue = {
    match: any;
    value: any;
};
export type FindValueByIndex<TKey, TValue extends ReadonlyArray<any>> = TKey extends `${infer TIndex extends number}` ? TValue[TIndex] : never;
export type FindValueByKey<TKey, TValue> = TValue extends ReadonlyArray<any> ? FindValueByIndex<TKey, TValue> : TValue[TKey & keyof TValue];
export type CreateMatchAndValue<TMatch, TValue> = TValue extends any ? {
    match: TMatch;
    value: TValue;
} : never;
export type NextMatchAndValue<TKey, TMatchAndValue extends AnyMatchAndValue> = TMatchAndValue extends any ? CreateMatchAndValue<TMatchAndValue['match'], FindValueByKey<TKey, TMatchAndValue['value']>> : never;
export type IsMatchKeyOf<TValue> = TValue extends ReadonlyArray<any> ? number extends TValue['length'] ? `${number}` : keyof TValue & `${number}` : TValue extends object ? keyof TValue & string : never;
export type IsMatchPath<TParentPath extends string, TMatchAndValue extends AnyMatchAndValue> = `${TParentPath}${IsMatchKeyOf<TMatchAndValue['value']>}`;
export type IsMatchResult<TKey, TMatchAndValue extends AnyMatchAndValue> = TMatchAndValue extends any ? TKey extends keyof TMatchAndValue['value'] ? TMatchAndValue['match'] : never : never;
export type IsMatchParse<TPath, TMatchAndValue extends AnyMatchAndValue, TParentPath extends string = ''> = TPath extends `${string}.${string}` ? TPath extends `${infer TFirst}.${infer TRest}` ? IsMatchParse<TRest, NextMatchAndValue<TFirst, TMatchAndValue>, `${TParentPath}${TFirst}.`> : never : {
    path: IsMatchPath<TParentPath, TMatchAndValue>;
    result: IsMatchResult<TPath, TMatchAndValue>;
};
export type IsMatch<TMatch, TPath> = IsMatchParse<TPath, TMatch extends any ? {
    match: TMatch;
    value: TMatch;
} : never>;
/**
 * Narrows matches based on a path
 * @experimental
 */
export declare const isMatch: <TMatch, TPath extends string>(match: TMatch, path: Constrain<TPath, IsMatch<TMatch, TPath>["path"]>) => match is IsMatch<TMatch, TPath>["result"];
export interface DefaultRouteMatchExtensions {
    scripts?: unknown;
    links?: unknown;
    headScripts?: unknown;
    meta?: unknown;
    styles?: unknown;
}
export interface RouteMatchExtensions extends DefaultRouteMatchExtensions {
}
export interface RouteMatch<out TRouteId, out TFullPath, out TAllParams, out TFullSearchSchema, out TLoaderData, out TAllContext, out TLoaderDeps> extends RouteMatchExtensions {
    id: string;
    routeId: TRouteId;
    fullPath: TFullPath;
    index: number;
    pathname: string;
    params: TAllParams;
    _strictParams: TAllParams;
    status: 'pending' | 'success' | 'error' | 'redirected' | 'notFound';
    isFetching: false | 'beforeLoad' | 'loader';
    error: unknown;
    paramsError: unknown;
    searchError: unknown;
    updatedAt: number;
    _nonReactive: {
        loadPromise?: ControlledPromise<void>;
        displayPendingPromise?: Promise<void>;
        minPendingPromise?: ControlledPromise<void>;
        dehydrated?: boolean;
    };
    loaderData?: TLoaderData;
    context: TAllContext;
    search: TFullSearchSchema;
    _strictSearch: TFullSearchSchema;
    fetchCount: number;
    abortController: AbortController;
    cause: 'preload' | 'enter' | 'stay';
    loaderDeps: TLoaderDeps;
    preload: boolean;
    invalid: boolean;
    headers?: Record<string, string>;
    globalNotFound?: boolean;
    staticData: StaticDataRouteOption;
    /** This attribute is not reactive */
    ssr?: SSROption;
    _forcePending?: boolean;
    _displayPending?: boolean;
}
export interface PreValidationErrorHandlingRouteMatch<TRouteId, TFullPath, TAllParams, TFullSearchSchema> {
    id: string;
    routeId: TRouteId;
    fullPath: TFullPath;
    index: number;
    pathname: string;
    search: {
        status: 'success';
        value: TFullSearchSchema;
    } | {
        status: 'error';
        error: unknown;
    };
    params: {
        status: 'success';
        value: TAllParams;
    } | {
        status: 'error';
        error: unknown;
    };
    staticData: StaticDataRouteOption;
    ssr?: boolean | 'data-only';
}
export type MakePreValidationErrorHandlingRouteMatchUnion<TRouter extends AnyRouter = RegisteredRouter, TRoute extends AnyRoute = ParseRoute<TRouter['routeTree']>> = TRoute extends any ? PreValidationErrorHandlingRouteMatch<TRoute['id'], TRoute['fullPath'], TRoute['types']['allParams'], TRoute['types']['fullSearchSchema']> : never;
export type MakeRouteMatchFromRoute<TRoute extends AnyRoute> = RouteMatch<TRoute['types']['id'], TRoute['types']['fullPath'], TRoute['types']['allParams'], TRoute['types']['fullSearchSchema'], TRoute['types']['loaderData'], TRoute['types']['allContext'], TRoute['types']['loaderDeps']>;
export type MakeRouteMatch<TRouteTree extends AnyRoute = RegisteredRouter['routeTree'], TRouteId = RouteIds<TRouteTree>, TStrict extends boolean = true> = RouteMatch<TRouteId, RouteById<TRouteTree, TRouteId>['types']['fullPath'], TStrict extends false ? AllParams<TRouteTree> : RouteById<TRouteTree, TRouteId>['types']['allParams'], TStrict extends false ? FullSearchSchema<TRouteTree> : RouteById<TRouteTree, TRouteId>['types']['fullSearchSchema'], TStrict extends false ? AllLoaderData<TRouteTree> : RouteById<TRouteTree, TRouteId>['types']['loaderData'], TStrict extends false ? AllContext<TRouteTree> : RouteById<TRouteTree, TRouteId>['types']['allContext'], RouteById<TRouteTree, TRouteId>['types']['loaderDeps']>;
export type AnyRouteMatch = RouteMatch<any, any, any, any, any, any, any>;
export type MakeRouteMatchUnion<TRouter extends AnyRouter = RegisteredRouter, TRoute extends AnyRoute = ParseRoute<TRouter['routeTree']>> = TRoute extends any ? RouteMatch<TRoute['id'], TRoute['fullPath'], TRoute['types']['allParams'], TRoute['types']['fullSearchSchema'], TRoute['types']['loaderData'], TRoute['types']['allContext'], TRoute['types']['loaderDeps']> : never;
/**
 * The `MatchRouteOptions` type is used to describe the options that can be used when matching a route.
 *
 * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/MatchRouteOptionsType#matchrouteoptions-type)
 */
export interface MatchRouteOptions {
    /**
     * If `true`, will match against pending location instead of the current location.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/MatchRouteOptionsType#pending-property)
     */
    pending?: boolean;
    /**
     * If `true`, will match against the current location with case sensitivity.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/MatchRouteOptionsType#casesensitive-property)
     *
     * @deprecated Declare case sensitivity in the route definition instead, or globally for all routes using the `caseSensitive` option on the router.
     */
    caseSensitive?: boolean;
    /**
     * If `true`, will match against the current location's search params using a deep inclusive check. e.g. `{ a: 1 }` will match for a current location of `{ a: 1, b: 2 }`.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/MatchRouteOptionsType#includesearch-property)
     */
    includeSearch?: boolean;
    /**
     * If `true`, will match against the current location using a fuzzy match. e.g. `/posts` will match for a current location of `/posts/123`.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/MatchRouteOptionsType#fuzzy-property)
     */
    fuzzy?: boolean;
}


// @filename: /node_modules/@tanstack/router-core/dist/esm/root.d.ts
/** Stable identifier used for the root route in a route tree. */
export declare const rootRouteId = "__root__";
export type RootRouteId = typeof rootRouteId;


// @filename: /node_modules/@tanstack/router-core/dist/esm/RouterProvider.d.ts
import { NavigateOptions, ToOptions } from './link.js';
import { ParsedLocation } from './location.js';
import { RoutePaths } from './routeInfo.js';
import { RegisteredRouter, ViewTransitionOptions } from './router.js';
export interface MatchLocation {
    to?: string | number | null;
    fuzzy?: boolean;
    caseSensitive?: boolean;
    from?: string;
}
export interface CommitLocationOptions {
    replace?: boolean;
    resetScroll?: boolean;
    hashScrollIntoView?: boolean | ScrollIntoViewOptions;
    viewTransition?: boolean | ViewTransitionOptions;
    /**
     * @deprecated All navigations use transitions under the hood now
     **/
    startTransition?: boolean;
    ignoreBlocker?: boolean;
}
export type NavigateFn = <TRouter extends RegisteredRouter, TTo extends string | undefined, TFrom extends RoutePaths<TRouter['routeTree']> | string = string, TMaskFrom extends RoutePaths<TRouter['routeTree']> | string = TFrom, TMaskTo extends string = ''>(opts: NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>) => Promise<void>;
export type BuildLocationFn = <TRouter extends RegisteredRouter, TTo extends string | undefined, TFrom extends RoutePaths<TRouter['routeTree']> | string = string, TMaskFrom extends RoutePaths<TRouter['routeTree']> | string = TFrom, TMaskTo extends string = ''>(opts: ToOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & {
    leaveParams?: boolean;
    _includeValidateSearch?: boolean;
    _isNavigate?: boolean;
}) => ParsedLocation;


// @filename: /node_modules/@tanstack/router-core/dist/esm/ssr/serializer/RawStream.d.ts
import { PluginInfo, SerovalNode } from 'seroval';
/**
 * Hint for RawStream encoding strategy during SSR serialization.
 * - 'binary': Always use base64 encoding (best for binary data like files, images)
 * - 'text': Try UTF-8 first, fallback to base64 (best for text-heavy data like RSC payloads)
 */
export type RawStreamHint = 'binary' | 'text';
/**
 * Options for RawStream configuration.
 */
export interface RawStreamOptions {
    /**
     * Encoding hint for SSR serialization.
     * - 'binary' (default): Always use base64 encoding
     * - 'text': Try UTF-8 first, fallback to base64 for invalid UTF-8 chunks
     */
    hint?: RawStreamHint;
}
/**
 * Marker class for ReadableStream<Uint8Array> that should be serialized
 * with base64 encoding (SSR) or binary framing (server functions).
 *
 * Wrap your binary streams with this to get efficient serialization:
 * ```ts
 * // For binary data (files, images, etc.)
 * return { data: new RawStream(file.stream()) }
 *
 * // For text-heavy data (RSC payloads, etc.)
 * return { data: new RawStream(rscStream, { hint: 'text' }) }
 * ```
 */
export declare class RawStream {
    readonly stream: ReadableStream<Uint8Array>;
    readonly hint: RawStreamHint;
    constructor(stream: ReadableStream<Uint8Array>, options?: RawStreamOptions);
}
/**
 * Callback type for RPC plugin to register raw streams with multiplexer
 */
export type OnRawStreamCallback = (streamId: number, stream: ReadableStream<Uint8Array>) => void;
export interface RawStreamSSRNode extends PluginInfo {
    hint: SerovalNode;
    factory: SerovalNode;
    stream: SerovalNode;
}
export interface RawStreamRPCNode extends PluginInfo {
    streamId: SerovalNode;
}
/**
 * SSR Plugin - uses base64 or UTF-8+base64 encoding for chunks, delegates to seroval's stream mechanism.
 * Used during SSR when serializing to JavaScript code for HTML injection.
 *
 * Supports two modes based on RawStream hint:
 * - 'binary': Always base64 encode (default)
 * - 'text': Try UTF-8 first, fallback to base64 for invalid UTF-8
 */
export declare const RawStreamSSRPlugin: import('seroval').Plugin<RawStream, RawStreamSSRNode>;
/**
 * Creates an RPC plugin instance that registers raw streams with a multiplexer.
 * Used for server function responses where we want binary framing.
 * Note: RPC always uses binary framing regardless of hint.
 *
 * @param onRawStream Callback invoked when a RawStream is encountered during serialization
 */
export declare function createRawStreamRPCPlugin(onRawStream: OnRawStreamCallback): import('seroval').Plugin<RawStream, RawStreamRPCNode>;
/**
 * Creates a deserialize-only plugin for client-side stream reconstruction.
 * Used in serverFnFetcher to wire up streams from frame decoder.
 *
 * @param getOrCreateStream Function to get/create a stream by ID from frame decoder
 */
export declare function createRawStreamDeserializePlugin(getOrCreateStream: (id: number) => ReadableStream<Uint8Array>): import('seroval').Plugin<any, RawStreamRPCNode>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/ssr/serializer/transformer.d.ts
import { Plugin, PluginInfo, SerovalNode } from 'seroval';
import { RegisteredConfigType, RegisteredSsr, SSROption } from '../../router.js';
import { LooseReturnType } from '../../utils.js';
import { AnyRoute, ResolveAllSSR } from '../../route.js';
import { RawStream } from './RawStream.js';
declare const TSR_SERIALIZABLE: unique symbol;
export type TSR_SERIALIZABLE = typeof TSR_SERIALIZABLE;
export type TsrSerializable = {
    [TSR_SERIALIZABLE]: true;
};
export interface DefaultSerializable {
    number: number;
    string: string;
    boolean: boolean;
    null: null;
    undefined: undefined;
    bigint: bigint;
    Date: Date;
    Uint8Array: Uint8Array;
    RawStream: RawStream;
    TsrSerializable: TsrSerializable;
    void: void;
}
export interface SerializableExtensions extends DefaultSerializable {
}
export type Serializable = SerializableExtensions[keyof SerializableExtensions];
export type UnionizeSerializationAdaptersInput<TAdapters extends ReadonlyArray<AnySerializationAdapter>> = TAdapters[number]['~types']['input'];
/**
 * Create a strongly-typed serialization adapter for SSR hydration.
 * Use to register custom types with the router serializer.
 */
export declare function createSerializationAdapter<TInput = unknown, TOutput = unknown, const TExtendsAdapters extends ReadonlyArray<AnySerializationAdapter> | never = never>(opts: CreateSerializationAdapterOptions<TInput, TOutput, TExtendsAdapters>): SerializationAdapter<TInput, TOutput, TExtendsAdapters>;
export interface CreateSerializationAdapterOptions<TInput, TOutput, TExtendsAdapters extends ReadonlyArray<AnySerializationAdapter> | never> {
    key: string;
    extends?: TExtendsAdapters;
    test: (value: unknown) => value is TInput;
    toSerializable: (value: TInput) => ValidateSerializable<TOutput, Serializable | UnionizeSerializationAdaptersInput<TExtendsAdapters>>;
    fromSerializable: (value: TOutput) => TInput;
}
export type ValidateSerializable<T, TSerializable> = T extends TSerializable ? T : T extends (...args: Array<any>) => any ? SerializationError<'Function may not be serializable'> : T extends RegisteredReadableStream ? SerializationError<'JSX is not be serializable'> : T extends ReadonlyArray<any> ? ValidateSerializableArray<T, TSerializable> : T extends Promise<any> ? ValidateSerializablePromise<T, TSerializable> : T extends ReadableStream<any> ? ValidateReadableStream<T, TSerializable> : T extends Set<any> ? ValidateSerializableSet<T, TSerializable> : T extends Map<any, any> ? ValidateSerializableMap<T, TSerializable> : T extends AsyncGenerator<any, any> ? ValidateSerializableAsyncGenerator<T, TSerializable> : T extends object ? ValidateSerializableMapped<T, TSerializable> : SerializationError<'Type may not be serializable'>;
export type ValidateSerializableAsyncGenerator<T, TSerializable> = T extends AsyncGenerator<infer T, infer TReturn, infer TNext> ? AsyncGenerator<ValidateSerializable<T, TSerializable>, ValidateSerializable<TReturn, TSerializable>, TNext> : never;
export type ValidateSerializablePromise<T, TSerializable> = T extends Promise<infer TAwaited> ? Promise<ValidateSerializable<TAwaited, TSerializable>> : never;
export type ValidateReadableStream<T, TSerializable> = T extends ReadableStream<infer TStreamed> ? ReadableStream<ValidateSerializable<TStreamed, TSerializable>> : never;
export type ValidateSerializableSet<T, TSerializable> = T extends Set<infer TItem> ? Set<ValidateSerializable<TItem, TSerializable>> : never;
export type ValidateSerializableMap<T, TSerializable> = T extends Map<infer TKey, infer TValue> ? Map<ValidateSerializable<TKey, TSerializable>, ValidateSerializable<TValue, TSerializable>> : never;
export type ValidateSerializableArray<T, TSerializable> = T extends readonly [
    any,
    ...Array<any>
] ? ValidateSerializableMapped<T, TSerializable> : T extends Array<infer U> ? Array<ValidateSerializable<U, TSerializable>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<ValidateSerializable<U, TSerializable>> : never;
export type ValidateSerializableMapped<T, TSerializable> = {
    [K in keyof T]: ValidateSerializable<T[K], TSerializable>;
};
declare const SERIALIZATION_ERROR: unique symbol;
export interface SerializationError<in out TMessage extends string> {
    [SERIALIZATION_ERROR]: TMessage;
}
export interface SerializationAdapter<TInput, TOutput, TExtendsAdapters extends ReadonlyArray<AnySerializationAdapter>> {
    '~types': SerializationAdapterTypes<TInput, TOutput, TExtendsAdapters>;
    key: string;
    extends?: TExtendsAdapters;
    test: (value: unknown) => value is TInput;
    toSerializable: (value: TInput) => TOutput;
    fromSerializable: (value: TOutput) => TInput;
}
export interface SerializationAdapterTypes<TInput, TOutput, TExtendsAdapters extends ReadonlyArray<AnySerializationAdapter>> {
    input: TInput | UnionizeSerializationAdaptersInput<TExtendsAdapters>;
    output: TOutput;
    extends: TExtendsAdapters;
}
export type AnySerializationAdapter = SerializationAdapter<any, any, any>;
export interface AdapterNode extends PluginInfo {
    v: SerovalNode;
}
/** Create a Seroval plugin for server-side serialization only. */
export declare function makeSsrSerovalPlugin(serializationAdapter: AnySerializationAdapter, options: {
    didRun: boolean;
}): Plugin<any, AdapterNode>;
/** Create a Seroval plugin for client/server symmetric (de)serialization. */
export declare function makeSerovalPlugin(serializationAdapter: AnySerializationAdapter): Plugin<any, AdapterNode>;
export type ValidateSerializableInput<TRegister, T> = ValidateSerializable<T, RegisteredSerializableInput<TRegister>>;
export type RegisteredSerializableInput<TRegister> = (unknown extends RegisteredSerializationAdapters<TRegister> ? never : RegisteredSerializationAdapters<TRegister> extends ReadonlyArray<AnySerializationAdapter> ? RegisteredSerializationAdapters<TRegister>[number]['~types']['input'] : never) | Serializable;
export type RegisteredSerializationAdapters<TRegister> = RegisteredConfigType<TRegister, 'serializationAdapters'>;
export type RegisteredSSROption<TRegister> = unknown extends RegisteredConfigType<TRegister, 'defaultSsr'> ? SSROption : RegisteredConfigType<TRegister, 'defaultSsr'>;
export type ValidateSerializableLifecycleResult<TRegister, TParentRoute extends AnyRoute, TSSR, TFn> = false extends RegisteredSsr<TRegister> ? any : ValidateSerializableLifecycleResultSSR<TRegister, TParentRoute, TSSR, TFn> extends infer TInput ? TInput : never;
export type ValidateSerializableLifecycleResultSSR<TRegister, TParentRoute extends AnyRoute, TSSR, TFn> = ResolveAllSSR<TParentRoute, TSSR> extends false ? any : RegisteredSSROption<TRegister> extends false ? any : ValidateSerializableInput<TRegister, LooseReturnType<TFn>>;
export type RegisteredReadableStream = unknown extends SerializerExtensions['ReadableStream'] ? never : SerializerExtensions['ReadableStream'];
export interface DefaultSerializerExtensions {
    ReadableStream: unknown;
}
export interface SerializerExtensions extends DefaultSerializerExtensions {
}
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/route.d.ts
import { LazyRoute } from './fileRoute.js';
import { NotFoundError } from './not-found.js';
import { RedirectFnRoute } from './redirect.js';
import { NavigateOptions, ParsePathParams } from './link.js';
import { ParsedLocation } from './location.js';
import { AnyRouteMatch, MakePreValidationErrorHandlingRouteMatchUnion, MakeRouteMatchFromRoute, MakeRouteMatchUnion, RouteMatch } from './Matches.js';
import { RootRouteId } from './root.js';
import { ParseRoute, RouteById, RouteIds, RoutePaths } from './routeInfo.js';
import { AnyRouter, Register, RegisteredRouter, SSROption } from './router.js';
import { BuildLocationFn, NavigateFn } from './RouterProvider.js';
import { Assign, Awaitable, Constrain, Expand, IntersectAssign, LooseAsyncReturnType, LooseReturnType, NoInfer } from './utils.js';
import { AnySchema, AnyStandardSchemaValidator, AnyValidator, AnyValidatorAdapter, AnyValidatorObj, DefaultValidator, ResolveSearchValidatorInput, ResolveValidatorOutput, StandardSchemaValidator, ValidatorAdapter, ValidatorFn, ValidatorObj } from './validators.js';
import { ValidateSerializableLifecycleResult } from './ssr/serializer/transformer.js';
export type AnyPathParams = {};
export type SearchSchemaInput = {
    __TSearchSchemaInput__: 'TSearchSchemaInput';
};
export type AnyContext = {};
export interface RouteContext {
}
export type PreloadableObj = {
    preload?: () => Promise<void>;
};
export type RoutePathOptions<TCustomId, TPath> = {
    path: TPath;
} | {
    id: TCustomId;
};
export interface StaticDataRouteOption {
}
export type RoutePathOptionsIntersection<TCustomId, TPath> = {
    path: TPath;
    id: TCustomId;
};
export type SearchFilter<TInput, TResult = TInput> = (prev: TInput) => TResult;
export type SearchMiddlewareMeta = {
    removed?: Map<string, unknown>;
    removedAny?: Set<string>;
    defaulted?: Map<string, unknown>;
    explicit?: unknown;
};
export type SearchMiddlewareContext<TSearchSchema> = {
    search: TSearchSchema;
    next: (newSearch: TSearchSchema) => TSearchSchema;
    meta?: SearchMiddlewareMeta;
};
export type SearchMiddleware<TSearchSchema> = (ctx: SearchMiddlewareContext<TSearchSchema>) => TSearchSchema;
export type ResolveId<TParentRoute, TCustomId extends string, TPath extends string> = TParentRoute extends {
    id: infer TParentId extends string;
} ? RoutePrefix<TParentId, string extends TCustomId ? TPath : TCustomId> : RootRouteId;
export type InferFullSearchSchema<TRoute> = TRoute extends {
    types: {
        fullSearchSchema: infer TFullSearchSchema;
    };
} ? TFullSearchSchema : {};
export type InferFullSearchSchemaInput<TRoute> = TRoute extends {
    types: {
        fullSearchSchemaInput: infer TFullSearchSchemaInput;
    };
} ? TFullSearchSchemaInput : {};
export type InferAllParams<TRoute> = TRoute extends {
    types: {
        allParams: infer TAllParams;
    };
} ? TAllParams : {};
export type InferAllContext<TRoute> = unknown extends TRoute ? TRoute : TRoute extends {
    types: {
        allContext: infer TAllContext;
    };
} ? TAllContext : {};
export type ResolveSearchSchemaFnInput<TSearchValidator> = TSearchValidator extends (input: infer TSearchSchemaInput) => any ? TSearchSchemaInput extends SearchSchemaInput ? Omit<TSearchSchemaInput, keyof SearchSchemaInput> : ResolveSearchSchemaFn<TSearchValidator> : AnySchema;
export type ResolveSearchSchemaInput<TSearchValidator> = TSearchValidator extends AnyStandardSchemaValidator ? NonNullable<TSearchValidator['~standard']['types']>['input'] : TSearchValidator extends AnyValidatorAdapter ? TSearchValidator['types']['input'] : TSearchValidator extends AnyValidatorObj ? ResolveSearchSchemaFnInput<TSearchValidator['parse']> : ResolveSearchSchemaFnInput<TSearchValidator>;
export type ResolveSearchSchemaFn<TSearchValidator> = TSearchValidator extends (...args: any) => infer TSearchSchema ? TSearchSchema : AnySchema;
export type ResolveSearchSchema<TSearchValidator> = unknown extends TSearchValidator ? TSearchValidator : TSearchValidator extends AnyStandardSchemaValidator ? NonNullable<TSearchValidator['~standard']['types']>['output'] : TSearchValidator extends AnyValidatorAdapter ? TSearchValidator['types']['output'] : TSearchValidator extends AnyValidatorObj ? ResolveSearchSchemaFn<TSearchValidator['parse']> : ResolveSearchSchemaFn<TSearchValidator>;
export type ResolveRequiredParams<TPath extends string, T> = {
    [K in ParsePathParams<TPath>['required']]: T;
};
export type ResolveOptionalParams<TPath extends string, T> = {
    [K in ParsePathParams<TPath>['optional']]?: T | undefined;
};
export type ResolveParams<TPath extends string, T = string> = ResolveRequiredParams<TPath, T> & ResolveOptionalParams<TPath, T>;
export type ParseParamsFn<in out TPath extends string, in out TParams> = (rawParams: Expand<ResolveParams<TPath>>) => TParams | false;
type ValidateParsedParams<TPath extends string, TParams> = [TParams] extends [
    ResolveParams<TPath, any>
] ? unknown : never;
export type StringifyParamsFn<in out TPath extends string, in out TParams> = (params: TParams) => ResolveParams<TPath>;
export type ParamsOptions<in out TPath extends string, in out TParams> = {
    params?: {
        parse?: ParseParamsFn<TPath, TParams> & ValidateParsedParams<TPath, TParams>;
        /**
         * When multiple route candidates use `params.parse` during matching,
         * higher priorities are tried first.
         *
         * @default 0
         */
        priority?: number;
        stringify?: StringifyParamsFn<TPath, TParams>;
    };
    /**
    @deprecated Use params.parse instead
    */
    parseParams?: ParseParamsFn<TPath, TParams> & ValidateParsedParams<TPath, TParams>;
    /**
    @deprecated Use params.stringify instead
    */
    stringifyParams?: StringifyParamsFn<TPath, TParams>;
};
interface RequiredStaticDataRouteOption {
    staticData: StaticDataRouteOption;
}
interface OptionalStaticDataRouteOption {
    staticData?: StaticDataRouteOption;
}
export type UpdatableStaticRouteOption = {} extends StaticDataRouteOption ? OptionalStaticDataRouteOption : RequiredStaticDataRouteOption;
export type MetaDescriptor = {
    charSet: 'utf-8';
} | {
    title: string;
} | {
    name: string;
    content: string;
} | {
    property: string;
    content: string;
} | {
    httpEquiv: string;
    content: string;
} | {
    'script:ld+json': LdJsonObject;
} | {
    tagName: 'meta' | 'link';
    [name: string]: string;
} | Record<string, unknown>;
type LdJsonObject = {
    [Key in string]: LdJsonValue;
} & {
    [Key in string]?: LdJsonValue | undefined;
};
type LdJsonArray = Array<LdJsonValue> | ReadonlyArray<LdJsonValue>;
type LdJsonPrimitive = string | number | boolean | null;
type LdJsonValue = LdJsonPrimitive | LdJsonObject | LdJsonArray;
export type RouteLinkEntry = {};
export type SearchValidator<TInput, TOutput> = ValidatorObj<TInput, TOutput> | ValidatorFn<TInput, TOutput> | ValidatorAdapter<TInput, TOutput> | StandardSchemaValidator<TInput, TOutput> | undefined;
export type AnySearchValidator = SearchValidator<any, any>;
export type DefaultSearchValidator = SearchValidator<Record<string, unknown>, AnySchema>;
export type RoutePrefix<TPrefix extends string, TPath extends string> = string extends TPath ? RootRouteId : TPath extends string ? TPrefix extends RootRouteId ? TPath extends '/' ? '/' : `/${TrimPath<TPath>}` : `${TPrefix}/${TPath}` extends '/' ? '/' : `/${TrimPathLeft<`${TrimPathRight<TPrefix>}/${TrimPath<TPath>}`>}` : never;
export type TrimPath<T extends string> = '' extends T ? '' : TrimPathRight<TrimPathLeft<T>>;
export type TrimPathLeft<T extends string> = T extends `${RootRouteId}/${infer U}` ? TrimPathLeft<U> : T extends `/${infer U}` ? TrimPathLeft<U> : T;
export type TrimPathRight<T extends string> = T extends '/' ? '/' : T extends `${infer U}/` ? TrimPathRight<U> : T;
export type ContextReturnType<TContextFn> = unknown extends TContextFn ? TContextFn : LooseReturnType<TContextFn> extends never ? AnyContext : LooseReturnType<TContextFn>;
export type ContextAsyncReturnType<TContextFn> = unknown extends TContextFn ? TContextFn : LooseAsyncReturnType<TContextFn> extends never ? AnyContext : LooseAsyncReturnType<TContextFn>;
export type ResolveRouteContext<TRouteContextFn, TBeforeLoadFn> = Assign<ContextReturnType<TRouteContextFn>, ContextAsyncReturnType<TBeforeLoadFn>>;
export type ResolveRouteLoaderFn<TLoaderFn> = TLoaderFn extends {
    handler: infer THandler;
} ? THandler : TLoaderFn;
export type RouteLoaderObject<TRegister, TParentRoute extends AnyRoute = AnyRoute, TId extends string = string, TParams = {}, TLoaderDeps = {}, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TServerMiddlewares = unknown, THandlers = undefined> = {
    handler: RouteLoaderFn<TRegister, TParentRoute, TId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn, TServerMiddlewares, THandlers>;
    staleReloadMode?: LoaderStaleReloadMode;
};
export type ResolveLoaderData<TLoaderFn> = unknown extends TLoaderFn ? TLoaderFn : LooseAsyncReturnType<ResolveRouteLoaderFn<TLoaderFn>> extends never ? undefined : LooseAsyncReturnType<ResolveRouteLoaderFn<TLoaderFn>>;
export type ResolveFullSearchSchema<TParentRoute extends AnyRoute, TSearchValidator> = unknown extends TParentRoute ? ResolveValidatorOutput<TSearchValidator> : IntersectAssign<InferFullSearchSchema<TParentRoute>, ResolveValidatorOutput<TSearchValidator>>;
export type ResolveFullSearchSchemaInput<TParentRoute extends AnyRoute, TSearchValidator> = IntersectAssign<InferFullSearchSchemaInput<TParentRoute>, ResolveSearchValidatorInput<TSearchValidator>>;
export type ResolveAllParamsFromParent<TParentRoute extends AnyRoute, TParams> = Assign<InferAllParams<TParentRoute>, TParams>;
export type RouteContextParameter<TParentRoute extends AnyRoute, TRouterContext> = unknown extends TParentRoute ? TRouterContext : Assign<TRouterContext, InferAllContext<TParentRoute>>;
export type BeforeLoadContextParameter<TParentRoute extends AnyRoute, TRouterContext, TRouteContextFn> = Assign<RouteContextParameter<TParentRoute, TRouterContext>, ContextReturnType<TRouteContextFn>>;
export type ResolveAllContext<TParentRoute extends AnyRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn> = Assign<BeforeLoadContextParameter<TParentRoute, TRouterContext, TRouteContextFn>, ContextAsyncReturnType<TBeforeLoadFn>>;
export interface FullSearchSchemaOption<in out TParentRoute extends AnyRoute, in out TSearchValidator> {
    search: Expand<ResolveFullSearchSchema<TParentRoute, TSearchValidator>>;
}
export interface RemountDepsOptions<in out TRouteId, in out TFullSearchSchema, in out TAllParams, in out TLoaderDeps> {
    routeId: TRouteId;
    search: TFullSearchSchema;
    params: TAllParams;
    loaderDeps: TLoaderDeps;
}
export type MakeRemountDepsOptionsUnion<TRouteTree extends AnyRoute = RegisteredRouter['routeTree']> = ParseRoute<TRouteTree> extends infer TRoute extends AnyRoute ? TRoute extends any ? RemountDepsOptions<TRoute['id'], TRoute['types']['fullSearchSchema'], TRoute['types']['allParams'], TRoute['types']['loaderDeps']> : never : never;
export interface RouteTypes<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps, in out TLoaderFn, in out TChildren, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> {
    parentRoute: TParentRoute;
    path: TPath;
    to: TrimPathRight<TFullPath>;
    fullPath: TFullPath;
    customId: TCustomId;
    id: TId;
    searchSchema: ResolveValidatorOutput<TSearchValidator>;
    searchSchemaInput: ResolveSearchValidatorInput<TSearchValidator>;
    searchValidator: TSearchValidator;
    fullSearchSchema: ResolveFullSearchSchema<TParentRoute, TSearchValidator>;
    fullSearchSchemaInput: ResolveFullSearchSchemaInput<TParentRoute, TSearchValidator>;
    params: TParams;
    allParams: ResolveAllParamsFromParent<TParentRoute, TParams>;
    routerContext: TRouterContext;
    routeContext: ResolveRouteContext<TRouteContextFn, TBeforeLoadFn>;
    routeContextFn: TRouteContextFn;
    beforeLoadFn: TBeforeLoadFn;
    allContext: ResolveAllContext<TParentRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn>;
    children: TChildren;
    loaderData: ResolveLoaderData<TLoaderFn>;
    loaderDeps: TLoaderDeps;
    fileRouteTypes: TFileRouteTypes;
    ssr: ResolveSSR<TSSR>;
    allSsr: ResolveAllSSR<TParentRoute, TSSR>;
}
export type ResolveSSR<TSSR> = TSSR extends (...args: ReadonlyArray<any>) => any ? LooseReturnType<TSSR> : TSSR;
export type ResolveAllSSR<TParentRoute extends AnyRoute, TSSR> = unknown extends TParentRoute ? ResolveSSR<TSSR> : unknown extends TSSR ? TParentRoute['types']['allSsr'] : ResolveSSR<TSSR>;
export type ResolveFullPath<TParentRoute extends AnyRoute, TPath extends string, TPrefixed = RoutePrefix<TParentRoute['fullPath'], TPath>> = TPrefixed extends RootRouteId ? '/' : TPrefixed;
export interface RouteExtensions<in out TId, in out TFullPath> {
    id: TId;
    fullPath: TFullPath;
}
export type RouteLazyFn<TRoute extends AnyRoute> = (lazyFn: () => Promise<LazyRoute<TRoute>>) => TRoute;
export type RouteAddChildrenFn<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps extends Record<string, any>, in out TLoaderFn, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> = <const TNewChildren>(children: Constrain<TNewChildren, ReadonlyArray<AnyRoute> | Record<string, AnyRoute>>) => Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TNewChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
export type RouteAddFileChildrenFn<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps extends Record<string, any>, in out TLoaderFn, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> = <const TNewChildren>(children: TNewChildren) => Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TNewChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
export type RouteAddFileTypesFn<TRegister, TParentRoute extends AnyRoute, TPath extends string, TFullPath extends string, TCustomId extends string, TId extends string, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps extends Record<string, any>, TLoaderFn, TChildren, TSSR, TServerMiddlewares, THandlers> = <TNewFileRouteTypes>() => Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TNewFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
export interface Route<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps extends Record<string, any>, in out TLoaderFn, in out TChildren, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> extends RouteExtensions<TId, TFullPath> {
    path: TPath;
    parentRoute: TParentRoute;
    children?: TChildren;
    types: RouteTypes<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    options: RouteOptions<TRegister, TParentRoute, TId, TCustomId, TFullPath, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares, THandlers>;
    isRoot: TParentRoute extends AnyRoute ? true : false;
    lazyFn?: () => Promise<LazyRoute<Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>>>;
    rank: number;
    to: TrimPathRight<TFullPath>;
    init: (opts: {
        originalIndex: number;
    }) => void;
    update: (options: UpdatableRouteOptions<TParentRoute, TCustomId, TFullPath, TParams, TSearchValidator, TLoaderFn, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn>) => this;
    lazy: RouteLazyFn<Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>>;
    addChildren: RouteAddChildrenFn<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    _addFileChildren: RouteAddFileChildrenFn<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    _addFileTypes: RouteAddFileTypesFn<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TSSR, TServerMiddlewares, THandlers>;
    /**
     * Create a redirect with `from` automatically set to this route's path.
     * Enables relative redirects like `Route.redirect({ to: './overview' })`.
     * @param opts Redirect options (same as `redirect()` but without `from`)
     * @returns A redirect Response that can be thrown from loaders/beforeLoad
     * @link https://tanstack.com/router/latest/docs/framework/react/api/router/redirectFunction
     */
    redirect: RedirectFnRoute<TFullPath>;
}
export type AnyRoute = Route<any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any>;
export type AnyRouteWithContext<TContext> = AnyRoute & {
    types: {
        allContext: TContext;
    };
};
export type RouteOptions<TRegister, TParentRoute extends AnyRoute = AnyRoute, TId extends string = string, TCustomId extends string = string, TFullPath extends string = string, TPath extends string = string, TSearchValidator = undefined, TParams = AnyPathParams, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TSSR = unknown, TServerMiddlewares = unknown, THandlers = undefined> = BaseRouteOptions<TRegister, TParentRoute, TId, TCustomId, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares, THandlers> & UpdatableRouteOptions<NoInfer<TParentRoute>, NoInfer<TCustomId>, NoInfer<TFullPath>, NoInfer<TParams>, NoInfer<TSearchValidator>, NoInfer<TLoaderFn>, NoInfer<TLoaderDeps>, NoInfer<TRouterContext>, NoInfer<TRouteContextFn>, NoInfer<TBeforeLoadFn>>;
export type RouteContextFn<in out TParentRoute extends AnyRoute, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteId> = (ctx: RouteContextOptions<TParentRoute, TSearchValidator, TParams, TRouterContext, TRouteId>) => any;
export type FileBaseRouteOptions<TRegister, TParentRoute extends AnyRoute = AnyRoute, TId extends string = string, TPath extends string = string, TSearchValidator = undefined, TParams = {}, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TRemountDepsFn = AnyContext, TSSR = unknown, TServerMiddlewares = unknown, THandlers = undefined> = ParamsOptions<TPath, TParams> & FilebaseRouteOptionsInterface<TRegister, TParentRoute, TId, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TRemountDepsFn, TSSR, TServerMiddlewares, THandlers>;
export interface FilebaseRouteOptionsInterface<TRegister, TParentRoute extends AnyRoute = AnyRoute, TId extends string = string, TPath extends string = string, TSearchValidator = undefined, TParams = {}, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TRemountDepsFn = AnyContext, TSSR = unknown, TServerMiddlewares = unknown, THandlers = undefined> {
    validateSearch?: Constrain<TSearchValidator, AnyValidator, DefaultValidator>;
    shouldReload?: boolean | ((match: LoaderFnContext<TRegister, TParentRoute, TId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn, TServerMiddlewares, THandlers>) => any);
    context?: Constrain<TRouteContextFn, (ctx: RouteContextOptions<TParentRoute, TParams, TRouterContext, TLoaderDeps, TId>) => any>;
    ssr?: Constrain<TSSR, undefined | SSROption | ((ctx: SsrContextOptions<TParentRoute, TSearchValidator, TParams>) => Awaitable<undefined | SSROption>)>;
    beforeLoad?: Constrain<TBeforeLoadFn, (ctx: BeforeLoadContextOptions<TRegister, TParentRoute, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TId, TServerMiddlewares, THandlers>) => ValidateSerializableLifecycleResult<TRegister, TParentRoute, TSSR, TBeforeLoadFn>>;
    loaderDeps?: (opts: FullSearchSchemaOption<TParentRoute, TSearchValidator>) => TLoaderDeps;
    remountDeps?: Constrain<TRemountDepsFn, (opt: RemountDepsOptions<TId, ResolveFullSearchSchema<TParentRoute, TSearchValidator>, Expand<ResolveAllParamsFromParent<TParentRoute, TParams>>, TLoaderDeps>) => any>;
    loader?: Constrain<TLoaderFn, RouteLoaderFn<TRegister, TParentRoute, TId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn, TServerMiddlewares, THandlers> | RouteLoaderObject<TRegister, TParentRoute, TId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn, TServerMiddlewares, THandlers>>;
}
export type BaseRouteOptions<TRegister, TParentRoute extends AnyRoute = AnyRoute, TId extends string = string, TCustomId extends string = string, TPath extends string = string, TSearchValidator = undefined, TParams = {}, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TSSR = unknown, TServerMiddlewares = unknown, THandlers = undefined> = RoutePathOptions<TCustomId, TPath> & FileBaseRouteOptions<TRegister, TParentRoute, TId, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, AnyContext, TSSR, TServerMiddlewares, THandlers> & {
    getParentRoute: () => TParentRoute;
};
export interface ContextOptions<in out TParentRoute extends AnyRoute, in out TParams, in out TRouteId> {
    abortController: AbortController;
    preload: boolean;
    params: Expand<ResolveAllParamsFromParent<TParentRoute, TParams>>;
    location: ParsedLocation;
    /**
     * @deprecated Use `throw redirect({ to: '/somewhere' })` instead
     **/
    navigate: NavigateFn;
    buildLocation: BuildLocationFn;
    cause: 'preload' | 'enter' | 'stay';
    matches: Array<MakeRouteMatchUnion>;
    routeId: TRouteId;
}
export interface RouteContextOptions<in out TParentRoute extends AnyRoute, in out TParams, in out TRouterContext, in out TLoaderDeps, in out TRouteId> extends ContextOptions<TParentRoute, TParams, TRouteId> {
    deps: TLoaderDeps;
    context: Expand<RouteContextParameter<TParentRoute, TRouterContext>>;
}
export interface SsrContextOptions<in out TParentRoute extends AnyRoute, in out TSearchValidator, in out TParams> {
    params: {
        status: 'success';
        value: Expand<ResolveAllParamsFromParent<TParentRoute, TParams>>;
    } | {
        status: 'error';
        error: unknown;
    };
    search: {
        status: 'success';
        value: Expand<ResolveFullSearchSchema<TParentRoute, TSearchValidator>>;
    } | {
        status: 'error';
        error: unknown;
    };
    location: ParsedLocation;
    matches: Array<MakePreValidationErrorHandlingRouteMatchUnion>;
}
export interface BeforeLoadContextOptions<in out TRegister, in out TParentRoute extends AnyRoute, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TRouteId, in out TServerMiddlewares, in out THandlers> extends ContextOptions<TParentRoute, TParams, TRouteId>, FullSearchSchemaOption<TParentRoute, TSearchValidator> {
    context: Expand<BeforeLoadContextParameter<TParentRoute, TRouterContext, TRouteContextFn>>;
}
type AssetFnContextOptions<in out TRouteId, in out TFullPath, in out TParentRoute extends AnyRoute, in out TParams, in out TSearchValidator, in out TLoaderFn, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps> = {
    ssr?: {
        nonce?: string;
    };
    matches: Array<RouteMatch<TRouteId, TFullPath, ResolveAllParamsFromParent<TParentRoute, TParams>, ResolveFullSearchSchema<TParentRoute, TSearchValidator>, ResolveLoaderData<TLoaderFn>, ResolveAllContext<TParentRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn>, TLoaderDeps>>;
    match: RouteMatch<TRouteId, TFullPath, ResolveAllParamsFromParent<TParentRoute, TParams>, ResolveFullSearchSchema<TParentRoute, TSearchValidator>, ResolveLoaderData<TLoaderFn>, ResolveAllContext<TParentRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn>, TLoaderDeps>;
    params: ResolveAllParamsFromParent<TParentRoute, TParams>;
    loaderData?: ResolveLoaderData<TLoaderFn>;
};
export interface DefaultUpdatableRouteOptionsExtensions {
    component?: unknown;
    errorComponent?: unknown;
    notFoundComponent?: unknown;
    pendingComponent?: unknown;
}
export interface UpdatableRouteOptionsExtensions extends DefaultUpdatableRouteOptionsExtensions {
}
export interface UpdatableRouteOptions<in out TParentRoute extends AnyRoute, in out TRouteId, in out TFullPath, in out TParams, in out TSearchValidator, in out TLoaderFn, in out TLoaderDeps, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn> extends UpdatableStaticRouteOption, UpdatableRouteOptionsExtensions {
    /**
     * If true, this route will be matched as case-sensitive
     *
     * @default false
     */
    caseSensitive?: boolean;
    /**
     * If true, this route will be forcefully wrapped in a suspense boundary
     */
    wrapInSuspense?: boolean;
    pendingMs?: number;
    pendingMinMs?: number;
    staleTime?: number;
    gcTime?: number;
    preload?: boolean;
    preloadStaleTime?: number;
    preloadGcTime?: number;
    search?: {
        middlewares?: Array<SearchMiddleware<ResolveFullSearchSchema<TParentRoute, TSearchValidator>>>;
    };
    /**
    @deprecated Use search.middlewares instead
    */
    preSearchFilters?: Array<SearchFilter<ResolveFullSearchSchema<TParentRoute, TSearchValidator>>>;
    /**
    @deprecated Use search.middlewares instead
    */
    postSearchFilters?: Array<SearchFilter<ResolveFullSearchSchema<TParentRoute, TSearchValidator>>>;
    onCatch?: (error: Error) => void;
    onError?: (err: any) => void;
    onEnter?: (match: RouteMatch<TRouteId, TFullPath, ResolveAllParamsFromParent<TParentRoute, TParams>, ResolveFullSearchSchema<TParentRoute, TSearchValidator>, ResolveLoaderData<TLoaderFn>, ResolveAllContext<TParentRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn>, TLoaderDeps>) => void;
    onStay?: (match: RouteMatch<TRouteId, TFullPath, ResolveAllParamsFromParent<TParentRoute, TParams>, ResolveFullSearchSchema<TParentRoute, TSearchValidator>, ResolveLoaderData<TLoaderFn>, ResolveAllContext<TParentRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn>, TLoaderDeps>) => void;
    onLeave?: (match: RouteMatch<TRouteId, TFullPath, ResolveAllParamsFromParent<TParentRoute, TParams>, ResolveFullSearchSchema<TParentRoute, TSearchValidator>, ResolveLoaderData<TLoaderFn>, ResolveAllContext<TParentRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn>, TLoaderDeps>) => void;
    headers?: (ctx: AssetFnContextOptions<TRouteId, TFullPath, TParentRoute, TParams, TSearchValidator, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps>) => Awaitable<Record<string, string> | undefined>;
    head?: (ctx: AssetFnContextOptions<TRouteId, TFullPath, TParentRoute, TParams, TSearchValidator, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps>) => Awaitable<{
        links?: AnyRouteMatch['links'];
        scripts?: AnyRouteMatch['headScripts'];
        meta?: AnyRouteMatch['meta'];
        styles?: AnyRouteMatch['styles'];
    }>;
    scripts?: (ctx: AssetFnContextOptions<TRouteId, TFullPath, TParentRoute, TParams, TSearchValidator, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps>) => Awaitable<AnyRouteMatch['scripts']>;
    codeSplitGroupings?: Array<Array<'loader' | 'component' | 'pendingComponent' | 'notFoundComponent' | 'errorComponent'>>;
}
export type RouteLoaderFn<in out TRegister, in out TParentRoute extends AnyRoute = AnyRoute, in out TId extends string = string, in out TParams = {}, in out TLoaderDeps = {}, in out TRouterContext = {}, in out TRouteContextFn = AnyContext, in out TBeforeLoadFn = AnyContext, in out TServerMiddlewares = unknown, in out THandlers = undefined> = (match: LoaderFnContext<TRegister, TParentRoute, TId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn, TServerMiddlewares, THandlers>) => any;
export type LoaderStaleReloadMode = 'background' | 'blocking';
export type RouteLoaderEntry<TRegister, TParentRoute extends AnyRoute = AnyRoute, TId extends string = string, TParams = {}, TLoaderDeps = {}, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TServerMiddlewares = unknown, THandlers = undefined> = RouteLoaderFn<TRegister, TParentRoute, TId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn, TServerMiddlewares, THandlers> | RouteLoaderObject<TRegister, TParentRoute, TId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn, TServerMiddlewares, THandlers>;
export interface LoaderFnContext<in out TRegister = unknown, in out TParentRoute extends AnyRoute = AnyRoute, in out TId extends string = string, in out TParams = {}, in out TLoaderDeps = {}, in out TRouterContext = {}, in out TRouteContextFn = AnyContext, in out TBeforeLoadFn = AnyContext, in out TServerMiddlewares = unknown, in out THandlers = undefined> {
    abortController: AbortController;
    preload: boolean;
    params: Expand<ResolveAllParamsFromParent<TParentRoute, TParams>>;
    deps: TLoaderDeps;
    context: Expand<ResolveAllContext<TParentRoute, TRouterContext, TRouteContextFn, TBeforeLoadFn>>;
    location: ParsedLocation;
    /**
     * @deprecated Use `throw redirect({ to: '/somewhere' })` instead
     **/
    navigate: (opts: NavigateOptions<AnyRouter>) => Promise<void> | void;
    parentMatchPromise: TId extends RootRouteId ? never : Promise<MakeRouteMatchFromRoute<TParentRoute>>;
    cause: 'preload' | 'enter' | 'stay';
    route: AnyRoute;
}
export interface DefaultRootRouteOptionsExtensions {
    shellComponent?: unknown;
}
export interface RootRouteOptionsExtensions extends DefaultRootRouteOptionsExtensions {
}
export interface RootRouteOptions<TRegister = unknown, TSearchValidator = undefined, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TSSR = unknown, TServerMiddlewares = unknown, THandlers = undefined> extends Omit<RouteOptions<TRegister, any, // TParentRoute
RootRouteId, // TId
RootRouteId, // TCustomId
'', // TFullPath
'', // TPath
TSearchValidator, {}, // TParams
TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares, THandlers>, 'path' | 'id' | 'getParentRoute' | 'caseSensitive' | 'parseParams' | 'stringifyParams' | 'params'>, RootRouteOptionsExtensions {
}
export type RouteConstraints = {
    TParentRoute: AnyRoute;
    TPath: string;
    TFullPath: string;
    TCustomId: string;
    TId: string;
    TSearchSchema: AnySchema;
    TFullSearchSchema: AnySchema;
    TParams: Record<string, any>;
    TAllParams: Record<string, any>;
    TParentContext: AnyContext;
    TRouteContext: RouteContext;
    TAllContext: AnyContext;
    TRouterContext: AnyContext;
    TChildren: unknown;
    TRouteTree: AnyRoute;
};
export type RouteTypesById<TRouter extends AnyRouter, TId> = RouteById<TRouter['routeTree'], TId>['types'];
export type RouteMask<TRouteTree extends AnyRoute> = {
    routeTree: TRouteTree;
    from: RoutePaths<TRouteTree>;
    to?: any;
    params?: any;
    search?: any;
    hash?: any;
    state?: any;
    unmaskOnReload?: boolean;
};
/**
 * @deprecated Use `ErrorComponentProps` instead.
 */
export type ErrorRouteProps = {
    error: unknown;
    info?: {
        componentStack: string;
    };
    reset: () => void;
};
export type ErrorComponentProps<TError = Error> = {
    error: TError;
    info?: {
        componentStack: string;
    };
    reset: () => void;
};
export type NotFoundRouteProps = {
    data?: unknown;
    isNotFound: boolean;
    routeId: RouteIds<RegisteredRouter['routeTree']>;
};
export declare class BaseRoute<in out TRegister = Register, in out TParentRoute extends AnyRoute = AnyRoute, in out TPath extends string = '/', in out TFullPath extends string = ResolveFullPath<TParentRoute, TPath>, in out TCustomId extends string = string, in out TId extends string = ResolveId<TParentRoute, TCustomId, TPath>, in out TSearchValidator = undefined, in out TParams = ResolveParams<TPath>, in out TRouterContext = AnyContext, in out TRouteContextFn = AnyContext, in out TBeforeLoadFn = AnyContext, in out TLoaderDeps extends Record<string, any> = {}, in out TLoaderFn = undefined, in out TChildren = unknown, in out TFileRouteTypes = unknown, in out TSSR = unknown, in out TServerMiddlewares = unknown, in out THandlers = undefined> {
    isRoot: TParentRoute extends AnyRoute ? true : false;
    options: RouteOptions<TRegister, TParentRoute, TId, TCustomId, TFullPath, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares, THandlers>;
    parentRoute: TParentRoute;
    private _id;
    private _path;
    private _fullPath;
    private _to;
    get to(): TrimPathRight<TFullPath>;
    get id(): TId;
    get path(): TPath;
    get fullPath(): TFullPath;
    children?: TChildren;
    originalIndex?: number;
    rank: number;
    lazyFn?: () => Promise<LazyRoute<Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>>>;
    constructor(options?: RouteOptions<TRegister, TParentRoute, TId, TCustomId, TFullPath, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares, THandlers>);
    types: RouteTypes<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    init: (opts: {
        originalIndex: number;
    }) => void;
    addChildren: RouteAddChildrenFn<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    _addFileChildren: RouteAddFileChildrenFn<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    _addFileTypes: RouteAddFileTypesFn<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TSSR, TServerMiddlewares, THandlers>;
    updateLoader: <TNewLoaderFn>(options: {
        loader: Constrain<TNewLoaderFn, RouteLoaderFn<TRegister, TParentRoute, TCustomId, TParams, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn>>;
    }) => BaseRoute<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TNewLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    update: (options: UpdatableRouteOptions<TParentRoute, TCustomId, TFullPath, TParams, TSearchValidator, TLoaderFn, TLoaderDeps, TRouterContext, TRouteContextFn, TBeforeLoadFn>) => this;
    lazy: RouteLazyFn<Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>>;
    /**
     * Create a redirect with `from` automatically set to this route's fullPath.
     * Enables relative redirects like `Route.redirect({ to: './overview' })`.
     * @param opts Redirect options (same as `redirect()` but without `from`)
     * @returns A redirect Response that can be thrown from loaders/beforeLoad
     * @link https://tanstack.com/router/latest/docs/framework/react/api/router/redirectFunction
     */
    redirect: RedirectFnRoute<TFullPath>;
}
export declare class BaseRouteApi<TId, TRouter extends AnyRouter = RegisteredRouter> {
    id: TId;
    constructor({ id }: {
        id: TId;
    });
    notFound: (opts?: NotFoundError) => NotFoundError;
    /**
     * Create a redirect with `from` automatically set to this route's path.
     * Enables relative redirects like `routeApi.redirect({ to: './overview' })`.
     * @param opts Redirect options (same as `redirect()` but without `from`)
     * @returns A redirect Response that can be thrown from loaders/beforeLoad
     * @link https://tanstack.com/router/latest/docs/framework/react/api/router/redirectFunction
     */
    redirect: RedirectFnRoute<RouteTypesById<TRouter, TId>['fullPath']>;
}
export interface RootRoute<in out TRegister, in out TSearchValidator = undefined, in out TRouterContext = {}, in out TRouteContextFn = AnyContext, in out TBeforeLoadFn = AnyContext, in out TLoaderDeps extends Record<string, any> = {}, in out TLoaderFn = undefined, in out TChildren = unknown, in out TFileRouteTypes = unknown, in out TSSR = unknown, in out TServerMiddlewares = unknown, in out THandlers = undefined> extends Route<TRegister, any, // TParentRoute
'/', // TPath
'/', // TFullPath
string, // TCustomId
RootRouteId, // TId
TSearchValidator, // TSearchValidator
{}, // TParams
TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, // TChildren
TFileRouteTypes, TSSR, TServerMiddlewares, THandlers> {
}
export declare class BaseRootRoute<in out TRegister = Register, in out TSearchValidator = undefined, in out TRouterContext = {}, in out TRouteContextFn = AnyContext, in out TBeforeLoadFn = AnyContext, in out TLoaderDeps extends Record<string, any> = {}, in out TLoaderFn = undefined, in out TChildren = unknown, in out TFileRouteTypes = unknown, in out TSSR = unknown, in out TServerMiddlewares = unknown, in out THandlers = undefined> extends BaseRoute<TRegister, any, // TParentRoute
'/', // TPath
'/', // TFullPath
string, // TCustomId
RootRouteId, // TId
TSearchValidator, // TSearchValidator
{}, // TParams
TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, // TChildren
TFileRouteTypes, TSSR, TServerMiddlewares, THandlers> {
    constructor(options?: RootRouteOptions<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TSSR, TServerMiddlewares, THandlers>);
}
export interface RouteLike {
    id: string;
    isRoot?: boolean;
    path?: string;
    fullPath: string;
    rank?: number;
    parentRoute?: RouteLike;
    children?: Array<RouteLike>;
    options?: {
        caseSensitive?: boolean;
    };
}
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/validators.d.ts
import { SearchSchemaInput } from './route.js';
export interface StandardSchemaValidatorProps<TInput, TOutput> {
    readonly types?: StandardSchemaValidatorTypes<TInput, TOutput> | undefined;
    readonly validate: AnyStandardSchemaValidate;
}
export interface StandardSchemaValidator<TInput, TOutput> {
    readonly '~standard': StandardSchemaValidatorProps<TInput, TOutput>;
}
export type AnyStandardSchemaValidator = StandardSchemaValidator<any, any>;
export interface StandardSchemaValidatorTypes<TInput, TOutput> {
    readonly input: TInput;
    readonly output: TOutput;
}
export interface AnyStandardSchemaValidateSuccess {
    readonly value: any;
    readonly issues?: undefined;
}
export interface AnyStandardSchemaValidateFailure {
    readonly issues: ReadonlyArray<AnyStandardSchemaValidateIssue>;
}
export interface AnyStandardSchemaValidateIssue {
    readonly message: string;
}
export interface AnyStandardSchemaValidateInput {
    readonly value: any;
}
export type AnyStandardSchemaValidate = (value: unknown) => (AnyStandardSchemaValidateSuccess | AnyStandardSchemaValidateFailure) | Promise<AnyStandardSchemaValidateSuccess | AnyStandardSchemaValidateFailure>;
export interface ValidatorObj<TInput, TOutput> {
    parse: ValidatorFn<TInput, TOutput>;
}
export type AnyValidatorObj = ValidatorObj<any, any>;
export interface ValidatorAdapter<TInput, TOutput> {
    types: {
        input: TInput;
        output: TOutput;
    };
    parse: (input: unknown) => TOutput;
}
export type AnyValidatorAdapter = ValidatorAdapter<any, any>;
export type AnyValidatorFn = ValidatorFn<any, any>;
export type ValidatorFn<TInput, TOutput> = (input: TInput) => TOutput;
export type Validator<TInput, TOutput> = ValidatorObj<TInput, TOutput> | ValidatorFn<TInput, TOutput> | ValidatorAdapter<TInput, TOutput> | StandardSchemaValidator<TInput, TOutput> | undefined;
export type AnyValidator = Validator<any, any>;
export type AnySchema = {};
export type DefaultValidator = Validator<Record<string, unknown>, AnySchema>;
export type ResolveSearchValidatorInputFn<TValidator> = TValidator extends (input: infer TSchemaInput) => any ? TSchemaInput extends SearchSchemaInput ? Omit<TSchemaInput, keyof SearchSchemaInput> : ResolveValidatorOutputFn<TValidator> : AnySchema;
export type ResolveSearchValidatorInput<TValidator> = TValidator extends AnyStandardSchemaValidator ? NonNullable<TValidator['~standard']['types']>['input'] : TValidator extends AnyValidatorAdapter ? TValidator['types']['input'] : TValidator extends AnyValidatorObj ? ResolveSearchValidatorInputFn<TValidator['parse']> : ResolveSearchValidatorInputFn<TValidator>;
export type ResolveValidatorInputFn<TValidator> = TValidator extends (input: infer TInput) => any ? TInput : undefined;
export type ResolveValidatorInput<TValidator> = TValidator extends AnyStandardSchemaValidator ? NonNullable<TValidator['~standard']['types']>['input'] : TValidator extends AnyValidatorAdapter ? TValidator['types']['input'] : TValidator extends AnyValidatorObj ? ResolveValidatorInputFn<TValidator['parse']> : ResolveValidatorInputFn<TValidator>;
export type ResolveValidatorOutputFn<TValidator> = TValidator extends (...args: any) => infer TSchema ? TSchema : AnySchema;
export type ResolveValidatorOutput<TValidator> = unknown extends TValidator ? TValidator : TValidator extends AnyStandardSchemaValidator ? NonNullable<TValidator['~standard']['types']>['output'] : TValidator extends AnyValidatorAdapter ? TValidator['types']['output'] : TValidator extends AnyValidatorObj ? ResolveValidatorOutputFn<TValidator['parse']> : ResolveValidatorOutputFn<TValidator>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/location.d.ts
import { ParsedHistoryState } from '@tanstack/history';
import { AnySchema } from './validators.js';
export interface ParsedLocation<TSearchObj extends AnySchema = {}> {
    /**
     * The full path of the location, including pathname, search, and hash.
     * Does not include the origin. Is the equivalent of calling
     * `url.replace(url.origin, '')`
     */
    href: string;
    /**
     * @description The pathname of the location, including the leading slash.
     */
    pathname: string;
    /**
     * The parsed search parameters of the location in object form.
     */
    search: TSearchObj;
    /**
     * The search string of the location, including the leading question mark.
     */
    searchStr: string;
    /**
     * The in-memory state of the location as it *may* exist in the browser's history.
     */
    state: ParsedHistoryState;
    /**
     * The hash of the location, excluding the leading hash character.
     * (e.g., '123' instead of '#123')
     */
    hash: string;
    /**
     * The masked location of the location.
     */
    maskedLocation?: ParsedLocation<TSearchObj>;
    /**
     * Whether to unmask the location on reload.
     */
    unmaskOnReload?: boolean;
    /**
     * @private
     * @description The public href of the location.
     * If a rewrite is applied, the `href` property will be the rewritten URL.
     */
    publicHref: string;
    /**
     * @private
     * @description Whether the publicHref is external (different origin from rewrite).
     */
    external: boolean;
}


// @filename: /node_modules/@tanstack/router-core/dist/esm/load-matches.d.ts
import { ParsedLocation } from './location.js';
import { AnyRoute } from './route.js';
import { AnyRouteMatch, MakeRouteMatch } from './Matches.js';
import { AnyRouter, UpdateMatchFn } from './router.js';
export declare function loadMatches(arg: {
    router: AnyRouter;
    location: ParsedLocation;
    matches: Array<AnyRouteMatch>;
    preload?: boolean;
    forceStaleReload?: boolean;
    onReady?: () => Promise<void>;
    updateMatch: UpdateMatchFn;
    sync?: boolean;
}): Promise<Array<MakeRouteMatch>>;
export type RouteComponentType = 'component' | 'errorComponent' | 'pendingComponent' | 'notFoundComponent';
export declare function loadRouteChunk(route: AnyRoute, componentTypesToLoad?: Array<RouteComponentType>): Promise<void> | undefined;
export declare function routeNeedsPreload(route: AnyRoute): boolean;
export declare const componentTypes: Array<RouteComponentType>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/lru-cache.d.ts
export type LRUCache<TKey, TValue> = {
    get: (key: TKey) => TValue | undefined;
    set: (key: TKey, value: TValue) => void;
    clear: () => void;
};
export declare function createLRUCache<TKey, TValue>(max: number): LRUCache<TKey, TValue>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/new-process-route-tree.d.ts
import { LRUCache } from './lru-cache.js';
export declare const SEGMENT_TYPE_PATHNAME = 0;
export declare const SEGMENT_TYPE_PARAM = 1;
export declare const SEGMENT_TYPE_WILDCARD = 2;
export declare const SEGMENT_TYPE_OPTIONAL_PARAM = 3;
declare const SEGMENT_TYPE_INDEX = 4;
declare const SEGMENT_TYPE_PATHLESS = 5;
/**
 * All the kinds of segments that can be present in a route path.
 */
export type SegmentKind = typeof SEGMENT_TYPE_PATHNAME | typeof SEGMENT_TYPE_PARAM | typeof SEGMENT_TYPE_WILDCARD | typeof SEGMENT_TYPE_OPTIONAL_PARAM;
/**
 * All the kinds of segments that can be present in the segment tree.
 */
type ExtendedSegmentKind = SegmentKind | typeof SEGMENT_TYPE_INDEX | typeof SEGMENT_TYPE_PATHLESS;
type ParsedSegment = Uint16Array & {
    /** segment type (0 = pathname, 1 = param, 2 = wildcard, 3 = optional param) */
    0: SegmentKind;
    /** index of the end of the prefix */
    1: number;
    /** index of the start of the value */
    2: number;
    /** index of the end of the value */
    3: number;
    /** index of the start of the suffix */
    4: number;
    /** index of the end of the segment */
    5: number;
};
/**
 * Populates the `output` array with the parsed representation of the given `segment` string.
 *
 * Usage:
 * ```ts
 * let output
 * let cursor = 0
 * while (cursor < path.length) {
 *   output = parseSegment(path, cursor, output)
 *   const end = output[5]
 *   cursor = end + 1
 * ```
 *
 * `output` is stored outside to avoid allocations during repeated calls. It doesn't need to be typed
 * or initialized, it will be done automatically.
 */
export declare function parseSegment(
/** The full path string containing the segment. */
path: string, 
/** The starting index of the segment within the path. */
start: number, 
/** A Uint16Array (length: 6) to populate with the parsed segment data. */
output?: Uint16Array): ParsedSegment;
type StaticSegmentNode<T extends RouteLike> = SegmentNode<T> & {
    kind: typeof SEGMENT_TYPE_PATHNAME | typeof SEGMENT_TYPE_PATHLESS | typeof SEGMENT_TYPE_INDEX;
};
type DynamicSegmentNode<T extends RouteLike> = SegmentNode<T> & {
    kind: typeof SEGMENT_TYPE_PARAM | typeof SEGMENT_TYPE_WILDCARD | typeof SEGMENT_TYPE_OPTIONAL_PARAM;
    prefix?: string;
    suffix?: string;
    caseSensitive: boolean;
};
type AnySegmentNode<T extends RouteLike> = StaticSegmentNode<T> | DynamicSegmentNode<T>;
type SegmentNode<T extends RouteLike> = {
    kind: ExtendedSegmentKind;
    pathless: Array<StaticSegmentNode<T>> | null;
    /** Exact index segment (highest priority) */
    index: StaticSegmentNode<T> | null;
    /** Static segments (2nd priority) */
    static: Map<string, StaticSegmentNode<T>> | null;
    /** Case insensitive static segments (3rd highest priority) */
    staticInsensitive: Map<string, StaticSegmentNode<T>> | null;
    /** Dynamic segments ($param) */
    dynamic: Array<DynamicSegmentNode<T>> | null;
    /** Optional dynamic segments ({-$param}) */
    optional: Array<DynamicSegmentNode<T>> | null;
    /** Wildcard segments ($ - lowest priority) */
    wildcard: Array<DynamicSegmentNode<T>> | null;
    /** Terminal route (if this path can end here) */
    route: T | null;
    /** The full path for this segment node (will only be valid on leaf nodes) */
    fullPath: string;
    parent: AnySegmentNode<T> | null;
    depth: number;
    /** route.options.params.parse function, set on the last node of the route */
    parse: null | ((params: Record<string, string>) => unknown);
    /** route.options.params.priority ?? 0 */
    priority: number;
};
type RouteLike = {
    id?: string;
    path?: string;
    children?: Array<RouteLike>;
    parentRoute?: RouteLike;
    isRoot?: boolean;
    options?: {
        caseSensitive?: boolean;
        parseParams?: (params: Record<string, string>) => unknown;
        params?: {
            parse?: (params: Record<string, string>) => unknown;
            priority?: number;
        };
    };
} & ({
    fullPath: string;
    from?: never;
} | {
    fullPath?: never;
    from: string;
});
export type ProcessedTree<TTree extends Extract<RouteLike, {
    fullPath: string;
}>, TFlat extends Extract<RouteLike, {
    from: string;
}>, TSingle extends Extract<RouteLike, {
    from: string;
}>> = {
    /** a representation of the `routeTree` as a segment tree */
    segmentTree: AnySegmentNode<TTree>;
    /** a mini route tree generated from the flat `routeMasks` list */
    masksTree: AnySegmentNode<TFlat> | null;
    /** @deprecated keep until v2 so that `router.matchRoute` can keep not caring about the actual route tree */
    singleCache: LRUCache<string, AnySegmentNode<TSingle>>;
    /** a cache of route matches from the `segmentTree` */
    matchCache: LRUCache<string, RouteMatch<TTree> | null>;
    /** a cache of route matches from the `masksTree` */
    flatCache: LRUCache<string, ReturnType<typeof findMatch<TFlat>>> | null;
};
export declare function processRouteMasks<TRouteLike extends Extract<RouteLike, {
    from: string;
}>>(routeList: Array<TRouteLike>, processedTree: ProcessedTree<any, TRouteLike, any>): void;
/**
 * Take an arbitrary list of routes, create a tree from them (if it hasn't been created already), and match a path against it.
 */
export declare function findFlatMatch<T extends Extract<RouteLike, {
    from: string;
}>>(
/** The path to match. */
path: string, 
/** The `processedTree` returned by the initial `processRouteTree` call. */
processedTree: ProcessedTree<any, T, any>): {
    route: T;
    /**
     * The raw (unparsed) params extracted from the path.
     * This will be the exhaustive list of all params defined in the route's path.
     */
    rawParams: Record<string, string>;
} | null;
/**
 * @deprecated keep until v2 so that `router.matchRoute` can keep not caring about the actual route tree
 */
export declare function findSingleMatch(from: string, caseSensitive: boolean, fuzzy: boolean, path: string, processedTree: ProcessedTree<any, any, {
    from: string;
}>): {
    route: {
        from: string;
    };
    /**
     * The raw (unparsed) params extracted from the path.
     * This will be the exhaustive list of all params defined in the route's path.
     */
    rawParams: Record<string, string>;
} | null;
type RouteMatch<T extends Extract<RouteLike, {
    fullPath: string;
}>> = {
    route: T;
    rawParams: Record<string, string>;
    branch: ReadonlyArray<T>;
};
export declare function findRouteMatch<T extends Extract<RouteLike, {
    fullPath: string;
}>>(
/** The path to match against the route tree. */
path: string, 
/** The `processedTree` returned by the initial `processRouteTree` call. */
processedTree: ProcessedTree<T, any, any>, 
/** If `true`, allows fuzzy matching (partial matches), i.e. which node in the tree would have been an exact match if the `path` had been shorter? */
fuzzy?: boolean): RouteMatch<T> | null;
/** Trim trailing slashes (except preserving root '/'). */
export declare function trimPathRight(path: string): string;
export interface ProcessRouteTreeResult<TRouteLike extends Extract<RouteLike, {
    fullPath: string;
}> & {
    id: string;
}> {
    /** Should be considered a black box, needs to be provided to all matching functions in this module. */
    processedTree: ProcessedTree<TRouteLike, any, any>;
    /** A lookup map of routes by their unique IDs. */
    routesById: Record<string, TRouteLike>;
    /** A lookup map of routes by their trimmed full paths. */
    routesByPath: Record<string, TRouteLike>;
}
/**
 * Processes a route tree into a segment trie for efficient path matching.
 * Also builds lookup maps for routes by ID and by trimmed full path.
 */
export declare function processRouteTree<TRouteLike extends Extract<RouteLike, {
    fullPath: string;
}> & {
    id: string;
}>(
/** The root of the route tree to process. */
routeTree: TRouteLike, 
/** Whether matching should be case sensitive by default (overridden by individual route options). */
caseSensitive?: boolean, 
/** Optional callback invoked for each route during processing. */
initRoute?: (route: TRouteLike, index: number) => void): ProcessRouteTreeResult<TRouteLike>;
declare function findMatch<T extends RouteLike>(path: string, segmentTree: AnySegmentNode<T>, fuzzy?: boolean): {
    route: T;
    /**
     * The raw (unparsed) params extracted from the path.
     * This will be the exhaustive list of all params defined in the route's path.
     */
    rawParams: Record<string, string>;
} | null;
export declare function buildRouteBranch<T extends RouteLike>(route: T): T[];
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/searchParams.d.ts
import { AnySchema } from './validators.js';
/** Default `parseSearch` that strips leading '?' and JSON-parses values. */
export declare const defaultParseSearch: (searchStr: string) => AnySchema;
/** Default `stringifySearch` using JSON.stringify for complex values. */
export declare const defaultStringifySearch: (search: Record<string, any>) => string;
/**
 * Build a `parseSearch` function using a provided JSON-like parser.
 *
 * The returned function strips a leading `?`, decodes values, and attempts to
 * JSON-parse string values using the given `parser`.
 *
 * @param parser Function to parse a string value (e.g. `JSON.parse`).
 * @returns A `parseSearch` function compatible with `Router` options.
 * @link https://tanstack.com/router/latest/docs/framework/react/guide/custom-search-param-serialization
 */
export declare function parseSearchWith(parser: (str: string) => any): (searchStr: string) => AnySchema;
/**
 * Build a `stringifySearch` function using a provided serializer.
 *
 * Non-primitive values are serialized with `stringify`. If a `parser` is
 * supplied, string values that are parseable are re-serialized to ensure
 * symmetry with `parseSearch`.
 *
 * @param stringify Function to serialize a value (e.g. `JSON.stringify`).
 * @param parser Optional parser to detect parseable strings.
 * @returns A `stringifySearch` function compatible with `Router` options.
 * @link https://tanstack.com/router/latest/docs/framework/react/guide/custom-search-param-serialization
 */
export declare function stringifySearchWith(stringify: (search: any) => string, parser?: (str: string) => any): (search: Record<string, any>) => string;
export type SearchSerializer = (searchObj: Record<string, any>) => string;
export type SearchParser = (searchStr: string) => Record<string, any>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/manifest.d.ts
export type AssetCrossOrigin = 'anonymous' | 'use-credentials';
export type ScriptFormat = 'module' | 'iife';
export declare const DEV_STYLES_ATTR = "data-tanstack-router-dev-styles";
export type AssetCrossOriginConfig = AssetCrossOrigin | Partial<Record<'script' | 'stylesheet', AssetCrossOrigin>>;
export type ManifestAssetLink = string | {
    href: string;
    crossOrigin?: AssetCrossOrigin;
};
export declare function getAssetCrossOrigin(assetCrossOrigin: AssetCrossOriginConfig | undefined, kind: 'script' | 'stylesheet'): AssetCrossOrigin | undefined;
export declare function getManifestScriptFormat(manifest: {
    scriptFormat?: ScriptFormat;
} | undefined): ScriptFormat;
export declare function getScriptPreloadAttrs(manifest: {
    scriptFormat?: ScriptFormat;
} | undefined, link: ManifestAssetLink, assetCrossOrigin?: AssetCrossOriginConfig): {
    rel: 'modulepreload' | 'preload';
    as?: 'script';
    href: string;
    crossOrigin?: AssetCrossOrigin;
};
export declare function resolveManifestAssetLink(link: ManifestAssetLink): {
    href: string;
    crossOrigin?: AssetCrossOrigin;
};
export type Manifest = {
    scriptFormat?: ScriptFormat;
    inlineStyle?: ManifestInlineCss;
    routes: Record<string, ManifestRoute>;
};
export type ServerManifest = {
    scriptFormat?: ScriptFormat;
    inlineCss?: ServerManifestInlineCss;
    routes: Record<string, ServerManifestRoute>;
};
export type ServerManifestInlineCss = {
    styles: Record<string, string>;
    templates?: Record<string, InlineCssTemplate>;
};
export type InlineCssTemplate = {
    strings: Array<string>;
    urls: Array<string>;
};
export type ManifestRoute = {
    filePath?: string;
    preloads?: Array<ManifestAssetLink>;
    scripts?: Array<ManifestScript>;
    css?: Array<ManifestCssLink>;
};
export type ServerManifestRoute = ManifestRoute;
export type ManifestRouteAssets = Pick<ManifestRoute, 'preloads' | 'scripts' | 'css'>;
export type RouterManagedTitleTag = {
    tag: 'title';
    attrs?: Record<string, any>;
    children: string;
};
export type RouterManagedMetaTag = {
    tag: 'meta';
    attrs?: Record<string, any>;
    children?: never;
};
export type RouterManagedLinkTag = {
    tag: 'link';
    attrs?: Record<string, any>;
    children?: never;
};
export type RouterManagedScriptTag = {
    tag: 'script';
    attrs?: Record<string, any>;
    children?: string;
};
export type ManifestScript = Omit<RouterManagedScriptTag, 'tag'>;
export type RouterManagedStyleTag = {
    tag: 'style';
    attrs?: Record<string, any>;
    children?: string;
    inlineCss?: true;
};
export type RouterManagedTag = RouterManagedTitleTag | RouterManagedMetaTag | RouterManagedLinkTag | RouterManagedScriptTag | RouterManagedStyleTag;
export declare function appendUniqueUserTags(target: Array<RouterManagedTag>, tags: Array<RouterManagedTag>): void;
export type ManifestCssLink = string | {
    href: string;
    crossOrigin?: AssetCrossOrigin;
    [DEV_STYLES_ATTR]?: true;
};
export type ManifestInlineCss = {
    attrs?: Record<string, any>;
    children?: string;
};
export type RouterManagedInlineCssTag = RouterManagedStyleTag & {
    inlineCss: true;
};
export declare function getStylesheetHref(asset: ManifestCssLink): string;
export declare function resolveManifestCssLink(link: ManifestCssLink): {
    href: string;
    crossOrigin?: AssetCrossOrigin;
    "data-tanstack-router-dev-styles"?: true;
};
export declare function createInlineCssStyleAsset(css: string): ManifestInlineCss;
export declare function createInlineCssPlaceholderAsset(): ManifestInlineCss;


// @filename: /node_modules/@tanstack/router-core/dist/esm/stores.d.ts
import { AnyRoute } from './route.js';
import { RouterState } from './router.js';
import { FullSearchSchema } from './routeInfo.js';
import { ParsedLocation } from './location.js';
import { AnyRedirect } from './redirect.js';
import { AnyRouteMatch } from './Matches.js';
export interface RouterReadableStore<TValue> {
    get: () => TValue;
}
export interface RouterWritableStore<TValue> extends RouterReadableStore<TValue> {
    set: ((updater: (prev: TValue) => TValue) => void) & ((value: TValue) => void);
}
export type RouterBatchFn = (fn: () => void) => void;
export type MutableStoreFactory = <TValue>(initialValue: TValue) => RouterWritableStore<TValue>;
export type ReadonlyStoreFactory = <TValue>(read: () => TValue) => RouterReadableStore<TValue>;
export type GetStoreConfig = (opts: {
    isServer?: boolean;
}) => StoreConfig;
export type StoreConfig = {
    createMutableStore: MutableStoreFactory;
    createReadonlyStore: ReadonlyStoreFactory;
    batch: RouterBatchFn;
    init?: (stores: RouterStores<AnyRoute>) => void;
};
type MatchStore = RouterWritableStore<AnyRouteMatch> & {
    routeId?: string;
};
type ReadableStore<TValue> = RouterReadableStore<TValue>;
/** SSR non-reactive createMutableStore */
export declare function createNonReactiveMutableStore<TValue>(initialValue: TValue): RouterWritableStore<TValue>;
/** SSR non-reactive createReadonlyStore */
export declare function createNonReactiveReadonlyStore<TValue>(read: () => TValue): RouterReadableStore<TValue>;
export interface RouterStores<in out TRouteTree extends AnyRoute> {
    status: RouterWritableStore<RouterState<TRouteTree>['status']>;
    loadedAt: RouterWritableStore<number>;
    isLoading: RouterWritableStore<boolean>;
    isTransitioning: RouterWritableStore<boolean>;
    location: RouterWritableStore<ParsedLocation<FullSearchSchema<TRouteTree>>>;
    resolvedLocation: RouterWritableStore<ParsedLocation<FullSearchSchema<TRouteTree>> | undefined>;
    statusCode: RouterWritableStore<number>;
    redirect: RouterWritableStore<AnyRedirect | undefined>;
    matchesId: RouterWritableStore<Array<string>>;
    pendingIds: RouterWritableStore<Array<string>>;
    matches: ReadableStore<Array<AnyRouteMatch>>;
    pendingMatches: ReadableStore<Array<AnyRouteMatch>>;
    cachedMatches: ReadableStore<Array<AnyRouteMatch>>;
    firstId: ReadableStore<string | undefined>;
    hasPending: ReadableStore<boolean>;
    matchRouteDeps: ReadableStore<{
        locationHref: string;
        resolvedLocationHref: string | undefined;
        status: RouterState<TRouteTree>['status'];
    }>;
    __store: RouterReadableStore<RouterState<TRouteTree>>;
    matchStores: Map<string, MatchStore>;
    pendingMatchStores: Map<string, MatchStore>;
    cachedMatchStores: Map<string, MatchStore>;
    /**
     * Get a computed store that resolves a routeId to its current match state.
     * Returns the same cached store instance for repeated calls with the same key.
     * The computed depends on matchesId + the individual match store, so
     * subscribers are only notified when the resolved match state changes.
     */
    getRouteMatchStore: (routeId: string) => RouterReadableStore<AnyRouteMatch | undefined>;
    setMatches: (nextMatches: Array<AnyRouteMatch>) => void;
    setPending: (nextMatches: Array<AnyRouteMatch>) => void;
    setCached: (nextMatches: Array<AnyRouteMatch>) => void;
}
export declare function createRouterStores<TRouteTree extends AnyRoute>(initialState: RouterState<TRouteTree>, config: StoreConfig): RouterStores<TRouteTree>;
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/router.d.ts
import { loadRouteChunk } from './load-matches.js';
import { LRUCache } from './lru-cache.js';
import { ProcessRouteTreeResult, ProcessedTree } from './new-process-route-tree.js';
import { SearchParser, SearchSerializer } from './searchParams.js';
import { AnyRedirect, ResolvedRedirect } from './redirect.js';
import { HistoryAction, HistoryLocation, HistoryState, ParsedHistoryState, RouterHistory } from '@tanstack/history';
import { Awaitable, Constrain, ControlledPromise, NoInfer, NonNullableUpdater, PickAsRequired, Updater } from './utils.js';
import { ParsedLocation } from './location.js';
import { AnyContext, AnyRoute, AnyRouteWithContext, LoaderStaleReloadMode, MakeRemountDepsOptionsUnion, RouteLike, RouteMask } from './route.js';
import { FullSearchSchema, RouteById, RoutePaths, RoutesById, RoutesByPath } from './routeInfo.js';
import { AnyRouteMatch, MakeRouteMatchUnion, MatchRouteOptions } from './Matches.js';
import { BuildLocationFn, CommitLocationOptions, NavigateFn } from './RouterProvider.js';
import { Manifest, ManifestRouteAssets, RouterManagedTag } from './manifest.js';
import { AnySchema } from './validators.js';
import { NavigateOptions, ResolveRelativePath, ToOptions } from './link.js';
import { AnySerializationAdapter, ValidateSerializableInput } from './ssr/serializer/transformer.js';
import { GetStoreConfig, RouterStores } from './stores.js';
export type ControllablePromise<T = any> = Promise<T> & {
    resolve: (value: T) => void;
    reject: (value?: any) => void;
};
export type InjectedHtmlEntry = Promise<string>;
export interface Register {
}
export type RegisteredSsr<TRegister = Register> = TRegister extends {
    ssr: infer TSSR;
} ? TSSR : false;
export type RegisteredRouter<TRegister = Register> = TRegister extends {
    router: infer TRouter;
} ? TRouter : AnyRouter;
export type RegisteredConfigType<TRegister, TKey> = TRegister extends {
    config: infer TConfig;
} ? TConfig extends {
    '~types': infer TTypes;
} ? TKey extends keyof TTypes ? TTypes[TKey] : unknown : unknown : unknown;
export type DefaultRemountDepsFn<TRouteTree extends AnyRoute> = (opts: MakeRemountDepsOptionsUnion<TRouteTree>) => any;
export interface DefaultRouterOptionsExtensions {
}
export interface RouterOptionsExtensions extends DefaultRouterOptionsExtensions {
}
export type SSROption = boolean | 'data-only';
export interface RouterOptions<TRouteTree extends AnyRoute, TTrailingSlashOption extends TrailingSlashOption, TDefaultStructuralSharingOption extends boolean = false, TRouterHistory extends RouterHistory = RouterHistory, TDehydrated = undefined> extends RouterOptionsExtensions {
    /**
     * The history object that will be used to manage the browser history.
     *
     * If not provided, a new createBrowserHistory instance will be created and used.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#history-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/history-types)
     */
    history?: TRouterHistory;
    /**
     * A function that will be used to stringify search params when generating links.
     *
     * @default defaultStringifySearch
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#stringifysearch-method)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/custom-search-param-serialization)
     */
    stringifySearch?: SearchSerializer;
    /**
     * A function that will be used to parse search params when parsing the current location.
     *
     * @default defaultParseSearch
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#parsesearch-method)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/custom-search-param-serialization)
     */
    parseSearch?: SearchParser;
    /**
     * If `false`, routes will not be preloaded by default in any way.
     *
     * If `'intent'`, routes will be preloaded by default when the user hovers over a link or a `touchstart` event is detected on a `<Link>`.
     *
     * If `'viewport'`, routes will be preloaded by default when they are within the viewport.
     *
     * @default false
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpreload-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/preloading)
     */
    defaultPreload?: false | 'intent' | 'viewport' | 'render';
    /**
     * The delay in milliseconds that a route must be hovered over or touched before it is preloaded.
     *
     * @default 50
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpreloaddelay-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/preloading#preload-delay)
     */
    defaultPreloadDelay?: number;
    /**
     * The default `preloadIntentProximity` a route should use if no preloadIntentProximity is provided.
     *
     * @default 0
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpreloadintentproximity-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/preloading#preload-intent-proximity)
     */
    defaultPreloadIntentProximity?: number;
    /**
     * The default `pendingMs` a route should use if no pendingMs is provided.
     *
     * @default 1000
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpendingms-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#avoiding-pending-component-flash)
     */
    defaultPendingMs?: number;
    /**
     * The default `pendingMinMs` a route should use if no pendingMinMs is provided.
     *
     * @default 500
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpendingminms-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#avoiding-pending-component-flash)
     */
    defaultPendingMinMs?: number;
    /**
     * The default `staleTime` a route should use if no staleTime is provided. This is the time in milliseconds that a route will be considered fresh.
     *
     * @default 0
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultstaletime-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#key-options)
     */
    defaultStaleTime?: number;
    /**
     * The default stale reload mode a route loader should use if no `loader.staleReloadMode` is provided.
     *
     * `'background'` preserves the current stale-while-revalidate behavior.
     * `'blocking'` waits for stale loader reloads to complete before resolving navigation.
     *
     * @default 'background'
     */
    defaultStaleReloadMode?: LoaderStaleReloadMode;
    /**
     * The default `preloadStaleTime` a route should use if no preloadStaleTime is provided.
     *
     * @default 30_000 `(30 seconds)`
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpreloadstaletime-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/preloading)
     */
    defaultPreloadStaleTime?: number;
    /**
     * The default `defaultPreloadGcTime` a route should use if no preloadGcTime is provided.
     *
     * @default 1_800_000 `(30 minutes)`
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpreloadgctime-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/preloading)
     */
    defaultPreloadGcTime?: number;
    /**
     * If `true`, route navigations will called using `document.startViewTransition()`.
     *
     * If the browser does not support this api, this option will be ignored.
     *
     * See [MDN](https://developer.mozilla.org/en-US/docs/Web/API/Document/startViewTransition) for more information on how this function works.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultviewtransition-property)
     */
    defaultViewTransition?: boolean | ViewTransitionOptions;
    /**
     * The default `hashScrollIntoView` a route should use if no hashScrollIntoView is provided while navigating
     *
     * See [MDN](https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView) for more information on `ScrollIntoViewOptions`.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaulthashscrollintoview-property)
     */
    defaultHashScrollIntoView?: boolean | ScrollIntoViewOptions;
    /**
     * @default 'fuzzy'
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#notfoundmode-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/not-found-errors#the-notfoundmode-option)
     */
    notFoundMode?: 'root' | 'fuzzy';
    /**
     * The default `gcTime` a route should use if no gcTime is provided.
     *
     * @default 1_800_000 `(30 minutes)`
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultgctime-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#key-options)
     */
    defaultGcTime?: number;
    /**
     * If `true`, all routes will be matched as case-sensitive.
     *
     * @default false
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#casesensitive-property)
     */
    caseSensitive?: boolean;
    /**
     *
     * The route tree that will be used to configure the router instance.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#routetree-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/routing/route-trees)
     */
    routeTree?: TRouteTree;
    /**
     * The basepath for then entire router. This is useful for mounting a router instance at a subpath.
     * ```
     * @default '/'
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#basepath-property)
     */
    basepath?: string;
    /**
     * The root context that will be provided to all routes in the route tree.
     *
     * This can be used to provide a context to all routes in the tree without having to provide it to each route individually.
     *
     * Optional or required if the root route was created with [`createRootRouteWithContext()`](https://tanstack.com/router/latest/docs/framework/react/api/router/createRootRouteWithContextFunction).
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#context-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/router-context)
     */
    context?: InferRouterContext<TRouteTree>;
    additionalContext?: any;
    /**
     * A function that will be called when the router is dehydrated.
     *
     * The return value of this function will be serialized and stored in the router's dehydrated state.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#dehydrate-method)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/external-data-loading#critical-dehydrationhydration)
     */
    dehydrate?: () => Awaitable<Constrain<TDehydrated, ValidateSerializableInput<Register, TDehydrated>>>;
    /**
     * A function that will be called when the router is hydrated.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#hydrate-method)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/external-data-loading#critical-dehydrationhydration)
     */
    hydrate?: (dehydrated: TDehydrated) => Awaitable<void>;
    /**
     * An array of route masks that will be used to mask routes in the route tree.
     *
     * Route masking is when you display a route at a different path than the one it is configured to match, like a modal popup that when shared will unmask to the modal's content instead of the modal's context.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#routemasks-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/route-masking)
     */
    routeMasks?: Array<RouteMask<TRouteTree>>;
    /**
     * If `true`, route masks will, by default, be removed when the page is reloaded.
     *
     * This can be overridden on a per-mask basis by setting the `unmaskOnReload` option on the mask, or on a per-navigation basis by setting the `unmaskOnReload` option in the `Navigate` options.
     *
     * @default false
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#unmaskonreload-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/route-masking#unmasking-on-page-reload)
     */
    unmaskOnReload?: boolean;
    /**
     * Use `notFoundComponent` instead.
     *
     * @deprecated
     * See https://tanstack.com/router/v1/docs/guide/not-found-errors#migrating-from-notfoundroute for more info.
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#notfoundroute-property)
     */
    notFoundRoute?: AnyRoute;
    /**
     * Configures how trailing slashes are treated.
     *
     * - `'always'` will add a trailing slash if not present
     * - `'never'` will remove the trailing slash if present
     * - `'preserve'` will not modify the trailing slash.
     *
     * @default 'never'
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#trailingslash-property)
     */
    trailingSlash?: TTrailingSlashOption;
    /**
     * While usually automatic, sometimes it can be useful to force the router into a server-side state, e.g. when using the router in a non-browser environment that has access to a global.document object.
     *
     * @default typeof document !== 'undefined'
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#isserver-property)
     */
    isServer?: boolean;
    /**
     * @default false
     */
    isShell?: boolean;
    /**
     * @default false
     */
    isPrerendering?: boolean;
    /**
     * The default `ssr` a route should use if no `ssr` is provided.
     *
     * @default true
     */
    defaultSsr?: SSROption;
    search?: {
        /**
         * Configures how unknown search params (= not returned by any `validateSearch`) are treated.
         *
         * @default false
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#search.strict-property)
         */
        strict?: boolean;
    };
    /**
     * Configures whether structural sharing is enabled by default for fine-grained selectors.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultstructuralsharing-property)
     */
    defaultStructuralSharing?: TDefaultStructuralSharingOption;
    /**
     * Configures which URI characters are allowed in path params that would ordinarily be escaped by encodeURIComponent.
     *
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#pathparamsallowedcharacters-property)
     * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/path-params#allowed-characters)
     */
    pathParamsAllowedCharacters?: Array<';' | ':' | '@' | '&' | '=' | '+' | '$' | ','>;
    defaultRemountDeps?: DefaultRemountDepsFn<TRouteTree>;
    /**
     * If `true`, scroll restoration will be enabled
     *
     * @default false
     */
    scrollRestoration?: boolean | ((opts: {
        location: ParsedLocation;
    }) => boolean);
    /**
     * A function that will be called to get the key for the scroll restoration cache.
     *
     * @default (location) => location.href
     */
    getScrollRestorationKey?: (location: ParsedLocation) => string;
    /**
     * The default behavior for scroll restoration.
     *
     * @default 'auto'
     */
    scrollRestorationBehavior?: ScrollBehavior;
    /**
     * An array of selectors that will be used to scroll to the top of the page in addition to `window`
     *
     * @default ['window']
     */
    scrollToTopSelectors?: Array<string | (() => Element | null | undefined)>;
    /**
     * When `true`, disables the global catch boundary that normally wraps all route matches.
     * This allows unhandled errors to bubble up to top-level error handlers in the browser.
     *
     * Useful for testing tools (like Storybook Test Runner), error reporting services,
     * and debugging scenarios where you want errors to reach the browser's global error handlers.
     *
     * @default false
     */
    disableGlobalCatchBoundary?: boolean;
    /**
     * An array of URL protocols to allow in links, redirects, and navigation.
     * Absolute URLs with protocols not in this list will be rejected.
     *
     * @default DEFAULT_PROTOCOL_ALLOWLIST (http:, https:, mailto:, tel:)
     * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#protocolallowlist-property)
     */
    protocolAllowlist?: Array<string>;
    serializationAdapters?: ReadonlyArray<AnySerializationAdapter>;
    /**
     * Configures how the router will rewrite the location between the actual href and the internal href of the router.
     *
     * @default undefined
     * @description You can provide a custom rewrite pair (in/out).
     * This is useful for shifting data from the origin to the path (for things like subdomain routing), or other advanced use cases.
     */
    rewrite?: LocationRewrite;
    origin?: string;
    ssr?: {
        nonce?: string;
    };
}
export type LocationRewrite = {
    /**
     * A function that will be called to rewrite the URL before it is interpreted by the router from the history instance.
     *
     * @default undefined
     */
    input?: LocationRewriteFunction;
    /**
     * A function that will be called to rewrite the URL before it is committed to the actual history instance from the router.
     *
     * @default undefined
     */
    output?: LocationRewriteFunction;
};
/**
 * A function that will be called to rewrite the URL.
 *
 * @param url The URL to rewrite.
 * @returns The rewritten URL (as a URL instance or full href string) or undefined if no rewrite is needed.
 */
export type LocationRewriteFunction = ({ url, }: {
    url: URL;
}) => undefined | string | URL;
export interface RouterState<in out TRouteTree extends AnyRoute = AnyRoute, in out TRouteMatch = MakeRouteMatchUnion> {
    status: 'pending' | 'idle';
    loadedAt: number;
    isLoading: boolean;
    isTransitioning: boolean;
    matches: Array<TRouteMatch>;
    location: ParsedLocation<FullSearchSchema<TRouteTree>>;
    resolvedLocation?: ParsedLocation<FullSearchSchema<TRouteTree>>;
    statusCode: number;
    redirect?: AnyRedirect;
}
export interface BuildNextOptions {
    to?: string | number | null;
    params?: true | Updater<unknown>;
    search?: true | Updater<unknown>;
    hash?: true | Updater<string>;
    state?: true | NonNullableUpdater<ParsedHistoryState, HistoryState>;
    mask?: {
        to?: string | number | null;
        params?: true | Updater<unknown>;
        search?: true | Updater<unknown>;
        hash?: true | Updater<string>;
        state?: true | NonNullableUpdater<ParsedHistoryState, HistoryState>;
        unmaskOnReload?: boolean;
    };
    from?: string;
    href?: string;
    _fromLocation?: ParsedLocation;
    unsafeRelative?: 'path';
    _isNavigate?: boolean;
}
type NavigationEventInfo = {
    fromLocation?: ParsedLocation;
    toLocation: ParsedLocation;
    pathChanged: boolean;
    hrefChanged: boolean;
    hashChanged: boolean;
};
export interface RouterEvents {
    onBeforeNavigate: {
        type: 'onBeforeNavigate';
    } & NavigationEventInfo;
    onBeforeLoad: {
        type: 'onBeforeLoad';
    } & NavigationEventInfo;
    onLoad: {
        type: 'onLoad';
    } & NavigationEventInfo;
    onResolved: {
        type: 'onResolved';
    } & NavigationEventInfo;
    onBeforeRouteMount: {
        type: 'onBeforeRouteMount';
    } & NavigationEventInfo;
    onRendered: {
        type: 'onRendered';
    } & NavigationEventInfo;
}
export type RouterEvent = RouterEvents[keyof RouterEvents];
export type ListenerFn<TEvent extends RouterEvent> = (event: TEvent) => void;
export type RouterListener<TRouterEvent extends RouterEvent> = {
    eventType: TRouterEvent['type'];
    fn: ListenerFn<TRouterEvent>;
};
export type SubscribeFn = <TType extends keyof RouterEvents>(eventType: TType, fn: ListenerFn<RouterEvents[TType]>) => () => void;
export interface MatchRoutesOpts {
    preload?: boolean;
    throwOnError?: boolean;
    dest?: BuildNextOptions;
}
export type InferRouterContext<TRouteTree extends AnyRoute> = TRouteTree['types']['routerContext'];
export type RouterContextOptions<TRouteTree extends AnyRoute> = AnyContext extends InferRouterContext<TRouteTree> ? {
    context?: InferRouterContext<TRouteTree>;
} : {
    context: InferRouterContext<TRouteTree>;
};
export type RouterConstructorOptions<TRouteTree extends AnyRoute, TTrailingSlashOption extends TrailingSlashOption, TDefaultStructuralSharingOption extends boolean, TRouterHistory extends RouterHistory, TDehydrated extends Record<string, any>> = Omit<RouterOptions<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>, 'context' | 'serializationAdapters' | 'defaultSsr'> & RouterContextOptions<TRouteTree>;
export type PreloadRouteFn<TRouteTree extends AnyRoute, TTrailingSlashOption extends TrailingSlashOption, TDefaultStructuralSharingOption extends boolean, TRouterHistory extends RouterHistory> = <TFrom extends RoutePaths<TRouteTree> | string = string, TTo extends string | undefined = undefined, TMaskFrom extends RoutePaths<TRouteTree> | string = TFrom, TMaskTo extends string = ''>(opts: NavigateOptions<RouterCore<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory>, TFrom, TTo, TMaskFrom, TMaskTo> & {}) => Promise<Array<AnyRouteMatch> | undefined>;
export type MatchRouteFn<TRouteTree extends AnyRoute, TTrailingSlashOption extends TrailingSlashOption, TDefaultStructuralSharingOption extends boolean, TRouterHistory extends RouterHistory> = <TFrom extends RoutePaths<TRouteTree> = '/', TTo extends string | undefined = undefined, TResolved = ResolveRelativePath<TFrom, NoInfer<TTo>>>(location: ToOptions<RouterCore<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory>, TFrom, TTo>, opts?: MatchRouteOptions) => false | RouteById<TRouteTree, TResolved>['types']['allParams'];
export type UpdateFn<TRouteTree extends AnyRoute, TTrailingSlashOption extends TrailingSlashOption, TDefaultStructuralSharingOption extends boolean, TRouterHistory extends RouterHistory, TDehydrated extends Record<string, any>> = (newOptions: RouterConstructorOptions<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>) => void;
export type InvalidateFn<TRouter extends AnyRouter> = (opts?: {
    filter?: (d: MakeRouteMatchUnion<TRouter>) => boolean;
    sync?: boolean;
    forcePending?: boolean;
}) => Promise<void>;
export type ParseLocationFn<TRouteTree extends AnyRoute> = (locationToParse: HistoryLocation, previousLocation?: ParsedLocation<FullSearchSchema<TRouteTree>>) => ParsedLocation<FullSearchSchema<TRouteTree>>;
export type GetMatchRoutesFn = (pathname: string) => {
    matchedRoutes: ReadonlyArray<AnyRoute>;
    /** exhaustive params, still in their string form */
    routeParams: Record<string, string>;
    foundRoute: AnyRoute | undefined;
    parseError?: unknown;
};
export type EmitFn = (routerEvent: RouterEvent) => void;
export type LoadFn = (opts?: {
    sync?: boolean;
    action?: {
        type: HistoryAction;
    };
}) => Promise<void>;
export type CommitLocationFn = ({ viewTransition, ignoreBlocker, ...next }: ParsedLocation & CommitLocationOptions) => Promise<void>;
export type StartTransitionFn = (fn: () => void) => void;
export interface MatchRoutesFn {
    (pathname: string, locationSearch?: AnySchema, opts?: MatchRoutesOpts): Array<MakeRouteMatchUnion>;
    /**
     * @deprecated use the following signature instead
     */
    (next: ParsedLocation, opts?: MatchRoutesOpts): Array<AnyRouteMatch>;
    (pathnameOrNext: string | ParsedLocation, locationSearchOrOpts?: AnySchema | MatchRoutesOpts, opts?: MatchRoutesOpts): Array<AnyRouteMatch>;
}
export type GetMatchFn = (matchId: string) => AnyRouteMatch | undefined;
export type UpdateMatchFn = (id: string, updater: (match: AnyRouteMatch) => AnyRouteMatch) => void;
export type LoadRouteChunkFn = (route: AnyRoute) => Promise<Array<void>>;
export type ResolveRedirect = (err: AnyRedirect) => ResolvedRedirect;
export type ClearCacheFn<TRouter extends AnyRouter> = (opts?: {
    filter?: (d: MakeRouteMatchUnion<TRouter>) => boolean;
}) => void;
export interface ServerSsr {
    /** Framework-only: injects router-owned HTML into the SSR stream. */
    injectHtml: (html: string) => void;
    /** Framework-only: injects a router-owned script tag into the SSR stream. */
    injectScript: (script: string) => void;
    isDehydrated: () => boolean;
    isSerializationFinished: () => boolean;
    /** Framework-only: atomically reserves the pass-through stream path if safe. */
    reserveStreamFastPath: () => boolean;
    /** Framework-only. */
    onInjectedHtml: (listener: () => void) => () => void;
    /** Framework-only. */
    onRenderFinished: (listener: () => void) => void;
    /** Framework-only. */
    setRenderFinished: () => void;
    /** Framework-only. */
    cleanup: () => void;
    /**
     * Register a listener invoked when the SSR request lifecycle ends (success,
     * error, abort, or stream lifetime expiry). Use to tear down per-request
     * resources whose references would otherwise pin the router (e.g. query
     * cache subscriptions, gcTime timers, abort controllers).
     *
     * Listeners run synchronously and exactly once. Errors are caught and logged.
     */
    onCleanup: (listener: () => void) => void;
    /** Framework-only. */
    onSerializationFinished: (listener: () => void) => () => void;
    /** Framework-only. */
    dehydrate: (opts?: {
        requestAssets?: ManifestRouteAssets;
    }) => Promise<void>;
    /** Framework-only. */
    takeBufferedScripts: () => RouterManagedTag | undefined;
    /** Framework-only: takes buffered router-owned HTML. */
    takeBufferedHtml: () => string | undefined;
    /** Framework-only. */
    liftScriptBarrier: () => void;
}
export interface RouterSsrLifecycle {
    onServerSsrAttach?: Array<(serverSsr: ServerSsr) => void>;
}
export type AnyRouterWithContext<TContext> = RouterCore<AnyRouteWithContext<TContext>, any, any, any, any>;
export type AnyRouter = RouterCore<any, any, any, any, any>;
export interface ViewTransitionOptions {
    types: Array<string> | ((locationChangeInfo: {
        fromLocation?: ParsedLocation;
        toLocation: ParsedLocation;
        pathChanged: boolean;
        hrefChanged: boolean;
        hashChanged: boolean;
    }) => Array<string> | false);
}
/**
 * Convert an unknown error into a minimal, serializable object.
 * Includes name and message (and stack in development).
 */
export declare function defaultSerializeError(err: unknown): {
    name: string;
    message: string;
} | {
    data: unknown;
};
/** Options for configuring trailing-slash behavior. */
export declare const trailingSlashOptions: {
    readonly always: "always";
    readonly never: "never";
    readonly preserve: "preserve";
};
export type TrailingSlashOption = (typeof trailingSlashOptions)[keyof typeof trailingSlashOptions];
/**
 * Compute whether path, href or hash changed between previous and current
 * resolved locations.
 */
export declare function getLocationChangeInfo(location: ParsedLocation, resolvedLocation?: ParsedLocation): {
    fromLocation: ParsedLocation<{}> | undefined;
    toLocation: ParsedLocation<{}>;
    pathChanged: boolean;
    hrefChanged: boolean;
    hashChanged: boolean;
};
export declare const locationHistoryActions: WeakMap<ParsedLocation<{}>, HistoryAction>;
export type CreateRouterFn = <TRouteTree extends AnyRoute, TTrailingSlashOption extends TrailingSlashOption = 'never', TDefaultStructuralSharingOption extends boolean = false, TRouterHistory extends RouterHistory = RouterHistory, TDehydrated extends Record<string, any> = Record<string, any>>(options: undefined extends number ? 'strictNullChecks must be enabled in tsconfig.json' : RouterConstructorOptions<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>) => RouterCore<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>;
declare global {
    var __TSR_CACHE__: {
        routeTree: AnyRoute;
        processRouteTreeResult: ProcessRouteTreeResult<AnyRoute>;
        resolvePathCache: LRUCache<string, string>;
    } | undefined;
}
/**
 * Core, framework-agnostic router engine that powers TanStack Router.
 *
 * Provides navigation, matching, loading, preloading, caching and event APIs
 * used by framework adapters (React/Solid). Prefer framework helpers like
 * `createRouter` in app code.
 *
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/RouterType
 */
export declare class RouterCore<in out TRouteTree extends AnyRoute, in out TTrailingSlashOption extends TrailingSlashOption, in out TDefaultStructuralSharingOption extends boolean, in out TRouterHistory extends RouterHistory = RouterHistory, in out TDehydrated extends Record<string, any> = Record<string, any>> {
    tempLocationKey: string | undefined;
    _scroll: {
        next: boolean;
        restoring?: boolean;
        restoration?: boolean;
        reset?: boolean;
    };
    shouldViewTransition?: boolean | ViewTransitionOptions;
    isViewTransitionTypesSupported?: boolean;
    subscribers: Set<RouterListener<RouterEvent>>;
    viewTransitionPromise?: ControlledPromise<true>;
    stores: RouterStores<TRouteTree>;
    private getStoreConfig;
    batch: (fn: () => void) => void;
    options: PickAsRequired<RouterOptions<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>, 'stringifySearch' | 'parseSearch' | 'context'>;
    history: TRouterHistory;
    rewrite?: LocationRewrite;
    origin?: string;
    latestLocation: ParsedLocation<FullSearchSchema<TRouteTree>>;
    pendingBuiltLocation?: ParsedLocation<FullSearchSchema<TRouteTree>>;
    basepath: string;
    routeTree: TRouteTree;
    routesById: RoutesById<TRouteTree>;
    routesByPath: RoutesByPath<TRouteTree>;
    processedTree: ProcessedTree<TRouteTree, any, any>;
    resolvePathCache: LRUCache<string, string>;
    private routeBranchCache;
    isServer: boolean;
    pathParamsDecoder?: (encoded: string) => string;
    protocolAllowlist: Set<string>;
    /**
     * @deprecated Use the `createRouter` function instead
     */
    constructor(options: RouterConstructorOptions<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>, getStoreConfig: GetStoreConfig);
    startTransition: StartTransitionFn;
    isShell(): boolean;
    isPrerendering(): boolean;
    update: UpdateFn<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>;
    get state(): RouterState<TRouteTree>;
    updateLatestLocation: () => void;
    buildRouteTree: () => ProcessRouteTreeResult<TRouteTree>;
    setRoutes({ routesById, routesByPath, processedTree, }: ProcessRouteTreeResult<TRouteTree>): void;
    /**
     * Subscribe to router lifecycle events like `onBeforeNavigate`, `onLoad`,
     * `onResolved`, etc. Returns an unsubscribe function.
     *
     * @link https://tanstack.com/router/latest/docs/framework/react/api/router/RouterEventsType
     */
    subscribe: SubscribeFn;
    emit: EmitFn;
    /**
     * Parse a HistoryLocation into a strongly-typed ParsedLocation using the
     * current router options, rewrite rules and search parser/stringifier.
     */
    parseLocation: ParseLocationFn<TRouteTree>;
    /** Resolve a path using the router's trailing-slash policy. */
    resolvePathWithBase: (from: string, path: string) => string;
    private getRouteBranch;
    get looseRoutesById(): Record<string, AnyRoute>;
    matchRoutes: MatchRoutesFn;
    private getParentContext;
    private matchRoutesInternal;
    getMatchedRoutes: GetMatchRoutesFn;
    /**
     * Lightweight route matching for buildLocation.
     * Only computes fullPath, accumulated search, and params - skipping expensive
     * operations like AbortController, ControlledPromise, loaderDeps, and full match objects.
     */
    private matchRoutesLightweight;
    cancelMatch: (id: string) => void;
    cancelMatches: () => void;
    /**
     * Build the next ParsedLocation from navigation options without committing.
     * Resolves `to`/`from`, params/search/hash/state, applies search validation
     * and middlewares, and returns a stable, stringified location object.
     *
     * @link https://tanstack.com/router/latest/docs/framework/react/api/router/RouterType#buildlocation-method
     */
    buildLocation: BuildLocationFn;
    commitLocationPromise: undefined | ControlledPromise<void>;
    /**
     * Commit a previously built location to history (push/replace), optionally
     * using view transitions and scroll restoration options.
     */
    commitLocation: CommitLocationFn;
    /** Convenience helper: build a location from options, then commit it. */
    buildAndCommitLocation: ({ replace, resetScroll, hashScrollIntoView, viewTransition, ignoreBlocker, href, ...rest }?: BuildNextOptions & CommitLocationOptions) => Promise<void>;
    /**
     * Imperatively navigate using standard `NavigateOptions`. When `reloadDocument`
     * or an absolute `href` is provided, performs a full document navigation.
     * Otherwise, builds and commits a client-side location.
     *
     * @link https://tanstack.com/router/latest/docs/framework/react/api/router/NavigateOptionsType
     */
    navigate: NavigateFn;
    latestLoadPromise: undefined | Promise<void>;
    beforeLoad: () => void;
    load: LoadFn;
    startViewTransition: (fn: () => Promise<void>) => void;
    updateMatch: UpdateMatchFn;
    getMatch: GetMatchFn;
    /**
     * Invalidate the current matches and optionally force them back into a pending state.
     *
     * - Marks all matches that pass the optional `filter` as `invalid: true`.
     * - If `forcePending` is true, or a match is currently in `'error'` or `'notFound'` status,
     *   its status is reset to `'pending'` and its `error` cleared so that the loader is re-run
     *   on the next `load()` call (eg. after HMR or a manual invalidation).
     */
    invalidate: InvalidateFn<RouterCore<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>>;
    getParsedLocationHref: (location: ParsedLocation) => string;
    resolveRedirect: (redirect: AnyRedirect) => AnyRedirect;
    clearCache: ClearCacheFn<this>;
    clearExpiredCache: () => void;
    loadRouteChunk: typeof loadRouteChunk;
    preloadRoute: PreloadRouteFn<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory>;
    matchRoute: MatchRouteFn<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory>;
    ssr?: {
        manifest: Manifest | undefined;
    };
    serverSsr?: ServerSsr;
    serverSsrLifecycle?: RouterSsrLifecycle;
    hasNotFoundMatch: () => boolean;
}
/** Error thrown when search parameter validation fails. */
export declare class SearchParamError extends Error {
}
/** Error thrown when path parameter parsing/validation fails. */
export declare class PathParamError extends Error {
}
/**
 * Lazily import a module function and forward arguments to it, retaining
 * parameter and return types for the selected export key.
 */
export declare function lazyFn<T extends Record<string, (...args: Array<any>) => any>, TKey extends keyof T = 'default'>(fn: () => Promise<T>, key?: TKey): (...args: Parameters<T[TKey]>) => Promise<Awaited<ReturnType<T[TKey]>>>;
/** Create an initial RouterState from a parsed location. */
export declare function getInitialRouterState(location: ParsedLocation): RouterState<any>;
/**
 * Build the matched route chain and extract params for a pathname.
 * Falls back to the root route if no specific route is found.
 */
export declare function getMatchedRoutes<TRouteLike extends RouteLike>({ pathname, routesById, processedTree, }: {
    pathname: string;
    routesById: Record<string, TRouteLike>;
    processedTree: ProcessedTree<any, any, any>;
}): {
    matchedRoutes: readonly TRouteLike[];
    routeParams: Record<string, string>;
    foundRoute: TRouteLike | undefined;
};
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/global.d.ts
import { AnyRouter } from './router.js';
declare global {
    interface Window {
        __TSR_ROUTER__?: AnyRouter;
    }
}
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/defer.d.ts
import { defaultSerializeError } from './router.js';
/**
 * Well-known symbol used by {@link defer} to tag a promise with
 * its deferred state. Consumers can read `promise[TSR_DEFERRED_PROMISE]`
 * to access `status`, `data`, or `error`.
 */
export declare const TSR_DEFERRED_PROMISE: unique symbol;
export type DeferredPromiseState<T> = {
    status: 'pending';
    data?: T;
    error?: unknown;
} | {
    status: 'success';
    data: T;
} | {
    status: 'error';
    data?: T;
    error: unknown;
};
export type DeferredPromise<T> = Promise<T> & {
    [TSR_DEFERRED_PROMISE]: DeferredPromiseState<T>;
};
/**
 * Wrap a promise with a deferred state for use with `<Await>` and `useAwaited`.
 *
 * The returned promise is augmented with internal state (status/data/error)
 * so UI can read progress or suspend until it settles.
 *
 * @param _promise The promise to wrap.
 * @param options Optional config. Provide `serializeError` to customize how
 * errors are serialized for transfer.
 * @returns The same promise with attached deferred metadata.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/deferFunction
 */
export declare function defer<T>(_promise: Promise<T>, options?: {
    serializeError?: typeof defaultSerializeError;
}): DeferredPromise<T>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/invariant.d.ts
export declare function invariant(): never;


// @filename: /node_modules/@tanstack/router-core/dist/esm/path.d.ts
import { LRUCache } from './lru-cache.js';
/** Join path segments, cleaning duplicate slashes between parts. */
export declare function joinPaths(paths: Array<string | undefined>): string;
/** Remove repeated slashes from a path string. */
export declare function cleanPath(path: string): string;
/** Trim leading slashes (except preserving root '/'). */
export declare function trimPathLeft(path: string): string;
/** Trim trailing slashes (except preserving root '/'). */
export declare function trimPathRight(path: string): string;
/** Trim both leading and trailing slashes. */
export declare function trimPath(path: string): string;
/** Remove a trailing slash from value when appropriate for comparisons. */
export declare function removeTrailingSlash(value: string, basepath: string): string;
/**
 * Compare two pathnames for exact equality after normalizing trailing slashes
 * relative to the provided `basepath`.
 */
export declare function exactPathTest(pathName1: string, pathName2: string, basepath: string): boolean;
interface ResolvePathOptions {
    base: string;
    to: string;
    trailingSlash?: 'always' | 'never' | 'preserve';
    cache?: LRUCache<string, string>;
}
/**
 * Resolve a destination path against a base, honoring trailing-slash policy
 * and supporting relative segments (`.`/`..`) and absolute `to` values.
 */
export declare function resolvePath({ base, to, trailingSlash, cache, }: ResolvePathOptions): string;
/**
 * Create a pre-compiled decode config from allowed characters.
 * This should be called once at router initialization.
 */
export declare function compileDecodeCharMap(pathParamsAllowedCharacters: ReadonlyArray<string>): (encoded: string) => string;
interface InterpolatePathOptions {
    path?: string;
    params: Record<string, unknown>;
    /**
     * A function that decodes a path parameter value.
     * Obtained from `compileDecodeCharMap(pathParamsAllowedCharacters)`.
     */
    decoder?: (encoded: string) => string;
}
type InterPolatePathResult = {
    interpolatedPath: string;
    usedParams: Record<string, unknown>;
    isMissingParams: boolean;
};
/**
 * Interpolate params and wildcards into a route path template.
 *
 * - Encodes params safely (configurable allowed characters)
 * - Supports `{-$optional}` segments, `{prefix{$id}suffix}` and `{$}` wildcards
 */
export declare function interpolatePath({ path, params, decoder, ...rest }: InterpolatePathOptions): InterPolatePathResult;
export {};


// @filename: /node_modules/@tanstack/router-core/dist/esm/qss.d.ts
/**
 * Program is a reimplementation of the `qss` package:
 * Copyright (c) Luke Edwards luke.edwards05@gmail.com, MIT License
 * https://github.com/lukeed/qss/blob/master/license.md
 *
 * This reimplementation uses modern browser APIs
 * (namely URLSearchParams) and TypeScript while still
 * maintaining the original functionality and interface.
 *
 * Update: this implementation has also been mangled to
 * fit exactly our use-case (single value per key in encoding).
 */
/**
 * Encodes an object into a query string.
 * @param obj - The object to encode into a query string.
 * @param stringify - An optional custom stringify function.
 * @returns The encoded query string.
 * @example
 * ```
 * // Example input: encode({ token: 'foo', key: 'value' })
 * // Expected output: "token=foo&key=value"
 * ```
 */
export declare function encode(obj: Record<string, any>, stringify?: (value: any) => string): string;
/**
 * Decodes a query string into an object.
 * @param str - The query string to decode.
 * @returns The decoded key-value pairs in an object format.
 * @example
 * // Example input: decode("token=foo&key=value")
 * // Expected output: { "token": "foo", "key": "value" }
 */
export declare function decode(str: any): any;


// @filename: /node_modules/@tanstack/router-core/dist/esm/config.d.ts
import { SSROption } from './router.js';
import { AnySerializationAdapter } from './ssr/serializer/transformer.js';
export interface RouterConfigOptions<in out TSerializationAdapters, in out TDefaultSsr> {
    serializationAdapters?: TSerializationAdapters;
    defaultSsr?: TDefaultSsr;
}
export interface RouterConfig<in out TSerializationAdapters, in out TDefaultSsr> {
    '~types': RouterConfigTypes<TSerializationAdapters, TDefaultSsr>;
    serializationAdapters: TSerializationAdapters;
    defaultSsr: TDefaultSsr | undefined;
}
export interface RouterConfigTypes<in out TSerializationAdapters, in out TDefaultSsr> {
    serializationAdapters: TSerializationAdapters;
    defaultSsr: TDefaultSsr;
}
export declare const createRouterConfig: <const TSerializationAdapters extends ReadonlyArray<AnySerializationAdapter> = [], TDefaultSsr extends SSROption = SSROption>(options: RouterConfigOptions<TSerializationAdapters, TDefaultSsr>) => RouterConfig<TSerializationAdapters, TDefaultSsr>;
export type AnyRouterConfig = RouterConfig<any, any>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/searchMiddleware.d.ts
import { NoInfer, PickOptional } from './utils.js';
import { SearchMiddleware } from './route.js';
import { IsRequiredParams } from './link.js';
/**
 * Retain specified search params across navigations.
 *
 * If `keys` is `true`, retain all current params. Otherwise, copy only the
 * listed keys from the current search into the next search.
 *
 * @param keys `true` to retain all, or a list of keys to retain.
 * @returns A search middleware suitable for route `search.middlewares`.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/retainSearchParamsFunction
 */
export declare function retainSearchParams<TSearchSchema extends object>(keys: Array<keyof TSearchSchema> | true): SearchMiddleware<TSearchSchema>;
/**
 * Remove optional or default-valued search params from navigations.
 *
 * - Pass `true` (only if there are no required search params) to strip all.
 * - Pass an array to always remove those optional keys.
 * - Pass an object of default values; keys equal (deeply) to the defaults are removed.
 *
 * @returns A search middleware suitable for route `search.middlewares`.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/stripSearchParamsFunction
 */
export declare function stripSearchParams<TSearchSchema, TOptionalProps = PickOptional<NoInfer<TSearchSchema>>, const TValues = Partial<NoInfer<TSearchSchema>> | Array<keyof TOptionalProps>, const TInput = IsRequiredParams<TSearchSchema> extends never ? TValues | true : TValues>(input: NoInfer<TInput>): SearchMiddleware<TSearchSchema>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/structuralSharing.d.ts
import { Constrain } from './utils.js';
export interface OptionalStructuralSharing<TStructuralSharing, TConstraint> {
    readonly structuralSharing?: Constrain<TStructuralSharing, TConstraint> | undefined;
}


// @filename: /node_modules/@tanstack/router-core/dist/esm/useRouteContext.d.ts
import { AllContext, RouteById } from './routeInfo.js';
import { AnyRouter } from './router.js';
import { Expand, StrictOrFrom } from './utils.js';
export interface UseRouteContextBaseOptions<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TSelected> {
    select?: (search: ResolveUseRouteContext<TRouter, TFrom, TStrict>) => TSelected;
}
export type UseRouteContextOptions<TRouter extends AnyRouter, TFrom extends string | undefined, TStrict extends boolean, TSelected> = StrictOrFrom<TRouter, TFrom, TStrict> & UseRouteContextBaseOptions<TRouter, TFrom, TStrict, TSelected>;
export type ResolveUseRouteContext<TRouter extends AnyRouter, TFrom, TStrict extends boolean> = TStrict extends false ? AllContext<TRouter['routeTree']> : Expand<RouteById<TRouter['routeTree'], TFrom>['types']['allContext']>;
export type UseRouteContextResult<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TSelected> = unknown extends TSelected ? ResolveUseRouteContext<TRouter, TFrom, TStrict> : TSelected;


// @filename: /node_modules/@tanstack/router-core/dist/esm/useSearch.d.ts
import { FullSearchSchema, RouteById } from './routeInfo.js';
import { AnyRouter } from './router.js';
import { Expand } from './utils.js';
export type UseSearchResult<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TSelected> = unknown extends TSelected ? ResolveUseSearch<TRouter, TFrom, TStrict> : TSelected;
export type ResolveUseSearch<TRouter extends AnyRouter, TFrom, TStrict extends boolean> = TStrict extends false ? FullSearchSchema<TRouter['routeTree']> : Expand<RouteById<TRouter['routeTree'], TFrom>['types']['fullSearchSchema']>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/useParams.d.ts
import { AllParams, RouteById } from './routeInfo.js';
import { AnyRouter } from './router.js';
import { Expand } from './utils.js';
export type ResolveUseParams<TRouter extends AnyRouter, TFrom, TStrict extends boolean> = TStrict extends false ? AllParams<TRouter['routeTree']> : Expand<RouteById<TRouter['routeTree'], TFrom>['types']['allParams']>;
export type UseParamsResult<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TSelected> = unknown extends TSelected ? ResolveUseParams<TRouter, TFrom, TStrict> : TSelected;


// @filename: /node_modules/@tanstack/router-core/dist/esm/useNavigate.d.ts
import { NavigateOptions } from './link.js';
import { RegisteredRouter } from './router.js';
export type UseNavigateResult<TDefaultFrom extends string> = <TRouter extends RegisteredRouter, TTo extends string | undefined, TFrom extends string = TDefaultFrom, TMaskFrom extends string = TFrom, TMaskTo extends string = ''>({ from, ...rest }: NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>) => Promise<void>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/useLoaderDeps.d.ts
import { RouteById } from './routeInfo.js';
import { AnyRouter } from './router.js';
import { Expand } from './utils.js';
export type ResolveUseLoaderDeps<TRouter extends AnyRouter, TFrom> = Expand<RouteById<TRouter['routeTree'], TFrom>['types']['loaderDeps']>;
export type UseLoaderDepsResult<TRouter extends AnyRouter, TFrom, TSelected> = unknown extends TSelected ? ResolveUseLoaderDeps<TRouter, TFrom> : TSelected;


// @filename: /node_modules/@tanstack/router-core/dist/esm/useLoaderData.d.ts
import { AllLoaderData, RouteById } from './routeInfo.js';
import { AnyRouter } from './router.js';
import { Expand } from './utils.js';
export type ResolveUseLoaderData<TRouter extends AnyRouter, TFrom, TStrict extends boolean> = TStrict extends false ? AllLoaderData<TRouter['routeTree']> : Expand<RouteById<TRouter['routeTree'], TFrom>['types']['loaderData']>;
export type UseLoaderDataResult<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TSelected> = unknown extends TSelected ? ResolveUseLoaderData<TRouter, TFrom, TStrict> : TSelected;


// @filename: /node_modules/@tanstack/router-core/dist/esm/scroll-restoration.d.ts
import { AnyRouter } from './router.js';
import { ParsedLocation } from './location.js';
export type ScrollRestorationEntry = {
    scrollX: number;
    scrollY: number;
};
export type ScrollRestorationOptions = {
    getKey?: (location: ParsedLocation) => string;
    scrollBehavior?: ScrollToOptions['behavior'];
};
export declare const storageKey = "tsr-scroll-restoration-v1_3";
/**
 * The default `getKey` function for `useScrollRestoration`.
 * It returns the `key` from the location state or the `href` of the location.
 *
 * The `location.href` is used as a fallback to support the use case where the location state is not available like the initial render.
 */
export declare const defaultGetScrollRestorationKey: (location: ParsedLocation) => string;
export declare function getElementScrollRestorationEntry(router: AnyRouter, options: ({
    id: string;
    getElement?: () => Window | Element | undefined | null;
} | {
    id?: string;
    getElement: () => Window | Element | undefined | null;
}) & {
    getKey?: (location: ParsedLocation) => string;
}): ScrollRestorationEntry | undefined;
export declare function setupScrollRestoration(router: AnyRouter, force?: boolean): void;


// @filename: /node_modules/@tanstack/router-core/dist/esm/typePrimitives.d.ts
import { FromPathOption, NavigateOptions, PathParamOptions, SearchParamOptions, ToPathOption } from './link.js';
import { RedirectOptions } from './redirect.js';
import { RouteIds } from './routeInfo.js';
import { AnyRouter, RegisteredRouter } from './router.js';
import { UseParamsResult } from './useParams.js';
import { UseSearchResult } from './useSearch.js';
import { Constrain, ConstrainLiteral } from './utils.js';
export type ValidateFromPath<TRouter extends AnyRouter = RegisteredRouter, TFrom = string> = FromPathOption<TRouter, TFrom>;
export type ValidateToPath<TRouter extends AnyRouter = RegisteredRouter, TTo extends string | undefined = undefined, TFrom extends string = string> = ToPathOption<TRouter, TFrom, TTo>;
export type ValidateSearch<TRouter extends AnyRouter = RegisteredRouter, TTo extends string | undefined = undefined, TFrom extends string = string> = SearchParamOptions<TRouter, TFrom, TTo>;
export type ValidateParams<TRouter extends AnyRouter = RegisteredRouter, TTo extends string | undefined = undefined, TFrom extends string = string> = PathParamOptions<TRouter, TFrom, TTo>;
/**
 * @private
 */
export type InferFrom<TOptions, TDefaultFrom extends string = string> = TOptions extends {
    from: infer TFrom extends string;
} ? TFrom : TDefaultFrom;
/**
 * @private
 */
export type InferTo<TOptions> = TOptions extends {
    to: infer TTo extends string;
} ? TTo : undefined;
/**
 * @private
 */
export type InferMaskTo<TOptions> = TOptions extends {
    mask: {
        to: infer TTo extends string;
    };
} ? TTo : '';
export type InferMaskFrom<TOptions> = TOptions extends {
    mask: {
        from: infer TFrom extends string;
    };
} ? TFrom : string;
export type ValidateNavigateOptions<TRouter extends AnyRouter = RegisteredRouter, TOptions = unknown, TDefaultFrom extends string = string> = Constrain<TOptions, NavigateOptions<TRouter, InferFrom<TOptions, TDefaultFrom>, InferTo<TOptions>, InferMaskFrom<TOptions>, InferMaskTo<TOptions>>>;
export type ValidateNavigateOptionsArray<TRouter extends AnyRouter = RegisteredRouter, TOptions extends ReadonlyArray<any> = ReadonlyArray<unknown>, TDefaultFrom extends string = string> = {
    [K in keyof TOptions]: ValidateNavigateOptions<TRouter, TOptions[K], TDefaultFrom>;
};
export type ValidateRedirectOptions<TRouter extends AnyRouter = RegisteredRouter, TOptions = unknown, TDefaultFrom extends string = string> = Constrain<TOptions, RedirectOptions<TRouter, InferFrom<TOptions, TDefaultFrom>, InferTo<TOptions>, InferMaskFrom<TOptions>, InferMaskTo<TOptions>>>;
export type ValidateRedirectOptionsArray<TRouter extends AnyRouter = RegisteredRouter, TOptions extends ReadonlyArray<any> = ReadonlyArray<unknown>, TDefaultFrom extends string = string> = {
    [K in keyof TOptions]: ValidateRedirectOptions<TRouter, TOptions[K], TDefaultFrom>;
};
export type ValidateId<TRouter extends AnyRouter = RegisteredRouter, TId extends string = string> = ConstrainLiteral<TId, RouteIds<TRouter['routeTree']>>;
/**
 * @private
 */
export type InferStrict<TOptions> = TOptions extends {
    strict: infer TStrict extends boolean;
} ? TStrict : true;
/**
 * @private
 */
export type InferShouldThrow<TOptions> = TOptions extends {
    shouldThrow: infer TShouldThrow extends boolean;
} ? TShouldThrow : true;
/**
 * @private
 */
export type InferSelected<TOptions> = TOptions extends {
    select: (...args: Array<any>) => infer TSelected;
} ? TSelected : unknown;
export type ValidateUseSearchResult<TOptions, TRouter extends AnyRouter = RegisteredRouter> = UseSearchResult<TRouter, InferFrom<TOptions>, InferStrict<TOptions>, InferSelected<TOptions>>;
export type ValidateUseParamsResult<TOptions, TRouter extends AnyRouter = RegisteredRouter> = Constrain<TOptions, UseParamsResult<TRouter, InferFrom<TOptions>, InferStrict<TOptions>, InferSelected<TOptions>>>;


// @filename: /node_modules/@tanstack/router-core/dist/esm/ssr/serializer/seroval-plugins.d.ts
import { RawStream } from './RawStream.js';
import { Plugin } from 'seroval';
export declare const defaultSerovalPlugins: (Plugin<Error, any> | Plugin<RawStream, any> | Plugin<ReadableStream<any>, any>)[];


// @filename: /node_modules/@tanstack/router-core/dist/esm/rewrite.d.ts
import { LocationRewrite } from './router.js';
/** Compose multiple rewrite pairs into a single in/out rewrite. */
export declare function composeRewrites(rewrites: Array<LocationRewrite>): {
    input: ({ url }: {
        url: URL;
    }) => URL;
    output: ({ url }: {
        url: URL;
    }) => URL;
};
/** Create a rewrite pair that strips/adds a basepath on input/output. */
export declare function rewriteBasepath(opts: {
    basepath: string;
    caseSensitive?: boolean;
}): {
    input: ({ url }: {
        url: URL;
    }) => URL;
    output: ({ url }: {
        url: URL;
    }) => URL;
};
/** Execute a location input rewrite if provided. */
export declare function executeRewriteInput(rewrite: LocationRewrite | undefined, url: URL): URL;
/** Execute a location output rewrite if provided. */
export declare function executeRewriteOutput(rewrite: LocationRewrite | undefined, url: URL): URL;


// @filename: /node_modules/@tanstack/router-core/dist/esm/index.d.ts
export * from './global.js';
export { TSR_DEFERRED_PROMISE, defer } from './defer.js';
export type { DeferredPromiseState, DeferredPromise } from './defer.js';
export { invariant } from './invariant.js';
export { preloadWarning } from './link.js';
export type { IsRequiredParams, AddTrailingSlash, RemoveTrailingSlashes, AddLeadingSlash, RemoveLeadingSlashes, ActiveOptions, LinkOptionsProps, ResolveCurrentPath, ResolveParentPath, ResolveRelativePath, FindDescendantToPaths, InferDescendantToPaths, RelativeToPath, RelativeToParentPath, RelativeToCurrentPath, AbsoluteToPath, RelativeToPathAutoComplete, NavigateOptions, ToOptions, ToMaskOptions, ToSubOptions, ResolveRoute, SearchParamOptions, PathParamOptions, ToPathOption, LinkOptions, MakeOptionalPathParams, FromPathOption, MakeOptionalSearchParams, MaskOptions, ToSubOptionsProps, RequiredToOptions, } from './link.js';
export type { RouteToPath, TrailingSlashOptionByRouter, ParseRoute, CodeRouteToPath, RouteIds, FullSearchSchema, FullSearchSchemaInput, AllParams, RouteById, AllContext, RoutePaths, RoutesById, RoutesByPath, AllLoaderData, RouteByPath, } from './routeInfo.js';
export type { InferFileRouteTypes, FileRouteTypes, FileRoutesByPath, CreateFileRoute, LazyRoute, LazyRouteOptions, CreateLazyFileRoute, } from './fileRoute.js';
export type { ParsedLocation } from './location.js';
export type { Manifest, ServerManifest, ManifestRoute, ManifestRouteAssets, ServerManifestRoute, ManifestCssLink, ManifestInlineCss, ServerManifestInlineCss, InlineCssTemplate, ManifestScript, RouterManagedTag, RouterManagedTitleTag, RouterManagedMetaTag, RouterManagedInlineCssTag, RouterManagedScriptTag, RouterManagedLinkTag, RouterManagedStyleTag, AssetCrossOrigin, AssetCrossOriginConfig, ManifestAssetLink, ScriptFormat, } from './manifest.js';
export { DEV_STYLES_ATTR, appendUniqueUserTags, createInlineCssStyleAsset, getAssetCrossOrigin, getManifestScriptFormat, getScriptPreloadAttrs, getStylesheetHref, resolveManifestAssetLink, resolveManifestCssLink, } from './manifest.js';
export { isMatch } from './Matches.js';
export type { AnyMatchAndValue, FindValueByIndex, FindValueByKey, CreateMatchAndValue, NextMatchAndValue, IsMatchKeyOf, IsMatchPath, IsMatchResult, IsMatchParse, IsMatch, RouteMatch, RouteMatchExtensions, MakeRouteMatchUnion, MakeRouteMatch, AnyRouteMatch, MakeRouteMatchFromRoute, MatchRouteOptions, } from './Matches.js';
export { joinPaths, cleanPath, trimPathLeft, trimPathRight, trimPath, removeTrailingSlash, exactPathTest, resolvePath, interpolatePath, } from './path.js';
export { encode, decode } from './qss.js';
export { rootRouteId } from './root.js';
export type { RootRouteId } from './root.js';
export { BaseRoute, BaseRouteApi, BaseRootRoute } from './route.js';
export type { AnyPathParams, SearchSchemaInput, AnyContext, RouteContext, PreloadableObj, RoutePathOptions, StaticDataRouteOption, RoutePathOptionsIntersection, SearchFilter, SearchMiddlewareContext, SearchMiddleware, ResolveId, InferFullSearchSchema, InferFullSearchSchemaInput, InferAllParams, InferAllContext, MetaDescriptor, RouteLinkEntry, SearchValidator, AnySearchValidator, DefaultSearchValidator, ErrorRouteProps, ErrorComponentProps, NotFoundRouteProps, ResolveParams, ParseParamsFn, StringifyParamsFn, ParamsOptions, UpdatableStaticRouteOption, ContextReturnType, ContextAsyncReturnType, ResolveRouteContext, ResolveLoaderData, RoutePrefix, TrimPath, TrimPathLeft, TrimPathRight, ResolveSearchSchemaFnInput, ResolveSearchSchemaInput, ResolveSearchSchemaFn, ResolveSearchSchema, ResolveFullSearchSchema, ResolveFullSearchSchemaInput, ResolveAllContext, BeforeLoadContextParameter, RouteContextParameter, ResolveAllParamsFromParent, AnyRoute, Route, RouteTypes, FullSearchSchemaOption, RemountDepsOptions, MakeRemountDepsOptionsUnion, ResolveFullPath, AnyRouteWithContext, RouteOptions, FileBaseRouteOptions, BaseRouteOptions, UpdatableRouteOptions, LoaderStaleReloadMode, RouteLoaderFn, RouteLoaderEntry, LoaderFnContext, RouteContextFn, ContextOptions, RouteContextOptions, SsrContextOptions, BeforeLoadContextOptions, RootRouteOptions, RootRouteOptionsExtensions, UpdatableRouteOptionsExtensions, RouteConstraints, RouteTypesById, RouteMask, RouteExtensions, RouteLazyFn, RouteAddChildrenFn, RouteAddFileChildrenFn, RouteAddFileTypesFn, ResolveOptionalParams, ResolveRequiredParams, RootRoute, FilebaseRouteOptionsInterface, } from './route.js';
export { createNonReactiveMutableStore, createNonReactiveReadonlyStore, } from './stores.js';
export type { RouterBatchFn, RouterReadableStore, GetStoreConfig, RouterStores, RouterWritableStore, } from './stores.js';
export { defaultSerializeError, getLocationChangeInfo, RouterCore, lazyFn, SearchParamError, PathParamError, getInitialRouterState, getMatchedRoutes, trailingSlashOptions, } from './router.js';
export type { ViewTransitionOptions, TrailingSlashOption, Register, AnyRouter, AnyRouterWithContext, RegisteredRouter, RouterState, BuildNextOptions, RouterListener, RouterEvent, ListenerFn, RouterEvents, MatchRoutesOpts, RouterOptionsExtensions, DefaultRemountDepsFn, PreloadRouteFn, MatchRouteFn, RouterContextOptions, RouterOptions, RouterConstructorOptions, UpdateFn, ParseLocationFn, InvalidateFn, ControllablePromise, InjectedHtmlEntry, EmitFn, LoadFn, GetMatchFn, SubscribeFn, UpdateMatchFn, CommitLocationFn, GetMatchRoutesFn, MatchRoutesFn, StartTransitionFn, LoadRouteChunkFn, ClearCacheFn, CreateRouterFn, SSROption, } from './router.js';
export * from './config.js';
export type { MatchLocation, CommitLocationOptions, NavigateFn, BuildLocationFn, } from './RouterProvider.js';
export { retainSearchParams, stripSearchParams } from './searchMiddleware.js';
export { defaultParseSearch, defaultStringifySearch, parseSearchWith, stringifySearchWith, } from './searchParams.js';
export type { SearchSerializer, SearchParser } from './searchParams.js';
export type { OptionalStructuralSharing } from './structuralSharing.js';
export { functionalUpdate, hasKeys, replaceEqualDeep, isPlainObject, isPlainArray, deepEqual, createControlledPromise, isModuleNotFoundError, DEFAULT_PROTOCOL_ALLOWLIST, escapeHtml, isDangerousProtocol, buildDevStylesUrl, } from './utils.js';
export type { NoInfer, IsAny, PickAsRequired, PickRequired, PickOptional, WithoutEmpty, Expand, DeepPartial, MakeDifferenceOptional, IsUnion, IsNonEmptyObject, Assign, IntersectAssign, Timeout, Updater, NonNullableUpdater, StringLiteral, ThrowOrOptional, ThrowConstraint, ControlledPromise, ExtractObjects, PartialMergeAllObject, MergeAllPrimitive, ExtractPrimitives, PartialMergeAll, Constrain, ConstrainLiteral, UnionToIntersection, MergeAllObjects, MergeAll, ValidateJSON, StrictOrFrom, LooseReturnType, LooseAsyncReturnType, Awaitable, } from './utils.js';
export type { StandardSchemaValidatorProps, StandardSchemaValidator, AnyStandardSchemaValidator, StandardSchemaValidatorTypes, AnyStandardSchemaValidateSuccess, AnyStandardSchemaValidateFailure, AnyStandardSchemaValidateIssue, AnyStandardSchemaValidateInput, AnyStandardSchemaValidate, ValidatorObj, AnyValidatorObj, ValidatorAdapter, AnyValidatorAdapter, AnyValidatorFn, ValidatorFn, Validator, AnyValidator, AnySchema, DefaultValidator, ResolveSearchValidatorInputFn, ResolveSearchValidatorInput, ResolveValidatorInputFn, ResolveValidatorInput, ResolveValidatorOutputFn, ResolveValidatorOutput, } from './validators.js';
export type { UseRouteContextBaseOptions, UseRouteContextOptions, UseRouteContextResult, } from './useRouteContext.js';
export type { UseSearchResult, ResolveUseSearch } from './useSearch.js';
export type { UseParamsResult, ResolveUseParams } from './useParams.js';
export type { UseNavigateResult } from './useNavigate.js';
export type { UseLoaderDepsResult, ResolveUseLoaderDeps } from './useLoaderDeps.js';
export type { UseLoaderDataResult, ResolveUseLoaderData } from './useLoaderData.js';
export type { Redirect, RedirectOptions, RedirectOptionsRoute, RedirectFnRoute, ResolvedRedirect, AnyRedirect, } from './redirect.js';
export { redirect, isRedirect, isResolvedRedirect, parseRedirect, } from './redirect.js';
export type { NotFoundError } from './not-found.js';
export { isNotFound, notFound } from './not-found.js';
export { defaultGetScrollRestorationKey, getElementScrollRestorationEntry, storageKey, setupScrollRestoration, } from './scroll-restoration.js';
export type { ScrollRestorationOptions, ScrollRestorationEntry, } from './scroll-restoration.js';
export type { ValidateFromPath, ValidateToPath, ValidateSearch, ValidateParams, InferFrom, InferTo, InferMaskTo, InferMaskFrom, ValidateNavigateOptions, ValidateNavigateOptionsArray, ValidateRedirectOptions, ValidateRedirectOptionsArray, ValidateId, InferStrict, InferShouldThrow, InferSelected, ValidateUseSearchResult, ValidateUseParamsResult, } from './typePrimitives.js';
export type { AnySerializationAdapter, SerializationAdapter, ValidateSerializableInput, SerializerExtensions, ValidateSerializable, RegisteredSerializableInput, SerializableExtensions, DefaultSerializable, Serializable, TSR_SERIALIZABLE, TsrSerializable, SerializationError, } from './ssr/serializer/transformer.js';
export { createSerializationAdapter, makeSerovalPlugin, makeSsrSerovalPlugin, } from './ssr/serializer/transformer.js';
export { defaultSerovalPlugins } from './ssr/serializer/seroval-plugins.js';
export { RawStream, createRawStreamRPCPlugin, createRawStreamDeserializePlugin, } from './ssr/serializer/RawStream.js';
export type { OnRawStreamCallback, RawStreamHint, RawStreamOptions, } from './ssr/serializer/RawStream.js';
export { composeRewrites, executeRewriteInput } from './rewrite.js';
export type { LocationRewrite, LocationRewriteFunction } from './router.js';


// @filename: /node_modules/@tanstack/react-router/dist/esm/awaited.d.ts
import * as React from 'react';
export type AwaitOptions<T> = {
    promise: Promise<T>;
};
/** Suspend until a deferred promise resolves or rejects and return its data. */
export declare function useAwaited<T>({ promise: _promise }: AwaitOptions<T>): T;
/**
 * Component that suspends on a deferred promise and renders its child with
 * the resolved value. Optionally provides a Suspense fallback.
 */
export declare function Await<T>(props: AwaitOptions<T> & {
    fallback?: React.ReactNode;
    children: (result: T) => React.ReactNode;
}): import("react/jsx-runtime").JSX.Element;


// @filename: /node_modules/@tanstack/react-router/dist/esm/structuralSharing.d.ts
import { AnyRouter, Constrain, OptionalStructuralSharing, ValidateJSON } from '@tanstack/router-core';
export type DefaultStructuralSharingEnabled<TRouter extends AnyRouter> = boolean extends TRouter['options']['defaultStructuralSharing'] ? false : NonNullable<TRouter['options']['defaultStructuralSharing']>;
export interface RequiredStructuralSharing<TStructuralSharing, TConstraint> {
    readonly structuralSharing: Constrain<TStructuralSharing, TConstraint>;
}
export type StructuralSharingOption<TRouter extends AnyRouter, TSelected, TStructuralSharing> = unknown extends TSelected ? OptionalStructuralSharing<TStructuralSharing, boolean> : unknown extends TRouter['routeTree'] ? OptionalStructuralSharing<TStructuralSharing, boolean> : TSelected extends ValidateJSON<TSelected> ? OptionalStructuralSharing<TStructuralSharing, boolean> : DefaultStructuralSharingEnabled<TRouter> extends true ? RequiredStructuralSharing<TStructuralSharing, false> : OptionalStructuralSharing<TStructuralSharing, false>;
export type StructuralSharingEnabled<TRouter extends AnyRouter, TStructuralSharing> = boolean extends TStructuralSharing ? DefaultStructuralSharingEnabled<TRouter> : TStructuralSharing;
export type ValidateSelected<TRouter extends AnyRouter, TSelected, TStructuralSharing> = StructuralSharingEnabled<TRouter, TStructuralSharing> extends true ? ValidateJSON<TSelected> : TSelected;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useLoaderData.d.ts
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
import { AnyRouter, RegisteredRouter, ResolveUseLoaderData, StrictOrFrom, UseLoaderDataResult } from '@tanstack/router-core';
export interface UseLoaderDataBaseOptions<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TSelected, TStructuralSharing> {
    select?: (match: ResolveUseLoaderData<TRouter, TFrom, TStrict>) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
}
export type UseLoaderDataOptions<TRouter extends AnyRouter, TFrom extends string | undefined, TStrict extends boolean, TSelected, TStructuralSharing> = StrictOrFrom<TRouter, TFrom, TStrict> & UseLoaderDataBaseOptions<TRouter, TFrom, TStrict, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>;
export type UseLoaderDataRoute<out TId> = <TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseLoaderDataBaseOptions<TRouter, TId, true, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>) => UseLoaderDataResult<TRouter, TId, true, TSelected>;
/**
 * Read and select the current route's loader data with type‑safety.
 *
 * Options:
 * - `from`/`strict`: Choose which route's data to read and strictness
 * - `select`: Map the loader data to a derived value
 * - `structuralSharing`: Enable structural sharing for stable references
 *
 * @returns The loader data (or selected value) for the matched route.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useLoaderDataHook
 */
export declare function useLoaderData<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string | undefined = undefined, TStrict extends boolean = true, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts: UseLoaderDataOptions<TRouter, TFrom, TStrict, TSelected, TStructuralSharing>): UseLoaderDataResult<TRouter, TFrom, TStrict, TSelected>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useMatch.d.ts
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
import { AnyRouter, MakeRouteMatch, MakeRouteMatchUnion, RegisteredRouter, StrictOrFrom, ThrowConstraint, ThrowOrOptional } from '@tanstack/router-core';
export interface UseMatchBaseOptions<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TThrow extends boolean, TSelected, TStructuralSharing extends boolean> {
    select?: (match: MakeRouteMatch<TRouter['routeTree'], TFrom, TStrict>) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
    shouldThrow?: TThrow;
}
export type UseMatchRoute<out TFrom> = <TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseMatchBaseOptions<TRouter, TFrom, true, true, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>) => UseMatchResult<TRouter, TFrom, true, TSelected>;
export type UseMatchOptions<TRouter extends AnyRouter, TFrom extends string | undefined, TStrict extends boolean, TThrow extends boolean, TSelected, TStructuralSharing extends boolean> = StrictOrFrom<TRouter, TFrom, TStrict> & UseMatchBaseOptions<TRouter, TFrom, TStrict, TThrow, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>;
export type UseMatchResult<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TSelected> = unknown extends TSelected ? TStrict extends true ? MakeRouteMatch<TRouter['routeTree'], TFrom, TStrict> : MakeRouteMatchUnion<TRouter> : TSelected;
/**
 * Read and select the nearest or targeted route match.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useMatchHook
 */
export declare function useMatch<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string | undefined = undefined, TStrict extends boolean = true, TThrow extends boolean = true, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts: UseMatchOptions<TRouter, TFrom, TStrict, ThrowConstraint<TStrict, TThrow>, TSelected, TStructuralSharing>): ThrowOrOptional<UseMatchResult<TRouter, TFrom, TStrict, TSelected>, TThrow>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useLoaderDeps.d.ts
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
import { AnyRouter, RegisteredRouter, ResolveUseLoaderDeps, StrictOrFrom, UseLoaderDepsResult } from '@tanstack/router-core';
export interface UseLoaderDepsBaseOptions<TRouter extends AnyRouter, TFrom, TSelected, TStructuralSharing> {
    select?: (deps: ResolveUseLoaderDeps<TRouter, TFrom>) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
}
export type UseLoaderDepsOptions<TRouter extends AnyRouter, TFrom extends string | undefined, TSelected, TStructuralSharing> = StrictOrFrom<TRouter, TFrom> & UseLoaderDepsBaseOptions<TRouter, TFrom, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>;
export type UseLoaderDepsRoute<out TId> = <TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseLoaderDepsBaseOptions<TRouter, TId, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>) => UseLoaderDepsResult<TRouter, TId, TSelected>;
/**
 * Read and select the current route's loader dependencies object.
 *
 * Options:
 * - `from`: Choose which route's loader deps to read
 * - `select`: Map the deps to a derived value
 * - `structuralSharing`: Enable structural sharing for stable references
 *
 * @returns The loader deps (or selected value) for the matched route.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useLoaderDepsHook
 */
export declare function useLoaderDeps<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string | undefined = undefined, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts: UseLoaderDepsOptions<TRouter, TFrom, TSelected, TStructuralSharing>): UseLoaderDepsResult<TRouter, TFrom, TSelected>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useParams.d.ts
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
import { AnyRouter, RegisteredRouter, ResolveUseParams, StrictOrFrom, ThrowConstraint, ThrowOrOptional, UseParamsResult } from '@tanstack/router-core';
export interface UseParamsBaseOptions<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TThrow extends boolean, TSelected, TStructuralSharing> {
    select?: (params: ResolveUseParams<TRouter, TFrom, TStrict>) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
    shouldThrow?: TThrow;
}
export type UseParamsOptions<TRouter extends AnyRouter, TFrom extends string | undefined, TStrict extends boolean, TThrow extends boolean, TSelected, TStructuralSharing> = StrictOrFrom<TRouter, TFrom, TStrict> & UseParamsBaseOptions<TRouter, TFrom, TStrict, TThrow, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>;
export type UseParamsRoute<out TFrom> = <TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseParamsBaseOptions<TRouter, TFrom, true, true, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>) => UseParamsResult<TRouter, TFrom, true, TSelected>;
/**
 * Access the current route's path parameters with type-safety.
 *
 * Options:
 * - `from`/`strict`: Specify the matched route and whether to enforce strict typing
 * - `select`: Project the params object to a derived value for memoized renders
 * - `structuralSharing`: Enable structural sharing for stable references
 * - `shouldThrow`: Throw if the route is not found in strict contexts
 *
 * @returns The params object (or selected value) for the matched route.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useParamsHook
 */
export declare function useParams<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string | undefined = undefined, TStrict extends boolean = true, TThrow extends boolean = true, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts: UseParamsOptions<TRouter, TFrom, TStrict, ThrowConstraint<TStrict, TThrow>, TSelected, TStructuralSharing>): ThrowOrOptional<UseParamsResult<TRouter, TFrom, TStrict, TSelected>, TThrow>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useSearch.d.ts
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
import { AnyRouter, RegisteredRouter, ResolveUseSearch, StrictOrFrom, ThrowConstraint, ThrowOrOptional, UseSearchResult } from '@tanstack/router-core';
export interface UseSearchBaseOptions<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TThrow extends boolean, TSelected, TStructuralSharing> {
    select?: (state: ResolveUseSearch<TRouter, TFrom, TStrict>) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
    shouldThrow?: TThrow;
}
export type UseSearchOptions<TRouter extends AnyRouter, TFrom, TStrict extends boolean, TThrow extends boolean, TSelected, TStructuralSharing> = StrictOrFrom<TRouter, TFrom, TStrict> & UseSearchBaseOptions<TRouter, TFrom, TStrict, TThrow, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>;
export type UseSearchRoute<out TFrom> = <TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseSearchBaseOptions<TRouter, TFrom, true, true, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>) => UseSearchResult<TRouter, TFrom, true, TSelected>;
/**
 * Read and select the current route's search parameters with type-safety.
 *
 * Options:
 * - `from`/`strict`: Control which route's search is read and how strictly it's typed
 * - `select`: Map the search object to a derived value for render optimization
 * - `structuralSharing`: Enable structural sharing for stable references
 * - `shouldThrow`: Throw when the route is not found (strict contexts)
 *
 * @returns The search object (or selected value) for the matched route.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useSearchHook
 */
export declare function useSearch<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string | undefined = undefined, TStrict extends boolean = true, TThrow extends boolean = true, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts: UseSearchOptions<TRouter, TFrom, TStrict, ThrowConstraint<TStrict, TThrow>, TSelected, TStructuralSharing>): ThrowOrOptional<UseSearchResult<TRouter, TFrom, TStrict, TSelected>, TThrow>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useRouteContext.d.ts
import { AnyRouter, RegisteredRouter, UseRouteContextBaseOptions, UseRouteContextOptions, UseRouteContextResult } from '@tanstack/router-core';
export type UseRouteContextRoute<out TFrom> = <TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown>(opts?: UseRouteContextBaseOptions<TRouter, TFrom, true, TSelected>) => UseRouteContextResult<TRouter, TFrom, true, TSelected>;
export declare function useRouteContext<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string | undefined = undefined, TStrict extends boolean = true, TSelected = unknown>(opts: UseRouteContextOptions<TRouter, TFrom, TStrict, TSelected>): UseRouteContextResult<TRouter, TFrom, TStrict, TSelected>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/typePrimitives.d.ts
import { AnyRouter, Constrain, InferFrom, InferMaskFrom, InferMaskTo, InferSelected, InferShouldThrow, InferStrict, InferTo, RegisteredRouter } from '@tanstack/router-core';
import { LinkComponentProps } from './link.js';
import { UseParamsOptions } from './useParams.js';
import { UseSearchOptions } from './useSearch.js';
export type ValidateLinkOptions<TRouter extends AnyRouter = RegisteredRouter, TOptions = unknown, TDefaultFrom extends string = string, TComp = 'a'> = Constrain<TOptions, LinkComponentProps<TComp, TRouter, InferFrom<TOptions, TDefaultFrom>, InferTo<TOptions>, InferMaskFrom<TOptions>, InferMaskTo<TOptions>>>;
/**
 * @private
 */
export type InferStructuralSharing<TOptions> = TOptions extends {
    structuralSharing: infer TStructuralSharing;
} ? TStructuralSharing : unknown;
export type ValidateUseSearchOptions<TOptions, TRouter extends AnyRouter = RegisteredRouter> = Constrain<TOptions, UseSearchOptions<TRouter, InferFrom<TOptions>, InferStrict<TOptions>, InferShouldThrow<TOptions>, InferSelected<TOptions>, InferStructuralSharing<TOptions>>>;
export type ValidateUseParamsOptions<TOptions, TRouter extends AnyRouter = RegisteredRouter> = Constrain<TOptions, UseParamsOptions<TRouter, InferFrom<TOptions>, InferStrict<TOptions>, InferShouldThrow<TOptions>, InferSelected<TOptions>, InferSelected<TOptions>>>;
export type ValidateLinkOptionsArray<TRouter extends AnyRouter = RegisteredRouter, TOptions extends ReadonlyArray<any> = ReadonlyArray<unknown>, TDefaultFrom extends string = string, TComp = 'a'> = {
    [K in keyof TOptions]: ValidateLinkOptions<TRouter, TOptions[K], TDefaultFrom, TComp>;
};


// @filename: /node_modules/@tanstack/react-router/dist/esm/link.d.ts
import { AnyRouter, Constrain, LinkOptions, RegisteredRouter, RoutePaths } from '@tanstack/router-core';
import { ReactNode } from 'react';
import { ValidateLinkOptions, ValidateLinkOptionsArray } from './typePrimitives.js';
import * as React from 'react';
/**
 * Build anchor-like props for declarative navigation and preloading.
 *
 * Returns stable `href`, event handlers and accessibility props derived from
 * router options and active state. Used internally by `Link` and custom links.
 *
 * Options cover `to`, `params`, `search`, `hash`, `state`, `preload`,
 * `activeProps`, `inactiveProps`, and more.
 *
 * @returns React anchor props suitable for `<a>` or custom components.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useLinkPropsHook
 */
export declare function useLinkProps<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string = string, const TTo extends string | undefined = undefined, const TMaskFrom extends string = TFrom, const TMaskTo extends string = ''>(options: UseLinkPropsOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>, forwardedRef?: React.ForwardedRef<Element>): React.ComponentPropsWithRef<'a'>;
type UseLinkReactProps<TComp> = TComp extends keyof React.JSX.IntrinsicElements ? React.JSX.IntrinsicElements[TComp] : TComp extends React.ComponentType<any> ? React.ComponentPropsWithoutRef<TComp> & React.RefAttributes<React.ComponentRef<TComp>> : never;
export type UseLinkPropsOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends RoutePaths<TRouter['routeTree']> | string = string, TTo extends string | undefined = '.', TMaskFrom extends RoutePaths<TRouter['routeTree']> | string = TFrom, TMaskTo extends string = '.'> = ActiveLinkOptions<'a', TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & UseLinkReactProps<'a'>;
export type ActiveLinkOptions<TComp = 'a', TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = '.', TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = LinkOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & ActiveLinkOptionProps<TComp>;
type ActiveLinkProps<TComp> = Partial<LinkComponentReactProps<TComp> & {
    [key: `data-${string}`]: unknown;
}>;
export interface ActiveLinkOptionProps<TComp = 'a'> {
    /**
     * A function that returns additional props for the `active` state of this link.
     * These props override other props passed to the link (`style`'s are merged, `className`'s are concatenated)
     */
    activeProps?: ActiveLinkProps<TComp> | (() => ActiveLinkProps<TComp>);
    /**
     * A function that returns additional props for the `inactive` state of this link.
     * These props override other props passed to the link (`style`'s are merged, `className`'s are concatenated)
     */
    inactiveProps?: ActiveLinkProps<TComp> | (() => ActiveLinkProps<TComp>);
}
export type LinkProps<TComp = 'a', TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = '.', TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = ActiveLinkOptions<TComp, TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & LinkPropsChildren;
export interface LinkPropsChildren {
    children?: React.ReactNode | ((state: {
        isActive: boolean;
        isTransitioning: boolean;
    }) => React.ReactNode);
}
type LinkComponentReactProps<TComp> = Omit<UseLinkReactProps<TComp>, keyof CreateLinkProps>;
export type LinkComponentProps<TComp = 'a', TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = '.', TMaskFrom extends string = TFrom, TMaskTo extends string = '.'> = LinkComponentReactProps<TComp> & LinkProps<TComp, TRouter, TFrom, TTo, TMaskFrom, TMaskTo>;
export type CreateLinkProps = LinkProps<any, any, string, string, string, string>;
export type LinkComponent<in out TComp, in out TDefaultFrom extends string = string> = <TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string = TDefaultFrom, const TTo extends string | undefined = undefined, const TMaskFrom extends string = TFrom, const TMaskTo extends string = ''>(props: LinkComponentProps<TComp, TRouter, TFrom, TTo, TMaskFrom, TMaskTo>) => React.ReactElement;
export interface LinkComponentRoute<in out TDefaultFrom extends string = string> {
    defaultFrom: TDefaultFrom;
    <TRouter extends AnyRouter = RegisteredRouter, const TTo extends string | undefined = undefined, const TMaskTo extends string = ''>(props: LinkComponentProps<'a', TRouter, this['defaultFrom'], TTo, this['defaultFrom'], TMaskTo>): React.ReactElement;
}
/**
 * Creates a typed Link-like component that preserves TanStack Router's
 * navigation semantics and type-safety while delegating rendering to the
 * provided host component.
 *
 * Useful for integrating design system anchors/buttons while keeping
 * router-aware props (eg. `to`, `params`, `search`, `preload`).
 *
 * @param Comp The host component to render (eg. a design-system Link/Button)
 * @returns A router-aware component with the same API as `Link`.
 * @link https://tanstack.com/router/latest/docs/framework/react/guide/custom-link
 */
export declare function createLink<const TComp>(Comp: Constrain<TComp, any, (props: CreateLinkProps) => ReactNode>): LinkComponent<TComp>;
/**
 * A strongly-typed anchor component for declarative navigation.
 * Handles path, search, hash and state updates with optional route preloading
 * and active-state styling.
 *
 * Props:
 * - `preload`: Controls route preloading (eg. 'intent', 'render', 'viewport', true/false)
 * - `preloadDelay`: Delay in ms before preloading on hover
 * - `activeProps`/`inactiveProps`: Additional props merged when link is active/inactive
 * - `resetScroll`/`hashScrollIntoView`: Control scroll behavior on navigation
 * - `viewTransition`/`startTransition`: Use View Transitions/React transitions for navigation
 * - `ignoreBlocker`: Bypass registered blockers
 *
 * @returns An anchor-like element that navigates without full page reloads.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/linkComponent
 */
export declare const Link: LinkComponent<'a'>;
export type LinkOptionsFnOptions<TOptions, TComp, TRouter extends AnyRouter = RegisteredRouter> = TOptions extends ReadonlyArray<any> ? ValidateLinkOptionsArray<TRouter, TOptions, string, TComp> : ValidateLinkOptions<TRouter, TOptions, string, TComp>;
export type LinkOptionsFn<TComp> = <const TOptions, TRouter extends AnyRouter = RegisteredRouter>(options: LinkOptionsFnOptions<TOptions, TComp, TRouter>) => TOptions;
/**
 * Validate and reuse navigation options for `Link`, `navigate` or `redirect`.
 * Accepts a literal options object and returns it typed for later spreading.
 * @example
 * const opts = linkOptions({ to: '/dashboard', search: { tab: 'home' } })
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/linkOptions
 */
export declare const linkOptions: LinkOptionsFn<'a'>;
export {};
/**
 * Type-check a literal object for use with `Link`, `navigate` or `redirect`.
 * Use to validate and reuse navigation options across your app.
 * @example
 * const opts = linkOptions({ to: '/dashboard', search: { tab: 'home' } })
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/linkOptions
 */


// @filename: /node_modules/@tanstack/react-router/dist/esm/route.d.ts
import { BaseRootRoute, BaseRoute, BaseRouteApi, AnyContext, AnyRoute, AnyRouter, ConstrainLiteral, ErrorComponentProps, NotFoundError, NotFoundRouteProps, Register, RegisteredRouter, ResolveFullPath, ResolveId, ResolveParams, RootRoute as RootRouteCore, RootRouteId, RootRouteOptions, RouteConstraints, Route as RouteCore, RouteIds, RouteMask, RouteOptions, RouteTypesById, RouterCore, ToMaskOptions, UseNavigateResult } from '@tanstack/router-core';
import { default as React } from 'react';
import { UseLoaderDataRoute } from './useLoaderData.js';
import { UseMatchRoute } from './useMatch.js';
import { UseLoaderDepsRoute } from './useLoaderDeps.js';
import { UseParamsRoute } from './useParams.js';
import { UseSearchRoute } from './useSearch.js';
import { UseRouteContextRoute } from './useRouteContext.js';
import { LinkComponentRoute } from './link.js';
declare module '@tanstack/router-core' {
    interface UpdatableRouteOptionsExtensions {
        component?: RouteComponent;
        errorComponent?: false | null | undefined | ErrorRouteComponent;
        notFoundComponent?: NotFoundRouteComponent;
        pendingComponent?: RouteComponent;
    }
    interface RootRouteOptionsExtensions {
        shellComponent?: ({ children, }: {
            children: React.ReactNode;
        }) => React.ReactNode;
    }
    interface RouteExtensions<in out TId extends string, in out TFullPath extends string> {
        useMatch: UseMatchRoute<TId>;
        useRouteContext: UseRouteContextRoute<TId>;
        useSearch: UseSearchRoute<TId>;
        useParams: UseParamsRoute<TId>;
        useLoaderDeps: UseLoaderDepsRoute<TId>;
        useLoaderData: UseLoaderDataRoute<TId>;
        useNavigate: () => UseNavigateResult<TFullPath>;
        Link: LinkComponentRoute<TFullPath>;
    }
}
/**
 * Returns a route-specific API that exposes type-safe hooks pre-bound
 * to a single route ID. Useful for consuming a route's APIs from files
 * where the route object isn't directly imported (e.g. code-split files).
 *
 * @param id Route ID string literal for the target route.
 * @returns A `RouteApi` instance bound to the given route ID.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/getRouteApiFunction
 */
export declare function getRouteApi<const TId, TRouter extends AnyRouter = RegisteredRouter>(id: ConstrainLiteral<TId, RouteIds<TRouter['routeTree']>>): RouteApi<TId, TRouter>;
export declare class RouteApi<TId, TRouter extends AnyRouter = RegisteredRouter> extends BaseRouteApi<TId, TRouter> {
    /**
     * @deprecated Use the `getRouteApi` function instead.
     */
    constructor({ id }: {
        id: TId;
    });
    useMatch: UseMatchRoute<TId>;
    useRouteContext: UseRouteContextRoute<TId>;
    useSearch: UseSearchRoute<TId>;
    useParams: UseParamsRoute<TId>;
    useLoaderDeps: UseLoaderDepsRoute<TId>;
    useLoaderData: UseLoaderDataRoute<TId>;
    useNavigate: () => UseNavigateResult<RouteTypesById<TRouter, TId>["fullPath"]>;
    notFound: (opts?: NotFoundError) => NotFoundError;
    Link: LinkComponentRoute<RouteTypesById<TRouter, TId>['fullPath']>;
}
export declare class Route<in out TRegister = unknown, in out TParentRoute extends RouteConstraints['TParentRoute'] = AnyRoute, in out TPath extends RouteConstraints['TPath'] = '/', in out TFullPath extends RouteConstraints['TFullPath'] = ResolveFullPath<TParentRoute, TPath>, in out TCustomId extends RouteConstraints['TCustomId'] = string, in out TId extends RouteConstraints['TId'] = ResolveId<TParentRoute, TCustomId, TPath>, in out TSearchValidator = undefined, in out TParams = ResolveParams<TPath>, in out TRouterContext = AnyContext, in out TRouteContextFn = AnyContext, in out TBeforeLoadFn = AnyContext, in out TLoaderDeps extends Record<string, any> = {}, in out TLoaderFn = undefined, in out TChildren = unknown, in out TFileRouteTypes = unknown, in out TSSR = unknown, in out TServerMiddlewares = unknown, in out THandlers = undefined> extends BaseRoute<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers> implements RouteCore<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers> {
    /**
     * @deprecated Use the `createRoute` function instead.
     */
    constructor(options?: RouteOptions<TRegister, TParentRoute, TId, TCustomId, TFullPath, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares, THandlers>);
    useMatch: UseMatchRoute<TId>;
    useRouteContext: UseRouteContextRoute<TId>;
    useSearch: UseSearchRoute<TId>;
    useParams: UseParamsRoute<TId>;
    useLoaderDeps: UseLoaderDepsRoute<TId>;
    useLoaderData: UseLoaderDataRoute<TId>;
    useNavigate: () => UseNavigateResult<TFullPath>;
    Link: LinkComponentRoute<TFullPath>;
}
/**
 * Creates a non-root Route instance for code-based routing.
 *
 * Use this to define a route that will be composed into a route tree
 * (typically via a parent route's `addChildren`). If you're using file-based
 * routing, prefer `createFileRoute`.
 *
 * @param options Route options (path, component, loader, context, etc.).
 * @returns A Route instance to be attached to the route tree.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRouteFunction
 */
export declare function createRoute<TRegister = unknown, TParentRoute extends RouteConstraints['TParentRoute'] = AnyRoute, TPath extends RouteConstraints['TPath'] = '/', TFullPath extends RouteConstraints['TFullPath'] = ResolveFullPath<TParentRoute, TPath>, TCustomId extends RouteConstraints['TCustomId'] = string, TId extends RouteConstraints['TId'] = ResolveId<TParentRoute, TCustomId, TPath>, TSearchValidator = undefined, TParams = ResolveParams<TPath>, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TChildren = unknown, TSSR = unknown, const TServerMiddlewares = unknown>(options: RouteOptions<TRegister, TParentRoute, TId, TCustomId, TFullPath, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, AnyContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares>): Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, AnyContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TSSR, TServerMiddlewares>;
export type AnyRootRoute = RootRoute<any, any, any, any, any, any, any, any, any, any, any>;
/**
 * Creates a root route factory that requires a router context type.
 *
 * Use when your root route expects `context` to be provided to `createRouter`.
 * The returned function behaves like `createRootRoute` but enforces a context type.
 *
 * @returns A factory function to configure and return a root route.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRootRouteWithContextFunction
 */
export declare function createRootRouteWithContext<TRouterContext extends {}>(): <TRegister = Register, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TSearchValidator = undefined, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TSSR = unknown, TServerMiddlewares = unknown>(options?: RootRouteOptions<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TSSR, TServerMiddlewares>) => RootRoute<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, unknown, unknown, TSSR, TServerMiddlewares, undefined>;
/**
 * @deprecated Use the `createRootRouteWithContext` function instead.
 */
export declare const rootRouteWithContext: typeof createRootRouteWithContext;
export declare class RootRoute<in out TRegister = unknown, in out TSearchValidator = undefined, in out TRouterContext = {}, in out TRouteContextFn = AnyContext, in out TBeforeLoadFn = AnyContext, in out TLoaderDeps extends Record<string, any> = {}, in out TLoaderFn = undefined, in out TChildren = unknown, in out TFileRouteTypes = unknown, in out TSSR = unknown, in out TServerMiddlewares = unknown, in out THandlers = undefined> extends BaseRootRoute<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers> implements RootRouteCore<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers> {
    /**
     * @deprecated `RootRoute` is now an internal implementation detail. Use `createRootRoute()` instead.
     */
    constructor(options?: RootRouteOptions<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TSSR, TServerMiddlewares, THandlers>);
    useMatch: UseMatchRoute<RootRouteId>;
    useRouteContext: UseRouteContextRoute<RootRouteId>;
    useSearch: UseSearchRoute<RootRouteId>;
    useParams: UseParamsRoute<RootRouteId>;
    useLoaderDeps: UseLoaderDepsRoute<RootRouteId>;
    useLoaderData: UseLoaderDataRoute<RootRouteId>;
    useNavigate: () => UseNavigateResult<"/">;
    Link: LinkComponentRoute<'/'>;
}
/**
 * Creates a root Route instance used to build your route tree.
 *
 * Typically paired with `createRouter({ routeTree })`. If you need to require
 * a typed router context, use `createRootRouteWithContext` instead.
 *
 * @param options Root route options (component, error, pending, etc.).
 * @returns A root route instance.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRootRouteFunction
 */
export declare function createRootRoute<TRegister = Register, TSearchValidator = undefined, TRouterContext = {}, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TSSR = unknown, const TServerMiddlewares = unknown, THandlers = undefined>(options?: RootRouteOptions<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TSSR, TServerMiddlewares, THandlers>): RootRoute<TRegister, TSearchValidator, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, unknown, unknown, TSSR, TServerMiddlewares, THandlers>;
export declare function createRouteMask<TRouteTree extends AnyRoute, TFrom extends string, TTo extends string>(opts: {
    routeTree: TRouteTree;
} & ToMaskOptions<RouterCore<TRouteTree, 'never', boolean>, TFrom, TTo>): RouteMask<TRouteTree>;
export interface DefaultRouteTypes<TProps> {
    component: ((props: TProps) => any) | React.LazyExoticComponent<(props: TProps) => any>;
}
export interface RouteTypes<TProps> extends DefaultRouteTypes<TProps> {
}
export type AsyncRouteComponent<TProps> = RouteTypes<TProps>['component'] & {
    preload?: () => Promise<void>;
};
export type RouteComponent = AsyncRouteComponent<{}>;
export type ErrorRouteComponent = AsyncRouteComponent<ErrorComponentProps>;
export type NotFoundRouteComponent = RouteTypes<NotFoundRouteProps>['component'];
export declare class NotFoundRoute<TRegister, TParentRoute extends AnyRootRoute, TRouterContext = AnyContext, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TSearchValidator = undefined, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TChildren = unknown, TSSR = unknown, TServerMiddlewares = unknown> extends Route<TRegister, TParentRoute, '/404', '/404', '404', '404', TSearchValidator, {}, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TSSR, TServerMiddlewares> {
    constructor(options: Omit<RouteOptions<TRegister, TParentRoute, string, string, string, string, TSearchValidator, {}, TLoaderDeps, TLoaderFn, TRouterContext, TRouteContextFn, TBeforeLoadFn, TSSR, TServerMiddlewares>, 'caseSensitive' | 'parseParams' | 'stringifyParams' | 'path' | 'id' | 'params'>);
}


// @filename: /node_modules/@tanstack/react-router/dist/esm/CatchBoundary.d.ts
import { ErrorRouteComponent } from './route.js';
import { ErrorInfo } from 'react';
import * as React from 'react';
export declare function CatchBoundary(props: {
    getResetKey: () => number | string;
    children: React.ReactNode;
    errorComponent?: ErrorRouteComponent;
    onCatch?: (error: Error, errorInfo: ErrorInfo) => void;
}): import("react/jsx-runtime").JSX.Element;
export declare function ErrorComponent({ error }: {
    error: any;
}): import("react/jsx-runtime").JSX.Element;


// @filename: /node_modules/@tanstack/react-router/dist/esm/ClientOnly.d.ts
import { default as React } from 'react';
export interface ClientOnlyProps {
    /**
     * The children to render when the JS is loaded.
     */
    children: React.ReactNode;
    /**
     * The fallback component to render if the JS is not yet loaded.
     */
    fallback?: React.ReactNode;
}
/**
 * Render the children only after the JS has loaded client-side. Use an optional
 * fallback component if the JS is not yet loaded.
 *
 * @example
 * Render a Chart component if JS loads, renders a simple FakeChart
 * component server-side or if there is no JS. The FakeChart can have only the
 * UI without the behavior or be a loading spinner or skeleton.
 *
 * ```tsx
 * return (
 *   <ClientOnly fallback={<FakeChart />}>
 *     <Chart />
 *   </ClientOnly>
 * )
 * ```
 */
export declare function ClientOnly({ children, fallback }: ClientOnlyProps): import("react/jsx-runtime").JSX.Element;
/**
 * Return a boolean indicating if the JS has been hydrated already.
 * When doing Server-Side Rendering, the result will always be false.
 * When doing Client-Side Rendering, the result will always be false on the
 * first render and true from then on. Even if a new component renders it will
 * always start with true.
 *
 * @example
 * ```tsx
 * // Disable a button that needs JS to work.
 * let hydrated = useHydrated()
 * return (
 *   <button type="button" disabled={!hydrated} onClick={doSomethingCustom}>
 *     Click me
 *   </button>
 * )
 * ```
 * @returns True if the JS has been hydrated already, false otherwise.
 */
export declare function useHydrated(): boolean;


// @filename: /node_modules/@tanstack/react-router/dist/esm/utils.d.ts
import * as React from 'react';
/**
 * React.use if available (React 19+), undefined otherwise.
 * Use dynamic lookup to avoid Webpack compilation errors with React 18.
 */
export declare const reactUse: (<T>(usable: Promise<T> | React.Context<T>) => T) | undefined;
export declare function useStableCallback<T extends (...args: Array<any>) => any>(fn: T): T;
export declare const useLayoutEffect: typeof React.useLayoutEffect;
/**
 * Taken from https://www.developerway.com/posts/implementing-advanced-use-previous-hook#part3
 */
export declare function usePrevious<T>(value: T): T | null;
/**
 * React hook to wrap `IntersectionObserver`.
 *
 * This hook will create an `IntersectionObserver` and observe the ref passed to it.
 *
 * When the intersection changes, the callback will be called with the `IntersectionObserverEntry`.
 *
 * @param ref - The ref to observe
 * @param intersectionObserverOptions - The options to pass to the IntersectionObserver
 * @param options - The options to pass to the hook
 * @param callback - The callback to call when the intersection changes
 * @returns The IntersectionObserver instance
 * @example
 * ```tsx
 * const MyComponent = () => {
 * const ref = React.useRef<HTMLDivElement>(null)
 * useIntersectionObserver(
 *  ref,
 *  (entry) => { doSomething(entry) },
 *  { rootMargin: '10px' },
 *  { disabled: false }
 * )
 * return <div ref={ref} />
 * ```
 */
export declare function useIntersectionObserver<T extends Element>(ref: React.RefObject<T | null>, callback: (entry: IntersectionObserverEntry | undefined) => void, intersectionObserverOptions?: IntersectionObserverInit, options?: {
    disabled?: boolean;
}): void;
/**
 * React hook to take a `React.ForwardedRef` and returns a `ref` that can be used on a DOM element.
 *
 * @param ref - The forwarded ref
 * @returns The inner ref returned by `useRef`
 * @example
 * ```tsx
 * const MyComponent = React.forwardRef((props, ref) => {
 *  const innerRef = useForwardedRef(ref)
 *  return <div ref={innerRef} />
 * })
 * ```
 */
export declare function useForwardedRef<T>(ref?: React.ForwardedRef<T>): React.RefObject<T | null>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/fileRoute.d.ts
import { UseParamsRoute } from './useParams.js';
import { UseMatchRoute } from './useMatch.js';
import { UseSearchRoute } from './useSearch.js';
import { AnyContext, AnyRoute, AnyRouter, Constrain, ConstrainLiteral, FileBaseRouteOptions, FileRoutesByPath, LazyRouteOptions, Register, RegisteredRouter, ResolveParams, Route, RouteById, RouteConstraints, RouteIds, RouteLoaderEntry, UpdatableRouteOptions, UseNavigateResult } from '@tanstack/router-core';
import { UseLoaderDepsRoute } from './useLoaderDeps.js';
import { UseLoaderDataRoute } from './useLoaderData.js';
import { UseRouteContextRoute } from './useRouteContext.js';
/**
 * Creates a file-based Route factory for a given path.
 *
 * Used by TanStack Router's file-based routing to associate a file with a
 * route. The returned function accepts standard route options. In normal usage
 * the `path` string is inserted and maintained by the `tsr` generator.
 *
 * @param path File path literal for the route (usually auto-generated).
 * @returns A function that accepts Route options and returns a Route instance.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createFileRouteFunction
 */
export declare function createFileRoute<TFilePath extends keyof FileRoutesByPath, TParentRoute extends AnyRoute = FileRoutesByPath[TFilePath]['parentRoute'], TId extends RouteConstraints['TId'] = FileRoutesByPath[TFilePath]['id'], TPath extends RouteConstraints['TPath'] = FileRoutesByPath[TFilePath]['path'], TFullPath extends RouteConstraints['TFullPath'] = FileRoutesByPath[TFilePath]['fullPath']>(path?: TFilePath): FileRoute<TFilePath, TParentRoute, TId, TPath, TFullPath>['createRoute'];
/**
  @deprecated It's no longer recommended to use the `FileRoute` class directly.
  Instead, use `createFileRoute('/path/to/file')(options)` to create a file route.
*/
export declare class FileRoute<TFilePath extends keyof FileRoutesByPath, TParentRoute extends AnyRoute = FileRoutesByPath[TFilePath]['parentRoute'], TId extends RouteConstraints['TId'] = FileRoutesByPath[TFilePath]['id'], TPath extends RouteConstraints['TPath'] = FileRoutesByPath[TFilePath]['path'], TFullPath extends RouteConstraints['TFullPath'] = FileRoutesByPath[TFilePath]['fullPath']> {
    path?: TFilePath | undefined;
    silent?: boolean;
    constructor(path?: TFilePath | undefined, _opts?: {
        silent: boolean;
    });
    createRoute: <TRegister = Register, TSearchValidator = undefined, TParams = ResolveParams<TPath>, TRouteContextFn = AnyContext, TBeforeLoadFn = AnyContext, TLoaderDeps extends Record<string, any> = {}, TLoaderFn = undefined, TChildren = unknown, TSSR = unknown, const TMiddlewares = unknown, THandlers = undefined>(options?: FileBaseRouteOptions<TRegister, TParentRoute, TId, TPath, TSearchValidator, TParams, TLoaderDeps, TLoaderFn, AnyContext, TRouteContextFn, TBeforeLoadFn, AnyContext, TSSR, TMiddlewares, THandlers> & UpdatableRouteOptions<TParentRoute, TId, TFullPath, TParams, TSearchValidator, TLoaderFn, TLoaderDeps, AnyContext, TRouteContextFn, TBeforeLoadFn>) => Route<TRegister, TParentRoute, TPath, TFullPath, TFilePath, TId, TSearchValidator, TParams, AnyContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, unknown, TSSR, TMiddlewares, THandlers>;
}
/**
  @deprecated It's recommended not to split loaders into separate files.
  Instead, place the loader function in the main route file via `createFileRoute`.
*/
export declare function FileRouteLoader<TFilePath extends keyof FileRoutesByPath, TRoute extends FileRoutesByPath[TFilePath]['preLoaderRoute']>(_path: TFilePath): <TLoaderFn>(loaderFn: Constrain<TLoaderFn, RouteLoaderEntry<Register, TRoute['parentRoute'], TRoute['types']['id'], TRoute['types']['params'], TRoute['types']['loaderDeps'], TRoute['types']['routerContext'], TRoute['types']['routeContextFn'], TRoute['types']['beforeLoadFn']>>) => TLoaderFn;
declare module '@tanstack/router-core' {
    interface LazyRoute<in out TRoute extends AnyRoute> {
        useMatch: UseMatchRoute<TRoute['id']>;
        useRouteContext: UseRouteContextRoute<TRoute['id']>;
        useSearch: UseSearchRoute<TRoute['id']>;
        useParams: UseParamsRoute<TRoute['id']>;
        useLoaderDeps: UseLoaderDepsRoute<TRoute['id']>;
        useLoaderData: UseLoaderDataRoute<TRoute['id']>;
        useNavigate: () => UseNavigateResult<TRoute['fullPath']>;
    }
}
export declare class LazyRoute<TRoute extends AnyRoute> {
    options: {
        id: string;
    } & LazyRouteOptions;
    constructor(opts: {
        id: string;
    } & LazyRouteOptions);
    useMatch: UseMatchRoute<TRoute['id']>;
    useRouteContext: UseRouteContextRoute<TRoute['id']>;
    useSearch: UseSearchRoute<TRoute['id']>;
    useParams: UseParamsRoute<TRoute['id']>;
    useLoaderDeps: UseLoaderDepsRoute<TRoute['id']>;
    useLoaderData: UseLoaderDataRoute<TRoute['id']>;
    useNavigate: () => UseNavigateResult<TRoute["fullPath"]>;
}
/**
 * Creates a lazily-configurable code-based route stub by ID.
 *
 * Use this for code-splitting with code-based routes. The returned function
 * accepts only non-critical route options like `component`, `pendingComponent`,
 * `errorComponent`, and `notFoundComponent` which are applied when the route
 * is matched.
 *
 * @param id Route ID string literal to associate with the lazy route.
 * @returns A function that accepts lazy route options and returns a `LazyRoute`.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createLazyRouteFunction
 */
export declare function createLazyRoute<TRouter extends AnyRouter = RegisteredRouter, TId extends string = string, TRoute extends AnyRoute = RouteById<TRouter['routeTree'], TId>>(id: ConstrainLiteral<TId, RouteIds<TRouter['routeTree']>>): (opts: LazyRouteOptions) => LazyRoute<TRoute>;
/**
 * Creates a lazily-configurable file-based route stub by file path.
 *
 * Use this for code-splitting with file-based routes (eg. `.lazy.tsx` files).
 * The returned function accepts only non-critical route options like
 * `component`, `pendingComponent`, `errorComponent`, and `notFoundComponent`.
 *
 * @param id File path literal for the route file.
 * @returns A function that accepts lazy route options and returns a `LazyRoute`.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createLazyFileRouteFunction
 */
export declare function createLazyFileRoute<TFilePath extends keyof FileRoutesByPath, TRoute extends FileRoutesByPath[TFilePath]['preLoaderRoute']>(id: TFilePath): (opts: LazyRouteOptions) => LazyRoute<TRoute>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/history.d.ts
declare module '@tanstack/history' {
    interface HistoryState {
        __tempLocation?: HistoryLocation;
        __tempKey?: string;
        __hashScrollIntoViewOptions?: boolean | ScrollIntoViewOptions;
    }
}
export {};


// @filename: /node_modules/@tanstack/react-router/dist/esm/lazyRouteComponent.d.ts
import { AsyncRouteComponent } from './route.js';
/**
 * Wrap a dynamic import to create a route component that supports
 * `.preload()` and friendly reload-on-module-missing behavior.
 *
 * @param importer Function returning a module promise
 * @param exportName Named export to use (default: `default`)
 * @returns A lazy route component compatible with TanStack Router
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/lazyRouteComponentFunction
 */
export declare function lazyRouteComponent<T extends Record<string, any>, TKey extends keyof T = 'default'>(importer: () => Promise<T>, exportName?: TKey): T[TKey] extends (props: infer TProps) => any ? AsyncRouteComponent<TProps> : never;


// @filename: /node_modules/@tanstack/react-router/dist/esm/Matches.d.ts
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
import { AnyRouter, DeepPartial, Expand, MakeOptionalPathParams, MakeOptionalSearchParams, MakeRouteMatchUnion, MaskOptions, MatchRouteOptions, RegisteredRouter, ResolveRoute, ToSubOptionsProps } from '@tanstack/router-core';
import * as React from 'react';
declare module '@tanstack/router-core' {
    interface RouteMatchExtensions {
        meta?: Array<React.JSX.IntrinsicElements['meta'] | undefined>;
        links?: Array<React.JSX.IntrinsicElements['link'] | undefined>;
        scripts?: Array<React.JSX.IntrinsicElements['script'] | undefined>;
        styles?: Array<React.JSX.IntrinsicElements['style'] | undefined>;
        headScripts?: Array<React.JSX.IntrinsicElements['script'] | undefined>;
    }
}
/**
 * Internal component that renders the router's active match tree with
 * suspense, error, and not-found boundaries. Rendered by `RouterProvider`.
 */
export declare function Matches(): import("react/jsx-runtime").JSX.Element;
export type UseMatchRouteOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = undefined, TMaskFrom extends string = TFrom, TMaskTo extends string = ''> = ToSubOptionsProps<TRouter, TFrom, TTo> & DeepPartial<MakeOptionalSearchParams<TRouter, TFrom, TTo>> & DeepPartial<MakeOptionalPathParams<TRouter, TFrom, TTo>> & MaskOptions<TRouter, TMaskFrom, TMaskTo> & MatchRouteOptions;
/**
 * Create a matcher function for testing locations against route definitions.
 *
 * The returned function accepts standard navigation options (`to`, `params`,
 * `search`, etc.) and returns either `false` (no match) or the matched params
 * object when the route matches the current or pending location.
 *
 * Useful for conditional rendering and active UI states.
 *
 * @returns A `matchRoute(options)` function that returns `false` or params.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useMatchRouteHook
 */
export declare function useMatchRoute<TRouter extends AnyRouter = RegisteredRouter>(): <const TFrom extends string = string, const TTo extends string | undefined = undefined, const TMaskFrom extends string = TFrom, const TMaskTo extends string = "">(opts: UseMatchRouteOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>) => false | Expand<ResolveRoute<TRouter, TFrom, TTo>["types"]["allParams"]>;
export type MakeMatchRouteOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = undefined, TMaskFrom extends string = TFrom, TMaskTo extends string = ''> = UseMatchRouteOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo> & {
    children?: ((params?: Expand<ResolveRoute<TRouter, TFrom, TTo>['types']['allParams']>) => React.ReactNode) | React.ReactNode;
};
/**
 * Component that conditionally renders its children based on whether a route
 * matches the provided `from`/`to` options. If `children` is a function, it
 * receives the matched params object.
 *
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/matchRouteComponent
 */
export declare function MatchRoute<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string = string, const TTo extends string | undefined = undefined, const TMaskFrom extends string = TFrom, const TMaskTo extends string = ''>(props: MakeMatchRouteOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>): any;
export interface UseMatchesBaseOptions<TRouter extends AnyRouter, TSelected, TStructuralSharing> {
    select?: (matches: Array<MakeRouteMatchUnion<TRouter>>) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
}
export type UseMatchesResult<TRouter extends AnyRouter, TSelected> = unknown extends TSelected ? Array<MakeRouteMatchUnion<TRouter>> : TSelected;
export declare function useMatches<TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseMatchesBaseOptions<TRouter, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>): UseMatchesResult<TRouter, TSelected>;
/**
 * Read the full array of active route matches or select a derived subset.
 *
 * Useful for debugging, breadcrumbs, or aggregating metadata across matches.
 *
 * @returns The array of matches (or the selected value).
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useMatchesHook
 */
/**
 * Read the full array of active route matches or select a derived subset.
 *
 * Useful for debugging, breadcrumbs, or aggregating metadata across matches.
 *
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useMatchesHook
 */
export declare function useParentMatches<TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseMatchesBaseOptions<TRouter, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>): UseMatchesResult<TRouter, TSelected>;
/**
 * Read the array of active route matches that are children of the current
 * match (or selected parent) in the match tree.
 */
export declare function useChildMatches<TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseMatchesBaseOptions<TRouter, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>): UseMatchesResult<TRouter, TSelected>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/Match.d.ts
import * as React from 'react';
export declare const Match: React.MemoExoticComponent<({ matchId, }: {
    matchId: string;
}) => import("react/jsx-runtime").JSX.Element>;
export declare const MatchInner: React.MemoExoticComponent<({ matchId, }: {
    matchId: string;
}) => any>;
/**
 * Render the next child match in the route tree. Typically used inside
 * a route component to render nested routes.
 *
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/outletComponent
 */
export declare const Outlet: React.MemoExoticComponent<() => import("react/jsx-runtime").JSX.Element | null>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/router.d.ts
import { RouterCore, AnyRoute, CreateRouterFn, RouterConstructorOptions, TrailingSlashOption } from '@tanstack/router-core';
import { RouterHistory } from '@tanstack/history';
import { ErrorRouteComponent, NotFoundRouteComponent, RouteComponent } from './route.js';
declare module '@tanstack/router-core' {
    interface RouterOptionsExtensions {
        /**
         * The default `component` a route should use if no component is provided.
         *
         * @default Outlet
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultcomponent-property)
         */
        defaultComponent?: RouteComponent;
        /**
         * The default `errorComponent` a route should use if no error component is provided.
         *
         * @default ErrorComponent
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaulterrorcomponent-property)
         * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#handling-errors-with-routeoptionserrorcomponent)
         */
        defaultErrorComponent?: ErrorRouteComponent;
        /**
         * The default `pendingComponent` a route should use if no pending component is provided.
         *
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultpendingcomponent-property)
         * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#showing-a-pending-component)
         */
        defaultPendingComponent?: RouteComponent;
        /**
         * The default `notFoundComponent` a route should use if no notFound component is provided.
         *
         * @default NotFound
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultnotfoundcomponent-property)
         * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/not-found-errors#default-router-wide-not-found-handling)
         */
        defaultNotFoundComponent?: NotFoundRouteComponent;
        /**
         * A component that will be used to wrap the entire router.
         *
         * This is useful for providing a context to the entire router.
         *
         * Only non-DOM-rendering components like providers should be used, anything else will cause a hydration error.
         *
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#wrap-property)
         */
        Wrap?: (props: {
            children: any;
        }) => React.JSX.Element;
        /**
         * A component that will be used to wrap the inner contents of the router.
         *
         * This is useful for providing a context to the inner contents of the router where you also need access to the router context and hooks.
         *
         * Only non-DOM-rendering components like providers should be used, anything else will cause a hydration error.
         *
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#innerwrap-property)
         */
        InnerWrap?: (props: {
            children: any;
        }) => React.JSX.Element;
        /**
         * The default `onCatch` handler for errors caught by the Router ErrorBoundary
         *
         * @link [API Docs](https://tanstack.com/router/latest/docs/framework/react/api/router/RouterOptionsType#defaultoncatch-property)
         * @link [Guide](https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#handling-errors-with-routeoptionsoncatch)
         */
        defaultOnCatch?: (error: Error, errorInfo: React.ErrorInfo) => void;
    }
}
/**
 * Creates a new Router instance for React.
 *
 * Pass the returned router to `RouterProvider` to enable routing.
 * Notable options: `routeTree` (your route definitions) and `context`
 * (required if the root route was created with `createRootRouteWithContext`).
 *
 * @param options Router options used to configure the router.
 * @returns A Router instance to be provided to `RouterProvider`.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRouterFunction
 */
export declare const createRouter: CreateRouterFn;
export declare class Router<in out TRouteTree extends AnyRoute, in out TTrailingSlashOption extends TrailingSlashOption = 'never', in out TDefaultStructuralSharingOption extends boolean = false, in out TRouterHistory extends RouterHistory = RouterHistory, in out TDehydrated extends Record<string, any> = Record<string, any>> extends RouterCore<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated> {
    constructor(options: RouterConstructorOptions<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>);
}


// @filename: /node_modules/@tanstack/react-router/dist/esm/RouterProvider.d.ts
import { AnyRouter, RegisteredRouter, RouterOptions } from '@tanstack/router-core';
import * as React from 'react';
/**
 * Low-level provider that places the router into React context and optionally
 * updates router options from props. Most apps should use `RouterProvider`.
 */
export declare function RouterContextProvider<TRouter extends AnyRouter = RegisteredRouter, TDehydrated extends Record<string, any> = Record<string, any>>({ router, children, ...rest }: RouterProps<TRouter, TDehydrated> & {
    children: React.ReactNode;
}): import("react/jsx-runtime").JSX.Element;
/**
 * Top-level component that renders the active route matches and provides the
 * router to the React tree via context.
 *
 * Accepts the same options as `createRouter` via props to update the router
 * instance after creation.
 *
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/createRouterFunction
 */
export declare function RouterProvider<TRouter extends AnyRouter = RegisteredRouter, TDehydrated extends Record<string, any> = Record<string, any>>({ router, ...rest }: RouterProps<TRouter, TDehydrated>): import("react/jsx-runtime").JSX.Element;
export type RouterProps<TRouter extends AnyRouter = RegisteredRouter, TDehydrated extends Record<string, any> = Record<string, any>> = Omit<RouterOptions<TRouter['routeTree'], NonNullable<TRouter['options']['trailingSlash']>, NonNullable<TRouter['options']['defaultStructuralSharing']>, TRouter['history'], TDehydrated>, 'context'> & {
    router: TRouter;
    context?: Partial<RouterOptions<TRouter['routeTree'], NonNullable<TRouter['options']['trailingSlash']>, NonNullable<TRouter['options']['defaultStructuralSharing']>, TRouter['history'], TDehydrated>['context']>;
};


// @filename: /node_modules/@tanstack/react-router/dist/esm/ScrollRestoration.d.ts
import { ParsedLocation, ScrollRestorationEntry, ScrollRestorationOptions } from '@tanstack/router-core';
/**
 * @deprecated Use the `scrollRestoration` router option instead.
 */
export declare function ScrollRestoration(_props: ScrollRestorationOptions): null;
export declare function useElementScrollRestoration(options: ({
    id: string;
    getElement?: () => Window | Element | undefined | null;
} | {
    id?: string;
    getElement: () => Window | Element | undefined | null;
}) & {
    getKey?: (location: ParsedLocation) => string;
}): ScrollRestorationEntry | undefined;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useBlocker.d.ts
import { HistoryAction } from '@tanstack/history';
import { AnyRoute, AnyRouter, ParseRoute, RegisteredRouter } from '@tanstack/router-core';
import * as React from 'react';
type ShouldBlockFnLocation<out TRouteId, out TFullPath, out TAllParams, out TFullSearchSchema> = {
    routeId: TRouteId;
    fullPath: TFullPath;
    pathname: string;
    params: TAllParams;
    search: TFullSearchSchema;
};
type MakeShouldBlockFnLocationUnion<TRouter extends AnyRouter = RegisteredRouter, TRoute extends AnyRoute = ParseRoute<TRouter['routeTree']>> = TRoute extends any ? ShouldBlockFnLocation<TRoute['id'], TRoute['fullPath'], TRoute['types']['allParams'], TRoute['types']['fullSearchSchema']> : never;
type BlockerResolver<TRouter extends AnyRouter = RegisteredRouter> = {
    status: 'blocked';
    current: MakeShouldBlockFnLocationUnion<TRouter>;
    next: MakeShouldBlockFnLocationUnion<TRouter>;
    action: HistoryAction;
    proceed: () => void;
    reset: () => void;
} | {
    status: 'idle';
    current: undefined;
    next: undefined;
    action: undefined;
    proceed: undefined;
    reset: undefined;
};
type ShouldBlockFnArgs<TRouter extends AnyRouter = RegisteredRouter> = {
    current: MakeShouldBlockFnLocationUnion<TRouter>;
    next: MakeShouldBlockFnLocationUnion<TRouter>;
    action: HistoryAction;
};
export type ShouldBlockFn<TRouter extends AnyRouter = RegisteredRouter> = (args: ShouldBlockFnArgs<TRouter>) => boolean | Promise<boolean>;
export type UseBlockerOpts<TRouter extends AnyRouter = RegisteredRouter, TWithResolver extends boolean = boolean> = {
    shouldBlockFn: ShouldBlockFn<TRouter>;
    enableBeforeUnload?: boolean | (() => boolean);
    disabled?: boolean;
    withResolver?: TWithResolver;
};
type LegacyBlockerFn = () => Promise<any> | any;
type LegacyBlockerOpts = {
    blockerFn?: LegacyBlockerFn;
    condition?: boolean | any;
};
export declare function useBlocker<TRouter extends AnyRouter = RegisteredRouter, TWithResolver extends boolean = false>(opts: UseBlockerOpts<TRouter, TWithResolver>): TWithResolver extends true ? BlockerResolver<TRouter> : void;
/**
 * @deprecated Use the shouldBlockFn property instead
 */
export declare function useBlocker(blockerFnOrOpts?: LegacyBlockerOpts): BlockerResolver;
/**
 * @deprecated Use the UseBlockerOpts object syntax instead
 */
export declare function useBlocker(blockerFn?: LegacyBlockerFn, condition?: boolean | any): BlockerResolver;
export declare function Block<TRouter extends AnyRouter = RegisteredRouter, TWithResolver extends boolean = boolean>(opts: PromptProps<TRouter, TWithResolver>): React.ReactNode;
/**
 *  @deprecated Use the UseBlockerOpts property instead
 */
export declare function Block(opts: LegacyPromptProps): React.ReactNode;
type LegacyPromptProps = {
    blockerFn?: LegacyBlockerFn;
    condition?: boolean | any;
    children?: React.ReactNode | ((params: BlockerResolver) => React.ReactNode);
};
type PromptProps<TRouter extends AnyRouter = RegisteredRouter, TWithResolver extends boolean = boolean, TParams = TWithResolver extends true ? BlockerResolver<TRouter> : void> = UseBlockerOpts<TRouter, TWithResolver> & {
    children?: React.ReactNode | ((params: TParams) => React.ReactNode);
};
export {};


// @filename: /node_modules/@tanstack/react-router/dist/esm/useNavigate.d.ts
import { AnyRouter, FromPathOption, NavigateOptions, RegisteredRouter, UseNavigateResult } from '@tanstack/router-core';
/**
 * Imperative navigation hook.
 *
 * Returns a stable `navigate(options)` function to change the current location
 * programmatically. Prefer the `Link` component for user-initiated navigation,
 * and use this hook from effects, callbacks, or handlers where imperative
 * navigation is required.
 *
 * Options:
 * - `from`: Optional route base used to resolve relative `to` paths.
 *
 * @returns A function that accepts `NavigateOptions`.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useNavigateHook
 */
export declare function useNavigate<TRouter extends AnyRouter = RegisteredRouter, TDefaultFrom extends string = string>(_defaultOpts?: {
    from?: FromPathOption<TRouter, TDefaultFrom>;
}): UseNavigateResult<TDefaultFrom>;
/**
 * Component that triggers a navigation when rendered. Navigation executes
 * in an effect after mount/update.
 *
 * Props are the same as `NavigateOptions` used by `navigate()`.
 *
 * @returns null
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/navigateComponent
 */
export declare function Navigate<TRouter extends AnyRouter = RegisteredRouter, const TFrom extends string = string, const TTo extends string | undefined = undefined, const TMaskFrom extends string = TFrom, const TMaskTo extends string = ''>(props: NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>): null;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useRouter.d.ts
import { AnyRouter, RegisteredRouter } from '@tanstack/router-core';
/**
 * Access the current TanStack Router instance from React context.
 * Must be used within a `RouterProvider`.
 *
 * Options:
 * - `warn`: Log a warning if no router context is found (default: true).
 *
 * @returns The registered router instance.
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useRouterHook
 */
export declare function useRouter<TRouter extends AnyRouter = RegisteredRouter>(opts?: {
    warn?: boolean;
}): TRouter;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useRouterState.d.ts
import { AnyRouter, RegisteredRouter, RouterState } from '@tanstack/router-core';
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
export type UseRouterStateOptions<TRouter extends AnyRouter, TSelected, TStructuralSharing> = {
    router?: TRouter;
    select?: (state: RouterState<TRouter['routeTree']>) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
} & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>;
export type UseRouterStateResult<TRouter extends AnyRouter, TSelected> = unknown extends TSelected ? RouterState<TRouter['routeTree']> : TSelected;
/**
 * Subscribe to the router's state store with optional selection and
 * structural sharing for render optimization.
 *
 * Options:
 * - `select`: Project the full router state to a derived slice
 * - `structuralSharing`: Replace-equal semantics for stable references
 * - `router`: Read state from a specific router instance instead of context
 *
 * @returns The selected router state (or the full state by default).
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useRouterStateHook
 */
export declare function useRouterState<TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseRouterStateOptions<TRouter, TSelected, TStructuralSharing>): UseRouterStateResult<TRouter, TSelected>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useLocation.d.ts
import { StructuralSharingOption, ValidateSelected } from './structuralSharing.js';
import { AnyRouter, RegisteredRouter, RouterState } from '@tanstack/router-core';
export interface UseLocationBaseOptions<TRouter extends AnyRouter, TSelected, TStructuralSharing extends boolean = boolean> {
    select?: (state: RouterState<TRouter['routeTree']>['location']) => ValidateSelected<TRouter, TSelected, TStructuralSharing>;
}
export type UseLocationResult<TRouter extends AnyRouter, TSelected> = unknown extends TSelected ? RouterState<TRouter['routeTree']>['location'] : TSelected;
/**
 * Read the current location from the router state with optional selection.
 * Useful for subscribing to just the pieces of location you care about.
 *
 * Options:
 * - `select`: Project the `location` object to a derived value
 * - `structuralSharing`: Enable structural sharing for stable references
 *
 * @returns The current location (or selected value).
 * @link https://tanstack.com/router/latest/docs/framework/react/api/router/useLocationHook
 */
export declare function useLocation<TRouter extends AnyRouter = RegisteredRouter, TSelected = unknown, TStructuralSharing extends boolean = boolean>(opts?: UseLocationBaseOptions<TRouter, TSelected, TStructuralSharing> & StructuralSharingOption<TRouter, TSelected, TStructuralSharing>): UseLocationResult<TRouter, TSelected>;


// @filename: /node_modules/@tanstack/react-router/dist/esm/useCanGoBack.d.ts
export declare function useCanGoBack(): boolean;


// @filename: /node_modules/@tanstack/react-router/dist/esm/not-found.d.ts
import { ErrorInfo } from 'react';
import { NotFoundError } from '@tanstack/router-core';
import * as React from 'react';
export declare function CatchNotFound(props: {
    fallback?: (error: NotFoundError) => React.ReactElement;
    onCatch?: (error: Error, errorInfo: ErrorInfo) => void;
    children: React.ReactNode;
}): import("react/jsx-runtime").JSX.Element;
export declare function DefaultGlobalNotFound(): import("react/jsx-runtime").JSX.Element;


// @filename: /node_modules/@tanstack/react-router/dist/esm/ScriptOnce.d.ts
/**
 * Server-only helper to emit a script tag exactly once during SSR.
 */
export declare function ScriptOnce({ children }: {
    children: string;
}): import("react/jsx-runtime").JSX.Element | null;


// @filename: /node_modules/@tanstack/react-router/dist/esm/Asset.d.ts
import { RouterManagedTag } from '@tanstack/router-core';
import * as React from 'react';
export declare function Asset(asset: RouterManagedTag & {
    nonce?: string;
    preventScriptHoist?: boolean;
}): React.ReactElement | null;


// @filename: /node_modules/@tanstack/react-router/dist/esm/HeadContent.d.ts
import { AssetCrossOriginConfig } from '@tanstack/router-core';
export interface HeadContentProps {
    assetCrossOrigin?: AssetCrossOriginConfig;
}
/**
 * Render route-managed head tags (title, meta, links, styles, head scripts).
 * Place inside the document head of your app shell.
 * @link https://tanstack.com/router/latest/docs/framework/react/guide/document-head-management
 */
export declare function HeadContent(props: HeadContentProps): import("react/jsx-runtime").JSX.Element;


// @filename: /node_modules/@tanstack/react-router/dist/esm/headContentUtils.d.ts
import { AssetCrossOriginConfig, RouterManagedTag } from '@tanstack/router-core';
/**
 * Build the list of head/link/meta/script tags to render for active matches.
 * Used internally by `HeadContent`.
 */
export declare const useTags: (assetCrossOrigin?: AssetCrossOriginConfig) => RouterManagedTag[];


// @filename: /node_modules/@tanstack/react-router/dist/esm/Scripts.d.ts
/**
 * Render body script tags collected from route matches and SSR manifests.
 * Should be placed near the end of the document body.
 */
export declare const Scripts: () => import("react/jsx-runtime").JSX.Element;


// @filename: /node_modules/@tanstack/react-router/dist/esm/ssr/serializer.d.ts
import type * as React from 'react';
declare module '@tanstack/router-core' {
    interface SerializerExtensions {
        ReadableStream: React.JSX.Element;
    }
}


// @filename: /node_modules/@tanstack/react-router/dist/esm/index.d.ts
export { defer, isMatch, joinPaths, cleanPath, trimPathLeft, trimPathRight, trimPath, resolvePath, interpolatePath, rootRouteId, defaultParseSearch, defaultStringifySearch, parseSearchWith, stringifySearchWith, functionalUpdate, replaceEqualDeep, isPlainObject, isPlainArray, deepEqual, createControlledPromise, retainSearchParams, stripSearchParams, createSerializationAdapter, } from '@tanstack/router-core';
export type { AnyRoute, DeferredPromiseState, DeferredPromise, ParsedLocation, RemoveTrailingSlashes, RemoveLeadingSlashes, ActiveOptions, ResolveRelativePath, RootRouteId, AnyPathParams, ResolveParams, ResolveOptionalParams, ResolveRequiredParams, SearchSchemaInput, AnyContext, RouteContext, PreloadableObj, RoutePathOptions, StaticDataRouteOption, RoutePathOptionsIntersection, UpdatableStaticRouteOption, MetaDescriptor, RouteLinkEntry, ParseParamsFn, SearchFilter, ResolveId, InferFullSearchSchema, InferFullSearchSchemaInput, ErrorRouteProps, ErrorComponentProps, NotFoundRouteProps, TrimPath, TrimPathLeft, TrimPathRight, StringifyParamsFn, ParamsOptions, InferAllParams, InferAllContext, LooseReturnType, LooseAsyncReturnType, ContextReturnType, ContextAsyncReturnType, ResolveLoaderData, ResolveRouteContext, SearchSerializer, SearchParser, SearchMiddleware, TrailingSlashOption, Manifest, RouterManagedTag, ControlledPromise, Constrain, Expand, MergeAll, Assign, IntersectAssign, ResolveValidatorInput, ResolveValidatorOutput, Register, AnyValidator, DefaultValidator, ValidatorFn, AnySchema, AnyValidatorAdapter, AnyValidatorFn, AnyValidatorObj, ResolveValidatorInputFn, ResolveValidatorOutputFn, ResolveSearchValidatorInput, ResolveSearchValidatorInputFn, Validator, ValidatorAdapter, ValidatorObj, FileRoutesByPath, RouteById, RootRouteOptions, CreateFileRoute, SerializationAdapter, AnySerializationAdapter, SerializableExtensions, } from '@tanstack/router-core';
export { createHistory, createBrowserHistory, createHashHistory, createMemoryHistory, } from '@tanstack/history';
export type { BlockerFn, HistoryLocation, RouterHistory, ParsedPath, HistoryState, } from '@tanstack/history';
export { useAwaited, Await } from './awaited.js';
export type { AwaitOptions } from './awaited.js';
export { CatchBoundary, ErrorComponent } from './CatchBoundary.js';
export { ClientOnly, useHydrated } from './ClientOnly.js';
export { reactUse, useLayoutEffect } from './utils.js';
export { FileRoute, createFileRoute, FileRouteLoader, LazyRoute, createLazyRoute, createLazyFileRoute, } from './fileRoute.js';
export * from './history.js';
export { lazyRouteComponent } from './lazyRouteComponent.js';
export { useLinkProps, createLink, Link, linkOptions } from './link.js';
export type { InferDescendantToPaths, RelativeToPath, RelativeToParentPath, RelativeToCurrentPath, AbsoluteToPath, RelativeToPathAutoComplete, NavigateOptions, ToOptions, ToMaskOptions, ToSubOptions, ResolveRoute, SearchParamOptions, PathParamOptions, ToPathOption, LinkOptions, MakeOptionalPathParams, FileRouteTypes, RouteContextParameter, BeforeLoadContextParameter, ResolveAllContext, ResolveAllParamsFromParent, ResolveFullSearchSchema, ResolveFullSearchSchemaInput, RouteIds, NavigateFn, BuildLocationFn, FullSearchSchemaOption, MakeRemountDepsOptionsUnion, RemountDepsOptions, ResolveFullPath, AnyRouteWithContext, AnyRouterWithContext, CommitLocationOptions, MatchLocation, UseNavigateResult, AnyRedirect, Redirect, RedirectOptions, ResolvedRedirect, MakeRouteMatch, MakeRouteMatchUnion, RouteMatch, AnyRouteMatch, RouteContextFn, RouteContextOptions, BeforeLoadContextOptions, ContextOptions, RouteOptions, FileBaseRouteOptions, BaseRouteOptions, UpdatableRouteOptions, RouteLoaderFn, LoaderFnContext, LazyRouteOptions, AnyRouter, RegisteredRouter, RouterContextOptions, ControllablePromise, InjectedHtmlEntry, RouterOptions, RouterState, ListenerFn, BuildNextOptions, RouterConstructorOptions, RouterEvents, RouterEvent, RouterListener, RouteConstraints, RouteMask, MatchRouteOptions, CreateLazyFileRoute, } from '@tanstack/router-core';
export type { UseLinkPropsOptions, ActiveLinkOptions, LinkProps, LinkComponent, LinkComponentProps, CreateLinkProps, } from './link.js';
export { Matches, useMatchRoute, MatchRoute, useMatches, useParentMatches, useChildMatches, } from './Matches.js';
export type { UseMatchRouteOptions, MakeMatchRouteOptions } from './Matches.js';
export { Match, Outlet } from './Match.js';
export { useMatch } from './useMatch.js';
export { useLoaderDeps } from './useLoaderDeps.js';
export { useLoaderData } from './useLoaderData.js';
export { redirect, isRedirect, createRouterConfig, DEFAULT_PROTOCOL_ALLOWLIST, } from '@tanstack/router-core';
export { RouteApi, getRouteApi, Route, createRoute, RootRoute, rootRouteWithContext, createRootRoute, createRootRouteWithContext, createRouteMask, NotFoundRoute, } from './route.js';
export type { AnyRootRoute, AsyncRouteComponent, RouteComponent, ErrorRouteComponent, NotFoundRouteComponent, DefaultRouteTypes, RouteTypes, } from './route.js';
export { createRouter, Router } from './router.js';
export { lazyFn, SearchParamError } from '@tanstack/router-core';
export { RouterProvider, RouterContextProvider } from './RouterProvider.js';
export type { RouterProps } from './RouterProvider.js';
export { useElementScrollRestoration, ScrollRestoration, } from './ScrollRestoration.js';
export type { UseBlockerOpts, ShouldBlockFn } from './useBlocker.js';
export { useBlocker, Block } from './useBlocker.js';
export { useNavigate, Navigate } from './useNavigate.js';
export { useParams } from './useParams.js';
export { useSearch } from './useSearch.js';
export { useRouteContext } from './useRouteContext.js';
export { useRouter } from './useRouter.js';
export { useRouterState } from './useRouterState.js';
export { useLocation } from './useLocation.js';
export { useCanGoBack } from './useCanGoBack.js';
export { CatchNotFound, DefaultGlobalNotFound } from './not-found.js';
export { notFound, isNotFound } from '@tanstack/router-core';
export type { NotFoundError } from '@tanstack/router-core';
export type { ValidateLinkOptions, InferStructuralSharing, ValidateUseSearchOptions, ValidateUseParamsOptions, ValidateLinkOptionsArray, } from './typePrimitives.js';
export type { ValidateFromPath, ValidateToPath, ValidateSearch, ValidateParams, InferFrom, InferTo, InferMaskTo, InferMaskFrom, ValidateNavigateOptions, ValidateNavigateOptionsArray, ValidateRedirectOptions, ValidateRedirectOptionsArray, ValidateId, InferStrict, InferShouldThrow, InferSelected, ValidateUseSearchResult, ValidateUseParamsResult, SerializerExtensions, RegisteredSerializableInput, Serializable, } from '@tanstack/router-core';
export { ScriptOnce } from './ScriptOnce.js';
export { Asset } from './Asset.js';
export { HeadContent } from './HeadContent.js';
export { useTags } from './headContentUtils.js';
export { Scripts } from './Scripts.js';
export type * from './ssr/serializer.js';
export { composeRewrites } from '@tanstack/router-core';
export type { LocationRewrite, LocationRewriteFunction, } from '@tanstack/router-core';


// @filename: /node_modules/react/index.d.ts
declare module "react" { export type ReactNode = any; export type CSSProperties = any; export type Ref<T = any> = any; export type ReactElement = any; export type ComponentType<P = any> = any; export type ComponentProps<T> = any; export function createElement(...args: any[]): any; const React: any; export default React; }

// @filename: /node_modules/react/jsx-runtime.d.ts
export namespace JSX { interface IntrinsicElements { [name: string]: any } } export const jsx: any; export const jsxs: any; export const Fragment: any;

// @filename: /node_modules/@tanstack/react-store/package.json
{"name":"@tanstack/react-store","type":"module","types":"index.d.ts","exports":{".":{"import":{"types":"./index.d.ts"}}}}

// @filename: /node_modules/@tanstack/react-store/index.d.ts
export declare function useStore(...args: any[]): any; export declare class Store<T = any> { constructor(...args: any[]); state: T; }

// @filename: /repro.ts
import { createRootRoute, createRoute, createRouter, redirect } from "@tanstack/react-router";

const rootRoute = createRootRoute();
const agentRoute = createRoute({ getParentRoute: () => rootRoute, path: "/agents/$agentName" });
const agentIndexRoute = createRoute({ getParentRoute: () => agentRoute, path: "/" });
const sessionRoute = createRoute({ getParentRoute: () => agentRoute, path: "/$sessionId" });
const routeTree = rootRoute.addChildren([agentRoute.addChildren([agentIndexRoute, sessionRoute])]);
const router = createRouter({ routeTree });
declare module "@tanstack/react-router" {
    interface Register { router: typeof router; }
}

redirect({
    from: "/agents/$agentName/",
    to: "/agents/$agentName/$sessionId",
    params: { agentName: "a", sessionId: "s" },
});

agentIndexRoute.redirect({
    to: "/agents/$agentName/$sessionId",
    params: { agentName: "a", sessionId: "s" },
});
