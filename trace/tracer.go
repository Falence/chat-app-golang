package trace

import "io"

// Tracer is the interface that describes an object capable of
// tracing events throughout code
type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return nil
}