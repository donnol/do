package do

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	prefix      = "http://localhost"
	httpsPrefix = "https://localhost"
)

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

func TestHTTPSProxy(t *testing.T) {
	t.Skip()

	// 为了避免错误：http server doesn't support hijacking connection
	// From: https://stackoverflow.com/questions/67770829/go-http-request-falls-back-to-http2-even-when-force-attempt-is-set-to-false
	os.Setenv("GODEBUG", "http2client=0")

	type args struct {
		localAddr  string
		remoteAddr string
		msg        []byte
		cert, key  string
		opt        *HTTPSProxyOption
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
				localAddr:  ":55688",
				remoteAddr: ":55799",
				msg:        []byte("hello" + strconv.Itoa(rand.Int())),
				cert:       "./testdata/cert/server.crt",
				key:        "./testdata/cert/server.key",
				opt: &HTTPSProxyOption{
					CertFile:   "./testdata/cert2/server.crt",
					KeyFile:    "./testdata/cert2/server.key",
					CaCertFile: "./testdata/cert3/server.crt",
					CaKeyFile:  "./testdata/cert3/server.key",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go startHTTPSServer(tt.args.remoteAddr, tt.args.msg, tt.args.cert, tt.args.key)

			go func() {
				if err := HTTPSProxy(tt.args.localAddr, tt.args.opt); (err != nil) != tt.wantErr {
					t.Errorf("HTTPProxy() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			time.Sleep(time.Millisecond * 200)

			{
				client := NewHTTPClient(HTTPClientSkipVerify())
				// From https://github.com/golang/go/issues/22554
				// Go's HTTP client does not support sending a request with
				// the CONNECT method. See the documentation on Transport for
				// details.
				r, err := SendHTTPRequest(client, http.MethodConnect, httpsPrefix+tt.args.localAddr+"/", nil, nil, CodeIs200, RawExtractor)
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

func startHTTPSServer(addr string, ret []byte, cert, key string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(ret)
	})

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Println("Starting https server on", addr)
	Must(s.ListenAndServeTLS(cert, key))
}

func TestHTTPSProxyByBing(t *testing.T) {
	t.Skip()

	addr := "localhost:55689"
	go func() {
		if err := HTTPSProxy(addr, &HTTPSProxyOption{
			CertFile:   "./testdata/cert2/server.crt",
			KeyFile:    "./testdata/cert2/server.key",
			CaCertFile: "./testdata/cert3/server.crt",
			CaKeyFile:  "./testdata/cert3/server.key",
		}); err != nil {
			t.Errorf("HTTPProxy() error = %v", err)
		}
	}()
	time.Sleep(time.Millisecond * 200)

	r := Must1(http.NewRequest(http.MethodGet, "https://www.bing.com/hp/api/model", nil))

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse("https://" + addr)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp := Must1(client.Do(r))
	defer resp.Body.Close()

	data := Must1(io.ReadAll(resp.Body))
	var res result
	Must(json.Unmarshal(data, &res))
	Assert(t, res.BgQuality, 50)
}

func TestHTTPSProxyByBaidu(t *testing.T) {
	t.Skip()

	addr := "localhost:55690"
	go func() {
		if err := HTTPSProxy(addr, &HTTPSProxyOption{
			CertFile:   "./testdata/cert2/server.crt",
			KeyFile:    "./testdata/cert2/server.key",
			CaCertFile: "./testdata/cert3/server.crt",
			CaKeyFile:  "./testdata/cert3/server.key",
		}); err != nil {
			t.Errorf("HTTPProxy() error = %v", err)
		}
	}()
	time.Sleep(time.Millisecond * 200)

	r := Must1(http.NewRequest(http.MethodGet, "https://www.baidu.com", nil))

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse("https://" + addr)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp := Must1(client.Do(r))
	defer resp.Body.Close()

	data := Must1(io.ReadAll(resp.Body))
	Assert(t, bytes.Contains(data, []byte("关于")), true, "resp.body is %s", data)
}
