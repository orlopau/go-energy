package sunspec_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/orlopau/go-energy/pkg/sunspec"
	"math"
	"testing"
)

type dummyModelScanner struct {
	models map[uint16]uint16
	scans  uint
}

func (s *dummyModelScanner) Scan() (map[uint16]uint16, error) {
	s.scans++
	return s.models, nil
}

func TestCachedModelScanner_GetAddress(t *testing.T) {
	dummyData := map[uint16]uint16{
		100: 101,
		105: 108,
	}

	s := &dummyModelScanner{models: dummyData}
	modelScanner := sunspec.CachedModelConverter{ModelScanner: s}

	for m, ax := range dummyData {
		a, err := modelScanner.GetAddress(m)
		if err != nil {
			t.Fatal(err)
		}

		if a != ax {
			t.Fatalf("model did not match")
		}
	}

	if s.scans > 1 {
		t.Fatalf("more than 1 scan")
	}
}

func TestCachedModelScanner_HasModel(t *testing.T) {
	s := &dummyModelScanner{models: map[uint16]uint16{
		100: 101,
		105: 108,
	}}
	modelScanner := sunspec.CachedModelConverter{ModelScanner: s}

	tests := []struct {
		model  uint16
		exists bool
	}{
		{100, true},
		{105, true},
		{666, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test model %v", tt.model), func(t *testing.T) {
			has, err := modelScanner.HasModel(tt.model)
			if err != nil {
				t.Fatal(err)
			}
			if has != tt.exists {
				t.Fatalf("expected %v exists to be %v", tt.model, tt.exists)
			}
		})
	}

	if s.scans > 1 {
		t.Fatalf("more than 1 scan")
	}
}

type dummyAddressReader struct {
	uints   map[uint16]uint64
	floats  map[uint16]float64
	strings map[uint16]string
	ints    map[uint16]int64
}

// ReadInto actually only reads from the uint64 slice.
func (d *dummyAddressReader) ReadInto(address uint16, data interface{}) error {
	buffer := bytes.NewBuffer(make([]byte, 0))
	var err error
	switch data.(type) {
	case *uint64:
		err = binary.Write(buffer, binary.BigEndian, d.uints[address])
	case *uint32:
		err = binary.Write(buffer, binary.BigEndian, uint32(d.uints[address]))
	case *uint16:
		err = binary.Write(buffer, binary.BigEndian, uint16(d.uints[address]))
	case *int64:
		err = binary.Write(buffer, binary.BigEndian, d.ints[address])
	case *int32:
		err = binary.Write(buffer, binary.BigEndian, int32(d.ints[address]))
	case *int16:
		err = binary.Write(buffer, binary.BigEndian, int16(d.ints[address]))
	}

	if err != nil {
		return err
	}

	err = binary.Read(buffer, binary.BigEndian, data)

	if err != nil {
		return err
	}

	return nil
}

func (d *dummyAddressReader) ReadString(address, words uint16) (string, error) {
	v, ok := d.strings[address]
	if !ok {
		return "", fmt.Errorf("couldn'tProvider retrieve string for address %v", address)
	}

	return v, nil
}

func (d dummyAddressReader) getInt(address uint16) (int64, error) {
	v, ok := d.ints[address]
	if !ok {
		return 0, fmt.Errorf("couldn'tProvider retrieve float for address %v", address)
	}

	return v, nil
}

func (d *dummyAddressReader) ReadInt16(address uint16) (int16, error) {
	data, err := d.getInt(address)
	return int16(data), err
}

func (d *dummyAddressReader) ReadInt32(address uint16) (int32, error) {
	data, err := d.getInt(address)
	return int32(data), err
}

func (d *dummyAddressReader) ReadInt64(address uint16) (int64, error) {
	data, err := d.getInt(address)
	return data, err
}

func (d dummyAddressReader) getFloat(address uint16) (float64, error) {
	v, ok := d.floats[address]
	if !ok {
		return 0, fmt.Errorf("couldn'tProvider retrieve float for address %v", address)
	}

	return v, nil
}

func (d *dummyAddressReader) ReadFloat32(address uint16) (float32, error) {
	data, err := d.getFloat(address)
	return float32(data), err
}

func (d *dummyAddressReader) ReadFloat64(address uint16) (float64, error) {
	data, err := d.getFloat(address)
	return data, err
}

func (d dummyAddressReader) getUInt(address uint16) (uint64, error) {
	v, ok := d.uints[address]
	if !ok {
		return 0, fmt.Errorf("couldn'tProvider retrieve uint for address %v", address)
	}

	return v, nil
}

func (d *dummyAddressReader) ReadUint16(address uint16) (uint16, error) {
	data, err := d.getUInt(address)
	return uint16(data), err
}

func (d *dummyAddressReader) ReadUint32(address uint16) (uint32, error) {
	data, err := d.getUInt(address)
	return uint32(data), err
}

func (d *dummyAddressReader) ReadUint64(address uint16) (uint64, error) {
	data, err := d.getUInt(address)
	return data, err
}

func TestAddressModelScanner_Scan(t *testing.T) {
	registers := map[uint16]uint64{
		40000: 0x53756e53,
		40002: 1,
		40003: 66,
		40070: 11,
		40071: 13,
		40085: 12,
		40086: 98,
		40185: math.MaxUint16,
	}

	scanner := &sunspec.AddressModelScanner{Reader: &dummyAddressReader{uints: registers}}

	models, err := scanner.Scan()
	if err != nil {
		t.Fatal(err)
	}

	modelsx := map[uint16]uint16{
		1:  40002,
		11: 40070,
		12: 40085,
	}

	if fmt.Sprint(models) != fmt.Sprint(modelsx) {
		t.Fatalf("expected %v, got %v", modelsx, models)
	}
}

type dummyModelConverter struct {
	models map[uint16]uint16
}

func (d *dummyModelConverter) GetAddress(model uint16) (uint16, error) {
	v, ok := d.models[model]
	if !ok {
		return 0, fmt.Errorf("couldn't find address for model %v", model)
	}

	return v, nil
}

func (d *dummyModelConverter) HasModel(model uint16) (bool, error) {
	_, ok := d.models[model]
	return ok, nil
}

func TestModelReader_ReadPointUint(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			uints: map[uint16]uint64{2: 123},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{1: 1},
		},
	}

	p16, err := m.ReadPointUint16(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p16 != 123 {
		t.Fatalf("uint16 values do not match")
	}

	p32, err := m.ReadPointUint32(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p32 != 123 {
		t.Fatalf("uint32 values do not match")
	}

	p64, err := m.ReadPointUint64(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p64 != 123 {
		t.Fatalf("uint64 values do not match")
	}
}

func TestModelReader_ReadPointInt(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			ints: map[uint16]int64{2: 123},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{1: 1},
		},
	}

	p16, err := m.ReadPointInt16(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p16 != 123 {
		t.Fatalf("int16 values do not match")
	}

	p32, err := m.ReadPointInt32(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p32 != 123 {
		t.Fatalf("int32 values do not match")
	}

	p64, err := m.ReadPointInt64(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p64 != 123 {
		t.Fatalf("int64 values do not match")
	}
}

func TestModelReader_ReadPointFloat(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			floats: map[uint16]float64{2: 123.123},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{1: 1},
		},
	}

	p16, err := m.ReadPointFloat32(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p16 != 123.123 {
		t.Fatalf("float32 values do not match")
	}

	p32, err := m.ReadPointFloat64(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if p32 != 123.123 {
		t.Fatalf("float64 values do not match")
	}
}

func TestModelReader_ReadString(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			strings: map[uint16]string{2: "hallo!"},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{1: 1},
		},
	}

	s, err := m.ReadString(1, 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	if s != "hallo!" {
		t.Fatalf("float32 values do not match")
	}
}

func TestModelReader_ReadInto(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			uints: map[uint16]uint64{2: 1337},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{1: 1},
		},
	}

	var v uint64
	err := m.ReadInto(1, 1, &v)
	if err != nil {
		t.Fatal(err)
	}

	if 1337 != v {
		t.Fatalf("want %v, got %v", 1337, v)
	}
}

func TestModelReader_HasModel(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{1: 1},
		},
	}

	hasModel, err := m.HasModel(1)
	if err != nil {
		t.Fatal(err)
	}

	if !hasModel {
		t.Fatalf("want true, got false")
	}

	hasModel, err = m.HasModel(111)
	if err != nil {
		t.Fatal(err)
	}

	if hasModel {
		t.Fatalf("want false, got true")
	}
}
