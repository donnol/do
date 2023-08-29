package do

// BatchRun handle data per batch
func BatchRun[T any](s []T, batchNum int, handler func([]T) error) (err error) {
	batch := make([]T, 0, batchNum)
	for _, item := range s {
		item := item
		batch = append(batch, item)

		if len(batch) < batchNum {
			continue
		}

		// 达到一批了则执行
		if err := handler(batch); err != nil {
			return err
		}
		batch = make([]T, 0, batchNum)
	}
	// 剩余未满一批的
	if len(batch) > 0 {
		if err = handler(batch); err != nil {
			return err
		}
	}

	return
}

// StreamRun handle data by stream, if batchNum is >0, run with batch
func StreamRun[T any](s chan T, batchNum int, handler func([]T) error) (err error) {
	batch := make([]T, 0, batchNum)
	for e := range s {
		e := e
		if batchNum > 0 {
			batch = append(batch, e)
			if len(batch) < batchNum {
				continue
			}
		} else {
			batch = []T{e}
		}

		if err = handler(batch); err != nil {
			return err
		}
		batch = make([]T, 0, batchNum)
	}
	// 剩余未满一批的
	if len(batch) > 0 {
		if err = handler(batch); err != nil {
			return err
		}
	}

	return
}
