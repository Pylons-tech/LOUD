package screen

import (
	"fmt"

	"github.com/mgutz/ansi"
)

func (screen *GameScreen) redrawBorders() {
	PrintString(ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor)))
	screen.drawBox(1, 1, screen.Width()-1, screen.Height()-1)
	drawHorizontalLine(1, 3, screen.Width())
	drawVerticalLine(screen.leftRightBorderX(), 3, screen.Height())
	drawHorizontalLine(1, screen.situationCmdBorderY(), screen.leftInnerWidth()+1)
	drawHorizontalLine(1, screen.situationInputBorderY(), screen.leftInnerWidth()+1)
}
