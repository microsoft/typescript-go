// === Document Symbols ===
```ts
// @FileName: /navbar_contains_no_duplicates.ts
declare module Windows {
    export module Foundation {
        export var A;
        export class Test {
            public wow();
        }
    }
}

declare module Windows {
    export module Foundation {
        export var B;
        export module Test {
            export function Boom(): number;
        }
    }
}

class ABC {
    public foo() {
        return 3;
    }
}

module ABC {
    export var x = 3;
}
```

# Symbols

```json
[
	{
		"name": "Windows",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 7,
				"character": 1
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
		"children": [
			{
				"name": "Foundation",
				"kind": "Namespace",
				"range": {
					"start": {
						"line": 1,
						"character": 4
					},
					"end": {
						"line": 6,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 18
					},
					"end": {
						"line": 1,
						"character": 28
					}
				},
				"children": [
					{
						"name": "A",
						"kind": "Variable",
						"range": {
							"start": {
								"line": 2,
								"character": 19
							},
							"end": {
								"line": 2,
								"character": 20
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 19
							},
							"end": {
								"line": 2,
								"character": 20
							}
						},
						"children": []
					},
					{
						"name": "Test",
						"kind": "Class",
						"range": {
							"start": {
								"line": 3,
								"character": 8
							},
							"end": {
								"line": 5,
								"character": 9
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 21
							},
							"end": {
								"line": 3,
								"character": 25
							}
						},
						"children": [
							{
								"name": "wow",
								"kind": "Method",
								"range": {
									"start": {
										"line": 4,
										"character": 12
									},
									"end": {
										"line": 4,
										"character": 25
									}
								},
								"selectionRange": {
									"start": {
										"line": 4,
										"character": 19
									},
									"end": {
										"line": 4,
										"character": 22
									}
								},
								"children": []
							}
						]
					},
					{
						"name": "B",
						"kind": "Variable",
						"range": {
							"start": {
								"line": 11,
								"character": 19
							},
							"end": {
								"line": 11,
								"character": 20
							}
						},
						"selectionRange": {
							"start": {
								"line": 11,
								"character": 19
							},
							"end": {
								"line": 11,
								"character": 20
							}
						},
						"children": []
					},
					{
						"name": "Test",
						"kind": "Namespace",
						"range": {
							"start": {
								"line": 12,
								"character": 8
							},
							"end": {
								"line": 14,
								"character": 9
							}
						},
						"selectionRange": {
							"start": {
								"line": 12,
								"character": 22
							},
							"end": {
								"line": 12,
								"character": 26
							}
						},
						"children": [
							{
								"name": "Boom",
								"kind": "Function",
								"range": {
									"start": {
										"line": 13,
										"character": 12
									},
									"end": {
										"line": 13,
										"character": 43
									}
								},
								"selectionRange": {
									"start": {
										"line": 13,
										"character": 28
									},
									"end": {
										"line": 13,
										"character": 32
									}
								},
								"children": []
							}
						]
					}
				]
			}
		]
	},
	{
		"name": "ABC",
		"kind": "Class",
		"range": {
			"start": {
				"line": 18,
				"character": 0
			},
			"end": {
				"line": 22,
				"character": 1
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
		"children": [
			{
				"name": "foo",
				"kind": "Method",
				"range": {
					"start": {
						"line": 19,
						"character": 4
					},
					"end": {
						"line": 21,
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
						"character": 14
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "ABC",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 24,
				"character": 0
			},
			"end": {
				"line": 26,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 24,
				"character": 7
			},
			"end": {
				"line": 24,
				"character": 10
			}
		},
		"children": [
			{
				"name": "x",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 25,
						"character": 15
					},
					"end": {
						"line": 25,
						"character": 20
					}
				},
				"selectionRange": {
					"start": {
						"line": 25,
						"character": 15
					},
					"end": {
						"line": 25,
						"character": 16
					}
				},
				"children": []
			}
		]
	}
]
```