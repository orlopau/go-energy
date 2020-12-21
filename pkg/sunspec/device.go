package sunspec

import (
	"github.com/orlopau/go-energy/pkg/modbus"
	"github.com/pkg/errors"
)

type ModbusDevice struct {
	*ModelReader
	client *modbus.Client
}

// Connect connects to a SunSpec modbus TCP device.
func Connect(addr string) (*ModbusDevice, error) {
	client, err := modbus.Connect(addr)
	if err != nil {
		return nil, err
	}

	device := newDevice(client)
	return device, nil
}

// newDevice creates a new device using an addressReaderCloser.
func newDevice(client *modbus.Client) *ModbusDevice {
	scanner := &AddressModelScanner{Reader: client}
	converter := &CachedModelConverter{
		ModelScanner: scanner,
	}

	m := &ModelReader{
		Reader:    client,
		Converter: converter,
	}

	return &ModbusDevice{ModelReader: m, client: client}
}

func (d *ModbusDevice) AutoSetDeviceAddress() error {
	d.SetDeviceAddress(126)

	addr, err := d.GetAnyPoint(PointDeviceAddress)
	if err != nil {
		return errors.Wrap(err, "auto setup of device address")
	}

	d.SetDeviceAddress(byte(addr))
	return nil
}

func (d *ModbusDevice) SetDeviceAddress(deviceAddr byte) {
	d.client.SetSlaveID(deviceAddr)
}
