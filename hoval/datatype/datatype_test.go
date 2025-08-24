package datatype_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/stretchr/testify/assert"
)

func TestToBytes(t *testing.T) {
	// Test cases table
	testCases := []struct {
		name     string
		dataType datatype.Type
		value    float64
		decimal  int
		expected []byte
		hasError bool
	}{
		// Unsigned tests
		{"U8 conversion", datatype.U8, 20.0, 0, []byte{0x14}, false},
		{"U16 conversion", datatype.U16, 20.20, 2, []byte{0x07, 0xe4}, false},
		{"U32 conversion", datatype.U32, 20.20, 2, []byte{0x0, 0x0, 0x07, 0xe4}, false},

		// Signed tests
		{"S8 conversion", datatype.S8, -10.0, 0, []byte{0xf6}, false},
		{"S16 conversion", datatype.S16, -10.74, 2, []byte{0xfb, 0xce}, false},
		{"S32 conversion", datatype.S32, -10.74, 2, []byte{0xff, 0xff, 0xfb, 0xce}, false},

		// List test
		{"List conversion", datatype.List, 5, 0, []byte{0x05}, false},

		// Invalid type test
		{"Invalid type", "INVALID", 10.0, 0, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := datatype.ToBytes(tc.dataType, tc.value, tc.decimal)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	// Test cases table
	testCases := []struct {
		name     string
		dataType datatype.Type
		data     []byte
		decimal  int
		expected float64
		hasError bool
	}{
		// Unsigned tests
		{"U8 conversion", datatype.U8, []byte{0x14}, 0, 20.0, false},
		{"U16 conversion", datatype.U16, []byte{0x07, 0xe4}, 2, 20.20, false},
		{"U32 conversion", datatype.U32, []byte{0x0, 0x0, 0x07, 0xe4}, 2, 20.20, false},

		// Signed tests
		{"S8 conversion", datatype.S8, []byte{0xf6}, 0, -10.0, false},
		{"S16 conversion", datatype.S16, []byte{0xfb, 0xce}, 2, -10.74, false},
		{"S32 conversion", datatype.S32, []byte{0xff, 0xff, 0xfb, 0xce}, 2, -10.74, false},

		// List test
		{"List conversion", datatype.List, []byte{0x05}, 0, 5.0, false},

		// Empty data test
		{"Empty data", datatype.U8, []byte{}, 0, 0.0, false},

		// Invalid type test
		{"Invalid type", "INVALID", []byte{0x01}, 0, 0.0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := datatype.FromBytes(tc.dataType, tc.data, tc.decimal)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestRoundFloat(t *testing.T) {
	// Test cases table
	testCases := []struct {
		name      string
		value     float64
		precision uint
		expected  float64
	}{
		{"Round to 2 decimal places", 10.345, 2, 10.35},
		{"Round down to 2 decimal places", 10.344, 2, 10.34},
		{"Round to 0 decimal places", 10.5, 0, 11.0},
		{"Round negative number", -10.345, 2, -10.35},
		{"Round negative number down", -10.344, 2, -10.34},
		{"Round with higher precision", 10.12345678, 6, 10.123457},
		{"Round with default precision", 10.345, datatype.RoundPrecision, 10.35},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := datatype.RoundFloat(tc.value, tc.precision)
			assert.Equal(t, tc.expected, result)
		})
	}
}
