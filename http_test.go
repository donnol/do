package do

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type (
	result struct {
		BgQuality int `json:"BgQuality"`
	}

	resultWithHeader struct {
		BgQuality int `json:"BgQuality"`

		header http.Header
	}
)

func (r *resultWithHeader) Extract(h http.Header) {
	// fmt.Printf("h: %+v\n", h)
	r.header = h
}

func TestDoHTTPRequest(t *testing.T) {
	type args[R any] struct {
		client        *http.Client
		method        string
		link          string
		body          io.Reader
		header        http.Header
		codeChecker   CodeChecker
		extractResult ResultExtractor[result]
	}
	type testCase[R any] struct {
		name    string
		args    args[R]
		want    result
		wantErr bool
	}
	tests := []testCase[result]{
		// TODO: Add test cases.
		{
			name: "1",
			args: args[result]{
				method:        http.MethodGet,
				link:          "https://www.bing.com/hp/api/model",
				body:          nil,
				codeChecker:   CodeIs200,
				extractResult: JSONExtractor[result],
			},
			want: result{BgQuality: 50},
		},
		{
			name: "nil",
			args: args[result]{
				method:        http.MethodGet,
				link:          "https://www.bing.com/hp/api/model",
				body:          nil,
				codeChecker:   nil,
				extractResult: nil,
			},
			want: result{BgQuality: 50},
		},
		{
			name: "header",
			args: args[result]{
				method: http.MethodGet,
				link:   "https://www.bing.com/hp/api/model",
				body:   nil,
				header: http.Header{
					"User-Agent": []string{`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36`},
				},
				codeChecker: func(code int) error {
					if code != http.StatusOK {
						return fmt.Errorf("code is %d", code)
					}
					return nil
				},
				extractResult: JSONExtractor[result],
			},
			want: result{BgQuality: 50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SendHTTPRequest(tt.args.client, tt.args.method, tt.args.link, tt.args.body, tt.args.header, tt.args.codeChecker, tt.args.extractResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoHTTPRequest() = %v, want %v", got, tt.want)
			}

			{
				got, err := SendHTTPRequest(NewHTTPClient(HTTPClientWithTimeout(2*time.Second), HTTPClientSkipVerify()), tt.args.method, tt.args.link, tt.args.body, tt.args.header, tt.args.codeChecker, tt.args.extractResult)
				if (err != nil) != tt.wantErr {
					t.Errorf("DoHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("DoHTTPRequest() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSendHTTPRequest_Raw(t *testing.T) {
	r, err := SendHTTPRequest(nil, http.MethodGet, "https://www.bing.com/hp/api/model", nil, nil, CodeIs200, RawExtractor)
	if err != nil {
		t.Error(err)
	}
	if len(r) == 0 {
		t.Errorf("bad case, bytes is empty")
	}
}

func TestDoHTTPRequest_RespHeader(t *testing.T) {
	type args[R any] struct {
		client        *http.Client
		method        string
		link          string
		body          io.Reader
		header        http.Header
		codeChecker   CodeChecker
		extractResult ResultExtractor[*resultWithHeader]
	}
	type testCase[R any] struct {
		name    string
		args    args[R]
		want    resultWithHeader
		wantErr bool
	}
	tests := []testCase[result]{
		// TODO: Add test cases.
		{
			name: "1",
			args: args[result]{
				method:        http.MethodGet,
				link:          "https://www.bing.com/hp/api/model",
				body:          nil,
				codeChecker:   CodeIs200,
				extractResult: JSONExtractor[*resultWithHeader],
			},
			want: resultWithHeader{BgQuality: 50},
		},
		{
			name: "nil",
			args: args[result]{
				method:        http.MethodGet,
				link:          "https://www.bing.com/hp/api/model",
				body:          nil,
				codeChecker:   nil,
				extractResult: nil,
			},
			want: resultWithHeader{BgQuality: 50},
		},
		{
			name: "header",
			args: args[result]{
				method: http.MethodGet,
				link:   "https://www.bing.com/hp/api/model",
				body:   nil,
				header: http.Header{
					"User-Agent": []string{`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36`},
				},
				codeChecker: func(code int) error {
					if code != http.StatusOK {
						return fmt.Errorf("code is %d", code)
					}
					return nil
				},
				extractResult: JSONExtractor[*resultWithHeader],
			},
			want: resultWithHeader{BgQuality: 50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SendHTTPRequest(tt.args.client, tt.args.method, tt.args.link, tt.args.body, tt.args.header, tt.args.codeChecker, tt.args.extractResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.BgQuality, tt.want.BgQuality) {
				t.Errorf("DoHTTPRequest() = %v, want %v", got, tt.want)
			}
			if len(got.header) == 0 {
				t.Errorf("DoHTTPRequest() header is empty")
			}
		})
	}
}

type resultCheck struct {
	BgQuality int `json:"BgQuality"`
}

func (r resultCheck) Check() error {
	if r.BgQuality == 50 {
		return fmt.Errorf("bad BgQuality")
	}
	return nil
}

func TestDoHTTPRequest_CheckResult(t *testing.T) {
	type args[R any] struct {
		client        *http.Client
		method        string
		link          string
		body          io.Reader
		header        http.Header
		codeChecker   CodeChecker
		extractResult ResultExtractor[resultCheck]
	}
	type testCase[R any] struct {
		name    string
		args    args[R]
		wantErr bool
	}
	tests := []testCase[result]{
		{
			name: "1",
			args: args[result]{
				method:        http.MethodGet,
				link:          "https://www.bing.com/hp/api/model",
				body:          nil,
				codeChecker:   CodeIs200,
				extractResult: JSONExtractor[resultCheck],
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SendHTTPRequest(tt.args.client, tt.args.method, tt.args.link, tt.args.body, tt.args.header, tt.args.codeChecker, tt.args.extractResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			Assert(t, err.Error(), "bad BgQuality")
		})
	}
}
