package transport

const MaxDataLength = 8

type Data [MaxDataLength]byte

type Message struct {
	ID     uint32
	Length uint8
	Data   Data
}

func NewMessage(
	id uint32,
	length uint8,
	data Data,
) *Message {
	return &Message{
		ID:     id,
		Length: length,
		Data:   data,
	}
}
