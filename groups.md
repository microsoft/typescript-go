# Symbols Diff Categorization

*473 `.symbols.diff` files under `testdata/baselines/reference/submodule/` categorized by type of binding change between Corsa (Go) and Strada (JS).*

Note: Many files exhibit multiple categories of change simultaneously. Such files appear under each relevant group.

---


# Inherited method resolves to base class instead of overriding derived class
## `C.foo` → `B.foo` for methods — symbol resolves to base class
## **Needs Investigation** — if the method is overridden in C, it should resolve to C's declaration, not B's.
conformance/typeFromPropertyAssignment23.symbols.diff

# Inherited property resolves to derived class instead of declaring base class
## `Element.textContent` → `TextElement.textContent` — property symbol attributed to derived class
## **Needs Investigation** — property should resolve to where it's originally declared (the base), not the derived class that inherits it.
conformance/thisPropertyAssignmentInherited.symbols.diff

# Duplicate/overload declarations no longer merge into a single symbol with multiple Decls
## Separate overload or duplicate method signatures each get their own symbol
## **Needs Investigation** — overloads in TS should share a single symbol; splitting them could break declaration emit and signature resolution for overloaded functions.
conformance/callSignaturesWithParameterInitializers2.symbols.diff
conformance/stringLiteralTypesInImplementationSignatures2.symbols.diff
conformance/symbolDeclarationEmit12.symbols.diff
conformance/symbolProperty44.symbols.diff
compiler/methodSignatureHandledDeclarationKindForSymbol.symbols.diff

# `import("...").X` property access on import types no longer resolves to member symbol
## Property access expressions on `import()` type references lose their symbol resolution
## **Needs Investigation** — losing symbol resolution for `import("...").X` means go-to-definition and hover will fail on JSDoc `import()` type references, a common JS pattern.
conformance/jsdocImportTypeReferenceToCommonjsModule.symbols.diff
conformance/jsdocImportTypeReferenceToESModule.symbols.diff

# Package self-name property access no longer resolves
## Accessing exports through a package's own name (e.g., `me.thing`) loses symbol resolution
## **Needs Investigation** — package self-name imports are a supported Node.js feature; losing resolution here could break editor navigation for monorepo/self-referencing packages.
compiler/nodeNextPackageSelfNameWithOutDir.symbols.diff
compiler/nodeNextPackageSelfNameWithOutDirDeclDir.symbols.diff

# Index signature symbols now use bracket notation instead of dot notation
## `Parent.__index` → `Parent[__index]` in symbol display names. Applies to JSX IntrinsicElements, static index signatures, and other index-signature-bearing types.
## **Benign** — cosmetic display name change; bracket notation is arguably more correct for index signatures.
conformance/tsxEmit1.symbols.diff
conformance/tsxEmit2.symbols.diff
conformance/tsxReactEmit1.symbols.diff
conformance/tsxReactEmit2.symbols.diff
conformance/tsxReactEmit4.symbols.diff
conformance/tsxReactEmit5.symbols.diff
conformance/tsxReactEmit6.symbols.diff
conformance/tsxReactEmit7.symbols.diff
conformance/tsxReactEmitEntities.symbols.diff
conformance/tsxReactEmitWhitespace.symbols.diff
conformance/tsxReactEmitWhitespace2.symbols.diff
conformance/tsxSpreadChildren.symbols.diff
conformance/tsxSpreadChildrenInvalidType(jsx=react,target=es2015).symbols.diff
conformance/tsxSpreadInvalidType.symbols.diff
conformance/tsxElementResolution2.symbols.diff
conformance/tsxElementResolution3.symbols.diff
conformance/tsxFragmentPreserveEmit.symbols.diff
conformance/tsxFragmentReactEmit.symbols.diff
conformance/inlineJsxAndJsxFragPragma.symbols.diff
conformance/inlineJsxAndJsxFragPragmaOverridesCompilerOptions.symbols.diff
conformance/inlineJsxFactoryDeclarations.symbols.diff
conformance/inlineJsxFactoryDeclarationsLocalTypes.symbols.diff
conformance/inlineJsxFactoryLocalTypeGlobalFallback.symbols.diff
conformance/inlineJsxFactoryOverridesCompilerOption.symbols.diff
conformance/inlineJsxFactoryWithFragmentIsError.symbols.diff
conformance/jsxParsingError1.symbols.diff
conformance/jsxParsingError2.symbols.diff
conformance/jsxParsingError3.symbols.diff
conformance/jsxParsingError4(strict=false).symbols.diff
conformance/jsxParsingError4(strict=true).symbols.diff
conformance/staticIndexSignature1.symbols.diff
conformance/staticIndexSignature2.symbols.diff
conformance/staticIndexSignature4.symbols.diff
conformance/staticIndexSignature6.symbols.diff
conformance/noPropertyAccessFromIndexSignature1.symbols.diff
conformance/noUncheckedIndexedAccessDestructuring.symbols.diff
conformance/propertyAccessStringIndexSignature(noimplicitany=false).symbols.diff
conformance/propertyAccessStringIndexSignature(noimplicitany=true).symbols.diff
conformance/thisTypeInFunctions2.symbols.diff
conformance/typeGuardOfFromPropNameInUnionType.symbols.diff
compiler/jsxChildrenGenericContextualTypes.symbols.diff
compiler/jsxIntrinsicDeclaredUsingTemplateLiteralTypeSignatures.symbols.diff
compiler/jsxNamespaceReexports.symbols.diff
compiler/propertyAccessOfReadonlyIndexSignature.symbols.diff
compiler/reactJsxReactResolvedNodeNext.symbols.diff
compiler/reactJsxReactResolvedNodeNextEsm.symbols.diff
compiler/jsxLocalNamespaceIndexSignatureNoCrash.symbols.diff
conformance/controlFlowStringIndex.symbols.diff

# Property access on index signatures now resolves to the `[__index]` symbol
## New symbol resolution lines are added where previously none existed (e.g. `c.x` now resolves to `T[__index]`)
## **Benign (improvement)** — Corsa is providing more symbol resolution info than Strada did; strictly more correct.
conformance/keyofAndIndexedAccess2.symbols.diff
conformance/noUncheckedIndexedAccessDestructuring.symbols.diff
conformance/objectSpreadIndexSignature.symbols.diff
compiler/propertyAccessOfReadonlyIndexSignature.symbols.diff

# `@enum` JSDoc tag no longer creates an additional declaration on the variable
## Where Strada added the `@enum` comment position as a Decl, Corsa only includes the variable declaration itself
## **Benign (intentional)** — CHANGES.md lists `@enum` as a removed Closure feature; the tag no longer has specific node parsing, so extra Decl entries are expected to disappear.
conformance/enumTag.symbols.diff
conformance/enumTagCircularReference.symbols.diff
conformance/enumTagImported.symbols.diff
conformance/enumTagOnExports.symbols.diff
conformance/enumTagUseBeforeDefCrash.symbols.diff
conformance/enumMergeWithExpando.symbols.diff
conformance/exportedAliasedEnumTag.symbols.diff
conformance/exportedEnumTypeAndValue.symbols.diff
conformance/jsDeclarationsEnumTag(target=es2015).symbols.diff
conformance/jsDeclarationsParameterTagReusesInputNodeInEmit1.symbols.diff
compiler/jsEnumCrossFileExport.symbols.diff
compiler/jsEnumTagOnObjectFrozen.symbols.diff
compiler/jsFileESModuleWithEnumTag.symbols.diff

# Unicode-escaped property names use bracket notation in symbol display
## e.g. `Class.x\u0078` → `Class[x\u0078]`, including private identifiers like `#x\u0078`
## **Benign** — cosmetic display change; bracket notation for non-standard identifiers is reasonable.
compiler/escapedIdentifiers.symbols.diff
compiler/unicodeEscapesInNames01(target=es2015).symbols.diff
compiler/unicodeEscapesInNames01(target=esnext).symbols.diff

# Constructor function `this` and `this.property` bindings removed
## Corsa no longer binds `this` to the constructor function symbol or creates property symbols from `this.x = ...` assignments in JS constructor functions
## **Benign (intentional)** — CHANGES.md documents that constructor function prototype assignment is removed; `@class`/`@constructor` is ignored; users should use `class` syntax instead.
conformance/constructorFunctions.symbols.diff
conformance/constructorFunctions2.symbols.diff
conformance/constructorFunctions3.symbols.diff
conformance/constructorFunctionsStrict.symbols.diff
conformance/constructorFunctionMergeWithClass.symbols.diff
conformance/constructorFunctionMethodTypeParameters.symbols.diff
conformance/constructorTagOnNestedBinaryExpression.symbols.diff
conformance/classCanExtendConstructorFunction.symbols.diff
conformance/propertiesOfGenericConstructorFunctions.symbols.diff
conformance/privateConstructorFunction.symbols.diff
conformance/thisPropertyAssignment.symbols.diff
conformance/thisPropertyAssignmentCircular.symbols.diff
conformance/thisTypeOfConstructorFunctions.symbols.diff
conformance/prototypePropertyAssignmentMergeAcrossFiles.symbols.diff
conformance/typeFromJSConstructor.symbols.diff
conformance/typeFromJSInitializer.symbols.diff
conformance/typeFromPropertyAssignment2.symbols.diff
conformance/typeFromPropertyAssignment3.symbols.diff
conformance/typeFromPropertyAssignment9.symbols.diff
conformance/typeFromPropertyAssignment9_1.symbols.diff
conformance/typeFromPropertyAssignment19.symbols.diff
conformance/typeFromPropertyAssignment20.symbols.diff
conformance/typeFromPropertyAssignment22.symbols.diff
conformance/typeFromPropertyAssignment27.symbols.diff
conformance/typeFromPropertyAssignment28.symbols.diff
conformance/typeFromPropertyAssignment40.symbols.diff
conformance/typeFromParamTagForFunction.symbols.diff
conformance/typeFromPrototypeAssignment.symbols.diff
conformance/typeFromPrototypeAssignment2.symbols.diff
conformance/typeFromPrototypeAssignment3.symbols.diff
conformance/typeFromPrototypeAssignment4.symbols.diff
conformance/typedefCrossModule.symbols.diff
conformance/assignmentToVoidZero2.symbols.diff
conformance/chainedPrototypeAssignment.symbols.diff
conformance/checkExportsObjectAssignPrototypeProperty.symbols.diff
conformance/exportNestedNamespaces.symbols.diff
conformance/inferringClassMembersFromAssignments2.symbols.diff
conformance/inferringClassMembersFromAssignments6.symbols.diff
conformance/callbackCrossModule.symbols.diff
conformance/moduleExportAlias2.symbols.diff
conformance/moduleExportNestedNamespaces.symbols.diff
conformance/moduleExportWithExportPropertyAssignment4.symbols.diff
compiler/functionExpressionNames.symbols.diff
compiler/javascriptDefinePropertyPrototypeNonConstructor.symbols.diff
compiler/jsdocFunctionClassPropertiesDeclaration.symbols.diff
compiler/jsDeclarationsGlobalFileConstFunction.symbols.diff
compiler/jsDeclarationsGlobalFileConstFunctionNamed.symbols.diff
compiler/jsFunctionWithPrototypeNoErrorTruncationNoCrash.symbols.diff
compiler/objectPropertyAsClass.symbols.diff
compiler/thisInObjectJs.symbols.diff
conformance/jsDeclarationsClassMethod.symbols.diff
conformance/jsDeclarationsExportAssignedConstructorFunction(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionClassesCjsExportAssignment(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionLikeClasses(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionLikeClasses2(target=es2015).symbols.diff
conformance/jsdocConstructorFunctionTypeReference.symbols.diff
conformance/jsdocFunctionType.symbols.diff
conformance/jsdocPrototypePropertyAccessWithType.symbols.diff
conformance/jsdocTemplateConstructorFunction.symbols.diff
conformance/jsdocTemplateConstructorFunction2.symbols.diff
conformance/jsdocTemplateTag4.symbols.diff
conformance/jsdocTemplateTag5.symbols.diff
conformance/jsdocTypeFromChainedAssignment.symbols.diff
conformance/jsdocTypeTagCast.symbols.diff

# Expando property declarations no longer merge with the base function/variable
## Fewer Decl entries on the host symbol; property assignments may lose their symbol or use a different parent name. This affects `foo.prop = ...` style assignments.
## **Benign (intentional)** — CHANGES.md documents expando declaration changes extensively: nested undeclared expandos, fallback initialisers, and prototype assignments are all removed/simplified. Fewer Decl entries is expected.
conformance/nullPropertyName.symbols.diff
conformance/newTargetNarrowing.symbols.diff
conformance/nestedPrototypeAssignment.symbols.diff
conformance/defaultPropertyAssignedClassWithPrototype.symbols.diff
conformance/circularMultipleAssignmentDeclaration.symbols.diff
conformance/checkJsdocTypeTagOnObjectProperty1.symbols.diff
conformance/checkSpecialPropertyAssignments.symbols.diff
conformance/classCanExtendConstructorFunction.symbols.diff
conformance/constructorFunctionMergeWithClass.symbols.diff
conformance/constructorFunctionMethodTypeParameters.symbols.diff
conformance/constructorFunctions2.symbols.diff
conformance/constructorFunctions3.symbols.diff
conformance/constructorFunctionsStrict.symbols.diff
conformance/contextualTypedSpecialAssignment.symbols.diff
conformance/exportDefaultNamespace.symbols.diff
conformance/exportNestedNamespaces.symbols.diff
conformance/globalMergeWithCommonJSAssignmentDeclaration.symbols.diff
conformance/commonJSImportNestedClassTypeReference.symbols.diff
conformance/inferringClassMembersFromAssignments2.symbols.diff
conformance/inferringClassMembersFromAssignments6.symbols.diff
conformance/inferringClassStaticMembersFromAssignments.symbols.diff
conformance/propertiesOfGenericConstructorFunctions.symbols.diff
conformance/propertyAssignmentUseParentType1.symbols.diff
conformance/propertyAssignmentUseParentType2.symbols.diff
conformance/propertyAssignmentUseParentType3.symbols.diff
conformance/prototypePropertyAssignmentMergeAcrossFiles2.symbols.diff
conformance/prototypePropertyAssignmentMergedTypeReference.symbols.diff
conformance/prototypePropertyAssignmentMergeWithInterfaceMethod.symbols.diff
conformance/spellingUncheckedJS.symbols.diff
conformance/thisTypeOfConstructorFunctions.symbols.diff
conformance/typedefTagExtraneousProperty.symbols.diff
conformance/typeFromJSConstructor.symbols.diff
conformance/typeFromJSInitializer.symbols.diff
conformance/typeFromPropertyAssignment.symbols.diff
conformance/typeFromPropertyAssignment4.symbols.diff
conformance/typeFromPropertyAssignment5.symbols.diff
conformance/typeFromPropertyAssignment6.symbols.diff
conformance/typeFromPropertyAssignment7.symbols.diff
conformance/typeFromPropertyAssignment8.symbols.diff
conformance/typeFromPropertyAssignment8_1.symbols.diff
conformance/typeFromPropertyAssignment9.symbols.diff
conformance/typeFromPropertyAssignment9_1.symbols.diff
conformance/typeFromPropertyAssignment10.symbols.diff
conformance/typeFromPropertyAssignment10_1.symbols.diff
conformance/typeFromPropertyAssignment11.symbols.diff
conformance/typeFromPropertyAssignment12.symbols.diff
conformance/typeFromPropertyAssignment13.symbols.diff
conformance/typeFromPropertyAssignment14.symbols.diff
conformance/typeFromPropertyAssignment15.symbols.diff
conformance/typeFromPropertyAssignment16.symbols.diff
conformance/typeFromPropertyAssignment17.symbols.diff
conformance/typeFromPropertyAssignment18.symbols.diff
conformance/typeFromPropertyAssignment19.symbols.diff
conformance/typeFromPropertyAssignment20.symbols.diff
conformance/typeFromPropertyAssignment22.symbols.diff
conformance/typeFromPropertyAssignment24.symbols.diff
conformance/typeFromPropertyAssignment25.symbols.diff
conformance/typeFromPropertyAssignment26.symbols.diff
conformance/typeFromPropertyAssignment27.symbols.diff
conformance/typeFromPropertyAssignment28.symbols.diff
conformance/typeFromPropertyAssignment29.symbols.diff
conformance/typeFromPropertyAssignment30.symbols.diff
conformance/typeFromPropertyAssignment31.symbols.diff
conformance/typeFromPropertyAssignment32.symbols.diff
conformance/typeFromPropertyAssignment33.symbols.diff
conformance/typeFromPropertyAssignment34.symbols.diff
conformance/typeFromPropertyAssignment35.symbols.diff
conformance/typeFromPropertyAssignment36.symbols.diff
conformance/typeFromPropertyAssignment37.symbols.diff
conformance/typeFromPropertyAssignment38.symbols.diff
conformance/typeFromPropertyAssignment39.symbols.diff
conformance/typeFromPropertyAssignment40.symbols.diff
conformance/typeFromPropertyAssignmentOutOfOrder(target=es2015).symbols.diff
conformance/typeFromPropertyAssignmentWithExport.symbols.diff
conformance/typeFromPrototypeAssignment.symbols.diff
conformance/typeFromPrototypeAssignment2.symbols.diff
conformance/typeFromPrototypeAssignment3.symbols.diff
conformance/typeFromPrototypeAssignment4.symbols.diff
conformance/typeTagCircularReferenceOnConstructorFunction.symbols.diff
conformance/typeTagPrototypeAssignment.symbols.diff
conformance/typedefCrossModule2.symbols.diff
conformance/typedefCrossModule3.symbols.diff
conformance/assignmentToVoidZero2.symbols.diff
conformance/chainedPrototypeAssignment.symbols.diff
conformance/checkExportsObjectAssignPrototypeProperty.symbols.diff
conformance/moduleExportAlias.symbols.diff
conformance/moduleExportAlias2.symbols.diff
conformance/moduleExportAlias4.symbols.diff
conformance/moduleExportAlias5.symbols.diff
conformance/moduleExportAliasImported.symbols.diff
conformance/moduleExportAliasUnknown.symbols.diff
conformance/moduleExportAssignment.symbols.diff
conformance/moduleExportAssignment2.symbols.diff
conformance/moduleExportAssignment4.symbols.diff
conformance/moduleExportAssignment5.symbols.diff
conformance/moduleExportAssignment7.symbols.diff
conformance/moduleExportNestedNamespaces.symbols.diff
conformance/moduleExportPropertyAssignmentDefault.symbols.diff
conformance/moduleExportsElementAccessAssignment.symbols.diff
conformance/moduleExportWithExportPropertyAssignment.symbols.diff
conformance/moduleExportWithExportPropertyAssignment2.symbols.diff
conformance/moduleExportWithExportPropertyAssignment3.symbols.diff
conformance/moduleExportWithExportPropertyAssignment4.symbols.diff
conformance/privateNameJsBadDeclaration(target=es2015).symbols.diff
conformance/jsdocConstructorFunctionTypeReference.symbols.diff
conformance/jsdocImplements_class.symbols.diff
conformance/jsdocPrototypePropertyAccessWithType.symbols.diff
conformance/jsdocReadonly.symbols.diff
conformance/jsdocTemplateClass.symbols.diff
conformance/jsdocTemplateConstructorFunction.symbols.diff
conformance/jsdocTemplateConstructorFunction2.symbols.diff
conformance/jsdocTemplateTag4.symbols.diff
conformance/jsdocTemplateTag5.symbols.diff
conformance/jsdocTypeFromChainedAssignment.symbols.diff
conformance/jsdocTypeReferenceToMergedClass.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport5.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport6.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport7.symbols.diff
conformance/jsContainerMergeJsContainer.symbols.diff
conformance/jsContainerMergeTsDeclaration.symbols.diff
conformance/jsContainerMergeTsDeclaration2.symbols.diff
conformance/jsContainerMergeTsDeclaration3.symbols.diff
conformance/jsDeclarationsClassLikeHeuristic(target=es2015).symbols.diff
conformance/jsDeclarationsClassMethod.symbols.diff
conformance/jsDeclarationsClassStatic(target=es2015).symbols.diff
conformance/jsDeclarationsClassStatic2.symbols.diff
conformance/jsDeclarationsClassStaticMethodAugmentation(target=es2015).symbols.diff
conformance/jsDeclarationsCrossfileMerge(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassExpressionAnonymousWithSub(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassExpressionShadowing(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassInstance3(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedConstructorFunction(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedConstructorFunctionWithSub(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignmentExpressionPlusSecondary(target=es2015).symbols.diff
conformance/jsDeclarationsExportForms(target=es2015).symbols.diff
conformance/jsDeclarationsExportSubAssignments(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionClassesCjsExportAssignment(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionKeywordProp(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionKeywordPropExhaustive(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionLikeClasses2(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionPrototypeStatic(target=es2015).symbols.diff
conformance/jsDeclarationsFunctions.symbols.diff
conformance/jsDeclarationsFunctionsCjs.symbols.diff
conformance/jsDeclarationsFunctionWithDefaultAssignedMember.symbols.diff
conformance/jsDeclarationsParameterTagReusesInputNodeInEmit1.symbols.diff
conformance/jsDeclarationsParameterTagReusesInputNodeInEmit2.symbols.diff
conformance/jsDeclarationsReactComponents(target=es2015).symbols.diff
conformance/jsDeclarationsTypeReferences3(target=es2015).symbols.diff
compiler/expandoFunctionNestedAssigments.symbols.diff
compiler/isolatedDeclarationErrors.symbols.diff
compiler/isolatedDeclarationErrorsExpandoFunctions.symbols.diff
compiler/isolatedDeclarationLazySymbols.symbols.diff
compiler/jsDeclarationsGlobalFileConstFunction.symbols.diff
compiler/jsDeclarationsGlobalFileConstFunctionNamed.symbols.diff
compiler/jsElementAccessNoContextualTypeCrash.symbols.diff
compiler/jsFunctionWithPrototypeNoErrorTruncationNoCrash.symbols.diff
compiler/jsxComponentTypeErrors.symbols.diff
compiler/jsxDeclarationsWithEsModuleInteropNoCrash.symbols.diff
compiler/lateBoundFunctionMemberAssignmentDeclarations.symbols.diff
compiler/noParameterReassignmentIIFEAnnotated.symbols.diff
compiler/noParameterReassignmentJSIIFE.symbols.diff
compiler/targetTypeTest1.symbols.diff
compiler/topLevelBlockExpando.symbols.diff
compiler/tsxStatelessComponentDefaultProps.symbols.diff
compiler/wellKnownSymbolExpando.symbols.diff

# CommonJS `exports`/`module.exports` symbol names and structure changed
## e.g. `Symbol("./module", ...)` → `Symbol(exports, ...)`, `Symbol(export=, ...)` → `Symbol(module, ...)`. The `exports` identifier itself is now always `Symbol(exports, ...)`.
## **Benign (intentional)** — CHANGES.md documents substantial CommonJS changes (mixing module.exports disallowed, aliasing removed, etc.). Symbol naming/structure changes follow from the rewritten CJS support.
conformance/assignmentToVoidZero1.symbols.diff
conformance/assignmentToVoidZero2.symbols.diff
conformance/binderUninitializedModuleExportsAssignment.symbols.diff
conformance/callbackCrossModule.symbols.diff
conformance/chainedPrototypeAssignment.symbols.diff
conformance/checkExportsObjectAssignProperty.symbols.diff
conformance/checkExportsObjectAssignPrototypeProperty.symbols.diff
conformance/checkObjectDefineProperty.symbols.diff
conformance/checkOtherObjectAssignProperty.symbols.diff
conformance/commonJSAliasedExport.symbols.diff
conformance/commonJSImportClassTypeReference.symbols.diff
conformance/commonJSImportExportedClassExpression.symbols.diff
conformance/commonJSImportNestedClassTypeReference.symbols.diff
conformance/commonJSReexport.symbols.diff
conformance/constructorFunctions2.symbols.diff
conformance/contextualTypedSpecialAssignment.symbols.diff
conformance/enumTagOnExports.symbols.diff
conformance/enumTagOnExports2.symbols.diff
conformance/exportedAliasedEnumTag.symbols.diff
conformance/exportNestedNamespaces.symbols.diff
conformance/exportNestedNamespaces2.symbols.diff
conformance/exportPropertyAssignmentNameResolution.symbols.diff
conformance/globalMergeWithCommonJSAssignmentDeclaration.symbols.diff
conformance/importAliasModuleExports.symbols.diff
conformance/jsdocTypeFromChainedAssignment2.symbols.diff
conformance/jsdocTypeFromChainedAssignment3.symbols.diff
conformance/jsdocTypeReferenceExports.symbols.diff
conformance/jsdocImportType.symbols.diff
conformance/jsdocImportType2.symbols.diff
conformance/jsdocImportTypeReferenceToClassAlias.symbols.diff
conformance/jsdocTypeReferenceToImportOfClassExpression.symbols.diff
conformance/jsdocTypeReferenceToImportOfFunctionExpression.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport1.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport2.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport3.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport4.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport5.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport6.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport7.symbols.diff
conformance/moduleExportAlias.symbols.diff
conformance/moduleExportAlias2.symbols.diff
conformance/moduleExportAlias3.symbols.diff
conformance/moduleExportAlias4.symbols.diff
conformance/moduleExportAlias5.symbols.diff
conformance/moduleExportAliasElementAccessExpression.symbols.diff
conformance/moduleExportAliasExports.symbols.diff
conformance/moduleExportAliasImported.symbols.diff
conformance/moduleExportAliasUnknown.symbols.diff
conformance/moduleExportAssignment.symbols.diff
conformance/moduleExportAssignment2.symbols.diff
conformance/moduleExportAssignment3.symbols.diff
conformance/moduleExportAssignment4.symbols.diff
conformance/moduleExportAssignment5.symbols.diff
conformance/moduleExportAssignment7.symbols.diff
conformance/moduleExportDuplicateAlias.symbols.diff
conformance/moduleExportDuplicateAlias2.symbols.diff
conformance/moduleExportDuplicateAlias3.symbols.diff
conformance/moduleExportNestedNamespaces.symbols.diff
conformance/moduleExportPropertyAssignmentDefault.symbols.diff
conformance/moduleExportsAliasLoop1.symbols.diff
conformance/moduleExportsAliasLoop2.symbols.diff
conformance/moduleExportsElementAccessAssignment.symbols.diff
conformance/moduleExportWithExportPropertyAssignment.symbols.diff
conformance/moduleExportWithExportPropertyAssignment2.symbols.diff
conformance/moduleExportWithExportPropertyAssignment3.symbols.diff
conformance/moduleExportWithExportPropertyAssignment4.symbols.diff
conformance/nestedDestructuringOfRequire.symbols.diff
conformance/nodeModulesAllowJsCjsFromJs(module=node16).symbols.diff
conformance/nodeModulesAllowJsCjsFromJs(module=node18).symbols.diff
conformance/nodeModulesAllowJsCjsFromJs(module=node20).symbols.diff
conformance/nodeModulesAllowJsCjsFromJs(module=nodenext).symbols.diff
conformance/nodeModulesAllowJsExportAssignment(module=node16).symbols.diff
conformance/nodeModulesAllowJsExportAssignment(module=node18).symbols.diff
conformance/nodeModulesAllowJsExportAssignment(module=node20).symbols.diff
conformance/nodeModulesAllowJsExportAssignment(module=nodenext).symbols.diff
conformance/paramTagOnCallExpression.symbols.diff
conformance/paramTagTypeResolution.symbols.diff
conformance/reExportJsFromTs.symbols.diff
conformance/requireOfESWithPropertyAccess.symbols.diff
conformance/requireTwoPropertyAccesses.symbols.diff
conformance/typeFromParamTagForFunction.symbols.diff
conformance/typeFromPropertyAssignment17.symbols.diff
conformance/typeFromPropertyAssignment37.symbols.diff
conformance/typedefCrossModule.symbols.diff
conformance/typedefCrossModule2.symbols.diff
conformance/typedefCrossModule3.symbols.diff
conformance/typedefCrossModule4.symbols.diff
conformance/typeTagModuleExports.symbols.diff
conformance/untypedModuleImport_allowJs.symbols.diff
conformance/jsDeclarationsClassExtendsVisibility(target=es2015).symbols.diff
conformance/jsDeclarationsClassStatic(target=es2015).symbols.diff
conformance/jsDeclarationsCommonjsRelativePath.symbols.diff
conformance/jsDeclarationsComputedNames(target=es2015).symbols.diff
conformance/jsDeclarationsCrossfileMerge(target=es2015).symbols.diff
conformance/jsDeclarationsDocCommentsOnConsts(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassExpression(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassExpressionAnonymous(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassExpressionAnonymousWithSub(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassExpressionShadowing(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassInstance1(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassInstance2(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedClassInstance3(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedConstructorFunction(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedConstructorFunctionWithSub(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignedVisibility(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignmentExpressionPlusSecondary(target=es2015).symbols.diff
conformance/jsDeclarationsExportAssignmentWithKeywordName(target=es2015).symbols.diff
conformance/jsDeclarationsExportDefinePropertyEmit.symbols.diff
conformance/jsDeclarationsExportDoubleAssignmentInClosure.symbols.diff
conformance/jsDeclarationsExportedClassAliases.symbols.diff
conformance/jsDeclarationsExportForms(target=es2015).symbols.diff
conformance/jsDeclarationsExportSubAssignments(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionClassesCjsExportAssignment(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionPrototypeStatic(target=es2015).symbols.diff
conformance/jsDeclarationsFunctionsCjs.symbols.diff
conformance/jsDeclarationsFunctionWithDefaultAssignedMember.symbols.diff
conformance/jsDeclarationsImportAliasExposedWithinNamespaceCjs.symbols.diff
conformance/jsDeclarationsJson(target=es2015).symbols.diff
conformance/jsDeclarationsPackageJson(target=es2015).symbols.diff
conformance/jsDeclarationsParameterTagReusesInputNodeInEmit1.symbols.diff
conformance/jsDeclarationsParameterTagReusesInputNodeInEmit2.symbols.diff
conformance/jsDeclarationsReexportAliasesEsModuleInterop(target=es2015).symbols.diff
conformance/jsDeclarationsReexportedCjsAlias(target=es2015).symbols.diff
conformance/jsDeclarationsReferenceToClassInstanceCrossFile.symbols.diff
conformance/jsDeclarationsTypeAliases.symbols.diff
conformance/jsDeclarationsTypedefAndImportTypes.symbols.diff
conformance/jsDeclarationsTypedefAndLatebound.symbols.diff
conformance/jsDeclarationsTypedefPropertyAndExportAssignment.symbols.diff
conformance/jsDeclarationsTypeReassignmentFromDeclaration.symbols.diff
conformance/jsDeclarationsTypeReassignmentFromDeclaration2.symbols.diff
conformance/jsDeclarationsTypeReferences(target=es2015).symbols.diff
conformance/jsDeclarationsTypeReferences2(target=es2015).symbols.diff
conformance/jsDeclarationsTypeReferences3(target=es2015).symbols.diff
compiler/functionExpressionNames.symbols.diff
compiler/importHelpersCommonJSJavaScript(verbatimmodulesyntax=false).symbols.diff
compiler/importHelpersCommonJSJavaScript(verbatimmodulesyntax=true).symbols.diff
compiler/importNonExportedMember12.symbols.diff
compiler/javascriptCommonjsModule.symbols.diff
compiler/javascriptImportDefaultBadExport.symbols.diff
compiler/jsDeclarationEmitExportAssignedArray.symbols.diff
compiler/jsDeclarationEmitExportAssignedFunctionWithExtraTypedefsMembers.symbols.diff
compiler/jsEnumTagOnObjectFrozen.symbols.diff
compiler/jsExportAssignmentNonMutableLocation.symbols.diff
compiler/jsExportMemberMergedWithModuleAugmentation.symbols.diff
compiler/jsExportMemberMergedWithModuleAugmentation2.symbols.diff
compiler/jsExportMemberMergedWithModuleAugmentation3.symbols.diff
compiler/jsFileClassPropertyInitalizationInObjectLiteral.symbols.diff
compiler/jsFileCompilationBindDeepExportsAssignment.symbols.diff
compiler/jsFileCompilationExternalPackageError.symbols.diff
compiler/moduleExportsTypeNoExcessPropertyCheckFromContainedLiteral.symbols.diff
compiler/modulePreserve4.symbols.diff
compiler/moduleResolution_explicitNodeModulesImport.symbols.diff
compiler/pushTypeGetTypeOfAlias.symbols.diff
compiler/requireOfJsonFileInJsFile.symbols.diff
compiler/resolveNameWithNamspace.symbols.diff
compiler/truthinessCallExpressionCoercion4.symbols.diff

# `import =` alias resolves to the target's own symbol name
## e.g. `import m4 = b.a` — the `.a` now resolves to `Symbol(m4, ...)` instead of `Symbol(m3.a, ...)`
## **Benign** — resolving to the alias's own name rather than the qualified target is a reasonable simplification; the symbol still points to the right thing.
conformance/commonJSImportNotAsPrimaryExpression.symbols.diff
conformance/exportsAndImports4-es6.symbols.diff
conformance/exportsAndImports4(target=es2015).symbols.diff
conformance/importEquals3.symbols.diff
conformance/importStatementsInterfaces.symbols.diff
compiler/dottedModuleName2.symbols.diff
compiler/exportDefaultProperty.symbols.diff
compiler/importedEnumMemberMergedWithExportedAliasIsError.symbols.diff
compiler/importOnAliasedIdentifiers.symbols.diff
compiler/isolatedModulesReExportType.symbols.diff
compiler/moduledecl(target=es2015).symbols.diff
compiler/nodeNextCjsNamespaceImportDefault1.symbols.diff
compiler/nodeNextCjsNamespaceImportDefault2.symbols.diff

# Declaration position `(file, --, --)` changed to actual numeric positions
## Symbols in tslib and other well-known declaration files now report concrete line/column instead of `--`
## **Benign (improvement)** — reporting actual positions is strictly more informative than `--`.
compiler/importHelpersCommonJSJavaScript(verbatimmodulesyntax=false).symbols.diff
compiler/importHelpersCommonJSJavaScript(verbatimmodulesyntax=true).symbols.diff
compiler/importHelpersES6.symbols.diff
compiler/importHelpersInAmbientContext(target=es2015).symbols.diff
compiler/importHelpersVerbatimModuleSyntax.symbols.diff
compiler/importHelpersWithExportStarAs(esmoduleinterop=true,module=commonjs).symbols.diff
compiler/importHelpersWithExportStarAs(esmoduleinterop=true,module=es2015).symbols.diff
compiler/importHelpersWithExportStarAs(esmoduleinterop=true,module=es2020).symbols.diff
compiler/importHelpersWithImportOrExportDefault(esmoduleinterop=true,module=commonjs).symbols.diff
compiler/importHelpersWithImportOrExportDefault(esmoduleinterop=true,module=es2015).symbols.diff
compiler/importHelpersWithImportOrExportDefault(esmoduleinterop=true,module=es2020).symbols.diff
compiler/importHelpersWithImportStarAs(esmoduleinterop=true,module=commonjs).symbols.diff
compiler/importHelpersWithImportStarAs(esmoduleinterop=true,module=es2015).symbols.diff
compiler/importHelpersWithImportStarAs(esmoduleinterop=true,module=es2020).symbols.diff
compiler/jsNoImplicitAnyNoCascadingReferenceErrors.symbols.diff
compiler/modulePreserveImportHelpers.symbols.diff
compiler/tslibMissingHelper.symbols.diff
compiler/tslibMultipleMissingHelper.symbols.diff
compiler/tslibNotFoundDifferentModules.symbols.diff

# `this.property` without initializer in class constructor no longer creates a property symbol
## e.g. `this.id;` (bare expression statement with @type annotation) no longer binds
## **Benign (intentional)** — CHANGES.md explicitly documents this: "A this-property expression with a type annotation in the constructor no longer creates a property." Users should use class field declarations or provide an initializer.
conformance/callbackTag2.symbols.diff
conformance/constructorTagOnObjectLiteralMethod.symbols.diff
conformance/assignmentToVoidZero2.symbols.diff
compiler/objectPropertyAsClass.symbols.diff
compiler/thisInObjectJs.symbols.diff
compiler/unusedTypeParameters_templateTag2.symbols.diff

# Namespace-qualified symbol names simplified
## e.g. `valueZ.otherType` → `otherType` — the parent qualifier is dropped from the symbol name
## **Benign** — cosmetic change to display name; the symbol identity is still correct, just shown without the parent prefix.
conformance/arbitraryModuleNamespaceIdentifiers_module(module=commonjs).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=es2020).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=es2022).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=es6).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=esnext).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=node16).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=node18).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=node20).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=nodenext).symbols.diff
conformance/arbitraryModuleNamespaceIdentifiers_module(module=preserve).symbols.diff
conformance/extendsClause.symbols.diff
conformance/nodeModulesDeclarationEmitWithPackageExports(module=node16).symbols.diff
conformance/nodeModulesDeclarationEmitWithPackageExports(module=node18).symbols.diff
conformance/nodeModulesDeclarationEmitWithPackageExports(module=node20).symbols.diff
conformance/nodeModulesDeclarationEmitWithPackageExports(module=nodenext).symbols.diff
conformance/verbatimModuleSyntaxRestrictionsCJS.symbols.diff
conformance/jsdocImplements_class.symbols.diff
conformance/jsdocTemplateTag4.symbols.diff
conformance/jsdocTemplateTag5.symbols.diff
compiler/defaultDeclarationEmitShadowedNamedCorrectly.symbols.diff
compiler/isolatedModulesReExportType.symbols.diff
compiler/jsEnumCrossFileExport.symbols.diff

# `@overload` JSDoc tags now contribute Decl entries to the host function/method symbol
## The host function symbol now includes declaration entries from @overload JSDoc tags
## **Benign (improvement)** — overload signatures should be part of the function's declarations; this aligns with TS behavior for overloads.
compiler/jsFileFunctionOverloads.symbols.diff
compiler/jsFileFunctionOverloads2.symbols.diff
compiler/jsFileMethodOverloads.symbols.diff
compiler/jsFileMethodOverloads2.symbols.diff
compiler/jsFileMethodOverloads3.symbols.diff
conformance/overloadTag1.symbols.diff

# `this` binding changed from `Symbol(this)` to `Symbol((Missing), ...)` in JS functions with @this tags
## JS functions using `@this` or JSDoc `this` type annotations now bind `this` to a `(Missing)` placeholder symbol
## **Clear bug** — `(Missing)` is a sentinel for "failed to resolve"; `this` should resolve to a real symbol, not a placeholder. Likely a binding gap for `@this` tag handling.
compiler/thisInFunctionCallJs.symbols.diff
conformance/inferThis.symbols.diff
conformance/constructorTagWithThisTag.symbols.diff
conformance/thisPrototypeMethodCompoundAssignmentJs.symbols.diff
conformance/thisTag1.symbols.diff

# File path normalization in symbol baseline headers
## `./` prefix dropped from file paths, double-slashes collapsed, or Windows path separators normalized
## **Benign** — trivial path normalization differences in test baselines; no semantic impact.
compiler/mergeSymbolReexportedTypeAliasInstantiation.symbols.diff
compiler/mergeSymbolReexportInterface.symbols.diff
compiler/mergeSymbolRexportFunction.symbols.diff
compiler/taggedTemplateWithoutDeclaredHelper(target=es2015).symbols.diff
compiler/tripleSlashReferenceAbsoluteWindowsPath.symbols.diff
compiler/typeGuardNarrowsIndexedAccessOfKnownProperty8.symbols.diff
compiler/uniqueSymbolJs.symbols.diff
conformance/callbackTag4.symbols.diff
conformance/exportSpecifiers_js.symbols.diff
conformance/importAttributes10.symbols.diff
conformance/importAttributes11.symbols.diff
conformance/importAttributes9.symbols.diff
conformance/importSpecifiers_js.symbols.diff
conformance/typeSatisfactionWithDefaultExport.symbols.diff
compiler/jsDeclarationEmitExportedClassWithExtends.symbols.diff

# `Object.defineProperty` key symbol now uses bracket-string or quoted form
## e.g. `m1.thing` → `m1["thing"]` or `Symbol(a)` → `Symbol("a")` for string-literal property keys in defineProperty calls
## **Benign** — cosmetic display difference; bracket notation with string keys is arguably more accurate for computed/defineProperty properties.
conformance/checkExportsObjectAssignProperty.symbols.diff
conformance/checkObjectDefineProperty.symbols.diff
conformance/checkOtherObjectAssignProperty.symbols.diff
conformance/jsDeclarationsExportDefinePropertyEmit.symbols.diff
conformance/jsDeclarationsGetterSetter.symbols.diff
compiler/jsExpandoObjectDefineProperty.symbols.diff

# `require("...")` string-literal source-module symbol entry removed
## The symbol entry for the string literal argument to `require()` is no longer emitted
## **Benign** — the string literal inside `require()` never needed its own symbol; removing it is a cleanup.
conformance/jsDeclarationsTypeReassignmentFromDeclaration2.symbols.diff
conformance/jsDeclarationsTypeReferences(target=es2015).symbols.diff
conformance/jsDeclarationsTypeReferences2(target=es2015).symbols.diff
conformance/jsDeclarationsTypeReferences4(target=es2015).symbols.diff
conformance/jsdocTypeReferenceToImport.symbols.diff

# `this` inside prototype object-literal methods now binds to `__object` instead of the constructor
## When `Foo.prototype = { method() { this... } }`, `this` resolves to the object-literal `__object` symbol
## **Benign (intentional)** — CHANGES.md documents: "Assigning to the `prototype` property of a function no longer makes it a constructor function." So `this` correctly resolves to the object literal, not a constructor instance.
conformance/typeFromPrototypeAssignment.symbols.diff
conformance/typeFromPrototypeAssignment2.symbols.diff
conformance/typeFromPrototypeAssignment3.symbols.diff
conformance/typeFromPrototypeAssignment4.symbols.diff
conformance/typeFromContextualThisType.symbols.diff
compiler/jsFunctionWithPrototypeNoErrorTruncationNoCrash.symbols.diff

# `.constructor` property access now resolves to `Object.constructor` from lib
## Accessing `.constructor` on a class type now resolves to the inherited `Object.constructor` symbol
## **Benign (improvement)** — resolving `.constructor` to the inherited lib symbol is correct; Strada was likely just not resolving it at all.
conformance/typesWithPrivateConstructor.symbols.diff
conformance/typesWithProtectedConstructor.symbols.diff

# Import specifier in re-exports resolves to alias target symbol name
## In `import { x as y }`, the symbol for `x` now reports as the target name `y`
## **Benign** — the specifier `x` is the alias source; resolving it to the local alias name `y` is a display choice, not a semantic error.
conformance/nodeModulesAllowJsSynchronousCallErrors(module=node16).symbols.diff
conformance/nodeModulesAllowJsSynchronousCallErrors(module=node18).symbols.diff
conformance/nodeModulesAllowJsSynchronousCallErrors(module=node20).symbols.diff
conformance/nodeModulesAllowJsSynchronousCallErrors(module=nodenext).symbols.diff
conformance/nodeModulesSynchronousCallErrors(module=node16).symbols.diff
conformance/nodeModulesSynchronousCallErrors(module=node18).symbols.diff
conformance/nodeModulesSynchronousCallErrors(module=node20).symbols.diff
conformance/nodeModulesSynchronousCallErrors(module=nodenext).symbols.diff

# Test baseline reduction (removed extra file outputs)
## Corsa drops extraneous file sections from baseline output
## **Benign** — no semantic change, just cleaner test output.
compiler/augmentExportEquals2.symbols.diff

# Private field symbol uses dot notation instead of bracket notation
## e.g. `K[#𝑚]` → `K.#𝑚` — private fields use standard dot notation in symbol display
## **Benign** — dot notation is the standard way to write private fields in JS/TS.
compiler/extendedUnicodePlaneIdentifiers.symbols.diff

# False merging of JS exports with TS module augmentation removed
## JS export member no longer incorrectly merges its symbol with a TS module augmentation declaration
## **Benign (improvement)** — removes an incorrect cross-language symbol merge.
compiler/jsExportMemberMergedWithModuleAugmentation.symbols.diff

# Prototype assignment no longer creates constructor function symbol
## `Foo.prototype.bar` resolves to `Function.prototype` rather than a custom symbol
## **Benign (intentional)** — documented in CHANGES.md; prototype assignment no longer creates constructors.
compiler/jsFileMethodOverloads4.symbols.diff

# Spurious duplicate Decl entries removed from late-bound assignments
## Removes extra Decl entries that were incorrectly added to property symbols
## **Benign (improvement)** — removes false duplicates in symbol declaration lists.
compiler/lateBoundAssignmentCandidateJS3.symbols.diff
compiler/jsdocTypedefNoCrash.symbols.diff
conformance/lateBoundAssignmentDeclarationSupport7.symbols.diff

# Misspelled JSDoc typedef no longer creates false user-code declaration
## `Animation` now only shows lib.dom declarations, not a spurious user-code Decl from a misspelled typedef
## **Benign (improvement)** — removes incorrect declaration that shadowed the real lib type.
compiler/misspelledJsDocTypedefTags.symbols.diff

# Arbitrary module namespace identifiers now resolve correctly
## `Symbol((Missing))` → `Symbol("invalid 2")` — string-named exports now get their actual name
## **Benign (improvement)** — `(Missing)` was a failure sentinel; the actual name is now resolved.
conformance/arbitraryModuleNamespaceIdentifiers_syntax.symbols.diff

# Previously-missing symbols for destructuring pattern bindings now emitted
## Destructured parameter names like `additionalFiles` and `a` now have symbol entries
## **Benign (improvement)** — Corsa tracks more symbols than Strada did for destructuring.
conformance/destructuringParameterDeclaration9(strict=true).symbols.diff

# Prototype reference now uses specific class instead of generic `Function.prototype`
## e.g. `Function.prototype` → `MobileDetect.prototype` or `A.prototype` — more precise type
## **Benign (improvement)** — resolving to the specific class is more accurate.
conformance/fixSignatureCaching.symbols.diff
conformance/privateNameBadDeclaration(target=es2015).symbols.diff

# `globalThis`/`window` user declarations removed
## User-code property assignments to `globalThis` no longer add Decl entries to the global symbol
## **Benign (intentional)** — CHANGES.md documents `this` alias for globalThis removal.
conformance/globalThisPropertyAssignment.symbols.diff

# Class member inference corrected for derived classes
## `Base.p` → `Derived.p` — property assigned in derived class now correctly attributed to derived
## **Benign (improvement)** — the property is assigned in the derived class, so it belongs there.
conformance/inferringClassMembersFromAssignments3.symbols.diff
conformance/inferringClassMembersFromAssignments4.symbols.diff
conformance/inferringClassMembersFromAssignments7.symbols.diff

# `this.x` property references now tracked in JSDoc augmentation contexts
## Previously-untracked `this.x` references in `@augments`-tagged classes now have symbol entries
## **Benign (improvement)** — more complete symbol resolution.
conformance/jsdocAugmentsMissingType.symbols.diff

# JSDoc type syntax in TypeScript files corrected
## Symbol names for parsed JSDoc types updated to match corrected parser behavior
## **Benign (intentional)** — matches CHANGES.md parser changes for JSDoc in TS.
conformance/jsdocDisallowedInTypescript.symbols.diff

# Minor symbol reference formatting changes
## Removed one `s.toString` symbol reference; minor cleanup
## **Benign** — trivial formatting difference with no semantic impact.
conformance/jsdocParseParenthesizedJSDocParameter.symbols.diff

# `this` correctly resolves to class type in `@this`-typed context
## `this` symbol resolves to `Foo` class rather than bare `this` reference
## **Benign (improvement)** — `@this` specifies the type, so resolution to that type is correct.
conformance/jsdocThisType.symbols.diff

# Declaration positions corrected in import alias declarations
## Column numbers in Decl position tuples corrected (e.g. `9,50` → `9,4`)
## **Benign (improvement)** — more accurate source positions.
conformance/jsDeclarationsImportAliasExposedWithinNamespace.symbols.diff
conformance/jsDeclarationsImportAliasExposedWithinNamespaceCjs.symbols.diff

# Duplicate Decl entries for conflicting private names removed
## When a private method and field conflict (already an error), duplicate Decl entries are eliminated
## **Benign** — these are error cases; cleaner symbol output for invalid code.
conformance/privateNameDuplicateField.symbols.diff

# False private field tracking removed for invalid assignment
## Invalid private name assignment outside class no longer creates a symbol entry
## **Benign (improvement)** — shouldn't track symbols for code that's semantically invalid.
conformance/privateNameJsBadAssignment.symbols.diff

# Additional Decl entries for multi-declaration function symbols
## `flatMap` function with multiple declarations now shows all Decl entries on the symbol
## **Benign** — more complete declaration tracking.
conformance/templateInsideCallback.symbols.diff

# Expando/class-expression symbol cleanup
## Removed duplicate declarations; fixed class-expression inner name usage
## **Benign (intentional)** — expando cleanup consistent with documented CHANGES.md expando changes.
conformance/typeFromPropertyAssignment.symbols.diff
