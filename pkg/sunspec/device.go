package sunspec

import (
	"github.com/orlopau/go-energy/pkg/modbus"
)

type addressReaderCloser interface {
	addressReader
}

// Device represents a SunSpec device connected via modbus.
type Device struct {
	ModelReader
}

// Connect connects to a SunSpec modbus TCP device.
func Connect(slaveId byte, addr string) (*Device, error) {
	client, err := modbus.Connect(addr, slaveId)
	if err != nil {
		return nil, err
	}

	device := newDevice(client)
	return device, nil
}

// newDevice creates a new device using an addressReaderCloser.
func newDevice(arc addressReaderCloser) *Device {
	scanner := &AddressModelScanner{Reader: arc}
	converter := &CachedModelConverter{
		ModelScanner: scanner,
	}

	m := ModelReader{
		Reader:    arc,
		Converter: converter,
	}

	d := &Device{
		ModelReader: m,
	}

	return d
}
