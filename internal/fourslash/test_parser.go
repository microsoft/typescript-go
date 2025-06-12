package fourslash

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/testrunner"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// Inserted in source files by surrounding desired text
// in a range with `[|` and `|]`. For example,
//
// [|text in range|]
//
// is a range with `text in range` "selected".
type MarkerRange struct {
	*Marker
	Range   core.TextRange
	LSRange lsproto.Range
}

type Marker struct {
	Filename   string
	Position   int
	LSPosition lsproto.Position
	Name       string
	Data       map[string]interface{}
}

type TestData struct {
	Files           []*TestFileInfo
	MarkerPositions map[string]*Marker
	Markers         []*Marker
	Symlinks        map[string]string
	GlobalOptions   map[string]string
	Ranges          []*MarkerRange
}

type testFileWithMarkers struct {
	file    *TestFileInfo
	markers []*Marker
	ranges  []*MarkerRange
}

func ParseTestData(t *testing.T, contents string, fileName string) TestData {
	// List of all the subfiles we've parsed out
	var files []*TestFileInfo

	markerPositions := make(map[string]*Marker)
	var markers []*Marker
	var ranges []*MarkerRange
	filesWithMarker, symlinks, _, globalOptions := testrunner.ParseTestFilesAndSymlinks(
		contents,
		fileName,
		parseFileContent,
	)

	hasTSConfig := false
	for _, file := range filesWithMarker {
		files = append(files, file.file)
		hasTSConfig = hasTSConfig || isConfigFile(file.file.Filename)

		markers = append(markers, file.markers...)
		ranges = append(ranges, file.ranges...)
		for _, marker := range file.markers {
			if _, ok := markerPositions[marker.Name]; ok {
				t.Fatalf("Duplicate marker name: %s", marker.Name)
			}
			markerPositions[marker.Name] = marker
		}

	}

	if hasTSConfig && len(globalOptions) > 0 {
		t.Fatalf("It is not allowed to use global options along with config files.")
	}

	return TestData{
		Files:           files,
		MarkerPositions: markerPositions,
		Markers:         markers,
		Symlinks:        symlinks,
		GlobalOptions:   globalOptions,
		Ranges:          ranges,
	}
}

func isConfigFile(filename string) bool {
	filename = strings.ToLower(filename)
	return strings.HasSuffix(filename, "tsconfig.json") || strings.HasSuffix(filename, "jsconfig.json")
}

type locationInformation struct {
	position       int
	sourcePosition int
	sourceLine     int
	sourceColumn   int
}

type rangeLocationInformation struct {
	locationInformation
	marker *Marker
}

type TestFileInfo struct {
	Filename string
	// The contents of the file (with markers, etc stripped out)
	Content string
	emit    bool
}

// FileName implements ls.Script.
func (t *TestFileInfo) FileName() string {
	return t.Filename
}

// Text implements ls.Script.
func (t *TestFileInfo) Text() string {
	return t.Content
}

var _ ls.Script = (*TestFileInfo)(nil)

const emitThisFileOption = "emitthisfile"

type parserState int

const (
	stateNone parserState = iota
	stateInSlashStarMarker
	stateInObjectMarker
)

func parseFileContent(filename string, content string, fileOptions map[string]string) *testFileWithMarkers {
	filename = tspath.GetNormalizedAbsolutePath(filename, "/")

	// The file content (minus metacharacters) so far
	var output strings.Builder

	var markers []*Marker

	/// A stack of the open range markers that are still unclosed
	openRanges := []rangeLocationInformation{}
	/// A list of closed ranges we've collected so far
	localRanges := []*MarkerRange{}

	// The total number of metacharacters removed from the file (so far)
	difference := 0

	// One-based current position data
	line := 1
	column := 1

	// The current marker (or maybe multi-line comment?) we're parsing, possibly
	var openMarker locationInformation

	// The latest position of the start of an unflushed plain text area
	lastNormalCharPosition := 0

	flush := func(lastSafeCharIndex int) {
		if lastSafeCharIndex != -1 {
			output.WriteString(content[lastNormalCharPosition:lastSafeCharIndex])
		} else {
			output.WriteString(content[lastNormalCharPosition:])
		}
	}

	state := stateNone
	previousCharacter, i := utf8.DecodeRuneInString(content)
	var size int
	var currentCharacter rune
	for ; i < len(content); i = i + size {
		currentCharacter, size = utf8.DecodeRuneInString(content[i:])
		switch state {
		case stateNone:
			if previousCharacter == '[' && currentCharacter == '|' {
				// found a range start
				openRanges = append(openRanges, rangeLocationInformation{
					locationInformation: locationInformation{
						position:       (i - 1) - difference,
						sourcePosition: i - 1,
						sourceLine:     line,
						sourceColumn:   column,
					},
				})
				// copy all text up to marker position
				flush(i - 1)
				lastNormalCharPosition = i + 1
				difference += 2
			} else if previousCharacter == '|' && currentCharacter == ']' {
				// found a range end
				if len(openRanges) == 0 {
					reportError(filename, line, column, "Found range end with no matching start.")
				}
				rangeStart := openRanges[len(openRanges)-1]
				openRanges = openRanges[:len(openRanges)-1]

				closedRange := &MarkerRange{Range: core.NewTextRange(rangeStart.position, (i-1)-difference)}
				if rangeStart.marker != nil {
					closedRange.Marker = rangeStart.marker
				} else {
					closedRange.Marker = &Marker{Filename: filename}
				}

				localRanges = append(localRanges, closedRange)

				// copy all text up to range marker position
				flush(i - 1)
				lastNormalCharPosition = i + 1
				difference += 2
			} else if previousCharacter == '/' && currentCharacter == '*' {
				// found a possible marker start
				state = stateInSlashStarMarker
				openMarker = locationInformation{
					position:       (i - 1) - difference,
					sourcePosition: i - 1,
					sourceLine:     line,
					sourceColumn:   column - 1,
				}
			} else if previousCharacter == '{' && currentCharacter == '|' {
				// found an object marker start
				state = stateInObjectMarker
				openMarker = locationInformation{
					position:       (i - 1) - difference,
					sourcePosition: i - 1,
					sourceLine:     line,
					sourceColumn:   column,
				}
				flush(i - 1)
			}
		case stateInObjectMarker:
			// Object markers are only ever terminated by |} and have no content restrictions
			if previousCharacter == '|' && currentCharacter == '}' {
				objectMarkerData := strings.TrimSpace(content[openMarker.sourcePosition+2 : i-1])
				marker := getObjectMarker(filename, openMarker, objectMarkerData)

				if len(openRanges) > 0 {
					openRanges[len(openRanges)-1].marker = marker
				}
				markers = append(markers, marker)

				// Set the current start to point to the end of the current marker to ignore its text
				lastNormalCharPosition = i + 1
				difference += i + 1 - openMarker.sourcePosition

				// Reset the state
				openMarker = locationInformation{}
				state = stateNone
			}
		case stateInSlashStarMarker:
			if previousCharacter == '*' && currentCharacter == '/' {
				// Record the marker
				// start + 2 to ignore the */, -1 on the end to ignore the * (/ is next)
				markerNameText := strings.TrimSpace(content[openMarker.sourcePosition+2 : i-1])
				marker := &Marker{
					Filename: filename,
					Position: openMarker.position,
					Name:     markerNameText,
				}
				if len(openRanges) > 0 {
					openRanges[len(openRanges)-1].marker = marker
				}
				markers = append(markers, marker)

				// Set the current start to point to the end of the current marker to ignore its text
				flush(openMarker.sourcePosition)
				lastNormalCharPosition = i + 1
				difference += i + 1 - openMarker.sourcePosition

				// Reset the state
				openMarker = locationInformation{}
				state = stateNone
			} else if !(stringutil.IsDigit(currentCharacter) ||
				stringutil.IsASCIILetter(currentCharacter) ||
				currentCharacter == '$' ||
				currentCharacter == '_') { // Invalid marker character
				if currentCharacter == '*' && i < len(content)-1 && content[i+1] == '/' {
					// The marker is about to be closed, ignore the 'invalid' char
				} else {
					// We've hit a non-valid marker character, so we were actually in a block comment
					// Bail out the text we've gathered so far back into the output
					flush(i)
					lastNormalCharPosition = i
					openMarker = locationInformation{}
					state = stateNone
				}
			}
		}
		if currentCharacter == '\n' && previousCharacter == '\r' {
			// Ignore trailing \n after \r
			continue
		} else if currentCharacter == '\n' || currentCharacter == '\r' {
			line++
			column = 1
			continue
		}
		column++
		previousCharacter = currentCharacter
	}

	// Add the remaining text
	flush(-1)

	outputString := output.String()
	// Set LS positions for markers
	lineMap := ls.ComputeLineStarts(outputString)
	converters := ls.NewConverters(lsproto.PositionEncodingKindUTF8, func(_ string) *ls.LineMap {
		return lineMap
	})

	emit := fileOptions[emitThisFileOption] == "true"

	testFileInfo := &TestFileInfo{
		Filename: filename,
		Content:  outputString,
		emit:     emit,
	}

	for _, marker := range markers {
		marker.LSPosition = converters.PositionToLineAndCharacter(testFileInfo, core.TextPos(marker.Position))
	}

	return &testFileWithMarkers{
		file:    testFileInfo,
		markers: markers,
		ranges:  localRanges,
	}
}

func getObjectMarker(fileName string, location locationInformation, text string) *Marker {
	// Attempt to parse the marker value as JSON
	var v interface{}
	e := json.Unmarshal([]byte("{ "+text+" }"), &v)

	if e != nil {
		reportError(fileName, location.sourceLine, location.sourceColumn, "Unable to parse marker text "+text)
		return nil
	}
	markerValue, ok := v.(map[string]interface{})
	if !ok || len(markerValue) == 0 {
		reportError(fileName, location.sourceLine, location.sourceColumn, "Object markers can not be empty")
		return nil
	}

	marker := &Marker{
		Filename: fileName,
		Position: location.position,
		Data:     markerValue,
	}

	// Object markers can be anonymous
	if markerValue["name"] != nil {
		if name, ok := markerValue["name"].(string); ok && name != "" {
			marker.Name = name
		}
	}

	return marker
}

func reportError(fileName string, line int, col int, message string) {
	// !!! not implemented
	// errorMessage := fileName + "(" + string(line) + "," + string(col) + "): " + message;
}
