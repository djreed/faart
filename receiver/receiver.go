package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
	"github.com/djreed/faart/shared"
)

var (
	readTimeout = time.Duration(2000 * time.Millisecond)

	datagrams = make(shared.DataMap)
	dataChan  = shared.NewDataChan()
	ackChan   = shared.NewAddressedAckChan()
	finalChan = shared.NewErrChan()
)

func receiver() error {
	localAddr := new(net.UDPAddr)
	conn, err := net.ListenUDP("udp4", localAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	splitAddr := strings.Split(conn.LocalAddr().String(), ":")
	log.ERR.Printf("[bound] %s", splitAddr[len(splitAddr)-1])

	go HandleDatagrams(conn, finalChan)
	go SendAcks(conn, ackChan)

	for {
		select {
		case <-finalChan:
			log.ERR.Printf("[completed]\n")
			byteData := flattenMapData(datagrams)
			printData(byteData)
			return nil
		}
	}
}

func HandleDatagrams(conn *net.UDPConn, doneChan shared.ErrChannel) {
	var lastPacketReceived time.Time
	for {
		datagram := packet.NewDatagram()
		read, retAddr, _ := conn.ReadFromUDP(datagram)
		if read > 0 {
			lastPacketReceived = time.Now()
			needAck, finalPacket := AcceptDatagram(datagram)
			ack := packet.CreateAck(datagram)
			ackPacket := packet.AddressedAck{Addr: retAddr, Ack: ack}
			if needAck {
				ackChan <- ackPacket
			}
			if finalPacket {
				// TODO: what if the final ack doesn't make it
				// TODO: What if we just send a ton of ACKs
				ackChan <- ackPacket
				ackChan <- ackPacket
				ackChan <- ackPacket
				ackChan <- ackPacket
				ackChan <- ackPacket
				ackChan <- ackPacket
				doneChan <- nil
			}
		} else {
			if time.Since(lastPacketReceived) > readTimeout {
				doneChan <- nil
			}
		}
	}
}

const RECV_TEMPLATE = "[recv data] %d (%d) %s\n"

func AcceptDatagram(datagram packet.Datagram) (bool, bool) {
	if datagram.Headers().Done() {
		return true, true
	}

	_, existing := datagrams[datagram.Headers().Sequence()]
	if !existing {
		if !datagram.Validate() {
			log.ERR.Printf("[recv corrupt packet]\n")
			return false, false
		}
		datagrams[datagram.Headers().Sequence()] = datagram
		log.ERR.Printf(RECV_TEMPLATE, datagram.Headers().Offset(), datagram.Headers().Length(), shared.ACCEPTED_OUT_ORDER)
	} else {
		log.ERR.Printf(RECV_TEMPLATE, datagram.Headers().Offset(), datagram.Headers().Length(), shared.IGNORED)
	}
	return true, false
}

func SendAcks(conn *net.UDPConn, ackChan shared.AddressedAckChannel) {
	for {
		select {
		case ack := <-ackChan:
			shared.SendAck(conn, ack.Addr, ack.Ack)
			continue
		}
	}
}

func flattenMapData(datagramMap shared.DataMap) []byte {
	packetCount := len(datagramMap)
	finalPacket := datagramMap[packet.SeqID(packetCount-1)]
	fileSize := uint32(finalPacket.Headers().Offset()) + uint32(finalPacket.Headers().Length())
	data := make([]byte, fileSize)
	for i := packet.SeqID(0); i < packet.SeqID(len(datagramMap)); i++ {
		offset := uint32(datagramMap[i].Headers().Offset())
		length := uint32(datagramMap[i].Headers().Length())
		datagramPacket := datagramMap[i].Packet()[0:length]
		copy(data[offset:offset+length], datagramPacket[:])
	}
	return data
}

func printData(data []byte) {
	decompressedData, err := packet.Decompress(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", decompressedData)
}
