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

	data "github.com/Pylons-tech/LOUD/data"
	screen "github.com/Pylons-tech/LOUD/screen"
)

var terminalCloseSignal chan os.Signal = make(chan os.Signal, 2)

func SetupLoggingFile(f *os.File) {
	log.Println("Starting to save log into file")
	log.SetOutput(f)
	log.Println("Starting")
}

func SetupScreenAndEvents(world data.World, logFile *os.File) {
	args := os.Args
	username := ""
	log.Println("args SetupScreenAndEvents", args)
	if len(args) < 2 {
		log.Fatal("you didn't configure username when running!")
	} else {
		username = args[1]
	}
	user := world.GetUser(username)

	SetupLoggingFile(logFile)

	screenInstance := screen.NewScreen(world, user)

	log.Println("setting up screen and events")

	tick := time.Tick(300 * time.Millisecond)
	regRefreshTick := time.Tick(5 * time.Second)

	if data.AutomateInput {
		screenInstance.SetScreenStatus(screen.RSLT_SWITCH_USER)
		time.AfterFunc(2*time.Second, func() {

		automateloop:
			for {
				log.Println("<-automateTick")
				switch screenInstance.GetScreenStatus() {
				case screen.RSLT_CREATE_COOKBOOK:
					if screenInstance.GetTxFailReason() != "" {
						data.SomethingWentWrongMsg = "create cookbook failed, " + screenInstance.GetTxFailReason()
						break automateloop
					}
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 122, // "z" 122 Switch user
					})
				case screen.RSLT_GET_PYLONS:
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 106, // "j" 106 Create cookbook
					})
				case screen.RSLT_SWITCH_USER:
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 121, // "y" 121 get initial pylons
					})
					data.AutomateRunCnt += 1
					log.Printf("Running %dth automation task", data.AutomateRunCnt)
				}
				time.Sleep(2 * time.Second)
			}
		})
	}

	// Setup terminal close handler
	signal.Notify(terminalCloseSignal, os.Interrupt, syscall.SIGTERM)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	screenInstance.Render()

	go func() {
		for {
			select {
			case <-regRefreshTick:
				screenInstance.FakeSync()
			case <-terminalCloseSignal:
				screenInstance.Reset()
				os.Exit(0)
			case <-tick:
				screenInstance.Render()
			}
		}
	}()
eventloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				if screenInstance.IsEndGameConfirmScreen() {
					screenInstance.SaveGame()
					screenInstance.Reset()
					break eventloop
				} else {
					screenInstance.HandleInputKey(ev)
				}
			default:
				screenInstance.HandleInputKey(ev)
			}
		case termbox.EventResize:
			logMessage := fmt.Sprintf("Handling TermBox Resize Event (%d, %d) at %s", ev.Width, ev.Height, time.Now().UTC().Format(time.RFC3339))
			log.Println(logMessage)

			screenInstance.SetScreenSize(ev.Width, ev.Height)
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

// ServeGame runs the main game loop.
func ServeGame(logFile *os.File) {
	rand.Seed(time.Now().Unix())

	world := data.LoadWorldFromDB("./world.db")
	defer world.Close()

	SetupScreenAndEvents(world, logFile)
}
