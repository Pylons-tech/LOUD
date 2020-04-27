package screen

func (screen *GameScreen) Width() int {
	return screen.screenSize.Width
}

func (screen *GameScreen) Height() int {
	return screen.screenSize.Height
}

func (screen *GameScreen) leftRightBorderX() int {
	return screen.Width()/2 - 1
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
	return screen.Height() / 2
}

func (screen *GameScreen) situationInnerStartY() int {
	return 2
}

func (screen *GameScreen) situationInnerEndY() int {
	return screen.situationCmdBorderY() - 1
}

func (screen *GameScreen) situationInnerHeight() int {
	return screen.situationCmdBorderY() - 2
}

func (screen *GameScreen) cmdInnerStartY() int {
	return screen.situationCmdBorderY() + 1
}
func (screen *GameScreen) cmdInnerEndY() int {
	return screen.Height() - 4
}
func (screen *GameScreen) cmdInnerHeight() int {
	return screen.Height() - screen.cmdInnerStartY() - 3
}

func (screen *GameScreen) cmdInputBorderY() int {
	return screen.Height() - 2
}
