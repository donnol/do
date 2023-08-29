package do

import (
	"database/sql"
	"testing"

	"github.com/donnol/do"
	"github.com/jmoiron/sqlx"
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

type batchFinderOfUser struct {
	ids []uint64
}

func (f *batchFinderOfUser) Query() (query string, args []any) {
	query = `select * from user where id IN (?)`
	args = append(args, f.ids)
	return
}

func (f *batchFinderOfUser) Batch() [][]any {
	r := make([][]any, 0, 3)
	for _, id := range f.ids {
		id := id
		r = append(r, []any{id})
	}
	return r
}

func (f *batchFinderOfUser) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *UserForDB, fields []any) {
	r = &UserForDB{}

	fields = do.FieldsByColumnType(r, colTypes, nil)

	return
}

func TestFindWithBatch(t *testing.T) {
	finder := &batchFinderOfUser{
		ids: []uint64{1, 2, 3},
	}
	u := []UserForDB{}
	if err := do.FindWithBatch(tdb, finder, &u); err != nil {
		q, a := finder.Query()
		t.Error(err, q, a, finder.Batch())
		t.Log(sqlx.In(q, a))
	}
	if len(u) != 1 && u[0].Id != 1 {
		t.Errorf("bad result: %+v", u)
	}
}
