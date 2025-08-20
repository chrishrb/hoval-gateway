package datatype

func StringToBytes(data string) ([]byte, error) {
	return []byte(data), nil
}

func StringFromBytes(data []byte) (string, error) {
	return string(data), nil
}
