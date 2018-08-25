package trace

import (
	"fmt"
	"io"
)

type tracer struct {
	out io.Writer
}

// Tracer 인터페이스와 일치
func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// Tracer is the interface that describes an object capable of tracing events throughout code
type Tracer interface {
	Trace(...interface{}) // Trace method accepts zero or more arguments of any type
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to Trace
func Off() Tracer {
	return &nilTracer{}
}
