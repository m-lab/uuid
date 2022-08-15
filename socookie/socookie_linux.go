package socookie

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	// defined in socket.h in the linux kernel
	syscallSoCookie = 57 // syscall.SO_COOKIE does not exist in golang 1.11
)

// get implements Get for Linux systems.
func get(file *os.File) (uint64, error) {
	var cookie uint64
	cookieLen := uint32(unsafe.Sizeof(cookie))
	rawConn, err := file.SyscallConn()
	if err != nil {
		return 0, err
	}
	var errno syscall.Errno
	err = rawConn.Control(func(fd uintptr) {
		// GetsockoptInt does not work for 64 bit integers, which is what the UUID is.
		// So we crib from the GetsockoptInt implementation and ndt-server/tcpinfox,
		// and call the syscall manually.
		_, _, errno = syscall.Syscall6(
			uintptr(syscall.SYS_GETSOCKOPT),
			fd,
			uintptr(syscall.SOL_SOCKET),
			uintptr(syscallSoCookie),
			uintptr(unsafe.Pointer(&cookie)),
			uintptr(unsafe.Pointer(&cookieLen)),
			uintptr(0))
	})
	if err != nil {
		return 0, err
	}
	if errno != 0 {
		return 0, fmt.Errorf("error in Getsockopt. Errno=%d", errno)
	}
	return cookie, nil
}
