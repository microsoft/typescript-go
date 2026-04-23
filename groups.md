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








