package meter

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	startIdentifier            = "SMA"
	protocolID          uint16 = 0x6069
	measTypeEnergyMeter uint8  = 0x08
	measTypeAverage     uint8  = 0x04
	measTypeVersion     uint8  = 0
	channelInternal     uint8  = 0
	channelOther        uint8  = 144
)

type OBISIdentifier struct {
	Channel, MeasVal, MeasType, Tariff uint8
}

type SoftwareVersion struct {
	Major, Minor, Build, Revision uint8
}

type EnergyMeterTelegram struct {
	SusyID          uint16
	SerialNo        uint32
	MeasuringTime   uint32
	Obis            map[OBISIdentifier]uint64
	SoftwareVersion SoftwareVersion
}

func DecodeTelegram(data []byte) (*EnergyMeterTelegram, error) {
	startIndex := bytes.Index(data, []byte(startIdentifier))
	if startIndex == -1 {
		return nil, fmt.Errorf("couldn't find start index in the datagram")
	}

	buf := bytes.NewBuffer(data[startIndex:])

	buf.Next(16)

	var id uint16
	err := binary.Read(buf, binary.BigEndian, &id)
	if err != nil {
		return nil, err
	}
	if id != protocolID {
		return nil, fmt.Errorf("expected %d as protocol identifier but got %d", protocolID, binary.BigEndian.Uint16(data[16:18]))
	}

	em := &EnergyMeterTelegram{
		Obis: make(map[OBISIdentifier]uint64),
	}

	if err := binary.Read(buf, binary.BigEndian, &em.SusyID); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &em.SerialNo); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &em.MeasuringTime); err != nil {
		return nil, err
	}

	// buffer len greater 4 -> another obis identifier and value
	for buf.Len() > 4 {
		obis := OBISIdentifier{}
		if err := binary.Read(buf, binary.BigEndian, &obis); err != nil {
			return nil, err
		}

		switch obis.Channel {
		case channelInternal:
			if obis.MeasType == measTypeEnergyMeter {
				var val uint64
				if err := binary.Read(buf, binary.BigEndian, &val); err != nil {
					return nil, err
				}
				em.Obis[obis] = val
			} else if obis.MeasType == measTypeAverage {
				var val uint32
				if err := binary.Read(buf, binary.BigEndian, &val); err != nil {
					return nil, err
				}
				em.Obis[obis] = uint64(val)
			} else {
				return nil, fmt.Errorf("unexpected measurement type: %d", obis.MeasType)
			}
		case channelOther:
			if obis.MeasType == measTypeVersion {
				if err := binary.Read(buf, binary.BigEndian, &em.SoftwareVersion); err != nil {
					return nil, err
				}
			}
		default:
			return nil, fmt.Errorf("unexpected channel: %d", obis.Channel)
		}
	}

	return em, nil
}
