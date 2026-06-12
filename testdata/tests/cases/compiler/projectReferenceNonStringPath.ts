// @filename: /src/tsconfig.json
{
    "references": [{ "path": true }, { "path": "./other", "circular": "yes" }]
}

// @filename: /src/index.ts
export const x = 1;
