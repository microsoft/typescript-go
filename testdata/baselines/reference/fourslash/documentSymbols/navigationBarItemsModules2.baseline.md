// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsModules2.ts
namespace Test.A { }

namespace Test.B {
    class Foo { }
}
```

# Symbols

```json
[
	{
		"name": "Test.A",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 20
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 10
			},
			"end": {
				"line": 0,
				"character": 16
			}
		},
		"children": []
	},
	{
		"name": "Test.B",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 2,
				"character": 0
			},
			"end": {
				"line": 4,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 2,
				"character": 10
			},
			"end": {
				"line": 2,
				"character": 16
			}
		},
		"children": [
			{
				"name": "Foo",
				"kind": "Class",
				"range": {
					"start": {
						"line": 3,
						"character": 4
					},
					"end": {
						"line": 3,
						"character": 17
					}
				},
				"selectionRange": {
					"start": {
						"line": 3,
						"character": 10
					},
					"end": {
						"line": 3,
						"character": 13
					}
				},
				"children": []
			}
		]
	}
]
```