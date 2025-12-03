// === Document Symbols ===
```js
// @FileName: /foo.js
var b = 1;
function A() {}; 
A.prototype.a = function() { };
A.b = function() { };
b = 2
/* Comment */
A.prototype.c = function() { }
var b = 2
A.prototype.d = function() { }
```

# Symbols

```json
[
	{
		"name": "b",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 0,
				"character": 4
			},
			"end": {
				"line": 0,
				"character": 9
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
		"children": []
	},
	{
		"name": "A",
		"kind": "Function",
		"range": {
			"start": {
				"line": 1,
				"character": 0
			},
			"end": {
				"line": 1,
				"character": 15
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 9
			},
			"end": {
				"line": 1,
				"character": 10
			}
		},
		"children": [
			{
				"name": "a",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 2,
						"character": 0
					},
					"end": {
						"line": 2,
						"character": 13
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 12
					},
					"end": {
						"line": 2,
						"character": 13
					}
				},
				"children": [
					{
						"name": "a",
						"kind": "Function",
						"range": {
							"start": {
								"line": 2,
								"character": 16
							},
							"end": {
								"line": 2,
								"character": 30
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 16
							},
							"end": {
								"line": 2,
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
						"line": 3,
						"character": 0
					},
					"end": {
						"line": 3,
						"character": 3
					}
				},
				"selectionRange": {
					"start": {
						"line": 3,
						"character": 2
					},
					"end": {
						"line": 3,
						"character": 3
					}
				},
				"children": [
					{
						"name": "b",
						"kind": "Function",
						"range": {
							"start": {
								"line": 3,
								"character": 6
							},
							"end": {
								"line": 3,
								"character": 20
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 6
							},
							"end": {
								"line": 3,
								"character": 6
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "c",
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
						"name": "c",
						"kind": "Function",
						"range": {
							"start": {
								"line": 6,
								"character": 16
							},
							"end": {
								"line": 6,
								"character": 30
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
			},
			{
				"name": "d",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 8,
						"character": 0
					},
					"end": {
						"line": 8,
						"character": 13
					}
				},
				"selectionRange": {
					"start": {
						"line": 8,
						"character": 12
					},
					"end": {
						"line": 8,
						"character": 13
					}
				},
				"children": [
					{
						"name": "d",
						"kind": "Function",
						"range": {
							"start": {
								"line": 8,
								"character": 16
							},
							"end": {
								"line": 8,
								"character": 30
							}
						},
						"selectionRange": {
							"start": {
								"line": 8,
								"character": 16
							},
							"end": {
								"line": 8,
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
		"name": "b",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 7,
				"character": 4
			},
			"end": {
				"line": 7,
				"character": 9
			}
		},
		"selectionRange": {
			"start": {
				"line": 7,
				"character": 4
			},
			"end": {
				"line": 7,
				"character": 5
			}
		},
		"children": []
	}
]
```