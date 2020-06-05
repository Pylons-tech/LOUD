package loud

import (
	"bytes"
	"strings"

	"github.com/vmihailenco/msgpack"
)

// MSGPack packs to msgpack using JSON rules
func MSGPack(target interface{}) ([]byte, error) {
	var outBuffer bytes.Buffer

	writer := msgpack.NewEncoder(&outBuffer)
	writer.UseJSONTag(true)
	err := writer.Encode(target)

	return outBuffer.Bytes(), err
}

// MSGUnpack unpacks from msgpack using JSON rules
func MSGUnpack(inBytes []byte, outItem interface{}) error {
	var inBuffer = bytes.NewBuffer(inBytes)

	reader := msgpack.NewDecoder(inBuffer)
	reader.UseJSONTag(true)
	err := reader.Decode(outItem)

	return err
}

// ChunkString split string by specific size
func ChunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}

// ChunkText split text including line break into specific size
func ChunkText(bigtext string, width int) []string {
	basicLines := strings.Split(bigtext, "\n")
	infoLines := []string{}
	for _, text := range basicLines {
		infoLines = append(infoLines, ChunkString(text, width)...)
	}
	return infoLines
}
