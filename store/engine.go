package store

import "github.com/chrishrb/hoval-gateway/hoval"

type DatapointStore interface {
	SetDatapoint(name string, dp *hoval.Datapoint)
	LookupDatapointByName(name string) *hoval.Datapoint
	LookupDatapointByIdentifier(fg hoval.FunctionGroup, fn uint8, dpID uint16) *hoval.Datapoint
}

type DeviceStore interface {
	SetDevice(address uint32, device *hoval.Device)
	LookupDevice(address uint32) *hoval.Device
	ListDevices() []*hoval.Device
}

type Engine interface {
	DatapointStore
	DeviceStore
}
