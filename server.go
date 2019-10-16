package main

import (
	"bytes"
	"compress/gzip"
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

	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}
	defer pc.Close()

	doneChan := make(chan error, 1)

	go func() {
		for {
			datagram := packet.NewDatagram()
			n, addr, err := pc.ReadFrom(datagram.Packet)
			if err != nil {
				doneChan <- err
				return
			}

			log.OUT.Printf("data-received: bytes=%d from=%s\n", n, addr.String())

			gzReader, err := gzip.NewReader(bytes.NewBuffer(datagram.Data()))
			if err != nil {
				log.OUT.Panic(err)
			}

			bytesOfData := make([]byte, packet.DATA_SIZE)
			gzReader.Read(bytesOfData)

			log.OUT.Printf("data-decoded: bytes=%d from=%s -> %s\n", n, addr.String(), bytesOfData)
		}
	}()

	select {
	case <-ctx.Done():
		log.OUT.Panic("Server Cancelled from Context")
	case err = <-doneChan:
		log.OUT.Panic(err)
	}

	return
}
