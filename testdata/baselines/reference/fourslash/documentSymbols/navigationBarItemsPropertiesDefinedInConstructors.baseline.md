// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsPropertiesDefinedInConstructors.ts
class List<T> {
    constructor(public a: boolean, private b: T, readonly c: string, d: number) {
        var local = 0;
    }
}
```

# Symbols

```json
[
	{
		"name": "List",
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
				"character": 10
			}
		},
		"children": [
			{
				"name": "constructor",
				"kind": "Constructor",
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
						"character": 4
					}
				},
				"children": [
					{
						"name": "local",
						"kind": "Variable",
						"range": {
							"start": {
								"line": 2,
								"character": 12
							},
							"end": {
								"line": 2,
								"character": 21
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 12
							},
							"end": {
								"line": 2,
								"character": 17
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "a",
				"kind": "Property",
				"range": {
					"start": {
						"line": 1,
						"character": 16
					},
					"end": {
						"line": 1,
						"character": 33
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 23
					},
					"end": {
						"line": 1,
						"character": 24
					}
				},
				"children": []
			},
			{
				"name": "b",
				"kind": "Property",
				"range": {
					"start": {
						"line": 1,
						"character": 35
					},
					"end": {
						"line": 1,
						"character": 47
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 43
					},
					"end": {
						"line": 1,
						"character": 44
					}
				},
				"children": []
			},
			{
				"name": "c",
				"kind": "Property",
				"range": {
					"start": {
						"line": 1,
						"character": 49
					},
					"end": {
						"line": 1,
						"character": 67
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 58
					},
					"end": {
						"line": 1,
						"character": 59
					}
				},
				"children": []
			}
		]
	}
]
```