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


// @filename: /node_modules/@tanstack/router-core/dist/esm/root.d.ts
/** Stable identifier used for the root route in a route tree. */
export declare const rootRouteId = "__root__";
export type RootRouteId = typeof rootRouteId;


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
