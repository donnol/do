package do

import (
	"testing"

	"github.com/donnol/do"
)

func TestBatch(t *testing.T) {
	finder := &finderOfUser{
		id: 1,
	}
	err := do.Batch(tdb, finder, 10, func(r []UserForDB) error {
		if len(r) != 1 {
			t.Errorf("bad result len, len(r) %v != %v", len(r), 1)
		}
		for _, r := range r {
			if r.Id != 1 {
				t.Errorf("bad case of id, %v != %v", r.Id, 1)
			}
			if r.Name != "jd" {
				t.Errorf("bad case of name, %v != %v", r.Name, "jd")
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestBatchConcurrent(t *testing.T) {
	finder := &finderOfUser{
		id: 1,
	}
	err := do.BatchConcurrent(tdb, finder, 10, func(r []UserForDB) error {
		if len(r) != 1 {
			t.Errorf("bad result len, len(r) %v != %v", len(r), 1)
		}
		for _, r := range r {
			if r.Id != 1 {
				t.Errorf("bad case of id, %v != %v", r.Id, 1)
			}
			if r.Name != "jd" {
				t.Errorf("bad case of name, %v != %v", r.Name, "jd")
			}
		}
		return nil
	}, 2)
	if err != nil {
		t.Error(err)
	}
}
