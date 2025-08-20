package can

type connectionDetails struct {
	device string
}

type Opt[T any] func(h *T)

func WithCANDevice[T Sender | Consumer](device string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Sender:
			x.device = device
		case *Consumer:
			x.device = device
		}
	}
}
