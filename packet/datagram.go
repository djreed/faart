package packet

const (
	// Maximal UDP datagram size size
	DATAGRAM_SIZE = 1472 // 1500 - 20 (IP) - 8 (UDP) = 1472 bytes

	SEQUENCE_POINTER = 0 // Where in the slice does the sequence begin
	SEQUENCE_SIZE    = 4 // Sequence ID = 4 bytes

	OFFSET_POINTER = 4 // Where in the slice does the offset begin
	OFFSET_SIZE    = 4 // Offset into file = 4 bytes

	CHECKSUM_POINTER = 8  // Where in the slice does the Checksum begin
	CHECKSUM_SIZE    = 16 // MD5 Hash = 16 bytes

	HEADER_SIZE = 24 // 4 + 4 + 16 = 24 bytes

	DATA_SIZE = 1448 // 1472 - 24 = 1448
)

//////////////////////////////////////////////
/* The Grand High Poobah of Data Structures */
//////////////////////////////////////////////

// datagram data is just a byte slice
type Datagram struct {
	Packet []byte
}

func (dg *Datagram) Headers() Header {
	return Header(dg.Packet[:HEADER_SIZE])

}

func (dg *Datagram) Data() Data {
	return Data(dg.Packet[HEADER_SIZE:])
}

func NewDatagram() *Datagram {
	slice := make([]byte, DATAGRAM_SIZE)
	return &Datagram{Packet: slice}
}

//////////////////////
/* Datagram Headers */
//////////////////////

type Header []byte

// Sequence ID of the current packet
func (h Header) Sequence() []byte {
	return h[SEQUENCE_POINTER : SEQUENCE_POINTER+SEQUENCE_SIZE]
}

// Copy in the Sequence ID of this packet
func (h Header) SetSequence(sequence []byte) {
	copy(h[SEQUENCE_POINTER:SEQUENCE_POINTER+SEQUENCE_SIZE], sequence)
}

// Sequence ID of the current packet
func (h Header) Offset() []byte {
	return h[OFFSET_POINTER : OFFSET_POINTER+OFFSET_SIZE]
}

// Copy in the offset into the file this data represents
func (h Header) SetOffset(offset []byte) {
	copy(h[OFFSET_POINTER:OFFSET_POINTER+OFFSET_SIZE], offset)
}

// Checksum hash of Data (MD5)
func (h Header) Checksum() []byte {
	return h[CHECKSUM_POINTER : CHECKSUM_POINTER+CHECKSUM_SIZE]
}

// Copy in checksum data
func (h Header) SetChecksum(checksum []byte) {
	copy(h[CHECKSUM_POINTER:CHECKSUM_POINTER+CHECKSUM_SIZE], checksum)
}

///////////////////
/* Datagram Data */
///////////////////

// GZipped file contents
type Data []byte

// Sequence ID of the current packet
func (d Data) Body() []byte {
	return d
}
