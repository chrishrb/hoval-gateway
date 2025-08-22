package datatype

import (
	"errors"
	"fmt"
	"math"
)

type Type string

const (
	U8   Type = "U8"
	U16  Type = "U16"
	U32  Type = "U32"
	S8   Type = "S8"
	S16  Type = "S16"
	S32  Type = "S32"
	List Type = "LIST"
)

const RoundPrecision = 2 // Default rounding precision for float values

var (
	ErrInvalidDataType      = errors.New("invalid data type")
	ErrInvalidMessageLength = errors.New("invalid message length")
)

func ToBytes(t Type, data any, decimal int) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}

	// Handle unsigned data
	if t == U8 || t == U16 || t == U32 {
		d, ok := data.(float64)
		if !ok {
			return nil, errors.New("data must be a float64 for unsigned types")
		}
		return UnsignedToBytes(t, d, decimal)
	}

	// Handle signed data
	if t == S8 || t == S16 || t == S32 {
		d, ok := data.(float64)
		if !ok {
			return nil, errors.New("data must be a float64 for signed types")
		}
		return SignedToBytes(t, d, decimal)
	}

	// Handle list data
	if t == List {
		d, ok := data.(uint8)
		if !ok {
			return nil, errors.New("data must be a uint8 for list type")
		}
		return ListToBytes(d)
	}

	return nil, fmt.Errorf("invalid data type: %v", t)
}

func FromBytes(t Type, data []byte, decimal int) (any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Handle unsigned data
	if t == U8 || t == U16 || t == U32 {
		return UnsignedFromBytes(t, data, decimal)
	}

	// Handle signed data
	if t == S8 || t == S16 || t == S32 {
		return SignedFromBytes(t, data, decimal)
	}

	// Handle list data
	if t == List {
		return ListFromBytes(data)
	}

	return nil, fmt.Errorf("invalid data type: %v", t)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
