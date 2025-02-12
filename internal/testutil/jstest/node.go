package jstest

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/repo"
)

const loaderScript = `import script from "./script.js";
process.stdout.write(JSON.stringify(await script(...process.argv.slice(2))));`

const tsLoaderScript = `import script from "./script.js";
import * as ts from "./typescript.js";
process.stdout.write(JSON.stringify(await script(ts, ...process.argv.slice(2))));`

// EvalNodeScript imports a Node.js script that exports a single function,
// calls it with the provided arguments, and unmarshals the JSON-stringified
// awaited return value into T. The script's function can either be exported
// with `module.exports =` or `export default`.
func EvalNodeScript[T any](t testing.TB, script string, dir string, args ...string) (result T, err error) {
	return evalNodeScript[T](t, script, loaderScript, dir, args...)
}

// EvalNodeScriptWithTS is like EvalNodeScript, but provides the TypeScript
// library to the script as the first argument.
func EvalNodeScriptWithTS[T any](t testing.TB, script string, dir string, args ...string) (result T, err error) {
	if dir == "" {
		dir = t.TempDir()
	}
	tsDest := filepath.Join(dir, "typescript.js")
	tsSrc := filepath.Join(repo.RootPath, "node_modules/typescript/lib/typescript.js")
	tsText, err := os.ReadFile(tsSrc)
	if err != nil {
		return result, err
	}
	if err = os.WriteFile(tsDest, tsText, 0o644); err != nil {
		return result, err
	}
	return evalNodeScript[T](t, script, tsLoaderScript, dir, args...)
}

func evalNodeScript[T any](t testing.TB, script string, loader string, dir string, args ...string) (result T, err error) {
	t.Helper()
	exe := getNodeExe(t)
	scriptPath := dir + "/script.js"
	if err = os.WriteFile(scriptPath, []byte(script), 0o644); err != nil {
		return result, err
	}
	loaderPath := dir + "/loader.js"
	if err = os.WriteFile(loaderPath, []byte(loader), 0o644); err != nil {
		return result, err
	}

	execArgs := make([]string, 0, 1+len(args))
	execArgs = append(execArgs, loaderPath)
	execArgs = append(execArgs, args...)
	execCmd := exec.Command(exe, execArgs...)
	execCmd.Dir = dir
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return result, fmt.Errorf("failed to run node: %w\n%s", err, output)
	}

	if err = json.Unmarshal(output, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON output: %w", err)
	}

	return result, nil
}

func getNodeExe(t testing.TB) string {
	t.Helper()

	const exeName = "node"
	exe, err := exec.LookPath(exeName)
	if err != nil {
		t.Skipf("%s not found: %v", exeName, err)
	}
	return exe
}
