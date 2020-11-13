package discovery

import (
	"encoding/hex"
	"errors"
	"net"
	"os"
	"time"
)

const (
	multicastAddress = "239.12.255.254:9522"
	discoveryPayload = "534d4100000402a0ffffffff0000002000000000"
)

// send sends the discovery payload message to the specified udp connection.
func send(conn net.PacketConn, addr net.Addr) error {
	// send multicast for discovery
	bytes, err := hex.DecodeString(discoveryPayload)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(bytes, addr)
	return err
}

// listen listens for the returned discovery messages and returns the addresses of the senders.
func listen(conn net.PacketConn, timeout time.Duration) ([]net.Addr, error) {
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	results := make([]net.Addr, 0)
	for {
		buf := make([]byte, 2500)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				break
			} else {
				return nil, err
			}
		}
		results = append(results, addr)
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
	defer func() {
		e := conn.Close()
		if e != nil {
			panic(e)
		}
	}()

	if err := send(conn, addr); err != nil {
		return nil, err
	}

	addrs, err := listen(conn, timeout)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}