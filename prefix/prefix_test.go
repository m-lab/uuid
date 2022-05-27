package prefix

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/m-lab/go/rtx"
)

func TestMain(m *testing.M) {
	os.Setenv("POD_NAME", "pod-x9lnt")
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

func TestGenerateWithBadPodName(t *testing.T) {
	osLookupEnv = func(e string) (string, bool) {
		return "", false
	}
	defer func() {
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
	if err == nil || s == "" {
		t.Error("Should have had a non-nil error and non-empty returned string")
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
