// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsMultilineStringIdentifiers2.ts
function f(p1: () => any, p2: string) { }
f(() => { }, `line1\
line2\
line3`);

class c1 {
    const a = ' ''line1\
        line2';
}

f(() => { }, `unterminated backtick 1
unterminated backtick 2
unterminated backtick 3
```

# Symbols

```json
[
	{
		"name": "f",
		"kind": "Function",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 41
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
		"name": "f() callback",
		"kind": "Function",
		"range": {
			"start": {
				"line": 1,
				"character": 2
			},
			"end": {
				"line": 1,
				"character": 11
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 2
			},
			"end": {
				"line": 1,
				"character": 2
			}
		},
		"children": []
	},
	{
		"name": "c1",
		"kind": "Class",
		"range": {
			"start": {
				"line": 5,
				"character": 0
			},
			"end": {
				"line": 8,
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
				"character": 8
			}
		},
		"children": [
			{
				"name": "a",
				"kind": "Property",
				"range": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 6,
						"character": 17
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 10
					},
					"end": {
						"line": 6,
						"character": 11
					}
				},
				"children": []
			},
			{
				"name": "\"line1        line2\"",
				"kind": "Property",
				"range": {
					"start": {
						"line": 6,
						"character": 17
					},
					"end": {
						"line": 7,
						"character": 15
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 17
					},
					"end": {
						"line": 7,
						"character": 14
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "f() callback",
		"kind": "Function",
		"range": {
			"start": {
				"line": 10,
				"character": 2
			},
			"end": {
				"line": 10,
				"character": 11
			}
		},
		"selectionRange": {
			"start": {
				"line": 10,
				"character": 2
			},
			"end": {
				"line": 10,
				"character": 2
			}
		},
		"children": []
	}
]
```