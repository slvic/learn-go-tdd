package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	httpserver "github.com/vkenrik117/httpserver"
)

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp File %v", err)
	}

	_, err = tmpfile.Write([]byte(initialData))
	if err != nil {
		t.Errorf("tmpfile.Write returned an error: %v", err)
	}

	removeFile := func() {
		err = tmpfile.Close()
		if err != nil {
			t.Errorf("tmpfile.Close returned an error: %v", err)
		}
		err = os.Remove(tmpfile.Name())
		if err != nil {
			t.Errorf("os.Remove returned an error: %v", err)
		}
	}

	return tmpfile, removeFile
}

func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := httpserver.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()

		want := []httpserver.Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		got = store.GetLeague()

		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := httpserver.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetPlayerScore("Chris")
		want := 33
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := httpserver.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := httpserver.NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		assertScoreEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := httpserver.NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := httpserver.NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := []httpserver.Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}

func assertScoreEquals(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
