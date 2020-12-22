// Provides functionality for reading and decoding SMA energymeter data packets.
package meter

import (
	"fmt"
	"io"
	"net"
)

const (
	multicastIP = "239.12.255.254:9522"
)

type energyMeterConnection interface {
	ReadFromUDP(b []byte) (int, *net.UDPAddr, error)
	io.Closer
}

type EnergyMeter struct {
	Conn energyMeterConnection
}

// Listen opens a multicast socket to listen for energymeter messages.
//
// Returns an EnergyMeter representing the opened connection.
func Listen() (*EnergyMeter, error) {
	addr, err := net.ResolveUDPAddr("udp", multicastIP)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &EnergyMeter{Conn: l}, nil
}

// Close closes the opened connection.
func (t *EnergyMeter) Close() error {
	if t.Conn == nil {
		return fmt.Errorf("no connection to close")
	}
	return t.Conn.Close()
}

// ReadTelegram reads and decodes an energymeter telegram from the connection.
//
// The method blocks until a telegram is received.
func (t *EnergyMeter) ReadTelegram() (*EnergyMeterTelegram, error) {
	if t.Conn == nil {
		return nil, fmt.Errorf("connection not opened")
	}

	b := make([]byte, 8192)
	_, _, err := t.Conn.ReadFromUDP(b)
	if err != nil {
		return nil, err
	}

	telegram, err := DecodeTelegram(b)
	if err != nil {
		return nil, err
	}

	return telegram, nil
}
