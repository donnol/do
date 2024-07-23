package do

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFieldWithAlias(t *testing.T) {
	type args struct {
		field        string
		alias        string
		defaultAlias string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "empty",
			args: args{
				field:        "name",
				alias:        "",
				defaultAlias: "",
			},
			want: "name",
		},
		{
			name: "default",
			args: args{
				field:        "name",
				alias:        "",
				defaultAlias: "s",
			},
			want: "s.name",
		},
		{
			name: "alias",
			args: args{
				field:        "name",
				alias:        "r",
				defaultAlias: "s",
			},
			want: "r.name",
		},
		{
			name: "alias1",
			args: args{
				field:        "name",
				alias:        "r",
				defaultAlias: "",
			},
			want: "r.name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FieldWithAlias(tt.args.field, tt.args.alias, tt.args.defaultAlias); got != tt.want {
				t.Errorf("FieldWithAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithWhereInt(t *testing.T) {
	type args struct {
		t     int
		cond  func(field string, value interface{}) string
		field string
		value []int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "zero",
			args: args{
				t: 0,
				cond: func(field string, value interface{}) string {
					return fmt.Sprintf(`%s = %v`, field, value)
				},
				field: "name",
				value: []int{},
			},
			want: "",
		},
		{
			name: "not-zero",
			args: args{
				t: 1,
				cond: func(field string, value interface{}) string {
					return fmt.Sprintf(`%s = %v`, field, value)
				},
				field: "name",
				value: []int{},
			},
			want: "name = 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithWhere(tt.args.t, tt.args.cond, tt.args.field, tt.args.value...); got != tt.want {
				t.Errorf("WithWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithWhereTime(t *testing.T) {
	type args struct {
		t     time.Time
		cond  func(field string, value interface{}) string
		field string
		value []time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "zero",
			args: args{
				t: time.Time{},
				cond: func(field string, value interface{}) string {
					return fmt.Sprintf(`%s = %v`, field, value)
				},
				field: "created",
				value: []time.Time{},
			},
			want: "",
		},
		{
			name: "not-zero",
			args: args{
				t: time.Date(2023, 10, 26, 0, 0, 0, 0, time.Local),
				cond: func(field string, value interface{}) string {
					return fmt.Sprintf(`%s = '%v'`, field, value)
				},
				field: "created",
				value: []time.Time{},
			},
			want: "created = '2023-10-26 00:00:00 ",
		},
		{
			name: "not-zero-value",
			args: args{
				t: time.Date(2023, 10, 26, 0, 0, 0, 0, time.Local),
				cond: func(field string, value interface{}) string {
					return fmt.Sprintf(`%s = '%v'`, field, value)
				},
				field: "created",
				value: []time.Time{
					time.Date(2023, 10, 27, 0, 0, 0, 0, time.Local),
				},
			},
			want: "created = '2023-10-27 00:00:00 ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithWhere(tt.args.t, tt.args.cond, tt.args.field, tt.args.value...); !strings.Contains(got, tt.want) {
				t.Errorf("WithWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}
