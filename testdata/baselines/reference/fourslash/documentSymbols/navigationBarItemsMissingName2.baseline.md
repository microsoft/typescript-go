// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsMissingName2.ts
/**
 * This is a class.
 */
class /* But it has no name! */ {
    foo() {}
}
```

# Symbols

```json
[
	{
		"name": "<class>",
		"kind": "Class",
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
				"character": 0
			},
			"end": {
				"line": 3,
				"character": 0
			}
		},
		"children": [
			{
				"name": "foo",
				"kind": "Method",
				"range": {
					"start": {
						"line": 4,
						"character": 4
					},
					"end": {
						"line": 4,
						"character": 12
					}
				},
				"selectionRange": {
					"start": {
						"line": 4,
						"character": 4
					},
					"end": {
						"line": 4,
						"character": 7
					}
				},
				"children": []
			}
		]
	}
]
```