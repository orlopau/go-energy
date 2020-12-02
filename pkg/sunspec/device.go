package sunspec

import (
	"github.com/orlopau/go-energy/pkg/modbus"
	"io"
)

type addressReaderCloser interface {
	AddressReader
	io.Closer
}

// Device represents a SunSpec device connected via modbus.
type Device struct {
	io.Closer
	ModelReader
}

// Connect connects to a SunSpec modbus TCP device.
func Connect(slaveId byte, addr string) (Device, error) {
	client, err := modbus.Connect(addr, slaveId)
	if err != nil {
		return Device{}, err
	}

	device := newDevice(client)
	return device, nil
}

// newDevice creates a new device using an addressReaderCloser.
func newDevice(arc addressReaderCloser) Device {
	scanner := &AddressModelScanner{UIntReader: arc}
	converter := &CachedModelConverter{
		ModelScanner: scanner,
	}

	m := ModelReader{
		AddressReader:  arc,
		ModelConverter: converter,
	}

	d := Device{
		Closer:      arc,
		ModelReader: m,
	}

	return d
}
