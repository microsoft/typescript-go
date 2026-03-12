package napi

// #cgo CFLAGS: -I${SRCDIR}/packaged/include
// #cgo linux LDFLAGS: -Wl,-undefined,dynamic_lookup
// #cgo darwin LDFLAGS: -Wl,-undefined,dynamic_lookup
// #cgo windows,amd64 LDFLAGS: -L${SRCDIR}/packaged/lib/windows_amd64 -lnode
// #include "node/node_api.h"
// napi_value goNapiInit(napi_env env, napi_value exports);
// NAPI_MODULE(tsgo_napi, goNapiInit)
import "C"
