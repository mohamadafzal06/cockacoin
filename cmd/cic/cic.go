package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/mohamadafzal06/cokacoin"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	path := os.Getenv("PATH_TO_DIR")
	store, _ := cokacoin.NewFolderStore(path)

	bc, err := cokacoin.NewBlockchain(2, store)
	if err != nil {
		log.Fatal(err.Error())
	}

	bc.Add("Hello")
	bc.Add("Another")

	if err := bc.Validate(); err != nil {
		log.Fatalf(err.Error())
	}

	bc.Print()

}
