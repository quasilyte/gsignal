package gsignal

// ConnectOneShot will call given function at most once, when event is emitted.
// It may call that function zero times if conn is disposed before any event is emitted.
// A nil connection will keep the function connected until the event is triggered.
func ConnectOneShot[T any](event *Event[T], conn connection, slot func(T)) {
	oneshot := &oneshotConnector[T]{conn: conn}
	event.Connect(oneshot, func(arg T) {
		oneshot.fired = true
		slot(arg)
	})
}

type oneshotConnector[T any] struct {
	conn  connection
	fired bool
}

func (c *oneshotConnector[T]) IsDisposed() bool {
	if c.fired {
		return true
	}
	return c.conn != nil && c.conn.IsDisposed()
}
