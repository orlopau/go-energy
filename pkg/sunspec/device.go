package sunspec

import (
	"io"
)

type AddressReaderCloser interface {
	AddressReader
	io.Closer
}

type Device struct {
	io.Closer
	ModelReader
}

func NewDevice(arc AddressReaderCloser) Device {
	scanner := &AddressModelScanner{UIntReader: arc}
	converter := &CachedModelConverter{
		ModelScanner: scanner,
	}

	m := ModelReader{
		AddressReader:  arc,
		ModelConverter: converter,
	}

	d := Device{
		Closer:      arc,
		ModelReader: m,
	}

	return d
}
