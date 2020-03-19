package screen

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/nsf/termbox-go"
)

func (screen *GameScreen) HandleInputKey(input termbox.Event) {
	screen.lastInput = input
	Key := strings.ToUpper(string(input.Ch))
	log.Println("Handling Key \"", Key, "\"", input.Ch)
	if screen.HandleFirstClassInputKeys(input) {
		return
	}
	if screen.HandleSecondClassInputKeys(input) {
		return
	}
	if screen.HandleThirdClassInputKeys(input) {
		return
	}

	screen.Render()
}

func (screen *GameScreen) HandleInputKeyLocationSwitch(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarLctMap := map[string]loud.UserLocation{
		"F": loud.FOREST,
		"S": loud.SHOP,
		"H": loud.HOME,
		"T": loud.SETTINGS,
		"M": loud.MARKET,
		"D": loud.DEVELOP,
	}

	if newStus, ok := tarLctMap[Key]; ok {
		screen.user.SetLocation(newStus)
		screen.refreshed = false
		return true
	} else {
		return false
	}
}
func (screen *GameScreen) HandleInputKeyHomeEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarStusMap := map[string]ScreenStatus{
		"1": SELECT_DEFAULT_CHAR,
		"2": SELECT_DEFAULT_WEAPON,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.refreshed = false
		return true
	} else {
		return false
	}
}
func (screen *GameScreen) HandleInputKeyMarketEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarStusMap := map[string]ScreenStatus{
		"1": SHOW_LOUD_BUY_REQUESTS,
		"2": SHOW_LOUD_SELL_REQUESTS,
		"3": SHOW_BUY_SWORD_REQUESTS,
		"4": SHOW_SELL_SWORD_REQUESTS,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.refreshed = false
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) HandleInputKeySettingsEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarLangMap := map[string]string{
		"1": "en",
		"2": "es",
	}

	if newLang, ok := tarLangMap[Key]; ok {
		loud.GameLanguage = newLang
		screen.refreshed = false
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) HandleInputKeyForestEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]ScreenStatus{
		"1": SELECT_HUNT_ITEM,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.refreshed = false
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) HandleInputKeyShopEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]ScreenStatus{
		"1": SELECT_BUY_ITEM,
		"2": SELECT_SELL_ITEM,
		"3": SELECT_UPGRADE_ITEM,
		"4": SELECT_BUY_CHARACTER,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.refreshed = false
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) MoveToNextStep() {
	nextMapper := map[ScreenStatus]ScreenStatus{
		RESULT_HUNT_FINISH:                 SELECT_HUNT_ITEM,
		RESULT_BUY_LOUD_REQUEST_CREATION:   SHOW_LOUD_BUY_REQUESTS,
		RESULT_FULFILL_BUY_LOUD_REQUEST:    SHOW_LOUD_BUY_REQUESTS,
		RESULT_SELL_LOUD_REQUEST_CREATION:  SHOW_LOUD_SELL_REQUESTS,
		RESULT_FULFILL_SELL_LOUD_REQUEST:   SHOW_LOUD_SELL_REQUESTS,
		RESULT_SELL_SWORD_REQUEST_CREATION: SHOW_SELL_SWORD_REQUESTS,
		RESULT_FULFILL_SELL_SWORD_REQUEST:  SHOW_SELL_SWORD_REQUESTS,
		RESULT_BUY_SWORD_REQUEST_CREATION:  SHOW_BUY_SWORD_REQUESTS,
		RESULT_FULFILL_BUY_SWORD_REQUEST:   SHOW_BUY_SWORD_REQUESTS,
		RESULT_SELECT_DEF_CHAR:             SELECT_DEFAULT_CHAR,
		RESULT_SELECT_DEF_WEAPON:           SELECT_DEFAULT_WEAPON,
		RESULT_BUY_ITEM_FINISH:             SELECT_BUY_ITEM,
		RESULT_BUY_CHARACTER_FINISH:        SELECT_BUY_CHARACTER,
		RESULT_SELL_FINISH:                 SELECT_SELL_ITEM,
		RESULT_UPGRADE_FINISH:              SELECT_UPGRADE_ITEM,
	}
	if nextStatus, ok := nextMapper[screen.scrStatus]; ok {
		if screen.user.GetLocation() == loud.DEVELOP {
			screen.scrStatus = SHOW_LOCATION
		} else {
			screen.scrStatus = nextStatus
		}
	} else {
		screen.scrStatus = SHOW_LOCATION
	}
	screen.txFailReason = ""
	screen.refreshed = false
}

func (screen *GameScreen) MoveToPrevStep() {
	prevMapper := map[ScreenStatus]ScreenStatus{
		CREATE_BUY_LOUD_REQUEST_ENTER_LOUD_VALUE:    SHOW_LOUD_BUY_REQUESTS,
		CREATE_BUY_LOUD_REQUEST_ENTER_PYLON_VALUE:   CREATE_BUY_LOUD_REQUEST_ENTER_LOUD_VALUE,
		CREATE_SELL_LOUD_REQUEST_ENTER_LOUD_VALUE:   SHOW_LOUD_SELL_REQUESTS,
		CREATE_SELL_LOUD_REQUEST_ENTER_PYLON_VALUE:  CREATE_SELL_LOUD_REQUEST_ENTER_LOUD_VALUE,
		CREATE_SELL_SWORD_REQUEST_SELECT_SWORD:      SHOW_SELL_SWORD_REQUESTS,
		CREATE_SELL_SWORD_REQUEST_ENTER_PYLON_VALUE: CREATE_SELL_SWORD_REQUEST_SELECT_SWORD,
		CREATE_BUY_SWORD_REQUEST_SELECT_SWORD:       SHOW_BUY_SWORD_REQUESTS,
		CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE:  CREATE_BUY_SWORD_REQUEST_SELECT_SWORD,
	}
	if nextStatus, ok := prevMapper[screen.scrStatus]; ok {
		screen.scrStatus = nextStatus
	} else {
		screen.scrStatus = SHOW_LOCATION
	}
	screen.refreshed = false
}

func (screen *GameScreen) HandleFirstClassInputKeys(input termbox.Event) bool {
	// implement first class commands, eg. development input keys
	if screen.HandleInputKeyLocationSwitch(input) {
		return true
	}
	Key := strings.ToUpper(string(input.Ch))
	switch Key {
	case "J": // Create cookbook
		screen.SetScreenStatusAndRefresh(WAIT_CREATE_COOKBOOK)
		go func() {
			txhash, err := loud.CreateCookbook(screen.user)
			log.Println("ended sending request for creating cookbook")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RESULT_CREATE_COOKBOOK)
			} else {
				time.AfterFunc(1*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RESULT_CREATE_COOKBOOK)
				})
			}
		}()
	case "Z": // Switch user
		screen.SetScreenStatusAndRefresh(WAIT_SWITCH_USER)
		go func() {
			newUser := screen.world.GetUser(fmt.Sprintf("%d", time.Now().Unix()))
			orgLocation := screen.user.GetLocation()
			screen.SwitchUser(newUser)           // this is moving user back to home
			screen.user.SetLocation(orgLocation) // set the user back to original location
			screen.SetScreenStatusAndRefresh(RESULT_SWITCH_USER)
		}()
	case "Y": // get initial pylons
		screen.SetScreenStatusAndRefresh(WAIT_GET_PYLONS)
		log.Println("started sending request for getting extra pylons")
		go func() {
			txhash, err := loud.GetExtraPylons(screen.user)
			log.Println("ended sending request for getting extra pylons")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RESULT_GET_PYLONS)
			} else {
				time.AfterFunc(1*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RESULT_GET_PYLONS)
				})
			}
		}()
	case "I":
		screen.activeItem = loud.GetWeaponItemFromKey(screen.user, Key)
		screen.SetScreenStatusAndRefresh(WAIT_HUNT_PROCESS)
		log.Println("started sending request for hunting item")
		go func() {
			txhash, err := loud.Hunt(screen.user, loud.Item{}, true)
			log.Println("ended sending request for hunting item")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RESULT_HUNT_FINISH)
			} else {
				time.AfterFunc(1*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RESULT_HUNT_FINISH)
				})
			}
		}()
	case "E": // REFRESH
		screen.Resync()
		return true
	case "C": // CANCEL, GO BACK
		screen.MoveToPrevStep()
		return true
	default:
		return false
	}
	return true
}

func (screen *GameScreen) HandleSecondClassInputKeys(input termbox.Event) bool {
	// implement second class commands, eg. input processing for show_location section
	if screen.user.GetLocation() == loud.HOME {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			return screen.HandleInputKeyHomeEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.MARKET {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			return screen.HandleInputKeyMarketEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.SETTINGS {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			return screen.HandleInputKeySettingsEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.FOREST {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			return screen.HandleInputKeyForestEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.SHOP {
		switch screen.scrStatus {
		case SHOW_LOCATION:
			return screen.HandleInputKeyShopEntryPoint(input)
		}
	}
	return false
}

func (screen *GameScreen) HandleThirdClassInputKeys(input termbox.Event) bool {
	// implement thid class commands, eg. commands which are not processed by first, second classes
	Key := strings.ToUpper(string(input.Ch))
	if screen.InputActive() {
		switch input.Key {
		case termbox.KeyBackspace2:
			fallthrough
		case termbox.KeyBackspace:
			log.Println("Pressed Backspace")
			lastIdx := len(screen.inputText) - 1
			if lastIdx < 0 {
				lastIdx = 0
			}
			screen.SetInputTextAndRender(screen.inputText[:lastIdx])
			return true
		case termbox.KeyEnter:
			switch screen.scrStatus {
			case CREATE_BUY_LOUD_REQUEST_ENTER_LOUD_VALUE:
				screen.scrStatus = CREATE_BUY_LOUD_REQUEST_ENTER_PYLON_VALUE
				screen.loudEnterValue = screen.inputText
				screen.inputText = ""
			case CREATE_BUY_LOUD_REQUEST_ENTER_PYLON_VALUE:
				screen.scrStatus = WAIT_BUY_LOUD_REQUEST_CREATION
				screen.pylonEnterValue = screen.inputText
				screen.SetInputTextAndRender("")
				txhash, err := loud.CreateBuyLoudTradeRequest(screen.user, screen.loudEnterValue, screen.pylonEnterValue)
				log.Println("ended sending request for creating buy loud request")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.SetScreenStatusAndRefresh(RESULT_BUY_LOUD_REQUEST_CREATION)
				} else {
					time.AfterFunc(2*time.Second, func() {
						screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
						screen.SetScreenStatusAndRefresh(RESULT_BUY_LOUD_REQUEST_CREATION)
					})
				}
			case CREATE_SELL_LOUD_REQUEST_ENTER_LOUD_VALUE:
				screen.scrStatus = CREATE_SELL_LOUD_REQUEST_ENTER_PYLON_VALUE
				screen.loudEnterValue = screen.inputText
				screen.inputText = ""
			case CREATE_SELL_LOUD_REQUEST_ENTER_PYLON_VALUE:
				screen.scrStatus = WAIT_SELL_LOUD_REQUEST_CREATION
				screen.pylonEnterValue = screen.inputText
				screen.SetInputTextAndRender("")
				txhash, err := loud.CreateSellLoudTradeRequest(screen.user, screen.loudEnterValue, screen.pylonEnterValue)

				log.Println("ended sending request for creating buy loud request")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.SetScreenStatusAndRefresh(RESULT_SELL_LOUD_REQUEST_CREATION)
				} else {
					time.AfterFunc(2*time.Second, func() {
						screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
						screen.SetScreenStatusAndRefresh(RESULT_SELL_LOUD_REQUEST_CREATION)
					})
				}
			case CREATE_SELL_SWORD_REQUEST_ENTER_PYLON_VALUE:
				screen.scrStatus = WAIT_SELL_SWORD_REQUEST_CREATION
				screen.pylonEnterValue = screen.inputText
				screen.SetInputTextAndRender("")
				txhash, err := loud.CreateSellSwordTradeRequest(screen.user, screen.activeItem, screen.pylonEnterValue)
				log.Println("ended sending request for creating sword -> pylon request")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.SetScreenStatusAndRefresh(RESULT_SELL_SWORD_REQUEST_CREATION)
				} else {
					time.AfterFunc(2*time.Second, func() {
						screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
						screen.SetScreenStatusAndRefresh(RESULT_SELL_SWORD_REQUEST_CREATION)
					})
				}
			case CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE:
				screen.scrStatus = WAIT_BUY_SWORD_REQUEST_CREATION
				screen.pylonEnterValue = screen.inputText
				screen.SetInputTextAndRender("")
				txhash, err := loud.CreateBuySwordTradeRequest(screen.user, screen.activeItem, screen.pylonEnterValue)
				log.Println("ended sending request for creating sword -> pylon request")
				if err != nil {
					screen.txFailReason = err.Error()
					screen.SetScreenStatusAndRefresh(RESULT_BUY_SWORD_REQUEST_CREATION)
				} else {
					time.AfterFunc(2*time.Second, func() {
						screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
						screen.SetScreenStatusAndRefresh(RESULT_BUY_SWORD_REQUEST_CREATION)
					})
				}
			default:
				return false
			}
			return true
		default:
			if _, err := strconv.Atoi(Key); err == nil {
				// If user entered number, just use it
				screen.SetInputTextAndRender(screen.inputText + Key)
			}
			return false
		}
	} else {
		switch input.Key {
		case termbox.KeyArrowLeft:
		case termbox.KeyArrowRight:
		case termbox.KeyArrowUp:
			if screen.activeLine > 0 {
				screen.activeLine -= 1
			}
			return true
		case termbox.KeyArrowDown:
			screen.activeLine += 1
			return true
		}
		if input.Key == termbox.KeyEnter {
			return screen.HandleThirdClassKeyEnterEvent()
		}

		if input.Key == termbox.KeyBackspace2 || input.Key == termbox.KeyBackspace {
			screen.MoveToPrevStep()
		}

		switch Key {
		case "R": // CREATE ORDER
			if screen.user.GetLocation() == loud.MARKET {
				switch screen.scrStatus {
				case SHOW_LOUD_BUY_REQUESTS:
					screen.scrStatus = CREATE_BUY_LOUD_REQUEST_ENTER_LOUD_VALUE
				case SHOW_LOUD_SELL_REQUESTS:
					screen.scrStatus = CREATE_SELL_LOUD_REQUEST_ENTER_LOUD_VALUE
				case SHOW_SELL_SWORD_REQUESTS:
					screen.scrStatus = CREATE_SELL_SWORD_REQUEST_SELECT_SWORD
				case SHOW_BUY_SWORD_REQUESTS:
					screen.scrStatus = CREATE_BUY_SWORD_REQUEST_SELECT_SWORD
				}
				screen.refreshed = false
			}
		case "O": // GO ON
			screen.MoveToNextStep()
			return true
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9": // Numbers
			screen.refreshed = false
			switch screen.scrStatus {
			case SELECT_DEFAULT_CHAR:
				screen.activeLine = loud.GetIndexFromString(Key)
				characters := screen.user.InventoryCharacters()
				if len(characters) <= screen.activeLine || screen.activeLine < 0 {
					return false
				}
				screen.RunActiveCharacterSelect()
			case SELECT_DEFAULT_WEAPON:
				screen.activeLine = loud.GetIndexFromString(Key)
				items := screen.user.InventoryItems()
				if len(items) <= screen.activeLine || screen.activeLine < 0 {
					return false
				}
				screen.RunActiveWeaponSelect()
			case SELECT_BUY_ITEM:
				screen.activeItem = loud.GetToBuyItemFromKey(Key)
				if len(screen.activeItem.Name) == 0 {
					return false
				}
				screen.RunActiveItemBuy()
			case SELECT_BUY_CHARACTER:
				screen.activeItem = loud.GetToBuyCharacterFromKey(Key)
				if len(screen.activeItem.Name) == 0 {
					return false
				}
				screen.RunActiveCharacterBuy()
			case SELECT_HUNT_ITEM:
				screen.activeItem = loud.GetWeaponItemFromKey(screen.user, Key)
				screen.RunActiveItemHunt()
			case SELECT_SELL_ITEM:
				screen.activeItem = loud.GetToSellItemFromKey(screen.user, Key)
				if len(screen.activeItem.Name) == 0 {
					return false
				}
				screen.RunActiveItemSell()

			case SELECT_UPGRADE_ITEM:
				screen.activeItem = loud.GetToUpgradeItemFromKey(screen.user, Key)
				if len(screen.activeItem.Name) == 0 {
					return false
				}
				screen.RunActiveItemUpgrade()
			}
			return true
		}
	}
	return false
}

func (screen *GameScreen) HandleThirdClassKeyEnterEvent() bool {
	switch screen.user.GetLocation() {
	case loud.HOME, loud.MARKET, loud.SHOP, loud.FOREST:
		switch screen.scrStatus {
		case SHOW_LOUD_BUY_REQUESTS:
			screen.RunSelectedLoudBuyTrade()
		case SHOW_LOUD_SELL_REQUESTS:
			screen.RunSelectedLoudSellTrade()
		case SHOW_BUY_SWORD_REQUESTS:
			screen.RunSelectedSwordBuyTradeRequest()
		case SHOW_SELL_SWORD_REQUESTS:
			screen.RunSelectedSwordSellTradeRequest()
		case CREATE_SELL_SWORD_REQUEST_SELECT_SWORD:
			userItems := screen.user.InventoryItems()
			if len(userItems) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = userItems[screen.activeLine]
			screen.scrStatus = CREATE_SELL_SWORD_REQUEST_ENTER_PYLON_VALUE
			screen.inputText = ""
			screen.refreshed = false
		case CREATE_BUY_SWORD_REQUEST_SELECT_SWORD:
			if len(loud.WorldItems) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = loud.WorldItems[screen.activeLine]
			screen.scrStatus = CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE
			screen.inputText = ""
			screen.refreshed = false
		case SELECT_DEFAULT_CHAR:
			items := screen.user.InventoryCharacters()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveCharacterSelect()
			log.Println("SELECT_DEFAULT_CHAR", screen.activeItem)
		case SELECT_DEFAULT_WEAPON:
			items := screen.user.InventoryItems()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveWeaponSelect()
			log.Println("SELECT_DEFAULT_WEAPON", screen.activeItem)
		case SELECT_BUY_ITEM:
			items := loud.ShopItems
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemBuy()
			log.Println("SELECT_BUY_ITEM", screen.activeItem)
		case SELECT_BUY_CHARACTER:
			characs := loud.ShopCharacters
			if len(characs) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = characs[screen.activeLine]
			screen.RunActiveCharacterBuy()
			log.Println("SELECT_BUY_CHARACTER", screen.activeItem)
		case SELECT_HUNT_ITEM:
			items := screen.user.InventoryItems()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemHunt()
		case SELECT_SELL_ITEM:
			items := screen.user.InventoryItems()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemSell()
		case SELECT_UPGRADE_ITEM:
			items := screen.user.UpgradableItems()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemUpgrade()
		default:
			screen.MoveToNextStep()
			return false
		}
	default:
		screen.MoveToNextStep()
		return false
	}
	return true
}
