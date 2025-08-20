package datatype_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/stretchr/testify/assert"
)

func TestSignedToBytes(t *testing.T) {
	// S8
	data, err := datatype.SignedToBytes(datatype.S8, -20.0, 0)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xec}, data)

	// S16
	data, err = datatype.SignedToBytes(datatype.S16, -20.23, 2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xf8, 0x19}, data)

	data, err = datatype.SignedToBytes(datatype.S16, 20.23, 2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x07, 0xe7}, data)

	// S32
	data, err = datatype.SignedToBytes(datatype.S32, -20.20, 2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xff, 0xff, 0xf8, 0x1c}, data)

	data, err = datatype.SignedToBytes(datatype.S32, 20.20, 2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x0, 0x0, 0x07, 0xe4}, data)
}

func TestSignedFromBytes(t *testing.T) {
	// S8
	data, err := datatype.SignedFromBytes(datatype.S8, []byte{0x14}, 0)
	assert.NoError(t, err)
	assert.Equal(t, 20.0, data)

	data, err = datatype.SignedFromBytes(datatype.S8, []byte{0xec}, 0)
	assert.NoError(t, err)
	assert.Equal(t, -20.0, data)

	// S16
	data, err = datatype.SignedFromBytes(datatype.S16, []byte{0xf8, 0x19}, 2)
	assert.NoError(t, err)
	assert.Equal(t, -20.23, data)

	data, err = datatype.SignedFromBytes(datatype.S16, []byte{0x07, 0xe7}, 2)
	assert.NoError(t, err)
	assert.Equal(t, 20.23, data)

	// S32
	data, err = datatype.SignedFromBytes(datatype.S32, []byte{0xff, 0xff, 0xf8, 0x1c}, 2)
	assert.NoError(t, err)
	assert.Equal(t, -20.20, data)

	data, err = datatype.SignedFromBytes(datatype.S32, []byte{0x0, 0x0, 0x07, 0xe4}, 2)
	assert.NoError(t, err)
	assert.Equal(t, 20.20, data)
}
