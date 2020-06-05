package loud

import (
	"sort"
	"strings"

	"github.com/Pylons-tech/LOUD/log"
	pylonSDK "github.com/Pylons-tech/pylons_sdk/cmd/test"
)

var isSyncingFromNode = false

// SyncFromNode is a function to handle sync from node
func SyncFromNode(user User) {
	if isSyncingFromNode {
		return
	}
	isSyncingFromNode = true
	defer func() {
		isSyncingFromNode = false
	}()
	log.Println("SyncFromNode Function Body")
	log.Println("username=", user.GetUserName())
	log.Println("userinfo=", pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT()))
	accAddr := pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT())
	accInfo := pylonSDK.GetAccountInfoFromName(user.GetUserName(), GetTestingT())
	log.Println("accountInfo=", accInfo)

	user.SetGold(int(accInfo.Coins.AmountOf("loudcoin").Int64()))
	log.Println("gold=", accInfo.Coins.AmountOf("loudcoin").Int64())
	user.SetPylonAmount(int(accInfo.Coins.AmountOf("pylon").Int64()))
	user.SetAddress(accAddr)

	rawItems, _ := pylonSDK.ListItemsViaCLI(accInfo.Address.String())
	myItems := []Item{}
	myCharacters := []Character{}
	for _, rawItem := range rawItems {
		if rawItem.CookbookID != GameCookbookID {
			continue
		}
		XP, _ := rawItem.FindDouble("XP")
		Level, _ := rawItem.FindLong("level")
		GiantKill, _ := rawItem.FindLong("GiantKill") // 🗿
		Special, _ := rawItem.FindLong("Special")
		SpecialDragonKill, _ := rawItem.FindLong("SpecialDragonKill")
		UndeadDragonKill, _ := rawItem.FindLong("UndeadDragonKill")

		Name, _ := rawItem.FindString("Name")
		itemType, _ := rawItem.FindString("Type")
		Attack, _ := rawItem.FindDouble("attack")
		Value, _ := rawItem.FindLong("value")
		LastUpdate := rawItem.LastUpdate

		if itemType == "Character" {
			myCharacters = append(myCharacters, Character{
				Level:             Level,
				Name:              Name,
				ID:                rawItem.ID,
				XP:                XP,
				GiantKill:         GiantKill,
				Special:           Special,
				SpecialDragonKill: SpecialDragonKill,
				UndeadDragonKill:  UndeadDragonKill,
				LastUpdate:        LastUpdate,
			})
		} else {
			myItems = append(myItems, Item{
				Level:      Level,
				Name:       Name,
				Attack:     int(Attack),
				Value:      Value,
				ID:         rawItem.ID,
				LastUpdate: LastUpdate,
			})
		}
	}
	// Sort characters by dragon kills, giant kill, special, level and name
	sort.SliceStable(myCharacters, func(i, j int) bool {
		chi := myCharacters[i]
		chj := myCharacters[j]
		if chi.UndeadDragonKill != chj.UndeadDragonKill {
			return chi.UndeadDragonKill > chj.UndeadDragonKill
		} else if chi.SpecialDragonKill != chj.SpecialDragonKill {
			return chi.SpecialDragonKill > chj.SpecialDragonKill
		} else if chi.GiantKill != chj.GiantKill {
			return chi.GiantKill > chj.GiantKill
		} else if chi.Special != chj.Special {
			return chi.Special > chj.Special
		} else if chi.Level != chj.Level {
			return chi.Level > chj.Level
		} else {
			return chi.Name > chj.Name
		}
	})
	user.SetCharacters(myCharacters)

	// Sort items by attack and name
	sort.SliceStable(myItems, func(i, j int) bool {
		iti := myItems[i]
		itj := myItems[j]
		if iti.Attack != itj.Attack {
			return iti.Attack > itj.Attack
		}
		return iti.Name > itj.Name
	})
	user.SetItems(myItems)

	log.Println("myItems=", myItems)
	log.Println("myCharacters=", myCharacters)

	nBuyTrdReqs := []TrdReq{}
	nSellTrdReqs := []TrdReq{}
	nBuyItemTrdReqs := []ItemBuyTrdReq{}
	nSellItemTrdReqs := []ItemSellTrdReq{}
	nBuyCharacterTrdReqs := []CharacterBuyTrdReq{}
	nSellCharacterTrdReqs := []CharacterSellTrdReq{}
	rawTrades, _ := pylonSDK.ListTradeViaCLI("")
	for _, tradeItem := range rawTrades {
		if !tradeItem.Completed && !tradeItem.Disabled && strings.Contains(tradeItem.ExtraInfo, TextCreatedByLOUD) {
			inputCoin := ""
			if len(tradeItem.CoinInputs) > 0 {
				inputCoin = tradeItem.CoinInputs[0].Coin
			}
			loudOutputAmount := tradeItem.CoinOutputs.AmountOf("loudcoin").Int64()
			pylonOutputAmount := tradeItem.CoinOutputs.AmountOf("pylon").Int64()
			itemInputLen := len(tradeItem.ItemInputs)
			itemOutputLen := len(tradeItem.ItemOutputs)
			isMyTrdReq := tradeItem.Sender.String() == accAddr
			if inputCoin == "loudcoin" { // loud sell trade
				loudAmount := tradeItem.CoinInputs[0].Count

				nBuyTrdReqs = append(nBuyTrdReqs, TrdReq{
					ID:         tradeItem.ID,
					Amount:     int(loudAmount),
					Total:      int(pylonOutputAmount),
					Price:      float64(pylonOutputAmount) / float64(loudAmount),
					IsMyTrdReq: isMyTrdReq,
				})
			} else if loudOutputAmount > 0 { // loud buy trade
				if len(tradeItem.CoinInputs) > 0 {
					inputPylonAmount := tradeItem.CoinInputs[0].Count
					nSellTrdReqs = append(nSellTrdReqs, TrdReq{
						ID:         tradeItem.ID,
						Amount:     int(loudOutputAmount),
						Total:      int(inputPylonAmount),
						Price:      float64(inputPylonAmount) / float64(loudOutputAmount),
						IsMyTrdReq: isMyTrdReq,
					})
				}
			} else if itemInputLen > 0 { // buy item trade
				MinLevel := 0
				MaxLevel := 0
				firstItemInput := tradeItem.ItemInputs[0]
				if len(firstItemInput.Longs) > 0 {
					MinLevel = firstItemInput.Longs[0].MinValue
					MaxLevel = firstItemInput.Longs[0].MaxValue
				}
				Name := firstItemInput.Strings[0].Value
				if tradeItem.ExtraInfo == ItemBuyReqTrdInfo {
					tItem := ItemSpec{
						Level: [2]int{MinLevel, MaxLevel},
						Name:  Name,
					}
					nBuyItemTrdReqs = append(nBuyItemTrdReqs, ItemBuyTrdReq{
						ID:         tradeItem.ID,
						TItem:      tItem,
						Price:      int(pylonOutputAmount),
						IsMyTrdReq: isMyTrdReq,
					})
				} else if tradeItem.ExtraInfo == ChrBuyReqTrdInfo { // character buy request created by loud game
					Special := NoSpecial
					if len(firstItemInput.Longs) > 1 {
						Special = firstItemInput.Longs[1].MinValue
					}
					MinXP := 0.0
					MaxXP := 0.0
					if len(firstItemInput.Doubles) > 0 {
						MinXP = firstItemInput.Doubles[0].MinValue.Float()
						MaxXP = firstItemInput.Doubles[0].MaxValue.Float()
					}
					tCharacter := CharacterSpec{
						Special: Special,
						Level:   [2]int{MinLevel, MaxLevel},
						Name:    Name,
						XP:      [2]float64{MinXP, MaxXP},
					}
					nBuyCharacterTrdReqs = append(nBuyCharacterTrdReqs, CharacterBuyTrdReq{
						ID:         tradeItem.ID,
						TCharacter: tCharacter,
						Price:      int(pylonOutputAmount),
						IsMyTrdReq: isMyTrdReq,
					})
				}
			} else if itemOutputLen > 0 { // sell item trade
				firstItemOutput := tradeItem.ItemOutputs[0]
				inputPylonAmount := int64(0)
				if len(tradeItem.CoinInputs) > 0 {
					inputPylonAmount = tradeItem.CoinInputs[0].Count
				}
				level, _ := firstItemOutput.FindLong("level")
				name, _ := firstItemOutput.FindString("Name")
				special, _ := firstItemOutput.FindLong("Special")
				GiantKill, _ := firstItemOutput.FindLong("GiantKill")
				SpecialDragonKill, _ := firstItemOutput.FindLong("SpecialDragonKill")
				UndeadDragonKill, _ := firstItemOutput.FindLong("UndeadDragonKill")

				if tradeItem.ExtraInfo == ItemSellReqTrdInfo {
					tItem := Item{
						ID:    firstItemOutput.ID,
						Level: level,
						Name:  name,
					}
					nSellItemTrdReqs = append(nSellItemTrdReqs, ItemSellTrdReq{
						ID:         tradeItem.ID,
						TItem:      tItem,
						Price:      int(inputPylonAmount),
						IsMyTrdReq: isMyTrdReq,
					})
				} else if tradeItem.ExtraInfo == ChrSellReqTrdInfo { // character sell request created by loud game
					XP, _ := firstItemOutput.FindDouble("XP")
					tCharacter := Character{
						ID:                firstItemOutput.ID,
						Level:             level,
						Name:              name,
						XP:                XP,
						Special:           special,
						GiantKill:         GiantKill,
						SpecialDragonKill: SpecialDragonKill,
						UndeadDragonKill:  UndeadDragonKill,
					}
					nSellCharacterTrdReqs = append(nSellCharacterTrdReqs, CharacterSellTrdReq{
						ID:         tradeItem.ID,
						TCharacter: tCharacter,
						Price:      int(inputPylonAmount),
						IsMyTrdReq: isMyTrdReq,
					})
				}
			}
		}
	}
	// Sort and show by low price buy requests
	sort.SliceStable(nBuyTrdReqs, func(i, j int) bool {
		return nBuyTrdReqs[i].Price < nBuyTrdReqs[j].Price
	})
	// Sort and show by high price sell requests
	sort.SliceStable(nSellTrdReqs, func(i, j int) bool {
		return nSellTrdReqs[i].Price > nSellTrdReqs[j].Price
	})
	BuyTrdReqs = nBuyTrdReqs
	SellTrdReqs = nSellTrdReqs
	ItemBuyTrdReqs = nBuyItemTrdReqs
	ItemSellTrdReqs = nSellItemTrdReqs
	CharacterBuyTrdReqs = nBuyCharacterTrdReqs
	CharacterSellTrdReqs = nSellCharacterTrdReqs
	log.Println("BuyTrdReqs=", BuyTrdReqs)
	log.Println("SellTrdReqs=", SellTrdReqs)

	user.FixLoadedData()

	ds, _, err := pylonSDK.GetDaemonStatus()
	if err == nil {
		user.SetLatestBlockHeight(ds.SyncInfo.LatestBlockHeight)
	}
}
