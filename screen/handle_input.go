package screen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/LOUD/log"
	"github.com/atotto/clipboard"
	"github.com/nsf/termbox-go"
)

// HandleInputKey process keyboard input events
func (screen *GameScreen) HandleInputKey(input termbox.Event) {
	// initialize actionText since it's turning into a new command
	screen.actionText = ""

	// log input command
	Key := strings.ToUpper(string(input.Ch))
	logKey := Key
	switch input.Key {
	case termbox.KeyEsc:
		logKey = "Esc"
	case termbox.KeyBackspace2,
		termbox.KeyBackspace:
		logKey = "Backspace"
	case termbox.KeySpace:
		logKey = "Space"
	case termbox.KeyEnter:
		logKey = "Enter"
	}
	log.WithFields(log.Fields{
		"key":      logKey,
		"char_int": input.Ch,
	}).Infoln("Handling Key")

	if screen.IsWaitScreen() && !screen.IsWaitScreenCmd(input) {
		// restrict commands on wait screen
		return
	} else if screen.InputActive() {
		screen.HandleTypingModeInputKeys(input)
		screen.Render()
	} else if screen.HandleFirstClassInputKeys(input) {
		return
	} else if screen.HandleSecondClassInputKeys(input) {
		return
	} else if screen.HandleThirdClassInputKeys(input) {
		return
	}
}

// HandleInputKeyLocationSwitch process try to switch location with input events and returns false if it's not location switch key
func (screen *GameScreen) HandleInputKeyLocationSwitch(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarLctMap := map[string]loud.UserLocation{
		"F": loud.Forest,
		"S": loud.Shop,
		"H": loud.Home,
		"T": loud.Settings,
		"C": loud.PylonsCentral,
		"D": loud.Develop,
		"P": loud.Help,
	}

	if newLct, ok := tarLctMap[Key]; ok {
		if newLct == loud.Forest && screen.user.GetActiveCharacter() == nil {
			screen.actionText = loud.Sprintf("You can't go to forest without character")
			screen.Render()
		} else {
			screen.user.SetLocation(newLct)
			screen.SetScreenStatus(ShowLocation)
			screen.Render()
			return true
		}
	}
	return false
}

// HandleInputKeyHomeEntryPoint handles input key at home
func (screen *GameScreen) HandleInputKeyHomeEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarStusMap := map[string]PageStatus{
		"1": SelectActiveChr,
		"2": SelectRenameChr,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		if newStus == SelectRenameChr && len(screen.user.InventoryCharacters()) == 0 {
			screen.actionText = loud.Sprintf("You need a character for this action!")
			screen.Render()
			return true
		}
		screen.SetScreenStatus(newStus)
		screen.Render()
		return true
	}
	return false
}

// HandleInputKeyPylonsCentralEntryPoint handles input key at pylons central
func (screen *GameScreen) HandleInputKeyPylonsCentralEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarStusMap := map[string]PageStatus{
		"1": SelectBuyChr,
		"2": WaitByGoldWithPylons,
		"3": ShowGoldBuyTrdReqs,
		"4": ShowGoldSellTrdReqs,
		"5": ShowBuyItemTrdReqs,
		"6": ShowSellItemTrdReqs,
		"7": ShowBuyChrTrdReqs,
		"8": ShowSellChrTrdReqs,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		if newStus == WaitByGoldWithPylons {
			screen.RunTxProcess(WaitByGoldWithPylons, RsltByGoldWithPylons, func() (string, error) {
				return loud.BuyGoldWithPylons(screen.user)
			})
		} else {
			screen.SetScreenStatus(newStus)
			screen.Render()
		}
		return true
	}
	return false
}

// HandleInputKeySettingsEntryPoint handles input key at settings
func (screen *GameScreen) HandleInputKeySettingsEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarLangMap := map[string]string{
		"1": "en",
		"2": "es",
	}

	if newLang, ok := tarLangMap[Key]; ok {
		loud.GameLanguage = newLang
		screen.Render()
		return true
	}
	return false
}

// HandleInputKeyForestEntryPoint handles input key at forest entry point
func (screen *GameScreen) HandleInputKeyForestEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	monsterMap := map[string]string{
		"1": loud.TextRabbit,
		"2": loud.TextGoblin,
		"3": loud.TextWolf,
		"4": loud.TextTroll,
		"5": loud.TextGiant,
		"6": loud.TextDragonFire,
		"7": loud.TextDragonIce,
		"8": loud.TextDragonAcid,
		"9": loud.TextDragonUndead,
	}

	tarStusMap := map[string]PageStatus{
		"1": ConfirmHuntRabbits,
		"2": ConfirmFightGoblin,
		"3": ConfirmFightWolf,
		"4": ConfirmFightTroll,
		"5": ConfirmFightGiant,
		"6": ConfirmFightDragonFire,
		"7": ConfirmFightDragonIce,
		"8": ConfirmFightDragonAcid,
		"9": ConfirmFightDragonUndead,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		if fst, _ := screen.ForestStatusCheck(newStus); len(fst) > 0 {
			screen.actionText = fst
			screen.Render()
			return true
		}
		screen.user.SetFightMonster(monsterMap[Key])
		screen.SetScreenStatus(newStus)
		screen.Render()
		return true
	}
	return false
}

// HandleInputKeyShopEntryPoint handles input key for shop
func (screen *GameScreen) HandleInputKeyShopEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]PageStatus{
		"1": SelectBuyItem,
		"2": SelectSellItem,
		"3": SelectUpgradeItem,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.SetScreenStatus(newStus)
		if screen.activeLine < 0 {
			screen.activeLine = 0
		}
		screen.Render()
		return true
	}
	return false
}

// HandleInputKeyHelpEntryPoint handles input key for help page
func (screen *GameScreen) HandleInputKeyHelpEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]PageStatus{
		"1": HelpAbout,
		"2": HelpGameObjective,
		"3": HelpNavigation,
		"4": HelpPageLayout,
		"5": HelpGameRules,
		"6": HelpHowItWorks,
		"7": HelpPylonsCentral,
		"8": HelpUpcomingReleases,
		"9": HelpSupport,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.SetScreenStatus(newStus)
		screen.Render()
		return true
	}
	return false
}

// MoveToNextStep handles next step command when press enter
func (screen *GameScreen) MoveToNextStep() {
	activeCharacter := screen.user.GetActiveCharacter()

	switch screen.scrStatus {
	case ConfirmHuntRabbits:
		screen.RunHuntRabbits()
		return
	case ConfirmFightGoblin:
		screen.RunFightGoblin()
		return
	case ConfirmFightWolf:
		screen.RunFightWolf()
		return
	case ConfirmFightTroll:
		screen.RunFightTroll()
		return
	case ConfirmFightGiant:
		screen.RunFightGiant()
		return
	case ConfirmFightDragonFire:
		screen.RunFightDragonFire()
		return
	case ConfirmFightDragonIce:
		screen.RunFightDragonIce()
		return
	case ConfirmFightDragonAcid:
		screen.RunFightDragonAcid()
		return
	case ConfirmFightDragonUndead:
		screen.RunFightDragonUndead()
		return
	}
	nextMapper := map[PageStatus]PageStatus{
		RsltHuntRabbits:            ConfirmHuntRabbits,
		RsltFightGoblin:            ConfirmFightGoblin,
		RsltFightTroll:             ConfirmFightTroll,
		RsltFightWolf:              ConfirmFightWolf,
		RsltFightGiant:             ConfirmFightGiant,
		RsltFightDragonFire:        ConfirmFightDragonFire,
		RsltFightDragonIce:         ConfirmFightDragonIce,
		RsltFightDragonAcid:        ConfirmFightDragonAcid,
		RsltFightDragonUndead:      ConfirmFightDragonUndead,
		RsltBuyGoldTrdReqCreation:  ShowGoldBuyTrdReqs,
		RsltFulfillBuyGoldTrdReq:   ShowGoldBuyTrdReqs,
		RsltSellGoldTrdReqCreation: ShowGoldSellTrdReqs,
		RsltFulfillSellGoldTrdReq:  ShowGoldSellTrdReqs,
		RsltSellItemTrdReqCreation: ShowSellItemTrdReqs,
		RsltFulfillSellItemTrdReq:  ShowSellItemTrdReqs,
		RsltBuyItemTrdReqCreation:  ShowBuyItemTrdReqs,
		RsltFulfillBuyItemTrdReq:   ShowBuyItemTrdReqs,
		RsltSellChrTrdReqCreation:  ShowSellChrTrdReqs,
		RsltFulfillSellChrTrdReq:   ShowSellChrTrdReqs,
		RsltBuyChrTrdReqCreation:   ShowBuyChrTrdReqs,
		RsltCancelTrdReq:           ShowLocation,
		RsltFulfillBuyChrTrdReq:    ShowBuyChrTrdReqs,
		RsltRenameChr:              SelectRenameChr,
		RsltSelectActiveChr:        SelectActiveChr,
		RsltBuyItem:                SelectBuyItem,
		RsltBuyChr:                 SelectActiveChr,
		RsltSellItem:               SelectSellItem,
		RsltUpgradeItem:            SelectUpgradeItem,
	}
	if nextStatus, ok := nextMapper[screen.scrStatus]; ok {
		if screen.user.GetLocation() == loud.Develop {
			screen.SetScreenStatus(ShowLocation)
		} else if screen.user.GetLocation() == loud.Forest && activeCharacter == nil {
			// move back to home in forest if no active character
			screen.SetScreenStatus(ShowLocation)
		} else if nextStatus == ConfirmFightGiant && activeCharacter.Special != loud.NoSpecial {
			// go back to forest entrypoint when Special is not empty
			screen.SetScreenStatus(ShowLocation)
		} else if nextStatus == SelectActiveChr {
			screen.user.SetLocation(loud.Home)
			screen.SetScreenStatus(nextStatus)
		} else {
			screen.SetScreenStatus(nextStatus)
		}
	} else {
		screen.SetScreenStatus(ShowLocation)
	}
	screen.txFailReason = ""
	screen.Render()
}

// MoveToPrevStep handles input key for go back action when press backspace
func (screen *GameScreen) MoveToPrevStep() {
	activeCharacter := screen.user.GetActiveCharacter()

	prevMapper := map[PageStatus]PageStatus{
		CreateBuyGoldTrdReqEnterGoldValue:   ShowGoldBuyTrdReqs,
		CreateBuyGoldTrdReqEnterPylonValue:  CreateBuyGoldTrdReqEnterGoldValue,
		CreateSellGoldTrdReqEnterGoldValue:  ShowGoldSellTrdReqs,
		CreateSellGoldTrdReqEnterPylonValue: CreateSellGoldTrdReqEnterGoldValue,
		CreateSellItemTrdReqSelectItem:      ShowSellItemTrdReqs,
		CreateSellItemTrdReqEnterPylonValue: CreateSellItemTrdReqSelectItem,
		CreateBuyItemTrdReqSelectItem:       ShowBuyItemTrdReqs,
		CreateBuyItmTrdReqEnterPylonValue:   CreateBuyItemTrdReqSelectItem,
		CreateSellChrTrdReqSelChr:           ShowSellChrTrdReqs,
		CreateSellChrTrdReqEnterPylonValue:  CreateSellChrTrdReqSelChr,
		CreateBuyChrTrdReqSelectChr:         ShowBuyChrTrdReqs,
		CreateBuyChrTrdReqEnterPylonValue:   CreateBuyChrTrdReqSelectChr,
		SelectRenameChrEntNewName:           SelectRenameChr,
		RsltHuntRabbits:                     ConfirmHuntRabbits,
		RsltFightGoblin:                     ConfirmFightGoblin,
		RsltFightTroll:                      ConfirmFightTroll,
		RsltFightWolf:                       ConfirmFightWolf,
		RsltFightGiant:                      ConfirmFightGiant,
		RsltFightDragonFire:                 ConfirmFightDragonFire,
		RsltFightDragonIce:                  ConfirmFightDragonIce,
		RsltFightDragonAcid:                 ConfirmFightDragonAcid,
		RsltFightDragonUndead:               ConfirmFightDragonUndead,
		RsltBuyGoldTrdReqCreation:           ShowGoldBuyTrdReqs,
		RsltFulfillBuyGoldTrdReq:            ShowGoldBuyTrdReqs,
		RsltSellGoldTrdReqCreation:          ShowGoldSellTrdReqs,
		RsltFulfillSellGoldTrdReq:           ShowGoldSellTrdReqs,
		RsltSellItemTrdReqCreation:          ShowSellItemTrdReqs,
		RsltFulfillSellItemTrdReq:           ShowSellItemTrdReqs,
		RsltBuyItemTrdReqCreation:           ShowBuyItemTrdReqs,
		RsltFulfillBuyItemTrdReq:            ShowBuyItemTrdReqs,
		RsltSellChrTrdReqCreation:           ShowSellChrTrdReqs,
		RsltFulfillSellChrTrdReq:            ShowSellChrTrdReqs,
		RsltBuyChrTrdReqCreation:            ShowBuyChrTrdReqs,
		RsltCancelTrdReq:                    ShowLocation,
		RsltFulfillBuyChrTrdReq:             ShowBuyChrTrdReqs,
		RsltRenameChr:                       SelectRenameChr,
		RsltSelectActiveChr:                 SelectActiveChr,
		RsltBuyItem:                         SelectBuyItem,
		RsltBuyChr:                          SelectBuyChr,
		RsltSellItem:                        SelectSellItem,
		RsltUpgradeItem:                     SelectUpgradeItem,

		HelpAbout:            ShowLocation,
		HelpGameObjective:    ShowLocation,
		HelpNavigation:       ShowLocation,
		HelpPageLayout:       ShowLocation,
		HelpGameRules:        ShowLocation,
		HelpHowItWorks:       ShowLocation,
		HelpUpcomingReleases: ShowLocation,
		HelpSupport:          ShowLocation,
	}

	nxtStatus := ShowLocation
	if nextStatus, ok := prevMapper[screen.scrStatus]; ok {
		nxtStatus = nextStatus
	}

	switch nxtStatus {
	case CreateBuyGoldTrdReqEnterGoldValue,
		CreateSellGoldTrdReqEnterGoldValue:
		// set loud value previously entered
		screen.inputText = screen.goldEnterValue
	case ShowLocation:
		// move to home if it's somewhere else's entrypoint
		if screen.scrStatus == ShowLocation {
			screen.user.SetLocation(loud.Home)
		}
	case ConfirmFightGiant:
		if activeCharacter.Special != loud.NoSpecial {
			// go back to forest entrypoint when Special is not empty
			screen.SetScreenStatus(ShowLocation)
		}
	}

	if screen.user.GetLocation() == loud.Forest && activeCharacter == nil {
		// move back to home in forest if no active character
		screen.SetScreenStatus(ShowLocation)
		screen.user.SetLocation(loud.Home)
	}

	screen.SetScreenStatus(nxtStatus)
	screen.Render()
}

// HandleFirstClassInputKeys handles the keys that are level one
func (screen *GameScreen) HandleFirstClassInputKeys(input termbox.Event) bool {
	if input.Key == termbox.KeyEsc {
		switch screen.scrStatus {
		case ConfirmEndGame:
			screen.SetScreenStatus(ShowLocation)
		default:
			screen.SetScreenStatus(ConfirmEndGame)
		}
		screen.Render()
		return true
	}
	// implement first class commands, eg. development input keys
	if screen.HandleInputKeyLocationSwitch(input) {
		return true
	}
	Key := strings.ToUpper(string(input.Ch))
	switch Key {
	case "J": // Create cookbook
		if !loud.AutomateInput {
			return false
		}
		screen.RunTxProcess(WaitCreateCookbook, RsltCreateCookbook, func() (string, error) {
			return loud.CreateCookbook(screen.user)
		})
	case "Z": // Switch user
		screen.SetScreenStatusAndRefresh(WaitSwitchUser)
		go func() {
			newUser := screen.world.GetUser(fmt.Sprintf("%d", time.Now().Unix()))
			orgLocation := screen.user.GetLocation()
			screen.SwitchUser(newUser)           // this is moving user back to home
			screen.user.SetLocation(orgLocation) // set the user back to original location
			screen.SetScreenStatusAndRefresh(RsltSwitchUser)
		}()
	case "Y": // get initial pylons
		screen.RunTxProcess(WaitGetPylons, RsltGetPylons, func() (string, error) {
			return loud.GetExtraPylons(screen.user)
		})
	case "B": // DEV ITEMS GET (troll toes, goblin ear, wolf tail and drops of 3 special dragons)
		screen.RunTxProcess(WaitDevGetTestItems, RsltDevGetTestItems, func() (string, error) {
			return loud.DevGetTestItems(screen.user)
		})
	case "L": // copy last txhash to CLIPBOARD
		err := clipboard.WriteAll(screen.user.GetLastTxHash())
		return err == nil
	case "M": // copy user's cosmos address to CLIPBOARD
		err := clipboard.WriteAll(screen.user.GetAddress())
		return err == nil
	case "E": // REFRESH
		screen.Resync()
		return true
	default:
		return false
	}
	return true
}

// HandleSecondClassInputKeys handles the keys that are level 2
func (screen *GameScreen) HandleSecondClassInputKeys(input termbox.Event) bool {
	// implement second class commands, eg. input processing for show_location section
	if screen.user.GetLocation() == loud.Home {
		switch screen.scrStatus {
		case ShowLocation:
			return screen.HandleInputKeyHomeEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.PylonsCentral {
		switch screen.scrStatus {
		case ShowLocation:
			return screen.HandleInputKeyPylonsCentralEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.Settings {
		switch screen.scrStatus {
		case ShowLocation:
			return screen.HandleInputKeySettingsEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.Forest {
		switch screen.scrStatus {
		case ShowLocation:
			return screen.HandleInputKeyForestEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.Shop {
		switch screen.scrStatus {
		case ShowLocation:
			return screen.HandleInputKeyShopEntryPoint(input)
		}
	} else if screen.user.GetLocation() == loud.Help {
		switch screen.scrStatus {
		case ShowLocation:
			return screen.HandleInputKeyHelpEntryPoint(input)
		}
	}
	return false
}

// HandleThirdClassInputKeys handles the keys that are level 3
func (screen *GameScreen) HandleThirdClassInputKeys(input termbox.Event) bool {
	// implement thid class commands, eg. commands which are not processed by first, second classes
	Key := strings.ToUpper(string(input.Ch))
	switch input.Key {
	case termbox.KeyArrowLeft:
	case termbox.KeyArrowRight:
	case termbox.KeyArrowUp:
		if screen.activeLine > 0 {
			screen.activeLine--
		}
		return true
	case termbox.KeyArrowDown:
		screen.activeLine++
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
		if screen.user.GetLocation() == loud.PylonsCentral {
			newStatus := screen.scrStatus
			switch screen.scrStatus {
			case ShowGoldBuyTrdReqs:
				newStatus = CreateBuyGoldTrdReqEnterGoldValue
			case ShowGoldSellTrdReqs:
				newStatus = CreateSellGoldTrdReqEnterGoldValue
			case ShowSellItemTrdReqs:
				newStatus = CreateSellItemTrdReqSelectItem
			case ShowBuyItemTrdReqs:
				newStatus = CreateBuyItemTrdReqSelectItem
			case ShowSellChrTrdReqs:
				newStatus = CreateSellChrTrdReqSelChr
			case ShowBuyChrTrdReqs:
				newStatus = CreateBuyChrTrdReqSelectChr
			}
			screen.SetScreenStatus(newStatus)
			screen.inputText = ""
			screen.Render()
			return true
		}
	case "O": // GO ON
		screen.MoveToNextStep()
		return true
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9": // Numbers
		switch screen.scrStatus {
		case SelectActiveChr:
			screen.activeLine = loud.GetIndexFromString(Key)
			screen.RunActiveCharacterSelect(screen.activeLine)
		case SelectRenameChr:
			screen.activeLine = loud.GetIndexFromString(Key)
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.SetScreenStatus(SelectRenameChrEntNewName)
			screen.inputText = ""
			screen.Render()
		case SelectBuyItem:
			screen.activeItem = loud.GetToBuyItemFromKey(Key)
			if len(screen.activeItem.Name) == 0 {
				return false
			}
			screen.RunActiveItemBuy()
		case SelectBuyChr:
			screen.activeCharacter = loud.GetToBuyCharacterFromKey(Key)
			if len(screen.activeCharacter.Name) == 0 {
				return false
			}
			screen.RunActiveCharacterBuy()
		case SelectSellItem:
			screen.activeItem = loud.GetToSellItemFromKey(screen.user, Key)
			if len(screen.activeItem.Name) == 0 {
				return false
			}
			screen.RunActiveItemSell()

		case SelectUpgradeItem:
			screen.activeItem = loud.GetToUpgradeItemFromKey(screen.user, Key)
			if len(screen.activeItem.Name) == 0 {
				return false
			}
			screen.RunActiveItemUpgrade()
		}
		screen.Render()
		return true
	}
	return false
}

// HandleThirdClassKeyEnterEvent handles the keys that are level 3's enter event
func (screen *GameScreen) HandleThirdClassKeyEnterEvent() bool {
	switch screen.user.GetLocation() {
	case loud.Home, loud.PylonsCentral, loud.Shop, loud.Forest:
		switch screen.scrStatus {
		case ShowGoldBuyTrdReqs:
			screen.RunSelectedLoudBuyTrdReq()
		case ShowGoldSellTrdReqs:
			screen.RunSelectedLoudSellTrdReq()
		case ShowBuyItemTrdReqs:
			screen.RunSelectedItemBuyTrdReq()
		case ShowSellItemTrdReqs:
			screen.RunSelectedItemSellTrdReq()
		case ShowBuyChrTrdReqs:
			screen.RunSelectedCharacterBuyTrdReq()
		case ShowSellChrTrdReqs:
			screen.RunSelectedCharacterSellTrdReq()
		case CreateSellItemTrdReqSelectItem:
			userItems := screen.user.InventoryItems()
			if len(userItems) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = userItems[screen.activeLine]
			screen.SetScreenStatus(CreateSellItemTrdReqEnterPylonValue)
			screen.inputText = ""
			screen.Render()
		case CreateBuyItemTrdReqSelectItem:
			if len(loud.WorldItemSpecs) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItSpec = loud.WorldItemSpecs[screen.activeLine]
			screen.SetScreenStatus(CreateBuyItmTrdReqEnterPylonValue)
			screen.inputText = ""
			screen.Render()
		case CreateSellChrTrdReqSelChr:
			userCharacters := screen.user.InventoryCharacters()
			if len(userCharacters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = userCharacters[screen.activeLine]
			screen.SetScreenStatus(CreateSellChrTrdReqEnterPylonValue)
			screen.inputText = ""
			screen.Render()
		case CreateBuyChrTrdReqSelectChr:
			if len(loud.WorldCharacterSpecs) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeChSpec = loud.WorldCharacterSpecs[screen.activeLine]
			screen.SetScreenStatus(CreateBuyChrTrdReqEnterPylonValue)
			screen.inputText = ""
			screen.Render()
		case SelectActiveChr:
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunActiveCharacterSelect(screen.activeLine)
		case SelectRenameChr:
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.SetScreenStatus(SelectRenameChrEntNewName)
			screen.inputText = ""
			screen.Render()
		case SelectBuyItem:
			items := loud.ShopItems
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemBuy()
		case SelectBuyChr:
			characters := loud.ShopCharacters
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunActiveCharacterBuy()
		case SelectSellItem:
			items := screen.user.InventoryItems()
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemSell()
		case SelectUpgradeItem:
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

// HandleTypingModeInputKeys handles input keys for input active mode screens
func (screen *GameScreen) HandleTypingModeInputKeys(input termbox.Event) bool {
	switch input.Key {
	case termbox.KeyEsc:
		screen.MoveToPrevStep()
		return true
	case termbox.KeyBackspace2,
		termbox.KeyBackspace:

		lastIdx := len(screen.inputText) - 1
		if lastIdx < 0 {
			lastIdx = 0
		}
		screen.SetInputTextAndRender(screen.inputText[:lastIdx])
		return true
	case termbox.KeySpace:
		if screen.scrStatus == SelectRenameChrEntNewName {
			screen.SetInputTextAndRender(screen.inputText + " ")
			return true
		}
		return false
	case termbox.KeyEnter:
		switch screen.scrStatus {
		case SelectRenameChrEntNewName:
			screen.RunCharacterRename(screen.inputText)
		case CreateBuyGoldTrdReqEnterGoldValue:
			screen.SetScreenStatus(CreateBuyGoldTrdReqEnterPylonValue)
			screen.goldEnterValue = screen.inputText
			screen.inputText = ""
			screen.Render()
		case CreateBuyGoldTrdReqEnterPylonValue:
			screen.SetScreenStatus(WaitBuyGoldTrdReqCreation)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
			txhash, err := loud.CreateBuyGoldTrdReq(screen.user, screen.goldEnterValue, screen.pylonEnterValue)
			log.WithFields(log.Fields{
				"sent_request": "buy gold",
			}).Infoln("info log")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RsltBuyGoldTrdReqCreation)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RsltBuyGoldTrdReqCreation)
				})
			}
		case CreateSellGoldTrdReqEnterGoldValue:
			screen.SetScreenStatus(CreateSellGoldTrdReqEnterPylonValue)
			screen.Render()
			screen.goldEnterValue = screen.inputText
			screen.inputText = ""
		case CreateSellGoldTrdReqEnterPylonValue:
			screen.SetScreenStatus(WaitSellGoldTrdReqCreation)
			screen.Render()
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
			txhash, err := loud.CreateSellGoldTrdReq(screen.user, screen.goldEnterValue, screen.pylonEnterValue)
			log.WithFields(log.Fields{
				"sent_request": "sell gold",
			}).Infoln("info log")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RsltSellGoldTrdReqCreation)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RsltSellGoldTrdReqCreation)
				})
			}
		case CreateSellItemTrdReqEnterPylonValue:
			screen.SetScreenStatus(WaitSellItemTrdReqCreation)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
			txhash, err := loud.CreateSellItemTrdReq(screen.user, screen.activeItem, screen.pylonEnterValue)
			log.WithFields(log.Fields{
				"sent_request": "sell item",
			}).Infoln("info log")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RsltSellItemTrdReqCreation)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RsltSellItemTrdReqCreation)
				})
			}
		case CreateBuyItmTrdReqEnterPylonValue:
			screen.SetScreenStatus(WaitBuyItemTrdReqCreation)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
			txhash, err := loud.CreateBuyItemTrdReq(screen.user, screen.activeItSpec, screen.pylonEnterValue)
			log.WithFields(log.Fields{
				"sent_request": "buy item",
			}).Infoln("info log")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RsltBuyItemTrdReqCreation)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RsltBuyItemTrdReqCreation)
				})
			}

		case CreateSellChrTrdReqEnterPylonValue:
			screen.SetScreenStatus(WaitSellChrTrdReqCreation)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
			txhash, err := loud.CreateSellCharacterTrdReq(screen.user, screen.activeCharacter, screen.pylonEnterValue)
			log.WithFields(log.Fields{
				"sent_request": "sell character",
			}).Infoln("info log")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RsltSellChrTrdReqCreation)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RsltSellChrTrdReqCreation)
				})
			}
		case CreateBuyChrTrdReqEnterPylonValue:
			screen.SetScreenStatus(WaitBuyChrTrdReqCreation)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
			txhash, err := loud.CreateBuyCharacterTrdReq(screen.user, screen.activeChSpec, screen.pylonEnterValue)
			log.WithFields(log.Fields{
				"sent_request": "buy character",
			}).Infoln("info log")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.SetScreenStatusAndRefresh(RsltBuyChrTrdReqCreation)
			} else {
				time.AfterFunc(2*time.Second, func() {
					screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
					screen.SetScreenStatusAndRefresh(RsltBuyChrTrdReqCreation)
				})
			}
		default:
			return false
		}
		return true
	default:
		iChar := string(input.Ch)
		Key := strings.ToUpper(iChar)
		if screen.scrStatus == SelectRenameChrEntNewName {
			validNameStr := regexp.MustCompile(`^[a-zA-Z0-9\s$#@!%^&*()]$`)
			if validNameStr.MatchString(iChar) {
				screen.SetInputTextAndRender(screen.inputText + iChar)
				return true
			}
		} else if _, err := strconv.Atoi(Key); err == nil {
			// If user entered number, just use it
			screen.SetInputTextAndRender(screen.inputText + Key)
			return true
		}
		return false
	}
}
