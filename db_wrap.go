package do

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func WrapDB(
	ctx context.Context,
	driverName string,
	dataSourceName string,
	f func(
		ctx context.Context,
		conn *sql.DB,
	) error,
) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	if err := f(ctx, db); err != nil {
		return err
	}

	return nil
}

func WrapTx(
	ctx context.Context,
	db *sql.DB,
	f func(
		ctx context.Context,
		tx *sql.Tx,
	) error,
) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	var success bool
	defer func() {
		if !success {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = f(ctx, tx)
	if err != nil {
		return err
	}
	success = true

	return
}

func WrapTxV[R any](
	ctx context.Context,
	db *sql.DB,
	f func(
		ctx context.Context,
		tx *sql.Tx,
	) (R, error),
) (r R, err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	var success bool
	defer func() {
		if !success {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	r, err = f(ctx, tx)
	if err != nil {
		return
	}
	success = true

	return
}

func WrapSQLConn(
	ctx context.Context,
	db *sql.DB,
	f func(
		ctx context.Context,
		conn *sql.Conn,
	) error,
) error {

	// 获取连接，并在返回前释放
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close conn failed: %v", err)
		}
	}()

	if err := f(ctx, conn); err != nil {
		funcName := FuncName(2, true)
		return fmt.Errorf("[DB] call in %s failed: %w", funcName, err)
	}

	return nil
}

func WrapSQLConnV[R any](
	ctx context.Context,
	db *sql.DB,
	f func(
		ctx context.Context,
		conn *sql.Conn,
	) (R, error),
) (r R, err error) {

	// 获取连接，并在返回前释放
	conn, err := db.Conn(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close conn failed: %v", err)
		}
	}()

	r, err = f(ctx, conn)
	if err != nil {
		funcName := FuncName(2, true)
		return r, fmt.Errorf("[DB] call in %s failed: %w", funcName, err)
	}

	return
}

// WrapSQLQueryRows query by stmt and args, return values with dest
// only support one row
func WrapSQLQueryRows(
	ctx context.Context,
	db *sql.DB,
	stmt string,
	args []interface{},
	dest ...interface{},
) error {

	if err := WrapSQLConn(ctx, db, func(ctx context.Context, conn *sql.Conn) error {

		rows, err := conn.QueryContext(ctx, stmt, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(dest...)
			if err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// WrapConnFindAll query by stmt and args, return values with dest
// support many rows
func WrapConnFindAll[F Finder[R], R any](
	ctx context.Context,
	db *sql.DB,
	finder F,
	initial R,
) (r []R, err error) {

	if err = WrapSQLConn(ctx, db, func(ctx context.Context, conn *sql.Conn) error {
		err = FindList(conn, finder, &r)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return
	}

	return
}

func WrapTxFindAll[F Finder[R], R any](
	ctx context.Context,
	db *sql.DB,
	finder F,
	initial R,
) (r []R, err error) {

	if err = WrapTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		err = FindList(tx, finder, &r)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	return
}
