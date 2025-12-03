// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsFunctions.ts
function foo() {
    var x = 10;
    function bar() {
        var y = 10;
        function biz() {
            var z = 10;
        }
        function qux() {
            // A function with an empty body should not be top level
        }
    }
}

function baz() {
    var v = 10;
}
```

# Symbols

```json
[
	{
		"name": "foo",
		"kind": "Function",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 11,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 9
			},
			"end": {
				"line": 0,
				"character": 12
			}
		},
		"children": [
			{
				"name": "x",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 1,
						"character": 8
					},
					"end": {
						"line": 1,
						"character": 14
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 8
					},
					"end": {
						"line": 1,
						"character": 9
					}
				},
				"children": []
			},
			{
				"name": "bar",
				"kind": "Function",
				"range": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 10,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 13
					},
					"end": {
						"line": 2,
						"character": 16
					}
				},
				"children": [
					{
						"name": "y",
						"kind": "Variable",
						"range": {
							"start": {
								"line": 3,
								"character": 12
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
								"character": 13
							}
						},
						"children": []
					},
					{
						"name": "biz",
						"kind": "Function",
						"range": {
							"start": {
								"line": 4,
								"character": 8
							},
							"end": {
								"line": 6,
								"character": 9
							}
						},
						"selectionRange": {
							"start": {
								"line": 4,
								"character": 17
							},
							"end": {
								"line": 4,
								"character": 20
							}
						},
						"children": [
							{
								"name": "z",
								"kind": "Variable",
								"range": {
									"start": {
										"line": 5,
										"character": 16
									},
									"end": {
										"line": 5,
										"character": 22
									}
								},
								"selectionRange": {
									"start": {
										"line": 5,
										"character": 16
									},
									"end": {
										"line": 5,
										"character": 17
									}
								},
								"children": []
							}
						]
					},
					{
						"name": "qux",
						"kind": "Function",
						"range": {
							"start": {
								"line": 7,
								"character": 8
							},
							"end": {
								"line": 9,
								"character": 9
							}
						},
						"selectionRange": {
							"start": {
								"line": 7,
								"character": 17
							},
							"end": {
								"line": 7,
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
		"name": "baz",
		"kind": "Function",
		"range": {
			"start": {
				"line": 13,
				"character": 0
			},
			"end": {
				"line": 15,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 13,
				"character": 9
			},
			"end": {
				"line": 13,
				"character": 12
			}
		},
		"children": [
			{
				"name": "v",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 14,
						"character": 8
					},
					"end": {
						"line": 14,
						"character": 14
					}
				},
				"selectionRange": {
					"start": {
						"line": 14,
						"character": 8
					},
					"end": {
						"line": 14,
						"character": 9
					}
				},
				"children": []
			}
		]
	}
]
```