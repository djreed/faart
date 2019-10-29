package packet

var (
	// Maximal UDP Ack Size
	ACK_SIZE = OFFSET_SIZE
)

type Ack []byte

func NewAck() Ack {
	return make([]byte, ACK_SIZE)
}

func CreateAck(datagram Datagram) Ack {
	ack := NewAck()
	ack.SetOffset(datagram.Headers().Offset())
	return ack
}

// Sequence ID of the current packet
func (ack Ack) Offset() OffsetVal {
	return OffsetVal(bytesToUint32(ack[OFFSET_POINTER : OFFSET_POINTER+OFFSET_SIZE]))
}

// Copy in the offset into the file this data represents
func (ack Ack) SetOffset(offset OffsetVal) {
	copy(ack[OFFSET_POINTER:OFFSET_POINTER+OFFSET_SIZE], uint32ToBytes(uint32(offset)))
}
