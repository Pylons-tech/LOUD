package loud

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Pylons-tech/LOUD/log"
	pylonSDK "github.com/Pylons-tech/pylons_sdk/cmd/test"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const CR8BY_LOUD = "created by loud game"
const ITEM_BUYREQ_TRDINFO = "sword buy request created by loud game"
const CHAR_BUYREQ_TRDINFO = "character buy request created by loud game"
const ITEM_SELREQ_TRDINFO = "sword sell request created by loud game"
const CHAR_SELREQ_TRDINFO = "character sell request created by loud game"

func CreateCookbook(user User) (string, error) { // This is for afti develop mode automation test is only using
	t := GetTestingT()
	username := user.GetUserName()
	sdkAddr := GetSDKAddrFromUserName(username)

	ccbMsg := msgs.NewMsgCreateCookbook(
		"tst_cookbook_name",                  // cbType.Name,
		fmt.Sprintf("%d", time.Now().Unix()), // cbType.ID,
		"addghjkllsdfdggdgjkkk",              // cbType.Description,
		"asdfasdfasdf",                       // cbType.Developer,
		"1.0.0",                              // cbType.Version,
		"a@example.com",                      // cbType.SupportEmail,
		0,                                    // cbType.Level,
		5,                                    // cbType.CostPerBlock,
		sdkAddr,                              // cbType.Sender,
	)

	txhash, _ := SendTxMsg(user, ccbMsg)
	if AutomateInput {
		ok, err := CheckSignatureMatchWithAftiCli(t, txhash, user.GetPrivKey(), ccbMsg, username, false)
		if !ok || err != nil {
			log.Println("error checking afticli", ok, err)
			SomethingWentWrongMsg = "automation test failed, " + err.Error()
		}
	}
	return txhash, nil
}

func GetExtraPylons(user User) (string, error) {
	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())
	extraPylonsMsg := msgs.NewMsgGetPylons(types.PremiumTier.Fee, sdkAddr)
	return SendTxMsg(user, extraPylonsMsg)
}

func BuyGoldWithPylons(user User) (string, error) {
	return ExecuteRecipe(user, RCP_BUY_GOLD_WITH_PYLON, []string{})
}

func DevGetTestItems(user User) (string, error) {
	return ExecuteRecipe(user, RCP_GET_TEST_ITEMS, []string{})
}

func RunHuntRecipe(monsterName, rcpName string, user User) (string, error) {
	activeCharacter := user.GetActiveCharacter()
	activeCharacterID := ""
	if activeCharacter != nil {
		activeCharacterID = activeCharacter.ID
	} else {
		return "", errors.New("character is required to hunt rabbits!")
	}

	user.SetFightMonster(monsterName)
	activeWeapon := user.GetFightWeapon()
	itemIDs := []string{activeCharacterID}
	if activeWeapon != nil {
		itemIDs = []string{activeCharacterID, activeWeapon.ID}
	}

	return ExecuteRecipe(user, rcpName, itemIDs)
}

func HuntRabbits(user User) (string, error) {
	return RunHuntRecipe(TextRabbit, RCP_HUNT_RABBITS, user)
}

func FightTroll(user User) (string, error) {
	return RunHuntRecipe(TextTroll, RCP_FIGHT_TROLL, user)
}

func FightWolf(user User) (string, error) { // 🐺
	return RunHuntRecipe(TextWolf, RCP_FIGHT_WOLF, user)
}

func FightGoblin(user User) (string, error) { // 👺
	return RunHuntRecipe(TextGoblin, RCP_FIGHT_GOBLIN, user)
}

func FightGiant(user User) (string, error) { // 🗿
	return RunHuntRecipe(TextGiant, RCP_FIGHT_GIANT, user)
}

func FightDragonFire(user User) (string, error) { // 🦐
	return RunHuntRecipe(TextDragonFire, RCP_FIGHT_DRAGONFIRE, user)
}

func FightDragonIce(user User) (string, error) { // 🦈
	return RunHuntRecipe(TextDragonIce, RCP_FIGHT_DRAGONICE, user)
}

func FightDragonAcid(user User) (string, error) { // 🐊
	return RunHuntRecipe(TextDragonAcid, RCP_FIGHT_DRAGONACID, user)
}

func FightDragonUndead(user User) (string, error) { // 🐉
	return RunHuntRecipe(TextDragonUndead, RCP_FIGHT_DRAGONUNDEAD, user)
}

func BuyCharacter(user User, ch Character) (string, error) {
	rcpName := ""
	switch ch.Name {
	case TextTigerChr:
		rcpName = RCP_BUY_CHARACTER
	default:
		return "", errors.New("You are trying to buy character which is not in shop")
	}
	if ch.Price > user.GetPylonAmount() {
		return "", errors.New("You don't have enough pylon to buy this character")
	}
	return ExecuteRecipe(user, rcpName, []string{})
}

func RenameCharacter(user User, ch Character, newName string) (string, error) {
	t := GetTestingT()
	addr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(addr)
	renameMsg := msgs.NewMsgUpdateItemString(ch.ID, "Name", newName, sdkAddr)
	log.Println("started sending transaction", user.GetUserName(), renameMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, renameMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash, Sprintf("rename character from %s to %s", ch.Name, newName))
	log.Println("ended sending transaction")
	return txhash, nil
}

func Buy(user User, item Item) (string, error) {
	rcpName := ""
	itemIDs := []string{}
	switch item.Name {
	case WoodenSword:
		if item.Level == 1 {
			rcpName = RCP_BUY_WOODEN_SWORD
		}
	case CopperSword:
		if item.Level == 1 {
			rcpName = RCP_BUY_COPPER_SWORD
		}
	case SilverSword:
		if item.Level == 1 {
			rcpName = RCP_BUY_SILVER_SWORD
			itemIDs = []string{user.InventoryItemIDByName(GoblinEar)}
		}
	case BronzeSword:
		if item.Level == 1 {
			rcpName = RCP_BUY_BRONZE_SWORD
			itemIDs = []string{user.InventoryItemIDByName(WolfTail)}
		}
	case IronSword:
		if item.Level == 1 {
			rcpName = RCP_BUY_IRON_SWORD
			itemIDs = []string{user.InventoryItemIDByName(TrollToes)}
		}
	case AngelSword:
		if item.Level == 1 {
			rcpName = RCP_BUY_ANGEL_SWORD
			itemIDs = []string{
				user.InventoryItemIDByName(DropDragonFire),
				user.InventoryItemIDByName(DropDragonIce),
				user.InventoryItemIDByName(DropDragonAcid),
			}
		}
	default:
		return "", errors.New("You are trying to buy item which is not in shop")
	}
	if item.Price > user.GetGold() {
		return "", errors.New("You don't have enough gold to buy this item")
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

func Sell(user User, item Item) (string, error) {
	itemIDs := []string{item.ID}

	rcpName := ""
	if item.Value > 0 {
		rcpName = RCP_SELL_SWORD
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

func Upgrade(user User, item Item) (string, error) {
	itemIDs := []string{item.ID}
	rcpName := ""
	switch item.Name {
	case WoodenSword:
		if item.Level == 1 {
			rcpName = RCP_WOODEN_SWORD_UPG
		}
	case CopperSword:
		if item.Level == 1 {
			rcpName = RCP_COPPER_SWORD_UPG
		}
	}
	if item.GetUpgradePrice() > user.GetGold() {
		return "", errors.New("You don't have enough gold to upgrade this item")
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

func CreateBuyLoudTrdReq(user User, loudEnterValue string, pylonEnterValue string) (string, error) {
	loudValue, err := strconv.Atoi(loudEnterValue)
	if err != nil {
		return "", err
	}
	if loudValue == 0 {
		return "", errors.New("gold amount shouldn't be zero to be a valid trading")
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("loudcoin", int64(loudValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := CR8BY_LOUD

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

func CreateSellLoudTrdReq(user User, loudEnterValue string, pylonEnterValue string) (string, error) {
	loudValue, err := strconv.Atoi(loudEnterValue)
	if err != nil {
		return "", err
	}
	if loudValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("loudcoin", int64(loudValue))}
	extraInfo := CR8BY_LOUD

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

func CreateBuyItemTrdReq(user User, itspec ItemSpec, pylonEnterValue string) (string, error) {
	// trade creator will get sword from pylon

	itemInputs := GetItemInputsFromItemSpec(itspec)

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := ITEM_BUYREQ_TRDINFO

	createTrdMsg := msgs.NewMsgCreateTrade(
		nil,
		itemInputs,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

func CreateSellItemTrdReq(user User, activeItem Item, pylonEnterValue string) (string, error) {
	// trade creator will get pylon from sword

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))
	itemOutputList, err := GetItemOutputFromActiveItem(activeItem)
	if err != nil {
		return "", err
	}

	extraInfo := ITEM_SELREQ_TRDINFO

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		nil,
		itemOutputList,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

func CreateBuyCharacterTrdReq(user User, chspec CharacterSpec, pylonEnterValue string) (string, error) {
	// trade creator will get character from pylon

	itemInputs := GetItemInputsFromCharacterSpec(chspec)

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := CHAR_BUYREQ_TRDINFO

	createTrdMsg := msgs.NewMsgCreateTrade(
		nil,
		itemInputs,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

func CreateSellCharacterTrdReq(user User, activeCharacter Character, pylonEnterValue string) (string, error) {
	// trade creator will get pylon from character

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))
	itemOutputList, err := GetItemOutputFromActiveCharacter(activeCharacter)
	if err != nil {
		return "", err
	}

	extraInfo := CHAR_SELREQ_TRDINFO

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		nil,
		itemOutputList,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

func FulfillTrade(user User, tradeID string) (string, error) {
	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())
	ffTrdMsg := msgs.NewMsgFulfillTrade(tradeID, sdkAddr, []string{})

	return SendTxMsg(user, ffTrdMsg)
}

func CancelTrade(user User, tradeID string) (string, error) {
	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())
	ccTrdMsg := msgs.NewMsgDisableTrade(tradeID, sdkAddr)

	return SendTxMsg(user, ccTrdMsg)
}
