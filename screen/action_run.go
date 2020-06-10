package screen

import (
	"time"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/LOUD/log"
)

// RunTxProcess execute the screen status changes when running transaction
func (screen *GameScreen) RunTxProcess(waitStatus PageStatus, resultStatus PageStatus, fn func() (string, error)) {
	screen.SetScreenStatusAndRefresh(waitStatus)

	log.Debugln("started sending request for ", waitStatus)
	go func() {
		txhash, err := fn()
		log.Debugln("ended sending request for ", waitStatus)
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

// RunActiveCharacterSelect execute the active character selection
func (screen *GameScreen) RunActiveCharacterSelect(index int) {
	screen.user.SetActiveCharacterIndex(index)
	screen.SetScreenStatusAndRefresh(RsltSelectActiveChr)
}

// RunCharacterRename execute the character rename process
func (screen *GameScreen) RunCharacterRename(newName string) {
	screen.RunTxProcess(WaitRenameChr, RsltRenameChr, func() (string, error) {
		return loud.RenameCharacter(screen.user, screen.activeCharacter, newName)
	})
}

// RunActiveItemBuy execute the item buying process
func (screen *GameScreen) RunActiveItemBuy() {
	if !screen.user.HasPreItemForAnItem(screen.activeItem) {
		screen.txFailReason = loud.Sprintf("You don't have required item to make %s", screen.activeItem.Name)
		screen.SetScreenStatusAndRefresh(RsltBuyItem)
		return
	}
	screen.RunTxProcess(WaitBuyItem, RsltBuyItem, func() (string, error) {
		return loud.BuyItem(screen.user, screen.activeItem)
	})
}

// RunActiveCharacterBuy execute the character buying process
func (screen *GameScreen) RunActiveCharacterBuy() {
	screen.RunTxProcess(WaitBuyChr, RsltBuyChr, func() (string, error) {
		return loud.BuyCharacter(screen.user, screen.activeCharacter)
	})
}

// RunActiveItemSell execute the item sell process
func (screen *GameScreen) RunActiveItemSell() {
	screen.RunTxProcess(WaitSellItem, RsltSellItem, func() (string, error) {
		return loud.SellItem(screen.user, screen.activeItem)
	})
}

// RunActiveItemUpgrade execute the item upgrade process
func (screen *GameScreen) RunActiveItemUpgrade() {
	screen.RunTxProcess(WaitUpgradeItem, RsltUpgradeItem, func() (string, error) {
		return loud.UpgradeItem(screen.user, screen.activeItem)
	})
}

// RunHuntRabbits execute the hunt rabbit process
func (screen *GameScreen) RunHuntRabbits() {
	screen.RunTxProcess(WaitHuntRabbits, RsltHuntRabbits, func() (string, error) {
		return loud.HuntRabbits(screen.user)
	})
}

// RunFightGiant execute the giant fight process
func (screen *GameScreen) RunFightGiant() {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight giant without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(WaitFightGiant, RsltFightGiant, func() (string, error) {
		return loud.FightGiant(screen.user)
	})
}

// RunFightDragonFire execute the fight fire dragon process
func (screen *GameScreen) RunFightDragonFire() {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight fire dragon without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(WaitFightDragonFire, RsltFightDragonFire, func() (string, error) {
		return loud.FightDragonFire(screen.user)
	})
}

// RunFightDragonIce execute the fight ice dragon process
func (screen *GameScreen) RunFightDragonIce() {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight ice dragon without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(WaitFightDragonIce, RsltFightDragonIce, func() (string, error) {
		return loud.FightDragonIce(screen.user)
	})
}

// RunFightDragonAcid execute the fight acid dragon process
func (screen *GameScreen) RunFightDragonAcid() {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight acid dragon without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(WaitFightDragonAcid, RsltFightDragonAcid, func() (string, error) {
		return loud.FightDragonAcid(screen.user)
	})
}

// RunFightDragonUndead execute the fight undead dragon process
func (screen *GameScreen) RunFightDragonUndead() {
	if len(screen.user.InventoryAngelSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight undead dragon without angel sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(WaitFightDragonUndead, RsltFightDragonUndead, func() (string, error) {
		return loud.FightDragonUndead(screen.user)
	})
}

// RunFightTroll execute the fight troll process
func (screen *GameScreen) RunFightTroll() {
	screen.RunTxProcess(WaitFightTroll, RsltFightTroll, func() (string, error) {
		return loud.FightTroll(screen.user)
	})
}

// RunFightWolf execute the fight wolf process
func (screen *GameScreen) RunFightWolf() {
	screen.RunTxProcess(WaitFightWolf, RsltFightWolf, func() (string, error) {
		return loud.FightWolf(screen.user)
	})
}

// RunFightGoblin execute the fight goblin process
func (screen *GameScreen) RunFightGoblin() {
	screen.RunTxProcess(WaitFightGoblin, RsltFightGoblin, func() (string, error) {
		return loud.FightGoblin(screen.user)
	})
}

// RunSelectedLoudBuyTrdReq execute the gold buy trading process
func (screen *GameScreen) RunSelectedLoudBuyTrdReq() {
	if len(loud.BuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		// when activeLine is not refering to real request but when it is refering to nil request
		screen.txFailReason = loud.Localize("you haven't selected any buy request")
		screen.SetScreenStatusAndRefresh(RsltFulfillBuyGoldTrdReq)
	} else {
		screen.activeTrdReq = loud.BuyTrdReqs[screen.activeLine]
		if screen.user.GetGold() < screen.activeTrdReq.Amount {
			screen.actionText = loud.Sprintf("You don't have enough gold to fulfill this trade.")
			screen.Render()
		} else if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelTrdReq, RsltCancelTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else {
			screen.RunTxProcess(WaitFulfillBuyGoldTrdReq, RsltFulfillBuyGoldTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID)
			})
		}
	}
}

// RunSelectedLoudSellTrdReq execute the gold sell trading process
func (screen *GameScreen) RunSelectedLoudSellTrdReq() {
	if len(loud.SellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell request")
		screen.SetScreenStatusAndRefresh(RsltFulfillSellGoldTrdReq)
	} else {
		screen.activeTrdReq = loud.SellTrdReqs[screen.activeLine]
		if screen.user.GetPylonAmount() < screen.activeTrdReq.Total {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelTrdReq, RsltCancelTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else {
			screen.RunTxProcess(WaitFulfillSellGoldTrdReq, RsltFulfillSellGoldTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID)
			})
		}
	}
}

// RunSelectedItemBuyTrdReq execute the item buy trading process
func (screen *GameScreen) RunSelectedItemBuyTrdReq() {
	if len(loud.ItemBuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy item request")
		screen.SetScreenStatusAndRefresh(RsltFulfillBuyItemTrdReq)
	} else {
		atir := loud.ItemBuyTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = atir
		if len(screen.user.GetMatchedItems(atir.TItem)) == 0 {
			screen.actionText = loud.Sprintf("You don't have matched items to fulfill this trade.")
			screen.Render()
		} else if atir.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelTrdReq, RsltCancelTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, atir.ID)
			})
		} else {
			screen.RunTxProcess(WaitFulfillBuyItemTrdReq, RsltFulfillBuyItemTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, atir.ID)
			})
		}
	}
}

// RunSelectedItemSellTrdReq execute the item sell trading process
func (screen *GameScreen) RunSelectedItemSellTrdReq() {
	if len(loud.ItemSellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell item request")
		screen.SetScreenStatusAndRefresh(RsltFulfillSellItemTrdReq)
	} else {
		sstr := loud.ItemSellTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = sstr
		if screen.user.GetPylonAmount() < sstr.Price {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else if sstr.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelTrdReq, RsltCancelTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, sstr.ID)
			})
		} else {
			screen.RunTxProcess(WaitFulfillSellItemTrdReq, RsltFulfillSellItemTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, sstr.ID)
			})
		}
	}
}

// RunSelectedCharacterBuyTrdReq execute the character buy trading process
func (screen *GameScreen) RunSelectedCharacterBuyTrdReq() {
	if len(loud.CharacterBuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any buy character request")
		screen.SetScreenStatusAndRefresh(RsltFulfillBuyChrTrdReq)
	} else {
		cbtr := loud.CharacterBuyTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = cbtr
		if len(screen.user.GetMatchedCharacters(cbtr.TCharacter)) == 0 {
			screen.actionText = loud.Sprintf("You don't have matched characters to fulfill this trade.")
			screen.Render()
		} else if cbtr.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelTrdReq, RsltCancelTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, cbtr.ID)
			})
		} else {
			screen.RunTxProcess(WaitFulfillBuyChrTrdReq, RsltFulfillBuyChrTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, cbtr.ID)
			})
		}
	}
}

// RunSelectedCharacterSellTrdReq execute the character sell trading process
func (screen *GameScreen) RunSelectedCharacterSellTrdReq() {
	if len(loud.CharacterSellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell character request")
		screen.SetScreenStatusAndRefresh(RsltFulfillSellChrTrdReq)
	} else {
		cstr := loud.CharacterSellTrdReqs[screen.activeLine]
		screen.activeItemTrdReq = cstr
		if screen.user.GetPylonAmount() < cstr.Price {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else if cstr.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelTrdReq, RsltCancelTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, cstr.ID)
			})
		} else {
			screen.RunTxProcess(WaitFulfillSellChrTrdReq, RsltFulfillSellChrTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, cstr.ID)
			})
		}
	}
}
