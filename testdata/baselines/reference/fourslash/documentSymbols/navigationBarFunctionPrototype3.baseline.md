// === Document Symbols ===
```js
// @FileName: /foo.js
var A; 
A.prototype.a = function() { };
A.b = function() { };
```

# Symbols

```json
[
	{
		"name": "A",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 0,
				"character": 4
			},
			"end": {
				"line": 0,
				"character": 5
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
				"name": "a",
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
						"name": "a",
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
			},
			{
				"name": "b",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 2,
						"character": 0
					},
					"end": {
						"line": 2,
						"character": 3
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 2
					},
					"end": {
						"line": 2,
						"character": 3
					}
				},
				"children": [
					{
						"name": "b",
						"kind": "Function",
						"range": {
							"start": {
								"line": 2,
								"character": 6
							},
							"end": {
								"line": 2,
								"character": 20
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 6
							},
							"end": {
								"line": 2,
								"character": 6
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