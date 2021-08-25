package poker_test

import (
	"io/ioutil"
	"testing"

	poker "github.com/vkenrik117/poker"
)

func TestTape_Write(t *testing.T) {
	var err error

	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &poker.Tape{File: file}

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
