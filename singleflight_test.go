package do

import (
	"sync"
	"testing"
	"time"
)

func TestSingleFlight(t *testing.T) {
	for i := 0; i < 2; i++ {
		testSingleFlight(t)

		// 执行完一次后将key删掉
		ForgotKey("test")
		ForgotKey("book")
	}
}

func testSingleFlight(t *testing.T) {
	cond := sync.NewCond(new(sync.Mutex))

	tests := NewSlice[int]()
	// test key
	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		r, err := SingleFlight("test", func() (int, error) {
			return 1, nil
		})
		if err != nil {
			t.Error(err)
		}
		// t.Logf("test 1 r is: %v", r)
		tests.Append(r)
	}()
	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		r, err := SingleFlight("test", func() (int, error) {
			return 2, nil
		})
		if err != nil {
			t.Error(err)
		}
		// t.Logf("test 2 r is: %v", r)
		tests.Append(r)
	}()
	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		r, err := SingleFlight("test", func() (int, error) {
			return 3, nil
		})
		if err != nil {
			t.Error(err)
		}
		// t.Logf("test 3 r is: %v", r)
		tests.Append(r)
	}()

	books := NewSlice[string]()
	// book key
	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		r, err := SingleFlight("book", func() (string, error) {
			return "book1", nil
		})
		if err != nil {
			t.Error(err)
		}
		// t.Logf("book 1 r is: %v", r)
		books.Append(r)
	}()
	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		r, err := SingleFlight("book", func() (string, error) {
			return "book2", nil
		})
		if err != nil {
			t.Error(err)
		}
		// t.Logf("book 2 r is: %v", r)
		books.Append(r)
	}()
	go func() {
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		r, err := SingleFlight("book", func() (string, error) {
			return "book3", nil
		})
		if err != nil {
			t.Error(err)
		}
		// t.Logf("book 3 r is: %v", r)
		books.Append(r)
	}()

	time.Sleep(200 * time.Millisecond)
	cond.Broadcast()
	time.Sleep(200 * time.Millisecond)

	if tests.Index(0) != tests.Index(1) || tests.Index(1) != tests.Index(2) {
		t.Errorf("bad case: %+v", tests)
	}

	if books.Index(0) != books.Index(1) || books.Index(1) != books.Index(2) {
		t.Errorf("bad case: %+v", books)
	}
}
