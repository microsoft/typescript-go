package compiler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type testFile struct {
	FileName string `json:"name"`
	Contents string `json:"contents"`
}

type programTest struct {
	TestName      string     `json:"name"`
	Files         []testFile `json:"files"`
	ExpectedFiles []string   `json:"expectedFiles"`
}

func TestProgram(t *testing.T) {
	t.Parallel()
	testsFilePath := filepath.Join(repo.TestDataPath, "fixtures", "program", "program_test_cases.json")

	file, err := os.Open(testsFilePath)
	if err != nil {
		t.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	var programTestCases []programTest
	err = decoder.Decode(&programTestCases)
	if err != nil {
		t.Fatalf("Failed to decode test cases: %v", err)
	}

	for _, testCase := range programTestCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()
			fs := fstest.MapFS{}
			for _, testFile := range testCase.Files {
				fs[testFile.FileName] = &fstest.MapFile{
					Data: []byte(testFile.Contents),
				}
			}

			opts := core.CompilerOptions{}

			program := NewProgram(ProgramOptions{
				RootPath:       "c:/dev/src",
				Host:           NewCompilerHost(&opts, "c:/dev/src", vfstest.FromMapFS(fs, true)),
				Options:        &opts,
				SingleThreaded: false,
			})

			actualFiles := []string{}
			for _, file := range program.files {
				actualFiles = append(actualFiles, file.FileName())
			}

			assert.DeepEqual(t, testCase.ExpectedFiles, actualFiles)
		})
	}
}
