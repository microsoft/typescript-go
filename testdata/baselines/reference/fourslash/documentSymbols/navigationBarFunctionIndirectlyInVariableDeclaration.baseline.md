// === Document Symbols ===
```ts
// @FileName: /navigationBarFunctionIndirectlyInVariableDeclaration.ts
var a = {
    propA: function() {
        var c;
    }
};
var b;
b = {
    propB: function() {
    // function must not have an empty body to appear top level
        var d;
    }
};
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
				"character": 4
			},
			"end": {
				"line": 4,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 4
			},
			"end": {
				"line": 0,
				"character": 5
			}
		},
		"children": [
			{
				"name": "propA",
				"kind": "Property",
				"range": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 3,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 1,
						"character": 9
					}
				},
				"children": [
					{
						"name": "c",
						"kind": "Variable",
						"range": {
							"start": {
								"line": 2,
								"character": 12
							},
							"end": {
								"line": 2,
								"character": 13
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 12
							},
							"end": {
								"line": 2,
								"character": 13
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "b",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 5,
				"character": 4
			},
			"end": {
				"line": 5,
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
		"children": []
	},
	{
		"name": "propB",
		"kind": "Property",
		"range": {
			"start": {
				"line": 7,
				"character": 4
			},
			"end": {
				"line": 10,
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
				"character": 9
			}
		},
		"children": [
			{
				"name": "d",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 9,
						"character": 12
					},
					"end": {
						"line": 9,
						"character": 13
					}
				},
				"selectionRange": {
					"start": {
						"line": 9,
						"character": 12
					},
					"end": {
						"line": 9,
						"character": 13
					}
				},
				"children": []
			}
		]
	}
]
```