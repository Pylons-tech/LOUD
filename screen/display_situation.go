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

func (screen *GameScreen) GetDisabledFontByActiveLine(idx int) FontType {
	if screen.activeLine == idx {
		return GREY_BOLD
	}
	return GREY
}

func (screen *GameScreen) renderUserSituation() {

	// situation box start point (x, y)
	scrBox := screen.GetSituationBox()
	x := scrBox.X
	y := scrBox.Y
	w := scrBox.W
	h := scrBox.H

	infoLines := []string{}
	tableLines := TextLines{}
	desc := ""
	descfont := REGULAR
	activeWeapon := screen.user.GetFightWeapon()
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
			loud.HELP:     loud.Localize("help desc"),
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
		tableLines = screen.renderTRTable(
			loud.BuyTrdReqs, w,
			func(idx int, request interface{}) FontType {
				tr := request.(loud.TrdReq)
				if screen.user.GetGold() < tr.Amount {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, tr.IsMyTrdReq)
			})
	case SHW_LOUD_SELL_TRDREQS:
		tableLines = screen.renderTRTable(
			loud.SellTrdReqs, w,
			func(idx int, request interface{}) FontType {
				tr := request.(loud.TrdReq)
				if screen.user.GetPylonAmount() < tr.Total {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, tr.IsMyTrdReq)
			})
	case SHW_BUYITM_TRDREQS:
		tableLines = screen.renderITRTable(
			"Buy item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemBuyTrdReqs,
			w,
			func(idx int, request interface{}) FontType {
				itr := request.(loud.ItemBuyTrdReq)
				if len(screen.user.GetMatchedItems(itr.TItem)) == 0 {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, itr.IsMyTrdReq)
			})
	case SHW_SELLITM_TRDREQS:
		tableLines = screen.renderITRTable(
			"Sell item requests",
			[2]string{"Item", "Price (pylon)"},
			loud.ItemSellTrdReqs,
			w,
			func(idx int, request interface{}) FontType {
				isMyTrdReq, _, requestPrice := RequestInfo(request)
				if screen.user.GetPylonAmount() < requestPrice {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, isMyTrdReq)
			})
	case SHW_SELLCHR_TRDREQS:
		tableLines = screen.renderITRTable(
			"Sell character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterSellTrdReqs,
			w,
			func(idx int, request interface{}) FontType {
				isMyTrdReq, _, requestPrice := RequestInfo(request)
				if screen.user.GetPylonAmount() < requestPrice {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, isMyTrdReq)
			})
	case SHW_BUYCHR_TRDREQS:
		tableLines = screen.renderITRTable(
			"Buy character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterBuyTrdReqs,
			w,
			func(idx int, request interface{}) FontType {
				itr := request.(loud.CharacterBuyTrdReq)
				if len(screen.user.GetMatchedCharacters(itr.TCharacter)) == 0 {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, itr.IsMyTrdReq)
			})
	case CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to get (should be integer value)")
	case CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL:
		desc = loud.Localize("Please enter gold amount to buy (should be integer value)")
	case RENAME_CHAR_ENT_NEWNAME:
		desc = loud.Localize("Please enter new character's name")
	case CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL:
		desc = loud.Localize("Please enter gold amount to sell (should be integer value)")

	case CR8_SELLITM_TRDREQ_SEL_ITEM:
		tableLines = screen.renderITTable(
			"Select item to sell",
			"Item",
			screen.user.InventoryItems(),
			w, nil)
	case CR8_SELLCHR_TRDREQ_SEL_CHR:
		tableLines = screen.renderITTable(
			"Select character to sell",
			"Character",
			screen.user.InventoryCharacters(),
			w, nil)
	case CR8_SELLITM_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_SELLCHR_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_BUYITM_TRDREQ_SEL_ITEM:
		tableLines = screen.renderITTable(
			"Select item to buy",
			"Item",
			loud.WorldItemSpecs,
			w, nil)
	case CR8_BUYCHR_TRDREQ_SEL_CHR:
		tableLines = screen.renderITTable(
			"Select character specs to get",
			"Character",
			loud.WorldCharacterSpecs,
			w, nil)
	case CR8_BUYITM_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CR8_BUYCHR_TRDREQ_ENT_PYLVAL:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case SEL_ACTIVE_CHAR:
		tableLines = screen.renderITTable(
			"Please select active character",
			"Character",
			screen.user.InventoryCharacters(),
			w, nil)
	case SEL_RENAME_CHAR:
		tableLines = screen.renderITTable(
			"Please select character to rename",
			"Character",
			screen.user.InventoryCharacters(),
			w, nil)
	case SEL_BUYITM:
		tableLines = screen.renderITTable(
			"select buy item desc",
			"Shop items",
			loud.ShopItems,
			w,
			func(idx int, item interface{}) FontType {
				return screen.getFontOfShopItem(idx, item.(loud.Item))
			})
	case SEL_SELLITM:
		tableLines = screen.renderITTable(
			"select sell item desc",
			"Item",
			screen.user.InventorySellableItems(),
			w, nil)
	case SEL_UPGITM:
		tableLines = screen.renderITTable(
			"select upgrade item desc",
			"Item",
			screen.user.InventoryUpgradableItems(),
			w, nil)
	case SEL_BUYCHR:
		tableLines = screen.renderITTable(
			"select buy character desc",
			"Character",
			loud.ShopCharacters,
			w, nil)
	case CONFIRM_HUNT_RABBITS:
		desc = loud.Localize("rabbits outcome")
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
	case HELP_ABOUT:
		desc = loud.Localize("help about")
	case HELP_GAME_OBJECTIVE:
		desc = loud.Localize("help what you can do")
	case HELP_NAVIGATION:
		desc = loud.Localize("help navigation")
	case HELP_PAGE_LAYOUT:
		desc = loud.Localize("help page layout")
	case HELP_GAME_RULES:
		desc = loud.Localize("help game rules")
	case HELP_HOW_IT_WORKS:
		desc = loud.Localize("help how it works")
	case HELP_PYLONS_CENTRAL:
		tableLines = screen.tradeTableColorDesc(w)
	case HELP_UPCOMING_RELEASES:
		desc = loud.Localize("help upcoming releases")
	case HELP_SUPPORT:
		desc = loud.Localize("help Support")
	}

	if screen.IsResultScreen() {
		desc, descfont = screen.TxResultSituationDesc()
	}

	if screen.IsWaitScreen() {
		infoLines, tableLines = screen.TxWaitSituationDesc(w)
	}

	basicLines := loud.ChunkText(desc, w-2)

	colorfulLines := TextLines{}
	for _, chli := range basicLines {
		colorfulLines = colorfulLines.appendF(fillSpace(chli, w), descfont)
	}
	tableLines = append(colorfulLines, tableLines...)

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
			screen.getFont(line.font)(fillSpace(line.content, w))))
		if index+2 > int(screen.Height()) {
			break
		}
	}
	totalLen := infoLen + len(tableLines)

	screen.drawFill(x, y+totalLen, w, h-totalLen-1)
}

func monsterTextWithUnicode(monster string) string {
	unicodeMonsterTexts := map[string]string{
		loud.RABBIT:        "🐇 (rabbit)",
		loud.GOBLIN:        "👺 (goblin)",
		loud.WOLF:          "🐺 (wolf)",
		loud.TROLL:         "👻 (troll)",
		loud.GIANT:         "🗿 (giant)",
		loud.DRAGON_FIRE:   "🦐 (fire dragon)",
		loud.DRAGON_ICE:    "🦈 (ice dragon)",
		loud.DRAGON_ACID:   "🐊 (acid dragon)",
		loud.DRAGON_UNDEAD: "🐉 (undead dragon)",
	}
	if umt, ok := unicodeMonsterTexts[monster]; ok {
		return umt
	}
	return ""
}

func (screen *GameScreen) GetKilledByDesc() string {
	monsterDesc := monsterTextWithUnicode(screen.user.GetTargetMonster())
	return loud.Sprintf(
		"Your %s character was killed by %s accidently",
		formatCharacterP(screen.user.GetDeadCharacter()),
		monsterDesc,
	)
}

func (screen *GameScreen) GetLostWeaponDesc() string {
	monsterDesc := monsterTextWithUnicode(screen.user.GetTargetMonster())
	activeWeapon := screen.user.GetFightWeapon()
	activeWeaponName := ""
	activeWeaponLevel := 1

	if activeWeapon != nil {
		activeWeaponName = activeWeapon.Name
		activeWeaponLevel = activeWeapon.Level
	}
	return loud.Sprintf(
		"You have lost your %s %d weapon while fighting %s accidently",
		activeWeaponName, activeWeaponLevel,
		monsterDesc)
}

func (screen *GameScreen) TxResultSituationDesc() (string, FontType) {
	desc := ""
	font := REGULAR
	resDescMap := map[ScreenStatus]string{
		RSLT_BUY_LOUD_TRDREQ_CREATION:  "gold buy request creation",
		RSLT_SELL_LOUD_TRDREQ_CREATION: "gold sell request creation",
		RSLT_SEL_ACT_CHAR:              "selecting active character",
		RSLT_RENAME_CHAR:               "renaming character",
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
		case RSLT_RENAME_CHAR:
			desc = loud.Sprintf("You have successfully updated character's name to %s!", screen.inputText)
		case RSLT_BUYITM:
			desc = loud.Sprintf("You have bought a weapon from the shop")
			desc += "\n"
		case RSLT_BUYCHR:
			desc = loud.Sprintf("You have bought %s from Pylons Central", screen.activeCharacter.Name)
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RSLT_HUNT_RABBITS:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("Your %s character is dead while following rabbits accidently", formatCharacter(*screen.user.GetActiveCharacter()))
				font = RED
			} else {
				desc = loud.Sprintf("You did hunt rabbits and earned %d.", earnedAmount)
			}
		case RSLT_FIGHT_GOBLIN:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with goblin and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YELLOW
				case 4:
					bonusItem := screen.user.GetItemByID(respOutput[3].ItemID)
					if bonusItem == nil {
						font = RED_BOLD
						desc = loud.Localize("Something went wrong!\nReturned ItemID is not available on user's inventory.")
					} else {
						font = GREEN
						desc += loud.Sprintf("You got bonus item called %s", bonusItem.Name)
						desc += "\n"
						if bonusItem.Name == loud.GOBLIN_BOOTS { // GOBLIN_BOOTS
							desc += loud.Sprintf("You can sell goblin boots to earn gold or trade in pylons central!")
						} else { // GOBLIN_EAR
							desc += loud.Sprintf("You can make silver sword with Goblin ear at the shop!")
						}
					}
				}
			}
		case RSLT_FIGHT_TROLL:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with troll and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YELLOW
				case 4:
					bonusItem := screen.user.GetItemByID(respOutput[3].ItemID)
					if bonusItem == nil {
						font = RED_BOLD
						desc = loud.Localize("Something went wrong!\nReturned ItemID is not available on user's inventory.")
					} else {
						font = GREEN
						desc += loud.Sprintf("You got bonus item called %s", bonusItem.Name)
						desc += "\n"
						if bonusItem.Name == loud.TROLL_SMELLY_BONES { // TROLL_SMELLY_BONES
							desc += loud.Sprintf("You can sell troll's smelly boots at to earn gold or trade in pylons central!")
						} else { // TROLL_TOES
							desc += loud.Sprintf("You can make iron sword with Troll toes at the shop!")
						}
					}
				}
			}
		case RSLT_FIGHT_WOLF:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with wolf and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YELLOW
				case 4:
					bonusItem := screen.user.GetItemByID(respOutput[3].ItemID)
					if bonusItem == nil {
						font = RED_BOLD
						desc = loud.Localize("Something went wrong!\nReturned ItemID is not available on user's inventory.")
					} else {
						font = GREEN
						desc += loud.Sprintf("You got bonus item called %s", bonusItem.Name)
						desc += "\n"
						if bonusItem.Name == loud.WOLF_FUR { // WOLF_FUR
							desc += loud.Sprintf("You can sell wolf fur at to earn gold or trade in pylons central!")
						} else { // WOLF_TAIL
							desc += loud.Sprintf("You can make bronze sword with Wolf tail at the shop!")
						}
					}
				}
			}
		case RSLT_FIGHT_GIANT:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with giant and earned %d.", earnedAmount)
				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YELLOW
				case 3:
					activeCharacter := screen.user.GetActiveCharacter()
					if activeCharacter != nil && activeCharacter.Special != loud.NO_SPECIAL { // Got special from this fight
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
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with fire dragon and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
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
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with ice dragon and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
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
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with acid dragon and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
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
				desc = screen.GetKilledByDesc()
				font = RED
			} else {
				desc = loud.Sprintf("You did fight with undead dragon and earned %d.", earnedAmount)
				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YELLOW
				}
			}
		case RSLT_BUY_GOLD_WITH_PYLONS:
			earnedAmount, _ := screen.GetTxResponseOutput()
			desc = loud.Sprintf("Bought %d gold with %d pylons.", earnedAmount, 100)
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
			desc = loud.Sprintf("You sold %s for %d gold.", screen.activeItem.Name, earnedAmount)
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

func (screen *GameScreen) TxWaitSituationDesc(width int) ([]string, TextLines) {
	desc := ""
	monsterName := monsterTextWithUnicode(screen.user.GetTargetMonster())
	activeWeapon := screen.user.GetFightWeapon()
	activeWeaponName := ""
	activeWeaponLevel := 1

	if activeWeapon != nil {
		activeWeaponLevel = activeWeapon.Level
		activeWeaponName = activeWeapon.Name
	}
	W8_TO_END := "\n" + loud.Localize("Please wait for a moment to finish the process")
	switch screen.scrStatus {
	case W8_RENAME_CHAR:
		desc = loud.Sprintf("You are waiting to rename character from %s to %s.", screen.activeCharacter.Name, screen.inputText)
	case W8_BUY_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are waiting for gold buy request creation")
		desc += screen.buyLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_SELL_LOUD_TRDREQ_CREATION:
		desc = loud.Localize("You are waiting for gold sell request creation")
		desc += screen.sellLoudDesc(screen.loudEnterValue, screen.pylonEnterValue)
	case W8_BUYITM:
		desc = loud.Sprintf("You are buying %s at the shop", screen.activeItem.Name)
	case W8_BUYCHR:
		desc = loud.Sprintf("You are buying %s at the pylons central", formatCharacter(screen.activeCharacter))
		desc += W8_TO_END
	case W8_HUNT_RABBITS:
		desc = loud.Sprintf("You are hunting rabbits")
		desc += W8_TO_END
	case W8_FIGHT_TROLL,
		W8_FIGHT_WOLF,
		W8_FIGHT_GOBLIN,
		W8_FIGHT_GIANT,
		W8_FIGHT_DRAGONFIRE,
		W8_FIGHT_DRAGONICE,
		W8_FIGHT_DRAGONACID,
		W8_FIGHT_DRAGONUNDEAD:
		desc = loud.Sprintf("You are fighting %s monster with weapon %s level %d", monsterName, activeWeaponName, activeWeaponLevel)
	case W8_BUY_GOLD_WITH_PYLONS:
		desc = loud.Sprintf("spending %d pylon for %d gold", 100, 5000)
	case W8_DEV_GET_TEST_ITEMS:
		desc = loud.Sprintf("Getting dev test items from pylon")
		desc += W8_TO_END
	case W8_GET_PYLONS:
		desc = loud.Sprintf("You are waiting for getting pylons process")
	case W8_SWITCH_USER:
		desc = loud.Sprintf("You are waiting for switching to new user")
	case W8_CREATE_COOKBOOK:
		desc = loud.Sprintf("You are waiting for creating cookbook")
	case W8_SELLITM:
		item := screen.activeItem
		desc = loud.Sprintf("You are selling %s for %s gold", item.Name, item.GetSellPriceRange())
	case W8_UPGITM:
		desc = loud.Sprintf("You are upgrading %s", screen.activeItem.Name)
	case W8_SELLITM_TRDREQ_CREATION:
		desc = loud.Localize("You are waiting for item sell request creation")
		desc += screen.sellItemDesc(screen.activeItem, screen.pylonEnterValue)
	case W8_SELLCHR_TRDREQ_CREATION:
		desc = loud.Localize("You are waiting for character sell request creation")
		desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
	case W8_BUYITM_TRDREQ_CREATION:
		desc = loud.Localize("You are waiting for item buy request creation")
		desc += screen.buyItemSpecDesc(screen.activeItSpec, screen.pylonEnterValue)
	case W8_BUYCHR_TRDREQ_CREATION:
		desc = loud.Localize("You are waiting for character buy request creation")
		desc += screen.buyCharacterSpecDesc(screen.activeChSpec, screen.pylonEnterValue)
	case W8_CANCEL_TRDREQ:
		desc = loud.Localize("You are waiting for cancelling one of your trades")
	// For FULFILL trades, msg should be reversed, since user is opposite
	case W8_FULFILL_SELLITM_TRDREQ:
		request := screen.activeItemTrdReq.(loud.ItemSellTrdReq)
		desc = loud.Sprintf("You are buying item at %d", request.Price)
		desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_SELLCHR_TRDREQ:
		request := screen.activeItemTrdReq.(loud.CharacterSellTrdReq)
		desc = loud.Sprintf("You are buying character at %d.", request.Price)
		desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUYITM_TRDREQ:
		request := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
		desc = loud.Sprintf("You are selling item at %d.", request.Price)
		desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case W8_FULFILL_BUYCHR_TRDREQ:
		request := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
		desc = loud.Sprintf("You are selling character at %d.", request.Price)
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
	return loud.ChunkText(desc, width-2), TextLines{}.appendF(fillSpace("......", width), BLINK_BLUE_BOLD)
}
