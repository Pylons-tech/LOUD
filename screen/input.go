package screen

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/atotto/clipboard"
	"github.com/nsf/termbox-go"
)

func (screen *GameScreen) HandleInputKey(input termbox.Event) {
	screen.actionText = ""
	if screen.IsWaitScreen() {
		return
	}
	Key := strings.ToUpper(string(input.Ch))
	log.Println("Handling Key \"", Key, "\"", input.Ch)
	if screen.InputActive() {
		screen.HandleTypingModeInputKeys(input)
	} else {
		if screen.HandleFirstClassInputKeys(input) {
			return
		}
		if screen.HandleSecondClassInputKeys(input) {
			return
		}
		if screen.HandleThirdClassInputKeys(input) {
			return
		}
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
		"C": loud.PYLCNTRL,
		"D": loud.DEVELOP,
	}

	if newStus, ok := tarLctMap[Key]; ok {
		if newStus == loud.FOREST && screen.user.GetActiveCharacter() == nil {
			screen.actionText = "You can't go to forest without character"
			screen.FreshRender()
		} else {
			screen.user.SetLocation(newStus)
			screen.FreshRender()
			return true
		}
	}
	return false
}
func (screen *GameScreen) HandleInputKeyHomeEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarStusMap := map[string]ScreenStatus{
		"1": SEL_ACTIVE_CHAR,
		"2": SEL_ACTIVE_WEAPON,
		"3": SEL_HEALTH_RESTORE_CHAR,
		"4": SEL_RENAME_CHAR,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		switch newStus {
		case SEL_ACTIVE_CHAR:
			screen.activeLine = screen.user.GetActiveCharacterIndex()
		case SEL_ACTIVE_WEAPON:
			screen.activeLine = screen.user.GetActiveWeaponIndex()
		}
		screen.FreshRender()
		return true
	} else {
		return false
	}
}
func (screen *GameScreen) HandleInputKeyPylonsCentralEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarStusMap := map[string]ScreenStatus{
		"1": SEL_BUYCHR,
		"2": W8_BUY_GOLD_WITH_PYLONS,
		"3": SHW_LOUD_BUY_TRDREQS,
		"4": SHW_LOUD_SELL_TRDREQS,
		"5": SHW_BUYITM_TRDREQS,
		"6": SHW_SELLITM_TRDREQS,
		"7": SHW_BUYCHR_TRDREQS,
		"8": SHW_SELLCHR_TRDREQS,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		if newStus == W8_BUY_GOLD_WITH_PYLONS {
			screen.RunTxProcess(W8_BUY_GOLD_WITH_PYLONS, RSLT_BUY_GOLD_WITH_PYLONS, func() (string, error) {
				return loud.BuyGoldWithPylons(screen.user)
			})
		} else {
			screen.scrStatus = newStus
			screen.FreshRender()
		}
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
		screen.FreshRender()
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) HandleInputKeyForestEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]ScreenStatus{
		"1": CONFIRM_HUNT_RABBITS,
		"2": CONFIRM_FIGHT_GOBLIN,
		"3": CONFIRM_FIGHT_WOLF,
		"4": CONFIRM_FIGHT_TROLL,
		"5": CONFIRM_FIGHT_GIANT,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		activeWeapon := screen.user.GetActiveWeapon()
		switch newStus {
		case CONFIRM_FIGHT_GOBLIN,
			CONFIRM_FIGHT_WOLF,
			CONFIRM_FIGHT_TROLL:
			if activeWeapon == nil {
				screen.actionText = "You need a sword for this action!"
				screen.FreshRender()
				return true
			}
		case CONFIRM_FIGHT_GIANT:
			if activeWeapon == nil || activeWeapon.Name != loud.IRON_SWORD {
				screen.actionText = "You need an iron sword for this action!"
				screen.FreshRender()
				return true
			}
		}
		screen.scrStatus = newStus
		screen.FreshRender()
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) HandleInputKeyShopEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]ScreenStatus{
		"1": SEL_BUYITM,
		"2": SEL_SELLITM,
		"3": SEL_UPGITM,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.scrStatus = newStus
		screen.FreshRender()
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) MoveToNextStep() {
	switch screen.scrStatus {
	case CONFIRM_HUNT_RABBITS:
		screen.RunHuntRabbits()
		return
	case CONFIRM_FIGHT_GIANT:
		screen.RunFightGiant()
		return
	case CONFIRM_FIGHT_TROLL:
		screen.RunFightTroll()
		return
	case CONFIRM_FIGHT_WOLF:
		screen.RunFightWolf()
		return
	case CONFIRM_FIGHT_GOBLIN:
		screen.RunFightGoblin()
		return
	}
	nextMapper := map[ScreenStatus]ScreenStatus{
		RSLT_HUNT_RABBITS:              CONFIRM_HUNT_RABBITS,
		RSLT_FIGHT_GOBLIN:              CONFIRM_FIGHT_GOBLIN,
		RSLT_FIGHT_TROLL:               CONFIRM_FIGHT_TROLL,
		RSLT_FIGHT_WOLF:                CONFIRM_FIGHT_WOLF,
		RSLT_FIGHT_GIANT:               CONFIRM_FIGHT_GIANT,
		RSLT_BUY_LOUD_TRDREQ_CREATION:  SHW_LOUD_BUY_TRDREQS,
		RSLT_FULFILL_BUY_LOUD_TRDREQ:   SHW_LOUD_BUY_TRDREQS,
		RSLT_SELL_LOUD_TRDREQ_CREATION: SHW_LOUD_SELL_TRDREQS,
		RSLT_FULFILL_SELL_LOUD_TRDREQ:  SHW_LOUD_SELL_TRDREQS,
		RSLT_SELLITM_TRDREQ_CREATION:   SHW_SELLITM_TRDREQS,
		RSLT_FULFILL_SELLITM_TRDREQ:    SHW_SELLITM_TRDREQS,
		RSLT_BUYITM_TRDREQ_CREATION:    SHW_BUYITM_TRDREQS,
		RSLT_FULFILL_BUYITM_TRDREQ:     SHW_BUYITM_TRDREQS,
		RSLT_SELLCHR_TRDREQ_CREATION:   SHW_SELLCHR_TRDREQS,
		RSLT_FULFILL_SELLCHR_TRDREQ:    SHW_SELLCHR_TRDREQS,
		RSLT_BUYCHR_TRDREQ_CREATION:    SHW_BUYCHR_TRDREQS,
		RSLT_CANCEL_TRDREQ:             SHW_LOCATION,
		RSLT_FULFILL_BUYCHR_TRDREQ:     SHW_BUYCHR_TRDREQS,
		RSLT_HEALTH_RESTORE_CHAR:       SEL_HEALTH_RESTORE_CHAR,
		RSLT_RENAME_CHAR:               SEL_RENAME_CHAR,
		RSLT_SEL_ACT_CHAR:              SEL_ACTIVE_CHAR,
		RSLT_SEL_ACT_WEAPON:            SEL_ACTIVE_WEAPON,
		RSLT_BUYITM:                    SEL_BUYITM,
		RSLT_BUYCHR:                    SEL_BUYCHR,
		RSLT_SELLITM:                   SEL_SELLITM,
		RSLT_UPGITM:                    SEL_UPGITM,
	}
	if nextStatus, ok := nextMapper[screen.scrStatus]; ok {
		if screen.user.GetLocation() == loud.DEVELOP {
			screen.scrStatus = SHW_LOCATION
		} else {
			screen.scrStatus = nextStatus
		}
	} else {
		screen.scrStatus = SHW_LOCATION
	}
	screen.txFailReason = ""
	screen.FreshRender()
}

func (screen *GameScreen) MoveToPrevStep() {
	prevMapper := map[ScreenStatus]ScreenStatus{
		CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL:  SHW_LOUD_BUY_TRDREQS,
		CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL:  CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL,
		CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL: SHW_LOUD_SELL_TRDREQS,
		CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL: CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL,
		CR8_SELLITM_TRDREQ_SEL_ITEM:     SHW_SELLITM_TRDREQS,
		CR8_SELLITM_TRDREQ_ENT_PYLVAL:   CR8_SELLITM_TRDREQ_SEL_ITEM,
		CR8_BUYITM_TRDREQ_SEL_ITEM:      SHW_BUYITM_TRDREQS,
		CR8_BUYITM_TRDREQ_ENT_PYLVAL:    CR8_BUYITM_TRDREQ_SEL_ITEM,
		CR8_SELLCHR_TRDREQ_SEL_CHR:      SHW_SELLCHR_TRDREQS,
		CR8_SELLCHR_TRDREQ_ENT_PYLVAL:   CR8_SELLCHR_TRDREQ_SEL_CHR,
		CR8_BUYCHR_TRDREQ_SEL_CHR:       SHW_BUYCHR_TRDREQS,
		CR8_BUYCHR_TRDREQ_ENT_PYLVAL:    CR8_BUYCHR_TRDREQ_SEL_CHR,
	}

	nxtStatus := SHW_LOCATION
	if nextStatus, ok := prevMapper[screen.scrStatus]; ok {
		nxtStatus = nextStatus
	}

	// move to home if it's somewhere else's 1st step
	if nxtStatus == SHW_LOCATION && screen.scrStatus == SHW_LOCATION {
		screen.user.SetLocation(loud.HOME)
	}

	screen.scrStatus = nxtStatus
	screen.FreshRender()
}

func (screen *GameScreen) HandleFirstClassInputKeys(input termbox.Event) bool {
	// implement first class commands, eg. development input keys
	if screen.HandleInputKeyLocationSwitch(input) {
		return true
	}
	Key := strings.ToUpper(string(input.Ch))
	switch Key {
	case "J": // Create cookbook
		screen.RunTxProcess(W8_CREATE_COOKBOOK, RSLT_CREATE_COOKBOOK, func() (string, error) {
			return loud.CreateCookbook(screen.user)
		})
	case "Z": // Switch user
		screen.SetScreenStatusAndRefresh(W8_SWITCH_USER)
		go func() {
			newUser := screen.world.GetUser(fmt.Sprintf("%d", time.Now().Unix()))
			orgLocation := screen.user.GetLocation()
			screen.SwitchUser(newUser)           // this is moving user back to home
			screen.user.SetLocation(orgLocation) // set the user back to original location
			screen.SetScreenStatusAndRefresh(RSLT_SWITCH_USER)
		}()
	case "Y": // get initial pylons
		screen.RunTxProcess(W8_GET_PYLONS, RSLT_GET_PYLONS, func() (string, error) {
			return loud.GetExtraPylons(screen.user)
		})
	case "B": // DEV ITEMS GET (troll toes, goblin ear, wolf tail)
		screen.RunTxProcess(W8_DEV_GET_TEST_ITEMS, RSLT_DEV_GET_TEST_ITEMS, func() (string, error) {
			return loud.DevGetTestItems(screen.user)
		})
	case "L": // copy last transaction to CLIPBOARD
		clipboard.WriteAll(screen.user.GetLastTransaction())
	case "M": // copy last transaction to CLIPBOARD
		clipboard.WriteAll(screen.user.GetAddress())
	case "E": // REFRESH
		screen.Resync()
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
		case SHW_LOCATION:
			return screen.HandleInputKeyHomeEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.PYLCNTRL {
		switch screen.scrStatus {
		case SHW_LOCATION:
			return screen.HandleInputKeyPylonsCentralEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.SETTINGS {
		switch screen.scrStatus {
		case SHW_LOCATION:
			return screen.HandleInputKeySettingsEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.FOREST {
		switch screen.scrStatus {
		case SHW_LOCATION:
			return screen.HandleInputKeyForestEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.SHOP {
		switch screen.scrStatus {
		case SHW_LOCATION:
			return screen.HandleInputKeyShopEntryPoint(input)
		}
	}
	return false
}

func (screen *GameScreen) HandleThirdClassInputKeys(input termbox.Event) bool {
	// implement thid class commands, eg. commands which are not processed by first, second classes
	Key := strings.ToUpper(string(input.Ch))
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
		if screen.user.GetLocation() == loud.PYLCNTRL {
			switch screen.scrStatus {
			case SHW_LOUD_BUY_TRDREQS:
				screen.scrStatus = CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL
			case SHW_LOUD_SELL_TRDREQS:
				screen.scrStatus = CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL
			case SHW_SELLITM_TRDREQS:
				screen.scrStatus = CR8_SELLITM_TRDREQ_SEL_ITEM
			case SHW_BUYITM_TRDREQS:
				screen.scrStatus = CR8_BUYITM_TRDREQ_SEL_ITEM
			case SHW_SELLCHR_TRDREQS:
				screen.scrStatus = CR8_SELLCHR_TRDREQ_SEL_CHR
			case SHW_BUYCHR_TRDREQS:
				screen.scrStatus = CR8_BUYCHR_TRDREQ_SEL_CHR
			}
			screen.inputText = ""
			screen.FreshRender()
			return true
		}
	case "O": // GO ON
		screen.MoveToNextStep()
		return true
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9": // Numbers
		switch screen.scrStatus {
		case SEL_ACTIVE_CHAR:
			screen.activeLine = loud.GetIndexFromString(Key)
			screen.RunActiveCharacterSelect(screen.activeLine)
		case SEL_ACTIVE_WEAPON:
			screen.activeLine = loud.GetIndexFromString(Key)
			screen.RunActiveWeaponSelect(screen.activeLine)
		case SEL_HEALTH_RESTORE_CHAR:
			screen.activeLine = loud.GetIndexFromString(Key)
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunCharacterHealthRestore()
		case SEL_RENAME_CHAR:
			screen.activeLine = loud.GetIndexFromString(Key)
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.scrStatus = RENAME_CHAR_ENT_NEWNAME
			screen.inputText = ""
			screen.FreshRender()
		case SEL_BUYITM:
			screen.activeItem = loud.GetToBuyItemFromKey(Key)
			if len(screen.activeItem.Name) == 0 {
				return false
			}
			screen.RunActiveItemBuy()
		case SEL_BUYCHR:
			screen.activeCharacter = loud.GetToBuyCharacterFromKey(Key)
			if len(screen.activeCharacter.Name) == 0 {
				return false
			}
			screen.RunActiveCharacterBuy()
		case SEL_SELLITM:
			screen.activeItem = loud.GetToSellItemFromKey(screen.user, Key)
			if len(screen.activeItem.Name) == 0 {
				return false
			}
			screen.RunActiveItemSell()

		case SEL_UPGITM:
			screen.activeItem = loud.GetToUpgradeItemFromKey(screen.user, Key)
			if len(screen.activeItem.Name) == 0 {
				return false
			}
			screen.RunActiveItemUpgrade()
		}
		screen.FreshRender()
		return true
	}
	return false
}

func (screen *GameScreen) HandleThirdClassKeyEnterEvent() bool {
	switch screen.user.GetLocation() {
	case loud.HOME, loud.PYLCNTRL, loud.SHOP, loud.FOREST:
		switch screen.scrStatus {
		case SHW_LOUD_BUY_TRDREQS:
			screen.RunSelectedLoudBuyTrdReq()
		case SHW_LOUD_SELL_TRDREQS:
			screen.RunSelectedLoudSellTrdReq()
		case SHW_BUYITM_TRDREQS:
			screen.RunSelectedItemBuyTrdReq()
		case SHW_SELLITM_TRDREQS:
			screen.RunSelectedItemSellTrdReq()
		case SHW_BUYCHR_TRDREQS:
			screen.RunSelectedCharacterBuyTrdReq()
		case SHW_SELLCHR_TRDREQS:
			screen.RunSelectedCharacterSellTrdReq()
		case CR8_SELLITM_TRDREQ_SEL_ITEM:
			userItems := screen.user.InventoryItems()
			if len(userItems) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = userItems[screen.activeLine]
			screen.scrStatus = CR8_SELLITM_TRDREQ_ENT_PYLVAL
			screen.inputText = ""
			screen.FreshRender()
		case CR8_BUYITM_TRDREQ_SEL_ITEM:
			if len(loud.WorldItemSpecs) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItSpec = loud.WorldItemSpecs[screen.activeLine]
			screen.scrStatus = CR8_BUYITM_TRDREQ_ENT_PYLVAL
			screen.inputText = ""
			screen.FreshRender()
		case CR8_SELLCHR_TRDREQ_SEL_CHR:
			userCharacters := screen.user.InventoryCharacters()
			if len(userCharacters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = userCharacters[screen.activeLine]
			screen.scrStatus = CR8_SELLCHR_TRDREQ_ENT_PYLVAL
			screen.inputText = ""
			screen.FreshRender()
		case CR8_BUYCHR_TRDREQ_SEL_CHR:
			if len(loud.WorldCharacterSpecs) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeChSpec = loud.WorldCharacterSpecs[screen.activeLine]
			screen.scrStatus = CR8_BUYCHR_TRDREQ_ENT_PYLVAL
			screen.inputText = ""
			screen.FreshRender()
		case SEL_ACTIVE_CHAR:
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunActiveCharacterSelect(screen.activeLine)
		case SEL_ACTIVE_WEAPON:
			items := screen.user.InventorySwords()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveWeaponSelect(screen.activeLine)
		case SEL_HEALTH_RESTORE_CHAR:
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunCharacterHealthRestore()
		case SEL_RENAME_CHAR:
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.scrStatus = RENAME_CHAR_ENT_NEWNAME
			screen.inputText = ""
			screen.FreshRender()
		case SEL_BUYITM:
			items := loud.ShopItems
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemBuy()
			log.Println("SEL_BUYITM", screen.activeItem)
		case SEL_BUYCHR:
			characters := loud.ShopCharacters
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunActiveCharacterBuy()
			log.Println("SEL_BUYCHR", screen.activeCharacter)
		case SEL_SELLITM:
			items := screen.user.InventoryItems()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemSell()
		case SEL_UPGITM:
			items := screen.user.InventoryUpgradableItems()
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

func (screen *GameScreen) HandleTypingModeInputKeys(input termbox.Event) bool {
	switch input.Key {
	case termbox.KeyBackspace2,
		termbox.KeyBackspace:

		log.Println("Pressed Backspace")
		lastIdx := len(screen.inputText) - 1
		if lastIdx < 0 {
			lastIdx = 0
		}
		screen.SetInputTextAndRender(screen.inputText[:lastIdx])
		return true
	case termbox.KeyEnter:
		switch screen.scrStatus {
		case RENAME_CHAR_ENT_NEWNAME:
			screen.RunCharacterRename(screen.inputText)
		case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL:
			screen.scrStatus = CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL
			screen.loudEnterValue = screen.inputText
			screen.inputText = ""
			screen.FreshRender()
		case CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL:
			screen.scrStatus = W8_BUY_LOUD_TRDREQ_CREATION
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndFreshRender("")
			txhash, err := loud.CreateBuyLoudTrdReq(screen.user, screen.loudEnterValue, screen.pylonEnterValue)
			log.Println("ended sending request for creating buy loud request")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RSLT_BUY_LOUD_TRDREQ_CREATION)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RSLT_BUY_LOUD_TRDREQ_CREATION)
				})
			}
		case CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL:
			screen.scrStatus = CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL
			screen.FreshRender()
			screen.loudEnterValue = screen.inputText
			screen.inputText = ""
		case CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL:
			screen.scrStatus = W8_SELL_LOUD_TRDREQ_CREATION
			screen.FreshRender()
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
			txhash, err := loud.CreateSellLoudTrdReq(screen.user, screen.loudEnterValue, screen.pylonEnterValue)

			log.Println("ended sending request for creating buy loud request")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RSLT_SELL_LOUD_TRDREQ_CREATION)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RSLT_SELL_LOUD_TRDREQ_CREATION)
				})
			}
		case CR8_SELLITM_TRDREQ_ENT_PYLVAL:
			screen.scrStatus = W8_SELLITM_TRDREQ_CREATION
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndFreshRender("")
			txhash, err := loud.CreateSellItemTrdReq(screen.user, screen.activeItem, screen.pylonEnterValue)
			log.Println("ended sending request for creating sword -> pylon request")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RSLT_SELLITM_TRDREQ_CREATION)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RSLT_SELLITM_TRDREQ_CREATION)
				})
			}
		case CR8_BUYITM_TRDREQ_ENT_PYLVAL:
			screen.scrStatus = W8_BUYITM_TRDREQ_CREATION
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndFreshRender("")
			txhash, err := loud.CreateBuyItemTrdReq(screen.user, screen.activeItSpec, screen.pylonEnterValue)
			log.Println("ended sending request for creating sword -> pylon request")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RSLT_BUYITM_TRDREQ_CREATION)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RSLT_BUYITM_TRDREQ_CREATION)
				})
			}

		case CR8_SELLCHR_TRDREQ_ENT_PYLVAL:
			screen.scrStatus = W8_SELLCHR_TRDREQ_CREATION
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndFreshRender("")
			txhash, err := loud.CreateSellCharacterTrdReq(screen.user, screen.activeCharacter, screen.pylonEnterValue)
			log.Println("ended sending request for creating character -> pylon request")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RSLT_SELLCHR_TRDREQ_CREATION)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RSLT_SELLCHR_TRDREQ_CREATION)
				})
			}
		case CR8_BUYCHR_TRDREQ_ENT_PYLVAL:
			screen.scrStatus = W8_BUYCHR_TRDREQ_CREATION
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndFreshRender("")
			txhash, err := loud.CreateBuyCharacterTrdReq(screen.user, screen.activeChSpec, screen.pylonEnterValue)
			log.Println("ended sending request for creating character -> pylon request")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RSLT_BUYCHR_TRDREQ_CREATION)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RSLT_BUYCHR_TRDREQ_CREATION)
				})
			}
		default:
			return false
		}
		return true
	default:
		iChar := string(input.Ch)
		Key := strings.ToUpper(iChar)
		if screen.scrStatus == RENAME_CHAR_ENT_NEWNAME {
			screen.SetInputTextAndRender(screen.inputText + iChar)
		} else if _, err := strconv.Atoi(Key); err == nil {
			// If user entered number, just use it
			screen.SetInputTextAndRender(screen.inputText + Key)
		}
		return false
	}
}
