package packet

func Compress(baseData []byte) ([]byte, error) {
	// var b bytes.Buffer
	// w := zlib.NewWriter(&b)
	// w.Write(baseData)
	// w.Close()
	// return b.Bytes(), nil
	return baseData, nil
}

func Decompress(compressed []byte) ([]byte, error) {
	// var b = bytes.NewBuffer(compressed)
	// r, err := zlib.NewReader(b)
	// var targetBuffer bytes.Buffer
	// io.Copy(&targetBuffer, r)
	// r.Close()
	// return targetBuffer.Bytes(), err
	return compressed, nil
}
