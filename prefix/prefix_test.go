package prefix

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/m-lab/go/osx"
	"github.com/m-lab/go/rtx"
)

var podName string = "pod-x9lnt"

func TestMain(m *testing.M) {
	cleanupPodNameEnv := osx.MustSetenv("POD_NAME", podName)
	defer cleanupPodNameEnv()
	os.Exit(m.Run())
}

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

func TestGenerateWithPodNamePrefix(t *testing.T) {
	s, err := generate([]string{})
	if err != nil {
		t.Errorf("err should have been nil but got: %v", err)
	}
	if !strings.HasPrefix(s, podName) {
		t.Errorf("wanted prefix '%s', but got '%s'", podName, s)
	}
}

func TestGenerateWithoutPodNameEnvVar(t *testing.T) {
	osLookupEnv = func(e string) (string, bool) {
		return "", false
	}
	osHostname = func() (string, error) {
		return "GOODHOSTNAME", nil
	}
	defer func() {
		osLookupEnv = os.LookupEnv
		osHostname = os.Hostname
	}()

	s, err := generate([]string{})
	if err != nil {
		t.Errorf("err should have been nil but got: %v", err)
	}
	if !strings.HasPrefix(s, "GOODHOSTNAME") {
		t.Errorf("wanted prefix 'GOODHOSTNAME', but got '%s'", s)
	}
}

func TestGenerateWithBadHostName(t *testing.T) {
	osHostname = func() (string, error) {
		return "", errors.New("hostname error for testing")
	}
	osLookupEnv = func(e string) (string, bool) {
		return "", false
	}
	defer func() {
		osHostname = os.Hostname
		osLookupEnv = os.LookupEnv
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
		t.Error("expected an error")
	}
	if !strings.HasPrefix(s, "BADPREFIX") {
		t.Errorf("wanted value of 'BADPREFIX', but got '%s'", s)
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
	if err == nil || s == "" {
		t.Error("Should have had a non-nil error and non-empty returned string")
	}
}

func TestGenerateWithBadProcUptime(t *testing.T) {
	f, err := ioutil.TempFile("", "TestGenerateWithBadProcUptime")
	rtx.Must(err, "Could not create tempfile")
	defer os.Remove(f.Name())

	procUptime = f.Name()
	defer func() {
		procUptime = "/proc/uptime"
	}()

	rtx.Must(ioutil.WriteFile(f.Name(), []byte("123"), 0600), "Could not write to temp file")
	s, err := generate([]string{})
	if err == nil || s == "" {
		t.Error("Should have had a non-nil error and non-empty returned string")
	}

	rtx.Must(ioutil.WriteFile(f.Name(), []byte("this_should_not_parse 123"), 0600), "Could not write to temp file")
	s, err = generate([]string{})
	if err == nil || s == "" {
		t.Error("Should have had a non-nil error and non-empty returned string")
	}
}
