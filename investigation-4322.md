# Investigation: #4322 — checker depends on unimplemented aliasResolver method during auto-imports

## Symptom

Crash with stack:

```
internal/ls/autoimport.(*aliasResolver).IsSourceFileDefaultLibrary()  aliasresolver.go:219
internal/checker.(*Checker).elaborateElement()                        relater.go:609
internal/checker.(*Checker).elaborateObjectLiteral()                  relater.go:513
internal/checker.(*Checker).elaborateError()                          relater.go:462
internal/checker.(*Checker).checkTypeRelatedToAndOptionallyElaborate() relater.go:434
internal/checker.(*Checker).isSignatureApplicable()                   checker.go:9274
internal/checker.(*Checker).reportCallResolutionErrors()              checker.go:9626
...
internal/checker.(*Checker).checkDeclarationInitializer()             checker.go:16691
internal/checker.(*Checker).getTypeForVariableLikeDeclaration()
```

## Root cause

`internal/ls/autoimport/aliasresolver.go` provides a lightweight `checker.Program`
implementation (`aliasResolver`) used during auto-import export extraction
(`registry.go` `extractPackage` and the ambient-module second pass). Many of its
`checker.Program` methods are stubbed with `panic("unimplemented")` because they
were assumed never to be reached.

`IsSourceFileDefaultLibrary` is one such stub, but the checker DOES call it. In
`relater.go`, `elaborateElement` (object-literal error elaboration) calls
`c.program.IsSourceFileDefaultLibrary(...)` at lines 594 and 609 to decide whether
to attach "The expected type comes from property ... declared here" related info.

So whenever the checker, while extracting exports through an `aliasResolver`,
type-checks an expression that produces an object-literal assignability error
(e.g. `export const x = f({ a: 1 })` where `f` expects `{ a: string }`), the
elaboration path hits the stub and panics.

## Reproduction

`internal/ls/autoimport/aliasresolver_crash_test.go` builds an `aliasResolver`
over a single source file containing a failing call expression and runs the
checker over it. Before the fix this panics with the exact issue stack
(`IsSourceFileDefaultLibrary -> elaborateElement -> elaborateObjectLiteral`).

I also tried reproducing through the full session/registry API (building
auto-import buckets for a node_modules package with the offending code), but the
high-level extraction path only resolves symbols/aliases and rarely forces the
lazy type computation that triggers elaboration, so it did not reliably crash.
The white-box test exercises the identical crashing call chain deterministically.

## Fix

Implement `aliasResolver.IsSourceFileDefaultLibrary` to return `false` instead of
panicking. The `aliasResolver` only ever processes package source files (never
default-library files), so `false` is the correct answer; the related-info
diagnostics produced during extraction are discarded anyway.
