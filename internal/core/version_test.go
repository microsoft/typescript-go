package core

import "testing"

func TestVersionMajorMinorAndVersionConstants(t *testing.T) {
    t.Parallel()

    constants := []struct  {
        name string
        got string
        expected string
    }{
        {"VersionMajorMinor", VersionMajorMinor, "7.0"},
		{"Version", Version, "7.0.0-dev"},
    }

    for _, tt := range constants {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("expected %q, but got %q", tt.expected, tt.got)
			}
		})
	}
}
