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

	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func SetupLoggingFile(f *os.File) {
	log.Println("Starting to save log into file")
	log.SetOutput(f)
	log.Println("Starting")
}

func SetupScreenAndEvents(world World, logFile *os.File) {
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

	SetupLoggingFile(logFile)

	screen := NewScreen(world, user)

	logMessage := fmt.Sprintf("setting up screen and events at %s", time.Now().UTC().Format(time.RFC3339))
	log.Println(logMessage)

	tick := time.Tick(50 * time.Millisecond)
	daemonStatusRefreshTick := time.Tick(10 * time.Second)
	daemonFetchResult := make(chan *ctypes.ResultStatus)

	// Setup terminal close handler
	terminalCloseSignal := make(chan os.Signal, 2)
	signal.Notify(terminalCloseSignal, os.Interrupt, syscall.SIGTERM)

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
		case <-daemonStatusRefreshTick:
			go func() {
				ds, err := pylonSDK.GetDaemonStatus()
				if err != nil {
					log.Println("couldn't get daemon status", err)
				} else {
					log.Println("success getting daemon status", err)
					daemonFetchResult <- ds
				}
			}()
		case ds := <-daemonFetchResult:
			screen.UpdateBlockHeight(ds.SyncInfo.LatestBlockHeight)
		case <-terminalCloseSignal:
			screen.Reset()
			os.Exit(0)
		}
	}
}

// ServeGame runs the main game loop.
func ServeGame(logFile *os.File) {
	rand.Seed(time.Now().Unix())

	world := LoadWorldFromDB("./world.db")
	defer world.Close()

	SetupScreenAndEvents(world, logFile)
}
