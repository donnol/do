package do

import (
	"context"
	"net/url"
)

type ParamParser[T ParamData] struct {
	decoder Decoder[T]
}

func NewParamParser[T ParamData](decoder Decoder[T]) *ParamParser[T] {
	paramParser := &ParamParser[T]{
		decoder: decoder,
	}
	return paramParser
}

type ParamData interface {
	url.Values | []byte
}

type Decoder[T ParamData] interface {
	Decode(src T, v any) error
}

type DecodeFunc[T ParamData] func(src T, v any) error

func (f DecodeFunc[T]) Decode(src T, v any) error {
	return f(src, v)
}

// Parse parse data to v with decoder.
func (p *ParamParser[T]) Parse(data T, v any) error {
	err := p.decoder.Decode(data, v)
	if err != nil {
		return err
	}

	return nil
}

// ParseAndCheck parse data to v with decoder and check v if v implement interface{ Check(context.Context) error } or interface{ Check() error }.
func (p *ParamParser[T]) ParseAndCheck(ctx context.Context, data T, v any) error {
	if err := p.Parse(data, v); err != nil {
		return err
	}

	// check param
	var err error
	switch vv := v.(type) {
	case interface{ Check(context.Context) error }:
		err = vv.Check(ctx)
	case interface{ Check() error }:
		err = vv.Check()
	}
	if err != nil {
		return err
	}

	return nil
}
