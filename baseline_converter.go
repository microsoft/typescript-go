package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// TypeScript baseline format
type TSInlayHint struct {
	Text            string        `json:"text"`
	Position        int           `json:"position"`
	Kind            string        `json:"kind"`
	WhitespaceAfter bool          `json:"whitespaceAfter"`
	DisplayParts    []DisplayPart `json:"displayParts"`
}

type DisplayPart struct {
	Text string `json:"text"`
	Span *Span  `json:"span,omitempty"`
	File string `json:"file,omitempty"`
}

type Span struct {
	Start  int `json:"start"`
	Length int `json:"length"`
}

// Go port format
type GoInlayHint struct {
	Position     Position `json:"position"`
	Label        []Label  `json:"label"`
	Kind         int      `json:"kind"`
	PaddingRight bool     `json:"paddingRight"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Label struct {
	Value    string    `json:"value"`
	Location *Location `json:"location,omitempty"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// ConvertBaseline converts TypeScript baseline format to Go port format
func ConvertBaseline(input string) (string, error) {
	lines := strings.Split(input, "\n")
	var result strings.Builder

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Copy comment headers and code lines as-is
		if strings.HasPrefix(line, "//") || (!strings.HasPrefix(line, "{") && !strings.HasPrefix(line, "}") && strings.TrimSpace(line) != "") {
			result.WriteString(line + "\n")
			i++
			continue
		}

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			result.WriteString("\n")
			i++
			continue
		}

		// Check if this is the start of a JSON object
		if strings.HasPrefix(strings.TrimSpace(line), "{") {
			// Find the end of this JSON object
			jsonStart := i
			braceCount := 0
			jsonEnd := i

			for j := i; j < len(lines); j++ {
				for _, char := range lines[j] {
					if char == '{' {
						braceCount++
					} else if char == '}' {
						braceCount--
					}
				}
				if braceCount == 0 {
					jsonEnd = j
					break
				}
			}

			// Extract and parse the JSON
			jsonLines := lines[jsonStart : jsonEnd+1]
			jsonStr := strings.Join(jsonLines, "\n")

			var tsHint TSInlayHint
			if err := json.Unmarshal([]byte(jsonStr), &tsHint); err != nil {
				// If JSON parsing fails, just copy the lines as-is
				for _, jsonLine := range jsonLines {
					result.WriteString(jsonLine + "\n")
				}
				i = jsonEnd + 1
				continue
			}

			// Convert to Go format
			goHint := convertTSToGoHint(tsHint, lines)

			// Serialize to JSON with proper indentation
			goJSON, err := json.MarshalIndent(goHint, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to marshal Go hint: %v", err)
			}

			result.WriteString(string(goJSON) + "\n")
			i = jsonEnd + 1
		} else {
			result.WriteString(line + "\n")
			i++
		}
	}

	return result.String(), nil
}

func convertTSToGoHint(tsHint TSInlayHint, allLines []string) GoInlayHint {
	// Convert position from absolute character position to line/character
	pos := convertAbsolutePosition(tsHint.Position, allLines)

	// Convert kind from string to number
	kind := convertKind(tsHint.Kind)

	// Convert display parts to labels
	labels := make([]Label, len(tsHint.DisplayParts))
	for i, part := range tsHint.DisplayParts {
		label := Label{
			Value: part.Text,
		}

		// Add location if span is present
		if part.Span != nil && part.File != "" {
			// Convert file path and span to location
			uri := convertFilePathToURI(part.File)
			startPos := convertAbsolutePosition(part.Span.Start, allLines)
			endPos := convertAbsolutePosition(part.Span.Start+part.Span.Length, allLines)

			label.Location = &Location{
				URI: uri,
				Range: Range{
					Start: startPos,
					End:   endPos,
				},
			}
		}

		labels[i] = label
	}

	return GoInlayHint{
		Position:     pos,
		Label:        labels,
		Kind:         kind,
		PaddingRight: tsHint.WhitespaceAfter,
	}
}

func convertAbsolutePosition(absPos int, allLines []string) Position {
	// Find which line and character position corresponds to the absolute position
	currentPos := 0

	for lineNum, line := range allLines {
		lineLength := len(line) + 1 // +1 for newline character
		if currentPos+lineLength > absPos {
			return Position{
				Line:      lineNum,
				Character: absPos - currentPos,
			}
		}
		currentPos += lineLength
	}

	// If we can't find the position, return the last line
	if len(allLines) > 0 {
		return Position{
			Line:      len(allLines) - 1,
			Character: len(allLines[len(allLines)-1]),
		}
	}

	return Position{Line: 0, Character: 0}
}

func convertKind(kind string) int {
	switch kind {
	case "Parameter":
		return 2
	case "Type":
		return 1
	default:
		return 2 // Default to Parameter
	}
}

func convertFilePathToURI(filePath string) string {
	// Convert TypeScript test file path to Go test file path
	// Example: "/tests/cases/fourslash/inlayHintsInteractiveAnyParameter1.ts"
	// -> "file:///inlayHintsInteractiveAnyParameter1.ts"

	// Extract filename from path
	parts := strings.Split(filePath, "/")
	filename := parts[len(parts)-1]

	return "file:///" + filename
}

// ProcessBaselineFile processes an entire baseline file
func ProcessBaselineFile(input string) (string, error) {
	// First, we need to extract the code lines to properly calculate positions
	lines := strings.Split(input, "\n")

	// Find code lines (lines that are not comments or JSON)
	var codeLines []string
	jsonPattern := regexp.MustCompile(`^\s*[\{\}"]`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || jsonPattern.MatchString(trimmed) {
			continue
		}
		codeLines = append(codeLines, line)
	}

	return ConvertBaseline(input)
}

func main() {
	// Example usage
	input := `// === Inlay Hints ===
foo(1);
    ^
{
  "text": "",
  "position": 29,
  "kind": "Parameter",
  "whitespaceAfter": true,
  "displayParts": [
    {
      "text": "v",
      "span": {
        "start": 14,
        "length": 1
      },
      "file": "/tests/cases/fourslash/inlayHintsInteractiveAnyParameter1.ts"
    },
    {
      "text": ":"
    }
  ]
}`

	output, err := ProcessBaselineFile(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(output)
}
