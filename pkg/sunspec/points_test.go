package sunspec_test

import (
	"github.com/orlopau/go-energy/pkg/sunspec"
	"testing"
)

func TestModelReader_HasPoint(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			uints: map[uint16]uint64{1: 55},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{802: 1},
		},
	}

	hasPoint, err := m.HasPoint(sunspec.PointSoc)
	if err != nil {
		t.Fatal(err)
	}

	if !hasPoint {
		t.Fatalf("want true got false")
	}
}

func TestModelReader_HasPoint_NotImplemented(t *testing.T) {
	t.Skipf("Skipped because not implemented detections is not yet implemented!")
	// TODO implement "not implemented" detection

	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{802: 1},
		},
	}

	hasPoint, err := m.HasPoint(sunspec.PointSoc)
	if err != nil {
		t.Fatal(err)
	}

	if hasPoint {
		t.Fatalf("want false got true")
	}
}

func TestModelReader_GetPoint(t *testing.T) {
	const soc = 55
	const pow = 20
	const scale = 2
	const scaledPow = 2000

	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			uints: map[uint16]uint64{12: soc, 115: scale},
			ints:  map[uint16]int64{114: pow},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{802: 1, 101: 100},
		},
	}

	s, err := m.GetPoint(sunspec.PointSoc)
	if err != nil {
		t.Fatal(err)
	}

	if s != soc {
		t.Fatalf("want %v, got %v", soc, s)
	}

	p, err := m.GetPoint(sunspec.PointPower1Phase)
	if err != nil {
		t.Fatal(err)
	}

	if p != scaledPow {
		t.Fatalf("want %v, got %v", scaledPow, p)
	}
}
