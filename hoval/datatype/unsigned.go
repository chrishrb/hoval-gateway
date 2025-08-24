package datatype

import (
	"encoding/binary"
	"errors"
	"math"
)

func UnsignedToBytes(t Type, data float64, decimal int) ([]byte, error) {
	scaledValue := uint64(data * math.Pow10(decimal))

	switch t {
	case U8:
		if scaledValue > math.MaxUint8 {
			return nil, errors.New("value exceeds maximum for TypeU8")
		}
		return []byte{byte(uint64(scaledValue))}, nil
	case U16:
		if scaledValue > math.MaxUint16 {
			return nil, errors.New("value exceeds maximum for TypeU16")
		}
		result := make([]byte, 2)
		binary.BigEndian.PutUint16(result, uint16(scaledValue))
		return result, nil
	case U32:
		if scaledValue > math.MaxUint32 {
			return nil, errors.New("value exceeds maximum for TypeU32")
		}
		result := make([]byte, 4)
		binary.BigEndian.PutUint32(result, uint32(scaledValue))
		return result, nil
	}

	return nil, ErrInvalidDataType
}

func UnsignedFromBytes(t Type, data []byte, decimal int) (float64, error) {
	var val float64

	switch t {
	case U8:
		if len(data) != 1 {
			return 0, ErrInvalidMessageLength
		}
		val = float64(data[0]) / math.Pow10(decimal)
	case U16:
		if len(data) != 2 {
			return 0, ErrInvalidMessageLength
		}
		val = float64(binary.BigEndian.Uint16(data)) / math.Pow10(decimal)
	case U32:
		if len(data) != 4 {
			return 0, ErrInvalidMessageLength
		}
		val = float64(binary.BigEndian.Uint32(data)) / math.Pow10(decimal)
	default:
		return 0, ErrInvalidDataType
	}

	return RoundFloat(val, RoundPrecision), nil
}
