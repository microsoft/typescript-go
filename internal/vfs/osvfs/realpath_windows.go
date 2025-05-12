package osvfs

import (
	"errors"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

// This implementation is based on what Node's fs.realpath.native does, via libuv: https://github.com/libuv/libuv/blob/ec5a4b54f7da7eeb01679005c615fee9633cdb3b/src/win/fs.c#L2937

func realpath(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close() //nolint:errcheck

	// Works on directories too since https://go.dev/cl/405275.
	h := windows.Handle(f.Fd())

	// based on https://github.com/golang/go/blob/f4e3ec3dbe3b8e04a058d266adf8e048bab563f2/src/os/file_windows.go#L389

	const _VOLUME_NAME_DOS = 0

	buf := make([]uint16, 310) // https://github.com/microsoft/go-winio/blob/3c9576c9346a1892dee136329e7e15309e82fb4f/internal/stringbuffer/wstring.go#L13
	for {
		n, err := windows.GetFinalPathNameByHandle(h, &buf[0], uint32(len(buf)), _VOLUME_NAME_DOS)
		if err != nil {
			return "", err
		}
		if n < uint32(len(buf)) {
			break
		}
		buf = make([]uint16, n)
	}

	s := syscall.UTF16ToString(buf)
	if len(s) > 4 && s[:4] == `\\?\` {
		s = s[4:]
		if len(s) > 3 && s[:3] == `UNC` {
			// return path like \\server\share\...
			return `\` + s[3:], nil
		}
		return s, nil
	}

	return "", errors.New("GetFinalPathNameByHandle returned unexpected path: " + s)
}
