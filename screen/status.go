package screen

import "strings"

type ScreenStatus string

const (
	SHW_LOCATION ScreenStatus = "SHW_LOCATION"
	// at home
	SEL_ACTIVE_CHAR   = "SEL_ACTIVE_CHAR"
	RSLT_SEL_ACT_CHAR = "RSLT_SEL_ACT_CHAR"

	SEL_ACTIVE_WEAPON   = "SEL_ACTIVE_WEAPON"
	RSLT_SEL_ACT_WEAPON = "RSLT_SEL_ACT_WEAPON"

	SEL_HEALTH_RESTORE_CHAR  = "SEL_HEALTH_RESTORE_CHAR"
	W8_HEALTH_RESTORE_CHAR   = "W8_HEALTH_RESTORE_CHAR"
	RSLT_HEALTH_RESTORE_CHAR = "RSLT_HEALTH_RESTORE_CHAR"

	SEL_RENAME_CHAR         = "SEL_RENAME_CHAR"
	RENAME_CHAR_ENT_NEWNAME = "RENAME_CHAR_ENT_NEWNAME"
	W8_RENAME_CHAR          = "W8_RENAME_CHAR"
	RSLT_RENAME_CHAR        = "RSLT_RENAME_CHAR"
	// in shop
	SEL_SELLITM  = "SEL_SELLITM"
	W8_SELLITM   = "W8_SELLITM"
	RSLT_SELLITM = "RSLT_SELLITM"

	SEL_BUYITM  = "SEL_BUYITM"
	W8_BUYITM   = "W8_BUYITM"
	RSLT_BUYITM = "RSLT_BUYITM"

	SEL_UPGITM  = "SEL_UPGITM"
	W8_UPGITM   = "W8_UPGITM"
	RSLT_UPGITM = "RSLT_UPGITM"

	// in forest
	SEL_HUNT_RABBITS_ITEM = "SEL_HUNT_RABBITS_ITEM"
	W8_HUNT_RABBITS       = "W8_HUNT_RABBITS"
	RSLT_HUNT_RABBITS     = "RSLT_HUNT_RABBITS"

	SEL_FIGHT_GOBLIN_ITEM = "SEL_FIGHT_GOBLIN_ITEM"
	W8_FIGHT_GOBLIN       = "W8_FIGHT_GOBLIN"
	RSLT_FIGHT_GOBLIN     = "RSLT_FIGHT_GOBLIN"

	SEL_FIGHT_TROLL_ITEM = "SEL_FIGHT_TROLL_ITEM"
	W8_FIGHT_TROLL       = "W8_FIGHT_TROLL"
	RSLT_FIGHT_TROLL     = "RSLT_FIGHT_TROLL"

	SEL_FIGHT_WOLF_ITEM = "SEL_FIGHT_WOLF_ITEM"
	W8_FIGHT_WOLF       = "W8_FIGHT_WOLF"
	RSLT_FIGHT_WOLF     = "RSLT_FIGHT_WOLF"

	SEL_FIGHT_GIANT_ITEM = "SEL_FIGHT_GIANT_ITEM"
	W8_FIGHT_GIANT       = "W8_FIGHT_GIANT"
	RSLT_FIGHT_GIANT     = "RSLT_FIGHT_GIANT"

	// in develop
	W8_CREATE_COOKBOOK   = "W8_CREATE_COOKBOOK"
	RSLT_CREATE_COOKBOOK = "RSLT_CREATE_COOKBOOK"

	W8_SWITCH_USER   = "W8_SWITCH_USER"
	RSLT_SWITCH_USER = "RSLT_SWITCH_USER"

	W8_DEV_GET_TEST_ITEMS   = "W8_DEV_GET_TEST_ITEMS"
	RSLT_DEV_GET_TEST_ITEMS = "RSLT_DEV_GET_TEST_ITEMS"
	W8_GET_PYLONS           = "W8_GET_PYLONS"
	RSLT_GET_PYLONS         = "RSLT_GET_PYLONS"

	// in pylons central
	W8_BUY_GOLD_WITH_PYLONS   = "W8_BUY_GOLD_WITH_PYLONS"
	RSLT_BUY_GOLD_WITH_PYLONS = "RSLT_BUY_GOLD_WITH_PYLONS"

	SEL_BUYCHR  = "SEL_BUYCHR"
	W8_BUYCHR   = "W8_BUYCHR"
	RSLT_BUYCHR = "RSLT_BUYCHR"

	SHW_LOUD_BUY_TRDREQS           = "SHW_LOUD_BUY_TRDREQS"           // navigation using arrow and list should be sorted by price
	CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL = "CR8_BUY_LOUD_TRDREQ_ENT_LUDVAL" // enter value after switching enter mode
	CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL = "CR8_BUY_LOUD_TRDREQ_ENT_PYLVAL"
	W8_BUY_LOUD_TRDREQ_CREATION    = "W8_BUY_LOUD_TRDREQ_CREATION"
	RSLT_BUY_LOUD_TRDREQ_CREATION  = "RSLT_BUY_LOUD_TRDREQ_CREATION"
	W8_FULFILL_BUY_LOUD_TRDREQ     = "W8_FULFILL_BUY_LOUD_TRDREQ" // after done go to show loud buy requests
	RSLT_FULFILL_BUY_LOUD_TRDREQ   = "RSLT_FULFILL_BUY_LOUD_TRDREQ"

	SHW_LOUD_SELL_TRDREQS           = "SHW_LOUD_SELL_TRDREQS"
	CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL = "CR8_SELL_LOUD_TRDREQ_ENT_LUDVAL"
	CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL = "CR8_SELL_LOUD_TRDREQ_ENT_PYLVAL"
	W8_SELL_LOUD_TRDREQ_CREATION    = "W8_SELL_LOUD_TRDREQ_CREATION"
	RSLT_SELL_LOUD_TRDREQ_CREATION  = "RSLT_SELL_LOUD_TRDREQ_CREATION"
	W8_FULFILL_SELL_LOUD_TRDREQ     = "W8_FULFILL_SELL_LOUD_TRDREQ"
	RSLT_FULFILL_SELL_LOUD_TRDREQ   = "RSLT_FULFILL_SELL_LOUD_TRDREQ"

	SHW_SELLITM_TRDREQS           = "SHW_SELLITM_TRDREQS"
	CR8_SELLITM_TRDREQ_SEL_ITEM   = "CR8_SELLITM_TRDREQ_SEL_ITEM"
	CR8_SELLITM_TRDREQ_ENT_PYLVAL = "CR8_SELLITM_TRDREQ_ENT_PYLVAL"
	W8_SELLITM_TRDREQ_CREATION    = "W8_SELLITM_TRDREQ_CREATION"
	RSLT_SELLITM_TRDREQ_CREATION  = "RSLT_SELLITM_TRDREQ_CREATION"
	W8_FULFILL_SELLITM_TRDREQ     = "W8_FULFILL_SELLITM_TRDREQ"
	RSLT_FULFILL_SELLITM_TRDREQ   = "RSLT_FULFILL_SELLITM_TRDREQ"

	SHW_BUYITM_TRDREQS           = "SHW_BUYITM_TRDREQS"
	CR8_BUYITM_TRDREQ_SEL_ITEM   = "CR8_BUYITM_TRDREQ_SEL_ITEM"
	CR8_BUYITM_TRDREQ_ENT_PYLVAL = "CR8_BUYITM_TRDREQ_ENT_PYLVAL"
	W8_BUYITM_TRDREQ_CREATION    = "W8_BUYITM_TRDREQ_CREATION"
	RSLT_BUYITM_TRDREQ_CREATION  = "RSLT_BUYITM_TRDREQ_CREATION"
	W8_FULFILL_BUYITM_TRDREQ     = "W8_FULFILL_BUYITM_TRDREQ"
	RSLT_FULFILL_BUYITM_TRDREQ   = "RSLT_FULFILL_BUYITM_TRDREQ"

	SHW_SELLCHR_TRDREQS           = "SHW_SELLCHR_TRDREQS"
	CR8_SELLCHR_TRDREQ_SEL_CHR    = "CR8_SELLCHR_TRDREQ_SEL_CHR"
	CR8_SELLCHR_TRDREQ_ENT_PYLVAL = "CR8_SELLCHR_TRDREQ_ENT_PYLVAL"
	W8_SELLCHR_TRDREQ_CREATION    = "W8_SELLCHR_TRDREQ_CREATION"
	RSLT_SELLCHR_TRDREQ_CREATION  = "RSLT_SELLCHR_TRDREQ_CREATION"
	W8_FULFILL_SELLCHR_TRDREQ     = "W8_FULFILL_SELLCHR_TRDREQ"
	RSLT_FULFILL_SELLCHR_TRDREQ   = "RSLT_FULFILL_SELLCHR_TRDREQ"

	SHW_BUYCHR_TRDREQS           = "SHW_BUYCHR_TRDREQS"
	CR8_BUYCHR_TRDREQ_SEL_CHR    = "CR8_BUYCHR_TRDREQ_SEL_CHR"
	CR8_BUYCHR_TRDREQ_ENT_PYLVAL = "CR8_BUYCHR_TRDREQ_ENT_PYLVAL"
	W8_BUYCHR_TRDREQ_CREATION    = "W8_BUYCHR_TRDREQ_CREATION"
	RSLT_BUYCHR_TRDREQ_CREATION  = "RSLT_BUYCHR_TRDREQ_CREATION"
	W8_FULFILL_BUYCHR_TRDREQ     = "W8_FULFILL_BUYCHR_TRDREQ"
	RSLT_FULFILL_BUYCHR_TRDREQ   = "RSLT_FULFILL_BUYCHR_TRDREQ"

	W8_CANCEL_TRDREQ   = "W8_CANCEL_TRDREQ"
	RSLT_CANCEL_TRDREQ = "RSLT_CANCEL_TRDREQ"
)

func (status ScreenStatus) IsWaitScreen() bool {
	return strings.Contains(string(status), "W8_")
}

func (status ScreenStatus) IsResultScreen() bool {
	return strings.Contains(string(status), "RSLT_")
}
