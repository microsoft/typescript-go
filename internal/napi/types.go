package napi

// #include "node/node_api.h"
import "C"

import "fmt"

// napiStatus represents a NAPI status code.
type napiStatus C.napi_status

const (
	napiStatusOk napiStatus = C.napi_ok
)

func napiStatusToError(status C.napi_status) error {
	if napiStatus(status) == napiStatusOk {
		return nil
	}
	return fmt.Errorf("napi error: status %d", int(status))
}
