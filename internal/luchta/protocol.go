package luchta

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type Run struct {
	ID      string            `json:"id"`
	Command string            `json:"command"`
	Cwd     string            `json:"cwd"`
	Env     map[string]string `json:"env"`
}

type ResolveTask struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Package string   `json:"package"`
	Cwd     string   `json:"cwd"`
	Scripts []string `json:"scripts"`
	Mode    string   `json:"mode"`
}

// DecodeMessage parses one JSONL line into a *Run or *ResolveTask.
func DecodeMessage(line []byte) (any, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(line, &probe); err != nil {
		return nil, fmt.Errorf("malformed message: %w", err)
	}
	switch probe.Type {
	case "run":
		var r Run
		if err := json.Unmarshal(line, &r); err != nil {
			return nil, fmt.Errorf("malformed run: %w", err)
		}
		return &r, nil
	case "resolveTask":
		var r ResolveTask
		if err := json.Unmarshal(line, &r); err != nil {
			return nil, fmt.Errorf("malformed resolveTask: %w", err)
		}
		return &r, nil
	default:
		return nil, fmt.Errorf("unknown message type %q", probe.Type)
	}
}

// Writer serializes protocol responses onto an io.Writer (one JSON object per line).
type Writer struct {
	mu  sync.Mutex
	enc *json.Encoder
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{enc: json.NewEncoder(w)}
}

func (w *Writer) emit(v any) {
	w.mu.Lock()
	defer w.mu.Unlock()
	_ = w.enc.Encode(v) // json.Encoder.Encode appends '\n'
}

func (w *Writer) Log(id, stream, line string) {
	type logMsg struct {
		Type   string `json:"type"`
		ID     string `json:"id"`
		Stream string `json:"stream"`
		Line   string `json:"line"`
	}
	w.emit(logMsg{Type: "log", ID: id, Stream: stream, Line: line})
}

// Report attaches a file (with a MIME type) to the task result, per the luchta
// worker protocol's `report` message. luchta persists it in the task cache and
// pretty-prints known MIME types (e.g. SARIF) in `luchta logs`.
func (w *Writer) Report(id, filename, mimeType, content string) {
	type reportMsg struct {
		Type     string `json:"type"`
		ID       string `json:"id"`
		Filename string `json:"filename"`
		MimeType string `json:"mimeType"`
		Content  string `json:"content"`
	}
	w.emit(reportMsg{Type: "report", ID: id, Filename: filename, MimeType: mimeType, Content: content})
}

func (w *Writer) Done(id string, exitCode int, inputs, outputs []string) {
	type doneMsg struct {
		Type     string   `json:"type"`
		ID       string   `json:"id"`
		ExitCode int      `json:"exitCode"`
		Inputs   []string `json:"inputs,omitempty"`
		Outputs  []string `json:"outputs,omitempty"`
	}
	w.emit(doneMsg{Type: "done", ID: id, ExitCode: exitCode, Inputs: inputs, Outputs: outputs})
}

func (w *Writer) Resolved(id, decision string) {
	type resultMsg struct {
		Decision string `json:"decision"`
	}
	type resolvedMsg struct {
		Type   string    `json:"type"`
		ID     string    `json:"id"`
		Result resultMsg `json:"result"`
	}
	w.emit(resolvedMsg{Type: "resolved", ID: id, Result: resultMsg{Decision: decision}})
}
