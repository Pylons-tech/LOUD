package screen

import (
	"fmt"

	loud "github.com/Pylons-tech/LOUD/data"
)

func (screen *GameScreen) renderTRTable(requests []loud.TrdReq, width int, fontFunc func(int, loud.TrdReq) FontType) ([]string, []string) {
	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	startLine, endLine := getWindowFromActiveLine(screen.activeLine, screen.GetSituationBox().H-5, len(requests))
	tableLines := []string{}
	tableLines = append(tableLines, screen.regularFont()(fillSpace("╭────────────────────┬───────────────┬───────────────╮", width)))
	tableLines = append(tableLines, screen.renderTRLine("GOLD price (pylon)", "Amount (gold)", "Total (pylon)", REGULAR, width))
	if startLine == 0 {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("├────────────────────┼───────────────┼───────────────┤", width)))
	} else {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("├─────↓↓↓↓↓↓↓↓↓──────┼───↓↓↓↓↓↓↓↓↓───┼───↓↓↓↓↓↓↓↓↓───┤", width)))
	}
	for li, request := range requests[startLine:endLine] {
		font := screen.getFontOfTR(startLine+li, request.IsMyTrdReq)
		if fontFunc != nil {
			font = fontFunc(startLine+li, request)
		}
		tableLines = append(
			tableLines,
			screen.renderTRLine(
				fmt.Sprintf("%.4f", request.Price),
				fmt.Sprintf("%d", request.Amount),
				fmt.Sprintf("%d", request.Total),
				font, width),
		)
	}
	if endLine == len(requests) {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("╰────────────────────┴───────────────┴───────────────╯", width)))
	} else {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("╰──────↑↑↑↑↑↑↑↑↑─────┴───↑↑↑↑↑↑↑↑↑───┴───↑↑↑↑↑↑↑↑↑───╯", width)))
	}
	return []string{}, tableLines
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

func (screen *GameScreen) renderITRTable(header string, theads [2]string, requestsSlice interface{}, width int, fontFunc func(int, interface{}) FontType) ([]string, []string) {
	requests := InterfaceSlice(requestsSlice)

	infoLines := loud.ChunkText(loud.Localize(header), width)
	numHeaderLines := len(infoLines)

	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	startLine, endLine := getWindowFromActiveLine(
		screen.activeLine,
		screen.GetSituationBox().H-5-numHeaderLines,
		len(requests))
	tableLines := []string{}
	tableLines = append(tableLines, screen.regularFont()(fillSpace("╭────────────────────────────────────┬───────────────╮", width)))
	tableLines = append(tableLines, screen.renderItemTrdReqTableLine(theads[0], theads[1], REGULAR, width))
	if startLine == 0 {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("├────────────────────────────────────┼───────────────┤", width)))
	} else {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("├──────────↓↓↓↓↓↓↓↓↓↓↓↓↓↓────────────┼─────↓↓↓↓↓↓────┤", width)))
	}
	for li, request := range requests[startLine:endLine] {
		line := ""
		isMyTrdReq, requestItem, requestPrice := RequestInfo(request)
		font := screen.getFontOfTR(startLine+li, isMyTrdReq)
		if fontFunc != nil {
			font = fontFunc(startLine+li, request)
		}
		line = screen.renderItemTrdReqTableLine(
			fmt.Sprintf("%s  ", formatByStructType(requestItem)),
			fmt.Sprintf("%d", requestPrice),
			font, width)
		tableLines = append(tableLines, line)
	}
	if endLine == len(requests) {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("╰────────────────────────────────────┴───────────────╯", width)))
	} else {
		tableLines = append(tableLines, screen.regularFont()(fillSpace("╰──────────↑↑↑↑↑↑↑↑↑↑↑↑↑↑────────────┴─────↑↑↑↑↑↑────╯", width)))
	}
	return infoLines, tableLines
}

func (screen *GameScreen) renderITTable(header string, th string, itemSlice interface{}, width int, fontFunc func(int, interface{}) FontType) ([]string, []string) {
	items := InterfaceSlice(itemSlice)

	infoLines := loud.ChunkText(loud.Localize(header), width)
	numHeaderLines := len(infoLines)
	fmtFunc := screen.regularFont()

	if screen.activeLine >= len(items) {
		screen.activeLine = len(items) - 1
	}
	startLine, endLine := getWindowFromActiveLine(
		screen.activeLine,
		screen.GetSituationBox().H-5-numHeaderLines,
		len(items))

	tableLines := []string{}
	tableLines = append(tableLines, fmtFunc(fillSpace("╭────────────────────────────────────────────────────╮", width)))
	tableLines = append(tableLines, screen.renderItemTableLine(-1, th, REGULAR, width))
	if startLine == 0 {
		tableLines = append(tableLines, fmtFunc(fillSpace("├────────────────────────────────────────────────────┤", width)))
	} else {
		tableLines = append(tableLines, fmtFunc(fillSpace("├──────────────────↓↓↓↓↓↓↓↓↓↓↓↓↓↓────────────────────┤", width)))
	}
	for li, item := range items[startLine:endLine] {
		index := startLine + li
		font := screen.getFontByActiveIndex(index)
		if fontFunc != nil {
			font = fontFunc(startLine+li, item)
		}
		line := screen.renderItemTableLine(
			index,
			fmt.Sprintf("%s  ", formatByStructType(item)),
			font, width)
		tableLines = append(tableLines, line)
	}
	if endLine == len(items) {
		tableLines = append(tableLines, fmtFunc(fillSpace("╰────────────────────────────────────────────────────╯", width)))
	} else {
		tableLines = append(tableLines, fmtFunc(fillSpace("╰───────────────────↑↑↑↑↑↑↑↑↑↑↑↑↑↑───────────────────╯", width)))
	}
	return infoLines, tableLines
}
