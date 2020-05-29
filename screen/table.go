package screen

import (
	"fmt"

	loud "github.com/Pylons-tech/LOUD/data"
)

type FontFuncType func(int, interface{}) FontType

func tableHeaderBodySeparator(TABLE_SEPARATORS []string, showFull bool) string {
	if showFull {
		return TABLE_SEPARATORS[1]
	}
	return TABLE_SEPARATORS[2]
}

func tableBodyFooterSeparator(TABLE_SEPARATORS []string, showFull bool) string {
	if showFull {
		return TABLE_SEPARATORS[3]
	}
	return TABLE_SEPARATORS[4]
}

func (screen *GameScreen) calcTLFont(fontFunc FontFuncType, idx int, isMyTrdReq bool, request interface{}) FontType {
	if fontFunc != nil {
		return fontFunc(idx, request)
	}
	return screen.getFontOfTableLine(idx, isMyTrdReq)
}

func (screen *GameScreen) renderTRTable(requests []loud.TrdReq, fontFunc FontFuncType) TextLines {
	TABLE_SEPARATORS := []string{
		"╭────────────────────┬───────────────┬───────────────╮",
		"├────────────────────┼───────────────┼───────────────┤",
		"├─────↓↓↓↓↓↓↓↓↓──────┼───↓↓↓↓↓↓↓↓↓───┼───↓↓↓↓↓↓↓↓↓───┤",
		"╰────────────────────┴───────────────┴───────────────╯",
		"╰──────↑↑↑↑↑↑↑↑↑─────┴───↑↑↑↑↑↑↑↑↑───┴───↑↑↑↑↑↑↑↑↑───╯",
	}
	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	startLine, endLine := getWindowFromActiveLine(screen.activeLine, screen.GetSituationBox().H-5, len(requests))
	tableLines := TextLines{}
	tableLines = tableLines.
		append(TABLE_SEPARATORS[0]).
		append(screen.renderTRLine("GOLD price (pylon)", "Amount (gold)", "Total (pylon)")).
		append(tableHeaderBodySeparator(TABLE_SEPARATORS, startLine == 0))

	for li, request := range requests[startLine:endLine] {
		line := screen.renderTRLine(
			fmt.Sprintf("%.4f", request.Price),
			fmt.Sprintf("%d", request.Amount),
			fmt.Sprintf("%d", request.Total))
		font := screen.calcTLFont(fontFunc, startLine+li, request.IsMyTrdReq, request)
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(tableBodyFooterSeparator(TABLE_SEPARATORS, endLine == len(requests)))
	return tableLines
}

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

func (screen *GameScreen) renderITRTable(header string, theads [2]string, requestsSlice interface{}, width int, fontFunc FontFuncType) TextLines {
	requests := InterfaceSlice(requestsSlice)
	TABLE_SEPARATORS := []string{
		"╭────────────────────────────────────┬───────────────╮",
		"├────────────────────────────────────┼───────────────┤",
		"├──────────↓↓↓↓↓↓↓↓↓↓↓↓↓↓────────────┼─────↓↓↓↓↓↓────┤",
		"╰────────────────────────────────────┴───────────────╯",
		"╰──────────↑↑↑↑↑↑↑↑↑↑↑↑↑↑────────────┴─────↑↑↑↑↑↑────╯",
	}

	infoLines := loud.ChunkText(loud.Localize(header), width)
	numHeaderLines := len(infoLines)

	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	startLine, endLine := getWindowFromActiveLine(
		screen.activeLine,
		screen.GetSituationBox().H-5-numHeaderLines,
		len(requests))
	tableLines := TextLines{}
	tableLines = tableLines.
		append(infoLines...).
		append(TABLE_SEPARATORS[0]).
		append(screen.renderItemTrdReqTableLine(theads[0], theads[1])).
		append(tableHeaderBodySeparator(TABLE_SEPARATORS, startLine == 0))
	for li, request := range requests[startLine:endLine] {
		isMyTrdReq, requestItem, requestPrice := RequestInfo(request)
		line := screen.renderItemTrdReqTableLine(
			fmt.Sprintf("%s  ", formatByStructType(requestItem)),
			fmt.Sprintf("%d", requestPrice))
		font := screen.calcTLFont(fontFunc, startLine+li, isMyTrdReq, request)
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(tableBodyFooterSeparator(TABLE_SEPARATORS, endLine == len(requests)))
	return tableLines
}

func (screen *GameScreen) renderITTable(header string, th string, itemSlice interface{}, width int, fontFunc FontFuncType) TextLines {
	items := InterfaceSlice(itemSlice)
	TABLE_SEPARATORS := []string{
		"╭────────────────────────────────────────────────────╮",
		"├────────────────────────────────────────────────────┤",
		"├──────────────────↓↓↓↓↓↓↓↓↓↓↓↓↓↓────────────────────┤",
		"╰────────────────────────────────────────────────────╯",
		"╰───────────────────↑↑↑↑↑↑↑↑↑↑↑↑↑↑───────────────────╯",
	}
	infoLines := loud.ChunkText(loud.Localize(header), width)
	numHeaderLines := len(infoLines)

	if screen.activeLine >= len(items) {
		screen.activeLine = len(items) - 1
	}
	startLine, endLine := getWindowFromActiveLine(
		screen.activeLine,
		screen.GetSituationBox().H-5-numHeaderLines,
		len(items))

	tableLines := TextLines{}
	tableLines = tableLines.
		append(infoLines...).
		append(TABLE_SEPARATORS[0]).
		append(screen.renderItemTableLine(-1, th)).
		append(tableHeaderBodySeparator(TABLE_SEPARATORS, startLine == 0))

	for li, item := range items[startLine:endLine] {
		line := screen.renderItemTableLine(
			startLine+li,
			fmt.Sprintf("%s  ", formatByStructType(item)))
		font := screen.calcTLFont(fontFunc, startLine+li, false, item)
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(tableBodyFooterSeparator(TABLE_SEPARATORS, endLine == len(items)))
	return tableLines
}
