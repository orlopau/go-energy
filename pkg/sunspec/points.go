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
	PointSoc         = Point{model: 124, point: 8, t: uint16(0), unit: UnitPercentage}
	PointPower1Phase = Point{model: 101, point: 14, t: int16(0), scaled: true, unit: UnitWatts}
	PointPower2Phase = Point{model: 102, point: 14, t: int16(0), scaled: true, unit: UnitWatts}
	PointPower3Phase = Point{model: 103, point: 14, t: int16(0), scaled: true, unit: UnitWatts}
)

type Point struct {
	point, model uint16
	t            interface{}
	scaled       bool
	unit         string
}

func (p Point) String() string {
	return fmt.Sprintf("Point{point:%v,model:%v,scaled:%v,unit:%v}", p.point, p.model, p.scaled, p.unit)
}

// HasPoint returns true if the reader has the model of the specified Point.
func (r *ModelReader) HasPoint(p Point) (bool, error) {
	hasModel, err := r.Converter.HasModel(p.model)
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
	tmpVal := p.t

	var val float64
	var err error
	switch i := tmpVal.(type) {
	case float64:
		err = r.ReadInto(p.model, p.point, &i)
		val = i
	case float32:
		err = r.ReadInto(p.model, p.point, &i)
		val = float64(i)
	case uint16:
		err = r.ReadInto(p.model, p.point, &i)
		val = float64(i)
	case uint32:
		err = r.ReadInto(p.model, p.point, &i)
		val = float64(i)
	case int16:
		err = r.ReadInto(p.model, p.point, &i)
		val = float64(i)
	case int32:
		err = r.ReadInto(p.model, p.point, &i)
		val = float64(i)
	}

	if err != nil {
		return 0, err
	}

	if !p.scaled {
		return val, nil
	}

	size := uint16(math.Floor(float64(binary.Size(tmpVal)) / 2.0))

	var factor uint16
	err = r.ReadInto(p.model, p.point+size, &factor)
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
