// === Document Symbols ===
```js
// @FileName: /foo.js
A.prototype.a = function() { };
A.prototype.b = function() { };
function A() {}
```

# Symbols

```json
[
	{
		"name": "A",
		"kind": "Function",
		"range": {
			"start": {
				"line": 2,
				"character": 0
			},
			"end": {
				"line": 2,
				"character": 15
			}
		},
		"selectionRange": {
			"start": {
				"line": 2,
				"character": 9
			},
			"end": {
				"line": 2,
				"character": 10
			}
		},
		"children": [
			{
				"name": "a",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 0,
						"character": 0
					},
					"end": {
						"line": 0,
						"character": 13
					}
				},
				"selectionRange": {
					"start": {
						"line": 0,
						"character": 12
					},
					"end": {
						"line": 0,
						"character": 13
					}
				},
				"children": [
					{
						"name": "a",
						"kind": "Function",
						"range": {
							"start": {
								"line": 0,
								"character": 16
							},
							"end": {
								"line": 0,
								"character": 30
							}
						},
						"selectionRange": {
							"start": {
								"line": 0,
								"character": 16
							},
							"end": {
								"line": 0,
								"character": 16
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "b",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 1,
						"character": 0
					},
					"end": {
						"line": 1,
						"character": 13
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 12
					},
					"end": {
						"line": 1,
						"character": 13
					}
				},
				"children": [
					{
						"name": "b",
						"kind": "Function",
						"range": {
							"start": {
								"line": 1,
								"character": 16
							},
							"end": {
								"line": 1,
								"character": 30
							}
						},
						"selectionRange": {
							"start": {
								"line": 1,
								"character": 16
							},
							"end": {
								"line": 1,
								"character": 16
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