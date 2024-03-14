package do

import (
	"context"
)

type (
	C = context.Context
	E = error
)

type PipeFunc[I, O any] func(C, I) (O, E)

func PipeToLogic[I, O any](pf PipeFunc[I, O]) Logic[I, O] {
	return Logic[I, O](pf)
}

func PipeFromLogic[I, O any](logic Logic[I, O]) PipeFunc[I, O] {
	return PipeFunc[I, O](logic)
}

func PipeFrom[I, O any](f func(C, I) (O, E)) PipeFunc[I, O] {
	return PipeFunc[I, O](f)
}

type Piper[I, O any] interface {
	Run(C, I) (O, E)
}

func (pf PipeFunc[I, O]) Run(c C, i I) (r O, e E) {
	return pf(c, i)
}

var (
	_ Piper[struct{}, struct{}] = PipeFunc[struct{}, struct{}](nil)
)

func Pipes[B, D, A, R any](
	ctx C,
	b B,
	before Piper[B, D],
	do Piper[D, A],
	after Piper[A, R],
) (r R, e E) {
	// 1
	d, err := before.Run(ctx, b)
	if err != nil {
		return
	}

	// 2
	a, err := do.Run(ctx, d)
	if err != nil {
		return
	}

	// 3
	r, err = after.Run(ctx, a)
	if err != nil {
		return
	}

	return r, nil
}

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
	success PipeFunc[O, R],
	failed PipeFunc[E, R],
) (r R, err E) {
	var o O
	var berr E

	defer func() {
		if berr != nil {
			r, err = failed(ctx, berr)
		} else {
			r, err = success(ctx, o)
		}
	}()

	o, berr = do(ctx, param)

	return
}

type (
	EventFunc[I, O, R any] func(
		ctx C,
		param I,
		do PipeFunc[I, O],
		success PipeFunc[O, R],
		failed PipeFunc[E, R],
	) (r R, err E)

	EventEntity[I, O, R any] struct {
		Param   I
		Do      PipeFunc[I, O]
		Success PipeFunc[O, R]
		Failed  PipeFunc[E, R]

		Handler func(R, E)
	}
)

var (
	_ EventFunc[int, int, int] = Event[int, int, int]
)

func EventLoop[I, O, R any](ctx C, n int) (chan<- EventEntity[I, O, R], chan<- struct{}) {
	innerch := make(chan EventEntity[I, O, R], n)
	stopch := make(chan struct{}, 1)

	go func() {
		defer func() {
			close(innerch)
			close(stopch)
		}()

		for {
			select {
			case event := <-innerch:
				r, err := Event(ctx, event.Param, event.Do, event.Success, event.Failed)
				event.Handler(r, err)
			case <-stopch:
				return
			}
		}
	}()

	return (chan<- EventEntity[I, O, R])(innerch), (chan<- struct{})(stopch)
}
