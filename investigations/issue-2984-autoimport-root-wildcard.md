# Investigation: microsoft/typescript-go#2984

## Summary
Auto-import did not prioritize (or even offer) `package.json#imports` root wildcard specifiers (`#/*`) for deep files, even with `importModuleSpecifierPreference` set to `shortest`/`non-relative`.

## Reproduction
Added fourslash test:
- `internal/fourslash/tests/autoImport_issue2984_rootWildcardVitest_test.go`

Scenario:
- `package.json` contains:
  - `"#/*": { "vitest": "./src/*", "types": "./src/*", "node": "./build/*", "default": "./src/*" }`
- Importing file at `/feature/very/deep/path/consumer.ts`
- Exported symbol at `/src/domain/entities/entity.ts`

### Failing behavior before fix
For `ImportModuleSpecifierPreference: "shortest"`:
- **Actual top suggestion:** `../../../../src/domain/entities/entity`
- **Expected:** `#/domain/entities/entity.js` (non-relative package import alias)

Same behavior reproduced for `ImportModuleSpecifierPreference: "non-relative"`.

## Root cause
File: `internal/modulespecifiers/specifiers.go`
Function: `tryGetModuleNameFromPackageJsonImports`

The imports-key validation incorrectly rejected all keys starting with `"#/"`:

```go
if !strings.HasPrefix(k, "#") || k == "#" || strings.HasPrefix(k, "#/") {
    continue // invalid imports entry
}
```

This rejects valid root wildcard key `#/*`, so `tryGetModuleNameFromPackageJsonImports` returned `""`, forcing fallback to relative specifiers.

## Minimal patch
Allow `#/*` while keeping other invalid `#/...` keys rejected:

```go
if !strings.HasPrefix(k, "#") || k == "#" || strings.HasPrefix(k, "#/") && !strings.HasPrefix(k, "#/*") {
    continue // invalid imports entry
}
```

## Validation after patch
Targeted test:

```bash
go test ./internal/fourslash/tests -run TestAutoImport_issue2984_rootWildcardVitest -count=1 -v
```

Observed:
- `shortest entity specifiers: [#/domain/entities/entity.js]`
- `non-relative entity specifiers: [#/domain/entities/entity.js]`

This confirms root wildcard package import is now preferred over deep relative import.
