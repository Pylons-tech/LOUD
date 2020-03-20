package screen

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
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
	SetDaemonFetchingFlag(bool)
	SaveGame()
	UpdateBlockHeight(int64)
	SetScreenSize(int, int)
	HandleInputKey(termbox.Event)
	GetScreenStatus() ScreenStatus
	SetScreenStatus(ScreenStatus)
	GetTxFailReason() string
	Resync()
	Render()
	Reset()
}

type GameScreen struct {
	world                       loud.World
	user                        loud.User
	screenSize                  ssh.Window
	activeItem                  loud.Item
	activeCharacter             loud.Character
	activeLine                  int
	activeTradeRequest          loud.TradeRequest
	activeItemTradeRequest      loud.ItemTradeRequest
	activeCharacterTradeRequest loud.CharacterTradeRequest
	pylonEnterValue             string
	loudEnterValue              string
	inputText                   string
	refreshingDaemonStatus      bool
	syncingData                 bool
	blockHeight                 int64
	txFailReason                string
	txResult                    []byte
	refreshed                   bool
	scrStatus                   ScreenStatus
	colorCodeCache              map[string](func(string) string)
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

func (screen *GameScreen) SetDaemonFetchingFlag(flag bool) {
	screen.refreshingDaemonStatus = flag
}

func (screen *GameScreen) UpdateBlockHeight(blockHeight int64) {
	screen.blockHeight = blockHeight
	screen.FreshRender()
}

func (screen *GameScreen) SetInputTextAndRender(text string) {
	screen.inputText = text
	screen.Render()
}

func (screen *GameScreen) pylonIcon() string {
	return screen.drawProgressMeter(1, 1, 117, bgcolor, 1)
}

func (screen *GameScreen) loudIcon() string {
	return screen.drawProgressMeter(1, 1, 208, bgcolor, 1)
}

func (screen *GameScreen) buyLoudDesc(loudValue interface{}, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		screen.loudIcon(),
		fmt.Sprintf("%v", loudValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellLoudDesc(loudValue interface{}, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.loudIcon(),
		fmt.Sprintf("%v", loudValue),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) buySwordDesc(activeItem loud.Item, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		fmt.Sprintf("%s", formatItem(activeItem)),
	}, "")
	return desc
}

func (screen *GameScreen) buyCharacterDesc(activeCharacter loud.Character, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		fmt.Sprintf("%s", formatCharacter(activeCharacter)),
	}, "")
	return desc
}

func (screen *GameScreen) sellSwordDesc(activeItem loud.Item, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		fmt.Sprintf("%s", formatItem(activeItem)),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellCharacterDesc(activeCharacter loud.Character, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		fmt.Sprintf("%s", formatCharacter(activeCharacter)),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) tradeTableColorDesc() []string {
	var infoLines = []string{}
	infoLines = append(infoLines, "white     ➝ other's request")
	infoLines = append(infoLines, screen.blueBoldFont()("bluebold")+"  ➝ selected request")
	infoLines = append(infoLines, screen.brownBoldFont()("brownbold")+" ➝ my request + selected")
	infoLines = append(infoLines, screen.brownFont()("brown")+"     ➝ my request")
	infoLines = append(infoLines, "\n")
	return infoLines
}

func (screen *GameScreen) redrawBorders() {
	io.WriteString(os.Stdout, ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor)))
	screen.drawBox(1, 1, screen.screenSize.Width-1, screen.screenSize.Height-1)
	screen.drawVerticalLine(screen.screenSize.Width/2-2, 1, screen.screenSize.Height)

	y := screen.screenSize.Height
	if y < 20 {
		y = 5
	} else {
		y = (y / 2) - 2
	}
	screen.drawHorizontalLine(1, y+2, screen.screenSize.Width/2-3)
	screen.drawHorizontalLine(1, screen.screenSize.Height-2, screen.screenSize.Width/2-3)
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

func (screen *GameScreen) renderTradeRequestTableLine(text1 string, text2 string, text3 string, isActiveLine bool, isDisabledLine bool) string {
	calcText := "│" + centerText(text1, " ", 20) + "│" + centerText(text2, " ", 15) + "│" + centerText(text3, " ", 15) + "│"
	if isActiveLine && isDisabledLine {
		onColor := screen.brownBoldFont()
		return onColor(calcText)
	} else if isActiveLine {
		onColor := screen.blueBoldFont()
		return onColor(calcText)
	} else if isDisabledLine {
		onColor := screen.brownFont()
		return onColor(calcText)
	}
	return calcText
}

func (screen *GameScreen) renderTradeRequestTable(requests []loud.TradeRequest) []string {
	infoLines := []string{}
	infoLines = append(infoLines, "╭────────────────────┬───────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ LOUD price (pylon) │ Amount (loud) │ Total (pylon) │")
	infoLines = append(infoLines, screen.renderTradeRequestTableLine("LOUD price (pylon)", "Amount (loud)", "Total (pylon)", false, false))
	infoLines = append(infoLines, "├────────────────────┼───────────────┼───────────────┤")
	numLines := screen.screenSize.Height/2 - 7
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
			screen.renderTradeRequestTableLine(
				fmt.Sprintf("%.4f", request.Price),
				fmt.Sprintf("%d", request.Amount),
				fmt.Sprintf("%d", request.Total),
				startLine+li == activeLine,
				request.IsMyTradeRequest,
			),
		)
	}
	infoLines = append(infoLines, "╰────────────────────┴───────────────┴───────────────╯")
	return infoLines
}

func (screen *GameScreen) renderItemTradeRequestTableLine(text1 string, text2 string, isActiveLine bool, isDisabledLine bool) string {
	calcText := "│" + centerText(text1, " ", 36) + "│" + centerText(text2, " ", 15) + "│"
	if isActiveLine && isDisabledLine {
		onColor := screen.brownBoldFont()
		return onColor(calcText)
	} else if isActiveLine {
		onColor := screen.blueBoldFont()
		return onColor(calcText)
	} else if isDisabledLine {
		onColor := screen.brownFont()
		return onColor(calcText)
	}
	return calcText
}

func (screen *GameScreen) renderItemTradeRequestTable(header string, requests []loud.ItemTradeRequest) []string {
	infoLines := strings.Split(header, "\n")
	numHeaderLines := len(infoLines)
	infoLines = append(infoLines, "╭────────────────────────────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ Item                │ Price (pylon) │")
	infoLines = append(infoLines, screen.renderItemTradeRequestTableLine("Item", "Price (pylon)", false, false))
	infoLines = append(infoLines, "├────────────────────────────────────┼───────────────┤")
	numLines := screen.screenSize.Height/2 - 7 - numHeaderLines
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
			screen.renderItemTradeRequestTableLine(
				fmt.Sprintf("%s  ", formatItem(request.TItem)),
				fmt.Sprintf("%d", request.Price),
				startLine+li == activeLine,
				request.IsMyTradeRequest,
			),
		)
	}
	infoLines = append(infoLines, "╰────────────────────────────────────┴───────────────╯")
	return infoLines
}

func (screen *GameScreen) renderCharacterTradeRequestTable(header string, requests []loud.CharacterTradeRequest) []string {
	infoLines := strings.Split(header, "\n")
	numHeaderLines := len(infoLines)
	infoLines = append(infoLines, "╭────────────────────────────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ Character                │ Price (pylon) │")
	infoLines = append(infoLines, screen.renderItemTradeRequestTableLine("Character", "Price (pylon)", false, false))
	infoLines = append(infoLines, "├────────────────────────────────────┼───────────────┤")
	numLines := screen.screenSize.Height/2 - 7 - numHeaderLines
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
			screen.renderItemTradeRequestTableLine(
				fmt.Sprintf("%s  ", formatCharacter(request.TCharacter)),
				fmt.Sprintf("%d", request.Price),
				startLine+li == activeLine,
				request.IsMyTradeRequest,
			),
		)
	}
	infoLines = append(infoLines, "╰────────────────────────────────────┴───────────────╯")
	return infoLines
}

func (screen *GameScreen) renderItemTableLine(text1 string, isActiveLine bool) string {
	calcText := "│" + centerText(text1, " ", 52) + "│"
	if isActiveLine {
		onColor := screen.blueBoldFont()
		return onColor(calcText)
	}
	return calcText
}

func (screen *GameScreen) renderItemTable(header string, items []loud.Item) []string {
	infoLines := strings.Split(header, "\n")
	numHeaderLines := len(infoLines)
	infoLines = append(infoLines, "╭────────────────────────────────────────────────────╮")
	// infoLines = append(infoLines, "│ Item                            │")
	infoLines = append(infoLines, screen.renderItemTableLine("Item", false))
	infoLines = append(infoLines, "├────────────────────────────────────────────────────┤")
	numLines := screen.screenSize.Height/2 - 7 - numHeaderLines
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
		infoLines = append(
			infoLines,
			screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatItem(item)),
				startLine+li == activeLine,
			),
		)
	}
	infoLines = append(infoLines, "╰────────────────────────────────────────────────────╯")
	return infoLines
}

func (screen *GameScreen) renderCharacterTable(header string, characters []loud.Character) []string {
	infoLines := strings.Split(header, "\n")
	numHeaderLines := len(infoLines)
	infoLines = append(infoLines, "╭────────────────────────────────────────────────────╮")
	// infoLines = append(infoLines, "│ Item                            │")
	infoLines = append(infoLines, screen.renderItemTableLine("Character", false))
	infoLines = append(infoLines, "├────────────────────────────────────────────────────┤")
	numLines := screen.screenSize.Height/2 - 7 - numHeaderLines
	if screen.activeLine >= len(characters) {
		screen.activeLine = len(characters) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(characters) {
		endLine = len(characters)
	}
	for li, character := range characters[startLine:endLine] {
		infoLines = append(
			infoLines,
			screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatCharacter(character)),
				startLine+li == activeLine,
			),
		)
	}
	infoLines = append(infoLines, "╰────────────────────────────────────────────────────╯")
	return infoLines
}

func (screen *GameScreen) drawVerticalLine(x, y, height int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))
	for i := 1; i < height; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x), color))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┬", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┴", cursor.MoveTo(y+height, x), color))
}

func (screen *GameScreen) drawHorizontalLine(x, y, width int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))
	for i := 1; i < width; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s─", cursor.MoveTo(y, x+i), color))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s├", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┤", cursor.MoveTo(y, x+width), color))
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
	screen.refreshed = false
}

func (screen *GameScreen) colorFunc(color string) func(string) string {
	_, ok := screen.colorCodeCache[color]

	if !ok {
		screen.colorCodeCache[color] = ansi.ColorFunc(color)
	}

	return screen.colorCodeCache[color]
}

func (screen *GameScreen) IsResultScreen() bool {
	return strings.Contains(string(screen.scrStatus), "RESULT_")
}
func (screen *GameScreen) InputActive() bool {
	switch screen.scrStatus {
	case CREATE_BUY_LOUD_REQUEST_ENTER_LOUD_VALUE:
		return true
	case CREATE_BUY_LOUD_REQUEST_ENTER_PYLON_VALUE:
		return true
	case CREATE_SELL_LOUD_REQUEST_ENTER_LOUD_VALUE:
		return true
	case CREATE_SELL_LOUD_REQUEST_ENTER_PYLON_VALUE:
		return true
	case CREATE_SELL_SWORD_REQUEST_ENTER_PYLON_VALUE:
		return true
	case CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE:
		return true
	case CREATE_SELL_CHARACTER_REQUEST_ENTER_PYLON_VALUE:
		return true
	case CREATE_BUY_CHARACTER_REQUEST_ENTER_PYLON_VALUE:
		return true
	}
	return false
}

func (screen *GameScreen) renderInputValue() {
	inputBoxWidth := uint32(screen.screenSize.Width/2) - 2
	inputWidth := inputBoxWidth - 9
	move := cursor.MoveTo(screen.screenSize.Height-1, 2)

	chatFunc := screen.colorFunc(fmt.Sprintf("231:%v", bgcolor))
	chat := chatFunc("INPUT▶ ")
	fmtString := fmt.Sprintf("%%-%vs", inputWidth)

	if screen.InputActive() {
		chatFunc = screen.colorFunc(fmt.Sprintf("0+b:%v", bgcolor-1))
	}

	fixedChat := truncateLeft(screen.inputText, int(inputWidth))

	inputText := fmt.Sprintf("%s%s%s", move, chat, chatFunc(fmt.Sprintf(fmtString, fixedChat)))

	io.WriteString(os.Stdout, inputText)
}

func (screen *GameScreen) renderCharacterSheet() {
	var HP uint64 = 10
	var MaxHP uint64 = 10
	bgcolor := uint64(bgcolor)
	warning := ""
	if float32(HP) < float32(MaxHP)*.25 {
		bgcolor = 124
		warning = loud.Localize("health low warning")
	} else if float32(HP) < float32(MaxHP)*.1 {
		bgcolor = 160
		warning = loud.Localize("health critical warning")
	}

	x := screen.screenSize.Width/2 - 1
	width := (screen.screenSize.Width - x)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))

	infoLines := []string{
		centerText(fmt.Sprintf("%v", screen.user.GetUserName()), " ", width),
		centerText(warning, "─", width),
		screen.pylonIcon() + fmtFunc(truncateRight(fmt.Sprintf(" %s: %v", "Pylon", screen.user.GetPylonAmount()), width-1)),
		screen.loudIcon() + fmtFunc(truncateRight(fmt.Sprintf(" %s: %v", loud.Localize("gold"), screen.user.GetGold()), width-1)),
		screen.drawProgressMeter(HP, MaxHP, 196, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" HP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 225, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" XP: %v/%v", HP, 10), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 208, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" AP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 117, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" RP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 76, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" MP: %v/%v", HP, MaxHP), width-10)),
	}

	infoLines = append(infoLines, centerText(loud.Localize("inventory items"), "─", width))
	items := screen.user.InventoryItems()
	for idx, item := range items {
		itemInfo := truncateRight(fmt.Sprintf("%s", formatItem(item)), width)
		if idx == screen.user.GetDefaultItemIndex() {
			itemInfo = screen.blueBoldFont()(itemInfo)
		}
		infoLines = append(infoLines, itemInfo)
	}

	infoLines = append(infoLines, centerText(loud.Localize("inventory chracters"), "─", width))
	characters := screen.user.InventoryCharacters()
	for idx, character := range characters {
		characterInfo := truncateRight(fmt.Sprintf("%s", formatCharacter(character)), width)
		if idx == screen.user.GetDefaultCharacterIndex() {
			characterInfo = screen.blueBoldFont()(characterInfo)
		}
		infoLines = append(infoLines, characterInfo)
	}
	infoLines = append(infoLines, centerText(" ❦ ", "─", width))

	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(2+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}

	nodeLines := []string{
		centerText(loud.Localize("pylons network status"), " ", width),
		centerText(screen.user.GetLastTransaction(), " ", width),
	}

	blockHeightText := centerText(loud.Localize("block height")+": "+strconv.FormatInt(screen.blockHeight, 10), " ", width)
	if screen.refreshingDaemonStatus {
		nodeLines = append(nodeLines, screen.blueBoldFont()(blockHeightText))
	} else {
		nodeLines = append(nodeLines, blockHeightText)
	}
	nodeLines = append(nodeLines, centerText(" ❦ ", "─", width))

	for index, line := range nodeLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(2+len(infoLines)+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}

	lastLine := len(infoLines) + len(nodeLines) + 1
	screen.drawFill(x, lastLine+1, width, screen.screenSize.Height-(lastLine+2))
}

func (screen *GameScreen) RunActiveCharacterSelect() {
	screen.user.SetDefaultCharacterIndex(screen.activeLine)
	screen.SetScreenStatusAndRefresh(RESULT_SELECT_DEF_CHAR)
}

func (screen *GameScreen) RunActiveWeaponSelect() {
	screen.user.SetDefaultItemIndex(screen.activeLine)
	screen.SetScreenStatusAndRefresh(RESULT_SELECT_DEF_WEAPON)
}

func (screen *GameScreen) RunActiveItemBuy() {
	screen.RunTxProcess(WAIT_BUY_ITEM_PROCESS, RESULT_BUY_ITEM_FINISH, func() (string, error) {
		return loud.Buy(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveCharacterBuy() {
	screen.RunTxProcess(WAIT_BUY_CHARACTER_PROCESS, RESULT_BUY_CHARACTER_FINISH, func() (string, error) {
		return loud.BuyCharacter(screen.user, screen.activeCharacter)
	})
}

func (screen *GameScreen) RunActiveItemSell() {
	screen.RunTxProcess(WAIT_SELL_PROCESS, RESULT_SELL_FINISH, func() (string, error) {
		return loud.Sell(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemUpgrade() {
	screen.RunTxProcess(WAIT_UPGRADE_PROCESS, RESULT_UPGRADE_FINISH, func() (string, error) {
		return loud.Upgrade(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemHunt() {
	screen.RunTxProcess(WAIT_HUNT_PROCESS, RESULT_HUNT_FINISH, func() (string, error) {
		return loud.Hunt(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunSelectedLoudBuyTrade() {
	if len(loud.BuyTradeRequests) <= screen.activeLine || screen.activeLine < 0 {
		// when activeLine is not refering to real request but when it is refering to nil request
		screen.txFailReason = loud.Localize("you haven't selected any buy request")
		screen.SetScreenStatusAndRefresh(RESULT_FULFILL_BUY_LOUD_REQUEST)
	} else {
		screen.activeTradeRequest = loud.BuyTradeRequests[screen.activeLine]
		screen.RunTxProcess(WAIT_FULFILL_BUY_LOUD_REQUEST, RESULT_FULFILL_BUY_LOUD_REQUEST, func() (string, error) {
			return loud.FulfillTrade(screen.user, screen.activeTradeRequest.ID)
		})
	}
}

func (screen *GameScreen) RunSelectedLoudSellTrade() {
	if len(loud.SellTradeRequests) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell request")
		screen.SetScreenStatusAndRefresh(RESULT_FULFILL_SELL_LOUD_REQUEST)
	} else {
		screen.activeTradeRequest = loud.SellTradeRequests[screen.activeLine]
		screen.RunTxProcess(WAIT_FULFILL_SELL_LOUD_REQUEST, RESULT_FULFILL_SELL_LOUD_REQUEST, func() (string, error) {
			return loud.FulfillTrade(screen.user, screen.activeTradeRequest.ID)
		})
	}
}

func (screen *GameScreen) RunSelectedSwordBuyTradeRequest() {
	if len(loud.SwordBuyTradeRequests) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy item request")
		screen.SetScreenStatusAndRefresh(RESULT_FULFILL_BUY_SWORD_REQUEST)
	} else {
		screen.activeItemTradeRequest = loud.SwordBuyTradeRequests[screen.activeLine]
		screen.RunTxProcess(WAIT_FULFILL_BUY_SWORD_REQUEST, RESULT_FULFILL_BUY_SWORD_REQUEST, func() (string, error) {
			return loud.FulfillTrade(screen.user, screen.activeItemTradeRequest.ID)
		})
	}
}

func (screen *GameScreen) RunSelectedSwordSellTradeRequest() {
	if len(loud.SwordSellTradeRequests) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell item request")
		screen.SetScreenStatusAndRefresh(RESULT_FULFILL_SELL_SWORD_REQUEST)
	} else {
		screen.activeItemTradeRequest = loud.SwordSellTradeRequests[screen.activeLine]
		screen.RunTxProcess(WAIT_FULFILL_SELL_SWORD_REQUEST, RESULT_FULFILL_SELL_SWORD_REQUEST, func() (string, error) {
			return loud.FulfillTrade(screen.user, screen.activeItemTradeRequest.ID)
		})
	}
}

func (screen *GameScreen) RunSelectedCharacterBuyTradeRequest() {
	if len(loud.CharacterBuyTradeRequests) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy character request")
		screen.SetScreenStatusAndRefresh(RESULT_FULFILL_BUY_CHARACTER_REQUEST)
	} else {
		screen.activeCharacterTradeRequest = loud.CharacterBuyTradeRequests[screen.activeLine]
		screen.RunTxProcess(WAIT_FULFILL_BUY_CHARACTER_REQUEST, RESULT_FULFILL_BUY_CHARACTER_REQUEST, func() (string, error) {
			return loud.FulfillTrade(screen.user, screen.activeCharacterTradeRequest.ID)
		})
	}
}

func (screen *GameScreen) RunSelectedCharacterSellTradeRequest() {
	if len(loud.CharacterSellTradeRequests) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell character request")
		screen.SetScreenStatusAndRefresh(RESULT_FULFILL_SELL_CHARACTER_REQUEST)
	} else {
		screen.activeCharacterTradeRequest = loud.CharacterSellTradeRequests[screen.activeLine]
		screen.RunTxProcess(WAIT_FULFILL_SELL_CHARACTER_REQUEST, RESULT_FULFILL_SELL_CHARACTER_REQUEST, func() (string, error) {
			return loud.FulfillTrade(screen.user, screen.activeCharacterTradeRequest.ID)
		})
	}
}

func (screen *GameScreen) SetScreenStatusAndRefresh(newStatus ScreenStatus) {
	screen.SetScreenStatus(newStatus)
	screen.FreshRender()
}

func (screen *GameScreen) FreshRender() {
	screen.refreshed = false
	screen.Render()
}

func (screen *GameScreen) RunTxProcess(waitStatus ScreenStatus, resultStatus ScreenStatus, fn func() (string, error)) {
	screen.SetScreenStatusAndRefresh(waitStatus)

	log.Println("started sending request for ", waitStatus)
	go func() {
		txhash, err := fn()
		log.Println("ended sending request for ", waitStatus)
		if err != nil {
			screen.txFailReason = err.Error()
			screen.SetScreenStatusAndRefresh(resultStatus)
		} else {
			time.AfterFunc(1*time.Second, func() {
				screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
				screen.SetScreenStatusAndRefresh(resultStatus)
			})
		}
	}()
}

func (screen *GameScreen) Render() {
	if len(loud.SomethingWentWrongMsg) > 0 {
		clear := cursor.ClearEntireScreen()
		dead := loud.Localize("Something went wrong, please close using Esc key and see loud.log")
		move := cursor.MoveTo(screen.screenSize.Height/2, screen.screenSize.Width/2-utf8.RuneCountInString(dead)/2)
		io.WriteString(os.Stdout, clear+move+dead)

		detailedErrorMsg := loud.Localize("detailed error: " + loud.SomethingWentWrongMsg)
		move = cursor.MoveTo(screen.screenSize.Height/2+3, screen.screenSize.Width/2-utf8.RuneCountInString(dead)/2)
		io.WriteString(os.Stdout, move+detailedErrorMsg)
		screen.refreshed = false
		return
	}
	if screen.scrStatus == "" {
		screen.scrStatus = SHOW_LOCATION
	}
	var HP uint64 = 10

	if screen.screenSize.Height < 20 || screen.screenSize.Width < 60 {
		clear := cursor.ClearEntireScreen()
		move := cursor.MoveTo(1, 1)
		io.WriteString(os.Stdout,
			fmt.Sprintf("%s%s%s", clear, move, loud.Localize("screen size warning")))
		return
	} else if HP == 0 {
		clear := cursor.ClearEntireScreen()
		dead := loud.Localize("dead")
		move := cursor.MoveTo(screen.screenSize.Height/2, screen.screenSize.Width/2-utf8.RuneCountInString(dead)/2)
		io.WriteString(os.Stdout, clear+move+dead)
		screen.refreshed = false
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
