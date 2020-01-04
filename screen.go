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
	"golang.org/x/crypto/ssh/terminal"
	"github.com/nsf/termbox-go"
)

// Screen represents a UI screen.
type Screen interface {
	SetScreenSize(int, int)
	HandleInputKey(termbox.Key)
	Render()
	Reset()
}

type GameScreen struct {
	world          World
	user           User
	screenSize     ssh.Window
	refreshed      bool
	colorCodeCache map[string](func(string) string)
}

const allowMouseInputAndHideCursor string = "\x1b[?1003h\x1b[?25l"
const resetScreen string = "\x1bc"
const ellipsis = "…"
const hpon = "◆"
const hpoff = "◇"
const bgcolor = 232

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

func (screen *GameScreen) drawFill(x, y, width, height int) {
	color := ansi.ColorCode(fmt.Sprintf("0:%v", bgcolor))

	midString := fmt.Sprintf("%%s%%s%%%vs", (width))
	for i := 0; i <= height; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf(midString, cursor.MoveTo(y+i, x), color, " "))
	}
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

func (screen *GameScreen) renderCharacterSheet(slotKeys map[string]func()) {
	var HP uint64 = 10
	var MaxHP uint64 = 10
	bgcolor := uint64(bgcolor)
	warning := ""
	if float32(HP) < float32(MaxHP)*.25 {
		bgcolor = 124
		warning = " (Health low) "
	} else if float32(HP) < float32(MaxHP)*.1 {
		bgcolor = 160
		warning = " (Health CRITICAL) "
	}

	x := screen.screenSize.Width/2 - 1
	width := (screen.screenSize.Width - x)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))

	infoLines := []string{
		centerText(fmt.Sprintf("%v the %v", "eugen", "Unworthy"), " ", width),
		centerText(warning, "─", width),
		truncateRight(fmt.Sprintf("%s (%v, %v)", "Delaware", 111, 222), width),
		truncateRight(fmt.Sprintf("Charge: %v/%v", "0", "1"), width),
		screen.drawProgressMeter(HP, MaxHP, 196, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" HP: %v/%v", HP, MaxHP), width-10)),
		screen.drawProgressMeter(HP, 6, 225, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" XP: %v/%v", HP, 10), width-10)),
		screen.drawProgressMeter(HP, MaxHP, 208, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" AP: %v/%v", HP, MaxHP), width-10)),
		screen.drawProgressMeter(HP, MaxHP, 117, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" RP: %v/%v", HP, MaxHP), width-10)),
		screen.drawProgressMeter(HP, MaxHP, 76, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" MP: %v/%v", HP, MaxHP), width-10))}

	infoLines = append(infoLines, centerText(" ❦ ", "─", width))

	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(2+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}

	lastLine := len(infoLines) + 1
	screen.drawFill(x, lastLine+1, width, screen.screenSize.Height-(lastLine+2))
}

func (screen *GameScreen) HandleInputKey(input termbox.Key) {
	screen.Render()
}

func (screen *GameScreen) Render() {
	var HP uint64 = 10

	// screen.user.Reload()

	if screen.screenSize.Height < 20 || screen.screenSize.Width < 60 {
		clear := cursor.ClearEntireScreen()
		move := cursor.MoveTo(1, 1)
		io.WriteString(os.Stdout,
			fmt.Sprintf("%s%sScreen is too small. Make your terminal larger. (60x20 minimum)", clear, move))
		return
	} else if HP == 0 {
		clear := cursor.ClearEntireScreen()
		dead := "You died. Respawning..."
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

	var slotKeys map[string]func()

	screen.renderCharacterSheet(slotKeys)
}

func (screen *GameScreen) Reset() {
	io.WriteString(os.Stdout, fmt.Sprintf("%s👋\n", resetScreen))
}

// NewScreen manages the window rendering for game
func NewScreen(world World, user User) Screen {
	width, height, _ := terminal.GetSize(0)

	window := ssh.Window{
		Width:  width,
		Height: height,
	}

	screen := GameScreen{
		world:          world,
		user:           user,
		screenSize:     window,
		colorCodeCache: make(map[string](func(string) string))}

	return &screen
}
