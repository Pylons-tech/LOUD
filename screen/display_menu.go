package screen

import (
	"fmt"
	"io"
	"os"

	"github.com/ahmetb/go-cursor"

	loud "github.com/Pylons-tech/LOUD/data"
)

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
		loud.HOME:     "home",
		loud.FOREST:   "forest",
		loud.SHOP:     "shop",
		loud.PYLCNTRL: "pylons central",
		loud.SETTINGS: "settings",
		loud.DEVELOP:  "develop",
		loud.HELP:     "help",
	}

	locations := []loud.UserLocation{
		loud.HOME,
		loud.FOREST,
		loud.SHOP,
		loud.PYLCNTRL,
		loud.SETTINGS,
		loud.DEVELOP,
		loud.HELP,
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
		text:     loud.Localize(EXIT_GAME_ESC_CMD),
		isActive: screen.scrStatus == CONFIRM_ENDGAME,
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
			io.WriteString(os.Stdout, splitText)
		} else {
			text = centerText(md.text, " ", md.width)
		}
		menuText := fmt.Sprintf("%s%s", move, menuFont(text))
		io.WriteString(os.Stdout, menuText)
	}
}
