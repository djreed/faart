package shared

import "time"

var (
	// How long to wait between packet sends
	SEND_PACKET_TIMEOUT = time.Duration(0)

	// How long to wait before re-queueing packets
	QUEUE_DATA_TIMEOUT = time.Duration(1500 * time.Millisecond)

	// How long to wait before re-queueing the FIN
	QUEUE_FIN_TIMEOUT = time.Duration(100 * time.Millisecond)

	// How long to wait before completing if sender has not received
	// an ACK to its FIN
	SEND_FIN_TIMEOUT = 10 * QUEUE_FIN_TIMEOUT

	// How long to wait on the receiver before completing if no data received
	RECV_READ_TIMEOUT = time.Duration(6000 * time.Millisecond)
)
