package screen

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/ahmetb/go-cursor"
	"github.com/gliderlabs/ssh"
	"github.com/mgutz/ansi"
	"github.com/nsf/termbox-go"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/LOUD/log"
)

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

func (screen *GameScreen) FakeSync() {
	screen.UpdateFakeBlockHeight(screen.fakeBlockHeight + 1)
	screen.Render()
}

func (screen *GameScreen) GetTxFailReason() string {
	return screen.txFailReason
}

func (screen *GameScreen) SwitchUser(newUser loud.User) {
	screen.user = newUser
}

func (screen *GameScreen) drawProgressMeter(min, max, fgcolor, bgcolor, width uint64) string {
	var blink bool
	if min > max {
		min = max
		blink = true
	}
	proportion := float64(float64(min) / float64(max))
	if math.IsNaN(proportion) {
		proportion = 0.0
	} else if proportion < 0.05 {
		blink = true
	}
	onWidth := uint64(float64(width) * proportion)
	offWidth := uint64(float64(width) * (1.0 - proportion))

	onColor := screen.colorFunc(fmt.Sprintf("%v:%v", fgcolor, bgcolor))
	offColor := onColor

	if blink {
		onColor = screen.colorFunc(fmt.Sprintf("%v+B:%v", fgcolor, bgcolor))
	}

	if (onWidth + offWidth) > width {
		onWidth = width
		offWidth = 0
	} else if (onWidth + offWidth) < width {
		onWidth += width - (onWidth + offWidth)
	}

	on := ""
	off := ""

	for i := 0; i < int(onWidth); i++ {
		on += hpon
	}

	for i := 0; i < int(offWidth); i++ {
		off += hpoff
	}

	return onColor(on) + offColor(off)
}

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

func (screen *GameScreen) IsResultScreen() bool {
	return screen.scrStatus.IsResultScreen()
}

func (screen *GameScreen) IsHelpScreen() bool {
	return screen.scrStatus.IsHelpScreen()
}

func (screen *GameScreen) IsWaitScreen() bool {
	return screen.scrStatus.IsWaitScreen()
}

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

func (screen *GameScreen) IsEndGameConfirmScreen() bool {
	return screen.scrStatus == CONFIRM_ENDGAME
}

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

func (screen *GameScreen) SetScreenStatusAndRefresh(newStatus ScreenStatus) {
	screen.SetScreenStatus(newStatus)
	screen.Render()
}

func (screen *GameScreen) FreshRender() {
	screen.refreshed = false
	screen.Render()
}

func (screen *GameScreen) GetScreenStatus() ScreenStatus {
	return screen.scrStatus
}

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
		activeLine := screen.user.GetActiveCharacterIndex()
		if len(screen.user.InventoryCharacters()) > 0 && activeLine == -1 {
			activeLine = 0
		}
		screen.activeLine = activeLine
	case SEL_BUYCHR:
		screen.activeLine = 0
	case SHW_LOUD_BUY_TRDREQS:
		// TODO should have activeLine modifier by the array and use this across the app
		if len(loud.BuyTrdReqs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case SHW_LOUD_SELL_TRDREQS:
		if len(loud.SellTrdReqs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case SHW_BUYITM_TRDREQS:
		if len(loud.ItemBuyTrdReqs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case SHW_SELLITM_TRDREQS:
		if len(loud.ItemSellTrdReqs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case SHW_BUYCHR_TRDREQS:
		if len(loud.CharacterBuyTrdReqs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case SHW_SELLCHR_TRDREQS:
		if len(loud.CharacterSellTrdReqs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case CR8_BUYCHR_TRDREQ_SEL_CHR:
		if len(loud.WorldCharacterSpecs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case CR8_SELLCHR_TRDREQ_SEL_CHR:
		if len(screen.user.InventoryCharacters()) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case CR8_BUYITM_TRDREQ_SEL_ITEM:
		if len(loud.WorldItemSpecs) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	case CR8_SELLITM_TRDREQ_SEL_ITEM:
		if len(screen.user.InventoryItems()) > 0 && screen.activeLine == -1 {
			screen.activeLine = 0
		}
	}
}

func (screen *GameScreen) Reset() {
	io.WriteString(os.Stdout, fmt.Sprintf("%s👋\n", resetScreen))
}

func (screen *GameScreen) SaveGame() {
	screen.user.Save()
}

func (screen *GameScreen) UpdateFakeBlockHeight(h int64) {
	screen.fakeBlockHeight = h
	screen.Render()
}

func (screen *GameScreen) BlockSince(h int64) uint64 {
	return uint64(screen.fakeBlockHeight - h)
}

func (screen *GameScreen) SetInputTextAndRender(text string) {
	screen.inputText = text
	screen.Render()
}
