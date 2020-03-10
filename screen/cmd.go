package screen

import (
	"fmt"
	"io"
	"os"
	"strings"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderUserCommands() {

	infoLines := []string{}
	switch screen.scrStatus {
	case SHOW_LOCATION:
		cmdMap := map[loud.UserLocation]string{
			loud.HOME:     "home",
			loud.FOREST:   "forest",
			loud.SHOP:     "shop",
			loud.MARKET:   "market",
			loud.SETTINGS: "settings",
			loud.DEVELOP:  "develop",
		}
		cmdString := loud.Localize(cmdMap[screen.user.GetLocation()])
		infoLines = strings.Split(cmdString, "\n")
		for _, loc := range []loud.UserLocation{loud.HOME, loud.FOREST, loud.SHOP, loud.MARKET, loud.SETTINGS, loud.DEVELOP} {
			if loc != screen.user.GetLocation() {
				infoLines = append(infoLines, loud.Localize("go to "+cmdMap[loc]))
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
		for idx, item := range loud.ShopItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d  ", idx+1, loud.Localize(item.Name), item.Level)+screen.loudIcon()+fmt.Sprintf(" %d", item.Price))
		}
		infoLines = append(infoLines,
			"Select ( ↵ )",
			loud.Localize("C)ancel"))
	case SELECT_BUY_CHARACTER:
		for idx, item := range loud.ShopCharacters {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d  ", idx+1, loud.Localize(item.Name), item.Level)+screen.pylonIcon()+fmt.Sprintf(" %d", item.Price))
		}
		infoLines = append(infoLines,
			"Select ( ↵ )",
			loud.Localize("C)ancel"))
	case SELECT_SELL_ITEM:
		userItems := screen.user.InventoryItems()
		for idx, item := range userItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d  ", idx+1, loud.Localize(item.Name), item.Level)+screen.loudIcon()+fmt.Sprintf(" %d", item.GetSellPrice()))
		}
		infoLines = append(infoLines,
			"Select ( ↵ )",
			loud.Localize("C)ancel"))
	case SELECT_HUNT_ITEM:
		infoLines = append(infoLines, loud.Localize("No item"))
		infoLines = append(infoLines, loud.Localize("Get I)nitial Coin"))
		for idx, item := range screen.user.InventoryItems() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d", idx+1, loud.Localize(item.Name), item.Level))
		}
		infoLines = append(infoLines,
			"Select ( ↵ )",
			loud.Localize("C)ancel"))
	case SELECT_UPGRADE_ITEM:
		for idx, item := range screen.user.UpgradableItems() {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d ", idx+1, loud.Localize(item.Name), item.Level)+screen.loudIcon()+fmt.Sprintf(" %d", item.GetUpgradePrice()))
		}
		infoLines = append(infoLines,
			"Select ( ↵ )",
			loud.Localize("C)ancel"))
	default:
		if screen.IsResultScreen() { // eg. RESULT_BUY_LOUD_REQUEST_CREATION
			infoLines = append(infoLines, loud.Localize("Go) on( ↵ )"))
		} else if screen.InputActive() { // eg. CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE
			infoLines = append(infoLines,
				loud.Localize("Finish Enter ( ↵ )"),
				loud.Localize("Go bac)k( ⌫ )"))
		}
	}

	infoLines = append(infoLines, "\n")
	if screen.syncingData {
		infoLines = append(infoLines, screen.blueBoldFont()(loud.Localize("Re)fresh Status")))
	} else {
		infoLines = append(infoLines, loud.Localize("Re)fresh Status"))
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
