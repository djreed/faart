package main

import (
	"context"
	"net"
	"strings"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufferSize = 1500

func receiver(ctx context.Context) (err error) {
	conn, err := net.ListenPacket("udp4", "")
	if err != nil {
		return
	}
	defer conn.Close()

	splitAddr := strings.Split(conn.LocalAddr().String(), ":")
	log.ERR.Printf("[bound] %s", splitAddr[len(splitAddr)-1])

	done := make(chan error, 1)
	data := make(chan []byte)
	totalData := make([]byte, 0)

	go handleConn(done, conn, data)

	for {
		select {
		case byteData := <-data:
			totalData = append(totalData, byteData...)
		case <-ctx.Done():
			return
		case err = <-done:
			log.OUT.Printf("[completed]\n")
			return
		}
	}
}

func handleConn(done chan error, conn net.PacketConn, data chan []byte) {
	for {
		datagram := packet.NewDatagram()
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

		data <- decompressedData

		log.ERR.Printf("[recv data] %d (%d) ACCEPTED (TODO ORDER)\n", datagram.Headers().Offset(), len(datagram.Packet()))
	}
}
