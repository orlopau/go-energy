package sunspec

import (
	"github.com/orlopau/go-energy/pkg/modbus"
)

type addressReaderCloser interface {
	addressReader
}

// Connect connects to a SunSpec modbus TCP device.
func Connect(slaveId byte, addr string) (*ModelReader, error) {
	client, err := modbus.Connect(addr, slaveId)
	if err != nil {
		return nil, err
	}

	device := newDevice(client)
	return device, nil
}

// newDevice creates a new device using an addressReaderCloser.
func newDevice(arc addressReaderCloser) *ModelReader {
	scanner := &AddressModelScanner{Reader: arc}
	converter := &CachedModelConverter{
		ModelScanner: scanner,
	}

	m := ModelReader{
		Reader:    arc,
		Converter: converter,
	}

	return &m
}
