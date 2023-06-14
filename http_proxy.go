package do

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HTTPProxyOption struct {
	Director       func(req *http.Request)
	ModifyResponse func(r *http.Response) error
	ErrorHandler   func(w http.ResponseWriter, r *http.Request, err error)
}

// HTTPProxy listen localAddr and transfer any request to remoteAddr. We can use handlers to specify one custom func to transfer data.
func HTTPProxy(localAddr, remoteAddr string, opt *HTTPProxyOption) (err error) {
	url, err := url.Parse(remoteAddr)
	if err != nil {
		return err
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	if opt != nil {
		if opt.Director != nil {
			rp.Director = opt.Director
		}
		if opt.ModifyResponse != nil {
			rp.ModifyResponse = opt.ModifyResponse
		}
		if opt.ErrorHandler != nil {
			rp.ErrorHandler = opt.ErrorHandler
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rp.ServeHTTP(w, r)
	})

	s := &http.Server{
		Addr:    localAddr,
		Handler: mux,
	}
	return s.ListenAndServe()
}
