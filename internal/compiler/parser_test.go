package compiler

import (
	"os"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	fileName := "../../_submodules/TypeScript/src/compiler/checker.ts"
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		b.Error(err)
	}
	sourceText := string(bytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseSourceFile(fileName, sourceText, ScriptTargetESNext)
	}
}
