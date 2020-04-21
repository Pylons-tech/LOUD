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
	SEL_CMD     = "Select ( ↵ )"
	GO_BACK_CMD = "Go back( ⌫ ) - Backspace Key"
)

func appendSelectGoBackCmds(infoLines []string) []string {
	return append(infoLines,
		SEL_CMD,
		loud.Localize(GO_BACK_CMD))
}

func (screen *GameScreen) renderUserCommands() {

	infoLines := []string{}
	switch screen.scrStatus {
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
		for _, loc := range []loud.UserLocation{loud.HOME, loud.FOREST, loud.SHOP, loud.PYLCNTRL, loud.SETTINGS, loud.DEVELOP} {
			if loc != screen.user.GetLocation() {
				infoLines = append(infoLines, loud.Localize("go to "+cmdMap[loc]))
			}
		}
	case SHW_LOUD_BUY_TRDREQS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Sell loud to fulfill selected request( ↵ )",
			"Create an order to buy loud(R)",
			GO_BACK_CMD)
	case SHW_LOUD_SELL_TRDREQS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Buy loud to fulfill selected request( ↵ )",
			"Create an order to sell loud(R)",
			GO_BACK_CMD)
	case SHW_BUYITM_TRDREQS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Sell item to fulfill selected request( ↵ )",
			"Create an order to buy item(R)",
			GO_BACK_CMD)
	case SHW_SELLITM_TRDREQS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Buy item to fulfill selected request( ↵ )",
			"Create an order to sell item(R)",
			GO_BACK_CMD)
	case SHW_BUYCHR_TRDREQS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Sell character to fulfill selected request( ↵ )",
			"Create an order to buy character(R)",
			GO_BACK_CMD)
	case SHW_SELLCHR_TRDREQS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Buy character to fulfill selected request( ↵ )",
			"Create an order to sell character(R)",
			GO_BACK_CMD)

	case CR8_BUYCHR_TRDREQ_SEL_CHR,
		CR8_SELLCHR_TRDREQ_SEL_CHR,
		CR8_SELLITM_TRDREQ_SEL_ITEM,
		CR8_BUYITM_TRDREQ_SEL_ITEM:
		infoLines = append(infoLines,
			SEL_CMD,
			GO_BACK_CMD)
	case SEL_ACTIVE_CHAR,
		SEL_HEALTH_RESTORE_CHAR,
		SEL_RENAME_CHAR:
		for idx, char := range screen.user.InventoryCharacters() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s  ", idx+1, formatCharacter(char)))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_ACTIVE_WEAPON:
		for idx, item := range screen.user.InventorySwords() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s  ", idx+1, formatItem(item)))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_BUYITM:
		for idx, item := range loud.ShopItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s  ", idx+1, formatItem(item))+screen.loudIcon()+fmt.Sprintf(" %d", item.Price))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_BUYCHR:
		for idx, item := range loud.ShopCharacters {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s  ", idx+1, formatCharacter(item))+screen.pylonIcon()+fmt.Sprintf(" %d", item.Price))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_SELLITM:
		userItems := screen.user.InventorySellableItems()
		for idx, item := range userItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s  ", idx+1, formatItem(item))+screen.loudIcon()+fmt.Sprintf(" %s", item.GetSellPriceRange()))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_HUNT_RABBITS_ITEM:
		infoLines = append(infoLines, loud.Localize("No item"))
		for idx, item := range screen.user.InventorySwords() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s", idx+1, formatItem(item)))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_FIGHT_GOBLIN_ITEM,
		SEL_FIGHT_TROLL_ITEM,
		SEL_FIGHT_WOLF_ITEM:
		for idx, item := range screen.user.InventorySwords() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s", idx+1, formatItem(item)))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_FIGHT_GIANT_ITEM:
		for idx, item := range screen.user.InventoryIronSwords() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s", idx+1, formatItem(item)))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	case SEL_UPGITM:
		for idx, item := range screen.user.InventoryUpgradableItems() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s ", idx+1, formatItem(item))+screen.loudIcon()+fmt.Sprintf(" %d", item.GetUpgradePrice()))
		}
		infoLines = appendSelectGoBackCmds(infoLines)
	default:
		if screen.IsResultScreen() { // eg. RSLT_BUY_LOUD_TRDREQ_CREATION
			infoLines = append(infoLines, loud.Localize("Go) on( ↵ )"))
		} else if screen.InputActive() { // eg. CR8_BUYITM_TRDREQ_ENT_PYLVAL
			infoLines = append(infoLines,
				loud.Localize("Finish Enter ( ↵ )"),
				loud.Localize(GO_BACK_CMD))
		}
	}

	infoLines = append(infoLines, "\n")
	refreshCmdTxt := loud.Localize("Re)fresh Status")
	if screen.syncingData {
		infoLines = append(infoLines, screen.blueBoldFont()(refreshCmdTxt))
	} else {
		infoLines = append(infoLines, refreshCmdTxt)
	}

	// box start point (x, y)
	x := 2
	y := screen.screenSize.Height/2 + 1

	bgcolor := uint64(bgcolor)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))
	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}
}
