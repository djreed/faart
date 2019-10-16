package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sample = []byte("Hello World")

func TestCompressOutput(t *testing.T) {
	compressed, err := Compress(sample)
	assert.NotNil(t, compressed)
	assert.Nil(t, err)
}

func TestDecompressOutput(t *testing.T) {
	compressed, _ := Compress(sample)
	decompressed, err := Decompress(compressed)
	assert.NotNil(t, decompressed)
	assert.Nil(t, err)
}

func TestCompressionInvertible(t *testing.T) {
	compressed, _ := Compress(sample)
	decompressed, _ := Decompress(compressed)
	assert.Equal(t, sample, decompressed)
}
