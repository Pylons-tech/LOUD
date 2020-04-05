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
	case SHOW_LOCATION:
		locationDescMap := map[loud.UserLocation]string{
			loud.HOME:     loud.Localize("home desc"),
			loud.FOREST:   loud.Localize("forest desc"),
			loud.SHOP:     loud.Localize("shop desc"),
			loud.MARKET:   loud.Localize("market desc"),
			loud.SETTINGS: loud.Localize("settings desc"),
			loud.DEVELOP:  loud.Localize("develop desc"),
		}
		desc = locationDescMap[screen.user.GetLocation()]
	case SHOW_LOUD_BUY_REQUESTS:
		infoLines = screen.renderTRTable(loud.BuyTradeRequests)
	case SHOW_LOUD_SELL_REQUESTS:
		infoLines = screen.renderTRTable(loud.SellTradeRequests)
	case SHOW_BUY_SWORD_REQUESTS:
		infoLines = screen.renderITRTable(
			"Buy sword requests",
			[2]string{"Item", "Price (pylon)"},
			loud.SwordBuyTradeRequests)
	case SHOW_SELL_SWORD_REQUESTS:
		infoLines = screen.renderITRTable(
			"Sell sword requests",
			[2]string{"Item", "Price (pylon)"},
			loud.SwordSellTradeRequests)
	case SHOW_SELL_CHARACTER_REQUESTS:
		infoLines = screen.renderITRTable(
			"Sell character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterSellTradeRequests)
	case SHOW_BUY_CHARACTER_REQUESTS:
		infoLines = screen.renderITRTable(
			"Buy character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterBuyTradeRequests)
	case CREATE_BUY_LOUD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CREATE_SELL_LOUD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to get (should be integer value)" // TODO should add Localize
	case CREATE_BUY_LOUD_REQUEST_ENTER_LOUD_VALUE:
		desc = "Please enter loud amount to buy (should be integer value)" // TODO should add Localize
	case CREATE_SELL_LOUD_REQUEST_ENTER_LOUD_VALUE:
		desc = "Please enter loud amount to sell (should be integer value)" // TODO should add Localize
	case CREATE_SELL_SWORD_REQUEST_SELECT_SWORD:
		infoLines = screen.renderITTable("Select sword to sell", "Item", screen.user.InventoryItems())
	case CREATE_SELL_CHARACTER_REQUEST_SELECT_CHARACTER:
		infoLines = screen.renderITTable("Select character to sell", "Character", screen.user.InventoryCharacters())
	case CREATE_SELL_SWORD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CREATE_SELL_CHARACTER_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CREATE_BUY_SWORD_REQUEST_SELECT_SWORD:
		infoLines = screen.renderITTable("Select sword to buy", "Item", loud.WorldItemSpecs)
	case CREATE_BUY_CHARACTER_REQUEST_SELECT_CHARACTER:
		infoLines = screen.renderITTable("Select character specs to get", "Character", loud.WorldCharacterSpecs)
	case CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CREATE_BUY_CHARACTER_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case SELECT_DEFAULT_CHAR:
		infoLines = screen.renderITTable("Please select default character", "Character", screen.user.InventoryCharacters())
	case SELECT_HEALTH_RESTORE_CHAR:
		infoLines = screen.renderITTable("Please select character to restore health", "Character", screen.user.InventoryCharacters())
	case SELECT_DEFAULT_WEAPON:
		infoLines = screen.renderITTable("Please select default weapon", "Item", screen.user.InventorySwords())
	case SELECT_BUY_ITEM:
		infoLines = screen.renderITTable("select buy item desc", "Item", loud.ShopItems)
	case SELECT_BUY_CHARACTER:
		infoLines = screen.renderITTable("select buy character desc", "Character", loud.ShopCharacters)
	case SELECT_SELL_ITEM:
		infoLines = screen.renderITTable("select sell item desc", "Item", screen.user.InventoryItems())
	case SELECT_HUNT_ITEM:
		infoLines = screen.renderITTable("select hunt item desc", "Item", screen.user.InventorySwords())
	case SELECT_FIGHT_GOBLIN_ITEM:
		infoLines = screen.renderITTable("select fight goblin item desc", "Item", screen.user.InventorySwords())
	case SELECT_FIGHT_WOLF_ITEM:
		infoLines = screen.renderITTable("select fight wolf item desc", "Item", screen.user.InventorySwords())
	case SELECT_FIGHT_TROLL_ITEM:
		infoLines = screen.renderITTable("select fight troll item desc", "Item", screen.user.InventorySwords())
	case SELECT_FIGHT_GIANT_ITEM:
		infoLines = screen.renderITTable("select fight giant item desc", "Item", screen.user.InventoryIronSwords())
	case SELECT_UPGRADE_ITEM:
		infoLines = screen.renderITTable("select upgrade item desc", "Item", screen.user.InventoryUpgradableItems())
	}

	if strings.HasPrefix(string(screen.scrStatus), "RESULT_") {
		desc = screen.TxResultSituationDesc()
	}

	if strings.HasPrefix(string(screen.scrStatus), "WAIT_") {
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
		RESULT_BUY_LOUD_REQUEST_CREATION:       "loud buy request creation",
		RESULT_SELL_LOUD_REQUEST_CREATION:      "loud sell request creation",
		RESULT_SELECT_DEF_CHAR:                 "selecting default character",
		RESULT_HEALTH_RESTORE_CHAR:             "selecting character to restore health",
		RESULT_SELECT_DEF_WEAPON:               "selecting default weapon",
		RESULT_BUY_ITEM_FINISH:                 "buy item",
		RESULT_BUY_CHARACTER_FINISH:            "buy character",
		RESULT_HUNT_FINISH:                     "hunt",
		RESULT_GET_INITIAL_COIN:                "get initial coin",
		RESULT_GET_PYLONS:                      "get pylon",
		RESULT_SWITCH_USER:                     "switch user",
		RESULT_CREATE_COOKBOOK:                 "create cookbook",
		RESULT_SELL_FINISH:                     "sell item",
		RESULT_UPGRADE_FINISH:                  "upgrade item",
		RESULT_SELL_SWORD_REQUEST_CREATION:     "sell sword request creation",
		RESULT_BUY_SWORD_REQUEST_CREATION:      "buy sword request creation",
		RESULT_SELL_CHARACTER_REQUEST_CREATION: "sell character request creation",
		RESULT_BUY_CHARACTER_REQUEST_CREATION:  "buy character request creation",
		RESULT_FULFILL_BUY_LOUD_REQUEST:        "sell loud", // for fullfill direction is reversed
		RESULT_FULFILL_SELL_LOUD_REQUEST:       "buy loud",
		RESULT_FULFILL_SELL_SWORD_REQUEST:      "buy sword",
		RESULT_FULFILL_SELL_CHARACTER_REQUEST:  "buy character",
		RESULT_FULFILL_BUY_SWORD_REQUEST:       "sell sword",
		RESULT_FULFILL_BUY_CHARACTER_REQUEST:   "sell character",
	}
	if screen.txFailReason != "" {
		desc = resDescMap[screen.scrStatus] + ": " + loud.Localize(screen.txFailReason)
	} else {
		switch screen.scrStatus {
		case RESULT_BUY_LOUD_REQUEST_CREATION:
			desc = loud.Localize("loud buy request was successfully created")
			desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		case RESULT_SELL_LOUD_REQUEST_CREATION:
			desc = loud.Localize("loud sell request was successfully created")
			desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		case RESULT_SELECT_DEF_CHAR:
			desc = loud.Localize("You have successfully set default character!")
		case RESULT_HEALTH_RESTORE_CHAR:
			desc = loud.Localize("You have successfully restored character's health!")
		case RESULT_SELECT_DEF_WEAPON:
			desc = loud.Localize("You have successfully set default weapon!")
		case RESULT_BUY_ITEM_FINISH:
			desc = loud.Sprintf("You have bought %s from the shop", formatItem(screen.activeItem))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RESULT_BUY_CHARACTER_FINISH:
			desc = loud.Sprintf("You have bought %s from the shop", formatCharacter(screen.activeCharacter))
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RESULT_HUNT_FINISH:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon"}
			desc = loud.Sprintf("You did hunt animals and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
		case RESULT_FIGHT_GOBLIN_FINISH:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon", loud.GOBLIN_EAR}
			desc = loud.Sprintf("You did fight with goblin and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
		case RESULT_FIGHT_TROLL_FINISH:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon", loud.TROLL_TOES}
			desc = loud.Sprintf("You did fight with troll and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
		case RESULT_FIGHT_WOLF_FINISH:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon", loud.WOLF_TAIL}
			desc = loud.Sprintf("You did fight with wolf and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
		case RESULT_FIGHT_GIANT_FINISH:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			earnedAmount := int64(0)
			json.Unmarshal(screen.txResult, &respOutput)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			resultTexts := []string{"gold", "character", "weapon"}
			desc = loud.Sprintf("You did fight with giant and earned %d. Detailed result: %+v", earnedAmount, resultTexts[:len(respOutput)])
		case RESULT_GET_INITIAL_COIN:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			json.Unmarshal(screen.txResult, &respOutput)
			earnedAmount := int64(0)
			if len(respOutput) > 0 {
				earnedAmount = respOutput[0].Amount
			}
			desc = loud.Sprintf("Got initial gold from pylons. Amount is %d.", earnedAmount)
		case RESULT_DEV_GET_TEST_ITEMS:
			respOutput := []handlers.ExecuteRecipeSerialize{}
			json.Unmarshal(screen.txResult, &respOutput)
			resultTexts := []string{loud.WOLF_TAIL, loud.TROLL_TOES, loud.GOBLIN_EAR}
			desc = loud.Sprintf("Finished getting developer test items. Detailed result %+v", resultTexts[:len(respOutput)])
		case RESULT_GET_PYLONS:
			desc = loud.Localize("You got extra pylons for loud game")
		case RESULT_SWITCH_USER:
			desc = loud.Sprintf("You switched user to %s", screen.user.GetUserName())
		case RESULT_CREATE_COOKBOOK:
			desc = loud.Localize("You created a new cookbook for a new game build")
		case RESULT_SELL_FINISH:
			desc = loud.Sprintf("You sold %s for gold.", formatItem(screen.activeItem))
		case RESULT_UPGRADE_FINISH:
			desc = loud.Sprintf("You have upgraded %s to get better hunt result", screen.activeItem.Name)
		case RESULT_SELL_SWORD_REQUEST_CREATION:
			desc = loud.Localize("sword sell request was successfully created")
			desc += screen.sellSwordDesc(screen.activeItem, screen.pylonEnterValue)
		case RESULT_SELL_CHARACTER_REQUEST_CREATION:
			desc = loud.Localize("character sell request was successfully created")
			desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
		case RESULT_BUY_SWORD_REQUEST_CREATION:
			desc = loud.Localize("sword buy request was successfully created")
			desc += screen.buySwordDesc(screen.activeItem, screen.pylonEnterValue)
		case RESULT_BUY_CHARACTER_REQUEST_CREATION:
			desc = loud.Localize("character buy request was successfully created")
			desc += screen.buyCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
		case RESULT_FULFILL_BUY_LOUD_REQUEST:
			request := screen.activeTradeRequest
			desc = loud.Localize("you have sold loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.sellLoudDesc(request.Amount, request.Total)
		case RESULT_FULFILL_SELL_LOUD_REQUEST:
			request := screen.activeTradeRequest
			desc = loud.Localize("you have bought loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.buyLoudDesc(request.Amount, request.Total)
		case RESULT_FULFILL_SELL_SWORD_REQUEST:
			request := screen.activeItemTradeRequest.(loud.ItemSellTradeRequest)
			desc = loud.Localize("you have bought sword successfully from sword/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.buySwordDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RESULT_FULFILL_SELL_CHARACTER_REQUEST:
			request := screen.activeItemTradeRequest.(loud.CharacterSellTradeRequest)
			desc = loud.Localize("you have bought character successfully from character/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		case RESULT_FULFILL_BUY_SWORD_REQUEST:
			request := screen.activeItemTradeRequest.(loud.ItemBuyTradeRequest)
			desc = loud.Localize("you have sold sword successfully from sword/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.sellSwordSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RESULT_FULFILL_BUY_CHARACTER_REQUEST:
			request := screen.activeItemTradeRequest.(loud.CharacterBuyTradeRequest)
			desc = loud.Localize("you have sold character successfully from character/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		}
	}
	return desc
}

func (screen *GameScreen) TxWaitSituationDesc() string {
	desc := ""
	WAIT_PROCESS_TO_END := "\n" + loud.Localize("Please wait for a moment to finish the process")
	switch screen.scrStatus {
	case WAIT_BUY_LOUD_REQUEST_CREATION:
		desc = loud.Localize("You are now waiting for loud buy request creation")
		desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case WAIT_SELL_LOUD_REQUEST_CREATION:
		desc = loud.Localize("You are now waiting for loud sell request creation")
		desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case WAIT_BUY_ITEM_PROCESS:
		desc = loud.Sprintf("You are now buying %s at the shop", formatItem(screen.activeItem))
		desc += WAIT_PROCESS_TO_END
	case WAIT_BUY_CHARACTER_PROCESS:
		desc = loud.Sprintf("You are now buying %s at the shop", formatCharacter(screen.activeCharacter))
		desc += WAIT_PROCESS_TO_END
	case WAIT_HUNT_PROCESS:
		if len(screen.activeItem.Name) > 0 {
			desc = loud.Sprintf("You are now hunting with %s", formatItem(screen.activeItem))
		} else {
			desc = loud.Localize("You are now hunting without weapon")
		}
		desc += WAIT_PROCESS_TO_END
	case WAIT_FIGHT_GIANT_PROCESS:
		desc = loud.Sprintf("You are now fighting with giant with %s", formatItem(screen.activeItem))
	case WAIT_FIGHT_GOBLIN_PROCESS:
		desc = loud.Sprintf("You are now fighting with goblin with %s", formatItem(screen.activeItem))
	case WAIT_FIGHT_TROLL_PROCESS:
		desc = loud.Sprintf("You are now fighting with troll with %s", formatItem(screen.activeItem))
	case WAIT_FIGHT_WOLF_PROCESS:
		desc = loud.Sprintf("You are now fighting with wolf with %s", formatItem(screen.activeItem))
	case WAIT_GET_INITIAL_COIN:
		desc = loud.Localize("Getting initial gold from pylon")
		desc += WAIT_PROCESS_TO_END
	case WAIT_DEV_GET_TEST_ITEMS:
		desc = loud.Localize("Getting dev test items from pylon")
		desc += WAIT_PROCESS_TO_END
	case WAIT_HEALTH_RESTORE_CHAR:
		desc = loud.Localize("Waiting for Health restoring")
		desc += WAIT_PROCESS_TO_END
	case WAIT_GET_PYLONS:
		desc = loud.Localize("You are waiting for getting pylons process")
	case WAIT_SWITCH_USER:
		desc = loud.Localize("You are waiting for switching to new user")
	case WAIT_CREATE_COOKBOOK:
		desc = loud.Localize("You are waiting for creating cookbook")
	case WAIT_SELL_PROCESS:
		desc = loud.Sprintf("You are now selling %s for gold", formatItem(screen.activeItem))
		desc += WAIT_PROCESS_TO_END
	case WAIT_UPGRADE_PROCESS:
		desc = loud.Sprintf("You are now upgrading %s", loud.Localize(screen.activeItem.Name))
		desc += WAIT_PROCESS_TO_END
	case WAIT_SELL_SWORD_REQUEST_CREATION:
		desc = loud.Localize("You are now waiting for sword sell request creation")
		desc += screen.sellSwordDesc(screen.activeItem, screen.pylonEnterValue)
	case WAIT_SELL_CHARACTER_REQUEST_CREATION:
		desc = loud.Localize("You are now waiting for character sell request creation")
		desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
	case WAIT_BUY_SWORD_REQUEST_CREATION:
		desc = loud.Localize("You are now waiting for sword buy request creation")
		desc += screen.buySwordDesc(screen.activeItem, screen.pylonEnterValue)
	case WAIT_BUY_CHARACTER_REQUEST_CREATION:
		desc = loud.Localize("You are now waiting for character buy request creation")
		desc += screen.buyCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
	// For FULFILL trades, msg should be reversed, since user is opposite
	case WAIT_FULFILL_SELL_SWORD_REQUEST:
		request := screen.activeItemTradeRequest.(loud.ItemSellTradeRequest)
		desc = loud.Sprintf("You are now buying sword at %d", request.Price)
		desc += screen.buySwordDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case WAIT_FULFILL_SELL_CHARACTER_REQUEST:
		request := screen.activeItemTradeRequest.(loud.CharacterSellTradeRequest)
		desc = loud.Localize("you are now buying character ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case WAIT_FULFILL_BUY_SWORD_REQUEST:
		request := screen.activeItemTradeRequest.(loud.ItemBuyTradeRequest)
		desc = loud.Localize("you are now selling sword ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.sellSwordSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case WAIT_FULFILL_BUY_CHARACTER_REQUEST:
		request := screen.activeItemTradeRequest.(loud.CharacterBuyTradeRequest)
		desc = loud.Localize("you are now selling character ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case WAIT_FULFILL_BUY_LOUD_REQUEST:
		request := screen.activeTradeRequest
		desc = loud.Localize("you are now selling loud for pylon") + fmt.Sprintf(" at %.4f.\n", request.Price)
		desc += screen.sellLoudDesc(request.Amount, request.Total)
	case WAIT_FULFILL_SELL_LOUD_REQUEST:
		request := screen.activeTradeRequest
		desc = loud.Localize("you are now buying loud from pylon") + fmt.Sprintf(" at %.4f.\n", request.Price)
		desc += screen.buyLoudDesc(request.Amount, request.Total)
	}
	return desc
}
