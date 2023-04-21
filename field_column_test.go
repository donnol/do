package do

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestFieldsByColumnType(t *testing.T) {
	type args struct {
		t           any
		colTypes    []*sql.ColumnType
		fieldMapper func(string) string
	}
	tests := []struct {
		name       string
		args       args
		wantFields []any
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				t:           &User{},
				colTypes:    []*sql.ColumnType{},
				fieldMapper: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFields := FieldsByColumnType(tt.args.t, tt.args.colTypes, tt.args.fieldMapper); !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Errorf("FieldsByColumnType() = %v, want %v", gotFields, tt.wantFields)
			}
		})
	}
}
