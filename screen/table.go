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
		var requestItem interface{}
		var requestPrice int
		var isMyTrdReq bool
		switch request.(type) {
		case loud.ItemBuyTrdReq:
			itr := request.(loud.ItemBuyTrdReq)
			isMyTrdReq, requestItem, requestPrice = itr.IsMyTrdReq, itr.TItem, itr.Price
		case loud.ItemSellTrdReq:
			itr := request.(loud.ItemSellTrdReq)
			isMyTrdReq, requestItem, requestPrice = itr.IsMyTrdReq, itr.TItem, itr.Price
		case loud.CharacterBuyTrdReq:
			itr := request.(loud.CharacterBuyTrdReq)
			isMyTrdReq, requestItem, requestPrice = itr.IsMyTrdReq, itr.TCharacter, itr.Price
		case loud.CharacterSellTrdReq:
			itr := request.(loud.CharacterSellTrdReq)
			isMyTrdReq, requestItem, requestPrice = itr.IsMyTrdReq, itr.TCharacter, itr.Price
		}
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
		line := ""
		switch item.(type) {
		case loud.Item:
			itemT := item.(loud.Item)
			index := startLine + li
			font := screen.getFontByActiveIndex(index)
			if fontFunc != nil {
				font = fontFunc(startLine+li, item)
			}
			line = screen.renderItemTableLine(
				index,
				fmt.Sprintf("%s  ", formatItem(itemT)),
				font,
				width,
			)
		case loud.Character:
			itemT := item.(loud.Character)
			index := startLine + li
			font := screen.getFontByActiveIndex(index)
			if fontFunc != nil {
				font = fontFunc(startLine+li, item)
			}
			line = screen.renderItemTableLine(
				index,
				fmt.Sprintf("%s  ", formatCharacter(itemT)),
				font,
				width,
			)
		case loud.ItemSpec:
			itemT := item.(loud.ItemSpec)
			index := startLine + li
			font := screen.getFontByActiveIndex(index)
			if fontFunc != nil {
				font = fontFunc(startLine+li, item)
			}
			line = screen.renderItemTableLine(
				index,
				fmt.Sprintf("%s  ", formatItemSpec(itemT)),
				font,
				width,
			)
		case loud.CharacterSpec:
			itemT := item.(loud.CharacterSpec)
			index := startLine + li
			font := screen.getFontByActiveIndex(index)
			if fontFunc != nil {
				font = fontFunc(startLine+li, item)
			}
			line = screen.renderItemTableLine(
				index,
				fmt.Sprintf("%s  ", formatCharacterSpec(itemT)),
				font,
				width,
			)
		}
		tableLines = append(tableLines, line)
	}
	if endLine == len(items) {
		tableLines = append(tableLines, fmtFunc(fillSpace("╰────────────────────────────────────────────────────╯", width)))
	} else {
		tableLines = append(tableLines, fmtFunc(fillSpace("╰───────────────────↑↑↑↑↑↑↑↑↑↑↑↑↑↑───────────────────╯", width)))
	}
	return infoLines, tableLines
}
