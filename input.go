package loud

import (
	"bufio"
	"context"
)

type inputEvent struct {
	inputString string
	err         error
}

func handleKeys(reader *bufio.Reader, stringChannel chan<- inputEvent, cancel context.CancelFunc) {
}
