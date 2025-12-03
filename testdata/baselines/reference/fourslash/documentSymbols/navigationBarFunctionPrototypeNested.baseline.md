// === Document Symbols ===
```js
// @FileName: /foo.js
function A() {}
A.B = function () {  } 
A.B.prototype.d = function () {  }  
Object.defineProperty(A.B.prototype, "x", {
    get() {}
})
A.prototype.D = function () {  } 
A.prototype.D.prototype.d = function () {  } 
```

# Symbols

```json
[
	{
		"name": "A",
		"kind": "Function",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 15
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
				"name": "B",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 1,
						"character": 0
					},
					"end": {
						"line": 1,
						"character": 3
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 2
					},
					"end": {
						"line": 1,
						"character": 3
					}
				},
				"children": [
					{
						"name": "B",
						"kind": "Function",
						"range": {
							"start": {
								"line": 1,
								"character": 6
							},
							"end": {
								"line": 1,
								"character": 22
							}
						},
						"selectionRange": {
							"start": {
								"line": 1,
								"character": 6
							},
							"end": {
								"line": 1,
								"character": 6
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "D",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 6,
						"character": 0
					},
					"end": {
						"line": 6,
						"character": 13
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 12
					},
					"end": {
						"line": 6,
						"character": 13
					}
				},
				"children": [
					{
						"name": "D",
						"kind": "Function",
						"range": {
							"start": {
								"line": 6,
								"character": 16
							},
							"end": {
								"line": 6,
								"character": 32
							}
						},
						"selectionRange": {
							"start": {
								"line": 6,
								"character": 16
							},
							"end": {
								"line": 6,
								"character": 16
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "d",
		"kind": "Function",
		"range": {
			"start": {
				"line": 2,
				"character": 18
			},
			"end": {
				"line": 2,
				"character": 34
			}
		},
		"selectionRange": {
			"start": {
				"line": 2,
				"character": 18
			},
			"end": {
				"line": 2,
				"character": 18
			}
		},
		"children": []
	},
	{
		"name": "get",
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
	},
	{
		"name": "d",
		"kind": "Function",
		"range": {
			"start": {
				"line": 7,
				"character": 28
			},
			"end": {
				"line": 7,
				"character": 44
			}
		},
		"selectionRange": {
			"start": {
				"line": 7,
				"character": 28
			},
			"end": {
				"line": 7,
				"character": 28
			}
		},
		"children": []
	}
]
```