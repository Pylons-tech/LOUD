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
	REFRESH_CMD      = "Re)fresh Status"
)

func (tl TextLines) appendDeselectCmd() TextLines {
	return append(tl, fmt.Sprintf("0) %s", loud.Localize("Deselect")))
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
		tl = append(tl, fmt.Sprintf("%d) %s  ", idx+1, fn(item)))
	}
	return tl
}

func (tl TextLines) appendRefreshCmd(screen *GameScreen) TextLines {
	refreshCmdTxt := loud.Localize(REFRESH_CMD)
	if screen.syncingData {
		tl = append(tl, screen.blueBoldFont()(refreshCmdTxt))
	} else {
		tl = append(tl, refreshCmdTxt)
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

	infoLines := TextLines{}
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
		infoLines = strings.Split(cmdString, "\n")
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
			append(screen.tradeTableColorDesc()...).
			appendT(
				"Sell loud to fulfill selected request( ↵ )",
				"place order to buy loud(R)",
				GO_BACK_CMD)
	case SHW_LOUD_SELL_TRDREQS:
		infoLines = infoLines.
			append(screen.tradeTableColorDesc()...).
			appendT(
				"Buy loud to fulfill selected request( ↵ )",
				"place order to sell loud(R)",
				GO_BACK_CMD)
	case SHW_BUYITM_TRDREQS:
		infoLines = infoLines.
			append(screen.tradeTableColorDesc()...).
			appendT(
				"Sell item to fulfill selected request( ↵ )",
				"place order to buy item(R)",
				GO_BACK_CMD)
	case SHW_SELLITM_TRDREQS:
		infoLines = infoLines.
			append(screen.tradeTableColorDesc()...).
			appendT(
				"Buy item to fulfill selected request( ↵ )",
				"place order to sell item(R)",
				GO_BACK_CMD)
	case SHW_BUYCHR_TRDREQS:
		infoLines = infoLines.
			append(screen.tradeTableColorDesc()...).
			appendT(
				"Sell character to fulfill selected request( ↵ )",
				"place order to buy character(R)",
				GO_BACK_CMD)
	case SHW_SELLCHR_TRDREQS:
		infoLines = infoLines.
			append(screen.tradeTableColorDesc()...).
			appendT(
				"Buy character to fulfill selected request( ↵ )",
				"place order to sell character(R)",
				GO_BACK_CMD)

	case CR8_BUYCHR_TRDREQ_SEL_CHR,
		CR8_SELLCHR_TRDREQ_SEL_CHR,
		CR8_SELLITM_TRDREQ_SEL_ITEM,
		CR8_BUYITM_TRDREQ_SEL_ITEM:
		infoLines = infoLines.
			appendT(
				SEL_CMD,
				GO_BACK_CMD)
	case SEL_HEALTH_RESTORE_CHAR,
		SEL_RENAME_CHAR:
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
			appendSelectCmds(
				loud.ShopItems,
				func(it interface{}) string {
					item := it.(loud.Item)
					return formatItem(item) + screen.goldIcon() + fmt.Sprintf(" %d", item.Price)
				}).
			appendSelectGoBackCmds()
	case SEL_BUYCHR:
		infoLines = infoLines.
			appendSelectCmds(
				loud.ShopCharacters,
				func(it interface{}) string {
					char := it.(loud.Character)
					return formatCharacter(char) + screen.pylonIcon() + fmt.Sprintf(" %d", char.Price)
				}).
			appendSelectGoBackCmds()
	case SEL_SELLITM:
		infoLines = infoLines.
			appendSelectCmds(
				screen.user.InventorySellableItems(),
				func(it interface{}) string {
					item := it.(loud.Item)
					return formatItem(item) + screen.goldIcon() + fmt.Sprintf(" %s", item.GetSellPriceRange())
				}).
			appendSelectGoBackCmds()
	case SEL_UPGITM:
		infoLines = infoLines.
			appendSelectCmds(
				screen.user.InventoryUpgradableItems(),
				func(it interface{}) string {
					item := it.(loud.Item)
					return formatItem(item) + screen.goldIcon() + fmt.Sprintf(" %d", item.GetUpgradePrice())
				}).
			appendSelectGoBackCmds()
	case CONFIRM_HUNT_RABBITS,
		CONFIRM_FIGHT_GOBLIN,
		CONFIRM_FIGHT_TROLL,
		CONFIRM_FIGHT_WOLF,
		CONFIRM_FIGHT_GIANT:
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

	infoLines = infoLines.append("\n")
	infoLines = infoLines.appendRefreshCmd(screen)
	infoLines = infoLines.appendEndGameCmd(screen)

	// cmd box start point (x, y)
	x := 2
	y := screen.cmdInnerStartY()
	w := screen.leftInnerWidth()
	h := screen.cmdInnerHeight()

	bgcolor := uint64(bgcolor)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))
	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x), fmtFunc(fillRightWithSpace(line, w))))
		if index+2 > int(screen.Height()) {
			break
		}
	}

	screen.drawFill(x, y+len(infoLines), w, h-len(infoLines))
}
