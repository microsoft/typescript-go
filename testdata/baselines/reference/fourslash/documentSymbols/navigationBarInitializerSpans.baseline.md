// === Document Symbols ===
```ts
// @FileName: /navigationBarInitializerSpans.ts
// get the name for the navbar from the variable name rather than the function name
const x = () => { var a; };
const f = function f() { var b; };
const y = { z: function z() { var c; } };
```

# Symbols

```json
[
	{
		"name": "x",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 1,
				"character": 6
			},
			"end": {
				"line": 1,
				"character": 26
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
				"kind": "Variable",
				"range": {
					"start": {
						"line": 1,
						"character": 22
					},
					"end": {
						"line": 1,
						"character": 23
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 22
					},
					"end": {
						"line": 1,
						"character": 23
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "f",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 2,
				"character": 6
			},
			"end": {
				"line": 2,
				"character": 33
			}
		},
		"selectionRange": {
			"start": {
				"line": 2,
				"character": 6
			},
			"end": {
				"line": 2,
				"character": 7
			}
		},
		"children": [
			{
				"name": "b",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 2,
						"character": 29
					},
					"end": {
						"line": 2,
						"character": 30
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 29
					},
					"end": {
						"line": 2,
						"character": 30
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "y",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 3,
				"character": 6
			},
			"end": {
				"line": 3,
				"character": 40
			}
		},
		"selectionRange": {
			"start": {
				"line": 3,
				"character": 6
			},
			"end": {
				"line": 3,
				"character": 7
			}
		},
		"children": [
			{
				"name": "z",
				"kind": "Property",
				"range": {
					"start": {
						"line": 3,
						"character": 12
					},
					"end": {
						"line": 3,
						"character": 38
					}
				},
				"selectionRange": {
					"start": {
						"line": 3,
						"character": 12
					},
					"end": {
						"line": 3,
						"character": 13
					}
				},
				"children": [
					{
						"name": "c",
						"kind": "Variable",
						"range": {
							"start": {
								"line": 3,
								"character": 34
							},
							"end": {
								"line": 3,
								"character": 35
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 34
							},
							"end": {
								"line": 3,
								"character": 35
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