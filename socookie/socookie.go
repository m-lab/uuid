// Package socookie allows to get a socket's cookie.
package socookie

import (
	"os"
)

// Get returns the cookie (the UUID) associated with a socket. For a given
// boot of a given hostname, this UUID is guaranteed to be unique (until the
// host receives more than 2^64 connections without rebooting).
//
// The implementation for Linux uses kernel support for getting the
// cookie. The implementation for other platforms does not and therefore
// cannot be relied upon; use it at your own risk.
func Get(file *os.File) (uint64, error) {
	return get(file)
}
