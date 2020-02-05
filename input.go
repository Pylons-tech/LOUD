package loud

import (
	"log"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

func (screen *GameScreen) HandleInputKey(input termbox.Event) {
	screen.lastInput = input
	Key := string(input.Ch)
	log.Println("Handling Key \"", Key, "\"")
	if screen.InputActive() {
		switch input.Key {
		case termbox.KeyBackspace2:
			fallthrough
		case termbox.KeyBackspace:
			log.Println("Pressed Backspace")
			screen.SetInputTextAndRender(screen.inputText[:len(screen.inputText)-1])
			return
		case termbox.KeyEnter:
			switch screen.scrStatus {
			case CREATE_BUY_LOUD_ORDER_ENTER_LOUD_VALUE:
				screen.scrStatus = CREATE_BUY_LOUD_ORDER_ENTER_PYLON_VALUE
				screen.loudEnterValue = screen.inputText
				screen.inputText = ""
			case CREATE_BUY_LOUD_ORDER_ENTER_PYLON_VALUE:
				screen.scrStatus = WAIT_BUY_LOUD_ORDER_CREATION
				screen.pylonEnterValue = screen.inputText
				screen.SetInputTextAndRender("")
				txhash, err := CreateBuyLoudOrder(screen.user, screen.loudEnterValue, screen.pylonEnterValue)
				log.Println("ended sending request for creating buy loud order")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_BUY_LOUD_ORDER_CREATION
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(2*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_BUY_LOUD_ORDER_CREATION
						screen.refreshed = false
						screen.Render()
					})
				}
			case CREATE_SELL_LOUD_ORDER_ENTER_LOUD_VALUE:
				screen.scrStatus = CREATE_SELL_LOUD_ORDER_ENTER_PYLON_VALUE
				screen.loudEnterValue = screen.inputText
				screen.inputText = ""
			case CREATE_SELL_LOUD_ORDER_ENTER_PYLON_VALUE:
				screen.scrStatus = WAIT_SELL_LOUD_ORDER_CREATION
				screen.pylonEnterValue = screen.inputText
				screen.SetInputTextAndRender("")
				txhash, err := CreateSellLoudOrder(screen.user, screen.loudEnterValue, screen.pylonEnterValue)

				log.Println("ended sending request for creating buy loud order")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_SELL_LOUD_ORDER_CREATION
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(2*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_SELL_LOUD_ORDER_CREATION
						screen.refreshed = false
						screen.Render()
					})
				}
			}
		default:
			if _, err := strconv.Atoi(Key); err == nil {
				// If user entered number, just use it
				screen.SetInputTextAndRender(screen.inputText + Key)
			}
			return
		}
		screen.refreshed = false
	} else {
		// TODO should check current location, scrStatus and then after that check Key, rather than checking Key first
		switch input.Key {
		case termbox.KeyArrowLeft:
		case termbox.KeyArrowRight:
		case termbox.KeyArrowUp:
			if screen.activeLine > 0 {
				screen.activeLine -= 1
			}
		case termbox.KeyArrowDown:
			screen.activeLine += 1
		}
		if input.Key == termbox.KeyEnter {
			if screen.user.GetLocation() == MARKET {
				switch screen.scrStatus {
				case SHOW_LOUD_BUY_ORDERS:
					screen.RunSelectedLoudBuyTrade()
				case SHOW_LOUD_SELL_ORDERS:
					screen.RunSelectedLoudSellTrade()
				}
			}
		}

		switch Key {
		case "H": // HOME
			fallthrough
		case "h":
			screen.user.SetLocation(HOME)
			screen.refreshed = false
		case "F": // FOREST
			fallthrough
		case "f":
			screen.user.SetLocation(FOREST)
			screen.refreshed = false
		case "S": // SHOP
			fallthrough
		case "s":
			screen.user.SetLocation(SHOP)
			screen.refreshed = false
		case "M": // MARKET
			fallthrough
		case "m":
			screen.user.SetLocation(MARKET)
			screen.refreshed = false
		case "T": // SETTINGS
			fallthrough
		case "t":
			screen.user.SetLocation(SETTINGS)
			screen.refreshed = false
		case "G":
			fallthrough
		case "g":
			gameLanguage = "en"
			screen.refreshed = false
		case "A":
			fallthrough
		case "a":
			gameLanguage = "es"
			screen.refreshed = false
		case "C": // CANCEL
			fallthrough
		case "c":
			screen.scrStatus = SHOW_LOCATION
			screen.refreshed = false
		case "O": // GO ON, GO BACK, CREATE ORDER
			fallthrough
		case "o":
			if screen.user.GetLocation() == MARKET {
				switch screen.scrStatus {
				case SHOW_LOUD_BUY_ORDERS:
					screen.scrStatus = CREATE_BUY_LOUD_ORDER_ENTER_LOUD_VALUE
				case SHOW_LOUD_SELL_ORDERS:
					screen.scrStatus = CREATE_SELL_LOUD_ORDER_ENTER_LOUD_VALUE
				case RESULT_BUY_LOUD_ORDER_CREATION:
					fallthrough
				case RESULT_FULFILL_BUY_LOUD_ORDER:
					screen.scrStatus = SHOW_LOUD_BUY_ORDERS
				case RESULT_SELL_LOUD_ORDER_CREATION:
					fallthrough
				case RESULT_FULFILL_SELL_LOUD_ORDER:
					screen.scrStatus = SHOW_LOUD_SELL_ORDERS
				default:
					screen.scrStatus = SHOW_LOCATION
				}
				screen.txFailReason = ""
				screen.refreshed = false
			} else {
				screen.txFailReason = ""
				screen.scrStatus = SHOW_LOCATION
				screen.refreshed = false
			}
		case "U": // HUNT
			fallthrough
		case "u":
			screen.scrStatus = SELECT_HUNT_ITEM
			screen.refreshed = false
		case "B": // BUY
			fallthrough
		case "b": // BUY
			if screen.user.GetLocation() == SHOP {
				screen.scrStatus = SELECT_BUY_ITEM
				screen.refreshed = false
			} else if screen.user.GetLocation() == MARKET {
				if screen.scrStatus == SHOW_LOCATION {
					screen.scrStatus = SHOW_LOUD_BUY_ORDERS
					screen.refreshed = false
				} else if screen.scrStatus == SHOW_LOUD_BUY_ORDERS {
					screen.RunSelectedLoudBuyTrade()
				}
			}
		case "E": // SELL
			fallthrough
		case "e":
			if screen.user.GetLocation() == SHOP {
				screen.scrStatus = SELECT_SELL_ITEM
				screen.refreshed = false
			} else if screen.user.GetLocation() == MARKET {
				if screen.scrStatus == SHOW_LOCATION {
					screen.scrStatus = SHOW_LOUD_SELL_ORDERS
					screen.refreshed = false
				} else if screen.scrStatus == SHOW_LOUD_SELL_ORDERS {
					screen.RunSelectedLoudSellTrade()
				}
			}
		case "P": // UPGRADE ITEM
			fallthrough
		case "p":
			screen.scrStatus = SELECT_UPGRADE_ITEM
			screen.refreshed = false
		case "N": // Go hunt with no weapon
			fallthrough
		case "n":
			fallthrough
		case "I":
			fallthrough
		case "i":
			fallthrough
		case "1": // SELECT 1st item
			fallthrough
		case "2": // SELECT 2nd item
			fallthrough
		case "3": // SELECT 3rd item
			fallthrough
		case "4": // SELECT 4th item
			fallthrough
		case "5": // SELECT 5rd item
			fallthrough
		case "6": // SELECT 6rd item
			fallthrough
		case "7": // SELECT 7rd item
			fallthrough
		case "8": // SELECT 8rd item
			fallthrough
		case "9": // SELECT 9rd item
			screen.refreshed = false
			switch screen.scrStatus {
			case SELECT_BUY_ITEM:
				screen.activeItem = GetToBuyItemFromKey(Key)
				if len(screen.activeItem.Name) == 0 {
					return
				}
				screen.scrStatus = WAIT_BUY_PROCESS
				screen.refreshed = false
				screen.Render()
				log.Println("started sending request for buying item")
				txhash, err := Buy(screen.user, Key)
				log.Println("ended sending request for buying item")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_BUY_FINISH
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(1*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_BUY_FINISH
						screen.refreshed = false
						screen.Render()
					})
				}
			case SELECT_HUNT_ITEM:
				screen.activeItem = GetWeaponItemFromKey(screen.user, Key)
				screen.scrStatus = WAIT_HUNT_PROCESS
				screen.refreshed = false
				screen.Render()
				log.Println("started sending request for hunting item")
				txhash, err := Hunt(screen.user, Key)
				log.Println("ended sending request for hunting item")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_HUNT_FINISH
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(1*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_HUNT_FINISH
						screen.refreshed = false
						screen.Render()
					})
				}
			case SELECT_SELL_ITEM:
				screen.activeItem = GetToSellItemFromKey(screen.user, Key)
				if len(screen.activeItem.Name) == 0 {
					return
				}
				screen.scrStatus = WAIT_SELL_PROCESS
				screen.refreshed = false
				screen.Render()
				log.Println("started sending request for selling item")
				txhash, err := Sell(screen.user, Key)
				log.Println("ended sending request for selling item")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_SELL_FINISH
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(1*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_SELL_FINISH
						screen.refreshed = false
						screen.Render()
					})
				}
			case SELECT_UPGRADE_ITEM:
				screen.activeItem = GetToUpgradeItemFromKey(screen.user, Key)
				if len(screen.activeItem.Name) == 0 {
					return
				}
				screen.scrStatus = WAIT_UPGRADE_PROCESS
				screen.refreshed = false
				screen.Render()
				log.Println("started sending request for upgrading item")
				txhash, err := Upgrade(screen.user, Key)
				log.Println("ended sending request for upgrading item")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_UPGRADE_FINISH
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(1*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_UPGRADE_FINISH
						screen.refreshed = false
						screen.Render()
					})
				}
			}
		}
	}
	screen.Render()
}
