package sunspec

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	UnitWatts      = "W"
	UnitPercentage = "%"
)

var (
	PointSoc         = Point{Model: 124, Point: 8, T: uint16(0), Unit: UnitPercentage}
	PointPower1Phase = Point{Model: 101, Point: 14, T: int16(0), Scaled: true, Unit: UnitWatts}
	PointPower2Phase = Point{Model: 102, Point: 14, T: int16(0), Scaled: true, Unit: UnitWatts}
	PointPower3Phase = Point{Model: 103, Point: 14, T: int16(0), Scaled: true, Unit: UnitWatts}
)

type Point struct {
	Point, Model uint16
	T            interface{}
	Scaled       bool
	Unit         string
}

func (p Point) String() string {
	return fmt.Sprintf("Point{point:%v,model:%v,scaled:%v,unit:%v}", p.Point, p.Model, p.Scaled, p.Unit)
}

// HasPoint returns true if the reader has the model of the specified Point.
func (r *ModelReader) HasPoint(p Point) (bool, error) {
	hasModel, err := r.Converter.HasModel(p.Model)
	if err != nil {
		return false, err
	}

	if !hasModel {
		return false, nil
	}

	return true, nil
}

// GetPoint reads a Point from a SunSpec reader.
func (r *ModelReader) GetPoint(p Point) (float64, error) {
	tmpVal := p.T

	var val float64
	var err error
	switch i := tmpVal.(type) {
	case float64:
		err = r.ReadInto(p.Model, p.Point, &i)
		val = i
	case float32:
		err = r.ReadInto(p.Model, p.Point, &i)
		val = float64(i)
	case uint16:
		err = r.ReadInto(p.Model, p.Point, &i)
		val = float64(i)
	case uint32:
		err = r.ReadInto(p.Model, p.Point, &i)
		val = float64(i)
	case int16:
		err = r.ReadInto(p.Model, p.Point, &i)
		val = float64(i)
	case int32:
		err = r.ReadInto(p.Model, p.Point, &i)
		val = float64(i)
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

func (r *ModelReader) GetAnyPoint(ps ...Point) (float64, error) {
	for _, v := range ps {
		p, err := r.GetPoint(v)
		if err == nil {
			return p, nil
		}
	}

	return 0, fmt.Errorf("did not find any of these points %v", ps)
}
