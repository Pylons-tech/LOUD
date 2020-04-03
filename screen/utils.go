package screen

import (
	"fmt"
	"reflect"
	"unicode/utf8"

	loud "github.com/Pylons-tech/LOUD/data"
)

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

func formatItem(item loud.Item) string {
	return fmt.Sprintf("%s Lv%d", loud.Localize(item.Name), item.Level)
}

func formatItemSpec(itemSpec loud.ItemSpec) string {
	return fmt.Sprintf("%s Lv%d-%d", loud.Localize(itemSpec.Name), itemSpec.Level[0], itemSpec.Level[1])
}

func formatCharacter(ch loud.Character) string {
	return fmt.Sprintf("%s Lv%d XP=%.0f HP=%d", loud.Localize(ch.Name), ch.Level, ch.XP, ch.HP)
}

func formatCharacterSpec(chs loud.CharacterSpec) string {
	return fmt.Sprintf("%s Lv%d-%d XP=%.0f-%.0f HP=%d-%d MaxHP=%d-%d", loud.Localize(chs.Name),
		chs.Level[0], chs.Level[1],
		chs.XP[0], chs.XP[1],
		chs.HP[0], chs.HP[1],
		chs.MaxHP[0], chs.MaxHP[1],
	)
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
	calcText := "│" + centerText(text1, " ", 20) + "│" + centerText(text2, " ", 15) + "│" + centerText(text3, " ", 15) + "│"
	if isActiveLine && isDisabledLine {
		onColor := screen.brownBoldFont()
		return onColor(calcText)
	} else if isActiveLine {
		onColor := screen.blueBoldFont()
		return onColor(calcText)
	} else if isDisabledLine {
		onColor := screen.brownFont()
		return onColor(calcText)
	}
	return calcText
}

func (screen *GameScreen) renderItemTableLine(text1 string, isActiveLine bool) string {
	calcText := "│" + centerText(text1, " ", 52) + "│"
	if isActiveLine {
		onColor := screen.blueBoldFont()
		return onColor(calcText)
	}
	return calcText
}

func (screen *GameScreen) renderItemTradeRequestTableLine(text1 string, text2 string, isActiveLine bool, isDisabledLine bool) string {
	calcText := "│" + centerText(text1, " ", 36) + "│" + centerText(text2, " ", 15) + "│"
	if isActiveLine && isDisabledLine {
		onColor := screen.brownBoldFont()
		return onColor(calcText)
	} else if isActiveLine {
		onColor := screen.blueBoldFont()
		return onColor(calcText)
	} else if isDisabledLine {
		onColor := screen.brownFont()
		return onColor(calcText)
	}
	return calcText
}
