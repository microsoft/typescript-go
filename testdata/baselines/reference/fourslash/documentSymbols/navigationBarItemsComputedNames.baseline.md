// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsComputedNames.ts
const enum E {
	A = 'A',
}
const a = '';

class C {
    [a]() {
        return 1;
    }

    [E.A]() {
        return 1;
    }

    [1]() {
        return 1;
    },

    ["foo"]() {
        return 1;
    },
}
```

# Symbols

```json
[
	{
		"name": "E",
		"kind": "Enum",
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
				"character": 11
			},
			"end": {
				"line": 0,
				"character": 12
			}
		},
		"children": [
			{
				"name": "A",
				"kind": "EnumMember",
				"range": {
					"start": {
						"line": 1,
						"character": 1
					},
					"end": {
						"line": 1,
						"character": 8
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 1
					},
					"end": {
						"line": 1,
						"character": 2
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
				"line": 3,
				"character": 6
			},
			"end": {
				"line": 3,
				"character": 12
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
		"children": []
	},
	{
		"name": "C",
		"kind": "Class",
		"range": {
			"start": {
				"line": 5,
				"character": 0
			},
			"end": {
				"line": 21,
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
				"name": "[a]",
				"kind": "Method",
				"range": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 8,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 6,
						"character": 7
					}
				},
				"children": []
			},
			{
				"name": "[E.A]",
				"kind": "Method",
				"range": {
					"start": {
						"line": 10,
						"character": 4
					},
					"end": {
						"line": 12,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 10,
						"character": 4
					},
					"end": {
						"line": 10,
						"character": 9
					}
				},
				"children": []
			},
			{
				"name": "1",
				"kind": "Method",
				"range": {
					"start": {
						"line": 14,
						"character": 4
					},
					"end": {
						"line": 16,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 14,
						"character": 4
					},
					"end": {
						"line": 14,
						"character": 7
					}
				},
				"children": []
			},
			{
				"name": "\"foo\"",
				"kind": "Method",
				"range": {
					"start": {
						"line": 18,
						"character": 4
					},
					"end": {
						"line": 20,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 18,
						"character": 4
					},
					"end": {
						"line": 18,
						"character": 11
					}
				},
				"children": []
			}
		]
	}
]
```