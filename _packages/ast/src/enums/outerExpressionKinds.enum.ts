//
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!! THIS FILE IS AUTO-GENERATED — DO NOT EDIT !!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
// Source: internal/ast/utilities.go
// Regenerate: npx hereby generate:enums
//
export enum OuterExpressionKinds {
    Parentheses = 1 << 0,
    TypeAssertions = 1 << 1,
    NonNullAssertions = 1 << 2,
    PartiallyEmittedExpressions = 1 << 3,
    ExpressionsWithTypeArguments = 1 << 4,
    Satisfies = 1 << 5,
    ExcludeJSDocTypeAssertion = 1 << 6,
    Assertions = TypeAssertions | NonNullAssertions | Satisfies,
    All = Parentheses | Assertions | PartiallyEmittedExpressions | ExpressionsWithTypeArguments,
}
