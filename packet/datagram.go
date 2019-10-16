package packet

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

const (
	// Maximal UDP datagram size size
	DATAGRAM_SIZE = 1472 // 1500 - 20 (IP) - 8 (UDP) = 1472 bytes

	// Sequence ID for this datagram
	SEQUENCE_POINTER = 0
	SEQUENCE_SIZE    = 4 // UInt32 = 4 bytes

	// Offset in the destination file
	OFFSET_POINTER = SEQUENCE_POINTER + SEQUENCE_SIZE
	OFFSET_SIZE    = 4 // UInt32 = 4 bytes

	// Checksum of the Packet's data
	CHECKSUM_POINTER = OFFSET_POINTER + OFFSET_SIZE
	CHECKSUM_SIZE    = 20 // SHA-1 Hash = 20 bytes

	// Length of the Packet
	LENGTH_POINTER = CHECKSUM_POINTER + CHECKSUM_SIZE
	LENGTH_SIZE    = 4 // UInt32 = 4 bytes

	HEADER_SIZE = LENGTH_POINTER + LENGTH_SIZE // 20 + 4 + 4 = 28

	PACKET_SIZE = DATAGRAM_SIZE - HEADER_SIZE // 1472 - 28 = 1444
)

//////////////////////////////////////////////
/* The Grand High Poobah of Data Structures */
//////////////////////////////////////////////

// datagram data is just a byte slice
type Datagram []byte

func NewDatagram() Datagram {
	return make([]byte, DATAGRAM_SIZE)
}

func CreateDatagram(seq uint32, offset uint32, data []byte) Datagram {
	dg := NewDatagram()

	dg.Headers().SetSequence(seq)
	dg.Headers().SetOffset(offset)

	copy(dg.Packet(), data)

	dataChecksum := CalculateChecksum(dg.Packet())
	dg.Headers().SetChecksum(dataChecksum)
	dg.Headers().SetLength(uint32(len(dg.Packet())))

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
func (h Header) Sequence() uint32 {
	return bytesToUint32(h[SEQUENCE_POINTER : SEQUENCE_POINTER+SEQUENCE_SIZE])
}

// Copy in the Sequence ID of this packet
func (h Header) SetSequence(sequence uint32) {
	copy(h[SEQUENCE_POINTER:SEQUENCE_POINTER+SEQUENCE_SIZE], uint32ToBytes(sequence))
}

// Sequence ID of the current packet
func (h Header) Offset() uint32 {
	return bytesToUint32(h[OFFSET_POINTER : OFFSET_POINTER+OFFSET_SIZE])
}

// Copy in the offset into the file this data represents
func (h Header) SetOffset(offset uint32) {
	copy(h[OFFSET_POINTER:OFFSET_POINTER+OFFSET_SIZE], uint32ToBytes(offset))
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
func (h Header) Length() uint32 {
	return bytesToUint32(h[LENGTH_POINTER : LENGTH_POINTER+LENGTH_SIZE])
}

// Copy in checksum data
func (h Header) SetLength(length uint32) {
	copy(h[LENGTH_POINTER:LENGTH_POINTER+LENGTH_SIZE], uint32ToBytes(length))
}

/////////////////////
/* Datagram Packet */
/////////////////////

// GZipped file contents
type Packet []byte
