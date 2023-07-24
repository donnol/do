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
