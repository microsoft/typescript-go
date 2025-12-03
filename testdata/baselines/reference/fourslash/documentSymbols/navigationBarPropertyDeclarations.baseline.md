// === Document Symbols ===
```ts
// @FileName: /navigationBarPropertyDeclarations.ts
class A {
    public A1 = class {
        public x = 1;
        private y() {}
        protected z() {}
    }

    public A2 = {
        x: 1,
        y() {},
        z() {}
    }

    public A3 = function () {}
    public A4 = () => {}
    public A5 = 1;
    public A6 = "A6";

    public ["A7"] = class {
        public x = 1;
        private y() {}
        protected z() {}
    }

    public [1] = {
        x: 1,
        y() {},
        z() {}
    }

    public [1 + 1] = 1;
}
```

# Symbols

```json
[
	{
		"name": "A",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 31,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 6
			},
			"end": {
				"line": 0,
				"character": 7
			}
		},
		"children": [
			{
				"name": "A1",
				"kind": "Property",
				"range": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 5,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 11
					},
					"end": {
						"line": 1,
						"character": 13
					}
				},
				"children": [
					{
						"name": "x",
						"kind": "Property",
						"range": {
							"start": {
								"line": 2,
								"character": 8
							},
							"end": {
								"line": 2,
								"character": 21
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 15
							},
							"end": {
								"line": 2,
								"character": 16
							}
						},
						"children": []
					},
					{
						"name": "y",
						"kind": "Method",
						"range": {
							"start": {
								"line": 3,
								"character": 8
							},
							"end": {
								"line": 3,
								"character": 22
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 16
							},
							"end": {
								"line": 3,
								"character": 17
							}
						},
						"children": []
					},
					{
						"name": "z",
						"kind": "Method",
						"range": {
							"start": {
								"line": 4,
								"character": 8
							},
							"end": {
								"line": 4,
								"character": 24
							}
						},
						"selectionRange": {
							"start": {
								"line": 4,
								"character": 18
							},
							"end": {
								"line": 4,
								"character": 19
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "A2",
				"kind": "Property",
				"range": {
					"start": {
						"line": 7,
						"character": 4
					},
					"end": {
						"line": 11,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 7,
						"character": 11
					},
					"end": {
						"line": 7,
						"character": 13
					}
				},
				"children": [
					{
						"name": "x",
						"kind": "Property",
						"range": {
							"start": {
								"line": 8,
								"character": 8
							},
							"end": {
								"line": 8,
								"character": 12
							}
						},
						"selectionRange": {
							"start": {
								"line": 8,
								"character": 8
							},
							"end": {
								"line": 8,
								"character": 9
							}
						},
						"children": []
					},
					{
						"name": "y",
						"kind": "Method",
						"range": {
							"start": {
								"line": 9,
								"character": 8
							},
							"end": {
								"line": 9,
								"character": 14
							}
						},
						"selectionRange": {
							"start": {
								"line": 9,
								"character": 8
							},
							"end": {
								"line": 9,
								"character": 9
							}
						},
						"children": []
					},
					{
						"name": "z",
						"kind": "Method",
						"range": {
							"start": {
								"line": 10,
								"character": 8
							},
							"end": {
								"line": 10,
								"character": 14
							}
						},
						"selectionRange": {
							"start": {
								"line": 10,
								"character": 8
							},
							"end": {
								"line": 10,
								"character": 9
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "A3",
				"kind": "Property",
				"range": {
					"start": {
						"line": 13,
						"character": 4
					},
					"end": {
						"line": 13,
						"character": 30
					}
				},
				"selectionRange": {
					"start": {
						"line": 13,
						"character": 11
					},
					"end": {
						"line": 13,
						"character": 13
					}
				},
				"children": []
			},
			{
				"name": "A4",
				"kind": "Property",
				"range": {
					"start": {
						"line": 14,
						"character": 4
					},
					"end": {
						"line": 14,
						"character": 24
					}
				},
				"selectionRange": {
					"start": {
						"line": 14,
						"character": 11
					},
					"end": {
						"line": 14,
						"character": 13
					}
				},
				"children": []
			},
			{
				"name": "A5",
				"kind": "Property",
				"range": {
					"start": {
						"line": 15,
						"character": 4
					},
					"end": {
						"line": 15,
						"character": 18
					}
				},
				"selectionRange": {
					"start": {
						"line": 15,
						"character": 11
					},
					"end": {
						"line": 15,
						"character": 13
					}
				},
				"children": []
			},
			{
				"name": "A6",
				"kind": "Property",
				"range": {
					"start": {
						"line": 16,
						"character": 4
					},
					"end": {
						"line": 16,
						"character": 21
					}
				},
				"selectionRange": {
					"start": {
						"line": 16,
						"character": 11
					},
					"end": {
						"line": 16,
						"character": 13
					}
				},
				"children": []
			},
			{
				"name": "\"A7\"",
				"kind": "Property",
				"range": {
					"start": {
						"line": 18,
						"character": 4
					},
					"end": {
						"line": 22,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 18,
						"character": 11
					},
					"end": {
						"line": 18,
						"character": 17
					}
				},
				"children": [
					{
						"name": "x",
						"kind": "Property",
						"range": {
							"start": {
								"line": 19,
								"character": 8
							},
							"end": {
								"line": 19,
								"character": 21
							}
						},
						"selectionRange": {
							"start": {
								"line": 19,
								"character": 15
							},
							"end": {
								"line": 19,
								"character": 16
							}
						},
						"children": []
					},
					{
						"name": "y",
						"kind": "Method",
						"range": {
							"start": {
								"line": 20,
								"character": 8
							},
							"end": {
								"line": 20,
								"character": 22
							}
						},
						"selectionRange": {
							"start": {
								"line": 20,
								"character": 16
							},
							"end": {
								"line": 20,
								"character": 17
							}
						},
						"children": []
					},
					{
						"name": "z",
						"kind": "Method",
						"range": {
							"start": {
								"line": 21,
								"character": 8
							},
							"end": {
								"line": 21,
								"character": 24
							}
						},
						"selectionRange": {
							"start": {
								"line": 21,
								"character": 18
							},
							"end": {
								"line": 21,
								"character": 19
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "1",
				"kind": "Property",
				"range": {
					"start": {
						"line": 24,
						"character": 4
					},
					"end": {
						"line": 28,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 24,
						"character": 11
					},
					"end": {
						"line": 24,
						"character": 14
					}
				},
				"children": [
					{
						"name": "x",
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
					},
					{
						"name": "y",
						"kind": "Method",
						"range": {
							"start": {
								"line": 26,
								"character": 8
							},
							"end": {
								"line": 26,
								"character": 14
							}
						},
						"selectionRange": {
							"start": {
								"line": 26,
								"character": 8
							},
							"end": {
								"line": 26,
								"character": 9
							}
						},
						"children": []
					},
					{
						"name": "z",
						"kind": "Method",
						"range": {
							"start": {
								"line": 27,
								"character": 8
							},
							"end": {
								"line": 27,
								"character": 14
							}
						},
						"selectionRange": {
							"start": {
								"line": 27,
								"character": 8
							},
							"end": {
								"line": 27,
								"character": 9
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "[1 + 1]",
				"kind": "Property",
				"range": {
					"start": {
						"line": 30,
						"character": 4
					},
					"end": {
						"line": 30,
						"character": 23
					}
				},
				"selectionRange": {
					"start": {
						"line": 30,
						"character": 11
					},
					"end": {
						"line": 30,
						"character": 18
					}
				},
				"children": []
			}
		]
	}
]
```