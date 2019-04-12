// +build !linux

package socookie

import (
	"os"
	"sync/atomic"
)

// cookieGen is the counter we use to emulate SO_COOKIE.
var cookieGen uint64

// get implements Get for non Linux systems.
func get(file *os.File) (uint64, error) {
	return atomic.AddUint64(&cookieGen, uint64(1)), nil
}
