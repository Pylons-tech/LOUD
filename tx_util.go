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
	"LOUD's Lv1 wooden sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey205e63ef7-1cea-4430-8a88-139eae46da38",
	"LOUD's Lv2 copper sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey20b7499c3-8059-43af-a2ed-e7b6ccb599bc",
	"LOUD's Lv1 copper sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey219fe38a7-b82a-4111-b7ef-2de769dd8a82",
	"LOUD's Copper sword lv1 to lv2 upgrade recipe": "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey227e8d321-efab-4c3b-93f2-3994c80bfc9d",
	"LOUD's Lv2 wooden sword sell recipe":           "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey249bfd83a-e21b-4820-9460-234880d4de0b",
	"LOUD's hunt with lv1 copper sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey24bafc255-afa9-40ba-af5c-a34dd76a6b7d",
	"LOUD's Wooden sword lv1 buy recipe":            "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey25b41c08c-11f6-44ff-aaac-abb4cc874e05",
	"LOUD's hunt with lv2 copper sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey25d0f36ca-7dad-4e06-a120-4e50133cdb8e",
	"LOUD's Copper sword lv1 buy recipe":            "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey26633cc72-3b2f-4134-9d0d-0c577afce9b4",
	"LOUD's hunt with lv1 wooden sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey285b1179c-23c1-4289-b7bc-f713c9304bd9",
	"LOUD's hunt without sword recipe":              "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey28b422c95-9adf-4f1f-84e6-486706d0e8f4",
	"LOUD's Wooden sword lv1 to lv2 upgrade recipe": "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey28d30f11a-f4ad-4ea3-b3ee-8814d025ae06",
	"LOUD's hunt with lv2 wooden sword recipe":      "cosmos19vlpdf25cxh0w2s80z44r9ktrgzncf7zsaqey2f51e5c1b-1af9-44f5-b1f0-aeced8b6d144",
}

func SyncFromNode(user User) {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT

	rawItems, _ := pylonSDK.ListItemsViaCLI("eugen")
	items := []Item{}
	for _, rawItem := range rawItems {
		Level, _ := rawItem.FindLong("Level")
		Name, _ := rawItem.FindString("Name")
		items = append(items, Item{
			Level: Level,
			Name:  Name,
			ID:    rawItem.ID,
		})
	}
	user.SetItems(items)
	log.Println("SyncFromNode items=", items)
	accInfo := pylonSDK.GetAccountInfoFromName("eugen", t)
	user.SetGold(int(accInfo.Coins.AmountOf("loudcoin").Int64()))
	log.Println("SyncFromNode gold=", accInfo.Coins.AmountOf("loudcoin").Int64())
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

func Hunt(user User, key string) string {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT

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
	rcpName := "LOUD's hunt without sword recipe"
	switch useItem.Name {
	case "Wooden sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's hunt with lv1 wooden sword recipe"
		} else {
			rcpName = "LOUD's hunt with lv2 wooden sword recipe"
		}
	case "Copper sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's hunt with lv1 copper sword recipe"
		} else {
			rcpName = "LOUD's hunt with lv2 copper sword recipe"
		}
	}
	rcpID := RcpIDs[rcpName]
	eugenAddr := pylonSDK.GetAccountAddr("eugen", nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)
	// execMsg := msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, []string{})
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, "eugen", false)
	user.SetLastTransaction(txhash)

	return txhash
}

func Buy(user User, key string) string {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT

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
	rcpID := RcpIDs[rcpName]
	eugenAddr := pylonSDK.GetAccountAddr("eugen", nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)
	// execMsg := msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, []string{})
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, "eugen", false)
	user.SetLastTransaction(txhash)

	return txhash
}

func Sell(user User, key string) string {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT

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
	rcpID := RcpIDs[rcpName]
	eugenAddr := pylonSDK.GetAccountAddr("eugen", nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)
	// execMsg := msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, []string{})
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, "eugen", false)
	user.SetLastTransaction(txhash)

	return txhash
}

func Upgrade(user User, key string) string {
	orgT := originT.T{}
	newT := testing.NewT(&orgT)
	t := &newT

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
	rcpID := RcpIDs[rcpName]
	eugenAddr := pylonSDK.GetAccountAddr("eugen", nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)
	// execMsg := msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, []string{})
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, "eugen", false)
	user.SetLastTransaction(txhash)

	return txhash
}
