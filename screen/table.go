package screen

import (
	"fmt"

	loud "github.com/Pylons-tech/LOUD/data"
)

// FontFuncType is a function to convert an object to font type
type FontFuncType func(int, interface{}) FontType

func tableHeaderBodySeparator(TableSeparators []string, showFull bool) string {
	if showFull {
		return TableSeparators[1]
	}
	return TableSeparators[2]
}

func tableBodyFooterSeparator(TableSeparators []string, showFull bool) string {
	if showFull {
		return TableSeparators[3]
	}
	return TableSeparators[4]
}

func (screen *GameScreen) calcTLFont(fontFunc FontFuncType, idx int, disabled bool, request interface{}) FontType {
	if fontFunc != nil {
		return fontFunc(idx, request)
	}
	return screen.getFontOfTableLine(idx, disabled)
}

// TableHeightWindow determines table startLine and endLine from existing header
func (screen *GameScreen) TableHeightWindow(header string, rawArrayInterFace interface{}, width int) ([]string, int, int) {
	rawArray := InterfaceSlice(rawArrayInterFace)
	infoLines := loud.ChunkText(loud.Localize(header), width)
	numHeaderLines := len(infoLines)

	if screen.activeLine >= len(rawArray) {
		screen.activeLine = len(rawArray) - 1
	}
	startLine, endLine := getWindowFromActiveLine(
		screen.activeLine,
		screen.GetSituationBox().H-5-numHeaderLines,
		len(rawArray))
	return infoLines, startLine, endLine
}

// TableHeader convert table header params into TextLines
func (screen *GameScreen) TableHeader(titleLines []string, TableSeparators []string, headerVisual string, startLine int) TextLines {
	return TextLines{}.
		append(titleLines...).
		append(TableSeparators[0]).
		append(headerVisual).
		append(tableHeaderBodySeparator(TableSeparators, startLine == 0))
}

// TableFooter convert table footer params into TextLines
func (screen *GameScreen) TableFooter(TableSeparators []string, endLine int, rawArrayInterFace interface{}) string {
	rawArray := InterfaceSlice(rawArrayInterFace)
	return tableBodyFooterSeparator(TableSeparators, endLine == len(rawArray))
}

// renderTRTable convert trade request table params into TextLines
func (screen *GameScreen) renderTRTable(rawArray []loud.TrdReq, width int, fontFunc FontFuncType) TextLines {
	TableSeparators := []string{
		"╭────────────────────┬───────────────┬───────────────╮",
		"├────────────────────┼───────────────┼───────────────┤",
		"├─────↓↓↓↓↓↓↓↓↓──────┼───↓↓↓↓↓↓↓↓↓───┼───↓↓↓↓↓↓↓↓↓───┤",
		"╰────────────────────┴───────────────┴───────────────╯",
		"╰──────↑↑↑↑↑↑↑↑↑─────┴───↑↑↑↑↑↑↑↑↑───┴───↑↑↑↑↑↑↑↑↑───╯",
	}

	_, startLine, endLine := screen.TableHeightWindow("", rawArray, width)

	tableLines := screen.TableHeader(
		[]string{},
		TableSeparators,
		screen.renderTRLine("GOLD price (pylon)", "Amount (gold)", "Total (pylon)"),
		startLine,
	)

	for li, request := range rawArray[startLine:endLine] {
		line := screen.renderTRLine(
			fmt.Sprintf("%.4f", request.Price),
			fmt.Sprintf("%d", request.Amount),
			fmt.Sprintf("%d", request.Total))
		font := screen.calcTLFont(fontFunc, startLine+li, request.IsMyTrdReq, request)
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(screen.TableFooter(TableSeparators, endLine, rawArray))
	return tableLines
}

func (screen *GameScreen) renderITRTable(header string, theads [2]string, rawArrayInterFace interface{}, width int, fontFunc FontFuncType) TextLines {
	rawArray := InterfaceSlice(rawArrayInterFace)
	TableSeparators := []string{
		"╭────────────────────────────────────┬───────────────╮",
		"├────────────────────────────────────┼───────────────┤",
		"├──────────↓↓↓↓↓↓↓↓↓↓↓↓↓↓────────────┼─────↓↓↓↓↓↓────┤",
		"╰────────────────────────────────────┴───────────────╯",
		"╰──────────↑↑↑↑↑↑↑↑↑↑↑↑↑↑────────────┴─────↑↑↑↑↑↑────╯",
	}

	infoLines, startLine, endLine := screen.TableHeightWindow(header, rawArrayInterFace, width)
	tableLines := screen.TableHeader(
		infoLines,
		TableSeparators,
		screen.renderItemTrdReqTableLine(theads[0], theads[1]),
		startLine,
	)

	for li, request := range rawArray[startLine:endLine] {
		isMyTrdReq, requestItem, requestPrice := RequestInfo(request)
		line := screen.renderItemTrdReqTableLine(
			fmt.Sprintf("%s  ", formatByStructType(requestItem)),
			fmt.Sprintf("%d", requestPrice))
		font := screen.calcTLFont(fontFunc, startLine+li, isMyTrdReq, request)
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(screen.TableFooter(TableSeparators, endLine, rawArray))
	return tableLines
}

func (screen *GameScreen) renderITTable(header string, th string, rawArrayInterFace interface{}, width int, fontFunc FontFuncType) TextLines {
	rawArray := InterfaceSlice(rawArrayInterFace)
	TableSeparators := []string{
		"╭────────────────────────────────────────────────────╮",
		"├────────────────────────────────────────────────────┤",
		"├──────────────────↓↓↓↓↓↓↓↓↓↓↓↓↓↓────────────────────┤",
		"╰────────────────────────────────────────────────────╯",
		"╰───────────────────↑↑↑↑↑↑↑↑↑↑↑↑↑↑───────────────────╯",
	}
	infoLines, startLine, endLine := screen.TableHeightWindow(header, rawArrayInterFace, width)
	tableLines := screen.TableHeader(
		infoLines,
		TableSeparators,
		screen.renderItemTableLine(-1, th),
		startLine,
	)

	for li, item := range rawArray[startLine:endLine] {
		line := screen.renderItemTableLine(
			startLine+li,
			fmt.Sprintf("%s  ", formatByStructType(item)))
		font := screen.calcTLFont(fontFunc, startLine+li, false, item)
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(screen.TableFooter(TableSeparators, endLine, rawArray))
	return tableLines
}
