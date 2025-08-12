package execute_test

import (
	"testing"
)

func TestBuildCommandLine(t *testing.T) {
	t.Parallel()
	testCases := []*tscInput{
		{
			subScenario:     "help",
			files:           FileMap{},
			commandLineArgs: []string{"--build", "--help"},
		},
	}

	for _, test := range testCases {
		test.run(t, "commandLine")
	}
}
