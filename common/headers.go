package common

const (
	SEQUENCE_POINTER = 0
	SEQUENCE_SIZE    = 4 // Sequence ID = 4 bytes

	OFFSET_POINTER = SEQUENCE_SIZE + SEQUENCE_POINTER
	OFFSET_SIZE    = 4 // Offset into file = 4 bytes

	CHECKSUM_POINTER = OFFSET_SIZE + SEQUENCE_SIZE
	CHECKSUM_SIZE    = 16 // MD5 Hash = 16 bytes
	CHECKSUM_PREV    = OFFSET_SIZE

	HEADER_SIZE = SEQUENCE_SIZE + OFFSET_SIZE + CHECKSUM_SIZE // Sum of all the header sizes = 24 bytes
)

//////////////////////
/* Datagram Headers */
//////////////////////

type Header struct {
	// [ SEQUENCE, OFFSET, CHECKSUM ]
	headers [HEADER_SIZE]byte
}

// Sequence ID of the current packet
func (h *Header) Sequence() [SEQUENCE_SIZE]byte {
	var sequence [SEQUENCE_SIZE]byte
	copy(sequence[:], h.headers[:OFFSET_POINTER])
	return sequence
}

// Sequence ID of the current packet
func (h *Header) Offset() [OFFSET_SIZE]byte {
	var offset [OFFSET_SIZE]byte
	copy(offset[:], h.headers[OFFSET_POINTER:CHECKSUM_POINTER])
	return offset
}

// Checksum hash of Data (MD5)
func (h *Header) Checksum() [CHECKSUM_SIZE]byte {
	var checksum [CHECKSUM_SIZE]byte
	copy(checksum[:], h.headers[CHECKSUM_POINTER:])
	return checksum
}
