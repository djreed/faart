package main

import (
	"context"
	"net"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufferSize = 1500

func server(ctx context.Context, address string) (err error) {
	log.OUT.Println("Starting Server")

	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}
	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		datagram := packet.NewDatagram()
		_, _, err := conn.ReadFrom(datagram)
		if err != nil {
			doneChan <- err
			return
		}

		decompressedData, err := packet.Decompress(datagram.Packet())
		if err != nil {
			panic(err)
		}

		log.OUT.Printf("Server Data = '%s', Valid = %v\n", datagram.Packet(), datagram.Validate())
		log.OUT.Printf("Decompressed = '%s'\n", decompressedData)
	}()

	select {
	case <-ctx.Done():
		log.OUT.Panic("Server Cancelled from Context")
	case err = <-doneChan:
		log.OUT.Panic(err)
	}

	return
}
