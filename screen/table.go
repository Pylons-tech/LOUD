package screen

import (
	"fmt"

	loud "github.com/Pylons-tech/LOUD/data"
)

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

func (screen *GameScreen) renderTRTable(requests []loud.TrdReq, fontFunc func(int, loud.TrdReq) FontType) TextLines {
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
		font := screen.getFontOfTR(startLine+li, request.IsMyTrdReq)
		if fontFunc != nil {
			font = fontFunc(startLine+li, request)
		}
		line := screen.renderTRLine(
			fmt.Sprintf("%.4f", request.Price),
			fmt.Sprintf("%d", request.Amount),
			fmt.Sprintf("%d", request.Total))
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

func (screen *GameScreen) renderITRTable(header string, theads [2]string, requestsSlice interface{}, width int, fontFunc func(int, interface{}) FontType) TextLines {
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
		line := ""
		isMyTrdReq, requestItem, requestPrice := RequestInfo(request)
		font := screen.getFontOfTR(startLine+li, isMyTrdReq)
		if fontFunc != nil {
			font = fontFunc(startLine+li, request)
		}
		line = screen.renderItemTrdReqTableLine(
			fmt.Sprintf("%s  ", formatByStructType(requestItem)),
			fmt.Sprintf("%d", requestPrice))
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(tableBodyFooterSeparator(TABLE_SEPARATORS, endLine == len(requests)))
	return tableLines
}

func (screen *GameScreen) renderITTable(header string, th string, itemSlice interface{}, width int, fontFunc func(int, interface{}) FontType) TextLines {
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
		index := startLine + li
		font := screen.getFontByActiveIndex(index)
		if fontFunc != nil {
			font = fontFunc(startLine+li, item)
		}
		line := screen.renderItemTableLine(
			index,
			fmt.Sprintf("%s  ", formatByStructType(item)))
		tableLines = tableLines.appendF(line, font)
	}
	tableLines = tableLines.
		append(tableBodyFooterSeparator(TABLE_SEPARATORS, endLine == len(items)))
	return tableLines
}
