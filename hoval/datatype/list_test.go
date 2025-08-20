package datatype_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/stretchr/testify/assert"
)

func TestListToBytes(t *testing.T) {
	data, err := datatype.ListToBytes(5)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x05}, data)
}

func TestListFromBytes(t *testing.T) {
	data, err := datatype.ListFromBytes([]byte{0x05})
	assert.NoError(t, err)
	assert.Equal(t, uint8(5), data)

	data, err = datatype.ListFromBytes([]byte{0x20, 0x20})
	assert.Error(t, err)
}
