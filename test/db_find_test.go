package do

import (
	"database/sql"
	"testing"

	"github.com/donnol/do"
)

// === 对应表结构 ===

type UserForDB struct {
	Id   uint64
	Name string
}

// === 查询 ===

type finderOfUser struct {
	id uint64
}

func (f *finderOfUser) Query() (query string, args []any) {
	query = `select * from user where id = ?`
	args = append(args, f.id)
	return
}

func (f *finderOfUser) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *UserForDB, fields []any) {
	r = &UserForDB{}

	// fields如果数量很多，写起来也很麻烦
	// 怎么根据column type得到fields呢？
	// fields = append(fields,
	// 	&r.Id,
	// 	&r.Name,
	// )

	fields = do.FieldsByColumnType(r, colTypes, nil)

	return
}

// 这样可以省去一个initial参数，但是不能推断出R的实际类型，就需要手动传入类型参数，这样就需要写多个代码，因此暂时不使用这个函数签名
func FindAll1[S do.Storer, F do.Finder[R], R any](db S, finder F) (r []R, err error) {
	return
}

func TestFindAll(t *testing.T) {
	// var finder1 Finder[*UserForDB] = &finderOfUser{}
	// FindAll1(tdb, finder1) // Error: cannot infer R (/home/jd/Project/jd/tools/db/find.go:37:38)

	finder := &finderOfUser{
		id: 1,
	}
	r, err := do.FindAll(tdb, finder, (UserForDB{}))
	if err != nil {
		t.Error(err)
	}
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
}

type finderOfUserOnlyId struct {
	id uint64
}

func (f *finderOfUserOnlyId) Query() (query string, args []any) {
	query = `select id from user where id = ?`
	args = append(args, f.id)
	return
}

func (f *finderOfUserOnlyId) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *UserForDB, fields []any) {
	r = &UserForDB{}
	fields = append(fields,
		&r.Id,
	)
	return
}

func TestFindAll_Id(t *testing.T) {
	finder := &finderOfUserOnlyId{
		id: 1,
	}
	r, err := do.FindAll(tdb, finder, (UserForDB{}))
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("bad result len, len(r) %v != %v", len(r), 1)
	}
	for _, r := range r {
		if r.Id != 1 {
			t.Errorf("bad case of id, %v != %v", r.Id, 1)
		}
		if r.Name != "" {
			t.Errorf("bad case of name, %v != %v", r.Name, "")
		}
	}
}

type finderOfUserOnlyTotal struct {
	id uint64
}

func (f *finderOfUserOnlyTotal) Query() (query string, args []any) {
	query = `select count(id) as count from user where id = ?`
	args = append(args, f.id)
	return
}

func (f *finderOfUserOnlyTotal) NewScanObjAndFields(colTypes []*sql.ColumnType) (r *int64, fields []any) {
	var total int64
	r = &total
	fields = append(fields,
		r,
	)
	return
}

func TestFindAll_Total(t *testing.T) {
	finder := &finderOfUserOnlyTotal{
		id: 1,
	}
	r, err := do.FindAll(tdb, finder, (int64)(0))
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("bad result len1, len(r) %v != %v", len(r), 1)
	}
	for _, r := range r {
		if r != 1 {
			t.Errorf("bad case of id, %v != %v", r, 1)
		}
	}
}

func TestFindList(t *testing.T) {
	finder := &finderOfUser{
		id: 1,
	}
	var r []UserForDB
	err := do.FindList(tdb, finder, &r)
	if err != nil {
		t.Error(err)
	}
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
}

func TestFindFirst(t *testing.T) {
	finder := &finderOfUserOnlyTotal{
		id: 1,
	}
	var total int64
	err := do.FindFirst(tdb, finder, &total)
	if err != nil {
		t.Error(err)
	}
	if total != 1 {
		t.Errorf("bad result, total %v != %v", total, 1)
	}
}

func TestFindListByFunc(t *testing.T) {
	var r []UserForDB
	err := do.FindListByFunc(tdb, findUserById(1), &r)
	if err != nil {
		t.Error(err)
	}
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

	// FindFunc实现了Finder接口
	{
		var r []UserForDB
		err := do.FindList(tdb, findUserById(1), &r)
		if err != nil {
			t.Error(err)
		}
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
	}
}

func findUserById(id uint64) do.FindFunc[UserForDB] {
	return func() (query string, args []any, genObj func(colTypes []*sql.ColumnType) (r *UserForDB, fields []any)) {
		query = `select * from user where id = ?`
		args = append(args, id)

		genObj = func(colTypes []*sql.ColumnType) (r *UserForDB, fields []any) {
			r = &UserForDB{}
			fields = do.FieldsByColumnType(r, colTypes, nil)
			return
		}
		return
	}
}
