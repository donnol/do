package do

import (
	"database/sql"
	"path/filepath"

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
