package do

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	w := NewWorker(10)
	w.Start()

	if err := w.Push(Job{
		doctx: func(ctx context.Context) error {
			panic("terriable")
		},
	}); err != nil {
		t.Error(err)
	}
	if err := w.Push(Job{
		doctx: func(ctx context.Context) error {
			for i := 0; i < 10; i++ {
				_ = i
			}
			return nil
		},
	}); err != nil {
		t.Error(err)
	}
	for i := 10; i < 1000; i++ {
		tmp := i
		if err := w.Push(Job{
			doctx: func(ctx context.Context) error {
				_ = tmp
				return nil
			},
		}); err != nil {
			t.Error(err)
		}
	}
	if err := w.Push(Job{}); err != ErrNilJobDo {
		t.Error(err)
	}

	w.Stop()

	if err := w.Push(Job{
		doctx: func(ctx context.Context) error {
			log.Printf("Push after stop.")
			return nil
		},
	}); err != ErrWorkerIsStop {
		t.Error(err)
	}
}

func TestWorkerBuffer(t *testing.T) {
	w := NewWorker(10)
	w.Start()

	for i := 1; i <= 1000; i++ {
		tmp := i
		if err := w.Push(Job{
			doctx: func(ctx context.Context) error {
				_ = tmp
				time.Sleep(50 * time.Millisecond)

				return nil
			},
		}); err != nil {
			t.Error(err)
		}

		// cause DATA RACE
		// bl := w.jobChan.BufLen()
		// if bl != 0 {
		// 	_ = bl
		// }
	}
	t.Log("wait...")

	w.Stop()

	t.Log("finish.")
}

func TestWorkerWithTimeout(t *testing.T) {
	w := NewWorker(0)
	w.Start()

	job := NewJob(func(ctx context.Context) error {
		// 使用for select监听ctx.Done，在业务逻辑执行间隔检查是否超时
		for {
			select {
			case <-ctx.Done():
				return fmt.Errorf("timeout exceed")
			default:
				// 模拟业务逻辑 -- 为了留出间隔检查超时，如果是大任务需要分批执行；
				log.Printf("runing...")
				time.Sleep(500 * time.Millisecond)

				// 但是，如果执行的任务拆分不了呢，它就是一个长时间的任务呢？
				// 那就没机会检查ctx.Done()了，也就没机会退出执行。
				// 那也只能让它一直执行下去了，直到它结束了。
			}
		}
	}, 1*time.Second, func(err error) {
		if err.Error() != "timeout exceed" {
			t.Errorf("not exist timeout error")
		}
	})
	if err := w.Push(*job); err != nil {
		t.Error(err)
	}

	w.Stop()
}
