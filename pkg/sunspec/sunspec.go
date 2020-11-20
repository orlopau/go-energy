package sunspec

import (
	"fmt"
	"math"
)

var (
	// sunsBaseAddresses contain the base modbus registers where the sunsIdentifier must be present
	// if the device is SunSpec compatible.
	sunsBaseAddresses = [...]uint16{
		40000,
		50000,
		00000,
	}
)

const (
	// sunsIdentifier is the well-known identifier present at one of the sunsBaseAddresses
	// if the device is SunSpec compatible.
	sunsIdentifier uint32 = 0x53756e53
)

type SunSpecReader struct {
	ModbusReader ModbusReader
	Models       map[uint16]uint16
}

// Scan scans the devices SunSpec Models using the sunsBaseAddresses and stores their addresses.
//
// This function must be executed prior to reading SunSpec module data.
//
// The register specified by the offset must point to the SunSpec Common Model ID.
// For further information, consult the documentation provided by https://sunspec.org/
func (t *SunSpecReader) Scan() error {
	// scan all base addresses for sunspec identifier
	var offset uint16
	for _, address := range sunsBaseAddresses {
		val, err := t.ModbusReader.ReadUint32(address)
		if err != nil {
			return err
		}
		if val == sunsIdentifier {
			offset = address + 2
			break
		}
	}

	models := make(map[uint16]uint16)
	for {
		modelID, err := t.ModbusReader.ReadUint16(offset)
		if err != nil {
			return err
		}

		if modelID == ^uint16(0) {
			break
		}

		models[modelID] = offset

		l, err := t.ModbusReader.ReadUint16(offset + 1)
		if err != nil {
			return err
		}

		offset += l + 2
	}

	t.Models = models
	return nil
}

// getModelAddress returns the modbus address for the given SunSpec model id
func (t *SunSpecReader) getModelAddress(model uint16) (uint16, error) {
	address, ok := t.Models[model]
	if !ok {
		return 0, fmt.Errorf("the scanned modules do not contain this model with id %d", model)
	}
	return address, nil
}

func (t *SunSpecReader) ReadPointUint16(model, point uint16) (uint16, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadUint16(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MaxUint16 {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (t *SunSpecReader) ReadPointUint32(model, point uint16) (uint32, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadUint32(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MaxUint32 {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (t *SunSpecReader) ReadPointUint64(model, point uint16) (uint64, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadUint64(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MaxUint64 {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (t *SunSpecReader) ReadPointInt16(model, point uint16) (int16, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadInt16(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MinInt16 {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (t *SunSpecReader) ReadPointInt32(model, point uint16) (int32, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadInt32(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MinInt32 {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (t *SunSpecReader) ReadPointInt64(model, point uint16) (int64, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadInt64(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MinInt64 {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (t *SunSpecReader) ReadPointFloat32(model, point uint16) (float32, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadFloat32(address + point)
	if err != nil {
		return 0, err
	}

	if math.IsNaN(float64(val)) {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (t *SunSpecReader) ReadPointFloat64(model, point uint16) (float64, error) {
	address, err := t.getModelAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := t.ModbusReader.ReadFloat64(address + point)
	if err != nil {
		return 0, err
	}

	if math.IsNaN(val) {
		return 0, fmt.Errorf("data point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}
