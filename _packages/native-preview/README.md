# TypeScript (Native Preview)

This package provides a preview build of [the native port of TypeScript](https://devblogs.microsoft.com/typescript/typescript-native-port/).
Not all features are implemented yet.

This package is intended for testing and experimentation.
It will eventually be replaced by the official TypeScript package.

## Usage

Use the `tsgo` command just like you would use `tsc`:

```sh
npx tsgo --help
```

You can also check the version or run against a project:

```sh
npx tsgo --version
npx tsgo -p tsconfig.json
```

## VS Code Integration

A preview VS Code extension is available on the Marketplace:

- https://marketplace.visualstudio.com/items?itemName=TypeScriptTeam.native-preview

To enable the preview language service in VS Code, add this setting:

```json
{
  "js/ts.experimental.useTsgo": true
}
```

## Issues and Feedback

The native port of TypeScript is still in progress.
We expect many gaps, but are seeking experimentation and feedback.
If you have found differences that are not yet known, we encourage you to leave feedback on [the issue tracker](https://github.com/microsoft/typescript-go/issues).
