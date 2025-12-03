// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsMultilineStringIdentifiers1.ts
declare module "Multiline\r\nMadness" {
}

declare module "Multiline\
Madness" {
}
declare module "MultilineMadness" {}

declare module "Multiline\
Madness2" {
}

interface Foo {
    "a1\\\r\nb";
    "a2\
    \
    b"(): Foo;
}

class Bar implements Foo {
    'a1\\\r\nb': Foo;

    'a2\
    \
    b'(): Foo {
        return this;
    }
}
```

# Symbols

```json
[
	{
		"name": "\"Multiline\\r\\nMadness\"",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 1,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 15
			},
			"end": {
				"line": 0,
				"character": 37
			}
		},
		"children": []
	},
	{
		"name": "\"MultilineMadness\"",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 3,
				"character": 0
			},
			"end": {
				"line": 5,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 3,
				"character": 15
			},
			"end": {
				"line": 4,
				"character": 8
			}
		},
		"children": []
	},
	{
		"name": "\"MultilineMadness2\"",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 8,
				"character": 0
			},
			"end": {
				"line": 10,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 8,
				"character": 15
			},
			"end": {
				"line": 9,
				"character": 9
			}
		},
		"children": []
	},
	{
		"name": "Foo",
		"kind": "Interface",
		"range": {
			"start": {
				"line": 12,
				"character": 0
			},
			"end": {
				"line": 17,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 12,
				"character": 10
			},
			"end": {
				"line": 12,
				"character": 13
			}
		},
		"children": [
			{
				"name": "\"a1\\\\\\r\\nb\"",
				"kind": "Property",
				"range": {
					"start": {
						"line": 13,
						"character": 4
					},
					"end": {
						"line": 13,
						"character": 16
					}
				},
				"selectionRange": {
					"start": {
						"line": 13,
						"character": 4
					},
					"end": {
						"line": 13,
						"character": 15
					}
				},
				"children": []
			},
			{
				"name": "\"a2        b\"",
				"kind": "Method",
				"range": {
					"start": {
						"line": 14,
						"character": 4
					},
					"end": {
						"line": 16,
						"character": 14
					}
				},
				"selectionRange": {
					"start": {
						"line": 14,
						"character": 4
					},
					"end": {
						"line": 16,
						"character": 6
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "Bar",
		"kind": "Class",
		"range": {
			"start": {
				"line": 19,
				"character": 0
			},
			"end": {
				"line": 27,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 19,
				"character": 6
			},
			"end": {
				"line": 19,
				"character": 9
			}
		},
		"children": [
			{
				"name": "\"a1\\\\\\r\\nb\"",
				"kind": "Property",
				"range": {
					"start": {
						"line": 20,
						"character": 4
					},
					"end": {
						"line": 20,
						"character": 21
					}
				},
				"selectionRange": {
					"start": {
						"line": 20,
						"character": 4
					},
					"end": {
						"line": 20,
						"character": 15
					}
				},
				"children": []
			},
			{
				"name": "\"a2        b\"",
				"kind": "Method",
				"range": {
					"start": {
						"line": 22,
						"character": 4
					},
					"end": {
						"line": 26,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 22,
						"character": 4
					},
					"end": {
						"line": 24,
						"character": 6
					}
				},
				"children": []
			}
		]
	}
]
```