package packet

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

const (
	// Maximal UDP datagram size size
	DATAGRAM_SIZE = 1472 // 1500 - 20 (IP) - 8 (UDP) = 1472 bytes

	SEQUENCE_POINTER = 0 // Where in the slice does the sequence begin
	SEQUENCE_SIZE    = 4 // Sequence ID = 4 bytes

	OFFSET_POINTER = SEQUENCE_POINTER + SEQUENCE_SIZE // Where in the slice does the offset begin
	OFFSET_SIZE    = 4                                // Offset into file = 4 bytes

	CHECKSUM_POINTER = OFFSET_POINTER + OFFSET_SIZE // Where in the slice does the Checksum begin
	CHECKSUM_SIZE    = 20                           // SHA-1 Hash = 20 bytes

	HEADER_SIZE = CHECKSUM_POINTER + CHECKSUM_SIZE // 20 + 4 + 4 = 28

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

/* BEGIN THE RITUAL
___________________6666666___________________
____________66666__________66666_____________
_________6666___________________666__________
_______666__6____________________6_666_______
_____666_____66_______________666____66______
____66_______66666_________66666______666____
___66_________6___66_____66___66_______666___
__66__________66____6666_____66_________666__
_666___________66__666_66___66___________66__
_66____________6666_______6666___________666_
_66___________6666_________6666__________666_
_66________666_________________666_______666_
_66_____666______66_______66______666____666_
_666__666666666666666666666666666666666__66__
__66_______________6____66______________666__
___66______________66___66_____________666___
____66______________6__66_____________666____
_______666___________666___________666_______
_________6666_________6_________666__________
____________66666_____6____66666_____________
___________________6666666___________________
*/

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

/////////////////////
/* Datagram Packet */
/////////////////////

// GZipped file contents
type Packet []byte
