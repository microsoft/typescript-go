package ls_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/testutil/lstestutil"
	"gotest.tools/v3/assert"
)

type definitionTestCase struct {
	name     string
	files    map[string]string
	expected map[string][]ls.Location
}

func TestDefinition(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}

	testCases := []definitionTestCase{
		{
			name: "localFunction",
			files: map[string]string{
				mainFileName: `
function localFunction() { }
/*localFunction*/localFunction();`,
			},
			expected: map[string][]ls.Location{
				"localFunction": {{
					FileName: mainFileName,
					Range:    core.NewTextRange(9, 22),
				}},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			runDefinitionTest(t, testCase.files, testCase.expected)
		})
	}
}

func runDefinitionTest(t *testing.T, files map[string]string, expected map[string][]ls.Location) {
	parsedFiles := make(map[string]string)
	var markerPositions map[string]*lstestutil.Marker
	for fileName, content := range files {
		if fileName == mainFileName {
			testData := lstestutil.ParseTestData("", content, fileName)
			markerPositions = testData.MarkerPositions
			parsedFiles[fileName] = testData.Files[0].Content // !!! Assumes no usage of @filename
		} else {
			parsedFiles[fileName] = content
		}
	}
	languageService := createLanguageService(mainFileName, parsedFiles)
	for markerName, expectedResult := range expected {
		marker, ok := markerPositions[markerName]
		if !ok {
			t.Fatalf("No marker found for '%s'", markerName)
		}
		locations := languageService.ProvideDefinitions(mainFileName, marker.Position)
		assert.DeepEqual(t, locations, expectedResult, cmp.AllowUnexported(core.TextRange{}))
	}
}
