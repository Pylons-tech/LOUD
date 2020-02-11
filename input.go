package loud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

func (screen *GameScreen) HandleInputKeyLocationSwitch(input termbox.Event) {
	Key := strings.ToUpper(string(input.Ch))

	tarLctMap := map[string]UserLocation{
		"F": FOREST,
		"S": SHOP,
		"H": HOME,
		"T": SETTINGS,
		"M": MARKET,
		"D": DEVELOP,
	}

	if newStus, ok := tarLctMap[Key]; ok {
		screen.user.SetLocation(newStus)
		screen.refreshed = false
	}
}
func (screen *GameScreen) HandleInputKeyMarketEntryPoint(input termbox.Event) {
	Key := string(input.Ch)

	tarStusMap := map[string]ScreenStatus{
		"1": SHOW_LOUD_BUY_ORDERS,
		"2": SHOW_LOUD_SELL_ORDERS,
		"3": SHOW_PYLON_SWORD_ORDERS,
		"4": SHOW_SWORD_PYLON_ORDERS,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.refreshed = false
	} else {
		screen.HandleInputKeyLocationSwitch(input)
	}
}

func (screen *GameScreen) HandleInputKeySettingsEntryPoint(input termbox.Event) {
	Key := string(input.Ch)

	tarLangMap := map[string]string{
		"1": "en",
		"2": "es",
	}

	if newLang, ok := tarLangMap[Key]; ok {
		gameLanguage = newLang
		screen.refreshed = false
	} else {
		screen.HandleInputKeyLocationSwitch(input)
	}
}

func (screen *GameScreen) HandleInputKeyForestEntryPoint(input termbox.Event) {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]ScreenStatus{
		"1": SELECT_HUNT_ITEM,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.refreshed = false
	} else {
		screen.HandleInputKeyLocationSwitch(input)
	}
}

func (screen *GameScreen) HandleInputKeyShopEntryPoint(input termbox.Event) {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]ScreenStatus{
		"1": SELECT_BUY_ITEM,
		"2": SELECT_SELL_ITEM,
		"3": SELECT_UPGRADE_ITEM,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.refreshed = false
	} else {
		screen.HandleInputKeyLocationSwitch(input)
	}
}

func (screen *GameScreen) HandleInputKey(input termbox.Event) {
	screen.lastInput = input
	Key := strings.ToUpper(string(input.Ch))
	log.Println("Handling Key \"", Key, "\"", input.Ch)
	if screen.user.GetLocation() == MARKET {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			screen.HandleInputKeyMarketEntryPoint(input)
			return
		}
	} else if screen.user.GetLocation() == SETTINGS {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			screen.HandleInputKeySettingsEntryPoint(input)
			return
		}
	} else if screen.user.GetLocation() == FOREST {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			screen.HandleInputKeyForestEntryPoint(input)
			return
		}
	} else if screen.user.GetLocation() == SHOP {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			screen.HandleInputKeyShopEntryPoint(input)
			return
		}
	}
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
				case SHOW_PYLON_SWORD_ORDERS:
					screen.RunSelectedSwordBuyOrder()
				case SHOW_SWORD_PYLON_ORDERS:
					screen.RunSelectedSwordSellOrder()
				}
			}
		}

		switch Key {
		case "J": // Create cookbook
			screen.scrStatus = WAIT_CREATE_COOKBOOK
			screen.refreshed = false
			screen.Render()
			go func() {
				txhash, err := CreateCookbook(screen.user)
				log.Println("ended sending request for creating cookbook")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_CREATE_COOKBOOK
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(1*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_CREATE_COOKBOOK
						screen.refreshed = false
						screen.Render()
					})
				}
			}()
		case "Z": // Switch user
			screen.scrStatus = WAIT_SWITCH_USER
			screen.refreshed = false
			screen.Render()
			go func() {
				newUser := screen.world.GetUser(fmt.Sprintf("%d", time.Now().Unix()))
				orgLocation := screen.user.GetLocation()
				screen.SwitchUser(newUser)           // this is moving user back to home
				screen.user.SetLocation(orgLocation) // set the user back to original location

				screen.scrStatus = RESULT_SWITCH_USER
				screen.refreshed = false
				screen.Render()
			}()
		case "H": // HOME
			screen.user.SetLocation(HOME)
			screen.refreshed = false
		case "F": // FOREST
			screen.user.SetLocation(FOREST)
			screen.refreshed = false
		case "S": // SHOP
			screen.user.SetLocation(SHOP)
			screen.refreshed = false
		case "M": // MARKET
			screen.user.SetLocation(MARKET)
			screen.refreshed = false
		case "T": // SETTINGS
			screen.user.SetLocation(SETTINGS)
			screen.refreshed = false
		case "D": // DEVELOP
			screen.user.SetLocation(DEVELOP)
			screen.refreshed = false
		case "C": // CANCEL
			screen.scrStatus = SHOW_LOCATION
			screen.refreshed = false
		case "O": // GO ON, GO BACK, CREATE ORDER
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
		case "Y": // get initial pylons
			screen.scrStatus = WAIT_GET_PYLONS
			screen.refreshed = false
			screen.Render()
			log.Println("started sending request for getting extra pylons")
			go func() {
				txhash, err := GetExtraPylons(screen.user)
				log.Println("ended sending request for getting extra pylons")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.scrStatus = RESULT_GET_PYLONS
					screen.refreshed = false
					screen.Render()
				} else {
					time.AfterFunc(1*time.Second, func() {
						screen.txResult, screen.txFailReason = ProcessTxResult(screen.user, txhash)
						screen.scrStatus = RESULT_GET_PYLONS
						screen.refreshed = false
						screen.Render()
					})
				}
			}()
		case "N": // Go hunt with no weapon
			fallthrough
		case "I":
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
				go func() {
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
				}()
			case SELECT_HUNT_ITEM:
				screen.activeItem = GetWeaponItemFromKey(screen.user, Key)
				screen.scrStatus = WAIT_HUNT_PROCESS
				screen.refreshed = false
				screen.Render()
				log.Println("started sending request for hunting item")
				go func() {
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
				}()
			case SELECT_SELL_ITEM:
				screen.activeItem = GetToSellItemFromKey(screen.user, Key)
				if len(screen.activeItem.Name) == 0 {
					return
				}
				screen.scrStatus = WAIT_SELL_PROCESS
				screen.refreshed = false
				screen.Render()
				log.Println("started sending request for selling item")
				go func() {
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
				}()
			case SELECT_UPGRADE_ITEM:
				screen.activeItem = GetToUpgradeItemFromKey(screen.user, Key)
				if len(screen.activeItem.Name) == 0 {
					return
				}
				screen.scrStatus = WAIT_UPGRADE_PROCESS
				screen.refreshed = false
				screen.Render()
				log.Println("started sending request for upgrading item")
				go func() {
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
				}()
			}
		}
	}
	screen.Render()
}
