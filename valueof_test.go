package do

import (
	"reflect"
	"testing"
)

func TestValueOf(t *testing.T) {
	var v = struct {
		Name string
	}{
		Name: "jd",
	}
	var m = map[string]string{
		"name": "jc",
	}

	type args struct {
		v     any
		field string
	}
	tests := []struct {
		name   string
		args   args
		wantF  string
		wantOk bool
	}{
		{
			name: "v",
			args: args{
				v:     v,
				field: "Name",
			},
			wantF:  "jd",
			wantOk: true,
		},
		{
			name: "&v",
			args: args{
				v:     &v,
				field: "Name",
			},
			wantF:  "jd",
			wantOk: true,
		},
		{
			name: "v-nofield",
			args: args{
				v:     v,
				field: "name",
			},
			wantF:  "",
			wantOk: false,
		},
		{
			name: "m",
			args: args{
				v:     m,
				field: "name",
			},
			wantF:  "jc",
			wantOk: true,
		},
		{
			name: "&m",
			args: args{
				v:     &m,
				field: "name",
			},
			wantF:  "jc",
			wantOk: true,
		},
		{
			name: "m-nofield",
			args: args{
				v:     m,
				field: "Name",
			},
			wantF:  "",
			wantOk: false,
		},
		{
			name: "notstruct&notmap",
			args: args{
				v:     1,
				field: "Name",
			},
			wantF:  "",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotF, gotOk := ValueOf[string](tt.args.v, tt.args.field)
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("ValueOf() gotF = %v, want %v", gotF, tt.wantF)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ValueOf() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
