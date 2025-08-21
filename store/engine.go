package store

import "github.com/chrishrb/hoval-gateway/hoval"

type DeviceStore interface {
	SetDevice(address uint32, device *hoval.Device)
	LookupDevice(address uint32) *hoval.Device
	ListDevices() []*hoval.Device
}

type Engine interface {
	DeviceStore
}
