package screen

import (
	"fmt"
	"io"
	"reflect"
	"unicode/utf8"

	"os"

	"github.com/ahmetb/go-cursor"
	"github.com/mgutz/ansi"

	loud "github.com/Pylons-tech/LOUD/data"
)

type TextLines []string

func (tl TextLines) append(elems ...string) TextLines {
	return append(tl, elems...)
}

func (tl TextLines) appendT(elems ...string) TextLines {
	elemsT := []string{}
	for _, el := range elems {
		elemsT = append(elemsT, loud.Localize(el))
	}
	return append(tl, elemsT...)
}

func NumberOfLetters(message string) int {
	// customUnicodes := map[string]int{
	// 	"💰": 1,
	//  "🔶": 1,
	//  "🔷": 1,
	// 	"🥺": 1,
	// 	"🗡️": 1,
	// 	"🦘": 1,
	// 	"⟳": 1,
	// 	"📋": 1,
	//  "🥇": 1,
	//  "❦": 1,
	//  "↓": 1,
	// }
	return utf8.RuneCountInString(message)
}

func truncateRight(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	return string([]rune(message)[0:width-1]) + ellipsis
}

func truncateLeft(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	strLen := utf8.RuneCountInString(message)
	return ellipsis + string([]rune(message)[strLen-width:strLen-1])
}

func justifyRight(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	strLen := utf8.RuneCountInString(message)
	return ellipsis + string([]rune(message)[strLen-width:strLen-1])
}

func centerText(message, pad string, width int) string {
	if utf8.RuneCountInString(message) > width {
		return truncateRight(message, width)
	}
	leftover := width - utf8.RuneCountInString(message)
	left := leftover / 2
	right := leftover - left

	if pad == "" {
		pad = " "
	}

	leftString := ""
	for utf8.RuneCountInString(leftString) <= left && utf8.RuneCountInString(leftString) <= right {
		leftString += pad
	}

	return fmt.Sprintf("%s%s%s", string([]rune(leftString)[0:left]), message, string([]rune(leftString)[0:right]))
}

func fillSpace(message string, width int) string {
	msgLen := utf8.RuneCountInString(message)
	// msgLen := len(message)
	if msgLen > width {
		return truncateRight(message, width)
	}
	leftover := width - msgLen

	rightString := ""
	rightLen := 0
	for rightLen < leftover {
		rightString += " "
		rightLen = utf8.RuneCountInString(rightString)
		// rightLen = len(rightString)
	}
	return message + rightString
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
	if item.Attack > 0 {
		itemStr += fmt.Sprintf(" attack=%d", item.Attack)
	}
	if len(item.PreItem) > 0 {
		itemStr += fmt.Sprintf(" piece=\"%s\"", item.PreItem)
	}
	return itemStr
}

func carryItemDesc(item *loud.Item) string {
	if item == nil {
		return ""
	} else {
		return "Carry: " + formatItem(*item)
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
	attackStr := formatIntRange(itemSpec.Attack)
	if len(attackStr) > 0 {
		itemStr += fmt.Sprintf(" attack=%s", attackStr)
	}
	return itemStr
}

func formatCharacter(ch loud.Character) string {
	chStr := loud.Localize(ch.Name)
	if ch.GiantKill > 0 {
		chStr = fmt.Sprintf("🥇x%d %s", ch.GiantKill, chStr)
	}
	if ch.Level > 0 {
		chStr += fmt.Sprintf(" Lv%d", ch.Level)
	}
	if ch.XP > 0 {
		chStr += fmt.Sprintf(" XP=%.0f", ch.XP)
	}
	if ch.HP > 0 {
		chStr += fmt.Sprintf(" HP=%d", ch.HP)
	}
	if ch.MaxHP > 0 {
		chStr += fmt.Sprintf(" MaxHP=%d", ch.MaxHP)
	}
	return chStr
}

func formatCharacterSpec(chs loud.CharacterSpec) string {
	chStr := loud.Localize(chs.Name)
	lvlStr := formatIntRange(chs.Level)
	if len(lvlStr) > 0 {
		chStr += fmt.Sprintf(" Lv%s", lvlStr)
	}
	xpStr := formatFloat64Range(chs.XP)
	if len(xpStr) > 0 {
		chStr += fmt.Sprintf(" XP=%s", xpStr)
	}
	hpStr := formatIntRange(chs.HP)
	if len(hpStr) > 0 {
		chStr += fmt.Sprintf(" HP=%s", hpStr)
	}
	maxHpStr := formatIntRange(chs.MaxHP)
	if len(maxHpStr) > 0 {
		chStr += fmt.Sprintf(" MaxHP=%s", maxHpStr)
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

func (screen *GameScreen) renderTRLine(text1 string, text2 string, text3 string, isActiveLine bool, isDisabledLine bool) string {
	text1 = loud.Localize(text1)
	text2 = loud.Localize(text2)
	text3 = loud.Localize(text3)

	calcText := "│" + centerText(text1, " ", 20) + "│" + centerText(text2, " ", 15) + "│" + centerText(text3, " ", 15) + "│"
	onColor := screen.regularFont()
	if isActiveLine && isDisabledLine {
		onColor = screen.brownBoldFont()
	} else if isActiveLine {
		onColor = screen.blueBoldFont()
	} else if isDisabledLine {
		onColor = screen.brownFont()
	}
	return onColor(calcText)
}

func (screen *GameScreen) renderItemTableLine(text1 string, isActiveLine bool) string {
	calcText := "│" + centerText(loud.Localize(text1), " ", 52) + "│"
	onColor := screen.regularFont()
	if isActiveLine {
		onColor = screen.blueBoldFont()
	}
	return onColor(calcText)
}

func (screen *GameScreen) renderItemTrdReqTableLine(text1 string, text2 string, isActiveLine bool, isDisabledLine bool) string {
	text1 = loud.Localize(text1)
	text2 = loud.Localize(text2)
	calcText := "│" + centerText(text1, " ", 36) + "│" + centerText(text2, " ", 15) + "│"
	onColor := screen.regularFont()
	if isActiveLine && isDisabledLine {
		onColor = screen.brownBoldFont()
	} else if isActiveLine {
		onColor = screen.blueBoldFont()
	} else if isDisabledLine {
		onColor = screen.brownFont()
	}
	return onColor(calcText)
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
