// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsExports.ts
export { a } from "a";

export { b as B } from "a" 

export import e = require("a");

export * from "a"; // no bindings here
```

# Symbols

```json
[
	{
		"name": "a",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 0,
				"character": 9
			},
			"end": {
				"line": 0,
				"character": 10
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 9
			},
			"end": {
				"line": 0,
				"character": 10
			}
		},
		"children": []
	},
	{
		"name": "B",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 2,
				"character": 9
			},
			"end": {
				"line": 2,
				"character": 15
			}
		},
		"selectionRange": {
			"start": {
				"line": 2,
				"character": 14
			},
			"end": {
				"line": 2,
				"character": 15
			}
		},
		"children": []
	},
	{
		"name": "e",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 4,
				"character": 0
			},
			"end": {
				"line": 4,
				"character": 31
			}
		},
		"selectionRange": {
			"start": {
				"line": 4,
				"character": 14
			},
			"end": {
				"line": 4,
				"character": 15
			}
		},
		"children": []
	}
]
```