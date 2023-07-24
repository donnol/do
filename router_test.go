package do

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type routeRegister struct {
}

func (r *routeRegister) Handle(method, path string, handler http.HandlerFunc) {

	http.Handle(path, handler)
}

type RouteDefaultHandler[P, R any] struct {
}

func (rh *RouteDefaultHandler[P, R]) Parse(req *http.Request, p *P) error {
	switch req.Method {
	case http.MethodGet, http.MethodDelete:
		// TODO:
	case http.MethodPost, http.MethodPut:
		data, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		json.Unmarshal(data, p)
	}
	return nil
}

func (rh *RouteDefaultHandler[P, R]) Write(w http.ResponseWriter, r R, err error) {
	w.WriteHeader(200)
	if err != nil {
		w.Write([]byte(`{"msg": "failed"}`))
		return
	}
	data, err := json.Marshal(r)
	if err != nil {
		w.Write([]byte(`{"msg": "json encode failed"}`))
		return
	}
	w.Write(data)
}

type (
	rparam struct {
		Id string `json:"id"`
	}
	rresult struct {
		Name string `json:"name"`
	}
)

func TestRegisterRouter(t *testing.T) {
	type args struct {
		g      *routeRegister
		rh     *RouteDefaultHandler[rparam, rresult]
		method string
		path   string
		f      func(context.Context, rparam) (rresult, error)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				g:      &routeRegister{},
				rh:     &RouteDefaultHandler[rparam, rresult]{},
				method: http.MethodPost,
				path:   "/user",
				f: func(ctx context.Context, p rparam) (r rresult, err error) {
					r.Name = p.Id + "-" + strconv.Itoa(rand.Intn(100))
					return
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterRouter(tt.args.g, tt.args.rh, tt.args.method, tt.args.path, tt.args.f)

			go func() {
				http.ListenAndServe(":12332", http.DefaultServeMux)
			}()
			time.Sleep(time.Millisecond * 200)

			buf := new(bytes.Buffer)
			p := rparam{Id: "1000"}
			err := json.NewEncoder(buf).Encode(p)
			if err != nil {
				t.Error(err)
			}
			r, err := SendHTTPRequest(nil, http.MethodPost, "http://localhost:12332/user", buf, nil, CodeIs200, JSONExtractor[rresult])
			if err != nil {
				t.Error(err)
			}

			Assert(t, r.Name[:5], "1000-")
		})
	}
}
