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
	"LOUD's hunt with lv2 wooden sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2050ec6df-1ad4-418a-83f1-e40253fc1199",
	"LOUD's Copper sword lv1 buy recipe":            "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey217240bde-3b24-46f6-83e4-e44445c68c7e",
	"LOUD's Lv1 copper sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey21acc52a0-413b-4903-98b2-4c96d8bb43e2",
	"LOUD's Lv2 copper sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2298b995f-bb1d-460b-b3c3-c00dd2505fd8",
	"LOUD's hunt with lv1 wooden sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey231718069-49ce-4067-9ded-c008df4318d5",
	"LOUD's hunt with lv1 copper sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey23613b882-0854-4913-bde4-73af72d45ba3",
	"LOUD's hunt without sword recipe":              "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey256550368-4108-4929-94e3-6f6b17502b46",
	"LOUD's get initial coin recipe":                "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey262561db4-5aec-44bf-a461-74246b52ad1b",
	"LOUD's Lv1 wooden sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2822fe7de-a514-4c42-960b-ab42b63864c6",
	"LOUD's Wooden sword lv1 to lv2 upgrade recipe": "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey28b553aa3-e0ec-4dd9-8adf-2f80dd71d88c",
	"LOUD's Lv2 wooden sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2be3cc622-518c-4dec-8c55-3275df1faf76",
	"LOUD's Wooden sword lv1 buy recipe":            "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2d8ba2de2-3bc7-4c0d-89fd-e409ddf97205",
	"LOUD's Copper sword lv1 to lv2 upgrade recipe": "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2e2a922bb-0b7f-4d59-bd84-4419e7a9c8ff",
	"LOUD's hunt with lv2 copper sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2e818407d-d670-44c5-b94f-17a45fbd2e93",
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
