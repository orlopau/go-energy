package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var byteOrder = binary.BigEndian

type Mockbus struct {
	holdingRegisters []byte
}

// NewMockbus creates a new mockbus instance.
//
// The parameter specifies the number of registers that can be used starting from 0.
func NewMockbus(i int) *Mockbus {
	return &Mockbus{holdingRegisters: make([]byte, i*2)}
}

func (m *Mockbus) AddHoldingRegisterEntry(addr uint16, data interface{}) error {
	var buf bytes.Buffer
	if err := binary.Write(&buf, byteOrder, data); err != nil {
		return err
	}

	bs := buf.Bytes()
	if len(bs)%2 != 0 {
		return fmt.Errorf("invalid data length, bytes must be multiple of two")
	}

	start := int(addr) * 2
	end := start + len(bs)

	for k, v := range m.holdingRegisters[start:end] {
		if v != 0 {
			return fmt.Errorf("adding this entry would override data at byte %v", k)
		}
	}

	copy(m.holdingRegisters[start:end], bs)

	return nil
}

func (m *Mockbus) AddHoldingRegisterEntries(entries map[uint16]interface{}) error {
	for address, v := range entries {
		err := m.AddHoldingRegisterEntry(address, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Mockbus) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	start := address * 2
	end := start + quantity*2

	if len(m.holdingRegisters) < int(end) {
		return nil, fmt.Errorf("register does not exist")
	}

	return m.holdingRegisters[start:end], nil
}

func (m *Mockbus) ReadHoldingRegistersUint(address, quantity uint16) ([]uint16, error) {
	regs, err := m.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return nil, err
	}

	uints := make([]uint16, len(regs)/2)
	for i := range uints {
		start := i * 2
		end := start + 2
		uints[i] = byteOrder.Uint16(regs[start:end])
	}

	return uints, nil
}
