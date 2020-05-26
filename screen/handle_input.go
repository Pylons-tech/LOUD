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

func (screen *GameScreen) HandleInputKey(input termbox.Event) {
	// initialize actionText since it's turning into a new command
	screen.actionText = ""

	// log input command
	Key := strings.ToUpper(string(input.Ch))
	log.Println("Handling Key \"", Key, "\"", input.Ch)

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

func (screen *GameScreen) HandleInputKeyLocationSwitch(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarLctMap := map[string]loud.UserLocation{
		"F": loud.FOREST,
		"S": loud.SHOP,
		"H": loud.HOME,
		"T": loud.SETTINGS,
		"C": loud.PYLCNTRL,
		"D": loud.DEVELOP,
		"P": loud.HELP,
	}

	if newLct, ok := tarLctMap[Key]; ok {
		if newLct == loud.FOREST && screen.user.GetActiveCharacter() == nil {
			screen.actionText = loud.Sprintf("You can't go to forest without character")
			screen.Render()
		} else {
			screen.user.SetLocation(newLct)
			screen.SetScreenStatus(SHW_LOCATION)
			screen.Render()
			return true
		}
	}
	return false
}
func (screen *GameScreen) HandleInputKeyHomeEntryPoint(input termbox.Event) bool {
	Key := string(input.Ch)

	tarStusMap := map[string]ScreenStatus{
		"1": SEL_ACTIVE_CHAR,
		"2": SEL_RENAME_CHAR,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		if newStus == SEL_RENAME_CHAR && len(screen.user.InventoryCharacters()) == 0 {
			screen.actionText = loud.Sprintf("You need a character for this action!")
			screen.Render()
			return true
		}
		screen.SetScreenStatus(newStus)
		screen.Render()
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
			screen.SetScreenStatus(newStus)
			screen.Render()
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
		screen.Render()
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) ForestStatusCheck(newStus ScreenStatus) (string, string) {
	activeCharacter := screen.user.GetActiveCharacter()
	if activeCharacter == nil {
		return loud.Sprintf("You need a character for this action!"), loud.Sprintf("no character!")
	}
	switch newStus {
	case CONFIRM_FIGHT_GIANT:
		if activeCharacter == nil || activeCharacter.Special != loud.NO_SPECIAL {
			return loud.Sprintf("You need no special character for this action!"), loud.Sprintf("no non-special character!")
		}
	case CONFIRM_FIGHT_DRAGONFIRE:
		if activeCharacter == nil || activeCharacter.Special != loud.FIRE_SPECIAL {
			return loud.Sprintf("You need a fire character for this action!"), loud.Sprintf("no fire character!")
		}
	case CONFIRM_FIGHT_DRAGONICE:
		if activeCharacter == nil || activeCharacter.Special != loud.ICE_SPECIAL {
			return loud.Sprintf("You need a ice character for this action!"), loud.Sprintf("no ice character!")
		}
	case CONFIRM_FIGHT_DRAGONACID:
		if activeCharacter == nil || activeCharacter.Special != loud.ACID_SPECIAL {
			return loud.Sprintf("You need a acid character for this action!"), loud.Sprintf("no acid character!")
		}
	}
	switch newStus {
	case CONFIRM_FIGHT_GOBLIN,
		CONFIRM_FIGHT_WOLF,
		CONFIRM_FIGHT_TROLL:
		if len(screen.user.InventorySwords()) == 0 {
			return loud.Sprintf("You need a sword for this action!"), loud.Sprintf("no sword!")
		}
	case CONFIRM_FIGHT_GIANT,
		CONFIRM_FIGHT_DRAGONFIRE,
		CONFIRM_FIGHT_DRAGONICE,
		CONFIRM_FIGHT_DRAGONACID:
		if len(screen.user.InventoryIronSwords()) == 0 {
			return loud.Sprintf("You need an iron sword for this action!"), loud.Sprintf("no iron sword!")
		}
	case CONFIRM_FIGHT_DRAGONUNDEAD:
		if len(screen.user.InventoryAngelSwords()) == 0 {
			return loud.Sprintf("You need an angel sword for this action!"), loud.Sprintf("no angel sword!")
		}
	}
	return "", ""
}

func (screen *GameScreen) HandleInputKeyForestEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	monsterMap := map[string]string{
		"1": loud.RABBIT,
		"2": loud.GOBLIN,
		"3": loud.WOLF,
		"4": loud.TROLL,
		"5": loud.GIANT,
		"6": loud.DRAGON_FIRE,
		"7": loud.DRAGON_ICE,
		"8": loud.DRAGON_ACID,
		"9": loud.DRAGON_UNDEAD,
	}

	tarStusMap := map[string]ScreenStatus{
		"1": CONFIRM_HUNT_RABBITS,
		"2": CONFIRM_FIGHT_GOBLIN,
		"3": CONFIRM_FIGHT_WOLF,
		"4": CONFIRM_FIGHT_TROLL,
		"5": CONFIRM_FIGHT_GIANT,
		"6": CONFIRM_FIGHT_DRAGONFIRE,
		"7": CONFIRM_FIGHT_DRAGONICE,
		"8": CONFIRM_FIGHT_DRAGONACID,
		"9": CONFIRM_FIGHT_DRAGONUNDEAD,
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
		screen.SetScreenStatus(newStus)
		if screen.activeLine < 0 {
			screen.activeLine = 0
		}
		screen.Render()
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) HandleInputKeyHelpEntryPoint(input termbox.Event) bool {
	Key := strings.ToUpper(string(input.Ch))

	tarStusMap := map[string]ScreenStatus{
		"1": HELP_ABOUT,
		"2": HELP_GAME_OBJECTIVE,
		"3": HELP_NAVIGATION,
		"4": HELP_PAGE_LAYOUT,
		"5": HELP_GAME_RULES,
		"6": HELP_HOW_IT_WORKS,
		"7": HELP_PYLONS_CENTRAL,
		"8": HELP_UPCOMING_RELEASES,
		"9": HELP_SUPPORT,
	}

	if newStus, ok := tarStusMap[Key]; ok {
		screen.SetScreenStatus(newStus)
		screen.Render()
		return true
	} else {
		return false
	}
}

func (screen *GameScreen) MoveToNextStep() {
	activeCharacter := screen.user.GetActiveCharacter()

	switch screen.scrStatus {
	case CONFIRM_HUNT_RABBITS:
		screen.RunHuntRabbits()
		return
	case CONFIRM_FIGHT_GOBLIN:
		screen.RunFightGoblin()
		return
	case CONFIRM_FIGHT_WOLF:
		screen.RunFightWolf()
		return
	case CONFIRM_FIGHT_TROLL:
		screen.RunFightTroll()
		return
	case CONFIRM_FIGHT_GIANT:
		screen.RunFightGiant()
		return
	case CONFIRM_FIGHT_DRAGONFIRE:
		screen.RunFightDragonFire()
		return
	case CONFIRM_FIGHT_DRAGONICE:
		screen.RunFightDragonIce()
		return
	case CONFIRM_FIGHT_DRAGONACID:
		screen.RunFightDragonAcid()
		return
	case CONFIRM_FIGHT_DRAGONUNDEAD:
		screen.RunFightDragonUndead()
		return
	}
	nextMapper := map[ScreenStatus]ScreenStatus{
		RSLT_HUNT_RABBITS:              CONFIRM_HUNT_RABBITS,
		RSLT_FIGHT_GOBLIN:              CONFIRM_FIGHT_GOBLIN,
		RSLT_FIGHT_TROLL:               CONFIRM_FIGHT_TROLL,
		RSLT_FIGHT_WOLF:                CONFIRM_FIGHT_WOLF,
		RSLT_FIGHT_GIANT:               CONFIRM_FIGHT_GIANT,
		RSLT_FIGHT_DRAGONFIRE:          CONFIRM_FIGHT_DRAGONFIRE,
		RSLT_FIGHT_DRAGONICE:           CONFIRM_FIGHT_DRAGONICE,
		RSLT_FIGHT_DRAGONACID:          CONFIRM_FIGHT_DRAGONACID,
		RSLT_FIGHT_DRAGONUNDEAD:        CONFIRM_FIGHT_DRAGONUNDEAD,
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
		RSLT_RENAME_CHAR:               SEL_RENAME_CHAR,
		RSLT_SEL_ACT_CHAR:              SEL_ACTIVE_CHAR,
		RSLT_BUYITM:                    SEL_BUYITM,
		RSLT_BUYCHR:                    SEL_ACTIVE_CHAR,
		RSLT_SELLITM:                   SEL_SELLITM,
		RSLT_UPGITM:                    SEL_UPGITM,
	}
	if nextStatus, ok := nextMapper[screen.scrStatus]; ok {
		if screen.user.GetLocation() == loud.DEVELOP {
			screen.SetScreenStatus(SHW_LOCATION)
		} else if screen.user.GetLocation() == loud.FOREST && activeCharacter == nil {
			// move back to home in forest if no active character
			screen.SetScreenStatus(SHW_LOCATION)
		} else if nextStatus == CONFIRM_FIGHT_GIANT && activeCharacter.Special != loud.NO_SPECIAL {
			// go back to forest entrypoint when Special is not empty
			screen.SetScreenStatus(SHW_LOCATION)
		} else if nextStatus == SEL_ACTIVE_CHAR {
			screen.user.SetLocation(loud.HOME)
			screen.SetScreenStatus(nextStatus)
		} else {
			screen.SetScreenStatus(nextStatus)
		}
	} else {
		screen.SetScreenStatus(SHW_LOCATION)
	}
	screen.txFailReason = ""
	screen.Render()
}

func (screen *GameScreen) MoveToPrevStep() {
	activeCharacter := screen.user.GetActiveCharacter()

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
		RENAME_CHAR_ENT_NEWNAME:         SEL_RENAME_CHAR,
		RSLT_HUNT_RABBITS:               CONFIRM_HUNT_RABBITS,
		RSLT_FIGHT_GOBLIN:               CONFIRM_FIGHT_GOBLIN,
		RSLT_FIGHT_TROLL:                CONFIRM_FIGHT_TROLL,
		RSLT_FIGHT_WOLF:                 CONFIRM_FIGHT_WOLF,
		RSLT_FIGHT_GIANT:                CONFIRM_FIGHT_GIANT,
		RSLT_FIGHT_DRAGONFIRE:           CONFIRM_FIGHT_DRAGONFIRE,
		RSLT_FIGHT_DRAGONICE:            CONFIRM_FIGHT_DRAGONICE,
		RSLT_FIGHT_DRAGONACID:           CONFIRM_FIGHT_DRAGONACID,
		RSLT_FIGHT_DRAGONUNDEAD:         CONFIRM_FIGHT_DRAGONUNDEAD,
		RSLT_BUY_LOUD_TRDREQ_CREATION:   SHW_LOUD_BUY_TRDREQS,
		RSLT_FULFILL_BUY_LOUD_TRDREQ:    SHW_LOUD_BUY_TRDREQS,
		RSLT_SELL_LOUD_TRDREQ_CREATION:  SHW_LOUD_SELL_TRDREQS,
		RSLT_FULFILL_SELL_LOUD_TRDREQ:   SHW_LOUD_SELL_TRDREQS,
		RSLT_SELLITM_TRDREQ_CREATION:    SHW_SELLITM_TRDREQS,
		RSLT_FULFILL_SELLITM_TRDREQ:     SHW_SELLITM_TRDREQS,
		RSLT_BUYITM_TRDREQ_CREATION:     SHW_BUYITM_TRDREQS,
		RSLT_FULFILL_BUYITM_TRDREQ:      SHW_BUYITM_TRDREQS,
		RSLT_SELLCHR_TRDREQ_CREATION:    SHW_SELLCHR_TRDREQS,
		RSLT_FULFILL_SELLCHR_TRDREQ:     SHW_SELLCHR_TRDREQS,
		RSLT_BUYCHR_TRDREQ_CREATION:     SHW_BUYCHR_TRDREQS,
		RSLT_CANCEL_TRDREQ:              SHW_LOCATION,
		RSLT_FULFILL_BUYCHR_TRDREQ:      SHW_BUYCHR_TRDREQS,
		RSLT_RENAME_CHAR:                SEL_RENAME_CHAR,
		RSLT_SEL_ACT_CHAR:               SEL_ACTIVE_CHAR,
		RSLT_BUYITM:                     SEL_BUYITM,
		RSLT_BUYCHR:                     SEL_BUYCHR,
		RSLT_SELLITM:                    SEL_SELLITM,
		RSLT_UPGITM:                     SEL_UPGITM,

		HELP_ABOUT:             SHW_LOCATION,
		HELP_GAME_OBJECTIVE:    SHW_LOCATION,
		HELP_NAVIGATION:        SHW_LOCATION,
		HELP_PAGE_LAYOUT:       SHW_LOCATION,
		HELP_GAME_RULES:        SHW_LOCATION,
		HELP_HOW_IT_WORKS:      SHW_LOCATION,
		HELP_UPCOMING_RELEASES: SHW_LOCATION,
		HELP_SUPPORT:           SHW_LOCATION,
	}

	nxtStatus := SHW_LOCATION
	if nextStatus, ok := prevMapper[screen.scrStatus]; ok {
		nxtStatus = nextStatus
	}

	switch nxtStatus {
	case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL,
		CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL:
		// set loud value previously entered
		screen.inputText = screen.loudEnterValue
	case SHW_LOCATION:
		// move to home if it's somewhere else's entrypoint
		if screen.scrStatus == SHW_LOCATION {
			screen.user.SetLocation(loud.HOME)
		}
	case CONFIRM_FIGHT_GIANT:
		if activeCharacter.Special != loud.NO_SPECIAL {
			// go back to forest entrypoint when Special is not empty
			screen.SetScreenStatus(SHW_LOCATION)
		}
	}

	if screen.user.GetLocation() == loud.FOREST && activeCharacter == nil {
		// move back to home in forest if no active character
		screen.SetScreenStatus(SHW_LOCATION)
		screen.user.SetLocation(loud.HOME)
	}

	screen.SetScreenStatus(nxtStatus)
	screen.Render()
}

func (screen *GameScreen) HandleFirstClassInputKeys(input termbox.Event) bool {
	if input.Key == termbox.KeyEsc {
		switch screen.scrStatus {
		case CONFIRM_ENDGAME:
			screen.SetScreenStatus(SHW_LOCATION)
		default:
			screen.SetScreenStatus(CONFIRM_ENDGAME)
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
	case "B": // DEV ITEMS GET (troll toes, goblin ear, wolf tail and drops of 3 special dragons)
		screen.RunTxProcess(W8_DEV_GET_TEST_ITEMS, RSLT_DEV_GET_TEST_ITEMS, func() (string, error) {
			return loud.DevGetTestItems(screen.user)
		})
	case "L": // copy last txhash to CLIPBOARD
		clipboard.WriteAll(screen.user.GetLastTxHash())
	case "M": // copy user's cosmos address to CLIPBOARD
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
	} else if screen.user.GetLocation() == loud.HELP {
		switch screen.scrStatus {
		case SHW_LOCATION:
			return screen.HandleInputKeyHelpEntryPoint(input)
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
			newStatus := screen.scrStatus
			switch screen.scrStatus {
			case SHW_LOUD_BUY_TRDREQS:
				newStatus = CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL
			case SHW_LOUD_SELL_TRDREQS:
				newStatus = CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL
			case SHW_SELLITM_TRDREQS:
				newStatus = CR8_SELLITM_TRDREQ_SEL_ITEM
			case SHW_BUYITM_TRDREQS:
				newStatus = CR8_BUYITM_TRDREQ_SEL_ITEM
			case SHW_SELLCHR_TRDREQS:
				newStatus = CR8_SELLCHR_TRDREQ_SEL_CHR
			case SHW_BUYCHR_TRDREQS:
				newStatus = CR8_BUYCHR_TRDREQ_SEL_CHR
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
		case SEL_ACTIVE_CHAR:
			screen.activeLine = loud.GetIndexFromString(Key)
			screen.RunActiveCharacterSelect(screen.activeLine)
		case SEL_RENAME_CHAR:
			screen.activeLine = loud.GetIndexFromString(Key)
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.SetScreenStatus(RENAME_CHAR_ENT_NEWNAME)
			screen.inputText = ""
			screen.Render()
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
		screen.Render()
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
			screen.SetScreenStatus(CR8_SELLITM_TRDREQ_ENT_PYLVAL)
			screen.inputText = ""
			screen.Render()
		case CR8_BUYITM_TRDREQ_SEL_ITEM:
			if len(loud.WorldItemSpecs) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItSpec = loud.WorldItemSpecs[screen.activeLine]
			screen.SetScreenStatus(CR8_BUYITM_TRDREQ_ENT_PYLVAL)
			screen.inputText = ""
			screen.Render()
		case CR8_SELLCHR_TRDREQ_SEL_CHR:
			userCharacters := screen.user.InventoryCharacters()
			if len(userCharacters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = userCharacters[screen.activeLine]
			screen.SetScreenStatus(CR8_SELLCHR_TRDREQ_ENT_PYLVAL)
			screen.inputText = ""
			screen.Render()
		case CR8_BUYCHR_TRDREQ_SEL_CHR:
			if len(loud.WorldCharacterSpecs) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeChSpec = loud.WorldCharacterSpecs[screen.activeLine]
			screen.SetScreenStatus(CR8_BUYCHR_TRDREQ_ENT_PYLVAL)
			screen.inputText = ""
			screen.Render()
		case SEL_ACTIVE_CHAR:
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunActiveCharacterSelect(screen.activeLine)
		case SEL_RENAME_CHAR:
			characters := screen.user.InventoryCharacters()
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.SetScreenStatus(RENAME_CHAR_ENT_NEWNAME)
			screen.inputText = ""
			screen.Render()
		case SEL_BUYITM:
			items := loud.ShopItems
			if len(items) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeItem = items[screen.activeLine]
			screen.RunActiveItemBuy()
		case SEL_BUYCHR:
			characters := loud.ShopCharacters
			if len(characters) <= screen.activeLine || screen.activeLine < 0 {
				return false
			}
			screen.activeCharacter = characters[screen.activeLine]
			screen.RunActiveCharacterBuy()
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
	case termbox.KeyEsc:
		screen.MoveToPrevStep()
		return true
	case termbox.KeyBackspace2,
		termbox.KeyBackspace:

		log.Println("Pressed Backspace")
		lastIdx := len(screen.inputText) - 1
		if lastIdx < 0 {
			lastIdx = 0
		}
		screen.SetInputTextAndRender(screen.inputText[:lastIdx])
		return true
	case termbox.KeySpace:
		log.Println("Pressed Space")
		if screen.scrStatus == RENAME_CHAR_ENT_NEWNAME {
			screen.SetInputTextAndRender(screen.inputText + " ")
			return true
		}
		return false
	case termbox.KeyEnter:
		switch screen.scrStatus {
		case RENAME_CHAR_ENT_NEWNAME:
			screen.RunCharacterRename(screen.inputText)
		case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL:
			screen.SetScreenStatus(CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL)
			screen.loudEnterValue = screen.inputText
			screen.inputText = ""
			screen.Render()
		case CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL:
			screen.SetScreenStatus(W8_BUY_LOUD_TRDREQ_CREATION)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
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
			screen.SetScreenStatus(CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL)
			screen.Render()
			screen.loudEnterValue = screen.inputText
			screen.inputText = ""
		case CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL:
			screen.SetScreenStatus(W8_SELL_LOUD_TRDREQ_CREATION)
			screen.Render()
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
			screen.SetScreenStatus(W8_SELLITM_TRDREQ_CREATION)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
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
			screen.SetScreenStatus(W8_BUYITM_TRDREQ_CREATION)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
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
			screen.SetScreenStatus(W8_SELLCHR_TRDREQ_CREATION)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
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
			screen.SetScreenStatus(W8_BUYCHR_TRDREQ_CREATION)
			screen.pylonEnterValue = screen.inputText
			screen.SetInputTextAndRender("")
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
