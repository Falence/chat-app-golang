package trace

// Tracer is the interface that describes an object capable of
// tracing events throughout code
type Tracer interface {
	Trace(...interface{})
}