package sunspec

import (
	"fmt"
	"math"
)

const (
	// sunsIdentifier is the well-known identifier present at one of the sunsBaseAddresses
	// if the device is SunSpec compatible.
	sunsIdentifier uint32 = 0x53756e53
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

// AddressReader is a interface wrapping methods for reading types from addresses.
type AddressReader interface {
	ReadUint16(address uint16) (uint16, error)
	ReadUint32(address uint16) (uint32, error)
	ReadUint64(address uint16) (uint64, error)
	ReadInt16(address uint16) (int16, error)
	ReadInt32(address uint16) (int32, error)
	ReadInt64(address uint16) (int64, error)
	ReadFloat32(address uint16) (float32, error)
	ReadFloat64(address uint16) (float64, error)
	ReadString(address, words uint16) (string, error)
}

type ModelConverter interface {
	GetAddress(model uint16) (uint16, error)
	HasModel(model uint16) (bool, error)
}

type ModelScanner interface {
	Scan() (map[uint16]uint16, error)
}

// ModelReader provides functionality reading SunSpec models and points.
type ModelReader struct {
	AddressReader  AddressReader
	ModelConverter ModelConverter
}

// CachedModelConverter implements ModelConverter by lazily scanning the SunSpec device and caching.
//
// The models are cached until Scan is executed again.
type CachedModelConverter struct {
	ModelScanner ModelScanner
	models       map[uint16]uint16
}

// AddressModelScanner implements ModelScanner scanning the device using the SunSpec specification.
type AddressModelScanner struct {
	UIntReader interface {
		ReadUint16(address uint16) (uint16, error)
		ReadUint32(address uint16) (uint32, error)
	}
}

// Scan scans the devices SunSpec models using the sunsBaseAddresses and stores them in the cache.
//
// The register specified by the offset must point to the SunSpec Common Model ID.
// For further information, consult the documentation provided by https://sunspec.org/
func (s *AddressModelScanner) Scan() (map[uint16]uint16, error) {
	// scan all base addresses for SunSpec identifier
	var offset uint16
	for _, address := range sunsBaseAddresses {
		val, err := s.UIntReader.ReadUint32(address)
		if err != nil {
			return nil, fmt.Errorf("couldn't read base address: %w", err)
		}

		if val == sunsIdentifier {
			offset = address + 2
			break
		}
	}

	models := make(map[uint16]uint16)

	for {
		modelID, err := s.UIntReader.ReadUint16(offset)
		if err != nil {
			return nil, err
		}

		if modelID == ^uint16(0) {
			break
		}

		models[modelID] = offset

		l, err := s.UIntReader.ReadUint16(offset + 1)
		if err != nil {
			return nil, err
		}

		offset += l + 2
	}

	return models, nil
}

func (c *CachedModelConverter) verifyModels() error {
	if c.models != nil {
		return nil
	}

	models, err := c.ModelScanner.Scan()
	if err != nil {
		return err
	}

	c.models = models

	return nil
}

// GetAddress retrieves the starting address of a SunSpec model.
func (c *CachedModelConverter) GetAddress(model uint16) (uint16, error) {
	err := c.verifyModels()
	if err != nil {
		return 0, err
	}

	address, ok := c.models[model]
	if !ok {
		return 0, fmt.Errorf("can't find model")
	}

	return address, nil
}

// HasModel checks if the SunSpec device implements a given model.
func (c *CachedModelConverter) HasModel(model uint16) (bool, error) {
	err := c.verifyModels()
	if err != nil {
		return false, err
	}

	_, ok := c.models[model]
	return ok, nil
}

func (r *ModelReader) ReadPointUint16(model, point uint16) (uint16, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadUint16(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MaxUint16 {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadPointUint32(model, point uint16) (uint32, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadUint32(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MaxUint32 {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadPointUint64(model, point uint16) (uint64, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadUint64(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MaxUint64 {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadPointInt16(model, point uint16) (int16, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadInt16(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MinInt16 {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadPointInt32(model, point uint16) (int32, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadInt32(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MinInt32 {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadPointInt64(model, point uint16) (int64, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadInt64(address + point)
	if err != nil {
		return 0, err
	}

	if val == math.MinInt64 {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadPointFloat32(model, point uint16) (float32, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadFloat32(address + point)
	if err != nil {
		return 0, err
	}

	if math.IsNaN(float64(val)) {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadPointFloat64(model, point uint16) (float64, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return 0, err
	}

	val, err := r.AddressReader.ReadFloat64(address + point)
	if err != nil {
		return 0, err
	}

	if math.IsNaN(val) {
		return 0, fmt.Errorf("uints point with id %d in model with id %d is not implemented", point, model)
	}

	return val, nil
}

func (r *ModelReader) ReadString(model, point, words uint16) (string, error) {
	address, err := r.ModelConverter.GetAddress(model)
	if err != nil {
		return "", err
	}

	val, err := r.AddressReader.ReadString(address+point, words)
	if err != nil {
		return "", err
	}

	return val, nil
}
