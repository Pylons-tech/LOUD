package loud

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
	"github.com/nsf/termbox-go"
)

func SetupLoggingFile() {
	f, err := os.OpenFile("loud.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.Println("Starting to save log into file")
	log.SetOutput(f)

	log.Println("Starting")
}

func SetupScreenAndEvents(world World) {
	args := os.Args
	username := ""
	log.Println("args SetupScreenAndEvents", args)
	if len(args) < 2 {
		log.Println("you didn't configure username when running, using default username \"eugen\"")
		username = "eugen"
	} else {
		username = args[1]
	}
	user := world.GetUser(username)

	SetupLoggingFile()

	screen := NewScreen(world, user)

	logMessage := fmt.Sprintf("setting up screen and events at %s", time.Now().UTC().Format(time.RFC3339))
	log.Println(logMessage)

	tick := time.Tick(50 * time.Millisecond)
	daemonStatusTick := time.Tick(10 * time.Second)

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
				screen.SaveGame()
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
		case <-daemonStatusTick:
			ds, err := pylonSDK.GetDaemonStatus()
			if err != nil {
				log.Println("couldn't get daemon status", err)
			} else {
				log.Println("success getting daemon status", err)
				screen.UpdateBlockHeight(ds.SyncInfo.LatestBlockHeight)
			}
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
