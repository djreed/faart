package main

import (
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"time"

	"github.com/djreed/faart/log"
	"github.com/djreed/faart/packet"
)

var timeout = time.Duration(10 * time.Second)

func sender(ctx context.Context, address string, reader io.Reader) (err error) {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return
	}

	defer conn.Close()

	done := make(chan error, 1)

	go handleConn(done, conn, reader)

	select {
	case <-ctx.Done():
		return
	case err = <-done:
		log.OUT.Printf("[completed]\n")
		return
	}
}

func handleConn(done chan error, conn io.Writer, reader io.Reader) {
	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		done <- err
		return
	}

	compressedData, err := packet.Compress(decompressedData)
	if err != nil {
		panic(err)
	}

	seqBase := rand.Int()

	for order := 0; (order * packet.PACKET_SIZE) < len(compressedData); order++ {
		offset := uint32(order * packet.PACKET_SIZE)
		packetSequence := uint32(seqBase + order)

		datagram := packet.CreateDatagram(packetSequence, offset, compressedData)

		_, err = conn.Write(datagram)
		if err != nil {
			log.OUT.Panic(err)
		}

		log.ERR.Printf("[send data] %d (%d)\n", datagram.Headers().Offset(), len(datagram.Packet()))
	}

	done <- nil
}
