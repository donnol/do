package do

import (
	"context"
	"net/http"
)

type RouteRegister interface {
	Handle(method, path string, handlers http.HandlerFunc)
}

type RouteHandler[P, R any] interface {
	Parse(req *http.Request, p *P) error
	Write(w http.ResponseWriter, r R, err error)
}

// RegisterRouter register router to RouteRegister with http.HandlerFunc
func RegisterRouter[H RouteRegister, RH RouteHandler[P, R], P, R any](
	g H,
	rh RH,
	method, path string,
	f func(context.Context, P) (R, error),
) {
	g.Handle(method, path, func(w http.ResponseWriter, req *http.Request) {
		var (
			p   P
			r   R
			err error
		)
		defer func() {
			// 返回
			rh.Write(w, r, err)
		}()

		// 参数
		err = rh.Parse(req, &p)
		if err != nil {
			return
		}

		// 业务
		r, err = f(req.Context(), p)
	})
}

type HandlerFunc[T any] func(T)

func (h HandlerFunc[T]) HTTPHandlerFunc(t T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h == nil {
			return
		}

		h(t)
	}
}

type Route[T any] struct {
	Method  string
	Path    string
	Comment string
	Opt     *RouteOption
	Handler HandlerFunc[T]
	Childs  []*Route[T]
}

func (r *Route[T]) WithChilds(childs ...*Route[T]) *Route[T] {
	r.Childs = append(r.Childs, childs...)
	return r
}

func (r *Route[T]) SetChilds(childs ...*Route[T]) *Route[T] {
	r.Childs = childs
	return r
}

func (r *Route[T]) WithOption(opt *RouteOption) *Route[T] {
	r.Opt = opt
	return r
}

func NewRoute[T any](method, path, comment string, h HandlerFunc[T], childs ...*Route[T]) *Route[T] {
	return &Route[T]{
		Method:  method,
		Path:    path,
		Comment: comment,
		Opt:     &RouteOption{NeedLogin: true},
		Handler: h,
		Childs:  childs,
	}
}

type (
	RouteOption struct {
		NeedLogin bool // 是否需要登录

		ParamFormat  string // json | xml
		ResultFormat string // json | xml
		UseBody      bool   // 使用body传递参数
	}
)
