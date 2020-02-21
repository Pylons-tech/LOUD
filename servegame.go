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
	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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
		log.Println("you didn't configure username when running, using default username \"eugen\"")
		username = "eugen"
	} else {
		username = args[1]
	}
	user := world.GetUser(username)

	SetupLoggingFile(logFile)

	screenInstance := screen.NewScreen(world, user)

	log.Println("setting up screen and events")

	tick := time.Tick(50 * time.Millisecond)
	daemonStatusRefreshTick := time.Tick(10 * time.Second)
	daemonFetchResult := make(chan *ctypes.ResultStatus)

	if data.AutomateInput {
		screenInstance.SetScreenStatus(screen.RESULT_SWITCH_USER)
		time.AfterFunc(2*time.Second, func() {

		automateloop:
			for {
				log.Println("<-automateTick")
				switch screenInstance.GetScreenStatus() {
				case screen.RESULT_CREATE_COOKBOOK:
					if screenInstance.GetTxFailReason() != "" {
						data.SomethingWentWrongMsg = "create cookbook failed, " + screenInstance.GetTxFailReason()
						break automateloop
					}
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 122, // "z" 122 Switch user
					})
				case screen.RESULT_GET_PYLONS:
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 106, // "j" 106 Create cookbook
					})
				case screen.RESULT_SWITCH_USER:
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

eventloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				screenInstance.SaveGame()
				screenInstance.Reset()
				break eventloop
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
		select {
		case <-tick:
			screenInstance.Render()
			continue
		case <-daemonStatusRefreshTick:
			go func() {
				screenInstance.SetDaemonFetchingFlag(true)
				screenInstance.Render()
				ds, err := pylonSDK.GetDaemonStatus()
				if err != nil {
					log.Println("couldn't get daemon status", err)
				} else {
					log.Println("success getting daemon status", err)
					daemonFetchResult <- ds
				}
				screenInstance.Resync()
			}()
		case ds := <-daemonFetchResult:
			screenInstance.SetDaemonFetchingFlag(false)
			screenInstance.UpdateBlockHeight(ds.SyncInfo.LatestBlockHeight)
		case <-terminalCloseSignal:
			screenInstance.Reset()
			break eventloop
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
