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
	switch screen.scrStatus {
	case SHW_LOCATION:
		locationDescMap := map[loud.UserLocation]string{
			loud.HOME:     loud.Localize("home desc"),
			loud.FOREST:   loud.Localize("forest desc"),
			loud.SHOP:     loud.Localize("shop desc"),
			loud.MARKET:   loud.Localize("market desc"),
			loud.SETTINGS: loud.Localize("settings desc"),
			loud.DEVELOP:  loud.Localize("develop desc"),
		}
		desc = locationDescMap[screen.user.GetLocation()]
	case SHW_LOUD_BUY_TRDREQS:
		infoLines = screen.renderTRTable(loud.BuyTradeRequests)
	case SHW_LOUD_SELL_TRDREQS:
		infoLines = screen.renderTRTable(loud.SellTradeRequests)
	case SHW_BUYITM_TRDREQS:
		infoLines = screen.renderITRTable(
			"Buy item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemBuyTradeRequests)
	case SHW_SELLITM_TRDREQS:
		infoLines = screen.renderITRTable(
			"Sell item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemSellTradeRequests)
	case SHW_SELLCHR_TRDREQS:
		infoLines = screen.renderITRTable(
			"Sell character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterSellTradeRequests)
	case SHW_BUYCHR_TRDREQS:
		infoLines = screen.renderITRTable(
			"Buy character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterBuyTradeRequests)
	case CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL:
		desc = "Please enter pylon amount to get (should be integer value)" // TODO should add Localize
	case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL:
		desc = "Please enter loud amount to buy (should be integer value)" // TODO should add Localize
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
	case SEL_DEFAULT_CHAR:
		infoLines = screen.renderITTable("Please select default character", "Character", screen.user.InventoryCharacters())
	case SEL_HEALTH_RESTORE_CHAR:
		infoLines = screen.renderITTable("Please select character to restore health", "Character", screen.user.InventoryCharacters())
	case SEL_DEFAULT_WEAPON:
		infoLines = screen.renderITTable("Please select default weapon", "Item", screen.user.InventorySwords())
	case SEL_BUYITM:
		infoLines = screen.renderITTable("select buy item desc", "Item", loud.ShopItems)
	case SEL_BUYCHR:
		infoLines = screen.renderITTable("select buy character desc", "Character", loud.ShopCharacters)
	case SEL_SELLITM:
		infoLines = screen.renderITTable("select sell item desc", "Item", screen.user.InventoryItems())
	case SEL_HUNT_ITEM:
		infoLines = screen.renderITTable("select hunt item desc", "Item", screen.user.InventorySwords())
	case SEL_FIGHT_GOBLIN_ITEM:
		infoLines = screen.renderITTable("select fight goblin item desc", "Item", screen.user.InventorySwords())
	case SEL_FIGHT_WOLF_ITEM:
		infoLines = screen.renderITTable("select fight wolf item desc", "Item", screen.user.InventorySwords())
	case SEL_FIGHT_TROLL_ITEM:
		infoLines = screen.renderITTable("select fight troll item desc", "Item", screen.user.InventorySwords())
	case SEL_FIGHT_GIANT_ITEM:
		infoLines = screen.renderITTable("select fight giant item desc", "Item", screen.user.InventoryIronSwords())
	case SEL_UPGITM:
		infoLines = screen.renderITTable("select upgrade item desc", "Item", screen.user.InventoryUpgradableItems())
	}

	if strings.HasPrefix(string(screen.scrStatus), "RSLT_") {
		desc = screen.TxResultSituationDesc()
	}

	if strings.HasPrefix(string(screen.scrStatus), "W8_") {
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
		RSLT_SEL_DEF_CHAR:              "selecting default character",
		RSLT_HEALTH_RESTORE_CHAR:       "selecting character to restore health",
		RSLT_SEL_DEF_WEAPON:            "selecting default weapon",
		RSLT_BUYITM:                    "buy item",
		RSLT_BUYCHR:                    "buy character",
		RSLT_HUNT:                      "hunt",
		RSLT_GET_INITIAL_COIN:          "get initial coin",
		RSLT_GET_PYLONS:                "get pylon",
		RSLT_SWITCH_USER:               "switch user",
		RSLT_CREATE_COOKBOOK:           "create cookbook",
		RSLT_SELLITM:                   "sell item",
		RSLT_UPGITM:                    "upgrade item",
		RSLT_SELLITM_TRDREQ_CREATION:   "sell item request creation",
		RSLT_BUYITM_TRDREQ_CREATION:    "buy item request creation",
		RSLT_SELLCHR_TRDREQ_CREATION:   "sell character request creation",
		RSLT_BUYCHR_TRDREQ_CREATION:    "buy character request creation",
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
		case RSLT_SEL_DEF_CHAR:
			desc = loud.Localize("You have successfully set default character!")
		case RSLT_HEALTH_RESTORE_CHAR:
			desc = loud.Localize("You have successfully restored character's health!")
		case RSLT_SEL_DEF_WEAPON:
			desc = loud.Localize("You have successfully set default weapon!")
		case RSLT_BUYITM:
			desc = loud.Sprintf("You have bought %s from the shop", formatItem(screen.activeItem))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RSLT_BUYCHR:
			desc = loud.Sprintf("You have bought %s from the shop", formatCharacter(screen.activeCharacter))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RSLT_HUNT:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon"}
			desc = loud.Sprintf("You did hunt animals and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
			switch len(respOutput) {
			case 0:
				desc += "\nYour character is dead during hunt accidently"
			case 2:
				desc += "\nYou have lost your weapon accidently"
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
		case RSLT_GET_INITIAL_COIN:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			json.Unmarshal(screen.txResult, &respOutput)
			earnedAmount := int64(0)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			desc = loud.Sprintf("Got initial gold from pylons. Amount is %d.", earnedAmount)
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
			desc = loud.Sprintf("You sold %s for gold.", formatItem(screen.activeItem))
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
			desc += screen.buyItemDesc(screen.activeItem, screen.pylonEnterValue)
		case RSLT_BUYCHR_TRDREQ_CREATION:
			desc = loud.Localize("character buy request was successfully created")
			desc += screen.buyCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
		case RSLT_FULFILL_BUY_LOUD_TRDREQ:
			request := screen.activeTradeRequest
			desc = loud.Localize("you have sold loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.sellLoudDesc(request.Amount, request.Total)
		case RSLT_FULFILL_SELL_LOUD_TRDREQ:
			request := screen.activeTradeRequest
			desc = loud.Localize("you have bought loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.buyLoudDesc(request.Amount, request.Total)
		case RSLT_FULFILL_SELLITM_TRDREQ:
			request := screen.activeItemTradeRequest.(loud.ItemSellTradeRequest)
			desc = loud.Localize("you have bought item successfully from item/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_SELLCHR_TRDREQ:
			request := screen.activeItemTradeRequest.(loud.CharacterSellTradeRequest)
			desc = loud.Localize("you have bought character successfully from character/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_BUYITM_TRDREQ:
			request := screen.activeItemTradeRequest.(loud.ItemBuyTradeRequest)
			desc = loud.Localize("you have sold item successfully from item/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RSLT_FULFILL_BUYCHR_TRDREQ:
			request := screen.activeItemTradeRequest.(loud.CharacterBuyTradeRequest)
			desc = loud.Localize("you have sold character successfully from character/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		}
	}
	return desc
}

func (screen *GameScreen) TxWaitSituationDesc() string {
	desc := ""
	W8_PROC_TO_END := "\n" + loud.Localize("Please wait for a moment to finish the process")
	switch screen.scrStatus {
	case W8_BUY_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for loud buy request creation")
		desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_SELL_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for loud sell request creation")
		desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_BUYITM_PROC:
		desc = loud.Sprintf("You are now buying %s at the shop", formatItem(screen.activeItem))
		desc += W8_PROC_TO_END
	case W8_BUYCHR_PROC:
		desc = loud.Sprintf("You are now buying %s at the shop", formatCharacter(screen.activeCharacter))
		desc += W8_PROC_TO_END
	case W8_HUNT_PROC:
		if len(screen.activeItem.Name) > 0 {
			desc = loud.Sprintf("You are now hunting with %s", formatItem(screen.activeItem))
		} else {
			desc = loud.Localize("You are now hunting without weapon")
		}
		desc += W8_PROC_TO_END
	case W8_FIGHT_GIANT_PROC:
		desc = loud.Sprintf("You are now fighting with giant with %s", formatItem(screen.activeItem))
	case W8_FIGHT_GOBLIN_PROC:
		desc = loud.Sprintf("You are now fighting with goblin with %s", formatItem(screen.activeItem))
	case W8_FIGHT_TROLL_PROC:
		desc = loud.Sprintf("You are now fighting with troll with %s", formatItem(screen.activeItem))
	case W8_FIGHT_WOLF_PROC:
		desc = loud.Sprintf("You are now fighting with wolf with %s", formatItem(screen.activeItem))
	case W8_GET_INITIAL_COIN:
		desc = loud.Localize("Getting initial gold from pylon")
		desc += W8_PROC_TO_END
	case W8_DEV_GET_TEST_ITEMS:
		desc = loud.Localize("Getting dev test items from pylon")
		desc += W8_PROC_TO_END
	case W8_HEALTH_RESTORE_CHAR:
		desc = loud.Localize("Waiting for Health restoring")
		desc += W8_PROC_TO_END
	case W8_GET_PYLONS:
		desc = loud.Localize("You are waiting for getting pylons process")
	case W8_SWITCH_USER:
		desc = loud.Localize("You are waiting for switching to new user")
	case W8_CREATE_COOKBOOK:
		desc = loud.Localize("You are waiting for creating cookbook")
	case W8_SELLITM_PROC:
		desc = loud.Sprintf("You are now selling %s for gold", formatItem(screen.activeItem))
		desc += W8_PROC_TO_END
	case W8_UPGITM_PROC:
		desc = loud.Sprintf("You are now upgrading %s", loud.Localize(screen.activeItem.Name))
		desc += W8_PROC_TO_END
	case W8_SELLITM_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for item sell request creation")
		desc += screen.sellItemDesc(screen.activeItem, screen.pylonEnterValue)
	case W8_SELLCHR_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for character sell request creation")
		desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
	case W8_BUYITM_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for item buy request creation")
		desc += screen.buyItemDesc(screen.activeItem, screen.pylonEnterValue)
	case W8_BUYCHR_TRDREQ_CREATION:
		desc = loud.Localize("You are now waiting for character buy request creation")
		desc += screen.buyCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
	// For FULFILL trades, msg should be reversed, since user is opposite
	case W8_FULFILL_SELLITM_TRDREQ:
		request := screen.activeItemTradeRequest.(loud.ItemSellTradeRequest)
		desc = loud.Sprintf("You are now buying item at %d", request.Price)
		desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_SELLCHR_TRDREQ:
		request := screen.activeItemTradeRequest.(loud.CharacterSellTradeRequest)
		desc = loud.Localize("you are now buying character ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUYITM_TRDREQ:
		request := screen.activeItemTradeRequest.(loud.ItemBuyTradeRequest)
		desc = loud.Localize("you are now selling item ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUYCHR_TRDREQ:
		request := screen.activeItemTradeRequest.(loud.CharacterBuyTradeRequest)
		desc = loud.Localize("you are now selling character ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUY_LOUD_TRDREQ:
		request := screen.activeTradeRequest
		desc = loud.Localize("you are now selling loud for pylon") + fmt.Sprintf(" at %.4f.\n", request.Price)
		desc += screen.sellLoudDesc(request.Amount, request.Total)
	case W8_FULFILL_SELL_LOUD_TRDREQ:
		request := screen.activeTradeRequest
		desc = loud.Localize("you are now buying loud from pylon") + fmt.Sprintf(" at %.4f.\n", request.Price)
		desc += screen.buyLoudDesc(request.Amount, request.Total)
	}
	return desc
}
