package do

import (
	"reflect"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				s: "abc",
			},
			want: []byte{97, 98, 99},
		},
		{
			name: "2",
			args: args{
				s: "请问",
			},
			want: []byte{232, 175, 183, 233, 151, 174},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				b: []byte{97, 98, 99},
			},
			want: "abc",
		},
		{
			name: "2",
			args: args{
				b: []byte{232, 175, 183, 233, 151, 174},
			},
			want: "请问",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToString(tt.args.b); got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
