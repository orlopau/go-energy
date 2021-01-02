package modbus_test

import (
	"encoding/binary"
	"github.com/orlopau/go-energy/pkg/modbus"
	"math"
	"testing"
)

func TestMockbus_ReadHoldingRegisters(t *testing.T) {
	tt := map[string]struct {
		in   map[uint16]uint32
		wErr bool
	}{
		"simple": {
			in: map[uint16]uint32{
				0: math.MaxUint32 - 200,
			},
			wErr: false,
		},
		"multiple writes": {
			in: map[uint16]uint32{
				0:  math.MaxUint32 - 200,
				20: math.MaxUint32 - 200,
			},
			wErr: false,
		},
		"overlapping registers": {
			in: map[uint16]uint32{
				20: math.MaxUint32 - 200,
				21: math.MaxUint32 - 200,
			},
			wErr: true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mockbus := modbus.NewMockbus(30)

			for addr, v := range tc.in {
				err := mockbus.AddHoldingRegisterEntry(addr, v)
				if err != nil {
					if tc.wErr {
						return
					}
					t.Fatal(err)
				}
			}

			for addr, v := range tc.in {
				regs, err := mockbus.ReadHoldingRegisters(addr, 2)
				if err != nil {
					t.Fatal(err)
				}

				u := binary.BigEndian.Uint32(regs)
				if v != u {
					t.Fatalf("expected %v, got %v", v, u)
				}
			}

			if tc.wErr {
				t.Fatal("expected error")
			}
		})
	}
}

func TestMockbus_ReadHoldingRegistersUint(t *testing.T) {
	tt := map[string]struct {
		in map[uint16]uint16
	}{
		"simple": {
			in: map[uint16]uint16{
				0: math.MaxUint16 - 200,
			},
		},
		"multiple writes": {
			in: map[uint16]uint16{
				0: math.MaxUint16 - 200,
				1: math.MaxUint16 - 200,
				2: math.MaxUint16 - 200,
				8: math.MaxUint16 - 200,
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mockbus := modbus.NewMockbus(20)

			for addr, v := range tc.in {
				err := mockbus.AddHoldingRegisterEntry(addr, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			for addr, v := range tc.in {
				uints, err := mockbus.ReadHoldingRegistersUint(addr, 1)
				if err != nil {
					t.Fatal(err)
				}
				u := uints[0]
				if v != u {
					t.Fatalf("expected %v, got %v", v, u)
				}
			}
		})
	}
}

func TestMockbus_AddHoldingRegisterEntries(t *testing.T) {
	tt := map[string]struct {
		in map[uint16]interface{}
	}{
		"simple": {
			in: map[uint16]interface{}{
				0: uint16(math.MaxUint16 - 200),
			},
		},
		"multiple writes": {
			in: map[uint16]interface{}{
				0: uint16(math.MaxUint16 - 200),
				1: uint16(math.MaxUint16 - 200),
				2: uint16(math.MaxUint16 - 200),
				8: uint16(math.MaxUint16 - 200),
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mockbus := modbus.NewMockbus(20)

			err := mockbus.AddHoldingRegisterEntries(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			for addr, v := range tc.in {
				uints, err := mockbus.ReadHoldingRegistersUint(addr, 1)
				if err != nil {
					t.Fatal(err)
				}
				u := uints[0]
				if v != u {
					t.Fatalf("expected %v, got %v", v, u)
				}
			}
		})
	}
}
