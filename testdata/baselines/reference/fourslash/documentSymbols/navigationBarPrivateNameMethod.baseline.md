// === Document Symbols ===
```ts
// @FileName: /navigationBarPrivateNameMethod.ts
class A {
  #foo() {
    class B {
      #bar() {
         function baz () {
         }
      }
    }
  }
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
				"line": 9,
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
				"name": "#foo",
				"kind": "Method",
				"range": {
					"start": {
						"line": 1,
						"character": 2
					},
					"end": {
						"line": 8,
						"character": 3
					}
				},
				"selectionRange": {
					"start": {
						"line": 1,
						"character": 2
					},
					"end": {
						"line": 1,
						"character": 6
					}
				},
				"children": [
					{
						"name": "B",
						"kind": "Class",
						"range": {
							"start": {
								"line": 2,
								"character": 4
							},
							"end": {
								"line": 7,
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
								"name": "#bar",
								"kind": "Method",
								"range": {
									"start": {
										"line": 3,
										"character": 6
									},
									"end": {
										"line": 6,
										"character": 7
									}
								},
								"selectionRange": {
									"start": {
										"line": 3,
										"character": 6
									},
									"end": {
										"line": 3,
										"character": 10
									}
								},
								"children": [
									{
										"name": "baz",
										"kind": "Function",
										"range": {
											"start": {
												"line": 4,
												"character": 9
											},
											"end": {
												"line": 5,
												"character": 10
											}
										},
										"selectionRange": {
											"start": {
												"line": 4,
												"character": 18
											},
											"end": {
												"line": 4,
												"character": 21
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
	}
]
```