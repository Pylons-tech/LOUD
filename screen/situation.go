package screen

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/pylons/x/pylons/handlers"
	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderUserSituation() {
	infoLines := []string{}
	desc := ""
	activeWeapon := screen.user.GetActiveWeapon()
	switch screen.scrStatus {
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
			} else {
				HP := uint64(activeCharacter.HP)
				MaxHP := uint64(activeCharacter.MaxHP)
				HP = min(HP+screen.BlockSince(activeCharacter.LastUpdate), MaxHP)
				if float32(HP) < float32(MaxHP)*.25 {
					desc = loud.Localize("home desc with low HP")
				}
			}
		}
	case SHW_LOUD_BUY_TRDREQS:
		infoLines = screen.renderTRTable(loud.BuyTrdReqs)
	case SHW_LOUD_SELL_TRDREQS:
		infoLines = screen.renderTRTable(loud.SellTrdReqs)
	case SHW_BUYITM_TRDREQS:
		infoLines = screen.renderITRTable(
			"Buy item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemBuyTrdReqs)
	case SHW_SELLITM_TRDREQS:
		infoLines = screen.renderITRTable(
			"Sell item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemSellTrdReqs)
	case SHW_SELLCHR_TRDREQS:
		infoLines = screen.renderITRTable(
			"Sell character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterSellTrdReqs)
	case SHW_BUYCHR_TRDREQS:
		infoLines = screen.renderITRTable(
			"Buy character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterBuyTrdReqs)
	case CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to get (should be integer value)" // TODO should add Localize
	case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL:
		desc = "Please enter loud amount to buy (should be integer value)" // TODO should add Localize
	case RENAME_CHAR_ENT_NEWNAME:
		desc = "Please enter new character's name - it's costing pylons per letter."
	case CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL:
		desc = "Please enter loud amount to sell (should be integer value)" // TODO should add Localize
	case CR8_SELLITM_TRDREQ_SEL_ITEM:
		infoLines = screen.renderITTable("Select item to sell", "Item", screen.user.InventoryItems())
	case CR8_SELLCHR_TRDREQ_SEL_CHR:
		infoLines = screen.renderITTable("Select character to sell", "Character", screen.user.InventoryCharacters())
	case CR8_SELLITM_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CR8_SELLCHR_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CR8_BUYITM_TRDREQ_SEL_ITEM:
		infoLines = screen.renderITTable("Select item to buy", "Item", loud.WorldItemSpecs)
	case CR8_BUYCHR_TRDREQ_SEL_CHR:
		infoLines = screen.renderITTable("Select character specs to get", "Character", loud.WorldCharacterSpecs)
	case CR8_BUYITM_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CR8_BUYCHR_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case SEL_ACTIVE_CHAR:
		infoLines = screen.renderITTable("Please select active character", "Character", screen.user.InventoryCharacters())
	case SEL_HEALTH_RESTORE_CHAR:
		infoLines = screen.renderITTable("Please select character to restore health", "Character", screen.user.InventoryCharacters())
	case SEL_RENAME_CHAR:
		infoLines = screen.renderITTable("Please select character to rename", "Character", screen.user.InventoryCharacters())
	case SEL_ACTIVE_WEAPON:
		infoLines = screen.renderITTable("Please select active weapon", "Item", screen.user.InventorySwords())
	case SEL_BUYITM:
		infoLines = screen.renderITTable("select buy item desc", "Item", loud.ShopItems)
	case SEL_SELLITM:
		infoLines = screen.renderITTable("select sell item desc", "Item", screen.user.InventorySellableItems())
	case SEL_UPGITM:
		infoLines = screen.renderITTable("select upgrade item desc", "Item", screen.user.InventoryUpgradableItems())
	case SEL_BUYCHR:
		infoLines = screen.renderITTable("select buy character desc", "Character", loud.ShopCharacters)
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
	}

	if screen.IsResultScreen() {
		desc = screen.TxResultSituationDesc()
	}

	if screen.IsWaitScreen() {
		desc = screen.TxWaitSituationDesc()
	}

	basicLines := strings.Split(desc, "\n")

	for _, line := range basicLines {
		infoLines = append(infoLines, loud.ChunkString(line, screen.screenSize.Width/2-4)...)
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

func (screen *GameScreen) TxResultSituationDesc() string {
	desc := ""
	resDescMap := map[ScreenStatus]string{
		RSLT_BUY_LOUD_TRDREQ_CREATION:  "loud buy request creation",
		RSLT_SELL_LOUD_TRDREQ_CREATION: "loud sell request creation",
		RSLT_SEL_ACT_CHAR:              "selecting active character",
		RSLT_HEALTH_RESTORE_CHAR:       "selecting character to restore health",
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
		desc = resDescMap[screen.scrStatus] + ": " + loud.Localize(screen.txFailReason)
	} else {
		switch screen.scrStatus {
		case RSLT_BUY_LOUD_TRDREQ_CREATION:
			desc = loud.Localize("loud buy request was successfully created")
			desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		case RSLT_SELL_LOUD_TRDREQ_CREATION:
			desc = loud.Localize("loud sell request was successfully created")
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
		case RSLT_HEALTH_RESTORE_CHAR:
			desc = loud.Localize("You have successfully restored character's health!")
		case RSLT_RENAME_CHAR:
			desc = loud.Sprintf("You have successfully updated character's name to %s!", screen.inputText)
		case RSLT_BUYITM:
			desc = loud.Sprintf("You have bought %s from the shop", formatItem(screen.activeItem))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RSLT_BUYCHR:
			desc = loud.Sprintf("You have bought %s from the shop", formatCharacter(screen.activeCharacter))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RSLT_HUNT_RABBITS:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon"}
			desc = loud.Sprintf("You did hunt rabbits and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
			switch len(respOutput) {
			case 0:
				desc += "\nYour character is dead during hunt accidently"
			case 2:
				if len(screen.activeItem.Name) > 0 {
					desc += "\nYou have lost your weapon accidently"
				}
			}
		case RSLT_FIGHT_GOBLIN:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon", loud.GOBLIN_EAR}
			desc = loud.Sprintf("You did fight with goblin and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
			switch len(respOutput) {
			case 0:
				desc += "\nYour character is dead during fighting goblin accidently"
			case 2:
				desc += "\nYou have lost your weapon accidently"
			case 4:
				desc += fmt.Sprintf("\nYou got bonus item called %s", loud.GOBLIN_EAR)
			}
		case RSLT_FIGHT_TROLL:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon", loud.TROLL_TOES}
			desc = loud.Sprintf("You did fight with troll and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
			switch len(respOutput) {
			case 0:
				desc += "\nYour character is dead during fighting troll accidently"
			case 2:
				desc += "\nYou have lost your weapon accidently"
			case 4:
				desc += fmt.Sprintf("\nYou got bonus item called %s", loud.TROLL_TOES)
			}
		case RSLT_FIGHT_WOLF:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon", loud.WOLF_TAIL}
			desc = loud.Sprintf("You did fight with wolf and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
			switch len(respOutput) {
			case 0:
				desc += "\nYour character is dead during fighting wolf accidently"
			case 2:
				desc += "\nYou have lost your weapon accidently"
			case 4:
				desc += fmt.Sprintf("\nYou got bonus item called %s", loud.WOLF_TAIL)
			}
		case RSLT_FIGHT_GIANT:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon"}
			desc = loud.Sprintf("You did fight with giant and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
			switch len(respOutput) {
			case 0:
				desc += "\nYour character is dead during fighting wolf accidently"
			case 2:
				desc += "\nYou have lost your weapon accidently"
			}
		case RSLT_BUY_GOLD_WITH_PYLONS:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			json.Unmarshal(screen.txResult, &respOutput)
			earnedAmount := int64(0)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			desc = loud.Sprintf("Bought gold with pylons. Amount is %d.", earnedAmount)
		case RSLT_DEV_GET_TEST_ITEMS:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			json.Unmarshal(screen.txResult, &respOutput)
			resultTexts := []string{loud.WOLF_TAIL, loud.TROLL_TOES, loud.GOBLIN_EAR}
			desc = loud.Sprintf("Finished getting developer test items. Detailed result %+v", resultTexts[:len(respOutput)])
		case RSLT_GET_PYLONS:
			desc = loud.Localize("You got extra pylons for loud game")
		case RSLT_SWITCH_USER:
			desc = loud.Sprintf("You switched user to %s", screen.user.GetUserName())
		case RSLT_CREATE_COOKBOOK:
			desc = loud.Localize("You created a new cookbook for a new game build")
		case RSLT_SELLITM:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
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
			desc = loud.Localize("you have sold loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.sellLoudDesc(request.Amount, request.Total)
		case RSLT_FULFILL_SELL_LOUD_TRDREQ:
			request := screen.activeTrdReq
			desc = loud.Localize("you have bought loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.buyLoudDesc(request.Amount, request.Total)
		case RSLT_FULFILL_SELLITM_TRDREQ:
			request := screen.activeItemTrdReq.(loud.ItemSellTrdReq)
			desc = loud.Localize("you have bought item successfully from item/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_SELLCHR_TRDREQ:
			request := screen.activeItemTrdReq.(loud.CharacterSellTrdReq)
			desc = loud.Localize("you have bought character successfully from character/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_BUYITM_TRDREQ:
			request := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
			desc = loud.Localize("you have sold item successfully from item/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_BUYCHR_TRDREQ:
			request := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
			desc = loud.Localize("you have sold character successfully from character/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		}
	}
	return desc
}

func (screen *GameScreen) TxWaitSituationDesc() string {
	desc := ""
	activeWeapon := screen.user.GetActiveWeapon()
	W8_TO_END := "\n" + loud.Localize("Please wait for a moment to finish the process")
	switch screen.scrStatus {
	case W8_RENAME_CHAR:
		desc = loud.Sprintf("You are now waiting to rename character from %s to %s.", screen.activeCharacter.Name, screen.inputText)
	case W8_BUY_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for loud buy request creation")
		desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_SELL_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for loud sell request creation")
		desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_BUYITM:
		desc = loud.Sprintf("You are now buying %s at the shop", formatItem(screen.activeItem))
		desc += W8_TO_END
	case W8_BUYCHR:
		desc = loud.Sprintf("You are now buying %s at the shop", formatCharacter(screen.activeCharacter))
		desc += W8_TO_END
	case W8_HUNT_RABBITS:
		if activeWeapon != nil {
			desc = loud.Sprintf("You are now hunting rabbits with %s", formatItem(*activeWeapon))
		} else {
			desc = loud.Sprintf("You are now hunting rabbits without weapon")
		}
		desc += W8_TO_END
	case W8_FIGHT_GIANT:
		desc = loud.Sprintf("You are now fighting with giant with %s", formatItem(*activeWeapon))
	case W8_FIGHT_GOBLIN:
		desc = loud.Sprintf("You are now fighting with goblin with %s", formatItem(*activeWeapon))
	case W8_FIGHT_TROLL:
		desc = loud.Sprintf("You are now fighting with troll with %s", formatItem(*activeWeapon))
	case W8_FIGHT_WOLF:
		desc = loud.Sprintf("You are now fighting with wolf with %s", formatItem(*activeWeapon))
	case W8_BUY_GOLD_WITH_PYLONS:
		desc = loud.Localize("Buying gold with pylon")
		desc += W8_TO_END
	case W8_DEV_GET_TEST_ITEMS:
		desc = loud.Localize("Getting dev test items from pylon")
		desc += W8_TO_END
	case W8_HEALTH_RESTORE_CHAR:
		desc = loud.Localize("Waiting for Health restoring")
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
		desc = loud.Sprintf("you are now selling loud for pylon at %.4f.", request.Price)
		desc += screen.sellLoudDesc(request.Amount, request.Total)
	case W8_FULFILL_SELL_LOUD_TRDREQ:
		request := screen.activeTrdReq
		desc = loud.Sprintf("you are now buying loud from pylon at %.4f.", request.Price)
		desc += screen.buyLoudDesc(request.Amount, request.Total)
	}
	desc += "\n"
	onColor := screen.colorFunc(fmt.Sprintf("%v+B:%v", 117, bgcolor))
	return desc + onColor("......")
}
