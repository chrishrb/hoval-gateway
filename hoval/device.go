package hoval

type Device struct {
	Address uint32
}

func NewDevice(address uint32) *Device {
	return &Device{
		Address: address,
	}
}

func (d *Device) GetName() string {
	if d.Address >= 1 && d.Address <= 16 {
		return "WEZ"
	} else if d.Address >= 65 && d.Address <= 80 {
		return "SOL"
	} else if d.Address >= 129 && d.Address <= 144 {
		return "PS"
	} else if d.Address >= 193 && d.Address <= 205 {
		return "FW"
	} else if d.Address >= 257 && d.Address <= 272 {
		return "HW/WW"
	} else if d.Address >= 385 && d.Address <= 400 {
		return "MWA"
	} else if d.Address >= 449 && d.Address <= 464 {
		return "GLT"
	} else if d.Address >= 513 && d.Address <= 528 {
		return "HV"
	} else if d.Address >= 1025 && d.Address <= 1087 {
		return "BM"
	} else if d.Address >= 1153 && d.Address <= 1160 {
		return "GW"
	} else if d.Address == 2047 {
		return "Broadcast"
	} else {
		return "Unknown"
	}
}
