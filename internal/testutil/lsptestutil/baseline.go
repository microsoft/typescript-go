package lsptestutil

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type LocationMarker interface {
	FileName() string
	LSPos() lsproto.Position
}

type BaselineLocationsOptions struct {
	// markerInfo
	Marker     LocationMarker // location
	MarkerName string         // name of the marker to be printed in baseline

	EndMarker string

	StartMarkerPrefix func(span lsproto.Location) *string
	EndMarkerSuffix   func(span lsproto.Location) *string

	AdditionalLocation *lsproto.Location

	OpenFiles map[string]string
}

func GetBaselineForLocationsWithFileContents(
	f vfs.FS,
	spans []lsproto.Location,
	options BaselineLocationsOptions,
) string {
	locationsByFile := collections.GroupBy(spans, func(span lsproto.Location) lsproto.DocumentUri { return span.Uri })
	rangesByFile := collections.MultiMap[lsproto.DocumentUri, lsproto.Range]{}
	for file, locs := range locationsByFile.M {
		for _, loc := range locs {
			rangesByFile.Add(file, loc.Range)
		}
	}
	return GetBaselineForGroupedLocationsWithFileContents(
		f,
		&rangesByFile,
		options,
	)
}

func GetBaselineForGroupedLocationsWithFileContents(
	f vfs.FS,
	groupedRanges *collections.MultiMap[lsproto.DocumentUri, lsproto.Range],
	options BaselineLocationsOptions,
) string {
	// We must always print the file containing the marker,
	// but don't want to print it twice at the end if it already
	// found in a file with ranges.
	foundMarker := false
	foundAdditionalLocation := false

	baselineEntries := []string{}
	err := f.WalkDir("/", func(path string, d vfs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if !d.Type().IsRegular() {
			return nil
		}

		fileName := lsconv.FileNameToDocumentURI(path)
		ranges := groupedRanges.Get(fileName)
		if len(ranges) == 0 {
			return nil
		}

		content, ok := readFile(f, path, options)
		if !ok {
			// !!! error?
			return nil
		}

		if options.Marker != nil && options.Marker.FileName() == path {
			foundMarker = true
		}

		if options.AdditionalLocation != nil && options.AdditionalLocation.Uri == fileName {
			foundAdditionalLocation = true
		}

		baselineEntries = append(baselineEntries, getBaselineContentForFile(path, content, ranges, nil, options))
		return nil
	})

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic("walkdir error during fourslash baseline: " + err.Error())
	}

	// In Strada, there is a bug where we only ever add additional spans to baselines if we haven't
	// already added the file to the baseline.
	if options.AdditionalLocation != nil && !foundAdditionalLocation {
		fileName := options.AdditionalLocation.Uri.FileName()
		if content, ok := readFile(f, fileName, options); ok {
			baselineEntries = append(
				baselineEntries,
				getBaselineContentForFile(fileName, content, []lsproto.Range{options.AdditionalLocation.Range}, nil, options),
			)
			if options.Marker != nil && options.Marker.FileName() == fileName {
				foundMarker = true
			}
		}
	}

	if !foundMarker && options.Marker != nil {
		// If we didn't find the marker in any file, we need to add it.
		markerFileName := options.Marker.FileName()
		if content, ok := readFile(f, markerFileName, options); ok {
			baselineEntries = append(baselineEntries, getBaselineContentForFile(markerFileName, content, nil, nil, options))
		}
	}

	// !!! skipDocumentContainingOnlyMarker

	return strings.Join(baselineEntries, "\n\n")
}

func readFile(f vfs.FS, fileName string, options BaselineLocationsOptions) (string, bool) {
	if content, ok := options.OpenFiles[fileName]; ok {
		return content, ok
	}
	return f.ReadFile(fileName)
}

type baselineDetail struct {
	pos            lsproto.Position
	positionMarker string
	span           *lsproto.Range
	kind           string
}

func getBaselineContentForFile(
	fileName string,
	content string,
	spansInFile []lsproto.Range,
	spanToContextId map[lsproto.Range]int,
	options BaselineLocationsOptions,
) string {
	details := []*baselineDetail{}
	detailPrefixes := map[*baselineDetail]string{}
	detailSuffixes := map[*baselineDetail]string{}
	canDetermineContextIdInline := true
	uri := lsconv.FileNameToDocumentURI(fileName)

	if options.Marker != nil && options.Marker.FileName() == fileName {
		details = append(details, &baselineDetail{pos: options.Marker.LSPos(), positionMarker: options.MarkerName})
	}

	for _, span := range spansInFile {
		textSpanIndex := len(details)
		details = append(details,
			&baselineDetail{pos: span.Start, positionMarker: "[|", span: &span, kind: "textStart"},
			&baselineDetail{pos: span.End, positionMarker: core.OrElse(options.EndMarker, "|]"), span: &span, kind: "textEnd"},
		)

		if options.StartMarkerPrefix != nil {
			startPrefix := options.StartMarkerPrefix(lsproto.Location{Uri: uri, Range: span})
			if startPrefix != nil {
				// Special case: if this span starts at the same position as the provided marker,
				// we want the span's prefix to appear before the marker name.
				// i.e. We want `/*START PREFIX*/A: /*RENAME*/[|ARENAME|]`,
				// not `/*RENAME*//*START PREFIX*/A: [|ARENAME|]`
				if options.Marker != nil && fileName == options.Marker.FileName() && span.Start == options.Marker.LSPos() {
					_, ok := detailPrefixes[details[0]]
					debug.Assert(!ok, "Expected only single prefix at marker location")
					detailPrefixes[details[0]] = *startPrefix
				} else {
					detailPrefixes[details[textSpanIndex]] = *startPrefix
				}
			}
		}

		if options.EndMarkerSuffix != nil {
			endSuffix := options.EndMarkerSuffix(lsproto.Location{Uri: uri, Range: span})
			if endSuffix != nil {
				detailSuffixes[details[textSpanIndex+1]] = *endSuffix
			}
		}
	}

	slices.SortStableFunc(details, func(d1, d2 *baselineDetail) int {
		return lsproto.ComparePositions(d1.pos, d2.pos)
	})
	// !!! if canDetermineContextIdInline

	textWithContext := newTextWithContext(fileName, content)

	// Our preferred way to write marker is
	// /*MARKER*/[| some text |]
	// [| some /*MARKER*/ text |]
	// [| some text |]/*MARKER*/
	// Stable sort should handle first two cases but with that marker will be before rangeEnd if locations match
	// So we will defer writing marker in this case by checking and finding index of rangeEnd if same
	var deferredMarkerIndex *int

	for index, detail := range details {
		if detail.span == nil && deferredMarkerIndex == nil {
			// If this is marker position and its same as textEnd and/or contextEnd we want to write marker after those
			for matchingEndPosIndex := index + 1; matchingEndPosIndex < len(details); matchingEndPosIndex++ {
				// Defer after the location if its same as rangeEnd
				if details[matchingEndPosIndex].pos == detail.pos && strings.HasSuffix(details[matchingEndPosIndex].kind, "End") {
					deferredMarkerIndex = ptrTo(matchingEndPosIndex)
				}
				// Dont defer further than already determined
				break
			}
			// Defer writing marker position to deffered marker index
			if deferredMarkerIndex != nil {
				continue
			}
		}
		textWithContext.add(detail)
		textWithContext.pos = detail.pos
		// Prefix
		prefix := detailPrefixes[detail]
		if prefix != "" {
			textWithContext.newContent.WriteString(prefix)
		}
		textWithContext.newContent.WriteString(detail.positionMarker)
		if detail.span != nil {
			switch detail.kind {
			case "textStart":
				var text string
				if contextId, ok := spanToContextId[*detail.span]; ok {
					isAfterContextStart := false
					for textStartIndex := index - 1; textStartIndex >= 0; textStartIndex-- {
						textStartDetail := details[textStartIndex]
						if textStartDetail.kind == "contextStart" && textStartDetail.span == detail.span {
							isAfterContextStart = true
							break
						}
						// Marker is ok to skip over
						if textStartDetail.span != nil {
							break
						}
					}
					// Skip contextId on span thats surrounded by context span immediately
					if !isAfterContextStart {
						if text == "" {
							text = fmt.Sprintf(`contextId: %v`, contextId)
						} else {
							text = fmt.Sprintf(`contextId: %v`, contextId) + `, ` + text
						}
					}
				}
				if text != "" {
					textWithContext.newContent.WriteString(`{ ` + text + ` |}`)
				}
			case "contextStart":
				if canDetermineContextIdInline {
					spanToContextId[*detail.span] = len(spanToContextId)
				}
			}

			if deferredMarkerIndex != nil && *deferredMarkerIndex == index {
				// Write the marker
				textWithContext.newContent.WriteString(options.MarkerName)
				deferredMarkerIndex = nil
				detail = details[0] // Marker detail
			}
		}
		if suffix, ok := detailSuffixes[detail]; ok {
			textWithContext.newContent.WriteString(suffix)
		}
	}
	textWithContext.add(nil)
	if textWithContext.newContent.Len() != 0 {
		textWithContext.readableContents.WriteString("\n")
		textWithContext.readableJsoncBaseline(textWithContext.newContent.String())
	}
	return textWithContext.readableContents.String()
}

var LineSplitter = regexp.MustCompile(`\r?\n`)

type textWithContext struct {
	nLinesContext int // number of context lines to write to baseline

	readableContents *strings.Builder // builds what will be returned to be written to baseline

	newContent *strings.Builder // helper; the part of the original file content to write between details
	pos        lsproto.Position
	isLibFile  bool
	fileName   string
	content    string // content of the original file
	lineStarts *lsconv.LSPLineMap
	converters *lsconv.Converters

	// posLineInfo
	posInfo  *lsproto.Position
	lineInfo int
}

// implements ls.Script
func (t *textWithContext) FileName() string {
	return t.fileName
}

// implements ls.Script
func (t *textWithContext) Text() string {
	return t.content
}

func newTextWithContext(fileName string, content string) *textWithContext {
	t := &textWithContext{
		nLinesContext: 4,

		readableContents: &strings.Builder{},

		isLibFile:  regexp.MustCompile(`lib.*\.d\.ts$`).MatchString(fileName),
		newContent: &strings.Builder{},
		pos:        lsproto.Position{Line: 0, Character: 0},
		fileName:   fileName,
		content:    content,
		lineStarts: lsconv.ComputeLSPLineStarts(content),
	}

	t.converters = lsconv.NewConverters(lsproto.PositionEncodingKindUTF8, func(_ string) *lsconv.LSPLineMap {
		return t.lineStarts
	})
	t.readableContents.WriteString("// === " + fileName + " ===")
	return t
}

func (t *textWithContext) add(detail *baselineDetail) {
	if t.content == "" && detail == nil {
		panic("Unsupported")
	}
	if detail == nil || (detail.kind != "textEnd" && detail.kind != "contextEnd") {
		// Calculate pos to location number of lines
		posLineIndex := t.lineInfo
		if t.posInfo == nil || *t.posInfo != t.pos {
			posLineIndex = t.lineStarts.ComputeIndexOfLineStart(t.converters.LineAndCharacterToPosition(t, t.pos))
		}

		locationLineIndex := len(t.lineStarts.LineStarts) - 1
		if detail != nil {
			locationLineIndex = t.lineStarts.ComputeIndexOfLineStart(t.converters.LineAndCharacterToPosition(t, detail.pos))
			t.posInfo = &detail.pos
			t.lineInfo = locationLineIndex
		}

		nLines := 0
		if t.newContent.Len() != 0 {
			nLines += t.nLinesContext + 1
		}
		if detail != nil {
			nLines += t.nLinesContext + 1
		}
		// first nLinesContext and last nLinesContext
		if locationLineIndex-posLineIndex > nLines {
			if t.newContent.Len() != 0 {
				var skippedString string
				if t.isLibFile {
					skippedString = "--- (line: --) skipped ---\n"
				} else {
					skippedString = fmt.Sprintf(`// --- (line: %v) skipped ---`, posLineIndex+t.nLinesContext+1)
				}

				t.readableContents.WriteString("\n")
				t.readableJsoncBaseline(t.newContent.String() + t.sliceOfContent(
					t.getIndex(t.pos),
					t.getIndex(t.lineStarts.LineStarts[posLineIndex+t.nLinesContext]),
				) + skippedString)

				if detail != nil {
					t.readableContents.WriteString("\n")
				}
				t.newContent.Reset()
			}
			if detail != nil {
				if t.isLibFile {
					t.newContent.WriteString("--- (line: --) skipped ---\n")
				} else {
					t.newContent.WriteString(fmt.Sprintf("--- (line: %v) skipped ---\n", locationLineIndex-t.nLinesContext+1))
				}
				t.newContent.WriteString(t.sliceOfContent(
					t.getIndex(t.lineStarts.LineStarts[locationLineIndex-t.nLinesContext+1]),
					t.getIndex(detail.pos),
				))
			}
			return
		}
	}
	if detail == nil {
		t.newContent.WriteString(t.sliceOfContent(t.getIndex(t.pos), nil))
	} else {
		t.newContent.WriteString(t.sliceOfContent(t.getIndex(t.pos), t.getIndex(detail.pos)))
	}
}

func (t *textWithContext) readableJsoncBaseline(text string) {
	for i, line := range LineSplitter.Split(text, -1) {
		if i > 0 {
			t.readableContents.WriteString("\n")
		}
		t.readableContents.WriteString(`// ` + line)
	}
}

func (t *textWithContext) sliceOfContent(start *int, end *int) string {
	if start == nil || *start < 0 {
		start = ptrTo(0)
	}

	if end == nil || *end > len(t.content) {
		end = ptrTo(len(t.content))
	}

	if *start > *end {
		return ""
	}

	return t.content[*start:*end]
}

func (t *textWithContext) getIndex(i any) *int {
	switch i := i.(type) {
	case *int:
		return i
	case int:
		return ptrTo(i)
	case core.TextPos:
		return ptrTo(int(i))
	case *core.TextPos:
		return ptrTo(int(*i))
	case lsproto.Position:
		return t.getIndex(t.converters.LineAndCharacterToPosition(t, i))
	case *lsproto.Position:
		return t.getIndex(t.converters.LineAndCharacterToPosition(t, *i))
	}
	panic(fmt.Sprintf("getIndex: unsupported type %T", i))
}
