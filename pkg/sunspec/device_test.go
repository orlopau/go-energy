package sunspec_test

import (
	"github.com/orlopau/go-energy/pkg/sunspec"
	"testing"
)

type dummyARC struct {
	*dummyAddressReader
	closed bool
}

func (d *dummyARC) Close() error {
	d.closed = true
	return nil
}

func TestNewDevice(t *testing.T) {
	dar := &dummyAddressReader{}

	dARC := &dummyARC{
		dummyAddressReader: dar,
	}

	_ = sunspec.NewDevice(dARC)
}
