// Package uuid provides functions that create a consistent globally unique
// UUID for a given TCP socket.  The package defines a new command-line flag
// `-uuid-prefix-file`, and that file and its contents should be set up prior
// to invoking any command which uses this library.
package uuid

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"syscall"
	"unsafe"
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
	// UUIDPrefixFile is a command-line flag to hold the filename which contains
	// the UUID prefix. Ideally it will be something like "/var/local/uuid/prefix",
	// and the contents of the named file will be a string like
	// "host.example.com_45353453".
	UUIDPrefixFile = flag.String("uuid-prefix-file", "/var/local/uuid/prefix",
		"The file holding the UUID prefix for sockets created in this network namespace.")

	// Only calculate these once - they never change. These are the contents of a
	// flag-specified file which holds this configuration information.
	cachedPrefixString   string
	cachedPrefixError    error
	cachedPrefixInitOnce sync.Once
)

// getPrefix returns a prefix string which contains the hostname and approximate
// boot time of the machine, which globally uniquely identifies the socket uuid
// namespace, assuming the the namespace was not created more than once in the
// span of two seconds.
//
// The return values of this function should be cached because that pair should
// be constant for a given instance of the program, unless the boot time changes
// (how?) or the hostname changes (why?) while this program is running.
func getPrefix(proxyFile string) (string, error) {
	contents, err := ioutil.ReadFile(proxyFile)
	return string(contents), err
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
	return FromCookie(cookie)
}

// FromCookie returns a string that is a globally unique identifier for the
// passed-in socket cookie.
//
// This function will never return the empty string, but the returned string
// value should only be used if the error is nil.
func FromCookie(cookie uint64) (string, error) {
	cachedPrefixInitOnce.Do(
		func() {
			// We can't do this setup in init() because the value of the flag needs to be parsed
			// from the command line. So we do it in this function, which should only be
			// called once.
			cachedPrefixString, cachedPrefixError = getPrefix(*UUIDPrefixFile)
		},
	)
	if cachedPrefixError != nil {
		return invalidUUID, cachedPrefixError
	}
	return fmt.Sprintf("%s_%016X", cachedPrefixString, cookie), nil
}
