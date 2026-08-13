// @loadExternalPlugins: true
// @noTypesAndSymbols: true

// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "target": "es2020",
        "module": "esnext",
        "moduleResolution": "bundler",
        "strict": true
    },
    "contentMappers": [
        { "package": "mapper", "extensions": [".invalidrange"] },
        { "package": "mapper", "extensions": [".invalidpolicy"] },
        { "package": "mapper", "extensions": [".expectmissing"] },
        { "package": "mapper", "extensions": [".overlap"] }
    ]
}

// @Filename: /node_modules/mapper/package.json
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": {
        "exec": ["compiler-test-mapper"],
        "extensions": {
            ".invalidrange": ".ts",
            ".invalidpolicy": ".ts",
            ".expectmissing": ".ts",
            ".overlap": ".ts"
        },
        "compilerOptions": ["target", "jsx"]
    }
}

// @Filename: /invalidRange.invalidrange
// @box-invalid-directive: invalid-range

// @Filename: /invalidPolicy.invalidpolicy
// @box-invalid-directive: invalid-policy

// @Filename: /expectWithoutUnusedDiagnostic.expectmissing
// @box-invalid-directive: expect-without-unused-diagnostic

// @Filename: /overlap.overlap
// @box-invalid-directive: overlap