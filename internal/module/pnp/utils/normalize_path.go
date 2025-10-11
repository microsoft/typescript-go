package utils

import (
	"os"
	"strings"
)

func NormalizePath(original string) string {
	origPortable := toPortablePath(original)

	rooted := strings.HasPrefix(origPortable, "/")

	body := origPortable
	if rooted {
		body = strings.TrimPrefix(body, "/")
	}

	var parts []string
	if body != "" {
		parts = strings.FieldsFunc(body, func(r rune) bool { return r == '/' || r == '\\' })
	}

	out := make([]string, 0, len(parts))

	for _, comp := range parts {
		switch comp {
		case "", ".":
		case "..":
			switch {
			case rooted && len(out) == 0:
			case len(out) == 0 || out[len(out)-1] == "..":
				out = append(out, "..")
			default:
				out = out[:len(out)-1]
			}
		default:
			out = append(out, comp)
		}
	}

	if rooted {
		if len(out) == 0 {
			return fromPortablePath("/")
		}
		out = append([]string{""}, out...)
	}

	if len(out) == 0 {
		return fromPortablePath(".")
	}

	str := strings.Join(out, "/")

	hasTrailing := strings.HasSuffix(origPortable, "/") ||
		strings.HasSuffix(original, string(os.PathSeparator))
	if hasTrailing && !strings.HasSuffix(str, "/") {
		str += "/"
	}

	return fromPortablePath(str)
}
