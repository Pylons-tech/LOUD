package screen

import "strings"

// PageStatus is a struct to manage screen's status
type PageStatus string

const (
	// ShowLocation describes show location param
	ShowLocation PageStatus = "SHOW_LOCATION"

	// ConfirmEndGame describes the page status for ending game confirmation page
	ConfirmEndGame = "CONFIRM_ENDGAME"

	// SelectActiveChr describes the page status for active character selection
	SelectActiveChr = "SEL_ACTIVE_CHAR"
	// RsltSelectActiveChr describes the page status for active character selection result
	RsltSelectActiveChr = "RSLT_SEL_ACT_CHAR"

	// SelectRenameChr describes the page status for rename character selection
	SelectRenameChr = "SEL_RENAME_CHAR"
	// SelectRenameChrEntNewName describes the page status for rename character selection enter new name
	SelectRenameChrEntNewName = "RENAME_CHAR_ENT_NEWNAME"
	// WaitRenameChr describes the status for renaming character wait screen
	WaitRenameChr = "W8_RENAME_CHAR"
	// RsltRenameChr describes the status for renaming character result screen
	RsltRenameChr = "RSLT_RENAME_CHAR"

	// SelectSellItem describes the status for sell item selection page
	SelectSellItem = "SEL_SELLITM"
	// WaitSellItem describes the status for sell item waiter page
	WaitSellItem = "W8_SELLITM"
	// RsltSellItem describes the status for sell item result page
	RsltSellItem = "RSLT_SELLITM"

	// SelectBuyItem describes the status for buying item page
	SelectBuyItem = "SEL_BUYITM"
	// WaitBuyItem describes the status for buying item wait page
	WaitBuyItem = "W8_BUYITM"
	// RsltBuyItem describes the status for buying item result page
	RsltBuyItem = "RSLT_BUYITM"

	// SelectUpgradeItem describes the status for upgrade item selection page
	SelectUpgradeItem = "SEL_UPGITM"
	// WaitUpgradeItem describes the status for upgrade item wait page
	WaitUpgradeItem = "W8_UPGITM"
	// RsltUpgradeItem describes the status for upgrade item result page
	RsltUpgradeItem = "RSLT_UPGITM"

	// ConfirmHuntRabbits describes the confirm hunt rabbits page
	ConfirmHuntRabbits = "CONFIRM_HUNT_RABBITS"
	// WaitHuntRabbits describes the confirm hunt rabbits page
	WaitHuntRabbits = "W8_HUNT_RABBITS"
	// RsltHuntRabbits describes the confirm hunt rabbits page
	RsltHuntRabbits = "RSLT_HUNT_RABBITS"

	// ConfirmFightGoblin describes the confirm fight goblin page
	ConfirmFightGoblin = "CONFIRM_FIGHT_GOBLIN"
	// WaitFightGoblin describes the wait fight goblin page
	WaitFightGoblin = "W8_FIGHT_GOBLIN"
	// RsltFightGoblin describes the wait fight goblin page
	RsltFightGoblin = "RSLT_FIGHT_GOBLIN"

	// ConfirmFightTroll describes the confirmation page for fighting troll
	ConfirmFightTroll = "CONFIRM_FIGHT_TROLL"
	// WaitFightTroll describes the wait page for fighting troll
	WaitFightTroll = "W8_FIGHT_TROLL"
	// RsltFightTroll describes the result page for fighting troll
	RsltFightTroll = "RSLT_FIGHT_TROLL"

	// ConfirmFightWolf describes the confirmation page for fighting wolf
	ConfirmFightWolf = "CONFIRM_FIGHT_WOLF"
	// WaitFightWolf describes the wait page for fighting wolf
	WaitFightWolf = "W8_FIGHT_WOLF"
	// RsltFightWolf describes the result page for fighting wolf
	RsltFightWolf = "RSLT_FIGHT_WOLF"

	// ConfirmFightGiant describes the confirmation page for fighting giant
	ConfirmFightGiant = "CONFIRM_FIGHT_GIANT"
	// WaitFightGiant describes the wait page for fighting giant
	WaitFightGiant = "W8_FIGHT_GIANT"
	// RsltFightGiant describes the result page for fighting giant
	RsltFightGiant = "RSLT_FIGHT_GIANT"

	// ConfirmFightDragonFire describes the confirmation page for fighting fire dragon
	ConfirmFightDragonFire = "CONFIRM_FIGHT_DRAGONFIRE"
	// WaitFightDragonFire describes the wait page for fighting fire dragon
	WaitFightDragonFire = "W8_FIGHT_DRAGONFIRE"
	// RsltFightDragonFire describes the result page for fighting fire dragon
	RsltFightDragonFire = "RSLT_FIGHT_DRAGONFIRE"

	// ConfirmFightDragonIce describes the confirmation page for fighting ice dragon
	ConfirmFightDragonIce = "CONFIRM_FIGHT_DRAGONICE"
	// WaitFightDragonIce describes the wait page for fighting ice dragon
	WaitFightDragonIce = "W8_FIGHT_DRAGONICE"
	// RsltFightDragonIce describes the result page for fighting ice dragon
	RsltFightDragonIce = "RSLT_FIGHT_DRAGONICE"

	// ConfirmFightDragonAcid describes the confirmation page for fighting acid dragon
	ConfirmFightDragonAcid = "CONFIRM_FIGHT_DRAGONACID"
	// WaitFightDragonAcid describes the wait page for fighting acid dragon
	WaitFightDragonAcid = "W8_FIGHT_DRAGONACID"
	// RsltFightDragonAcid describes the result page for fighting acid dragon
	RsltFightDragonAcid = "RSLT_FIGHT_DRAGONACID"

	// ConfirmFightDragonUndead describes the confirmation page for fighting undead dragon
	ConfirmFightDragonUndead = "CONFIRM_FIGHT_DRAGONUNDEAD"
	// WaitFightDragonUndead describes the wait page for fighting undead dragon
	WaitFightDragonUndead = "W8_FIGHT_DRAGONUNDEAD"
	// RsltFightDragonUndead describes the result page for fighting undead dragon
	RsltFightDragonUndead = "RSLT_FIGHT_DRAGONUNDEAD"

	// WaitCreateCookbook describes the wait page for creating cookbook
	WaitCreateCookbook = "W8_CREATE_COOKBOOK"
	// RsltCreateCookbook describes the result page for creating cookbook
	RsltCreateCookbook = "RSLT_CREATE_COOKBOOK"

	// WaitSwitchUser describes the wait page for switching user
	WaitSwitchUser = "W8_SWITCH_USER"
	// RsltSwitchUser describes the result page for switching user
	RsltSwitchUser = "RSLT_SWITCH_USER"

	// WaitDevGetTestItems describes the wait page for getting developer test items
	WaitDevGetTestItems = "W8_DEV_GET_TEST_ITEMS"
	// RsltDevGetTestItems describes the result page for getting developer test items
	RsltDevGetTestItems = "RSLT_DEV_GET_TEST_ITEMS"
	// WaitGetPylons describes the wait page for getting pylons
	WaitGetPylons = "W8_GET_PYLONS"
	// RsltGetPylons describes the result page for getting pylons
	RsltGetPylons = "RSLT_GET_PYLONS"

	// WaitByGoldWithPylons describes the wait page to buy gold with pylons
	WaitByGoldWithPylons = "W8_BUY_GOLD_WITH_PYLONS"
	// RsltByGoldWithPylons describes the result page to buy gold with pylons
	RsltByGoldWithPylons = "RSLT_BUY_GOLD_WITH_PYLONS"

	// SelectBuyChr describes the select page for buying character
	SelectBuyChr = "SEL_BUYCHR"
	// WaitBuyChr describes the wait page for buying character
	WaitBuyChr = "W8_BUYCHR"
	// RsltBuyChr describes the result page for buying character
	RsltBuyChr = "RSLT_BUYCHR"

	// ShowGoldBuyTrdReqs describes the listing page for buying gold trade requests
	ShowGoldBuyTrdReqs = "SHW_LOUD_BUY_TRDREQS"
	// CreateBuyGoldTrdReqEnterGoldValue describes the gold buy trade request creation page's enter gold amount
	CreateBuyGoldTrdReqEnterGoldValue = "CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL"
	// CreateBuyGoldTrdReqEnterPylonValue describes the gold buy trade request creation page's enter pylon amount
	CreateBuyGoldTrdReqEnterPylonValue = "CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL"
	// WaitBuyGoldTrdReqCreation describes wait page for the gold buy trade request creation
	WaitBuyGoldTrdReqCreation = "W8_BUY_LOUD_TRDREQ_CREATION"
	// RsltBuyGoldTrdReqCreation describes result page for the gold buy trade request creation
	RsltBuyGoldTrdReqCreation = "RSLT_BUY_LOUD_TRDREQ_CREATION"
	// WaitFulfillBuyGoldTrdReq describes the wait page for fulfilling the buying gold trade request
	WaitFulfillBuyGoldTrdReq = "W8_FULFILL_BUY_LOUD_TRDREQ"
	// RsltFulfillBuyGoldTrdReq describes the result page for fulfilling the buying gold trade request
	RsltFulfillBuyGoldTrdReq = "RSLT_FULFILL_BUY_LOUD_TRDREQ"
	// WaitCancelBuyGoldTrdReq describes the wait page to cancel trade request
	WaitCancelBuyGoldTrdReq = "W8_CANCEL_BUY_GOLD_TRDREQ"
	// RsltCancelBuyGoldTrdReq describes the result page to cancel trade request
	RsltCancelBuyGoldTrdReq = "RSLT_CANCEL_BUY_GOLD_TRDREQ"

	// ShowGoldSellTrdReqs describes the listing page for selling gold trade requests
	ShowGoldSellTrdReqs = "SHW_LOUD_SELL_TRDREQS"
	// CreateSellGoldTrdReqEnterGoldValue describes the gold sell trade request creation page's enter gold amount
	CreateSellGoldTrdReqEnterGoldValue = "CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL"
	// CreateSellGoldTrdReqEnterPylonValue describes the gold sell trade request creation page's enter pylon amount
	CreateSellGoldTrdReqEnterPylonValue = "CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL"
	// WaitSellGoldTrdReqCreation describes wait page for the gold sell trade request creation
	WaitSellGoldTrdReqCreation = "W8_SELL_LOUD_TRDREQ_CREATION"
	// RsltSellGoldTrdReqCreation describes result page for the gold sell trade request creation
	RsltSellGoldTrdReqCreation = "RSLT_SELL_LOUD_TRDREQ_CREATION"
	// WaitFulfillSellGoldTrdReq describes the wait page for fulfilling the selling gold trade request
	WaitFulfillSellGoldTrdReq = "W8_FULFILL_SELL_LOUD_TRDREQ"
	// RsltFulfillSellGoldTrdReq describes the result page for fulfilling the selling gold trade request
	RsltFulfillSellGoldTrdReq = "RSLT_FULFILL_SELL_LOUD_TRDREQ"
	// WaitCancelSellGoldTrdReq describes the wait page to cancel trade request
	WaitCancelSellGoldTrdReq = "W8_CANCEL_SELL_GOLD_TRDREQ"
	// RsltCancelSellGoldTrdReq describes the result page to cancel trade request
	RsltCancelSellGoldTrdReq = "RSLT_CANCEL_SELL_GOLD_TRDREQ"

	// ShowSellItemTrdReqs describes the listing page for selling item trade requests
	ShowSellItemTrdReqs = "SHW_SELLITM_TRDREQS"
	// CreateSellItemTrdReqSelectItem describes the item sell trade request creation page's select item
	CreateSellItemTrdReqSelectItem = "CR8_SELLITM_TRDREQ_SEL_ITEM"
	// CreateSellItemTrdReqEnterPylonValue describes the item sell trade request creation page's enter pylon amount
	CreateSellItemTrdReqEnterPylonValue = "CR8_SELLITM_TRDREQ_ENT_PYLVAL"
	// WaitSellItemTrdReqCreation describes wait page for the item sell trade request creation
	WaitSellItemTrdReqCreation = "W8_SELLITM_TRDREQ_CREATION"
	// RsltSellItemTrdReqCreation describes result page for the item sell trade request creation
	RsltSellItemTrdReqCreation = "RSLT_SELLITM_TRDREQ_CREATION"
	// WaitFulfillSellItemTrdReq describes the wait page for fulfilling the selling item trade request
	WaitFulfillSellItemTrdReq = "W8_FULFILL_SELLITM_TRDREQ"
	// RsltFulfillSellItemTrdReq describes the result page for fulfilling the selling item trade request
	RsltFulfillSellItemTrdReq = "RSLT_FULFILL_SELLITM_TRDREQ"
	// WaitCancelSellItemTrdReq describes the wait page to cancel trade request
	WaitCancelSellItemTrdReq = "W8_CANCEL_SELLITM_TRDREQ"
	// RsltCancelSellItemTrdReq describes the result page to cancel trade request
	RsltCancelSellItemTrdReq = "RSLT_CANCEL_SELLITM_TRDREQ"

	// ShowBuyItemTrdReqs describes the listing page for buying item trade requests
	ShowBuyItemTrdReqs = "SHW_BUYITM_TRDREQS"
	// SelectFitBuyItemTrdReq describes the page for selecting buy item trade request
	SelectFitBuyItemTrdReq = "SEL_FIT_BUYITM_TRDREQ"
	// CreateBuyItemTrdReqSelectItem describes the item buying trade request creation page's select item
	CreateBuyItemTrdReqSelectItem = "CR8_BUYITM_TRDREQ_SEL_ITEM"
	// CreateBuyItmTrdReqEnterPylonValue describes the item buying trade request creation page's enter pylon amount
	CreateBuyItmTrdReqEnterPylonValue = "CR8_BUYITM_TRDREQ_ENT_PYLVAL"
	// WaitBuyItemTrdReqCreation describes wait page for the item buying trade request creation
	WaitBuyItemTrdReqCreation = "W8_BUYITM_TRDREQ_CREATION"
	// RsltBuyItemTrdReqCreation describes result page for the item buying trade request creation
	RsltBuyItemTrdReqCreation = "RSLT_BUYITM_TRDREQ_CREATION"
	// WaitFulfillBuyItemTrdReq describes the wait page for fulfilling the buying item trade request
	WaitFulfillBuyItemTrdReq = "W8_FULFILL_BUYITM_TRDREQ"
	// RsltFulfillBuyItemTrdReq describes the result page for fulfilling the buying item trade request
	RsltFulfillBuyItemTrdReq = "RSLT_FULFILL_BUYITM_TRDREQ"
	// WaitCancelBuyItemTrdReq describes the wait page to cancel trade request
	WaitCancelBuyItemTrdReq = "W8_CANCEL_BUYITM_TRDREQ"
	// RsltCancelBuyItemTrdReq describes the result page to cancel trade request
	RsltCancelBuyItemTrdReq = "RSLT_CANCEL_BUYITM_TRDREQ"

	// ShowSellChrTrdReqs describes the listing page for selling character trade requests
	ShowSellChrTrdReqs = "SHW_SELLCHR_TRDREQS"
	// CreateSellChrTrdReqSelChr describes the character sell trade request creation page's select character
	CreateSellChrTrdReqSelChr = "CR8_SELLCHR_TRDREQ_SEL_CHR"
	// CreateSellChrTrdReqEnterPylonValue describes the character sell trade request creation page's enter pylon amount
	CreateSellChrTrdReqEnterPylonValue = "CR8_SELLCHR_TRDREQ_ENT_PYLVAL"
	// WaitSellChrTrdReqCreation describes wait page for the character sell trade request creation
	WaitSellChrTrdReqCreation = "W8_SELLCHR_TRDREQ_CREATION"
	// RsltSellChrTrdReqCreation describes result page for the character sell trade request creation
	RsltSellChrTrdReqCreation = "RSLT_SELLCHR_TRDREQ_CREATION"
	// WaitFulfillSellChrTrdReq describes the wait page for fulfilling the sell character trade request
	WaitFulfillSellChrTrdReq = "W8_FULFILL_SELLCHR_TRDREQ"
	// RsltFulfillSellChrTrdReq describes the result page for fulfilling the sell character trade request
	RsltFulfillSellChrTrdReq = "RSLT_FULFILL_SELLCHR_TRDREQ"
	// WaitCancelSellChrTrdReq describes the wait page to cancel trade request
	WaitCancelSellChrTrdReq = "W8_CANCEL_SELLCHR_TRDREQ"
	// RsltCancelSellChrTrdReq describes the result page to cancel trade request
	RsltCancelSellChrTrdReq = "RSLT_CANCEL_SELLCHR_TRDREQ"

	// ShowBuyChrTrdReqs describes the listing page for buying character trade requests
	ShowBuyChrTrdReqs = "SHW_BUYCHR_TRDREQS"
	// SelectFitBuyChrTrdReq describes the page for selecting buy item trade request
	SelectFitBuyChrTrdReq = "SEL_FIT_BUYCHR_TRDREQ"
	// CreateBuyChrTrdReqSelectChr describes the character buy trade request creation page's select character
	CreateBuyChrTrdReqSelectChr = "CR8_BUYCHR_TRDREQ_SEL_CHR"
	// CreateBuyChrTrdReqEnterPylonValue describes the character buy trade request creation page's enter pylon amount
	CreateBuyChrTrdReqEnterPylonValue = "CR8_BUYCHR_TRDREQ_ENT_PYLVAL"
	// WaitBuyChrTrdReqCreation describes wait page for the character buy trade request creation
	WaitBuyChrTrdReqCreation = "W8_BUYCHR_TRDREQ_CREATION"
	// RsltBuyChrTrdReqCreation describes result page for the character buy trade request creation
	RsltBuyChrTrdReqCreation = "RSLT_BUYCHR_TRDREQ_CREATION"
	// WaitFulfillBuyChrTrdReq describes the wait page for fulfilling the buy character trade request
	WaitFulfillBuyChrTrdReq = "W8_FULFILL_BUYCHR_TRDREQ"
	// RsltFulfillBuyChrTrdReq describes the result page for fulfilling the buy character trade request
	RsltFulfillBuyChrTrdReq = "RSLT_FULFILL_BUYCHR_TRDREQ"
	// WaitCancelBuyChrTrdReq describes the wait page to cancel trade request
	WaitCancelBuyChrTrdReq = "W8_CANCEL_BUYCHR_TRDREQ"
	// RsltCancelBuyChrTrdReq describes the result page to cancel trade request
	RsltCancelBuyChrTrdReq = "RSLT_CANCEL_BUYCHR_TRDREQ"

	// HelpAbout describes the help page about the game
	HelpAbout = "HELP_ABOUT"
	// HelpGameObjective describes the help page for game objective
	HelpGameObjective = "HELP_GAME_OBJECTIVE"
	// HelpNavigation describes the help page for navigation
	HelpNavigation = "HELP_NAVIGATION"
	// HelpPageLayout describes the help page for page layout
	HelpPageLayout = "HELP_PAGE_LAYOUT"
	// HelpGameRules describes the help page for game rules
	HelpGameRules = "HELP_GAME_RULES"
	// HelpHowItWorks describes the help page for how it works section
	HelpHowItWorks = "HELP_HOW_IT_WORKS"
	// HelpPylonsCentral describes the help page for pylons central
	HelpPylonsCentral = "HELP_PYLONS_CENTRAL"
	// HelpUpcomingReleases describes the help page for upcoming releases
	HelpUpcomingReleases = "HELP_UPCOMING_RELEASES"
	// HelpSupport describes the help page for support
	HelpSupport = "HELP_SUPPORT"
)

// IsWaitScreen returns if page status is related to wait
func (status PageStatus) IsWaitScreen() bool {
	return strings.Contains(string(status), "W8_")
}

// IsResultScreen returns if page status is result screen
func (status PageStatus) IsResultScreen() bool {
	return strings.Contains(string(status), "RSLT_")
}

// IsHelpScreen returns if page status is help screen
func (status PageStatus) IsHelpScreen() bool {
	return strings.Contains(string(status), "HELP_")
}
