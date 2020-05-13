package screen

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/handlers"
	"github.com/ahmetb/go-cursor"
)

func devDetailedResultDesc(res []string) string {
	resT := []string{}
	for _, it := range res {
		resT = append(resT, loud.Localize(it))
	}
	return fmt.Sprintf("\n%s:\n  %s\n", loud.Localize("Detailed result"), strings.Join(resT, "\n  "))
}

func (screen *GameScreen) GetTxResponseOutput() (int64, []handlers.ExecuteRecipeSerialize) {
	respOutput := []handlers.ExecuteRecipeSerialize{}
	earnedAmount := int64(0)
	json.Unmarshal(screen.txResult, &respOutput)
	if len(respOutput) > 0 {
		earnedAmount = respOutput[0].Amount
	}
	return earnedAmount, respOutput
}

func (screen *GameScreen) renderUserSituation() {

	// situation box start point (x, y)
	scrBox := screen.GetSituationBox()
	x := scrBox.X
	y := scrBox.Y
	w := scrBox.W
	h := scrBox.H

	infoLines := []string{}
	tableLines := []string{}
	desc := ""
	descfont := REGULAR
	activeWeapon := screen.user.GetActiveWeapon()
	switch screen.scrStatus {
	case CONFIRM_ENDGAME:
		desc = loud.Localize("Are you really gonna end game?")
	case SHW_LOCATION:
		locationDescMap := map[loud.UserLocation]string{
			loud.HOME:     loud.Localize("home desc"),
			loud.FOREST:   loud.Localize("forest desc"),
			loud.SHOP:     loud.Localize("shop desc"),
			loud.PYLCNTRL: loud.Localize("pylons central desc"),
			loud.SETTINGS: loud.Localize("settings desc"),
			loud.DEVELOP:  loud.Localize("develop desc"),
		}
		desc = locationDescMap[screen.user.GetLocation()]
		if screen.user.GetLocation() == loud.HOME {
			activeCharacter := screen.user.GetActiveCharacter()
			if activeCharacter == nil {
				desc = loud.Localize("home desc without character")
			} else if screen.user.GetPylonAmount() == 0 {
				desc = loud.Localize("home desc without pylon")
			}
		}
	case SHW_LOUD_BUY_TRDREQS:
		infoLines, tableLines = screen.renderTRTable(loud.BuyTrdReqs, w)
	case SHW_LOUD_SELL_TRDREQS:
		infoLines, tableLines = screen.renderTRTable(loud.SellTrdReqs, w)
	case SHW_BUYITM_TRDREQS:
		infoLines, tableLines = screen.renderITRTable(
			"Buy item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemBuyTrdReqs,
			w)
	case SHW_SELLITM_TRDREQS:
		infoLines, tableLines = screen.renderITRTable(
			"Sell item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemSellTrdReqs,
			w)
	case SHW_SELLCHR_TRDREQS:
		infoLines, tableLines = screen.renderITRTable(
			"Sell character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterSellTrdReqs,
			w)
	case SHW_BUYCHR_TRDREQS:
		infoLines, tableLines = screen.renderITRTable(
			"Buy character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterBuyTrdReqs,
			w)
	case CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to get (should be integer value)")
	case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL:
		desc = loud.Localize("Please enter gold amount to buy (should be integer value)")
	case RENAME_CHAR_ENT_NEWNAME:
		desc = loud.Localize("Please enter new character's name - it's costing pylons per letter.")
	case CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL:
		desc = loud.Localize("Please enter gold amount to sell (should be integer value)")

	case CR8_SELLITM_TRDREQ_SEL_ITEM:
		infoLines, tableLines = screen.renderITTable(
			"Select item to sell",
			"Item",
			screen.user.InventoryItems(),
			w)
	case CR8_SELLCHR_TRDREQ_SEL_CHR:
		infoLines, tableLines = screen.renderITTable(
			"Select character to sell",
			"Character",
			screen.user.InventoryCharacters(),
			w)
	case CR8_SELLITM_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_SELLCHR_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_BUYITM_TRDREQ_SEL_ITEM:
		infoLines, tableLines = screen.renderITTable(
			"Select item to buy",
			"Item",
			loud.WorldItemSpecs,
			w)
	case CR8_BUYCHR_TRDREQ_SEL_CHR:
		infoLines, tableLines = screen.renderITTable(
			"Select character specs to get",
			"Character",
			loud.WorldCharacterSpecs,
			w)
	case CR8_BUYITM_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_BUYCHR_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case SEL_ACTIVE_CHAR:
		infoLines, tableLines = screen.renderITTable(
			"Please select active character",
			"Character",
			screen.user.InventoryCharacters(),
			w)
	case SEL_RENAME_CHAR:
		infoLines, tableLines = screen.renderITTable(
			"Please select character to rename",
			"Character",
			screen.user.InventoryCharacters(),
			w)
	case SEL_ACTIVE_WEAPON:
		infoLines, tableLines = screen.renderITTable(
			"Please select active weapon",
			"Item",
			screen.user.InventorySwords(),
			w)
	case SEL_BUYITM:
		infoLines, tableLines = screen.renderITTable(
			"select buy item desc",
			"Item",
			loud.ShopItems,
			w)
	case SEL_SELLITM:
		infoLines, tableLines = screen.renderITTable(
			"select sell item desc",
			"Item",
			screen.user.InventorySellableItems(), w)
	case SEL_UPGITM:
		infoLines, tableLines = screen.renderITTable(
			"select upgrade item desc",
			"Item",
			screen.user.InventoryUpgradableItems(), w)
	case SEL_BUYCHR:
		infoLines, tableLines = screen.renderITTable(
			"select buy character desc",
			"Character",
			loud.ShopCharacters, w)
	case CONFIRM_HUNT_RABBITS:
		if activeWeapon != nil {
			desc = loud.Localize("rabbits with a sword outcome")
			desc += carryItemDesc(activeWeapon)
		} else {
			desc = loud.Localize("rabbits without sword outcome")
		}
	case CONFIRM_FIGHT_GOBLIN:
		desc = loud.Localize("goblin outcome")
		desc += carryItemDesc(activeWeapon)
	case CONFIRM_FIGHT_WOLF:
		desc = loud.Localize("wolf outcome")
		desc += carryItemDesc(activeWeapon)
	case CONFIRM_FIGHT_TROLL:
		desc = loud.Localize("troll outcome")
		desc += carryItemDesc(activeWeapon)
	case CONFIRM_FIGHT_GIANT:
		desc = loud.Localize("giant outcome")
		desc += carryItemDesc(activeWeapon)
	case CONFIRM_FIGHT_DRAGONFIRE:
		desc = loud.Localize("fire dragon outcome")
		desc += carryItemDesc(activeWeapon)
	case CONFIRM_FIGHT_DRAGONICE:
		desc = loud.Localize("ice dragon outcome")
		desc += carryItemDesc(activeWeapon)
	case CONFIRM_FIGHT_DRAGONACID:
		desc = loud.Localize("acid dragon outcome")
		desc += carryItemDesc(activeWeapon)
	case CONFIRM_FIGHT_DRAGONUNDEAD:
		desc = loud.Localize("undead dragon outcome")
		desc += carryItemDesc(activeWeapon)
	}

	if screen.IsResultScreen() {
		desc, descfont = screen.TxResultSituationDesc()
	}

	if screen.IsWaitScreen() {
		infoLines, tableLines = screen.TxWaitSituationDesc(w)
	}

	basicLines := strings.Split(desc, "\n")

	for _, line := range basicLines {
		chunkedlines := loud.ChunkString(line, screen.leftInnerWidth()-2)
		if descfont == REGULAR {
			infoLines = append(infoLines, chunkedlines...)
		} else {
			chunkedColorfulLines := []string{}
			for _, chli := range chunkedlines {
				chunkedColorfulLines = append(chunkedColorfulLines, screen.getFont(descfont)(fillSpace(chli, w)))
			}
			tableLines = append(chunkedColorfulLines, tableLines...)
		}
	}

	fmtFunc := screen.regularFont()
	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x),
			fmtFunc(fillSpace(line, w))))
		if index+2 > int(screen.Height()) {
			break
		}
	}
	infoLen := len(infoLines)

	for index, line := range tableLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+infoLen+index, x),
			line))
		if index+2 > int(screen.Height()) {
			break
		}
	}
	totalLen := infoLen + len(tableLines)

	screen.drawFill(x, y+totalLen, w, h-totalLen-1)
}

func (screen *GameScreen) TxResultSituationDesc() (string, FontType) {
	desc := ""
	font := REGULAR
	resDescMap := map[ScreenStatus]string{
		RSLT_BUY_LOUD_TRDREQ_CREATION:  "loud buy request creation",
		RSLT_SELL_LOUD_TRDREQ_CREATION: "loud sell request creation",
		RSLT_SEL_ACT_CHAR:              "selecting active character",
		RSLT_RENAME_CHAR:               "renaming character",
		RSLT_SEL_ACT_WEAPON:            "selecting active weapon",
		RSLT_BUYITM:                    "buy item",
		RSLT_BUYCHR:                    "buy character",
		RSLT_HUNT_RABBITS:              "hunt rabbits",
		RSLT_BUY_GOLD_WITH_PYLONS:      "buy gold with pylons",
		RSLT_GET_PYLONS:                "get pylon",
		RSLT_SWITCH_USER:               "switch user",
		RSLT_CREATE_COOKBOOK:           "create cookbook",
		RSLT_SELLITM:                   "sell item",
		RSLT_UPGITM:                    "upgrade item",
		RSLT_SELLITM_TRDREQ_CREATION:   "sell item request creation",
		RSLT_BUYITM_TRDREQ_CREATION:    "buy item request creation",
		RSLT_SELLCHR_TRDREQ_CREATION:   "sell character request creation",
		RSLT_BUYCHR_TRDREQ_CREATION:    "buy character request creation",
		RSLT_CANCEL_TRDREQ:             "cancel trade",
		RSLT_FULFILL_BUY_LOUD_TRDREQ:   "sell loud", // for fullfill direction is reversed
		RSLT_FULFILL_SELL_LOUD_TRDREQ:  "buy loud",
		RSLT_FULFILL_SELLITM_TRDREQ:    "buy item",
		RSLT_FULFILL_SELLCHR_TRDREQ:    "buy character",
		RSLT_FULFILL_BUYITM_TRDREQ:     "sell item",
		RSLT_FULFILL_BUYCHR_TRDREQ:     "sell character",
	}
	if screen.txFailReason != "" {
		desc = loud.Localize(resDescMap[screen.scrStatus]+" failure reason") + ": " + loud.Localize(screen.txFailReason)
		font = RED_BOLD
	} else {
		switch screen.scrStatus {
		case RSLT_BUY_LOUD_TRDREQ_CREATION:
			desc = loud.Localize("gold buy request was successfully created")
			desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		case RSLT_SELL_LOUD_TRDREQ_CREATION:
			desc = loud.Localize("gold sell request was successfully created")
			desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		case RSLT_SEL_ACT_CHAR:
			if screen.user.GetActiveCharacter() == nil {
				desc = loud.Localize("You have successfully unset the active character!")
			} else {
				desc = loud.Localize("You have successfully set the active character!")
			}
		case RSLT_SEL_ACT_WEAPON:
			if screen.user.GetActiveWeapon() == nil {
				desc = loud.Localize("You have successfully unset the active weapon!")
			} else {
				desc = loud.Localize("You have successfully set the active weapon!")
			}
		case RSLT_RENAME_CHAR:
			desc = loud.Sprintf("You have successfully updated character's name to %s!", screen.inputText)
		case RSLT_BUYITM:
			desc = loud.Sprintf("You have bought %s from the shop", formatItem(screen.activeItem))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RSLT_BUYCHR:
			desc = loud.Sprintf("You have bought %s from Pylons Central", formatCharacter(screen.activeCharacter))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RSLT_HUNT_RABBITS:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("Your character is dead while following rabbits accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did hunt rabbits and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon"}
				// desc = devDetailedResultDesc(resultTexts[:resLen])
				if resLen == 2 && screen.user.GetLastTxMetaData() == loud.RCP_HUNT_RABBITS_YESWORD {
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				}
			}
		case RSLT_FIGHT_GOBLIN:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by goblin accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with goblin and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon", loud.GOBLIN_EAR}
				// desc += devDetailedResultDesc(resultTexts[:resLen])

				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.GOBLIN_EAR)
					desc += "\n"
					desc += loud.Sprintf("You can make silver sword with Goblin ear at the shop!")
					font = GREEN
				}
			}
		case RSLT_FIGHT_TROLL:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by troll accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with troll and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon", loud.TROLL_TOES}
				// desc = devDetailedResultDesc(resultTexts[:resLen])
				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.TROLL_TOES)
					desc += "\n"
					desc += loud.Sprintf("You can make iron sword with Troll toes at the shop!")
					font = GREEN
				}
			}
		case RSLT_FIGHT_WOLF:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by wolf accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with wolf and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon", loud.WOLF_TAIL}
				// desc = devDetailedResultDesc(resultTexts[:resLen])
				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.WOLF_TAIL)
					desc += "\n"
					desc += loud.Sprintf("You can make bronze sword with Wolf tail at the shop!")
					font = GREEN
				}
			}
		case RSLT_FIGHT_GIANT:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by giant accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with giant and earned %d.", earnedAmount)
				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				case 3:
					activeCharacter := screen.user.GetActiveCharacter()
					if activeCharacter.Special != loud.NO_SPECIAL { // Got special from this fight
						desc += loud.Sprintf("You got %s (special) from the giant!!", formatSpecial(activeCharacter.Special))
						desc += "\n"
						desc += loud.Sprintf("You can now fight with %s with this character!", formatSpecialDragon(activeCharacter.Special))
						font = GREEN
					}
				}
			}
		case RSLT_FIGHT_DRAGONFIRE:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by fire dragon accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with fire dragon and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon", loud.DROP_DRAGONFIRE}
				// desc = devDetailedResultDesc(resultTexts[:resLen])
				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.DROP_DRAGONFIRE)
					desc += "\n"
					desc += loud.Sprintf("Once you have drops from 3 special dragons, you can create angel sword at the shop!")
					font = GREEN
				}
			}
		case RSLT_FIGHT_DRAGONICE:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by ice dragon accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with ice dragon and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon", loud.DROP_DRAGONICE}
				// desc = devDetailedResultDesc(resultTexts[:resLen])
				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.DROP_DRAGONICE)
					desc += "\n"
					desc += loud.Sprintf("Once you have drops from 3 special dragons, you can create angel sword at the shop!")
					font = GREEN
				}
			}
		case RSLT_FIGHT_DRAGONACID:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by acid dragon accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with acid dragon and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon", loud.DROP_DRAGONACID}
				// desc = devDetailedResultDesc(resultTexts[:resLen])
				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.DROP_DRAGONACID)
					desc += "\n"
					desc += loud.Sprintf("Once you have drops from 3 special dragons, you can create angel sword at the shop!")
					font = GREEN
				}
			}
		case RSLT_FIGHT_DRAGONUNDEAD:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("You were killed by undead dragon accidently")
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with undead dragon and earned %d.", earnedAmount)
				// resultTexts := []string{"gold", "character", "weapon"}
				// desc = devDetailedResultDesc(resultTexts[:resLen])
				switch resLen {
				case 2:
					desc += loud.Sprintf("You have lost your weapon accidently")
					font = YELLOW
				}
			}
		case RSLT_BUY_GOLD_WITH_PYLONS:
			earnedAmount, _ := screen.GetTxResponseOutput()
			desc = loud.Sprintf("Bought gold with pylons. Amount is %d.", earnedAmount)
		case RSLT_DEV_GET_TEST_ITEMS:
			_, respOutput := screen.GetTxResponseOutput()
			resultTexts := []string{
				loud.WOLF_TAIL,
				loud.TROLL_TOES,
				loud.GOBLIN_EAR,
				loud.DROP_DRAGONFIRE,
				loud.DROP_DRAGONICE,
				loud.DROP_DRAGONACID,
				"Ruppell's Fox",
				"Gentoo penguin",
				"Colorado River toad",
			}
			desc = loud.Sprintf("Finished getting developer test items.")
			desc += devDetailedResultDesc(resultTexts[:len(respOutput)])
		case RSLT_GET_PYLONS:
			desc = loud.Localize("You got extra pylons for LOUD game")
		case RSLT_SWITCH_USER:
			desc = loud.Sprintf("You switched user to %s", screen.user.GetUserName())
		case RSLT_CREATE_COOKBOOK:
			desc = loud.Localize("You created a new cookbook for a new game build")
		case RSLT_SELLITM:
			earnedAmount, _ := screen.GetTxResponseOutput()
			desc = loud.Sprintf("You sold %s for %d gold.", formatItem(screen.activeItem), earnedAmount)
		case RSLT_UPGITM:
			desc = loud.Sprintf("You have upgraded %s to get better hunt result", screen.activeItem.Name)
		case RSLT_SELLITM_TRDREQ_CREATION:
			desc = loud.Localize("item sell request was successfully created")
			desc += screen.sellItemDesc(screen.activeItem, screen.pylonEnterValue)
		case RSLT_SELLCHR_TRDREQ_CREATION:
			desc = loud.Localize("character sell request was successfully created")
			desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
		case RSLT_BUYITM_TRDREQ_CREATION:
			desc = loud.Localize("item buy request was successfully created")
			desc += screen.buyItemSpecDesc(screen.activeItSpec, screen.pylonEnterValue)
		case RSLT_BUYCHR_TRDREQ_CREATION:
			desc = loud.Localize("character buy request was successfully created")
			desc += screen.buyCharacterSpecDesc(screen.activeChSpec, screen.pylonEnterValue)
		case RSLT_CANCEL_TRDREQ:
			desc = loud.Localize("successfully cancelled trade request")
		case RSLT_FULFILL_BUY_LOUD_TRDREQ:
			request := screen.activeTrdReq
			desc = loud.Sprintf("you have sold gold successfully from coin market at %.4f", request.Price)
			desc += screen.sellLoudDesc(request.Amount, request.Total)
		case RSLT_FULFILL_SELL_LOUD_TRDREQ:
			request := screen.activeTrdReq
			desc = loud.Sprintf("you have bought gold successfully from coin market at %.4f", request.Price)
			desc += screen.buyLoudDesc(request.Amount, request.Total)
		case RSLT_FULFILL_SELLITM_TRDREQ:
			request := screen.activeItemTrdReq.(loud.ItemSellTrdReq)
			desc = loud.Localize("you have bought item successfully from item/pylon market")
			desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_SELLCHR_TRDREQ:
			request := screen.activeItemTrdReq.(loud.CharacterSellTrdReq)
			desc = loud.Localize("you have bought character successfully from character/pylon market")
			desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_BUYITM_TRDREQ:
			request := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
			desc = loud.Localize("you have sold item successfully from item/pylon market")
			desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_BUYCHR_TRDREQ:
			request := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
			desc = loud.Localize("you have sold character successfully from character/pylon market")
			desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		}
	}
	return desc, font
}

func (screen *GameScreen) TxWaitSituationDesc(width int) ([]string, []string) {
	desc := ""
	activeWeapon := screen.user.GetActiveWeapon()
	W8_TO_END := "\n" + loud.Localize("Please wait for a moment to finish the process")
	switch screen.scrStatus {
	case W8_RENAME_CHAR:
		desc = loud.Sprintf("You are now waiting to rename character from %s to %s.", screen.activeCharacter.Name, screen.inputText)
	case W8_BUY_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for gold buy request creation")
		desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_SELL_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for gold sell request creation")
		desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_BUYITM:
		desc = loud.Sprintf("You are now buying %s at the shop", formatItem(screen.activeItem))
		desc += W8_TO_END
	case W8_BUYCHR:
		desc = loud.Sprintf("You are now buying %s at the shop", formatCharacter(screen.activeCharacter))
		desc += W8_TO_END
	case W8_HUNT_RABBITS:
		if activeWeapon != nil {
			desc = loud.Sprintf("You are now hunting rabbits with %s", formatItemP(activeWeapon))
		} else {
			desc = loud.Sprintf("You are now hunting rabbits without weapon")
		}
		desc += W8_TO_END
	case W8_FIGHT_GIANT:
		desc = loud.Sprintf("You are now fighting with giant with %s", formatItemP(activeWeapon))
	case W8_FIGHT_DRAGONFIRE:
		desc = loud.Sprintf("You are now fighting with fire dragon with %s", formatItemP(activeWeapon))
	case W8_FIGHT_DRAGONICE:
		desc = loud.Sprintf("You are now fighting with ice dragon with %s", formatItemP(activeWeapon))
	case W8_FIGHT_DRAGONACID:
		desc = loud.Sprintf("You are now fighting with acid dragon with %s", formatItemP(activeWeapon))
	case W8_FIGHT_DRAGONUNDEAD:
		desc = loud.Sprintf("You are now fighting with undead dragon with %s", formatItemP(activeWeapon))
	case W8_FIGHT_GOBLIN:
		desc = loud.Sprintf("You are now fighting with goblin with %s", formatItemP(activeWeapon))
	case W8_FIGHT_TROLL:
		desc = loud.Sprintf("You are now fighting with troll with %s", formatItemP(activeWeapon))
	case W8_FIGHT_WOLF:
		desc = loud.Sprintf("You are now fighting with wolf with %s", formatItemP(activeWeapon))
	case W8_BUY_GOLD_WITH_PYLONS:
		desc = loud.Localize("Buying gold with pylon")
		desc += W8_TO_END
	case W8_DEV_GET_TEST_ITEMS:
		desc = loud.Localize("Getting dev test items from pylon")
		desc += W8_TO_END
	case W8_GET_PYLONS:
		desc = loud.Localize("You are waiting for getting pylons process")
	case W8_SWITCH_USER:
		desc = loud.Localize("You are waiting for switching to new user")
	case W8_CREATE_COOKBOOK:
		desc = loud.Localize("You are waiting for creating cookbook")
	case W8_SELLITM:
		desc = loud.Sprintf("You are now selling %s for gold", formatItem(screen.activeItem))
		desc += W8_TO_END
	case W8_UPGITM:
		desc = loud.Sprintf("You are now upgrading %s", loud.Localize(screen.activeItem.Name))
		desc += W8_TO_END
	case W8_SELLITM_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for item sell request creation")
		desc += screen.sellItemDesc(screen.activeItem, screen.pylonEnterValue)
	case W8_SELLCHR_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for character sell request creation")
		desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
	case W8_BUYITM_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for item buy request creation")
		desc += screen.buyItemSpecDesc(screen.activeItSpec, screen.pylonEnterValue)
	case W8_BUYCHR_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for character buy request creation")
		desc += screen.buyCharacterSpecDesc(screen.activeChSpec, screen.pylonEnterValue)
	case W8_CANCEL_TRDREQ:
		desc = loud.Localize("You are now waiting for cancelling one of your trades")
	// For FULFILL trades, msg should be reversed, since user is opposite
	case W8_FULFILL_SELLITM_TRDREQ:
		request := screen.activeItemTrdReq.(loud.ItemSellTrdReq)
		desc = loud.Sprintf("You are now buying item at %d", request.Price)
		desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_SELLCHR_TRDREQ:
		request := screen.activeItemTrdReq.(loud.CharacterSellTrdReq)
		desc = loud.Sprintf("you are now buying character at %d.", request.Price)
		desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUYITM_TRDREQ:
		request := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
		desc = loud.Sprintf("you are now selling item at %d.", request.Price)
		desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUYCHR_TRDREQ:
		request := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
		desc = loud.Sprintf("you are now selling character at %d.", request.Price)
		desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUY_LOUD_TRDREQ:
		request := screen.activeTrdReq
		desc = loud.Sprintf("Making pylons from gold")
		desc += screen.sellLoudDesc(request.Amount, request.Total)
	case W8_FULFILL_SELL_LOUD_TRDREQ:
		request := screen.activeTrdReq
		desc = loud.Sprintf("Making gold from pylons")
		desc += screen.buyLoudDesc(request.Amount, request.Total)
	}
	desc += "\n"
	return strings.Split(desc, "\n"), []string{
		screen.blinkBlueBoldFont()(fillSpace("......", width)),
	}
}
