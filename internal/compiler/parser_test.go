package compiler

import (
	"testing"
)

func BenchmarkParse(b *testing.B) {
	srcCompilerCheckerTS.SkipIfNotExist(b)

	fileName := srcCompilerCheckerTS.Path()
	sourceText := srcCompilerCheckerTS.ReadFile(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseSourceFile(fileName, sourceText, ScriptTargetESNext)
	}
}
