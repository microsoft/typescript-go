package luchta

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
			w.Resolved(m.ID, "accept")
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

func handleRun(ctx context.Context, w *Writer, run *Run) {
	defer func() {
		if r := recover(); r != nil {
			w.Log(run.ID, "stderr", fmt.Sprintf("panic: %v", r))
			w.Done(run.ID, 1, nil, nil)
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
	w.Done(run.ID, res.ExitCode, res.Inputs, res.Outputs)
}
