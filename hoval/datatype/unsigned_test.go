package datatype_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/stretchr/testify/assert"
)

func TestUnsignedToBytes(t *testing.T) {
	// U8
	data, err := datatype.UnsignedToBytes(datatype.U8, 20.0, 0)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x14}, data)

	// U16
	data, err = datatype.UnsignedToBytes(datatype.U16, 20.20, 2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x07, 0xe4}, data)

	// U32
	data, err = datatype.UnsignedToBytes(datatype.U32, 20.20, 2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x0, 0x0, 0x07, 0xe4}, data)
}

func TestUnsignedFromBytes(t *testing.T) {
	// U8
	data, err := datatype.UnsignedFromBytes(datatype.U8, []byte{0x14}, 0)
	assert.NoError(t, err)
	assert.Equal(t, 20.0, data)

	// U16
	data, err = datatype.UnsignedFromBytes(datatype.U16, []byte{0x07, 0xe4}, 2)
	assert.NoError(t, err)
	assert.Equal(t, 20.20, data)

	// U32
	data, err = datatype.UnsignedFromBytes(datatype.U32, []byte{0x0, 0x0, 0x07, 0xe4}, 2)
	assert.NoError(t, err)
	assert.Equal(t, 20.20, data)
}
