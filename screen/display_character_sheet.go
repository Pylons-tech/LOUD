package screen

import (
	"fmt"
	"io"
	"os"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderCharacterSheet() {

	// update blockHeight from newly synced data
	if lbh := screen.user.GetLatestBlockHeight(); lbh > screen.blockHeight {
		screen.blockHeight = lbh
		screen.fakeBlockHeight = lbh
	}

	characters := screen.user.InventoryCharacters()
	activeCharacter := screen.user.GetActiveCharacter()
	activeCharacterRestBlocks := uint64(0)
	if activeCharacter != nil {
		activeCharacterRestBlocks = screen.BlockSince(activeCharacter.LastUpdate)
	}
	activeWeapon := screen.user.GetActiveWeapon()

	x := screen.rightInnerStartX()
	width := screen.rightInnerWidth()

	charBkColor := uint64(bgcolor)
	warning := ""

	charFunc := screen.colorFunc(fmt.Sprintf("255:%v", charBkColor))
	fmtFunc := screen.regularFont()

	infoLines := []string{
		fmtFunc(centerText(fmt.Sprintf("%v", screen.user.GetUserName()), " ", width)),
		fmtFunc(centerText(loud.Localize("inventory"), "─", width)),
		fmtFunc(fillSpace(fmt.Sprintf("💰 %v", screen.user.GetGold()), width)),
		fmtFunc(fillSpace("", width)),
	}

	items := screen.user.InventoryItems()
	for _, item := range items {
		itemInfo := fillSpace(formatItem(item), width)
		if activeWeapon != nil && item.ID == activeWeapon.ID {
			itemInfo = screen.blueBoldFont()(itemInfo)
		} else {
			itemInfo = fmtFunc(itemInfo)
		}
		infoLines = append(infoLines, itemInfo)
	}

	for idx, character := range characters {
		characterInfo := fillSpace(formatCharacter(character), width)
		if idx == screen.user.GetActiveCharacterIndex() {
			characterInfo = screen.blueBoldFont()(characterInfo)
		} else {
			characterInfo = fmtFunc(characterInfo)
		}
		infoLines = append(infoLines, characterInfo)
	}

	infoLines = append(infoLines,
		charFunc(centerText(fmt.Sprintf(" %s%s", loud.Localize("Active Character"), warning), "─", width)),
		// screen.drawProgressMeter(HP, MaxHP, 196, bgcolor, 10)+charFunc(fillSpace(fmt.Sprintf(" HP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 225, bgcolor, 10) + charFunc(truncateRight(fmt.Sprintf(" XP: %v/%v", HP, 10), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 208, bgcolor, 10) + charFunc(truncateRight(fmt.Sprintf(" AP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 117, bgcolor, 10) + charFunc(truncateRight(fmt.Sprintf(" RP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 76, bgcolor, 10) + charFunc(truncateRight(fmt.Sprintf(" MP: %v/%v", HP, MaxHP), width-10)),
	)
	if activeCharacter != nil {
		infoLines = append(infoLines,
			charFunc(fillSpace(formatCharacterP(activeCharacter), width)),
			charFunc(fillSpace(fmt.Sprintf("%s: %d", loud.Localize("rest blocks"), activeCharacterRestBlocks), width)),
		)
	}
	if activeWeapon != nil {
		infoLines = append(infoLines,
			fmtFunc(centerText(fmt.Sprintf(" %s ", loud.Localize("Active Weapon")), "─", width)),
			fmtFunc(fillSpace(formatItemP(activeWeapon), width)),
		)
	}

	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(2+index, x),
			line))
		if index+2 > int(screen.Height()) {
			break
		}
	}

	nodeLines := []string{
		fmtFunc(centerText(" "+loud.Localize("pylons network status")+" ", "─", width)),
		fmtFunc(fillSpace(fmt.Sprintf("%s: %s 📋(M)", loud.Localize("Address"), truncateRight(screen.user.GetAddress(), 32)), width)),
		fmtFunc(fillSpace(fmt.Sprintf("%s %s: %v", screen.pylonIcon(), "Pylon", screen.user.GetPylonAmount()), width)),
	}

	if len(screen.user.GetLastTxHash()) > 0 {
		txHashT := fmt.Sprintf("%s: %s 📋(L)", loud.Localize("Last TxHash"), truncateRight(screen.user.GetLastTxHash(), 32))
		nodeLines = append(nodeLines, fmtFunc(fillSpace(txHashT, width)))
	}

	blockHeightText := fillSpace(fmt.Sprintf("%s ⟳ (E): %d(%d)", loud.Localize("block height"), screen.blockHeight, screen.fakeBlockHeight), width)
	if screen.syncingData {
		nodeLines = append(nodeLines, screen.blueBoldFont()(blockHeightText))
	} else {
		nodeLines = append(nodeLines, fmtFunc(blockHeightText))
	}
	nodeLines = append(nodeLines, fmtFunc(centerText(" ❦ ", "─", width)))

	for index, line := range nodeLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(2+len(infoLines)+index, x),
			line))
		if index+2 > int(screen.Height()) {
			break
		}
	}

	lastLine := len(infoLines) + len(nodeLines) + 1
	screen.drawFill(x, lastLine+1, width, screen.Height()-(lastLine+2))
}
