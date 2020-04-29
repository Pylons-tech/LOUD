package screen

import (
	"fmt"
)

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
