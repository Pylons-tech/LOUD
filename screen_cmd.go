package loud

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderUserCommands() {

	infoLines := []string{}
	switch screen.scrStatus {
	case SHOW_LOCATION:
		cmdMap := map[UserLocation]string{
			HOME:     "home",
			FOREST:   "forest",
			SHOP:     "shop",
			MARKET:   "market",
			SETTINGS: "settings",
			DEVELOP:  "develop",
		}
		cmdString := localize(cmdMap[screen.user.GetLocation()])
		infoLines = strings.Split(cmdString, "\n")
		for _, loc := range []UserLocation{HOME, FOREST, SHOP, MARKET, SETTINGS, DEVELOP} {
			if loc != screen.user.GetLocation() {
				infoLines = append(infoLines, localize("go to "+cmdMap[loc]))
			}
		}
	case SHOW_LOUD_BUY_REQUESTS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Sell loud to fulfill selected request( ↵ )",
			"Create an order to buy loud(R)",
			"Go bac)k( ⌫ )")
	case SHOW_LOUD_SELL_REQUESTS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Buy loud to fulfill selected request( ↵ )",
			"Create an order to sell loud(R)",
			"Go bac)k( ⌫ )")
	case SHOW_BUY_SWORD_REQUESTS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Sell item to fulfill selected request( ↵ )",
			"Create an order to buy item(R)",
			"Go bac)k( ⌫ )")
	case SHOW_SELL_SWORD_REQUESTS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)
		infoLines = append(infoLines,
			"Buy item to fulfill selected request( ↵ )",
			"Create an order to sell item(R)",
			"Go bac)k( ⌫ )")

	case CREATE_SELL_SWORD_REQUEST_SELECT_SWORD:
		fallthrough
	case CREATE_BUY_SWORD_REQUEST_SELECT_SWORD:
		infoLines = append(infoLines,
			"Select ( ↵ )",
			"Go bac)k( ⌫ )")
	case SELECT_BUY_ITEM:
		for idx, item := range shopItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d  ", idx+1, localize(item.Name), item.Level)+screen.loudIcon()+fmt.Sprintf(" %d", item.Price))
		}
		infoLines = append(infoLines, localize("C)ancel"))
	case SELECT_SELL_ITEM:
		userItems := screen.user.InventoryItems()
		for idx, item := range userItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d  ", idx+1, localize(item.Name), item.Level)+screen.loudIcon()+fmt.Sprintf(" %d", item.GetSellPrice()))
		}
		infoLines = append(infoLines, localize("C)ancel"))
	case SELECT_HUNT_ITEM:
		userWeaponItems := screen.user.InventoryItems()
		infoLines = append(infoLines, localize("N)o item"))
		for idx, item := range userWeaponItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d", idx+1, localize(item.Name), item.Level))
		}
		infoLines = append(infoLines,
			localize("C)ancel"))
	case SELECT_UPGRADE_ITEM:
		userUpgradeItems := screen.user.UpgradableItems()
		for idx, item := range userUpgradeItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d ", idx+1, localize(item.Name), item.Level)+screen.loudIcon()+fmt.Sprintf(" %d", item.GetUpgradePrice()))
		}
		infoLines = append(infoLines, localize("C)ancel"))
	default:
		if screen.IsResultScreen() { // eg. RESULT_BUY_LOUD_REQUEST_CREATION
			infoLines = append(infoLines, localize("Go) on( ↵ )"))
		} else if screen.InputActive() { // eg. CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE
			infoLines = append(infoLines,
				localize("Finish Enter ( ↵ )"),
				localize("Go bac)k( ⌫ )"))
		}
	}

	infoLines = append(infoLines, "\n")
	if screen.syncingData {
		infoLines = append(infoLines, screen.blueBoldFont()(localize("Re)fresh Status")))
	} else {
		infoLines = append(infoLines, localize("Re)fresh Status"))
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
