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

	cf "github.com/Pylons-tech/LOUD/config"
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

// GameCookbookID is CookbookID of current game
var GameCookbookID = "LOUD-v0.1.0-1589853709"

const (
	// RcpBuyGoldWithPylon is a recipe to buy gold with pylons
	RcpBuyGoldWithPylon = "LOUD's buy gold with pylons recipe"
	// RcpBuyCharacter is a recipe to buy character
	RcpBuyCharacter = "LOUD's Get Character recipe"
	// RcpSellSword is a recipe to sell item
	RcpSellSword = "LOUD's item sell recipe"
	// RcpCopperSwordUpgrade is a recipe to upgrade copper sword
	RcpCopperSwordUpgrade = "LOUD's Copper sword lv1 to lv2 upgrade recipe"
	// RcpWoodenSwordUpgrade is a recipe to upgrade wooden sword
	RcpWoodenSwordUpgrade = "LOUD's Wooden sword lv1 to lv2 upgrade recipe"
	// RcpBuyWoodenSword is a recipe to buy wooden sword
	RcpBuyWoodenSword = "LOUD's Wooden sword lv1 buy recipe"
	// RcpBuyCopperSword is a recipe to buy copper sword
	RcpBuyCopperSword = "LOUD's Copper sword lv1 buy recipe"
	// RcpBuyBronzeSword is a recipe to make bronze sword
	RcpBuyBronzeSword = "LOUD's Bronze sword lv1 make recipe"
	// RcpBuySilverSword is a recipe to make silver sword
	RcpBuySilverSword = "LOUD's Silver sword lv1 make recipe"
	// RcpBuyIronSword is a recipe to make iron sword
	RcpBuyIronSword = "LOUD's Iron sword lv1 make recipe"
	// RcpBuyAngelSword is a recipe to make angel sword
	RcpBuyAngelSword = "LOUD's Angel sword lv1 make recipe"

	// RcpHuntRabbits is a recipe to hunt rabbits
	RcpHuntRabbits = "LOUD's hunt rabbits without sword recipe"
	// RcpFightGoblin is a recipe to fight goblin
	RcpFightGoblin = "LOUD's fight with goblin with a sword recipe"
	// RcpFightWolf is a recipe to fight wolf
	RcpFightWolf = "LOUD's fight with wolf with a sword recipe"
	// RcpFightTroll is a recipe to fight troll
	RcpFightTroll = "LOUD's fight with troll with a sword recipe"
	// RcpFightGiant is a recipe to fight giant
	RcpFightGiant = "LOUD's fight with giant with a sword recipe" // 🗿
	// RcpFightDragonFire is a recipe to fight fire dragon
	RcpFightDragonFire = "LOUD's fight with fire dragon with an iron sword recipe"
	// RcpFightDragonIce is a recipe to fight ice dragon
	RcpFightDragonIce = "LOUD's fight with ice dragon with an iron sword recipe"
	// RcpFightDragonAcid is a recipe to fight acid dragon
	RcpFightDragonAcid = "LOUD's fight with acid dragon with an iron sword recipe"
	// RcpFightDragonUndead is a recipe to fight undead dragon
	RcpFightDragonUndead = "LOUD's fight with undead dragon with an angel sword recipe"

	// RcpGetTestItems is a recipe to get test items for development purposes
	RcpGetTestItems = "LOUD's Dev Get Test Items recipe"
)

// RcpIDs convert recipe name to id
var RcpIDs map[string]string = map[string]string{
	RcpBuyGoldWithPylon:   "LOUD-buy-gold-from-pylons-recipe-v0.1.0-1589853709",
	RcpBuyCharacter:       "LOUD-get-character-recipe-v0.1.0-1589853709",
	RcpSellSword:          "LOUD-sell-an-item-recipe-v0.1.0-1589853709",
	RcpCopperSwordUpgrade: "LOUD-upgrade-copper-sword-lv1-to-lv2-recipe-v0.1.0-1589853709",
	RcpWoodenSwordUpgrade: "LOUD-upgrade-wooden-sword-lv1-to-lv2-recipe-v0.1.0-1589853709",
	RcpBuyWoodenSword:     "LOUD-wooden-sword-lv1-buy-recipe-v0.1.0-1589853709",
	RcpBuyCopperSword:     "LOUD-copper-sword-lv1-buy-recipe-v0.1.0-1589853709",
	RcpBuyBronzeSword:     "LOUD-bronze-sword-lv1-make-recipe-v0.1.0-1589853709",
	RcpBuySilverSword:     "LOUD-silver-sword-lv1-make-recipe-v0.1.0-1589853709",
	RcpBuyIronSword:       "LOUD-iron-sword-lv1-make-recipe-v0.1.0-1589853709",
	RcpBuyAngelSword:      "LOUD-angel-sword-lv1-make-recipe-v0.1.0-1589853709",

	RcpHuntRabbits:       "LOUD-hunt-rabbits-with-no-weapon-recipe-v0.1.0-1589853709",
	RcpFightGiant:        "LOUD-fight-giant-with-iron-sword-recipe-v0.1.0-1589853709",
	RcpFightGoblin:       "LOUD-fight-goblin-with-a-sword-recipe-v0.1.0-1589853709",
	RcpFightTroll:        "LOUD-fight-troll-with-a-sword-recipe-v0.1.0-1589853709",
	RcpFightWolf:         "LOUD-fight-wolf-with-a-sword-recipe-v0.1.0-1589853709",
	RcpFightDragonFire:   "LOUD-fight-fire-dragon-with-iron-sword-recipe-v0.1.0-1589853709",
	RcpFightDragonIce:    "LOUD-fight-ice-dragon-with-iron-sword-recipe-v0.1.0-1589853709",
	RcpFightDragonAcid:   "LOUD-fight-acid-dragon-with-iron-sword-recipe-v0.1.0-1589853709",
	RcpFightDragonUndead: "LOUD-fight-undead-dragon-with-angel-sword-recipe-v0.1.0-1589853709",

	RcpGetTestItems: "LOUD-dev-get-test-items-recipe-v0.1.0-1589853709",
}

// Remote mode
var customNode string
var restEndpoint string
var maxWaitBlock int64
var useRestTx bool = false

// AutomateInput refers to automatic keyboard input event generation
var AutomateInput bool = false

// AutomateRunCnt refers to automatic keyboard input event generation count
var AutomateRunCnt int = 0

func init() {
	cfg, cferr := cf.ReadConfig()
	useRestTx = cfg.Terminal.UseRestTx
	AutomateInput = cfg.Terminal.AutomateInput

	if cferr == nil {
		restEndpoint = cfg.SDK.RestEndpoint
		customNode = cfg.SDK.CliEndpoint
		maxWaitBlock = cfg.SDK.MaxWaitBlock
	} else {
		log.Fatal("Couldn't load configuration file, log=\"", cferr, "\"")
	}

	pylonSDK.CLIOpts.CustomNode = customNode
	pylonSDK.CLIOpts.MaxWaitBlock = maxWaitBlock
	if useRestTx {
		pylonSDK.CLIOpts.RestEndpoint = restEndpoint
	}
	log.Infoln("initialized pylonSDK to connect to customNode", customNode, "and to custom rest endpoint", useRestTx)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// RunSHCmd is a function to run sh command and return response
func RunSHCmd(args []string) ([]byte, error) {
	cmd := exec.Command("/bin/sh", args...)
	res, err := cmd.CombinedOutput()
	log.Debugln("Running command /bin/sh", args)
	return res, err
}

// CheckSignatureMatchWithAftiCli is a function to check java code(afti)'s signature and pylonscli's signature
func CheckSignatureMatchWithAftiCli(t *testing.T, txhash string, privKey string, msgValue sdk.Msg, signer string, isBech32Addr bool) (bool, error) {

	_, err := pylonSDK.WaitAndGetTxData(txhash, pylonSDK.GetMaxWaitBlock(), t)
	if err != nil {
		t.Fatal("Error waiting for transaction by hash")
	}
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
	var nonce uint64

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
	err = ioutil.WriteFile(rawTxFile, output, 0644)
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

	log.Debugln("RunSHCmd output, err=", string(aftiOutput), err)
	cliTxOutput, _, err := pylonSDK.RunPylonsCli([]string{"query", "tx", txhash}, "")
	if err != nil {
		log.Debugln("txhash=", txhash, "txoutput=", string(cliTxOutput), "queryerr=", err)
	}

	// use regexp to find signature from cli command response
	re := regexp.MustCompile(`"signature":.*"(.*)"`)
	cliTxSign := re.FindSubmatch([]byte(cliTxOutput))
	aftiTxSign := re.FindSubmatch([]byte(aftiOutput))

	if len(cliTxSign) < 2 {
		log.Warnln("couldn't get pyloncli signature from", string(cliTxOutput))
		return false, errors.New("couldn't get pyloncli signature")
	} else if len(aftiTxSign) < 2 {
		log.Warnln("couldn't get afticli signature from", string(aftiOutput))
		return false, errors.New("couldn't get afticli signature")
	} else {
		pylonCliSignature := string(cliTxSign[1])
		aftiSignature := string(aftiTxSign[1])
		log.Infof("comparing afticli and pyloncli \n%s\nand\n %s", pylonCliSignature, aftiSignature)
	}
	log.Infoln("where")
	log.Infoln("msg=", string(output))
	log.Infoln("username=", originSigner)
	log.Infoln("Bech32Addr=", signer)
	log.Infoln("privKey=", privKey)
	log.Infoln("account-number=", strconv.FormatUint(accInfo.GetAccountNumber(), 10))
	log.Infoln("sequence", strconv.FormatUint(nonce, 10))

	if string(cliTxSign[1]) != string(aftiTxSign[1]) {
		return false, errors.New("comparison different afticli and pyloncli ")
	}

	pylonSDK.CleanFile(rawTxFile, t)

	return true, nil
}

// GetInitialPylons is a function to get initial pylons from faucet
func GetInitialPylons(username string) (string, error) {
	addr := pylonSDK.GetAccountAddr(username, GetTestingT())
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	log.Debugln("GetInitialPylons => sdkAddr, err", sdkAddr, err)

	// this code is making the account to useable by doing get-pylons
	txModel, err := pylonSDK.GenTxWithMsg([]sdk.Msg{msgs.NewMsgGetPylons(types.PremiumTier.Fee, sdkAddr)})
	if err != nil {
		return "", err
	}
	output, err := pylonSDK.GetAminoCdc().MarshalJSON(txModel)
	if err != nil {
		return "", err
	}
	tmpDir, err := ioutil.TempDir("", "pylons")
	if err != nil {
		return "", err
	}

	rawTxFile := filepath.Join(tmpDir, "raw_tx_get_pylons_"+addr+".json")
	err = ioutil.WriteFile(rawTxFile, output, 0644)
	if err != nil {
		return "", err
	}

	// pylonscli tx sign raw_tx_get_pylons_eugen.json --account-number 0 --sequence 0 --offline --from eugen
	txSignArgs := []string{"tx", "sign", rawTxFile,
		"--from", addr,
		"--offline",
		"--chain-id", "pylonschain",
		"--sequence", "0",
		"--account-number", "0",
	}
	signedTx, _, err := pylonSDK.RunPylonsCli(txSignArgs, "11111111\n")
	if err != nil {
		return "", err
	}

	postBodyJSON := make(map[string]interface{})
	err = json.Unmarshal(signedTx, &postBodyJSON)
	if err != nil {
		log.Fatal("Error unmarshalling signedTx into postBody JSON")
	}
	postBodyJSON["tx"] = postBodyJSON["value"]
	postBodyJSON["value"] = nil
	postBodyJSON["mode"] = "sync"
	postBody, err := json.Marshal(postBodyJSON)

	log.Debugln("postBody", string(postBody))

	if err != nil {
		return "", err
	}
	resp, err := http.Post(restEndpoint+"/txs", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return "", err
	}

	var result map[string]string

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal("Error decoding api response to result")
	}
	defer resp.Body.Close()
	log.Debugln("get_pylons_api_response", result)
	return result["txhash"], nil
}

// ComputePrivKeyFromMnemonic calculates private key from mnemonic
func ComputePrivKeyFromMnemonic(mnemonic string) (string, string) {
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		os.Exit(1)
	}

	// This priv get code came from dbKeybase.CreateMnemonic function of cosmos-sdk
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	// hdPath := hd.NewFundraiserParams(0, 0).String()
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, "44'/118'/0'/0/0")
	if err != nil {
		os.Exit(1)
	}
	priv := secp256k1.PrivKeySecp256k1(derivedPriv)

	privKeyHex := hex.EncodeToString(priv[:])
	cosmosAddr := sdk.AccAddress(priv.PubKey().Address().Bytes()).String()
	return privKeyHex, cosmosAddr
}

// InitPylonAccount initialize an account on local and get initial balance from faucet
func InitPylonAccount(username string) string {
	log.Debugln("InitPylonAccount has started")
	var privKey string
	// "pylonscli keys add ${username}"
	addResult, _, err := pylonSDK.RunPylonsCli([]string{
		"keys", "add", username,
	}, "11111111\n11111111\n")

	log.Debugln("addResult, err := pylonSDK.RunPylonsCli", string(addResult), "---", err)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Warnln("pylonscli is not globally installed on your machine")
			SomethingWentWrongMsg = "pylonscli is not globally installed on your machine"
		} else {
			log.Infoln("using existing account for", username)
			usr, _ := user.Current()
			pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
			err = os.MkdirAll(pylonsDir, os.ModePerm)
			if err != nil {
				log.Fatal("Couldn't create dirs for ~/.pylons")
			}
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
			log.Debugln("privKey=", privKey)
		}
	} else {
		addedKeyResInterface := make(map[string]string)
		err = json.Unmarshal(addResult, &addedKeyResInterface)
		if err != nil {
			log.Fatal("Error unmarshalling into key result interface")
		}

		// mnemonic key from the pylonscli add result
		mnemonic := addedKeyResInterface["mnemonic"]
		log.Debugln("using mnemonic: ", mnemonic)

		privKey, _ = ComputePrivKeyFromMnemonic(mnemonic) // get privKey and cosmosAddr

		addResult, err = json.Marshal(addedKeyResInterface)
		if err != nil {
			log.Fatal("Couldn't marshal added Keys result")
		}

		usr, _ := user.Current()
		pylonsDir := filepath.Join(usr.HomeDir, ".pylons")
		err = os.MkdirAll(pylonsDir, os.ModePerm)
		if err != nil {
			log.Fatal("Couldn't create dirs for ~/.pylons")
		}
		keyFile := filepath.Join(pylonsDir, username+".json")
		if ioutil.WriteFile(keyFile, addResult, 0644) != nil {
			log.Fatal("error writing file to ~/.pylons directory")
		}
		log.Debugln("privKey=", privKey)
		log.Infoln("created new account for", username, "and saved to ~/.pylons/"+username+".json")
	}
	addr := pylonSDK.GetAccountAddr(username, GetTestingT())
	accBytes, _, err := pylonSDK.RunPylonsCli([]string{"query", "account", addr}, "")
	log.Debugln("query account for", addr, "result", string(accBytes), err)
	if err != nil {
		log.Debugln("err.Error()", err.Error())
		if strings.Contains(string(accBytes), "dial tcp [::1]:26657: connect: connection refused") { // Daemon is off
			log.Fatalln("Daemon refused to connect, please check daemon is running!")
		} else { // account does not exist
			txhash, err := GetInitialPylons(username)
			if err != nil {
				log.Fatalln("txhash, err := GetInitialPylons", txhash, err)
			}
			log.Debugln("ran command for new account on remote chain and waiting for next block ...", addr)
			if pylonSDK.WaitForNextBlock() != nil {
				return "error waiting for block"
			}
		}
	} else {
		log.Infoln("using existing account on remote chain", addr)
	}

	// Remove nonce file
	log.Debugln("start removing nonce file")
	nonceRootDir := "./"
	nonceFile := filepath.Join(nonceRootDir, "nonce.json")
	err = os.Remove(nonceFile)
	log.Debugln("remove nonce file result", err)
	log.Debugln("InitPylonAccount has ended with privKey", privKey)
	return privKey
}

// LogFullTxResultByHash implements log for a transaction hash
func LogFullTxResultByHash(txhash string) {
	output, _, err := pylonSDK.RunPylonsCli([]string{"query", "tx", txhash}, "")

	log.Debugln("txhash=", txhash, "txoutput=", string(output), "queryerr=", err)
}

// ProcessTxResult is a function to handle result of a transaction made
func ProcessTxResult(user User, txhash string) ([]byte, string) {
	t := GetTestingT()

	resp := handlers.ExecuteRecipeResp{}

	txHandleResBytes, err := pylonSDK.WaitAndGetTxData(txhash, pylonSDK.GetMaxWaitBlock(), t)
	if err != nil {
		errString := fmt.Sprintf("error getting tx result bytes %+v", err)
		log.Warnln(errString)
		LogFullTxResultByHash(txhash)
		return []byte{}, errString
	}
	LogFullTxResultByHash(txhash)
	hmrErrMsg := pylonSDK.GetHumanReadableErrorFromTxHash(txhash, t)
	if len(hmrErrMsg) > 0 {
		errString := fmt.Sprintf("txhash=%s hmrErrMsg=%s", txhash, hmrErrMsg)
		log.Warnln(errString)
		return []byte{}, errString
	}
	SyncFromNode(user)

	err = pylonSDK.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
	if err != nil {
		errString := fmt.Sprintf("failed to parse transaction result; maybe this is get_pylons then ignore. txhash=%s", txhash)
		log.Warnln(errString)
		return []byte{}, errString
	}
	log.Debugln("ProcessTxResult ==>", resp.Message, string(resp.Output))
	return resp.Output, ""
}

// GetTestingT is a function to convert testing.T to cusomized testing.T
func GetTestingT() *testing.T {
	newT := testing.NewT(nil)
	t := &newT
	return t
}

// ExecuteRecipe is a function to execute recipe by name and input items
func ExecuteRecipe(user User, rcpName string, itemIDs []string) (string, error) {
	t := GetTestingT()
	if len(rcpName) == 0 {
		return "", errors.New("Recipe Name does not exist")
	}
	rcpID, ok := RcpIDs[rcpName]
	if !ok {
		return "", errors.New("RecipeID does not exist for rcpName=" + rcpName)
	}
	addr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(addr)
	execMsg := msgs.NewMsgExecuteRecipe(rcpID, sdkAddr, itemIDs)
	log.Debugln("started sending transaction", user.GetUserName(), execMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, execMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash, rcpName)
	log.Debugln("ended sending transaction")
	return txhash, nil
}

// GetIndexFromString is a function to convert 1-9 string to 0-8 index
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

// GetToBuyItemFromKey returns which item to buy when user provide 1-9 key
func GetToBuyItemFromKey(key string) Item {
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(ShopItems) {
		useItem = ShopItems[itemKey]
	}
	return useItem
}

// GetToBuyCharacterFromKey returns which character to buy when user provide 1-9 key
func GetToBuyCharacterFromKey(key string) Character {
	useCharacter := Character{}
	cKey := GetIndexFromString(key)
	if cKey >= 0 && cKey < len(ShopCharacters) {
		useCharacter = ShopCharacters[cKey]
	}
	return useCharacter
}

// GetToSellItemFromKey returns which item to sell from 1-9 key
func GetToSellItemFromKey(user User, key string) Item {
	items := user.InventorySellableItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

// GetToUpgradeItemFromKey returns which item to upgrade from 1-9 key
func GetToUpgradeItemFromKey(user User, key string) Item {
	items := user.InventoryUpgradableItems()
	useItem := Item{}
	itemKey := GetIndexFromString(key)
	if itemKey >= 0 && itemKey < len(items) {
		useItem = items[itemKey]
	}
	return useItem
}

// GetItemInputsFromItemSpec calculate ItemInput from ItemSpec
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

// GetItemOutputFromActiveItem calculate ItemOutput from ActiveItem
func GetItemOutputFromActiveItem(activeItem Item) (types.ItemList, error) {
	var itemOutputs types.ItemList
	io, err := pylonSDK.GetItemByGUID(activeItem.ID)
	itemOutputs = append(itemOutputs, io)
	return itemOutputs, err
}

// GetItemInputsFromCharacterSpec calculate ItemInputs from CharacterSpec
func GetItemInputsFromCharacterSpec(chspec CharacterSpec) types.ItemInputList {
	// TODO should make this to express all the required fields like GiantKill, SpecialDragonKill, UndeadDragonKill
	// But for now expressing only XP, level, Name and Special as it's the main requirement.
	// If possible, we can try removing XP, and level too.

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
	if chspec.Special != NoSpecial {
		ii.Longs = append(ii.Longs, types.LongInputParam{
			Key: "Special", MinValue: chspec.Special, MaxValue: chspec.Special})
	}
	itemInputs = append(itemInputs, ii)
	return itemInputs
}

// GetItemOutputFromActiveCharacter calculate ItemOutput from ActiveCharacter
func GetItemOutputFromActiveCharacter(activeCharacter Character) (types.ItemList, error) {
	var itemOutputs types.ItemList
	io, err := pylonSDK.GetItemByGUID(activeCharacter.ID)
	itemOutputs = append(itemOutputs, io)
	return itemOutputs, err
}

// GetSDKAddrFromUserName convert key to sdk address
func GetSDKAddrFromUserName(username string) sdk.AccAddress {
	addr := pylonSDK.GetAccountAddr(username, nil)
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		log.Fatal("sdkAddr, err := sdk.AccAddressFromBech32(addr)", sdkAddr, err)
	}
	return sdkAddr
}

// SendTxMsg returns transaction from a user
func SendTxMsg(user User, txMsg sdk.Msg) (string, error) {
	t := GetTestingT()
	log.Debugln("started sending transaction", user.GetUserName(), txMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, txMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash, txMsg.Type())
	log.Debugln("ended sending transaction")
	return txhash, nil
}
