package screen

import (
	"fmt"
)

type FontType string

const (
	REGULAR           FontType = ""
	GREY                       = "grey"
	BROWN                      = "brown"
	BLINK_BLUE_BOLD            = "blink_blue_bold"
	INPUT_ACTIVE_FONT          = "input_active_font"
	BROWN_BOLD                 = "brown_bold"
	BLUE_BOLD                  = "blue_bold"
)

func (screen *GameScreen) getFont(ft FontType) func(string) string {
	switch ft {
	case REGULAR:
		return screen.regularFont()
	case GREY:
		return screen.greyFont()
	case BROWN:
		return screen.brownFont()
	case BLINK_BLUE_BOLD:
		return screen.blinkBlueBoldFont()
	case INPUT_ACTIVE_FONT:
		return screen.inputActiveFont()
	case BROWN_BOLD:
		return screen.brownBoldFont()
	case BLUE_BOLD:
		return screen.blueBoldFont()
	default:
		return screen.regularFont()
	}
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

func (screen *GameScreen) regularFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))
}

func (screen *GameScreen) greyFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v:%v", 181, 232))
}

func (screen *GameScreen) blinkBlueBoldFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+B:%v", 117, bgcolor))
}

func (screen *GameScreen) inputActiveFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("0+b:%v", bgcolor-1))
}
