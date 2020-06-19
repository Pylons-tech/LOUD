package screen

import (
	"fmt"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/ahmetb/go-cursor"
)

const (
	// SelectCommand express the select command text
	SelectCommand = "Select ( ↵ )"
	// UpdownCommand express table navigation command text
	UpdownCommand = "Table navigation (↑↓)"
	// GoonEnterCommand express go on enter command text
	GoonEnterCommand = "Go on ( ↵ )"
	// FinishEnterCommand express finish enter command text
	FinishEnterCommand = "Finish Enter ( ↵ )"
	// GobackCommand express go back command text
	GobackCommand = "Go back ( ⌫ ) - Backspace Key"
	// GobackEscCommand express go back by escape command
	GobackEscCommand = "Go back ( Esc )"
	// ExitGameEscCommand express exit game by escape key command
	ExitGameEscCommand = "Exit Game ( Esc )"
)

func (tl TextLines) appendSelectGoBackCmds() TextLines {
	return tl.appendT(
		SelectCommand,
		UpdownCommand,
		GobackCommand)
}

func (tl TextLines) appendGoOnBackCmds() TextLines {
	return tl.appendT(
		GoonEnterCommand,
		GobackCommand)
}

var maxShortcutItemCommands = 3

func getWindowFromActiveLine(activeLine, windowSize, maxWindow int) (int, int) {
	if windowSize > maxWindow {
		windowSize = maxWindow
	}
	if activeLine >= maxWindow {
		activeLine = maxWindow - 1
	}
	startLine := activeLine - windowSize/2
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + windowSize
	if endLine >= maxWindow {
		startLine -= endLine - maxWindow
		endLine = maxWindow
	}
	return startLine, endLine
}

func (tl TextLines) appendCustomFontSelectCmds(itemsSlice interface{}, activeLine int, fn func(int, interface{}) TextLine) TextLines {
	moreText := loud.Sprintf("use arrows for more")
	items := InterfaceSlice(itemsSlice)

	startLine, endLine := getWindowFromActiveLine(activeLine, maxShortcutItemCommands, len(items))
	if startLine != 0 {
		tl = tl.append("..." + " " + moreText)
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
		tl = tl.append("..." + " " + moreText)
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
	case ConfirmEndGame:
		infoLines = infoLines.
			appendT(
				GobackEscCommand,
				GoonEnterCommand)
	case ShowLocation:
		cmdMap := map[loud.UserLocation]string{
			loud.Home:          "home",
			loud.Forest:        "forest",
			loud.Shop:          "shop",
			loud.PylonsCentral: "pylons central",
			loud.Settings:      "settings",
			loud.Develop:       "develop",
			loud.Help:          "help",
		}
		cmdString := loud.Localize(cmdMap[screen.user.GetLocation()])
		infoLines = infoLines.
			append(loud.ChunkText(cmdString, w)...)

		if screen.user.GetLocation() == loud.Forest {
			forestStusMap := map[int]PageStatus{
				0: ConfirmHuntRabbits,
				1: ConfirmFightGoblin,
				2: ConfirmFightWolf,
				3: ConfirmFightTroll,
				4: ConfirmFightGiant,
				5: ConfirmFightDragonFire,
				6: ConfirmFightDragonIce,
				7: ConfirmFightDragonAcid,
				8: ConfirmFightDragonUndead,
			}

			for k, v := range forestStusMap {
				if _, fst := screen.ForestStatusCheck(v); len(fst) > 0 {
					infoLines[k].content += ": " + fst
					infoLines[k].font = GreyFont
				}
			}
		}
		if screen.user.GetLocation() == loud.Settings && len(infoLines) > 2 {
			switch loud.GameLanguage {
			case "en":
				infoLines[1].font = BlueBoldFont
			case "es":
				infoLines[2].font = BlueBoldFont
			}
		}
		if screen.user.GetLocation() == loud.Home {
			if len(infoLines) > 1 && len(screen.user.InventoryCharacters()) == 0 {
				// make it grey when no character's there
				infoLines[1].content += ": " + loud.Sprintf("no character!")
				infoLines[1].font = GreyFont
			}
		}
	case ShowGoldBuyTrdReqs:
		infoLines = infoLines.
			appendT(
				"Sell gold to fulfill selected request( ↵ )",
				"Place order to buy gold(R)",
				GobackCommand)
	case ShowGoldSellTrdReqs:
		infoLines = infoLines.
			appendT(
				"Buy gold to fulfill selected request( ↵ )",
				"place order to sell gold(R)",
				GobackCommand)
	case ShowBuyItemTrdReqs:
		infoLines = infoLines.
			appendT(
				"Sell item to fulfill selected request( ↵ )",
				"Place order to buy item(R)",
				GobackCommand)
	case SelectFitBuyItemTrdReq:
		infoLines = infoLines.
			appendT(
				"Sell item to fulfill selected request( ↵ )",
				"Place order to buy item(R)",
				GobackCommand)
	case ShowSellItemTrdReqs:
		infoLines = infoLines.
			appendT(
				"Buy item to fulfill selected request( ↵ )",
				"Place order to sell item(R)",
				GobackCommand)
	case ShowBuyChrTrdReqs:
		infoLines = infoLines.
			appendT(
				"Sell character to fulfill selected request( ↵ )",
				"Place order to buy character(R)",
				GobackCommand)
	case SelectFitBuyChrTrdReq:
		infoLines = infoLines.
			appendT(
				"Sell character to fulfill selected request( ↵ )",
				"Place order to buy character(R)",
				GobackCommand)
	case ShowSellChrTrdReqs:
		infoLines = infoLines.
			appendT(
				"Buy character to fulfill selected request( ↵ )",
				"Place order to sell character(R)",
				GobackCommand)

	case CreateBuyChrTrdReqSelectChr,
		CreateSellChrTrdReqSelChr,
		CreateSellItemTrdReqSelectItem,
		CreateBuyItemTrdReqSelectItem:
		infoLines = infoLines.appendSelectGoBackCmds()
	case SelectRenameChr:
		infoLines = infoLines.
			appendCustomFontSelectCmdsScreenCharacters(screen).
			appendSelectGoBackCmds()
	case SelectActiveChr:
		infoLines = infoLines.
			append(fmt.Sprintf("0) %s", loud.Localize("No character selection"))).
			appendCustomFontSelectCmdsScreenCharacters(screen).
			appendSelectGoBackCmds()
	case SelectBuyItem:
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
	case SelectBuyChr:
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
	case SelectSellItem:
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
	case SelectUpgradeItem:
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
	case ConfirmHuntRabbits,
		ConfirmFightGoblin,
		ConfirmFightTroll,
		ConfirmFightWolf,
		ConfirmFightGiant,
		ConfirmFightDragonFire,
		ConfirmFightDragonIce,
		ConfirmFightDragonAcid,
		ConfirmFightDragonUndead:
		infoLines = infoLines.
			appendGoOnBackCmds()
	default:
		if screen.IsHelpScreen() {
			infoLines = infoLines.
				appendT(GobackCommand)
		} else if screen.IsResultScreen() { // eg. RsltBuyGoldTrdReqCreation
			infoLines = infoLines.appendT(GoonEnterCommand)
		} else if screen.InputActive() { // eg. CreateBuyItmTrdReqEnterPylonValue
			infoLines = infoLines.
				appendT(
					FinishEnterCommand,
					GobackEscCommand)
		}
	}

	for index, line := range infoLines {
		lineFont := screen.getFont(line.font)
		PrintString(fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x),
			lineFont(fillSpace(line.content, w))))
	}

	infoLen := len(infoLines)

	for index, line := range tableLines {
		PrintString(fmt.Sprintf("%s%s",
			cursor.MoveTo(y+infoLen+index, x),
			line))
	}
	totalLen := infoLen + len(tableLines)

	screen.drawFill(x, y+totalLen, w, h-totalLen-1)
}
