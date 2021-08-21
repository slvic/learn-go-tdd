package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type FileSystemPlayerStore struct {
	Database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialisePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		Database: json.NewEncoder(&Tape{file}),
		league:   league,
	}, nil
}

func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		err = db.Close()
	}

	if err != nil {
		return nil, nil, fmt.Errorf("could not close the connection: %v", err)
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, closeFunc, nil
}

func initialisePlayerDBFile(file *os.File) (err error) {
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("could not Seek the file: %v", err)
	}

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		_, err = file.Write([]byte("[]"))
		if err != nil {
			return fmt.Errorf("could not Write into the file when the file is empty: %v", err)
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("could not Seek the file when the file is empty: %v", err)
		}

	}

	return nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	err := f.Database.Encode(f.league)
	if err != nil {
		fmt.Printf("could not Encode league into JSON: %v", err)
	}
}