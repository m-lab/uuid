package uuid

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/m-lab/go/rtx"
)

// The tests here are whitebox tests because we want to specify the uuidPrefix
// directly instead of using the command-line flag.

func TestUUID(t *testing.T) {
	uuidPrefix = []byte("host.example.com_1552945174")

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
	uuidPrefix = []byte("host.example.com_1552945174")

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
	uuidPrefix = []byte("host.example.com_1552945174")

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
