package screen

import (
	"encoding/json"
	"fmt"
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

// GetTxResponseOutput returns parsed transaction output
func (screen *GameScreen) GetTxResponseOutput() (int64, []handlers.ExecuteRecipeSerialize) {
	respOutput := []handlers.ExecuteRecipeSerialize{}
	earnedAmount := int64(0)
	err := json.Unmarshal(screen.txResult, &respOutput)
	if err != nil {
		return 0, respOutput
	}
	if len(respOutput) > 0 {
		earnedAmount = respOutput[0].Amount
	}
	return earnedAmount, respOutput
}

// GetDisabledFontByActiveLine returns bold or not for grey line
func (screen *GameScreen) GetDisabledFontByActiveLine(idx int) FontType {
	if screen.activeLine == idx {
		return GreyBoldFont
	}
	return GreyFont
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
	descfont := RegularFont
	activeWeapon := screen.user.GetFightWeapon()
	switch screen.scrStatus {
	case ConfirmEndGame:
		desc = loud.Localize("Are you really gonna end game?")
	case ShowLocation:
		locationDescMap := map[loud.UserLocation]string{
			loud.Home:          loud.Localize("home desc"),
			loud.Forest:        loud.Localize("forest desc"),
			loud.Shop:          loud.Localize("shop desc"),
			loud.PylonsCentral: loud.Localize("pylons central desc"),
			loud.Settings:      loud.Localize("settings desc"),
			loud.Develop:       loud.Localize("develop desc"),
			loud.Help:          loud.Localize("help desc"),
		}
		desc = locationDescMap[screen.user.GetLocation()]
		if screen.user.GetLocation() == loud.Home {
			activeCharacter := screen.user.GetActiveCharacter()
			if activeCharacter == nil {
				desc = loud.Localize("home desc without character")
			} else if screen.user.GetPylonAmount() == 0 {
				desc = loud.Localize("home desc without pylon")
			}
		}
	case ShowGoldBuyTrdReqs:
		tableLines = screen.renderTRTable(
			loud.BuyTrdReqs, w,
			func(idx int, request interface{}) FontType {
				tr := request.(loud.TrdReq)
				if screen.user.GetGold() < tr.Amount {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, tr.IsMyTrdReq)
			})
	case ShowGoldSellTrdReqs:
		tableLines = screen.renderTRTable(
			loud.SellTrdReqs, w,
			func(idx int, request interface{}) FontType {
				tr := request.(loud.TrdReq)
				if screen.user.GetPylonAmount() < tr.Total {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, tr.IsMyTrdReq)
			})
	case ShowBuyItemTrdReqs:
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
	case SelectFitBuyItemTrdReq:
		atir := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
		matchingItems := screen.user.GetMatchedItems(atir.TItem)
		tableLines = screen.renderITTable(
			"Select item to sell",
			"Item",
			matchingItems,
			w, nil)
	case ShowSellItemTrdReqs:
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
	case ShowSellChrTrdReqs:
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
	case ShowBuyChrTrdReqs:
		// log.Debugln("InventoryCharacters", screen.user.InventoryCharacters())
		tableLines = screen.renderITRTable(
			"Buy character requests",
			[2]string{"Character", "Price (pylon)"},
			loud.CharacterBuyTrdReqs,
			w,
			func(idx int, request interface{}) FontType {
				itr := request.(loud.CharacterBuyTrdReq)
				// log.Debugln("GetMatchedCharacters",
				// 	len(screen.user.GetMatchedCharacters(itr.TCharacter)),
				// 	request.(loud.CharacterBuyTrdReq),
				// 	screen.user.InventoryCharacters())
				if len(screen.user.GetMatchedCharacters(itr.TCharacter)) == 0 {
					return screen.GetDisabledFontByActiveLine(idx)
				}
				return screen.getFontOfTableLine(idx, itr.IsMyTrdReq)
			})
	case SelectFitBuyChrTrdReq:
		cbtr := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
		matchingChrs := screen.user.GetMatchedCharacters(cbtr.TCharacter)
		tableLines = screen.renderITTable(
			"Select character to sell",
			"Character",
			matchingChrs,
			w, nil)
	case CreateBuyGoldTrdReqEnterPylonValue:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CreateSellGoldTrdReqEnterPylonValue:
		desc = loud.Localize("Please enter pylon amount to get (should be integer value)")
	case CreateBuyGoldTrdReqEnterGoldValue:
		desc = loud.Localize("Please enter gold amount to buy (should be integer value)")
	case SelectRenameChrEntNewName:
		desc = loud.Localize("Please enter new character's name")
	case CreateSellGoldTrdReqEnterGoldValue:
		desc = loud.Localize("Please enter gold amount to sell (should be integer value)")

	case CreateSellItemTrdReqSelectItem:
		tableLines = screen.renderITTable(
			"Select item to sell",
			"Item",
			screen.user.InventoryItems(),
			w, nil)
	case CreateSellChrTrdReqSelChr:
		tableLines = screen.renderITTable(
			"Select character to sell",
			"Character",
			screen.user.InventoryCharacters(),
			w, nil)
	case CreateSellItemTrdReqEnterPylonValue:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CreateSellChrTrdReqEnterPylonValue:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CreateBuyItemTrdReqSelectItem:
		tableLines = screen.renderITTable(
			"Select item to buy",
			"Item",
			loud.WorldItemSpecs,
			w, nil)
	case CreateBuyChrTrdReqSelectChr:
		tableLines = screen.renderITTable(
			"Select character specs to get",
			"Character",
			loud.WorldCharacterSpecs,
			w, nil)
	case CreateBuyItmTrdReqEnterPylonValue:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case CreateBuyChrTrdReqEnterPylonValue:
		desc = loud.Localize("Please enter pylon amount to use (should be integer value)")
	case SelectActiveChr:
		tableLines = screen.renderITTable(
			"Please select active character",
			"Character",
			screen.user.InventoryCharacters(),
			w, nil)
	case SelectRenameChr:
		tableLines = screen.renderITTable(
			"Please select character to rename",
			"Character",
			screen.user.InventoryCharacters(),
			w, nil)
	case SelectBuyItem:
		tableLines = screen.renderITTable(
			"select buy item desc",
			"Shop items",
			loud.ShopItems,
			w,
			func(idx int, item interface{}) FontType {
				return screen.getFontOfShopItem(idx, item.(loud.Item))
			})
	case SelectSellItem:
		tableLines = screen.renderITTable(
			"select sell item desc",
			"Item",
			screen.user.InventorySellableItems(),
			w, nil)
	case SelectUpgradeItem:
		tableLines = screen.renderITTable(
			"select upgrade item desc",
			"Item",
			screen.user.InventoryUpgradableItems(),
			w, nil)
	case SelectBuyChr:
		tableLines = screen.renderITTable(
			"select buy character desc",
			"Character",
			loud.ShopCharacters,
			w, nil)
	case ConfirmHuntRabbits:
		desc = loud.Localize("rabbits outcome")
	case ConfirmFightGoblin:
		desc = loud.Localize("goblin outcome")
		desc += carryItemDesc(activeWeapon)
	case ConfirmFightWolf:
		desc = loud.Localize("wolf outcome")
		desc += carryItemDesc(activeWeapon)
	case ConfirmFightTroll:
		desc = loud.Localize("troll outcome")
		desc += carryItemDesc(activeWeapon)
	case ConfirmFightGiant:
		desc = loud.Localize("giant outcome")
		desc += carryItemDesc(activeWeapon)
	case ConfirmFightDragonFire:
		desc = loud.Localize("fire dragon outcome")
		desc += carryItemDesc(activeWeapon)
	case ConfirmFightDragonIce:
		desc = loud.Localize("ice dragon outcome")
		desc += carryItemDesc(activeWeapon)
	case ConfirmFightDragonAcid:
		desc = loud.Localize("acid dragon outcome")
		desc += carryItemDesc(activeWeapon)
	case ConfirmFightDragonUndead:
		desc = loud.Localize("undead dragon outcome")
		desc += carryItemDesc(activeWeapon)
	case HelpAbout:
		desc = loud.Localize("help about")
	case HelpGameObjective:
		desc = loud.Localize("help what you can do")
	case HelpNavigation:
		desc = loud.Localize("help navigation")
	case HelpPageLayout:
		desc = loud.Localize("help page layout")
	case HelpGameRules:
		desc = loud.Localize("help game rules")
	case HelpHowItWorks:
		desc = loud.Localize("help how it works")
	case HelpPylonsCentral:
		tableLines = screen.tradeTableColorDesc(w)
	case HelpUpcomingReleases:
		desc = loud.Localize("help upcoming releases")
	case HelpSupport:
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
		PrintString(fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x),
			fmtFunc(fillSpace(line, w))))
		if index+2 > int(screen.Height()) {
			break
		}
	}
	infoLen := len(infoLines)

	for index, line := range tableLines {
		PrintString(fmt.Sprintf("%s%s",
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
		loud.TextRabbit:       "🐇 (rabbit)",
		loud.TextGoblin:       "👺 (goblin)",
		loud.TextWolf:         "🐺 (wolf)",
		loud.TextTroll:        "👻 (troll)",
		loud.TextGiant:        "🗿 (giant)",
		loud.TextDragonFire:   "🦐 (fire dragon)",
		loud.TextDragonIce:    "🦈 (ice dragon)",
		loud.TextDragonAcid:   "🐊 (acid dragon)",
		loud.TextDragonUndead: "🐉 (undead dragon)",
	}
	if umt, ok := unicodeMonsterTexts[monster]; ok {
		return umt
	}
	return ""
}

// GetKilledByDesc returns who killed and who was killed description
func (screen *GameScreen) GetKilledByDesc() string {
	monsterDesc := monsterTextWithUnicode(screen.user.GetTargetMonster())
	return loud.Sprintf(
		"Your %s character was killed by %s accidently",
		formatCharacterP(screen.user.GetDeadCharacter()),
		monsterDesc,
	)
}

// GetLostWeaponDesc returns description text when weapon is lost
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

// TxResultSituationDesc returns transaction result as user friendly text
func (screen *GameScreen) TxResultSituationDesc() (string, FontType) {
	desc := ""
	font := RegularFont
	resDescMap := map[PageStatus]string{
		RsltBuyGoldTrdReqCreation:  "gold buy request creation",
		RsltSellGoldTrdReqCreation: "gold sell request creation",
		RsltSelectActiveChr:        "selecting active character",
		RsltRenameChr:              "renaming character",
		RsltBuyItem:                "buy item",
		RsltBuyChr:                 "buy character",
		RsltHuntRabbits:            "hunt rabbits",
		RsltByGoldWithPylons:       "buy gold with pylons",
		RsltGetPylons:              "get pylon",
		RsltSwitchUser:             "switch user",
		RsltCreateCookbook:         "create cookbook",
		RsltSellItem:               "sell item",
		RsltUpgradeItem:            "upgrade item",
		RsltSellItemTrdReqCreation: "sell item request creation",
		RsltBuyItemTrdReqCreation:  "buy item request creation",
		RsltSellChrTrdReqCreation:  "sell character request creation",
		RsltBuyChrTrdReqCreation:   "buy character request creation",
		RsltCancelTrdReq:           "cancel trade",
		RsltFulfillBuyGoldTrdReq:   "sell loud", // for fullfill direction is reversed
		RsltFulfillSellGoldTrdReq:  "buy loud",
		RsltFulfillSellItemTrdReq:  "buy item",
		RsltFulfillSellChrTrdReq:   "buy character",
		RsltFulfillBuyItemTrdReq:   "sell item",
		RsltFulfillBuyChrTrdReq:    "sell character",
	}
	if screen.txFailReason != "" {
		desc = loud.Localize(resDescMap[screen.scrStatus]+" failure reason") + ": " + loud.Localize(screen.txFailReason)
		font = RedBoldFont
	} else {
		switch screen.scrStatus {
		case RsltBuyGoldTrdReqCreation:
			desc = loud.Localize("gold buy request was successfully created")
			desc += screen.buyLoudDesc(screen.goldEnterValue, screen.pylonEnterValue)
		case RsltSellGoldTrdReqCreation:
			desc = loud.Localize("gold sell request was successfully created")
			desc += screen.sellLoudDesc(screen.goldEnterValue, screen.pylonEnterValue)
		case RsltSelectActiveChr:
			if screen.user.GetActiveCharacter() == nil {
				desc = loud.Localize("You have successfully unset the active character!")
			} else {
				desc = loud.Localize("You have successfully set the active character!")
			}
		case RsltRenameChr:
			desc = loud.Sprintf("You have successfully updated character's name to %s!", screen.inputText)
		case RsltBuyItem:
			desc = loud.Sprintf("You have bought a weapon from the shop")
			desc += "\n"
		case RsltBuyChr:
			desc = loud.Sprintf("You have bought %s from Pylons Central", screen.activeCharacter.Name)
			desc += "\n"
			desc += loud.Localize("Please use it for hunting")
		case RsltHuntRabbits:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = loud.Sprintf("Your %s character is dead while following rabbits accidently", formatCharacter(*screen.user.GetActiveCharacter()))
				font = RedFont
			} else {
				desc = loud.Sprintf("You did hunt rabbits and earned %d.", earnedAmount)
			}
		case RsltFightGoblin:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with goblin and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				case 4:
					bonusItem := screen.user.GetItemByID(respOutput[3].ItemID)
					if bonusItem == nil {
						font = RedBoldFont
						desc = loud.Localize("Something went wrong!\nReturned ItemID is not available on user's inventory.")
					} else {
						font = GreenFont
						desc += loud.Sprintf("You got bonus item called %s", bonusItem.Name)
						desc += "\n"
						if bonusItem.Name == loud.GoblinBoots { // GOBLIN_BOOTS
							desc += loud.Sprintf("You can sell goblin boots to earn gold or trade in pylons central!")
						} else { // GOBLIN_EAR
							desc += loud.Sprintf("You can make silver sword with Goblin ear at the shop!")
						}
					}
				}
			}
		case RsltFightTroll:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with troll and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				case 4:
					bonusItem := screen.user.GetItemByID(respOutput[3].ItemID)
					if bonusItem == nil {
						font = RedBoldFont
						desc = loud.Localize("Something went wrong!\nReturned ItemID is not available on user's inventory.")
					} else {
						font = GreenFont
						desc += loud.Sprintf("You got bonus item called %s", bonusItem.Name)
						desc += "\n"
						if bonusItem.Name == loud.TrollSmellyBones { // troll smelly bones
							desc += loud.Sprintf("You can sell troll's smelly boots at to earn gold or trade in pylons central!")
						} else { // TROLL_TOES
							desc += loud.Sprintf("You can make iron sword with Troll toes at the shop!")
						}
					}
				}
			}
		case RsltFightWolf:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with wolf and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				case 4:
					bonusItem := screen.user.GetItemByID(respOutput[3].ItemID)
					if bonusItem == nil {
						font = RedBoldFont
						desc = loud.Localize("Something went wrong!\nReturned ItemID is not available on user's inventory.")
					} else {
						font = GreenFont
						desc += loud.Sprintf("You got bonus item called %s", bonusItem.Name)
						desc += "\n"
						if bonusItem.Name == loud.WolfFur { // wolf fur
							desc += loud.Sprintf("You can sell wolf fur at to earn gold or trade in pylons central!")
						} else { // WOLF_TAIL
							desc += loud.Sprintf("You can make bronze sword with Wolf tail at the shop!")
						}
					}
				}
			}
		case RsltFightGiant:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with giant and earned %d.", earnedAmount)
				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				case 3:
					activeCharacter := screen.user.GetActiveCharacter()
					if activeCharacter != nil && activeCharacter.Special != loud.NoSpecial { // Got special from this fight
						desc += loud.Sprintf("You got %s (special) from the giant!!", formatSpecial(activeCharacter.Special))
						desc += "\n"
						desc += loud.Sprintf("You can now fight with %s with this character!", formatSpecialDragon(activeCharacter.Special))
						font = GreenFont
					}
				}
			}
		case RsltFightDragonFire:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with fire dragon and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.DropDragonFire)
					desc += "\n"
					desc += loud.Sprintf("Once you have drops from 3 special dragons, you can create angel sword at the shop!")
					font = GreenFont
				}
			}
		case RsltFightDragonIce:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with ice dragon and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.DropDragonIce)
					desc += "\n"
					desc += loud.Sprintf("Once you have drops from 3 special dragons, you can create angel sword at the shop!")
					font = GreenFont
				}
			}
		case RsltFightDragonAcid:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with acid dragon and earned %d.", earnedAmount)

				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				case 4:
					desc += loud.Sprintf("You got bonus item called %s", loud.DropDragonAcid)
					desc += "\n"
					desc += loud.Sprintf("Once you have drops from 3 special dragons, you can create angel sword at the shop!")
					font = GreenFont
				}
			}
		case RsltFightDragonUndead:
			earnedAmount, respOutput := screen.GetTxResponseOutput()
			resLen := len(respOutput)
			if resLen == 0 {
				desc = screen.GetKilledByDesc()
				font = RedFont
			} else {
				desc = loud.Sprintf("You did fight with undead dragon and earned %d.", earnedAmount)
				switch resLen {
				case 2:
					desc += screen.GetLostWeaponDesc()
					font = YelloFont
				}
			}
		case RsltByGoldWithPylons:
			earnedAmount, _ := screen.GetTxResponseOutput()
			desc = loud.Sprintf("Bought %d gold with %d pylons.", earnedAmount, 100)
		case RsltDevGetTestItems:
			_, respOutput := screen.GetTxResponseOutput()
			resultTexts := []string{
				loud.WolfTail,
				loud.TrollToes,
				loud.GoblinBoots,
				loud.DropDragonFire,
				loud.DropDragonIce,
				loud.DropDragonAcid,
				"Ruppell's Fox",
				"Gentoo penguin",
				"Colorado River toad",
			}
			desc = loud.Sprintf("Finished getting developer test items.")
			desc += devDetailedResultDesc(resultTexts[:len(respOutput)])
		case RsltGetPylons:
			desc = loud.Localize("You got extra pylons for LOUD game")
		case RsltSwitchUser:
			desc = loud.Sprintf("You switched user to %s", screen.user.GetUserName())
		case RsltCreateCookbook:
			desc = loud.Localize("You created a new cookbook for a new game build")
		case RsltSellItem:
			earnedAmount, _ := screen.GetTxResponseOutput()
			desc = loud.Sprintf("You sold %s for %d gold.", screen.activeItem.Name, earnedAmount)
		case RsltUpgradeItem:
			desc = loud.Sprintf("You have upgraded %s to get better hunt result", screen.activeItem.Name)
		case RsltSellItemTrdReqCreation:
			desc = loud.Localize("item sell request was successfully created")
			desc += screen.sellItemDesc(screen.activeItem, screen.pylonEnterValue)
		case RsltSellChrTrdReqCreation:
			desc = loud.Localize("character sell request was successfully created")
			desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
		case RsltBuyItemTrdReqCreation:
			desc = loud.Localize("item buy request was successfully created")
			desc += screen.buyItemSpecDesc(screen.activeItSpec, screen.pylonEnterValue)
		case RsltBuyChrTrdReqCreation:
			desc = loud.Localize("character buy request was successfully created")
			desc += screen.buyCharacterSpecDesc(screen.activeChSpec, screen.pylonEnterValue)
		case RsltCancelTrdReq:
			desc = loud.Localize("successfully cancelled trade request")
		case RsltFulfillBuyGoldTrdReq:
			request := screen.activeTrdReq
			desc = loud.Sprintf("you have sold gold successfully from coin market at %.4f", request.Price)
			desc += screen.sellLoudDesc(request.Amount, request.Total)
		case RsltFulfillSellGoldTrdReq:
			request := screen.activeTrdReq
			desc = loud.Sprintf("you have bought gold successfully from coin market at %.4f", request.Price)
			desc += screen.buyLoudDesc(request.Amount, request.Total)
		case RsltFulfillSellItemTrdReq:
			request := screen.activeItemTrdReq.(loud.ItemSellTrdReq)
			desc = loud.Localize("you have bought item successfully from item/pylon market")
			desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RsltFulfillSellChrTrdReq:
			request := screen.activeItemTrdReq.(loud.CharacterSellTrdReq)
			desc = loud.Localize("you have bought character successfully from character/pylon market")
			desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		case RsltFulfillBuyItemTrdReq:
			request := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
			desc = loud.Localize("you have sold item successfully from item/pylon market")
			desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
		case RsltFulfillBuyChrTrdReq:
			request := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
			desc = loud.Localize("you have sold character successfully from character/pylon market")
			desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
		}
	}
	return desc, font
}

// TxWaitSituationDesc returns wait text for a pending transaction
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
	WaitToEnd := "\n" + loud.Localize("Please wait for a moment to finish the process")
	switch screen.scrStatus {
	case WaitRenameChr:
		desc = loud.Sprintf("You are waiting to rename character from %s to %s.", screen.activeCharacter.Name, screen.inputText)
	case WaitBuyGoldTrdReqCreation:
		desc = loud.Localize("You are waiting for gold buy request creation")
		desc += screen.buyLoudDesc(screen.goldEnterValue, screen.pylonEnterValue)
	case WaitSellGoldTrdReqCreation:
		desc = loud.Localize("You are waiting for gold sell request creation")
		desc += screen.sellLoudDesc(screen.goldEnterValue, screen.pylonEnterValue)
	case WaitBuyItem:
		desc = loud.Sprintf("You are buying %s at the shop", screen.activeItem.Name)
	case WaitBuyChr:
		desc = loud.Sprintf("You are buying %s at the pylons central", formatCharacter(screen.activeCharacter))
		desc += WaitToEnd
	case WaitHuntRabbits:
		desc = loud.Sprintf("You are hunting rabbits")
		desc += WaitToEnd
	case WaitFightTroll,
		WaitFightWolf,
		WaitFightGoblin,
		WaitFightGiant,
		WaitFightDragonFire,
		WaitFightDragonIce,
		WaitFightDragonAcid,
		WaitFightDragonUndead:
		desc = loud.Sprintf("You are fighting %s monster with weapon %s level %d", monsterName, activeWeaponName, activeWeaponLevel)
	case WaitByGoldWithPylons:
		desc = loud.Sprintf("spending %d pylon for %d gold", 100, 5000)
	case WaitDevGetTestItems:
		desc = loud.Sprintf("Getting dev test items from pylon")
		desc += WaitToEnd
	case WaitGetPylons:
		desc = loud.Sprintf("You are waiting for getting pylons process")
	case WaitSwitchUser:
		desc = loud.Sprintf("You are waiting for switching to new user")
	case WaitCreateCookbook:
		desc = loud.Sprintf("You are waiting for creating cookbook")
	case WaitSellItem:
		item := screen.activeItem
		desc = loud.Sprintf("You are selling %s for %s gold", item.Name, item.GetSellPriceRange())
	case WaitUpgradeItem:
		desc = loud.Sprintf("You are upgrading %s", screen.activeItem.Name)
	case WaitSellItemTrdReqCreation:
		desc = loud.Localize("You are waiting for item sell request creation")
		desc += screen.sellItemDesc(screen.activeItem, screen.pylonEnterValue)
	case WaitSellChrTrdReqCreation:
		desc = loud.Localize("You are waiting for character sell request creation")
		desc += screen.sellCharacterDesc(screen.activeCharacter, screen.pylonEnterValue)
	case WaitBuyItemTrdReqCreation:
		desc = loud.Localize("You are waiting for item buy request creation")
		desc += screen.buyItemSpecDesc(screen.activeItSpec, screen.pylonEnterValue)
	case WaitBuyChrTrdReqCreation:
		desc = loud.Localize("You are waiting for character buy request creation")
		desc += screen.buyCharacterSpecDesc(screen.activeChSpec, screen.pylonEnterValue)
	case WaitCancelTrdReq:
		desc = loud.Localize("You are waiting for cancelling one of your trades")
	// For FULFILL trades, msg should be reversed, since user is opposite
	case WaitFulfillSellItemTrdReq:
		request := screen.activeItemTrdReq.(loud.ItemSellTrdReq)
		desc = loud.Sprintf("You are buying item at %d", request.Price)
		desc += screen.buyItemDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case WaitFulfillSellChrTrdReq:
		request := screen.activeItemTrdReq.(loud.CharacterSellTrdReq)
		desc = loud.Sprintf("You are buying character at %d.", request.Price)
		desc += screen.buyCharacterDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case WaitFulfillBuyItemTrdReq:
		request := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)
		desc = loud.Sprintf("You are selling item at %d.", request.Price)
		desc += screen.sellItemSpecDesc(request.TItem, fmt.Sprintf("%d", request.Price))
	case WaitFulfillBuyChrTrdReq:
		request := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)
		desc = loud.Sprintf("You are selling character at %d.", request.Price)
		desc += screen.sellCharacterSpecDesc(request.TCharacter, fmt.Sprintf("%d", request.Price))
	case WaitFulfillBuyGoldTrdReq:
		request := screen.activeTrdReq
		desc = loud.Sprintf("Making pylons from gold")
		desc += screen.sellLoudDesc(request.Amount, request.Total)
	case WaitFulfillSellGoldTrdReq:
		request := screen.activeTrdReq
		desc = loud.Sprintf("Making gold from pylons")
		desc += screen.buyLoudDesc(request.Amount, request.Total)
	}
	desc += "\n"
	return loud.ChunkText(desc, width-2), TextLines{}.appendF(fillSpace("......", width), BlinkBlueFont)
}
