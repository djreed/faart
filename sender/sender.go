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
	acked     = make(shared.AckMap)

	dataChan  = shared.NewDataChan()
	ackChan   = shared.NewAckChan()
	completed = make(shared.ErrChannel, 1)
	state     = SENDING
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
				doneID := packet.SeqID(len(datagrams) + 1)
				finalDatagram := packet.CreateDatagram(doneID, packet.OffsetVal(0), []byte{})
				finalDatagram.Headers().SetDone(true)
				datagrams[doneID] = finalDatagram
				acked[doneID] = false
				go QueuePacketTimeout(shared.FIN_TIMEOUT, finalDatagram, dataChan)

				state = VALIDATING_END
				go func() {
					time.Sleep(shared.FIN_TIMEOUT_WAIT)
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
		datagram := packet.CreateDatagram(seqID, offset, compressedData[offset:])
		datagrams[seqID] = datagram
	}

	dataChan = make(chan packet.Datagram, len(datagrams))
	ackChan = make(chan packet.Ack, len(datagrams))

	go SendData(conn, dataChan)
	go QueueData(conn, dataChan)
	go QueueAcks(conn, ackChan)
	go HandleAcks(ackChan)
}

func QueuePacketTimeout(timeout time.Duration, datagram packet.Datagram, dataChan shared.DataChannel) {
	repetitions := 1
	for {
		if acked[datagram.Headers().Sequence()] {
			return
		} else {
			dataChan <- datagram
			time.Sleep(time.Duration(repetitions) * timeout)
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
				continue
			}
		}
	}
}

func QueueData(conn io.Writer, datachan shared.DataChannel) {
	for _, datagram := range datagrams {
		go QueuePacketTimeout(shared.SEND_PACKET_TIMEOUT, datagram, dataChan)
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
			acked[ack.Sequence()] = true
			if doneSending() {
				completed <- nil
			}
		}
	}
}

func doneSending() bool {
	return len(datagrams) == len(acked)
}
