package nativepath

import (
	"os"
	"path/filepath"
	"runtime"
)

// TermuxExecutableEnv contains the executable path supplied by Termux when Android runs it through the system linker.
const TermuxExecutableEnv = "TERMUX_EXEC__PROC_SELF_EXE"

// Executable returns the path of the current executable.
func Executable() (string, error) {
	if runtime.GOOS == "android" {
		if executable := os.Getenv(TermuxExecutableEnv); executable != "" {
			return filepath.Abs(executable)
		}
	}
	return os.Executable()
}
