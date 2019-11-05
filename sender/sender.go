package main

import (
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
	"github.com/djreed/faart/shared"
)

const (
	SENDING        = 0
	VALIDATING_END = 1
)

var (
	// For now, don't wait on each send
	datagrams = make(shared.DataMap)
	dataChan  = shared.NewDataChan()
	ackChan   = shared.NewAckChan()
	completed = make(shared.ErrChannel, 1)
	state     = SENDING
	maxCount  = -1
)

func sender(address string, reader io.Reader) error {
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

	for {
		select {
		case err = <-completed:
			switch state {
			case SENDING:
				doneID := packet.SeqID(maxCount + 1)
				finalDatagram := packet.CreateDatagram(doneID, packet.OffsetVal(0), []byte{})
				finalDatagram.Headers().SetDone(true)
				datagrams[doneID] = finalDatagram
				dataChan <- finalDatagram
				state = VALIDATING_END
				go func() {
					time.Sleep(shared.SEND_FIN_TIMEOUT)
					completed <- nil
				}()
				continue

			case VALIDATING_END:
				log.ERR.Printf("[completed]\n")
				return err
			}
		}
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
		doneDatagram := packet.CreateDatagram(seqID, offset, compressedData[offset:])
		datagrams[seqID] = doneDatagram
		dataChan <- doneDatagram
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
		if state == SENDING {
			time.Sleep(shared.QUEUE_DATA_TIMEOUT)
		} else if state == VALIDATING_END {
			time.Sleep(shared.QUEUE_FIN_TIMEOUT)
		}
	}
}

func SendData(conn io.Writer, dataChan shared.DataChannel) {
	for {
		select {
		case datagram := <-dataChan:
			if err := shared.SendDatagram(conn, datagram); err != nil {
				// TODO Handle errors
				completed <- nil
				return
			} else {
				log.ERR.Printf("[send data] %d (%d)\n", datagram.Headers().Offset(), datagram.Headers().Length())
				time.Sleep(shared.SEND_PACKET_TIMEOUT)
				continue
			}
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
			if doneSending() {
				completed <- nil
			}
		}
	}
}

func doneSending() bool {
	return len(datagrams) == 0
}
