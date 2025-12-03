// === Document Symbols ===
```js
// @FileName: /file.js
const _sym = Symbol("_sym");
class MyClass {
    constructor() {
        // Dynamic assignment properties can't show up in navigation,
        // as they're not syntactic members
        // Additonally, late bound members are always filtered out, besides
        this[_sym] = "ok";
    }

    method() {
        this[_sym] = "yep";
        const x = this[_sym];
    }
}
```

# Symbols

```json
[
	{
		"name": "_sym",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 0,
				"character": 6
			},
			"end": {
				"line": 0,
				"character": 27
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 6
			},
			"end": {
				"line": 0,
				"character": 10
			}
		},
		"children": []
	},
	{
		"name": "MyClass",
		"kind": "Class",
		"range": {
			"start": {
				"line": 1,
				"character": 0
			},
			"end": {
				"line": 13,
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
				"character": 13
			}
		},
		"children": [
			{
				"name": "constructor",
				"kind": "Constructor",
				"range": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 7,
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
			},
			{
				"name": "method",
				"kind": "Method",
				"range": {
					"start": {
						"line": 9,
						"character": 4
					},
					"end": {
						"line": 12,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 9,
						"character": 4
					},
					"end": {
						"line": 9,
						"character": 10
					}
				},
				"children": [
					{
						"name": "x",
						"kind": "Variable",
						"range": {
							"start": {
								"line": 11,
								"character": 14
							},
							"end": {
								"line": 11,
								"character": 28
							}
						},
						"selectionRange": {
							"start": {
								"line": 11,
								"character": 14
							},
							"end": {
								"line": 11,
								"character": 15
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