package screen

type ScreenBox struct {
	X int
	Y int
	W int // inner width
	H int // inner height
}

func (screen *GameScreen) Width() int {
	return screen.screenSize.Width
}

func (screen *GameScreen) Height() int {
	return screen.screenSize.Height
}

func (screen *GameScreen) GetMenuBox() ScreenBox {
	return ScreenBox{
		X: 2,
		Y: 2,
		W: screen.Width() - 2,
		H: 1,
	}
}

func (screen *GameScreen) GetSituationBox() ScreenBox {
	y := screen.situationCmdBorderY() + 1
	return ScreenBox{
		X: 2,
		Y: y,
		W: screen.leftInnerWidth(),
		H: screen.Height() - y - 3,
	}
}

func (screen *GameScreen) GetCmdBox() ScreenBox {
	return ScreenBox{
		X: 2,
		Y: 4,
		W: screen.leftInnerWidth(),
		H: screen.situationCmdBorderY() - 4,
	}
}

func (screen *GameScreen) GetCharacterSheetBox() ScreenBox {
	return ScreenBox{
		X: screen.rightInnerStartX(),
		Y: 4,
		W: screen.rightInnerWidth(),
		H: screen.Height() - 4,
	}
}

func (screen *GameScreen) leftRightBorderX() int {
	return screen.Width() - 40
}

func (screen *GameScreen) leftInnerWidth() int {
	return screen.leftRightBorderX() - 2
}

func (screen *GameScreen) rightInnerStartX() int {
	return screen.leftRightBorderX() + 1
}

func (screen *GameScreen) rightInnerEndX() int {
	return screen.Width() - 1
}

func (screen *GameScreen) rightInnerWidth() int {
	return screen.Width() - screen.leftInnerWidth() - 3
}

func (screen *GameScreen) situationCmdBorderY() int {
	return 14
}

func (screen *GameScreen) situationInputBorderY() int {
	return screen.Height() - 2
}
