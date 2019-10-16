package common

const (
	DATA_SIZE = DATAGRAM_LIMIT - HEADER_SIZE // Size of the data block (limit minus what's used by Flags and Headers)
)

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
