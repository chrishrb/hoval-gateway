package transport

// Receiver is the interface implemented in order to receive messages
type Receiver interface {
	Connect(errCh chan error)
}
