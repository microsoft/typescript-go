// === Document Symbols ===
```ts
// @FileName: /navigationBarWithLocalVariables.ts
function x(){
	const x = Object()
	x.foo = ""
}
```

# Symbols

```json
[
	{
		"name": "x",
		"kind": "Function",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 3,
				"character": 1
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
		"children": [
			{
				"name": "x",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 1,
						"character": 7
					},
					"end": {
						"line": 1,
						"character": 19
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 7
					},
					"end": {
						"line": 1,
						"character": 8
					}
				},
				"children": []
			}
		]
	}
]
```