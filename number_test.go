package do

import (
	"reflect"
	"testing"
)

func TestSplitInt(t *testing.T) {
	type args struct {
		n uint64
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{
			name: "0",
			args: args{
				n: 0,
			},
			want: []uint64{0},
		},
		{
			name: "1",
			args: args{
				n: 1,
			},
			want: []uint64{1},
		},
		{
			name: "12",
			args: args{
				n: 12,
			},
			want: []uint64{1, 2},
		},
		{
			name: "1234",
			args: args{
				n: 1234,
			},
			want: []uint64{1, 2, 3, 4},
		},
		{
			name: "20060102000001",
			args: args{
				n: 20060102000001,
			},
			want: []uint64{2, 0, 0, 6, 0, 1, 0, 2, 0, 0, 0, 0, 0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitUint(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinInt(t *testing.T) {
	type args struct {
		parts []uint64
	}
	tests := []struct {
		name  string
		args  args
		wantN uint64
	}{
		{
			name: "0",
			args: args{
				parts: []uint64{0},
			},
			wantN: 0,
		},
		{
			name: "1",
			args: args{
				parts: []uint64{1},
			},
			wantN: 1,
		},
		{
			name: "12",
			args: args{
				parts: []uint64{1, 2},
			},
			wantN: 12,
		},
		{
			name: "1234",
			args: args{
				parts: []uint64{1, 2, 3, 4},
			},
			wantN: 1234,
		},
		{
			name: "120034",
			args: args{
				parts: []uint64{1, 2, 0, 0, 3, 4},
			},
			wantN: 120034,
		},
		{
			name: "20060102000001",
			args: args{
				parts: []uint64{2, 0, 0, 6, 0, 1, 0, 2, 0, 0, 0, 0, 0, 1},
			},
			wantN: 20060102000001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := JoinUint(tt.args.parts); gotN != tt.wantN {
				t.Errorf("JoinInt() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestPow10(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name  string
		args  args
		wantR uint64
	}{
		{
			name: "-1",
			args: args{
				n: -1,
			},
			wantR: 0,
		},
		{
			name: "0",
			args: args{
				n: 0,
			},
			wantR: 1,
		},
		{
			name: "1",
			args: args{
				n: 1,
			},
			wantR: 10,
		},
		{
			name: "2",
			args: args{
				n: 2,
			},
			wantR: 100,
		},
		{
			name: "10",
			args: args{
				n: 10,
			},
			wantR: 10000000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := Pow10(tt.args.n); gotR != tt.wantR {
				t.Errorf("Pow10() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
