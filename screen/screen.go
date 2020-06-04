package screen

import (
	"fmt"
	"io"
	"os"

	"github.com/ahmetb/go-cursor"
	"github.com/gliderlabs/ssh"
	"github.com/nsf/termbox-go"

	loud "github.com/Pylons-tech/LOUD/data"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

const allowMouseInputAndHideCursor string = "\x1b[?1003h\x1b[?25l"
const resetScreen string = "\x1bc"
const ellipsis = "…"

// const hpon = "◆"
// const hpoff = "◇"
const bgcolor = 232

// Screen represents a UI screen.
type Screen interface {
	SaveGame()
	IsEndGameConfirmScreen() bool
	UpdateFakeBlockHeight(int64)
	SetScreenSize(int, int)
	HandleInputKey(termbox.Event)
	GetScreenStatus() PageStatus
	SetScreenStatus(PageStatus)
	GetTxFailReason() string
	FakeSync()
	Resync()
	Render()
	Reset()
}

// GameScreen is a struct to manage screen of game
type GameScreen struct {
	world            loud.World
	user             loud.User
	screenSize       ssh.Window
	activeItem       loud.Item
	activeItSpec     loud.ItemSpec
	activeCharacter  loud.Character
	activeChSpec     loud.CharacterSpec
	activeLine       int
	activeTrdReq     loud.TrdReq
	activeItemTrdReq interface{}
	pylonEnterValue  string
	loudEnterValue   string
	actionText       string
	inputText        string
	syncingData      bool
	blockHeight      int64
	fakeBlockHeight  int64
	txFailReason     string
	txResult         []byte
	refreshed        bool
	scrStatus        PageStatus
	colorCodeCache   map[string](func(string) string)
}

// NewScreen manages the window rendering for game
func NewScreen(world loud.World, user loud.User) Screen {
	width, _ := terminal.Width()
	height, _ := terminal.Height()

	window := ssh.Window{
		Width:  int(width),
		Height: int(height),
	}

	screen := GameScreen{
		world:          world,
		user:           user,
		screenSize:     window,
		colorCodeCache: make(map[string](func(string) string))}

	return &screen
}

// Render method is for printing on terminal
func (screen *GameScreen) Render() {
	if len(loud.SomethingWentWrongMsg) > 0 {
		clear := cursor.ClearEntireScreen()
		dead := loud.Localize("Something went wrong, please close using Esc key and see loud.log")
		move := cursor.MoveTo(screen.Height()/2, screen.Width()/2-NumberOfSpaces(dead)/2)
		io.WriteString(os.Stdout, clear+move+dead)

		detailedErrorMsg := fmt.Sprintf("%s: %s", loud.Localize("detailed error"), loud.SomethingWentWrongMsg)
		move = cursor.MoveTo(screen.Height()/2+3, screen.Width()/2-NumberOfSpaces(dead)/2)
		io.WriteString(os.Stdout, move+detailedErrorMsg)
		screen.refreshed = false
		return
	}
	if screen.scrStatus == "" {
		screen.SetScreenStatus(ShowLocation)
	}

	if screen.Height() < 38 || screen.Width() < 120 {
		clear := cursor.ClearEntireScreen()
		move := cursor.MoveTo(1, 1)
		io.WriteString(os.Stdout,
			fmt.Sprintf("%s%s%s", clear, move, loud.Localize("screen size warning")))
		return
	}

	if !screen.refreshed {
		clear := cursor.ClearEntireScreen() + allowMouseInputAndHideCursor
		io.WriteString(os.Stdout, clear)
		screen.redrawBorders()
		screen.refreshed = true
	}

	screen.renderUserCommands()
	screen.renderUserSituation()
	screen.renderCharacterSheet()
	screen.renderInputValue()
	screen.renderMenu()
}
