package packet

import (
	"bytes"
	"encoding/binary"
	"math"
	"net"
	"unsafe"
)

type AddressedDatagram struct {
	Datagram Datagram
	Addr     *net.UDPAddr
}

const (
	// Maximal UDP datagram size size
	DATAGRAM_SIZE = 1472 // MTU - IP - UDP = 1500 - 20 - 8 = 1472 bytes

	// Sequence ID
	SEQUENCE_POINTER = 0
	SEQUENCE_SIZE    = 4

	// File Offset
	OFFSET_POINTER = SEQUENCE_POINTER + SEQUENCE_SIZE
	OFFSET_SIZE    = 4

	// Checksum of the Packet's data
	CHECKSUM_POINTER = OFFSET_POINTER + OFFSET_SIZE
	CHECKSUM_SIZE    = 20

	// Length of the Packet
	LENGTH_POINTER = CHECKSUM_POINTER + CHECKSUM_SIZE
	LENGTH_SIZE    = 4

	// // Whether you're completely done with the data being send
	DONE_FLAG_POINTER = LENGTH_POINTER + LENGTH_SIZE
	DONE_FLAG_SIZE    = 1

	// Packet count of the total file
	COUNT_POINTER = DONE_FLAG_POINTER + DONE_FLAG_SIZE
	COUNT_SIZE    = 4

	// HEADER_SIZE = LENGTH_POINTER + LENGTH_SIZE
	HEADER_SIZE = COUNT_POINTER + COUNT_SIZE

	PACKET_SIZE = DATAGRAM_SIZE - HEADER_SIZE
)

//////////////////////////////////////////////
/* The Grand High Poobah of Data Structures */
//////////////////////////////////////////////

// datagram data is just a byte slice
type Datagram []byte

type SeqID uint32
type OffsetVal uint32
type ByteData []byte
type PacketLen uint32
type DoneFlag bool
type PacketCount uint32

func NewDatagram() Datagram {
	return make([]byte, DATAGRAM_SIZE)
}

func CreateDatagram(sequence SeqID, offset OffsetVal, packet ByteData, packetCount PacketCount) Datagram {
	dg := NewDatagram()

	dg.Headers().SetSequence(sequence)
	dg.Headers().SetOffset(offset)

	copy(dg.Packet(), packet[:])

	dataChecksum := CalculateChecksum(dg.Packet())
	dg.Headers().SetChecksum(dataChecksum)

	packetSize := PacketLen(math.Min(float64(len(packet)), PACKET_SIZE))
	dg.Headers().SetLength(packetSize)
	dg.Headers().SetDone(false)
	dg.Headers().SetCount(packetCount)

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

func bytesToBool(b []byte) bool {
	return bytes.Contains(b, []byte{1})
}

func boolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

type Header []byte

// Sequence ID of the current packet
func (h Header) Sequence() SeqID {
	return SeqID(bytesToUint32(h[SEQUENCE_POINTER : SEQUENCE_POINTER+SEQUENCE_SIZE]))
}

func (h Header) SetSequence(seq SeqID) {
	copy(h[SEQUENCE_POINTER:SEQUENCE_POINTER+SEQUENCE_SIZE], uint32ToBytes(uint32(seq)))
}

// Sequence ID of the current packet
func (h Header) Offset() OffsetVal {
	return OffsetVal(bytesToUint32(h[OFFSET_POINTER : OFFSET_POINTER+OFFSET_SIZE]))
}

func (h Header) SetOffset(offset OffsetVal) {
	copy(h[OFFSET_POINTER:OFFSET_POINTER+OFFSET_SIZE], uint32ToBytes(uint32(offset)))
}

// Checksum hash of Data (MD5)
func (h Header) Checksum() []byte {
	return h[CHECKSUM_POINTER : CHECKSUM_POINTER+CHECKSUM_SIZE]
}
func (h Header) SetChecksum(checksum []byte) {
	copy(h[CHECKSUM_POINTER:CHECKSUM_POINTER+CHECKSUM_SIZE], checksum)
}

// Length of packet contents
func (h Header) Length() PacketLen {
	return PacketLen(bytesToUint32(h[LENGTH_POINTER : LENGTH_POINTER+LENGTH_SIZE]))
}
func (h Header) SetLength(length PacketLen) {
	copy(h[LENGTH_POINTER:LENGTH_POINTER+LENGTH_SIZE], uint32ToBytes(uint32(length)))
}

// Whether this is the last packet
func (h Header) Done() DoneFlag {
	return DoneFlag(bytesToBool((h[DONE_FLAG_POINTER : DONE_FLAG_POINTER+DONE_FLAG_SIZE])))
}
func (h Header) SetDone(done DoneFlag) {
	copy(h[DONE_FLAG_POINTER:DONE_FLAG_POINTER+DONE_FLAG_SIZE], boolToBytes(bool(done)))
}

// Length of packet contents
func (h Header) Count() PacketCount {
	return PacketCount(bytesToUint32(h[COUNT_POINTER : COUNT_POINTER+COUNT_SIZE]))
}
func (h Header) SetCount(count PacketCount) {
	copy(h[COUNT_POINTER:COUNT_POINTER+COUNT_SIZE], uint32ToBytes(uint32(count)))
}

/////////////////////
/* Datagram Packet */
/////////////////////

// Packet bytes
type Packet []byte
