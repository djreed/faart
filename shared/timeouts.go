package shared

import "time"

var (
	// How long to wait between packet sends
	SEND_PACKET_WAIT = time.Duration(1 * time.Microsecond)

	// How long to wait before re-queueing packets
	SEND_PACKET_TIMEOUT = time.Duration(200 * time.Millisecond)

	// How long to wait before re-queueing the FIN
	FIN_TIMEOUT = time.Duration(100 * time.Millisecond)

	// How long to wait before completing if sender has not received
	// an ACK to its FIN
	FIN_TIMEOUT_WAIT = 10 * FIN_TIMEOUT

	// How long to wait on the receiver before completing if no data received
	RECV_READ_TIMEOUT = time.Duration(5000 * time.Millisecond)
)
