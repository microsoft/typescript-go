package utils

import (
	"regexp"
	"runtime"
)

var reWindowsPath = regexp.MustCompile(`^([a-zA-Z]:.*)$`)
var reUNCWindowsPath = regexp.MustCompile(`^[\/\\][\/\\](\.[\/\\])?(.*)$`)
var rePortablePath = regexp.MustCompile(`^\/([a-zA-Z]:.*)$`)
var reUNCPortablePath = regexp.MustCompile(`^\/unc\/(\.dot\/)?(.*)$`)

func toPortablePath(s string) string {
	if runtime.GOOS != "windows" {
		return s
	}
	if m := reWindowsPath.FindStringSubmatch(s); m != nil {
		return "/" + m[1]
	}
	if m := reUNCWindowsPath.FindStringSubmatch(s); m != nil {
		if m[1] != "" {
			return "/unc/.dot/" + m[2]
		}
		return "/unc/" + m[2]
	}
	return s
}

func fromPortablePath(s string) string {
	if runtime.GOOS != "windows" {
		return s
	}
	// "/C:..." → "C:..."
	if m := rePortablePath.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	// "/unc/(.dot/)?rest"
	if m := reUNCPortablePath.FindStringSubmatch(s); m != nil {
		if m[1] != "" {
			return `\\.\` + m[2]
		}
		return `\\` + m[2]
	}
	return s
}
