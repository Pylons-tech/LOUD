package screen

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/ahmetb/go-cursor"
	"github.com/gliderlabs/ssh"
	"github.com/mgutz/ansi"
	"github.com/nsf/termbox-go"

	loud "github.com/Pylons-tech/LOUD/data"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

const allowMouseInputAndHideCursor string = "\x1b[?1003h\x1b[?25l"
const resetScreen string = "\x1bc"
const ellipsis = "…"
const hpon = "◆"
const hpoff = "◇"
const bgcolor = 232

// Screen represents a UI screen.
type Screen interface {
	SaveGame()
	IsEndGameConfirmScreen() bool
	UpdateFakeBlockHeight(int64)
	SetScreenSize(int, int)
	HandleInputKey(termbox.Event)
	GetScreenStatus() ScreenStatus
	SetScreenStatus(ScreenStatus)
	GetTxFailReason() string
	FakeSync()
	Resync()
	Render()
	Reset()
}

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
	scrStatus        ScreenStatus
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

func (screen *GameScreen) SwitchUser(newUser loud.User) {
	screen.user = newUser
}

func (screen *GameScreen) Resync() {
	screen.syncingData = true
	screen.FreshRender()
	go func() {
		log.Println("start syncing from node")
		loud.SyncFromNode(screen.user)
		log.Println("end syncing from node")
		screen.syncingData = false
		screen.FreshRender()
	}()
}

func (screen *GameScreen) FakeSync() {
	screen.UpdateFakeBlockHeight(screen.fakeBlockHeight + 1)
	screen.FreshRender()
}

func (screen *GameScreen) GetTxFailReason() string {
	return screen.txFailReason
}

func (screen *GameScreen) GetScreenStatus() ScreenStatus {
	return screen.scrStatus
}

func (screen *GameScreen) SetScreenStatus(newStatus ScreenStatus) {
	screen.scrStatus = newStatus
}

func (screen *GameScreen) Reset() {
	io.WriteString(os.Stdout, fmt.Sprintf("%s👋\n", resetScreen))
}

func (screen *GameScreen) SaveGame() {
	screen.user.Save()
}

func (screen *GameScreen) UpdateFakeBlockHeight(h int64) {
	screen.fakeBlockHeight = h
	screen.FreshRender()
}

func (screen *GameScreen) BlockSince(h int64) uint64 {
	return uint64(screen.fakeBlockHeight - h)
}

func (screen *GameScreen) SetInputTextAndRender(text string) {
	screen.inputText = text
	screen.Render()
}

func (screen *GameScreen) SetInputTextAndFreshRender(text string) {
	screen.refreshed = false
	screen.SetInputTextAndRender(text)
}

func (screen *GameScreen) pylonIcon() string {
	return screen.drawProgressMeter(1, 1, 117, bgcolor, 1)
}

func (screen *GameScreen) goldIcon() string {
	return "💰"
}

func (screen *GameScreen) buyLoudDesc(loudValue interface{}, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		screen.goldIcon(),
		fmt.Sprintf("%v", loudValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellLoudDesc(loudValue interface{}, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.goldIcon(),
		fmt.Sprintf("%v", loudValue),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) buyItemDesc(activeItem loud.Item, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatItem(activeItem),
	}, "")
	return desc
}

func (screen *GameScreen) buyItemSpecDesc(itemSpec loud.ItemSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatItemSpec(itemSpec),
	}, "")
	return desc
}

func (screen *GameScreen) buyCharacterDesc(activeCharacter loud.Character, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatCharacter(activeCharacter),
	}, "")
	return desc
}

func (screen *GameScreen) buyCharacterSpecDesc(charSpec loud.CharacterSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatCharacterSpec(charSpec),
	}, "")
	return desc
}

func (screen *GameScreen) sellItemDesc(activeItem loud.Item, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatItem(activeItem),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellCharacterDesc(activeCharacter loud.Character, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatCharacter(activeCharacter),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellItemSpecDesc(activeItem loud.ItemSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatItemSpec(activeItem),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellCharacterSpecDesc(activeCharacter loud.CharacterSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatCharacterSpec(activeCharacter),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) tradeTableColorDesc() []string {
	var infoLines = []string{}

	infoLines = append(infoLines, loud.Localize("white trade line desc"))
	infoLines = append(infoLines, screen.blueBoldFont()(loud.Localize("bluebold trade line desc")))
	infoLines = append(infoLines, screen.brownBoldFont()(loud.Localize("brownbold trade line desc")))
	infoLines = append(infoLines, screen.brownFont()(loud.Localize("brown trade line desc")))
	infoLines = append(infoLines, "\n")
	return infoLines
}

func (screen *GameScreen) redrawBorders() {
	io.WriteString(os.Stdout, ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor)))
	screen.drawBox(1, 1, screen.Width()-1, screen.Height()-1)
	drawVerticalLine(screen.leftRightBorderX(), 1, screen.Height())
	drawHorizontalLine(1, screen.situationCmdBorderY(), screen.leftInnerWidth())
	drawHorizontalLine(1, screen.cmdInputBorderY(), screen.leftInnerWidth())
}

func (screen *GameScreen) blueBoldFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+bh:%v", 117, 232))
}

func (screen *GameScreen) brownBoldFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+bh:%v", 181, 232))
}

func (screen *GameScreen) brownFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v:%v", 181, 232))
}

func (screen *GameScreen) renderTRTable(requests []loud.TrdReq) []string {
	infoLines := []string{}
	infoLines = append(infoLines, "╭────────────────────┬───────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ LOUD price (pylon) │ Amount (loud) │ Total (pylon) │")
	infoLines = append(infoLines, screen.renderTRLine("LOUD price (pylon)", "Amount (loud)", "Total (pylon)", false, false))
	infoLines = append(infoLines, "├────────────────────┼───────────────┼───────────────┤")
	numLines := screen.situationInnerHeight() - 5
	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(requests) {
		endLine = len(requests)
	}
	for li, request := range requests[startLine:endLine] {
		infoLines = append(
			infoLines,
			screen.renderTRLine(
				fmt.Sprintf("%.4f", request.Price),
				fmt.Sprintf("%d", request.Amount),
				fmt.Sprintf("%d", request.Total),
				startLine+li == activeLine,
				request.IsMyTrdReq,
			),
		)
	}
	infoLines = append(infoLines, "╰────────────────────┴───────────────┴───────────────╯")
	return infoLines
}

func (screen *GameScreen) renderITRTable(title string, theads [2]string, requestsSlice interface{}) []string {
	requests := InterfaceSlice(requestsSlice)
	infoLines := strings.Split(loud.Localize(title), "\n")
	numHeaderLines := len(infoLines)
	infoLines = append(infoLines, "╭────────────────────────────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ Item                │ Price (pylon) │")
	infoLines = append(infoLines, screen.renderItemTrdReqTableLine(theads[0], theads[1], false, false))
	infoLines = append(infoLines, "├────────────────────────────────────┼───────────────┤")
	numLines := screen.situationInnerHeight() - 5 - numHeaderLines
	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(requests) {
		endLine = len(requests)
	}
	for li, request := range requests[startLine:endLine] {
		line := ""
		switch request.(type) {
		case loud.ItemBuyTrdReq:
			itr := request.(loud.ItemBuyTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatItemSpec(itr.TItem)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		case loud.ItemSellTrdReq:
			itr := request.(loud.ItemSellTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatItem(itr.TItem)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		case loud.CharacterBuyTrdReq:
			itr := request.(loud.CharacterBuyTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatCharacterSpec(itr.TCharacter)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		case loud.CharacterSellTrdReq:
			itr := request.(loud.CharacterSellTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatCharacter(itr.TCharacter)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		}
		infoLines = append(infoLines, line)
	}
	infoLines = append(infoLines, "╰────────────────────────────────────┴───────────────╯")
	return infoLines
}

func (screen *GameScreen) renderITTable(header string, th string, itemSlice interface{}) []string {
	items := InterfaceSlice(itemSlice)
	infoLines := strings.Split(loud.Localize(header), "\n")
	numHeaderLines := len(infoLines)
	infoLines = append(infoLines, "╭────────────────────────────────────────────────────╮")
	// infoLines = append(infoLines, "│ Item                            │")
	infoLines = append(infoLines, screen.renderItemTableLine(th, false))
	infoLines = append(infoLines, "├────────────────────────────────────────────────────┤")
	numLines := screen.situationInnerHeight() - 5 - numHeaderLines
	if screen.activeLine >= len(items) {
		screen.activeLine = len(items) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(items) {
		endLine = len(items)
	}
	for li, item := range items[startLine:endLine] {
		line := ""
		switch item.(type) {
		case loud.Item:
			itemT := item.(loud.Item)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatItem(itemT)),
				startLine+li == activeLine,
			)
		case loud.Character:
			itemT := item.(loud.Character)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatCharacter(itemT)),
				startLine+li == activeLine,
			)
		case loud.ItemSpec:
			itemT := item.(loud.ItemSpec)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatItemSpec(itemT)),
				startLine+li == activeLine,
			)
		case loud.CharacterSpec:
			itemT := item.(loud.CharacterSpec)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatCharacterSpec(itemT)),
				startLine+li == activeLine,
			)
		}
		infoLines = append(infoLines, line)
	}
	infoLines = append(infoLines, "╰────────────────────────────────────────────────────╯")
	return infoLines
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

	midString := fmt.Sprintf("%%s%%s%%%vs", (width))
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

func (screen *GameScreen) renderInputValue() {
	inputBoxWidth := uint32(screen.leftInnerWidth())
	inputWidth := inputBoxWidth - 9
	move := cursor.MoveTo(screen.Height()-1, 2)

	chatFunc := screen.colorFunc(fmt.Sprintf("231:%v", bgcolor))
	chat := chatFunc("👉👉👉 ")
	fmtString := fmt.Sprintf("%%-%vs", inputWidth)

	if screen.InputActive() {
		chatFunc = screen.colorFunc(fmt.Sprintf("0+b:%v", bgcolor-1))
	}

	fixedChat := truncateLeft(screen.inputText, int(inputWidth))
	inputText := fmt.Sprintf("%s%s%s", move, chat, chatFunc(fmt.Sprintf(fmtString, fixedChat)))

	if !screen.InputActive() {
		inputText = fmt.Sprintf("%s%s", move, chatFunc(screen.actionText))
	}

	io.WriteString(os.Stdout, inputText)
}

func (screen *GameScreen) SetScreenStatusAndRefresh(newStatus ScreenStatus) {
	screen.SetScreenStatus(newStatus)
	screen.FreshRender()
}

func (screen *GameScreen) FreshRender() {
	screen.refreshed = false
	screen.Render()
}

func (screen *GameScreen) Render() {
	if len(loud.SomethingWentWrongMsg) > 0 {
		clear := cursor.ClearEntireScreen()
		dead := loud.Localize("Something went wrong, please close using Esc key and see loud.log")
		move := cursor.MoveTo(screen.Height()/2, screen.Width()/2-utf8.RuneCountInString(dead)/2)
		io.WriteString(os.Stdout, clear+move+dead)

		detailedErrorMsg := fmt.Sprintf("%s: %s", loud.Localize("detailed error"), loud.SomethingWentWrongMsg)
		move = cursor.MoveTo(screen.Height()/2+3, screen.Width()/2-utf8.RuneCountInString(dead)/2)
		io.WriteString(os.Stdout, move+detailedErrorMsg)
		screen.refreshed = false
		return
	}
	if screen.scrStatus == "" {
		screen.scrStatus = SHW_LOCATION
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
}
