// === Document Symbols ===
```ts
// @FileName: /navigationBarAnonymousClassAndFunctionExpressions.ts
global.cls = class { };
(function() {
    const x = () => {
        // Presence of inner function causes x to be a top-level function.
        function xx() {}
    };
    const y = {
        // This is not a top-level function (contains nothing, but shows up in childItems of its parent.)
        foo: function() {}
    };
    (function nest() {
        function moreNest() {}
    })();
})();
(function() { // Different anonymous functions are not merged
    // These will only show up as childItems.
    function z() {}
    console.log(function() {})
    describe("this", 'function', `is a function`, `with template literal ${"a"}`, () => {});
    [].map(() => {});
})
(function classes() {
    // Classes show up in top-level regardless of whether they have names or inner declarations.
    const cls2 = class { };
    console.log(class cls3 {});
    (class { });
})
```

# Symbols

```json
[
	{
		"name": "cls",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 13
			},
			"end": {
				"line": 0,
				"character": 22
			}
		},
		"selectionRange": {
			"start": {
				"line": 0,
				"character": 13
			},
			"end": {
				"line": 0,
				"character": 13
			}
		},
		"children": []
	},
	{
		"name": "<function>",
		"kind": "Function",
		"range": {
			"start": {
				"line": 1,
				"character": 1
			},
			"end": {
				"line": 13,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 1
			},
			"end": {
				"line": 1,
				"character": 1
			}
		},
		"children": [
			{
				"name": "x",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 2,
						"character": 10
					},
					"end": {
						"line": 5,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 10
					},
					"end": {
						"line": 2,
						"character": 11
					}
				},
				"children": [
					{
						"name": "xx",
						"kind": "Function",
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
								"character": 17
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
				"name": "y",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 6,
						"character": 10
					},
					"end": {
						"line": 9,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 10
					},
					"end": {
						"line": 6,
						"character": 11
					}
				},
				"children": [
					{
						"name": "foo",
						"kind": "Property",
						"range": {
							"start": {
								"line": 8,
								"character": 8
							},
							"end": {
								"line": 8,
								"character": 26
							}
						},
						"selectionRange": {
							"start": {
								"line": 8,
								"character": 8
							},
							"end": {
								"line": 8,
								"character": 11
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "nest",
				"kind": "Function",
				"range": {
					"start": {
						"line": 10,
						"character": 5
					},
					"end": {
						"line": 12,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 10,
						"character": 14
					},
					"end": {
						"line": 10,
						"character": 18
					}
				},
				"children": [
					{
						"name": "moreNest",
						"kind": "Function",
						"range": {
							"start": {
								"line": 11,
								"character": 8
							},
							"end": {
								"line": 11,
								"character": 30
							}
						},
						"selectionRange": {
							"start": {
								"line": 11,
								"character": 17
							},
							"end": {
								"line": 11,
								"character": 25
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "<function>",
		"kind": "Function",
		"range": {
			"start": {
				"line": 14,
				"character": 1
			},
			"end": {
				"line": 20,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 14,
				"character": 1
			},
			"end": {
				"line": 14,
				"character": 1
			}
		},
		"children": [
			{
				"name": "z",
				"kind": "Function",
				"range": {
					"start": {
						"line": 16,
						"character": 4
					},
					"end": {
						"line": 16,
						"character": 19
					}
				},
				"selectionRange": {
					"start": {
						"line": 16,
						"character": 13
					},
					"end": {
						"line": 16,
						"character": 14
					}
				},
				"children": []
			},
			{
				"name": "console.log() callback",
				"kind": "Function",
				"range": {
					"start": {
						"line": 17,
						"character": 16
					},
					"end": {
						"line": 17,
						"character": 29
					}
				},
				"selectionRange": {
					"start": {
						"line": 17,
						"character": 16
					},
					"end": {
						"line": 17,
						"character": 16
					}
				},
				"children": []
			},
			{
				"name": "describe() callback",
				"kind": "Function",
				"range": {
					"start": {
						"line": 18,
						"character": 82
					},
					"end": {
						"line": 18,
						"character": 90
					}
				},
				"selectionRange": {
					"start": {
						"line": 18,
						"character": 82
					},
					"end": {
						"line": 18,
						"character": 82
					}
				},
				"children": []
			},
			{
				"name": "map() callback",
				"kind": "Function",
				"range": {
					"start": {
						"line": 19,
						"character": 11
					},
					"end": {
						"line": 19,
						"character": 19
					}
				},
				"selectionRange": {
					"start": {
						"line": 19,
						"character": 11
					},
					"end": {
						"line": 19,
						"character": 11
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "classes",
		"kind": "Function",
		"range": {
			"start": {
				"line": 21,
				"character": 1
			},
			"end": {
				"line": 26,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 21,
				"character": 10
			},
			"end": {
				"line": 21,
				"character": 17
			}
		},
		"children": [
			{
				"name": "cls2",
				"kind": "Variable",
				"range": {
					"start": {
						"line": 23,
						"character": 10
					},
					"end": {
						"line": 23,
						"character": 26
					}
				},
				"selectionRange": {
					"start": {
						"line": 23,
						"character": 10
					},
					"end": {
						"line": 23,
						"character": 14
					}
				},
				"children": []
			},
			{
				"name": "cls3",
				"kind": "Class",
				"range": {
					"start": {
						"line": 24,
						"character": 16
					},
					"end": {
						"line": 24,
						"character": 29
					}
				},
				"selectionRange": {
					"start": {
						"line": 24,
						"character": 22
					},
					"end": {
						"line": 24,
						"character": 26
					}
				},
				"children": []
			},
			{
				"name": "<class>",
				"kind": "Class",
				"range": {
					"start": {
						"line": 25,
						"character": 5
					},
					"end": {
						"line": 25,
						"character": 14
					}
				},
				"selectionRange": {
					"start": {
						"line": 25,
						"character": 5
					},
					"end": {
						"line": 25,
						"character": 5
					}
				},
				"children": []
			}
		]
	}
]
```