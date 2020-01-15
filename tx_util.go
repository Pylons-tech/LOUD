package loud

import (
	"encoding/json"
	"log"
	originT "testing"

	fixtureSDK "github.com/MikeSofaer/pylons/cmd/fixtures_test"
	testing "github.com/MikeSofaer/pylons/cmd/fixtures_test/evtesting"
	pylonSDK "github.com/MikeSofaer/pylons/cmd/test"
	"github.com/MikeSofaer/pylons/x/pylons/handlers"
	"github.com/MikeSofaer/pylons/x/pylons/msgs"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Weapon int

const (
	NO_WEAPON Weapon = iota
	WOODEN_SWORD_LV1
	WOODEN_SWORD_LV2
	COPPER_SWORD_LV1
	COPPER_SWORD_LV2
)

var RcpIDs map[string]string = map[string]string{
	"LOUD's Copper sword lv1 buy recipe":            "LOUD-copper-sword-lv1-buy-recipe-v0.0.0-1579053457",
	"LOUD's get initial coin recipe":                "LOUD-get-initial-coin-recipe-v0.0.0-1579053457",
	"LOUD's hunt with lv1 copper sword recipe":      "LOUD-hunt-with-copper-sword-lv1-recipe-v0.0.0-1579053457",
	"LOUD's hunt with lv2 copper sword recipe":      "LOUD-hunt-with-copper-sword-lv2-recipe-v0.0.0-1579053457",
	"LOUD's hunt without sword recipe":              "LOUD-hunt-with-no-weapon-recipe-v0.0.0-1579053457",
	"LOUD's hunt with lv1 wooden sword recipe":      "LOUD-hunt-with-wooden-sword-lv1-recipe-v0.0.0-1579053457",
	"LOUD's hunt with lv2 wooden sword recipe":      "LOUD-hunt-with-wooden-sword-lv2-recipe-v0.0.0-1579053457",
	"LOUD's Lv1 copper sword sell recipe":           "LOUD-sell-copper-sword-lv1-recipe-v0.0.0-1579053457",
	"LOUD's Lv2 copper sword sell recipe":           "LOUD-sell-copper-sword-lv2-recipe-v0.0.0-1579053457",
	"LOUD's Lv1 wooden sword sell recipe":           "LOUD-sell-wooden-sword-lv1-recipe-v0.0.0-1579053457",
	"LOUD's Lv2 wooden sword sell recipe":           "LOUD-sell-wooden-sword-lv2-recipe-v0.0.0-1579053457",
	"LOUD's Copper sword lv1 to lv2 upgrade recipe": "LOUD-upgrade-copper-sword-lv1-to-lv2-recipe-v0.0.0-1579053457",
	"LOUD's Wooden sword lv1 to lv2 upgrade recipe": "LOUD-upgrade-wooden-sword-lv1-to-lv2-recipe-v0.0.0-1579053457",
	"LOUD's Wooden sword lv1 buy recipe":            "LOUD-wooden-sword-lv1-buy-recipe-v0.0.0-1579053457",
}

func SyncFromNode(user User) {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT

	accInfo := pylonSDK.GetAccountInfoFromName("eugen", t)
	user.SetGold(int(accInfo.Coins.AmountOf("loudcoin").Int64()))
	log.Println("SyncFromNode gold=", accInfo.Coins.AmountOf("loudcoin").Int64())

	rawItems, _ := pylonSDK.ListItemsViaCLI(accInfo.Address.String())
	items := []Item{}
	for _, rawItem := range rawItems {
		Level, _ := rawItem.FindLong("level")
		Name, _ := rawItem.FindString("Name")
		items = append(items, Item{
			Level: Level,
			Name:  Name,
			ID:    rawItem.ID,
		})
	}
	user.SetItems(items)
	log.Println("SyncFromNode items=", items)
}

func ProcessTxResult(user User, txhash string) handlers.ExecuteRecipeSerialize {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT

	txHandleResBytes, err := pylonSDK.WaitAndGetTxData(txhash, 3, t)
	pylonSDK.ErrValidation(t, "error getting tx result bytes %+v", err)

	fixtureSDK.CheckErrorOnTx(txhash, t)
	resp := handlers.ExecuteRecipeResp{}
	respOutput := handlers.ExecuteRecipeSerialize{}
	err = pylonSDK.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
	if err != nil {
		log.Println("failed to parse transaction result txhash=", txhash)
	}

	json.Unmarshal(resp.Output, &respOutput)
	log.Println("ProcessTxResult::txResp", resp.Message, respOutput)
	SyncFromNode(user)
	return respOutput
}

func GetTestingT() *testing.T {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT
	return t
}

func ExecuteRecipe(user User, rcpName string, itemIDs []string) string {
	t := GetTestingT()

	rcpID := RcpIDs[rcpName]
	eugenAddr := pylonSDK.GetAccountAddr("eugen", nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)
	// execMsg := msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, itemIDs)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, "eugen", false)
	user.SetLastTransaction(txhash)
	return txhash
}

func GetWeaponItemFromKey(user User, key string) Item {
	items := user.InventoryItems()
	useItem := Item{}
	switch key {
	case "1": // SELECT 1st item
		useItem = items[0]
	case "2": // SELECT 2nd item
		useItem = items[1]
	case "3": // SELECT 3rd item
		useItem = items[2]
	case "4": // SELECT 4th item
		useItem = items[3]
	}
	return useItem
}

func Hunt(user User, key string) string {
	rcpName := "LOUD's hunt without sword recipe"

	useItem := GetWeaponItemFromKey(user, key)
	itemIDs := []string{}
	switch key {
	case "I": // get initial coin
		fallthrough
	case "i":
		rcpName = "LOUD's get initial coin recipe"
	}

	switch useItem.Name {
	case "Wooden sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's hunt with lv1 wooden sword recipe"
		} else {
			rcpName = "LOUD's hunt with lv2 wooden sword recipe"
		}
		itemIDs = []string{useItem.ID}
	case "Copper sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's hunt with lv1 copper sword recipe"
		} else {
			rcpName = "LOUD's hunt with lv2 copper sword recipe"
		}
		itemIDs = []string{useItem.ID}
	}

	return ExecuteRecipe(user, rcpName, itemIDs)
}

func GetToBuyItemFromKey(key string) Item {
	useItem := Item{}
	switch key {
	case "1": // SELECT 1st item
		useItem = shopItems[0]
	case "2": // SELECT 2nd item
		useItem = shopItems[1]
	case "3": // SELECT 3rd item
		useItem = shopItems[2]
	case "4": // SELECT 4th item
		useItem = shopItems[3]
	}
	return useItem
}
func Buy(user User, key string) string {
	useItem := GetToBuyItemFromKey(key)
	rcpName := ""
	switch useItem.Name {
	case "Wooden sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's Wooden sword lv1 buy recipe"
		}
	case "Copper sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's Copper sword lv1 buy recipe"
		}
	}
	return ExecuteRecipe(user, rcpName, []string{})
}

func GetToSellItemFromKey(user User, key string) Item {
	items := user.InventoryItems()
	useItem := Item{}
	switch key {
	case "1": // SELECT 1st item
		useItem = items[0]
	case "2": // SELECT 2nd item
		useItem = items[1]
	case "3": // SELECT 3rd item
		useItem = items[2]
	case "4": // SELECT 4th item
		useItem = items[3]
	}
	return useItem
}

func Sell(user User, key string) string {
	useItem := GetToSellItemFromKey(user, key)
	itemIDs := []string{useItem.ID}

	rcpName := ""
	switch useItem.Name {
	case "Wooden sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's Lv1 wooden sword sell recipe"
		} else {
			rcpName = "LOUD's Lv2 wooden sword sell recipe"
		}
	case "Copper sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's Lv1 copper sword sell recipe"
		} else {
			rcpName = "LOUD's Lv2 copper sword sell recipe"
		}
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

func GetToUpgradeItemFromKey(user User, key string) Item {
	items := user.UpgradableItems()
	useItem := Item{}
	switch key {
	case "1": // SELECT 1st item
		useItem = items[0]
	case "2": // SELECT 2nd item
		useItem = items[1]
	case "3": // SELECT 3rd item
		useItem = items[2]
	case "4": // SELECT 4th item
		useItem = items[3]
	}
	return useItem
}

func Upgrade(user User, key string) string {
	useItem := GetToUpgradeItemFromKey(user, key)
	itemIDs := []string{useItem.ID}
	rcpName := ""
	switch useItem.Name {
	case "Wooden sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's Wooden sword lv1 to lv2 upgrade recipe"
		}
	case "Copper sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's Copper sword lv1 to lv2 upgrade recipe"
		}
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}
