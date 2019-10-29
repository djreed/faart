package shared

import (
	"io"
	"net"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
)

type (
	DataChannel chan packet.Datagram
	AckChannel  chan packet.Ack
	DoneChan    chan error
)

type DataMap map[packet.SeqID]packet.Datagram

const (
	ACCEPTED_IN_ORDER  = "ACCEPTED (in-order)"
	ACCEPTED_OUT_ORDER = "ACCEPTED (out-of-order)"
	IGNORED            = "IGNORED"

	ACK_BUFFER      = 1024
	DATAGRAM_BUFFER = 1024
)

func NewDataChan() DataChannel {
	return make(DataChannel, DATAGRAM_BUFFER)
}

func NewAckChan() AckChannel {
	return make(AckChannel, ACK_BUFFER)
}

func SendDatagram(conn io.Writer, datagram packet.Datagram) {
	if _, err := conn.Write(datagram); err != nil {
		log.OUT.Panic(err)
	}
}

func SendAck(conn *net.UDPConn, target *net.UDPAddr, ack packet.Ack) {
	if _, err := conn.WriteToUDP(ack, target); err != nil {
		log.OUT.Panic(err)
	}
}
