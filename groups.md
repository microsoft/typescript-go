# JS Baseline Diff Categorization

354 `.js.diff` files analyzed under `testdata/baselines/reference/submodule/`.

---

# Declaration emit now produces .d.ts output that was previously absent
## https://github.com/microsoft/typescript-go/issues/3094
## Corsa emits declaration files for tests that previously had no .d.ts baseline output
compiler/arrayFakeFlatNoCrashInferenceDeclarations.js.diff
compiler/declarationEmitCommonJsModuleReferencedType.js.diff
compiler/declarationEmitComputedPropertyNameSymbol1.js.diff
compiler/declarationEmitComputedPropertyNameSymbol2.js.diff
compiler/declarationEmitIsolatedDeclarationErrorNotEmittedForNonEmittedFile.js.diff
compiler/declarationEmitMappedTypeTemplateTypeofSymbol.js.diff
compiler/declarationEmitMixinPrivateProtected.js.diff
compiler/declarationEmitObjectAssignedDefaultExport.js.diff
compiler/declarationEmitPrivatePromiseLikeInterface.js.diff
compiler/declarationEmitReadonlyComputedProperty.js.diff
compiler/declarationEmitReexportedSymlinkReference3.js.diff
compiler/declarationEmitUnsafeImportSymbolName.js.diff
compiler/declarationEmitUsingTypeAlias1.js.diff
compiler/declarationEmitVarInElidedBlock.js.diff
compiler/emitClassExpressionInDeclarationFile2.js.diff
compiler/globalThisDeclarationEmit.js.diff
compiler/hugeDeclarationOutputGetsTruncatedWithError.js.diff
compiler/isolatedDeclarationErrors.js.diff
compiler/isolatedDeclarationErrorsAugmentation.js.diff
compiler/isolatedDeclarationErrorsClasses.js.diff
compiler/isolatedDeclarationErrorsClassesExpressions.js.diff
compiler/isolatedDeclarationErrorsDefault.js.diff
compiler/isolatedDeclarationErrorsEnums.js.diff
compiler/isolatedDeclarationErrorsExpandoFunctions.js.diff
compiler/isolatedDeclarationErrorsExpressions.js.diff
compiler/isolatedDeclarationErrorsFunctionDeclarations.js.diff
compiler/isolatedDeclarationErrorsObjects.js.diff
compiler/isolatedDeclarationErrorsReturnTypes.js.diff
compiler/isolatedDeclarationLazySymbols.js.diff
compiler/isolatedDeclarationsAddUndefined.js.diff
compiler/isolatedDeclarationsAllowJs.js.diff
compiler/moduleAugmentationInAmbientModule2.js.diff
compiler/moduleAugmentationInAmbientModule3.js.diff
compiler/moduleAugmentationInAmbientModule4.js.diff
compiler/privateFieldsInClassExpressionDeclaration.js.diff
conformance/declarationFiles.js.diff
conformance/jsDeclarationsTypeReassignmentFromDeclaration2.js.diff
conformance/jsDeclarationsExportAssignedConstructorFunctionWithSub(target=es2015).js.diff
conformance/legacyNodeModulesExportsSpecifierGenerationConditions.js.diff
conformance/nodeModulesExportsBlocksSpecifierResolution(module=node16).js.diff
conformance/nodeModulesExportsBlocksSpecifierResolution(module=node18).js.diff
conformance/nodeModulesExportsBlocksSpecifierResolution(module=node20).js.diff
conformance/nodeModulesExportsBlocksSpecifierResolution(module=nodenext).js.diff
conformance/nodeModulesExportsSourceTs(module=node16).js.diff
conformance/nodeModulesExportsSourceTs(module=node18).js.diff
conformance/nodeModulesExportsSourceTs(module=node20).js.diff
conformance/nodeModulesExportsSourceTs(module=nodenext).js.diff






# Comma expressions in emitted JS are now parenthesized for correctness
## TODO: Accept this (parenthesizing comma expressions prevents precedence bugs)
## `a, b` in computed properties and JSX context now emitted as `(a, b)`
conformance/jsxCheckJsxNoTypeArgumentsAllowed.js.diff
conformance/jsxInvalidEsprimaTestSuite.js.diff
conformance/jsxParsingError1.js.diff
conformance/parserComputedPropertyName35.js.diff
conformance/tsxErrorRecovery2.js.diff

# Unnecessary parentheses removed from emitted JS and declarations
## TODO: Accept this (removing redundant parentheses is a safe simplification)
## Redundant parenthesization stripped from instantiation expressions, optional chain emit, and declaration types
compiler/declarationEmitPromise.js.diff
compiler/declarationEmitUsingTypeAlias2.js.diff
compiler/optionalChainWithInstantiationExpression2(target=es2019).js.diff
conformance/importWithTypeArguments.js.diff
conformance/instantiationExpressionErrors.js.diff
conformance/instantiationExpressions.js.diff

# Parenthesization added/changed in declaration emit for complex types
## TODO: Accept this (added parentheses are required for correct parsing of these types)
## `keyof infer U extends number` → `keyof (infer U extends number)`, `T | T & undefined` → `T | (T & undefined)`, `readonly readonly string[]` → `readonly (readonly string[])`
conformance/inferTypesWithExtends1.js.diff
conformance/spreadObjectOrFalsy.js.diff
conformance/readonlyArraysAndTuples.js.diff

# Getter/setter member ordering changed (getter emitted before setter)
## TODO: Accept this (cosmetic ordering; consistent getter-before-setter is reasonable)
## In declaration emit, getters now consistently appear before setters
compiler/lateBoundAssignmentCandidateJS1.js.diff
conformance/jsDeclarationsGetterSetter.js.diff
conformance/thisPropertyAssignmentInherited.js.diff
conformance/jsDeclarationsReusesExistingTypeAnnotations.js.diff

# Object literal accessors now emit separate `get`/`set` declarations instead of collapsing to a property
## TODO: Accept this (preserving accessor form is more faithful and retains get/set type asymmetry)
## Declaration emit preserves accessor form rather than merging to a single property type
compiler/declarationEmitObjectLiteralAccessors1.js.diff
compiler/declarationEmitObjectLiteralAccessorsJs1.js.diff

# Namespace/expando declarations restructured (namespaces → `declare const` objects)
## TODO: Accept this (Documented in CHANGES.md in section "Expando declarations")
## Const objects and expandos no longer emit as namespaces; instead emit as `declare const` with object literal types
compiler/jsDeclarationsWithDefaultAsNamespaceLikeMerge.js.diff
conformance/expandoOnAlias.js.diff
conformance/jsDeclarationsConstsAsNamespacesWithReferences.js.diff
conformance/jsDeclarationsImportAliasExposedWithinNamespace.js.diff
conformance/jsDeclarationsTypeReferences3(target=es2015).js.diff
conformance/jsdocTemplateTagNameResolution.js.diff
conformance/requireOfESWithPropertyAccess.js.diff
conformance/nestedDestructuringOfRequire.js.diff

# Expando function namespace property declarations each get their own `declare namespace` block
## TODO: Accept this (separate namespace blocks merge to the same result in TypeScript)
## Each property assignment on a function gets a separate namespace block instead of one merged block
conformance/jsDeclarationsFunctionKeywordProp(target=es2015).js.diff
conformance/jsDeclarationsFunctionKeywordPropExhaustive(target=es2015).js.diff
conformance/nullPropertyName.js.diff
conformance/typeFromPropertyAssignment29.js.diff

# Expando function namespace removed entirely (nested assignments no longer produce declarations)
## TODO: Accept this (Documented in CHANGES.md in section "Expando declarations")
## Expando namespace with nested property assignments is completely removed from declaration emit
compiler/expandoFunctionNestedAssigments.js.diff

# Namespace-to-const restructuring with DtsFileErrors (redeclaration errors)
## TODO: Log a bug (produces invalid .d.ts with TS2451 redeclaration errors)
## Namespace properties restructured to `declare const` patterns, sometimes causing redeclaration errors
conformance/typeFromPropertyAssignment39.js.diff

# Namespace converted to type alias in declaration emit
## TODO: Accept this (Documented in CHANGES.md in section "JSDoc Tags and Types")
## `namespace Dotted { type Name = number }` → `type Dotted = number`
conformance/jsDeclarationsImportNamespacedType.js.diff

# `@enum` tag no longer produces namespace+type pairs; emits as `declare const` with JSDoc preserved
## TODO: Accept this (Documented in CHANGES.md in section "JSDoc Tags and Types")
## @enum types now emit as `declare const` object types instead of `type + namespace` declarations
conformance/jsDeclarationsEnumTag(target=es2015).js.diff

# Variables instantiated from classes with private/protected constructors now typed correctly instead of `any`
## TODO: Accept this (type accuracy improvement; `any` was incorrect)
## Type inference improvement: `new C()` inside the class now returns `C` type, not `any`
conformance/classConstructorAccessibility.js.diff
conformance/classConstructorAccessibility2.js.diff
conformance/typesWithPrivateConstructor.js.diff
conformance/typesWithProtectedConstructor.js.diff

# Invalid type annotations, modifiers, and type parameters now properly stripped in JS output
## TODO: Accept this (invalid syntax should not appear in JS output; stripping it is more correct)
## Constructor return types, getter/setter type params, and export modifiers on constructors are now correctly removed
conformance/parserConstructorDeclaration10.js.diff
conformance/parserConstructorDeclaration3.js.diff
conformance/parserConstructorDeclaration9.js.diff
conformance/parserGetAccessorWithTypeParameters1(target=es2015).js.diff
conformance/parserSetAccessorWithTypeAnnotation1(target=es2015).js.diff
conformance/parserSetAccessorWithTypeParameters1(target=es2015).js.diff
conformance/typeGuardFunctionErrors.js.diff

# Decorator metadata reflects stable union type ordering
## TODO: Log a bug (decorator metadata is runtime-accessible; changing `Object` → `String` affects `Reflect.getMetadata`)
## `design:type` for `T | null` changed from `Object` to the non-null type (e.g. `String`, `Class1`)
compiler/metadataOfUnionWithNull(strictnullchecks=true).js.diff
compiler/metadataReferencedWithinFilteredUnion(strictnullchecks=true,target=es2015).js.diff

# void-zero (`void 0`) expando assignments no longer ignored
## TODO: Accept this (Documented in CHANGES.md in section "Expandos")
## `x = void 0` now creates a property with type `undefined` instead of being silently ignored
conformance/assignmentToVoidZero1.js.diff
conformance/assignmentToVoidZero2.js.diff

# JSDoc type resolution changes: `Object<K,V>` → `Record`, `?` → `any|null`, `function()` → `Function`
## TODO: Accept this (Documented in CHANGES.md in section "JSDoc Tags and Types")
## JSDoc-specific type syntax now resolves differently in declaration emit
conformance/jsDeclarationsJSDocRedirectedLookups.js.diff
conformance/jsDeclarationsReusesExistingNodesMappingJSDocTypes.js.diff
conformance/jsDeclarationsRestArgsWithThisTypeInJSDocFunction(target=es2015).js.diff
conformance/jsDeclarationsMissingTypeParameters(target=es2015).js.diff

# Missing generic type arguments no longer auto-filled with `any` (produces TS2314 errors)
## TODO: Log a bug (?)
## `Array` without type arguments now emits bare `Array` instead of `Array<any>`
conformance/jsDeclarationsMissingGenerics(target=es2015).js.diff

# Optional property types in JS declarations drop redundant `| undefined`
## TODO: Log a bug
## `b?: number | undefined` simplified to `b?: number`
conformance/jsDeclarationsOptionalTypeLiteralProps1.js.diff
conformance/jsDeclarationsOptionalTypeLiteralProps2.js.diff

# Nullable rest tuple member syntax changed from `...?T` to `...T | null`
## TODO: Accept this (`...T | null` is standard syntax; `...?T` was a non-standard printer shorthand)
## Named tuple member `...?string[]` now prints as `...string[] | null`
conformance/namedTupleMembersErrors.js.diff
conformance/restTupleElements1.js.diff

# Re-export statements split from single combined line into separate per-specifier lines
## TODO: Accept this (split re-exports are semantically equivalent; same names exported from same source)
## `export { default, foo, bar } from "x"` → three separate `export { ... } from "x"` statements
conformance/nodeModulesAllowJsImportHelpersCollisions3(module=node16,target=es2015).js.diff
conformance/nodeModulesAllowJsImportHelpersCollisions3(module=node18,target=es2015).js.diff
conformance/nodeModulesAllowJsImportHelpersCollisions3(module=node20,target=es2015).js.diff
conformance/nodeModulesAllowJsImportHelpersCollisions3(module=nodenext,target=es2015).js.diff

# Output file section ordering changed (deterministic but different from Strada)
## TODO: Accept this (both orderings are deterministic; the difference is cosmetic)
## JS and .d.ts output file sections appear in a different deterministic order
conformance/nodeModulesAllowJsPackageImports(module=node16).js.diff
conformance/nodeModulesAllowJsPackageImports(module=node18).js.diff
conformance/nodeModulesAllowJsPackageImports(module=node20).js.diff
conformance/nodeModulesAllowJsPackageImports(module=nodenext).js.diff

# Output file ordering changed + declaration types resolve to `any` instead of `typeof` references
## TODO: Log a bug
## File sections reordered AND package export types resolve differently in declaration emit
conformance/nodeModulesDeclarationEmitWithPackageExports(module=node16).js.diff
conformance/nodeModulesDeclarationEmitWithPackageExports(module=node18).js.diff
conformance/nodeModulesDeclarationEmitWithPackageExports(module=node20).js.diff
conformance/nodeModulesDeclarationEmitWithPackageExports(module=nodenext).js.diff

# `export = a` reordered after declaration; ESM side-effect import `import "fs"` replaces `export {}`
## TODO: Log a bug
## CJS declaration ordering + ESM side-effect emit changed
conformance/nodeModulesAllowJsExportAssignment(module=node16).js.diff
conformance/nodeModulesAllowJsExportAssignment(module=node18).js.diff
conformance/nodeModulesAllowJsExportAssignment(module=node20).js.diff
conformance/nodeModulesAllowJsExportAssignment(module=nodenext).js.diff

# BOM character handling fixed
## TODO: Accept this (bug fix; Strada was incorrectly mangling the BOM)
## BOM is now correctly preserved in output instead of being mangled, and spurious TS1127 errors are removed
compiler/emitBOM.js.diff

# `new new Date` now emits `new (new Date)` with explicit parentheses
## TODO: Accept this (explicit parentheses ensure correct `new` expression precedence)
## Parenthesization added for `new` expression precedence correctness
compiler/newOperator.js.diff

# Import alias `import I = M` now additionally emits `var I = M` in JS output
## TODO: Accept this (emitting the variable assignment is more correct for runtime usage)
## Namespace import aliases in wrong context now emit JS variable assignment
compiler/moduleElementsInWrongContext.js.diff
compiler/moduleElementsInWrongContext2.js.diff

# Duplicate variable declarations removed from JS emit
## TODO: Accept this (redundant declarations are dead code; removing them is a safe cleanup)
## Redundant `var console;` or `let x;` declarations no longer emitted when namespace and variable share a name
compiler/module_augmentExistingVariable.js.diff
compiler/nameCollisions.js.diff

# Empty namespace IIFE no longer emitted
## TODO: Accept this (empty IIFEs are dead code with no runtime effect)
## `namespace global {}` inside a class no longer generates an empty IIFE in JS output
compiler/nestedGlobalNamespaceInClass.js.diff

# JSX factory call changed from indirect `(0, _a.jsx)(...)` to direct `_jsx(...)` call
## TODO: Log a bug
## JSX emit no longer uses comma-expression indirection for factory calls
compiler/commentsOnJSXExpressionsArePreserved(jsx=react-jsx,module=commonjs,moduledetection=legacy).js.diff

# Inline JSDoc parameter comments no longer emitted in function type declarations
## TODO: Accept this (comments are cosmetic; type signatures are identical)
## Comments like `/** fooFunctionValue param */` stripped from declaration output
compiler/commentsFunction(target=es2015).js.diff

# JSDoc `@typedef` descriptions preserved inline instead of being converted to per-property comments
## TODO: Accept this (comment format change only; the type structure is identical)
## Full `@typedef`/`@property` JSDoc blocks are now preserved with the type alias
conformance/jsDeclarationsTypedefDescriptionsPreserved.js.diff

# `declare` modifier on import correctly elides the JS emit
## TODO: Accept this (`declare` imports are type-only and should not produce runtime code)
## `import = b` with `declare` modifier no longer emits `var a = b` in JS
## Input was declare import a = b;
compiler/declareModifierOnImport1.js.diff
## Input was declare export import a = x.c;
compiler/importDeclWithDeclareModifier.js.diff

# `import {} from "./a"` added for type-only import side effects
## TODO: Accept this (empty import has no effect on module exports or type signatures)
## Type-only imports with side effects now emit an empty import statement
conformance/bundlerSyntaxRestrictions(module=esnext).js.diff
conformance/bundlerSyntaxRestrictions(module=preserve).js.diff

# Unused `const _super = Object.create(null, {})` removed from async generator emit
## TODO: Accept this (dead code elimination; the `_super` binding was unused)
## Dead code in async methods with super access eliminated
conformance/asyncMethodWithSuper_es6.js.diff

# `accessor` keyword erroneously added to constructor in error recovery emit
## TODO: Accept this (input is invalid; both outputs are wrong, so the difference is immaterial)
## Invalid `accessor` modifier leaks onto `constructor()` in JS output
conformance/autoAccessorDisallowedModifiers(target=esnext).js.diff

# `null`-initialized variable inferred as `any` instead of `null` type
## TODO: Log a bug
## `var l11 = null` now has type `any` instead of `null`
compiler/letDeclarations.js.diff

# Mapped type with missing type now emits `any` instead of empty/missing in declaration
## TODO: Accept this (`any` is a safer fallback for incomplete types than empty/missing)
## Error recovery for incomplete mapped types improved
compiler/mappedTypeNoTypeNoCrash.js.diff

# @overload handling changes in JS declaration emit
## TODO: Log a bug (overload signatures change or break; some lose type parameters producing TS2304 errors)
## JSDoc @overload functions emit differently: JSDoc comments removed/repeated, some overloads collapsed, function → const arrow
compiler/jsFileAlternativeUseOfOverloadTag.js.diff
compiler/jsFileFunctionOverloads.js.diff
compiler/jsFileFunctionOverloads2.js.diff
compiler/jsFileMethodOverloads.js.diff
compiler/jsFileMethodOverloads2.js.diff
compiler/jsFileMethodOverloads5.js.diff
conformance/overloadTag1.js.diff
conformance/overloadTag2.js.diff

# @callback/@overload tags reordered and restructured in declaration emit
## TODO: Log a bug (`(...args: string[])` becomes `(...args: string)`, losing the array type on rest params)
## @callback/@overload type definitions reorganized (placement, optional property simplification, function→const)
conformance/callbackTagNestedParameter.js.diff
conformance/callbackTagVariadicType.js.diff

# `@satisfies` tag handling changes: functions become const declarations with enriched JSDoc
## TODO: Log a bug (fn4–fn6 parameter types change, e.g. `b: number` → `b: never`; visible type differences)
## Functions with @satisfies become `declare const` with arrow types and added @param/@satisfies annotations
conformance/checkJsdocSatisfiesTag15.js.diff

# @callback typedef now emitted as exported type in JS declaration emit
## @callback constructs on constructors restructured with `export declare class` and explicit constructor
## TODO: Accept this (same type structure; callback typedef and class exports are equivalent)
## Documented in CHANGES.md in section "JSDoc Tags"
conformance/callbackOnConstructor.js.diff

# Labeled statement followed by export declaration now preserves original text
## TODO: Log a bug
## `export const title` in labeled statement context no longer emits as bare semicolon
conformance/labeledStatementExportDeclarationNoCrash1(module=commonjs).js.diff

# `export {}` statement moved from end to beginning of JS output file
## TODO: Accept this (position of the module marker is cosmetic; no runtime effect)
## Side-effect-only module marker repositioned
compiler/thisInObjectJs.js.diff

# CJS exports pattern changes in declaration emit
## TODO: Log a bug (some files lose exports entirely, e.g. `jsDeclarationsExportDoubleAssignmentInClosure` emits only `export {}`)
## Various CJS-specific declaration emit restructuring
compiler/jsExportAssignmentNonMutableLocation.js.diff
conformance/commonJSImportClassTypeReference.js.diff
conformance/commonJSImportNestedClassTypeReference.js.diff
conformance/jsDeclarationsExportDoubleAssignmentInClosure.js.diff

## TODO: Log a bug (Jake has a fix)
conformance/exportNonInitializedVariablesInIfThenStatementNoCrash1(module=commonjs).js.diff

# Weird testcases that don't correspond to anything real or are no-ops
## TODO: Accept this
compiler/augmentExportEquals2.js.diff
compiler/importDeclWithExportModifierAndExportAssignment.js.diff
conformance/defaultExportsCannotMerge04(target=es2015).js.diff

# Numeric string property key uses computed property syntax `["404"]` instead of plain `"404"`
## TODO: Accept this (`"404"` and `["404"]` declare the same property; consumers access it identically)
## Declaration emit changes how numeric-string-keyed properties are represented
compiler/declarationEmitPropertyNumericStringKey.js.diff

# Non-literal computed property `[fieldName]` changed to index signature `[x: string]` in declaration emit
## TODO: Accept this (fieldName is not a legal const to use in a computed position, so this is an error case)
## Dynamic computed property keys resolved to index signatures
compiler/declarationEmitSimpleComputedNames1.js.diff

# `typeof` references resolved to concrete types in declaration emit
## TODO: Accept this (concrete types are equivalent to what `typeof` resolves to)
## `typeof a` rest parameters and symbol references resolved to tuple/symbol types instead of typeof
compiler/declarationEmitTypeofRest.js.diff

# Namespace-qualified names simplified to unqualified where context allows
## TODO: Accept this (unqualified names are correct when already in scope; simplification reduces noise)
## `me.Things<me.Props>` → `Things<Props>` when namespace context makes qualification unnecessary
compiler/defaultDeclarationEmitShadowedNamedCorrectly.js.diff

# Declaration emit uses `import("...")` type instead of namespace-qualified references
## TODO: Accept this (`import("./file").B` resolves to the same type as the namespace-qualified `A.B`)
## `A.B` changed to `import("./file").B` in some declaration contexts
compiler/es5ExportEqualsDts(target=es2015).js.diff
compiler/declarationEmitNameConflicts.js.diff

# Properties with `__proto__`-like names now quoted in declaration emit
## TODO: Accept this (it's quoted in the originating sourcefile)
compiler/escapedReservedCompilerNamedIdentifier.js.diff

# Nested `/*elided*/` type annotation comments removed at deeper nesting levels
## TODO: Accept this (`/*elided*/` comments are internal artifacts with no semantic meaning)
## Declaration emit no longer preserves elided comments in nested class expressions
compiler/emitClassExpressionInDeclarationFile.js.diff

# `noImplicitThis` functions returning `this` now fully expand the object type in declarations
## TODO: Log a bug (return type changes from `any` to expanded object type; consumers see a different type contract)
## Functions returning `this` now emit the full resolved type instead of truncated `/*elided*/ any`
compiler/noImplicitThisBigThis.js.diff

# JSDoc `@template` tag default now emits `= ` (empty) instead of `= any` for invalid defaults
## TODO: Accept this (both are invalid defaults; neither is correct, so the difference is immaterial)
## Invalid @template defaults handled differently
conformance/jsdocTemplateTagDefault.js.diff

# @template-tagged function now emits as `const` arrow type instead of `function` declaration
## TODO: Accept this (`declare function f<T>()` and `declare const f: <T>() => R` are equivalent callable types)
## Function assigned to `const` with @template tag uses arrow function type in declaration
conformance/instantiateTemplateTagTypeParameterOnVariableStatement.js.diff

# JSDoc module namepath `module:A` now emits as `module` (unresolved) instead of `any`
## TODO: Log a bug (`any` → unresolved `module` identifier causes type errors for consumers)
## Module namepath JSDoc type syntax no longer resolves
conformance/jsDeclarationsModuleReferenceHasEmit(target=es2015).js.diff

# `definiteAssignment` shorthand in object literal loses optional modifier
## TODO: Log a bug (removing `?` makes a previously-optional member required; breaks consumers omitting it)
## `a?()` becomes `a()` in declaration emit for definite assignment assertions with object shorthand
conformance/definiteAssignmentAssertionsWithObjectShortHand.js.diff

# Generic function parameter inference change
## TODO: Log a bug
## `x3` declaration changes from `any[]` to `any[][]`
conformance/genericFunctionParameters.js.diff

# JS emit line-break reformatting: multi-line function call arguments collapsed to single line
## TODO: Accept this (whitespace-only change; no semantic difference)
## `(0, renderer_1.dom)(...)` factory call arguments joined from multiple lines into one line
conformance/inlineJsxFactoryDeclarationsLocalTypes.js.diff

# `@typedef` hoisted as direct type exports in JS declaration emit
## TODO: Accept this (same type exported; just hoisted to top level rather than nested)
## Typedef declarations are now top-level exports instead of nested in namespaces
conformance/typedefOnSemicolonClassElement.js.diff

# Declaration emit trailing comma changes in destructuring parameters
## TODO: Log a bug
## Extra trailing commas after omitted elements removed: `[, z, ,]` → `[, z, ]`
compiler/declarationEmitDestructuring5.js.diff

# Late-bound computed property `[prop]` preserved in declaration emit instead of resolved to plain `prop`
## TODO: Log a bug (`prop` vs `[prop]` are different access patterns; also property → method signature changes)
## Symbol-keyed computed properties keep their computed syntax
compiler/lateBoundAssignmentCandidateJS2.js.diff
compiler/lateBoundMethodNameAssigmentJS.js.diff

# Symbol getter/setter declaration emit change: new method overload added and getter return type changed
## TODO: Log a bug (getter return type changed from `I` to `undefined`; new method overload added — visible type changes)
## `[Symbol.toPrimitive]` gains a new `(x: I): void` method overload and getter return type changes from `I` to `undefined`
conformance/symbolDeclarationEmit12.js.diff

# Declaration emit restructuring within namespaces: `let` → `const`, exports separated, `declare` added
## TODO: Accept this (`let` vs `const` and `declare` keyword are cosmetic in .d.ts; same exports)
## Namespace member declarations restructured (let→const, export split out, `declare` keyword added)
conformance/jsDeclarationsEnums(target=es2015).js.diff
conformance/jsDeclarationsMultipleExportFromMerge(target=es2015).js.diff
conformance/jsDeclarationsTypeReferences4(target=es2015).js.diff

# `export as namespace` added or restructured in declaration emit
## TODO: Log a bug (new `export as namespace GLO` changes module's global visibility; adds a new public contract)
## Global namespace export declarations added or moved
conformance/jsDeclarationsExportFormsErr(target=es2015).js.diff

# Complex multi-faceted JS declaration emit restructuring (multiple overlapping changes)
## TODO: Accept this (overlapping cosmetic changes — `declare` keyword, ordering, `import(...)` form — same exported types)
## Files where multiple categories of change overlap; ordering + declare keyword + structure all changed
compiler/jsxDeclarationsWithEsModuleInteropNoCrash.js.diff
compiler/jsDocDeclarationEmitDoesNotUseNodeModulesPathWithoutError.js.diff
compiler/reuseTypeAnnotationImportTypeInGlobalThisTypeArgument.js.diff
compiler/declarationEmitResolveTypesIfNotReusable.js.diff
conformance/jsDeclarationsClasses(target=es2015).js.diff
conformance/jsDeclarationsClassExtendsVisibility(target=es2015).js.diff
conformance/jsDeclarationsClassImplementsGenericsSerialization(target=es2015).js.diff
conformance/jsDeclarationsClassStatic(target=es2015).js.diff
conformance/jsDeclarationsClassStaticMethodAugmentation(target=es2015).js.diff
conformance/jsDeclarationsCrossfileMerge(target=es2015).js.diff
conformance/jsDeclarationsDefaultsErr(target=es2015).js.diff
conformance/jsDeclarationsFunctions.js.diff
conformance/jsDeclarationsFunctionWithDefaultAssignedMember.js.diff
conformance/jsDeclarationsFunctionJSDoc(target=es2015).js.diff
conformance/jsDeclarationsClassAccessor.js.diff
conformance/jsDeclarationsExportAssignedClassInstance3(target=es2015).js.diff
conformance/jsDeclarationsInterfaces(target=es2015).js.diff
conformance/jsDeclarationsPackageJson(target=es2015).js.diff
conformance/jsDeclarationsReactComponents(target=es2015).js.diff
conformance/jsDeclarationsNonIdentifierInferredNames.js.diff
conformance/jsDeclarationsTypedefFunction.js.diff
conformance/jsDeclarationsUniqueSymbolUsage.js.diff
conformance/jsDeclarationsParameterTagReusesInputNodeInEmit1.js.diff
conformance/jsDeclarationsParameterTagReusesInputNodeInEmit2.js.diff
conformance/typeTagOnFunctionReferencesGeneric.js.diff
conformance/plainJSGrammarErrors.js.diff
conformance/jsDeclarationsExportForms(target=es2015).js.diff
conformance/commonJSAliasedExport.js.diff
conformance/jsDeclarationsFunctionsCjs.js.diff
conformance/jsDeclarationsExportDefinePropertyEmit.js.diff
conformance/jsDeclarationsExportSubAssignments(target=es2015).js.diff

# Declaration emit simplification of class extending built-in Array
## TODO: Log a bug
## `extends Array<any>` simplified to `extends Array`, inherited constructor overloads removed
compiler/javascriptThisAssignmentInStaticBlock.js.diff

# JS emit line-break reformatting in private name static method
## TODO: Accept this (whitespace/comment positioning change; no semantic difference)
## Comment `// Error` repositioned to next line in private field access method call
conformance/privateNameStaticMethod.js.diff

# Declaration emit type parameter renaming bug (regression)
## TODO: Log a bug
## Corsa incorrectly generates `T_1` where `T` is undefined, producing TS2304 DtsFileErrors
compiler/declarationEmitShadowing.js.diff

# JSDoc-in-TypeScript parsing changes cause broken emit
## TODO: Accept this (caused by documented deprecation of @function tag)
## JSDoc type annotations in `.ts` files now cause parsing issues resulting in broken JS output
conformance/jsdocDisallowedInTypescript.js.diff
