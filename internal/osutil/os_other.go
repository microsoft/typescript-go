//go:build !android

package osutil

import "os"

func args() []string {
	return os.Args //nolint:forbidigo
}

func executable() (string, error) {
	return os.Executable() //nolint:forbidigo
}
