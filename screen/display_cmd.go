package screen

import (
	"fmt"
	"io"
	"os"
	"strings"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/ahmetb/go-cursor"
)

const (
	SEL_CMD          = "Select ( ↵ )"
	GO_ON_ENTER_CMD  = "Go on ( ↵ )"
	FINISH_ENTER_CMD = "Finish Enter ( ↵ )"
	GO_BACK_CMD      = "Go back ( ⌫ ) - Backspace Key"
	GO_BACK_ESC_CMD  = "Go back ( Esc )"
	END_GAME_ESC_CMD = "End Game ( Esc )"
)

func (tl TextLines) appendDeselectCmd() TextLines {
	return tl.append(fmt.Sprintf("0) %s", loud.Localize("Deselect")))
}

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

func (tl TextLines) appendSelectCmds(itemsSlice interface{}, fn func(interface{}) string) TextLines {
	items := InterfaceSlice(itemsSlice)
	for idx, item := range items {
		tl = tl.append(fmt.Sprintf("%d) %s  ", idx+1, fn(item)))
	}
	return tl
}

func (tl TextLines) appendCustomFontSelectCmds(itemsSlice interface{}, fn func(interface{}) TextLine) TextLines {
	items := InterfaceSlice(itemsSlice)
	for idx, item := range items {
		fni := fn(item)
		tl = append(tl, TextLine{
			content: fmt.Sprintf("%d) %s  ", idx+1, fni.content),
			font:    fni.font,
		})
	}
	return tl
}

func (tl TextLines) appendEndGameCmd(screen *GameScreen) TextLines {
	if screen.scrStatus != CONFIRM_ENDGAME && !screen.InputActive() {
		tl = tl.appendT(END_GAME_ESC_CMD)
	}
	return tl
}

func (screen *GameScreen) renderUserCommands() {
	// cmd box start point (x, y)
	x := 2
	y := screen.cmdInnerStartY()
	w := screen.leftInnerWidth()
	h := screen.cmdInnerHeight()

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
		}
		cmdString := loud.Localize(cmdMap[screen.user.GetLocation()])
		infoLines = infoLines.
			append(strings.Split(cmdString, "\n")...)

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

		for _, loc := range []loud.UserLocation{
			loud.HOME,
			loud.FOREST,
			loud.SHOP,
			loud.PYLCNTRL,
			loud.SETTINGS,
			loud.DEVELOP,
		} {
			if loc != screen.user.GetLocation() {
				infoLines = infoLines.
					append(loud.Localize("go to " + cmdMap[loc]))
			}
		}
	case SHW_LOUD_BUY_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Sell gold to fulfill selected request( ↵ )",
				"Place order to buy gold(R)",
				GO_BACK_CMD)
		tableLines = screen.tradeTableColorDesc(w)
	case SHW_LOUD_SELL_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Buy gold to fulfill selected request( ↵ )",
				"place order to sell gold(R)",
				GO_BACK_CMD)
		tableLines = screen.tradeTableColorDesc(w)
	case SHW_BUYITM_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Sell item to fulfill selected request( ↵ )",
				"Place order to buy item(R)",
				GO_BACK_CMD)
		tableLines = screen.tradeTableColorDesc(w)
	case SHW_SELLITM_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Buy item to fulfill selected request( ↵ )",
				"Place order to sell item(R)",
				GO_BACK_CMD)
		tableLines = screen.tradeTableColorDesc(w)
	case SHW_BUYCHR_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Sell character to fulfill selected request( ↵ )",
				"Place order to buy character(R)",
				GO_BACK_CMD)
		tableLines = screen.tradeTableColorDesc(w)
	case SHW_SELLCHR_TRDREQS:
		infoLines = infoLines.
			appendT(
				"Buy character to fulfill selected request( ↵ )",
				"Place order to sell character(R)",
				GO_BACK_CMD)
		tableLines = screen.tradeTableColorDesc(w)

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
			appendSelectCmds(
				screen.user.InventoryCharacters(),
				func(it interface{}) string {
					return formatCharacter(it.(loud.Character))
				}).
			appendSelectGoBackCmds()
	case SEL_ACTIVE_CHAR:
		infoLines = infoLines.
			appendDeselectCmd().
			appendSelectCmds(
				screen.user.InventoryCharacters(),
				func(it interface{}) string {
					return formatCharacter(it.(loud.Character))
				}).
			appendSelectGoBackCmds()
	case SEL_ACTIVE_WEAPON:
		infoLines = infoLines.
			appendDeselectCmd().
			appendSelectCmds(
				screen.user.InventorySwords(),
				func(it interface{}) string {
					return formatItem(it.(loud.Item))
				}).
			appendSelectGoBackCmds()
	case SEL_BUYITM:
		infoLines = infoLines.
			appendCustomFontSelectCmds(
				loud.ShopItems,
				func(it interface{}) TextLine {
					item := it.(loud.Item)
					preitemOk := screen.user.HasPreItemForAnItem(item)
					goldEnough := item.Price <= screen.user.GetGold()
					font := REGULAR
					bonusText := ""
					contentStr := ""
					if !preitemOk {
						font = GREY
						bonusText = fmt.Sprintf(": %s", loud.Localize("no material"))
					}
					if !goldEnough {
						font = GREY
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
			appendSelectCmds(
				loud.ShopCharacters,
				func(it interface{}) string {
					char := it.(loud.Character)
					return fmt.Sprintf("%s  %s%d", formatCharacter(char), screen.pylonIcon(), char.Price)
				}).
			appendSelectGoBackCmds()
	case SEL_SELLITM:
		infoLines = infoLines.
			appendSelectCmds(
				screen.user.InventorySellableItems(),
				func(it interface{}) string {
					item := it.(loud.Item)
					return formatItem(item) + fmt.Sprintf("💰 %s", item.GetSellPriceRange())
				}).
			appendSelectGoBackCmds()
	case SEL_UPGITM:
		infoLines = infoLines.
			appendSelectCmds(
				screen.user.InventoryUpgradableItems(),
				func(it interface{}) string {
					item := it.(loud.Item)
					return formatItem(item) + fmt.Sprintf("💰 %d", item.GetUpgradePrice())
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
		if screen.IsResultScreen() { // eg. RSLT_BUY_LOUD_TRDREQ_CREATION
			infoLines = infoLines.appendT(GO_ON_ENTER_CMD)
		} else if screen.InputActive() { // eg. CR8_BUYITM_TRDREQ_ENT_PYLVAL
			infoLines = infoLines.
				appendT(
					FINISH_ENTER_CMD,
					GO_BACK_ESC_CMD)
		}
	}

	infoLines = infoLines.append("") // same as enter command
	infoLines = infoLines.appendEndGameCmd(screen)

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

	screen.drawFill(x, y+totalLen, w, h-totalLen)
}
