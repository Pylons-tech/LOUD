package loud

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Pylons-tech/LOUD/log"
	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"
	pylonSDK "github.com/Pylons-tech/pylons_sdk/cmd/test"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/handlers"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tyler-smith/go-bip39"
)

const (
	RCP_BUY_GOLD_WITH_PYLON = "LOUD's buy gold with pylons recipe"
	RCP_BUY_CHARACTER       = "LOUD's Get Character recipe"
	RCP_SELL_SWORD          = "LOUD's sword sell recipe"
	RCP_COPPER_SWORD_UPG    = "LOUD's Copper sword lv1 to lv2 upgrade recipe"
	RCP_WOODEN_SWORD_UPG    = "LOUD's Wooden sword lv1 to lv2 upgrade recipe"
	RCP_BUY_WOODEN_SWORD    = "LOUD's Wooden sword lv1 buy recipe"
	RCP_BUY_COPPER_SWORD    = "LOUD's Copper sword lv1 buy recipe"
	RCP_BUY_BRONZE_SWORD    = "LOUD's Bronze sword lv1 make recipe"
	RCP_BUY_IRON_SWORD      = "LOUD's Iron sword lv1 make recipe"
	RCP_BUY_SILVER_SWORD    = "LOUD's Silver sword lv1 make recipe"

	RCP_HUNT_RABBITS_NOSWORD = "LOUD's hunt rabbits without sword recipe"
	RCP_HUNT_RABBITS_YESWORD = "LOUD's hunt rabbits with a sword recipe"
	RCP_FIGHT_GIANT          = "LOUD's fight with giant with a sword recipe"
	RCP_FIGHT_GOBLIN         = "LOUD's fight with goblin with a sword recipe"
	RCP_FIGHT_TROLL          = "LOUD's fight with troll with a sword recipe"
	RCP_FIGHT_WOLF           = "LOUD's fight with wolf with a sword recipe"

	RCP_GET_TEST_ITEMS = "LOUD's Dev Get Test Items recipe"
)

var RcpIDs map[string]string = map[string]string{
	RCP_BUY_GOLD_WITH_PYLON: "LOUD-buy-gold-from-pylons-recipe-v0.0.1-1579652622",
	RCP_BUY_CHARACTER:       "LOUD-get-character-recipe-v0.0.0-1583801800",
	RCP_SELL_SWORD:          "LOUD-sell-a-sword-recipe-v0.0.0-1583631194",
	RCP_COPPER_SWORD_UPG:    "LOUD-upgrade-copper-sword-lv1-to-lv2-recipe-v0.0.0-1579053457",
	RCP_WOODEN_SWORD_UPG:    "LOUD-upgrade-wooden-sword-lv1-to-lv2-recipe-v0.0.0-1579053457",
	RCP_BUY_WOODEN_SWORD:    "LOUD-wooden-sword-lv1-buy-recipe-v0.0.0-1579053457",
	RCP_BUY_COPPER_SWORD:    "LOUD-copper-sword-lv1-buy-recipe-v0.0.0-1579053457",
	RCP_BUY_BRONZE_SWORD:    "LOUD-bronze-sword-lv1-make-recipe-v0.0.0-1579053457",
	RCP_BUY_IRON_SWORD:      "LOUD-iron-sword-lv1-make-recipe-v0.0.0-1579053457",
	RCP_BUY_SILVER_SWORD:    "LOUD-silver-sword-lv1-make-recipe-v0.0.0-1579053457",

	RCP_HUNT_RABBITS_NOSWORD: "LOUD-hunt-rabbits-with-no-weapon-recipe-v0.0.0-1579053457",
	RCP_HUNT_RABBITS_YESWORD: "LOUD-hunt-rabbits-with-a-sword-recipe-v0.0.0-1583631194",
	RCP_FIGHT_GIANT:          "LOUD-fight-giant-with-iron-sword-recipe-v0.0.0-1583631194",
	RCP_FIGHT_GOBLIN:         "LOUD-fight-goblin-with-a-sword-recipe-v0.0.0-1583631194",
	RCP_FIGHT_TROLL:          "LOUD-fight-troll-with-a-sword-recipe-v0.0.0-1583631194",
	RCP_FIGHT_WOLF:           "LOUD-fight-wolf-with-a-sword-recipe-v0.0.0-1583631194",

	RCP_GET_TEST_ITEMS: "LOUD-dev-get-test-items-recipe-v0.0.0-1583801800",
	RCP_RESTORE_HEALTH: "LOUD-health-restore-recipe-v0.0.1-1579652622",
}

// Remote mode
var customNode string = "35.223.7.2:26657"
var restEndpoint string = "http://35.238.123.59:80"

// Local mode
var customNodeLocal string = "localhost:26657"
var restEndpointLocal string = "http://localhost:1317"

var useRestTx bool = false
var useLocalDm bool = false
var AutomateInput bool = false
var AutomateRunCnt int = 0

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
				AutomateInput = true
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

func RunSHCmd(args []string) ([]byte, error) {
	cmd := exec.Command("/bin/sh", args...)
	res, err := cmd.CombinedOutput()
	log.Println("Running command /bin/sh", args)
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

	originSigner := signer
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
		nonce = nonceMap[signer] - 1
	} else {
		return false, errors.New("nonce file does not exist :(")
	}

	output, err := pylonSDK.GetAminoCdc().MarshalJSON(msgValue)
	t.MustNil(err)

	rawTxFile := filepath.Join(tmpDir, "raw_tx_"+strconv.FormatUint(nonce, 10)+".json")
	ioutil.WriteFile(rawTxFile, output, 0644)
	if err != nil {
		return false, err
	}

	t.Log("TX sign with nonce=", nonce)
	// sh txutil.sh <op> <privkey> <account number> <sequence> <msg>
	txSignArgs := []string{
		"./artifacts_txutil.sh",
		"SIGNED_TX",
		privKey,
		strconv.FormatUint(accInfo.GetAccountNumber(), 10),
		strconv.FormatUint(nonce, 10),
		rawTxFile,
	}
	aftiOutput, err := RunSHCmd(txSignArgs)
	if err != nil {
		return false, err
	}

	log.Println("RunSHCmd output, err=", string(aftiOutput), err)
	cliTxOutput, err := pylonSDK.RunPylonsCli([]string{"query", "tx", txhash}, "")
	if err != nil {
		log.Println("txhash=", txhash, "txoutput=", string(cliTxOutput), "queryerr=", err)
	}

	// use regexp to find signature from cli command response
	re := regexp.MustCompile(`"signature":.*"(.*)"`)
	cliTxSign := re.FindSubmatch([]byte(cliTxOutput))
	aftiTxSign := re.FindSubmatch([]byte(aftiOutput))

	if len(cliTxSign) < 2 {
		log.Println("couldn't get pyloncli signature from", string(cliTxOutput))
		return false, errors.New("couldn't get pyloncli signature")
	} else if len(aftiTxSign) < 2 {
		log.Println("couldn't get afticli signature from", string(aftiOutput))
		return false, errors.New("couldn't get afticli signature")
	} else {
		pylonCliSignature := string(cliTxSign[1])
		aftiSignatue := string(aftiTxSign[1])
		log.Println("comparing afticli and pyloncli ;)", pylonCliSignature, "\nand\n", aftiSignatue)
	}
	log.Println("where")
	log.Println("msg=", string(output))
	log.Println("username=", originSigner)
	log.Println("Bech32Addr=", signer)
	log.Println("privKey=", privKey)
	log.Println("account-number=", strconv.FormatUint(accInfo.GetAccountNumber(), 10))
	log.Println("sequence", strconv.FormatUint(nonce, 10))

	if string(cliTxSign[1]) != string(aftiTxSign[1]) {
		return false, errors.New("comparison different afticli and pyloncli ")
	}

	pylonSDK.CleanFile(rawTxFile, t)

	return true, nil
}

func GetInitialPylons(username string) (string, error) {
	addr := pylonSDK.GetAccountAddr(username, GetTestingT())
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	log.Println("GetInitialPylons => sdkAddr, err", sdkAddr, err)

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

func ComputePrivKeyFromMnemonic(mnemonic string) (string, string) {
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		os.Exit(1)
	}

	// This priv get code came from dbKeybase.CreateMnemonic function of cosmos-sdk
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	hdPath := hd.NewFundraiserParams(0, 0).String()
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath)
	if err != nil {
		os.Exit(1)
	}
	priv := secp256k1.PrivKeySecp256k1(derivedPriv)

	privKeyHex := hex.EncodeToString(priv[:])
	cosmosAddr := sdk.AccAddress(priv.PubKey().Address().Bytes()).String()
	return privKeyHex, cosmosAddr
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
			SomethingWentWrongMsg = "pylonscli is not globally installed on your machine"
		} else {
			log.Println("using existing account for", username)
			usr, _ := user.Current()
			pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
			os.MkdirAll(pylonsDir, os.ModePerm)
			keyFile := filepath.Join(pylonsDir, username+".json")
			addResult, err = ioutil.ReadFile(keyFile)
			if err != nil && AutomateInput {
				log.Fatal("Couldn't get private key from ", username, ".json")
			}
			addedKeyResInterface := make(map[string]string)
			err = json.Unmarshal(addResult, &addedKeyResInterface)
			if err != nil && AutomateInput {
				log.Fatal("Couldn't parse file for", username, ".json", err.Error())
			}
			privKey = addedKeyResInterface["privkey"]
			log.Println("privKey=", privKey)
		}
	} else {
		addedKeyResInterface := make(map[string]string)
		json.Unmarshal(addResult, &addedKeyResInterface)

		// mnemonic key from the pylonscli add result
		mnemonic := addedKeyResInterface["mnemonic"]
		log.Println("using mnemonic: ", mnemonic)

		privKey, _ = ComputePrivKeyFromMnemonic(mnemonic) // get privKey and cosmosAddr

		addResult, err = json.Marshal(addedKeyResInterface)

		usr, _ := user.Current()
		pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
		os.MkdirAll(pylonsDir, os.ModePerm)
		keyFile := filepath.Join(pylonsDir, username+".json")
		ioutil.WriteFile(keyFile, addResult, 0644)
		log.Println("privKey=", privKey)
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
	hmrErrMsg, _ := pylonSDK.GetHumanReadableErrorFromTxHash(txhash, t)
	if len(hmrErrMsg) > 0 {
		errString := fmt.Sprintf("txhash=%s hmrErrMsg=%s", txhash, hmrErrMsg)
		log.Println(errString)
		return []byte{}, errString
	}
	SyncFromNode(user)

	err = pylonSDK.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
	if err != nil {
		errString := fmt.Sprintf("failed to parse transaction result; maybe this is get_pylons then ignore. txhash=%s", txhash)
		log.Println(errString)
		return []byte{}, errString
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
	user.SetLastTransaction(txhash, rcpName)
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
	items := user.InventorySwords()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func GetSwordItemFromKey(user User, key string) Item {
	items := user.InventorySwords()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func GetIronSwordItemFromKey(user User, key string) Item {
	items := user.InventoryIronSwords()
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
	if itemKey >= 0 && itemKey < len(ShopItems) {
		useItem = ShopItems[itemKey]
	}
	return useItem
}

func GetToBuyCharacterFromKey(key string) Character {
	useCharacter := Character{}
	cKey := GetIndexFromString(key)
	if cKey >= 0 && cKey < len(ShopCharacters) {
		useCharacter = ShopCharacters[cKey]
	}
	return useCharacter
}

func GetToSellItemFromKey(user User, key string) Item {
	items := user.InventorySellableItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func GetToUpgradeItemFromKey(user User, key string) Item {
	items := user.InventoryUpgradableItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

func GetItemInputsFromItemSpec(itspec ItemSpec) types.ItemInputList {
	var itemInputs types.ItemInputList

	ii := types.ItemInput{
		Doubles: nil,
		Longs: types.LongInputParamList{
			types.LongInputParam{Key: "level", MinValue: itspec.Level[0], MaxValue: itspec.Level[1]},
		},
		Strings: types.StringInputParamList{
			types.StringInputParam{Key: "Name", Value: itspec.Name},
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

func GetItemInputsFromCharacterSpec(chspec CharacterSpec) types.ItemInputList {
	var itemInputs types.ItemInputList

	ii := types.ItemInput{
		Doubles: types.DoubleInputParamList{
			types.DoubleInputParam{Key: "XP", MinValue: types.ToFloatString(chspec.XP[0]), MaxValue: types.ToFloatString(chspec.XP[1])},
		},
		Longs: types.LongInputParamList{
			types.LongInputParam{Key: "level", MinValue: chspec.Level[0], MaxValue: chspec.Level[1]},
		},
		Strings: types.StringInputParamList{
			types.StringInputParam{Key: "Name", Value: chspec.Name},
		},
	}
	itemInputs = append(itemInputs, ii)
	return itemInputs
}

func GetItemOutputFromActiveCharacter(activeCharacter Character) (types.ItemList, error) {
	var itemOutputs types.ItemList
	io, err := pylonSDK.GetItemByGUID(activeCharacter.ID)
	itemOutputs = append(itemOutputs, io)
	return itemOutputs, err
}

func GetSDKAddrFromUserName(username string) sdk.AccAddress {
	addr := pylonSDK.GetAccountAddr(username, nil)
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		log.Fatal("sdkAddr, err := sdk.AccAddressFromBech32(addr)", sdkAddr, err)
	}
	return sdkAddr
}

func SendTxMsg(user User, txMsg sdk.Msg) (string, error) {
	t := GetTestingT()
	log.Println("started sending transaction", user.GetUserName(), txMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, txMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash, txMsg.Type())
	log.Println("ended sending transaction")
	return txhash, nil
}
