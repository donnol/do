package do

import (
	"database/sql"
	"testing"
)

type sqlResult struct {
	id int64
	n  int64
}

func (r *sqlResult) LastInsertId() (int64, error) {
	return r.id, nil
}

func (r *sqlResult) RowsAffected() (int64, error) {
	return r.n, nil
}

func TestHandleResult(t *testing.T) {
	type args struct {
		r sql.Result
	}
	tests := []struct {
		name    string
		args    args
		wantId  int64
		wantN   int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "exist",
			args: args{
				r: &sqlResult{n: 1, id: 1},
			},
			wantId:  1,
			wantN:   1,
			wantErr: false,
		},
		{
			name: "not-exist",
			args: args{
				r: &sqlResult{n: 0, id: 0},
			},
			wantId:  0,
			wantN:   0,
			wantErr: false,
		},
		{
			name: "id-exist-n-not-exist",
			args: args{
				r: &sqlResult{n: 0, id: 1},
			},
			wantId:  0,
			wantN:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotN, err := HandleResult(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("HandleResult() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotN != tt.wantN {
				t.Errorf("HandleResult() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
