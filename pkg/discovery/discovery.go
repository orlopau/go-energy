// Provides SMA device discovery functionality.
package discovery

import (
	"bytes"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"time"
)

const (
	multicastAddress       = "239.12.255.254:9522"
	discoveryPayload       = "534d4100000402a0ffffffff0000002000000000"
	discoveryVerifyPayload = "534d4100000402A000000001000200000001"
)

// send sends the discovery payload message to the specified udp connection.
func send(conn net.PacketConn, addr net.Addr) error {
	// send multicast for discovery
	p, err := hex.DecodeString(discoveryPayload)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(p, addr)
	return err
}

// listen listens for the returned discovery messages and returns the addresses of the senders.
func listen(conn net.PacketConn, timeout time.Duration) ([]net.Addr, error) {
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	verifyPayload, err := hex.DecodeString(discoveryVerifyPayload)
	if err != nil {
		return nil, err
	}

	results := make([]net.Addr, 0)
	for {
		buf := make([]byte, 2500)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				break
			} else {
				return nil, err
			}
		}

		var contains bool
		for _, v := range results {
			if addr == v {
				contains = true
				break
			}
		}

		if bytes.Compare(buf[0:n], verifyPayload) == 0 && !contains {
			results = append(results, addr)
		}
	}

	return results, nil
}

// DiscoverInverters discovers inverters connected to the network at the specified interface.
//
// The function sends a multicast discover request and waits the specified duration for responses.
func DiscoverInverters(ifi *net.Interface, timeout time.Duration) ([]net.Addr, error) {
	addr, err := net.ResolveUDPAddr("udp", multicastAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return nil, err
	}

	if err := send(conn, addr); err != nil {
		return nil, err
	}

	addrs, err := listen(conn, timeout)
	if err != nil {
		return nil, err
	}

	e := conn.Close()
	if e != nil {
		panic(e)
	}

	return addrs, nil
}
