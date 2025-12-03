// === Document Symbols ===
```ts
// @FileName: /navigationBarGetterAndSetter.ts
class X {
    get x() {}
    set x(value) {
        // Inner declaration should make the setter top-level.
        function f() {}
    }
}
```

# Symbols

```json
[
	{
		"name": "X",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 6,
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
						"character": 14
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 8
					},
					"end": {
						"line": 1,
						"character": 9
					}
				},
				"children": []
			},
			{
				"name": "x",
				"kind": "Property",
				"range": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 5,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 8
					},
					"end": {
						"line": 2,
						"character": 9
					}
				},
				"children": [
					{
						"name": "f",
						"kind": "Function",
						"range": {
							"start": {
								"line": 4,
								"character": 8
							},
							"end": {
								"line": 4,
								"character": 23
							}
						},
						"selectionRange": {
							"start": {
								"line": 4,
								"character": 17
							},
							"end": {
								"line": 4,
								"character": 18
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