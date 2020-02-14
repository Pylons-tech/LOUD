package main

// import (
// 	"log"
// 	"time"
// )

// func main() {
// 	var tick <-chan time.Time
// 	tick = time.Tick(3 * time.Second)

// 	for {
// 		select {
// 		case <-tick:
// 			log.Println("<-tick")
// 		}
// 	}
// }

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
	loud.ServeGame(f)
}
