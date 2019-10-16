package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
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
		for {
			bytesOfData, err := ioutil.ReadAll(reader)
			if err != nil {
				doneChan <- err
				return
			}

			log.OUT.Printf("data-read: bytes=%d -> %s\n", len(bytesOfData), bytesOfData)

			compressedData := new(bytes.Buffer)
			compressedWriter := gzip.NewWriter(compressedData)
			compressedWriter.Write(bytesOfData)
			compressedWriter.Close()

			datagram := packet.NewDatagram()
			copy(datagram.Data(), compressedData.Bytes())

			n, err := conn.Write(datagram.Packet)
			if err != nil {
				log.OUT.Panic(err)
			}

			log.OUT.Printf("data-written: bytes=%d\n", n)
		}
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
