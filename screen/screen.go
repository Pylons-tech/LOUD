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
	inputText        string
	syncingData      bool
	blockHeight      int64
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

func (screen *GameScreen) UpdateBlockHeight(blockHeight int64) {
	screen.blockHeight = blockHeight
	screen.FreshRender()
}

func (screen *GameScreen) BlockSince(baseBlockHeight int64) uint64 {
	return uint64(screen.blockHeight - baseBlockHeight)
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

func (screen *GameScreen) renderTRTable(requests []loud.TrdReq) []string {
	infoLines := []string{}
	infoLines = append(infoLines, "╭────────────────────┬───────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ LOUD price (pylon) │ Amount (loud) │ Total (pylon) │")
	infoLines = append(infoLines, screen.renderTRLine("LOUD price (pylon)", "Amount (loud)", "Total (pylon)", false, false))
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
	infoLines := strings.Split(title, "\n")
	numHeaderLines := len(infoLines)
	infoLines = append(infoLines, "╭────────────────────────────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ Item                │ Price (pylon) │")
	infoLines = append(infoLines, screen.renderItemTrdReqTableLine(theads[0], theads[1], false, false))
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
	return screen.scrStatus.IsResultScreen()
}

func (screen *GameScreen) IsWaitScreen() bool {
	return screen.scrStatus.IsWaitScreen()
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
	var HP uint64 = 0
	var MaxHP uint64 = 0

	if lbh := screen.user.GetLatestBlockHeight(); lbh > screen.blockHeight {
		screen.blockHeight = lbh
	}

	characters := screen.user.InventoryCharacters()
	dfc := screen.user.GetDefaultCharacterIndex()
	if dfc >= 0 && dfc < len(characters) {
		DFC := characters[dfc]
		HP = uint64(DFC.HP)
		MaxHP = uint64(DFC.MaxHP)
		HP = min(HP+screen.BlockSince(DFC.LastUpdate), MaxHP)
	}
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
		itemInfo := truncateRight(formatItem(item), width)
		if idx == screen.user.GetDefaultItemIndex() {
			itemInfo = screen.blueBoldFont()(itemInfo)
		}
		infoLines = append(infoLines, itemInfo)
	}

	infoLines = append(infoLines, centerText(loud.Localize("inventory chracters"), "─", width))
	for idx, character := range characters {
		characterInfo := truncateRight(formatCharacter(character), width)
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
		centerText(loud.Localize("pylons network status")+" [Copy to Clipboard(L)]", " ", width),
		centerText(screen.user.GetLastTransaction(), " ", width),
	}

	blockHeightText := centerText(loud.Localize("block height")+": "+strconv.FormatInt(screen.blockHeight, 10), " ", width)
	if screen.syncingData {
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
	screen.SetScreenStatusAndRefresh(RSLT_SEL_DEF_CHAR)
}

func (screen *GameScreen) RunActiveWeaponSelect() {
	screen.user.SetDefaultItemIndex(screen.activeLine)
	screen.SetScreenStatusAndRefresh(RSLT_SEL_DEF_WEAPON)
}

func (screen *GameScreen) RunCharacterHealthRestore() {
	screen.RunTxProcess(W8_HEALTH_RESTORE_CHAR, RSLT_HEALTH_RESTORE_CHAR, func() (string, error) {
		return loud.RestoreHealth(screen.user, screen.activeCharacter)
	})
}

func (screen *GameScreen) RunCharacterRename(newName string) {
	screen.RunTxProcess(W8_RENAME_CHAR, RSLT_RENAME_CHAR, func() (string, error) {
		return loud.RenameCharacter(screen.user, screen.activeCharacter, newName)
	})
}

func (screen *GameScreen) RunActiveItemBuy() {
	screen.RunTxProcess(W8_BUYITM, RSLT_BUYITM, func() (string, error) {
		return loud.Buy(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveCharacterBuy() {
	screen.RunTxProcess(W8_BUYCHR, RSLT_BUYCHR, func() (string, error) {
		return loud.BuyCharacter(screen.user, screen.activeCharacter)
	})
}

func (screen *GameScreen) RunActiveItemSell() {
	screen.RunTxProcess(W8_SELLITM, RSLT_SELLITM, func() (string, error) {
		return loud.Sell(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemUpgrade() {
	screen.RunTxProcess(W8_UPGITM, RSLT_UPGITM, func() (string, error) {
		return loud.Upgrade(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemHunt() {
	screen.RunTxProcess(W8_HUNT, RSLT_HUNT, func() (string, error) {
		return loud.Hunt(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemFightGiant() {
	screen.RunTxProcess(W8_FIGHT_GIANT, RSLT_FIGHT_GIANT, func() (string, error) {
		return loud.FightGiant(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemFightTroll() {
	screen.RunTxProcess(W8_FIGHT_TROLL, RSLT_FIGHT_TROLL, func() (string, error) {
		return loud.FightTroll(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemFightWolf() {
	screen.RunTxProcess(W8_FIGHT_WOLF, RSLT_FIGHT_WOLF, func() (string, error) {
		return loud.FightWolf(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemFightGoblin() {
	screen.RunTxProcess(W8_FIGHT_GOBLIN, RSLT_FIGHT_GOBLIN, func() (string, error) {
		return loud.FightGoblin(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunSelectedLoudBuyTrdReq() {
	if len(loud.BuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		// when activeLine is not refering to real request but when it is refering to nil request
		screen.txFailReason = loud.Localize("you haven't selected any buy request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_BUY_LOUD_TRDREQ)
	} else {
		screen.activeTrdReq = loud.BuyTrdReqs[screen.activeLine]
		if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_BUY_LOUD_TRDREQ, RSLT_FULFILL_BUY_LOUD_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedLoudSellTrdReq() {
	if len(loud.SellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_SELL_LOUD_TRDREQ)
	} else {
		screen.activeTrdReq = loud.SellTrdReqs[screen.activeLine]
		if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_SELL_LOUD_TRDREQ, RSLT_FULFILL_SELL_LOUD_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedItemBuyTrdReq() {
	if len(loud.ItemBuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy item request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_BUYITM_TRDREQ)
	} else {
		atir := loud.ItemBuyTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = atir
		if atir.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, atir.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_BUYITM_TRDREQ, RSLT_FULFILL_BUYITM_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, atir.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedItemSellTrdReq() {
	if len(loud.ItemSellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell item request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_SELLITM_TRDREQ)
	} else {
		sstr := loud.ItemSellTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = sstr
		if sstr.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, sstr.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_SELLITM_TRDREQ, RSLT_FULFILL_SELLITM_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, sstr.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedCharacterBuyTrdReq() {
	if len(loud.CharacterBuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy character request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_BUYCHR_TRDREQ)
	} else {
		cbtr := loud.CharacterBuyTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = cbtr
		if cbtr.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, cbtr.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_BUYCHR_TRDREQ, RSLT_FULFILL_BUYCHR_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, cbtr.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedCharacterSellTrdReq() {
	if len(loud.CharacterSellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell character request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_SELLCHR_TRDREQ)
	} else {
		cstr := loud.CharacterSellTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = cstr
		if cstr.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, cstr.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_SELLCHR_TRDREQ, RSLT_FULFILL_SELLCHR_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, cstr.ID)
			})
		}
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
		screen.scrStatus = SHW_LOCATION
	}
	var HP uint64 = 10

	if screen.screenSize.Height < 38 || screen.screenSize.Width < 120 {
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
