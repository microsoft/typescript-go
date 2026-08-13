package osutil

import (
	"os"
	"path/filepath"
)

const termuxExecutableEnv = "TERMUX_EXEC__PROC_SELF_EXE"

func args() []string {
	args := os.Args                                            //nolint:forbidigo
	if os.Getenv(termuxExecutableEnv) != "" && len(args) > 1 { //nolint:forbidigo
		// Termux launches non-cgo binaries through Android's linker, leaving the executable path in argv[1].
		if executable, err := Executable(); err == nil && args[1] == executable {
			return append([]string{args[0]}, args[2:]...)
		}
	}
	return args
}

func executable() (string, error) {
	// Under Termux, /proc/self/exe identifies Android's linker rather than the executable it loaded.
	if executable := os.Getenv(termuxExecutableEnv); executable != "" { //nolint:forbidigo
		return filepath.Abs(executable) //nolint:forbidigo
	}
	return os.Executable() //nolint:forbidigo
}
