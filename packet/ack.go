package packet

import "net"

var (
	// Maximal UDP Ack Size
	ACK_SIZE = SEQUENCE_SIZE + OFFSET_SIZE
)

type AddressedAck struct {
	Ack
	Addr *net.UDPAddr
}

type Ack []byte

func NewAck() Ack {
	return make([]byte, ACK_SIZE)
}

func CreateAck(datagram Datagram) Ack {
	ack := NewAck()
	ack.SetSequence(datagram.Headers().Sequence())
	ack.SetOffset(datagram.Headers().Offset())
	return ack
}

func (ack Ack) Sequence() SeqID {
	return SeqID(bytesToUint32(ack[SEQUENCE_POINTER : SEQUENCE_POINTER+SEQUENCE_SIZE]))
}
func (ack Ack) SetSequence(seq SeqID) {
	copy(ack[SEQUENCE_POINTER:SEQUENCE_POINTER+SEQUENCE_SIZE], uint32ToBytes(uint32(seq)))
}

func (ack Ack) Offset() OffsetVal {
	return OffsetVal(bytesToUint32(ack[OFFSET_POINTER : OFFSET_POINTER+OFFSET_SIZE]))
}
func (ack Ack) SetOffset(offset OffsetVal) {
	copy(ack[OFFSET_POINTER:OFFSET_POINTER+OFFSET_SIZE], uint32ToBytes(uint32(offset)))
}
