package service_test

import "github.com/chrishrb/hoval-gateway/hoval/datapoint"

type DatapointProviderMock struct{}

func (m DatapointProviderMock) GetByUnit(unit string) []datapoint.Datapoint {
	return nil
}

func (m DatapointProviderMock) GetByFunction(fg, fn uint8, dpID uint16) *datapoint.Datapoint {
	if fg == 50 && fn == 0 && dpID == 40650 {
		return &datapoint.Datapoint{
			DatapointName:  "TestDatapoint",
			FunctionGroup:  50,
			FunctionNumber: 0,
			DatapointID:    40650,
			TypeName:       "LIST",
			Writable:       true,
		}
	}
	return nil
}

func (m DatapointProviderMock) GetByName(name string) *datapoint.Datapoint {
	if name == "TestDatapoint" {
		return &datapoint.Datapoint{
			DatapointName:  "TestDatapoint",
			FunctionGroup:  50,
			FunctionNumber: 0,
			DatapointID:    40650,
			TypeName:       "LIST",
		}
	}
	return nil
}
