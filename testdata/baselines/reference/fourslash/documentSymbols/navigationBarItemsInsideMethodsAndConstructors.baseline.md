// === Document Symbols ===
```ts
// @FileName: /navigationBarItemsInsideMethodsAndConstructors.ts
class Class {
    constructor() {
        function LocalFunctionInConstructor() {}
        interface LocalInterfaceInConstrcutor {}
        enum LocalEnumInConstructor { LocalEnumMemberInConstructor }
    }

    method() {
        function LocalFunctionInMethod() {
            function LocalFunctionInLocalFunctionInMethod() {}
        }
        interface LocalInterfaceInMethod {}
        enum LocalEnumInMethod { LocalEnumMemberInMethod }
    }

    emptyMethod() { } // Non child functions method should not be duplicated
}
```

# Symbols

```json
[
	{
		"name": "Class",
		"kind": "Class",
		"range": {
			"start": {
				"line": 0,
				"character": 0
			},
			"end": {
				"line": 16,
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
				"character": 11
			}
		},
		"children": [
			{
				"name": "constructor",
				"kind": "Constructor",
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
						"character": 4
					},
					"end": {
						"line": 1,
						"character": 4
					}
				},
				"children": [
					{
						"name": "LocalFunctionInConstructor",
						"kind": "Function",
						"range": {
							"start": {
								"line": 2,
								"character": 8
							},
							"end": {
								"line": 2,
								"character": 48
							}
						},
						"selectionRange": {
							"start": {
								"line": 2,
								"character": 17
							},
							"end": {
								"line": 2,
								"character": 43
							}
						},
						"children": []
					},
					{
						"name": "LocalInterfaceInConstrcutor",
						"kind": "Interface",
						"range": {
							"start": {
								"line": 3,
								"character": 8
							},
							"end": {
								"line": 3,
								"character": 48
							}
						},
						"selectionRange": {
							"start": {
								"line": 3,
								"character": 18
							},
							"end": {
								"line": 3,
								"character": 45
							}
						},
						"children": []
					},
					{
						"name": "LocalEnumInConstructor",
						"kind": "Enum",
						"range": {
							"start": {
								"line": 4,
								"character": 8
							},
							"end": {
								"line": 4,
								"character": 68
							}
						},
						"selectionRange": {
							"start": {
								"line": 4,
								"character": 13
							},
							"end": {
								"line": 4,
								"character": 35
							}
						},
						"children": [
							{
								"name": "LocalEnumMemberInConstructor",
								"kind": "EnumMember",
								"range": {
									"start": {
										"line": 4,
										"character": 38
									},
									"end": {
										"line": 4,
										"character": 66
									}
								},
								"selectionRange": {
									"start": {
										"line": 4,
										"character": 38
									},
									"end": {
										"line": 4,
										"character": 66
									}
								},
								"children": []
							}
						]
					}
				]
			},
			{
				"name": "method",
				"kind": "Method",
				"range": {
					"start": {
						"line": 7,
						"character": 4
					},
					"end": {
						"line": 13,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 7,
						"character": 4
					},
					"end": {
						"line": 7,
						"character": 10
					}
				},
				"children": [
					{
						"name": "LocalFunctionInMethod",
						"kind": "Function",
						"range": {
							"start": {
								"line": 8,
								"character": 8
							},
							"end": {
								"line": 10,
								"character": 9
							}
						},
						"selectionRange": {
							"start": {
								"line": 8,
								"character": 17
							},
							"end": {
								"line": 8,
								"character": 38
							}
						},
						"children": [
							{
								"name": "LocalFunctionInLocalFunctionInMethod",
								"kind": "Function",
								"range": {
									"start": {
										"line": 9,
										"character": 12
									},
									"end": {
										"line": 9,
										"character": 62
									}
								},
								"selectionRange": {
									"start": {
										"line": 9,
										"character": 21
									},
									"end": {
										"line": 9,
										"character": 57
									}
								},
								"children": []
							}
						]
					},
					{
						"name": "LocalInterfaceInMethod",
						"kind": "Interface",
						"range": {
							"start": {
								"line": 11,
								"character": 8
							},
							"end": {
								"line": 11,
								"character": 43
							}
						},
						"selectionRange": {
							"start": {
								"line": 11,
								"character": 18
							},
							"end": {
								"line": 11,
								"character": 40
							}
						},
						"children": []
					},
					{
						"name": "LocalEnumInMethod",
						"kind": "Enum",
						"range": {
							"start": {
								"line": 12,
								"character": 8
							},
							"end": {
								"line": 12,
								"character": 58
							}
						},
						"selectionRange": {
							"start": {
								"line": 12,
								"character": 13
							},
							"end": {
								"line": 12,
								"character": 30
							}
						},
						"children": [
							{
								"name": "LocalEnumMemberInMethod",
								"kind": "EnumMember",
								"range": {
									"start": {
										"line": 12,
										"character": 33
									},
									"end": {
										"line": 12,
										"character": 56
									}
								},
								"selectionRange": {
									"start": {
										"line": 12,
										"character": 33
									},
									"end": {
										"line": 12,
										"character": 56
									}
								},
								"children": []
							}
						]
					}
				]
			},
			{
				"name": "emptyMethod",
				"kind": "Method",
				"range": {
					"start": {
						"line": 15,
						"character": 4
					},
					"end": {
						"line": 15,
						"character": 21
					}
				},
				"selectionRange": {
					"start": {
						"line": 15,
						"character": 4
					},
					"end": {
						"line": 15,
						"character": 15
					}
				},
				"children": []
			}
		]
	}
]
```