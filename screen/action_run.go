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
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight giant without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(W8_FIGHT_GIANT, RSLT_FIGHT_GIANT, func() (string, error) {
		return loud.FightGiant(screen.user)
	})
}

func (screen *GameScreen) RunFightDragonFire() {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight fire dragon without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(W8_FIGHT_DRAGONFIRE, RSLT_FIGHT_DRAGONFIRE, func() (string, error) {
		return loud.FightDragonFire(screen.user)
	})
}

func (screen *GameScreen) RunFightDragonIce() {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight ice dragon without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(W8_FIGHT_DRAGONICE, RSLT_FIGHT_DRAGONICE, func() (string, error) {
		return loud.FightDragonIce(screen.user)
	})
}

func (screen *GameScreen) RunFightDragonAcid() {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight acid dragon without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(W8_FIGHT_DRAGONACID, RSLT_FIGHT_DRAGONACID, func() (string, error) {
		return loud.FightDragonAcid(screen.user)
	})
}

func (screen *GameScreen) RunFightDragonUndead() {
	if len(screen.user.InventoryAngelSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight undead dragon without angel sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(W8_FIGHT_DRAGONUNDEAD, RSLT_FIGHT_DRAGONUNDEAD, func() (string, error) {
		return loud.FightDragonUndead(screen.user)
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
		if screen.user.GetGold() < screen.activeTrdReq.Amount {
			screen.actionText = loud.Sprintf("You don't have enough gold to fulfill this trade.")
			screen.Render()
		} else if screen.activeTrdReq.IsMyTrdReq {
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
		if screen.user.GetPylonAmount() < screen.activeTrdReq.Total {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else if screen.activeTrdReq.IsMyTrdReq {
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
		if len(screen.user.GetMatchedItems(atir.TItem)) == 0 {
			screen.actionText = loud.Sprintf("You don't have matched items to fulfill this trade.")
			screen.Render()
		} else if atir.IsMyTrdReq {
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
		if screen.user.GetPylonAmount() < sstr.Price {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else if sstr.IsMyTrdReq {
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
		if len(screen.user.GetMatchedCharacters(cbtr.TCharacter)) == 0 {
			screen.actionText = loud.Sprintf("You don't have matched characters to fulfill this trade.")
			screen.Render()
		} else if cbtr.IsMyTrdReq {
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
		if screen.user.GetPylonAmount() < cstr.Price {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else if cstr.IsMyTrdReq {
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
