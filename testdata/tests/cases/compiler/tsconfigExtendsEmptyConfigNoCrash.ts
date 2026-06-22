// https://github.com/microsoft/typescript-go/issues/4265
// @noEmit: true

// @filename: tsconfig.json
{
  "extends": "./base.json",
  "files": ["./main.ts"]
}

// @filename: base.json

// @filename: main.ts
export const x = 1;
