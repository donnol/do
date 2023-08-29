package do

import (
	"context"
	"fmt"
)

// Batch process data find from Storer in batches
func Batch[S Storer, F Finder[R], R any](db S, finder F, batchNum int, handler func([]R) error) (err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	batch := make([]R, 0, batchNum)
	for rows.Next() {
		t, fields := finder.NewScanObjAndFields(colTypes)
		if err = rows.Scan(fields...); err != nil {
			return
		}

		batch = append(batch, *t)
		if batchNum > 0 && len(batch) >= batchNum {
			if err = handler(batch); err != nil {
				err = fmt.Errorf("batch handle failed %w", err)
				return
			}

			batch = make([]R, 0, batchNum)
		}
	}
	if err = rows.Err(); err != nil {
		return
	}

	if len(batch) != 0 {
		if err = handler(batch); err != nil {
			err = fmt.Errorf("batch handle failed %w", err)
			return
		}
	}

	return
}

// BatchConcurrent batch process data concurrently
func BatchConcurrent[S Storer, F Finder[R], R any](db S, finder F, batchNum int, handler func([]R) error, concNum int) (err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}

	batchWorker := NewWorker(concNum)
	batchWorker.Start()
	defer batchWorker.Stop()

	batch := make([]R, 0, batchNum)
	for rows.Next() {
		t, fields := finder.NewScanObjAndFields(colTypes)
		if err = rows.Scan(fields...); err != nil {
			return
		}

		batch = append(batch, *t)
		if batchNum > 0 && len(batch) >= batchNum {
			batchCopy := batch
			if err = batchWorker.Push(*NewJob(
				DoWithCtx(func(ctx context.Context) error {
					if err = handler(batchCopy); err != nil {
						err = fmt.Errorf("batch handle failed %w", err)
						return err
					}
					return nil
				}),
				0,
				nil,
			)); err != nil {
				return
			}

			batch = make([]R, 0, batchNum)
		}
	}
	if err = rows.Err(); err != nil {
		return
	}

	if len(batch) != 0 {
		batchCopy := batch
		if err = batchWorker.Push(*NewJob(
			DoWithCtx(func(ctx context.Context) error {
				if err = handler(batchCopy); err != nil {
					err = fmt.Errorf("batch handle failed %w", err)
					return err
				}
				return nil
			}),
			0,
			nil,
		)); err != nil {
			return
		}
	}

	return
}

type Batcher interface {
	Batch() [][]any
}

type BatchQueryer interface {
	Queryer
	Batcher
}

type BatchFinder[T any] interface {
	Finder[T]
	Batcher
}

// FindWithBatch use batchFunc to split args to little batch, for example: args is 1, [1, 2, 3], split to 3 batch is: 1, [1]; 1, [2]; 3, [3], the slice become little while the others is not change
func FindWithBatch[S Storer, F BatchFinder[R], R any](db S, finder F, res *[]R) (err error) {
	query, _ := finder.Query()

	// 如果args参数里存在数组或切片，则分批获取
	for _, bargs := range finder.Batch() {
		rows, err := db.QueryContext(context.TODO(), query, bargs...)
		if err != nil {
			return err
		}
		defer rows.Close()

		colTypes, err := rows.ColumnTypes()
		if err != nil {
			return err
		}
		for rows.Next() {
			obj, fields := finder.NewScanObjAndFields(colTypes)
			if err = rows.Scan(fields...); err != nil {
				return err
			}
			// PrintFields(fields)

			*res = append(*res, *obj)
		}
		if err = rows.Err(); err != nil {
			return err
		}
	}

	return
}

// ExecWithBatch exec with batch
func ExecWithBatch[S Storer, Q BatchQueryer](db S, q Q) (ra, lid int64, err error) {
	query, _ := q.Query()

	for _, bargs := range q.Batch() {
		r, err := db.ExecContext(context.TODO(), query, bargs...)
		if err != nil {
			return 0, 0, err
		}
		lid, err = r.LastInsertId()
		if err != nil {
			return 0, 0, err
		}
		tra, err := r.RowsAffected()
		if err != nil {
			return 0, 0, err
		}
		ra += tra
	}

	return
}
