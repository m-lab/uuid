package socookie

import (
	"flag"
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/m-lab/uuid/prefix"

	"github.com/m-lab/go/flagx"
)

const (
	// defined in socket.h in the linux kernel
	syscallSoCookie = 57 // syscall.SO_COOKIE does not exist in golang 1.11
)

// get implements Get for Linux systems.
func get(file *os.File) (uint64, error) {
	var cookie uint64
	cookieLen := uint32(unsafe.Sizeof(cookie))
	// GetsockoptInt does not work for 64 bit integers, which is what the UUID is.
	// So we crib from the GetsockoptInt implementation and ndt-server/tcpinfox,
	// and call the syscall manually.
	_, _, errno := syscall.Syscall6(
		uintptr(syscall.SYS_GETSOCKOPT),
		uintptr(int(file.Fd())),
		uintptr(syscall.SOL_SOCKET),
		uintptr(syscallSoCookie),
		uintptr(unsafe.Pointer(&cookie)),
		uintptr(unsafe.Pointer(&cookieLen)),
		uintptr(0))

	if errno != 0 {
		return 0, fmt.Errorf("Error in Getsockopt. Errno=%d", errno)
	}
	return cookie, nil
}
