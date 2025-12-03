// === Document Symbols ===
```ts
// @FileName: /file1.ts
module a {
    function foo() {}
}
module b {
    function foo() {}
}
module a {
    function bar() {}
}
```

# Symbols

```json
[
	{
		"name": "a",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 2,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 7
			},
			"end": {
				"line": 0,
				"character": 8
			}
		},
		"children": [
			{
				"name": "foo",
				"kind": "Function",
				"range": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 1,
						"character": 21
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 13
					},
					"end": {
						"line": 1,
						"character": 16
					}
				},
				"children": []
			},
			{
				"name": "bar",
				"kind": "Function",
				"range": {
					"start": {
						"line": 7,
						"character": 4
					},
					"end": {
						"line": 7,
						"character": 21
					}
				},
				"selectionRange": {
					"start": {
						"line": 7,
						"character": 13
					},
					"end": {
						"line": 7,
						"character": 16
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "b",
		"kind": "Namespace",
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
				"character": 7
			},
			"end": {
				"line": 3,
				"character": 8
			}
		},
		"children": [
			{
				"name": "foo",
				"kind": "Function",
				"range": {
					"start": {
						"line": 4,
						"character": 4
					},
					"end": {
						"line": 4,
						"character": 21
					}
				},
				"selectionRange": {
					"start": {
						"line": 4,
						"character": 13
					},
					"end": {
						"line": 4,
						"character": 16
					}
				},
				"children": []
			}
		]
	}
]
```



// === Document Symbols ===
```ts
// @FileName: /file2.ts
module a {}
function a() {}
```

# Symbols

```json
[
	{
		"name": "a",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 11
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 7
			},
			"end": {
				"line": 0,
				"character": 8
			}
		},
		"children": []
	},
	{
		"name": "a",
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
		"children": []
	}
]
```



// === Document Symbols ===
```ts
// @FileName: /file3.ts
module a {
    interface A {
        foo: number;
    }
}
module a {
    interface A {
        bar: number;
    }
}
```

# Symbols

```json
[
	{
		"name": "a",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 4,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 7
			},
			"end": {
				"line": 0,
				"character": 8
			}
		},
		"children": [
			{
				"name": "A",
				"kind": "Interface",
				"range": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 3,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 14
					},
					"end": {
						"line": 1,
						"character": 15
					}
				},
				"children": [
					{
						"name": "foo",
						"kind": "Property",
						"range": {
							"start": {
								"line": 2,
								"character": 8
							},
							"end": {
								"line": 2,
								"character": 20
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 8
							},
							"end": {
								"line": 2,
								"character": 11
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "A",
				"kind": "Interface",
				"range": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 8,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 14
					},
					"end": {
						"line": 6,
						"character": 15
					}
				},
				"children": [
					{
						"name": "bar",
						"kind": "Property",
						"range": {
							"start": {
								"line": 7,
								"character": 8
							},
							"end": {
								"line": 7,
								"character": 20
							}
						},
						"selectionRange": {
							"start": {
								"line": 7,
								"character": 8
							},
							"end": {
								"line": 7,
								"character": 11
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



// === Document Symbols ===
```ts
// @FileName: /file4.ts
module A { export var x; }
module A.B { export var y; }
```

# Symbols

```json
[
	{
		"name": "A",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 26
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 7
			},
			"end": {
				"line": 0,
				"character": 8
			}
		},
		"children": [
			{
				"name": "x",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 0,
						"character": 22
					},
					"end": {
						"line": 0,
						"character": 23
					}
				},
				"selectionRange": {
					"start": {
						"line": 0,
						"character": 22
					},
					"end": {
						"line": 0,
						"character": 23
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "A.B",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 1,
				"character": 0
			},
			"end": {
				"line": 1,
				"character": 28
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 7
			},
			"end": {
				"line": 1,
				"character": 10
			}
		},
		"children": [
			{
				"name": "y",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 1,
						"character": 24
					},
					"end": {
						"line": 1,
						"character": 25
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 24
					},
					"end": {
						"line": 1,
						"character": 25
					}
				},
				"children": []
			}
		]
	}
]
```