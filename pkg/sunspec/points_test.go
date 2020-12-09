package sunspec_test

import (
	"github.com/orlopau/go-energy/pkg/sunspec"
	"testing"
)

var testPoint = sunspec.Point{Point: 11, Model: 802, T: uint16(0)}

func TestModelReader_HasPoint(t *testing.T) {
	m := &sunspec.ModelReader{
		Reader: &dummyAddressReader{
			uints: map[uint16]uint64{12: 55},
		},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{802: 1},
		},
	}

	hasPoint, _, err := m.HasAnyPoint(testPoint)
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
		Reader: &dummyAddressReader{},
		Converter: &dummyModelConverter{
			models: map[uint16]uint16{802: 1},
		},
	}

	hasPoint, _, err := m.HasAnyPoint(sunspec.PointSoc)
	if err != nil {
		t.Fatal(err)
	}

	if hasPoint {
		t.Fatalf("want false got true")
	}
}

var testPoint2 = sunspec.Point{Model: 101, Point: 14, T: int16(0), Scaled: true, Unit: sunspec.UnitWatts}

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

	// test unscaled
	s, err := m.GetAnyPoint(testPoint)
	if err != nil {
		t.Fatal(err)
	}

	if s != soc {
		t.Fatalf("want %v, got %v", soc, s)
	}

	// test scaled
	p, err := m.GetAnyPoint(testPoint2)
	if err != nil {
		t.Fatal(err)
	}

	if p != scaledPow {
		t.Fatalf("want %v, got %v", scaledPow, p)
	}
}
