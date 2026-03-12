package napi

/*
#include <stdlib.h>
#include <string.h>
#include "node/node_api.h"
napi_value goNapiCallbackTrampoline(napi_env env, napi_callback_info info);

static inline napi_status go_napi_create_function(
	napi_env env,
	const char* utf8name,
	size_t length,
	napi_callback cb,
	uintptr_t handle,
	napi_value* result)
{
	return napi_create_function(env, utf8name, length, cb, (void*)handle, result);
}

static inline napi_status go_napi_get_cb_info(
	napi_env env,
	napi_callback_info cbinfo,
	size_t* argc,
	napi_value* argv,
	napi_value* this_arg,
	uintptr_t* handle)
{
	return napi_get_cb_info(env, cbinfo, argc, argv, this_arg, (void**)handle);
}
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

// InitFunc is the function signature for module initialization.
type InitFunc func(env Env, exports Value) (Value, error)

var initFn InitFunc

// RegisterInit registers the module initialization function.
// This must be called before the module is loaded by Node.js.
func RegisterInit(fn InitFunc) {
	if initFn != nil {
		panic("RegisterInit already called")
	}
	initFn = fn
}

//export goNapiInit
func goNapiInit(env C.napi_env, exports C.napi_value) C.napi_value {
	e := Env{raw: env}
	if initFn == nil {
		panic("RegisterInit not called")
	}
	result, err := initFn(e, Value{raw: exports})
	if err != nil {
		panic(err)
	}
	return result.raw
}

// Env wraps a napi_env handle.
type Env struct {
	raw C.napi_env
}

// Value wraps a napi_value handle.
type Value struct {
	raw C.napi_value
}

// IsExceptionPending returns true if an exception is pending on the environment.
func (e Env) IsExceptionPending() (bool, error) {
	var result C.bool
	status := C.napi_is_exception_pending(e.raw, &result)
	if err := napiStatusToError(status); err != nil {
		return false, err
	}
	return bool(result), nil
}

// GetBooleanValue returns a NAPI boolean value.
func (e Env) GetBooleanValue(value bool) (Value, error) {
	var result C.napi_value
	status := C.napi_get_boolean(e.raw, C.bool(value), &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// BooleanValueToBool converts a NAPI boolean value to a Go bool.
func (e Env) BooleanValueToBool(value Value) (bool, error) {
	var result C.bool
	status := C.napi_get_value_bool(e.raw, value.raw, &result)
	if err := napiStatusToError(status); err != nil {
		return false, err
	}
	return bool(result), nil
}

// GetGlobal returns the global object.
func (e Env) GetGlobal() (Value, error) {
	var result C.napi_value
	status := C.napi_get_global(e.raw, &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// GetNullValue returns the null value.
func (e Env) GetNullValue() (Value, error) {
	var result C.napi_value
	status := C.napi_get_null(e.raw, &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// GetUndefinedValue returns the undefined value.
func (e Env) GetUndefinedValue() (Value, error) {
	var result C.napi_value
	status := C.napi_get_undefined(e.raw, &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// StringToStringValue converts a Go string to a NAPI string value.
func (e Env) StringToStringValue(str string) (Value, error) {
	var result C.napi_value
	status := C.napi_create_string_utf8(e.raw, (*C.char)(unsafe.Pointer(unsafe.StringData(str))), C.size_t(len(str)), &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// StringValueToString converts a NAPI string value to a Go string.
func (e Env) StringValueToString(str Value) (string, error) {
	var size C.size_t
	status := C.napi_get_value_string_utf8(e.raw, str.raw, nil, 0, &size)
	if err := napiStatusToError(status); err != nil {
		return "", err
	}

	sizePlusOne := size + 1
	buffer := make([]byte, sizePlusOne)

	status = C.napi_get_value_string_utf8(e.raw, str.raw, (*C.char)(unsafe.Pointer(unsafe.SliceData(buffer))), sizePlusOne, nil)
	if err := napiStatusToError(status); err != nil {
		return "", err
	}
	return string(buffer[:size]), nil
}

// BytesToBuffer creates a NAPI Buffer from a Go byte slice. The data is copied.
func (e Env) BytesToBuffer(data []byte) (Value, error) {
	var result C.napi_value
	var buf unsafe.Pointer
	status := C.napi_create_buffer(e.raw, C.size_t(len(data)), &buf, &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	if len(data) > 0 {
		C.memcpy(buf, unsafe.Pointer(unsafe.SliceData(data)), C.size_t(len(data)))
	}
	return Value{raw: result}, nil
}

// BufferToBytes extracts the byte slice from a NAPI Buffer.
func (e Env) BufferToBytes(value Value) ([]byte, error) {
	var data unsafe.Pointer
	var length C.size_t
	status := C.napi_get_buffer_info(e.raw, value.raw, &data, &length)
	if err := napiStatusToError(status); err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, nil
	}
	return unsafe.Slice((*byte)(data), length), nil
}

// ThrowError throws a JavaScript error with the given message.
func (e Env) ThrowError(msg string) error {
	cMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsg))
	status := C.napi_throw_error(e.raw, nil, cMsg)
	return napiStatusToError(status)
}

// SetNamedProperty sets a property on an object.
func (e Env) SetNamedProperty(object Value, name string, value Value) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	status := C.napi_set_named_property(e.raw, object.raw, cName, value.raw)
	return napiStatusToError(status)
}

// CallbackFunction is the function signature for NAPI callbacks.
// It receives a single argument and returns a single value.
type CallbackFunction func(env Env, args []Value) Value

//export goNapiCallbackTrampoline
func goNapiCallbackTrampoline(env C.napi_env, info C.napi_callback_info) C.napi_value {
	var argc C.size_t = 16
	var argv [16]C.napi_value
	var handle C.uintptr_t
	var thisArg C.napi_value
	status := C.go_napi_get_cb_info(env, info, &argc, &argv[0], &thisArg, &handle)
	if err := napiStatusToError(status); err != nil {
		panic(err)
	}

	fn := cgo.Handle(handle).Value().(CallbackFunction)
	args := make([]Value, argc)
	for i := range args {
		args[i] = Value{raw: argv[i]}
	}
	return fn(Env{raw: env}, args).raw
}

// CreateFunction creates a NAPI function from a Go callback.
func (e Env) CreateFunction(name string, fn CallbackFunction) (Value, error) {
	handle := cgo.NewHandle(fn)

	var result C.napi_value
	status := C.go_napi_create_function(
		e.raw,
		(*C.char)(unsafe.Pointer(unsafe.StringData(name))),
		C.size_t(len(name)),
		(*[0]byte)(C.goNapiCallbackTrampoline),
		C.uintptr_t(handle),
		&result,
	)

	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// CallFunction calls a JS function with the given arguments.
// recv is the `this` value (use GetUndefinedValue for global calls).
func (e Env) CallFunction(recv Value, fn Value, args []Value) (Value, error) {
	var result C.napi_value
	argc := C.size_t(len(args))
	var argv *C.napi_value
	if len(args) > 0 {
		cArgs := make([]C.napi_value, len(args))
		for i, a := range args {
			cArgs[i] = a.raw
		}
		argv = &cArgs[0]
	}
	status := C.napi_call_function(e.raw, recv.raw, fn.raw, argc, argv, &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// Ref wraps a napi_ref handle to prevent garbage collection of a JS value.
type Ref struct {
	raw C.napi_ref
}

// CreateReference creates a strong reference to a JS value.
func (e Env) CreateReference(value Value) (Ref, error) {
	var ref C.napi_ref
	status := C.napi_create_reference(e.raw, value.raw, 1, &ref)
	if err := napiStatusToError(status); err != nil {
		return Ref{}, err
	}
	return Ref{raw: ref}, nil
}

// GetReferenceValue retrieves the JS value from a reference.
func (e Env) GetReferenceValue(ref Ref) (Value, error) {
	var result C.napi_value
	status := C.napi_get_reference_value(e.raw, ref.raw, &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}

// DeleteReference deletes a reference, allowing the value to be GC'd.
func (e Env) DeleteReference(ref Ref) error {
	status := C.napi_delete_reference(e.raw, ref.raw)
	return napiStatusToError(status)
}

// TypeOf returns the JavaScript type of a value as a string.
func (e Env) TypeOf(value Value) (string, error) {
	var result C.napi_valuetype
	status := C.napi_typeof(e.raw, value.raw, &result)
	if err := napiStatusToError(status); err != nil {
		return "", err
	}
	switch result {
	case C.napi_undefined:
		return "undefined", nil
	case C.napi_null:
		return "null", nil
	case C.napi_boolean:
		return "boolean", nil
	case C.napi_number:
		return "number", nil
	case C.napi_string:
		return "string", nil
	case C.napi_symbol:
		return "symbol", nil
	case C.napi_object:
		return "object", nil
	case C.napi_function:
		return "function", nil
	case C.napi_external:
		return "external", nil
	default:
		return "unknown", nil
	}
}

// GetNamedProperty gets a named property from a JS object.
func (e Env) GetNamedProperty(object Value, name string) (Value, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var result C.napi_value
	status := C.napi_get_named_property(e.raw, object.raw, cName, &result)
	if err := napiStatusToError(status); err != nil {
		return Value{}, err
	}
	return Value{raw: result}, nil
}
