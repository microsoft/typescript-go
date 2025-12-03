// === Document Symbols ===
```ts
// @FileName: /navigationBarMerging_grandchildren.ts
// Should not merge grandchildren with property assignments
const o = {
    a: {
        m() {},
    },
    b: {
        m() {},
    },
}
```

# Symbols

```json
[
	{
		"name": "o",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 1,
				"character": 6
			},
			"end": {
				"line": 8,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 6
			},
			"end": {
				"line": 1,
				"character": 7
			}
		},
		"children": [
			{
				"name": "a",
				"kind": "Property",
				"range": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 4,
						"character": 5
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
				"children": [
					{
						"name": "m",
						"kind": "Method",
						"range": {
							"start": {
								"line": 3,
								"character": 8
							},
							"end": {
								"line": 3,
								"character": 14
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 8
							},
							"end": {
								"line": 3,
								"character": 9
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "b",
				"kind": "Property",
				"range": {
					"start": {
						"line": 5,
						"character": 4
					},
					"end": {
						"line": 7,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 5,
						"character": 4
					},
					"end": {
						"line": 5,
						"character": 5
					}
				},
				"children": [
					{
						"name": "m",
						"kind": "Method",
						"range": {
							"start": {
								"line": 6,
								"character": 8
							},
							"end": {
								"line": 6,
								"character": 14
							}
						},
						"selectionRange": {
							"start": {
								"line": 6,
								"character": 8
							},
							"end": {
								"line": 6,
								"character": 9
							}
						},
						"children": []
					}
				]
			}
		]
	}
]
```