package inmemory

import "github.com/chrishrb/hoval-gateway/hoval"

func (s *Store) SetDevice(address uint32, device *hoval.Device) {
	s.Lock()
	defer s.Unlock()

	s.devices[address] = device
}

func (s *Store) LookupDevice(address uint32) *hoval.Device {
	s.Lock()
	defer s.Unlock()

	device, exists := s.devices[address]
	if !exists {
		return nil
	}
	return device
}

func (s *Store) ListDevices() []*hoval.Device {
	s.Lock()
	defer s.Unlock()

	devices := make([]*hoval.Device, 0, len(s.devices))
	for _, device := range s.devices {
		devices = append(devices, device)
	}
	return devices
}
