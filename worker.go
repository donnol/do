package do

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/smallnest/chanx"
)

const (
	errCount = 100

	defaultCount = 50
)

var logger = log.New(os.Stdout, "[Worker]", log.LstdFlags|log.Lshortfile)

var (
	ErrWorkerIsStop = errors.New("Worker is stop")
	ErrNilJobDo     = errors.New("Job do field is nil")
)

type Job struct {
	doctx DoWithCtx

	timeout time.Duration

	errorHandler ErrorHandler
}

type DoWithCtx func(ctx context.Context) error

type ErrorHandler func(error)

func NewJob(do DoWithCtx, timeout time.Duration, eh ErrorHandler) *Job {
	return &Job{
		doctx:        do,
		timeout:      timeout,
		errorHandler: eh,
	}
}

func (job *Job) run(ctx context.Context) error {
	if job.doctx == nil {
		return ErrNilJobDo
	}

	if err := job.doctx(ctx); err != nil {
		if job.errorHandler != nil {
			job.errorHandler(err)
		} else {
			return err
		}
	}

	return nil
}

type Worker struct {
	// all chan must have make, read, write, close operate
	limitChan chan struct{} // control the goroutine number
	stopChan  chan struct{}
	jobChan   *chanx.UnboundedChan[Job]
	errChan   chan error

	wg   *sync.WaitGroup
	stop bool
}

// NewWorker new a worker with limit number
func NewWorker(n int) *Worker {
	if n <= 0 {
		n = defaultCount
	}
	return &Worker{
		limitChan: make(chan struct{}, n),
		stopChan:  make(chan struct{}),
		jobChan:   chanx.NewUnboundedChan[Job](context.Background(), n),
		errChan:   make(chan error, errCount),
		wg:        new(sync.WaitGroup),
	}
}

func (w *Worker) Start() {
	go w.handleError()
	go w.start()

	logger.Printf("Start.\n")
}

func (w *Worker) start() {
	for {
		select {
		case job, ok := <-w.jobChan.Out: // 有工作
			if !ok {
				continue
			}

			w.do(job)

		case <-w.stopChan:
			w.close()
			return
		}
	}
}

func (w *Worker) do(job Job) {
	// 占据一个坑
	w.limitChan <- struct{}{}

	// 开始工作
	go func(job Job) {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("job: %+v\n", r)
			}

			// 释放一个坑
			<-w.limitChan
			w.wg.Done()
		}()

		// 执行
		ctx := context.Background()
		if job.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, job.timeout)
			defer cancel()
		}

		if err := job.run(ctx); err != nil {
			w.errChan <- err
		}
	}(job)
}

func (w *Worker) handleError() {
	for err := range w.errChan {
		logger.Printf("err is %v\n", err)
	}
}

func (w *Worker) Stop() {
	w.stop = true

	w.wait()

	w.stopChan <- struct{}{}

	logger.Printf("Stop.\n")
}

func (w *Worker) close() {
	// close管道时，有可能panic
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("close: %+v\n", r)
		}
	}()

	close(w.stopChan)
	close(w.errChan)
	close(w.jobChan.In)
	close(w.limitChan)
}

func (w *Worker) wait() {
	w.wg.Wait()
}

func (w *Worker) Push(job Job) error {
	if w.stop {
		return ErrWorkerIsStop
	}
	if job.doctx == nil {
		return ErrNilJobDo
	}

	w.jobChan.In <- job
	w.wg.Add(1)

	return nil
}
