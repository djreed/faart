package packet

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"

	"github.com/djreed/faart/log"
)

const (
	COMPRESS_LEVEL = flate.BestCompression
)

func Compress(baseData []byte) ([]byte, error) {
	log.ERR.Printf("Compressing %d bytes", len(baseData))
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(baseData)
	w.Close()
	return b.Bytes(), nil
}

func Decompress(compressed []byte) ([]byte, error) {
	log.ERR.Printf("Decompressing %d bytes", len(compressed))
	var b = bytes.NewBuffer(compressed)
	r, err := gzip.NewReader(b)
	var targetBuffer bytes.Buffer
	io.Copy(&targetBuffer, r)
	r.Close()
	return targetBuffer.Bytes(), err
}
