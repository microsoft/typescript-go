// @skipLibCheck: true
// @noEmit: true
// @noTypesAndSymbols: true
// @filename: /node_modules/@tanstack/react-router/dist/esm/index.d.ts
import { BaseRootRoute, BaseRoute, RouteOptions } from "@tanstack/router-core";
export declare class Route<in out TRegister = unknown, in out TParentRoute extends RouteConstraints["TParentRoute"] = AnyRoute, in out TPath extends RouteConstraints["TPath"] = "/", in out TFullPath extends RouteConstraints["TFullPath"] = ResolveFullPath<TParentRoute, TPath>, in out TId extends RouteConstraints["TId"] = ResolveId<TParentRoute, TPath>, in out TLoaderDeps extends Record<string, any> = {}, in out THandlers = undefined> extends BaseRoute<TRegister, TParentRoute, TPath, TFullPath, THandlers> {
}
export declare function createRoute<TRegister = unknown, TPath extends RouteConstraints["TPath"] = "/", TFullPath extends RouteConstraints["TFullPath"] = ResolveFullPath<TParentRoute, TPath>, const TServerMiddlewares = unknown>(options: RouteOptions<TRegister, TServerMiddlewares>): Route<TRegister, TParentRoute, TPath, TServerMiddlewares>;
export declare class RootRoute<in out TRegister = unknown, in out TRouterContext = {}, in out TLoaderDeps extends Record<string, any> = {}, in out THandlers = undefined> extends BaseRootRoute<TRegister, THandlers> {
}
export declare function createRootRoute<TRegister = Register, TRouterContext = {}, THandlers = undefined>(options?: RootRouteOptions<TRegister, THandlers>): RootRoute<TRegister, THandlers>;
export declare const createRouter: CreateRouterFn;
// @filename: /node_modules/@tanstack/react-router/package.json
{"name":"@tanstack/react-router","type":"module","types":"dist/esm/index.d.ts","exports":{".":{"import":{"types":"./dist/esm/index.d.ts"}}}}
// @filename: /node_modules/@tanstack/router-core/dist/esm/fileRoute.d.ts
export interface FileRouteTypes {
}
// @filename: /node_modules/@tanstack/router-core/dist/esm/index.d.ts
export { BaseRoute, BaseRootRoute } from "./route.js";
// @filename: /node_modules/@tanstack/router-core/dist/esm/link.d.ts
export type NavigateOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = ".", TMaskFrom extends string = TFrom, TMaskTo extends string = "."> = ToOptions<TRouter, TFrom, TMaskTo>;
export type ToOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = ".", TMaskTo extends string = "."> = ToSubOptions<TRouter, TFrom, TTo>;
export type ToSubOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = "."> = ToSubOptionsProps<TRouter, TFrom, TTo> & PathParamOptions<TRouter, TFrom, TTo>;
export interface OptionalToOptions<in out TRouter extends AnyRouter, in out TFrom extends string, in out TTo extends string | undefined> {
}
export type MakeToRequired<TRouter extends AnyRouter, TFrom extends string, TTo extends string | undefined> = string extends TFrom ? string extends TTo ? OptionalToOptions<TRouter, TTo> : OptionalToOptions<TRouter, TTo> : OptionalToOptions<TRouter, TFrom, TTo>;
export type ToSubOptionsProps<TRouter extends AnyRouter = RegisteredRouter, TFrom extends RoutePaths<TRouter["routeTree"]> | string = string, TTo extends string | undefined = "."> = MakeToRequired<TRouter, TFrom, TTo>;
export interface MakeOptionalPathParams<in out TRouter extends AnyRouter, in out TFrom, in out TTo> {
    params?: true | (ParamsReducer<TRouter, "PATH", TFrom, TTo> & {});
}
export type IsRequired<TRouter extends AnyRouter, TParamVariant extends ParamVariant, TFrom, TTo> = ResolveRelativePath<TFrom, TTo> extends infer TPath ? undefined extends TPath ? never : TPath extends CatchAllPaths<TRouter> ? never : IsRequiredParams<ResolveRelativeToParams<TRouter, TTo>> : never;
export type PathParamOptions<TRouter extends AnyRouter, TFrom, TTo> = IsRequired<TRouter, "PATH", TFrom, TTo> extends never ? MakeOptionalPathParams<TRouter, TFrom, TTo> : MakeRequiredPathParams<TRouter, TTo>;
// @filename: /node_modules/@tanstack/router-core/dist/esm/redirect.d.ts
import { NavigateOptions } from "./link.js";
export type RedirectOptions<TRouter extends AnyRouter = RegisteredRouter, TFrom extends string = string, TTo extends string | undefined = undefined, TMaskFrom extends string = TFrom, TMaskTo extends string = "."> = {} & NavigateOptions<TRouter, TFrom, TTo, TMaskFrom, TMaskTo>;
export type RedirectOptionsRoute<TDefaultFrom extends string = string, TRouter extends AnyRouter = RegisteredRouter, TTo extends string | undefined = undefined, TMaskTo extends string = ""> = Omit<RedirectOptions<TRouter, TDefaultFrom, TTo, TDefaultFrom, TMaskTo>, "from">;
export interface RedirectFnRoute<in out TDefaultFrom extends string = string> {
    <TRouter extends AnyRouter = RegisteredRouter, const TTo extends string | undefined = undefined, const TMaskTo extends string = "">(opts: RedirectOptionsRoute<TDefaultFrom, TRouter, TTo, TMaskTo>): Redirect<TRouter, TMaskTo>;
}
// @filename: /node_modules/@tanstack/router-core/dist/esm/route.d.ts
export type InferAllParams<TRoute> = TRoute extends {
    };
} ? TAllParams : {};
export interface RouteTypes<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps, in out TLoaderFn, in out TChildren, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> {
}
export interface Route<in out TRegister, in out TParentRoute extends AnyRoute, in out TPath extends string, in out TFullPath extends string, in out TCustomId extends string, in out TId extends string, in out TSearchValidator, in out TParams, in out TRouterContext, in out TRouteContextFn, in out TBeforeLoadFn, in out TLoaderDeps extends Record<string, any>, in out TLoaderFn, in out TChildren, in out TFileRouteTypes, in out TSSR, in out TServerMiddlewares, in out THandlers> extends RouteExtensions<TId, TFullPath> {
}
export declare class BaseRoute<in out TRegister = Register, in out TPath extends string = "/", in out TFullPath extends string = ResolveFullPath<TParentRoute, TPath>, in out TId extends string = ResolveId<TParentRoute, TPath>, in out TParams = ResolveParams<TPath>, in out THandlers = undefined> {
    redirect: RedirectFnRoute<TFullPath>;
}
export declare class BaseRootRoute<in out TRegister = Register, in out THandlers = undefined> extends BaseRoute<TRegister, THandlers> {
}
export type ParseRoute<TRouteTree, TAcc = TRouteTree> = TRouteTree extends {
    types: {
    };
export type CodeRoutesByPath<TRouteTree extends AnyRoute> = ParseRoute<TRouteTree> extends infer TRoutes extends AnyRoute ? {
} : never;
export interface Register {
}
export type RegisteredRouter<TRegister = Register> = TRegister extends {
} ? TRouter : AnyRouter;
export interface RouterOptions<TRouteTree extends AnyRoute, TDehydrated = undefined> extends RouterOptionsExtensions {
}
export declare class RouterCore<in out TRouteTree extends AnyRoute, in out TDehydrated extends Record<string, any>> {
}
export type MakeDifferenceOptional<TLeft, TRight> = keyof TLeft & keyof TRight extends never ? TRight : Omit<TRight, keyof TRight> & {
// @filename: /node_modules/@tanstack/router-core/package.json
{"name":"@tanstack/router-core","type":"module","types":"dist/esm/index.d.ts","exports":{".":{"import":{"types":"./dist/esm/index.d.ts"}}}}
// @filename: /repro.ts
import { createRootRoute, createRoute, createRouter } from "@tanstack/react-router";
const agentRoute = createRoute({
});
const agentIndexRoute = createRoute({
});
const sessionRoute = createRoute({
});
const router = createRouter({
});
declare module "@tanstack/react-router" {
    interface Register {
    }
}
agentIndexRoute.redirect({
    params: {
    },
});
