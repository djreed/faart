package main

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufferSize = 1500

type dataPacket struct {
	seq    uint32
	offset uint32
	data   []byte
}

func receiver(ctx context.Context) (err error) {
	conn, err := net.ListenPacket("udp4", "")
	if err != nil {
		return
	}
	defer conn.Close()

	splitAddr := strings.Split(conn.LocalAddr().String(), ":")
	log.ERR.Printf("[bound] %s", splitAddr[len(splitAddr)-1])

	done := make(chan error, 1)
	data := make(chan dataPacket)
	totalData := make([]byte, 0)

	go handleConn(done, conn, data)

	for {
		select {
		case packetData := <-data:
			totalData = addDataToExisting(totalData, packetData)
		case <-ctx.Done():
			return
		case err = <-done:
			log.OUT.Printf("[completed]\n")
			log.OUT.Print(string(totalData))
			return
		}
	}
}

func handleConn(done chan error, conn net.PacketConn, data chan dataPacket) {
	for {
		datagram := packet.NewDatagram()

		conn.SetDeadline(time.Now().Add(time.Duration(5 * time.Second)))
		_, _, err := conn.ReadFrom(datagram)
		if err != nil {
			done <- err
			return
		}

		if !datagram.Validate() {
			log.ERR.Printf("[recv corrupt packet]\n")
			continue
		}

		decompressedData, err := packet.Decompress(datagram.Packet())
		if err != nil {
			panic(err)
		}

		pack := dataPacket{
			seq:    datagram.Headers().Sequence(),
			offset: datagram.Headers().Offset(),
			data:   decompressedData,
		}
		data <- pack

		log.ERR.Printf("[recv data] %d (%d) ACCEPTED (TODO ORDER)\n", datagram.Headers().Offset(), len(datagram.Packet()))
	}
}

func addDataToExisting(existing []byte, data dataPacket) []byte {
	dataEnd := int(data.offset) + len(data.data) - 1
	missing := dataEnd - len(existing)
	if missing > 0 {
		existing = append(existing, make([]byte, missing)...)
	}
	copy(existing[data.offset:dataEnd], data.data)
	return existing
}
