package main

import (
	"fmt"
	poker "github.com/vkenrik117/httpserver"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, closeStore, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer closeStore()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	poker.NewCLI(store, os.Stdin).PlayPoker()
}
