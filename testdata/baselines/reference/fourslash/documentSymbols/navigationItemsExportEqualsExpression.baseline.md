// === Document Symbols ===
```ts
// @FileName: /navigationItemsExportEqualsExpression.ts
export = function () {}
export = function () {
    return class Foo {
    }
}

export = () => ""
export = () => {
    return class Foo {
    }
}

export = function f1() {}
export = function f2() {
    return class Foo {
    }
}

const abc = 12;
export = abc;
export = class AB {}
export = {
    a: 1,
    b: 1,
    c: {
        d: 1
    }
}
```

# Symbols

```json
[
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 23
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 23
			}
		},
		"children": [
			{
				"name": "export=",
				"kind": "Function",
				"range": {
					"start": {
						"line": 0,
						"character": 9
					},
					"end": {
						"line": 0,
						"character": 23
					}
				},
				"selectionRange": {
					"start": {
						"line": 0,
						"character": 9
					},
					"end": {
						"line": 0,
						"character": 9
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 1,
				"character": 0
			},
			"end": {
				"line": 4,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 0
			},
			"end": {
				"line": 4,
				"character": 1
			}
		},
		"children": [
			{
				"name": "export=",
				"kind": "Function",
				"range": {
					"start": {
						"line": 1,
						"character": 9
					},
					"end": {
						"line": 4,
						"character": 1
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 9
					},
					"end": {
						"line": 1,
						"character": 9
					}
				},
				"children": [
					{
						"name": "Foo",
						"kind": "Class",
						"range": {
							"start": {
								"line": 2,
								"character": 11
							},
							"end": {
								"line": 3,
								"character": 5
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 17
							},
							"end": {
								"line": 2,
								"character": 20
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 6,
				"character": 0
			},
			"end": {
				"line": 6,
				"character": 17
			}
		},
		"selectionRange": {
			"start": {
				"line": 6,
				"character": 0
			},
			"end": {
				"line": 6,
				"character": 17
			}
		},
		"children": [
			{
				"name": "export=",
				"kind": "Function",
				"range": {
					"start": {
						"line": 6,
						"character": 9
					},
					"end": {
						"line": 6,
						"character": 17
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 9
					},
					"end": {
						"line": 6,
						"character": 9
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 7,
				"character": 0
			},
			"end": {
				"line": 10,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 7,
				"character": 0
			},
			"end": {
				"line": 10,
				"character": 1
			}
		},
		"children": [
			{
				"name": "export=",
				"kind": "Function",
				"range": {
					"start": {
						"line": 7,
						"character": 9
					},
					"end": {
						"line": 10,
						"character": 1
					}
				},
				"selectionRange": {
					"start": {
						"line": 7,
						"character": 9
					},
					"end": {
						"line": 7,
						"character": 9
					}
				},
				"children": [
					{
						"name": "Foo",
						"kind": "Class",
						"range": {
							"start": {
								"line": 8,
								"character": 11
							},
							"end": {
								"line": 9,
								"character": 5
							}
						},
						"selectionRange": {
							"start": {
								"line": 8,
								"character": 17
							},
							"end": {
								"line": 8,
								"character": 20
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 12,
				"character": 0
			},
			"end": {
				"line": 12,
				"character": 25
			}
		},
		"selectionRange": {
			"start": {
				"line": 12,
				"character": 0
			},
			"end": {
				"line": 12,
				"character": 25
			}
		},
		"children": [
			{
				"name": "f1",
				"kind": "Function",
				"range": {
					"start": {
						"line": 12,
						"character": 9
					},
					"end": {
						"line": 12,
						"character": 25
					}
				},
				"selectionRange": {
					"start": {
						"line": 12,
						"character": 18
					},
					"end": {
						"line": 12,
						"character": 20
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 13,
				"character": 0
			},
			"end": {
				"line": 16,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 13,
				"character": 0
			},
			"end": {
				"line": 16,
				"character": 1
			}
		},
		"children": [
			{
				"name": "f2",
				"kind": "Function",
				"range": {
					"start": {
						"line": 13,
						"character": 9
					},
					"end": {
						"line": 16,
						"character": 1
					}
				},
				"selectionRange": {
					"start": {
						"line": 13,
						"character": 18
					},
					"end": {
						"line": 13,
						"character": 20
					}
				},
				"children": [
					{
						"name": "Foo",
						"kind": "Class",
						"range": {
							"start": {
								"line": 14,
								"character": 11
							},
							"end": {
								"line": 15,
								"character": 5
							}
						},
						"selectionRange": {
							"start": {
								"line": 14,
								"character": 17
							},
							"end": {
								"line": 14,
								"character": 20
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "abc",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 18,
				"character": 6
			},
			"end": {
				"line": 18,
				"character": 14
			}
		},
		"selectionRange": {
			"start": {
				"line": 18,
				"character": 6
			},
			"end": {
				"line": 18,
				"character": 9
			}
		},
		"children": []
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 19,
				"character": 0
			},
			"end": {
				"line": 19,
				"character": 13
			}
		},
		"selectionRange": {
			"start": {
				"line": 19,
				"character": 9
			},
			"end": {
				"line": 19,
				"character": 12
			}
		},
		"children": []
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 20,
				"character": 0
			},
			"end": {
				"line": 20,
				"character": 20
			}
		},
		"selectionRange": {
			"start": {
				"line": 20,
				"character": 0
			},
			"end": {
				"line": 20,
				"character": 20
			}
		},
		"children": [
			{
				"name": "AB",
				"kind": "Class",
				"range": {
					"start": {
						"line": 20,
						"character": 9
					},
					"end": {
						"line": 20,
						"character": 20
					}
				},
				"selectionRange": {
					"start": {
						"line": 20,
						"character": 15
					},
					"end": {
						"line": 20,
						"character": 17
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "export=",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 21,
				"character": 0
			},
			"end": {
				"line": 27,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 21,
				"character": 0
			},
			"end": {
				"line": 27,
				"character": 1
			}
		},
		"children": [
			{
				"name": "a",
				"kind": "Property",
				"range": {
					"start": {
						"line": 22,
						"character": 4
					},
					"end": {
						"line": 22,
						"character": 8
					}
				},
				"selectionRange": {
					"start": {
						"line": 22,
						"character": 4
					},
					"end": {
						"line": 22,
						"character": 5
					}
				},
				"children": []
			},
			{
				"name": "b",
				"kind": "Property",
				"range": {
					"start": {
						"line": 23,
						"character": 4
					},
					"end": {
						"line": 23,
						"character": 8
					}
				},
				"selectionRange": {
					"start": {
						"line": 23,
						"character": 4
					},
					"end": {
						"line": 23,
						"character": 5
					}
				},
				"children": []
			},
			{
				"name": "c",
				"kind": "Property",
				"range": {
					"start": {
						"line": 24,
						"character": 4
					},
					"end": {
						"line": 26,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 24,
						"character": 4
					},
					"end": {
						"line": 24,
						"character": 5
					}
				},
				"children": [
					{
						"name": "d",
						"kind": "Property",
						"range": {
							"start": {
								"line": 25,
								"character": 8
							},
							"end": {
								"line": 25,
								"character": 12
							}
						},
						"selectionRange": {
							"start": {
								"line": 25,
								"character": 8
							},
							"end": {
								"line": 25,
								"character": 9
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