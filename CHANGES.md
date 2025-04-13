# CHANGES.md

This file lists intentional changes made while optimizing the given code across several modules. It also highlights improvements made beyond optimization. Made by: [NopAngel](https://github.com/NopAngel/)

## Summary

- **Files Optimized:** 13
- **Additional Changes:** Modularization, readability enhancements, error handling improvements, and consistent use of immutability and best practices.

---

## Scanner

1. Node positions now use UTF8 offsets from the beginning of the file, replacing UTF16 offsets. Node positions in files with non-ASCII characters are now larger than previously.

---

## Parser

1. Source files no longer include an EndOfFile token as their last child.
2. Malformed `...T?` at the end of a tuple results in a parse error rather than a grammar error.
3. Malformed string ImportSpecifiers (`import x as "OOPS" from "y"`) now contain the string's text instead of an empty identifier.
4. Empty binding elements no longer have a separate `OmittedExpression` kind. These elements are now `BindingElement` with `nil Initialiser`, `Name`, and `DotDotDotToken`.
5. `ShorthandPropertyAssignment` no longer includes an `EqualsToken` as a child when paired with an `ObjectAssignmentInitializer`.
6. JSDoc nodes now include leading whitespace in their location.
7. Comments in JSDoc are always parsed as `JSDocText` nodes, where `string` is no longer part of the `comment` type.
8. Leading/trailing whitespace/asterisks, as well as initial `/**` in JSDocText nodes, are no longer erroneously included in their location.
9. `JSDocMemberName` is parsed as `QualifiedName`. Previously differentiated by type, `QualifiedName` now has a less restrictive type for its left child.

---

## JSDoc Types

JSDoc types are parsed in normal type annotation positions, showing grammar errors where applicable. Corsa no longer parses certain JSDoc types and tags:

1. No postfix `T?` and `T!` types. Prefix `?T` and `!T` types are still parsed, with `!T` retaining no semantics.
2. No Closure `function(string,string): void` types.
3. No standalone `?` type in JSDoc.
4. No JSDoc module namepaths: `module:folder/file.C`.

---

## JSDoc Tags

Previously specific nodes are now parsed as generic `JSDocTag` nodes:

1. `@class`
2. `@throws`
3. `@author`
4. `@enum`
