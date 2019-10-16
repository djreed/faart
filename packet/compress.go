package packet

import (
	"compress/flate"
)

const (
	COMPRESS_FOCUS = flate.BestCompression
)

func Compress(baseData []byte) ([]byte, error) {
	return baseData, nil
	// compressedData := bytes.NewBuffer(baseData)
	// writer, _ := flate.NewWriter(compressedData, COMPRESS_FOCUS)
	// _, err := writer.Write(baseData)
	// defer writer.Close()
	// writer.Flush()
	// return compressedData.Bytes(), err
}

func Decompress(compressed []byte) ([]byte, error) {
	return compressed, nil
	// compressedData := bytes.NewBuffer(compressed)
	// // decompressedData := bytes.NewBuffer(make([]byte, UPPER_BOUND))
	// reader, err := zlib.NewReader(compressedData)
	// if err != nil {
	// 	panic(err)
	// }
	// return ioutil.ReadAll(reader)
}
