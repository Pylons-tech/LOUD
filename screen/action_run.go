package screen

import (
	"time"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/LOUD/log"
)

func (screen *GameScreen) RunTxProcess(waitStatus ScreenStatus, resultStatus ScreenStatus, fn func() (string, error)) {
	screen.SetScreenStatusAndRefresh(waitStatus)

	log.Println("started sending request for ", waitStatus)
	go func() {
		txhash, err := fn()
		log.Println("ended sending request for ", waitStatus)
		if err != nil {
			screen.txFailReason = err.Error()
			screen.SetScreenStatusAndRefresh(resultStatus)
		} else {
			time.AfterFunc(1*time.Second, func() {
				screen.txResult, screen.txFailReason = loud.ProcessTxResult(screen.user, txhash)
				screen.SetScreenStatusAndRefresh(resultStatus)
			})
		}
	}()
}

func (screen *GameScreen) RunActiveCharacterSelect(index int) {
	screen.user.SetActiveCharacterIndex(index)
	screen.SetScreenStatusAndRefresh(RSLT_SEL_ACT_CHAR)
}

func (screen *GameScreen) RunActiveWeaponSelect(index int) {
	screen.user.SetActiveWeaponIndex(index)
	screen.SetScreenStatusAndRefresh(RSLT_SEL_ACT_WEAPON)
}

func (screen *GameScreen) RunCharacterHealthRestore() {
	screen.RunTxProcess(W8_HEALTH_RESTORE_CHAR, RSLT_HEALTH_RESTORE_CHAR, func() (string, error) {
		return loud.RestoreHealth(screen.user, screen.activeCharacter)
	})
}

func (screen *GameScreen) RunCharacterRename(newName string) {
	screen.RunTxProcess(W8_RENAME_CHAR, RSLT_RENAME_CHAR, func() (string, error) {
		return loud.RenameCharacter(screen.user, screen.activeCharacter, newName)
	})
}

func (screen *GameScreen) RunActiveItemBuy() {
	if !screen.user.HasPreItemForAnItem(screen.activeItem) {
		screen.txFailReason = loud.Sprintf("You don't have required item to make %s", screen.activeItem.Name)
		screen.SetScreenStatusAndRefresh(RSLT_BUYITM)
		return
	}
	screen.RunTxProcess(W8_BUYITM, RSLT_BUYITM, func() (string, error) {
		return loud.Buy(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveCharacterBuy() {
	screen.RunTxProcess(W8_BUYCHR, RSLT_BUYCHR, func() (string, error) {
		return loud.BuyCharacter(screen.user, screen.activeCharacter)
	})
}

func (screen *GameScreen) RunActiveItemSell() {
	screen.RunTxProcess(W8_SELLITM, RSLT_SELLITM, func() (string, error) {
		return loud.Sell(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunActiveItemUpgrade() {
	screen.RunTxProcess(W8_UPGITM, RSLT_UPGITM, func() (string, error) {
		return loud.Upgrade(screen.user, screen.activeItem)
	})
}

func (screen *GameScreen) RunHuntRabbits() {
	screen.RunTxProcess(W8_HUNT_RABBITS, RSLT_HUNT_RABBITS, func() (string, error) {
		return loud.HuntRabbits(screen.user)
	})
}

func (screen *GameScreen) RunFightGiant() {
	activeWeapon := screen.user.GetActiveWeapon()
	if activeWeapon == nil || activeWeapon.Name != loud.IRON_SWORD {
		screen.actionText = loud.Sprintf("You can't fight giant without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(W8_FIGHT_GIANT, RSLT_FIGHT_GIANT, func() (string, error) {
		return loud.FightGiant(screen.user)
	})
}

func (screen *GameScreen) RunFightTroll() {
	screen.RunTxProcess(W8_FIGHT_TROLL, RSLT_FIGHT_TROLL, func() (string, error) {
		return loud.FightTroll(screen.user)
	})
}

func (screen *GameScreen) RunFightWolf() {
	screen.RunTxProcess(W8_FIGHT_WOLF, RSLT_FIGHT_WOLF, func() (string, error) {
		return loud.FightWolf(screen.user)
	})
}

func (screen *GameScreen) RunFightGoblin() {
	screen.RunTxProcess(W8_FIGHT_GOBLIN, RSLT_FIGHT_GOBLIN, func() (string, error) {
		return loud.FightGoblin(screen.user)
	})
}

func (screen *GameScreen) RunSelectedLoudBuyTrdReq() {
	if len(loud.BuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		// when activeLine is not refering to real request but when it is refering to nil request
		screen.txFailReason = loud.Localize("you haven't selected any buy request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_BUY_LOUD_TRDREQ)
	} else {
		screen.activeTrdReq = loud.BuyTrdReqs[screen.activeLine]
		if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_BUY_LOUD_TRDREQ, RSLT_FULFILL_BUY_LOUD_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedLoudSellTrdReq() {
	if len(loud.SellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_SELL_LOUD_TRDREQ)
	} else {
		screen.activeTrdReq = loud.SellTrdReqs[screen.activeLine]
		if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_SELL_LOUD_TRDREQ, RSLT_FULFILL_SELL_LOUD_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedItemBuyTrdReq() {
	if len(loud.ItemBuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy item request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_BUYITM_TRDREQ)
	} else {
		atir := loud.ItemBuyTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = atir
		if atir.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, atir.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_BUYITM_TRDREQ, RSLT_FULFILL_BUYITM_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, atir.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedItemSellTrdReq() {
	if len(loud.ItemSellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell item request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_SELLITM_TRDREQ)
	} else {
		sstr := loud.ItemSellTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = sstr
		if sstr.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, sstr.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_SELLITM_TRDREQ, RSLT_FULFILL_SELLITM_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, sstr.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedCharacterBuyTrdReq() {
	if len(loud.CharacterBuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy character request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_BUYCHR_TRDREQ)
	} else {
		cbtr := loud.CharacterBuyTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = cbtr
		if cbtr.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, cbtr.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_BUYCHR_TRDREQ, RSLT_FULFILL_BUYCHR_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, cbtr.ID)
			})
		}
	}
}

func (screen *GameScreen) RunSelectedCharacterSellTrdReq() {
	if len(loud.CharacterSellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell character request")
		screen.SetScreenStatusAndRefresh(RSLT_FULFILL_SELLCHR_TRDREQ)
	} else {
		cstr := loud.CharacterSellTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = cstr
		if cstr.IsMyTrdReq {
			screen.RunTxProcess(W8_CANCEL_TRDREQ, RSLT_CANCEL_TRDREQ, func() (string, error) {
				return loud.CancelTrade(screen.user, cstr.ID)
			})
		} else {
			screen.RunTxProcess(W8_FULFILL_SELLCHR_TRDREQ, RSLT_FULFILL_SELLCHR_TRDREQ, func() (string, error) {
				return loud.FulfillTrade(screen.user, cstr.ID)
			})
		}
	}
}
