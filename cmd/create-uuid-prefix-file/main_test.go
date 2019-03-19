package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/m-lab/go/osx"
	"github.com/m-lab/go/rtx"
)

func TestMain(t *testing.T) {
	f, err := ioutil.TempFile("", "TestMain")
	rtx.Must(err, "Could not create tempfile")
	defer os.Remove(f.Name())

	cleanupEnv := osx.MustSetenv("FILENAME", f.Name())
	defer cleanupEnv()

	main()

	b, err := ioutil.ReadFile(f.Name())
	rtx.Must(err, "Could not read tempfile after creation.")
	s := string(b)
	if !strings.Contains(s, "_") || strings.Contains(s, "_unsafe_") || strings.ContainsAny(s, "\n\t ") {
		t.Errorf("Bad UUID %q in file %s", s, f.Name())
	}
}
