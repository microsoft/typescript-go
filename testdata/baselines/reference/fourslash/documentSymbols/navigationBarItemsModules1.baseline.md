// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsModules1.ts
declare module "X.Y.Z" {}

declare module 'X2.Y2.Z2' {}

declare module "foo";

module A.B.C {
    export var x;
}

module A.B {
    export var y;
}

module A {
    export var z;
}

module A {
    module B {
        module C {
            declare var x;
        }
    }
}
```

# Symbols

```json
[
	{
		"name": "\"X.Y.Z\"",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 0,
				"character": 25
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 15
			},
			"end": {
				"line": 0,
				"character": 22
			}
		},
		"children": []
	},
	{
		"name": "\"X2.Y2.Z2\"",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 2,
				"character": 0
			},
			"end": {
				"line": 2,
				"character": 28
			}
		},
		"selectionRange": {
			"start": {
				"line": 2,
				"character": 15
			},
			"end": {
				"line": 2,
				"character": 25
			}
		},
		"children": []
	},
	{
		"name": "\"foo\"",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 4,
				"character": 0
			},
			"end": {
				"line": 4,
				"character": 21
			}
		},
		"selectionRange": {
			"start": {
				"line": 4,
				"character": 15
			},
			"end": {
				"line": 4,
				"character": 20
			}
		},
		"children": []
	},
	{
		"name": "A.B.C",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 6,
				"character": 0
			},
			"end": {
				"line": 8,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 6,
				"character": 7
			},
			"end": {
				"line": 6,
				"character": 12
			}
		},
		"children": [
			{
				"name": "x",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 7,
						"character": 15
					},
					"end": {
						"line": 7,
						"character": 16
					}
				},
				"selectionRange": {
					"start": {
						"line": 7,
						"character": 15
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
		"name": "A.B",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 10,
				"character": 0
			},
			"end": {
				"line": 12,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 10,
				"character": 7
			},
			"end": {
				"line": 10,
				"character": 10
			}
		},
		"children": [
			{
				"name": "y",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 11,
						"character": 15
					},
					"end": {
						"line": 11,
						"character": 16
					}
				},
				"selectionRange": {
					"start": {
						"line": 11,
						"character": 15
					},
					"end": {
						"line": 11,
						"character": 16
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "A",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 14,
				"character": 0
			},
			"end": {
				"line": 16,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 14,
				"character": 7
			},
			"end": {
				"line": 14,
				"character": 8
			}
		},
		"children": [
			{
				"name": "z",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 15,
						"character": 15
					},
					"end": {
						"line": 15,
						"character": 16
					}
				},
				"selectionRange": {
					"start": {
						"line": 15,
						"character": 15
					},
					"end": {
						"line": 15,
						"character": 16
					}
				},
				"children": []
			},
			{
				"name": "B",
				"kind": "Namespace",
				"range": {
					"start": {
						"line": 19,
						"character": 4
					},
					"end": {
						"line": 23,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 19,
						"character": 11
					},
					"end": {
						"line": 19,
						"character": 12
					}
				},
				"children": [
					{
						"name": "C",
						"kind": "Namespace",
						"range": {
							"start": {
								"line": 20,
								"character": 8
							},
							"end": {
								"line": 22,
								"character": 9
							}
						},
						"selectionRange": {
							"start": {
								"line": 20,
								"character": 15
							},
							"end": {
								"line": 20,
								"character": 16
							}
						},
						"children": [
							{
								"name": "x",
								"kind": "Variable",
								"range": {
									"start": {
										"line": 21,
										"character": 24
									},
									"end": {
										"line": 21,
										"character": 25
									}
								},
								"selectionRange": {
									"start": {
										"line": 21,
										"character": 24
									},
									"end": {
										"line": 21,
										"character": 25
									}
								},
								"children": []
							}
						]
					}
				]
			}
		]
	}
]
```