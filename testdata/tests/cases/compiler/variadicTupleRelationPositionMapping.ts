// @strict: true
// @noEmit: true

type TrailingElement<
  A extends unknown[],
  Result extends [...A, boolean],
> = Result;

type MatchingTrailingElement<A extends unknown[]> = TrailingElement<
  A,
  [...A, boolean]
>;

type ExcessTrailingElements<A extends unknown[]> = TrailingElement<
  A,
  [...A, boolean, boolean, boolean, boolean]
>;

type TwoVariadicElements<
  A extends unknown[],
  B extends unknown[],
  Result extends [...A, boolean, ...B, boolean],
> = Result;

type MatchingTwoVariadicElements<
  A extends unknown[],
  B extends unknown[],
> = TwoVariadicElements<A, B, [...A, boolean, ...B, boolean]>;

type ExcessMiddleElement<
  A extends unknown[],
  B extends unknown[],
> = TwoVariadicElements<A, B, [...A, boolean, boolean, ...B, boolean]>;

type FixedPrefixAndSuffix<
  T extends unknown[],
  Result extends [boolean, ...T, boolean],
> = Result;

type MatchingFixedPrefixAndSuffix<T extends unknown[]> = FixedPrefixAndSuffix<
  T,
  [boolean, ...T, boolean]
>;

type MissingFixedPrefix<T extends unknown[]> = FixedPrefixAndSuffix<
  T,
  [...T, boolean]
>;
