package contentmapper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/zeebo/xxh3"
)

// ProtocolVersion is the content mapper protocol version this host speaks.
const ProtocolVersion = 1

// Content mapper protocol method names.
const (
	MethodInitialize   = "initialize"
	MethodOpenProject  = "openProject"
	MethodCloseProject = "closeProject"
	MethodTransform    = "transform"
)

// InitializeParams is the parameter object for the initialize request.
type InitializeParams struct {
	ProtocolVersion int `json:"protocolVersion"`
	// Locale is the BCP 47 locale to use for mapper-authored diagnostic messages, when configured.
	Locale string `json:"locale,omitempty"`
	// PositionEncodings lists the coordinate spaces the host accepts.
	PositionEncodings []PositionEncoding `json:"positionEncodings"`
}

// InitializeResult is the mapper's response to the initialize request.
type InitializeResult struct {
	ProtocolVersion int `json:"protocolVersion"`
	// PositionEncoding selects the coordinate space for all mappings and diagnostics.
	PositionEncoding PositionEncoding `json:"positionEncoding"`
	// DiagnosticSource is the prefix used for every mapper-authored diagnostic code.
	DiagnosticSource string `json:"diagnosticSource"`
}

// OpenProjectParams is the parameter object for the openProject request.
type OpenProjectParams struct {
	// ConfigFileName is the absolute project configuration file name, or empty when there is none.
	ConfigFileName string `json:"configFileName"`
	// ProjectHandle is an opaque, process-local handle assigned by the host.
	ProjectHandle string `json:"projectHandle"`
	// Options is the mapper entry's options from the project's contentMappers configuration.
	Options json.Value `json:"options,omitempty"`
	// CompilerOptions contains the project's effective compiler options.
	CompilerOptions json.Value `json:"compilerOptions"`
}

// OpenProjectResult is the mapper's response to an openProject request.
type OpenProjectResult struct {
	// ConfigIdentity is a stable fingerprint of all dynamic configuration that can affect transforms.
	ConfigIdentity string `json:"configIdentity"`
	// WatchedFiles are absolute files whose changes may alter ConfigIdentity or transform output.
	WatchedFiles []string `json:"watchedFiles,omitempty"`
}

// CloseProjectParams is the parameter object for the closeProject request.
type CloseProjectParams struct {
	ProjectHandle string `json:"projectHandle"`
}

// PositionEncoding is the coordinate space a mapper uses for mappings and diagnostics.
type PositionEncoding string

const (
	PositionEncodingUTF8  PositionEncoding = "utf-8"
	PositionEncodingUTF16 PositionEncoding = "utf-16"
)

// TransformParams is the parameter object for the transform request.
type TransformParams struct {
	// FileName is the absolute name of the foreign file being transformed.
	FileName string `json:"fileName"`
	// Content is the foreign file's source text.
	Content string `json:"content"`
	// Options is the mapper entry's options from the project's contentMappers configuration.
	Options json.Value `json:"options,omitempty"`
	// ProjectHandle identifies the dynamic mapper project configuration, when one is required.
	ProjectHandle string `json:"projectHandle,omitempty"`
	// CompilerOptions holds the values of the options the mapper declared in initialize, keyed by option
	// name and ordered by the mapper's declaration. It is an empty object when the mapper declared none.
	CompilerOptions *collections.OrderedMap[string, json.Value] `json:"compilerOptions"`
}

// MappedOutput is generated source text and its mapping to an original input.
type MappedOutput struct {
	// Text is the generated JavaScript or TypeScript source text.
	Text string `json:"text"`
	// Mappings is the span map's tuple-array JSON (see spanmap.Marshal), expressed in the selected
	// position encoding. Absent or empty means the output is fully synthesized.
	Mappings json.Value `json:"mappings,omitempty"`
}

type SupplementalOutput struct {
	MappedOutput
	// Extension is the transformed virtual source extension.
	Extension string `json:"extension"`
}

// TransformResult is the canonical output for one input file.
type TransformResult struct {
	MappedOutput
	// Diagnostics are mapper-authored errors expressed in original-source coordinates.
	Diagnostics []Diagnostic `json:"diagnostics,omitempty"`
	// Supplemental contains additional unnamed compiler inputs associated with this source file.
	Supplemental []SupplementalOutput `json:"supplemental,omitempty"`
}

// Diagnostic is an error reported by a mapper in original-source coordinates.
type Diagnostic struct {
	// MessageText is the diagnostic message.
	MessageText string `json:"messageText"`
	// Start and Length locate the diagnostic in the original content using the selected position encoding.
	Start  int   `json:"start"`
	Length int   `json:"length"`
	Code   int32 `json:"code,omitempty"`
}

// dialFunc establishes a running connection to a mapper. In production it spawns the mapper's process;
// tests substitute an in-memory connection. It returns the connection and a closer that tears it down.
type dialFunc func(ctx context.Context, mapper *Mapper, diagnosticLocale locale.Locale) (ipc.Conn, io.Closer, PositionEncoding, string, error)

// host manages one child process per mapper identity. It is the production implementation of Host.
type host struct {
	ctx    context.Context
	cancel context.CancelFunc
	stop   func() bool
	dial   dialFunc

	lifecycleMu      sync.RWMutex
	diagnosticLocale locale.Locale

	mu            sync.Mutex
	conns         map[string]*mapperConn
	projects      map[string]*projectEntry
	projectLeases map[string]*projectLease
	nextProjectID uint64
}

type projectEntry struct {
	mapper         *Mapper
	spec           ProjectSpec
	projectHandle  string
	opened         bool
	configIdentity string
	watchedFiles   []string
}

type mapperConn struct {
	conn   ipc.Conn
	closer io.Closer
	// err, when non-nil, records that this mapper failed to start; it is cached so we do not repeatedly
	// try (and fail) to spawn a broken mapper.
	err              error
	positionEncoding PositionEncoding
	diagnosticSource string
	// refs is the number of active Acquire calls retaining this identity.
	refs int
}

var _ Host = (*host)(nil)

// Spawner starts a child process, returning its stdio as an io.ReadWriteCloser (Read is the
// process's stdout, Write is its stdin) whose Close tears the process down. This seam keeps os/exec out
// of this package: production hosts spawn a real process, tests supply an in-process pipe.
type Spawner interface {
	Spawn(command []string, dir string) (io.ReadWriteCloser, error)
}

// SpawnerFunc adapts a spawn function to the Spawner interface.
type SpawnerFunc func(command []string, dir string) (io.ReadWriteCloser, error)

func (f SpawnerFunc) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	return f(command, dir)
}

// NewHost creates a Host that spawns each mapper's process via the given spawner and drives it over a
// JSON-RPC connection. The host's lifetime is bound to ctx: cancelling it (e.g. the CLI's signal context
// on SIGINT, or a build/watch session ending) tears every mapper process down, so owners of a session
// context need not close the host explicitly. Close does the same synchronously.
func NewHost(ctx context.Context, spawner Spawner, diagnosticLocale locale.Locale) Host {
	return newWithDial(ctx, diagnosticLocale, func(ctx context.Context, mapper *Mapper, diagnosticLocale locale.Locale) (ipc.Conn, io.Closer, PositionEncoding, string, error) {
		if len(mapper.Exec) == 0 {
			return nil, nil, "", "", fmt.Errorf("content mapper %q declares no command to run", mapper.Package)
		}
		rwc, err := spawner.Spawn(mapper.Exec, mapper.PackageDirectory)
		if err != nil {
			return nil, nil, "", "", err
		}
		conn := ipc.NewAsyncConn(rwc, rejectHandler{})
		go func() { _ = conn.Run(ctx) }()
		positionEncoding, diagnosticSource, err := handshake(ctx, conn, diagnosticLocale)
		if err != nil {
			_ = rwc.Close()
			return nil, nil, "", "", fmt.Errorf("content mapper %q failed to initialize: %w", mapper.Package, err)
		}
		return conn, rwc, positionEncoding, diagnosticSource, nil
	})
}

func newWithDial(ctx context.Context, diagnosticLocale locale.Locale, dial dialFunc) *host {
	hostCtx, cancel := context.WithCancel(ctx)
	h := &host{ctx: hostCtx, cancel: cancel, dial: dial, diagnosticLocale: diagnosticLocale, conns: make(map[string]*mapperConn), projects: make(map[string]*projectEntry), projectLeases: make(map[string]*projectLease)}
	h.stop = context.AfterFunc(ctx, func() { _ = h.Close() })
	return h
}

func (h *host) SetLocale(diagnosticLocale locale.Locale) {
	h.lifecycleMu.Lock()
	defer h.lifecycleMu.Unlock()
	if h.diagnosticLocale.String() == diagnosticLocale.String() {
		return
	}
	h.diagnosticLocale = diagnosticLocale

	h.mu.Lock()
	var closers []io.Closer
	for _, entry := range h.conns {
		if entry.closer != nil {
			closers = append(closers, entry.closer)
		}
		entry.conn = nil
		entry.closer = nil
		entry.err = nil
		entry.positionEncoding = ""
		entry.diagnosticSource = ""
	}
	for _, project := range h.projects {
		project.opened = false
	}
	h.mu.Unlock()
	for _, closer := range closers {
		_ = closer.Close()
	}
}

func (h *host) Project(spec ProjectSpec) Project {
	h.lifecycleMu.RLock()
	defer h.lifecycleMu.RUnlock()

	key := projectSpecKey(spec)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.projects == nil {
		return nil
	}
	if lease := h.projectLeases[key]; lease != nil {
		return lease.retainLocked()
	}
	lease := &projectLease{host: h, key: key, entries: make(map[*Mapper]string, len(spec.Mappers))}
	lease.refs = 1
	for _, mapper := range spec.Mappers {
		entryKey := fmt.Sprintf("%s:%d", mapper.Identity(), h.nextProjectID)
		h.nextProjectID++
		entry := &projectEntry{mapper: mapper, spec: spec, projectHandle: entryKey}
		h.projects[entryKey] = entry
		connEntry := h.conns[mapper.Identity()]
		if connEntry == nil {
			connEntry = &mapperConn{}
			h.conns[mapper.Identity()] = connEntry
		}
		connEntry.refs++
		lease.entries[mapper] = entryKey
	}
	h.projectLeases[key] = lease
	return lease
}

func projectSpecKey(spec ProjectSpec) string {
	var key strings.Builder
	fmt.Fprintf(&key, "%s\x00%p", spec.ConfigFileName, spec.CompilerOptions)
	for _, mapper := range spec.Mappers {
		fmt.Fprintf(&key, "\x00%p", mapper)
	}
	return key.String()
}

func combinedIdentity(mapper *Mapper, configIdentity string) string {
	buf := make([]byte, 0, len(mapper.Identity())+len(mapper.Options)+len(configIdentity)+2)
	buf = append(buf, mapper.Identity()...)
	buf = append(buf, 0)
	buf = append(buf, mapper.Options...)
	buf = append(buf, 0)
	buf = append(buf, configIdentity...)
	hash := xxh3.Hash128(buf).Bytes()
	return mapper.Identity() + ":" + hex.EncodeToString(hash[:])
}

func (h *host) openProjectLocked(ctx context.Context, entry *projectEntry) error {
	if entry.opened {
		return nil
	}
	conn, _, _, err := h.connForLocked(entry.mapper)
	if err != nil {
		return err
	}
	compilerOptions, err := json.Marshal(entry.spec.CompilerOptions)
	if err != nil {
		return err
	}
	raw, err := conn.Call(ctx, MethodOpenProject, OpenProjectParams{
		ConfigFileName:  entry.spec.ConfigFileName,
		ProjectHandle:   entry.projectHandle,
		Options:         entry.mapper.Options,
		CompilerOptions: compilerOptions,
	})
	if err != nil {
		return err
	}
	var result OpenProjectResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return &ProjectError{Kind: ProjectErrorKindMalformedResponse}
	}
	if entry.mapper.DynamicConfig && result.ConfigIdentity == "" {
		return &ProjectError{Kind: ProjectErrorKindMissingConfigIdentity}
	}
	entry.configIdentity = result.ConfigIdentity
	for _, fileName := range result.WatchedFiles {
		if !tspath.PathIsAbsolute(fileName) {
			return &ProjectError{Kind: ProjectErrorKindNonAbsoluteWatchedFile}
		}
	}
	entry.watchedFiles = slices.Clone(result.WatchedFiles)
	entry.opened = true
	return nil
}

func (h *host) Acquire(mappers []*Mapper) func() {
	seen := make(map[string]struct{}, len(mappers))
	identities := make([]string, 0, len(mappers))
	h.mu.Lock()
	if h.conns != nil {
		for _, mapper := range mappers {
			identity := mapper.Identity()
			if _, ok := seen[identity]; ok {
				continue
			}
			seen[identity] = struct{}{}
			identities = append(identities, identity)
			entry := h.conns[identity]
			if entry == nil {
				entry = &mapperConn{}
				h.conns[identity] = entry
			}
			entry.refs++
		}
	}
	h.mu.Unlock()
	return sync.OnceFunc(func() { h.release(identities) })
}

// Transform sends the file's content to the mapper's process and decodes the transformed result. The
// mapper receives the subset of the project's compiler options it declared in its manifest (an empty
// object if it declared none).
func (h *host) Transform(mapper *Mapper, request Request) (Result, error) {
	h.lifecycleMu.RLock()
	defer h.lifecycleMu.RUnlock()
	return h.transformLocked(mapper, request, "")
}

func (h *host) transformLocked(mapper *Mapper, request Request, projectHandle string) (Result, error) {
	conn, positionEncoding, diagnosticSource, err := h.connFor(mapper)
	if err != nil {
		return Result{}, NewTransformError(TransformErrorKindInitialize, err)
	}
	options, err := mapper.MarshalDeclaredOptions(request.CompilerOptions)
	if err != nil {
		return Result{}, NewTransformError(TransformErrorKindCompilerOptions, err)
	}
	raw, err := conn.Call(h.ctx, MethodTransform, TransformParams{
		FileName:        request.FileName,
		Content:         request.Content,
		Options:         mapper.Options,
		ProjectHandle:   projectHandle,
		CompilerOptions: options,
	})
	if err != nil {
		return Result{}, NewTransformError(TransformErrorKindRequest, err)
	}
	decoded, err := decodeTransformResult(raw, request.Content, positionEncoding, diagnosticSource)
	if err != nil {
		return Result{}, NewTransformError(TransformErrorKindResponse, err)
	}
	return decoded, nil
}

// Close shuts down every mapper process. It is safe to call more than once and is invoked automatically
// when the context passed to New is cancelled.
func (h *host) Close() error {
	h.lifecycleMu.Lock()
	defer h.lifecycleMu.Unlock()
	h.stop()
	h.cancel()
	h.mu.Lock()
	var closers []io.Closer
	for _, mc := range h.conns {
		if mc.closer != nil {
			closers = append(closers, mc.closer)
		}
	}
	h.conns = nil
	h.projects = nil
	h.projectLeases = nil
	h.mu.Unlock()
	var errs []error
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// connFor returns the connection for a mapper's identity, spawning its process on first use. Mappers
// sharing an identity share a single process.
func (h *host) connFor(mapper *Mapper) (ipc.Conn, PositionEncoding, string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.connForLocked(mapper)
}

func (h *host) connForLocked(mapper *Mapper) (ipc.Conn, PositionEncoding, string, error) {
	if h.conns == nil {
		return nil, "", "", errors.New("content mapper host is closed")
	}
	identity := mapper.Identity()
	entry := h.conns[identity]
	if entry == nil {
		entry = &mapperConn{}
		h.conns[identity] = entry
	}
	if entry.conn != nil || entry.err != nil {
		return entry.conn, entry.positionEncoding, entry.diagnosticSource, entry.err
	}
	conn, closer, positionEncoding, diagnosticSource, err := h.dial(h.ctx, mapper, h.diagnosticLocale)
	entry.conn = conn
	entry.closer = closer
	entry.err = err
	entry.positionEncoding = positionEncoding
	entry.diagnosticSource = diagnosticSource
	return conn, positionEncoding, diagnosticSource, err
}

type projectLease struct {
	host    *host
	key     string
	entries map[*Mapper]string
	refs    int
	once    sync.Once
}

type retainedProject struct {
	*projectLease
	once sync.Once
}

func (p *projectLease) retainLocked() Project {
	p.refs++
	return &retainedProject{projectLease: p}
}

func (p *retainedProject) Close() (err error) {
	p.once.Do(func() { err = p.projectLease.release() })
	return err
}

func (p *projectLease) Refresh() error {
	p.host.lifecycleMu.RLock()
	defer p.host.lifecycleMu.RUnlock()
	p.host.mu.Lock()
	defer p.host.mu.Unlock()
	var result error
	for _, key := range p.entries {
		entry := p.host.projects[key]
		if entry == nil || !entry.opened {
			continue
		}
		if connEntry := p.host.conns[entry.mapper.Identity()]; connEntry != nil && connEntry.conn != nil {
			_, err := connEntry.conn.Call(p.host.ctx, MethodCloseProject, CloseProjectParams{ProjectHandle: entry.projectHandle})
			result = errors.Join(result, err)
		}
		entry.opened = false
	}
	return result
}

func (p *projectLease) Identities() ([]string, error) {
	p.host.mu.Lock()
	defer p.host.mu.Unlock()
	identities := make([]string, 0, len(p.entries))
	for mapper, key := range p.entries {
		entry := p.host.projects[key]
		if mapper.DynamicConfig {
			if err := p.host.openProjectLocked(p.host.ctx, entry); err != nil {
				return nil, err
			}
			identities = append(identities, combinedIdentity(mapper, entry.configIdentity))
		} else {
			hash := mapper.TransformIdentity(entry.spec.CompilerOptions).Bytes()
			identities = append(identities, mapper.Identity()+":"+hex.EncodeToString(hash[:]))
		}
	}
	slices.Sort(identities)
	return identities, nil
}

func (p *projectLease) Identity(mapper *Mapper) (string, error) {
	p.host.mu.Lock()
	defer p.host.mu.Unlock()
	entry := p.host.projects[p.entries[mapper]]
	if entry == nil {
		return "", nil
	}
	if mapper.DynamicConfig {
		if err := p.host.openProjectLocked(p.host.ctx, entry); err != nil {
			return "", err
		}
		return combinedIdentity(mapper, entry.configIdentity), nil
	}
	hash := mapper.TransformIdentity(entry.spec.CompilerOptions).Bytes()
	return mapper.Identity() + ":" + hex.EncodeToString(hash[:]), nil
}

func (p *projectLease) WatchedFiles() ([]string, error) {
	p.host.mu.Lock()
	defer p.host.mu.Unlock()
	var files []string
	for _, key := range p.entries {
		entry := p.host.projects[key]
		if entry.mapper.DynamicConfig {
			if err := p.host.openProjectLocked(p.host.ctx, entry); err != nil {
				return nil, err
			}
		}
		files = append(files, entry.watchedFiles...)
	}
	slices.Sort(files)
	return slices.Compact(files), nil
}

func (p *projectLease) Transform(mapper *Mapper, request Request) (Result, error) {
	p.host.lifecycleMu.RLock()
	defer p.host.lifecycleMu.RUnlock()
	p.host.mu.Lock()
	entry := p.host.projects[p.entries[mapper]]
	if entry == nil {
		p.host.mu.Unlock()
		return Result{}, errors.New("content mapper project is closed")
	}
	handle := ""
	if mapper.DynamicConfig {
		if err := p.host.openProjectLocked(p.host.ctx, entry); err != nil {
			p.host.mu.Unlock()
			return Result{}, NewTransformError(TransformErrorKindProject, err)
		}
		handle = entry.projectHandle
	}
	p.host.mu.Unlock()
	return p.host.transformLocked(mapper, request, handle)
}

func (p *projectLease) Close() error {
	var result error
	p.once.Do(func() {
		result = p.release()
	})
	return result
}

func (p *projectLease) release() error {
	var result error
	{
		p.host.lifecycleMu.RLock()
		defer p.host.lifecycleMu.RUnlock()
		var releasedIdentities []string
		p.host.mu.Lock()
		p.refs--
		if p.refs < 0 {
			p.host.mu.Unlock()
			panic("content mapper project reference count below zero")
		}
		if p.refs != 0 {
			p.host.mu.Unlock()
			return nil
		}
		if p.host.projectLeases[p.key] == p {
			delete(p.host.projectLeases, p.key)
		}
		for _, key := range p.entries {
			entry := p.host.projects[key]
			if entry == nil {
				continue
			}
			if entry.opened {
				if connEntry := p.host.conns[entry.mapper.Identity()]; connEntry != nil && connEntry.conn != nil {
					_, err := connEntry.conn.Call(p.host.ctx, MethodCloseProject, CloseProjectParams{ProjectHandle: entry.projectHandle})
					result = errors.Join(result, err)
				}
			}
			delete(p.host.projects, key)
			releasedIdentities = append(releasedIdentities, entry.mapper.Identity())
		}
		p.host.mu.Unlock()
		p.host.release(releasedIdentities)
	}
	return result
}

func (h *host) release(identities []string) {
	var closers []io.Closer
	h.mu.Lock()
	if h.conns != nil {
		for _, identity := range identities {
			entry := h.conns[identity]
			if entry == nil {
				continue
			}
			entry.refs--
			if entry.refs == 0 {
				delete(h.conns, identity)
				if entry.closer != nil {
					closers = append(closers, entry.closer)
				}
			}
		}
	}
	h.mu.Unlock()
	for _, closer := range closers {
		_ = closer.Close()
	}
}

func handshake(ctx context.Context, conn ipc.Conn, diagnosticLocale locale.Locale) (PositionEncoding, string, error) {
	raw, err := conn.Call(ctx, MethodInitialize, InitializeParams{
		ProtocolVersion:   ProtocolVersion,
		Locale:            diagnosticLocale.String(),
		PositionEncodings: []PositionEncoding{PositionEncodingUTF8, PositionEncodingUTF16},
	})
	if err != nil {
		return "", "", err
	}
	var res InitializeResult
	if err := json.Unmarshal(raw, &res); err != nil {
		return "", "", err
	}
	if res.ProtocolVersion != ProtocolVersion {
		return "", "", &InitializeError{Kind: InitializeErrorKindProtocolVersion, ProtocolVersion: res.ProtocolVersion}
	}
	if res.PositionEncoding != PositionEncodingUTF8 && res.PositionEncoding != PositionEncodingUTF16 {
		return "", "", &InitializeError{Kind: InitializeErrorKindPositionEncoding, PositionEncoding: res.PositionEncoding}
	}
	if strings.TrimSpace(res.DiagnosticSource) == "" {
		return "", "", &InitializeError{Kind: InitializeErrorKindEmptyDiagnosticSource}
	}
	if strings.EqualFold(res.DiagnosticSource, "typescript") || strings.EqualFold(res.DiagnosticSource, "tsc") {
		return "", "", &InitializeError{Kind: InitializeErrorKindReservedDiagnosticSource, DiagnosticSource: res.DiagnosticSource}
	}
	nativeExtensions := core.Flatten(tspath.AllSupportedExtensionsWithJson)
	if slices.ContainsFunc(nativeExtensions, func(extension string) bool {
		return strings.EqualFold(res.DiagnosticSource, strings.TrimPrefix(extension, "."))
	}) {
		return "", "", &InitializeError{Kind: InitializeErrorKindReservedDiagnosticSource, DiagnosticSource: res.DiagnosticSource}
	}
	return res.PositionEncoding, res.DiagnosticSource, nil
}

func decodeTransformResult(raw json.Value, originalText string, positionEncoding PositionEncoding, diagnosticSource string) (Result, error) {
	var res TransformResult
	if err := json.Unmarshal(raw, &res); err != nil {
		return Result{}, err
	}
	mapped, originalPositions, err := decodeMappedOutput(res.MappedOutput, originalText, positionEncoding)
	if err != nil {
		return Result{}, err
	}
	result := Result{
		Text:     mapped.Text,
		Mappings: mapped.Mappings,
	}
	for _, supplemental := range res.Supplemental {
		if !IsSupportedVirtualExtension(supplemental.Extension) {
			return Result{}, &InvalidSupplementalVirtualExtensionError{Extension: supplemental.Extension}
		}
		mapped, _, err := decodeMappedOutput(supplemental.MappedOutput, originalText, positionEncoding)
		if err != nil {
			return Result{}, err
		}
		mapped.VirtualExtension = supplemental.Extension
		result.Supplemental = append(result.Supplemental, mapped)
	}
	for _, d := range res.Diagnostics {
		if d.Start < 0 || d.Length < 0 || d.Start > int(^uint(0)>>1)-d.Length {
			return Result{}, fmt.Errorf("invalid content mapper diagnostic range [%d, %d)", d.Start, d.Start+d.Length)
		}
		start, err := originalPositions.normalize(d.Start)
		if err != nil {
			return Result{}, fmt.Errorf("invalid content mapper diagnostic start: %w", err)
		}
		end, err := originalPositions.normalize(d.Start + d.Length)
		if err != nil {
			return Result{}, fmt.Errorf("invalid content mapper diagnostic end: %w", err)
		}
		result.Diagnostics = append(result.Diagnostics, ast.NewExternalDiagnostic(
			nil,
			core.NewTextRange(start, end),
			diagnosticSource,
			diagnostics.CategoryError,
			d.Code,
			d.MessageText,
		))
	}
	return result, nil
}

func decodeMappedOutput(output MappedOutput, originalText string, positionEncoding PositionEncoding) (MappedResult, *positionNormalizer, error) {
	result := MappedResult{
		Text: output.Text,
	}
	generatedPositions, err := newPositionNormalizer(output.Text, positionEncoding)
	if err != nil {
		return MappedResult{}, nil, err
	}
	originalPositions, err := newPositionNormalizer(originalText, positionEncoding)
	if err != nil {
		return MappedResult{}, nil, err
	}
	// A successful transform always carries a span map. Absent or empty mappings describe fully
	// synthesized output (no segment corresponds to the original), so decode to an empty map rather than
	// nil, which would mean "not content-mapped".
	if len(output.Mappings) > 0 {
		mappings, err := spanmap.Unmarshal(output.Mappings)
		if err != nil {
			return MappedResult{}, nil, err
		}
		result.Mappings, err = normalizeMappings(mappings, generatedPositions, originalPositions)
		if err != nil {
			return MappedResult{}, nil, err
		}
	} else {
		result.Mappings = spanmap.New(nil)
	}
	return result, originalPositions, nil
}

func normalizeMappings(mappings *spanmap.SpanMap, generatedPositions *positionNormalizer, originalPositions *positionNormalizer) (*spanmap.SpanMap, error) {
	segments := mappings.Segments()
	for i := range segments {
		segment := &segments[i]
		var err error
		segment.GenStart, err = generatedPositions.normalizeTextPos(segment.GenStart)
		if err != nil {
			return nil, fmt.Errorf("invalid content mapper mapping %d generated start: %w", i, err)
		}
		segment.GenEnd, err = generatedPositions.normalizeTextPos(segment.GenEnd)
		if err != nil {
			return nil, fmt.Errorf("invalid content mapper mapping %d generated end: %w", i, err)
		}
		segment.OrigStart, err = originalPositions.normalizeTextPos(segment.OrigStart)
		if err != nil {
			return nil, fmt.Errorf("invalid content mapper mapping %d original start: %w", i, err)
		}
		segment.OrigEnd, err = originalPositions.normalizeTextPos(segment.OrigEnd)
		if err != nil {
			return nil, fmt.Errorf("invalid content mapper mapping %d original end: %w", i, err)
		}
	}
	return spanmap.New(segments), nil
}

type positionNormalizer struct {
	text        string
	encoding    PositionEncoding
	positionMap *ast.PositionMap
	length      int
}

func newPositionNormalizer(text string, encoding PositionEncoding) (*positionNormalizer, error) {
	normalizer := &positionNormalizer{text: text, encoding: encoding}
	switch encoding {
	case PositionEncodingUTF8:
		normalizer.length = len(text)
	case PositionEncodingUTF16:
		normalizer.positionMap = ast.ComputePositionMap(text)
		normalizer.length = normalizer.positionMap.UTF8ToUTF16(len(text))
	default:
		return nil, fmt.Errorf("unsupported position encoding %q", encoding)
	}
	return normalizer, nil
}

func (n *positionNormalizer) normalizeTextPos(position core.TextPos) (core.TextPos, error) {
	normalized, err := n.normalize(int(position))
	return core.TextPos(normalized), err
}

func (n *positionNormalizer) normalize(position int) (int, error) {
	if position < 0 {
		return 0, fmt.Errorf("position %d is negative", position)
	}
	if position > n.length {
		return 0, fmt.Errorf("position %d exceeds %s length %d", position, n.encoding, n.length)
	}
	var bytePosition int
	switch n.encoding {
	case PositionEncodingUTF8:
		bytePosition = position
	case PositionEncodingUTF16:
		bytePosition = n.positionMap.UTF16ToUTF8(position)
	}
	if bytePosition < len(n.text) && !utf8.RuneStart(n.text[bytePosition]) {
		return 0, fmt.Errorf("position %d splits a Unicode code point", position)
	}
	return bytePosition, nil
}

// rejectHandler rejects any request initiated by the mapper. The content mapper protocol is currently
// parent-driven only; a request from the child is a protocol violation.
type rejectHandler struct{}

func (rejectHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	return nil, fmt.Errorf("content mapper sent an unexpected request: %s", method)
}

func (rejectHandler) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}
