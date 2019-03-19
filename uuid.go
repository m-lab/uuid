// Package uuid provides functions that create a consistent globally unique
// UUID for a given TCP socket.  The package defines a new command-line flag
// `-uuid-prefix-file`, and that file and its contents should be set up prior
// to invoking any command which uses this library.
package uuid

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

	// Whenever there is an error We return this value instead of the empty
	// string. We do this in an effort to detect when client code
	// accidentally uses the returned UUID even when it should not have.
	//
	// This is borne out of past experience, most notably an incident where
	// returning an empty string and an error condition caused the
	// resulting code to create a file named ".gz", which was (thanks to
	// shell-filename-hiding rules) a very hard bug to uncover.  If a file
	// is ever named "INVALID_UUID.gz", it will be much easier to detect
	// that there is a problem versus just ".gz"
	invalidUUID = "INVALID_UUID"
)

var (
	// uuidPrefix is the UUID prefix derived from a command-line flag specifying
	// the file which contains the UUID prefix. Ideally the file will be something
	// like "/var/local/uuid/prefix", and the contents of the file will be bytes
	// that resemble "host.example.com_45353453".
	//
	// By default it is a string generated in an unsafe manner that contains the
	// substring "_unsafe_" in it. This is not great, but it is better than the
	// default of the empty string, and it allows people to use code which uses
	// this library without having to set up a well-known location.
	uuidPrefix = flagx.FileBytes(prefix.UnsafeString())
)

func init() {
	flag.Var(&uuidPrefix, "uuid-prefix-file", "The file holding the UUID prefix for sockets created in this network namespace.")
}

// getCookie returns the cookie (the UUID) associated with a socket. For a given
// boot of a given hostname, this UUID is guaranteed to be unique (until the
// host receives more than 2^64 connections without rebooting).
func getCookie(file *os.File) (uint64, error) {
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

// FromTCPConn returns a string that is a globally unique identifier for the
// socket held by the passed-in TCPConn (assuming hostnames are unique).
//
// This function will never return the empty string, but the returned string
// value should only be used if the error is nil.
func FromTCPConn(t *net.TCPConn) (string, error) {
	file, err := t.File()
	if err != nil {
		return invalidUUID, err
	}
	defer file.Close()
	return FromFile(file)
}

// FromFile returns a string that is a globally unique identifier for the socket
// represented by the os.File pointer.
//
// This function will never return the empty string, but the returned string
// value should only be used if the error is nil.
func FromFile(file *os.File) (string, error) {
	cookie, err := getCookie(file)
	if err != nil {
		return invalidUUID, err
	}
	return FromCookie(cookie), nil
}

// FromCookie returns a string that is a globally unique identifier for the
// passed-in socket cookie.
func FromCookie(cookie uint64) string {
	return fmt.Sprintf("%s_%016X", string(uuidPrefix), cookie)
}
