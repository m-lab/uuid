package uuid

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/m-lab/go/rtx"
)

// The tests here are whitebox tests because the library performs caching, and
// we want to reset the cache for every test.

func TestFileAsProxyForBoottime(t *testing.T) {
	_, err := getPrefix("/this/file/does/not/exist")
	if err == nil {
		t.Error("Should have had an error on a non-existent file")
	}
}

func setup(prefix string) func() {
	f, err := ioutil.TempFile("", "test_prefix_file")
	rtx.Must(err, "Could not create tempfile")
	_, err = f.Write([]byte(prefix))
	rtx.Must(err, "Could not write prefix file")
	rtx.Must(f.Close(), "Could not close tempfile")
	*UUIDPrefixFile = f.Name()
	cachedPrefixInitOnce = sync.Once{}

	// Throw away this value. We just want to make sure all the caching has been
	// set up and the sync.Once struct has already been called. Also, this should
	// never return an error - if we see an error, then we screwed up the setup
	// somehow and the validity of any subsequent test is in question.
	_, err = FromCookie(0)
	rtx.Must(err, "The FromCookie call in setUp should never fail")

	return func() {
		os.Remove(f.Name())
	}
}

func TestErrorDoesntCauseNullUUID(t *testing.T) {
	cleanup := setup("host.example.com_1552945174")
	defer cleanup()

	cachedPrefixString = ""
	cachedPrefixError = errors.New("An error for testing")

	id, err := FromCookie(0)
	if err == nil {
		t.Error("Should have had an error")
	}
	if err != cachedPrefixError {
		t.Error("Error should have been", cachedPrefixError, "not", err)
	}
	if id == "" {
		t.Error("An error should not cause an empty-string uuid to be returned")
	}
}

func TestUUID(t *testing.T) {
	cleanup := setup("host.example.com_1552945174")
	defer cleanup()
	// We use the less-used TCP versions of Listen and Accept because we want to be
	// sure that we are getting a real TCP connection.
	localAddr, err := net.ResolveTCPAddr("tcp", "localhost:12345")
	rtx.Must(err, "No localhost")
	listener, err := net.ListenTCP("tcp", localAddr)
	rtx.Must(err, "Could not make TCP listener")
	local1, err := net.Dial("tcp", ":12345")
	defer local1.Close()
	local2, err := net.Dial("tcp", ":12345")
	defer local2.Close()
	rtx.Must(err, "Could not connect to myself")
	conn1, err := listener.AcceptTCP()
	rtx.Must(err, "Could not accept conn1")
	conn2, err := listener.AcceptTCP()
	rtx.Must(err, "Could not acceptc conn2")
	uuid1, err := FromTCPConn(conn1)
	rtx.Must(err, "Could not get uuid for conn1")
	uuid2, err := FromTCPConn(conn2)
	rtx.Must(err, "Could not get uuid for conn2")
	if uuid1 == uuid2 {
		t.Error("UUIDs must not be the same")
	}
	fmt.Println("UUIDs:", uuid1, uuid2)
	rtx.Must(err, "Could not get uuid")
	left1 := strings.LastIndex(uuid1, "_")
	left2 := strings.LastIndex(uuid2, "_")
	if left1 <= 0 || left2 <= 0 || uuid1[0:left1] != uuid2[0:left2] {
		t.Error("The left part of the UUIDs was not constant:", uuid1, uuid2)
	}
}

func TestFromFileError(t *testing.T) {
	cleanup := setup("host.example.com_1552945174")
	defer cleanup()
	f, err := ioutil.TempFile("", "TestFileError")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	id, err := FromFile(f)
	if err == nil {
		t.Error("Should have had an error")
	}
	if id == "" {
		t.Error("An error should not cause an empty-string uuid to be returned")
	}
}

func TestFromTCPConnError(t *testing.T) {
	cleanup := setup("host.example.com_1552945174")
	defer cleanup()
	localAddr, err := net.ResolveTCPAddr("tcp", "localhost:12346")
	rtx.Must(err, "No localhost")
	listener, err := net.ListenTCP("tcp", localAddr)
	rtx.Must(err, "Could not make TCP listener")
	local, err := net.Dial("tcp", ":12346")
	defer local.Close()
	rtx.Must(err, "Could not connect to myself")
	conn, err := listener.AcceptTCP()
	rtx.Must(err, "Could not accept conn1")
	conn.Close()
	local.Close()
	id, err := FromTCPConn(conn)
	if err == nil {
		t.Error("Should have had an error")
	}
	if id == "" {
		t.Error("Should not return the empty string on error")
	}
}
