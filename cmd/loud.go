package main

import (
	"log"
	"os"

	loud "github.com/Pylons-tech/LOUD"
)

func main() {
	f, err := os.OpenFile("loud.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println("Starting")

	loud.ServeGame()

}
