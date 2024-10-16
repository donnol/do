package do

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/donnol/do"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	createSQL = `
	create table if not exists user (
		id integer not null,
		name varchar(255) not null,
		primary key(id)
	);
	INSERT OR IGNORE INTO user (id, name) values (1, 'jd');
	`
)

var (
	tdb = func() *sql.DB {
		dir := "../testdata"
		do.Must(do.MkdirAllIfNotExist(dir))
		gdb := do.Must1(gorm.Open(sqlite.Open(filepath.Join(dir, "test.db"))))

		sqldb := do.Must1(gdb.DB())

		do.Must1(sqldb.Exec(createSQL))

		return sqldb
	}()
)

func TestWrapTxV(t *testing.T) {
	ctx := context.Background()
	r, err := do.WrapTxV[UserForDB](ctx, tdb, func(ctx context.Context, tx *sql.Tx) (UserForDB, error) {

		// 只拿id列
		finder := &finderOfUserOnlyId{
			id: 1,
		}
		var r UserForDB
		do.Must(do.FindFirst(tx, finder, &r))
		if r.Id != 1 {
			t.Errorf("bad case of id, %v != %v", r.Id, 1)
		}
		if r.Name != "" {
			t.Errorf("bad case of name, %v != %v", r.Name, "")
		}

		return r, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if r.Id != 1 {
		t.Errorf("bad case of id, %v != %v", r.Id, 1)
	}
	if r.Name != "" {
		t.Errorf("bad case of name, %v != %v", r.Name, "")
	}
}

func TestWrapSQLConnV(t *testing.T) {
	ctx := context.Background()
	r, err := do.WrapSQLConnV[UserForDB](ctx, tdb, func(ctx context.Context, tx *sql.Conn) (UserForDB, error) {

		// 只拿id列
		finder := &finderOfUserOnlyId{
			id: 1,
		}
		var r UserForDB
		do.Must(do.FindFirst(tx, finder, &r))
		if r.Id != 1 {
			t.Errorf("bad case of id, %v != %v", r.Id, 1)
		}
		if r.Name != "" {
			t.Errorf("bad case of name, %v != %v", r.Name, "")
		}

		return r, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", r)

	if r.Id != 1 {
		t.Errorf("bad case of id, %v != %v", r.Id, 1)
	}
	if r.Name != "" {
		t.Errorf("bad case of name, %v != %v", r.Name, "")
	}
}
