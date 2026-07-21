package lsp

import (
	"regexp"
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
)

// VS Code's telemetry pipeline redacts any string matching
// /(key|token|sig|secret|signature|password|passwd|pwd|android:value)[^a-zA-Z0-9]/i
// as `<REDACTED: Generic Secret>`, which trips on innocuous Go frames like
// `getSignatureHelp(`. Insert `X_X` after each trigger keyword that we know
// can appear in our sanitized output, when followed by punctuation we
// actually emit (`(`, `[`, `.`, `|`); reverse by removing the marker (replace
// `X_X` with the empty string) on the dashboard.
var genericSecretKeywordRegex = regexp.MustCompile(`(?i)(key|token|signature|sig|pwd)([(\[.|])`)

func defeatGenericSecretRegex(s string) string {
	return genericSecretKeywordRegex.ReplaceAllString(s, "${1}X_X${2}")
}

func sanitizeStackTrace(stack string) string {
	// TODO: should we just look for the first '(' and
	// just strip everything before the prior newline?
	startIndex := strings.Index(stack, "runtime/debug.Stack()")
	if startIndex >= 0 {
		stack = stack[startIndex:]
	} else {
		// For stacks that have already had recovery frames trimmed
		// (e.g., from PanicWithStack), find the first line containing our module.
		moduleIndex := strings.Index(stack, "typescript-go/internal")
		if moduleIndex < 0 {
			return ""
		}
		// Back up to the beginning of the line containing our module.
		lineStart := strings.LastIndex(stack[:moduleIndex], "\n")
		if lineStart >= 0 {
			stack = stack[lineStart+1:]
		}
	}

	result := &strings.Builder{}

	for lineNum, line := range core.Enumerate(strings.Lines(stack)) {
		if lineNum > 0 {
			result.WriteByte('\n')
		}

		i := 0
		// Skip whitespace
		for i < len(line) {
			if line[i] != ' ' && line[i] != '\t' {
				break
			}
			i++
		}

		result.WriteString(line[:i])

		line = line[i:]

		ourModuleIndex := strings.Index(line, "typescript-go/internal")
		if ourModuleIndex >= 0 {
			line = line[ourModuleIndex:]
			writeSanitizedModuleOrPath(line, result)
		} else {
			result.WriteString("(REDACTED FRAME)")
		}
	}

	return defeatGenericSecretRegex(result.String())
}

func writeSanitizedModuleOrPath(line string, result *strings.Builder) {
	// We don't expect things like \r, but it doesn't hurt to trim just in case.
	line = strings.TrimSpace(line)

	if plusHex := strings.Index(line, " +0x"); plusHex >= 0 {
		line = line[:plusHex]
	} else if inGoroutine := strings.LastIndex(line, " in goroutine "); inGoroutine >= 0 {
		line = line[:inGoroutine]
	}

	for segmentIndex, segment := range strings.Split(line, "/") {
		if segmentIndex > 0 {
			result.WriteString("|>")
		}

		// See if the string ends with ), and strip out all the arguments.
		if strings.HasSuffix(segment, ")") {
			openParenIndex := strings.LastIndexByte(segment, '(')
			if openParenIndex < 0 {
				// Closing parenthesis, but no opening - bail out.
				result.WriteString("???")
				continue
			}

			segment = segment[:openParenIndex]
			result.WriteString(segment)
			result.WriteString("()")
			continue
		}

		result.WriteString(segment)
	}
}
