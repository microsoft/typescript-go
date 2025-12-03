// === Document Symbols ===
```ts
// @FileName: /navigationBarAnonymousClassAndFunctionExpressions3.ts
describe('foo', () => {
    test(`a ${1} b ${2}`, () => {})
})

const a = 1;
const b = 2;
describe('foo', () => {
    test(`a ${a} b {b}`, () => {})
})
```

# Symbols

```json
[
	{
		"name": "describe() callback",
		"kind": "Function",
		"range": {
			"start": {
				"line": 0,
				"character": 16
			},
			"end": {
				"line": 2,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 16
			},
			"end": {
				"line": 0,
				"character": 16
			}
		},
		"children": [
			{
				"name": "test() callback",
				"kind": "Function",
				"range": {
					"start": {
						"line": 1,
						"character": 26
					},
					"end": {
						"line": 1,
						"character": 34
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 26
					},
					"end": {
						"line": 1,
						"character": 26
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "a",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 4,
				"character": 6
			},
			"end": {
				"line": 4,
				"character": 11
			}
		},
		"selectionRange": {
			"start": {
				"line": 4,
				"character": 6
			},
			"end": {
				"line": 4,
				"character": 7
			}
		},
		"children": []
	},
	{
		"name": "b",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 5,
				"character": 6
			},
			"end": {
				"line": 5,
				"character": 11
			}
		},
		"selectionRange": {
			"start": {
				"line": 5,
				"character": 6
			},
			"end": {
				"line": 5,
				"character": 7
			}
		},
		"children": []
	},
	{
		"name": "describe() callback",
		"kind": "Function",
		"range": {
			"start": {
				"line": 6,
				"character": 16
			},
			"end": {
				"line": 8,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 6,
				"character": 16
			},
			"end": {
				"line": 6,
				"character": 16
			}
		},
		"children": [
			{
				"name": "test() callback",
				"kind": "Function",
				"range": {
					"start": {
						"line": 7,
						"character": 25
					},
					"end": {
						"line": 7,
						"character": 33
					}
				},
				"selectionRange": {
					"start": {
						"line": 7,
						"character": 25
					},
					"end": {
						"line": 7,
						"character": 25
					}
				},
				"children": []
			}
		]
	}
]
```