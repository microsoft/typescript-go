//go:build linux || darwin

package osvfs

import "syscall"

func ignoringEINTR(fn func() error) error {
	for {
		err := fn()
		if err != syscall.EINTR { //nolint:errorlint // syscall functions return raw syscall.Errno, never wrapped
			return err
		}
	}
}

func ignoringEINTR2[T any](fn func() (T, error)) (T, error) {
	for {
		v, err := fn()
		if err != syscall.EINTR { //nolint:errorlint // syscall functions return raw syscall.Errno, never wrapped
			return v, err
		}
	}
}
