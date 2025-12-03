// === Document Symbols ===
```js
// @FileName: /foo.js
var A; 
A.prototype = { };
A.prototype = { m() {} };
A.prototype.a = function() { };
A.b = function() { };

var B; 
B["prototype"] = { };
B["prototype"] = { m() {} };
B["prototype"]["a"] = function() { };
B["b"] = function() { };
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
						"line": 3,
						"character": 0
					},
					"end": {
						"line": 3,
						"character": 13
					}
				},
				"selectionRange": {
					"start": {
						"line": 3,
						"character": 12
					},
					"end": {
						"line": 3,
						"character": 13
					}
				},
				"children": [
					{
						"name": "a",
						"kind": "Function",
						"range": {
							"start": {
								"line": 3,
								"character": 16
							},
							"end": {
								"line": 3,
								"character": 30
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 16
							},
							"end": {
								"line": 3,
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
						"line": 4,
						"character": 0
					},
					"end": {
						"line": 4,
						"character": 3
					}
				},
				"selectionRange": {
					"start": {
						"line": 4,
						"character": 2
					},
					"end": {
						"line": 4,
						"character": 3
					}
				},
				"children": [
					{
						"name": "b",
						"kind": "Function",
						"range": {
							"start": {
								"line": 4,
								"character": 6
							},
							"end": {
								"line": 4,
								"character": 20
							}
						},
						"selectionRange": {
							"start": {
								"line": 4,
								"character": 6
							},
							"end": {
								"line": 4,
								"character": 6
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "m",
		"kind": "Method",
		"range": {
			"start": {
				"line": 2,
				"character": 16
			},
			"end": {
				"line": 2,
				"character": 22
			}
		},
		"selectionRange": {
			"start": {
				"line": 2,
				"character": 16
			},
			"end": {
				"line": 2,
				"character": 17
			}
		},
		"children": []
	},
	{
		"name": "B",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 6,
				"character": 4
			},
			"end": {
				"line": 6,
				"character": 5
			}
		},
		"selectionRange": {
			"start": {
				"line": 6,
				"character": 4
			},
			"end": {
				"line": 6,
				"character": 5
			}
		},
		"children": [
			{
				"name": "\"a\"",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 9,
						"character": 0
					},
					"end": {
						"line": 9,
						"character": 19
					}
				},
				"selectionRange": {
					"start": {
						"line": 9,
						"character": 15
					},
					"end": {
						"line": 9,
						"character": 18
					}
				},
				"children": [
					{
						"name": "\"a\"",
						"kind": "Function",
						"range": {
							"start": {
								"line": 9,
								"character": 22
							},
							"end": {
								"line": 9,
								"character": 36
							}
						},
						"selectionRange": {
							"start": {
								"line": 9,
								"character": 22
							},
							"end": {
								"line": 9,
								"character": 22
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "\"b\"",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 10,
						"character": 0
					},
					"end": {
						"line": 10,
						"character": 6
					}
				},
				"selectionRange": {
					"start": {
						"line": 10,
						"character": 2
					},
					"end": {
						"line": 10,
						"character": 5
					}
				},
				"children": [
					{
						"name": "\"b\"",
						"kind": "Function",
						"range": {
							"start": {
								"line": 10,
								"character": 9
							},
							"end": {
								"line": 10,
								"character": 23
							}
						},
						"selectionRange": {
							"start": {
								"line": 10,
								"character": 9
							},
							"end": {
								"line": 10,
								"character": 9
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "m",
		"kind": "Method",
		"range": {
			"start": {
				"line": 8,
				"character": 19
			},
			"end": {
				"line": 8,
				"character": 25
			}
		},
		"selectionRange": {
			"start": {
				"line": 8,
				"character": 19
			},
			"end": {
				"line": 8,
				"character": 20
			}
		},
		"children": []
	}
]
```