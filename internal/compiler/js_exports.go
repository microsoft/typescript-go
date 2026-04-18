//go:build js && wasm

package compiler

import (
	"syscall/js"
)

// InitJSExports attaches internal compiler APIs to the provided exports map.
// This allows the js/wasm environment (e.g. ts-morph) to parse and mutate ASTs.
func InitJSExports(exports map[string]interface{}) {
	exports["parseSourceFile"] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Mock implementation: actual logic would use parser.ParseSourceFile
		if len(args) < 2 {
			return 0
		}
		// fileName := args[0].String()
		// sourceText := args[1].String()
		// sf := parser.ParseSourceFile(fileName, sourceText, ...)
		// return GlobalRegistry.Register(sf)
		return 1 // Mock Handle ID
	})
}
