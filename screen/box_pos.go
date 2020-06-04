package screen

// Box is a struct to manage regions on terminal
type Box struct {
	X int
	Y int
	W int // inner width
	H int // inner height
}

// Width returns whole screen width
func (screen *GameScreen) Width() int {
	return screen.screenSize.Width
}

// Height returns whole screen height
func (screen *GameScreen) Height() int {
	return screen.screenSize.Height
}

// GetMenuBox returns menubox region
func (screen *GameScreen) GetMenuBox() Box {
	return Box{
		X: 2,
		Y: 2,
		W: screen.Width() - 2,
		H: 1,
	}
}

// GetSituationBox returns situation box region
func (screen *GameScreen) GetSituationBox() Box {
	y := screen.situationCmdBorderY() + 1
	return Box{
		X: 2,
		Y: y,
		W: screen.leftInnerWidth(),
		H: screen.Height() - y - 3,
	}
}

// GetCmdBox returns command box region
func (screen *GameScreen) GetCmdBox() Box {
	return Box{
		X: 2,
		Y: 4,
		W: screen.leftInnerWidth(),
		H: screen.situationCmdBorderY() - 4,
	}
}

// GetCharacterSheetBox returns character box region
func (screen *GameScreen) GetCharacterSheetBox() Box {
	return Box{
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

// func (screen *GameScreen) rightInnerEndX() int {
// 	return screen.Width() - 1
// }

func (screen *GameScreen) rightInnerWidth() int {
	return screen.Width() - screen.leftInnerWidth() - 3
}

func (screen *GameScreen) situationCmdBorderY() int {
	return 14
}

func (screen *GameScreen) situationInputBorderY() int {
	return screen.Height() - 2
}
