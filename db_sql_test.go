package do

import (
	"reflect"
	"testing"
)

func TestSQLProcess(t *testing.T) {
	type args struct {
		funcs []SQLQueryFunc
	}
	tests := []struct {
		name      string
		args      args
		wantQuery SQLQuery
		wantArgs  SQLArgs
	}{
		// TODO: Add test cases.
		{
			name: "id=1",
			args: args{
				funcs: []SQLQueryFunc{
					func() (query SQLQuery, args SQLArgs) {
						query = "select * from user where id = ?"
						args = append(args, 1)
						return
					},
				},
			},
			wantQuery: "select * from user where id = ?",
			wantArgs:  []any{1},
		},
		{
			name: "id=2&name=jd",
			args: args{
				funcs: []SQLQueryFunc{
					func() (query SQLQuery, args SQLArgs) {
						query = "select %s from user where id = ? %s"
						args = append(args, 1)
						return
					},
					func() (query SQLQuery, args SQLArgs) {
						query = "id, name"
						return
					},
					func() (query SQLQuery, args SQLArgs) {
						query = "and name = ?"
						args = append(args, "jd")
						return
					},
				},
			},
			wantQuery: "select id, name from user where id = ? and name = ?",
			wantArgs:  []any{1, "jd"},
		},
		{
			name: "id=2&name=jd",
			args: args{
				funcs: []SQLQueryFunc{
					// 首个函数确定语句架子，后面的函数补充内容
					func() (query SQLQuery, args SQLArgs) {
						query = "select %s from user %s where id = ? %s %s %s"
						args = append(args, 1)
						return
					},
					func() (query SQLQuery, args SQLArgs) {
						query = "id, name"
						return
					},
					func() (query SQLQuery, args SQLArgs) {
						query = "left join org on org.id = user.org_id and org.valid = ?"
						args = append(args, 1)
						return
					},
					func() (query SQLQuery, args SQLArgs) {
						query = "and name = ?"
						args = append(args, "jd")
						return
					},
					func() (query SQLQuery, args SQLArgs) {
						query = "group by user.id"
						return
					},
					func() (query SQLQuery, args SQLArgs) {
						query = "limit 10 offset 0"
						return
					},
				},
			},
			wantQuery: "select id, name from user left join org on org.id = user.org_id and org.valid = ? where id = ? and name = ? group by user.id limit 10 offset 0",
			wantArgs:  []any{1, 1, "jd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotArgs := SQLProcess(tt.args.funcs...)
			if gotQuery != tt.wantQuery {
				t.Errorf("SQLProcess() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("SQLProcess() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}

			{
				gotQuery, gotArgs := SQLProcessRaw(tt.args.funcs...)
				if gotQuery != tt.wantQuery.Raw() {
					t.Errorf("SQLProcess() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
				}
				if !reflect.DeepEqual(gotArgs, tt.wantArgs.Raw()) {
					t.Errorf("SQLProcess() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
				}
			}
		})
	}
}
