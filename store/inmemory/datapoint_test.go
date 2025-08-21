package inmemory_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestSetDatapoint(t *testing.T) {
	store := inmemory.NewStore(nil)

	datapoint := &hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 1,
		DatapointID:    100,
		DatapointType:  datatype.U8,
		DecimalPlaces:  0,
	}

	// Test setting a datapoint
	store.SetDatapoint("test-datapoint", datapoint)

	// Verify the datapoint was stored
	retrievedDatapoint := store.LookupDatapointByName("test-datapoint")
	assert.NotNil(t, retrievedDatapoint)
	assert.Equal(t, hoval.FunctionGroup(50), retrievedDatapoint.FunctionGroup)
	assert.Equal(t, uint8(1), retrievedDatapoint.FunctionNumber)
	assert.Equal(t, uint16(100), retrievedDatapoint.DatapointID)
	assert.Equal(t, datatype.U8, retrievedDatapoint.DatapointType)
	assert.Equal(t, 0, retrievedDatapoint.DecimalPlaces)
}

func TestSetDatapointOverwrite(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Create two datapoints with different properties
	datapoint1 := &hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 1,
		DatapointID:    100,
		DatapointType:  datatype.U8,
		DecimalPlaces:  0,
	}
	datapoint2 := &hoval.Datapoint{
		FunctionGroup:  51,
		FunctionNumber: 2,
		DatapointID:    200,
		DatapointType:  datatype.S16,
		DecimalPlaces:  2,
	}

	// Set first datapoint
	store.SetDatapoint("test-datapoint", datapoint1)
	retrieved1 := store.LookupDatapointByName("test-datapoint")
	assert.Equal(t, hoval.FunctionGroup(50), retrieved1.FunctionGroup)
	assert.Equal(t, uint16(100), retrieved1.DatapointID)

	// Overwrite with second datapoint
	store.SetDatapoint("test-datapoint", datapoint2)
	retrieved2 := store.LookupDatapointByName("test-datapoint")
	assert.Equal(t, hoval.FunctionGroup(51), retrieved2.FunctionGroup)
	assert.Equal(t, uint16(200), retrieved2.DatapointID)
}

func TestLookupDatapointByName(t *testing.T) {
	store := inmemory.NewStore(nil)

	datapoint := &hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 1,
		DatapointID:    100,
		DatapointType:  datatype.U8,
		DecimalPlaces:  0,
	}

	// Test lookup of existing datapoint
	store.SetDatapoint("existing-datapoint", datapoint)
	retrievedDatapoint := store.LookupDatapointByName("existing-datapoint")
	assert.NotNil(t, retrievedDatapoint)
	assert.Equal(t, hoval.FunctionGroup(50), retrievedDatapoint.FunctionGroup)
	assert.Equal(t, uint8(1), retrievedDatapoint.FunctionNumber)
	assert.Equal(t, uint16(100), retrievedDatapoint.DatapointID)

	// Test lookup of non-existing datapoint
	nonExistentDatapoint := store.LookupDatapointByName("non-existent")
	assert.Nil(t, nonExistentDatapoint)
}

func TestLookupDatapointByNameEmptyStore(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Test lookup in empty store
	datapoint := store.LookupDatapointByName("any-name")
	assert.Nil(t, datapoint)
}

func TestLookupDatapointByIdentifier(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Create multiple datapoints with different identifiers
	datapoint1 := &hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 1,
		DatapointID:    100,
		DatapointType:  datatype.U8,
		DecimalPlaces:  0,
	}
	datapoint2 := &hoval.Datapoint{
		FunctionGroup:  51,
		FunctionNumber: 2,
		DatapointID:    200,
		DatapointType:  datatype.S16,
		DecimalPlaces:  1,
	}
	datapoint3 := &hoval.Datapoint{
		FunctionGroup:  52,
		FunctionNumber: 3,
		DatapointID:    300,
		DatapointType:  datatype.U32,
		DecimalPlaces:  2,
	}

	// Store the datapoints
	store.SetDatapoint("dp1", datapoint1)
	store.SetDatapoint("dp2", datapoint2)
	store.SetDatapoint("dp3", datapoint3)

	// Test lookup by exact identifier match
	found1 := store.LookupDatapointByIdentifier(50, 1, 100)
	assert.NotNil(t, found1)
	assert.Equal(t, hoval.FunctionGroup(50), found1.FunctionGroup)
	assert.Equal(t, uint8(1), found1.FunctionNumber)
	assert.Equal(t, uint16(100), found1.DatapointID)

	found2 := store.LookupDatapointByIdentifier(51, 2, 200)
	assert.NotNil(t, found2)
	assert.Equal(t, hoval.FunctionGroup(51), found2.FunctionGroup)
	assert.Equal(t, uint8(2), found2.FunctionNumber)
	assert.Equal(t, uint16(200), found2.DatapointID)

	// Test lookup with non-existing identifier
	notFound := store.LookupDatapointByIdentifier(1, 99, 999)
	assert.Nil(t, notFound)
}

func TestLookupDatapointByIdentifierPartialMatch(t *testing.T) {
	store := inmemory.NewStore(nil)

	datapoint := &hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 1,
		DatapointID:    100,
		DatapointType:  datatype.U8,
		DecimalPlaces:  0,
	}

	store.SetDatapoint("test-dp", datapoint)

	// Test partial matches (should not find anything)

	// Wrong function group
	notFound1 := store.LookupDatapointByIdentifier(1, 1, 100)
	assert.Nil(t, notFound1)

	// Wrong function number
	notFound2 := store.LookupDatapointByIdentifier(50, 2, 100)
	assert.Nil(t, notFound2)

	// Wrong datapoint ID
	notFound3 := store.LookupDatapointByIdentifier(50, 1, 200)
	assert.Nil(t, notFound3)
}

func TestLookupDatapointByIdentifierEmptyStore(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Test lookup in empty store
	datapoint := store.LookupDatapointByIdentifier(50, 1, 100)
	assert.Nil(t, datapoint)
}
