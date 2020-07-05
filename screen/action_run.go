package screen

import (
	"time"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/Pylons-tech/LOUD/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RunTxProcess execute the screen status changes when running transaction
func (screen *GameScreen) RunTxProcess(waitStatus PageStatus, resultStatus PageStatus, fn func() (string, error)) {
	log.WithFields(log.Fields{
		"func_start":    "RunTxProcess",
		"wait_status":   waitStatus,
		"result_status": resultStatus,
	}).Debugln("debug log")
	screen.SetScreenStatusAndRefresh(waitStatus)

	go func() {
		txhash, err := fn()
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

// RunActiveFriendRemove execute the friend remove action
func (screen *GameScreen) RunActiveFriendRemove(index int) {
	friends := screen.user.Friends()
	friends = append(friends[:index], friends[index+1:]...)
	screen.user.SetFriends(friends)
	screen.SetScreenStatusAndRefresh(RsltFriendRemove)
}

// RunFriendRegister execute the friend register action
func (screen *GameScreen) RunFriendRegister() {
	_, err := sdk.AccAddressFromBech32(screen.friendAddress)

	if err != nil {
		log.Println("Invalid friend address", err.Error())
		screen.actionText = loud.Sprintf("Invalid friend address \"%s\"", screen.friendAddress)
		screen.Render()
	} else {
		friends := screen.user.Friends()
		friends = append(friends, loud.Friend{
			Name:    screen.friendNameValue,
			Address: screen.friendAddress,
		})
		screen.user.SetFriends(friends)
		screen.SetScreenStatusAndRefresh(RsltFriendRegister)
	}
}

// RunCharacterRename execute the character rename process
func (screen *GameScreen) RunCharacterRename(newName string) {
	screen.RunTxProcess(WaitRenameChr, RsltRenameChr, func() (string, error) {
		return loud.RenameCharacter(screen.user, screen.activeCharacter, newName)
	})
}

// RunSendItem execute the process to send item
func (screen *GameScreen) RunSendItem() {
	screen.RunTxProcess(WaitSendItem, RsltSendItem, func() (string, error) {
		return loud.SendItem(screen.user, screen.activeFriend, screen.activeItem)
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
func (screen *GameScreen) RunFightGiant(tarBonus int) {
	if len(screen.user.InventoryIronSwords()) == 0 {
		screen.actionText = loud.Sprintf("You can't fight giant without iron sword.")
		screen.Render()
		return
	}
	screen.RunTxProcess(WaitFightGiant, RsltFightGiant, func() (string, error) {
		return loud.FightGiant(screen.user, tarBonus)
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

// RunSelectedBuyGoldTrdReq execute the gold buy trading process
func (screen *GameScreen) RunSelectedBuyGoldTrdReq() {
	if len(loud.BuyTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		// when activeLine is not refering to real request but when it is refering to nil request
		screen.txFailReason = loud.Localize("you haven't selected any buy request")
		screen.SetScreenStatusAndRefresh(RsltFulfillBuyGoldTrdReq)
	} else {
		screen.activeTrdReq = loud.BuyTrdReqs[screen.activeLine]
		if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelBuyGoldTrdReq, RsltCancelBuyGoldTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else if screen.user.GetGold() < screen.activeTrdReq.Amount {
			screen.actionText = loud.Sprintf("You don't have enough gold to fulfill this trade.")
			screen.Render()
		} else {
			screen.RunTxProcess(WaitFulfillBuyGoldTrdReq, RsltFulfillBuyGoldTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID, []string{})
			})
		}
	}
}

// RunSelectedSellGoldTrdReq execute the gold sell trading process
func (screen *GameScreen) RunSelectedSellGoldTrdReq() {
	if len(loud.SellTrdReqs) <= screen.activeLine || screen.activeLine < 0 {
		screen.txFailReason = loud.Localize("you haven't selected any sell request")
		screen.SetScreenStatusAndRefresh(RsltFulfillSellGoldTrdReq)
	} else {
		screen.activeTrdReq = loud.SellTrdReqs[screen.activeLine]
		if screen.activeTrdReq.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelSellGoldTrdReq, RsltCancelSellGoldTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, screen.activeTrdReq.ID)
			})
		} else if screen.user.GetPylonAmount() < screen.activeTrdReq.Total {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else {
			screen.RunTxProcess(WaitFulfillSellGoldTrdReq, RsltFulfillSellGoldTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, screen.activeTrdReq.ID, []string{})
			})
		}
	}
}

// RunSelectedItemBuyTrdReq execute the item buy trading process
func (screen *GameScreen) RunSelectedItemBuyTrdReq() {
	atir := screen.activeItemTrdReq.(loud.ItemBuyTrdReq)

	if atir.IsMyTrdReq {
		screen.RunTxProcess(WaitCancelBuyItemTrdReq, RsltCancelBuyItemTrdReq, func() (string, error) {
			return loud.CancelTrade(screen.user, atir.ID)
		})
	} else {
		matchingItems := screen.user.GetMatchedItems(atir.TItem)
		if len(matchingItems) <= screen.activeLine || screen.activeLine < 0 {
			screen.txFailReason = loud.Localize("You haven't selected any matched item for the request")
			screen.SetScreenStatusAndRefresh(RsltFulfillBuyItemTrdReq)
		} else {
			itemIDs := []string{matchingItems[screen.activeLine].ID}
			screen.RunTxProcess(WaitFulfillBuyItemTrdReq, RsltFulfillBuyItemTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, atir.ID, itemIDs)
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
		if sstr.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelSellItemTrdReq, RsltCancelSellItemTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, sstr.ID)
			})
		} else if screen.user.GetPylonAmount() < sstr.Price {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else {
			screen.RunTxProcess(WaitFulfillSellItemTrdReq, RsltFulfillSellItemTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, sstr.ID, []string{})
			})
		}
	}
}

// RunSelectedCharacterBuyTrdReq execute the character buy trading process
func (screen *GameScreen) RunSelectedCharacterBuyTrdReq() {
	cbtr := screen.activeItemTrdReq.(loud.CharacterBuyTrdReq)

	if cbtr.IsMyTrdReq {
		screen.RunTxProcess(WaitCancelBuyChrTrdReq, RsltCancelBuyChrTrdReq, func() (string, error) {
			return loud.CancelTrade(screen.user, cbtr.ID)
		})
	} else {
		matchingChrs := screen.user.GetMatchedCharacters(cbtr.TCharacter)
		if len(matchingChrs) <= screen.activeLine || screen.activeLine < 0 {
			screen.txFailReason = loud.Localize("You haven't selected any matched characters for the request")
			screen.SetScreenStatusAndRefresh(RsltFulfillBuyChrTrdReq)
		} else {
			itemIDs := []string{matchingChrs[screen.activeLine].ID}
			screen.RunTxProcess(WaitFulfillBuyChrTrdReq, RsltFulfillBuyChrTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, cbtr.ID, itemIDs)
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
		if cstr.IsMyTrdReq {
			screen.RunTxProcess(WaitCancelSellChrTrdReq, RsltCancelSellChrTrdReq, func() (string, error) {
				return loud.CancelTrade(screen.user, cstr.ID)
			})
		} else if screen.user.GetPylonAmount() < cstr.Price {
			screen.actionText = loud.Sprintf("You don't have enough pylons to fulfill this trade.")
			screen.Render()
		} else {
			screen.RunTxProcess(WaitFulfillSellChrTrdReq, RsltFulfillSellChrTrdReq, func() (string, error) {
				return loud.FulfillTrade(screen.user, cstr.ID, []string{})
			})
		}
	}
}
