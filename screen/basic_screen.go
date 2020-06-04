package screen

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ahmetb/go-cursor"
	"github.com/gliderlabs/ssh"
	"github.com/mgutz/ansi"
	"github.com/nsf/termbox-go"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/LOUD/log"
)

// Resync execute new sync from node
func (screen *GameScreen) Resync() {
	screen.syncingData = true
	screen.Render()
	go func() {
		log.Println("start syncing from node")
		loud.SyncFromNode(screen.user)
		log.Println("end syncing from node")
		screen.syncingData = false
		screen.Render()
	}()
}

// FakeSync update the expected block height, it's because syncing node in real time is not effective
func (screen *GameScreen) FakeSync() {
	screen.UpdateFakeBlockHeight(screen.fakeBlockHeight + 1)
	screen.Render()
}

// GetTxFailReason returns last transaction's failure reason
func (screen *GameScreen) GetTxFailReason() string {
	return screen.txFailReason
}

// SwitchUser change user to new user
func (screen *GameScreen) SwitchUser(newUser loud.User) {
	screen.user = newUser
}

// func (screen *GameScreen) drawProgressMeter(min, max, fgcolor, bgcolor, width uint64) string {
// 	var blink bool
// 	if min > max {
// 		min = max
// 		blink = true
// 	}
// 	proportion := float64(float64(min) / float64(max))
// 	if math.IsNaN(proportion) {
// 		proportion = 0.0
// 	} else if proportion < 0.05 {
// 		blink = true
// 	}
// 	onWidth := uint64(float64(width) * proportion)
// 	offWidth := uint64(float64(width) * (1.0 - proportion))

// 	onColor := screen.colorFunc(fmt.Sprintf("%v:%v", fgcolor, bgcolor))
// 	offColor := onColor

// 	if blink {
// 		onColor = screen.colorFunc(fmt.Sprintf("%v+B:%v", fgcolor, bgcolor))
// 	}

// 	if (onWidth + offWidth) > width {
// 		onWidth = width
// 		offWidth = 0
// 	} else if (onWidth + offWidth) < width {
// 		onWidth += width - (onWidth + offWidth)
// 	}

// 	on := ""
// 	off := ""

// 	for i := 0; i < int(onWidth); i++ {
// 		on += hpon
// 	}

// 	for i := 0; i < int(offWidth); i++ {
// 		off += hpoff
// 	}

// 	return onColor(on) + offColor(off)
// }

func (screen *GameScreen) drawFill(x, y, width, height int) {
	color := ansi.ColorCode(fmt.Sprintf("0:%v", bgcolor))

	midString := fmt.Sprintf("%%s%%s%%%vs", width)
	for i := 0; i <= height; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf(midString, cursor.MoveTo(y+i, x), color, " "))
	}
}

func (screen *GameScreen) drawBox(x, y, width, height int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))

	for i := 1; i < width; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s─", cursor.MoveTo(y, x+i), color))
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s─", cursor.MoveTo(y+height, x+i), color))
	}

	for i := 1; i < height; i++ {
		midString := fmt.Sprintf("%%s%%s│%%%vs│", (width - 1))
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x), color))
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x+width), color))
		io.WriteString(os.Stdout, fmt.Sprintf(midString, cursor.MoveTo(y+i, x), color, " "))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╭", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╰", cursor.MoveTo(y+height, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╮", cursor.MoveTo(y, x+width), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╯", cursor.MoveTo(y+height, x+width), color))
}

// SetScreenSize do handle the case user resize the terminal
func (screen *GameScreen) SetScreenSize(Width, Height int) {
	screen.screenSize = ssh.Window{
		Width:  Width,
		Height: Height,
	}
	screen.FreshRender()
}

func (screen *GameScreen) colorFunc(color string) func(string) string {
	_, ok := screen.colorCodeCache[color]

	if !ok {
		screen.colorCodeCache[color] = ansi.ColorFunc(color)
	}

	return screen.colorCodeCache[color]
}

// IsResultScreen returns if current screen is result screen
func (screen *GameScreen) IsResultScreen() bool {
	return screen.scrStatus.IsResultScreen()
}

// IsHelpScreen returns if current screen is help screen
func (screen *GameScreen) IsHelpScreen() bool {
	return screen.scrStatus.IsHelpScreen()
}

// IsWaitScreen returns if current screen is wait screen
func (screen *GameScreen) IsWaitScreen() bool {
	return screen.scrStatus.IsWaitScreen()
}

// IsWaitScreenCmd returns if input is processable on wait screen
func (screen *GameScreen) IsWaitScreenCmd(input termbox.Event) bool {
	if input.Key == termbox.KeyEsc {
		return true
	}
	Key := strings.ToUpper(string(input.Ch))
	switch Key {
	case "E", "M", "L": // Refresh, Cosmos address copy, TxHash copy
		return true
	}
	return false
}

// IsEndGameConfirmScreen returns if user is seeing end game confirmation page
func (screen *GameScreen) IsEndGameConfirmScreen() bool {
	return screen.scrStatus == CONFIRM_ENDGAME
}

// InputActive returns if screen's input area is active
func (screen *GameScreen) InputActive() bool {
	switch screen.scrStatus {
	case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL,
		CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL,
		CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL,
		CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL,
		CR8_SELLITM_TRDREQ_ENT_PYLVAL,
		CR8_BUYITM_TRDREQ_ENT_PYLVAL,
		CR8_SELLCHR_TRDREQ_ENT_PYLVAL,
		CR8_BUYCHR_TRDREQ_ENT_PYLVAL,
		RENAME_CHAR_ENT_NEWNAME:
		return true
	}
	return false
}

// SetScreenStatusAndRefresh set screen status and do refresh
func (screen *GameScreen) SetScreenStatusAndRefresh(newStatus ScreenStatus) {
	screen.SetScreenStatus(newStatus)
	screen.Render()
}

// FreshRender do fresh render and used when user resize screen or etc that's unusal
func (screen *GameScreen) FreshRender() {
	screen.refreshed = false
	screen.Render()
}

// GetScreenStatus returns screen status
func (screen *GameScreen) GetScreenStatus() ScreenStatus {
	return screen.scrStatus
}

// SelectDefaultActiveLine select activeLine to 0 when it's not set
func (screen *GameScreen) SelectDefaultActiveLine(arrayInterface interface{}) {
	array := InterfaceSlice(arrayInterface)
	if len(array) > 0 && screen.activeLine == -1 {
		screen.activeLine = 0
	}
}

// SetScreenStatus select the screen status and do the intercept operations while switch
func (screen *GameScreen) SetScreenStatus(newStatus ScreenStatus) {
	screen.scrStatus = newStatus

	switch newStatus {
	case RSLT_HUNT_RABBITS,
		RSLT_FIGHT_GOBLIN,
		RSLT_FIGHT_TROLL,
		RSLT_FIGHT_WOLF,
		RSLT_FIGHT_GIANT,
		RSLT_FIGHT_DRAGONFIRE,
		RSLT_FIGHT_DRAGONICE,
		RSLT_FIGHT_DRAGONACID,
		RSLT_FIGHT_DRAGONUNDEAD:
		_, respOutput := screen.GetTxResponseOutput()
		resLen := len(respOutput)
		if resLen == 0 { // it means character is dead
			screen.user.SetActiveCharacterIndex(-1)
		}
	case SEL_ACTIVE_CHAR, SEL_RENAME_CHAR:
		screen.activeLine = screen.user.GetActiveCharacterIndex()
		screen.SelectDefaultActiveLine(screen.user.InventoryCharacters())
	case SEL_BUYCHR:
		screen.activeLine = 0
	case SHW_LOUD_BUY_TRDREQS:
		screen.SelectDefaultActiveLine(loud.BuyTrdReqs)
	case SHW_LOUD_SELL_TRDREQS:
		screen.SelectDefaultActiveLine(loud.SellTrdReqs)
	case SHW_BUYITM_TRDREQS:
		screen.SelectDefaultActiveLine(loud.ItemBuyTrdReqs)
	case SHW_SELLITM_TRDREQS:
		screen.SelectDefaultActiveLine(loud.ItemSellTrdReqs)
	case SHW_BUYCHR_TRDREQS:
		screen.SelectDefaultActiveLine(loud.CharacterBuyTrdReqs)
	case SHW_SELLCHR_TRDREQS:
		screen.SelectDefaultActiveLine(loud.CharacterSellTrdReqs)
	case CR8_BUYCHR_TRDREQ_SEL_CHR:
		screen.SelectDefaultActiveLine(loud.WorldCharacterSpecs)
	case CR8_SELLCHR_TRDREQ_SEL_CHR:
		screen.SelectDefaultActiveLine(screen.user.InventoryCharacters())
	case CR8_BUYITM_TRDREQ_SEL_ITEM:
		screen.SelectDefaultActiveLine(loud.WorldItemSpecs)
	case CR8_SELLITM_TRDREQ_SEL_ITEM:
		screen.SelectDefaultActiveLine(screen.user.InventoryItems())
	}
}

// Reset reset the screen stdout mode
func (screen *GameScreen) Reset() {
	io.WriteString(os.Stdout, fmt.Sprintf("%s👋\n", resetScreen))
}

// SaveGame saves the game status into file
func (screen *GameScreen) SaveGame() {
	screen.user.Save()
}

// UpdateFakeBlockHeight update the estimation block height for visual
func (screen *GameScreen) UpdateFakeBlockHeight(h int64) {
	screen.fakeBlockHeight = h
	screen.Render()
}

// BlockSince returns block offset from current block to specific block in the past
func (screen *GameScreen) BlockSince(h int64) uint64 {
	return uint64(screen.fakeBlockHeight - h)
}

// SetInputTextAndRender set the input text and render
func (screen *GameScreen) SetInputTextAndRender(text string) {
	screen.inputText = text
	screen.Render()
}
