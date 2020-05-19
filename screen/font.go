package screen

import (
	"fmt"

	loud "github.com/Pylons-tech/LOUD/data"
)

type FontType string

const (
	REGULAR           FontType = ""
	GREY                       = "grey"
	BROWN                      = "brown"
	RED                        = "red"
	RED_BOLD                   = "red_bold"
	YELLOW                     = "yellow"
	GREEN                      = "green"
	BLINK_BLUE_BOLD            = "blink_blue_bold"
	INPUT_ACTIVE_FONT          = "input_active_font"
	BROWN_BOLD                 = "brown_bold"
	BLUE_BOLD                  = "blue_bold"
	GREY_BOLD                  = "grey_bold"
)

func (screen *GameScreen) getFont(ft FontType) func(string) string {
	switch ft {
	case REGULAR:
		return screen.regularFont()
	case GREY:
		return screen.greyFont()
	case BROWN:
		return screen.brownFont()
	case RED:
		return screen.redFont()
	case YELLOW:
		return screen.yellowFont()
	case GREEN:
		return screen.greenFont()
	case RED_BOLD:
		return screen.redBoldFont()
	case BLUE_BOLD:
		return screen.blueBoldFont()
	case GREY_BOLD:
		return screen.greyBoldFont()
	case BLINK_BLUE_BOLD:
		return screen.blinkBlueBoldFont()
	case INPUT_ACTIVE_FONT:
		return screen.inputActiveFont()
	case BROWN_BOLD:
		return screen.brownBoldFont()
	default:
		return screen.regularFont()
	}
}

func (screen *GameScreen) redFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v:%v", 196, bgcolor))
}

func (screen *GameScreen) yellowFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v:%v", 208, bgcolor))
}

func (screen *GameScreen) greenFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v:%v", 76, bgcolor))
}

func (screen *GameScreen) redBoldFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+bh:%v", 196, bgcolor))
}

func (screen *GameScreen) blueBoldFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+bh:%v", 117, 232))
}

func (screen *GameScreen) greyBoldFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+bh:%v", 181, 232))
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

func (screen *GameScreen) menuRegularFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+bh:%v", 255, 0))
}

func (screen *GameScreen) menuActiveFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+bh:%v", 255, 202))
}

func (screen *GameScreen) blinkBlueBoldFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("%v+B:%v", 117, bgcolor))
}

func (screen *GameScreen) inputActiveFont() func(string) string {
	return screen.colorFunc(fmt.Sprintf("0+b:%v", bgcolor-1))
}

func (screen *GameScreen) getFontByActiveIndex(idx int) FontType {
	activeLine := screen.activeLine
	font := REGULAR
	if activeLine == idx {
		font = BLUE_BOLD
	}
	return font
}

func (screen *GameScreen) getFontOfTR(idx int, isMyTR bool) FontType {
	font := REGULAR
	isActiveLine := screen.activeLine == idx
	isDisabledLine := isMyTR
	if isActiveLine && isDisabledLine {
		font = BROWN_BOLD
	} else if isActiveLine {
		font = BLUE_BOLD
	} else if isDisabledLine {
		font = BROWN
	}
	return font
}

func (screen *GameScreen) getFontOfShopItem(idx int, item loud.Item) FontType {
	font := REGULAR
	switch {
	case !screen.user.HasPreItemForAnItem(item): // ! preitem ok
		font = GREY
	case !(item.Price <= screen.user.GetGold()): // ! gold enough
		font = GREY
	}
	if idx == screen.activeLine {
		switch font {
		case REGULAR:
			font = BLUE_BOLD
		case GREY:
			font = GREY_BOLD
		}
	}
	return font
}
