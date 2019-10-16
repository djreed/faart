package common

const (
	DATAGRAM_LIMIT = 1472 // 1500 - 20 (IP) - 8 (UDP) = 1472 bytes
)

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
