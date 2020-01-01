package loud

import (
	"bytes"

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
