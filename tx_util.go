package loud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	testing "github.com/Pylons-tech/pylons/cmd/fixtures_test/evtesting"
	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
	"github.com/Pylons-tech/pylons/x/pylons/handlers"
	"github.com/Pylons-tech/pylons/x/pylons/msgs"
	"github.com/Pylons-tech/pylons/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var RcpIDs map[string]string = map[string]string{
	"LOUD's Copper sword lv1 buy recipe":            "LOUD-copper-sword-lv1-buy-recipe-v0.0.0-1579053457",
	"LOUD's get initial coin recipe":                "LOUD-get-initial-coin-recipe-v0.0.1-1579652622",
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
var customNodeLocal string = "localhost:26657"
var restEndpointLocal string = "http://localhost:1317"

var useRestTx bool = false
var useLocalDm bool = false

func init() {
	args := os.Args

	for _, arg := range args[2:len(args)] {
		switch arg {
		case "-locald":
			useLocalDm = true
		case "-userest":
			useRestTx = true
		}
	}
	if useLocalDm {
		customNode = customNodeLocal
		restEndpoint = restEndpointLocal
	}

	pylonSDK.CLIOpts.CustomNode = customNode
	if useRestTx {
		pylonSDK.CLIOpts.RestEndpoint = restEndpoint
	}
	log.Println("initing pylonSDK to customNode", customNode, "useRestTx=", useRestTx)
}

func SyncFromNode(user User) {
	log.Println("SyncFromNode username=", user.GetUserName())
	log.Println("SyncFromNode username=", pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT()))
	accAddr := pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT())
	accInfo := pylonSDK.GetAccountInfoFromName(user.GetUserName(), GetTestingT())
	log.Println("accountInfo Result=", accInfo)

	user.SetGold(int(accInfo.Coins.AmountOf("loudcoin").Int64()))
	user.SetPylonAmount(int(accInfo.Coins.AmountOf("pylon").Int64()))
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

	nBuyOrders := []Order{}
	nSellOrders := []Order{}
	rawTrades, _ := pylonSDK.ListTradeViaCLI("")
	for _, tradeItem := range rawTrades {
		if tradeItem.Completed == false && len(tradeItem.CoinInputs) > 0 {
			inputCoin := tradeItem.CoinInputs[0].Coin
			if inputCoin == "loudcoin" { // loud sell trade
				pylonAmount := tradeItem.CoinOutputs.AmountOf("pylon").Int64()
				loudAmount := tradeItem.CoinInputs[0].Count
				nSellOrders = append(nSellOrders, Order{
					ID:        tradeItem.ID,
					Amount:    int(loudAmount),
					Total:     int(pylonAmount),
					Price:     float64(pylonAmount) / float64(loudAmount),
					IsMyOrder: tradeItem.Sender.String() == accAddr,
				})
			} else { // loud buy trade
				loudAmount := tradeItem.CoinOutputs.AmountOf("loudcoin").Int64()
				pylonAmount := tradeItem.CoinInputs[0].Count
				nBuyOrders = append(nBuyOrders, Order{
					ID:        tradeItem.ID,
					Amount:    int(loudAmount),
					Total:     int(pylonAmount),
					Price:     float64(pylonAmount) / float64(loudAmount),
					IsMyOrder: tradeItem.Sender.String() == accAddr,
				})
			}
		}
	}
	// Sort and show by low price buy orders
	sort.SliceStable(nBuyOrders, func(i, j int) bool {
		return nBuyOrders[i].Price < nBuyOrders[j].Price
	})
	// Sort and show by high price sell orders
	sort.SliceStable(nSellOrders, func(i, j int) bool {
		return nSellOrders[i].Price > nSellOrders[j].Price
	})
	buyOrders = nBuyOrders
	sellOrders = nSellOrders
	log.Println("SyncFromNode buyOrders=", nBuyOrders)
	log.Println("SyncFromNode sellOrders=", sellOrders)
}

func GetExtraPylons(user User) (string, error) {
	t := GetTestingT()
	username := user.GetUserName()
	addr := pylonSDK.GetAccountAddr(username, t)
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	log.Println("sdkAddr, err := sdk.AccAddressFromBech32(addr)", sdkAddr, err)
	extraPylonsMsg := msgs.NewMsgGetPylons(types.PremiumTier.Fee, sdkAddr)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, extraPylonsMsg, username, false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func GetInitialPylons(username string) (string, error) {
	addr := pylonSDK.GetAccountAddr(username, GetTestingT())
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
	if err != nil {
		return "", err
	}

	postBodyJSON := make(map[string]interface{})
	json.Unmarshal(signedTx, &postBodyJSON)
	postBodyJSON["tx"] = postBodyJSON["value"]
	postBodyJSON["value"] = nil
	postBodyJSON["mode"] = "sync"
	postBody, err := json.Marshal(postBodyJSON)

	log.Println("postBody", string(postBody))

	if err != nil {
		return "", err
	}
	resp, err := http.Post(restEndpoint+"/txs", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return "", err
	}

	var result map[string]string

	json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()
	log.Println("get_pylons_api_response", result)
	return result["txhash"], nil
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
	accBytes, err := pylonSDK.RunPylonsCli([]string{"query", "account", addr}, "")
	log.Println("query account for", addr, "result", string(accBytes), err)
	if err != nil {
		log.Println("err.Error()", err.Error())
		if strings.Contains(string(accBytes), "dial tcp [::1]:26657: connect: connection refused") { // Daemon is off
			log.Println("Daemon refused to connect, please check daemon is running!")
			os.Exit(3)
		} else { // account does not exist
			txhash, err := GetInitialPylons(username)
			if err != nil {
				log.Fatalln("txhash, err := GetInitialPylons", txhash, err)
			}
			log.Println("ran command for new account on remote chain and waiting for next block ...", addr)
			pylonSDK.WaitForNextBlock()
		}
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

func LogFullTxResultByHash(txhash string) {
	output, err := pylonSDK.RunPylonsCli([]string{"query", "tx", txhash}, "")

	log.Println("txhash=", txhash, "txoutput=", string(output), "queryerr=", err)
}

func ProcessTxResult(user User, txhash string) ([]byte, string) {
	t := GetTestingT()

	resp := handlers.ExecuteRecipeResp{}

	txHandleResBytes, err := pylonSDK.WaitAndGetTxData(txhash, 3, t)
	if err != nil {
		errString := fmt.Sprintf("error getting tx result bytes %+v", err)
		log.Println(errString)
		LogFullTxResultByHash(txhash)
		return []byte{}, errString
	}
	LogFullTxResultByHash(txhash)
	hmrErrMsg := pylonSDK.GetHumanReadableErrorFromTxHash(txhash, t)
	if len(hmrErrMsg) > 0 {
		errString := fmt.Sprintf("txhash=%s hmrErrMsg=%s", txhash, hmrErrMsg)
		log.Println(errString)
		return []byte{}, errString
	}
	SyncFromNode(user)

	err = pylonSDK.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
	if err != nil {
		errString := fmt.Sprintf("failed to parse transaction result; maybe this is get_pylons? txhash=%s", txhash)
		log.Println(errString)
		return []byte{}, ""
	} else {
		log.Println("ProcessTxResult::txResp", resp.Message, string(resp.Output))
		return resp.Output, ""
	}
}

func GetTestingT() *testing.T {
	newT := testing.NewT(nil)
	t := &newT
	return t
}

func ExecuteRecipe(user User, rcpName string, itemIDs []string) (string, error) {
	t := GetTestingT()
	if len(rcpName) == 0 {
		return "", errors.New("Recipe Name does not exist!")
	}
	rcpID, ok := RcpIDs[rcpName]
	if !ok {
		return "", errors.New("RecipeID does not exist for rcpName=" + rcpName)
	}
	addr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(addr)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, itemIDs)
	log.Println("started sending transaction", user.GetUserName(), execMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
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
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func Hunt(user User, key string) (string, error) {
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
	if itemKey >= 0 && itemKey < len(shopItems) {
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
			rcpName = "LOUD's Wooden sword lv1 buy recipe"
		}
	case "Copper sword":
		if useItem.Level == 1 {
			rcpName = "LOUD's Copper sword lv1 buy recipe"
		}
	default:
		return "", errors.New("You are trying to buy something which is not in shop")
	}
	if useItem.Price > user.GetGold() {
		return "", errors.New("You don't have enough gold to buy this item")
	}
	return ExecuteRecipe(user, rcpName, []string{})
}

func GetToSellItemFromKey(user User, key string) Item {
	items := user.InventoryItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func Sell(user User, key string) (string, error) {
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
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func Upgrade(user User, key string) (string, error) {
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
	if useItem.GetUpgradePrice() > user.GetGold() {
		return "", errors.New("You don't have enough gold to upgrade this item")
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

func CreateSellLoudOrder(user User, loudEnterValue string, pylonEnterValue string) (string, error) {
	t := GetTestingT()
	loudValue, err := strconv.Atoi(loudEnterValue)
	if err != nil {
		return "", err
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}

	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, err := sdk.AccAddressFromBech32(eugenAddr)

	inputCoinList := types.GenCoinInputList("loudcoin", int64(loudValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := "created by loud game"

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	log.Println("started sending transaction", user.GetUserName(), createTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, createTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func CreateBuyLoudOrder(user User, loudEnterValue string, pylonEnterValue string) (string, error) {
	t := GetTestingT()
	loudValue, err := strconv.Atoi(loudEnterValue)
	if err != nil {
		return "", err
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}

	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("loudcoin", int64(loudValue))}
	extraInfo := "created by loud game"

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	log.Println("started sending transaction", user.GetUserName(), createTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, createTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func FulfillTrade(user User, tradeID string) (string, error) {
	t := GetTestingT()
	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)
	ffTrdMsg := msgs.NewMsgFulfillTrade(tradeID, sdkAddr, []string{})

	log.Println("started sending transaction", user.GetUserName(), ffTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, ffTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}
