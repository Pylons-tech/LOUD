package screen

import (
	"fmt"
	"strings"

	loud "github.com/Pylons-tech/LOUD/data"
)

func (screen *GameScreen) renderTRTable(requests []loud.TrdReq) ([]string, []string) {
	tableLines := []string{}
	tableLines = append(tableLines, screen.regularFont()("╭────────────────────┬───────────────┬───────────────╮"))
	tableLines = append(tableLines, screen.renderTRLine("GOLD price (pylon)", "Amount (gold)", "Total (pylon)", false, false))
	tableLines = append(tableLines, screen.regularFont()("├────────────────────┼───────────────┼───────────────┤"))
	numLines := screen.situationInnerHeight() - 5
	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(requests) {
		endLine = len(requests)
	}
	for li, request := range requests[startLine:endLine] {
		tableLines = append(
			tableLines,
			screen.renderTRLine(
				fmt.Sprintf("%.4f", request.Price),
				fmt.Sprintf("%d", request.Amount),
				fmt.Sprintf("%d", request.Total),
				startLine+li == activeLine,
				request.IsMyTrdReq,
			),
		)
	}
	tableLines = append(tableLines, screen.regularFont()("╰────────────────────┴───────────────┴───────────────╯"))
	return []string{}, tableLines
}

func (screen *GameScreen) renderITRTable(title string, theads [2]string, requestsSlice interface{}) ([]string, []string) {
	requests := InterfaceSlice(requestsSlice)
	infoLines := strings.Split(loud.Localize(title), "\n")
	numHeaderLines := len(infoLines)

	tableLines := []string{}
	tableLines = append(tableLines, screen.regularFont()("╭────────────────────────────────────┬───────────────╮"))
	tableLines = append(tableLines, screen.renderItemTrdReqTableLine(theads[0], theads[1], false, false))
	tableLines = append(tableLines, screen.regularFont()("├────────────────────────────────────┼───────────────┤"))
	numLines := screen.situationInnerHeight() - 5 - numHeaderLines
	if screen.activeLine >= len(requests) {
		screen.activeLine = len(requests) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(requests) {
		endLine = len(requests)
	}
	for li, request := range requests[startLine:endLine] {
		line := ""
		switch request.(type) {
		case loud.ItemBuyTrdReq:
			itr := request.(loud.ItemBuyTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatItemSpec(itr.TItem)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		case loud.ItemSellTrdReq:
			itr := request.(loud.ItemSellTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatItem(itr.TItem)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		case loud.CharacterBuyTrdReq:
			itr := request.(loud.CharacterBuyTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatCharacterSpec(itr.TCharacter)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		case loud.CharacterSellTrdReq:
			itr := request.(loud.CharacterSellTrdReq)
			line = screen.renderItemTrdReqTableLine(
				fmt.Sprintf("%s  ", formatCharacter(itr.TCharacter)),
				fmt.Sprintf("%d", itr.Price),
				startLine+li == activeLine,
				itr.IsMyTrdReq,
			)
		}
		tableLines = append(tableLines, line)
	}
	tableLines = append(tableLines, screen.regularFont()("╰────────────────────────────────────┴───────────────╯"))
	return infoLines, tableLines
}

func (screen *GameScreen) renderITTable(header string, th string, itemSlice interface{}) ([]string, []string) {
	items := InterfaceSlice(itemSlice)
	infoLines := strings.Split(loud.Localize(header), "\n")
	numHeaderLines := len(infoLines)
	numLines := screen.situationInnerHeight() - 5 - numHeaderLines
	fmtFunc := screen.regularFont()

	tableLines := []string{}
	tableLines = append(tableLines, fmtFunc("╭────────────────────────────────────────────────────╮"))
	tableLines = append(tableLines, screen.renderItemTableLine(th, false))
	tableLines = append(tableLines, fmtFunc("├────────────────────────────────────────────────────┤"))
	if screen.activeLine >= len(items) {
		screen.activeLine = len(items) - 1
	}
	activeLine := screen.activeLine
	startLine := activeLine - numLines + 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + numLines
	if endLine > len(items) {
		endLine = len(items)
	}
	for li, item := range items[startLine:endLine] {
		line := ""
		switch item.(type) {
		case loud.Item:
			itemT := item.(loud.Item)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatItem(itemT)),
				startLine+li == activeLine,
			)
		case loud.Character:
			itemT := item.(loud.Character)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatCharacter(itemT)),
				startLine+li == activeLine,
			)
		case loud.ItemSpec:
			itemT := item.(loud.ItemSpec)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatItemSpec(itemT)),
				startLine+li == activeLine,
			)
		case loud.CharacterSpec:
			itemT := item.(loud.CharacterSpec)
			line = screen.renderItemTableLine(
				fmt.Sprintf("%s  ", formatCharacterSpec(itemT)),
				startLine+li == activeLine,
			)
		}
		tableLines = append(tableLines, line)
	}
	tableLines = append(tableLines, fmtFunc("╰────────────────────────────────────────────────────╯"))
	return infoLines, tableLines
}
