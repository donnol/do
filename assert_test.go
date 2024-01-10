package do

import (
	"bytes"
	"fmt"
	"testing"
)

type myHandler struct {
	buf *bytes.Buffer
}

func (h *myHandler) Errorf(format string, args ...any) {
	h.buf.WriteString(fmt.Sprintf(format, args...))
}

func TestAssert(t *testing.T) {
	h := &myHandler{
		buf: new(bytes.Buffer),
	}

	type args struct {
		logger     AssertHandler
		l          int
		r          int
		msgAndArgs []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantBuf string
	}{
		{
			name: "equal",
			args: args{
				logger:     h,
				l:          1,
				r:          1,
				msgAndArgs: []any{},
			},
			wantErr: false,
		},
		{
			name: "not equal",
			args: args{
				logger: h,
				l:      0,
				r:      1,
				msgAndArgs: []any{
					"msg %s",
					"need help",
				},
			},
			wantErr: true,
			wantBuf: "[/home/jd/Project/jd/do/assert_test.go:61] Bad case, 0 != 1, msg need help",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Assert(tt.args.logger, tt.args.l, tt.args.r, tt.args.msgAndArgs...)
			if tt.wantErr && h.buf.Len() == 0 {
				t.Errorf("bad case, buf is empty")
			} else if tt.wantErr && h.buf.String() != tt.wantBuf {
				t.Errorf("bad case, buf is %s", h.buf)
			}
		})
	}
}

func TestAssertSlice(t *testing.T) {
	h := &myHandler{
		buf: new(bytes.Buffer),
	}
	type args struct {
		logger     AssertHandler
		l          []int
		r          []int
		msgAndArgs []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantBuf string
	}{
		{
			name: "equal",
			args: args{
				logger:     h,
				l:          []int{1},
				r:          []int{1},
				msgAndArgs: []any{},
			},
			wantErr: false,
		},
		{
			name: "not equal",
			args: args{
				logger: h,
				l:      []int{1},
				r:      []int{2},
				msgAndArgs: []any{
					"msg %s",
					"need help",
				},
			},
			wantErr: true,
			wantBuf: "Bad case, No.0: 1 != 2, msg need help",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertSlice(tt.args.logger, tt.args.l, tt.args.r, tt.args.msgAndArgs...)
			if tt.wantErr && h.buf.Len() == 0 {
				t.Errorf("bad case, buf is empty")
			} else if tt.wantErr && h.buf.String() != tt.wantBuf {
				t.Errorf("bad case, buf is %s", h.buf)
			}
		})
	}
}
