package prefix

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// Variables to aid in mocking for whitebox testing.
	osHostname  = os.Hostname
	osLookupEnv = os.LookupEnv
	procUptime  = "/proc/uptime"
)

// Generate creates a prefix and writes it to the specified file. This file
// should be stored in a well-known location, and this generation process should
// only occur once.
func Generate(filename string) error {
	contents, err := generate([]string{})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, []byte(contents), 0444)
	if err != nil {
		return err
	}
	return nil
}

// UnsafeString returns a prefix for the local system with the annotation
// "unsafe" indicating that the prefix is not guaranteed to be consistent.
//
// This function is intended to be used to provide default prefixes that are
// better than the empty string, but not all the way to "good".
func UnsafeString() string {
	s, _ := generate([]string{"unsafe"}) // ignore any errors. This is part of the "unsafe".
	return s
}

// generate creates the UUID prefix. It never returns the empty string, because
// poorly-coded library users will likely cause terrible problems if they use
// the returned string without checking the error condition and the returned
// string is the empty string.
func generate(extras []string) (string, error) {
	var err error
	// If it exists, prefer using the value of the environment variable POD_NAME
	// as the prefix for the UUID, else fall back on using the value of
	// os.Hostname().
	prefixValue, ok := osLookupEnv("POD_NAME")
	if !ok || prefixValue == "" {
		prefixValue, err = osHostname()
		if err != nil {
			return "BADHOSTNAME", err
		}
	}
	now := time.Now()
	uptimeBytes, err := ioutil.ReadFile(procUptime)
	if err != nil {
		return prefixValue + "_BADBOOTTIME", err
	}
	uptimePieces := strings.Split(string(uptimeBytes), " ")
	if len(uptimePieces) < 2 {
		return prefixValue + "_BADBOOTTIME", errors.New("could not tokenize /proc/uptime contents")
	}
	uptimeFloat, err := strconv.ParseFloat(uptimePieces[0], 64)
	if err != nil {
		return prefixValue + "_BADBOOTTIME", errors.New("could not parse /proc/uptime contents")
	}
	boottime := now.Add(-1 * time.Duration(uptimeFloat*1000000) * time.Microsecond).Unix()
	pieces := []string{
		prefixValue,
		fmt.Sprintf("%d", boottime),
	}
	pieces = append(pieces, extras...)
	return strings.Join(pieces, "_"), nil
}
