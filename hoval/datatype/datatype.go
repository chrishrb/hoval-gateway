package datatype

import (
	"errors"
)

type Type int

const (
	U8 Type = iota
	U16
	U32
	S8
	S16
	S32
	List
	String
)

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

	// Handle string data
	if t == String {
		d, ok := data.(string)
		if !ok {
			return nil, errors.New("data must be a string for string type")
		}
		return StringToBytes(d)
	}

	return nil, ErrInvalidDataType
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

	// Handle string data
	if t == String {
		return StringFromBytes(data)
	}

	return nil, ErrInvalidDataType
}
