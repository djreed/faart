package packet

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
)

const (
	COMPRESS_LEVEL = flate.BestCompression
)

func Compress(baseData []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(baseData)
	w.Close()
	return b.Bytes(), nil
}

func Decompress(compressed []byte) ([]byte, error) {
	var b = bytes.NewBuffer(compressed)
	r, err := gzip.NewReader(b)
	var targetBuffer bytes.Buffer
	io.Copy(&targetBuffer, r)
	r.Close()
	return targetBuffer.Bytes(), err
}
