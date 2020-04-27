package screen

import (
	"fmt"
	"io"
	"os"

	loud "github.com/Pylons-tech/LOUD/data"
	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderCharacterSheet() {
	var HP uint64 = 0
	var MaxHP uint64 = 0

	// update blockHeight from newly synced data
	if lbh := screen.user.GetLatestBlockHeight(); lbh > screen.blockHeight {
		screen.blockHeight = lbh
		screen.fakeBlockHeight = lbh
	}

	characters := screen.user.InventoryCharacters()
	activeCharacter := screen.user.GetActiveCharacter()
	activeCharacterRestBlocks := uint64(0)
	if activeCharacter != nil {
		HP = uint64(activeCharacter.HP)
		MaxHP = uint64(activeCharacter.MaxHP)
		activeCharacterRestBlocks = screen.BlockSince(activeCharacter.LastUpdate)
		HP = min(HP+activeCharacterRestBlocks, MaxHP)
	}
	activeWeapon := screen.user.GetActiveWeapon()

	x := screen.rightInnerStartX()
	width := screen.rightInnerWidth()

	infoLines := []string{
		centerText(fmt.Sprintf("%v", screen.user.GetUserName()), " ", width),
		centerText(loud.Localize("inventory"), "─", width),
		truncateRight(fmt.Sprintf("%s %v", screen.goldIcon(), screen.user.GetGold()), width-1),
		"",
	}

	items := screen.user.InventoryItems()
	for _, item := range items {
		itemInfo := truncateRight(formatItem(item), width)
		if activeWeapon != nil && item.ID == activeWeapon.ID {
			itemInfo = screen.blueBoldFont()(itemInfo)
		}
		infoLines = append(infoLines, itemInfo)
	}

	for idx, character := range characters {
		characterInfo := truncateRight(formatCharacter(character), width)
		if idx == screen.user.GetActiveCharacterIndex() {
			characterInfo = screen.blueBoldFont()(characterInfo)
		}
		infoLines = append(infoLines, characterInfo)
	}

	bgcolor := uint64(bgcolor)
	warning := ""
	if float32(HP) < float32(MaxHP)*.25 {
		bgcolor = 124
		warning = loud.Localize("health low warning")
	} else if float32(HP) < float32(MaxHP)*.1 {
		bgcolor = 160
		warning = loud.Localize("health critical warning")
	}
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))

	infoLines = append(infoLines,
		fmtFunc(centerText(fmt.Sprintf(" %s%s", loud.Localize("Active Character"), warning), "─", width)),
		screen.drawProgressMeter(HP, MaxHP, 196, bgcolor, 10)+fmtFunc(truncateRight(fmt.Sprintf(" HP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 225, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" XP: %v/%v", HP, 10), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 208, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" AP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 117, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" RP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 76, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" MP: %v/%v", HP, MaxHP), width-10)),
	)
	if activeCharacter != nil {
		infoLines = append(infoLines,
			fmtFunc(truncateRight(formatCharacter(*activeCharacter), width)),
			fmtFunc(truncateRight(fmt.Sprintf("%s: %d", loud.Localize("rest blocks"), activeCharacterRestBlocks), width)),
		)
	}
	if activeWeapon != nil {
		infoLines = append(infoLines,
			centerText(fmt.Sprintf(" %s ", loud.Localize("Active Weapon")), "─", width),
			fmtFunc(truncateRight(formatItem(*activeWeapon), width)),
		)
	}

	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(2+index, x),
			fillRightWithSpace(line, width)))
		if index+2 > int(screen.Height()) {
			break
		}
	}

	nodeLines := []string{
		centerText(" "+loud.Localize("pylons network status")+" ", "─", width),
		fmt.Sprintf("%s: %s 📋(M)", loud.Localize("Address"), truncateRight(screen.user.GetAddress(), 32)),
		screen.pylonIcon() + truncateRight(fmt.Sprintf(" %s: %v", "Pylon", screen.user.GetPylonAmount()), width-1),
	}

	if len(screen.user.GetLastTxHash()) > 0 {
		nodeLines = append(nodeLines, fmt.Sprintf("%s: %s 📋(L)", loud.Localize("Last TxHash"), truncateRight(screen.user.GetLastTxHash(), 32)))
	}

	blockHeightText := truncateRight(fmt.Sprintf("%s: %d(%d)", loud.Localize("block height"), screen.blockHeight, screen.fakeBlockHeight), width-1)
	if screen.syncingData {
		nodeLines = append(nodeLines, screen.blueBoldFont()(blockHeightText))
	} else {
		nodeLines = append(nodeLines, blockHeightText)
	}
	nodeLines = append(nodeLines, centerText(" ❦ ", "─", width))

	for index, line := range nodeLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(2+len(infoLines)+index, x),
			fillRightWithSpace(line, width)))
		if index+2 > int(screen.Height()) {
			break
		}
	}

	lastLine := len(infoLines) + len(nodeLines) + 1
	screen.drawFill(x, lastLine+1, width, screen.Height()-(lastLine+2))
}
