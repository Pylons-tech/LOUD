package screen

import (
	"fmt"
	"strings"

	"github.com/ahmetb/go-cursor"
	"github.com/gliderlabs/ssh"
	"github.com/mgutz/ansi"
	"github.com/nsf/termbox-go"

	loud "github.com/Pylons-tech/LOUD/data"
)

// Resync execute new sync from node
func (screen *GameScreen) Resync() {
	screen.syncingData = true
	screen.Render()
	go func() {
		loud.SyncFromNode(screen.user)
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
		PrintString(fmt.Sprintf(midString, cursor.MoveTo(y+i, x), color, " "))
	}
}

func (screen *GameScreen) drawBox(x, y, width, height int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))

	for i := 1; i < width; i++ {
		PrintString(fmt.Sprintf("%s%s─", cursor.MoveTo(y, x+i), color))
		PrintString(fmt.Sprintf("%s%s─", cursor.MoveTo(y+height, x+i), color))
	}

	for i := 1; i < height; i++ {
		midString := fmt.Sprintf("%%s%%s│%%%vs│", (width - 1))
		PrintString(fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x), color))
		PrintString(fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x+width), color))
		PrintString(fmt.Sprintf(midString, cursor.MoveTo(y+i, x), color, " "))
	}

	PrintString(fmt.Sprintf("%s%s╭", cursor.MoveTo(y, x), color))
	PrintString(fmt.Sprintf("%s%s╰", cursor.MoveTo(y+height, x), color))
	PrintString(fmt.Sprintf("%s%s╮", cursor.MoveTo(y, x+width), color))
	PrintString(fmt.Sprintf("%s%s╯", cursor.MoveTo(y+height, x+width), color))
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
	return screen.scrStatus == ConfirmEndGame
}

// InputActive returns if screen's input area is active
func (screen *GameScreen) InputActive() bool {
	switch screen.scrStatus {
	case CreateBuyGoldTrdReqEnterGoldValue,
		CreateBuyGoldTrdReqEnterPylonValue,
		CreateSellGoldTrdReqEnterGoldValue,
		CreateSellGoldTrdReqEnterPylonValue,
		CreateSellItemTrdReqEnterPylonValue,
		CreateBuyItmTrdReqEnterPylonValue,
		CreateSellChrTrdReqEnterPylonValue,
		CreateBuyChrTrdReqEnterPylonValue,
		SelectRenameChrEntNewName,
		FriendRegisterEnterName,
		FriendRegisterEnterAddress:
		return true
	}
	return false
}

// SetScreenStatusAndRefresh set screen status and do refresh
func (screen *GameScreen) SetScreenStatusAndRefresh(newStatus PageStatus) {
	screen.SetScreenStatus(newStatus)
	screen.Render()
}

// FreshRender do fresh render and used when user resize screen or etc that's unusal
func (screen *GameScreen) FreshRender() {
	screen.refreshed = false
	screen.Render()
}

// GetScreenStatus returns screen status
func (screen *GameScreen) GetScreenStatus() PageStatus {
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
func (screen *GameScreen) SetScreenStatus(newStatus PageStatus) {
	screen.scrStatus = newStatus

	switch newStatus {
	case RsltHuntRabbits,
		RsltFightGoblin,
		RsltFightTroll,
		RsltFightWolf,
		RsltFightGiant,
		RsltFightDragonFire,
		RsltFightDragonIce,
		RsltFightDragonAcid,
		RsltFightDragonUndead:
		_, respOutput := screen.GetTxResponseOutput()
		resLen := len(respOutput)
		if resLen == 0 { // it means character is dead
			screen.user.SetActiveCharacterIndex(-1)
		}
	case SelectActiveChr, SelectRenameChr:
		screen.activeLine = screen.user.GetActiveCharacterIndex()
		screen.SelectDefaultActiveLine(screen.user.UnlockedCharacters())
	case SendItemSelectFriend, SendCharacterSelectFriend:
		screen.SelectDefaultActiveLine(screen.user.Friends())
	case SendItemSelectItem:
		screen.SelectDefaultActiveLine(screen.user.UnlockedItems())
	case SendCharacterSelectCharacter:
		screen.SelectDefaultActiveLine(screen.user.UnlockedCharacters())
	case SelectBuyChr:
		screen.activeLine = 0
	case ShowGoldBuyTrdReqs:
		screen.SelectDefaultActiveLine(loud.BuyTrdReqs)
	case ShowGoldSellTrdReqs:
		screen.SelectDefaultActiveLine(loud.SellTrdReqs)
	case ShowBuyItemTrdReqs:
		screen.SelectDefaultActiveLine(loud.ItemBuyTrdReqs)
	case SelectFitBuyItemTrdReq:
		atir := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
		matchingItems := screen.user.GetMatchedItems(atir.TItem)
		screen.SelectDefaultActiveLine(matchingItems)
	case ShowSellItemTrdReqs:
		screen.SelectDefaultActiveLine(loud.ItemSellTrdReqs)
	case ShowBuyChrTrdReqs:
		screen.SelectDefaultActiveLine(loud.CharacterBuyTrdReqs)
	case SelectFitBuyChrTrdReq:
		cbtr := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
		matchingChrs := screen.user.GetMatchedCharacters(cbtr.TCharacter)
		screen.SelectDefaultActiveLine(matchingChrs)
	case ShowSellChrTrdReqs:
		screen.SelectDefaultActiveLine(loud.CharacterSellTrdReqs)
	case CreateBuyChrTrdReqSelectChr:
		screen.SelectDefaultActiveLine(loud.WorldCharacterSpecs)
	case CreateSellChrTrdReqSelChr:
		screen.SelectDefaultActiveLine(screen.user.UnlockedCharacters())
	case CreateBuyItemTrdReqSelectItem:
		screen.SelectDefaultActiveLine(loud.WorldItemSpecs)
	case CreateSellItemTrdReqSelectItem:
		screen.SelectDefaultActiveLine(screen.user.UnlockedItems())
	}
}

// Reset reset the screen stdout mode
func (screen *GameScreen) Reset() {
	PrintString(fmt.Sprintf("%s👋\n", resetScreen))
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
