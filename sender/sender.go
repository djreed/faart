package main

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
	"github.com/djreed/faart/shared"
)

var (
	queueTimeout = time.Duration(1 * time.Second)

	datagrams = make(shared.DataMap)

	dataChan = shared.NewDataChan()
	ackChan  = shared.NewAckChan()

	doneSendData = make(shared.DoneChan, 1)
	doneRecvAck  = make(shared.DoneChan, 1)

	completed = make(shared.DoneChan, 1)
)

func sender(ctx context.Context, address string, reader io.Reader) error {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return err
	}

	defer conn.Close()

	go handleConn(conn, reader)

	select {
	case err = <-completed:
		log.ERR.Printf("[completed]\n")
		return err
	}
}

func handleConn(conn *net.UDPConn, reader io.Reader) {
	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		completed <- err
		return
	}

	compressedData, err := packet.Compress(decompressedData)
	if err != nil {
		panic(err)
	}

	// Populate packet map
	for sequence := 0; sequence*packet.PACKET_SIZE < len(compressedData); sequence++ {
		seqID := packet.SeqID(sequence)
		offset := packet.OffsetVal(sequence * packet.PACKET_SIZE)
		datagram := packet.CreateDatagram(seqID, offset, compressedData[offset:])
		datagrams[seqID] = datagram
	}

	dataChan = make(chan packet.Datagram, len(datagrams))
	ackChan = make(chan packet.Ack, len(datagrams))

	go QueueData(dataChan)
	go SendData(conn, dataChan)
	go QueueAcks(conn, ackChan)
	go HandleAcks(ackChan)
}

func QueueData(dataChan shared.DataChannel) {
	for {
		for _, datagram := range datagrams {
			dataChan <- datagram
		}
		time.Sleep(queueTimeout)
	}
}

func SendData(conn io.Writer, dataChan shared.DataChannel) {
	for {
		select {
		case datagram := <-dataChan:
			shared.SendDatagram(conn, datagram)
			log.ERR.Printf("[send data] %d (%d)\n", datagram.Headers().Offset(), datagram.Headers().Length())
			continue
		case <-doneSendData:
			return
		}
	}
}

func QueueAcks(conn net.PacketConn, ackChan shared.AckChannel) {
	for {
		ack := packet.NewAck()
		read, _, _ := conn.ReadFrom(ack)
		if read > 0 {
			ackChan <- ack
		}
	}
}

func HandleAcks(ackChan shared.AckChannel) {
	for {
		select {
		case ack := <-ackChan:
			log.ERR.Printf("[recv ack] %d\n", ack.Offset())
			delete(datagrams, ack.Sequence())
			if len(datagrams) == 0 {
				doneSendData <- nil
				doneRecvAck <- nil
			}
		case <-doneRecvAck:
			completed <- nil
			return
		}
	}
}
