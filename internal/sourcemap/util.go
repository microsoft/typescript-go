package sourcemap

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	sourceMapCommentRegExp       = regexp.MustCompile(`^\/\/[@#] sourceMappingURL=(.+)\r?\n?$`)
	whitespaceOrMapCommentRegExp = regexp.MustCompile(`^\s*(\/\/[@#] .*)?$`)
)

// Tries to find the sourceMappingURL comment at the end of a file.
func TryGetSourceMappingURL(lineInfo *LineInfo) string {
	for index := lineInfo.LineCount() - 1; index >= 0; index-- {
		line := lineInfo.LineText(index)
		comment := sourceMapCommentRegExp.FindStringSubmatch(line)
		if comment != nil {
			return strings.TrimRightFunc(comment[1], unicode.IsSpace)
		} else if !whitespaceOrMapCommentRegExp.MatchString(line) {
			break
		}
	}
	return ""
}
