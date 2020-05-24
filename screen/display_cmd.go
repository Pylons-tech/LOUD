package screen

import (
	"fmt"
	"io"
	"os"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/ahmetb/go-cursor"
)

const (
	SEL_CMD           = "Select ( ↵ )"
	GO_ON_ENTER_CMD   = "Go on ( ↵ )"
	FINISH_ENTER_CMD  = "Finish Enter ( ↵ )"
	GO_BACK_CMD       = "Go back ( ⌫ ) - Backspace Key"
	GO_BACK_ESC_CMD   = "Go back ( Esc )"
	EXIT_GAME_ESC_CMD = "Exit Game ( Esc )"
)

func (tl TextLines) appendSelectGoBackCmds() TextLines {
	return tl.appendT(
		SEL_CMD,
		GO_BACK_CMD)
}

func (tl TextLines) appendGoOnBackCmds() TextLines {
	return tl.appendT(
		GO_ON_ENTER_CMD,
		GO_BACK_CMD)
}

var MAX_SHORTCUT_ITEM_CMDSEL = 3

func getWindowFromActiveLine(activeLine, window_size, max_window int) (int, int) {
	if window_size > max_window {
		window_size = max_window
	}
	if activeLine >= max_window {
		activeLine = max_window - 1
	}
	startLine := activeLine - window_size/2
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + window_size
	if endLine >= max_window {
		startLine -= endLine - max_window
		endLine = max_window
	}
	return startLine, endLine
}

func (tl TextLines) appendCustomFontSelectCmds(itemsSlice interface{}, activeLine int, fn func(int, interface{}) TextLine) TextLines {
	items := InterfaceSlice(itemsSlice)

	startLine, endLine := getWindowFromActiveLine(activeLine, MAX_SHORTCUT_ITEM_CMDSEL, len(items))
	if startLine != 0 {
		tl = tl.append("...")
	} else {
		tl = tl.append("")
	}
	for li, item := range items[startLine:endLine] {
		idx := li + startLine
		fni := fn(idx, item)
		tl = append(tl, TextLine{
			content: fmt.Sprintf("%d) %s  ", idx+1, fni.content),
			font:    fni.font,
		})
	}
	if endLine < len(items) {
		tl = tl.append("...")
	} else {
		tl = tl.append("")
	}
	return tl
}

func (tl TextLines) appendCustomFontSelectCmdsScreenCharacters(screen *GameScreen) TextLines {
	return tl.appendCustomFontSelectCmds(
		screen.user.InventoryCharacters(), screen.activeLine,
		func(idx int, it interface{}) TextLine {
			return TextLine{
				content: formatCharacter(it.(loud.Character)),
				font:    screen.getFontByActiveIndex(idx),
			}
		})
}

func (screen *GameScreen) renderUserCommands() {
	// cmd box start point (x, y)
	scrBox := screen.GetCmdBox()
	x := scrBox.X
	y := scrBox.Y
	w := scrBox.W
	h := scrBox.H

	infoLines := TextLines{}
	tableLines := []string{}
	switch screen.scrStatus {
	case CONFIRM_ENDGAME:
		infoLines = infoLines.
			appendT(
				GO_BACK_ESC_CMD,
				GO_ON_ENTER_CMD)
	case SHW_LOCATION:
		cmdMap := map[loud.UserLocation]string{
			loud.HOME:     "home",
			loud.FOREST:   "forest",
			loud.SHOP:     "shop",
			loud.PYLCNTRL: "pylons central",
			loud.SETTINGS: "settings",
			loud.DEVELOP:  "develop",
			loud.HELP:     "help",
		}
		cmdString := loud.Localize(cmdMap[screen.user.GetLocation()])
		infoLines = infoLines.
			append(loud.ChunkText(cmdString, w)...)

		if screen.user.GetLocation() == loud.FOREST {
			forestStusMap := map[int]ScreenStatus{
				0: CONFIRM_HUNT_RABBITS,
				1: CONFIRM_FIGHT_GOBLIN,
				2: CONFIRM_FIGHT_WOLF,
				3: CONFIRM_FIGHT_TROLL,
				4: CONFIRM_FIGHT_GIANT,
				5: CONFIRM_FIGHT_DRAGONFIRE,
				6: CONFIRM_FIGHT_DRAGONICE,
				7: CONFIRM_FIGHT_DRAGONACID,
				8: CONFIRM_FIGHT_DRAGONUNDEAD,
			}

			for k, v := range forestStusMap {
				if _, fst := screen.ForestStatusCheck(v); len(fst) > 0 {
					infoLines[k].content += ": " + fst
					infoLines[k].font = GREY
				}
			}
		}
		if screen.user.GetLocation() == loud.SETTINGS && len(infoLines) > 2 {
			switch loud.GameLanguage {
			case "en":
				infoLines[1].font = BLUE_BOLD
			case "es":
				infoLines[2].font = BLUE_BOLD
			}
		}
		if screen.user.GetLocation() == loud.HOME {
			if len(infoLines) > 1 {
				// make it grey when no character's there
				infoLines[1].content += ": " + loud.Sprintf("no character!")
				infoLines[1].font = GREY
			}
		}
	case SHW_LOUD_BUY_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Sell gold to fulfill selected request( ↵ )",
				"Place order to buy gold(R)",
				GO_BACK_CMD)
	case SHW_LOUD_SELL_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Buy gold to fulfill selected request( ↵ )",
				"place order to sell gold(R)",
				GO_BACK_CMD)
	case SHW_BUYITM_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Sell item to fulfill selected request( ↵ )",
				"Place order to buy item(R)",
				GO_BACK_CMD)
	case SHW_SELLITM_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Buy item to fulfill selected request( ↵ )",
				"Place order to sell item(R)",
				GO_BACK_CMD)
	case SHW_BUYCHR_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Sell character to fulfill selected request( ↵ )",
				"Place order to buy character(R)",
				GO_BACK_CMD)
	case SHW_SELLCHR_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Buy character to fulfill selected request( ↵ )",
				"Place order to sell character(R)",
				GO_BACK_CMD)

	case CR8_BUYCHR_TRDREQ_SEL_CHR,
		CR8_SELLCHR_TRDREQ_SEL_CHR,
		CR8_SELLITM_TRDREQ_SEL_ITEM,
		CR8_BUYITM_TRDREQ_SEL_ITEM:
		infoLines = infoLines.
			appendT(
				SEL_CMD,
				GO_BACK_CMD)
	case SEL_RENAME_CHAR:
		infoLines = infoLines.
			appendCustomFontSelectCmdsScreenCharacters(screen).
			appendSelectGoBackCmds()
	case SEL_ACTIVE_CHAR:
		infoLines = infoLines.
			append(fmt.Sprintf("0) %s", loud.Localize("No character selection"))).
			appendCustomFontSelectCmdsScreenCharacters(screen).
			appendSelectGoBackCmds()
	case SEL_BUYITM:
		infoLines = infoLines.
			appendCustomFontSelectCmds(
				loud.ShopItems, screen.activeLine,
				func(idx int, it interface{}) TextLine {
					item := it.(loud.Item)
					font := screen.getFontOfShopItem(idx, it.(loud.Item))
					bonusText := ""
					contentStr := ""
					switch {
					case !screen.user.HasPreItemForAnItem(item): // ! preitem ok
						bonusText = fmt.Sprintf(": %s", loud.Localize("no material"))
					case !(item.Price <= screen.user.GetGold()): // ! gold enough
						bonusText = fmt.Sprintf(": %s", loud.Localize("not enough gold"))
					}
					if len(item.PreItems) > 0 {
						contentStr = formatItem(item) + fmt.Sprintf("💰 %d + %s %s", item.Price, item.PreItemStr(), bonusText)
					} else {
						contentStr = formatItem(item) + fmt.Sprintf("💰 %d %s", item.Price, bonusText)
					}
					return TextLine{
						content: contentStr,
						font:    font,
					}
				}).
			appendSelectGoBackCmds()
	case SEL_BUYCHR:
		infoLines = infoLines.
			appendCustomFontSelectCmds(
				loud.ShopCharacters, screen.activeLine,
				func(idx int, it interface{}) TextLine {
					char := it.(loud.Character)
					return TextLine{
						content: fmt.Sprintf("%s  🔷 %d", formatCharacter(char), char.Price),
						font:    screen.getFontByActiveIndex(idx),
					}
				}).
			appendSelectGoBackCmds()
	case SEL_SELLITM:
		infoLines = infoLines.
			appendCustomFontSelectCmds(
				screen.user.InventorySellableItems(), screen.activeLine,
				func(idx int, it interface{}) TextLine {
					item := it.(loud.Item)
					return TextLine{
						content: formatItem(item) + fmt.Sprintf("💰 %s", item.GetSellPriceRange()),
						font:    screen.getFontByActiveIndex(idx),
					}
				}).
			appendSelectGoBackCmds()
	case SEL_UPGITM:
		infoLines = infoLines.
			appendCustomFontSelectCmds(
				screen.user.InventoryUpgradableItems(), screen.activeLine,
				func(idx int, it interface{}) TextLine {
					item := it.(loud.Item)
					return TextLine{
						content: formatItem(item) + fmt.Sprintf("💰 %d", item.GetUpgradePrice()),
						font:    screen.getFontByActiveIndex(idx),
					}
				}).
			appendSelectGoBackCmds()
	case CONFIRM_HUNT_RABBITS,
		CONFIRM_FIGHT_GOBLIN,
		CONFIRM_FIGHT_TROLL,
		CONFIRM_FIGHT_WOLF,
		CONFIRM_FIGHT_GIANT,
		CONFIRM_FIGHT_DRAGONFIRE,
		CONFIRM_FIGHT_DRAGONICE,
		CONFIRM_FIGHT_DRAGONACID,
		CONFIRM_FIGHT_DRAGONUNDEAD:
		infoLines = infoLines.
			appendGoOnBackCmds()
	default:
		if screen.IsHelpScreen() {
			infoLines = infoLines.
				appendT(GO_BACK_CMD)
		} else if screen.IsResultScreen() { // eg. RSLT_BUY_LOUD_TRDREQ_CREATION
			infoLines = infoLines.appendT(GO_ON_ENTER_CMD)
		} else if screen.InputActive() { // eg. CR8_BUYITM_TRDREQ_ENT_PYLVAL
			infoLines = infoLines.
				appendT(
					FINISH_ENTER_CMD,
					GO_BACK_ESC_CMD)
		}
	}

	for index, line := range infoLines {
		lineFont := screen.getFont(line.font)
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x),
			lineFont(fillSpace(line.content, w))))
	}

	infoLen := len(infoLines)

	for index, line := range tableLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+infoLen+index, x),
			line))
	}
	totalLen := infoLen + len(tableLines)

	screen.drawFill(x, y+totalLen, w, h-totalLen-1)
}
