// === Document Symbols ===
```ts
// @FileName: /navigationBarAnonymousClassAndFunctionExpressions2.ts
console.log(console.log(class Y {}, class X {}), console.log(class B {}, class A {}));
console.log(class Cls { meth() {} });
```

# Symbols

```json
[
	{
		"name": "Y",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 24
			},
			"end": {
				"line": 0,
				"character": 34
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 30
			},
			"end": {
				"line": 0,
				"character": 31
			}
		},
		"children": []
	},
	{
		"name": "X",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 36
			},
			"end": {
				"line": 0,
				"character": 46
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 42
			},
			"end": {
				"line": 0,
				"character": 43
			}
		},
		"children": []
	},
	{
		"name": "B",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 61
			},
			"end": {
				"line": 0,
				"character": 71
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 67
			},
			"end": {
				"line": 0,
				"character": 68
			}
		},
		"children": []
	},
	{
		"name": "A",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 73
			},
			"end": {
				"line": 0,
				"character": 83
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 79
			},
			"end": {
				"line": 0,
				"character": 80
			}
		},
		"children": []
	},
	{
		"name": "Cls",
		"kind": "Class",
		"range": {
			"start": {
				"line": 1,
				"character": 12
			},
			"end": {
				"line": 1,
				"character": 35
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 18
			},
			"end": {
				"line": 1,
				"character": 21
			}
		},
		"children": [
			{
				"name": "meth",
				"kind": "Method",
				"range": {
					"start": {
						"line": 1,
						"character": 24
					},
					"end": {
						"line": 1,
						"character": 33
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 24
					},
					"end": {
						"line": 1,
						"character": 28
					}
				},
				"children": []
			}
		]
	}
]
```