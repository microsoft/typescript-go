package nativepath

import (
	"os"
	"runtime"
)

const termuxExecutableEnv = "TERMUX_EXEC__PROC_SELF_EXE"

// Executable returns the path of the current executable.
func Executable() (string, error) {
	if runtime.GOOS == "android" {
		if executable := os.Getenv(termuxExecutableEnv); executable != "" {
			return executable, nil
		}
	}
	return os.Executable()
}
