package utils

import (
	"runtime"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"", "."},
		{"/", "/"},
		{"foo", "foo"},
		{"foo/bar", "foo/bar"},
		{"foo//bar", "foo/bar"},
		{"foo/./bar", "foo/bar"},
		{"foo/../bar", "bar"},
		{"foo/..//bar", "bar"},
		{"foo/bar/..", "foo"},
		{"foo/../../bar", "../bar"},
		{"../foo/../../bar", "../../bar"},
		{"./foo", "foo"},
		{"../foo", "../foo"},
		{"../D:/foo", "../D:/foo"},
		{"/foo/bar", "/foo/bar"},
		{"/foo/../../bar/baz", "/bar/baz"},
		{"/../foo/bar", "/foo/bar"},
		{"/../foo/bar//", "/foo/bar/"},
		{"/foo/bar/", "/foo/bar/"},
	}

	for _, tc := range tests {
		if got := NormalizePath(tc.in); got != tc.want {
			t.Errorf("NormalizePath(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormalizePathWindow(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-only test")
	}
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{`D:\foo\..\bar`, "D:/bar"},
		{`D:\foo\..\..\C:\bar\test`, "C:/bar/test"},
		{`\\server-name\foo\..\bar`, `\\server-name/bar`},
		{`\\server-name\foo\..\..\..\C:\bar\test`, "C:/bar/test"},
	}

	for _, tc := range tests {
		if got := NormalizePath(tc.in); got != tc.want {
			t.Errorf("NormalizePath(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}