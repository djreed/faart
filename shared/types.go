package shared

import (
	"io"
	"net"

	"github.com/djreed/faart/packet"
)

type (
	DataChannel         chan packet.Datagram
	AckChannel          chan packet.Ack
	AddressedAckChannel chan packet.AddressedAck
	ErrChannel          chan error
)

type DataMap map[packet.SeqID]packet.Datagram

const (
	ACCEPTED_IN_ORDER  = "ACCEPTED (in-order)"
	ACCEPTED_OUT_ORDER = "ACCEPTED (out-of-order)"
	IGNORED            = "IGNORED"

	ACK_BUFFER      = 4092
	DATAGRAM_BUFFER = 4092
)

func NewErrChan() ErrChannel {
	return make(ErrChannel, 0)
}

func NewDataChan() DataChannel {
	return make(DataChannel, DATAGRAM_BUFFER)
}

func NewAckChan() AckChannel {
	return make(AckChannel, ACK_BUFFER)
}

func NewAddressedAckChan() AddressedAckChannel {
	return make(AddressedAckChannel, ACK_BUFFER)
}

func SendDatagram(conn io.Writer, datagram packet.Datagram) error {
	_, err := conn.Write(datagram)
	return err
}

func SendAck(conn *net.UDPConn, target *net.UDPAddr, ack packet.Ack) error {
	_, err := conn.WriteToUDP(ack, target)
	return err
}
