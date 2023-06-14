package do

import (
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

const prefix = "http://localhost"

func TestHTTPProxy(t *testing.T) {
	type args struct {
		localAddr  string
		remoteAddr string
		msg        []byte
		opt        *HTTPProxyOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				localAddr:  ":55678",
				remoteAddr: ":55789",
				msg:        []byte("hello" + strconv.Itoa(rand.Int())),
				opt:        nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go startServer(tt.args.remoteAddr, tt.args.msg)

			go func() {
				if err := HTTPProxy(tt.args.localAddr, prefix+tt.args.remoteAddr, tt.args.opt); (err != nil) != tt.wantErr {
					t.Errorf("HTTPProxy() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			time.Sleep(time.Millisecond * 200)

			{
				r, err := SendHTTPRequest(nil, http.MethodGet, prefix+tt.args.localAddr, nil, nil, CodeIs200, RawExtractor)
				if err != nil {
					t.Error(err)
				}
				if string(r) != string(tt.args.msg) {
					t.Errorf("bad case, %s != %s", r, tt.args.msg)
				}
			}
		})
	}
}

func startServer(addr string, ret []byte) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(ret)
	})

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	Must(s.ListenAndServe())
}
