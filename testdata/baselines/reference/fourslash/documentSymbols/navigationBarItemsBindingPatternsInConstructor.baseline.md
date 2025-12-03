// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsBindingPatternsInConstructor.ts
class A {
    x: any
    constructor([a]: any) {
    }
}
class B {
    x: any;
    constructor( {a} = { a: 1 }) {
    }
}
```

# Symbols

```json
[
	{
		"name": "A",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 4,
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
						"character": 10
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 1,
						"character": 5
					}
				},
				"children": []
			},
			{
				"name": "constructor",
				"kind": "Constructor",
				"range": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 3,
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
						"character": 4
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "B",
		"kind": "Class",
		"range": {
			"start": {
				"line": 5,
				"character": 0
			},
			"end": {
				"line": 9,
				"character": 1
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
		"children": [
			{
				"name": "x",
				"kind": "Property",
				"range": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 6,
						"character": 11
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 6,
						"character": 5
					}
				},
				"children": []
			},
			{
				"name": "constructor",
				"kind": "Constructor",
				"range": {
					"start": {
						"line": 7,
						"character": 4
					},
					"end": {
						"line": 8,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 7,
						"character": 4
					},
					"end": {
						"line": 7,
						"character": 4
					}
				},
				"children": []
			}
		]
	}
]
```