package screen

import (
	"fmt"
	"io"
	"os"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderCharacterSheet() {

	scrBox := screen.GetCharacterSheetBox()
	// situation box start point (x, y)
	x := scrBox.X
	y := scrBox.Y
	w := scrBox.W
	h := scrBox.H

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

	charBkColor := uint64(bgcolor)
	warning := ""

	charFunc := screen.colorFunc(fmt.Sprintf("255:%v", charBkColor))
	fmtFunc := screen.regularFont()

	infoLines := []string{
		fmtFunc(centerText(fmt.Sprintf("%v", screen.user.GetUserName()), " ", w)),
		fmtFunc(centerText(loud.Localize("inventory"), "─", w)),
		fmtFunc(fillSpace(fmt.Sprintf("💰 %v", screen.user.GetGold()), w)),
		fmtFunc(fillSpace("", w)),
	}

	MAX_INVENTORY_LEN := h - 15

	for idx, character := range characters {
		if len(infoLines) > MAX_INVENTORY_LEN {
			infoLines = append(infoLines, fmtFunc(fillSpace("...", w)))
			break
		}
		characterInfo := fillSpace(formatCharacter(character), w)
		if idx == screen.user.GetActiveCharacterIndex() {
			characterInfo = screen.blueBoldFont()(characterInfo)
		} else {
			characterInfo = fmtFunc(characterInfo)
		}
		infoLines = append(infoLines, characterInfo)
	}

	items := screen.user.InventoryItems()
	for _, item := range items {
		if len(infoLines) > MAX_INVENTORY_LEN {
			infoLines = append(infoLines, fmtFunc(fillSpace("...", w)))
			break
		}
		itemInfo := fillSpace(formatItem(item), w)
		infoLines = append(infoLines, fmtFunc(itemInfo))
	}

	// HP := uint64(100)
	// MaxHP := uint64(100)
	infoLines = append(infoLines,
		charFunc(centerText(fmt.Sprintf(" %s%s", loud.Localize("Active Character"), warning), "─", w)),
		// screen.drawProgressMeter(HP, MaxHP, 196, bgcolor, 10)+charFunc(fillSpace(fmt.Sprintf(" HP: %v/%v", HP, MaxHP), w-10)),
		// screen.drawProgressMeter(HP, MaxHP, 225, bgcolor, 10)+charFunc(truncateRight(fmt.Sprintf(" XP: %v/%v", HP, 10), w-10)),
		// screen.drawProgressMeter(HP, MaxHP, 208, bgcolor, 10)+charFunc(truncateRight(fmt.Sprintf(" AP: %v/%v", HP, MaxHP), w-10)),
		// screen.drawProgressMeter(HP, MaxHP, 117, bgcolor, 10)+charFunc(truncateRight(fmt.Sprintf(" RP: %v/%v", HP, MaxHP), w-10)),
		// screen.drawProgressMeter(HP, MaxHP, 76, bgcolor, 10)+charFunc(truncateRight(fmt.Sprintf(" MP: %v/%v", HP, MaxHP), w-10)),
	)
	if activeCharacter != nil {
		infoLines = append(infoLines,
			charFunc(fillSpace(formatCharacterP(activeCharacter), w)),
			charFunc(fillSpace(fmt.Sprintf("%s: %d", loud.Localize("rest blocks"), activeCharacterRestBlocks), w)),
		)
	}

	lenInfoLines := len(infoLines)

	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x),
			line))
	}

	nodeLines := []string{
		fmtFunc(centerText(" "+loud.Localize("pylons network status")+" ", "─", w)),
		fmtFunc(fillSpace(fmt.Sprintf("%s: %s 📋 (M)", loud.Localize("Address"), truncateRight(screen.user.GetAddress(), 15)), w)),
		fmtFunc(fillSpace(fmt.Sprintf("%s %s: %v", screen.pylonIcon(), "Pylon", screen.user.GetPylonAmount()), w)),
	}

	if len(screen.user.GetLastTxHash()) > 0 {
		txHashT := fmt.Sprintf("%s: %s 📋 (L)", loud.Localize("Last TxHash"), truncateRight(screen.user.GetLastTxHash(), 15))
		nodeLines = append(nodeLines, fmtFunc(fillSpace(txHashT, w)))
	}

	blockHeightText := fillSpace(fmt.Sprintf("%s ⟳ (E): %d(%d)", loud.Localize("block height"), screen.blockHeight, screen.fakeBlockHeight), w)
	if screen.syncingData {
		nodeLines = append(nodeLines, screen.blueBoldFont()(blockHeightText))
	} else {
		nodeLines = append(nodeLines, fmtFunc(blockHeightText))
	}
	nodeLines = append(nodeLines, fmtFunc(centerText(" ❦ ", "─", w)))

	for index, line := range nodeLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+lenInfoLines+index, x),
			line))
	}

	lastLine := y + lenInfoLines + len(nodeLines) + 1
	screen.drawFill(x, lastLine+1, w, h-lastLine)
}
