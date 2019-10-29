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
const DATA_TIMEOUT = time.Duration(3 * time.Second)

var (
	datagrams = make(shared.DataMap)

	dataChan = shared.NewDataChan()
	ackChan  = make(chan AckPacket, shared.ACK_BUFFER)

	doneRecvData = make(shared.DoneChan, 1)
	doneSendAck  = make(shared.DoneChan, 1)

	completed = make(shared.DoneChan, 1)

	lastAck time.Time

	compressedSize = 0
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

	go HandleDatagrams(conn, dataChan)
	go SendAcks(conn, ackChan)

	timeout := time.After(STARTUP_WAIT)
	for {
		select {
		case <-dataChan:
			timeout = time.After(DATA_TIMEOUT)
			continue
		case <-timeout:
			log.ERR.Printf("[completed]\n")
			byteData := flattenMapData(datagrams)
			printData(byteData)
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

const RECV_TEMPLATE = "[recv data] %d (%d) %s\n"

func AcceptDatagram(datagram packet.Datagram) bool {
	_, existing := datagrams[datagram.Headers().Sequence()]
	if !existing {
		if !datagram.Validate() {
			log.ERR.Printf("[recv corrupt packet]\n")
			return false
		}
		datagrams[datagram.Headers().Sequence()] = datagram
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

func flattenMapData(datagramMap shared.DataMap) []byte {
	finalPacket := datagramMap[packet.SeqID(len(datagramMap)-1)]
	fileSize := uint32(finalPacket.Headers().Offset()) + uint32(finalPacket.Headers().Length())
	data := make([]byte, fileSize)
	for i := packet.SeqID(0); i < packet.SeqID(len(datagramMap)); i++ {
		offset := uint32(datagramMap[i].Headers().Offset())
		length := uint32(datagramMap[i].Headers().Length())
		datagramPacket := datagramMap[i].Packet()[0:length]
		copy(data[offset:], datagramPacket[:])
	}
	return data
}

func printData(data []byte) {
	decompressedData, err := packet.Decompress(data)
	if err != nil {
		panic(err)
	}

	log.LOG.Print(string(decompressedData)) // TODO reassemble data
}
