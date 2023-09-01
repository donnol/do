package do

import (
	"context"
	"fmt"
)

type ContextHelper[K ~struct{}, V any] struct {
	k K
}

func (h ContextHelper[K, V]) WithValue(ctx context.Context, v V) context.Context {
	return context.WithValue(ctx, h.k, v)
}

func (h ContextHelper[K, V]) Value(ctx context.Context) (v V, ok bool) {
	v, ok = ctx.Value(h.k).(V)
	return
}

func (h ContextHelper[K, V]) MustValue(ctx context.Context) (v V) {
	v, ok := h.Value(ctx)
	if !ok {
		panic(fmt.Errorf("context can't find value of %#v", h.k))
	}
	return
}
