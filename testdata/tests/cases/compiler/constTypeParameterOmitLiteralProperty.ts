// https://github.com/microsoft/typescript-go/issues/4391
// @moduleResolution: bundler
// @module: preserve
// @target: esnext
// @strict: true
// @skipLibCheck: true
// @noEmit: true
// @jsx: react-jsx
// @noTypesAndSymbols: true

// @filename: /node_modules/@tanstack/react-router/dist/esm/index.d.ts
export type { Register } from "@tanstack/router-core";
export { createRoute, createRootRoute, } from "./route.js";
export { createRouter } from "./router.js";

// @filename: /node_modules/@tanstack/react-router/dist/esm/route.d.ts
import { BaseRootRoute, BaseRoute, RouteOptions } from "@tanstack/router-core";
export declare class Route<in out TRegister = unknown, in out TParentRoute extends RouteConstraints["TParentRoute"] = AnyRoute, in out TPath extends RouteConstraints["TPath"] = "/", in out TFullPath extends RouteConstraints["TFullPath"] = ResolveFullPath<TParentRoute, TPath>, in out TId extends RouteConstraints["TId"] = ResolveId<TParentRoute, TPath>, in out TLoaderDeps extends Record<string, any> = {}, in out THandlers = undefined> extends BaseRoute<TRegister, TParentRoute, TPath, TFullPath, THandlers> {
}
export declare function createRoute<TRegister = unknown, TPath extends RouteConstraints["TPath"] = "/", TFullPath extends RouteConstraints["TFullPath"] = ResolveFullPath<TParentRoute, TPath>, const TServerMiddlewares = unknown>(options: RouteOptions<TRegister, TServerMiddlewares>): Route<TRegister, TParentRoute, TPath, TServerMiddlewares>;
export declare class RootRoute<in out TRegister = unknown, in out TRouterContext = {}, in out TLoaderDeps extends Record<string, any> = {}, in out THandlers = undefined> extends BaseRootRoute<TRegister, THandlers> {
}
export declare function createRootRoute<TRegister = Register, TRouterContext = {}, THandlers = undefined>(options?: RootRouteOptions<TRegister, THandlers>): RootRoute<TRegister, THandlers>;

// @filename: /node_modules/@tanstack/react-router/dist/esm/router.d.ts
import { CreateRouterFn } from "@tanstack/router-core";
export declare const createRouter: CreateRouterFn;

// @filename: /node_modules/@tanstack/react-router/package.json
{"name":"@tanstack/react-router","type":"module","types":"dist/esm/index.d.ts","exports":{".":{"import":{"types":"./dist/esm/index.d.ts"}}}}

// @filename: /node_modules/@tanstack/router-core/dist/esm/fileRoute.d.ts
export interface FileRouteTypes {
    to: any;
}
export type InferFileRouteTypes<TRouteTree extends AnyRoute> = unknown extends TRouteTree["types"]["fileRouteTypes"] ? never : TRouteTree["types"]["fileRouteTypes"] extends FileRouteTypes ? TRouteTree["types"]["fileRouteTypes"] : never;

// @filename: /node_modules/@tanstack/router-core/dist/esm/index.d.ts
export { BaseRoute, BaseRootRoute } from "./route.js";
export type { Register, CreateRouterFn, } from "./router.js";

// @filename: /node_modules/@tanstack/router-core/dist/esm/link.d.ts
import { RouteByPath, RouteToPath } from "./routeInfo.js";
import { ConstrainLiteral, Expand, MakeDifferenceOptional, Updater } from "./utils.js";
export type AddTrailingSlash<T> = T extends `${string}/` ? T : `${T & string}/`;
export type ResolveRelativePath<TFrom, TTo = "."> = string extends TFrom ? TTo : string extends TTo ? TFrom : undefined extends TTo ? TFrom : TTo extends string ? TFrom extends string ? TTo extends `/${string}` ? TTo : TTo extends `..${string}` ? ResolveParentPath<TFrom, TTo> : AddLeadingSlash<JoinPath<TFrom, TTo>> : never : never;
export type AbsoluteToPath<TRouter extends AnyRouter, TFrom extends string> = (string extends TFrom ? CurrentPath<TRouter> : TFrom extends `/` ? never : CurrentPath<TRouter>) | RouteToPath<TRouter>;
export type RelativeToPathAutoComplete<TRouter extends AnyRouter, TFrom extends string, TTo extends string> = string extends TTo ? string : string extends TFrom ? AbsoluteToPath<TRouter, TFrom> : TTo & `..${string}` extends never ? TTo & `.${string}` extends never ? AbsoluteToPath<TRouter, TFrom> : RelativeToCurrentPath<TRouter, TTo> : RelativeToParentPath<TRouter, TTo>;
export type NavigateOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = ".", TMaskFrom extends string = TFrom, TMaskTo extends string = "."> = ToOptions<TRouter, TFrom, TMaskTo>;
export type ToOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = ".", TMaskTo extends string = "."> = ToSubOptions<TRouter, TFrom, TTo>;
export type ToSubOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = "."> = ToSubOptionsProps<TRouter, TFrom, TTo> & PathParamOptions<TRouter, TFrom, TTo>;
export interface OptionalToOptions<in out TRouter extends AnyRouter, in out TFrom extends string, in out TTo extends string | undefined> {
    to?: ToPathOption<TRouter, TFrom, TTo>;
}
export type MakeToRequired<TRouter extends AnyRouter, TFrom extends string, TTo extends string | undefined> = string extends TFrom ? string extends TTo ? OptionalToOptions<TRouter, TTo> : OptionalToOptions<TRouter, TTo> : OptionalToOptions<TRouter, TFrom, TTo>;
export type ToSubOptionsProps<TRouter extends AnyRouter = RegisteredRouter, TFrom extends RoutePaths<TRouter["routeTree"]> | string = string, TTo extends string | undefined = "."> = MakeToRequired<TRouter, TFrom, TTo>;
export type ParamsReducerFn<in out TRouter extends AnyRouter, in out TParamVariant extends ParamVariant, in out TFrom, in out TTo> = (current: Expand<ResolveFromParams<TRouter, TFrom>>) => Expand<ResolveRelativeToParams<TRouter, TTo>>;
type ParamsReducer<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = Expand<ResolveRelativeToParams<TRouter, TParamVariant, TFrom, TTo>> | (ParamsReducerFn<TRouter, TParamVariant, TFrom, TTo> & {});
export type ResolveRoute<TRouter extends AnyRouter, TFrom, TPath = ResolveRelativePath<TFrom, TTo>> = TPath extends string ? TFrom extends TPath ? RouteByPath<TRouter["routeTree"], TPath> : RouteByToPath<TRouter, TPath> : never;
type ResolveToParamType<TParamVariant extends ParamVariant> = TParamVariant extends "PATH" ? "allParams" : "fullSearchSchemaInput";
export type ResolveToParams<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = ResolveRelativePath<TFrom, TTo> extends infer TPath ? undefined extends TPath ? never : string extends TPath ? ResolveAllToParams<TRouter, TParamVariant> : ResolveRoute<TRouter, TFrom, TTo>["types"][ResolveToParamType<TParamVariant>] : never;
type ResolveRelativeToParams<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo, TToParams = ResolveToParams<TRouter, TParamVariant, TFrom, TTo>> = TParamVariant extends "SEARCH" ? TToParams : string extends TFrom ? TToParams : MakeDifferenceOptional<ResolveFromParams<TRouter, TFrom>, TToParams>;
export interface MakeOptionalPathParams<in out TRouter extends AnyRouter, in out TFrom, in out TTo> {
    params?: true | (ParamsReducer<TRouter, "PATH", TFrom, TTo> & {});
}
export type IsRequired<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = ResolveRelativePath<TFrom, TTo> extends infer TPath ? undefined extends TPath ? never : TPath extends CatchAllPaths<TRouter> ? never : IsRequiredParams<ResolveRelativeToParams<TRouter, TTo>> : never;
export type PathParamOptions<TRouter extends AnyRouter, TFrom, TTo> = IsRequired<TRouter, "PATH", TFrom, TTo> extends never ? MakeOptionalPathParams<TRouter, TFrom, TTo> : MakeRequiredPathParams<TRouter, TTo>;
export type ToPathOption<TRouter extends AnyRouter = AnyRouter, TFrom extends string = string, TTo extends string | undefined = string> = ConstrainLiteral<TTo, RelativeToPathAutoComplete<TRouter, NoInfer<TFrom> extends string ? NoInfer<TFrom> : "", NoInfer<TTo> & string>>;

// @filename: /node_modules/@tanstack/router-core/dist/esm/redirect.d.ts
import { NavigateOptions } from "./link.js";
import { RegisteredRouter } from "./router.js";
export type RedirectOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = undefined, TMaskFrom extends string = TFrom, TMaskTo extends string = "."> = {} & NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>;
export type RedirectOptionsRoute<TDefaultFrom extends string = string, TRouter extends AnyRouter = RegisteredRouter, TTo extends string | undefined = undefined, TMaskTo extends string = ""> = Omit<RedirectOptions<TRouter, TDefaultFrom, TTo, TDefaultFrom, TMaskTo>, "from">;
export interface RedirectFnRoute<in out TDefaultFrom extends string = string> {
    <TRouter extends AnyRouter = RegisteredRouter, const TTo extends string | undefined = undefined, const TMaskTo extends string = "">(opts: RedirectOptionsRoute<TDefaultFrom, TRouter, TTo, TMaskTo>): Redirect<TRouter, TMaskTo>;
}

// @filename: /node_modules/@tanstack/router-core/dist/esm/route.d.ts
import { RedirectFnRoute } from "./redirect.js";
import { Assign, Constrain, NoInfer } from "./utils.js";
export type InferAllParams<TRoute> = TRoute extends {
    types: {
        allParams: infer TAllParams;
    };
} ? TAllParams : {};
export type ResolveAllParamsFromParent<TParentRoute extends AnyRoute, TParams> = Assign<InferAllParams<TParentRoute>, TParams>;
export interface RouteTypes<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps, in out TLoaderFn, in out TChildren, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> {
    allParams: ResolveAllParamsFromParent<TParentRoute, TParams>;
    children: TChildren;
}
export type RouteAddChildrenFn<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps extends Record<string, any>, in out TLoaderFn, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> = <const TNewChildren>(children: Constrain<TNewChildren, Record<string, AnyRoute>>) => Route<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TNewChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
export interface Route<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps extends Record<string, any>, in out TLoaderFn, in out TChildren, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> extends RouteExtensions<TId, TFullPath> {
    types: RouteTypes<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
}
export declare class BaseRoute<in out TRegister = Register, in out TPath extends string = "/", in out TFullPath extends string = ResolveFullPath<TParentRoute, TPath>, in out TId extends string = ResolveId<TParentRoute, TPath>, in out TParams = ResolveParams<TPath>, in out THandlers = undefined> {
    get fullPath(): TFullPath;
    types: RouteTypes<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TChildren, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    addChildren: RouteAddChildrenFn<TRegister, TParentRoute, TPath, TFullPath, TCustomId, TId, TSearchValidator, TParams, TRouterContext, TRouteContextFn, TBeforeLoadFn, TLoaderDeps, TLoaderFn, TFileRouteTypes, TSSR, TServerMiddlewares, THandlers>;
    redirect: RedirectFnRoute<TFullPath>;
}
export declare class BaseRootRoute<in out TRegister = Register, in out THandlers = undefined> extends BaseRoute<TRegister, THandlers> {
}

// @filename: /node_modules/@tanstack/router-core/dist/esm/routeInfo.d.ts
import { InferFileRouteTypes } from "./fileRoute.js";
import { AddTrailingSlash, RemoveTrailingSlashes } from "./link.js";
export type ParseRoute<TRouteTree, TAcc = TRouteTree> = TRouteTree extends {
    types: {
        children: infer TChildren;
    };
} ? unknown extends TChildren ? TAcc : TChildren extends ReadonlyArray<any> ? ParseRoute<TChildren[number], TChildren[number]> : ParseRoute<TChildren[keyof TChildren], TChildren[keyof TChildren]> : TAcc;
export type CodeRoutesByPath<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? {
    [K in TRoutes as K["fullPath"]]: K;
} : never;
export type RoutesByPath<TRouteTree extends AnyRoute> = InferFileRouteTypes<TRouteTree> extends never ? CodeRoutesByPath<TRouteTree> : InferFileRouteTypes<TRouteTree>["fileRoutesByFullPath"];
export type RouteByPath<TRouteTree extends AnyRoute, TPath> = Extract<RoutesByPath<TRouteTree>[TPath & keyof RoutesByPath<TRouteTree>], AnyRoute>;
export type TrailingSlashOptionByRouter<TRouter extends AnyRouter> = TrailingSlashOption extends TRouter["options"]["trailingSlash"] ? "never" : NonNullable<TRouter["options"]["trailingSlash"]>;
export type FileRouteToPath<TRouter extends AnyRouter, TTo = InferFileRouteTypes<TRouter["routeTree"]>["to"], TTrailingSlashOption = TrailingSlashOptionByRouter<TRouter>> = "never" extends TTrailingSlashOption ? TTo : "always" extends TTrailingSlashOption ? AddTrailingSlash<TTo> : TTo;
export type RouteToPath<TRouter extends AnyRouter> = unknown extends TRouter ? string : InferFileRouteTypes<TRouter["routeTree"]> extends never ? CodeRouteToPath<TRouter> : FileRouteToPath<TRouter>;

// @filename: /node_modules/@tanstack/router-core/dist/esm/router.d.ts
export interface Register {
}
export type RegisteredRouter<TRegister = Register> = TRegister extends {
    router: infer TRouter;
} ? TRouter : AnyRouter;
export interface RouterOptions<TRouteTree extends AnyRoute, TDehydrated = undefined> extends RouterOptionsExtensions {
    routeTree?: TRouteTree;
}
export type RouterConstructorOptions<TRouteTree extends AnyRoute, TTrailingSlashOption extends TrailingSlashOption, TDefaultStructuralSharingOption extends boolean, TRouterHistory extends RouterHistory, TDehydrated extends Record<string, any>> = Omit<RouterOptions<TRouteTree, TDehydrated>, "defaultSsr">;
export type CreateRouterFn = <TRouteTree extends AnyRoute, TDehydrated extends Record<string, any>>(options: undefined extends number ? "strictNullChecks must be enabled in tsconfig.json" : RouterConstructorOptions<TRouteTree, TTrailingSlashOption, TDefaultStructuralSharingOption, TRouterHistory, TDehydrated>) => RouterCore<TRouteTree, TDehydrated>;
export declare class RouterCore<in out TRouteTree extends AnyRoute, in out TDehydrated extends Record<string, any>> {
    routeTree: TRouteTree;
}

// @filename: /node_modules/@tanstack/router-core/dist/esm/utils.d.ts
export type Expand<T> = T extends object ? T extends infer O ? O extends Function ? O : {} : never : T;
export type MakeDifferenceOptional<TLeft, TRight> = keyof TLeft & keyof TRight extends never ? TRight : Omit<TRight, keyof TRight> & {
};
export type IsNonEmptyObject<T> = T extends object ? keyof T extends never ? false : true : false;
export type Assign<TLeft, TRight> = TLeft extends any ? TRight extends any ? IsNonEmptyObject<TLeft> extends false ? TRight : IsNonEmptyObject<TRight> extends false ? TLeft : keyof TLeft & keyof TRight extends never ? TLeft & TRight : Omit<TLeft, keyof TRight> & TRight : never : never;
export type Constrain<T, TDefault = TConstraint> = (T extends TConstraint ? T : never);
export type ConstrainLiteral<T, TConstraint, TDefault = TConstraint> = (T & TConstraint);

// @filename: /node_modules/@tanstack/router-core/package.json
{"name":"@tanstack/router-core","type":"module","types":"dist/esm/index.d.ts","exports":{".":{"import":{"types":"./dist/esm/index.d.ts"}}}}

// @filename: /repro.ts
import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router";
const rootRoute = createRootRoute();
const agentRoute = createRoute({
});
const agentIndexRoute = createRoute({
});
const sessionRoute = createRoute({
});
const routeTree = rootRoute.addChildren([agentRoute.addChildren([agentIndexRoute, sessionRoute])]);
const router = createRouter({
    routeTree
});
declare module "@tanstack/react-router" {
    interface Register {
        router: typeof router;
    }
}
agentIndexRoute.redirect({
    to: "/agents/$agentName/$sessionId",
    params: {
        sessionId: "s"
    },
});
