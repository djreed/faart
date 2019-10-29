package main

import (
	"net"
	"strings"
	"time"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
	"github.com/djreed/faart/shared"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufferSize = 1500
const STARTUP_WAIT = 30 * time.Second
const TIMEOUT_MS = 3000

const RECV_TEMPLATE = "[recv data] %d (%d) %s\n"

var (
	datagrams = make(shared.DataMap)

	dataChan = shared.NewDataChan()
	ackChan  = make(chan AckPacket, shared.ACK_BUFFER)

	doneRecvData = make(shared.DoneChan, 1)
	doneSendAck  = make(shared.DoneChan, 1)

	completed = make(shared.DoneChan, 1)

	lastAck time.Time
)

type AckPacket struct {
	addr *net.UDPAddr
	ack  packet.Ack
}

type AckPacketChan chan AckPacket

func receiver() error {
	localAddr := new(net.UDPAddr)
	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	splitAddr := strings.Split(conn.LocalAddr().String(), ":")
	log.ERR.Printf("[bound] %s", splitAddr[len(splitAddr)-1])
	totalData := make([]byte, 0)

	go HandleDatagrams(conn, dataChan)
	go SendAcks(conn, ackChan)

	timeout := time.After(STARTUP_WAIT)
	for {
		select {
		case datagram := <-dataChan:
			totalData = appendDatagram(totalData, datagram)
			timeout = time.After(time.Duration(TIMEOUT_MS * time.Millisecond))
			continue
		case <-timeout:
			log.ERR.Printf("[completed]\n")
			printData(totalData)
			return nil
		}
	}
}

func HandleDatagrams(conn *net.UDPConn, datachan shared.DataChannel) {
	for {
		datagram := packet.NewDatagram()
		read, retAddr, _ := conn.ReadFromUDP(datagram) // TODO error handling
		if read > 0 {
			success := AcceptDatagram(datagram)
			if success {
				datachan <- datagram
				ackChan <- AckPacket{addr: retAddr, ack: packet.CreateAck(datagram)}
			}
		}
	}
}

func AcceptDatagram(datagram packet.Datagram) bool {
	_, existing := datagrams[datagram.Headers().Offset()]
	if !existing {
		if !datagram.Validate() {
			log.ERR.Printf("[recv corrupt packet]\n")
			return false
		}
		datagrams[datagram.Headers().Offset()] = datagram
		log.ERR.Printf(RECV_TEMPLATE, datagram.Headers().Offset(), datagram.Headers().Length(), shared.ACCEPTED_OUT_ORDER)
	} else {
		log.ERR.Printf(RECV_TEMPLATE, datagram.Headers().Offset(), datagram.Headers().Length(), shared.IGNORED)
	}
	return !existing
}

func SendAcks(conn *net.UDPConn, ackChan AckPacketChan) {
	for {
		select {
		case ack := <-ackChan:
			shared.SendAck(conn, ack.addr, ack.ack)
			continue
		}
	}
}

func appendDatagram(existing []byte, datagram packet.Datagram) []byte {
	dataEnd := uint32(datagram.Headers().Offset()) + uint32(datagram.Headers().Length())
	missingLen := dataEnd - uint32(cap(existing))
	if missingLen > 0 {
		// Pad out existing bytes with 0s
		existing = append(existing, make([]byte, missingLen)...)
		copy(existing[datagram.Headers().Offset():dataEnd], datagram.Packet()[0:datagram.Headers().Length()])
	}
	return existing
}

func printData(data []byte) {
	decompressedData, err := packet.Decompress(data)
	if err != nil {
		panic(err)
	}
	log.OUT.Print(string(decompressedData)) // TODO reassemble data
}
