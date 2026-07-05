package luchta

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

const maxLineLength = 1 << 20 // matches luchta MAX_LINE_LENGTH

// Serve reads JSONL messages from in, dispatches each Run/ResolveTask, and writes
// protocol responses to out. Free-form parse errors go to errw (stderr). It returns
// when in reaches EOF, after all in-flight Runs complete.
func Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer) error {
	w := NewWriter(out)
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineLength)

	var wg sync.WaitGroup
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		msg, err := DecodeMessage([]byte(line))
		if err != nil {
			fmt.Fprintf(errw, "luchta-tsc-worker: %v\n", err)
			continue
		}
		switch m := msg.(type) {
		case *ResolveTask:
			handleResolveTask(w, m, errw)
		case *Run:
			wg.Add(1)
			go func(run *Run) {
				defer wg.Done()
				handleRun(ctx, w, run)
			}(m)
		}
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// mergeInputs combines declared inputs from task config with worker-specific inputs.
// The result is deduped and sorted for deterministic output.
func mergeInputs(declared, workerInputs []string) []string {
	seen := make(map[string]struct{}, len(declared)+len(workerInputs))
	for _, inp := range declared {
		seen[inp] = struct{}{}
	}
	for _, inp := range workerInputs {
		seen[inp] = struct{}{}
	}
	result := make([]string, 0, len(seen))
	for k := range seen {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func handleResolveTask(w *Writer, task *ResolveTask, errw io.Writer) {
	workerInputs, err := ResolveInputs(task.Cwd)
	if err != nil {
		fmt.Fprintf(errw, "luchta-tsc-worker: resolve inputs: %v\n", err)
		// On error, fall back to accept — let the run phase handle problems
		w.ResolvedAccept(task.ID)
		return
	}
	// Merge declared inputs with worker-specific inputs (deduped, sorted)
	merged := mergeInputs(task.Inputs, workerInputs)
	w.ResolvedModify(task.ID, merged)
}

func handleRun(ctx context.Context, w *Writer, run *Run) {
	defer func() {
		if r := recover(); r != nil {
			w.Log(run.ID, "stderr", fmt.Sprintf("panic: %v", r))
			w.Done(run.ID, 1, nil)
		}
	}()
	res := CompilePackage(ctx, run.Cwd)
	if res.InternalError != "" {
		w.Log(run.ID, "stderr", res.InternalError)
	}
	// Emit diagnostics as a structured SARIF report rather than free-form text,
	// so luchta can render IDE-clickable errors/warnings (see luchta PR #110).
	if len(res.Diagnostics) > 0 {
		w.Report(run.ID, sarifReportFilename, sarifMimeType, DiagnosticsToSARIF(res.Diagnostics))
	}
	w.Done(run.ID, res.ExitCode, res.Outputs)
}
