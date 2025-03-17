package api

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

type Callback int

const (
	CallbackDirectoryExists Callback = 1 << iota
	CallbackFileExists
	CallbackGetAccessibleEntries
	CallbackReadFile
	CallbackRealpath
)

type ServerOptions struct {
	In                 io.Reader
	Out                io.Writer
	Err                io.Writer
	Cwd                string
	NewLine            string
	DefaultLibraryPath string
}

var _ APIHost = (*Server)(nil)
var _ vfs.FS = (*Server)(nil)

type Server struct {
	r      *bufio.Reader
	w      *bufio.Writer
	stderr io.Writer

	cwd                string
	newLine            string
	fs                 vfs.FS
	defaultLibraryPath string

	callbackMu       sync.Mutex
	enabledCallbacks Callback
	logger           *project.Logger
	api              *API

	requestId int
}

func NewServer(options *ServerOptions) *Server {
	if options.Cwd == "" {
		panic("Cwd is required")
	}

	server := &Server{
		r:                  bufio.NewReader(options.In),
		w:                  bufio.NewWriter(options.Out),
		stderr:             options.Err,
		cwd:                options.Cwd,
		newLine:            options.NewLine,
		fs:                 bundled.WrapFS(osvfs.FS()),
		defaultLibraryPath: options.DefaultLibraryPath,
	}
	logger := project.NewLogger([]io.Writer{options.Err}, "", project.LogLevelVerbose)
	api := NewAPI(server, APIOptions{
		Logger: logger,
	})
	server.logger = logger
	server.api = api
	return server
}

// DefaultLibraryPath implements APIHost.
func (s *Server) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

// FS implements APIHost.
func (s *Server) FS() vfs.FS {
	return s
}

// GetCurrentDirectory implements APIHost.
func (s *Server) GetCurrentDirectory() string {
	return s.cwd
}

// NewLine implements APIHost.
func (s *Server) NewLine() string {
	return s.newLine
}

func (s *Server) Run() error {
	for {
		messageType, method, payload, err := s.readRequest("")
		if err != nil {
			return err
		}

		switch messageType {
		case "request":
			now := time.Now()
			result, err := s.handleRequest(method, payload)

			s.logger.PerfTrace(fmt.Sprintf("%s handled - %s", method, time.Since(now)))
			now = time.Now()
			if err != nil {
				if err := s.sendError(method, err); err != nil {
					return err
				}
			} else {
				if err := s.sendResponse(method, result); err != nil {
					return err
				}
				s.logger.PerfTrace(fmt.Sprintf("%s sent - %s", method, time.Since(now)))
			}
		default:
			return fmt.Errorf("%w: expected request, recieved: %s", ErrInvalidRequest, messageType)
		}
	}
}

func (s *Server) readRequest(expectedMethod string) (messageType string, method string, payload []byte, err error) {
	messageType, err = s.r.ReadString('\t')
	if err != nil {
		return "", "", nil, err
	}
	messageType = messageType[:len(messageType)-1]
	method, err = s.r.ReadString('\t')
	if err != nil {
		return "", "", nil, err
	}
	method = method[:len(method)-1]
	if expectedMethod != "" && method != expectedMethod {
		return "", "", nil, fmt.Errorf("%w: expected method %q, recieved %q", ErrInvalidRequest, expectedMethod, method)
	}
	var size uint32
	if err = binary.Read(s.r, binary.LittleEndian, &size); err != nil {
		return "", "", nil, fmt.Errorf("%w: expected payload size: %w", ErrInvalidRequest, err)
	}
	payload = make([]byte, size)
	bytesRead, err := io.ReadFull(s.r, payload)
	if err != nil {
		return "", "", nil, err
	}
	if bytesRead != int(size) {
		return "", "", nil, fmt.Errorf("%w: expected %d bytes, read %d", ErrInvalidRequest, size, bytesRead)
	}
	return messageType, method, payload, nil
}

func (s *Server) enableCallback(callback string) error {
	switch callback {
	case "directoryExists":
		s.enabledCallbacks |= CallbackDirectoryExists
	case "fileExists":
		s.enabledCallbacks |= CallbackFileExists
	case "getAccessibleEntries":
		s.enabledCallbacks |= CallbackGetAccessibleEntries
	case "readFile":
		s.enabledCallbacks |= CallbackReadFile
	case "realpath":
		s.enabledCallbacks |= CallbackRealpath
	default:
		return fmt.Errorf("unknown callback: %s", callback)
	}
	return nil
}

func (s *Server) handleRequest(method string, payload []byte) ([]byte, error) {
	s.requestId++
	switch method {
	case "configure":
		return nil, s.handleConfigure(payload)
	case "echo":
		return payload, nil
	default:
		return s.api.HandleRequest(s.requestId, method, payload)
	}
}

func (s *Server) handleConfigure(payload []byte) error {
	var params *ConfigureParams
	if err := json.Unmarshal(payload, &params); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}
	for _, callback := range params.Callbacks {
		if err := s.enableCallback(callback); err != nil {
			return err
		}
	}
	if params.LogFile != "" {
		s.logger.SetFile(params.LogFile)
	} else {
		s.logger.SetFile("")
	}
	return nil
}

func (s *Server) sendResponse(method string, result []byte) error {
	return s.write("response", method, result)
}

func (s *Server) sendError(method string, err error) error {
	payload, err := json.Marshal(err.Error())
	if err != nil {
		return err
	}
	return s.write("error", method, payload)
}

func (s *Server) write(messageType, method string, payload []byte) error {
	if _, err := s.w.WriteString(messageType); err != nil {
		return err
	}
	if _, err := s.w.WriteString("\t"); err != nil {
		return err
	}
	if _, err := s.w.WriteString(method); err != nil {
		return err
	}
	if err := s.w.WriteByte('\t'); err != nil {
		return err
	}
	if err := binary.Write(s.w, binary.LittleEndian, uint32(len(payload))); err != nil {
		return err
	}
	if _, err := s.w.Write(payload); err != nil {
		return err
	}
	return s.w.Flush()
}

func (s *Server) call(method string, payload any) ([]byte, error) {
	s.callbackMu.Lock()
	defer s.callbackMu.Unlock()
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	if err = s.write("call", method, jsonPayload); err != nil {
		return nil, err
	}

	messageType, _, responsePayload, err := s.readRequest(method)
	if err != nil {
		return nil, err
	}

	if string(messageType) != "call-response" && string(messageType) != "call-error" {
		return nil, fmt.Errorf("%w: expected call-response or call-error, recieved: %s", ErrInvalidRequest, messageType)
	}

	if string(messageType) == "call-error" {
		return nil, fmt.Errorf("%w: %s", ErrClientError, responsePayload)
	}

	return responsePayload, nil
}

// DirectoryExists implements vfs.FS.
func (s *Server) DirectoryExists(path string) bool {
	if s.enabledCallbacks&CallbackDirectoryExists != 0 {
		result, err := s.call("directoryExists", path)
		if err != nil {
			panic(err)
		}
		if len(result) > 0 {
			return string(result) == "true"
		}
	}
	return s.fs.DirectoryExists(path)
}

// FileExists implements vfs.FS.
func (s *Server) FileExists(path string) bool {
	if s.enabledCallbacks&CallbackFileExists != 0 {
		result, err := s.call("fileExists", path)
		if err != nil {
			panic(err)
		}
		if len(result) > 0 {
			return string(result) == "true"
		}
	}
	return s.fs.FileExists(path)
}

// GetAccessibleEntries implements vfs.FS.
func (s *Server) GetAccessibleEntries(path string) vfs.Entries {
	if s.enabledCallbacks&CallbackGetAccessibleEntries != 0 {
		result, err := s.call("getAccessibleEntries", path)
		if err != nil {
			panic(err)
		}
		if len(result) > 0 {
			var rawEntries *struct {
				Files       []string `json:"files"`
				Directories []string `json:"directories"`
			}
			if err := json.Unmarshal(result, &rawEntries); err != nil {
				panic(err)
			}
			if rawEntries != nil {
				return vfs.Entries{
					Files:       rawEntries.Files,
					Directories: rawEntries.Directories,
				}
			}
		}
	}
	return s.fs.GetAccessibleEntries(path)
}

// ReadFile implements vfs.FS.
func (s *Server) ReadFile(path string) (contents string, ok bool) {
	if s.enabledCallbacks&CallbackReadFile != 0 {
		data, err := s.call("readFile", path)
		if err != nil {
			panic(err)
		}
		if string(data) == "null" {
			return "", false
		}
		if len(data) > 0 {
			var result string
			if err := json.Unmarshal(data, &result); err != nil {
				panic(err)
			}
			return result, true
		}
	}
	return s.fs.ReadFile(path)
}

// Realpath implements vfs.FS.
func (s *Server) Realpath(path string) string {
	if s.enabledCallbacks&CallbackRealpath != 0 {
		data, err := s.call("realpath", path)
		if err != nil {
			panic(err)
		}
		if len(data) > 0 {
			var result string
			if err := json.Unmarshal(data, &result); err != nil {
				panic(err)
			}
			return result
		}
	}
	return s.fs.Realpath(path)
}

// UseCaseSensitiveFileNames implements vfs.FS.
func (s *Server) UseCaseSensitiveFileNames() bool {
	return s.fs.UseCaseSensitiveFileNames()
}

// WriteFile implements vfs.FS.
func (s *Server) WriteFile(path string, data string, writeByteOrderMark bool) error {
	return s.fs.WriteFile(path, data, writeByteOrderMark)
}

// WalkDir implements vfs.FS.
func (s *Server) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	panic("unimplemented")
}

// Stat implements vfs.FS.
func (s *Server) Stat(path string) vfs.FileInfo {
	panic("unimplemented")
}
