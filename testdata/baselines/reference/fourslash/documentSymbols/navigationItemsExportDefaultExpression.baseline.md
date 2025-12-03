// === Document Symbols ===
```ts
// @FileName: /navigationItemsExportDefaultExpression.ts
export default function () {}
export default function () {
    return class Foo {
    }
}

export default () => ""
export default () => {
    return class Foo {
    }
}

export default function f1() {}
export default function f2() {
    return class Foo {
    }
}

const abc = 12;
export default abc;
export default class AB {}
export default {
    a: 1,
    b: 1,
    c: {
        d: 1
    }
}

function foo(props: { x: number; y: number }) {}
export default foo({ x: 1, y: 1 });
```

# Symbols

```json
[
	{
		"name": "default",
		"kind": "Function",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 29
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 0
			}
		},
		"children": []
	},
	{
		"name": "default",
		"kind": "Function",
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
				"line": 1,
				"character": 0
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
	},
	{
		"name": "default",
		"kind": "Function",
		"range": {
			"start": {
				"line": 6,
				"character": 15
			},
			"end": {
				"line": 6,
				"character": 23
			}
		},
		"selectionRange": {
			"start": {
				"line": 6,
				"character": 15
			},
			"end": {
				"line": 6,
				"character": 15
			}
		},
		"children": []
	},
	{
		"name": "default",
		"kind": "Function",
		"range": {
			"start": {
				"line": 7,
				"character": 15
			},
			"end": {
				"line": 10,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 7,
				"character": 15
			},
			"end": {
				"line": 7,
				"character": 15
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
	},
	{
		"name": "f1",
		"kind": "Function",
		"range": {
			"start": {
				"line": 12,
				"character": 0
			},
			"end": {
				"line": 12,
				"character": 31
			}
		},
		"selectionRange": {
			"start": {
				"line": 12,
				"character": 24
			},
			"end": {
				"line": 12,
				"character": 26
			}
		},
		"children": []
	},
	{
		"name": "f2",
		"kind": "Function",
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
				"character": 24
			},
			"end": {
				"line": 13,
				"character": 26
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
		"name": "AB",
		"kind": "Class",
		"range": {
			"start": {
				"line": 20,
				"character": 0
			},
			"end": {
				"line": 20,
				"character": 26
			}
		},
		"selectionRange": {
			"start": {
				"line": 20,
				"character": 21
			},
			"end": {
				"line": 20,
				"character": 23
			}
		},
		"children": []
	},
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
	},
	{
		"name": "foo",
		"kind": "Function",
		"range": {
			"start": {
				"line": 29,
				"character": 0
			},
			"end": {
				"line": 29,
				"character": 48
			}
		},
		"selectionRange": {
			"start": {
				"line": 29,
				"character": 9
			},
			"end": {
				"line": 29,
				"character": 12
			}
		},
		"children": []
	},
	{
		"name": "x",
		"kind": "Property",
		"range": {
			"start": {
				"line": 30,
				"character": 21
			},
			"end": {
				"line": 30,
				"character": 25
			}
		},
		"selectionRange": {
			"start": {
				"line": 30,
				"character": 21
			},
			"end": {
				"line": 30,
				"character": 22
			}
		},
		"children": []
	},
	{
		"name": "y",
		"kind": "Property",
		"range": {
			"start": {
				"line": 30,
				"character": 27
			},
			"end": {
				"line": 30,
				"character": 31
			}
		},
		"selectionRange": {
			"start": {
				"line": 30,
				"character": 27
			},
			"end": {
				"line": 30,
				"character": 28
			}
		},
		"children": []
	}
]
```