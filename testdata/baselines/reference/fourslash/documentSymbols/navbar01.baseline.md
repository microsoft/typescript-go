// === Document Symbols ===
```ts
// @FileName: /navbar01.ts
// Interface
interface IPoint {
    getDist(): number;
    new(): IPoint;
    (): any;
    [x:string]: number;
    prop: string;
}

/// Module
module Shapes {
    // Class
    export class Point implements IPoint {
        constructor (public x: number, public y: number) { }

        // Instance member
        getDist() { return Math.sqrt(this.x * this.x + this.y * this.y); }

        // Getter
        get value(): number { return 0; }

        // Setter
        set value(newValue: number) { return; }

        // Static member
        static origin = new Point(0, 0);

        // Static method
        private static getOrigin() { return Point.origin;}
    }

    enum Values { value1, value2, value3 }
}

// Local variables
var p: IPoint = new Shapes.Point(3, 4);
var dist = p.getDist();
```

# Symbols

```json
[
	{
		"name": "IPoint",
		"kind": "Interface",
		"range": {
			"start": {
				"line": 1,
				"character": 0
			},
			"end": {
				"line": 7,
				"character": 1
			}
		},
		"selectionRange": {
			"start": {
				"line": 1,
				"character": 10
			},
			"end": {
				"line": 1,
				"character": 16
			}
		},
		"children": [
			{
				"name": "getDist",
				"kind": "Method",
				"range": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 2,
						"character": 22
					}
				},
				"selectionRange": {
					"start": {
						"line": 2,
						"character": 4
					},
					"end": {
						"line": 2,
						"character": 11
					}
				},
				"children": []
			},
			{
				"name": "new()",
				"kind": "Constructor",
				"range": {
					"start": {
						"line": 3,
						"character": 4
					},
					"end": {
						"line": 3,
						"character": 18
					}
				},
				"selectionRange": {
					"start": {
						"line": 3,
						"character": 4
					},
					"end": {
						"line": 3,
						"character": 4
					}
				},
				"children": []
			},
			{
				"name": "()",
				"kind": "Function",
				"range": {
					"start": {
						"line": 4,
						"character": 4
					},
					"end": {
						"line": 4,
						"character": 12
					}
				},
				"selectionRange": {
					"start": {
						"line": 4,
						"character": 4
					},
					"end": {
						"line": 4,
						"character": 4
					}
				},
				"children": []
			},
			{
				"name": "[]",
				"kind": "Property",
				"range": {
					"start": {
						"line": 5,
						"character": 4
					},
					"end": {
						"line": 5,
						"character": 23
					}
				},
				"selectionRange": {
					"start": {
						"line": 5,
						"character": 4
					},
					"end": {
						"line": 5,
						"character": 4
					}
				},
				"children": []
			},
			{
				"name": "prop",
				"kind": "Property",
				"range": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 6,
						"character": 17
					}
				},
				"selectionRange": {
					"start": {
						"line": 6,
						"character": 4
					},
					"end": {
						"line": 6,
						"character": 8
					}
				},
				"children": []
			}
		]
	},
	{
		"name": "Shapes",
		"kind": "Namespace",
		"range": {
			"start": {
				"line": 10,
				"character": 0
			},
			"end": {
				"line": 32,
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
				"character": 13
			}
		},
		"children": [
			{
				"name": "Point",
				"kind": "Class",
				"range": {
					"start": {
						"line": 12,
						"character": 4
					},
					"end": {
						"line": 29,
						"character": 5
					}
				},
				"selectionRange": {
					"start": {
						"line": 12,
						"character": 17
					},
					"end": {
						"line": 12,
						"character": 22
					}
				},
				"children": [
					{
						"name": "constructor",
						"kind": "Constructor",
						"range": {
							"start": {
								"line": 13,
								"character": 8
							},
							"end": {
								"line": 13,
								"character": 60
							}
						},
						"selectionRange": {
							"start": {
								"line": 13,
								"character": 8
							},
							"end": {
								"line": 13,
								"character": 8
							}
						},
						"children": []
					},
					{
						"name": "x",
						"kind": "Property",
						"range": {
							"start": {
								"line": 13,
								"character": 21
							},
							"end": {
								"line": 13,
								"character": 37
							}
						},
						"selectionRange": {
							"start": {
								"line": 13,
								"character": 28
							},
							"end": {
								"line": 13,
								"character": 29
							}
						},
						"children": []
					},
					{
						"name": "y",
						"kind": "Property",
						"range": {
							"start": {
								"line": 13,
								"character": 39
							},
							"end": {
								"line": 13,
								"character": 55
							}
						},
						"selectionRange": {
							"start": {
								"line": 13,
								"character": 46
							},
							"end": {
								"line": 13,
								"character": 47
							}
						},
						"children": []
					},
					{
						"name": "getDist",
						"kind": "Method",
						"range": {
							"start": {
								"line": 16,
								"character": 8
							},
							"end": {
								"line": 16,
								"character": 74
							}
						},
						"selectionRange": {
							"start": {
								"line": 16,
								"character": 8
							},
							"end": {
								"line": 16,
								"character": 15
							}
						},
						"children": []
					},
					{
						"name": "value",
						"kind": "Property",
						"range": {
							"start": {
								"line": 19,
								"character": 8
							},
							"end": {
								"line": 19,
								"character": 41
							}
						},
						"selectionRange": {
							"start": {
								"line": 19,
								"character": 12
							},
							"end": {
								"line": 19,
								"character": 17
							}
						},
						"children": []
					},
					{
						"name": "value",
						"kind": "Property",
						"range": {
							"start": {
								"line": 22,
								"character": 8
							},
							"end": {
								"line": 22,
								"character": 47
							}
						},
						"selectionRange": {
							"start": {
								"line": 22,
								"character": 12
							},
							"end": {
								"line": 22,
								"character": 17
							}
						},
						"children": []
					},
					{
						"name": "origin",
						"kind": "Property",
						"range": {
							"start": {
								"line": 25,
								"character": 8
							},
							"end": {
								"line": 25,
								"character": 40
							}
						},
						"selectionRange": {
							"start": {
								"line": 25,
								"character": 15
							},
							"end": {
								"line": 25,
								"character": 21
							}
						},
						"children": []
					},
					{
						"name": "getOrigin",
						"kind": "Method",
						"range": {
							"start": {
								"line": 28,
								"character": 8
							},
							"end": {
								"line": 28,
								"character": 58
							}
						},
						"selectionRange": {
							"start": {
								"line": 28,
								"character": 23
							},
							"end": {
								"line": 28,
								"character": 32
							}
						},
						"children": []
					}
				]
			},
			{
				"name": "Values",
				"kind": "Enum",
				"range": {
					"start": {
						"line": 31,
						"character": 4
					},
					"end": {
						"line": 31,
						"character": 42
					}
				},
				"selectionRange": {
					"start": {
						"line": 31,
						"character": 9
					},
					"end": {
						"line": 31,
						"character": 15
					}
				},
				"children": [
					{
						"name": "value1",
						"kind": "EnumMember",
						"range": {
							"start": {
								"line": 31,
								"character": 18
							},
							"end": {
								"line": 31,
								"character": 24
							}
						},
						"selectionRange": {
							"start": {
								"line": 31,
								"character": 18
							},
							"end": {
								"line": 31,
								"character": 24
							}
						},
						"children": []
					},
					{
						"name": "value2",
						"kind": "EnumMember",
						"range": {
							"start": {
								"line": 31,
								"character": 26
							},
							"end": {
								"line": 31,
								"character": 32
							}
						},
						"selectionRange": {
							"start": {
								"line": 31,
								"character": 26
							},
							"end": {
								"line": 31,
								"character": 32
							}
						},
						"children": []
					},
					{
						"name": "value3",
						"kind": "EnumMember",
						"range": {
							"start": {
								"line": 31,
								"character": 34
							},
							"end": {
								"line": 31,
								"character": 40
							}
						},
						"selectionRange": {
							"start": {
								"line": 31,
								"character": 34
							},
							"end": {
								"line": 31,
								"character": 40
							}
						},
						"children": []
					}
				]
			}
		]
	},
	{
		"name": "p",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 35,
				"character": 4
			},
			"end": {
				"line": 35,
				"character": 38
			}
		},
		"selectionRange": {
			"start": {
				"line": 35,
				"character": 4
			},
			"end": {
				"line": 35,
				"character": 5
			}
		},
		"children": []
	},
	{
		"name": "dist",
		"kind": "Variable",
		"range": {
			"start": {
				"line": 36,
				"character": 4
			},
			"end": {
				"line": 36,
				"character": 22
			}
		},
		"selectionRange": {
			"start": {
				"line": 36,
				"character": 4
			},
			"end": {
				"line": 36,
				"character": 8
			}
		},
		"children": []
	}
]
```