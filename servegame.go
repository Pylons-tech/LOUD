package loud

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsf/termbox-go"
)

func SetupScreenAndEvents(world World) {
	user := world.GetUser("eugen")
	screen := NewScreen(world, user)

	logMessage := fmt.Sprintf("setting up screen and events at %s", time.Now().UTC().Format(time.RFC3339))
	log.Println(logMessage)

	tick := time.Tick(50 * time.Millisecond)

	// Setup terminal close handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		screen.Reset()
		os.Exit(0)
	}()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	screen.Render()

eventloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				screen.Reset()
				break eventloop
			default:
				screen.HandleInputKey(ev)
			}
		case termbox.EventResize:
			logMessage := fmt.Sprintf("Handling TermBox Resize Event (%d, %d) at %s", ev.Width, ev.Height, time.Now().UTC().Format(time.RFC3339))
			log.Println(logMessage)

			screen.SetScreenSize(ev.Width, ev.Height)
		case termbox.EventError:
			panic(ev.Err)
		}
		select {
		case <-tick:
			screen.Render()
			continue
		}
	}
}

// ServeGame runs the main game loop.
func ServeGame() {
	rand.Seed(time.Now().Unix())

	world := LoadWorldFromDB("./world.db")
	defer world.Close()

	SetupScreenAndEvents(world)
}
