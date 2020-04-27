package screen

import (
	"fmt"
	"io"
	"os"

	"github.com/ahmetb/go-cursor"
)

func (screen *GameScreen) renderInputValue() {
	inputBoxWidth := uint32(screen.leftInnerWidth())
	inputWidth := inputBoxWidth - 9
	move := cursor.MoveTo(screen.Height()-1, 2)

	chatFunc := screen.regularFont()
	chat := chatFunc("👉👉👉 ")
	fmtString := fmt.Sprintf("%%-%vs", inputWidth)

	if screen.InputActive() {
		chatFunc = screen.inputActiveFont()
	}

	fixedChat := truncateLeft(screen.inputText, int(inputWidth))
	inputText := fmt.Sprintf("%s%s%s", move, chat, chatFunc(fmt.Sprintf(fmtString, fixedChat)))

	if !screen.InputActive() {
		inputText = fmt.Sprintf("%s%s", move, chatFunc(fillSpace(screen.actionText, int(inputBoxWidth))))
	}

	io.WriteString(os.Stdout, inputText)
}
