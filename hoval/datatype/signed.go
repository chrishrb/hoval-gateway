package datatype

import (
	"encoding/binary"
	"errors"
	"math"
)

func SignedToBytes(t Type, data float64, decimal int) ([]byte, error) {
	scaledValue := int64(data * math.Pow10(decimal))

	switch t {
	case S8:
		if scaledValue > math.MaxInt8 || scaledValue < math.MinInt8 {
			return nil, errors.New("value exceeds maximum for TypeS8")
		}
		return []byte{byte(uint64(scaledValue))}, nil
	case S16:
		if scaledValue > math.MaxInt16 || scaledValue < math.MinInt16 {
			return nil, errors.New("value exceeds maximum for TypeS16")
		}
		result := make([]byte, 2)
		binary.BigEndian.PutUint16(result, uint16(scaledValue))
		return result, nil
	case S32:
		if scaledValue > math.MaxInt32 || scaledValue < math.MinInt32 {
			return nil, errors.New("value exceeds maximum for TypeU32")
		}
		result := make([]byte, 4)
		binary.BigEndian.PutUint32(result, uint32(scaledValue))
		return result, nil
	}

	return nil, ErrInvalidDataType
}

func SignedFromBytes(t Type, data []byte, decimal int) (float64, error) {
	var val int64

	switch t {
	case S8:
		if len(data) != 1 {
			return 0, ErrInvalidMessageLength
		}
		val = int64(int8(data[0]))
	case S16:
		if len(data) != 2 {
			return 0, ErrInvalidMessageLength
		}
		val = int64(int16(binary.BigEndian.Uint16(data)))
	case S32:
		if len(data) != 4 {
			return 0, ErrInvalidMessageLength
		}
		val = int64(int32(binary.BigEndian.Uint32(data)))
	}

	return roundFloat(float64(val) * math.Pow10(-decimal), RoundPrecision), nil
}
