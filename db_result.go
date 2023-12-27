package do

import "database/sql"

// HandleResult check affected rows first, if it is not zero then get last insert id
func HandleResult(r sql.Result) (id int64, n int64, err error) {
	// Must check affected rows first before get last insert id
	// https://stackoverflow.com/questions/20704983/clearing-last-insert-id-before-inserting-to-tell-if-whats-returned-is-from-my
	n, err = r.RowsAffected()
	if err != nil {
		return
	}
	if n > 0 {
		id, err = r.LastInsertId()
		if err != nil {
			return
		}
	}
	return
}
