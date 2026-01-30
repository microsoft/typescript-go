/**
 * Base interfaces and types for the TypeScript API.
 *
 * This module provides the foundation for both sync and async implementations
 * of the TypeScript API client.
 */

export { type API, type APIOptions, type Project, type Symbol, SymbolFlags, type Type, TypeFlags } from "./api.ts";
export { type Identifiable, type ObjectFactories, ObjectRegistry, type ProjectLike, type ReleaseFunction } from "./objectRegistry.ts";
export type { MaybeAsync, MaybeAsyncArray } from "./types.ts";
