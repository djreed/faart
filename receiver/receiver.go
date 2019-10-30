package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
	"github.com/djreed/faart/shared"
)

var (
	datagrams = make(shared.DataMap)
	ackChan   = make(chan AckPacket, shared.ACK_BUFFER)
	doneData  = make(shared.DoneChan, 1)
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

	go HandleDatagrams(conn, doneData)
	go SendAcks(conn, ackChan)

	for {
		select {
		case <-doneData:
			log.ERR.Printf("[completed]\n")
			byteData := flattenMapData(datagrams)
			printData(byteData)
			return nil
		}
	}
}

func HandleDatagrams(conn *net.UDPConn, doneChan shared.DoneChan) {
	for {
		datagram := packet.NewDatagram()
		read, retAddr, _ := conn.ReadFromUDP(datagram)
		if read > 0 {
			needAck, final := AcceptDatagram(datagram)
			if needAck {
				ackChan <- AckPacket{addr: retAddr, ack: packet.CreateAck(datagram)}
			}
			if final {
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
	return !existing, false
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
