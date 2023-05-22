package do

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	CodeChecker func(code int) error
)

func CodeIs200(code int) error {
	if code != http.StatusOK {
		return fmt.Errorf("bad http code: %d", code)
	}
	return nil
}

type (
	ResultExtractor[R any] func(data []byte) (R, error)
)

func JSONExtractor[R any](data []byte) (R, error) {
	var r R
	if err := json.Unmarshal(data, &r); err != nil {
		return r, err
	}

	return r, nil
}

func XMLExtractor[R any](data []byte) (R, error) {
	var r R
	if err := xml.Unmarshal(data, &r); err != nil {
		return r, err
	}

	return r, nil
}

// SendHTTPRequest send http request and get result of type R. If you want to got resp header, the R should implement RespHeaderExtractor interface.
func SendHTTPRequest[R any](
	client *http.Client,
	method string,
	link string,
	body io.Reader,
	header http.Header,
	codeChecker CodeChecker,
	extractResult ResultExtractor[R],
) (R, error) {
	var r R

	if method == "" || link == "" {
		return r, fmt.Errorf("bad param: method or link is empty")
	}

	req, err := http.NewRequest(method, link, body)
	if err != nil {
		return r, err
	}
	for k, v := range header {
		for _, vv := range v {
			req.Header.Set(k, vv)
		}
	}

	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	if codeChecker == nil {
		codeChecker = CodeIs200
	}
	err = codeChecker(resp.StatusCode)
	if err != nil {
		return r, fmt.Errorf("check code failed: %v, data: %s", err, data)
	}

	if extractResult == nil {
		extractResult = JSONExtractor[R]
	}
	r, err = extractResult(data)
	if err != nil {
		return r, fmt.Errorf("extract result failed: %v, data: %s", err, data)
	}

	// with header
	h := resp.Header
	if e, ok := any(r).(RespHeaderExtractor); ok && e != nil {
		e.Extract(h)
	}

	return r, nil
}

type RespHeaderExtractor interface {
	Extract(h http.Header)
}
