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

type TextLine struct {
	content string
	font    FontType
}

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
			newSliceLen += 1
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
			newSliceLen += 1
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

func justifyRight(message string, width int) string {
	if NumberOfSpaces(message) < width {
		fmtString := fmt.Sprintf("%%%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	return fillSpaceLeft(ellipsis+SliceFromEnd(message, width-1), width)
}

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
	} else {
		return "Carry: " + formatItemP(item)
	}
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
	case loud.FIRE_SPECIAL:
		return "🔥"
	case loud.ICE_SPECIAL:
		return "🌊"
	case loud.ACID_SPECIAL:
		return "🥗"
	}
	return ""
}

func formatSpecialDragon(special int) string {
	switch special {
	case loud.FIRE_SPECIAL:
		return "Fire dragon"
	case loud.ICE_SPECIAL:
		return "Ice dragon"
	case loud.ACID_SPECIAL:
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
	if ch.Special != loud.NO_SPECIAL {
		chStr = formatSpecial(ch.Special) + " " + chStr // adding space for Sierra issue
	}
	if ch.GiantKill > 0 {
		chStr = fmt.Sprintf("🗿 x%d %s", ch.GiantKill, chStr)
	}
	if ch.SpecialDragonKill > 0 {
		switch ch.Special {
		case loud.FIRE_SPECIAL:
			chStr = fmt.Sprintf("🦐 x%d %s", ch.SpecialDragonKill, chStr)
		case loud.ICE_SPECIAL:
			chStr = fmt.Sprintf("🦈 x%d %s", ch.SpecialDragonKill, chStr)
		case loud.ACID_SPECIAL:
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
	if chs.Special != loud.NO_SPECIAL {
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

func (screen *GameScreen) renderTRLine(text1 string, text2 string, text3 string, font FontType, width int) string {
	text1 = loud.Localize(text1)
	text2 = loud.Localize(text2)
	text3 = loud.Localize(text3)

	calcText := "│" + centerText(text1, " ", 20) + "│" + centerText(text2, " ", 15) + "│" + centerText(text3, " ", 15) + "│"
	onColor := screen.getFont(font)
	return onColor(fillSpace(calcText, width))
}

func (screen *GameScreen) renderItemTableLine(index int, text1 string, font FontType, width int) string {
	text := loud.Localize(text1)
	if index >= 0 {
		text = fmt.Sprintf(" %d) %s", index+1, text)
	}
	calcText := "│" + fillSpace(text, 52) + "│"
	onColor := screen.getFont(font)
	return onColor(fillSpace(calcText, width))
}

func (screen *GameScreen) renderItemTrdReqTableLine(text1 string, text2 string, font FontType, width int) string {
	text1 = loud.Localize(text1)
	text2 = loud.Localize(text2)
	calcText := "│" + centerText(text1, " ", 36) + "│" + centerText(text2, " ", 15) + "│"
	onColor := screen.getFont(font)
	return onColor(fillSpace(calcText, width))
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
