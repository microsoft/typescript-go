# Investigation: Crash with Malformed package.json Imports

## Problem Summary
The language server crashes with the error:
```
panic handling request textDocument/diagnostic: should be able to extract TS extension from string that passes IsDeclarationFileName
```

This occurs in `/internal/checker/checker.go` at lines 14817 and 14833.

## Root Cause Analysis

### Reproduction Case
The crash occurs with a malformed `package.json` imports mapping:
```json
{
  "imports": {
    "#/*": {
      "types": "./src/*ts",    // Missing dot before ts! Should be "./src/*.ts"
      "default": "./dist/*js"
    }
  }
}
```

When importing `#/b.`, the following happens:

1. **Import statement**: `import * as b from "#/b."`
   - Note: The import literally ends with a dot, no extension (e.g., from fourslash marker position)

2. **Pattern matching**:
   - Pattern key: `"#/*"`
   - Pattern target: `"./src/*ts"` (malformed!)
   - The `*` matches `"b."` from the import

3. **Pattern substitution**:
   - Target `"./src/*ts"` becomes `"./src/b.ts"` (replacing `*` with `"b."`)
   - The file `/src/b.ts` is successfully resolved!

4. **ResolvedUsingTsExtension flag**:
   - Set to `true` in `/internal/module/resolver.go:1607`:
     ```go
     resolvedUsingTsExtension: packageJSONValue != "" && !strings.HasSuffix(packageJSONValue, extension)
     ```
   - This is true because `"./src/*ts"` doesn't end with `.ts`

5. **Checker validation**:
   - In `checker.go:14825`, the code sees `resolvedUsingTsExtension` is true
   - It assumes `moduleReference` must have a TS extension that can be extracted
   - `TryExtractTSExtension("#/b.")` returns empty string (because `"#/b."` doesn't end with a valid TS extension)
   - **PANIC!**

### Debug Output
```
DEBUG: moduleReference="#/b.", resolvedFileName="/src/b.ts"
DEBUG: IsDeclarationFileName=false, TryExtractTSExtension=""
```

### Key Insight
The code assumes that when `resolvedUsingTsExtension` is true, the `moduleReference` must contain a TS extension. However, this assumption is violated when:
- Package.json pattern mapping is used
- The mapped pattern is malformed (e.g., missing dots)
- The import specifier doesn't actually have a TS extension

## Why IsDeclarationFileName Returns False
Looking at the implementation in `/internal/tspath/extension.go`:
```go
func IsDeclarationFileName(fileName string) bool {
    return GetDeclarationFileExtension(fileName) != ""
}

func GetDeclarationFileExtension(fileName string) string {
    base := GetBaseFileName(fileName)
    for _, ext := range SupportedDeclarationExtensions {
        if strings.HasSuffix(base, ext) {
            return ext
        }
    }
    // ... more checks for .d.* patterns
    return ""
}
```

The string `"#/b."` doesn't end with any declaration extension (`.d.ts`, `.d.cts`, `.d.mts`), so `IsDeclarationFileName` correctly returns false.

## Why TryExtractTSExtension Returns Empty
Looking at the implementation:
```go
func TryExtractTSExtension(fileName string) string {
    for _, ext := range supportedTSExtensionsForExtractExtension {
        if FileExtensionIs(fileName, ext) {
            return ext
        }
    }
    return ""
}
```

The supported extensions are: `.d.ts`, `.d.cts`, `.d.mts`, `.ts`, `.tsx`, `.mts`, `.cts`.
The string `"#/b."` doesn't end with any of these, so it returns empty string.

## Proposed Fix

The panic occurs because the code assumes that if `resolvedUsingTsExtension` is true, then `moduleReference` must have a TS extension. This assumption is invalid.

**Solution**: When `TryExtractTSExtension` returns empty string, skip the error reporting instead of panicking. This makes sense because:
1. If the module reference doesn't actually have a TS extension, these specific errors don't apply
2. The resolution was successful (the file was found)
3. The malformed package.json is the real issue, not the import statement

### Specific Changes
In `/internal/checker/checker.go`, replace the panic at lines 14812-14817 and 14828-14833 with a simple check that skips error reporting when extraction fails:

```go
// Line 14809-14824 (first branch)
if resolvedModule.ResolvedUsingTsExtension && tspath.IsDeclarationFileName(moduleReference) {
    if ast.FindAncestor(location, ast.IsEmittableImport) != nil {
        tsExtension := tspath.TryExtractTSExtension(moduleReference)
        if tsExtension == "" {
            // Skip error reporting if we can't extract a TS extension.
            // This can happen with malformed package.json mappings.
            return  
        }
        c.error(
            errorNode,
            diagnostics.A_declaration_file_cannot_be_imported_without_import_type_Did_you_mean_to_import_an_implementation_file_0_instead,
            c.getSuggestedImportSource(moduleReference, tsExtension, mode),
        )
    }
}

// Line 14825-14840 (second branch)
else if resolvedModule.ResolvedUsingTsExtension && !c.compilerOptions.AllowImportingTsExtensionsFrom(importingSourceFile.FileName()) {
    if ast.FindAncestor(location, ast.IsEmittableImport) != nil {
        tsExtension := tspath.TryExtractTSExtension(moduleReference)
        if tsExtension == "" {
            // Skip error reporting if we can't extract a TS extension.
            // This can happen with malformed package.json mappings.
            return
        }
        c.error(
            errorNode,
            diagnostics.An_import_path_can_only_end_with_a_0_extension_when_allowImportingTsExtensions_is_enabled,
            tsExtension,
        )
    }
}
```

## Alternative Fixes Considered

### 1. Fix the pattern substitution in module resolver
Make the module resolver smarter about pattern substitution to avoid creating invalid module references. However, this is complex and might break other use cases.

### 2. Set ResolvedUsingTsExtension more carefully
Only set `resolvedUsingTsExtension` to true when the module reference actually has a TS extension. However, this might miss legitimate cases where the flag should be set.

### 3. Validate package.json patterns earlier
Add validation to reject malformed patterns like `"./src/*ts"`. However, this is a separate concern and doesn't fix the immediate crash.

## Recommendation
Proceed with the proposed fix (option 1 above) because:
- It's the simplest and safest fix
- It handles the edge case gracefully without breaking existing functionality
- It allows the resolution to succeed when the file is actually found
- Users will still see other errors related to the malformed package.json
