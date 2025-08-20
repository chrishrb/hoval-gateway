package hoval

import (
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
)

// Datapoint represents a datapoint in Hoval communication.
//
// Example:
// Name: Ventilation mode selection
// FunctionGroup: 50
// FunctionNumber: 0
// DatapointID: 40650
// Type: LIST
type Datapoint struct {
	FunctionGroup  FunctionGroup
	FunctionNumber uint8
	DatapointID    uint16
	DatapointType  datatype.Type
	DecimalPlaces  int
}

func NewDatapoint(
	functionGroup FunctionGroup,
	functionNumber uint8,
	datapointID uint16,
	datapointType datatype.Type,
	decimalPlaces int,
) *Datapoint {
	return &Datapoint{
		FunctionGroup:  functionGroup,
		FunctionNumber: functionNumber,
		DatapointID:    datapointID,
		DatapointType:  datapointType,
		DecimalPlaces:  decimalPlaces,
	}
}

type DatapointStore struct {
	datapoints map[string]Datapoint
}

func NewDatapointStore() *DatapointStore {
	return &DatapointStore{
		datapoints: make(map[string]Datapoint),
	}
}

func (ds *DatapointStore) Add(name string, dp Datapoint) {
	ds.datapoints[name] = dp
}

func (ds *DatapointStore) LookupByName(name string) *Datapoint {
	result, exists := ds.datapoints[name]
	if !exists {
		return nil
	}
	return &result
}

func (ds *DatapointStore) LookupByIdentifier(fg FunctionGroup, fn uint8, dpID uint16) *Datapoint {
	for _, dp := range ds.datapoints {
		if dp.FunctionGroup == fg && dp.FunctionNumber == fn && dp.DatapointID == dpID {
			return &dp
		}
	}
	return nil
}
