package screen

import (
	"fmt"

	"github.com/ahmetb/go-cursor"

	loud "github.com/Pylons-tech/LOUD/data"
)

// MenuDisplay describes the fields that are used for each menu item
type MenuDisplay struct {
	text     string
	isActive bool
	start    int
	width    int
	split    bool
}

func (screen *GameScreen) renderMenu() {
	scrBox := screen.GetMenuBox()
	x := scrBox.X
	y := scrBox.Y
	w := scrBox.W
	// h := scrBox.H

	cmdMap := map[loud.UserLocation]string{
		loud.Home:          "home",
		loud.Forest:        "forest",
		loud.Shop:          "shop",
		loud.PylonsCentral: "pylons central",
		loud.Friends:       "friends",
		loud.Settings:      "settings",
		loud.Develop:       "develop",
		loud.Help:          "help",
	}

	locations := []loud.UserLocation{
		loud.Home,
		loud.Forest,
		loud.Shop,
		loud.PylonsCentral,
		loud.Friends,
		loud.Settings,
		loud.Develop,
		loud.Help,
	}

	menuDisplays := []MenuDisplay{}
	// mw := w / (len(locations) + 1)
	mx := x
	for _, loc := range locations {
		md := MenuDisplay{
			text:     loud.Localize("go to " + cmdMap[loc]),
			isActive: loc == screen.user.GetLocation(),
			start:    x + mx,
			split:    true,
		}
		md.width = len(md.text) + 5
		md.start = mx
		menuDisplays = append(menuDisplays, md)

		mx += md.width
	}

	md := MenuDisplay{
		text:     loud.Localize(ExitGameEscCommand),
		isActive: screen.scrStatus == ConfirmEndGame,
	}
	md.width = len(md.text) + 4
	md.start = w - md.width

	menuDisplays = append(menuDisplays, md)

	for _, md := range menuDisplays {
		menuFont := screen.menuRegularFont()
		if md.isActive {
			menuFont = screen.menuActiveFont()
		}

		move := cursor.MoveTo(y, md.start)
		text := ""
		if md.split {
			text = centerText(md.text, " ", md.width-1)
			splitText := fmt.Sprintf("%s%s", cursor.MoveTo(y, md.start+md.width-1), screen.regularFont()("│"))
			PrintString(splitText)
		} else {
			text = centerText(md.text, " ", md.width)
		}
		menuText := fmt.Sprintf("%s%s", move, menuFont(text))
		PrintString(menuText)
	}
}
