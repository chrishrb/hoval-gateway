package inmemory_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestSetDevice(t *testing.T) {
	store := inmemory.NewStore(nil)

	device := &hoval.Device{
		Address: 100,
	}

	// Test setting a device
	store.SetDevice(1, device)

	// Verify the device was stored
	retrievedDevice := store.LookupDevice(1)
	assert.NotNil(t, retrievedDevice)
	assert.Equal(t, uint32(100), retrievedDevice.Address)
}

func TestSetDeviceOverwrite(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Create two devices with different addresses
	device1 := &hoval.Device{Address: 100}
	device2 := &hoval.Device{Address: 200}

	// Set first device
	store.SetDevice(1, device1)
	retrieved1 := store.LookupDevice(1)
	assert.Equal(t, uint32(100), retrieved1.Address)

	// Overwrite with second device
	store.SetDevice(1, device2)
	retrieved2 := store.LookupDevice(1)
	assert.Equal(t, uint32(200), retrieved2.Address)
}

func TestLookupDevice(t *testing.T) {
	store := inmemory.NewStore(nil)

	device := &hoval.Device{
		Address: 100,
	}

	// Test lookup of existing device
	store.SetDevice(1, device)
	retrievedDevice := store.LookupDevice(1)
	assert.NotNil(t, retrievedDevice)
	assert.Equal(t, uint32(100), retrievedDevice.Address)

	// Test lookup of non-existing device
	nonExistentDevice := store.LookupDevice(255)
	assert.Nil(t, nonExistentDevice)
}

func TestLookupDeviceEmptyStore(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Test lookup in empty store
	device := store.LookupDevice(1)
	assert.Nil(t, device)
}

func TestListDevices(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Test listing devices in empty store
	devices := store.ListDevices()
	assert.NotNil(t, devices)
	assert.Empty(t, devices)

	// Add some devices
	device1 := &hoval.Device{Address: 100}
	device2 := &hoval.Device{Address: 200}
	device3 := &hoval.Device{Address: 300}

	store.SetDevice(1, device1)
	store.SetDevice(2, device2)
	store.SetDevice(3, device3)

	// Test listing multiple devices
	devices = store.ListDevices()
	assert.NotNil(t, devices)
	assert.Len(t, devices, 3)

	// Verify all devices are present (order might vary since it's from a map)
	addresses := make([]uint32, len(devices))
	for i, device := range devices {
		addresses[i] = device.Address
	}
	assert.Contains(t, addresses, uint32(100))
	assert.Contains(t, addresses, uint32(200))
	assert.Contains(t, addresses, uint32(300))
}

func TestListDevicesAfterOverwrite(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Add devices
	device1 := &hoval.Device{Address: 100}
	device2 := &hoval.Device{Address: 200}
	device3 := &hoval.Device{Address: 300}

	store.SetDevice(1, device1)
	store.SetDevice(2, device2)
	store.SetDevice(1, device3) // Overwrite device at address 1

	// Should still have only 2 devices
	devices := store.ListDevices()
	assert.Len(t, devices, 2)

	// Verify the correct devices are present
	addresses := make([]uint32, len(devices))
	for i, device := range devices {
		addresses[i] = device.Address
	}
	assert.Contains(t, addresses, uint32(200))
	assert.Contains(t, addresses, uint32(300))
	assert.NotContains(t, addresses, uint32(100))
}

func TestDeviceStoreConcurrency(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Test concurrent access to ensure thread safety
	done := make(chan bool, 2)

	// Goroutine 1: Set devices
	go func() {
		for i := range uint32(10) {
			device := &hoval.Device{Address: uint32(i * 10)}
			store.SetDevice(i, device)
		}
		done <- true
	}()

	// Goroutine 2: List devices
	go func() {
		for range 10 {
			store.ListDevices()
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// Verify final state
	devices := store.ListDevices()
	assert.Len(t, devices, 10)
}

func TestDeviceNameMapping(t *testing.T) {
	store := inmemory.NewStore(nil)

	// Test devices with different address ranges to verify GetName() functionality
	testCases := []struct {
		address      uint32
		expectedName string
	}{
		{5, "WEZ"},     // 1-16
		{70, "SOL"},    // 65-80
		{135, "PS"},    // 129-144
		{200, "FW"},    // 193-205
		{260, "HW/WW"}, // 257-272
		{390, "MWA"},   // 385-400
		{450, "GLT"},   // 449-464
		{520, "HV"},    // 513-528
		{1030, "BM"},   // 1025-1087
		{1155, "GW"},   // 1153-1160
		{9999, "Unknown"},
	}

	for _, tc := range testCases {
		device := &hoval.Device{Address: tc.address}
		store.SetDevice(uint32(tc.address%256), device)

		retrievedDevice := store.LookupDevice(uint32(tc.address % 256))
		assert.NotNil(t, retrievedDevice)
		assert.Equal(t, tc.expectedName, retrievedDevice.GetName())
	}
}
