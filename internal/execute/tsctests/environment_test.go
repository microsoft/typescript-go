package tsctests

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"gotest.tools/v3/assert"
)

type programCaptureSystem struct {
	*TestSys
	programs []*incremental.Program
}

func (s *programCaptureSystem) OnProgram(program *incremental.Program) {
	s.TestSys.OnProgram(program)
	s.programs = append(s.programs, program)
}

func TestSingleThreadedEnvironmentVariable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		commandLineArgs []string
		files           FileMap
	}{
		{
			name:            "compile",
			commandLineArgs: []string{"--singleThreaded", "false"},
			files: FileMap{
				"/home/src/workspaces/project/index.ts":      `export const value = 1;`,
				"/home/src/workspaces/project/tsconfig.json": `{"compilerOptions":{"incremental":true,"noEmit":true}}`,
			},
		},
		{
			name:            "build",
			commandLineArgs: []string{"--build", "--singleThreaded", "false"},
			files: FileMap{
				"/home/src/workspaces/project/index.ts":      `export const value = 1;`,
				"/home/src/workspaces/project/tsconfig.json": `{"compilerOptions":{"composite":true}}`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			sys := &programCaptureSystem{TestSys: newTestSys(&tscInput{
				files: testCase.files,
				env:   map[string]string{"TS_SINGLE_THREADED": "1"},
			}, false)}

			result := execute.CommandLine(context.Background(), sys, testCase.commandLineArgs, sys)

			assert.Equal(t, result.Status, tsc.ExitStatusSuccess, sys.currentWrite.String())
			assert.Assert(t, len(sys.programs) > 0)
			for _, program := range sys.programs {
				assert.Assert(t, program.GetProgram().SingleThreaded())
			}
		})
	}
}
