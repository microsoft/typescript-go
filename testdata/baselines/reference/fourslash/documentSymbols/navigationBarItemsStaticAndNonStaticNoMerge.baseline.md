// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsStaticAndNonStaticNoMerge.ts
class C {
    static x;
    x;
}
```

# Symbols

```json
[
	{
		"name": "C",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 3,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 6
			},
			"end": {
				"line": 0,
				"character": 7
			}
		},
		"children": [
			{
				"name": "x",
				"kind": "Property",
				"range": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 1,
						"character": 13
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
			},
			{
				"name": "x",
				"kind": "Property",
				"range": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 2,
						"character": 6
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 2,
						"character": 5
					}
				},
				"children": []
			}
		]
	}
]
```