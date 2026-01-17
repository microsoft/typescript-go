// @strict: true
// @filename: lib.ts
interface Options {
  verbose?: boolean;
}

export const doSomething = (input: string, options?: Options): string => {
  return input;
};

// @filename: index.ts
// @declaration: true
// @outdir: dist
import { doSomething } from "./lib";

// Using Parameters<typeof fn>[1] to extract the options type
export const wrapper = (
  input: string,
  options?: Parameters<typeof doSomething>[1],
) => doSomething(input, options);
