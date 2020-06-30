package screen

import (
	"fmt"

	loud "github.com/Pylons-tech/LOUD/data"
)

// FontType is a struct to mange type of font
type FontType string

const (
	// RegularFont describes regular text font
	RegularFont FontType = ""
	// GreyFont describes grey text font
	GreyFont = "grey"
	// BrownFont describes brown text font
	BrownFont = "brown"
	// RedFont describes red text font
	RedFont = "red"
	// RedBoldFont describes red bold text font
	RedBoldFont = "red_bold"
	// YelloFont describes yellow text font
	YelloFont = "yellow"
	// GreenFont describes green text font
	GreenFont = "green"
	// BlinkBlueFont describes blinking blue bold text
	BlinkBlueFont = "blink_blue_bold"
	// InputActiveFont describes font for input enter text
	InputActiveFont = "input_active_font"
	// BrownBoldFont describes font for brown bold text
	BrownBoldFont = "brown_bold"
	// BlueBoldFont describes font for blue bold text
	BlueBoldFont = "blue_bold"
	// GreyBoldFont describes font for blue bold text
	GreyBoldFont = "grey_bold"
)

func (screen *GameScreen) getFont(ft FontType) func(string) string {
	switch ft {
	case RegularFont:
		return screen.regularFont()
	case GreyFont:
		return screen.greyFont()
	case BrownFont:
		return screen.brownFont()
	case RedFont:
		return screen.redFont()
	case YelloFont:
		return screen.yellowFont()
	case GreenFont:
		return screen.greenFont()
	case RedBoldFont:
		return screen.redBoldFont()
	case BlueBoldFont:
		return screen.blueBoldFont()
	case GreyBoldFont:
		return screen.greyBoldFont()
	case BlinkBlueFont:
		return screen.blinkBlueBoldFont()
	case InputActiveFont:
		return screen.inputActiveFont()
	case BrownBoldFont:
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
	font := RegularFont
	if activeLine == idx {
		font = BlueBoldFont
	}
	return font
}

func (screen *GameScreen) getFontOfTableLine(idx int, disabled bool) (FontType, string) {
	font, memo := RegularFont, "ok"
	isActiveLine := screen.activeLine == idx
	if isActiveLine && disabled {
		font, memo = BrownBoldFont, "disabled"
	} else if isActiveLine {
		font = BlueBoldFont
	} else if disabled {
		font, memo = BrownFont, "disabled"
	}
	return font, memo
}

func (screen *GameScreen) getFontOfShopItem(idx int, item loud.Item) (FontType, string) {
	font, memo := RegularFont, ""
	switch {
	case !screen.user.HasPreItemForAnItem(item): // ! preitem ok
		font, memo = GreyFont, "nopreitem"
	case !(item.Price <= screen.user.GetGold()): // ! gold enough
		font, memo = GreyFont, "goldlack"
	}
	if idx == screen.activeLine {
		switch font {
		case RegularFont:
			font = BlueBoldFont
		case GreyFont:
			font = GreyBoldFont
		}
	}
	return font, memo
}
