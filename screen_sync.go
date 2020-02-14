package loud

import (
	"log"
	"sort"

	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
)

func SyncFromNode(user User) {
	log.Println("SyncFromNode username=", user.GetUserName())
	log.Println("SyncFromNode userinfo=", pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT()))
	accAddr := pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT())
	accInfo := pylonSDK.GetAccountInfoFromName(user.GetUserName(), GetTestingT())
	log.Println("accountInfo Result=", accInfo)

	user.SetGold(int(accInfo.Coins.AmountOf("loudcoin").Int64()))
	user.SetPylonAmount(int(accInfo.Coins.AmountOf("pylon").Int64()))
	log.Println("SyncFromNode gold=", accInfo.Coins.AmountOf("loudcoin").Int64())

	rawItems, _ := pylonSDK.ListItemsViaCLI(accInfo.Address.String())
	myItems := []Item{}
	for _, rawItem := range rawItems {
		Level, _ := rawItem.FindLong("level")
		Name, _ := rawItem.FindString("Name")
		item := Item{
			Level: Level,
			Name:  Name,
			ID:    rawItem.ID,
		}
		myItems = append(myItems, item)
	}
	user.SetItems(myItems)
	log.Println("SyncFromNode myItems=", myItems)

	nBuyTradeRequests := []TradeRequest{}
	nSellTradeRequests := []TradeRequest{}
	nBuySwordTradeRequests := []ItemTradeRequest{}
	nSellSwordTradeRequests := []ItemTradeRequest{}
	rawTrades, _ := pylonSDK.ListTradeViaCLI("")
	for _, tradeItem := range rawTrades {
		if tradeItem.Completed == false {
			inputCoin := ""
			if len(tradeItem.CoinInputs) > 0 {
				inputCoin = tradeItem.CoinInputs[0].Coin
			}
			loudOutputAmount := tradeItem.CoinOutputs.AmountOf("loudcoin").Int64()
			pylonOutputAmount := tradeItem.CoinOutputs.AmountOf("pylon").Int64()
			itemInputLen := len(tradeItem.ItemInputs)
			itemOutputLen := len(tradeItem.ItemOutputs)
			isMyTradeRequest := tradeItem.Sender.String() == accAddr
			if inputCoin == "loudcoin" { // loud sell trade
				loudAmount := tradeItem.CoinInputs[0].Count

				nBuyTradeRequests = append(nBuyTradeRequests, TradeRequest{
					ID:               tradeItem.ID,
					Amount:           int(loudAmount),
					Total:            int(pylonOutputAmount),
					Price:            float64(pylonOutputAmount) / float64(loudAmount),
					IsMyTradeRequest: isMyTradeRequest,
				})
			} else if loudOutputAmount > 0 { // loud buy trade
				inputPylonAmount := tradeItem.CoinInputs[0].Count
				nSellTradeRequests = append(nSellTradeRequests, TradeRequest{
					ID:               tradeItem.ID,
					Amount:           int(loudOutputAmount),
					Total:            int(inputPylonAmount),
					Price:            float64(inputPylonAmount) / float64(loudOutputAmount),
					IsMyTradeRequest: isMyTradeRequest,
				})
			} else if itemInputLen > 0 { // buy sword trade
				tItem := Item{
					Level: tradeItem.ItemInputs[0].Longs[0].MinValue, // Level
					Name:  tradeItem.ItemInputs[0].Strings[0].Value,
				}
				nBuySwordTradeRequests = append(nBuySwordTradeRequests, ItemTradeRequest{
					ID:               tradeItem.ID,
					TItem:            tItem,
					Price:            int(pylonOutputAmount),
					IsMyTradeRequest: isMyTradeRequest,
				})
			} else if itemOutputLen > 0 { // sell sword trade
				inputPylonAmount := tradeItem.CoinInputs[0].Count
				level, _ := tradeItem.ItemOutputs[0].FindLong("level")
				name, _ := tradeItem.ItemOutputs[0].FindString("Name")
				tItem := Item{
					ID:    tradeItem.ItemOutputs[0].ID,
					Level: level,
					Name:  name,
				}
				nSellSwordTradeRequests = append(nSellSwordTradeRequests, ItemTradeRequest{
					ID:               tradeItem.ID,
					TItem:            tItem,
					Price:            int(inputPylonAmount),
					IsMyTradeRequest: isMyTradeRequest,
				})
			}
		}
	}
	// Sort and show by low price buy requests
	sort.SliceStable(nBuyTradeRequests, func(i, j int) bool {
		return nBuyTradeRequests[i].Price < nBuyTradeRequests[j].Price
	})
	// Sort and show by high price sell requests
	sort.SliceStable(nSellTradeRequests, func(i, j int) bool {
		return nSellTradeRequests[i].Price > nSellTradeRequests[j].Price
	})
	buyTradeRequests = nBuyTradeRequests
	sellTradeRequests = nSellTradeRequests
	swordBuyTradeRequests = nBuySwordTradeRequests
	swordSellTradeRequests = nSellSwordTradeRequests
	log.Println("SyncFromNode buyTradeRequests=", buyTradeRequests)
	log.Println("SyncFromNode sellTradeRequests=", sellTradeRequests)
}
