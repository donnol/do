package do

import (
	"bufio"
	"bytes"
	"context"
	"log"
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

	Assert(t, strings.Contains(output.String(), "panic stack:"), true, "output is %s", output)
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
	Assert(t, strings.Contains(output.String(), "panic stack:"), true, "output is %s", output)
}
