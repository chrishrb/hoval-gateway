package pubsub

type Receiver interface {
	Connect(errCh chan error)
}
