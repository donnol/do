package do

import (
	"context"
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
				time.Sleep(500 * time.Millisecond)

				return nil
			},
		}); err != nil {
			t.Error(err)
		}

		bl := w.jobChan.BufLen()
		if bl != 0 {
			_ = bl
		}
	}
	t.Log("finish")

	w.Stop()
}

func TestWorkerWithTimeout(t *testing.T) {
	w := NewWorker(0)
	w.Start()

	job := NewJob(func(ctx context.Context) error {
		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
		}
		return nil
	}, 5*time.Second, func(err error) {
		log.Printf("err is %v\n", err)
	})
	if err := w.Push(*job); err != nil {
		t.Error(err)
	}

	w.Stop()
}
