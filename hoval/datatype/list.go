package datatype

func ListToBytes(data uint8) ([]byte, error) {
	return []byte{data}, nil
}

func ListFromBytes(data []byte) (uint8, error) {
	if len(data) != 1 {
		return 0, ErrInvalidMessageLength
	}
	return uint8(data[0]), nil
}
