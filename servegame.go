package loud

import (
	// "regexp"
	"bufio"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nsf/termbox-go"

	cf "github.com/Pylons-tech/LOUD/config"
	data "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/LOUD/log"
	screen "github.com/Pylons-tech/LOUD/screen"
)

var terminalCloseSignal chan os.Signal = make(chan os.Signal, 2)

// SetupLoggingFile set custom file for logging output
func SetupLoggingFile(f *os.File) {
	log.WithFields(log.Fields{
		"func_start": "SetupLoggingFile",
	}).Debugln("function start")

	log.SetOutput(f)
	log.Debugln("You should see this log as a first line of log file")
}

// SetupScreenAndEvents setup screen object and events
func SetupScreenAndEvents(world data.World, logFile *os.File) {
	log.WithFields(log.Fields{
		"func_start": "SetupScreenAndEvents",
		"os.args":    os.Args,
	}).Debugln("function start")

	args := os.Args
	username := ""
	if len(args) < 2 {
		log.Infoln("you didn't configure username when running!")
		log.Infoln("Please enter your username!")
		// for {
		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
		username = strings.TrimSuffix(username, "\n")
		// For now, not put validation for username
		// break
		// var validUsername = regexp.MustCompile(`^[a-z][a-z0-9/]{2,63}$`)
		// isValid := validUsername.MatchString(username)
		// if isValid {
		// 	break
		// } else {
		// 	log.Infoln("username should consist of only a-z and 0-9. And first letter should be a-z.")
		// }
		// }
	} else {
		username = args[1]
	}

	log.WithFields(log.Fields{
		"set_username": "SetupScreenAndEvents",
		"username":     username,
		"username_len": len(username),
	}).Infoln("set the username")

	user := world.GetUser(username)

	SetupLoggingFile(logFile)

	screenInstance := screen.NewScreen(world, user)
	log.Debugln("setting up screen and events")

	tick := time.NewTicker(300 * time.Millisecond)
	config, cferr := cf.ReadConfig()
	if cferr != nil {
		log.WithFields(log.Fields{
			"log": cferr,
		}).Fatal("Couldn't load configuration file")
	}
	regRefreshTick := time.NewTicker(time.Duration(config.App.DaemonTimeoutCommit) * time.Second)

	if data.AutomateInput {
		screenInstance.SetScreenStatus(screen.RsltSwitchUser)
		time.AfterFunc(2*time.Second, func() {

		automateloop:
			for {
				log.Debugln("automation tick")
				switch screenInstance.GetScreenStatus() {
				case screen.RsltCreateCookbook:
					if screenInstance.GetTxFailReason() != "" {
						data.SomethingWentWrongMsg = "create cookbook failed, " + screenInstance.GetTxFailReason()
						break automateloop
					}
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 122, // "z" 122 Switch user
					})
				case screen.RsltGetPylons:
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 106, // "j" 106 Create cookbook
					})
				case screen.RsltSwitchUser:
					screenInstance.HandleInputKey(termbox.Event{
						Ch: 121, // "y" 121 get initial pylons
					})
					data.AutomateRunCnt++
					log.WithFields(log.Fields{
						"index": data.AutomateRunCnt,
					}).Infof("Running automation task")
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
			case <-regRefreshTick.C:
				screenInstance.FakeSync()
			case <-terminalCloseSignal:
				screenInstance.Reset()
				os.Exit(0)
			case <-tick.C:
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
			log.WithFields(log.Fields{
				"width":  ev.Width,
				"height": ev.Height,
			}).Infoln("Handling TermBox Resize Event")

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
