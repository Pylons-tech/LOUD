package main

import (
	"os"

	loud "github.com/Pylons-tech/LOUD"
	"github.com/Pylons-tech/LOUD/log"
)

func main() {
	f, err := os.OpenFile("loud.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Warnf("error opening file: %v", err)
		log.Debugln("just going on without using log file ...")
	}
	defer f.Close()
	loud.ServeGame(f)
}
