package do

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// Go run f in a new goroutine with defer recover
func TestGo(t *testing.T) {
	var output = new(bytes.Buffer)
	bw := bufio.NewWriter(output)
	log.SetOutput(bw)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	Go(context.Background(), 1, func(ctx context.Context, p int) {
		defer wg.Done()
		log.Println("param:", p)
	})

	wg.Add(1)
	Go(context.Background(), 1, func(ctx context.Context, p int) {
		defer wg.Done()
		log.Println("param:", p)

		panic(p)
	})

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	Must(bw.Flush())

	Assert(t, strings.Contains(output.String(), "panic: 1 \nstack:"), true, "output is %s", output)
}

func TestGoR(t *testing.T) {
	var output = new(bytes.Buffer)
	bw := bufio.NewWriter(output)
	log.SetOutput(bw)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	ch := GoR(context.Background(), 1, func(ctx context.Context, p int) string {
		defer wg.Done()
		log.Println("param:", p)

		return "r1"
	})

	wg.Add(1)
	GoR(context.Background(), 1, func(ctx context.Context, p int) string {
		defer wg.Done()
		log.Println("param:", p)

		panic(p)
	})

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	Must(bw.Flush())

	r := <-ch
	Assert(t, r, "r1")
	Assert(t, strings.Contains(output.String(), "panic: 1 \nstack:"), true, "output is %s", output)
}

func TestCallInDefRec(t *testing.T) {
	type args struct {
		c C
		p int
		f func(C, int) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				c: context.Background(),
				p: 1,
				f: func(ctx context.Context, p int) string {
					return strconv.Itoa(p)
				},
			},
			want: "1",
		},
		{
			name: "panic",
			args: args{
				c: context.Background(),
				p: 1,
				f: func(ctx context.Context, p int) string {
					panic(fmt.Errorf("panic test"))
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallInDefRec(tt.args.c, tt.args.p, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallInDefRec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallInDefRec2(t *testing.T) {
	type args struct {
		c C
		p int
		f func(C, int) (string, E)
	}
	tests := []struct {
		name  string
		args  args
		wantR string
		wantE E
	}{
		{
			name: "1",
			args: args{
				c: context.Background(),
				p: 1,
				f: func(ctx context.Context, p int) (string, error) {
					return strconv.Itoa(p), nil
				},
			},
			wantR: "1",
			wantE: nil,
		},
		{
			name: "1",
			args: args{
				c: context.Background(),
				p: 1,
				f: func(ctx context.Context, p int) (string, error) {
					panic(fmt.Errorf("bad case for panic test"))
				},
			},
			wantR: "",
			wantE: fmt.Errorf("failed: bad case for panic test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotE := CallInDefRec2(tt.args.c, tt.args.p, tt.args.f)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("CallInDefRec2() gotR = %v, want %v", gotR, tt.wantR)
			}
			if tt.wantE != nil && gotE.Error() != tt.wantE.Error() {
				t.Errorf("CallInDefRec2() gotE = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}
