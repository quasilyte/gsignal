package gsignal

// Void is a helper type that is recommended to be used instead of `struct{}`.
// You may want to use Void as the event type parameter when there is no
// useful data to be transmitted.
type Void struct{}

// Event is a slot-signal container.
// It holds all currect event listeners and invokes their callbacks
// when the associated event is triggered.
// An event is triggered when Emit() method is called.
//
// If you need 0 arguments callback, use Void type for the argument.
// If you need more than 1 argument in your callback, use tuple helper package.
// For example, a tuple.Value3[int, float, string] can be used to pass
// three arguments to your callback.
type Event[T any] struct {
	handlers []eventHandler[T]
}

// Reset disconnects all connected event listeners (slot functions).
// After this operation the Event object in its zero-like state, ready to be re-used.
func (e *Event[T]) Reset() {
	e.handlers = e.handlers[:0]
}

// Flush makes sure disposed event listeners are removed from the object.
//
// Normally, this happens as a part of every Emit, but some use cases
// may have rare [Emits] and frequent [Connect].
//
// Don't call this method unless you're certain that you need it.
func (e *Event[T]) Flush() {
	// This method is slightly faster than the self-append alternative.
	length := 0
	for _, h := range e.handlers {
		if h.c != nil && h.c.IsDisposed() {
			continue
		}
		e.handlers[length] = h
		length++
	}
	e.handlers = e.handlers[:length]
}

// Forward is a convenience wrapper over connecting to e and calling emit on e2
// with the same arguments.
// In other words, this method sets up a forwarding from e to e2.
// When e does Emit(), e2 would receive it and Emit() as well.
//
// The conn argument is used for the underlying Connect() call.
func (e *Event[T]) Forward(conn connection, e2 *Event[T]) {
	e.Connect(conn, func(arg T) {
		e2.Emit(arg)
	})
}

// Connect adds an event listener that will be called for every Emit called for this event.
// When connection is disposed, an associated callback will be unregistered.
// If this connection should be persistent, pass a nil value as conn.
// For a non-nil conn, it's possible to disconnect from event by using Disconnect method.
func (e *Event[T]) Connect(conn connection, slot func(arg T)) {
	e.handlers = append(e.handlers, eventHandler[T]{
		c: conn,
		f: slot,
	})
}

// Disconnect removes an event listener identified by this connection.
// Note that you can't disconnect a listener that was connected with nil connection object.
func (e *Event[T]) Disconnect(conn connection) {
	for i, h := range e.handlers {
		if h.c == conn {
			e.handlers[i].c = theRemovedConnection
			break
		}
	}
}

// Emit triggers the associated event and calls all active callbacks with provided argument.
func (e *Event[T]) Emit(arg T) {
	// This method is slightly faster than the self-append alternative.
	length := 0
	for _, h := range e.handlers {
		if h.c != nil && h.c.IsDisposed() {
			continue
		}
		h.f(arg)
		e.handlers[length] = h
		length++
	}
	e.handlers = e.handlers[:length]
}

// IsEmpty is a shorthand for NumConnections==0.
func (e *Event[T]) IsEmpty() bool {
	return len(e.handlers) == 0
}

// NumConnections reports the number of alive event connections.
func (e *Event[T]) NumConnections() int {
	return len(e.handlers)
}

type eventHandler[T any] struct {
	c connection
	f func(T)
}

type connection interface {
	IsDisposed() bool
}

type removedConnection struct{}

func (r *removedConnection) IsDisposed() bool { return true }

var theRemovedConnection = &removedConnection{}
