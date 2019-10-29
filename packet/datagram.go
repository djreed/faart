package packet

import (
	"bytes"
	"encoding/binary"
	"math"
	"unsafe"
)

const (
	// Maximal UDP datagram size size
	DATAGRAM_SIZE = 1472 // MTU - IP - UDP = 1500 - 20 - 8 = 1472 bytes

	// Offset in the destination file
	OFFSET_POINTER = 0
	OFFSET_SIZE    = 4

	// Checksum of the Packet's data
	CHECKSUM_POINTER = OFFSET_POINTER + OFFSET_SIZE
	CHECKSUM_SIZE    = 20

	// Length of the Packet
	LENGTH_POINTER = CHECKSUM_POINTER + CHECKSUM_SIZE
	LENGTH_SIZE    = 4

	// Whether you're completely done with the data being send
	DONE_FLAG_POINTER = LENGTH_POINTER + LENGTH_SIZE
	DONE_FLAG_SIZE    = 4

	HEADER_SIZE = DONE_FLAG_POINTER + DONE_FLAG_SIZE

	PACKET_SIZE = DATAGRAM_SIZE - HEADER_SIZE
)

//////////////////////////////////////////////
/* The Grand High Poobah of Data Structures */
//////////////////////////////////////////////

// datagram data is just a byte slice
type Datagram []byte

type OffsetVal uint32
type ByteData []byte
type PacketLen uint32
type DoneFlag uint32

const (
	DONE = iota
	MORE
)

func NewDatagram() Datagram {
	return make([]byte, DATAGRAM_SIZE)
}

func CreateDatagram(offset OffsetVal, packet ByteData) Datagram {
	dg := NewDatagram()

	dg.Headers().SetOffset(offset)

	copy(dg.Packet(), packet)

	dataChecksum := CalculateChecksum(dg.Packet())
	dg.Headers().SetChecksum(dataChecksum)

	packetSize := PacketLen(math.Min(float64(len(packet)), PACKET_SIZE))
	dg.Headers().SetLength(packetSize)
	dg.Headers().SetDone(MORE)

	return dg
}

func (dg Datagram) Headers() Header {
	return Header(dg[:HEADER_SIZE])
}

func (dg Datagram) Packet() Packet {
	return Packet(dg[HEADER_SIZE:])
}

func (dg Datagram) Validate() bool {
	headerChecksum := dg.Headers().Checksum()
	dataChecksum := CalculateChecksum(dg.Packet())
	return bytes.Equal(headerChecksum, dataChecksum)
}

//////////////////////
/* Datagram Headers */
//////////////////////

func bytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func uint32ToBytes(n uint32) []byte {
	return (*[4]byte)(unsafe.Pointer(&n))[:]
}

type Header []byte

// Sequence ID of the current packet
func (h Header) Offset() OffsetVal {
	return OffsetVal(bytesToUint32(h[OFFSET_POINTER : OFFSET_POINTER+OFFSET_SIZE]))
}

// Copy in the offset into the file this data represents
func (h Header) SetOffset(offset OffsetVal) {
	copy(h[OFFSET_POINTER:OFFSET_POINTER+OFFSET_SIZE], uint32ToBytes(uint32(offset)))
}

// Checksum hash of Data (MD5)
func (h Header) Checksum() []byte {
	return h[CHECKSUM_POINTER : CHECKSUM_POINTER+CHECKSUM_SIZE]
}

// Copy in checksum data
func (h Header) SetChecksum(checksum []byte) {
	copy(h[CHECKSUM_POINTER:CHECKSUM_POINTER+CHECKSUM_SIZE], checksum)
}

// Checksum hash of Data (MD5)
func (h Header) Length() PacketLen {
	return PacketLen(bytesToUint32(h[LENGTH_POINTER : LENGTH_POINTER+LENGTH_SIZE]))
}

// Copy in checksum data
func (h Header) SetLength(length PacketLen) {
	copy(h[LENGTH_POINTER:LENGTH_POINTER+LENGTH_SIZE], uint32ToBytes(uint32(length)))
}

// Checksum hash of Data (MD5)
func (h Header) Done() DoneFlag {
	return DoneFlag(bytesToUint32(h[DONE_FLAG_POINTER : DONE_FLAG_POINTER+DONE_FLAG_SIZE]))
}

// Copy in checksum data
func (h Header) SetDone(done DoneFlag) {
	copy(h[DONE_FLAG_POINTER:DONE_FLAG_POINTER+DONE_FLAG_SIZE], uint32ToBytes(uint32(done)))
}

/////////////////////
/* Datagram Packet */
/////////////////////

// GZipped file contents
type Packet []byte
