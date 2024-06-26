package sqlparser

import (
	"bytes"
	"reflect"
	"testing"
	"text/template"

	"github.com/andreyvit/diff"
	"github.com/donnol/do"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

func Test_parse(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name     string
		args     args
		wantCols []string
		wantErr  bool
	}{
		{
			name: "1",
			args: args{
				sql: "SELECT a, b FROM t",
			},
			wantCols: []string{"a", "b"},
		},
		{
			name: "2",
			args: args{
				sql: "SELECT a, b, c FROM t",
			},
			wantCols: []string{"a", "b", "c"},
		},
		{
			name: "3",
			args: args{
				sql: `create table user (
					id integer not null,
					name varchar(255) not null, 
					created_at datetime not null, 
					updated_at timestamp not null
				)`,
			},
			wantCols: []string{"id", "name", "created_at", "updated_at"},
		},
		{
			name: "4",
			args: args{
				sql: `update user set name = 'jd' where id = 1`,
			},
			wantCols: []string{"name", "id"},
		},
		{
			name: "5",
			args: args{
				sql: `insert into user (name) values ('jd')`,
			},
			wantCols: []string{"name"},
		},
		{
			name: "6",
			args: args{
				sql: `delete from user where id = 1`,
			},
			wantCols: []string{"id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r := extract(got)
			if !reflect.DeepEqual(r, tt.wantCols) {
				t.Errorf("parse() = %v, want %v", r, tt.wantCols)
			}
		})
	}
}

type UserTable struct{}
type _UserTableEnum struct {
}

func (UserTable) EnumHelper() _UserTableEnum {
	e := _UserTableEnum{}

	return e
}

var _ = func() struct{} {
	e := UserTable{}.EnumHelper()
	_ = e
	return struct{}{}
}()

func TestStruct_Gen(t *testing.T) {
	type fields struct {
		TableName string
		Name      string
		Comment   string
		Fields    []Field
	}
	type args struct {
		opt Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				Name:      "User",
				TableName: "user_table",
				Comment:   "用户表",
				Fields: []Field{
					{
						Name:    "id",
						DBField: "id",
						Type:    "UNSIGNED BIGINT",
						Tag:     "",
						Comment: "主键id",
					},
					{
						Name:    "name",
						DBField: "name",
						Type:    "varchar",
						Tag:     "",
						Comment: "名称",
					},
					{
						Name:    "created_at",
						DBField: "created_at",
						Type:    "datetime",
						Tag:     "",
						Comment: "创建时间",
					},
					{
						Name:    "updated_at",
						DBField: "updated_at",
						Type:    "timestamp",
						Tag:     "",
						Comment: "更新时间",
					},
				},
			},
			args: args{
				opt: Option{},
			},
			wantW: "\n" +
				"	// UserTable 用户表" + "\n" +
				"	type UserTable struct {" + "\n" +
				"		Id uint64 `json:\"id\" db:\"id\"` // 主键id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"		UpdatedAt time.Time `json:\"updatedAt\" db:\"updated_at\"` // 更新时间" + "\n" +
				"" + "\n" +
				"		useAlias bool" + "\n" +
				"		alias    string" + "\n" +
				"	}" + "\n	\n" +
				`	func (UserTable) TableName() string {
		return "user_table"
	}

	func (s UserTable) UseAlias(alias ...string) UserTable {
		s.useAlias = true
		if len(alias) > 0 {
			s.alias = alias[0]
		}
		return s
	}
	
	func (s UserTable) Columns() []string {
		return s.NameHelper().Columns()
	}
	
	func (s UserTable) Values() []any {
		return []any{
			s.Id,
			s.Name,
			s.CreatedAt,
			s.UpdatedAt,
		}
	}
	
	func (s *UserTable) ValuePtrs() []any {
		return []any{
			&s.Id,
			&s.Name,
			&s.CreatedAt,
			&s.UpdatedAt,
		}
	}
	
	func (s UserTable) Exists() bool {
		return s.Id != 0
	}
	
	type _UserTableNameHelper struct {
		Id string // field: id
		Name string // field: name
		CreatedAt string // field: created_at
		UpdatedAt string // field: updated_at
	}
	
	func (_UserTableNameHelper) Columns() []string {
		return []string{
			"id",
			"name",
			"created_at",
			"updated_at",
		}
	}
	
	func (s UserTable) NameHelper() _UserTableNameHelper {
		withAlias := func(field string, alias string) string {
			if alias == "" {
				return field
			}
			return alias + "." + field
		}
		alias := ""
		if s.useAlias {
			if s.alias != "" {
				alias = s.alias
			} else {
				alias = s.TableName()
			}
		}
		return _UserTableNameHelper{
			Id: withAlias("id", alias),
			Name: withAlias("name", alias),
			CreatedAt: withAlias("created_at", alias),
			UpdatedAt: withAlias("updated_at", alias),
		}
	}
	
	
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Struct{
				Name:      tt.fields.Name,
				TableName: tt.fields.TableName,
				Comment:   tt.fields.Comment,
				Fields:    tt.fields.Fields,
			}
			w := &bytes.Buffer{}
			if err := s.Gen(w, tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("Struct.Gen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Struct.Gen() = %v, want %v, diff: %s", gotW, tt.wantW, diff.LineDiff(gotW, tt.wantW))
			}
		})
	}
}

func TestParseCreateSQL(t *testing.T) {
	type args struct {
		sql string
		opt Option
	}
	tests := []struct {
		name  string
		args  args
		want  *Struct
		wantW string
	}{
		{
			name: "1",
			args: args{
				sql: `create table user (
id integer unsigned not null comment 'id',
name varchar(255) not null comment '名称', 
created_at datetime not null comment '创建时间', 
updated_at timestamp not null comment '更新时间'
) comment '用户表'`,
			},
			want: &Struct{
				Name:      "",
				TableName: "user",
				Comment:   "用户表",
				Fields: []Field{
					{Name: "id", DBField: "id", Type: "int unsigned", Tag: "", Comment: "id"},
					{Name: "name", DBField: "name", Type: "varchar", Tag: "", Comment: "名称"},
					{Name: "created_at", DBField: "created_at", Type: "datetime", Tag: "", Comment: "创建时间"},
					{Name: "updated_at", DBField: "updated_at", Type: "timestamp", Tag: "", Comment: "更新时间"},
				},
			},
			wantW: "\n" +
				"	// User 用户表" + "\n" +
				"	type User struct {" + "\n" +
				"		Id uint `json:\"id\" db:\"id\"` // id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"		UpdatedAt time.Time `json:\"updatedAt\" db:\"updated_at\"` // 更新时间" + "\n" +
				"" + "\n" +
				"		useAlias bool" + "\n" +
				"		alias    string" + "\n" +
				"	}" + "\n	\n" +
				`	func (User) TableName() string {
		return "user"
	}

	func (s User) UseAlias(alias ...string) User {
		s.useAlias = true
		if len(alias) > 0 {
			s.alias = alias[0]
		}
		return s
	}
	
	func (s User) Columns() []string {
		return s.NameHelper().Columns()
	}
	
	func (s User) Values() []any {
		return []any{
			s.Id,
			s.Name,
			s.CreatedAt,
			s.UpdatedAt,
		}
	}
	
	func (s *User) ValuePtrs() []any {
		return []any{
			&s.Id,
			&s.Name,
			&s.CreatedAt,
			&s.UpdatedAt,
		}
	}
	
	func (s User) Exists() bool {
		return s.Id != 0
	}
	
	type _UserNameHelper struct {
		Id string // field: id
		Name string // field: name
		CreatedAt string // field: created_at
		UpdatedAt string // field: updated_at
	}
	
	func (_UserNameHelper) Columns() []string {
		return []string{
			"id",
			"name",
			"created_at",
			"updated_at",
		}
	}
	
	func (s User) NameHelper() _UserNameHelper {
		withAlias := func(field string, alias string) string {
			if alias == "" {
				return field
			}
			return alias + "." + field
		}
		alias := ""
		if s.useAlias {
			if s.alias != "" {
				alias = s.alias
			} else {
				alias = s.TableName()
			}
		}
		return _UserNameHelper{
			Id: withAlias("id", alias),
			Name: withAlias("name", alias),
			CreatedAt: withAlias("created_at", alias),
			UpdatedAt: withAlias("updated_at", alias),
		}
	}
	
	
`,
		},
		{
			name: "ignoreField",
			args: args{
				sql: `create table user (
id integer unsigned not null comment 'id',
name varchar(255) not null comment '名称', 
created_at datetime not null comment '创建时间', 
updated_at timestamp not null comment '更新时间'
) comment '用户表'`,
				opt: Option{
					IgnoreField: []string{"updated_at"},
				},
			},
			want: &Struct{
				Name:      "",
				TableName: "user",
				Comment:   "用户表",
				Fields: []Field{
					{Name: "id", DBField: "id", Type: "int unsigned", Tag: "", Comment: "id"},
					{Name: "name", DBField: "name", Type: "varchar", Tag: "", Comment: "名称"},
					{Name: "created_at", DBField: "created_at", Type: "datetime", Tag: "", Comment: "创建时间"},
					{Name: "updated_at", DBField: "updated_at", Type: "timestamp", Tag: "", Comment: "更新时间"},
				},
			},
			wantW: "\n" +
				"	// User 用户表" + "\n" +
				"	type User struct {" + "\n" +
				"		Id uint `json:\"id\" db:\"id\"` // id" + "\n" +
				"		Name string `json:\"name\" db:\"name\"` // 名称" + "\n" +
				"		CreatedAt time.Time `json:\"createdAt\" db:\"created_at\"` // 创建时间" + "\n" +
				"" + "\n" +
				"		useAlias bool" + "\n" +
				"		alias    string" + "\n" +
				"	}" + "\n	\n" +
				`	func (User) TableName() string {
		return "user"
	}

	func (s User) UseAlias(alias ...string) User {
		s.useAlias = true
		if len(alias) > 0 {
			s.alias = alias[0]
		}
		return s
	}
	
	func (s User) Columns() []string {
		return s.NameHelper().Columns()
	}
	
	func (s User) Values() []any {
		return []any{
			s.Id,
			s.Name,
			s.CreatedAt,
		}
	}
	
	func (s *User) ValuePtrs() []any {
		return []any{
			&s.Id,
			&s.Name,
			&s.CreatedAt,
		}
	}
	
	func (s User) Exists() bool {
		return s.Id != 0
	}
	
	type _UserNameHelper struct {
		Id string // field: id
		Name string // field: name
		CreatedAt string // field: created_at
	}
	
	func (_UserNameHelper) Columns() []string {
		return []string{
			"id",
			"name",
			"created_at",
		}
	}
	
	func (s User) NameHelper() _UserNameHelper {
		withAlias := func(field string, alias string) string {
			if alias == "" {
				return field
			}
			return alias + "." + field
		}
		alias := ""
		if s.useAlias {
			if s.alias != "" {
				alias = s.alias
			} else {
				alias = s.TableName()
			}
		}
		return _UserNameHelper{
			Id: withAlias("id", alias),
			Name: withAlias("name", alias),
			CreatedAt: withAlias("created_at", alias),
		}
	}
	
	
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCreateSQL(tt.args.sql); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCreateSQL() = %+v, want %+v", got, tt.want)
			} else {
				buf := new(bytes.Buffer)
				if err := got.Gen(buf, tt.args.opt); err != nil {
					t.Error(err)
				}
				if buf.String() != tt.wantW {
					t.Errorf("Struct.Gen() = %v, want %v, diff: %s", buf.String(), tt.wantW, diff.LineDiff(buf.String(), tt.wantW))
				}
			}
		})
	}
}

func TestStruct_GenData(t *testing.T) {
	type args struct {
		sql string
		n   int64
		opt Option
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "full",
			args: args{
				sql: `create table user (
					id integer unsigned not null comment 'id',
					name varchar(255) not null comment '名称', 
					created_at datetime not null comment '创建时间', 
					updated_at timestamp not null comment '更新时间'
					) comment '用户表'`,
				n:   0,
				opt: doption,
			},
			wantW: "INSERT IGNORE INTO `user` (" + "\n" +
				"`id`," + "\n" +
				"`name`," + "\n" +
				"`created_at`," + "\n" +
				"`updated_at`" + "\n" +
				") VALUES (",
			wantErr: false,
		},
		{
			name: "ignore",
			args: args{
				sql: `create table user (
					id integer unsigned not null comment 'id',
					name varchar(255) not null comment '名称', 
					created_at datetime not null comment '创建时间', 
					updated_at timestamp not null comment '更新时间'
					) comment '用户表'`,
				n: 0,
				opt: Option{
					IgnoreField: []string{"updated_at"},
				},
			},
			wantW: "INSERT IGNORE INTO `user` (" + "\n" +
				"`id`," + "\n" +
				"`name`," + "\n" +
				"`created_at`" + "\n" +
				") VALUES (",
			wantErr: false,
		},
		{
			name: "ignore",
			args: args{
				sql: `create table user (
					id integer unsigned not null comment 'id',
					name varchar(255) not null comment '名称', 
					created_at datetime not null comment '创建时间', 
					updated_at timestamp not null comment '更新时间'
					) comment '用户表'`,
				n: 2,
				opt: Option{
					IgnoreField: []string{"updated_at"},
				},
			},
			wantW: "INSERT IGNORE INTO `user` (" + "\n" +
				"`id`," + "\n" +
				"`name`," + "\n" +
				"`created_at`" + "\n" +
				") VALUES (",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ParseCreateSQL(tt.args.sql)
			w := &bytes.Buffer{}
			if err := s.GenData(w, tt.args.n, tt.args.opt); (err != nil) != tt.wantErr {
				t.Errorf("Struct.GenData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 因为值是随机生成的，所以只比较前面部分
			if gotW := w.String(); len(gotW) <= len(tt.wantW) {
				t.Errorf("Struct.GenData() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestParseCreateSQLBatch(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want []*Struct
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				sql: `create table user (
					id integer not null,
					name varchar(255) not null, 
					created_at datetime not null, 
					updated_at timestamp not null
				);
				create table role (
					id integer not null,
					name varchar(255) not null, 
					created_at datetime not null, 
					updated_at timestamp not null
				);
				`,
			},
			want: []*Struct{
				{TableName: "user", Fields: []Field{
					{
						Name:    "id",
						Type:    "int",
						DBField: "id",
					},
					{
						Name:    "name",
						Type:    "varchar",
						DBField: "name",
					},
					{
						Name:    "created_at",
						Type:    "datetime",
						DBField: "created_at",
					},
					{
						Name:    "updated_at",
						Type:    "timestamp",
						DBField: "updated_at",
					},
				}},
				{TableName: "role", Fields: []Field{
					{
						Name:    "id",
						Type:    "int",
						DBField: "id",
					},
					{
						Name:    "name",
						Type:    "varchar",
						DBField: "name",
					},
					{
						Name:    "created_at",
						Type:    "datetime",
						DBField: "created_at",
					},
					{
						Name:    "updated_at",
						Type:    "timestamp",
						DBField: "updated_at",
					},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCreateSQLBatch(tt.args.sql)
			for i, one := range got {
				if !reflect.DeepEqual(*one, *tt.want[i]) {
					t.Errorf("ParseCreateSQLBatch() = %v, want %v", *one, *tt.want[i])
				}
			}
		})
	}
}

func Test_processFieldType(t *testing.T) {
	type args struct {
		fieldType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				fieldType: "integer",
			},
			want: "integer",
		},
		{
			name: "2",
			args: args{
				fieldType: "varchar(255)",
			},
			want: "varchar",
		},
		{
			name: "3",
			args: args{
				fieldType: "double(10,2)",
			},
			want: "double",
		},
		{
			name: "4",
			args: args{
				fieldType: "BIGINT UNSIGNED",
			},
			want: "BIGINT UNSIGNED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processFieldType(tt.args.fieldType); got != tt.want {
				t.Errorf("processFieldType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultToJSONObjectTmpl(t *testing.T) {
	tmpl, err := template.New("resultToJSONObjectTmpl").Parse(ResultToJSONObjectTmpl)
	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, ResultToJSONObject{
		Fields: []ResultToJSONObjectField{
			{
				JSONName:   "id",
				ColumnName: "id",
				NoComma:    false,
			},
			{
				JSONName:   "name",
				ColumnName: "name",
				NoComma:    false,
			},
			{
				JSONName:   "createTime",
				ColumnName: "create_time",
				NoComma:    true,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	want := `json_object(
		'id', id,
		'name', name,
		'createTime', create_time
	)`
	do.Assert(t, buf.String(), want, diff.LineDiff(buf.String(), want))
}

func TestResultToJSONObjectTmpl2(t *testing.T) {
	tmpl, err := template.New("resultToJSONObjectTmpl").Parse(ResultToJSONObjectTmpl)
	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, FromStructForTmpl(&StructForTmpl{
		Fields: []StructField{
			{
				DBField: "id",
			},
			{
				DBField: "name",
			},
			{
				DBField: "create_time",
			},
		},
	}))
	if err != nil {
		t.Error(err)
	}

	want := `json_object(
		'id', id,
		'name', name,
		'createTime', create_time
	)`
	do.Assert(t, buf.String(), want, diff.LineDiff(buf.String(), want))
}

type tstruct struct {
	modifyMap map[string]struct{}
}

func (t *tstruct) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.TableName:
		if _, ok := t.modifyMap["table|"+node.Name.L]; !ok {
			node.Name.L += " u"
			node.Name.O += " u"
			t.modifyMap["table|"+node.Name.L] = struct{}{}
		}
	case *ast.SelectStmt:
		node.From.TableRefs.Left.Accept(t)
		for _, field := range node.Fields.Fields {
			field.Expr.Accept(t)
		}
	case *ast.OrderByClause:
		for _, item := range node.Items {
			item.Expr.Accept(t)
		}
	case *ast.ColumnNameExpr:
		if _, ok := t.modifyMap["column|"+node.Name.Name.L]; !ok {
			node.Name.Name.L = "u." + node.Name.Name.L
			node.Name.Name.O = "u." + node.Name.Name.O
			t.modifyMap["column|"+node.Name.Name.L] = struct{}{}
		}
	case *ast.GroupByClause:
		for _, item := range node.Items {
			item.Expr.Accept(t)
		}
	case *ast.Assignment:
	case *ast.ByItem:
	case *ast.FieldList:
	case *ast.HavingClause:
	case *ast.AsOfClause:
	case *ast.Join:
	case *ast.Limit:
	case *ast.OnCondition:
	case *ast.SelectField:
	case *ast.TableRefsClause:
	case *ast.TableSource:
	case *ast.SetOprSelectList:
	case *ast.WildCardField:
	case *ast.WindowSpec:
	case *ast.PartitionByClause:
	case *ast.FrameClause:
	case *ast.FrameBound:
	}

	return in, false
}

func (t *tstruct) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func TestSetSelectAlias(t *testing.T) {
	s := &tstruct{
		modifyMap: make(map[string]struct{}),
	}

	sql := `select id, name from t_user where id = 1 group by id order by id desc limit 1`
	p := parser.New()

	nodes, _, err := p.ParseSQL(sql)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	for _, node := range nodes {
		node.Accept(s)
		node.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, buf))
		do.Assert(t, node.Text(), sql)
	}

	want := "SELECT `u.id`,`u.name` FROM `t_user u` WHERE `u.id`=1 GROUP BY `u.id` ORDER BY `u.id` DESC LIMIT 1"
	do.Assert(t, buf.String(), want)
}
