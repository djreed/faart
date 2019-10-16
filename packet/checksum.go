package packet

import "crypto/sha1"

// MD5 hashes are limited to 16 bytes
func CalculateChecksum(data []byte) []byte {
	var hash [CHECKSUM_SIZE]byte
	hash = sha1.Sum(data)
	return hash[:]
}
