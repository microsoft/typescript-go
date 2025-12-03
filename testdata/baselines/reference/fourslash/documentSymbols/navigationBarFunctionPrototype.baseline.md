// === Document Symbols ===
```js
// @FileName: /foo.js
function f() {}
f.prototype.x = 0;
f.y = 0;
f.prototype.method = function () {};
Object.defineProperty(f, 'staticProp', { 
    set: function() {}, 
    get: function(){
    } 
});
Object.defineProperty(f.prototype, 'name', { 
    set: function() {}, 
    get: function(){
    } 
}); 
```

# Symbols

```json
[
	{
		"name": "f",
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
				"name": "x",
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
				"children": []
			},
			{
				"name": "y",
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
				"children": []
			},
			{
				"name": "method",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 3,
						"character": 0
					},
					"end": {
						"line": 3,
						"character": 18
					}
				},
				"selectionRange": {
					"start": {
						"line": 3,
						"character": 12
					},
					"end": {
						"line": 3,
						"character": 18
					}
				},
				"children": [
					{
						"name": "method",
						"kind": "Function",
						"range": {
							"start": {
								"line": 3,
								"character": 21
							},
							"end": {
								"line": 3,
								"character": 35
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 21
							},
							"end": {
								"line": 3,
								"character": 21
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "set",
		"kind": "Property",
		"range": {
			"start": {
				"line": 5,
				"character": 4
			},
			"end": {
				"line": 5,
				"character": 22
			}
		},
		"selectionRange": {
			"start": {
				"line": 5,
				"character": 4
			},
			"end": {
				"line": 5,
				"character": 7
			}
		},
		"children": []
	},
	{
		"name": "get",
		"kind": "Property",
		"range": {
			"start": {
				"line": 6,
				"character": 4
			},
			"end": {
				"line": 7,
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
				"character": 7
			}
		},
		"children": []
	},
	{
		"name": "set",
		"kind": "Property",
		"range": {
			"start": {
				"line": 10,
				"character": 4
			},
			"end": {
				"line": 10,
				"character": 22
			}
		},
		"selectionRange": {
			"start": {
				"line": 10,
				"character": 4
			},
			"end": {
				"line": 10,
				"character": 7
			}
		},
		"children": []
	},
	{
		"name": "get",
		"kind": "Property",
		"range": {
			"start": {
				"line": 11,
				"character": 4
			},
			"end": {
				"line": 12,
				"character": 5
			}
		},
		"selectionRange": {
			"start": {
				"line": 11,
				"character": 4
			},
			"end": {
				"line": 11,
				"character": 7
			}
		},
		"children": []
	}
]
```