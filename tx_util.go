package loud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	originT "testing"

	fixtureSDK "github.com/MikeSofaer/pylons/cmd/fixtures_test"
	testing "github.com/MikeSofaer/pylons/cmd/fixtures_test/evtesting"
	pylonSDK "github.com/MikeSofaer/pylons/cmd/test"
	"github.com/MikeSofaer/pylons/x/pylons/handlers"
	"github.com/MikeSofaer/pylons/x/pylons/msgs"
	"github.com/MikeSofaer/pylons/x/pylons/types"
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

// Remote mode
var customNode string = "35.223.7.2:26657"
var restEndpoint string = "http://35.238.123.59:80"

// Local mode
// var customNode string = "localhost:26657"
// var restEndpoint string = "http://localhost:1317"

func init() {
	log.Println("initing pylonSDK to customNode", customNode)
	pylonSDK.CLIOpts.CustomNode = customNode
}

func SyncFromNode(user User) {
	log.Println("SyncFromNode username=", user.GetUserName())
	log.Println("SyncFromNode username=", pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT()))
	accInfo := pylonSDK.GetAccountInfoFromName(user.GetUserName(), GetTestingT())
	log.Println("accountInfo Result=", accInfo)

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

func GetInitialPylons(addr string) {

	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	log.Println("sdkAddr, err := sdk.AccAddressFromBech32(addr)", sdkAddr, err)

	// this code is making the account to useable by doing get-pylons
	txModel, err := pylonSDK.GenTxWithMsg([]sdk.Msg{msgs.NewMsgGetPylons(types.PremiumTier.Fee, sdkAddr)})
	output, err := pylonSDK.GetAminoCdc().MarshalJSON(txModel)

	tmpDir, err := ioutil.TempDir("", "pylons")

	rawTxFile := filepath.Join(tmpDir, "raw_tx_get_pylons_"+addr+".json")
	ioutil.WriteFile(rawTxFile, output, 0644)

	// pylonscli tx sign raw_tx_get_pylons_eugen.json --account-number 0 --sequence 0 --offline --from eugen
	txSignArgs := []string{"tx", "sign", rawTxFile,
		"--from", addr,
		"--offline",
		"--chain-id", "pylonschain",
		"--sequence", "0",
		"--account-number", "0",
	}
	signedTx, err := pylonSDK.RunPylonsCli(txSignArgs, "11111111\n")

	postBodyJSON := make(map[string]interface{})
	json.Unmarshal(signedTx, &postBodyJSON)
	postBodyJSON["tx"] = postBodyJSON["value"]
	postBodyJSON["value"] = nil
	postBodyJSON["mode"] = "sync"
	postBody, err := json.Marshal(postBodyJSON)

	log.Println("postBody", string(postBody))

	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(restEndpoint+"/txs", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()
	log.Println("get_pylons_api_response", result)
}

func InitPylonAccount(username string) {
	// "pylonscli keys add ${username}"
	addResult, err := pylonSDK.RunPylonsCli([]string{
		"keys", "add", username,
	}, "11111111\n11111111\n")

	log.Println("addResult, err := pylonSDK.RunPylonsCli", string(addResult), "---", err)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Println("pylonscli is not globally installed on your machine")
			os.Exit(1)
		} else {
			log.Println("using existing account for", username)
		}
	} else {
		usr, _ := user.Current()
		pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
		os.MkdirAll(pylonsDir, os.ModePerm)
		keyFile := filepath.Join(pylonsDir, username+".json")
		ioutil.WriteFile(keyFile, addResult, 0644)
		log.Println("created new account for", username, "and saved to ~/.pylons/"+username+".json")
	}
	addr := pylonSDK.GetAccountAddr(username, GetTestingT())
	// pylonSDK.CLIOpts.CustomNode = customNode
	accBytes, err := pylonSDK.RunPylonsCli([]string{"query", "account", addr}, "")
	log.Println("query account for", addr, "result", string(accBytes), err)
	if err != nil { // account does not exist
		GetInitialPylons(addr)
		log.Println("ran command for new account on remote chain and waiting for next block ...", addr)
		pylonSDK.WaitForNextBlock()
	} else {
		log.Println("using existing account on remote chain", addr)
	}

	// Remove nonce file
	log.Println("start removing nonce file")
	nonceRootDir := "./"
	nonceFile := filepath.Join(nonceRootDir, "nonce.json")
	err = os.Remove(nonceFile)
	log.Println("remove nonce file result", err)
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
	addr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(addr)
	// execMsg := msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, itemIDs)
	log.Println("started sending transaction", user.GetUserName(), execMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash
}

func GetIndexFromString(key string) int {
	switch key {
	case "1": // SELECT 1st item
		return 0
	case "2": // SELECT 2nd item
		return 1
	case "3": // SELECT 3rd item
		return 2
	case "4": // SELECT 4th item
		return 3
	case "5": // SELECT 5th item
		return 4
	case "6": // SELECT 6th item
		return 5
	case "7": // SELECT 7th item
		return 6
	case "8": // SELECT 8th item
		return 7
	case "9": // SELECT 9th item
		return 8
	}
	return -1
}

func GetWeaponItemFromKey(user User, key string) Item {
	items := user.InventoryItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 {
		useItem = items[itemKey]
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
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 {
		useItem = shopItems[itemKey]
	}
	return useItem
}
func Buy(user User, key string) (string, error) {
	useItem := GetToBuyItemFromKey(key)
	rcpName := ""
	switch useItem.Name {
	case "Wooden sword":
		if useItem.Level == 1 {
			if useItem.Price > user.GetGold() {
				return "", errors.New("You don't have enough funds to buy this item")
			}
			rcpName = "LOUD's Wooden sword lv1 buy recipe"
		}
	case "Copper sword":
		if useItem.Level == 1 {
			if useItem.Price > user.GetGold() {
				return "", errors.New("You don't have enough funds to buy this item")
			}
			rcpName = "LOUD's Copper sword lv1 buy recipe"
		}
	default:
		return "", errors.New("you are trying to buy something which is not in shop")
	}
	return ExecuteRecipe(user, rcpName, []string{}), nil
}

func GetToSellItemFromKey(user User, key string) Item {
	items := user.InventoryItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 {
		useItem = items[itemKey]
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
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 {
		useItem = items[itemKey]
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
