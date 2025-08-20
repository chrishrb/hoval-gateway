package datatype_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/stretchr/testify/assert"
)

func TestStringToBytes(t *testing.T) {
	// U8
	data, err := datatype.StringToBytes("hello world")
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64}, data)
}

func TestStringFromBytes(t *testing.T) {
	data, err := datatype.StringFromBytes([]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64})
	assert.NoError(t, err)
	assert.Equal(t, "hello world", data)
}
