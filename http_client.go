package do

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type (
	HTTPClientOption struct {
		timeout    time.Duration
		skipVerify bool
		transport  *http.Transport
	}
	HTTPClientSetter func(*HTTPClientOption)
)

func HTTPClientWithTimeout(timeout time.Duration) HTTPClientSetter {
	return func(ho *HTTPClientOption) {
		ho.timeout = timeout
	}
}

func HTTPClientSkipVerify() HTTPClientSetter {
	return func(ho *HTTPClientOption) {
		ho.skipVerify = true
	}
}

func HTTPClientWithTransport(tp *http.Transport) HTTPClientSetter {
	return func(ho *HTTPClientOption) {
		ho.transport = tp
	}
}

func NewHTTPClient(opts ...HTTPClientSetter) *http.Client {
	opt := &HTTPClientOption{}
	for _, set := range opts {
		set(opt)
	}

	defaultTransportDialContext := func(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
		return dialer.DialContext
	}

	ts := http.DefaultTransport
	if opt.transport != nil {
		ts = opt.transport
	} else if opt.skipVerify {
		ts = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: defaultTransportDialContext(&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}),
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	to := 10 * time.Second
	if opt.timeout != 0 {
		to = opt.timeout
	}
	client := &http.Client{
		Timeout:   to,
		Transport: ts,
	}

	return client
}
