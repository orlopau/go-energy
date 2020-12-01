package modbus

import (
	"bytes"
	"encoding/binary"
	"github.com/goburrow/modbus"
	"time"
)

type Client struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
}

func Connect(addr string, slaveId byte) (*Client, error) {
	handler := modbus.NewTCPClientHandler(addr)
	handler.Timeout = 20 * time.Second
	handler.SlaveId = slaveId
	err := handler.Connect()
	if err != nil {
		return nil, err
	}

	return &Client{handler: handler, client: modbus.NewClient(handler)}, nil
}

func (c *Client) Close() error {
	return c.handler.Close()
}

func (c *Client) ReadRegisterInto(address, quantity uint16, data interface{}) error {
	registers, err := c.client.ReadHoldingRegisters(address, quantity)
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

func (c *Client) ReadUint16(address uint16) (uint16, error) {
	var val uint16
	err := c.ReadRegisterInto(address, 1, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadUint32(address uint16) (uint32, error) {
	var val uint32
	err := c.ReadRegisterInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadUint64(address uint16) (uint64, error) {
	var val uint64
	err := c.ReadRegisterInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadInt16(address uint16) (int16, error) {
	var val int16
	err := c.ReadRegisterInto(address, 1, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadInt32(address uint16) (int32, error) {
	var val int32
	err := c.ReadRegisterInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadInt64(address uint16) (int64, error) {
	var val int64
	err := c.ReadRegisterInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadFloat32(address uint16) (float32, error) {
	var val float32
	err := c.ReadRegisterInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadFloat64(address uint16) (float64, error) {
	var val float64
	err := c.ReadRegisterInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadString(address, words uint16) (string, error) {
	registers, err := c.client.ReadHoldingRegisters(address, words)
	if err != nil {
		return "", err
	}
	return string(registers), nil
}
