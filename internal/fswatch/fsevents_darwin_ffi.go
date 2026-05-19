//go:build darwin && (amd64 || arm64)

package fswatch

import (
	"io"
	"math"
	"os"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ---------------------------------------------------------------------------
// fsevents_darwin_ffi.go: cgo-free macOS CoreFoundation / CoreServices FFI
//
// Provides Go access to Apple's FSEvents, CoreFoundation, and libdispatch
// frameworks entirely without cgo, following the pattern established by
// crypto/x509/internal/macos in the Go standard library.
//
// Imported symbols include CoreFoundation helpers (CFRelease,
// CFStringCreateWithCString, CFArrayCreate), libdispatch
// (dispatch_queue_create), and CoreServices FSEvents functions
// (FSEventStreamCreate, SetDispatchQueue, Start, Stop, Invalidate,
// Release).
//
// Each framework symbol has three parts:
//  1. //go:cgo_import_dynamic: tells the linker to import the C symbol
//     from a shared library (CoreFoundation.framework, CoreServices.framework,
//     or libSystem.B.dylib).
//  2. A TEXT trampoline in the .s file: a minimal assembly stub that JMPs
//     to the imported symbol. For simple functions this is a bare JMP; for
//     FSEventStreamCreate the trampoline also moves the float64 latency
//     argument from an integer register to a float register.
//  3. A GLOBL/DATA pair that exports the trampoline's ABI0 address as a Go
//     uintptr variable (·fse_X_trampoline_addr), which the Go wrapper
//     passes to runtime's syscall_syscall6.
//
//	┌──────────────────────────────────────────────────────────┐
//	│  Go wrapper: cfRelease(ref)                              │
//	│    syscall_syscall6(trampoline_addr, ref, ...)           │
//	│         │                                                │
//	│         ▼                                                │
//	│  ┌──────────────────────────────────┐                    │
//	│  │ .s trampoline (ABI0)             │                    │
//	│  │  fse_CFRelease_trampoline<>:     │                    │
//	│  │    JMP fse_CFRelease(SB)         │                    │
//	│  └─────────────┬────────────────────┘                    │
//	│                │                                         │
//	│                ▼                                         │
//	│  ┌──────────────────────────────────┐                    │
//	│  │ //go:cgo_import_dynamic          │                    │
//	│  │ CFRelease from CoreFoundation    │                    │
//	│  └──────────────────────────────────┘                    │
//	└──────────────────────────────────────────────────────────┘
//
// FSEvents callback synchronization (per-stream):
//
//	  GCD dispatch queue thread          Go goroutine (eventLoop)
//	  ─────────────────────────          ────────────────────────
//	  FSEvents fires C callback
//	  on a libdispatch OS thread
//	         │
//	  ┌──────▼──────────────────┐
//	  │ asm: retain CFArray     │
//	  │ paths, copy flags,      │
//	  │ allocate payload        │
//	  └──────┬──────────────────┘
//	         │
//	  write(eventPipeWrite, payload*) ─► read(eventFile) unblocks
//	         │                                │
//	         │  (GCD thread blocked)   fsEventsCallback(cb, payload)
//	         │                         classifies events,
//	         │                         frees payload,
//	         │                         posts to dirWatch.events
//	         │                                │
//	  read(donePipeRead)  ◄──────────── write(donePipeWrite)
//	         │
//	  asm: return to FSEvents
//	  (GCD thread released)
//
// The assembly callback never enters Go ABI; it stays entirely in C
// context. Two pipe pairs per stream (eventPipe and donePipe) synchronize
// the C dispatch queue thread with a dedicated Go event-loop goroutine.
// The Go side uses os.File.Read (integrated with netpoll/kqueue on macOS)
// so the goroutine parks efficiently without blocking an OS thread.
//
// streamCallback memory layout (must match assembly offsets):
//
//	offset  0: eventPipeWrite fd        ─┐ Read by asm to call
//	offset  8: donePipeRead fd          ─┘ write() / read()
//	offset 16: active callback flag
//	offset 24: closing flag
// ---------------------------------------------------------------------------

// Framework linker flags for the external linker.
// Note: //go:cgo_ldflag is only valid in cgo-generated code. The
// //go:cgo_import_dynamic directives below are sufficient: the Go
// linker records the framework paths in the Mach-O LC_LOAD_DYLIB
// commands automatically.

// Implemented in the runtime package (runtime/sys_darwin.go).
// These are the same linknames that golang.org/x/sys/unix uses.
//
//go:linkname syscall_syscall6 syscall.syscall6
func syscall_syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// ---------------------------------------------------------------------------
// CoreFoundation imports, trampoline addresses, and Go wrappers.
//
// Each function groups its //go:cgo_import_dynamic directive, its
// trampoline address variable (populated by GLOBL/DATA in the .s files),
// and its Go wrapper together.
// ---------------------------------------------------------------------------

//go:cgo_import_dynamic fse_CFRelease CFRelease "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFRelease_trampoline_addr uintptr

func cfRelease(ref uintptr) {
	_, _, _ = syscall_syscall6(fse_CFRelease_trampoline_addr, ref, 0, 0, 0, 0, 0)
}

//go:cgo_import_dynamic fse_CFStringCreateWithCString CFStringCreateWithCString "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFStringCreateWithCString_trampoline_addr uintptr

func cfStringCreate(allocator uintptr, cstr unsafe.Pointer, encoding uint32) uintptr {
	ret, _, _ := syscall_syscall6(fse_CFStringCreateWithCString_trampoline_addr, allocator, uintptr(cstr), uintptr(encoding), 0, 0, 0)
	runtime.KeepAlive(cstr)
	return ret
}

//go:cgo_import_dynamic fse_CFArrayCreate CFArrayCreate "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFArrayCreate_trampoline_addr uintptr

func cfArrayCreate(allocator uintptr, values unsafe.Pointer, count int, callbacks uintptr) uintptr {
	ret, _, _ := syscall_syscall6(fse_CFArrayCreate_trampoline_addr, allocator, uintptr(values), uintptr(count), callbacks, 0, 0)
	runtime.KeepAlive(values)
	return ret
}

//go:cgo_import_dynamic fse_CFArrayGetValueAtIndex CFArrayGetValueAtIndex "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFArrayGetValueAtIndex_trampoline_addr uintptr

func cfArrayGetValueAtIndex(array uintptr, index int) uintptr {
	ret, _, _ := syscall_syscall6(fse_CFArrayGetValueAtIndex_trampoline_addr, array, uintptr(index), 0, 0, 0, 0)
	return ret
}

// ----- NFC normalization helpers -----
//
// FSEvents reports paths using whatever bytes are stored on disk. APFS is
// normalization-insensitive for lookups (a file created as NFD opens fine
// under the NFC form, and vice versa) but it stores and reports the original
// bytes. The library normalizes every path that crosses the darwin boundary
// to Unicode NFC so that:
//   - WatchDirectory("/.../caf\u00e9") and WatchDirectory("/.../cafe\u0301")
//     coalesce to a single dir watch;
//   - WatchFile filters by exact-string compare in NFC always match;
//   - subscribers can compare event paths against their own NFC strings.
//
// All-ASCII inputs are bit-identical in NFC and NFD, so the hot path skips
// the FFI entirely. The rare non-ASCII case round-trips through CoreFoundation
// (UTF-8 → CFString → CFMutableString → CFStringNormalize → UTF-8) with no Go
// Unicode tables, no extra dependency.

const (
	cfStringNormalizationFormC = 2 // kCFStringNormalizationFormC
)

//go:cgo_import_dynamic fse_CFStringCreateMutableCopy CFStringCreateMutableCopy "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFStringCreateMutableCopy_trampoline_addr uintptr

func cfStringCreateMutableCopy(allocator uintptr, maxLength int, str uintptr) uintptr {
	ret, _, _ := syscall_syscall6(fse_CFStringCreateMutableCopy_trampoline_addr, allocator, uintptr(maxLength), str, 0, 0, 0)
	return ret
}

//go:cgo_import_dynamic fse_CFStringNormalize CFStringNormalize "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFStringNormalize_trampoline_addr uintptr

func cfStringNormalize(mutStr uintptr, form uintptr) {
	_, _, _ = syscall_syscall6(fse_CFStringNormalize_trampoline_addr, mutStr, form, 0, 0, 0, 0)
}

//go:cgo_import_dynamic fse_CFStringGetLength CFStringGetLength "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFStringGetLength_trampoline_addr uintptr

func cfStringGetLength(str uintptr) int {
	ret, _, _ := syscall_syscall6(fse_CFStringGetLength_trampoline_addr, str, 0, 0, 0, 0, 0)
	return int(ret)
}

//go:cgo_import_dynamic fse_CFStringGetMaximumSizeForEncoding CFStringGetMaximumSizeForEncoding "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFStringGetMaximumSizeForEncoding_trampoline_addr uintptr

func cfStringGetMaximumSizeForEncoding(length int, encoding uint32) int {
	ret, _, _ := syscall_syscall6(fse_CFStringGetMaximumSizeForEncoding_trampoline_addr, uintptr(length), uintptr(encoding), 0, 0, 0, 0)
	return int(ret)
}

//go:cgo_import_dynamic fse_CFStringGetCString CFStringGetCString "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

var fse_CFStringGetCString_trampoline_addr uintptr

func cfStringGetCString(str uintptr, buf unsafe.Pointer, bufSize int, encoding uint32) bool {
	ret, _, _ := syscall_syscall6(fse_CFStringGetCString_trampoline_addr, str, uintptr(buf), uintptr(bufSize), uintptr(encoding), 0, 0)
	runtime.KeepAlive(buf)
	return ret != 0
}

// isASCII reports whether every byte in s is below 0x80. Pure-ASCII paths
// are identical in every Unicode normalization form, so we can skip the
// CoreFoundation round-trip entirely, which is the overwhelming common case.
func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// cfStringToNFC returns the CFString at src as a NFC-normalized Go string.
// If normalization fails, it falls back to the unnormalized UTF-8 contents.
// Returns "" only if both the normalized and unnormalized conversions fail
// (e.g. src is not a CFString, or allocation fails).
func cfStringToNFC(src uintptr) string {
	if src == 0 {
		return ""
	}
	if s := cfStringNormalizedToGo(src); s != "" {
		return s
	}
	return cfStringToGo(src)
}

// cfStringNormalizedToGo returns the CFString at src as a NFC-normalized Go
// string, or "" on any failure.
func cfStringNormalizedToGo(src uintptr) string {
	mut := cfStringCreateMutableCopy(0, 0, src)
	if mut == 0 {
		return ""
	}
	defer cfRelease(mut)

	cfStringNormalize(mut, cfStringNormalizationFormC)
	return cfStringToGo(mut)
}

// cfStringToGo extracts the UTF-8 contents of the CFString at src as a Go
// string, or "" on failure.
func cfStringToGo(src uintptr) string {
	length := cfStringGetLength(src)
	bufSize := cfStringGetMaximumSizeForEncoding(length, cfStringEncodingUTF8) + 1
	buf := make([]byte, bufSize)
	if !cfStringGetCString(src, unsafe.Pointer(&buf[0]), bufSize, cfStringEncodingUTF8) {
		return ""
	}
	// CFStringGetCString writes a NUL terminator; trim it.
	n := 0
	for n < len(buf) && buf[n] != 0 {
		n++
	}
	return string(buf[:n])
}

// normalizeNFC returns s in Unicode NFC (canonical composed) form. ASCII
// inputs are returned unchanged. Non-ASCII inputs go through CoreFoundation;
// if any step fails (e.g. invalid UTF-8 from a corrupt path), the original
// string is returned so the caller still sees *something* rather than nothing.
func normalizeNFC(s string) string {
	if isASCII(s) {
		return s
	}

	cstr := append([]byte(s), 0)
	src := cfStringCreate(0, unsafe.Pointer(&cstr[0]), cfStringEncodingUTF8)
	runtime.KeepAlive(cstr)
	if src == 0 {
		return s
	}
	defer cfRelease(src)

	normalized := cfStringToNFC(src)
	if normalized == "" {
		return s
	}
	return normalized
}

// ---------------------------------------------------------------------------
// libdispatch imports.
// ---------------------------------------------------------------------------

//go:cgo_import_dynamic fse_dispatch_queue_create dispatch_queue_create "/usr/lib/libSystem.B.dylib"

var fse_dispatch_queue_create_trampoline_addr uintptr

func dispatchQueueCreate(label unsafe.Pointer) uintptr {
	ret, _, _ := syscall_syscall6(fse_dispatch_queue_create_trampoline_addr, uintptr(label), 0, 0, 0, 0, 0)
	runtime.KeepAlive(label)
	return ret
}

//go:cgo_import_dynamic fse_dispatch_release dispatch_release "/usr/lib/libSystem.B.dylib"

var fse_dispatch_release_trampoline_addr uintptr

func dispatchRelease(obj uintptr) {
	_, _, _ = syscall_syscall6(fse_dispatch_release_trampoline_addr, obj, 0, 0, 0, 0, 0)
}

// ---------------------------------------------------------------------------
// CoreServices / FSEvents imports, trampoline addresses, and Go wrappers.
// ---------------------------------------------------------------------------

//go:cgo_import_dynamic fse_FSEventStreamCreate FSEventStreamCreate "/System/Library/Frameworks/CoreServices.framework/Versions/A/CoreServices"

var fse_FSEventStreamCreate_trampoline_addr uintptr // arch-specific trampoline

func fsEventStreamCreate(allocator, callback uintptr, ctx unsafe.Pointer, paths uintptr, since uint64, latency float64) uintptr {
	// syscall_syscall6 only carries 6 integer args. The arch-specific
	// trampoline moves the latency bits from an integer register to the
	// float register and hardcodes flags =
	// kFSEventStreamCreateFlagUseCFTypes | kFSEventStreamCreateFlagFileEvents (0x11).
	ret, _, _ := syscall_syscall6(
		fse_FSEventStreamCreate_trampoline_addr,
		allocator,
		callback,
		uintptr(ctx),
		paths,
		uintptr(since),
		uintptr(math.Float64bits(latency)),
	)
	runtime.KeepAlive(ctx)
	return ret
}

//go:cgo_import_dynamic fse_FSEventStreamSetDispatchQueue FSEventStreamSetDispatchQueue "/System/Library/Frameworks/CoreServices.framework/Versions/A/CoreServices"

var fse_FSEventStreamSetDispatchQueue_trampoline_addr uintptr

func fsEventStreamSetDispatchQueue(stream, queue uintptr) {
	_, _, _ = syscall_syscall6(fse_FSEventStreamSetDispatchQueue_trampoline_addr, stream, queue, 0, 0, 0, 0)
}

//go:cgo_import_dynamic fse_FSEventStreamStart FSEventStreamStart "/System/Library/Frameworks/CoreServices.framework/Versions/A/CoreServices"

var fse_FSEventStreamStart_trampoline_addr uintptr

func fsEventStreamStart(stream uintptr) uint8 {
	r1, _, _ := syscall_syscall6(fse_FSEventStreamStart_trampoline_addr, stream, 0, 0, 0, 0, 0)
	return uint8(r1)
}

//go:cgo_import_dynamic fse_FSEventStreamFlushSync FSEventStreamFlushSync "/System/Library/Frameworks/CoreServices.framework/Versions/A/CoreServices"

var fse_FSEventStreamFlushSync_trampoline_addr uintptr

func fsEventStreamFlushSync(stream uintptr) {
	_, _, _ = syscall_syscall6(fse_FSEventStreamFlushSync_trampoline_addr, stream, 0, 0, 0, 0, 0)
}

//go:cgo_import_dynamic fse_FSEventStreamStop FSEventStreamStop "/System/Library/Frameworks/CoreServices.framework/Versions/A/CoreServices"

var fse_FSEventStreamStop_trampoline_addr uintptr

func fsEventStreamStop(stream uintptr) {
	_, _, _ = syscall_syscall6(fse_FSEventStreamStop_trampoline_addr, stream, 0, 0, 0, 0, 0)
}

//go:cgo_import_dynamic fse_FSEventStreamInvalidate FSEventStreamInvalidate "/System/Library/Frameworks/CoreServices.framework/Versions/A/CoreServices"

var fse_FSEventStreamInvalidate_trampoline_addr uintptr

func fsEventStreamInvalidate(stream uintptr) {
	_, _, _ = syscall_syscall6(fse_FSEventStreamInvalidate_trampoline_addr, stream, 0, 0, 0, 0, 0)
}

//go:cgo_import_dynamic fse_FSEventStreamRelease FSEventStreamRelease "/System/Library/Frameworks/CoreServices.framework/Versions/A/CoreServices"

var fse_FSEventStreamRelease_trampoline_addr uintptr

func fsEventStreamRelease(stream uintptr) {
	_, _, _ = syscall_syscall6(fse_FSEventStreamRelease_trampoline_addr, stream, 0, 0, 0, 0, 0)
}

// ---------------------------------------------------------------------------
// Direct callback assembly imports.
// ---------------------------------------------------------------------------

// These symbols are called directly by fsEventsCallbackASM and have no Go
// wrappers.
//go:cgo_import_dynamic fse_CFRetain CFRetain "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic fse_write write "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic fse_read read "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic fse___error __error "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic fse_malloc malloc "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic fse_memcpy memcpy "/usr/lib/libSystem.B.dylib"

// ---------------------------------------------------------------------------
// libSystem imports, trampoline addresses, and Go wrappers.
// ---------------------------------------------------------------------------

//go:cgo_import_dynamic fse_free free "/usr/lib/libSystem.B.dylib"

var fse_free_trampoline_addr uintptr

func libcFree(ptr uintptr) {
	if ptr != 0 {
		_, _, _ = syscall_syscall6(fse_free_trampoline_addr, ptr, 0, 0, 0, 0, 0)
	}
}

// ---------------------------------------------------------------------------
// Callback address.
// ---------------------------------------------------------------------------

// fsEventsCallbackAsmAddr is the address of the arch-specific callback
// function defined in fsevents_darwin_ffi_{amd64,arm64}.s.
var fsEventsCallbackAsmAddr uintptr

// ---------------------------------------------------------------------------
// Per-stream callback infrastructure
// ---------------------------------------------------------------------------

// streamCallback is the per-stream buffer shared between the C callback
// assembly and the Go event loop goroutine. The assembly receives a pointer
// to this struct as the FSEventStreamContext.info parameter and uses offset
// addressing to access the pipe fds and lifecycle flags.
//
// The struct layout must match the assembly (fsevents_darwin_ffi_{amd64,arm64}.s):
//
//	offset  0: eventPipeWrite fd
//	offset  8: donePipeRead fd
//	offset 16: active callback flag
//	offset 24: closing flag
type streamCallback struct {
	eventPipeWrite uintptr
	donePipeRead   uintptr
	activeCallback uintptr
	closing        uintptr

	// Go-only fields (not accessed by assembly, offset doesn't matter).
	eventFile     *os.File
	donePipeWrite int
	queue         uintptr // per-stream serial dispatch queue
	done          chan struct{}
	dirWatch      *dirWatch
	closed        atomic.Bool
}

type fsEventsCallbackPayload struct {
	numEvents uintptr
	paths     uintptr
	flags     uintptr
}

const closedFD = ^uintptr(0)

func (p *fsEventsCallbackPayload) close() {
	if p == nil {
		return
	}
	if p.paths != 0 {
		cfRelease(p.paths)
	}
	libcFree(p.flags)
	libcFree(uintptr(unsafe.Pointer(p)))
}

// newStreamCallback allocates a pinned streamCallback with its own pipe pair
// and per-stream serial dispatch queue, and starts a goroutine to process
// callbacks. The per-stream serial queue both serializes this stream's
// callbacks (the asm shim's activeCallback flag and pipe handshake are not
// re-entrant) and prevents cross-stream head-of-line blocking that a
// process-wide serial queue would cause.
func newStreamCallback(w *dirWatch) (*streamCallback, error) {
	var eventPipe, donePipe [2]int
	if err := unix.Pipe(eventPipe[:]); err != nil {
		return nil, err
	}
	unix.CloseOnExec(eventPipe[0])
	unix.CloseOnExec(eventPipe[1])
	if err := unix.Pipe(donePipe[:]); err != nil {
		unix.Close(eventPipe[0])
		unix.Close(eventPipe[1])
		return nil, err
	}
	unix.CloseOnExec(donePipe[0])
	unix.CloseOnExec(donePipe[1])

	label := []byte("typescript.fswatch.fsevents.stream\x00")
	queue := dispatchQueueCreate(unsafe.Pointer(&label[0]))
	runtime.KeepAlive(label)
	if queue == 0 {
		unix.Close(eventPipe[0])
		unix.Close(eventPipe[1])
		unix.Close(donePipe[0])
		unix.Close(donePipe[1])
		return nil, errStreamCreateNull
	}

	cb := &streamCallback{
		eventPipeWrite: uintptr(eventPipe[1]),
		donePipeRead:   uintptr(donePipe[0]),
		eventFile:      os.NewFile(uintptr(eventPipe[0]), "fsevents-event"),
		donePipeWrite:  donePipe[1],
		queue:          queue,
		done:           make(chan struct{}),
		dirWatch:       w,
	}
	go cb.eventLoop()
	return cb, nil
}

func (cb *streamCallback) beginClose() {
	atomic.StoreUintptr(&cb.closing, 1)
}

func (cb *streamCallback) waitIdle() {
	// Yield first: in the common case the asm callback unblocks on
	// donePipeRead's close immediately and we never even need to sleep.
	// If the asm goroutine has been preempted, escalate to a real sleep
	// so we don't pin a P spinning on the atomic load.
	const yieldRounds = 16
	for i := 0; atomic.LoadUintptr(&cb.activeCallback) != 0; i++ {
		if i < yieldRounds {
			runtime.Gosched()
			continue
		}
		time.Sleep(time.Millisecond)
	}
}

func (cb *streamCallback) unblockCallbacks() {
	cb.closeDonePipeRead()
}

func (cb *streamCallback) closeDonePipeRead() {
	fd := atomic.SwapUintptr(&cb.donePipeRead, closedFD)
	if fd != closedFD {
		unix.Close(int(fd))
	}
}

// close shuts down the event loop goroutine and releases resources.
func (cb *streamCallback) close() {
	unix.Close(int(cb.eventPipeWrite))
	<-cb.done
	cb.eventFile.Close()
	cb.closeDonePipeRead()
	unix.Close(cb.donePipeWrite)
	if cb.queue != 0 {
		dispatchRelease(cb.queue)
		cb.queue = 0
	}
}

// eventLoop runs on a dedicated goroutine for this stream. It reads signals
// from the callback assembly (via eventPipe), processes the event using the
// retained/copied payload, then signals completion (via donePipe).
// The eventFile.Read() call integrates with Go's netpoll (kqueue on macOS),
// so the goroutine parks without blocking an OS thread while idle.
func (cb *streamCallback) eventLoop() {
	defer close(cb.done)
	var payload *fsEventsCallbackPayload
	buf := unsafe.Slice((*byte)(unsafe.Pointer(&payload)), unsafe.Sizeof(payload))
	done := []byte{0}
	for {
		payload = nil
		if _, err := io.ReadFull(cb.eventFile, buf); err != nil {
			return // pipe closed or error → shutdown
		}

		fsEventsCallback(cb, payload)

		if _, err := ignoringEINTR(func() (int, error) { return unix.Write(cb.donePipeWrite, done) }); err != nil {
			return
		}
	}
}
