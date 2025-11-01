package fourslash

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
)

func (f *FourslashTest) addResultToBaseline(t *testing.T, command string, actual string) {
	b, ok := f.baselines[command]
	if !ok {
		f.baselines[command] = &strings.Builder{}
		b = f.baselines[command]
	}
	if b.Len() != 0 {
		b.WriteString("\n\n\n\n")
	}
	b.WriteString(`// === ` + command + " ===\n" + actual)
}

func (f *FourslashTest) writeToBaseline(command string, content string) {
	b, ok := f.baselines[command]
	if !ok {
		f.baselines[command] = &strings.Builder{}
		b = f.baselines[command]
	}
	b.WriteString(content)
}

func getBaselineFileName(t *testing.T, command string) string {
	return getBaseFileNameFromTest(t) + "." + getBaselineExtension(command)
}

func getBaselineExtension(command string) string {
	switch command {
	case "QuickInfo", "SignatureHelp", "Smart Selection":
		return "baseline"
	case "Auto Imports":
		return "baseline.md"
	case "findAllReferences", "goToDefinition", "findRenameLocations":
		return "baseline.jsonc"
	default:
		return "baseline.jsonc"
	}
}

func getBaselineOptions(command string) baseline.Options {
	subfolder := "fourslash/" + normalizeCommandName(command)
	switch command {
	case "Smart Selection":
		return baseline.Options{
			Subfolder:   subfolder,
			IsSubmodule: true,
		}
	case "findRenameLocations":
		return baseline.Options{
			Subfolder:   subfolder,
			IsSubmodule: true,
			DiffFixupOld: func(s string) string {
				var commandLines []string
				commandPrefix := regexp.MustCompile(`^// === ([a-z\sA-Z]*) ===`)
				testFilePrefix := "/tests/cases/fourslash"
				serverTestFilePrefix := "/server"
				contextSpanOpening := "<|"
				contextSpanClosing := "|>"
				oldPreference := "providePrefixAndSuffixTextForRename"
				newPreference := "useAliasesForRename"
				replacer := strings.NewReplacer(
					contextSpanOpening, "",
					contextSpanClosing, "",
					testFilePrefix, "",
					serverTestFilePrefix, "",
					oldPreference, newPreference,
				)
				lines := strings.Split(s, "\n")
				var isInCommand bool
				for _, line := range lines {
					if strings.HasPrefix(line, "// @findInStrings: ") || strings.HasPrefix(line, "// @findInComments: ") {
						continue
					}
					matches := commandPrefix.FindStringSubmatch(line)
					if len(matches) > 0 {
						commandName := matches[1]
						if commandName == command {
							isInCommand = true
						} else {
							isInCommand = false
						}
					}
					if isInCommand {
						fixedLine := replacer.Replace(line)
						commandLines = append(commandLines, fixedLine)
					}
				}
				return strings.Join(commandLines, "\n")
			},
		}
	default:
		return baseline.Options{
			Subfolder: subfolder,
		}
	}
}

func normalizeCommandName(command string) string {
	words := strings.Fields(command)
	command = strings.Join(words, "")
	return stringutil.LowerFirstChar(command)
}

func (f *FourslashTest) getBaselineForLocationsWithFileContents(spans []lsproto.Location, options lsptestutil.BaselineLocationsOptions) string {
	return lsptestutil.GetBaselineForLocationsWithFileContents(f.FS, spans, options)
}

type markerAndItem[T any] struct {
	Marker *Marker `json:"marker"`
	Item   T       `json:"item"`
}

func annotateContentWithTooltips[T comparable](
	t *testing.T,
	f *FourslashTest,
	markersAndItems []markerAndItem[T],
	opName string,
	getRange func(item T) *lsproto.Range,
	getTooltipLines func(item T, prev T) []string,
) string {
	barWithGutter := "| " + strings.Repeat("-", 70)

	// sort by file, then *backwards* by position in the file
	// so we can insert multiple times on a line without counting
	sorted := slices.Clone(markersAndItems)
	slices.SortFunc(sorted, func(a, b markerAndItem[T]) int {
		if c := cmp.Compare(a.Marker.FileName(), b.Marker.FileName()); c != 0 {
			return c
		}
		return -cmp.Compare(a.Marker.Position, b.Marker.Position)
	})

	filesToLines := collections.NewOrderedMapWithSizeHint[string, []string](1)
	var previous T
	for _, itemAndMarker := range sorted {
		marker := itemAndMarker.Marker
		item := itemAndMarker.Item

		textRange := getRange(item)
		if textRange == nil {
			start := marker.LSPosition
			end := start
			end.Character = end.Character + 1
			textRange = &lsproto.Range{Start: start, End: end}
		}

		if textRange.Start.Line != textRange.End.Line {
			t.Fatalf("Expected text range to be on a single line, got %v", textRange)
		}
		underline := strings.Repeat(" ", int(textRange.Start.Character)) +
			strings.Repeat("^", int(textRange.End.Character-textRange.Start.Character))

		fileName := marker.FileName()
		lines, ok := filesToLines.Get(fileName)
		if !ok {
			lines = lsptestutil.LineSplitter.Split(f.getScriptInfo(fileName).content, -1)
		}

		var tooltipLines []string
		if item != *new(T) {
			tooltipLines = getTooltipLines(item, previous)
		}
		if len(tooltipLines) == 0 {
			tooltipLines = []string{fmt.Sprintf("No %s at /*%s*/.", opName, *marker.Name)}
		}
		tooltipLines = core.Map(tooltipLines, func(line string) string {
			return "| " + line
		})

		linesToInsert := make([]string, len(tooltipLines)+3)
		linesToInsert[0] = underline
		linesToInsert[1] = barWithGutter
		copy(linesToInsert[2:], tooltipLines)
		linesToInsert[len(linesToInsert)-1] = barWithGutter

		lines = slices.Insert(
			lines,
			int(textRange.Start.Line+1),
			linesToInsert...,
		)
		filesToLines.Set(fileName, lines)

		previous = item
	}

	builder := strings.Builder{}
	seenFirst := false
	for fileName, lines := range filesToLines.Entries() {
		builder.WriteString(fmt.Sprintf("=== %s ===\n", fileName))
		for _, line := range lines {
			builder.WriteString("// ")
			builder.WriteString(line)
			builder.WriteByte('\n')
		}

		if seenFirst {
			builder.WriteString("\n\n")
		} else {
			seenFirst = true
		}
	}

	return builder.String()
}

func codeFence(lang string, code string) string {
	return "```" + lang + "\n" + code + "\n```"
}
