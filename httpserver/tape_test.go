package main_test

import (
	"io/ioutil"
	"testing"

	httpserver "github.com/vkenrik117/httpserver"
)

func TestTape_Write(t *testing.T) {
	var err error

	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &httpserver.Tape{file}

	_, err = tape.Write([]byte("abc"))
	if err != nil {
		t.Errorf("Tape.Write returned an error: %v", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		t.Errorf("File.Seek returned an error: %v", err)
	}
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
