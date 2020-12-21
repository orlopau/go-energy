package modbus

import (
	"bytes"
	"encoding/binary"
	"github.com/goburrow/modbus"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

const backoffDuration = 10 * time.Second

type registerReader interface {
	ReadHoldingRegisters(address uint16, quantity uint16) (results []byte, err error)
}

type Client struct {
	handler *modbus.TCPClientHandler
	client  registerReader
}

func Connect(addr string) (*Client, error) {
	handler := modbus.NewTCPClientHandler(addr)
	handler.Timeout = 20 * time.Second
	handler.IdleTimeout = 24 * time.Hour
	err := reconnect(handler)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to modbus")
	}

	return &Client{handler: handler, client: modbus.NewClient(handler)}, nil
}

func (c *Client) Close() error {
	return c.handler.Close()
}

func (c *Client) SetSlaveID(id byte) {
	c.handler.SlaveId = id
}

func (c *Client) ReadInto(address uint16, v interface{}) error {
	b := binary.Size(v)

	return c.readBytesInto(address, uint16(b), v)
}

func reconnect(handler *modbus.TCPClientHandler) error {
	err := handler.Close()
	if err != nil {
		return err
	}

	for {
		log.Printf("connecting to %v...", handler.Address)
		err := handler.Connect()
		if err == nil {
			log.Printf("connected to %v", handler.Address)
			return nil
		}

		log.Println(errors.Wrap(err, "couldn't connect to device").Error())
		<-time.After(backoffDuration)
	}
}

func (c *Client) readBytesInto(address, quantity uint16, data interface{}) error {
	for {
		registers, err := c.client.ReadHoldingRegisters(address, quantity)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			err := reconnect(c.handler)
			if err != nil {
				return err
			}
			continue
		}
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
}

func (c *Client) ReadUint16(address uint16) (uint16, error) {
	var val uint16
	err := c.readBytesInto(address, 1, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadUint32(address uint16) (uint32, error) {
	var val uint32
	err := c.readBytesInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadUint64(address uint16) (uint64, error) {
	var val uint64
	err := c.readBytesInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadInt16(address uint16) (int16, error) {
	var val int16
	err := c.readBytesInto(address, 1, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadInt32(address uint16) (int32, error) {
	var val int32
	err := c.readBytesInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadInt64(address uint16) (int64, error) {
	var val int64
	err := c.readBytesInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadFloat32(address uint16) (float32, error) {
	var val float32
	err := c.readBytesInto(address, 2, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadFloat64(address uint16) (float64, error) {
	var val float64
	err := c.readBytesInto(address, 4, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Client) ReadString(address, words uint16) (string, error) {
	var val string
	err := c.readBytesInto(address, words, val)
	if err != nil {
		return "", err
	}
	return val, nil
}
