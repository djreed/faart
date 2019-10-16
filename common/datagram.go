package common

const (
	DATAGRAM_LIMIT = 1472 // 1500 - 20 (IP) - 8 (UDP) = 1472 bytes

	SEQUENCE_SIZE = 4                                           // Sequence ID -> 32 bit integer = 4 bytes
	OFFSET_SIZE   = 4                                           // Offset into file -> 32 bit integer = 4 bytes
	CHECKSUM_SIZE = 16                                          // MD5 Hash -> 128 bits = 16 bytes
	HEADER_SIZE   = SEQUENCE_SIZE + OFFSET_SIZE + CHECKSUM_SIZE // Sum of all the header sizes = 24 bytes

	// DATA_SIZE = 1472 - 24 = 1448
	DATA_SIZE = DATAGRAM_LIMIT - HEADER_SIZE // Size of the data block (limit minus what's used by Flags and Headers)
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
	var sequenceID [SEQUENCE_SIZE]byte
	copy(sequenceID[:], h.headers[0:SEQUENCE_SIZE])
	return sequenceID
}

// Checksum hash of Data (MD5)
func (h *Header) Checksum() [CHECKSUM_SIZE]byte {
	var checksum [CHECKSUM_SIZE]byte
	copy(checksum[:], h.headers[SEQUENCE_SIZE:CHECKSUM_SIZE])
	return checksum
}

///////////////////
/* Datagram Data */
///////////////////

// GZipped file contents
type Data struct {
	data [DATA_SIZE]byte
}

// Sequence ID of the current packet
func (d *Data) Body() [DATA_SIZE]byte {
	var data [DATA_SIZE]byte
	copy(data[:], d.data[:])
	return data
}

//////////////////////////////////////////////
/* The Grand High Poobah of Data Structures */
//////////////////////////////////////////////

// datagram data is just a byte slice
type Datagram struct {
	packet [DATAGRAM_LIMIT]byte
}

func (dg *Datagram) Headers() *Header {
	var headers [HEADER_SIZE]byte
	copy(headers[:], dg.packet[0:HEADER_SIZE])
	return &Header{headers: headers}
}

func (dg *Datagram) Data() *Data {
	var data [DATA_SIZE]byte
	copy(data[:], dg.packet[HEADER_SIZE:])
	return &Data{data: data}
}
