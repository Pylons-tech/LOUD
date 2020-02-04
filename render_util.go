package loud

import (
	"fmt"
	"io"
	"math"
	"os"
	"unicode/utf8"

	"github.com/ahmetb/go-cursor"
	"github.com/gliderlabs/ssh"
	"github.com/mgutz/ansi"
	"github.com/nsf/termbox-go"

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
	Render()
	Reset()
}

type ScreenStatus int

type GameScreen struct {
	world           World
	user            User
	screenSize      ssh.Window
	activeItem      Item
	lastInput       termbox.Event
	activeLine      int
	pylonEnterValue string
	loudEnterValue  string
	inputText       string
	blockHeight     int64
	txFailReason    string
	txResult        []byte
	refreshed       bool
	scrStatus       ScreenStatus
	colorCodeCache  map[string](func(string) string)
}

const (
	SHOW_LOCATION ScreenStatus = iota
	// in shop
	SELECT_SELL_ITEM
	WAIT_SELL_PROCESS
	RESULT_SELL_FINISH

	SELECT_BUY_ITEM
	WAIT_BUY_PROCESS
	RESULT_BUY_FINISH

	SELECT_UPGRADE_ITEM
	WAIT_UPGRADE_PROCESS
	RESULT_UPGRADE_FINISH
	// in forest
	SELECT_HUNT_ITEM
	WAIT_HUNT_PROCESS
	RESULT_HUNT_FINISH
	// in market
	SELECT_MARKET // buy loud or sell loud

	SHOW_LOUD_BUY_ORDERS                   // navigation using arrow and list should be sorted by price
	CREATE_BUY_LOUD_ORDER_ENTER_LOUD_VALUE // enter value after switching enter mode
	CREATE_BUY_LOUD_ORDER_ENTER_PYLON_VALUE
	WAIT_BUY_LOUD_ORDER_CREATION
	RESULT_BUY_LOUD_ORDER_CREATION
	WAIT_FULFILL_BUY_LOUD_ORDER // after done go to show loud buy orders
	RESULT_FULFILL_BUY_LOUD_ORDER

	SHOW_LOUD_SELL_ORDERS
	CREATE_SELL_LOUD_ORDER_ENTER_LOUD_VALUE
	CREATE_SELL_LOUD_ORDER_ENTER_PYLON_VALUE
	WAIT_SELL_LOUD_ORDER_CREATION
	RESULT_SELL_LOUD_ORDER_CREATION
	WAIT_FULFILL_SELL_LOUD_ORDER
	RESULT_FULFILL_SELL_LOUD_ORDER
)

// NewScreen manages the window rendering for game
func NewScreen(world World, user User) Screen {
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

func (screen *GameScreen) Reset() {
	io.WriteString(os.Stdout, fmt.Sprintf("%s👋\n", resetScreen))
}

func (screen *GameScreen) SaveGame() {
	screen.user.Save()
}

func (screen *GameScreen) UpdateBlockHeight(blockHeight int64) {
	screen.blockHeight = blockHeight
	screen.refreshed = false
	screen.Render()
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

func (screen *GameScreen) renderOrderTableLine(text1 string, text2 string, text3 string, isActiveLine bool, isDisabledLine bool) string {
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

func (screen *GameScreen) renderOrderTable(orders []Order) []string {
	infoLines := []string{}
	infoLines = append(infoLines, "╭────────────────────┬───────────────┬───────────────╮")
	// infoLines = append(infoLines, "│ LOUD price (pylon) │ Amount (loud) │ Total (pylon) │")
	infoLines = append(infoLines, screen.renderOrderTableLine("LOUD price (pylon)", "Amount (loud)", "Total (pylon)", false, false))
	infoLines = append(infoLines, "├────────────────────┼───────────────┼───────────────┤")
	numLines := screen.screenSize.Height/2 - 7
	if screen.activeLine >= len(orders) {
		screen.activeLine = len(orders) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(orders) {
		endLine = len(orders)
	}
	for li, order := range orders[startLine:endLine] {
		infoLines = append(
			infoLines,
			screen.renderOrderTableLine(
				fmt.Sprintf("%.4f", order.Price),
				fmt.Sprintf("%d", order.Amount),
				fmt.Sprintf("%d", order.Total),
				startLine+li == activeLine,
				order.IsMyOrder,
			),
		)
	}
	infoLines = append(infoLines, "╰────────────────────┴───────────────┴───────────────╯")
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

func truncateRight(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	return string([]rune(message)[0:width-1]) + ellipsis
}

func truncateLeft(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	strLen := utf8.RuneCountInString(message)
	return ellipsis + string([]rune(message)[strLen-width:strLen-1])
}

func justifyRight(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	strLen := utf8.RuneCountInString(message)
	return ellipsis + string([]rune(message)[strLen-width:strLen-1])
}

func centerText(message, pad string, width int) string {
	if utf8.RuneCountInString(message) > width {
		return truncateRight(message, width)
	}
	leftover := width - utf8.RuneCountInString(message)
	left := leftover / 2
	right := leftover - left

	if pad == "" {
		pad = " "
	}

	leftString := ""
	for utf8.RuneCountInString(leftString) <= left && utf8.RuneCountInString(leftString) <= right {
		leftString += pad
	}

	return fmt.Sprintf("%s%s%s", string([]rune(leftString)[0:left]), message, string([]rune(leftString)[0:right]))
}
