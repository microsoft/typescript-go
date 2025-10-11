package utils

import "testing"

func TestNormalizePath(t *testing.T) {
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
