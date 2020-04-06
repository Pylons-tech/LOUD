package loud

import (
	"log"
	"sort"
	"strings"

	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
)

func SyncFromNode(user User) {
	log.Println("SyncFromNode Function Body")
	log.Println("username=", user.GetUserName())
	log.Println("userinfo=", pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT()))
	accAddr := pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT())
	accInfo := pylonSDK.GetAccountInfoFromName(user.GetUserName(), GetTestingT())
	log.Println("accountInfo=", accInfo)

	user.SetGold(int(accInfo.Coins.AmountOf("loudcoin").Int64()))
	user.SetPylonAmount(int(accInfo.Coins.AmountOf("pylon").Int64()))
	log.Println("gold=", accInfo.Coins.AmountOf("loudcoin").Int64())

	rawItems, _ := pylonSDK.ListItemsViaCLI(accInfo.Address.String())
	myItems := []Item{}
	myCharacters := []Character{}
	for _, rawItem := range rawItems {
		XP, _ := rawItem.FindDouble("XP")
		HP, _ := rawItem.FindLong("HP")
		MaxHP, _ := rawItem.FindLong("MaxHP")
		Level, _ := rawItem.FindLong("level")
		Name, _ := rawItem.FindString("Name")
		itemType, _ := rawItem.FindString("Type")
		Attack, _ := rawItem.FindDouble("attack")

		if itemType == "Character" {
			myCharacters = append(myCharacters, Character{
				Level: Level,
				Name:  Name,
				ID:    rawItem.ID,
				XP:    XP,
				HP:    HP,
				MaxHP: MaxHP,
			})
		} else {
			myItems = append(myItems, Item{
				Level:  Level,
				Name:   Name,
				Attack: int(Attack),
				ID:     rawItem.ID,
			})
		}
	}
	user.SetItems(myItems)
	user.SetCharacters(myCharacters)
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
		if tradeItem.Completed == false && strings.Contains(tradeItem.ExtraInfo, CR8BY_LOUD) {
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
				inputPylonAmount := tradeItem.CoinInputs[0].Count
				nSellTrdReqs = append(nSellTrdReqs, TrdReq{
					ID:         tradeItem.ID,
					Amount:     int(loudOutputAmount),
					Total:      int(inputPylonAmount),
					Price:      float64(inputPylonAmount) / float64(loudOutputAmount),
					IsMyTrdReq: isMyTrdReq,
				})
			} else if itemInputLen > 0 { // buy item trade
				MinLevel := tradeItem.ItemInputs[0].Longs[0].MinValue
				MaxLevel := tradeItem.ItemInputs[0].Longs[0].MaxValue
				Name := tradeItem.ItemInputs[0].Strings[0].Value
				if tradeItem.ExtraInfo == ITEM_BUYREQ_TRDINFO {
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
				} else if tradeItem.ExtraInfo == CHAR_BUYREQ_TRDINFO { // character buy request created by loud game
					MinXP := 0.0
					MaxXP := 0.0
					if len(tradeItem.ItemInputs[0].Doubles) > 0 {
						MinXP = tradeItem.ItemInputs[0].Doubles[0].MinValue.Float()
						MaxXP = tradeItem.ItemInputs[0].Doubles[0].MaxValue.Float()
					}
					tCharacter := CharacterSpec{
						Level: [2]int{MinLevel, MaxLevel},
						Name:  Name,
						XP:    [2]float64{MinXP, MaxXP},
					}
					nBuyCharacterTrdReqs = append(nBuyCharacterTrdReqs, CharacterBuyTrdReq{
						ID:         tradeItem.ID,
						TCharacter: tCharacter,
						Price:      int(pylonOutputAmount),
						IsMyTrdReq: isMyTrdReq,
					})
				}
			} else if itemOutputLen > 0 { // sell item trade
				inputPylonAmount := tradeItem.CoinInputs[0].Count
				level, _ := tradeItem.ItemOutputs[0].FindLong("level")
				name, _ := tradeItem.ItemOutputs[0].FindString("Name")
				if tradeItem.ExtraInfo == ITEM_SELREQ_TRDINFO {
					tItem := Item{
						ID:    tradeItem.ItemOutputs[0].ID,
						Level: level,
						Name:  name,
					}
					nSellItemTrdReqs = append(nSellItemTrdReqs, ItemSellTrdReq{
						ID:         tradeItem.ID,
						TItem:      tItem,
						Price:      int(inputPylonAmount),
						IsMyTrdReq: isMyTrdReq,
					})
				} else if tradeItem.ExtraInfo == CHAR_SELREQ_TRDINFO { // character sell request created by loud game
					XP, _ := tradeItem.ItemOutputs[0].FindDouble("XP")
					HP, _ := tradeItem.ItemOutputs[0].FindLong("HP")
					MaxHP, _ := tradeItem.ItemOutputs[0].FindLong("MaxHP")
					tCharacter := Character{
						ID:    tradeItem.ItemOutputs[0].ID,
						Level: level,
						Name:  name,
						XP:    XP,
						HP:    HP,
						MaxHP: MaxHP,
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
}
