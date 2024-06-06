package uuid

import (
    "os"
    "testing"
)

func TestSetUUIDPrefixFile(t *testing.T) {
    expected := []byte("host.example.com_1552945174")
    f, err :=  os.CreateTemp("", "")
    if err != nil {
        t.Error(err)
    }
    defer os.Remove(f.Name())
    if _, err = f.Write(expected); err != nil {
        t.Error(err)
    }
    if err = f.Close(); err != nil {
        t.Error(err)
    }

    err = SetUUIDPrefixFile(f.Name())
    if err != nil {
        t.Error(err)
    }
    if string(uuidPrefix) != string(expected) {
        t.Errorf("Expected '%s', got '%s'", string(expected), string(uuidPrefix))
    }
}

func TestSetUUIDPrefixFileError(t *testing.T) {
    err := SetUUIDPrefixFile("INVALID_FILENAME")
    if err == nil {
        t.Error("Expected error for invalid file, but got none.")
    }
}
