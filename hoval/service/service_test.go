package service_test

import "github.com/chrishrb/hoval-gateway/config"

type DatapointProviderMock struct{}

func (m DatapointProviderMock) GetByUnit(unit string) []config.Datapoint {
	return nil
}

func (m DatapointProviderMock) GetByFunction(fg, fn uint8, dpID uint16) *config.Datapoint {
	return &config.Datapoint{
		DatapointName:  "TestDatapoint",
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
		Type:           "List",
	}
}

func (m DatapointProviderMock) GetByName(name string) *config.Datapoint {
	return &config.Datapoint{
		DatapointName:  "TestDatapoint",
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
		Type:           "List",
	}
}
