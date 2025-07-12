package core

import (
	"runtime"
	"testing"

	"gotest.tools/v3/assert"
)

func TestConcurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		opts           *CompilerOptions
		numFiles       int
		singleThreaded bool
		checkerCount   int
	}{
		{"defaults", &CompilerOptions{}, 100, false, 4},
		{"default", &CompilerOptions{Concurrency: "default"}, 100, false, 4},
		{"auto", &CompilerOptions{Concurrency: "true"}, 100, false, 4},
		{"true", &CompilerOptions{Concurrency: "true"}, 100, false, 4},
		{"yes", &CompilerOptions{Concurrency: "yes"}, 100, false, 4},
		{"on", &CompilerOptions{Concurrency: "on"}, 100, false, 4},
		{"singleThreaded", &CompilerOptions{SingleThreaded: TSTrue}, 100, true, 1},
		{"single", &CompilerOptions{Concurrency: "single"}, 100, true, 1},
		{"none", &CompilerOptions{Concurrency: "none"}, 100, true, 1},
		{"false", &CompilerOptions{Concurrency: "false"}, 100, true, 1},
		{"no", &CompilerOptions{Concurrency: "no"}, 100, true, 1},
		{"off", &CompilerOptions{Concurrency: "off"}, 100, true, 1},
		{"max", &CompilerOptions{Concurrency: "max"}, 1000, false, runtime.GOMAXPROCS(0)},
		{"half", &CompilerOptions{Concurrency: "half"}, 1000, runtime.GOMAXPROCS(0)/2 == 1, runtime.GOMAXPROCS(0) / 2},
		{"checker-per-file", &CompilerOptions{Concurrency: "checker-per-file"}, 100, false, 100},
		{"more than files", &CompilerOptions{Concurrency: "1000"}, 100, false, 100},
		{"10", &CompilerOptions{Concurrency: "10"}, 100, false, 10},
		{"1", &CompilerOptions{Concurrency: "1"}, 100, true, 1},
		{"invalid", &CompilerOptions{Concurrency: "i dunno"}, 100, false, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := ParseConcurrency(tt.opts)
			singleThreaded := c.SingleThreaded()
			checkerCount := c.CheckerCount(tt.numFiles)
			assert.Equal(t, singleThreaded, tt.singleThreaded)
			assert.Equal(t, checkerCount, tt.checkerCount)
		})
	}

	t.Run("TestProgramConcurrency", func(t *testing.T) {
		t.Parallel()

		c, _ := TestProgramConcurrency()
		assert.Assert(t, c.CheckerCount(10000) > 0)
	})
}
