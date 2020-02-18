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
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	testing "github.com/Pylons-tech/pylons/cmd/fixtures_test/evtesting"
	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
	"github.com/Pylons-tech/pylons/x/pylons/handlers"
	"github.com/Pylons-tech/pylons/x/pylons/msgs"
	"github.com/Pylons-tech/pylons/x/pylons/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tyler-smith/go-bip39"
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
var automateInput bool = false
var automateRunCnt int = 0

func init() {
	args := os.Args

	if len(args) > 1 {
		for _, arg := range args[2:len(args)] {
			switch arg {
			case "-locald":
				useLocalDm = true
			case "-userest":
				useRestTx = true
			case "-automate":
				automateInput = true
			}
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func RunAftiCli(args []string) ([]byte, error) { // run pylonscli with specific params : helper function
	// cmd := exec.Command(path.Join(os.Getenv("GOPATH"), "/bin/artifacts/txutil.sh"), args...)
	cmd := exec.Command("/bin/sh", args...)
	res, err := cmd.CombinedOutput()
	log.Println("Running command ./artifacts_txutil.sh", args)
	return res, err
}

func CheckSignatureMatchWithAftiCli(t *testing.T, txhash string, privKey string, msgValue sdk.Msg, signer string, isBech32Addr bool) (bool, error) {
	pylonSDK.WaitAndGetTxData(txhash, 3, t)
	tmpDir, err := ioutil.TempDir("", "pylons")
	if err != nil {
		panic(err.Error())
	}
	nonceRootDir := "./"
	nonceFile := filepath.Join(nonceRootDir, "nonce.json")
	if !isBech32Addr {
		signer = pylonSDK.GetAccountAddr(signer, t)
	}

	accInfo := pylonSDK.GetAccountInfoFromAddr(signer, t)
	nonce := accInfo.Sequence

	nonceMap := make(map[string]uint64)

	if fileExists(nonceFile) {
		nonceBytes := pylonSDK.ReadFile(nonceFile, t)
		err := json.Unmarshal(nonceBytes, &nonceMap)
		if err != nil {
			return false, err
		}
		nonce = nonceMap[signer]
	} else {
		nonce = accInfo.GetSequence()
	}
	nonceMap[signer] = nonce
	nonceOutput, err := json.Marshal(nonceMap)
	t.MustNil(err)
	ioutil.WriteFile(nonceFile, nonceOutput, 0644)

	output, err := pylonSDK.GetAminoCdc().MarshalJSON(msgValue)
	t.MustNil(err)

	rawTxFile := filepath.Join(tmpDir, "raw_tx_"+strconv.FormatUint(nonce, 10)+".json")
	ioutil.WriteFile(rawTxFile, output, 0644)
	if err != nil {
		return false, err
	}

	t.Log("TX sign with nonce=", nonce)
	// pylonscli tx sign sample_transaction.json --account-number 2 --sequence 10 --offline --from eugen
	txSignArgs := []string{
		"./artifacts_txutil.sh",
		privKey,
		strconv.FormatUint(accInfo.GetAccountNumber(), 10),
		strconv.FormatUint(nonce, 10),
		rawTxFile,
		"SIGNED_TX",
		// sh txutil.sh <privkey> <account number> <sequence> <msg> <op>
	}
	aftiOutput, err := RunAftiCli(txSignArgs)
	if err != nil {
		return false, err
	}

	// TODO check afti's content and the signature returned by txhash here
	log.Println("RunAftiCli output, err=", string(output), err)
	cliTxOutput, err := pylonSDK.RunPylonsCli([]string{"query", "tx", txhash}, "")
	if err != nil {
		log.Println("txhash=", txhash, "txoutput=", string(cliTxOutput), "queryerr=", err)
	}

	// "signature": "ouyh4zAwNs22FB7I9x3rRhb4RDPT2/UmIPOUc89/Nb9uYslAxX09CTbEf+7K8o3fyDW4QERf7zoPzno1gg6RDg=="
	// var txQueryResp sdk.TxResponse
	// err = pylonSDK.GetAminoCdc().UnmarshalJSON(cliTxOutput, &txQueryResp)
	re := regexp.MustCompile(`"signature":.*"(.*)"`)
	cliTxSign := re.FindSubmatch([]byte(cliTxOutput))
	aftiTxSign := re.FindSubmatch([]byte(aftiOutput))

	log.Println("comparing ", string(cliTxSign[1]), "\nand\n", string(aftiTxSign[1]))
	log.Println("where username=",
		signer,
		"privKey=",
		privKey,
		"account-number=",
		strconv.FormatUint(accInfo.GetAccountNumber(), 10),
		"sequence", strconv.FormatUint(nonce, 10),
		string(output),
	)
	if string(cliTxSign[1]) != string(aftiTxSign[1]) {
		os.Exit(1)
	}

	pylonSDK.CleanFile(rawTxFile, t)

	return true, nil
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

func InitPylonAccount(username string) string {
	var privKey string
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
			usr, _ := user.Current()
			pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
			os.MkdirAll(pylonsDir, os.ModePerm)
			keyFile := filepath.Join(pylonsDir, username+".json")
			addResult, err = ioutil.ReadFile(keyFile)
			if err != nil {
				log.Fatal("Couldn't get private key from ", username, ".json")
			}
			addedKeyResInterface := make(map[string]string)
			json.Unmarshal(addResult, &addedKeyResInterface)
			privKey = addedKeyResInterface["privkey"]
			log.Println("using existing account for", username, "privKey=", privKey)

			// os.Exit(1)
		}
	} else {
		addedKeyResInterface := make(map[string]string)
		json.Unmarshal(addResult, &addedKeyResInterface)

		// Generate a mnemonic for memorization or user-friendly seeds
		entropy, _ := bip39.EntropyFromMnemonic(addedKeyResInterface["mnemonic"])
		mnemonic, _ := bip39.NewMnemonic(entropy)

		// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
		seed := bip39.NewSeed(mnemonic, "11111111")
		// seed, err := bip39.NewSeedWithErrorChecking(tests.TestMnemonic, "")

		// masterKey, ch := hd.ComputeMastersFromSeed(seed)
		masterKey, _ := hd.ComputeMastersFromSeed(seed)
		privKey = fmt.Sprintf("%x", masterKey)
		addedKeyResInterface["privkey"] = privKey

		addResult, err = json.Marshal(addedKeyResInterface)

		usr, _ := user.Current()
		pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
		os.MkdirAll(pylonsDir, os.ModePerm)
		keyFile := filepath.Join(pylonsDir, username+".json")
		ioutil.WriteFile(keyFile, addResult, 0644)
		log.Println("privKey=", privKey)
		log.Println("created new account for", username, "and saved to ~/.pylons/"+username+".json")

		// os.Exit(1)
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
	return privKey
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

func GetToBuyItemFromKey(key string) Item {
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(shopItems) {
		useItem = shopItems[itemKey]
	}
	return useItem
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

func GetToUpgradeItemFromKey(user User, key string) Item {
	items := user.UpgradableItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func GetItemInputsFromActiveItem(activeItem Item) types.ItemInputList {
	var itemInputs types.ItemInputList

	ii := types.ItemInput{
		Doubles: nil,
		Longs: types.LongInputParamList{
			types.LongInputParam{"level", activeItem.Level, activeItem.Level},
		},
		Strings: types.StringInputParamList{
			types.StringInputParam{"Name", activeItem.Name},
		},
	}
	itemInputs = append(itemInputs, ii)
	return itemInputs
}

func GetItemOutputFromActiveItem(activeItem Item) (types.ItemList, error) {
	var itemOutputs types.ItemList
	io, err := pylonSDK.GetItemByGUID(activeItem.ID)
	itemOutputs = append(itemOutputs, io)
	return itemOutputs, err
}
