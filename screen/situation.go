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
	waitProcessEnd := loud.Localize("wait process to end")
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
	case CREATE_BUY_LOUD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case CREATE_SELL_LOUD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to get (should be integer value)" // TODO should add Localize
	case CREATE_BUY_LOUD_REQUEST_ENTER_LOUD_VALUE:
		desc = "Please enter loud amount to buy (should be integer value)" // TODO should add Localize
	case CREATE_SELL_LOUD_REQUEST_ENTER_LOUD_VALUE:
		desc = "Please enter loud amount to sell (should be integer value)" // TODO should add Localize

	case SHOW_LOUD_BUY_REQUESTS:
		infoLines = screen.renderTradeRequestTable(loud.BuyTradeRequests)
	case SHOW_LOUD_SELL_REQUESTS:
		infoLines = screen.renderTradeRequestTable(loud.SellTradeRequests)
	case SELECT_BUY_ITEM:
		desc = loud.Localize("select buy item desc")
	case SELECT_SELL_ITEM:
		desc = loud.Localize("select sell item desc")
	case SELECT_HUNT_ITEM:
		desc = loud.Localize("select hunt item desc")
	case SELECT_UPGRADE_ITEM:
		desc = loud.Localize("select upgrade item desc")
	case WAIT_BUY_LOUD_REQUEST_CREATION:
		desc = loud.Localize("you are now waiting for loud buy request creation")
		desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case WAIT_SELL_LOUD_REQUEST_CREATION:
		desc = loud.Localize("you are now waiting for loud sell request creation")
		desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case WAIT_BUY_PROCESS:
		desc = fmt.Sprintf("%s %s Lv%d.\n%s", loud.Localize("wait buy process desc"), loud.Localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
	case WAIT_HUNT_PROCESS:
		if len(screen.activeItem.Name) > 0 {
			desc = fmt.Sprintf("%s %s Lv%d.\n%s", loud.Localize("wait hunt process desc"), loud.Localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
		} else {
			switch string(screen.lastInput.Ch) {
			case "I", "i":
				desc = fmt.Sprintf("%s\n%s", loud.Localize("Getting initial gold from pylon"), waitProcessEnd)
			default:
				desc = fmt.Sprintf("%s\n%s", loud.Localize("hunting without weapon"), waitProcessEnd)
			}
		}
	case WAIT_GET_PYLONS:
		desc = loud.Localize("You are waiting for getting pylons process")
	case WAIT_SWITCH_USER:
		desc = loud.Localize("You are waiting for switching to new user")
	case WAIT_CREATE_COOKBOOK:
		desc = loud.Localize("You are waiting for creating cookbook")
	case WAIT_SELL_PROCESS:
		desc = fmt.Sprintf("%s %s Lv%d.\n%s", loud.Localize("wait sell process desc"), loud.Localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
	case WAIT_UPGRADE_PROCESS:
		desc = fmt.Sprintf("%s %s.\n%s", loud.Localize("wait upgrade process desc"), loud.Localize(screen.activeItem.Name), waitProcessEnd)
	case RESULT_BUY_LOUD_REQUEST_CREATION:
		if screen.txFailReason != "" {
			desc = loud.Localize("loud buy request creation fail reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = loud.Localize("loud buy request was successfully created")
			desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		}
	case RESULT_SELL_LOUD_REQUEST_CREATION:
		if screen.txFailReason != "" {
			desc = loud.Localize("loud sell request creation fail reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = loud.Localize("loud sell request was successfully created")
			desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
		}
	case RESULT_BUY_FINISH:
		if screen.txFailReason != "" {
			desc = loud.Localize("buy failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("%s %s Lv%d.\n%s", loud.Localize("result buy finish desc"), loud.Localize(screen.activeItem.Name), screen.activeItem.Level, loud.Localize("use for hunting"))
		}
	case RESULT_HUNT_FINISH:
		if screen.txFailReason != "" {
			desc = loud.Localize("hunt failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			respOutput := handlers.ExecuteRecipeSerialize{}
			json.Unmarshal(screen.txResult, &respOutput)
			switch string(screen.lastInput.Ch) {
			case "I", "i":
				desc = fmt.Sprintf("%s %d.", loud.Localize("Got initial gold from pylons. Amount is"), respOutput.Amount)
			default:
				desc = fmt.Sprintf("%s %d.", loud.Localize("result hunt finish desc"), respOutput.Amount)
			}
		}
	case RESULT_GET_PYLONS:
		if screen.txFailReason != "" {
			desc = loud.Localize("get pylon failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("You got extra pylons for loud game")
		}
	case RESULT_SWITCH_USER:
		if screen.txFailReason != "" {
			desc = loud.Localize("switch user fail reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("You switched user to %s", screen.user.GetUserName())
		}
	case RESULT_CREATE_COOKBOOK:
		if screen.txFailReason != "" {
			desc = loud.Localize("create cookbook failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("You created a new cookbook for a new game build")
		}
	case RESULT_SELL_FINISH:
		if screen.txFailReason != "" {
			desc = loud.Localize("sell failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("%s %s Lv%d.", loud.Localize("result sell finish desc"), loud.Localize(screen.activeItem.Name), screen.activeItem.Level)
		}
	case RESULT_UPGRADE_FINISH:
		if screen.txFailReason != "" {
			desc = loud.Localize("upgrade failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = fmt.Sprintf("%s: %s.", loud.Localize("result upgrade finish desc"), loud.Localize(screen.activeItem.Name))
		}
	case SHOW_BUY_SWORD_REQUESTS:
		infoLines = screen.renderItemTradeRequestTable(loud.SwordBuyTradeRequests)
	case CREATE_SELL_SWORD_REQUEST_SELECT_SWORD:
		infoLines = screen.renderItemTable(screen.user.InventoryItems())
	case CREATE_SELL_SWORD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case WAIT_SELL_SWORD_REQUEST_CREATION:
		desc = loud.Localize("you are now waiting for sword sell request creation")
		desc += screen.sellSwordDesc(screen.activeItem, screen.pylonEnterValue)
	case RESULT_SELL_SWORD_REQUEST_CREATION:
		if screen.txFailReason != "" {
			desc = loud.Localize("sword sell request creation fail reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = loud.Localize("sword sell request was successfully created")
			desc += screen.sellSwordDesc(screen.activeItem, screen.pylonEnterValue)
		}
	case SHOW_SELL_SWORD_REQUESTS:
		infoLines = screen.renderItemTradeRequestTable(loud.SwordSellTradeRequests)
	case CREATE_BUY_SWORD_REQUEST_SELECT_SWORD:
		infoLines = screen.renderItemTable(loud.WorldItems)
	case CREATE_BUY_SWORD_REQUEST_ENTER_PYLON_VALUE:
		desc = "Please enter pylon amount to use (should be integer value)" // TODO should add Localize
	case WAIT_BUY_SWORD_REQUEST_CREATION:
		desc = loud.Localize("you are now waiting for sword buy request creation")
		desc += screen.buySwordDesc(screen.activeItem, screen.pylonEnterValue)
	case RESULT_BUY_SWORD_REQUEST_CREATION:
		if screen.txFailReason != "" {
			desc = loud.Localize("sword buy request creation fail reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			desc = loud.Localize("sword buy request was successfully created")
			desc += screen.buySwordDesc(screen.activeItem, screen.pylonEnterValue)
		}
	// For FULFILL trades, msg should be reversed, since user is opposite
	case WAIT_FULFILL_BUY_LOUD_REQUEST:
		request := screen.activeTradeRequest
		desc = loud.Localize("you are now selling loud for pylon") + fmt.Sprintf(" at %.4f.\n", request.Price)
		desc += screen.sellLoudDesc(request.Amount, request.Total)
	case WAIT_FULFILL_SELL_LOUD_REQUEST:
		request := screen.activeTradeRequest
		desc = loud.Localize("you are now buying loud from pylon") + fmt.Sprintf(" at %.4f.\n", request.Price)
		desc += screen.buyLoudDesc(request.Amount, request.Total)
	case RESULT_FULFILL_BUY_LOUD_REQUEST:
		if screen.txFailReason != "" {
			desc = loud.Localize("sell loud failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			request := screen.activeTradeRequest
			desc = loud.Localize("you have sold loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.sellLoudDesc(request.Amount, request.Total)
		}
	case RESULT_FULFILL_SELL_LOUD_REQUEST:
		if screen.txFailReason != "" {
			desc = loud.Localize("buy loud failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			request := screen.activeTradeRequest
			desc = loud.Localize("you have bought loud coin successfully from loud/pylon market") + fmt.Sprintf(" at %.4f.\n", request.Price)
			desc += screen.buyLoudDesc(request.Amount, request.Total)
		}
	case WAIT_FULFILL_SELL_SWORD_REQUEST:
		request := screen.activeItemTradeRequest
		desc = loud.Localize("you are now buying sword ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.buySwordDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case WAIT_FULFILL_BUY_SWORD_REQUEST:
		request := screen.activeItemTradeRequest
		desc = loud.Localize("you are now selling sword ") + fmt.Sprintf(" at %d.\n", request.Price)
		desc += screen.sellSwordDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case RESULT_FULFILL_SELL_SWORD_REQUEST:
		if screen.txFailReason != "" {
			desc = loud.Localize("buy sword failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			request := screen.activeItemTradeRequest
			desc = loud.Localize("you have bought sword successfully from sword/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.buySwordDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		}
	case RESULT_FULFILL_BUY_SWORD_REQUEST:
		if screen.txFailReason != "" {
			desc = loud.Localize("sell sword failed reason") + ": " + loud.Localize(screen.txFailReason)
		} else {
			request := screen.activeItemTradeRequest
			desc = loud.Localize("you have sold sword successfully from sword/pylon market") + fmt.Sprintf(" at %d.\n", request.Price)
			desc += screen.sellSwordDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		}
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
