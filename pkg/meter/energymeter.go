package meter

import (
	"fmt"
	"io"
	"net"
)

const (
	multicastIP = "239.12.255.254:9522"
)

type EnergyMeterConnection interface {
	ReadFromUDP(b []byte) (int, *net.UDPAddr, error)
	io.Closer
}

type EnergyMeterReader interface {
	ReadTelegram() (*EnergyMeterTelegram, error)
}

type energyMeter struct {
	Conn EnergyMeterConnection
}

// Listen opens a multicast socket to listen for energymeter messages.
//
// Messages can be read using ReadTelegram.
func Listen() (*energyMeter, error) {
	addr, err := net.ResolveUDPAddr("udp", multicastIP)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &energyMeter{Conn: l}, nil
}

func (t *energyMeter) Close() error {
	if t.Conn == nil {
		return fmt.Errorf("no connection to close")
	}
	return t.Conn.Close()
}

// ReadTelegram reads and decodes an energymeter telegram from the connection.
//
// The method blocks until a telegram is received.
func (t *energyMeter) ReadTelegram() (*EnergyMeterTelegram, error) {
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
