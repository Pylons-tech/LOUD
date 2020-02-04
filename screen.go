package loud

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Pylons-tech/pylons/x/pylons/handlers"
	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderUserCommands() {

	infoLines := []string{}
	switch screen.scrStatus {
	case SHOW_LOCATION:
		cmdMap := map[UserLocation]string{
			HOME:     localize("home"),
			FOREST:   localize("forest"),
			SHOP:     localize("shop"),
			MARKET:   localize("market"),
			SETTINGS: localize("settings"),
		}
		cmdString := cmdMap[screen.user.GetLocation()]
		infoLines = strings.Split(cmdString, "\n")
	case SHOW_LOUD_BUY_ORDERS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)

		infoLines = append(infoLines, "B)uy( ↵ )")
		infoLines = append(infoLines, "Create a buy o)rder")
		infoLines = append(infoLines, "Go bac)k")
	case SHOW_LOUD_SELL_ORDERS:
		infoLines = append(infoLines, screen.tradeTableColorDesc()...)

		infoLines = append(infoLines, "Se)ll( ↵ )")
		infoLines = append(infoLines, "Create sell o)rder")
		infoLines = append(infoLines, "Go bac)k")
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
		infoLines = append(infoLines, localize("Get I)nitial Coin"))
		infoLines = append(infoLines, localize("C)ancel"))
	case SELECT_UPGRADE_ITEM:
		userUpgradeItems := screen.user.UpgradableItems()
		for idx, item := range userUpgradeItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d ", idx+1, localize(item.Name), item.Level)+screen.loudIcon()+fmt.Sprintf(" %d", item.GetUpgradePrice()))
		}
		infoLines = append(infoLines, localize("C)ancel"))
	case RESULT_BUY_LOUD_ORDER_CREATION:
		infoLines = append(infoLines, localize("Go) on"))
	case RESULT_SELL_LOUD_ORDER_CREATION:
		infoLines = append(infoLines, localize("Go) on"))
	case RESULT_FULFILL_BUY_LOUD_ORDER:
		infoLines = append(infoLines, localize("Go) on"))
	case RESULT_FULFILL_SELL_LOUD_ORDER:
		infoLines = append(infoLines, localize("Go) on"))
	case RESULT_BUY_FINISH:
		fallthrough
	case RESULT_HUNT_FINISH:
		fallthrough
	case RESULT_SELL_FINISH:
		fallthrough
	case RESULT_UPGRADE_FINISH:
		infoLines = append(infoLines, localize("Go) on"))
	default:
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

func (screen *GameScreen) renderUserSituation() {
	infoLines := []string{}
	desc := ""
	waitProcessEnd := localize("wait process to end")
	switch screen.scrStatus {
	case SHOW_LOCATION:
		locationDescMap := map[UserLocation]string{
			HOME:     localize("home desc"),
			FOREST:   localize("forest desc"),
			SHOP:     localize("shop desc"),
			MARKET:   localize("market desc"),
			SETTINGS: localize("settings desc"),
		}
		desc = locationDescMap[screen.user.GetLocation()]
	case CREATE_BUY_LOUD_ORDER_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add localize
	case CREATE_SELL_LOUD_ORDER_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to get (should be integer value)" // TODO should add localize
	case CREATE_BUY_LOUD_ORDER_ENTER_LOUD_VALUE:
		desc = "Please enter loud amount to buy (should be integer value)" // TODO should add localize
	case CREATE_SELL_LOUD_ORDER_ENTER_LOUD_VALUE:
		desc = "Please enter loud amount to sell (should be integer value)" // TODO should add localize

	case SHOW_LOUD_BUY_ORDERS:
		infoLines = screen.renderOrderTable(buyOrders)
	case SHOW_LOUD_SELL_ORDERS:
		infoLines = screen.renderOrderTable(sellOrders)
	case SELECT_BUY_ITEM:
		desc = localize("select buy item desc")
	case SELECT_SELL_ITEM:
		desc = localize("select sell item desc")
	case SELECT_HUNT_ITEM:
		desc = localize("select hunt item desc")
	case SELECT_UPGRADE_ITEM:
		desc = localize("select upgrade item desc")
	case WAIT_FULFILL_BUY_LOUD_ORDER:
		order := screen.activeOrder
		desc = localize("you are now buying loud from pylon") + fmt.Sprintf(" at %.4f.\n", order.Price)
		desc += screen.buyLoudDesc(order.Amount, order.Total)
	case WAIT_FULFILL_SELL_LOUD_ORDER:
		order := screen.activeOrder
		desc = localize("you are now selling loud for pylon") + fmt.Sprintf(" at %.4f.\n", order.Price)
		desc += screen.sellLoudDesc(order.Amount, order.Total)
	case WAIT_BUY_LOUD_ORDER_CREATION:
		desc = localize("you are now waiting for loud buy order creation")
		desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case WAIT_SELL_LOUD_ORDER_CREATION:
		desc = localize("you are now waiting for loud sell order creation")
		desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case WAIT_BUY_PROCESS:
		desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("wait buy process desc"), localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
	case WAIT_HUNT_PROCESS:
		if len(screen.activeItem.Name) > 0 {
			desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("wait hunt process desc"), localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
		} else {
			switch string(screen.lastInput.Ch) {
			case "I":
				fallthrough
			case "i":
				desc = fmt.Sprintf("%s\n%s", localize("Getting initial gold from pylon"), waitProcessEnd)
			default:
				desc = fmt.Sprintf("%s\n%s", localize("hunting without weapon"), waitProcessEnd)
			}
		}
	case WAIT_SELL_PROCESS:
		desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("wait sell process desc"), localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
	case WAIT_UPGRADE_PROCESS:
		desc = fmt.Sprintf("%s %s.\n%s", localize("wait upgrade process desc"), localize(screen.activeItem.Name), waitProcessEnd)
	case RESULT_BUY_LOUD_ORDER_CREATION:
		if screen.txFailReason != "" {
			desc = localize("loud buy order creation fail reason") + ": " + localize(screen.txFailReason)
		} else {
			desc = localize("loud buy order was successfully created")
			desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		}
	case RESULT_SELL_LOUD_ORDER_CREATION:
		if screen.txFailReason != "" {
			desc = localize("sell buy order creation fail reason") + ": " + localize(screen.txFailReason)
		} else {
			desc = localize("loud sell order was successfully created")
			desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		}
	case RESULT_FULFILL_BUY_LOUD_ORDER:
		if screen.txFailReason != "" {
			desc = localize("buy loud failed reason") + ": " + localize(screen.txFailReason)
		} else {
			order := screen.activeOrder
			desc = localize("you have bought loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", order.Price)
			desc += screen.buyLoudDesc(order.Amount, order.Total)
		}
	case RESULT_FULFILL_SELL_LOUD_ORDER:
		if screen.txFailReason != "" {
			desc = localize("sell loud failed reason") + ": " + localize(screen.txFailReason)
		} else {
			order := screen.activeOrder
			desc = localize("you have sold loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", order.Price)
			desc += screen.sellLoudDesc(order.Amount, order.Total)
		}
	case RESULT_BUY_FINISH:
		if screen.txFailReason != "" {
			desc = localize("buy failed reason") + ": " + localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("result buy finish desc"), localize(screen.activeItem.Name), screen.activeItem.Level, localize("use for hunting"))
		}
	case RESULT_HUNT_FINISH:
		if screen.txFailReason != "" {
			desc = localize("hunt failed reason") + ": " + localize(screen.txFailReason)
		} else {
			respOutput := handlers.ExecuteRecipeSerialize{}
			json.Unmarshal(screen.txResult, &respOutput)
			switch string(screen.lastInput.Ch) {
			case "I":
				fallthrough
			case "i":
				desc = fmt.Sprintf("%s %d.", localize("Got initial gold from pylons. Amount is"), respOutput.Amount)
			default:
				desc = fmt.Sprintf("%s %d.", localize("result hunt finish desc"), respOutput.Amount)
			}
		}
	case RESULT_SELL_FINISH:
		if screen.txFailReason != "" {
			desc = localize("sell failed reason") + ": " + localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("%s %s Lv%d.", localize("result sell finish desc"), localize(screen.activeItem.Name), screen.activeItem.Level)
		}
	case RESULT_UPGRADE_FINISH:
		if screen.txFailReason != "" {
			desc = localize("upgrade failed reason") + ": " + localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("%s: %s.", localize("result upgrade finish desc"), localize(screen.activeItem.Name))
		}
	}

	basicLines := strings.Split(desc, "\n")

	for _, line := range basicLines {
		infoLines = append(infoLines, ChunkString(line, screen.screenSize.Width/2-4)...)
	}

	// box start point (x, y)
	x := 2
	y := 2

	bgcolor := uint64(bgcolor)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))
	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(y+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}
}

func (screen *GameScreen) InputActive() bool {
	switch screen.scrStatus {
	case CREATE_BUY_LOUD_ORDER_ENTER_LOUD_VALUE:
		return true
	case CREATE_BUY_LOUD_ORDER_ENTER_PYLON_VALUE:
		return true
	case CREATE_SELL_LOUD_ORDER_ENTER_LOUD_VALUE:
		return true
	case CREATE_SELL_LOUD_ORDER_ENTER_PYLON_VALUE:
		return true
	}
	return false
}

func (screen *GameScreen) renderInputValue() {
	inputWidth := uint32(screen.screenSize.Width/2) - 2
	move := cursor.MoveTo(screen.screenSize.Height-1, 2)

	fmtString := fmt.Sprintf("%%-%vs", inputWidth-7)

	chatFunc := screen.colorFunc(fmt.Sprintf("231:%v", bgcolor))
	chat := chatFunc("VALUE▶ ")

	if screen.InputActive() {
		chatFunc = screen.colorFunc(fmt.Sprintf("0+b:%v", bgcolor-1))
	}

	fixedChat := truncateLeft(screen.inputText, int(inputWidth-7))

	inputText := fmt.Sprintf("%s%s%s", move, chat, chatFunc(fmt.Sprintf(fmtString, fixedChat)))

	io.WriteString(os.Stdout, inputText)
}

func (screen *GameScreen) renderCharacterSheet() {
	var HP uint64 = 10
	var MaxHP uint64 = 10
	bgcolor := uint64(bgcolor)
	warning := ""
	if float32(HP) < float32(MaxHP)*.25 {
		bgcolor = 124
		warning = localize("health low warning")
	} else if float32(HP) < float32(MaxHP)*.1 {
		bgcolor = 160
		warning = localize("health critical warning")
	}

	x := screen.screenSize.Width/2 - 1
	width := (screen.screenSize.Width - x)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))

	infoLines := []string{
		centerText(fmt.Sprintf("%v", screen.user.GetUserName()), " ", width),
		centerText(warning, "─", width),
		screen.pylonIcon() + fmtFunc(truncateRight(fmt.Sprintf(" %s: %v", "Pylon", screen.user.GetPylonAmount()), width-1)),
		screen.loudIcon() + fmtFunc(truncateRight(fmt.Sprintf(" %s: %v", localize("gold"), screen.user.GetGold()), width-1)),
		screen.drawProgressMeter(HP, MaxHP, 196, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" HP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 225, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" XP: %v/%v", HP, 10), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 208, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" AP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 117, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" RP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 76, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" MP: %v/%v", HP, MaxHP), width-10)),
	}

	infoLines = append(infoLines, centerText(localize("inventory items"), "─", width))
	items := screen.user.InventoryItems()
	for _, item := range items {
		infoLines = append(infoLines, truncateRight(fmt.Sprintf("%s Lv%d", localize(item.Name), item.Level), width))
	}
	infoLines = append(infoLines, centerText(" ❦ ", "─", width))

	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(2+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}

	nodeLines := []string{
		centerText(localize("pylons network status"), " ", width),
		centerText(screen.user.GetLastTransaction(), " ", width),
	}

	blockHeightText := centerText(localize("block height")+": "+strconv.FormatInt(screen.blockHeight, 10), " ", width)
	if screen.refreshingDaemonStatus {
		nodeLines = append(nodeLines, screen.blueBoldFont()(blockHeightText))
	} else {
		nodeLines = append(nodeLines, blockHeightText)
	}
	nodeLines = append(nodeLines, centerText(" ❦ ", "─", width))

	for index, line := range nodeLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(2+len(infoLines)+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}

	lastLine := len(infoLines) + len(nodeLines) + 1
	screen.drawFill(x, lastLine+1, width, screen.screenSize.Height-(lastLine+2))
}

func (screen *GameScreen) RunSelectedLoudBuyTrade() {
	if len(buyOrders) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = localize("you haven't selected any buy order")
		screen.scrStatus = RESULT_FULFILL_BUY_LOUD_ORDER
		screen.refreshed = false
		screen.Render()
	} else {
		screen.scrStatus = WAIT_FULFILL_BUY_LOUD_ORDER
		screen.activeOrder = buyOrders[screen.activeLine]
		screen.refreshed = false
		screen.Render()
		txhash, err := FulfillTrade(screen.user, buyOrders[screen.activeLine].ID)

		log.Println("ended sending request for creating buy loud order")
		if err != nil {
			screen.txFailReason = err.Error()
			screen.scrStatus = RESULT_FULFILL_BUY_LOUD_ORDER
			screen.refreshed = false
			screen.Render()
		} else {
			time.AfterFunc(2*time.Second, func() {
				screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
				screen.scrStatus = RESULT_FULFILL_BUY_LOUD_ORDER
				screen.refreshed = false
				screen.Render()
			})
		}
	}
}

func (screen *GameScreen) RunSelectedLoudSellTrade() {
	if len(sellOrders) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = localize("you haven't selected any sell order")
		screen.scrStatus = RESULT_FULFILL_SELL_LOUD_ORDER
		screen.refreshed = false
		screen.Render()
	} else {
		screen.scrStatus = WAIT_FULFILL_SELL_LOUD_ORDER
		screen.activeOrder = sellOrders[screen.activeLine]
		screen.refreshed = false
		screen.Render()
		txhash, err := FulfillTrade(screen.user, sellOrders[screen.activeLine].ID)

		log.Println("ended sending request for creating sell loud order")
		if err != nil {
			screen.txFailReason = err.Error()
			screen.scrStatus = RESULT_FULFILL_SELL_LOUD_ORDER
			screen.refreshed = false
			screen.Render()
		} else {
			time.AfterFunc(2*time.Second, func() {
				screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
				screen.scrStatus = RESULT_FULFILL_SELL_LOUD_ORDER
				screen.refreshed = false
				screen.Render()
			})
		}
	}
}

func (screen *GameScreen) Render() {
	var HP uint64 = 10

	if screen.screenSize.Height < 20 || screen.screenSize.Width < 60 {
		clear := cursor.ClearEntireScreen()
		move := cursor.MoveTo(1, 1)
		io.WriteString(os.Stdout,
			fmt.Sprintf("%s%s%s", clear, move, localize("screen size warning")))
		return
	} else if HP == 0 {
		clear := cursor.ClearEntireScreen()
		dead := localize("dead")
		move := cursor.MoveTo(screen.screenSize.Height/2, screen.screenSize.Width/2-utf8.RuneCountInString(dead)/2)
		io.WriteString(os.Stdout, clear+move+dead)
		screen.refreshed = false
		return
	}

	if !screen.refreshed {
		clear := cursor.ClearEntireScreen() + allowMouseInputAndHideCursor
		io.WriteString(os.Stdout, clear)
		screen.redrawBorders()
		screen.refreshed = true
	}

	screen.renderUserCommands()
	screen.renderUserSituation()
	screen.renderCharacterSheet()
	screen.renderInputValue()
}
