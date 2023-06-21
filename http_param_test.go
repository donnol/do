package do

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func decodeValues(srv url.Values, v any) error {
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	rv = rv.Elem()
	rt = rt.Elem()
	for i := 0; i < rt.NumField(); i++ {
		fname := rt.Field(i).Tag.Get("form")

		fv := srv[fname][0]

		vf := rv.Field(i)
		vf.SetString(fv)
	}
	return nil
}

func decodeJSONBytes(srv []byte, v any) error {
	return json.Unmarshal(srv, v)
}

func decodeXMLBytes(srv []byte, v any) error {
	return xml.Unmarshal(srv, v)
}

var (
	valuesParser    = NewParamParser[url.Values](DecodeFunc[url.Values](decodeValues))
	bytesJSONParser = NewParamParser[[]byte](DecodeFunc[[]byte](decodeJSONBytes))
	bytesXMLParser  = NewParamParser[[]byte](DecodeFunc[[]byte](decodeXMLBytes))
	_               = bytesXMLParser
)

func ParseParam[T ParamData](data T, v any) error {
	var err error
	switch vv := any(data).(type) {
	case []byte:
		err = bytesJSONParser.ParseAndCheck(context.Background(), vv, v)
	case url.Values:
		err = valuesParser.ParseAndCheck(context.Background(), vv, v)
	}
	return err
}

type param struct {
	Name string `json:"name" form:"name"`
}

func (p param) Check() error {
	if p.Name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}

func TestParamParse(t *testing.T) {
	type args struct {
		t    string
		data any
		r    any
	}

	for _, tt := range []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "bytes",
			args: args{
				t:    "bytes",
				data: []byte(`{"name": "jd"}`),
				r:    &param{},
			},
			want:    &param{Name: "jd"},
			wantErr: false,
		},
		{
			name: "values",
			args: args{
				t:    "values",
				data: url.Values{"name": []string{"jd"}},
				r:    &param{},
			},
			want:    &param{Name: "jd"},
			wantErr: false,
		},
		{
			name: "values err",
			args: args{
				t:    "values",
				data: url.Values{"name": []string{""}},
				r:    &param{},
			},
			want:    &param{Name: ""},
			wantErr: true,
		},
	} {
		switch tt.args.t {
		case "bytes":
			if err := ParseParam(tt.args.data.([]byte), tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("got err: %v", err)
			} else if !reflect.DeepEqual(tt.args.r, tt.want) {
				t.Errorf("bad case: %v != %v", tt.args.r, tt.want)
			}
		case "values":
			if err := ParseParam(tt.args.data.(url.Values), tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("got err: %v", err)
			} else if !reflect.DeepEqual(tt.args.r, tt.want) {
				t.Errorf("bad case: %v != %v", tt.args.r, tt.want)
			}
		}
	}
}
