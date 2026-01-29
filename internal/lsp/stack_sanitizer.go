package lsp

import (
	"strings"
)

func sanitizeStackTrace(stack string) string {
	// TODO: should we just look for the first '(' and
	// just strip everything before the prior newline?
	startIndex := strings.Index(stack, "runtime/debug.Stack()")
	if startIndex < 0 {
		return ""
	}
	stack = stack[startIndex:]

	result := &strings.Builder{}

	// TODO: Use an iterator to avoid allocation.
	for lineNum, line := range strings.Split(stack, "\n") {
		if lineNum > 0 {
			result.WriteByte('\n')
		}

		i := 0
		// Skip whitespace
		for i = 0; i < len(line); i++ {
			if line[i] != ' ' && line[i] != '\t' {
				break
			}
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

	return result.String()
}

func writeSanitizedModuleOrPath(line string, result *strings.Builder) {
	// We don't expect things like \r, but it doesn't hurt to trim just in case.
	line = strings.TrimSpace(line)

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

const s = `goroutine 1196 [running]:
runtime/debug.Stack()
        /usr/local/go/src/runtime/debug/stack.go:26 +0x8e
github.com/microsoft/typescript-go/internal/lsp.(*Server).recover(0xc0001dae08, {0x14bc418, 0xc00bc60960}, 0xc00baf16e0)
        /workspaces/typescript-go/internal/lsp/server.go:777 +0x65
panic({0x1077b40?, 0x1abcb70?})
        /usr/local/go/src/runtime/panic.go:783 +0x136
github.com/microsoft/typescript-go/internal/ls.(*LanguageService).getCompletionData.func15()
        /workspaces/typescript-go/internal/ls/completions.go:1303 +0xfa
github.com/microsoft/typescript-go/internal/ls.(*LanguageService).getCompletionData.func18()
        /workspaces/typescript-go/internal/ls/completions.go:1548 +0x2df
github.com/microsoft/typescript-go/internal/ls.(*LanguageService).getCompletionData(0xc004b08240, {0x14bc418, 0xc00bc60a20}, 0xc0069ef908, 0xc000272008, 0x1b, 0xc002b28e00)
        /workspaces/typescript-go/internal/ls/completions.go:1581 +0x2b92
github.com/microsoft/typescript-go/internal/ls.(*LanguageService).getCompletionsAtPosition(0xc004b08240, {0x14bc418, 0xc00bc60a20}, 0xc000272008, 0x1b, 0x0)
        /workspaces/typescript-go/internal/ls/completions.go:347 +0x690
github.com/microsoft/typescript-go/internal/ls.(*LanguageService).ProvideCompletion(0xc004b08240, {0x14bc418, 0xc00bc60a20}, {0xc0092e02a0, 0x28}, {0x2, 0x4}, 0xc004580c30)
        /workspaces/typescript-go/internal/ls/completions.go:47 +0x207
github.com/microsoft/typescript-go/internal/lsp.(*Server).handleCompletion(0xc0001dae08, {0x14bc418, 0xc00bc60960}, 0xc004b08240, 0xc00baf14d0)
        /workspaces/typescript-go/internal/lsp/server.go:1102 +0xe5
github.com/microsoft/typescript-go/internal/lsp.registerLanguageServiceWithAutoImportsRequestHandler[...].func1({0x14bc418, 0xc00bc60960}, 0xc00baf16e0)
        /workspaces/typescript-go/internal/lsp/server.go:682 +0x32a
github.com/microsoft/typescript-go/internal/lsp.(*Server).handleRequestOrNotification(0xc0001dae08, {0x14bc418, 0xc00bc60960}, 0xc00baf16e0)
        /workspaces/typescript-go/internal/lsp/server.go:531 +0x11e
github.com/microsoft/typescript-go/internal/lsp.(*Server).dispatchLoop.func1()
        /workspaces/typescript-go/internal/lsp/server.go:414 +0x65
created by github.com/microsoft/typescript-go/internal/lsp.(*Server).dispatchLoop in goroutine 19
        /workspaces/typescript-go/internal/lsp/server.go:438 +0x60`

func main() {
	println(sanitizeStackTrace(s))
}
