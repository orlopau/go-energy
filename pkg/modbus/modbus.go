package modbus

import (
	"bytes"
	"encoding/binary"
	"github.com/goburrow/modbus"
)

type Reader struct {
	Client modbus.Client
}

func (t *Reader) ReadRegisterInto(address, quantity uint16, data interface{}) error {
	registers, err := t.Client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(registers)
	err = binary.Read(buf, binary.BigEndian, data)
	if err != nil {
		return err
	}

	return nil
}

func (t *Reader) ReadUint16(address uint16) (uint16, error) {
	var val uint16
	err := t.ReadRegisterInto(address, 1, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadUint32(address uint16) (uint32, error) {
	var val uint32
	err := t.ReadRegisterInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadUint64(address uint16) (uint64, error) {
	var val uint64
	err := t.ReadRegisterInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadInt16(address uint16) (int16, error) {
	var val int16
	err := t.ReadRegisterInto(address, 1, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadInt32(address uint16) (int32, error) {
	var val int32
	err := t.ReadRegisterInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadInt64(address uint16) (int64, error) {
	var val int64
	err := t.ReadRegisterInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadFloat32(address uint16) (float32, error) {
	var val float32
	err := t.ReadRegisterInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadFloat64(address uint16) (float64, error) {
	var val float64
	err := t.ReadRegisterInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (t *Reader) ReadString(address, words uint16) (string, error) {
	registers, err := t.Client.ReadHoldingRegisters(address, words)
	if err != nil {
		return "", err
	}
	return string(registers), nil
}
