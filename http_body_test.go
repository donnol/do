package do

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestMultipartBody(t *testing.T) {
	type args struct {
		field string
		name  string
		data  []byte
	}
	tests := []struct {
		name                string
		args                args
		wantFileContentType string
		wantBody            string
		wantErr             bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				field: "file",
				name:  "test.md",
				data: func() []byte {
					return []byte("abc")
				}(),
			},
			wantFileContentType: "multipart/form-data",
			wantBody:            "abc",
			wantErr:             false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			gotFileContentType, err := MultipartBody(body, tt.args.field, tt.args.name, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MultipartBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(gotFileContentType, tt.wantFileContentType) {
				t.Errorf("MultipartBody() = %v, want %v", gotFileContentType, tt.wantFileContentType)
			}
			if gotBody := body.String(); !strings.Contains(gotBody, tt.wantBody) {
				t.Errorf("MultipartBody() = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestSendFile(t *testing.T) {
	addr := ":9089"
	go start(addr)
	time.Sleep(time.Millisecond * 100)

	body := new(bytes.Buffer)
	ct, err := MultipartBody(body, "file", "testFileName", []byte("test data"))
	if err != nil {
		t.Error(err)
	}
	r, err := SendHTTPRequest(nil, http.MethodPost, "http://127.0.0.1"+addr+"/upload", body, http.Header{
		"Content-Type": []string{ct},
	}, CodeIs200, JSONExtractor[any])
	if err != nil {
		t.Error(err)
	}
	_ = r
}

func start(addr string) {
	Must(http.ListenAndServe(addr, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/upload" {
				mr := Must1(r.MultipartReader())
				form := Must1(mr.ReadForm(10 * 1024))
				value := form.Value["file"]
				if value[0] != "testFileName" {
					panic(fmt.Errorf("bad field value, %s != %s", value[0], "testFileName"))
				}
				for k, v := range form.File {
					if k != "file" {
						panic(fmt.Errorf("bad key"))
					}
					for _, fh := range v {
						f := Must1(fh.Open())
						defer f.Close()
						buf := make([]byte, fh.Size)
						Must1(f.Read(buf))
						if string(buf) != "test data" {
							panic(fmt.Errorf("bad file data: %s != %s", buf, "test data"))
						}
					}
				}

				w.Write([]byte(`{"name": "jd"}`))
			}
		},
	)))
}
