# Fast And Accurate Reliable Transfer

## Building

`make` within this directory, which outputs `3700recv` and `3700send`

## Running

Run `3700recv` to set up the receiving server and receive the bound port.

Run `[data] | 3700send [hostname:port]` to connect to the receiving server on the given port above.
Data will be read from 3700send on STDIN and sent to 3700recv until completed.

## Approach

I use selective acks (SACK) on the receiver's side, so that there is less of a packet
delay in the event of a dropped / timed out packet.

Timeouts are detected with a fixed-length timeout period (`shared/timeouts.go`).
On the sender's side this is `1000ms`. On the receiver's side, if no packets
are received within 10 sender timeouts, it's assumed that the connection can be
closed.

Once all data has been transmitted successfully (detected by the DONE flag on a packet),
the receiver blasts out 7 ACKS and quits, printing the data received to STDOUT.

The sender attempts to send data until it receives an ACK to its DONE, however if it
does not receive an ACK to 5 consecutive DONE sends, it will terminate as well.

## Problems

Because I'm using a fixed window size of len(data), this does not perform well
under a poor latency test (the `0.1 mb/s` bandwidth, `500 ms` latency test
can take up to 2 minutes to run in the worst case).

This is also true for the packet loss cases -- thanks to the second-long timeouts,
often the program will take one timeout cycle to complete successfully.

Both of these could be solved with RTT tracking and a dynamic timeout, however
experimenting with this functionality gave me worse performance in either poor
bandwidth or in packet loss, so I took average-to-bad performance in both rather
than sacrificing one for the other.

## Testing

Using https://github.com/tylertreat/comcast and a text file of Moby Dick I was
able to test my transfer system's success/failure/speed under varying network
speeds.

I have unit tests for some of the more technical behavior (notably checksums
and compression) but otherwise relied on end-to-end testing.
