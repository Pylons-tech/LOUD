package screen

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"os"

	"github.com/ahmetb/go-cursor"
	"github.com/mgutz/ansi"

	loud "github.com/Pylons-tech/LOUD/data"
)

// TextLine is a struct to manage a line on screen
type TextLine struct {
	content string
	font    FontType
}

// TextLines is a struct to manage bulk TextLine objects
type TextLines []TextLine

func (tl TextLines) append(elems ...string) TextLines {
	elemsT := []TextLine{}
	for _, el := range elems {
		elemsT = append(elemsT, TextLine{
			content: el,
			font:    "",
		})
	}
	return append(tl, elemsT...)
}

func (tl TextLines) appendF(elem string, font FontType) TextLines {
	return append(tl, TextLine{
		content: elem,
		font:    font,
	})
}

func (tl TextLines) appendT(elems ...string) TextLines {
	elemsT := []TextLine{}
	for _, el := range elems {
		elemsT = append(elemsT, TextLine{
			content: loud.Localize(el),
			font:    "",
		})
	}
	return append(tl, elemsT...)
}

// SliceFromStart cut the text from the start, max width is provided
func SliceFromStart(text string, width int) string {
	sliceLen := 0
	for {
		newSliceLen := sliceLen
		startWithCustomUnicode := false
		for k, v := range customUnicodes {
			if strings.HasPrefix(text[sliceLen:len(text)], k) {
				newSliceLen += len(v)
				startWithCustomUnicode = true
				break
			}
		}
		if !startWithCustomUnicode {
			newSliceLen++
		}
		if newSliceLen <= width {
			sliceLen = newSliceLen
		}
		if newSliceLen >= width || newSliceLen >= len(text) {
			break
		}
	}
	return text[:sliceLen]
}

// SliceFromEnd cut the text from the end, max width is provided
func SliceFromEnd(text string, width int) string {
	sliceLen := 0
	for {
		newSliceLen := sliceLen
		endWithCustomUnicode := false
		for k, v := range customUnicodes {
			if strings.HasSuffix(text[:len(text)-sliceLen], k) {
				newSliceLen += len(v)
				endWithCustomUnicode = true
				break
			}
		}
		if !endWithCustomUnicode {
			newSliceLen++
		}
		if newSliceLen <= width {
			sliceLen = newSliceLen
		}
		if newSliceLen >= width || newSliceLen >= len(text) {
			break
		}
	}
	return text[len(text)-sliceLen : len(text)]
}

func truncateRight(message string, width int) string {
	if NumberOfSpaces(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	return fillSpace(SliceFromStart(message, width-1)+ellipsis, width)
}

func truncateLeft(message string, width int) string {
	if NumberOfSpaces(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	return fillSpaceLeft(ellipsis+SliceFromEnd(message, width-1), width)
}

// func justifyRight(message string, width int) string {
// 	if NumberOfSpaces(message) < width {
// 		fmtString := fmt.Sprintf("%%%vs", width)

// 		return fmt.Sprintf(fmtString, message)
// 	}
// 	return fillSpaceLeft(ellipsis+SliceFromEnd(message, width-1), width)
// }

func centerText(message, pad string, width int) string {
	if NumberOfSpaces(message) > width {
		return fillSpace(SliceFromStart(message, width-1)+ellipsis, width)
	}
	leftover := width - NumberOfSpaces(message)
	left := leftover / 2
	right := leftover - left

	if pad == "" {
		pad = " "
	}

	leftString := ""
	for NumberOfSpaces(leftString) <= left && NumberOfSpaces(leftString) <= right {
		leftString += pad
	}

	return fmt.Sprintf("%s%s%s", string([]rune(leftString)[0:left]), message, string([]rune(leftString)[0:right]))
}

func fillSpaceLeft(message string, width int) string {
	msgLen := NumberOfSpaces(message)
	if msgLen > width {
		return fillSpace(SliceFromStart(message, width-1)+ellipsis, width)
	}
	leftover := width - msgLen

	fillString := ""
	fillLen := 0
	for fillLen < leftover {
		fillString += " "
		fillLen = NumberOfSpaces(fillString)
	}
	return fillString + message
}

func fillSpace(message string, width int) string {
	msgLen := NumberOfSpaces(message)
	if msgLen > width {
		return fillSpace(SliceFromStart(message, width-1)+ellipsis, width)
	}
	leftover := width - msgLen

	fillString := ""
	fillLen := 0
	for fillLen < leftover {
		fillString += " "
		fillLen = NumberOfSpaces(fillString)
	}
	return message + fillString
}

func drawVerticalLine(x, y, height int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))
	for i := 1; i < height; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x), color))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┬", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┴", cursor.MoveTo(y+height, x), color))
}

func drawHorizontalLine(x, y, width int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))
	for i := 1; i < width; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s─", cursor.MoveTo(y, x+i), color))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s├", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┤", cursor.MoveTo(y, x+width), color))
}

func formatItem(item loud.Item) string {
	itemStr := item.Name
	if item.Level > 0 {
		itemStr += fmt.Sprintf(" Lv%d", item.Level)
	}
	// if item.Attack > 0 {
	// 	itemStr += fmt.Sprintf(" attack=%d", item.Attack)
	// }
	return itemStr
}

func formatItemP(item *loud.Item) string {
	if item == nil {
		return ""
	}
	return formatItem(*item)
}

func carryItemDesc(item *loud.Item) string {
	if item == nil {
		return ""
	}
	return "Carry: " + formatItemP(item)
}

func formatIntRange(r [2]int) string {
	if r[0] == r[1] {
		if r[0] == 0 {
			return ""
		}
		return fmt.Sprintf("%d", r[0])
	}
	return fmt.Sprintf("%d-%d", r[0], r[1])
}

func formatFloat64Range(r [2]float64) string {
	if r[0] == r[1] {
		if r[0] == 0 {
			return ""
		}
		return fmt.Sprintf("%.0f", r[0])
	}
	return fmt.Sprintf("%.0f-%.0f", r[0], r[1])
}

func formatItemSpec(itemSpec loud.ItemSpec) string {
	itemStr := itemSpec.Name
	lvlStr := formatIntRange(itemSpec.Level)
	if len(lvlStr) > 0 {
		itemStr += fmt.Sprintf(" Lv%s", lvlStr)
	}
	// attackStr := formatIntRange(itemSpec.Attack)
	// if len(attackStr) > 0 {
	// 	itemStr += fmt.Sprintf(" attack=%s", attackStr)
	// }
	return itemStr
}

func formatSpecial(special int) string {
	switch special {
	case loud.FireSpecial:
		return "🔥"
	case loud.IceSpecial:
		return "🌊"
	case loud.AcidSpecial:
		return "🥗"
	}
	return ""
}

func formatSpecialDragon(special int) string {
	switch special {
	case loud.FireSpecial:
		return "Fire dragon"
	case loud.IceSpecial:
		return "Ice dragon"
	case loud.AcidSpecial:
		return "Acid dragon"
	}
	return ""
}

func formatBigNumber(number int) string {
	if number > 1000000 {
		return fmt.Sprintf("%dM", number/1000000)
	}
	if number > 1000 {
		return fmt.Sprintf("%dk", number/1000)
	}
	return fmt.Sprintf("%d", number)
}

func formatCharacter(ch loud.Character) string {
	chStr := loud.Localize(ch.Name)
	if ch.Special != loud.NoSpecial {
		chStr = formatSpecial(ch.Special) + " " + chStr // adding space for Sierra issue
	}
	if ch.GiantKill > 0 {
		chStr = fmt.Sprintf("🗿 x%d %s", ch.GiantKill, chStr)
	}
	if ch.SpecialDragonKill > 0 {
		switch ch.Special {
		case loud.FireSpecial:
			chStr = fmt.Sprintf("🦐 x%d %s", ch.SpecialDragonKill, chStr)
		case loud.IceSpecial:
			chStr = fmt.Sprintf("🦈 x%d %s", ch.SpecialDragonKill, chStr)
		case loud.AcidSpecial:
			chStr = fmt.Sprintf("🐊 x%d %s", ch.SpecialDragonKill, chStr)
		}
	}
	if ch.UndeadDragonKill > 0 {
		chStr = fmt.Sprintf("🐉 x%d %s", ch.UndeadDragonKill, chStr)
	}
	if ch.Level > 0 {
		chStr += fmt.Sprintf(" Lv%d", ch.Level)
	}
	if ch.XP > 0 {
		chStr += fmt.Sprintf(" XP=%s", formatBigNumber(int(ch.XP)))
	}
	return chStr
}

func formatCharacterP(ch *loud.Character) string {
	if ch == nil {
		return ""
	}
	return formatCharacter(*ch)
}

func formatCharacterSpec(chs loud.CharacterSpec) string {
	chStr := loud.Localize(chs.Name)
	if chs.Special != loud.NoSpecial {
		chStr = formatSpecial(chs.Special) + " " + chStr // adding space for Sierra issue
	}

	lvlStr := formatIntRange(chs.Level)
	if len(lvlStr) > 0 {
		chStr += fmt.Sprintf(" Lv%s", lvlStr)
	}
	xpStr := formatFloat64Range(chs.XP)
	if len(xpStr) > 0 {
		chStr += fmt.Sprintf(" XP=%s", xpStr)
	}
	return chStr
}

// InterfaceSlice returns array from interface object
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func formatByStructType(item interface{}) string {
	switch item.(type) {
	case loud.Item:
		return formatItem(item.(loud.Item))
	case loud.Character:
		return formatCharacter(item.(loud.Character))
	case loud.ItemSpec:
		return formatItemSpec(item.(loud.ItemSpec))
	case loud.CharacterSpec:
		return formatCharacterSpec(item.(loud.CharacterSpec))
	default:
		return "unrecognized struct type"
	}
}

func (screen *GameScreen) renderTRLine(text1 string, text2 string, text3 string) string {
	text1 = loud.Localize(text1)
	text2 = loud.Localize(text2)
	text3 = loud.Localize(text3)

	calcText := "│" + centerText(text1, " ", 20) + "│" + centerText(text2, " ", 15) + "│" + centerText(text3, " ", 15) + "│"
	return calcText
}

func (screen *GameScreen) renderItemTableLine(index int, text1 string) string {
	text := loud.Localize(text1)
	if index >= 0 {
		text = fmt.Sprintf(" %d) %s", index+1, text)
	}
	calcText := "│" + fillSpace(text, 52) + "│"
	return calcText
}

func (screen *GameScreen) renderItemTrdReqTableLine(text1 string, text2 string) string {
	text1 = loud.Localize(text1)
	text2 = loud.Localize(text2)
	calcText := "│" + centerText(text1, " ", 36) + "│" + centerText(text2, " ", 15) + "│"
	return calcText
}

// func min(a, b uint64) uint64 {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// RequestInfo returns basic info from several type of requests
func RequestInfo(request interface{}) (bool, interface{}, int) {
	switch request.(type) {
	case loud.ItemBuyTrdReq:
		itr := request.(loud.ItemBuyTrdReq)
		return itr.IsMyTrdReq, itr.TItem, itr.Price
	case loud.ItemSellTrdReq:
		itr := request.(loud.ItemSellTrdReq)
		return itr.IsMyTrdReq, itr.TItem, itr.Price
	case loud.CharacterBuyTrdReq:
		itr := request.(loud.CharacterBuyTrdReq)
		return itr.IsMyTrdReq, itr.TCharacter, itr.Price
	case loud.CharacterSellTrdReq:
		itr := request.(loud.CharacterSellTrdReq)
		return itr.IsMyTrdReq, itr.TCharacter, itr.Price
	}
	return false, loud.ItemBuyTrdReq{}, 0
}

// ForestStatusCheck checks if a player can fight target monster
func (screen *GameScreen) ForestStatusCheck(newStus PageStatus) (string, string) {
	activeCharacter := screen.user.GetActiveCharacter()
	if activeCharacter == nil {
		return loud.Sprintf("You need a character for this action!"), loud.Sprintf("no character!")
	}
	switch newStus {
	case ConfirmFightGiant:
		if activeCharacter == nil || activeCharacter.Special != loud.NoSpecial {
			return loud.Sprintf("You need no special character for this action!"), loud.Sprintf("no non-special character!")
		}
	case ConfirmFightDragonFire:
		if activeCharacter == nil || activeCharacter.Special != loud.FireSpecial {
			return loud.Sprintf("You need a fire character for this action!"), loud.Sprintf("no fire character!")
		}
	case ConfirmFightDragonIce:
		if activeCharacter == nil || activeCharacter.Special != loud.IceSpecial {
			return loud.Sprintf("You need a ice character for this action!"), loud.Sprintf("no ice character!")
		}
	case ConfirmFightDragonAcid:
		if activeCharacter == nil || activeCharacter.Special != loud.AcidSpecial {
			return loud.Sprintf("You need a acid character for this action!"), loud.Sprintf("no acid character!")
		}
	}
	switch newStus {
	case ConfirmFightGoblin,
		ConfirmFightWolf,
		ConfirmFightTroll:
		if len(screen.user.InventorySwords()) == 0 {
			return loud.Sprintf("You need a sword for this action!"), loud.Sprintf("no sword!")
		}
	case ConfirmFightGiant,
		ConfirmFightDragonFire,
		ConfirmFightDragonIce,
		ConfirmFightDragonAcid:
		if len(screen.user.InventoryIronSwords()) == 0 {
			return loud.Sprintf("You need an iron sword for this action!"), loud.Sprintf("no iron sword!")
		}
	case ConfirmFightDragonUndead:
		if len(screen.user.InventoryAngelSwords()) == 0 {
			return loud.Sprintf("You need an angel sword for this action!"), loud.Sprintf("no angel sword!")
		}
	}
	return "", ""
}
