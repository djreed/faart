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
	readTimeout = time.Duration(5000 * time.Millisecond)

	datagrams = make(shared.DataMap)
	dataChan  = shared.NewAddressedDataChan()
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

	go ReceiveDatagrams(conn, dataChan)
	go HandleDatagrams(conn, dataChan, finalChan)
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

func ReceiveDatagrams(conn *net.UDPConn, dataChan shared.AddressedDataChannel) {
	for {
		datagram := packet.NewDatagram()
		read, retAddr, err := conn.ReadFromUDP(datagram)
		log.ERR.Printf("Received Datagram: %d %v\n", read, err)
		if read > 0 && err == nil {
			addressedDatagram := packet.AddressedDatagram{Addr: retAddr, Datagram: datagram}
			dataChan <- addressedDatagram
		}
	}
}

func HandleDatagrams(conn *net.UDPConn, dataChan shared.AddressedDataChannel, doneChan shared.ErrChannel) {
	var lastPacketReceived time.Time
	for {
		select {
		case addressedDatagram := <-dataChan:
			lastPacketReceived = time.Now()
			needAck, finalPacket := AcceptDatagram(addressedDatagram.Datagram)
			ack := packet.CreateAck(addressedDatagram.Datagram)
			ackPacket := packet.AddressedAck{Addr: addressedDatagram.Addr, Ack: ack}
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
			continue

		default:
			if !lastPacketReceived.IsZero() && time.Since(lastPacketReceived) > readTimeout {
				doneChan <- nil
				return
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
