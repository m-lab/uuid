package prefix

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/m-lab/go/rtx"
)

func TestUnsafeString(t *testing.T) {
	s := UnsafeString()
	if !strings.Contains(s, "unsafe") {
		t.Error(s, "should contain the substring \"unsafe\"")
	}
}

func TestGenerate(t *testing.T) {
	f, err := ioutil.TempFile("", "TestGenerate")
	rtx.Must(err, "Could not create tempfile")
	defer os.Remove(f.Name())

	rtx.Must(Generate(f.Name()), "Could not generate prefix")

	contents, err := ioutil.ReadFile(f.Name())
	rtx.Must(err, "Could not read the tempfile")
	if !strings.Contains(string(contents), "_") {
		t.Error(string(contents), "should contain the substring \"_\"")
	}
}

func TestGenerateWithBadHostname(t *testing.T) {
	osHostname = func() (string, error) {
		return "", errors.New("hostname error for testing")
	}
	defer func() {
		osHostname = os.Hostname
	}()

	f, err := ioutil.TempFile("", "TestGenerateWithBadHostname")
	rtx.Must(err, "Could not create tempfile")
	defer os.Remove(f.Name())

	err = Generate(f.Name())
	if err == nil {
		t.Errorf("Should have had an error here")
	}

	s, err := generate([]string{})
	if err == nil {
		t.Error("Should have had an error")
	}
	if s == "" {
		t.Error("Errors should not have an empty return string")
	}
}

func TestGenerateWithBadDestinationDir(t *testing.T) {
	err := Generate("/this/directory/does/not/exist/nor_does_this_file")
	if err == nil {
		t.Errorf("Should have had an error with a bad filename")
	}
}

func TestGenerateWithNonexistentProcUptime(t *testing.T) {
	procUptime = "/this/file/does/not/exist"
	defer func() {
		procUptime = "/proc/uptime"
	}()
	s, err := generate([]string{})
	if err == nil {
		t.Error("Should have had an error")
	}
	if s == "" {
		t.Error("Errors should not have an empty return string")
	}
}
