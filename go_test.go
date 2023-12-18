package do

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"strings"
	"sync"
	"testing"
)

var output = new(bytes.Buffer)
var bw *bufio.Writer

func init() {
	bw = bufio.NewWriter(output)
	log.SetOutput(bw)
}

// Go run f in a new goroutine with defer recover
func TestGo(t *testing.T) {
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

	Must(bw.Flush())
	Assert(t, strings.Contains(output.String(), "panic stack:"), true, "output is %s", output)
}

func TestGoR(t *testing.T) {
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

	r := <-ch
	Assert(t, r, "r1")
	Assert(t, strings.Contains(output.String(), "panic stack:"), true, "output is %s", output)
}
