//go:build darwin && amd64

#include "textflag.h"

// fsevents_darwin_ffi_amd64.s: amd64 assembly for the FSEvents backend
//
// Contains two functions:
//
//  1. FSEventStreamCreate trampoline: shuffles the float64 latency arg
//     from R9 (integer register, where syscall6 puts it) into X0 (xmm0,
//     where the System V AMD64 ABI expects the first float argument),
//     and hardcodes the flags argument to 0x11
//     (kFSEventStreamCreateFlagUseCFTypes |
//     kFSEventStreamCreateFlagFileEvents).
//
//  2. fsEventsCallbackASM: the C-convention callback invoked by FSEvents
//     on a GCD dispatch queue thread. Retains/copies callback data into a
//     payload, writes the payload pointer to eventPipe to wake the Go event-loop
//     goroutine, then blocks on donePipe until processing is complete. Never
//     enters Go ABI; stays entirely in System V AMD64 calling convention.

// ---------------------------------------------------------------------------
// FSEventStreamCreate trampoline: shuffles the float64 latency argument.
//
// The runtime's syscall6 trampoline loads 6 args into registers:
//   DI=allocator  SI=callback  DX=ctx  CX=paths
//   R8=sinceWhen  R9=latency(bits)
//
// The C function expects latency in X0 (xmm0) and flags in R9.
// flags is always 0x11 (kFSEventStreamCreateFlagUseCFTypes |
// kFSEventStreamCreateFlagFileEvents), so we hardcode it.
// ---------------------------------------------------------------------------

TEXT fse_FSEventStreamCreate_trampoline<>(SB), NOSPLIT, $0-0
	MOVQ R9, X0
	MOVQ $0x11, R9
	JMP  fse_FSEventStreamCreate(SB)

GLOBL ·fse_FSEventStreamCreate_trampoline_addr(SB), RODATA, $8
DATA ·fse_FSEventStreamCreate_trampoline_addr(SB)/8, $fse_FSEventStreamCreate_trampoline<>(SB)

// ---------------------------------------------------------------------------
// FSEvents callback: called from a GCD dispatch queue with C convention.
//   DI=streamRef  SI=info  DX=numEvents  CX=paths  R8=flags  R9=ids
//
// `info` is a pointer to a streamCallback struct (see fsevents_darwin_ffi.go):
//   offset  0: eventPipeWrite fd    (8 bytes)
//   offset  8: donePipeRead fd      (8 bytes)
//   offset 16: active callback flag (8 bytes)
//   offset 24: closing flag         (8 bytes)
//
// Stays entirely in C context (no cgocallback). Saves args to the per-stream
// heap-allocated payload, writes its pointer to the stream's eventPipe to wake
// its Go event loop goroutine, then blocks on the stream's donePipe until
// processing is complete.
//
// NOFRAME: this function is entered from C, not Go. We manage the frame
// ourselves following the System V AMD64 ABI.
//
// Frame layout (104 bytes, 16-byte aligned):
//   On entry from C, RSP ≡ 8 mod 16 (return address pushed by CALL).
//   SUB $104 → RSP ≡ 8-104 = 0 mod 16, aligned for CALL into libc.
//   RSP+ 0: payload pointer bytes written to eventPipe
//   RSP+ 8: scratch byte for donePipe read
//   RSP+16: saved info pointer
//   RSP+24: saved numEvents
//   RSP+32: saved original flags pointer
//   RSP+40: retained CFArray paths
//   RSP+48: copied flags pointer
//   RSP+96: saved RBP  ← BP points here (C frame chain)
//   RSP+104: return address (pushed by C's CALL)
// ---------------------------------------------------------------------------

TEXT fsEventsCallbackASM<>(SB), NOSPLIT|NOFRAME, $0
	SUBQ $104, SP
	MOVQ BP, 96(SP)
	LEAQ 96(SP), BP

	MOVQ SI, 16(SP)    // info
	MOVQ DX, 24(SP)    // numEvents
	MOVQ R8, 32(SP)    // original flags

	// Dekker handshake with Go's stopStream: publish activeCallback=1, then
	// observe closing. MFENCE prevents the load of closing from being
	// reordered ahead of the store to active by the store buffer; without
	// it Go could observe activeCallback==0 in waitIdle even though this
	// callback is about to proceed past the closing check, leading to
	// fd-close races.
	MOVQ $1, (2*8)(SI)
	MFENCE
	MOVQ (3*8)(SI), AX
	CMPQ AX, $0
	JNE  done

	// Retain the CFArray paths because FSEvents owns the callback argument.
	MOVQ CX, DI
	XORL AX, AX
	CALL fse_CFRetain(SB)
	CMPQ AX, $0
	JEQ  done
	MOVQ AX, 40(SP)

	// Copy the flags array into C heap memory owned by the Go event loop.
	MOVQ 24(SP), DI
	SHLQ $2, DI
	XORL AX, AX
	CALL fse_malloc(SB)
	CMPQ AX, $0
	JEQ  releasePaths
	MOVQ AX, 48(SP)

	MOVQ AX, DI
	MOVQ 32(SP), SI
	MOVQ 24(SP), DX
	SHLQ $2, DX
	XORL AX, AX
	CALL fse_memcpy(SB)

	// Allocate and populate fsEventsCallbackPayload.
	MOVQ $24, DI
	XORL AX, AX
	CALL fse_malloc(SB)
	CMPQ AX, $0
	JEQ  freeFlags
	MOVQ AX, 0(SP)

	MOVQ 24(SP), CX
	MOVQ CX, (0*8)(AX)
	MOVQ 40(SP), CX
	MOVQ CX, (1*8)(AX)
	MOVQ 48(SP), CX
	MOVQ CX, (2*8)(AX)

	// write(info->eventPipeWrite, &payload, sizeof(payload)).
writeAgain:
	MOVQ 16(SP), AX      // reload info
	MOVQ (0*8)(AX), DI   // eventPipeWrite
	LEAQ 0(SP), SI       // buf (payload pointer)
	MOVQ $8, DX          // count
	XORL AX, AX          // no float args
	CALL fse_write(SB)
	CMPQ AX, $8
	JEQ  waitDone
	CMPQ AX, $-1
	JNE  freePayload
	XORL AX, AX          // no float args
	CALL fse___error(SB)
	MOVL (AX), AX
	CMPL AX, $4          // EINTR
	JEQ  writeAgain
	JMP  freePayload

// read(info->donePipeRead, &buf, 1): block until Go is done.
waitDone:
readAgain:
	MOVQ 16(SP), AX      // reload info
	MOVQ (1*8)(AX), DI   // donePipeRead
	LEAQ 8(SP), SI       // buf (scratch area)
	MOVQ $1, DX          // count
	XORL AX, AX          // no float args
	CALL fse_read(SB)
	CMPQ AX, $1
	JEQ  done
	CMPQ AX, $-1
	JNE  done
	XORL AX, AX          // no float args
	CALL fse___error(SB)
	MOVL (AX), AX
	CMPL AX, $4          // EINTR
	JEQ  readAgain
	JMP  done

freePayload:
	MOVQ 0(SP), DI
	XORL AX, AX
	CALL fse_free(SB)

freeFlags:
	MOVQ 48(SP), DI
	XORL AX, AX
	CALL fse_free(SB)

releasePaths:
	MOVQ 40(SP), DI
	XORL AX, AX
	CALL fse_CFRelease(SB)

	// Return 0.
done:
	MOVQ 16(SP), AX
	MOVQ $0, (2*8)(AX)
	XORL AX, AX
	MOVQ 96(SP), BP
	ADDQ $104, SP
	RET

GLOBL ·fsEventsCallbackAsmAddr(SB), RODATA, $8
DATA ·fsEventsCallbackAsmAddr(SB)/8, $fsEventsCallbackASM<>(SB)
