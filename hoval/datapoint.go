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
