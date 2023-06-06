package do

import (
	"context"
)

type (
	C = context.Context
	E = error
)

type PipeFunc[I, O any] func(C, I) (O, E)

// Pipe is a pipe run the PipeFuncs in order
func Pipe[B, D, A, R any](
	ctx C,
	b B,
	before PipeFunc[B, D],
	do PipeFunc[D, A],
	after PipeFunc[A, R],
) (r R, err E) {
	// 1
	d, err := before(ctx, b)
	if err != nil {
		return
	}

	// 2
	a, err := do(ctx, d)
	if err != nil {
		return
	}

	// 3
	r, err = after(ctx, a)
	if err != nil {
		return
	}

	return r, nil
}

// Event do something with input I, handle result with success or failed
func Event[I, O, R any](
	ctx C,
	param I,
	do PipeFunc[I, O],
	success func(O) R,
	failed func(E) R,
) (r R) {
	o, err := do(ctx, param)
	if err != nil {
		r = failed(err)
	} else {
		r = success(o)
	}
	return
}
