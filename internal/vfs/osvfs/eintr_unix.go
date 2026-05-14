//go:build !windows

package osvfs

import "syscall"

func ignoringEINTR(fn func() error) error {
	for {
		err := fn()
		if err != syscall.EINTR {
			return err
		}
	}
}

func ignoringEINTR2[T any](fn func() (T, error)) (T, error) {
	for {
		v, err := fn()
		if err != syscall.EINTR {
			return v, err
		}
	}
}
