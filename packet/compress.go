package packet

import (
	"bytes"
	"compress/flate"
)

const (
	COMPRESS_FOCUS = flate.BestCompression
)

func Compress(src []byte) ([]byte, error) {
	compressedData := new(bytes.Buffer)
	writer, _ := flate.NewWriter(compressedData, COMPRESS_FOCUS)
	writer.Write(src)
	writer.Close()
	return compressedData.Bytes(), nil
}

func Decompress(compressed []byte) ([]byte, error) {
	compressedData := bytes.NewBuffer(compressed)
	reader := flate.NewReader(compressedData)
	data := make([]byte, 0)
	for n, _ := reader.Read(data); n > 0; {
	}
	return data, nil
}
