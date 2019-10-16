package main

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
)

var timeout = time.Duration(10 * time.Second)

func client(ctx context.Context, address string, reader io.Reader) (err error) {
	log.OUT.Println("Starting Client")

	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		decompressedData := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
		if err != nil {
			doneChan <- err
			return
		}

		compressedData, err := packet.Compress(decompressedData)
		if err != nil {
			panic(err)
		}

		datagram := packet.CreateDatagram(1, 5, compressedData)

		_, err = conn.Write(datagram)
		if err != nil {
			log.OUT.Panic(err)
		}

		log.OUT.Printf("Original = '%s'\n", decompressedData)
		log.OUT.Printf("Client Data = '%s', Valid = %v\n", compressedData, datagram.Validate())
	}()

	select {
	case <-ctx.Done():
		log.OUT.Panic("Client Cancelled from Context")
	case err = <-doneChan:
		os.Exit(0)
		return
	}

	return
}
