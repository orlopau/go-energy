package sunspec

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"math"
)

const (
	UnitWatts      = "W"
	UnitPercentage = "%"
)

var (
	ErrPointNotImplemented = errors.New("point is not implemented")
)

const (
	// float values are not implemented if they are NaN
	notImplInt16  = math.MinInt16
	notImplUint16 = math.MaxUint16
	notImplInt32  = math.MinInt32
	notImplUint32 = math.MaxUint32
)

var (
	PointSoc           = Point{Model: 124, Point: 8, T: uint16(0), Unit: UnitPercentage}
	PointDeviceAddress = Point{Model: 1, Point: 66, T: uint16(0)}
	PointPower1Phase   = Point{Model: 101, Point: 14, T: int16(0), Scaled: true, Unit: UnitWatts}
	PointPower2Phase   = Point{Model: 102, Point: 14, T: int16(0), Scaled: true, Unit: UnitWatts}
	PointPower3Phase   = Point{Model: 103, Point: 14, T: int16(0), Scaled: true, Unit: UnitWatts}
)

// Point represents a SunSpec point.
type Point struct {
	Point, Model uint16
	// T is the type of the point.
	T            interface{}
	// Scaled must be set to true if the value is a scaled value with an additional register for scaling.
	Scaled       bool
	Unit         string
}

func (p Point) String() string {
	return fmt.Sprintf("Point{point:%v,model:%v,scaled:%v,unit:%v}", p.Point, p.Model, p.Scaled, p.Unit)
}

// hasPoint returns true if the reader has the model of the specified Point.
func (r *ModelReader) hasPoint(p Point) (bool, error) {
	hasModel, err := r.Converter.HasModel(p.Model)
	if err != nil {
		return false, err
	}

	if !hasModel {
		return false, nil
	}

	return true, nil
}

// HasAnyPoint checks if any of the specified Points is present on the device.
//
// If a point is present, it returns that point.
//
// If no point is present, it returns false.
func (r *ModelReader) HasAnyPoint(ps ...Point) (bool, Point, error) {
	for _, v := range ps {
		has, err := r.hasPoint(v)
		if err != nil {
			return false, Point{}, err
		}
		if has {
			return true, v, nil
		}
	}

	return false, Point{}, nil
}

// getPoint reads a Point from a SunSpec reader.
func (r *ModelReader) getPoint(p Point) (float64, error) {
	tmpVal := p.T

	var val float64
	var err error
	switch raw := tmpVal.(type) {
	case float64:
		err = r.ReadInto(p.Model, p.Point, &raw)
		if math.IsNaN(raw) {
			return 0, ErrPointNotImplemented
		}
		val = raw
	case float32:
		err = r.ReadInto(p.Model, p.Point, &raw)
		if math.IsNaN(float64(raw)) {
			return 0, ErrPointNotImplemented
		}
		val = float64(raw)
	case uint16:
		err = r.ReadInto(p.Model, p.Point, &raw)
		if raw == notImplUint16 {
			return 0, ErrPointNotImplemented
		}
		val = float64(raw)
	case uint32:
		err = r.ReadInto(p.Model, p.Point, &raw)
		if raw == notImplUint32 {
			return 0, ErrPointNotImplemented
		}
		val = float64(raw)
	case int16:
		err = r.ReadInto(p.Model, p.Point, &raw)
		if raw == notImplInt16 {
			return 0, ErrPointNotImplemented
		}
		val = float64(raw)
	case int32:
		err = r.ReadInto(p.Model, p.Point, &raw)
		if val == notImplInt32 {
			return 0, ErrPointNotImplemented
		}
		val = float64(raw)
	}

	if err != nil {
		return 0, err
	}

	if !p.Scaled {
		return val, nil
	}

	size := uint16(math.Floor(float64(binary.Size(tmpVal)) / 2.0))

	var factor uint16
	err = r.ReadInto(p.Model, p.Point+size, &factor)
	if err != nil {
		return 0, err
	}

	return val * math.Pow10(int(factor)), nil
}

// GetAnyPoint fetches the first available point and returns its value.
//
// Returns an error if none point of the given points is found.
func (r *ModelReader) GetAnyPoint(ps ...Point) (float64, error) {
	for _, v := range ps {
		p, err := r.getPoint(v)
		if err == nil {
			return p, nil
		}
		if !errors.Is(err, ErrPointNotImplemented) {
			return 0, err
		}
	}

	return 0, errors.Wrap(ErrPointNotImplemented, fmt.Sprintf("did not find any of these points %v", ps))
}
