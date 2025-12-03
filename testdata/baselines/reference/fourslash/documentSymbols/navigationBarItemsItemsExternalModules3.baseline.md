// === Document Symbols ===
```ts
// @FileName: /test/my fil	e.ts
export class Bar {
    public s: string;
}
export var x: number;
```

# Symbols

```json
[
	{
		"name": "Bar",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 2,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 13
			},
			"end": {
				"line": 0,
				"character": 16
			}
		},
		"children": [
			{
				"name": "s",
				"kind": "Property",
				"range": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 1,
						"character": 21
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 11
					},
					"end": {
						"line": 1,
						"character": 12
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "x",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 3,
				"character": 11
			},
			"end": {
				"line": 3,
				"character": 20
			}
		},
		"selectionRange": {
			"start": {
				"line": 3,
				"character": 11
			},
			"end": {
				"line": 3,
				"character": 12
			}
		},
		"children": []
	}
]
```